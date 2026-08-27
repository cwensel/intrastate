package cli

// RDR 0011 — oracles over the diff-base resolution the kernel-untouched and
// production-diff-scope guards depend on (`0011:S5`, MVV 9, REQ-94,
// REQ-119, REQ-120).
//
// The defect these close: all three guards used to resolve their base with
// a bare `git merge-base HEAD main` and `t.Skip` on error, which is
// fail-OPEN in both directions this repo's CI actually runs. On a shallow
// PR checkout (`actions/checkout@v6` at the default `fetch-depth: 1`) local
// `main` does not exist, `merge-base` errors, and every guard skips. On a
// push to `main` the merge base IS HEAD, so the diff is empty and each
// assertion passes however much the commit changed `internal/resolve`.
//
// Both directions report green, so no oracle written over the guards' own
// verdicts can see them. What is asserted here instead is the RESOLUTION
// DECISION, extracted as the pure `resolveDiffBase` and exercised over
// throwaway repositories that stage each environment by construction —
// the guard is checked in the shapes it was written for rather than only in
// whichever shape the suite happens to run in.

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// gitAt runs one git command in `dir` and fails the test on error. The
// fixtures below build real repositories rather than mocking git: the
// subject is what `git merge-base` answers over a given ref topology, and a
// fake would assert the fixture's model of git instead of git.
func gitAt(t *testing.T, dir string, args ...string) string {
	t.Helper()

	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.invalid",
		"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.invalid",
		"GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s in %s: %v\n%s", strings.Join(args, " "), dir, err, out)
	}
	return strings.TrimSpace(string(out))
}

// commitAt writes `name` and commits it, returning the new HEAD.
func commitAt(t *testing.T, dir, name string) string {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(name), 0o600); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	gitAt(t, dir, "add", name)
	gitAt(t, dir, "commit", "-m", name)
	return gitAt(t, dir, "rev-parse", "HEAD")
}

// newDiffBaseRepo returns an initialized repository on branch `main`
// carrying one commit.
func newDiffBaseRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	gitAt(t, dir, "init", "--initial-branch=main", ".")
	commitAt(t, dir, "root.txt")
	return dir
}

// TestReq94_ADiffBaseThatCannotResolveUnderCIIsAFailureNotASkip is RED
// against the pre-fix inline code, which could only `t.Skip` here — the
// shallow-PR shape, and the one this repo's CI actually runs.
func TestReq94_ADiffBaseThatCannotResolveUnderCIIsAFailureNotASkip(t *testing.T) {
	dir := newDiffBaseRepo(t)
	// Leave `main` behind: no candidate in the ladder names a commit.
	gitAt(t, dir, "checkout", "-b", "feature")
	commitAt(t, dir, "work.txt")
	gitAt(t, dir, "branch", "-D", "main")

	_, _, err := resolveDiffBase(dir, "true", "")
	if !errors.Is(err, errNoDiffBase) {
		t.Fatalf("resolveDiffBase over a repo with no `main` = %v; want "+
			"errNoDiffBase. A shallow PR checkout carries no local `main`, "+
			"and a diff claim with no branch point is an UNCHECKED claim, "+
			"not a satisfied one — the skip that used to happen here is "+
			"exactly how `0011:S5` went unenforced in every CI mode", err)
	}
}

// TestReq94_ADiffBaseEqualToHEADIsAFailureInEveryEnvironment is the
// vacuous-pass direction: on a push to the default branch HEAD is its own
// merge base, so `git diff <base>` is empty whatever the commit changed.
// Same idiom as DEV-5's vacuity guard, and unconditional — a base of HEAD
// is a broken harness, not an untouched kernel, in CI and locally alike.
func TestReq94_ADiffBaseEqualToHEADIsAFailureInEveryEnvironment(t *testing.T) {
	dir := newDiffBaseRepo(t)
	head := commitAt(t, dir, "kernel.txt")

	for _, ci := range []string{"", "true"} {
		base, _, err := resolveDiffBase(dir, ci, "")
		if !errors.Is(err, errDiffBaseIsHEAD) {
			t.Errorf("ci=%q: resolveDiffBase on `main` at HEAD (%s) = "+
				"(%q, %v); want errDiffBaseIsHEAD. Every diff taken "+
				"against HEAD is empty, so the assertion reports green "+
				"without checking anything", ci, head[:8], base, err)
		}
	}
}

// TestReq94_ADiffBaseResolvesToTheTrueMergeBaseOnAnOrdinaryBranch is the
// GREEN case, and the one that keeps a developer's own `make check`
// working: in a normal clone `main` exists and names the branch point.
func TestReq94_ADiffBaseResolvesToTheTrueMergeBaseOnAnOrdinaryBranch(t *testing.T) {
	dir := newDiffBaseRepo(t)
	branchPoint := commitAt(t, dir, "shared.txt")
	gitAt(t, dir, "checkout", "-b", "feature")
	commitAt(t, dir, "work.txt")

	base, ref, err := resolveDiffBase(dir, "true", "")
	if err != nil {
		t.Fatalf("resolveDiffBase on an ordinary branch: %v", err)
	}
	if base != branchPoint {
		t.Errorf("base = %q; want the branch point %q", base, branchPoint)
	}
	if ref != "main" {
		t.Errorf("ref = %q; want `main`. The local branch leads the ladder "+
			"deliberately: `origin/main` may be many commits stale in an "+
			"ordinary clone, and a stale base widens the diff to every file "+
			"landed since the last fetch — failing REQ-120's allowlist and "+
			"the kernel-diff oracles for no defect at all", ref)
	}
}

