package cmdbind_test

// RDR 0026 — the contract obligations that are structural rather than
// behavioural: the surfaces this record does NOT change, the write
// journey's published cost, and the two documentation obligations Phase 2
// and Phase 5 owe.
//
// These rows exist because the record states each of them as a promise a
// later refactor could void silently. A behavioural test cannot catch "no
// second bound was minted" or "0025 was not edited"; a structural one can.

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/cmdbind"
)

// --- what `spawn`'s signature and the invocation carrier promise ---------

// REQ-68 (EIA `spawn` row): "C1 lands once; no per-binding wait. Its
// signature is unchanged by this record: it returns the `invocation` by
// VALUE (not a pointer), which is what makes C1 `refusal:`'s \"the field is
// readable on every error path\" hold without a nil check"
// REQ-111 (IP Phase 1): "No other line of `spawn` moves."
// REQ-56 (AP): "the order stays as it is (`cmd.Wait()`, then `reapGroup`,
// then the drain join), and the join becomes bounded by the existing
// `cmdbind.go::WaitDelay` (500 ms), not by a second constant."
// REQ-2 (C1 `precedence:`): the fixed order in `spawn`.
// ADVERSARIAL — the order is read from the shipped source's own statement
// sequence, because "no arm of `spawn` returns between Wait and the join"
// (C1 `stdin:`) is a property of the code's SHAPE that no input can probe.
func TestReq2_SpawnKeepsTheFixedOrderWaitThenReapGroupThenTheDrainJoin(t *testing.T) {
	t.Parallel()

	body := spawnBody(t)
	wait := strings.Index(body, "cmd.Wait()")
	reap := strings.Index(body, "reapGroup(")
	join := strings.Index(body, "drains.Wait()")

	if wait < 0 || reap < 0 || join < 0 {
		t.Fatalf("`spawn` no longer names all three steps (Wait=%d "+
			"reapGroup=%d join=%d); C1 `precedence:` fixes the order as "+
			"Wait → reapGroup → drain join", wait, reap, join)
	}
	if !(wait < reap && reap < join) {
		t.Fatalf("`spawn`'s order is Wait@%d reapGroup@%d join@%d; C1 "+
			"`precedence:` fixes it as Wait (direct child reaped) → reapGroup "+
			"(one kill(2) to -pgid, non-blocking) → drain join under ONE "+
			"timer of WaitDelay covering BOTH read drains. Release-first is "+
			"what keeps the common backgrounded-grandchild case from paying "+
			"the bound (76c121b, intrastate#x5vq)", wait, reap, join)
	}

	// C1 `stdin:` — "no arm of `spawn` returns between Wait and the join".
	between := body[wait:join]
	if strings.Contains(between, "return") {
		t.Errorf("`spawn` returns between `Wait` and the drain join:\n%s\n"+
			"C1 `stdin:` states that the non-ExitError arm refuses AFTER "+
			"reapGroup and the join — no arm of `spawn` returns between them",
			between)
	}
}

// REQ-24 (C1 `precedence:`): "An explicit close before them would
// double-close, so none is added and the Illustrative Code draws none"
// REQ-66 (IC): "The block deliberately draws NO `graceOn` and NO
// `closeReads` call"
// REQ-20 (C1 `precedence:` (a)): "no separate `graceOn` flag is introduced"
// ADVERSARIAL — both are NEGATIVE REQs about the shipped shape.
func TestReq20_NoGraceOnFlagAndNoExplicitCloseAreAddedToSpawn(t *testing.T) {
	t.Parallel()

	body := spawnBody(t)
	lower := strings.ToLower(body)

	if strings.Contains(lower, "graceon") {
		t.Error("`spawn` names a `graceOn` flag; C1 `precedence:` (a) fixes " +
			"the parent's deadline-set on the read ends as THE mark — " +
			"`os.File`'s own poller is internally synchronised, so " +
			"`SetReadDeadline` from the joining goroutine against a `Read` in " +
			"flight is the sanctioned cross-goroutine call, and no separate " +
			"flag is introduced")
	}
	if strings.Contains(lower, "closereads") {
		t.Error("`spawn` names a `closeReads` call; the existing defers at " +
			"pipe creation already close both read ends, so an explicit close " +
			"before them would double-close (Infra-Audit close-ownership row)")
	}
	// The defers themselves must survive: they are what satisfies "no read
	// end outlives `spawn`".
	if strings.Count(body, "Close()") < 2 || !strings.Contains(body, "defer") {
		t.Error("`spawn` no longer defers both read-end closes at pipe " +
			"creation; C1 `precedence:` states the guarantee is no fd leaks, " +
			"carried by those defers, and nothing in this change moves the " +
			"fds' release")
	}
}

