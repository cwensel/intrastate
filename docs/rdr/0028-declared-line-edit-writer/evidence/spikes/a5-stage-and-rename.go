//go:build ignore

// a5-stage-and-rename spike: stage a replacement file beside the target and
// rename over it, exercising the actual os.Rename + os.Chmod path an
// implementation would use (not just shell `mv`).
//
// Usage:
//
//	a5-stage-and-rename <target> <new-content-file> -copy-mode   # copy ORIGINAL file's mode (A5's proposal)
//	a5-stage-and-rename <target> <new-content-file> -fixed-0600  # fixed 0600, like flowbind.go::save
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: a5-stage-and-rename <target> <new-content-file> -copy-mode|-fixed-0600")
		os.Exit(2)
	}
	target := os.Args[1]
	contentFile := os.Args[2]
	mode := os.Args[3]

	content, err := os.ReadFile(contentFile)
	must(err)

	var origMode os.FileMode = 0o600
	if mode == "-copy-mode" {
		fi, err := os.Stat(target)
		must(err)
		origMode = fi.Mode().Perm()
		fmt.Printf("original mode read via stat: %o\n", origMode)
	} else if mode != "-fixed-0600" {
		fmt.Fprintln(os.Stderr, "unknown mode flag:", mode)
		os.Exit(2)
	}

	dir := filepath.Dir(target)
	tmp, err := os.CreateTemp(dir, ".a5-stage-*")
	must(err)
	staged := tmp.Name()
	defer func() { _ = os.Remove(staged) }()

	_, err = tmp.Write(content)
	must(err)
	must(tmp.Close())

	must(os.Chmod(staged, origMode))

	// Resolve the target path through symlinks first, so a rename over a
	// symlink path lands on the REAL file rather than replacing the
	// symlink itself. This is what contract C3's "symlinks resolved" buys.
	resolved, err := filepath.EvalSymlinks(target)
	if err != nil {
		// Target may not exist yet as a symlink chain (fresh file) — fall
		// back to the raw path.
		resolved = target
	}

	must(os.Rename(staged, resolved))
	fmt.Printf("staged %s -> renamed over resolved target %s (mode=%s)\n", staged, resolved, mode)
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
