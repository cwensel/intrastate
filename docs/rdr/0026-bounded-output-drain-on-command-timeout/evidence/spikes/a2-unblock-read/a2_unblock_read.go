// Spike: RDR 0026 Critical Assumption A2.
//
// A2: the joining goroutine can END a pending blocking Read on an os.Pipe read
// end from another goroutine -- SetReadDeadline or Close on the read end -- and
// the drain goroutine then returns promptly with the bytes it had read.
//
// Every scenario keeps the WRITE end held open by a setsid(2) grandchild, so
// there is no EOF: the blocked Read can only end because of the deadline or the
// Close, never because the writer went away.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

// holderMain runs in the re-exec'd grandchild. fd 3 is the inherited pipe write
// end. It optionally writes a preamble, then sits on the open write end for the
// given lifetime so the parent's read end never sees EOF.
func holderMain() {
	preamble := os.Args[2]
	lifetimeMS, _ := strconv.Atoi(os.Args[3])
	if lifetimeMS <= 0 {
		// A zero lifetime would make the writer vanish instantly and every
		// scenario would see a spurious EOF. Refuse rather than lie.
		fmt.Fprintf(os.Stderr, "holder: refusing lifetime %dms\n", lifetimeMS)
		os.Exit(2)
	}
	w := os.NewFile(3, "pipe-write-end")
	if rf := os.Getenv("A2_READY_FILE"); rf != "" {
		_ = os.WriteFile(rf, []byte(strconv.Itoa(os.Getpid())), 0o600)
	}
	if preamble != "" {
		_, _ = w.Write([]byte(preamble))
	}
	// Hold the write end open, quiet, for lifetimeMS. No further writes.
	time.Sleep(time.Duration(lifetimeMS) * time.Millisecond)
	os.Exit(0)
}

// spawnHolder starts a setsid grandchild holding w, then closes the parent's
// copy of w so the read end's only remaining writer is the escaped process.
func spawnHolder(w *os.File, preamble string, lifetime time.Duration) (*exec.Cmd, error) {
	if lifetime < time.Second {
		return nil, fmt.Errorf("holder lifetime %v too short: the writer must outlive the measurement", lifetime)
	}
	self, err := os.Executable()
	if err != nil {
		return nil, err
	}
	readyFile, err := os.CreateTemp("", "a2-ready-*")
	if err != nil {
		return nil, err
	}
	readyPath := readyFile.Name()
	_ = readyFile.Close()
	_ = os.Remove(readyPath)
	cmd := exec.Command(self, "holder", preamble, strconv.Itoa(int(lifetime.Milliseconds())))
	cmd.Env = append(os.Environ(), "A2_READY_FILE="+readyPath)
	cmd.ExtraFiles = []*os.File{w} // becomes fd 3 in the child
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	cmd.Stdout, cmd.Stderr = nil, nil
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	// Parent drops its write end: only the grandchild holds a writer now.
	_ = w.Close()
	// Block until the grandchild confirms it is running and holding fd 3, so a
	// scenario can never mistake "writer already exited" for a real result.
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if b, err := os.ReadFile(readyPath); err == nil && len(b) > 0 {
			_ = os.Remove(readyPath)
			return cmd, nil
		}
		time.Sleep(5 * time.Millisecond)
	}
	killHolder(cmd)
	return nil, fmt.Errorf("holder did not signal ready")
}

func killHolder(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	// Setsid put it in its own group; kill the group to be sure nothing strays.
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	_ = cmd.Process.Kill()
	_, _ = cmd.Process.Wait()
}

func describe(err error) string {
	switch {
	case err == nil:
		return "<nil>"
	case errors.Is(err, os.ErrDeadlineExceeded):
		return fmt.Sprintf("%v [errors.Is(err, os.ErrDeadlineExceeded)=true, type=%T]", err, err)
	case errors.Is(err, os.ErrClosed):
		return fmt.Sprintf("%v [errors.Is(err, os.ErrClosed)=true, type=%T]", err, err)
	case errors.Is(err, io.EOF):
		return fmt.Sprintf("%v [errors.Is(err, io.EOF)=true -- WRITER LEAKED, scenario invalid]", err)
	default:
		return fmt.Sprintf("%v [type=%T]", err, err)
	}
}