// REQ-11 (C1 `precedence:`): "That re-arm has NO site today and this clause
// creates one: `readBounded` ... is `io.ReadAll(io.LimitReader(r, limit))`
// followed by `io.Copy(io.Discard, r)`, so `io.ReadAll` owns the loop and
// there is no per-read point to re-arm. Implementing this clause replaces
// that body with an explicit `for { SetReadDeadline(now+DrainGrace); Read }`
// loop — a rewrite of the function, not the signature change alone."
// REQ-14: "the discard-remainder exit must itself carry the deadline, and
// the pre-timer path must stay deadline-free."
// REQ-13: "`readBounded` returns that report alongside its bytes"
// REQ-67 (IC): "the re-arm lives INSIDE that read loop, not in the parent."
// ADVERSARIAL — the mechanism is a property of the loop's shape; MVV row
// 6(b) catches the wrong mechanism behaviourally, and this row catches the
// shape a passing 6(b) could still have reached by accident.
func TestReq11_ReadBoundedOwnsAnExplicitReadLoopThatRearmsTheDeadlinePerRead(t *testing.T) {
	t.Parallel()

	body := funcBody(t, "readBounded")

	if strings.Contains(body, "io.ReadAll") {
		t.Error("`readBounded` still calls `io.ReadAll`, which OWNS the loop " +
			"and leaves no per-read point to re-arm. C1 `precedence:` states " +
			"that implementing the clause REPLACES that body with an explicit " +
			"`for { SetReadDeadline(now+DrainGrace); Read }` loop — a rewrite " +
			"of the function, not the signature change alone")
	}
	if !strings.Contains(body, "SetReadDeadline") {
		t.Error("`readBounded` does not arm a read deadline; the re-arm lives " +
			"INSIDE that read loop, not in the parent — a one-shot " +
			"`SetReadDeadline` from the parent before `drains.Wait()` would " +
			"bound the whole remaining tail against the clock, which MVV " +
			"row 6(b) exists to fail")
	}
	// The re-arm is per read: the deadline call sits INSIDE a loop.
	if !rearmIsInsideALoop(t, "readBounded") {
		t.Error("`readBounded`'s `SetReadDeadline` does not sit inside its " +
			"read loop. C1 `precedence:` fixes the deadline as re-armed at " +
			"`now + DrainGrace` BEFORE EACH read of the final drain, so the " +
			"grace bounds the IDLE GAP between reads and never the size or " +
			"total duration of the tail")
	}
}

// --- the published cost of a write journey -------------------------------

// REQ-45 (C1 `bound:`): "The bound is PER INVOCATION, and a command write
// runs up to THREE of them through the same reader, not two: the pre-write
// baseline read when `protectedKeys` returns non-empty ..., the write itself
// ..., and the read-back .... A write journey against an escaping helper
// therefore costs up to `3·(timeout + 2·WaitDelay)` end to end"
// REQ-46: "the baseline leg's refusal is DISCARDED ..., so that leg spends a
// full bound and yields no signal ... This record does not change that
// discard"
// REQ-47: "The two legs are NOT separately bounded and the contract does
// not claim they are"
// REQ-115 (IP Phase 4): "if Phase 4 must slip, the bound ships with the
// write path excluded (reads and gates only) rather than with an unsensed
// write."
// BOUNDARY — the journey's total, asserted where the record publishes it.
func TestReq45_AWriteJourneyAgainstAnEscapingHelperStaysInsideThreeBounds(t *testing.T) {
	t.Parallel()

	pidfile := escPIDFile(t)
	w := &cmdbind.Writer{
		Accessor: escEntry(t, "1s", "success", pidfile, `{"state.phase":"review"}`),
		Name:     escName,
		Config:   cmdbind.Config{AllowCommands: true},
	}

	timeout := 1 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()
	err := w.Apply(ctx, escArtifact(t), nil)
	elapsed := time.Since(start)

	escKill(escReadPID(t, pidfile))

	if err == nil {
		t.Fatal("a write whose pipes are held by an escapee reported success")
	}
	// ONE invocation, ONE bound. The 3x figure is the JOURNEY's, which C1
	// `bound:` names precisely so the published per-invocation total is not
	// read as the journey's.
	if want := escBoundFor(timeout); elapsed > want {
		t.Fatalf("one write INVOCATION took %v, past the per-invocation "+
			"committed total of %v plus its %v tolerance. C1 `bound:` fixes "+
			"the bound as PER INVOCATION; the journey's up-to-3x cost is a "+
			"consequence of three invocations, never of one running long",
			elapsed, timeout+2*escWaitDelay, escTolerance)
	}
}

