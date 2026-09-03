//go:build ignore

// FX-tail-and-exit: the child of the A9 spike driver.
//
// It writes an envelope of exactly argv[2] bytes to stdout and EXITS. It
// spawns nothing, escapes nothing, and holds nothing: the write ends die
// with it, so the parent's read end WILL see EOF once the pipe drains.
//
// That is the whole point of A9. The tail is large (up to 1 MiB) but the
// writer is gone, so the only question the drain has to answer is whether a
// re-armed short deadline can walk a multi-read tail to EOF without one of
// those reads timing out. A held pipe would confound that; here there is no
// holder to blame.
//
// The child necessarily BLOCKS in write(2) until the parent starts draining
// whenever the envelope exceeds the pipe buffer (~64 KiB). The driver's A9
// case delays its drain on purpose, so at the moment the drain finally runs
// the pipe is full and the child is parked mid-write. That is the adversarial
// shape: the drain must then pull the remainder through in many reads.
//
// Modes (argv[1]):
//
//	tail   — write envelope(argv[2]) to stdout in one Write, then exit 0.
//	         The whole payload is handed over as fast as the pipe accepts it.
//
//	paced  — write envelope(argv[2]) to stdout in chunks of argv[3] bytes,
//	         sleeping argv[4] milliseconds BETWEEN chunks, then exit 0.
//
// The paced mode is what gives A9 its teeth. In tail mode a 1 MiB payload is
// already sitting in the pipe (or in a parked write) by the time the drain
// runs, so the drain walks it in a few milliseconds and NO deadline of any
// kind is under stress — the re-armed and single-absolute arms are then
// indistinguishable, and the spike would prove nothing.
//
// Paced mode separates the chunks by a real idle gap. Each gap is under
// DrainGrace, so a re-armed per-read deadline survives every one of them;
// but the gaps SUM to far more than DrainGrace, so a single absolute
// deadline set once at the start expires partway through and reports a
// false held-pipe. That contrast is the A9 result.
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "fixture: usage: fixture tail <size> | paced <size> <chunk> <gapms>")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "tail":
		size := atoi(os.Args[2])
		env := envelope(size)
		if _, err := os.Stdout.Write(env); err != nil {
			fmt.Fprintln(os.Stderr, "fixture: stdout write:", err)
			os.Exit(3)
		}
		// Exit of our own accord. Closing stdout here is what eventually
		// gives the parent's drain its EOF.
		os.Exit(0)

	case "paced":
		size := atoi(os.Args[2])
		chunk := atoi(os.Args[3])
		gap := time.Duration(atoi(os.Args[4])) * time.Millisecond
		env := envelope(size)
		for off := 0; off < len(env); off += chunk {
			end := off + chunk
			if end > len(env) {
				end = len(env)
			}
			if _, err := os.Stdout.Write(env[off:end]); err != nil {
				fmt.Fprintln(os.Stderr, "fixture: stdout write:", err)
				os.Exit(3)
			}
			if end < len(env) {
				time.Sleep(gap)
			}
		}
		os.Exit(0)

	default:
		fmt.Fprintln(os.Stderr, "fixture: unknown mode", os.Args[1])
		os.Exit(2)
	}
}

// envelope builds a deterministic payload of exactly size bytes that the
// driver reconstructs independently and compares byte-for-byte. Same shape as
// the a3-a6 spike's envelope, so the two spikes' payloads are comparable.
func envelope(size int) []byte {
	b := make([]byte, size)
	for i := range b {
		b[i] = byte('a' + (i % 26))
	}
	// Frame it so a truncation ANYWHERE is visible, not just at the tail.
	copy(b, "ENVELOPE-BEGIN:")
	if size >= 13 {
		copy(b[size-13:], ":ENVELOPE-END")
	}
	return b
}

func atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fixture: atoi:", err)
		os.Exit(2)
	}
	return n
}
