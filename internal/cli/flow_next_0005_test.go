package cli

// RDR 0005 — `flow next` (REQ-40..REQ-47, REQ-74, REQ-75, REQ-118), the
// narrowed read-accessor set it shares with `flow resolve` (REQ-35..REQ-39),
// and the gate-disposition mini-check that separates enumeration from
// selection (REQ-43, REQ-46).
//
// The load-bearing distinction across this file: `next` ENUMERATES. Every
// oracle that could be satisfied by a verb which merely selects is written
// so that selection-like behaviour FAILS it — a gate deny still exits 0,
// gates run only on opt-in, and a guard-excluded row's gate never runs.

import (
	"slices"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
)

// REQ-40: "flow next MUST return the legal recognized-outcome alphabet for
// the supplied state, plus candidate summaries containing source rule
// identity, required facts, unresolved guard/gate facts, and preview next
// tags, write targets, and clear keys when those can be read from
// normalized model data without evaluating missing facts."
// REQ-74: `next` data minimum — "`model` (id or path) and `revision`,
// `observed` tags, `owned` tags assembled from the readers, `readers[]` …,
// `outcomes[]` … and `candidates[]`."
// HAPPY PATH
func TestReq40And74_NextPayloadCarriesTheAlphabetAndCandidateSummaries(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")

	data := flowData(t, requireSuccess(t, "flow", "next", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--tag", "profile=mid", "--as=json"))

	for _, key := range []string{
		"model", "revision", "observed", "owned", "readers", "outcomes",
		"candidates",
	} {
		if _, ok := data[key]; !ok {
			t.Errorf("`next` payload carries no %q; the data minimum fixes "+
				"it. keys = %v", key, keysOf(data))
		}
	}

	// `outcomes[]` is the recognized-outcome ALPHABET: bare tags, nothing
	// else — RDR 0002's layout carries no per-outcome recognizer text.
	outcomes, ok := stringsAt(data, "outcomes")
	if !ok {
		t.Fatalf("`outcomes` is not an array of bare tags: %#v", data["outcomes"])
	}
	if len(outcomes) == 0 {
		t.Error("`outcomes` is empty; the fixture model declares three")
	}
	for _, o := range outcomes {
		if !slices.Contains([]string{"advance", "hold", "bail"}, o) {
			t.Errorf("`outcomes` carries %q, which the fixture model does "+
				"not declare", o)
		}
	}
}

// REQ-75: "Each candidate carries source rule identity, the outcome it
// belongs to, required facts, unresolved guard/gate facts, evaluated gate
// results when `--evaluate-gates` was given, and preview next tags, write
// targets, and clear keys when the normalized model exposes them without
// evaluating missing facts."
// REQ-118: "The contract requires candidate summaries read from normalized
// data and forbids evaluating missing facts".
// DOMAIN EDGE
func TestReq75And118_CandidateSummariesAreReadFromNormalizedModelData(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")

	data := flowData(t, requireSuccess(t, "flow", "next", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art), "--as=json"))

	candidates, ok := objectsAt(data, "candidates")
	if !ok {
		t.Fatalf("`candidates` is not an array of objects: %#v",
			data["candidates"])
	}
	if len(candidates) == 0 {
		t.Fatal("`candidates` is empty; the fixture model's rows are " +
			"reachable from the supplied state")
	}

	var sawAdvance bool
	for i, c := range candidates {
		if c["rule"] == nil || c["rule"] == "" {
			t.Errorf("candidate[%d] carries no source rule identity: %v",
				i, keysOf(c))
		}
		if c["outcome"] == nil || c["outcome"] == "" {
			t.Errorf("candidate[%d] does not name the outcome it belongs to",
				i)
		}
		if _, present := c["required"]; !present {
			t.Errorf("candidate[%d] carries no `required` facts", i)
		}
		if _, present := c["unresolved"]; !present {
			t.Errorf("candidate[%d] carries no unresolved guard/gate facts", i)
		}

		if c["rule"] == "advance-draft" {
			sawAdvance = true
			// `advance-draft` carries `clear = ["stale"]`, normalized to a
			// `<clear>` write. The preview splits the clear keys out
			// (REQ-40, ASSUMPTION A-5), so the sentinel never leaks into
			// the write targets as if it were a value.
			clears, ok := stringsAt(c, "clear")
			if !ok || !slices.Contains(clears, "stale") {
				t.Errorf("candidate `advance-draft` does not preview `stale` "+
					"as a clear key; the normalized row names it. got %#v",
					c["clear"])
			}
			if writes, ok := c["writes"].(map[string]any); ok {
				if v, held := writes["stale"]; held {
					t.Errorf("the `<clear>` sentinel leaked into the write "+
						"targets as %#v; clear keys are split OUT of the "+
						"preview writes", v)
				}
			}
		}
	}
	if !sawAdvance {
		t.Error("no candidate for rule `advance-draft` was reported; it is " +
			"reachable from `status = draft`")
	}
}

// REQ-41: "It MUST run gate accessors only when --evaluate-gates is given;
// otherwise it MUST list gate ids as unresolved facts."
// REQ-46: "gates not requested (`next` default)" → "not run; gate ids
// listed as unresolved facts".
// REQ-57: "`flow next` stays effect-free by default".
// REQ-123 / `0005:S1`: "`flow next` over the fixture model in `--as=json`
// and `--as=text`, with and without `--evaluate-gates`" — "without the flag
// no gate accessor runs and gate ids appear as unresolved facts".
// HAPPY PATH — the default arm.
func TestReq41And46And57_WithoutTheFlagGateIdsAreUnresolvedFactsNotResults(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")

	data := flowData(t, requireSuccess(t, "flow", "next", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art), "--as=json"))

	candidates, ok := objectsAt(data, "candidates")
	if !ok {
		t.Fatalf("`candidates` is not an array of objects: %#v",
			data["candidates"])
	}

	var checked bool
	for _, c := range candidates {
		if c["rule"] != "advance-draft" {
			continue
		}
		checked = true
		// The gate id appears as an UNRESOLVED fact …
		unresolved, ok := stringsAt(c, "unresolved")
		if !ok {
			t.Fatalf("candidate `unresolved` is not an array of ids: %#v",
				c["unresolved"])
		}
		if !slices.Contains(unresolved, "approval") {
			t.Errorf("gate id `approval` is not listed among the unresolved "+
				"facts %v; without --evaluate-gates the gate does not run "+
				"and its id is what the caller sees", unresolved)
		}
		// … and NO gate result is reported, because no gate ran.
		if gates, present := c["gates"]; present {
			if arr, ok := gates.([]any); ok && len(arr) > 0 {
				t.Errorf("candidate carries evaluated gate results %#v "+
					"without --evaluate-gates; `next` is effect-free by "+
					"default", gates)
			}
		}
	}
	if !checked {
		t.Fatal("the gated candidate `advance-draft` was not reported")
	}
}

