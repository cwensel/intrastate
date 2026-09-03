//go:build ignore

// Live spike for RDR 0026 assumptions A3 (FIFO / no-loss) and A6 (the total
// bound), against the FX-deadline-escape fixture.
//
// The driver mirrors internal/cli/cmdbind/cmdbind.go spawn() order:
//
//	exec.CommandContext(ctx, ...) with a context DEADLINE
//	setProcGroup (Setpgid: true)
//	cmd.Cancel = kill(-pgid, SIGKILL)
//	cmd.WaitDelay = WaitDelay (read from cmdbind: 500ms)
//	manual os.Pipe ends on cmd.Stdout / cmd.Stderr
//	Start, close parent write ends
//	two concurrent bounded drain goroutines
//	cmd.Wait()
//	reapGroup: kill(-pgid, SIGKILL)
//	BOUNDED drain join under one more WaitDelay timer
//	  on the timer: SetReadDeadline(now) on the READ ends
//	join, close
//
// The only deliberate deviation from the shipped code is that the drain join
// here is BOUNDED (the shipped code's drains.Wait() is unbounded); bounding
// it is exactly what RDR 0026 proposes and what A6 measures.
package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// WaitDelay mirrors cmdbind.WaitDelay — 500, in milliseconds
// (internal/cli/cmdbind/cmdbind.go:66).
const WaitDelayMS = 500

const waitDelay = WaitDelayMS * time.Millisecond

// StdoutCap mirrors cmdbind.StdoutCap.
const StdoutCap = 1 << 20

var fixture string

func main() {
	fixture = os.Args[1]

	fmt.Printf("=== RDR 0026 A3/A6 live spike ===\n")
	fmt.Printf("platform: %s/%s  go: %s\n", runtime.GOOS, runtime.GOARCH, runtime.Version())
	fmt.Printf("WaitDelay: %v (cmdbind.WaitDelay = %d ms)\n\n", waitDelay, WaitDelayMS)

	runEscapeCheck()
	runA3()
	runA6()
}

// --- escape-is-real ------------------------------------------------------

// runEscapeCheck proves the grandchild really LEFT the child's process
// group: the fixture reports both pgids on fd 3, and the driver then shows
// that the group SIGKILL the real spawn() sends did not reach it.
func runEscapeCheck() {
	fmt.Printf("--- escape-is-real check ---\n")
	r := runOne(shapeB, 1<<10, false)
	fmt.Printf("fixture fd-3 report: %s", r.report)
	fmt.Printf("grandchild alive AFTER kill(-pgid, SIGKILL)? %v\n",
		r.escapeeAliveAfterGroupKill)
	fmt.Printf("(the drain therefore saw no EOF from the pipe: bounded join fired = %v)\n\n",
		r.boundFired)
}

// --- A3 -------------------------------------------------------------------

type a3Case struct {
	name    string
	size    int
	delayed bool // hold the drain goroutine until AFTER the bound would fire
}

func runA3() {
	fmt.Printf("--- A3: FIFO / no-loss at the bound (shape-B) ---\n")
	fmt.Printf("Claim: at the moment the bound expires, the accumulated buffer\n")
	fmt.Printf("holds the direct child's COMPLETE envelope, byte-for-byte.\n\n")

	cases := []a3Case{
		{"1KiB      prompt-drain", 1 << 10, false},
		{"1KiB      delayed-drain", 1 << 10, true},
		{"63KiB     prompt-drain", 63 << 10, false},
		{"63KiB     delayed-drain", 63 << 10, true},
		{"128KiB    prompt-drain", 128 << 10, false},
		{"128KiB    delayed-drain", 128 << 10, true},
		{"1MiB      prompt-drain", 1 << 20, false},
		{"1MiB      delayed-drain", 1 << 20, true},
	}

	fmt.Printf("%-26s %-8s %-9s %-9s %-8s %-7s %s\n",
		"case", "wrote", "got", "complete", "childexit", "wrtret", "note")
	for _, c := range cases {
		r := runOne(shapeB, c.size, c.delayed)
		want := envelope(c.size)
		complete := bytes.Equal(r.stdout, want)
		note := ""
		if !complete {
			note = fmt.Sprintf("MISMATCH prefix-ok=%v", isPrefix(r.stdout, want))
		}
		if !r.childWriteReturned {
			note += " child-BLOCKED-in-write"
		}
		fmt.Printf("%-26s %-8d %-9d %-9v %-8s %-7v %s\n",
			c.name, c.size, len(r.stdout), complete, r.exitDesc,
			r.childWriteReturned, note)
	}
	fmt.Println()
}

// --- A6 -------------------------------------------------------------------