// ---- Point 1/2/3: blocked Read ends via SetReadDeadline, on an EMPTY pipe ----

func scenarioDeadlineOnBlockedRead() {
	fmt.Println("=== S1: SetReadDeadline unblocks a Read already blocked on an empty pipe ===")
	fmt.Println("    (points 1, 2, 3; also point 4b: empty pipe + live writer -> expect deadline error)")
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Println("  os.Pipe:", err)
		return
	}
	defer r.Close()
	holder, err := spawnHolder(w, "", 4*time.Second)
	if err != nil {
		fmt.Println("  spawnHolder:", err)
		return
	}
	defer killHolder(holder)
	fmt.Printf("  grandchild pid=%d holds the write end open (no preamble, pipe stays empty)\n", holder.Process.Pid)

	type res struct {
		n       int
		err     error
		blocked time.Duration
	}
	done := make(chan res, 1)
	go func() {
		buf := make([]byte, 4096)
		start := time.Now()
		n, err := r.Read(buf) // blocks: empty pipe, writer still open
		done <- res{n, err, time.Since(start)}
	}()

	// Let the reader genuinely park in the blocking Read.
	time.Sleep(300 * time.Millisecond)
	select {
	case got := <-done:
		fmt.Printf("  UNEXPECTED: Read returned before any deadline was set: n=%d err=%s\n", got.n, describe(got.err))
		return
	default:
	}
	fmt.Println("  reader is parked in Read (did not return within 300ms)")

	setAt := time.Now()
	sderr := r.SetReadDeadline(time.Now())
	fmt.Printf("  SetReadDeadline(time.Now()) returned: %s\n", describe(sderr))

	select {
	case got := <-done:
		fmt.Printf("  Read RETURNED %v after SetReadDeadline\n", time.Since(setAt))
		fmt.Printf("    n=%d err=%s\n", got.n, describe(got.err))
		fmt.Printf("    total time blocked in Read: %v\n", got.blocked)
		fmt.Printf("    VERDICT S1: unblocked=%v deadlineExceeded=%v\n",
			true, errors.Is(got.err, os.ErrDeadlineExceeded))
	case <-time.After(3 * time.Second):
		fmt.Println("  FAIL: Read did NOT return within 3s of SetReadDeadline -- A2 REFUTED for the deadline mechanism")
	}
	fmt.Println()
}

// ---- Point 4a: buffered bytes + live quiet writer, deadline set, then Read ----

func scenarioBufferedBytesAfterDeadline() {
	fmt.Println("=== S2: buffered bytes present, writer alive but quiet, deadline set BEFORE the Read ===")
	fmt.Println("    (point 4a -- load-bearing for C1 'bytes still buffered are delivered')")
	const preamble = "BUFFERED-PAYLOAD-0123456789"
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Println("  os.Pipe:", err)
		return
	}
	defer r.Close()
	holder, err := spawnHolder(w, preamble, 4*time.Second)
	if err != nil {
		fmt.Println("  spawnHolder:", err)
		return
	}
	defer killHolder(holder)
	fmt.Printf("  grandchild pid=%d wrote %d bytes then went quiet, write end STILL OPEN\n",
		holder.Process.Pid, len(preamble))

	// Give the grandchild time to land its bytes in the pipe buffer.
	time.Sleep(500 * time.Millisecond)

	sderr := r.SetReadDeadline(time.Now())
	fmt.Printf("  SetReadDeadline(time.Now()) returned: %s\n", describe(sderr))

	buf := make([]byte, 4096)
	start := time.Now()
	n, rerr := r.Read(buf)
	el := time.Since(start)
	fmt.Printf("  Read after the deadline: n=%d elapsed=%v err=%s\n", n, el, describe(rerr))
	if n > 0 {
		fmt.Printf("    bytes delivered: %q\n", string(buf[:n]))
	}
	fmt.Printf("    RESULT 4a: buffered bytes %s (expected %d, got %d)\n",
		map[bool]string{true: "DELIVERED", false: "LOST -- deadline error preempted them"}[n == len(preamble)],
		len(preamble), n)

	// Follow-up: with the deadline still in the past and the buffer now drained,
	// the next Read should be the deadline error (that is the drain's stop signal).
	start = time.Now()
	n2, rerr2 := r.Read(buf)
	fmt.Printf("  subsequent Read (buffer now empty, deadline still past): n=%d elapsed=%v err=%s\n",
		n2, time.Since(start), describe(rerr2))
	fmt.Println()
}