// REQ-42: "Under --evaluate-gates it MUST run the gates of every candidate
// row it reports and no others — a row whose guards already exclude it is
// not a candidate, so its gates MUST NOT run — and it MUST report each gate
// result on the candidate that carries it."
// REQ-46: "gate on a guard-excluded row" → "not run, not reported".
// `0005:MVV`: the gated-candidate / guard-excluded-row pair.
// REQ-130 / `0005:S8`: "under `--evaluate-gates` only the reported
// candidate's gate runs, its result rides that candidate, and a deny there
// still exits 0 (`next` enumerates, it does not select)".
// ADVERSARIAL
func TestReq42And46_OnlyReportedCandidatesGatesRunAndResultsRideTheCandidate(t *testing.T) {
	model := writeFlowModel(t, flowGatedNextModel)
	art := seedArtifact(t, model, "status=draft", "flag=true")

	data := flowData(t, requireSuccess(t, "flow", "next", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--evaluate-gates", "--as=json"))

	candidates, ok := objectsAt(data, "candidates")
	if !ok {
		t.Fatalf("`candidates` is not an array of objects: %#v",
			data["candidates"])
	}

	var reported, excluded bool
	for _, c := range candidates {
		switch c["rule"] {
		case "gated-reported":
			reported = true
			gates, ok := objectsAt(c, "gates")
			if !ok || len(gates) == 0 {
				t.Fatalf("the reported candidate carries no gate results "+
					"under --evaluate-gates: %#v", c["gates"])
			}
			var sawReportedGate bool
			for _, g := range gates {
				if g["id"] == "reported" {
					sawReportedGate = true
					if g["result"] == nil || g["result"] == "" {
						t.Error("the gate result on `gated-reported` carries " +
							"no `result` verdict")
					}
				}
				if g["id"] == "excluded" {
					t.Error("the EXCLUDED row's gate result rides the " +
						"reported candidate; each gate result MUST be " +
						"reported on the candidate that carries it")
				}
			}
			if !sawReportedGate {
				t.Error("the reported candidate's own gate `reported` has " +
					"no result")
			}
		case "gated-excluded":
			excluded = true
		}
	}
	if !reported {
		t.Error("candidate `gated-reported` was not reported; its guard " +
			"holds under flag=true")
	}
	if excluded {
		t.Error("candidate `gated-excluded` was reported; its guard excludes " +
			"it under flag=true, so it is not a candidate and its gate MUST " +
			"NOT run")
	}
}

// REQ-43: "A deny on one candidate constrains only that candidate: flow next
// reports it and still exits 0, because next enumerates rather than selects,
// and only flow resolve turns a deny into flow-gate-denied."
// REQ-46: "gate denies" → "exit 0; result on the candidate, no refusal".
// `0005:MVV`: "a deny there still exits 0".
// ADVERSARIAL — the single sharpest separation between `next` and `resolve`.
func TestReq43And46_ADenyUnderNextIsReportedAndStillExitsZero(t *testing.T) {
	model := writeFlowModel(t, flowGatedNextModel)
	// The artifact is seeded so the gate DENIES. How a fixture gate is made
	// to deny is the implementation's binding concern; the contract is that
	// whatever the verdict, `next` exits 0 and reports it.
	art := seedArtifact(t, model, "status=draft", "flag=true")

	stdout, _, err := runCmd(t, "flow", "next", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--evaluate-gates", "--as=json")
	if err != nil {
		t.Fatalf("`next --evaluate-gates` exited non-zero: %v\n"+
			"`next` ENUMERATES: a gate result — deny included — is reported "+
			"on the candidate and is never a refusal. Only `flow resolve` "+
			"turns a deny into `flow-gate-denied`.\nstdout:\n%s", err, stdout)
	}

	data := flowData(t, stdout)
	candidates, ok := objectsAt(data, "candidates")
	if !ok || len(candidates) == 0 {
		t.Fatalf("`candidates` is empty or malformed: %#v", data["candidates"])
	}
	// No refusal code leaked onto the success payload.
	for _, c := range candidates {
		gates, _ := objectsAt(c, "gates")
		for _, g := range gates {
			if code, ok := g["code"].(string); ok &&
				strings.HasPrefix(code, "flow-gate-") {
				t.Errorf("a gate result on a `next` candidate carries the "+
					"refusal code %q; under `next` a gate result is a "+
					"RESULT, never a refusal", code)
			}
		}
	}
}

// REQ-44: "It MUST NOT invent guard facts that were neither supplied, read,
// nor produced by a declared gate accessor."
// must not be materialized as a value.
// ADVERSARIAL — a fact absent from every channel must stay unresolved and
func TestReq44_NextInventsNoGuardFactAbsentFromEveryChannel(t *testing.T) {
	model := writeFlowModel(t, flowGatedNextModel)
	// `flag` is owned and NOT seeded, so no reader supplies it, no --tag
	// supplies it, and no gate produces it.
	art := seedArtifact(t, model, "status=draft")

	data := flowData(t, requireSuccess(t, "flow", "next", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art), "--as=json"))

	if owned, ok := data["owned"].(map[string]any); ok {
		if v, held := owned["flag"]; held {
			t.Errorf("`owned[flag]` = %#v was reported though no reader "+
				"supplied it; `next` MUST NOT invent a guard fact", v)
		}
	}

	// The rows guarded on `flag` must report it as UNRESOLVED rather than
	// being silently decided either way.
	candidates, _ := objectsAt(data, "candidates")
	for _, c := range candidates {
		unresolved, _ := stringsAt(c, "unresolved")
		if !slices.Contains(unresolved, "flag") {
			t.Errorf("candidate %v is guarded on `flag`, which nothing "+
				"supplied, but does not list it as an unresolved fact "+
				"(unresolved = %v)", c["rule"], unresolved)
		}
	}
}

// REQ-45: "It never calls a language model."
// REQ-111: "The command path should stay single-invocation deterministic".
// REQ-110: "The same request over the same model revision and the same
// artifact contents must produce the same success or refusal, excluding
// exit-3 environment failures."
// BOUNDARY — determinism is the observable proxy for "no model call".
func TestReq45And110And111_TheSameRequestProducesTheSameResultEveryRun(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	args := []string{"flow", "next", "--model", model, "--artifact", bind,
		"--tag", "profile=mid", "--as=json"}

	first := requireSuccess(t, args...)
	for i := range 4 {
		again := requireSuccess(t, args...)
		if again != first {
			t.Fatalf("run %d differs from run 0; the same request over the "+
				"same model revision and artifact contents must produce the "+
				"same result\nfirst:\n%s\nagain:\n%s", i+1, first, again)
		}
	}

	// The same holds for a refusal identity.
	rargs := []string{"flow", "resolve", "--model", model, "--artifact", bind,
		"--outcome", "not-an-outcome", "--as=json"}
	firstRefusal, _, _ := runCmd(t, rargs...)
	for i := range 4 {
		againRefusal, _, _ := runCmd(t, rargs...)
		if againRefusal != firstRefusal {
			t.Fatalf("refusal run %d differs from run 0; refusal identity is "+
				"deterministic too\nfirst:\n%s\nagain:\n%s",
				i+1, firstRefusal, againRefusal)
		}
	}
}

// REQ-35: "The invoked read-accessor set MUST be exactly those readers
// serving an owned key some candidate row of the requested model requires —
// not every declared reader."
// REQ-36: "A reader no candidate row needs MUST NOT run, and its unbound
// artifact role MUST NOT raise flow-artifact-missing; flow-artifact-missing
// is scoped to the roles the invoked set needs."
// REQ-38: "flow next and flow resolve MAY run declared read accessors over
// explicit --artifact role=path bindings and MUST assemble owned state only
// from them."
// `0005:MVV`: "That pair is the only assertion that distinguishes the
// narrowed invoked set from 'every declared reader'."
// DEV-1 reading (a): the union over ALL model rows.
// REQ-130 / `0005:S8`: "`next` and `resolve` neither invoke the unneeded
// reader nor raise `flow-artifact-missing`, while `read-state` on the same
// model invokes it and refuses the missing binding".
// ADVERSARIAL
func TestReq35And36And38_NextAndResolveSkipAReaderNoCandidateRowNeeds(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	// ONLY the `state` role is bound. `read.orphan` binds role `orphan`,
	// which stays UNBOUND — and `orphan` serves only `note`, which NO rule's
	// RequiresOwned names.
	bind := artifactBinding(flowStateRole, art)

	t.Run("next", func(t *testing.T) {
		stdout, _, err := runCmd(t, "flow", "next", "--model", model,
			"--artifact", bind, "--as=json")
		if err != nil {
			t.Fatalf("`next` refused with the orphan reader's role unbound: "+
				"%v\nA reader NO candidate row needs must not run, and its "+
				"unbound role must NOT raise flow-artifact-missing", err)
		}
		assertOrphanReaderNotInvoked(t, flowData(t, stdout))
	})

	t.Run("resolve", func(t *testing.T) {
		stdout, _, err := runCmd(t, "flow", "resolve", "--model", model,
			"--artifact", bind, "--outcome", "hold", "--as=json")
		if err != nil {
			t.Fatalf("`resolve` refused with the orphan reader's role "+
				"unbound: %v\nflow-artifact-missing is scoped to the roles "+
				"the INVOKED set needs", err)
		}
		assertOrphanReaderNotInvoked(t, flowData(t, stdout))
	})
}

// REQ-37: "flow read-state is the exception and runs every declared reader,
// because a diagnostic read has no candidate set to narrow by."
// `0005:MVV`: "`flow read-state` on the same model MUST invoke it and refuse
// the missing binding."
// DOMAIN EDGE — the other half of the discriminating pair.
func TestReq37_ReadStateRunsEveryDeclaredReaderAndRefusesTheMissingBinding(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")

	// The SAME single binding that `next` / `resolve` accept above must be
	// refused here, because `read-state` narrows by nothing.
	ce := requireRefusal(t, "flow-artifact-missing", 2,
		"flow", "read-state", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art), "--as=json")
	if ce.Param != flowOrphanRole {
		t.Errorf("param = %q; want the unbound ROLE %q — `read-state` runs "+
			"EVERY declared reader, so the orphan reader's role is needed",
			ce.Param, flowOrphanRole)
	}

	// With both roles bound it succeeds and reports both readers.
	data := flowData(t, requireSuccess(t, "flow", "read-state",
		"--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--artifact", artifactBinding(flowOrphanRole, art), "--as=json"))

	readers, ok := objectsAt(data, "readers")
	if !ok {
		t.Fatalf("`readers` is not an array of objects: %#v", data["readers"])
	}
	var sawState, sawOrphan bool
	for _, r := range readers {
		switch r["id"] {
		case "state":
			sawState = true
		case "orphan":
			sawOrphan = true
		}
	}
	if !sawState || !sawOrphan {
		t.Errorf("`read-state` reported readers %v; it runs EVERY declared "+
			"reader — both `state` and `orphan`", readers)
	}
}

// REQ-39: "They … MUST NOT run write accessors." (of `flow next` and
// `flow resolve`)
// REQ-59: "It MUST NOT invoke gate accessors or coerce gate allow, deny, or
// indeterminate results into tag values." (of `flow read-state`)
// owned state, read back through the CLI, is unchanged across the call.
// ADVERSARIAL — the observable proof that no write ran: the artifact's
func TestReq39And59_NextResolveAndReadStateNeverMutateTheArtifact(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	stateBind := artifactBinding(flowStateRole, art)
	orphanBind := artifactBinding(flowOrphanRole, art)

	before := requireSuccess(t, "flow", "read-state", "--model", model,
		"--artifact", stateBind, "--artifact", orphanBind, "--as=json")

	requireSuccess(t, "flow", "next", "--model", model,
		"--artifact", stateBind, "--as=json")
	requireSuccess(t, "flow", "resolve", "--model", model,
		"--artifact", stateBind, "--outcome", "hold", "--as=json")
	requireSuccess(t, "flow", "read-state", "--model", model,
		"--artifact", stateBind, "--artifact", orphanBind, "--as=json")

	after := requireSuccess(t, "flow", "read-state", "--model", model,
		"--artifact", stateBind, "--artifact", orphanBind, "--as=json")

	if before != after {
		t.Errorf("owned state changed across `next` / `resolve` / "+
			"`read-state`; none of them runs a write accessor\n"+
			"before:\n%s\nafter:\n%s", before, after)
	}
}

// REQ-59 / REQ-115: "`flow read-state` … does not invoke gate accessors,
// because gates return allow, deny, or indeterminate rather than tag
// values."
// ADVERSARIAL — a gate verdict must never appear as a tag value.
func TestReq59And115_ReadStateNeverCoercesAGateVerdictIntoATagValue(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")

	data := flowData(t, requireSuccess(t, "flow", "read-state",
		"--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--artifact", artifactBinding(flowOrphanRole, art), "--as=json"))

	// No reader entry names the model's gate, and no returned tag value is
	// a gate verdict.
	readers, _ := objectsAt(data, "readers")
	for _, r := range readers {
		if r["id"] == "approval" {
			t.Error("`read-state` reported the gate accessor `approval` as " +
				"a reader; it MUST NOT invoke gate accessors")
		}
		tags, ok := r["tags"].(map[string]any)
		if !ok {
			continue
		}
		for k, v := range tags {
			s, ok := v.(string)
			if !ok {
				continue
			}
			if slices.Contains([]string{"allow", "deny", "indeterminate"}, s) {
				t.Errorf("reader %v returned tag %q = %q, a GATE VERDICT "+
					"coerced into a tag value", r["id"], k, s)
			}
		}
	}
}

// REQ-121: "the kernel's unexported pre-selection guard filter also spells
// itself `gate`; it is a different mechanism from this RDR's post-selection
// accessor gates and must not be wired to them."
// REQ-113: "It does not own transition-model representation, guard
// semantics, accessor safety, graph lint, skill execution, or constrained
// decoding."
// rather than minting kinds of its own.
// BOUNDARY — the CLI reports the kernel's closed refusal-kind vocabulary
func TestReq113And121_TheCLIMapsTheKernelsClosedKindSetAndAddsNoKinds(t *testing.T) {
	// The mapping is `flow-<kind>` with underscores hyphenated (REQ-4),
	// over EXACTLY the kernel's closed five-kind set.
	want := map[string]resolve.RefusalKind{}
	for _, k := range resolve.RefusalKinds() {
		want["flow-"+strings.ReplaceAll(string(k), "_", "-")] = k
	}
	if len(want) != 5 {
		t.Fatalf("the kernel exposes %d refusal kinds; RDR 0005 maps a "+
			"closed set of five", len(want))
	}

	// Every code the CLI raises from a KERNEL refusal must be one of those
	// five — the CLI mints no kind of its own, and in particular does not
	// wire the kernel's unexported pre-selection guard filter (which also
	// spells itself `gate`) to this RDR's post-selection accessor gates.
	mvv := writeFlowModel(t, flowMVVModel)
	ambiguous := writeFlowModel(t, flowAmbiguousModel)
	guardUnev := writeFlowModel(t, flowGuardUnevaluableModel)
	mvvArt := seedArtifact(t, mvv, "status=draft")
	ambigArt := seedArtifact(t, ambiguous, "status=draft")
	guardArt := seedArtifact(t, guardUnev, "status=draft")

	kernelRefusals := []struct {
		kind resolve.RefusalKind
		args []string
	}{
		{resolve.KindUnmodeledOutcome, []string{"flow", "resolve",
			"--model", mvv,
			"--artifact", artifactBinding(flowStateRole, mvvArt),
			"--outcome", "not-an-outcome"}},
		{resolve.KindNoMatch, []string{"flow", "resolve", "--model", mvv,
			"--artifact", artifactBinding(flowStateRole, mvvArt),
			"--outcome", "bail"}},
		{resolve.KindAmbiguousMatch, []string{"flow", "resolve",
			"--model", ambiguous,
			"--artifact", artifactBinding(flowStateRole, ambigArt),
			"--outcome", "advance"}},
		{resolve.KindOwnedStateUnavailable, []string{"flow", "resolve",
			"--model", mvv,
			"--artifact", artifactBinding(flowStateRole, mvvArt),
			"--outcome", "advance"}},
		{resolve.KindGuardUnevaluable, []string{"flow", "resolve",
			"--model", guardUnev,
			"--artifact", artifactBinding(flowStateRole, guardArt),
			"--outcome", "advance"}},
	}

	for _, tc := range kernelRefusals {
		wantCode := "flow-" + strings.ReplaceAll(string(tc.kind), "_", "-")
		t.Run(string(tc.kind), func(t *testing.T) {
			requireRefusal(t, wantCode, 2, append(tc.args, "--as=json")...)
		})
	}
}

// assertOrphanReaderNotInvoked checks that `readers[]` — the accessor
// identities INVOKED — omits the reader no candidate row needs.
func assertOrphanReaderNotInvoked(t *testing.T, data map[string]any) {
	t.Helper()

	raw, ok := data["readers"]
	if !ok {
		t.Fatalf("payload carries no `readers[]`; it names the accessor "+
			"identities invoked. keys = %v", keysOf(data))
	}

	// `readers[]` may be rendered as ids or as objects carrying an id.
	if ids, ok := stringsAt(data, "readers"); ok {
		if slices.Contains(ids, "orphan") {
			t.Errorf("`readers[]` = %v names `orphan`; that reader serves "+
				"only `note`, which NO candidate row requires, so it MUST "+
				"NOT run", ids)
		}
		if !slices.Contains(ids, "state") {
			t.Errorf("`readers[]` = %v omits `state`; it serves the owned "+
				"keys the candidate rows require and MUST run", ids)
		}
		return
	}
	objs, ok := objectsAt(data, "readers")
	if !ok {
		t.Fatalf("`readers` is neither an id array nor an object array: %#v",
			raw)
	}
	var sawState bool
	for _, r := range objs {
		if r["id"] == "orphan" {
			t.Error("`readers[]` names `orphan`; that reader serves only " +
				"`note`, which NO candidate row requires, so it MUST NOT run")
		}
		if r["id"] == "state" {
			sawState = true
		}
	}
	if !sawState {
		t.Error("`readers[]` omits `state`; it serves the owned keys the " +
			"candidate rows require and MUST run")
	}
}