func runA6() {
	fmt.Printf("--- A6: the total bound, escaped-grandchild case ---\n")
	fmt.Printf("Claim: the invocation returns within timeout + 2*WaitDelay = %v.\n\n",
		timeout+2*waitDelay)

	const iters = 10
	for _, shape := range []int{shapeA, shapeB} {
		name := "shape-A (deadline, child silent)"
		size := 0
		if shape == shapeB {
			name = "shape-B (success, child writes 1 KiB then exits 0)"
			size = 1 << 10
		}
		var a, b, c []time.Duration
		var descs []string
		for i := 0; i < iters; i++ {
			r := runOne(shape, size, false)
			a = append(a, r.deadlineToWait)
			b = append(b, r.waitToJoin)
			c = append(c, r.total)
			descs = append(descs, r.exitDesc)
		}
		fmt.Printf("%s, n=%d\n", name, iters)
		fmt.Printf("  %-42s %-10s %-10s %-10s\n", "measure", "min", "median", "max")
		report("(a) ctx deadline -> Wait() returned", a)
		report("(b) Wait() returned -> bounded join done", b)
		report("(c) total elapsed (spawn -> all closed)", c)
		fmt.Printf("  budget: timeout(%v) + 2*WaitDelay(%v) = %v\n",
			timeout, 2*waitDelay, timeout+2*waitDelay)
		fmt.Printf("  (a) <= WaitDelay?   %v   (max %v)\n", max(a) <= waitDelay, max(a))
		fmt.Printf("  (b) <= WaitDelay?   %v   (max %v)\n", max(b) <= waitDelay, max(b))
		fmt.Printf("  (c) <= budget?      %v   (max %v, margin %v)\n",
			max(c) <= timeout+2*waitDelay, max(c), timeout+2*waitDelay-max(c))
		fmt.Printf("  child dispositions: %v\n\n", uniq(descs))
	}
}

func report(label string, ds []time.Duration) {
	s := append([]time.Duration(nil), ds...)
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
	fmt.Printf("  %-42s %-10v %-10v %-10v\n", label, s[0], s[len(s)/2], s[len(s)-1])
}

func max(ds []time.Duration) time.Duration {
	m := ds[0]
	for _, d := range ds {
		if d > m {
			m = d
		}
	}
	return m
}

