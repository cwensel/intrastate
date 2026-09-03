//go:build ignore

// Spike A10 for RDR 0026: the pollability probe is SAFE and DETECTING.
//
// The probe the RDR proposes runs once, at pipe creation:
//
//	if err := rd.SetReadDeadline(time.Now().Add(x)); err != nil { ...not pollable... }
//	_ = rd.SetReadDeadline(time.Time{})   // clear it again immediately
//
// Two independent properties have to hold, and a failure of EITHER sinks the
// design:
//
//	(a) DETECTING — on a file the runtime poller never registered,
//	    SetReadDeadline must FAIL, and fail recognisably, so the caller can
//	    fall back rather than silently rely on a deadline that will never
//	    fire. Two non-pollable shapes are constructed: a regular file, and
//	    os.NewFile over a syscall.Dup'd pipe fd (dup + NewFile bypasses the
//	    poller registration os.Pipe performs, giving a genuine pipe fd that
//	    is nonetheless unpollable).
//
//	(b) SAFE / NO-OP — on an ordinary pollable os.Pipe, probe-then-clear must
//	    leave read behaviour EXACTLY as it was. The weak way to test this is
//	    to read a payload and see it arrive; that would pass even if the
//	    clear had silently left a deadline armed, because the bytes were
//	    already there. The strong test, run here, is that after probe+clear a
//	    blocking Read on an EMPTY pipe with a LIVE writer still parks
//	    INDEFINITELY: verify it is still parked at 300 ms (a leaked deadline
//	    would have fired long before, the probe's own window being far
//	    shorter), then have the writer write and confirm the parked read
//	    wakes with the right bytes. Byte-for-byte round-trip equality with
//	    and without the probe is also checked.
package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"syscall"
	"time"
)

// probeWindow is the deadline the probe arms before clearing it. Short, and
// far under the 300 ms parked-read observation, so a leaked deadline could
// not hide.
const probeWindow = 10 * time.Millisecond

// parkObservation is how long the safety test watches a blocking read to
// confirm it is still parked. 30x the probe window.
const parkObservation = 300 * time.Millisecond

func main() {
	fmt.Printf("### RDR 0026 A10 spike: pollability probe ###\n")
	fmt.Printf("platform: %s/%s  go: %s\n", runtime.GOOS, runtime.GOARCH, runtime.Version())
	fmt.Printf("probe window: %v   parked-read observation: %v\n\n", probeWindow, parkObservation)

	partA()
	partB1()
	partB2()
}

// probe performs exactly the RDR's probe: arm, then clear. It reports the
// error from the ARM (the detecting signal) and from the CLEAR.
func probe(f *os.File) (armErr, clearErr error) {
	armErr = f.SetReadDeadline(time.Now().Add(probeWindow))
	clearErr = f.SetReadDeadline(time.Time{})
	return armErr, clearErr
}

func classify(err error) string {
	switch {
	case err == nil:
		return "<nil>"
	case errors.Is(err, os.ErrNoDeadline):
		return fmt.Sprintf("%q [errors.Is(err, os.ErrNoDeadline)=true, type=%T]", err.Error(), err)
	default:
		return fmt.Sprintf("%q [errors.Is(err, os.ErrNoDeadline)=false, type=%T]", err.Error(), err)
	}
}

// --- (a) DETECTING --------------------------------------------------------

