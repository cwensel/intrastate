package cli

// RDR 0023 — Phase 3b ADVERSARIAL oracles.
//
// Written against the RDR's own `Failure Modes` list, not against invented
// ones. Each test below attacks one of `0023:F1`..`0023:F7` on the axis the
// Phase 1 suite leaves open, and each says in its header whether it
// currently FAILS (a defect for the fixup phase) or PASSES (a regression
// guard closing an oracle the contract mandates but the build never wrote).
//
// The three attacked failure modes, and why these three:
//
//	F2  "a projected payload missing a plan-group field (mechanism bug)
//	     fails the ±-flag equality oracle and any consumer struct relying
//	     on `rule`/`emit`."
//	    Two independent surfaces reach it, and the shipped suite covers
//	    neither at the source. (ADV-1, ADV-2.)
//
//	F3  "a projection accidentally applied to a refusal path would change
//	     refusal bytes — A6 plus a refusal-unchanged oracle make it a test
//	     failure, not a field report."
//	    The shipped refusal differential runs three exit-2 refusals and
//	    stops short of the two that discriminate the SITE. (ADV-3.)
//
// The house rule the whole 0023 suite runs on holds here too: nothing below
// reads the projection's own code to learn what the projection does. Every
// expectation is an AUTHORED literal, the SAME REQUEST'S default-mode
// output, or the run's own observed side effects.

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
)