// ---- Point 4a variant: io.Copy-style drain loop, the shape the RDR uses ----

func scenarioDrainLoopUnderDeadline() {
	fmt.Println("=== S3: realistic drain loop -- writer streams, then deadline fires mid-drain ===")
	fmt.Println("    (does the loop keep the bytes it already accumulated?)")
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Println("  os.Pipe:", err)
		return
	}
	defer r.Close()
	const preamble = "CHUNK-A|CHUNK-B|CHUNK-C|"
	holder, err := spawnHolder(w, preamble, 4*time.Second)
	if err != nil {
		fmt.Println("  spawnHolder:", err)
		return
	}
	defer killHolder(holder)

	type res struct {
		acc string
		err error
		el  time.Duration
	}
	done := make(chan res, 1)
	go func() {
		var acc []byte
		buf := make([]byte, 8) // small: forces several Read calls
		start := time.Now()
		for {
			n, err := r.Read(buf)
			acc = append(acc, buf[:n]...) // accumulate BEFORE inspecting err
			if err != nil {
				done <- res{string(acc), err, time.Since(start)}
				return
			}
		}
	}()

	time.Sleep(500 * time.Millisecond) // let the chunks arrive and the loop park
	setAt := time.Now()
	_ = r.SetReadDeadline(time.Now())
	select {
	case got := <-done:
		fmt.Printf("  drain loop ended %v after SetReadDeadline\n", time.Since(setAt))
		fmt.Printf("    accumulated=%q (%d bytes, expected %d)\n", got.acc, len(got.acc), len(preamble))
		fmt.Printf("    terminating err=%s\n", describe(got.err))
		fmt.Printf("    RESULT S3: partial output %s\n",
			map[bool]string{true: "PRESERVED", false: "TRUNCATED/LOST"}[got.acc == preamble])
	case <-time.After(3 * time.Second):
		fmt.Println("  FAIL: drain loop did not end within 3s")
	}
	fmt.Println()
}

// ---- Point 6: the Close-on-read-end variant ----

func scenarioCloseOnBlockedRead() {
	fmt.Println("=== S4: Close() on the read end while a Read is blocked (alternative mechanism) ===")
	fmt.Println("    (point 6)")
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Println("  os.Pipe:", err)
		return
	}
	holder, err := spawnHolder(w, "", 4*time.Second)
	if err != nil {
		fmt.Println("  spawnHolder:", err)
		return
	}
	defer killHolder(holder)
	fmt.Printf("  grandchild pid=%d holds the write end open (pipe empty)\n", holder.Process.Pid)

	type res struct {
		n   int
		err error
	}
	done := make(chan res, 1)
	go func() {
		buf := make([]byte, 4096)
		n, err := r.Read(buf)
		done <- res{n, err}
	}()
	time.Sleep(300 * time.Millisecond)
	select {
	case got := <-done:
		fmt.Printf("  UNEXPECTED: Read returned before Close: n=%d err=%s\n", got.n, describe(got.err))
		return
	default:
	}
	fmt.Println("  reader is parked in Read")

	closeAt := time.Now()
	cerr := r.Close()
	fmt.Printf("  Close() returned after %v: %s\n", time.Since(closeAt), describe(cerr))
	select {
	case got := <-done:
		fmt.Printf("  Read RETURNED %v after Close: n=%d err=%s\n", time.Since(closeAt), got.n, describe(got.err))
	case <-time.After(3 * time.Second):
		fmt.Println("  FAIL: Read did NOT return within 3s of Close")
	}
	fmt.Println()
}