func partA() {
	fmt.Printf("--- A10(a): the probe DETECTS a non-pollable file ---\n")
	fmt.Printf("Claim: SetReadDeadline returns os.ErrNoDeadline on a file the runtime\n")
	fmt.Printf("poller never registered, so the caller can tell and fall back.\n\n")

	// Shape 1: an ordinary regular file. Never pollable on any platform.
	tmp, err := os.CreateTemp("", "a10-regular-*")
	must(err)
	tmpPath := tmp.Name()
	_, _ = tmp.WriteString("regular file contents")
	_ = tmp.Close()
	rf, err := os.Open(tmpPath)
	must(err)
	armErr, clearErr := probe(rf)
	fmt.Printf("shape 1: regular file (os.Open of a temp file)\n")
	fmt.Printf("  SetReadDeadline(arm):   %s\n", classify(armErr))
	fmt.Printf("  SetReadDeadline(clear): %s\n", classify(clearErr))
	fmt.Printf("  detected as non-pollable: %v\n", errors.Is(armErr, os.ErrNoDeadline))
	// A regular file still reads fine after a failed probe.
	b, rerr := io.ReadAll(rf)
	fmt.Printf("  read after failed probe: n=%d err=%v content=%q\n\n", len(b), rerr, string(b))
	_ = rf.Close()
	_ = os.Remove(tmpPath)

	// Shape 2: os.NewFile over a DUP of a pipe read-end fd. The fd is a
	// genuine pipe, but syscall.Dup + os.NewFile bypasses the poller
	// registration that os.Pipe performs, so it is unpollable.
	pr, pw, err := os.Pipe()
	must(err)
	dupFd, err := syscall.Dup(int(pr.Fd()))
	must(err)
	dup := os.NewFile(uintptr(dupFd), "dup-of-pipe-read-end")
	armErr2, clearErr2 := probe(dup)
	fmt.Printf("shape 2: os.NewFile over syscall.Dup of an os.Pipe read end\n")
	fmt.Printf("  (a genuine pipe fd, but never registered with the runtime poller)\n")
	fmt.Printf("  SetReadDeadline(arm):   %s\n", classify(armErr2))
	fmt.Printf("  SetReadDeadline(clear): %s\n", classify(clearErr2))
	fmt.Printf("  detected as non-pollable: %v\n", errors.Is(armErr2, os.ErrNoDeadline))
	// It still functions as a pipe: write, then read it back through the dup.
	_, _ = pw.Write([]byte("through-the-dup"))
	rb := make([]byte, 64)
	n, rerr2 := dup.Read(rb)
	fmt.Printf("  read through the dup after failed probe: n=%d err=%v content=%q\n",
		n, rerr2, string(rb[:n]))
	_ = dup.Close()
	_ = pr.Close()
	_ = pw.Close()

	// Contrast: the ordinary os.Pipe read end, which IS pollable.
	cr, cw, err := os.Pipe()
	must(err)
	armErr3, clearErr3 := probe(cr)
	fmt.Printf("\ncontrast: ordinary os.Pipe read end (expected pollable)\n")
	fmt.Printf("  SetReadDeadline(arm):   %s\n", classify(armErr3))
	fmt.Printf("  SetReadDeadline(clear): %s\n", classify(clearErr3))
	fmt.Printf("  detected as non-pollable: %v\n", errors.Is(armErr3, os.ErrNoDeadline))
	_ = cr.Close()
	_ = cw.Close()

	fmt.Printf("\n  os.ErrNoDeadline as declared by the stdlib: %q\n\n", os.ErrNoDeadline.Error())
}

// --- (b1) SAFE: a blocking read still parks indefinitely ------------------

type parkTrial struct {
	stillParkedAt300ms bool
	wokeAfterWrite     bool
	bytesCorrect       bool
	err                error
	wakeLatency        time.Duration
}

// parkedReadTrial probes-and-clears an ordinary pipe, parks a Read on the
// EMPTY pipe with the writer still open, and checks the read is still parked
// after parkObservation. A deadline leaked by the probe would have fired at
// probeWindow, 30x sooner. Then the writer writes, and the parked read must
// wake with exactly those bytes.
func parkedReadTrial(withProbe bool) parkTrial {
	var t parkTrial
	rd, wr, err := os.Pipe()
	must(err)
	defer rd.Close()
	defer wr.Close()

	if withProbe {
		if armErr, clearErr := probe(rd); armErr != nil || clearErr != nil {
			t.err = fmt.Errorf("probe failed on an ordinary pipe: arm=%v clear=%v", armErr, clearErr)
			return t
		}
	}

	const want = "WOKE-WITH-THESE-BYTES"
	type readOut struct {
		n   int
		err error
		at  time.Time
		buf []byte
	}
	done := make(chan readOut, 1)
	go func() {
		b := make([]byte, 64)
		n, rerr := rd.Read(b) // parks: pipe is empty, writer is open
		done <- readOut{n: n, err: rerr, at: time.Now(), buf: b}
	}()

	// Watch it stay parked. Nothing has been written, so the ONLY thing that
	// could return this read early is a deadline the probe failed to clear.
	select {
	case out := <-done:
		t.err = fmt.Errorf("read returned EARLY after %d bytes, err=%v -- probe perturbed it", out.n, out.err)
		return t
	case <-time.After(parkObservation):
		t.stillParkedAt300ms = true
	}

	wroteAt := time.Now()
	_, _ = wr.Write([]byte(want))
	select {
	case out := <-done:
		t.wokeAfterWrite = true
		t.wakeLatency = out.at.Sub(wroteAt)
		t.bytesCorrect = out.err == nil && string(out.buf[:out.n]) == want
		if out.err != nil {
			t.err = fmt.Errorf("read after write errored: %v", out.err)
		}
	case <-time.After(2 * time.Second):
		t.err = errors.New("read never woke after the write")
	}
	return t
}