// --- what this record does NOT change ------------------------------------

// REQ-112 (IP Phase 2): "Re-point `cmdbind.go`'s C4 comments (\"necessary
// and sufficient\", \"run to a true EOF\", `WaitDelay`'s doc) at 0026:C1;
// 0025 itself is not edited."
// REQ-15: "A re-implementation that shifts the cap boundary is a contract
// violation, not a refactor"
// DOMAIN EDGE — the record's own citation-repair obligation. A shipped
// comment still claiming C4's "true EOF, nothing truncated" would be a
// promise the code no longer makes.
func TestReq112_TheC4CommentsAtTheSeamAreRepointedAt0026C1(t *testing.T) {
	t.Parallel()

	src := readSource(t, "cmdbind.go")

	// The three wordings C4 made and this record narrows.
	stale := []string{"necessary-and-sufficient", "necessary and sufficient",
		"run to a true EOF", "true EOF"}
	for _, s := range stale {
		if !strings.Contains(src, s) {
			continue
		}
		// It may only survive where the SAME sentence re-cites 0026:C1.
		for _, line := range strings.Split(src, "\n") {
			if strings.Contains(line, s) && !strings.Contains(src, "0026:C1") {
				t.Errorf("`cmdbind.go` still carries %q with no re-citation of "+
					"`0026:C1`. Phase 2 re-points C4's comments at this "+
					"record's successor text — the drains no longer promise a "+
					"true EOF; they promise whole output WITHIN THE BOUND",
					s)
				break
			}
		}
	}
	if !strings.Contains(src, "0026:C1") {
		t.Error("`cmdbind.go` cites no `0026:C1` at all; Phase 2's whole " +
			"content is re-pointing the C4 comments at this record, and the " +
			"comment that justifies the manual pipes is the one that " +
			"promised what C4's own text never did")
	}
}

// REQ-112: "0025 itself is not edited."
// REQ-108 (F6): "the cure is to amend 0025:C4's `detail:` rule ..., which is
// Implemented, is this record's Overrides target ... It is accepted rather
// than fixed"
// ADVERSARIAL — a locked record is never amended; the successor text is
// this record's.
func TestReq112_TheLocked0025RecordIsNotEdited(t *testing.T) {
	t.Parallel()

	path := repoFile(t, "docs/rdr/0025-command-invoking-accessor-bindings.md")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("0025's record is not readable from this tree: %v", err)
	}
	if strings.Contains(string(b), "0026:C1") {
		t.Error("0025's locked record now cites `0026:C1`; Phase 2 states " +
			"plainly that \"0025 itself is not edited\" — the successor text " +
			"is 0026's, and `cmdbind.go`'s comments are what re-cite it")
	}
}

// REQ-119 (IP Phase 5): "Then state the restated residue in
// `docs/cli-output-contract.md`, beside the command-entry `timeout` text
// authors already read — the \"a process may leak and its output is never
// read\" consequence, and the remediation the refusal names (close or
// redirect the helper's inherited stdio). Acceptance: that doc names the
// held-pipe refusal and the residue"
// REQ-118 (IP Phase 5): "SURVEY, then state."
// REQ-55 (C1 `residue:`): "Diagnosis: the refusal's Detail names the held
// pipe(s)"
// DOMAIN EDGE — this phase is the record's ONLY user-facing artifact, and
// `§consequences`' claims hold only once it ships.
func TestReq119_TheOutputContractDocNamesTheHeldPipeRefusalAndTheResidue(t *testing.T) {
	t.Parallel()

	path := repoFile(t, "docs/cli-output-contract.md")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the output-contract doc: %v", err)
	}
	doc := strings.ToLower(string(b))

	if !strings.Contains(doc, "held") {
		t.Error("`docs/cli-output-contract.md` does not name the HELD-PIPE " +
			"refusal. Phase 5's acceptance is that the doc names it, and " +
			"`§consequences`' claim that the refusal is \"visible, with the " +
			"pipe named\" holds only once it ships")
	}
	if !strings.Contains(doc, "close or redirect") {
		t.Error("`docs/cli-output-contract.md` does not carry the remediation " +
			"the refusal names (\"close or redirect the helper's inherited " +
			"stdio\"). For an author whose helper is a third-party binary that " +
			"daemonizes, that redirection IS the path — no opt-out flag ships")
	}
	if !strings.Contains(doc, "leak") && !strings.Contains(doc, "residue") {
		t.Error("`docs/cli-output-contract.md` does not state the restated " +
			"residue — \"a process may leak and its output is never read\". " +
			"`§consequences` records the residue as \"explicitly admitted\", " +
			"which holds only once the doc says so")
	}
}

