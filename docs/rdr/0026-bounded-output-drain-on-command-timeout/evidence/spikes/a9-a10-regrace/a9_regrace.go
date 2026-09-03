//go:build ignore

// Spike A9 for RDR 0026: the per-read deadline re-arm bounds the IDLE GAP
// between reads, not the SIZE of the tail.
//
// Background. A prior spike (a2-unblock-read, S2/T1) established that a
// deadline of NOW short-circuits before the syscall: Read returns n=0 with
// os.ErrDeadlineExceeded even when bytes are sitting in the pipe buffer.
// The design therefore does not deadline at NOW. Its final drain re-arms a
// SHORT FUTURE deadline of now+DrainGrace BEFORE EACH READ:
//
//	for {
//	    rd.SetReadDeadline(time.Now().Add(DrainGrace))
//	    n, err := rd.Read(buf)
//	    ...
//	}
//
// The claim under test is that this bounds only the gap the drain will WAIT
// for the next byte, and places no ceiling on how many bytes the drain may
// ultimately recover. A 1 MiB tail through a ~64 KiB pipe buffer takes 16+
// reads; the drain's total wall-clock therefore FAR exceeds DrainGrace while
// every individual read completes well inside it. That total exceeding the
// grace while the drain still reaches EOF with every byte IS the result.
//
// The CONTROL is the design this replaces: one ABSOLUTE deadline of
// now+DrainGrace, set ONCE and never re-armed. On a tail that takes longer
// than DrainGrace to walk, that deadline expires partway through and the
// drain reports os.ErrDeadlineExceeded — a FALSE "held pipe" verdict on a
// pipe whose writer has already exited. The control failing is what makes
// the re-arm load-bearing rather than decorative.
//
// Adversarial ordering. In both arms the drain goroutine is started LATE, on
// purpose: the driver sleeps past the point a real timeout timer would have
// fired before it reads the first byte. So the drain begins with the pipe
// buffer full and the child parked in write(2) (for payloads over ~64 KiB),
// or with the child already exited and the whole tail queued (for payloads
// under it). Either way the drain is reading strictly after the deadline
// moment, which is the situation the RDR's final drain is in.
//
// Two child SHAPES, and only the second one has teeth:
//
//	BURST — the child writes the whole payload in one Write and exits. By
//	  the time the late drain runs, every byte is already queued (or is in a
//	  parked write that completes the instant the drain makes room), so the
//	  drain walks 1 MiB in a few MILLISECONDS. No deadline is under stress
//	  and the two arms are INDISTINGUISHABLE. Reported anyway, because a
//	  reader who assumes a big tail is automatically a slow tail should see
//	  that it is not: the burst rows are the honest null result.
//
//	PACED — the child writes the payload in chunks with a real idle gap
//	  between them. Each gap is under DrainGrace; the gaps SUM to many times
//	  DrainGrace. This is the shape that separates the arms, and it is the
//	  realistic one: a child that emits output as it works, rather than all
//	  at once at the end. Here the re-armed deadline survives every gap and
//	  reaches EOF with the whole payload, while the single absolute deadline
//	  expires partway through and reports a FALSE held-pipe.
//
// The distinction the spike is really testing is therefore GAP vs TOTAL, not
// SMALL vs LARGE. A tail is only dangerous to a single absolute deadline
// when it takes wall-clock time to arrive; size alone does not do it.
package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"time"
)

// DrainGrace mirrors the RDR's proposed constant: the short future deadline
// re-armed before each read of the final drain.
const DrainGraceMS = 50

const drainGrace = DrainGraceMS * time.Millisecond

// lateStart is how long the driver withholds the drain goroutine. It stands
// in for "the timeout already fired and we are now in the final drain": long
// enough that the child has filled the pipe and parked in write(2).
const lateStart = 200 * time.Millisecond

// readBuf is the drain's per-read buffer. 32 KiB, smaller than the pipe
// buffer, so a 1 MiB tail is guaranteed to take many reads.
const readBuf = 32 << 10

var fixture string