// --- ADV-1 -----------------------------------------------------------------
//
// F-2: "a projected payload missing a plan-group field (mechanism bug)
// fails the ±-flag equality oracle and any consumer struct relying on
// `rule`/`emit`." — (`0023:F2`)
//
// REQ-68: the C2 oracle asserts three things: "every struct field has
// exactly one entry; the entry set and the field set are equal (neither a
// field without an entry nor an entry without a field); and the keys a
// projected run actually emits equal the declaration's `plan` side."
// REQ-69: "A new field with no entry then fails at the first assertion
// rather than defaulting into either width."
// REQ-66: "the assignment is DECLARED (a standalone table keyed by field
// name, one entry per field, carrying `echo` or `plan`)"
//
// ADVERSARIAL — the DECLARATION ITSELF is the subject, and that is the gap.
//
// The shipped completeness oracle
// (`TestReq60And61And62And63And68And69_EveryPayloadFieldIsAssignedToExactlyOneC2Group`)
// asserts REQ-68's first two clauses against a map AUTHORED IN THE TEST
// FILE, and its third clause end-to-end against emitted keys. That is
// exactly right for C2's independence requirement — an oracle that read the
// production table to learn the assignment would be the tautology C2
// forbids. But it leaves the production carrier `resolvePayloadGroups`
// UNVERIFIED as a carrier: nothing asserts the SHIPPED table is itself
// complete over the struct.
//
// The asymmetry is what makes this reachable. `projectAwayEchoGroup` drives
// off `echoFieldNames()`, which is a map lookup by Go field name and
// silently yields the zero `payloadGroup` for any name the table does not
// carry. So:
//
//   - an ECHO field the table loses is silently NOT projected, and the
//     end-to-end assertion catches it (an echo key survives the flag);
//   - a PLAN field the table loses changes NOTHING observable, because the
//     projection never consults plan entries — the field defaults into the
//     projected width by falling out of the map, which is precisely
//     "defaulting into either width" REQ-69 forbids and precisely the
//     "bare omit-list whose complement is whatever else exists" REQ-60
//     forbids.
//
// This test asserts the shipped table against the shipped struct in BOTH
// directions, so a stale key or an unassigned field fails at the FIRST
// assertion (REQ-69) rather than surviving on the strength of the plan
// side's being unobservable.
//
// It is NOT a tautological read of the projection: the subject here is the
// DECLARATION's shape (total, bijective over the field set), never the
// membership it declares. Which side each field is on stays the authored
// census's business, asserted where the shipped suite already asserts it.
//
// STATUS: PASSES against the current implementation — today's table is
// complete — and is the SOLE discriminator of the gap. Verified by
// mutation: renaming one PLAN key in `resolvePayloadGroups` (leaving the
// struct and its JSON tags untouched) leaves the ENTIRE shipped suite
// green, including the C2 completeness oracle, and turns only this test
// red. The defect the fixup phase owes is the missing assertion, not a
// wrong table.
func TestAdversarialF2_TheShippedDeclarationTableIsTotalOverThePayloadStruct(t *testing.T) {
	rt := resolvePayloadType()

	declared := map[string]payloadGroup{}
	for name, side := range resolvePayloadGroups {
		declared[name] = side
	}

	// (1) every struct field has exactly one entry.
	fieldNames := make([]string, 0, rt.NumField())
	for i := range rt.NumField() {
		f := rt.Field(i)
		fieldNames = append(fieldNames, f.Name)

		side, held := declared[f.Name]
		if !held {
			t.Errorf("`resolvePayload.%s` (wire key %q) has NO entry in the "+
				"shipped declaration `resolvePayloadGroups`.\n"+
				"REQ-68 requires the declaration carry one entry per field and "+
				"REQ-69 requires an unassigned field to fail at the FIRST "+
				"assertion rather than defaulting into either width.\n"+
				"An unassigned PLAN field is INVISIBLE to every other oracle in "+
				"this suite: `echoFieldNames()` is a map lookup that yields the "+
				"zero group for an absent name, the projection only ever "+
				"consults the ECHO side, and the end-to-end key-set assertions "+
				"compare EMITTED keys — which a plan field reaches whether or "+
				"not it was ever declared. That is the \"bare omit-list whose "+
				"complement is whatever else exists\" `0023:C2` forbids by "+
				"name, reached through the carrier rather than through the "+
				"projection.",
				f.Name, wireKeyForField(f))
			continue
		}
		if side != groupEcho && side != groupPlan {
			t.Errorf("`resolvePayload.%s` is declared %q; the doctrine admits "+
				"exactly two groups (`echo`, `plan`) and no third value",
				f.Name, side)
		}
	}

	// (2) the entry set and the field set are EQUAL — neither a field
	// without an entry (above) nor an entry without a field (here). A stale
	// key is the drift arm a Go-identifier rename leaves behind: the JSON
	// tag is unchanged, so every WIRE-keyed oracle in this suite still sees
	// the same key set, while `echoFieldNames()` quietly drops a name.
	for name := range declared {
		if !slices.Contains(fieldNames, name) {
			t.Errorf("the shipped declaration `resolvePayloadGroups` carries an "+
				"entry for %q, which `resolvePayload` does not declare as a "+
				"field.\nREQ-68 fixes the equality in BOTH directions. A stale "+
				"entry is what a Go-identifier rename leaves behind — the JSON "+
				"tag survives, so every wire-keyed oracle stays green while "+
				"`echoFieldNames()` silently returns one fewer name and the "+
				"renamed field falls into whichever width its absence "+
				"implies.\nstruct fields = %v",
				name, fieldNames)
		}
	}

	// Vacuity guard. Both loops above are satisfied by an empty struct and
	// an empty table, so the pair proves nothing without a floor. 14 is the
	// count `decision_table_0010_test.go` pins (REQ-64) and the count
	// `0023:A2`'s conversion is required to preserve (REQ-90).
	if len(fieldNames) == 0 || len(declared) == 0 {
		t.Fatalf("vacuous: %d struct fields against %d declared entries; the "+
			"bidirectional equality above is satisfied by two empty sets",
			len(fieldNames), len(declared))
	}
}

