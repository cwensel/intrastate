// Follow-up spike for RDR 0026 / A2 point 4a.
//
// The first spike showed that with a deadline already in the past, Read returns
// os.ErrDeadlineExceeded with n=0 even though bytes sit in the pipe buffer.
// Question this program answers: are those bytes LOST, or merely DEFERRED --
// i.e. does clearing the deadline let a later Read still retrieve them?
// That difference decides whether the RDR's C1 "bytes still buffered are
// delivered" is achievable, and by what drain discipline.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

func holderMain() {
	preamble := os.Args[2]
	ms, _ := strconv.Atoi(os.Args[3])
	w := os.NewFile(3, "we")
	if preamble != "" {
		_, _ = w.Write([]byte(preamble))
	}
	if rf := os.Getenv("A2_READY_FILE"); rf != "" {
		_ = os.WriteFile(rf, []byte(strconv.Itoa(os.Getpid())), 0o600)
	}
	time.Sleep(time.Duration(ms) * time.Millisecond)
	os.Exit(0)
}

func spawnHolder(w *os.File, preamble string, lifetime time.Duration) (*exec.Cmd, error) {
	self, _ := os.Executable()
	rf, err := os.CreateTemp("", "a2b-ready-*")
	if err != nil {
		return nil, err
	}
	p := rf.Name()
	rf.Close()
	os.Remove(p)
	cmd := exec.Command(self, "holder", preamble, strconv.Itoa(int(lifetime.Milliseconds())))
	cmd.Env = append(os.Environ(), "A2_READY_FILE="+p)
	cmd.ExtraFiles = []*os.File{w}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	w.Close()
	dl := time.Now().Add(3 * time.Second)
	for time.Now().Before(dl) {
		if b, err := os.ReadFile(p); err == nil && len(b) > 0 {
			os.Remove(p)
			return cmd, nil
		}
		time.Sleep(5 * time.Millisecond)
	}
	return nil, fmt.Errorf("holder not ready")
}

func kill(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process != nil {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		cmd.Process.Kill()
		cmd.Process.Wait()
	}
}

func d(err error) string {
	switch {
	case err == nil:
		return "<nil>"
	case errors.Is(err, os.ErrDeadlineExceeded):
		return "i/o timeout [os.ErrDeadlineExceeded]"
	case errors.Is(err, os.ErrClosed):
		return "file already closed [os.ErrClosed]"
	default:
		return fmt.Sprintf("%v [%T]", err, err)
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "holder" {
		holderMain()
		return
	}
	fmt.Printf("A2 follow-up: are buffered bytes LOST or DEFERRED? -- %s %s/%s\n\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)

	const payload = "BUFFERED-PAYLOAD-0123456789"

	// --- T1: deadline in the past, then CLEAR it, then Read again. ---
	fmt.Println("=== T1: past deadline -> timeout, then SetReadDeadline(zero) -> can we still get the bytes? ===")
	r, w, _ := os.Pipe()
	h, err := spawnHolder(w, payload, 4*time.Second)
	if err != nil {
		fmt.Println("  spawn:", err)
		return
	}
	time.Sleep(400 * time.Millisecond)
	r.SetReadDeadline(time.Now())
	buf := make([]byte, 4096)
	n1, e1 := r.Read(buf)
	fmt.Printf("  Read with past deadline:      n=%d err=%s\n", n1, d(e1))
	// Now clear the deadline entirely and retry.
	if err := r.SetReadDeadline(time.Time{}); err != nil {
		fmt.Println("  clear deadline err:", err)
	}
	n2, e2 := r.Read(buf)
	fmt.Printf("  Read after clearing deadline: n=%d err=%s\n", n2, d(e2))
	if n2 > 0 {
		fmt.Printf("    recovered bytes: %q\n", string(buf[:n2]))
	}
	fmt.Printf("  => buffered bytes were %s\n",
		map[bool]string{true: "DEFERRED (recoverable after clearing the deadline)",
			false: "LOST (not recoverable)"}[n2 == len(payload)])
	r.Close()
	kill(h)
	fmt.Println()

	// --- T2: FUTURE deadline instead of a past one. Does Read drain first? ---
	fmt.Println("=== T2: deadline set to a FUTURE instant; bytes already buffered ===")
	r2, w2, _ := os.Pipe()
	h2, err := spawnHolder(w2, payload, 5*time.Second)
	if err != nil {
		fmt.Println("  spawn:", err)
		return
	}
	time.Sleep(400 * time.Millisecond)
	r2.SetReadDeadline(time.Now().Add(750 * time.Millisecond))
	n3, e3 := r2.Read(buf)
	fmt.Printf("  Read #1 (bytes buffered, deadline 750ms out): n=%d err=%s\n", n3, d(e3))
	if n3 > 0 {
		fmt.Printf("    delivered: %q\n", string(buf[:n3]))
	}
	// Buffer now empty, writer alive and quiet: this Read must hit the deadline.
	st := time.Now()
	n4, e4 := r2.Read(buf)
	fmt.Printf("  Read #2 (buffer empty, writer alive): n=%d after %v err=%s\n",
		n4, time.Since(st).Round(time.Millisecond), d(e4))
	fmt.Printf("  => a FUTURE deadline %s\n",
		map[bool]string{true: "DRAINS buffered bytes first, then times out (this is the C1-compatible discipline)",
			false: "did not deliver the buffered bytes"}[n3 == len(payload)])
	r2.Close()
	kill(h2)
	fmt.Println()

	// --- T3: does a past deadline lose bytes that arrive LATER? (ordering (b) sanity) ---
	fmt.Println("=== T3: past deadline set while pipe EMPTY, writer alive -- point 4b ===")
	r3, w3, _ := os.Pipe()
	h3, err := spawnHolder(w3, "", 3*time.Second)
	if err != nil {
		fmt.Println("  spawn:", err)
		return
	}
	r3.SetReadDeadline(time.Now())
	st = time.Now()
	n5, e5 := r3.Read(buf)
	fmt.Printf("  Read on empty pipe w/ past deadline: n=%d after %v err=%s\n",
		n5, time.Since(st).Round(time.Microsecond), d(e5))
	fmt.Printf("  => %s\n", map[bool]string{true: "deadline error as expected (point 4b CONFIRMED)",
		false: "unexpected"}[n5 == 0 && errors.Is(e5, os.ErrDeadlineExceeded)])
	r3.Close()
	kill(h3)
	fmt.Println()
	fmt.Println("done")
}