func main() {
	fixture = os.Args[1]

	fmt.Printf("### RDR 0026 A9 spike: per-read deadline re-arm ###\n")
	fmt.Printf("platform: %s/%s  go: %s\n", runtime.GOOS, runtime.GOARCH, runtime.Version())
	fmt.Printf("DrainGrace: %v   late-start delay: %v   read buffer: %d bytes\n",
		drainGrace, lateStart, readBuf)
	fmt.Printf("pipe buffer (measured): %d bytes\n\n", measurePipeBuffer())

	runSizes()
	runTrials()
}

// --- arm 1: re-armed per-read grace, arm 2: single absolute deadline ------

type mode int

const (
	modeRegrace  mode = iota // SetReadDeadline(now+grace) BEFORE EACH read
	modeAbsolute             // SetReadDeadline(now+grace) ONCE, never re-armed
)

func (m mode) String() string {
	if m == modeRegrace {
		return "re-armed"
	}
	return "absolute"
}

type result struct {
	got      []byte
	reads    int // reads that returned n>0
	calls    int // total Read calls, including the terminal one
	elapsed  time.Duration
	terminal string // "EOF" | "deadline" | other
	deadline bool   // terminated on os.ErrDeadlineExceeded
	complete bool   // bytes.Equal against an independently rebuilt envelope
	maxGap   time.Duration // longest observed gap between consecutive reads
}

// drain runs the final-drain loop against rd in the given mode and reports
// what it recovered, how many reads it took, and how it terminated.
func drain(rd *os.File, m mode, limit int) result {
	var r result
	var acc bytes.Buffer
	buf := make([]byte, readBuf)

	start := time.Now()
	if m == modeAbsolute {
		// The control: ONE deadline, set once, never re-armed. Whatever the
		// tail costs, the whole loop must finish inside this single window.
		_ = rd.SetReadDeadline(start.Add(drainGrace))
	}
	for {
		if m == modeRegrace {
			// The design: re-arm a fresh short FUTURE deadline before every
			// read, so the bound applies to the gap before THIS read's first
			// byte, not to the cumulative cost of the tail.
			_ = rd.SetReadDeadline(time.Now().Add(drainGrace))
		}
		readStart := time.Now()
		n, err := rd.Read(buf)
		// The per-read gap: how long THIS read waited. The claim is that
		// this is what DrainGrace bounds. Compare it against r.elapsed,
		// which is what DrainGrace does NOT bound.
		if gap := time.Since(readStart); gap > r.maxGap {
			r.maxGap = gap
		}
		r.calls++
		if n > 0 {
			r.reads++
			acc.Write(buf[:n])
		}
		if err != nil {
			r.terminal = describe(err)
			r.deadline = errors.Is(err, os.ErrDeadlineExceeded)
			break
		}
		if acc.Len() > limit {
			r.terminal = "limit-exceeded"
			break
		}
	}
	r.elapsed = time.Since(start)
	r.got = acc.Bytes()
	return r
}

func describe(err error) string {
	switch {
	case errors.Is(err, io.EOF):
		return "EOF"
	case errors.Is(err, os.ErrDeadlineExceeded):
		return "deadline"
	case errors.Is(err, os.ErrClosed):
		return "closed"
	default:
		return fmt.Sprintf("other(%v)", err)
	}
}

// shape selects how the child delivers the payload.
type shape struct {
	name   string
	paced  bool
	chunk  int           // paced only: bytes per Write
	gap    time.Duration // paced only: idle time between chunks
}

// burst hands the whole payload over in one Write.
var burst = shape{name: "burst"}

// paced delivers the payload in chunks separated by an idle gap that is
// comfortably UNDER DrainGrace, so a re-armed deadline clears every gap
// individually while their SUM dwarfs the grace.
func pacedShape(chunk int, gap time.Duration) shape {
	return shape{name: "paced", paced: true, chunk: chunk, gap: gap}
}