// --- ADV-2 -----------------------------------------------------------------
//
// F-2: "a projected payload missing a plan-group field (mechanism bug)
// fails the ±-flag equality oracle" — (`0023:F2`), reached through the
// clause `0023:F2` rests on:
//
// REQ-41: "The oracle MUST instead observe which readers ACTUALLY RAN in
// each run — recording execution at the reader invocation site (a counting
// or recording seam around the reader pass, the run's own side effects,
// never a re-call of the planning function) — and assert the two observed
// sets are equal."
// REQ-43: "The DEFAULT-mode run of the same request supplies the expected
// set. The discriminating property is that a projected run which skips the
// reader pass MUST turn this assertion red; an assertion that cannot
// distinguish that case does not satisfy this clause."
// REQ-37: "A later change that skips work whose only consumer is a
// projected-away field breaches this clause."
// REQ-78 (S1's discriminating control): "SKIP the reader pass under the
// flag and confirm the reader assertion goes red."
//
// ADVERSARIAL — the mandated oracle was never written, so the seam built to
// carry it is dead code.
//
// `internal/cli/flow_exec.go::readerExecutionHook` exists, sits at the
// per-reader `exec.Read` call C1 fixes as the position, and is flag-blind —
// all correct. But no shipped test installs it. The Phase 1 reader oracle
// (`TestReq35And37And40And41And42And43And44And46_...`) asserts a PROXY
// instead: ±-flag equality of `rule`/`next`/`writes`/`clear`, reasoning
// that a skipped reader pass would resolve over an empty owned view and
// select differently.
//
// The proxy is not the measurand REQ-41 fixes, and it is weaker in a way
// that matters: it only discriminates a skipped reader pass when the owned
// state that pass establishes CHANGES THE SELECTION. An implementation that
// skipped a reader under the flag on a request whose rule does not turn on
// that reader's keys would leave every plan-group field identical and the
// proxy green — while `readers` and `owned`, the two fields the pass yields
// together, are exactly the ones the flag projects away and so cannot be
// read back. That is REQ-37's breach ("skips work whose only consumer is a
// projected-away field") with nothing left to see it.
//
// This test asserts the measurand itself: the OBSERVED executed set from
// each run, collected at the seam, compared across the two widths.
//
// STATUS: PASSES against the current implementation, and is recorded as a
// REGRESSION GUARD rather than as a caught defect. The seam is correctly
// positioned at the `exec.Read` call and correctly flag-blind, so the
// implementation is right on this axis; what was missing is the ORACLE C1
// mandates, which left `readerExecutionHook` as dead code in the shipped
// build. Verified to discriminate by mutation: a flag-conditional skip of
// the reader pass turns the `reader-backed/owned-decides` arm red on the
// observed set. The `decision-table/owned-does-not-decide` arm carries no
// reader by construction and is the CONTROL that fixes the other arm's
// non-vacuity — it is not itself a discriminator, and this header does not
// claim it is.
func TestAdversarialF2_ObservedReaderExecutionIsEqualPlusOrMinusTheFlag(t *testing.T) {
	// Two shapes, chosen so the pair covers what the proxy cannot.
	for _, shape := range []struct {
		name string
		args func(t *testing.T) []string
		// ownedDecidesSelection records whether the proxy oracle could
		// stand in for this shape. On the decision-table shape it cannot,
		// which is the whole reason that arm is here.
		ownedDecidesSelection bool
		wantReaders           bool
	}{
		{
			// A reader-backed model whose owned state DECIDES the rule.
			// The proxy covers this one; it is here as the control that
			// the seam observes a NON-EMPTY set at all.
			name: "reader-backed/owned-decides",
			args: func(t *testing.T) []string {
				model := writeFlowModel(t, flowMVVModel)
				art := seedArtifact(t, model, "status=draft", "stale=x")
				return append(resolveArgs(model,
					artifactBinding(flowStateRole, art), "advance"), "--as=json")
			},
			ownedDecidesSelection: true,
			wantReaders:           true,
		},
		{
			// A DECISION TABLE: zero accessors, selection on observed tags
			// alone. Here the proxy is blind — every plan-group field is
			// identical whether or not a reader pass ran, because there is
			// no owned state to change the selection. The observed set is
			// the only measurand left, and it must still agree.
			name: "decision-table/owned-does-not-decide",
			args: func(t *testing.T) []string {
				return append(pricingCall(t), "--as=json")
			},
			ownedDecidesSelection: false,
			wantReaders:           false,
		},
	} {
		t.Run(shape.name, func(t *testing.T) {
			args := shape.args(t)

			// The seam is driven identically by both widths; a seam that
			// could see the flag would make the very execution it records
			// flag-dependent, and C1 fixes exactly one lexical site for the
			// flag, which is not this one.
			observe := func(extra ...string) []string {
				restore := readerExecutionHook
				var ran []string
				readerExecutionHook = func(readerID string) {
					ran = append(ran, readerID)
				}
				defer func() { readerExecutionHook = restore }()

				requireSuccess(t, append(slices.Clone(args), extra...)...)
				return ran
			}

			defaultRan := observe()
			projectedRan := observe("--" + planOnlyFlag)

			// Control: the seam is live. Without this the equality below is
			// satisfied by two runs that recorded nothing, which is the
			// vacuous form C1 rejects by name.
			if shape.wantReaders && len(defaultRan) == 0 {
				t.Fatalf("the DEFAULT run recorded no reader execution at the " +
					"seam. C1 fixes the measurand as OBSERVED READER " +
					"EXECUTION and the default run as the source of the " +
					"expected set; an empty expected set makes the " +
					"differential unfailable")
			}
			if !shape.wantReaders && len(defaultRan) != 0 {
				t.Fatalf("the decision-table shape recorded reader execution "+
					"%v; this arm exists because it runs NO reader, which is "+
					"what makes the plan-group proxy blind here. Without that "+
					"property the arm proves nothing", defaultRan)
			}

			// The assertion. Observed execution, not a recomputation:
			// `invokedReaders(model, outcome)` is a pure function of two
			// arguments the flag does not touch, so comparing it returns
			// equal sets on every implementation including one that never
			// runs a reader at all.
			if !slices.Equal(defaultRan, projectedRan) {
				t.Errorf("the OBSERVED executed reader set differs ± the "+
					"flag:\n  default   = %v\n  projected = %v\n"+
					"--%s is REPORT-ONLY: it MUST NOT change what is decided, "+
					"and a projected run that skips work whose only consumer "+
					"is a projected-away field (`readers`, `owned`) breaches "+
					"C1's report-only clause. This is the measurand C1 fixes "+
					"— the run's own side effects at the `exec.Read` site — "+
					"not the plan-group proxy, which on this shape cannot "+
					"distinguish the two runs at all (ownedDecidesSelection "+
					"= %v)",
					defaultRan, projectedRan, planOnlyFlag,
					shape.ownedDecidesSelection)
			}
		})
	}
}

