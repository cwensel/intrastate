package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type result struct {
	stdout, stderr []byte
	exit           int // -1 when no exit code (spawn error)
	spawnErr       error
	dur            time.Duration
}

// run invokes argv[0] with argv[1:], no shell. env==nil means inherit; env==[]string{} means scrubbed.
func run(argv []string, stdin string, env []string) result {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	if env != nil {
		cmd.Env = env
	}
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	var so, se bytes.Buffer
	cmd.Stdout, cmd.Stderr = &so, &se
	t := time.Now()
	err := cmd.Run()
	r := result{stdout: so.Bytes(), stderr: se.Bytes(), exit: -1, dur: time.Since(t)}
	var ee *exec.ExitError
	switch {
	case err == nil:
		r.exit = 0
	case errors.As(err, &ee):
		r.exit = ee.ExitCode()
	default:
		r.spawnErr = err
	}
	return r
}

func show(label string, argv []string, stdin string, env []string, r result) {
	fmt.Println("---")
	fmt.Printf("%s\n  argv=%q\n", label, argv)
	if stdin != "" {
		fmt.Printf("  stdin=%q\n", stdin)
	}
	if env != nil {
		fmt.Printf("  env=%q\n", redact(env))
	}
	if r.spawnErr != nil {
		fmt.Printf("  spawn_error=%v  is_ErrNotFound=%v  exit=(none)\n", r.spawnErr, errors.Is(r.spawnErr, exec.ErrNotFound))
	} else {
		fmt.Printf("  exit=%d\n", r.exit)
	}
	if len(r.stdout) > 200 {
		fmt.Printf("  stdout_len=%d stdout_head=%q\n", len(r.stdout), r.stdout[:60])
	} else {
		fmt.Printf("  stdout=%q\n", r.stdout)
	}
	if len(r.stderr) > 0 {
		s := r.stderr
		if len(s) > 600 {
			s = s[:600]
		}
		fmt.Printf("  stderr=%q\n", s)
	}
	// C3 envelope parse attempt (before exit classification)
	var obj map[string]string
	if err := json.Unmarshal(r.stdout, &obj); err != nil {
		fmt.Printf("  c3_read_envelope: NOT PARSEABLE (%v)\n", err)
	} else {
		keys := 0
		for range obj {
			keys++
		}
		fmt.Printf("  c3_read_envelope: OK keys=%d\n", keys)
	}
	fmt.Printf("  dur=%s\n", r.dur.Round(time.Microsecond))
}

// redact replaces the parent PATH value so the transcript does not carry the host's PATH.
func redact(env []string) []string {
	out := make([]string, len(env))
	for i, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			e = "PATH=<parent PATH>"
		}
		out[i] = e
	}
	return out
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func fileBytes(p string) string {
	b, _ := os.ReadFile(p)
	return fmt.Sprintf("%q", b)
}