// runOne spawns the fixture writing a size-byte envelope in the given shape,
// withholds the drain for lateStart, then drains in the requested mode.
func runOne(size int, m mode, sh shape) result {
	rd, wr, err := os.Pipe()
	must(err)

	var cmd *exec.Cmd
	if sh.paced {
		cmd = exec.Command(fixture, "paced", strconv.Itoa(size),
			strconv.Itoa(sh.chunk), strconv.Itoa(int(sh.gap.Milliseconds())))
	} else {
		cmd = exec.Command(fixture, "tail", strconv.Itoa(size))
	}
	cmd.Stdout = wr
	cmd.Stderr = nil
	must(cmd.Start())
	// Parent drops its copy of the write end: the child is the only writer,
	// so when the child exits the read end gets a real EOF.
	_ = wr.Close()

	// ADVERSARIAL: do not read for lateStart. A real timeout would have
	// fired by now; the pipe is full and the child is parked in write(2)
	// for any payload above the pipe buffer.
	time.Sleep(lateStart)

	r := drain(rd, m, 4<<20)
	_ = rd.Close()
	_ = cmd.Wait()

	r.complete = bytes.Equal(r.got, envelope(size))
	return r
}

// --- per-size sweep -------------------------------------------------------

var sizes = []struct {
	name string
	n    int
}{
	{"1KiB", 1 << 10},
	{"64KiB", 64 << 10},
	{"256KiB", 256 << 10},
	{"1MiB", 1 << 20},
}

// pacedFor gives each payload a chunking that produces a MULTI-SECOND
// delivery out of gaps that are each well under DrainGrace. 32 chunks with a
// 20 ms gap = ~620 ms of pure idle time, 12x DrainGrace.
func pacedFor(size int) shape {
	chunk := size / 32
	if chunk < 1 {
		chunk = 1
	}
	return pacedShape(chunk, 20*time.Millisecond)
}

func runSizes() {
	fmt.Printf("--- A9: drain started LATE, per-read grace re-armed vs set-once ---\n")
	fmt.Printf("Claim: every byte is recovered and the drain ends on EOF, even when the\n")
	fmt.Printf("drain's TOTAL wall-clock exceeds DrainGrace many times over, because the\n")
	fmt.Printf("grace bounds only the GAP before each read.\n\n")

	for _, sh := range []string{"burst", "paced"} {
		if sh == "burst" {
			fmt.Printf("SHAPE burst: child writes the whole payload in one Write, then exits.\n")
			fmt.Printf("(expected NULL result: the tail is already queued, so it drains in\n")
			fmt.Printf(" milliseconds and neither arm is under any deadline stress)\n\n")
		} else {
			fmt.Printf("SHAPE paced: child writes in 32 chunks with a 20ms idle gap between\n")
			fmt.Printf("them (~620ms of idle in total, %.0fx DrainGrace; each single gap is\n",
				float64(31*20*time.Millisecond)/float64(drainGrace))
			fmt.Printf("0.4x DrainGrace). This is the shape that separates the two arms.\n\n")
		}
		fmt.Printf("%-8s %-9s %-9s %-9s %-6s %-13s %-11s %-9s %s\n",
			"mode", "payload", "recovered", "complete", "reads", "elapsed", "max-gap", "terminal", "note")
		for _, s := range sizes {
			for _, m := range []mode{modeRegrace, modeAbsolute} {
				shp := burst
				if sh == "paced" {
					shp = pacedFor(s.n)
				}
				r := runOne(s.n, m, shp)
				note := ""
				if !r.complete {
					note = fmt.Sprintf("TRUNCATED %d/%d prefix-ok=%v",
						len(r.got), s.n, isPrefix(r.got, envelope(s.n)))
				}
				if r.deadline {
					note += " FALSE-HELD-REPORT"
				}
				if r.complete && r.elapsed > drainGrace {
					note = fmt.Sprintf("elapsed %.1fx grace, still EOF",
						float64(r.elapsed)/float64(drainGrace))
				}
				fmt.Printf("%-8s %-9d %-9d %-9v %-6d %-13v %-11v %-9s %s\n",
					m, s.n, len(r.got), r.complete, r.reads,
					r.elapsed.Round(time.Microsecond), r.maxGap.Round(time.Microsecond),
					r.terminal, note)
			}
		}
		fmt.Println()
	}
}