// --- ADV-3 -----------------------------------------------------------------
//
// F-3: "a projection accidentally applied to a refusal path would change
// refusal bytes — A6 plus a refusal-unchanged oracle make it a test
// failure, not a field report." — (`0023:F3`)
//
// REQ-33: the differential asserts "byte-identical refusal envelopes
// (CLIError, findings — never projected, and carrying no echo group to
// project)"
// REQ-32: the same differential asserts "identical exit codes"
// REQ-48: the projection "sits on the SUCCESS path only, AFTER the last
// respond.Fail return"
// REQ-50: "A defensive \"if refusing, skip projection\" guard is FORBIDDEN:
// refusal flag-blindness (A6) must hold structurally, because a refusal
// returns before the projection is reachable"
//
// ADVERSARIAL — the shipped refusal differential stops short of the two
// refusals that actually discriminate the SITE.
//
// `TestReq33And100And131_RefusalEnvelopesAreByteIdenticalPlusOrMinusTheFlag`
// runs three refusals: `flow-tag-invalid` (before the readers),
// `flow-unmodeled-outcome` (after the readers, before the gates) and
// `flow-ambiguous-match` (same band). All three are exit-2, and all three
// return from `runFlowResolve` well ABOVE the gate site. Two classes are
// uncovered and both are the ones C1's site rule is written about:
//
//  1. `flow-gate-denied` — the LAST `respond.Fail` in the function, the one
//     immediately preceding the projection branch. REQ-48 places the
//     projection "AFTER the last respond.Fail return", so this is the exact
//     boundary the clause draws, and it is the single refusal a site that
//     drifted one statement upward would newly reach. The three shipped arms
//     would all still pass such a drift, because a projection applied above
//     them would still be above their own returns.
//
//  2. The exit-3 GROUP. REQ-32 asserts "identical exit codes" and the
//     shipped arms exercise exit 2 only, so the assertion is pinned on one
//     group. An environment failure (`flow-accessor-timeout`,
//     `flow-read-incomplete`) refuses from inside the reader pass — the one
//     pass a report-width flag has any reason to want to skip (REQ-37) —
//     and it is the class where a flag-conditional shortcut would surface
//     as a CHANGED EXIT rather than as changed bytes.
//
// STATUS: PASSES against the current implementation. The site is correct:
// every refusal returns before the projection is reachable and no defensive
// guard exists. Recorded as a REGRESSION GUARD over the two uncovered
// refusal classes. Verified to discriminate by mutation: a flag-conditional
// skip of the reader pass turns all three arms red — including the exit-3
// arm, on the exit-code assertion the shipped exit-2-only differential
// cannot make.
func TestAdversarialF3_RefusalsAtTheSiteBoundaryAndInTheExit3GroupAreFlagBlind(t *testing.T) {
	denyModel := writeFlowModel(t, flowGateDenyModel)
	denyArt := seedArtifact(t, denyModel, "status=draft")

	failModel := writeFlowModel(t, flowGateFailModel)
	failArt := seedArtifact(t, failModel, "status=draft")

	for _, refusal := range []struct {
		name string
		args []string
		code string
		exit int
		why  string
	}{
		{
			// The LAST respond.Fail before the projection branch. REQ-48
			// draws the boundary here and nowhere else.
			name: "gate-denied/the-last-fail-before-the-site",
			args: resolveArgs(denyModel,
				artifactBinding(flowStateRole, denyArt), "deny-path"),
			code: "flow-gate-denied",
			exit: 2,
			why: "this is the last `respond.Fail` return in `runFlowResolve` " +
				"— the statement REQ-48's \"AFTER the last respond.Fail\" is " +
				"measured against. A projection site that drifted one " +
				"statement upward would reach THIS refusal first, and every " +
				"shipped refusal arm sits far enough above it to stay green " +
				"through that drift",
		},
		{
			// Same site band, the other verdict, so a site that discriminated
			// on the gate VERDICT rather than on the return would show here.
			name: "gate-indeterminate/the-same-band",
			args: resolveArgs(denyModel,
				artifactBinding(flowStateRole, denyArt), "indeterminate-path"),
			code: "flow-gate-indeterminate",
			exit: 2,
			why: "the gate band's other verdict; a projection reachable on a " +
				"gate refusal would not be verdict-selective, so both arms " +
				"must be flag-blind together",
		},
		{
			// The exit-3 group, refusing from INSIDE the reader pass — the
			// one pass REQ-37 names as the reachable surface for a
			// report-width shortcut.
			name: "accessor-unavailable/the-exit-3-group",
			args: resolveArgs(failModel,
				artifactBinding(flowStateRole, failArt), "advance"),
			code: "",
			exit: 3,
			why: "the shipped refusal differential runs exit-2 refusals only, " +
				"so REQ-32's \"identical exit codes\" is pinned on one group. " +
				"This refusal comes from inside the reader pass, which is the " +
				"work REQ-37 names as the reachable surface for a " +
				"report-width shortcut — a flag-conditional skip there " +
				"surfaces as a CHANGED EXIT, not as changed bytes",
		},
	} {
		t.Run(refusal.name, func(t *testing.T) {
			defOut, _, defErr := runCmd(t,
				append(slices.Clone(refusal.args), "--as=json")...)
			projOut, _, projErr := runCmd(t, append(slices.Clone(refusal.args),
				"--"+planOnlyFlag, "--as=json")...)

			// Control: this arm really is the refusal class it claims.
			// Without it the byte comparison passes on two identical
			// successes, or on two identical parse errors.
			if defErr == nil {
				t.Fatalf("the default run SUCCEEDED; this arm exists to "+
					"exercise a refusal (%s)\n%s", refusal.why, defOut)
			}
			if got := clierr.ExitCodeFor(defErr); got != refusal.exit {
				t.Fatalf("the default run exits %d; this arm exists to "+
					"exercise exit %d — %s", got, refusal.exit, refusal.why)
			}
			if refusal.code != "" {
				if got := clierr.ErrorCode(defErr); got != refusal.code {
					t.Fatalf("the default refusal code = %q; want %q",
						got, refusal.code)
				}
			}

			// REQ-32: identical exit codes, ± the flag.
			if got, want := clierr.ExitCodeFor(projErr),
				clierr.ExitCodeFor(defErr); got != want {
				t.Errorf("the run exits %d under the flag and %d by default. "+
					"--%s MUST NOT change HOW THE RUN FAILS: the projection "+
					"sits on the SUCCESS path only, after the last "+
					"`respond.Fail` return, and no defensive \"if refusing, "+
					"skip projection\" guard is permitted — refusal "+
					"flag-blindness holds STRUCTURALLY or not at all.\n%s",
					got, want, planOnlyFlag, refusal.why)
			}

			// REQ-33: byte-identical refusal envelopes, on the full emitted
			// line (`0023:A-7`) — the same unit the success-width clause
			// fixes.
			if got, want := strings.TrimRight(projOut, "\n"),
				strings.TrimRight(defOut, "\n"); got != want {
				t.Errorf("the refusal envelope differs ± the flag:\n"+
					"  default   = %s\n  projected = %s\n"+
					"`0023:A6`: the refusal envelope's members are "+
					"code/message/param/detail/hint/findings — there is NO "+
					"echo-group member to project, so --%s cannot change a "+
					"single refusal byte.\n%s",
					want, got, planOnlyFlag, refusal.why)
			}

			// The refusal is not merely equal but UNPROJECTED: no echo-group
			// key may have been deleted from it, and the envelope must not
			// have acquired one either. Asserted on both sides, because
			// "equal" is also satisfied by two identically-projected
			// refusals.
			for _, key := range echoGroup0023 {
				needle := `"` + key + `":`
				if strings.Contains(defOut, needle) {
					t.Errorf("the DEFAULT refusal envelope carries the "+
						"echo-group key %q; `0023:A6` grounds the whole "+
						"report-only claim on refusals having no request-echo "+
						"member to project", key)
				}
			}
		})
	}
}