func main() {
	tmp, err := os.MkdirTemp("", "a3spike")
	must(err)
	defer os.RemoveAll(tmp)
	spike, _ := os.Getwd()
	art := filepath.Join(tmp, "state.cfg")
	reset := func() { must(os.WriteFile(art, []byte("[state]\n\tphase = draft\n"), 0o644)) }
	reset()
	parentPath := os.Getenv("PATH")
	git, _ := exec.LookPath("git")
	fmt.Printf("artifact=%s\ngit=%s\n", art, git)

	fmt.Println("\n################ A3 direct-binding class")
	a := []string{"git", "config", "--file", art, "--get", "state.phase"}
	r := run(a, "", nil)
	show("R1 read: git config --get (single value)", a, "", nil, r)
	fmt.Printf("  RAW-mode bind: key=state.phase value=%q (TrimRight \\n)\n", strings.TrimRight(string(r.stdout), "\n"))

	a = []string{"git", "config", "--file", art, "--list"}
	show("R2a read: git config --list (raw k=v lines)", a, "", nil, run(a, "", nil))
	a = []string{filepath.Join(spike, "wrapper-list.sh"), art}
	show("R2b read: wrapper-list.sh (k=v -> JSON)", a, "", nil, run(a, "", nil))

	a = []string{"git", "config", "--file", art, "state.phase", "review"}
	show("R3a write: git config set (value in argv -- NOT expressible in fixed argv)", a, "", nil, run(a, "", nil))
	fmt.Printf("  artifact_after=%s\n", fileBytes(art))
	reset()
	a = []string{"git", "config", "--file", art, "state.phase"}
	stdin := `{"state.phase":"review"}`
	show("R3b write: git config with data on stdin only (does git read stdin?)", a, stdin, nil, run(a, stdin, nil))
	fmt.Printf("  artifact_after=%s\n", fileBytes(art))
	a = []string{filepath.Join(spike, "wrapper-write.sh"), art}
	show("R3c write: wrapper-write.sh (stdin JSON -> git config per key)", a, stdin, nil, run(a, stdin, nil))
	fmt.Printf("  artifact_after=%s\n", fileBytes(art))
	a = []string{"git", "config", "--file", art, "--get", "state.phase"}
	show("R3d read-back via R1", a, "", nil, run(a, "", nil))

	reset()
	a = []string{"cat", art}
	show("R4a P-1: cat {artifact} as WRITE with JSON stdin", a, stdin, nil, run(a, stdin, nil))
	fmt.Printf("  artifact_after=%s\n", fileBytes(art))
	a = []string{"tee", art}
	show("R4b P-1: tee {artifact} as WRITE with JSON stdin", a, stdin, nil, run(a, stdin, nil))
	fmt.Printf("  artifact_after=%s\n", fileBytes(art))
	a = []string{"git", "config", "--file", art, "--get", "state.phase"}
	show("R4c read-back after tee", a, "", nil, run(a, "", nil))
	reset()

	a = []string{filepath.Join(spike, "wrapper-big.sh"), art}
	show("R5 P-2: 2 MiB JSON envelope", a, "", nil, run(a, "", nil))

	a = []string{"git", "config", "--file", art, "--get", "missing.key"}
	show("R6 read-absent: git config --get missing.key", a, "", nil, run(a, "", nil))

	fmt.Println("\n################ A8 exit-code mapping")
	a = []string{"test", "-f", art}
	show("E1a test -f (exists)", a, "", nil, run(a, "", nil))
	a = []string{"test", "-f", filepath.Join(tmp, "nope")}
	show("E1b test -f (missing)", a, "", nil, run(a, "", nil))
	a = []string{"test", "--bogus", art}
	show("E1c test with bogus option (execution failure)", a, "", nil, run(a, "", nil))

	repo := filepath.Join(tmp, "repo")
	must(os.MkdirAll(repo, 0o755))
	for _, c := range [][]string{
		{"git", "-C", repo, "init", "-q"},
		{"git", "-C", repo, "-c", "user.name=s", "-c", "user.email=s@x", "-c", "commit.gpgsign=false", "commit", "-q", "--allow-empty", "-m", "init"},
	} {
		if rr := run(c, "", nil); rr.exit != 0 {
			fmt.Printf("setup %q failed: %s\n", c, rr.stderr)
		}
	}
	must(os.WriteFile(filepath.Join(repo, "f"), []byte("a\n"), 0o644))
	run([]string{"git", "-C", repo, "add", "f"}, "", nil)
	run([]string{"git", "-C", repo, "-c", "user.name=s", "-c", "user.email=s@x", "-c", "commit.gpgsign=false", "commit", "-q", "-m", "f"}, "", nil)
	a = []string{"git", "-C", repo, "diff", "--quiet"}
	show("E1d git diff --quiet (clean)", a, "", nil, run(a, "", nil))
	must(os.WriteFile(filepath.Join(repo, "f"), []byte("b\n"), 0o644))
	show("E1e git diff --quiet (dirty)", a, "", nil, run(a, "", nil))
	a = []string{"git", "-C", filepath.Join(tmp, "norepo"), "diff", "--quiet"}
	show("E1f git with bad path (execution failure)", a, "", nil, run(a, "", nil))
	a = []string{"git", "config", "--file", filepath.Join(tmp, "nofile.cfg"), "--get", "state.phase"}
	show("E1g git config --file <missing file> --get", a, "", nil, run(a, "", nil))
	a = []string{"definitely-not-a-binary-xyz", art}
	show("E1h missing executable (spawn error)", a, "", nil, run(a, "", nil))

	fmt.Println("\n################ A7 env policy")
	a = []string{"git", "config", "--file", art, "--get", "state.phase"}
	show("V1 Env=[] (scrubbed; bare argv0 'git')", a, "", []string{}, run(a, "", []string{}))
	a = []string{"git", "config", "--file", art, "--list", "--show-origin"}
	show("V1b Env=[] --list --show-origin (which files does git open?)", a, "", []string{}, run(a, "", []string{}))
	env := []string{"PATH=" + parentPath}
	a = []string{"git", "config", "--file", art, "--get", "state.phase"}
	show("V2 Env=[PATH]", a, "", env, run(a, "", env))

	home := filepath.Join(tmp, "home")
	must(os.MkdirAll(home, 0o755))
	must(os.WriteFile(filepath.Join(home, ".gitconfig"), []byte("[state]\n\tphase = FROM_HOME_GITCONFIG\n"), 0o644))
	env = []string{"PATH=" + parentPath, "HOME=" + home}
	show("V3a Env=[PATH,HOME=tmp w/ conflicting ~/.gitconfig] --get", a, "", env, run(a, "", env))
	env3 := append(append([]string{}, env...), "GIT_TRACE=1")
	show("V3b same + GIT_TRACE=1", a, "", env3, run(a, "", env3))
	a2 := []string{"git", "config", "--get", "state.phase"}
	show("V3c control: NO --file, HOME=tmp -> reads ~/.gitconfig", a2, "", env, run(a2, "", env))
	a3 := []string{"git", "config", "--file", art, "--list", "--show-origin"}
	show("V3d --file --list --show-origin under HOME=tmp", a3, "", env, run(a3, "", env))

	gcg := filepath.Join(tmp, "global.cfg")
	must(os.WriteFile(gcg, []byte("[state]\n\tphase = FROM_GIT_CONFIG_GLOBAL\n"), 0o644))
	env4 := []string{"PATH=" + parentPath, "HOME=" + home, "GIT_CONFIG_GLOBAL=" + gcg}
	show("V4a GIT_CONFIG_GLOBAL set, --file --get", a, "", env4, run(a, "", env4))
	show("V4a' GIT_CONFIG_GLOBAL set, NO --file --get (control)", a2, "", env4, run(a2, "", env4))
	env5 := []string{"PATH=" + parentPath, "HOME=" + home, "GIT_DIR=" + filepath.Join(repo, ".git")}
	show("V4b GIT_DIR set, --file --get", a, "", env5, run(a, "", env5))
	env6 := []string{"PATH=" + parentPath, "HOME=" + home, "GIT_CONFIG_COUNT=1", "GIT_CONFIG_KEY_0=state.phase", "GIT_CONFIG_VALUE_0=FROM_ENV_INJECTION"}
	show("V4c GIT_CONFIG_COUNT injection, --file --get", a, "", env6, run(a, "", env6))
	show("V4c' GIT_CONFIG_COUNT injection, --file --get-all", []string{"git", "config", "--file", art, "--get-all", "state.phase"}, "", env6, run([]string{"git", "config", "--file", art, "--get-all", "state.phase"}, "", env6))
	show("V4c'' GIT_CONFIG_COUNT injection, NO --file (control)", a2, "", env6, run(a2, "", env6))
	env7 := []string{"PATH=" + parentPath, "HOME=" + home, "GIT_CONFIG_PARAMETERS='state.phase=FROM_GIT_CONFIG_PARAMETERS'"}
	show("V4d GIT_CONFIG_PARAMETERS, --file --get", a, "", env7, run(a, "", env7))
	env8 := []string{"PATH=" + parentPath, "HOME=" + home, "GIT_CONFIG_SYSTEM=" + gcg, "GIT_CONFIG_NOSYSTEM=0"}
	show("V4e GIT_CONFIG_SYSTEM, --file --get", a, "", env8, run(a, "", env8))
	env9 := []string{"PATH=" + parentPath, "HOME=" + home, "GIT_EXEC_PATH=" + tmp}
	show("V4f GIT_EXEC_PATH=tmp (hijack helper dir), --file --get", a, "", env9, run(a, "", env9))
	// git config editor path: -e would launch $GIT_EDITOR; not run (interactive). Show env var reaches git via --show-scope of GIT_CONFIG_GLOBAL only.
	fmt.Println("---")
	fmt.Println("done")
}