func uniq(ss []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range ss {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

// --- the modelled invocation ---------------------------------------------

const (
	shapeA = iota
	shapeB
)

// timeout is the executor deadline. Short on purpose.
const timeout = 1500 * time.Millisecond

// escapeeSleep is how long the grandchild holds the pipes. It must outlast
// timeout + 2*WaitDelay comfortably, or the case under test evaporates.
const escapeeSleep = 8

// childSleep is shape-A's own sleep — past the deadline.
const childSleep = 6

type result struct {
	stdout             []byte
	stderr             []byte
	deadlineToWait     time.Duration // (a)
	waitToJoin         time.Duration // (b)
	total              time.Duration // (c)
	exitDesc           string
	childWriteReturned bool
	report             string
	boundFired         bool
	escapeeAliveAfterGroupKill bool
}

func runOne(shape, size int, delayDrain bool) result {
	var r result

	// The fd-3 report channel: the fixture tells us its pgid, the
	// grandchild's pgid, and whether its stdout write RETURNED.
	repR, repW, err := os.Pipe()
	must(err)

	var argv []string
	if shape == shapeA {
		argv = []string{"shape-A", strconv.Itoa(escapeeSleep), "0", strconv.Itoa(childSleep)}
	} else {
		argv = []string{"shape-B", strconv.Itoa(escapeeSleep), strconv.Itoa(size), "0"}
	}

	start := time.Now()
	deadline := start.Add(timeout)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// --- mirror cmdbind.spawn --------------------------------------------
	cmd := exec.CommandContext(ctx, fixture, argv...)
	cmd.Stdin = bytes.NewReader([]byte("{}"))
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		if kerr := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); kerr != nil {
			return cmd.Process.Kill()
		}
		return nil
	}
	cmd.WaitDelay = waitDelay

	outR, outW, err := os.Pipe()
	must(err)
	errR, errW, err := os.Pipe()
	must(err)
	cmd.Stdout = outW
	cmd.Stderr = errW
	cmd.ExtraFiles = []*os.File{repW}

	must(cmd.Start())
	_ = outW.Close()
	_ = errW.Close()
	_ = repW.Close()

	// The report reader is unbounded and independent — it is NOT part of the
	// modelled contract, only spike instrumentation.
	// The fd-3 report channel is held by BOTH the child and the escaped
	// grandchild, so io.ReadAll on it would not return until the grandchild
	// exits — far too late to probe whether the group kill reached it. Read
	// incrementally into a mutex-guarded buffer instead.
	var repMu sync.Mutex
	var repBuf bytes.Buffer
	go func() {
		b := make([]byte, 512)
		for {
			n, err := repR.Read(b)
			if n > 0 {
				repMu.Lock()
				repBuf.Write(b[:n])
				repMu.Unlock()
			}
			if err != nil {
				return
			}
		}
	}()
	snapshotReport := func() string {
		repMu.Lock()
		defer repMu.Unlock()
		return repBuf.String()
	}

	// drainGate releases the adversarial stdout drain only after the direct
	// child is reaped and the group is signalled.
	drainGate := make(chan struct{})

	// The two concurrent bounded drains, as in spawn().
	var stdout, stderr []byte
	var drains sync.WaitGroup
	drains.Add(2)
	go func() {
		defer drains.Done()
		if delayDrain {
			// ADVERSARIAL ORDERING: do not read at all until the direct
			// child has been reaped AND the group SIGKILL has been sent —
			// i.e. start reading inside the bounded-join window, with the
			// child long gone. If FIFO holds, every byte the child wrote is
			// still queued in the pipe waiting for this first read.
			//
			// It must start BEFORE the join's SetReadDeadline lands, or the
			// deadline pre-empts the first read and the case measures a
			// reader that never ran rather than the FIFO property.
			<-drainGate
		}
		stdout = readBounded(outR, StdoutCap+1)
	}()
	go func() {
		defer drains.Done()
		stderr = readBounded(errR, StdoutCap)
	}()

	werr := cmd.Wait()
	waitReturned := time.Now()
	r.deadlineToWait = waitReturned.Sub(deadline)

	// reapGroup — the group SIGKILL after the direct child is reaped.
	if cmd.Process != nil && cmd.Process.Pid > 1 {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

	// Release the adversarial drain HERE: the child exited some time ago,
	// the group has been signalled, and only the escaped grandchild still
	// holds the write ends. Anything the drain now reads it read strictly
	// after the writer was gone.
	close(drainGate)

	// --- the BOUNDED drain join (RDR 0026's proposal) ---------------------
	joined := make(chan struct{})
	go func() { drains.Wait(); close(joined) }()
	timer := time.NewTimer(waitDelay)
	select {
	case <-joined:
		timer.Stop()
	case <-timer.C:
		r.boundFired = true
		// The bound expired: unblock the readers by deadlining the READ
		// ends. Every byte already in the pipe has been delivered by now if
		// FIFO holds — that is the A3 claim under test.
		now := time.Now()
		_ = outR.SetReadDeadline(now)
		_ = errR.SetReadDeadline(now)
		<-joined
	}
	joinDone := time.Now()
	r.waitToJoin = joinDone.Sub(waitReturned)

	_ = outR.Close()
	_ = errR.Close()
	r.total = time.Since(start)

	// Cleanup: the escaped grandchild is in its OWN session, so the group
	// kill above cannot reach it. Kill it explicitly by the pid the fixture
	// reported — leave no stray sleepers.
	r.report = snapshotReport()
	// Is the escapee still alive despite the group SIGKILL above? Signal 0
	// probes existence without delivering anything. This runs BEFORE any
	// cleanup, so a true here is the escape proving itself: the group kill
	// aimed at the child's pgid could not reach a process that setsid()
	// moved into its own group and session.
	if pid := escapeePid(r.report); pid > 1 {
		r.escapeeAliveAfterGroupKill = syscall.Kill(pid, 0) == nil
	}
	killEscapee(r.report)
	_ = repR.Close()

	r.stdout = stdout
	r.stderr = stderr
	r.exitDesc = describe(werr)
	r.childWriteReturned = bytes.Contains([]byte(r.report), []byte("child-write-returned"))
	return r
}

func describe(err error) string {
	if err == nil {
		return "exit0"
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		if st, ok := ee.Sys().(syscall.WaitStatus); ok && st.Signaled() {
			return "sig:" + st.Signal().String()
		}
		return "exit" + strconv.Itoa(ee.ExitCode())
	}
	return err.Error()
}

func readBounded(r io.Reader, limit int) []byte {
	b, err := io.ReadAll(io.LimitReader(r, int64(limit)))
	if err != nil {
		return b
	}
	_, _ = io.Copy(io.Discard, r)
	return b
}

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

func isPrefix(got, want []byte) bool {
	return len(got) <= len(want) && bytes.Equal(got, want[:len(got)])
}

// killEscapee parses "grandchild pid=N" out of the fd-3 report and SIGKILLs
// it, so no sleeper survives the spike.
func escapeePid(rep string) int {
	const key = "grandchild pid="
	i := bytes.Index([]byte(rep), []byte(key))
	if i < 0 {
		return -1
	}
	rest := rep[i+len(key):]
	j := 0
	for j < len(rest) && rest[j] >= '0' && rest[j] <= '9' {
		j++
	}
	pid, err := strconv.Atoi(rest[:j])
	if err != nil {
		return -1
	}
	return pid
}

func killEscapee(rep string) {
	if pid := escapeePid(rep); pid > 1 {
		_ = syscall.Kill(pid, syscall.SIGKILL)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