// --- ADV-3 companion: the site's structural position -----------------------
//
// F-3 again, from the source side. REQ-48/REQ-50 make the site's position a
// STRUCTURAL property, and the behavioural arms above corroborate it without
// asserting it: they show that today's refusals are flag-blind, not that a
// refusal COULD NOT be reached below the projection.
//
// STATUS: PASSES. Regression guard on the one property the behavioural arms
// cannot express — that no `respond.Fail` return follows the projection
// call in `runFlowResolve`, which is what makes refusal flag-blindness
// structural rather than a fact about today's fixtures.
func TestAdversarialF3_NoRefusalReturnFollowsTheProjectionSite(t *testing.T) {
	src, err := readCLISource(t, "flow_resolve.go")
	if err != nil {
		t.Fatalf("read flow_resolve.go: %v", err)
	}

	body := functionBody0023(src, "func runFlowResolve(")
	if body == "" {
		t.Fatal("`runFlowResolve` not found in flow_resolve.go; the site " +
			"rule is stated about that function")
	}

	siteAt := strings.Index(body, "projectResolvePayload(")
	if siteAt < 0 {
		t.Fatal("`runFlowResolve` does not call the projection; C1 fixes the " +
			"site as \"applied to the verb-specific result BEFORE respond.OK\"")
	}

	// REQ-48: the site sits AFTER the last `respond.Fail` return.
	if lastFail := strings.LastIndex(body, "respond.Fail("); lastFail > siteAt {
		t.Errorf("a `respond.Fail` call follows the projection site in " +
			"`runFlowResolve`. The projection \"sits on the SUCCESS path " +
			"only, AFTER the last respond.Fail return\" — a refusal reachable " +
			"below the site would make refusal flag-blindness a fact about " +
			"the fixtures rather than a structural property, and `0023:A6` " +
			"rests on the structural reading")
	}

	// REQ-50: no defensive guard. A guard would both signal a misplaced site
	// and make the flag readable on a path C1 requires it cannot influence.
	guardBody := body[:siteAt]
	for _, guard := range []string{
		"if refus", "if ce != nil && planOnly", "skipProjection",
	} {
		if strings.Contains(guardBody, guard) {
			t.Errorf("`runFlowResolve` carries a %q guard around the "+
				"projection; a defensive \"if refusing, skip projection\" "+
				"guard is FORBIDDEN — refusal flag-blindness must hold "+
				"STRUCTURALLY, because a refusal returns before the "+
				"projection is reachable", guard)
		}
	}

	// REQ-49: the flag is read at exactly ONE lexical site, and it is not
	// this function. `runFlowResolve` hands the assembled payload to the
	// projection; it never reads the flag itself.
	if strings.Contains(body, planOnlyFlag) &&
		!strings.Contains(body, "projectResolvePayload(") {
		t.Errorf("`runFlowResolve` names %q outside the projection call; the "+
			"flag is read at exactly ONE lexical site — the projection branch "+
			"— and is never a parameter to payload assembly", planOnlyFlag)
	}
}

// readCLISource reads one non-test file from `internal/cli`.
//
// Reached by PATH rather than by symbol so a tree that has moved the file
// fails on the read with a readable message, instead of failing to compile
// the package and taking every shipped suite red with it.
func readCLISource(t *testing.T, name string) (string, error) {
	t.Helper()

	b, err := os.ReadFile(repoRootFor(t) + "/internal/cli/" + name)
	return string(b), err
}

// functionBody0023 returns the source text from decl to the first
// column-zero closing brace after it — one Go function body, without
// pulling in a parser the rest of this suite does not use.
func functionBody0023(src, decl string) string {
	at := strings.Index(src, decl)
	if at < 0 {
		return ""
	}
	rest := src[at:]
	if end := strings.Index(rest, "\n}\n"); end > 0 {
		return rest[:end]
	}
	return rest
}
