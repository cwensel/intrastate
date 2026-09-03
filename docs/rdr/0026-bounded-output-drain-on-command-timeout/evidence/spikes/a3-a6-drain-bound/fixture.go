//go:build ignore

// FX-deadline-escape: the child of the spike driver.
//
// It spawns a GRANDCHILD with SysProcAttr{Setsid: true} that inherits BOTH
// stdout and stderr and then sleeps, so the grandchild leaves the child's
// process group and holds the pipe write ends past any group-directed
// SIGKILL aimed at the child's group.
//
// macOS ships no setsid(1), so a shell fixture cannot produce this shape;
// the escape must be made from Go via SysProcAttr.
//
// Modes (argv[1]):
//
//	grandchild  — the re-exec'd escapee: report its pgid on fd 3, then sleep.
//	shape-A     — deadline shape: spawn escapee, write NOTHING, sleep past
//	              the deadline.
//	shape-B     — success shape: spawn escapee, write the whole envelope to
//	              stdout, exit 0. The escapee keeps holding both pipes.
//
// argv[2] for shape-A/shape-B is the sleep seconds for the escapee.
// argv[3] for shape-B is the envelope payload size in bytes.
// argv[4] for shape-A is the child's own sleep seconds.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "fixture: need a mode")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "grandchild":
		grandchild()
	case "shape-A":
		escapeeSleep := atoi(os.Args[2])
		selfSleep := atoi(os.Args[4])
		spawnEscapee(escapeeSleep)
		// Write NOTHING. Sleep past the deadline.
		time.Sleep(time.Duration(selfSleep) * time.Second)
		os.Exit(0)
	case "shape-B":
		escapeeSleep := atoi(os.Args[2])
		size := atoi(os.Args[3])
		spawnEscapee(escapeeSleep)
		// Write the WHOLE envelope, then exit 0 of our own accord.
		env := envelope(size)
		if _, err := os.Stdout.Write(env); err != nil {
			fmt.Fprintln(os.Stderr, "fixture: stdout write:", err)
			os.Exit(3)
		}
		// Mark on fd 3 that the write RETURNED, so the driver can tell a
		// child that completed its write from one still blocked in write(2).
		if f := os.NewFile(3, "report"); f != nil {
			fmt.Fprintf(f, "child-write-returned %d\n", len(env))
			_ = f.Close()
		}
		os.Exit(0)
	default:
		fmt.Fprintln(os.Stderr, "fixture: unknown mode", os.Args[1])
		os.Exit(2)
	}
}

// envelope builds a deterministic payload of exactly size bytes that the
// driver can reconstruct independently and compare byte-for-byte.
func envelope(size int) []byte {
	b := make([]byte, size)
	for i := range b {
		b[i] = byte('a' + (i % 26))
	}
	// Frame it so a truncation anywhere is visible, not just at the tail.
	copy(b, "ENVELOPE-BEGIN:")
	if size >= 13 {
		copy(b[size-13:], ":ENVELOPE-END")
	}
	return b
}

// spawnEscapee re-execs this same binary in "grandchild" mode with
// Setsid: true, handing it BOTH inherited output pipes.
func spawnEscapee(sleepSec int) {
	self, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "fixture: executable:", err)
		os.Exit(4)
	}
	cmd := exec.Command(self, "grandchild", strconv.Itoa(sleepSec))
	// BOTH ends inherited — this is the whole point of the fixture.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// fd 3 in the grandchild is the driver's report pipe (our fd 3).
	if f := os.NewFile(3, "report"); f != nil {
		cmd.ExtraFiles = []*os.File{f}
	}
	// Setsid LEAVES the child's process group (and its session), so a
	// kill(-pgid) aimed at the child's group does not reach it.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "fixture: start grandchild:", err)
		os.Exit(5)
	}
	// Report the ESCAPE evidence: our own pgid vs the grandchild's.
	selfPgid, _ := syscall.Getpgid(os.Getpid())
	gcPgid, gerr := syscall.Getpgid(cmd.Process.Pid)
	if f := os.NewFile(3, "report"); f != nil {
		fmt.Fprintf(f, "child pid=%d pgid=%d | grandchild pid=%d pgid=%d err=%v\n",
			os.Getpid(), selfPgid, cmd.Process.Pid, gcPgid, gerr)
	}
	// Do NOT Wait: the grandchild must outlive us.
}

// grandchild holds both inherited pipes open for its sleep, writing nothing
// to them, so the parent's drain sees no EOF while it lives.
func grandchild() {
	sleepSec := 5
	if len(os.Args) > 2 {
		sleepSec = atoi(os.Args[2])
	}
	// Belt and braces: die on our own even if the driver's cleanup misses.
	time.Sleep(time.Duration(sleepSec) * time.Second)
	os.Exit(0)
}

func atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fixture: atoi:", err)
		os.Exit(2)
	}
	return n
}