func partB1() {
	fmt.Printf("--- A10(b): the probe is a NO-OP on an ordinary pollable pipe ---\n")
	fmt.Printf("Claim: after probe+clear, a blocking Read on an EMPTY pipe with a live\n")
	fmt.Printf("writer still parks INDEFINITELY (no deadline error), then wakes on the\n")
	fmt.Printf("write with the right bytes. A leaked %v deadline would fire at %v,\n",
		probeWindow, probeWindow)
	fmt.Printf("%.0fx sooner than the %v observation window.\n\n",
		float64(parkObservation)/float64(probeWindow), parkObservation)

	const iters = 12
	for _, withProbe := range []bool{true, false} {
		label := "WITH probe+clear"
		if !withProbe {
			label = "WITHOUT probe (baseline)"
		}
		parked, woke, correct, failed := 0, 0, 0, 0
		var worstWake time.Duration
		var firstErr error
		for i := 0; i < iters; i++ {
			t := parkedReadTrial(withProbe)
			if t.stillParkedAt300ms {
				parked++
			}
			if t.wokeAfterWrite {
				woke++
			}
			if t.bytesCorrect {
				correct++
			}
			if t.err != nil {
				failed++
				if firstErr == nil {
					firstErr = t.err
				}
			}
			if t.wakeLatency > worstWake {
				worstWake = t.wakeLatency
			}
		}
		fmt.Printf("%s, n=%d\n", label, iters)
		fmt.Printf("  still parked at %v (no deadline fired): %d/%d\n", parkObservation, parked, iters)
		fmt.Printf("  woke on the write:                      %d/%d\n", woke, iters)
		fmt.Printf("  woke with the CORRECT bytes:            %d/%d\n", correct, iters)
		fmt.Printf("  trials with any error:                  %d/%d\n", failed, iters)
		fmt.Printf("  worst wake latency after the write:     %v\n", worstWake.Round(time.Microsecond))
		if firstErr != nil {
			fmt.Printf("  first error: %v\n", firstErr)
		}
		fmt.Println()
	}
}

// --- (b2) SAFE: full payload round-trips identically ----------------------

// roundTrip writes a size-byte envelope through a pipe from a goroutine and
// reads it all back, optionally probing the read end first. It reports
// whether the recovered bytes equal the envelope exactly.
func roundTrip(size int, withProbe bool) (ok bool, got int, err error) {
	rd, wr, perr := os.Pipe()
	if perr != nil {
		return false, 0, perr
	}
	defer rd.Close()

	if withProbe {
		if armErr, clearErr := probe(rd); armErr != nil || clearErr != nil {
			return false, 0, fmt.Errorf("probe: arm=%v clear=%v", armErr, clearErr)
		}
	}

	want := envelope(size)
	go func() {
		_, _ = wr.Write(want)
		_ = wr.Close() // gives the reader its EOF
	}()

	b, rerr := io.ReadAll(rd)
	if rerr != nil {
		return false, len(b), rerr
	}
	return bytes.Equal(b, want), len(b), nil
}

func partB2() {
	fmt.Printf("--- A10(b, cont): full payload round-trips byte-for-byte, probe or not ---\n\n")

	sizes := []struct {
		name string
		n    int
	}{
		{"1KiB", 1 << 10},
		{"64KiB", 64 << 10},
		{"1MiB", 1 << 20},
	}
	const iters = 12

	fmt.Printf("%-9s %-26s %-10s %-12s %s\n", "payload", "arm", "identical", "bytes", "err")
	for _, s := range sizes {
		for _, withProbe := range []bool{true, false} {
			label := "WITH probe+clear"
			if !withProbe {
				label = "WITHOUT probe (baseline)"
			}
			okCount, minGot := 0, -1
			var firstErr error
			for i := 0; i < iters; i++ {
				ok, got, err := roundTrip(s.n, withProbe)
				if ok {
					okCount++
				}
				if minGot < 0 || got < minGot {
					minGot = got
				}
				if err != nil && firstErr == nil {
					firstErr = err
				}
			}
			es := "<nil>"
			if firstErr != nil {
				es = firstErr.Error()
			}
			fmt.Printf("%-9s %-26s %-10s %-12d %s\n",
				s.name, label, fmt.Sprintf("%d/%d", okCount, iters), minGot, es)
		}
	}
	fmt.Println()
}

// --- helpers --------------------------------------------------------------

func envelope(size int) []byte {
	b := make([]byte, size)
	for i := range b {
		b[i] = byte('a' + (i % 26))
	}
	copy(b, "ENVELOPE-BEGIN:")
	if size >= 13 {
		copy(b[size-13:], ":ENVELOPE-END")
	}
	return b
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