// TestReq94_ADiffBaseFallsBackToOriginMainWhenNoLocalBranchExists covers a
// checkout carrying the remote ref but no local branch — a detached or
// single-branch clone.
func TestReq94_ADiffBaseFallsBackToOriginMainWhenNoLocalBranchExists(t *testing.T) {
	upstream := newDiffBaseRepo(t)
	branchPoint := commitAt(t, upstream, "shared.txt")

	dir := t.TempDir()
	gitAt(t, dir, "clone", "--no-local", upstream, ".")
	gitAt(t, dir, "checkout", "-b", "feature")
	commitAt(t, dir, "work.txt")
	// Drop the local branch; `origin/main` remains.
	gitAt(t, dir, "branch", "-D", "main")

	base, ref, err := resolveDiffBase(dir, "true", "")
	if err != nil {
		t.Fatalf("resolveDiffBase with only `origin/main`: %v", err)
	}
	if base != branchPoint {
		t.Errorf("base = %q; want the branch point %q", base, branchPoint)
	}
	if ref != "origin/main" {
		t.Errorf("ref = %q; want `origin/main`", ref)
	}
}

// TestReq94_ADiffBaseThatCannotResolveOffCIIsASkipSentinel keeps the
// local-developer path open: the same unresolvable repo returns the SAME
// sentinel, and only the caller's CI check decides fatal-vs-skip. The
// decision is one branch in `diffBase`, not a second resolution policy.
func TestReq94_ADiffBaseThatCannotResolveOffCIIsASkipSentinel(t *testing.T) {
	dir := newDiffBaseRepo(t)
	gitAt(t, dir, "checkout", "-b", "feature")
	commitAt(t, dir, "work.txt")
	gitAt(t, dir, "branch", "-D", "main")

	_, _, err := resolveDiffBase(dir, "", "")
	if !errors.Is(err, errNoDiffBase) {
		t.Fatalf("resolveDiffBase off CI = %v; want errNoDiffBase — the "+
			"sentinel is environment-independent so a developer on an odd "+
			"checkout stays unblocked while CI cannot pass without a base",
			err)
	}
	if errors.Is(err, errDiffBaseIsHEAD) {
		t.Error("the unresolvable case reports errDiffBaseIsHEAD; the two " +
			"sentinels route to different verdicts and must stay distinct")
	}
}

// TestReq94_TheGithubBaseRefExtendsTheLadderWhenItNamesAFetchedRef covers
// the PR shape where the workflow supplies the target branch by name and
// enough history has been fetched for it to resolve.
func TestReq94_TheGithubBaseRefExtendsTheLadderWhenItNamesAFetchedRef(t *testing.T) {
	dir := newDiffBaseRepo(t)
	branchPoint := commitAt(t, dir, "shared.txt")
	gitAt(t, dir, "branch", "release")
	gitAt(t, dir, "checkout", "-b", "feature")
	commitAt(t, dir, "work.txt")
	gitAt(t, dir, "branch", "-D", "main")

	base, ref, err := resolveDiffBase(dir, "true", "release")
	if err != nil {
		t.Fatalf("resolveDiffBase with GITHUB_BASE_REF=release: %v", err)
	}
	if base != branchPoint || ref != "release" {
		t.Errorf("(base, ref) = (%q, %q); want (%q, \"release\")",
			base, ref, branchPoint)
	}
}

// TestReq120_TheDiffBaseIsResolvedInExactlyOnePlace is the source-scan half.
// The defect was that THREE sites each re-derived the base inline, so a fix
// to one left the other two fail-open; the guard against that regressing is
// that `merge-base` appears once across the 0011 test surface. The repo
// already ships this idiom (`escape_shape_0009_test.go` scans shipped source
// for a contract artifact).
func TestReq120_TheDiffBaseIsResolvedInExactlyOnePlace(t *testing.T) {
	dir := filepath.Join(repoRootFor(t), "internal", "cli")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}

	var sites []string
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), "_0011_test.go") {
			continue
		}
		src, rerr := os.ReadFile(filepath.Join(dir, e.Name()))
		if rerr != nil {
			t.Fatalf("read %s: %v", e.Name(), rerr)
		}
		// Count the git INVOCATION, not the word: this file names
		// `merge-base` in prose and in its own failure message, and a scan
		// over the bare token would assert its own text. Parsing gives the
		// call expressions directly, so a call wrapped across lines — which
		// the surviving one is — counts once and a comment counts never.
		fset := token.NewFileSet()
		f, perr := parser.ParseFile(fset, filepath.Join(dir, e.Name()),
			src, 0)
		if perr != nil {
			t.Fatalf("parse %s: %v", e.Name(), perr)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || !isGitCommand(call) {
				return true
			}
			for _, arg := range call.Args {
				if lit, lok := arg.(*ast.BasicLit); lok &&
					lit.Kind == token.STRING &&
					lit.Value == `"merge-base"` {
					sites = append(sites, e.Name()+":"+strconv.Itoa(
						fset.Position(call.Pos()).Line))
				}
			}
			return true
		})
	}

	if len(sites) != 1 {
		t.Errorf("`merge-base` is invoked at %d sites across the 0011 test "+
			"surface (%v); want exactly 1. Three inline copies is what made "+
			"`0011:S5`'s kernel-untouched claim unenforced in every CI mode "+
			"— each skipped on its own, and fixing one left the others "+
			"fail-open. The ladder lives in `resolveDiffBase` alone",
			len(sites), sites)
	}
}

// isGitCommand reports whether `call` is `exec.Command("git", …)`.
func isGitCommand(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Command" {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok || pkg.Name != "exec" || len(call.Args) == 0 {
		return false
	}
	lit, ok := call.Args[0].(*ast.BasicLit)
	return ok && lit.Kind == token.STRING && lit.Value == `"git"`
}