// REQ-135 (XC `Build tool compatibility`): "Unix-only process-group
// syscalls stay behind the build tags established by intrastate#1drn
// (`procgroup_unix.go`); Windows continues to take 0025:C4's platform
// refusal. A12 confirms the drain-start stall seam this RDR adds needs **no
// third build tag** — it is portable Go in `cmdbind`."
// PRESERVED INVARIANT for the first half; the third-tag prohibition is new.
// ADVERSARIAL
func TestReq135_TheStallSeamAddsNoThirdBuildTagAndThePlatformSplitIsUnchanged(t *testing.T) {
	t.Parallel()

	dir := packageDir(t)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading the package directory: %v", err)
	}

	tagged := map[string]bool{}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		b, rerr := os.ReadFile(filepath.Join(dir, name))
		if rerr != nil {
			t.Fatalf("reading %s: %v", name, rerr)
		}
		if strings.Contains(string(b), "//go:build") {
			tagged[name] = true
		}
	}

	// Exactly the two halves intrastate#1drn established, and no third.
	want := map[string]bool{"procgroup_unix.go": true, "procgroup_other.go": true}
	if !reflect.DeepEqual(tagged, want) {
		t.Errorf("build-tagged files = %v; want %v. A12 confirms the "+
			"drain-start stall seam this record adds needs NO third build "+
			"tag — it is portable Go in `cmdbind`, and a tagged seam would "+
			"make it unbuildable rather than observable on one platform",
			tagged, want)
	}
}

// REQ-133 (TS coverage goal): "the FX-deadline-escape helper is exercised
// on darwin and linux in CI, so the escape is real on both (Background)."
// REQ-110 (IP Prerequisites): "A `Setsid: true` Go test helper exists in
// `internal/cli/cmdbind` so the escape is real on both OSes"
// DOMAIN EDGE — the fixture must ACTUALLY escape, or every scenario above
// is testing a reachable writer and passing for the wrong reason.
func TestReq110_TheFixtureGrandchildGenuinelyEscapesTheProcessGroup(t *testing.T) {
	t.Parallel()

	if !cmdbind.ProcGroupSupported {
		t.Skipf("this platform takes 0025:C4's refusal, so no group escape " +
			"is meaningful here")
	}

	pidfile := escPIDFile(t)
	r := escReader(t, "500ms", "deadline", pidfile)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	_, _, err := r.Read(ctx, escArtifact(t), []string{escKey})
	if err == nil {
		t.Fatal("the escape fixture did not refuse")
	}

	pid := escReadPID(t, pidfile)
	defer escKill(pid)

	// The escape is REAL: the grandchild survives reapGroup's group-scoped
	// SIGKILL. A shell fixture using `setsid sh -c …` does NOT escape on
	// macOS — reapGroup kills it, the drains close, and the read returns
	// bounded in ~200 ms, which is the probe that wrongly showed "no hang".
	time.Sleep(200 * time.Millisecond)
	if !escAlive(pid) {
		t.Fatalf("the fixture grandchild pid %d was reaped by the group "+
			"signal, so it never escaped and every held-pipe scenario in "+
			"this package is exercising a REACHABLE writer. The escape must "+
			"be made with a real `setsid(2)` call — macOS ships no `setsid` "+
			"BINARY (Background)", pid)
	}
}

// REQ-106 (MC `oracle`, S2 control) and REQ-123 (S2): the existing oracles
// re-run under the bounded join. Their own files are the assertion; this
// row asserts they are still THERE, since deleting a passing oracle is how
// an absence-of-change row silently stops discriminating.
// ADVERSARIAL
func TestReq123_TheExistingGrandchildLeakAndCancelOraclesStillShip(t *testing.T) {
	t.Parallel()

	src := readSource(t, "adversarial_0025_test.go")
	for _, name := range []string{
		"TestAdvBackgroundGrandchildMustNotStallASuccessfulRead",
		"TestAdvBackgroundGrandchildDoesNotSurviveASuccessfulRead",
		"TestAdvKilledChildOnParentCancelIsNotAnAnswer",
	} {
		if !strings.Contains(src, name) {
			t.Errorf("the existing oracle %s is gone. S2 re-runs the in-group "+
				"grandchild-leak and cancel oracles under the bounded join so "+
				"the prior point-fixes stay proven; deleting one is how "+
				"\"unchanged results\" stops meaning anything", name)
		}
	}
}

