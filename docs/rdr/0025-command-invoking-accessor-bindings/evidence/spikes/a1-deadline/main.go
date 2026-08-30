// a1spike: RDR 0025 assumption A1 -- can a declared timeout bound a command
// invocation including the child's process tree, a non-reading child's stdin
// pipe, and post-kill pipe drain?
//
// Tree tracking: 200ms after Start we snapshot every descendant of the child
// (by ppid walk over `ps`). After Wait returns (or the watchdog fires) we
// re-check which of those pids are still alive with the same command line.
// That catches orphans re-parented to launchd, which a pgid or marker scan
// would miss for the non-Setpgid scenarios.
package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	timeout    = 1 * time.Second
	waitDelay  = 500 * time.Millisecond
	watchdog   = 5 * time.Second
	snapshotAt = 200 * time.Millisecond
)

type fixings struct {
	setpgid   bool
	groupKill bool
	waitDelay time.Duration
}

var s2fix = fixings{setpgid: true, groupKill: true, waitDelay: waitDelay}

type proc struct {
	pid, ppid, pgid int
	cmd             string
}

func psAll() []proc {
	out, err := exec.Command("ps", "-axo", "pid=,ppid=,pgid=,command=").Output()
	if err != nil {
		panic(err)
	}
	var ps []proc
	for _, line := range strings.Split(string(out), "\n") {
		f := strings.Fields(line)
		if len(f) < 4 {
			continue
		}
		pid, _ := strconv.Atoi(f[0])
		ppid, _ := strconv.Atoi(f[1])
		pgid, _ := strconv.Atoi(f[2])
		ps = append(ps, proc{pid, ppid, pgid, strings.Join(f[3:], " ")})
	}
	return ps
}

// tree returns root plus all its descendants.
func tree(root int) []proc {
	all := psAll()
	var out []proc
	seen := map[int]bool{}
	var walk func(int)
	walk = func(p int) {
		for _, x := range all {
			if x.pid == p && !seen[p] {
				seen[p] = true
				out = append(out, x)
			}
			if x.ppid == p && !seen[x.pid] {
				walk(x.pid)
			}
		}
	}
	walk(root)
	return out
}

func alive(snap []proc) []proc {
	all := psAll()
	var out []proc
	for _, s := range snap {
		for _, x := range all {
			if x.pid == s.pid && x.cmd == s.cmd {
				out = append(out, x)
			}
		}
	}
	return out
}

func fmtProcs(ps []proc) string {
	var b strings.Builder
	b.WriteByte('[')
	for i, p := range ps {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d(ppid=%d,pgid=%d,%q)", p.pid, p.ppid, p.pgid, p.cmd)
	}
	b.WriteByte(']')
	return b.String()
}

var leftovers []proc

func runScenario(name string, fx fixings, stdin []byte, argv ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	if fx.setpgid {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}
	if fx.groupKill {
		cmd.Cancel = func() error {
			return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
	}
	cmd.WaitDelay = fx.waitDelay
	if stdin != nil {
		cmd.Stdin = bytes.NewReader(stdin)
	}

	// cmd.Output wires stdout to a pipe + copying goroutine (the hazard path).
	start := time.Now()
	type done struct {
		out []byte
		err error
	}
	ch := make(chan done, 1)
	started := make(chan struct{})
	go func() {
		// Output() calls Start internally; we need the pid for the snapshot,
		// so poll Process until set.
		go func() {
			for cmd.Process == nil {
				time.Sleep(time.Millisecond)
			}
			close(started)
		}()
		out, err := cmd.Output()
		ch <- done{out, err}
	}()

	var snap []proc
	select {
	case <-started:
		time.Sleep(snapshotAt)
		snap = tree(cmd.Process.Pid)
	case d := <-ch: // fast child (S6) may finish before we snapshot
		ch <- d
	}

	var err error
	var out string
	blocked := false
	select {
	case d := <-ch:
		err = d.err
		out = strings.TrimSpace(string(d.out))
	case <-time.After(watchdog - time.Since(start)):
		blocked = true
		err = errors.New("WATCHDOG: Wait did not return within 5s")
	}
	elapsed := time.Since(start)
	time.Sleep(100 * time.Millisecond) // let a just-killed tree get reaped
	surv := alive(snap)
	leftovers = append(leftovers, surv...)

	errStr := "<nil>"
	if err != nil {
		errStr = strconv.Quote(err.Error())
	}
	pids := make([]int, len(surv))
	for i, p := range surv {
		pids[i] = p.pid
	}
	fmt.Printf("%s elapsed=%dms err=%s waitErr_is_DeadlineExceeded=%t waitDelayHit=%t survivors=%d %v\n",
		name, elapsed.Milliseconds(), errStr,
		errors.Is(err, context.DeadlineExceeded), errors.Is(err, exec.ErrWaitDelay),
		len(surv), pids)
	fmt.Printf("    ctxErr=%v blockedPastWatchdog=%t tree@%s=%s\n",
		ctx.Err(), blocked, snapshotAt, fmtProcs(snap))
	if len(surv) > 0 {
		fmt.Printf("    survivors=%s\n", fmtProcs(surv))
	}
	if out != "" {
		fmt.Printf("    output=%q\n", out)
	}
}

func main() {
	fmt.Printf("go=%s os=%s arch=%s timeout=%s waitDelay=%s watchdog=%s\n---\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, timeout, waitDelay, watchdog)

	// NOTE: macOS /bin/sh execs a lone simple command in place (no grandchild),
	// so `sh -c 'sleep 30'` would not exercise the grandchild-holds-pipe case.
	// A trailing `; :` forces the shell to fork sleep as a real grandchild.
	// The marker in the shell string lands in sh's command line for pkill.
	mk := func(s string) string {
		return fmt.Sprintf("sleep 30; : A1SPIKE_%s_%d", s, os.Getpid())
	}
	big := bytes.Repeat([]byte("x"), 1<<20) // 1 MiB > 64 KiB pipe buffer

	runScenario("S1", fixings{}, nil, "sh", "-c", mk("S1"))
	runScenario("S2", s2fix, nil, "sh", "-c", mk("S2"))
	runScenario("S3", fixings{waitDelay: waitDelay}, nil, "sh", "-c", mk("S3"))
	runScenario("S4", s2fix, big, "sleep", "30")
	runScenario("S5", s2fix, nil, "sleep", "30")
	runScenario("S6", s2fix, nil, "echo", `{"k":"v"}`)

	fmt.Println("---")
	fmt.Printf("cleanup: killing %d leftover tree pids %s\n", len(leftovers), fmtProcs(leftovers))
	for _, p := range leftovers {
		_ = syscall.Kill(p.pid, syscall.SIGKILL)
	}
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("post-cleanup survivors=%d\n", len(alive(leftovers)))
}