// --- n>=10 trials at 1 MiB ------------------------------------------------

func runTrials() {
	const iters = 12
	const size = 1 << 20

	fmt.Printf("--- A9: n=%d trials at 1 MiB, both shapes (a flake would show here) ---\n\n", iters)

	for _, sh := range []string{"burst", "paced"} {
		for _, m := range []mode{modeRegrace, modeAbsolute} {
			var elapsed, gaps []time.Duration
			var reads []int
			complete, eof, dl := 0, 0, 0
			shortest := -1
			for i := 0; i < iters; i++ {
				shp := burst
				if sh == "paced" {
					shp = pacedFor(size)
				}
				r := runOne(size, m, shp)
				elapsed = append(elapsed, r.elapsed)
				gaps = append(gaps, r.maxGap)
				reads = append(reads, r.reads)
				if r.complete {
					complete++
				}
				if r.terminal == "EOF" {
					eof++
				}
				if r.deadline {
					dl++
				}
				if shortest < 0 || len(r.got) < shortest {
					shortest = len(r.got)
				}
			}
			sort.Slice(elapsed, func(i, j int) bool { return elapsed[i] < elapsed[j] })
			sort.Slice(gaps, func(i, j int) bool { return gaps[i] < gaps[j] })
			sort.Ints(reads)
			fmt.Printf("shape=%s mode=%s payload=%d n=%d\n", sh, m, size, iters)
			fmt.Printf("  complete (bytes.Equal):     %d/%d\n", complete, iters)
			fmt.Printf("  terminated on EOF:          %d/%d\n", eof, iters)
			fmt.Printf("  terminated on deadline:     %d/%d\n", dl, iters)
			fmt.Printf("  fewest bytes recovered:     %d of %d\n", shortest, size)
			fmt.Printf("  reads       min/median/max: %d / %d / %d\n",
				reads[0], reads[len(reads)/2], reads[len(reads)-1])
			fmt.Printf("  elapsed     min/median/max: %v / %v / %v\n",
				elapsed[0].Round(time.Microsecond),
				elapsed[len(elapsed)/2].Round(time.Microsecond),
				elapsed[len(elapsed)-1].Round(time.Microsecond))
			fmt.Printf("  max per-read GAP observed:  %v  (DrainGrace = %v)\n",
				gaps[len(gaps)-1].Round(time.Microsecond), drainGrace)
			fmt.Printf("  max TOTAL vs DrainGrace:    %.1fx (%v vs %v)\n",
				float64(elapsed[len(elapsed)-1])/float64(drainGrace),
				elapsed[len(elapsed)-1].Round(time.Microsecond), drainGrace)
			switch {
			case sh == "burst":
				fmt.Printf("  => NULL: tail already queued, no deadline stress, arms indistinguishable\n")
			case m == modeRegrace:
				fmt.Printf("  => grace bounded the GAP (max gap < grace), NOT the TOTAL (total >> grace):\n")
				fmt.Printf("     whole payload recovered, EOF every trial, zero false held reports\n")
			default:
				fmt.Printf("  => the single absolute deadline FAILED: this is what the re-arm avoids\n")
			}
			fmt.Println()
		}
	}
}

// --- helpers --------------------------------------------------------------

// measurePipeBuffer reports how many bytes a fresh pipe accepts before a
// write would block, so the reads-per-tail numbers can be read in context.
func measurePipeBuffer() int {
	rd, wr, err := os.Pipe()
	if err != nil {
		return -1
	}
	defer rd.Close()
	defer wr.Close()
	// Non-blocking probe: give the write end a past deadline so a write that
	// would block returns instead of hanging.
	total := 0
	chunk := make([]byte, 4096)
	for total < 4<<20 {
		_ = wr.SetWriteDeadline(time.Now().Add(50 * time.Millisecond))
		n, err := wr.Write(chunk)
		total += n
		if err != nil {
			break
		}
	}
	return total
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

func must(err error) {
	if err != nil {
		panic(err)
	}
}