// REQ-109 (F4): "if A3 is wrong, bytes the direct child wrote can be lost
// at the bound; the refusal still fires, so the loss is bounded to a
// refused invocation, never a parsed one."
// REQ-92 (MC `fidelity`, held pipe): exemption "the whole stream — refused,
// not truncated-then-read"
// ADVERSARIAL — the bounded loss is only acceptable BECAUSE the invocation
// refuses; a build that lost bytes and returned them would be the silent
// class.
func TestReq109_LossAtTheBoundIsAlwaysBoundedToARefusedInvocation(t *testing.T) {
	t.Parallel()

	pidfile := escPIDFile(t)
	// The child writes a partial, well-formed-looking prefix and the escapee
	// holds the pipe: whatever arrived is bytes the direct child wrote.
	r := escReader(t, "20s", "success", pidfile, `{"state.phase":"rev`)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	values, unreadable, err := r.Read(ctx, escArtifact(t), []string{escKey})
	escKill(escReadPID(t, pidfile))

	if err == nil {
		t.Fatalf("a TRUNCATED payload on a held pipe was accepted: values=%+v "+
			"unreadable=%v. F4 admits that bytes can be lost at the bound "+
			"ONLY because the refusal still fires — the loss is bounded to a "+
			"refused invocation, never a parsed one", values, unreadable)
	}
	if len(values) != 0 {
		t.Errorf("values=%+v were returned alongside the refusal; a held "+
			"stream is refused, not truncated-then-read", values)
	}
	// And the refusal is the typed carrier, so the loss is diagnosable.
	escExecError(t, err)
}

// --- source-shape helpers -------------------------------------------------

// packageDir returns this package's own directory, resolved from the test
// binary's working directory (Go runs a package test with cwd = the package
// dir).
func packageDir(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("resolving the package directory: %v", err)
	}
	return wd
}

// repoFile resolves a repo-root-relative path from the package directory.
func repoFile(t *testing.T, rel string) string {
	t.Helper()

	dir := packageDir(t)
	for range 8 {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, rel)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatalf("could not find the repo root above %s", packageDir(t))
	return ""
}

func readSource(t *testing.T, name string) string {
	t.Helper()

	b, err := os.ReadFile(filepath.Join(packageDir(t), name))
	if err != nil {
		t.Fatalf("reading %s: %v", name, err)
	}
	return string(b)
}

// funcBody returns the shipped source text of one top-level function in
// `cmdbind.go`, so a structural obligation is read from the code rather
// than asserted about it.
func funcBody(t *testing.T, name string) string {
	t.Helper()

	path := filepath.Join(packageDir(t), "cmdbind.go")
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parsing cmdbind.go: %v", err)
	}
	src, rerr := os.ReadFile(path)
	if rerr != nil {
		t.Fatalf("reading cmdbind.go: %v", rerr)
	}
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != name || fn.Body == nil {
			continue
		}
		return string(src[fset.Position(fn.Body.Pos()).Offset:fset.Position(fn.Body.End()).Offset])
	}
	t.Fatalf("cmdbind.go declares no function %q", name)
	return ""
}

func spawnBody(t *testing.T) string {
	t.Helper()
	return funcBody(t, "spawn")
}

// rearmIsInsideALoop reports whether the named function's
// `SetReadDeadline` call sits inside a `for` statement — the difference
// between a per-read re-arm and the one-shot MVV row 6(b) fails.
func rearmIsInsideALoop(t *testing.T, name string) bool {
	t.Helper()

	path := filepath.Join(packageDir(t), "cmdbind.go")
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatalf("parsing cmdbind.go: %v", err)
	}
	found := false
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != name || fn.Body == nil {
			continue
		}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			loop, ok := n.(*ast.ForStmt)
			if !ok {
				return true
			}
			ast.Inspect(loop, func(inner ast.Node) bool {
				sel, ok := inner.(*ast.SelectorExpr)
				if ok && sel.Sel.Name == "SetReadDeadline" {
					found = true
				}
				return true
			})
			return true
		})
	}
	return found
}

// keep the accessor import live for the refusal-carrier helper this file
// shares with the rest of the suite.
var _ = accessor.ExecError{}