// ---- Point 6 variant: does Close lose bytes that were already buffered? ----

func scenarioCloseWithBufferedBytes() {
	fmt.Println("=== S5: Close() on the read end with UNREAD buffered bytes ===")
	fmt.Println("    (point 6 vs point 4 -- what the Close mechanism costs you)")
	const preamble = "BUFFERED-PAYLOAD-0123456789"
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Println("  os.Pipe:", err)
		return
	}
	holder, err := spawnHolder(w, preamble, 4*time.Second)
	if err != nil {
		fmt.Println("  spawnHolder:", err)
		return
	}
	defer killHolder(holder)
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("  %d bytes sit unread in the pipe; now Close() the read end\n", len(preamble))
	_ = r.Close()
	buf := make([]byte, 4096)
	n, rerr := r.Read(buf)
	fmt.Printf("  Read after Close: n=%d err=%s\n", n, describe(rerr))
	fmt.Printf("    RESULT S5: buffered bytes are %s by Close\n",
		map[bool]string{true: "STILL READABLE", false: "UNREACHABLE"}[n > 0])
	fmt.Println()
}

// ---- Point 5: ErrNoDeadline / non-pollable file ----

func scenarioNoDeadline() {
	fmt.Println("=== S6: SetReadDeadline on a file whose poller registration did not happen ===")
	fmt.Println("    (point 5 -- failure mode F5 trigger)")

	// A regular file is the natural non-pollable case: kindOpenFile, not
	// registered with the runtime poller.
	f, err := os.CreateTemp("", "a2-nodeadline-*")
	if err != nil {
		fmt.Println("  CreateTemp:", err)
		return
	}
	defer os.Remove(f.Name())
	defer f.Close()
	derr := f.SetReadDeadline(time.Now().Add(time.Second))
	fmt.Printf("  regular file  SetReadDeadline -> %s\n", describe(derr))
	fmt.Printf("    errors.Is(err, os.ErrNoDeadline) = %v\n", errors.Is(derr, os.ErrNoDeadline))

	// The pipe case for contrast: a pollable os.Pipe read end must accept it.
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Println("  os.Pipe:", err)
		return
	}
	defer r.Close()
	defer w.Close()
	perr := r.SetReadDeadline(time.Now().Add(time.Second))
	fmt.Printf("  os.Pipe read end SetReadDeadline -> %s\n", describe(perr))
	fmt.Printf("    errors.Is(err, os.ErrNoDeadline) = %v\n", errors.Is(perr, os.ErrNoDeadline))

	// A pipe fd re-adopted via os.NewFile is still kindPipe (newFile stats the
	// fd), so this is NOT expected to be non-pollable -- recorded to show what a
	// re-adopted pipe actually does rather than to assert F5.
	r2, w2, err := os.Pipe()
	if err == nil {
		defer r2.Close()
		defer w2.Close()
		dup, derr2 := syscall.Dup(int(r2.Fd()))
		if derr2 == nil {
			nf := os.NewFile(uintptr(dup), "readopted-pipe")
			e := nf.SetReadDeadline(time.Now().Add(time.Second))
			fmt.Printf("  os.NewFile(dup of pipe read end) SetReadDeadline -> %s\n", describe(e))
			fmt.Printf("    errors.Is(err, os.ErrNoDeadline) = %v\n", errors.Is(e, os.ErrNoDeadline))
			nf.Close()
		}
	}
	fmt.Printf("  os.ErrNoDeadline value as declared by the stdlib: %q\n", os.ErrNoDeadline.Error())
	fmt.Println()
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "holder" {
		holderMain()
		return
	}
	fmt.Printf("RDR 0026 / A2 spike -- %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("started %s\n\n", time.Now().Format(time.RFC3339))

	scenarioDeadlineOnBlockedRead()
	scenarioBufferedBytesAfterDeadline()
	scenarioDrainLoopUnderDeadline()
	scenarioCloseOnBlockedRead()
	scenarioCloseWithBufferedBytes()
	scenarioNoDeadline()

	fmt.Println("done -- all grandchildren killed via their process groups")
}
