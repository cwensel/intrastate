package graphlint_test

// RDR 0006 — finding shape, the code taxonomy and severity tiers
// (REQ-71..REQ-82), emission completeness and deterministic ordering
// (REQ-83..REQ-89), and the finding-identity tuple (REQ-114, REQ-127).

import (
	"slices"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
	"github.com/newcoinc/intrastate/internal/graphlint"
)

// REQ-71: "Every blocking finding MUST carry a stable code, model
// identity, severity, human-readable message, and the source rule/context
// id or source span when the normalized model can provide one."
// SC-5.
// HAPPY PATH
func TestReq71_EveryBlockingFindingCarriesTheRequiredIdentityFields(t *testing.T) {
	// A multi-defect model, so the assertion runs over several codes.
	r := lint(t, optionalGuardDecls, multiDefectBody)

	blocking := r.Blocking()
	if len(blocking) == 0 {
		t.Fatalf("the multi-defect fixture produced no blocking findings; "+
			"report:%s", render(r))
	}
	for _, f := range blocking {
		if f.Code == "" {
			t.Errorf("blocking finding carries no stable code: %+v", f)
		}
		if !slices.Contains(graphlint.BlockingCodes(), f.Code) {
			t.Errorf("blocking finding carries the code %q, outside the "+
				"blocking taxonomy %v", f.Code, graphlint.BlockingCodes())
		}
		if f.Model == "" {
			t.Errorf("%s carries no model identity: %+v", f.Code, f)
		}
		if f.Severity != graphlint.SeverityBlocking {
			t.Errorf("%s severity = %q; want %q", f.Code, f.Severity,
				graphlint.SeverityBlocking)
		}
		if strings.TrimSpace(f.Message) == "" {
			t.Errorf("%s carries no human-readable message: %+v", f.Code, f)
		}
	}
}

// REQ-13: "The finding's `model` field is the `[model].id` from the model
// itself, never the path."
// BOUNDARY
func TestReq13_FindingModelFieldIsTheModelIDNeverThePath(t *testing.T) {
	m := mustLoad(t, source(optionalGuardDecls, multiDefectBody))
	r := graphlint.Run(graphlint.NewRequest(m))

	if len(r.Findings) == 0 {
		t.Fatalf("the fixture produced no findings; report:%s", render(r))
	}
	for _, f := range r.Findings {
		if f.Model != m.ID {
			t.Errorf("%s model = %q; want the `[model].id` %q", f.Code, f.Model, m.ID)
		}
		// The loader was handed the source id "fixture.toml"; a finding
		// carrying it took the PATH instead of the model id.
		if strings.Contains(f.Model, ".toml") || strings.Contains(f.Model, "/") {
			t.Errorf("%s model = %q; the field is the model id, never the "+
				"path", f.Code, f.Model)
		}
	}
}

// REQ-127 / SC-5: "a finding whose normalized row lacks a rule id but has
// a source span or graph element id … remains actionable and stable by
// carrying the available source span or graph element id in the identity
// fields."
// INPUT EDGE
func TestReq127_FindingWithoutARuleIDCarriesASpanOrElementID(t *testing.T) {
	// The traversal-scoped and model-scoped findings are the arms with no
	// rule id to carry: the missing-root dangling finding names a
	// declaration, not a rule.
	const noRoot = `
terminal = ["done"]

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
	r := lint(t, statusOnlyDecls, noRoot)
	f := requireCode(t, r, graphlint.CodeDanglingEdge)

	for _, got := range f {
		if got.Rule != "" {
			continue
		}
		// No rule id, so a span or a graph element id must stand in.
		if got.Span == "" && got.Element == "" {
			t.Errorf("%s carries neither a rule id, a source span, nor a "+
				"graph element id, so it is neither actionable nor stable: "+
				"%+v", got.Code, got)
		}
	}
}

// REQ-72: "A finding attributed to one guard atom MUST carry that atom's
// `Key`, `Operator`, `Literal`, and `Block`; a finding scoped to an escape
// population MUST carry the failure class."
// REQ-82 / SC-3: "no arm of this code may carry the code alone."
// BOUNDARY
func TestReq72And82_AtomAttributedAndEscapeScopedFindingsCarryTheirFields(t *testing.T) {
	// The atom-attributed arm: a withheld claim names the refusing atom.
	const withheld = `
terminal = ["done"]

[initial]
status = "a"
opt = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "can-refuse"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
eq = "p"
[rule.write]
status = "b"
`
	r := lint(t, optionalGuardDecls, withheld)
	for _, f := range withCode(r, graphlint.CodeUnprovableCoverage) {
		if f.Reason != graphlint.ReasonRowCanRefuse {
			continue
		}
		if f.Key == "" || f.Operator == "" || f.Literal == "" || f.Block == "" {
			t.Errorf("the atom-attributed %s finding is missing atom fields: "+
				"key=%q operator=%q literal=%q block=%q",
				f.Code, f.Key, f.Operator, f.Literal, f.Block)
		}
	}

	// SC-3: the DIMENSION arm carries no atom, but must still name the
	// dimension — no arm may carry the code alone.
	const nonFinite = `
terminal = ["done"]

[initial]
status = "a"
free = "anything"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "reads-free"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.free]
eq = "x"
[rule.write]
status = "b"
`
	dim := lint(t, nonFiniteDecls, nonFinite)
	for _, f := range requireCode(t, dim, graphlint.CodeUnprovableCoverage) {
		if f.Dimension == "" && f.Key == "" && f.Rule == "" {
			t.Errorf("a %s finding carries the code alone — no dimension, no "+
				"key, no row: %+v", f.Code, f)
		}
	}

	// The escape-scoped arm: two escape rows sharing a class.
	const escPair = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "e1"
escape = ["no_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"

[[rule]]
id = "e2"
escape = ["no_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
`
	esc := lint(t, statusOnlyDecls, escPair)
	for _, f := range requireCode(t, esc, graphlint.CodeOverlap) {
		if f.Class == "" {
			t.Errorf("the escape-scoped %s finding carries no failure class: "+
				"%+v", f.Code, f)
		}
	}
}

// REQ-73: "The blocking code set is exactly: `graph-dangling-edge`,
// `graph-dead-end`, `graph-overlap`, `graph-coverage-gap`,
// `graph-unprovable-coverage`, `graph-single-valued-state`,
// `graph-always-present-owned`, `graph-owned-before-write`,
// `graph-terminal-escape`, `graph-product-too-large`."
// BOUNDARY
func TestReq73_TheBlockingCodeSetIsExactlyTheTenNamed(t *testing.T) {
	want := []string{
		"graph-always-present-owned",
		"graph-coverage-gap",
		"graph-dangling-edge",
		"graph-dead-end",
		"graph-overlap",
		"graph-owned-before-write",
		"graph-product-too-large",
		"graph-single-valued-state",
		"graph-terminal-escape",
		"graph-unprovable-coverage",
	}
	got := slices.Clone(graphlint.BlockingCodes())
	slices.Sort(got)

	if !slices.Equal(got, want) {
		t.Errorf("blocking code set =\n  %v\nwant exactly\n  %v", got, want)
	}
	for _, c := range want {
		if !graphlint.IsBlocking(c) {
			t.Errorf("%q is not reported blocking", c)
		}
	}
}

// REQ-74: "The advisory tier is closed at `graph-coverage-closed-by-escape`,
// `graph-redundant-row`, `graph-unreachable-rule`, and
// `graph-vacuous-atom`."
// BOUNDARY
func TestReq74_TheAdvisoryTierIsClosedAtExactlyFourMembers(t *testing.T) {
	want := []string{
		"graph-coverage-closed-by-escape",
		"graph-redundant-row",
		"graph-unreachable-rule",
		"graph-vacuous-atom",
	}
	got := slices.Clone(graphlint.AdvisoryCodes())
	slices.Sort(got)

	if !slices.Equal(got, want) {
		t.Errorf("advisory tier =\n  %v\nwant exactly\n  %v", got, want)
	}
	for _, c := range want {
		if graphlint.IsBlocking(c) {
			t.Errorf("%q is reported blocking; the advisory tier never "+
				"changes the success disposition", c)
		}
	}
	// The two tiers are disjoint.
	for _, c := range got {
		if slices.Contains(graphlint.BlockingCodes(), c) {
			t.Errorf("%q appears in both tiers", c)
		}
	}
}

// ASSUMPTION-3 / REQ-71: severity is two-valued, `blocking` / `info`.
// BOUNDARY
func TestReqSeverityVocabularyIsExactlyBlockingAndInfo(t *testing.T) {
	got := slices.Clone(graphlint.Severities())
	slices.Sort(got)
	want := []string{"blocking", "info"}

	if !slices.Equal(got, want) {
		t.Errorf("severities = %v; want exactly %v — no third tier exists",
			got, want)
	}
}

// REQ-75: "Advisory findings MUST NOT change the success disposition."
// HAPPY PATH
func TestReq75_AdvisoryFindingsDoNotChangeTheSuccessDisposition(t *testing.T) {
	// A model carrying every advisory code and no blocking defect.
	r := lint(t, twoStateDecls, advisoryBody)

	requireClean(t, r)
	if len(r.Advisory()) == 0 {
		t.Fatalf("the advisory fixture produced no advisory findings, so the "+
			"assertion is vacuous; report:%s", render(r))
	}
	for _, f := range r.Advisory() {
		if graphlint.IsBlocking(f.Code) {
			t.Errorf("Advisory() returned the blocking code %q", f.Code)
		}
		if f.Severity != graphlint.SeverityInfo {
			t.Errorf("%s severity = %q; want %q", f.Code, f.Severity,
				graphlint.SeverityInfo)
		}
	}
}

// REQ-77: "A redundant row is one whose accepted assignments are a proper
// subset of a sibling's in the same group" — code `graph-redundant-row`,
// `info`.
// SC-14.
// DOMAIN EDGE
func TestReq77_RedundantRowIsAProperSubsetOfASiblingsAssignments(t *testing.T) {
	// `narrow` accepts `flag = true` only; `wide` carries no guard and so
	// accepts every assignment. `narrow` is a proper subset of `wide`.
	const body = `
terminal = ["done"]

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "wide"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "narrow"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"
`
	r := lint(t, twoStateDecls, body)
	if !namesRule(r, graphlint.CodeRedundantRow, "narrow") {
		t.Errorf("no %s finding names the row whose accepted assignments are "+
			"a proper subset of a sibling's; report:%s",
			graphlint.CodeRedundantRow, render(r))
	}

	// Distinct from overlap: two rows neither of which subsumes the other
	// take `graph-overlap`, which is BLOCKING, not `graph-redundant-row`.
	if graphlint.IsBlocking(graphlint.CodeRedundantRow) {
		t.Errorf("%s is blocking; a redundant row is advisory and distinct "+
			"from overlap", graphlint.CodeRedundantRow)
	}
}

// REQ-78: "an unreachable rule is one no reachable owned-state node
// satisfies, including RDR 0002's dead-rule case" — a predicate set
// requiring `recognized` to be absent. Code `graph-unreachable-rule`,
// `info`.
// SC-14.
// DOMAIN EDGE
func TestReq78_UnreachableRuleCoversNoReachableNodeAndTheDeadRuleCase(t *testing.T) {
	// The first arm: a match pattern no reachable owned-state satisfies.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b", "orphan"]
single_valued = true
required = true
`
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "orphaned"
[rule.match.status]
eq = "orphan"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
	r := lint(t, decls, body)
	if !namesRule(r, graphlint.CodeUnreachableRule, "orphaned") {
		t.Errorf("no %s finding names the row no reachable owned-state node "+
			"satisfies; report:%s", graphlint.CodeUnreachableRule, render(r))
	}
	if graphlint.IsBlocking(graphlint.CodeUnreachableRule) {
		t.Errorf("%s is blocking; it is advisory", graphlint.CodeUnreachableRule)
	}

	// The reachable sibling is NOT reported, so the check discriminates.
	for _, f := range withCode(r, graphlint.CodeUnreachableRule) {
		if f.Rule == "advance" {
			t.Errorf("%s names the reachable row `advance`; report:%s",
				graphlint.CodeUnreachableRule, render(r))
		}
	}
}

// REQ-79: "a vacuous atom is an `exists` atom over an always-present key,
// which RDR 0003 requires be reported as well-formed-but-vacuous rather
// than rejected." Code `graph-vacuous-atom`, `info`.
// DOMAIN EDGE
func TestReq79_ExistsOverAnAlwaysPresentKeyIsAVacuousAtomNotARejection(t *testing.T) {
	// `always` carries the explicit always-present marker, so `exists` over
	// it is always true — well-formed but vacuous.
	const body = `
terminal = ["done"]

[initial]
status = "a"
always = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "vacuous"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.always]
exists = true
[rule.write]
status = "b"
`
	r := lint(t, alwaysPresentOwnedDecls, body)
	f := requireCode(t, r, graphlint.CodeVacuousAtom)

	for _, got := range f {
		if got.Severity != graphlint.SeverityInfo {
			t.Errorf("%s severity = %q; want %q — a report, not a rejection",
				got.Code, got.Severity, graphlint.SeverityInfo)
		}
		if got.Key != "always" {
			t.Errorf("%s does not name the always-present key `always`: %+v",
				got.Code, got)
		}
	}
	// Reported, never rejected.
	requireClean(t, r)
}

// REQ-80: "**`graph-unprovable-coverage` carries a `reason`
// discriminator.** … The finding MUST carry a stable `reason` from this
// closed, append-only set" — `dimension-not-finite`,
// `tag-not-single-valued`, `row-can-refuse`.
// ASSUMPTION-4: the discriminator is carried only on this code.
// BOUNDARY
func TestReq80_UnprovableCoverageCarriesAReasonFromTheClosedSet(t *testing.T) {
	want := []string{"dimension-not-finite", "row-can-refuse", "tag-not-single-valued"}
	got := slices.Clone(graphlint.Reasons())
	slices.Sort(got)

	if !slices.Equal(got, want) {
		t.Errorf("reason set = %v; want exactly %v", got, want)
	}

	// Every emitted finding of this code carries a reason from the set,
	// and no OTHER code carries one.
	r := lint(t, optionalGuardDecls, multiDefectBody)
	for _, f := range r.Findings {
		if f.Code == graphlint.CodeUnprovableCoverage {
			if !slices.Contains(want, f.Reason) {
				t.Errorf("%s carries reason %q, outside the closed set %v",
					f.Code, f.Reason, want)
			}
			continue
		}
		if f.Reason != "" {
			t.Errorf("%s carries the reason discriminator %q; the "+
				"discriminator is scoped to %s alone",
				f.Code, f.Reason, graphlint.CodeUnprovableCoverage)
		}
	}
}

// REQ-81: "`row-can-refuse` is *not* a model defect — it is an honest
// withholding — and the message MUST say so rather than reading as an
// authoring error."
// DOMAIN EDGE
func TestReq81_RowCanRefuseMessageReadsAsAWithholdingNotAnError(t *testing.T) {
	const body = `
terminal = ["done"]

[initial]
status = "a"
opt = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "can-refuse"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
eq = "p"
[rule.write]
status = "b"
`
	r := lint(t, optionalGuardDecls, body)

	var checked bool
	for _, f := range withCode(r, graphlint.CodeUnprovableCoverage) {
		if f.Reason != graphlint.ReasonRowCanRefuse {
			continue
		}
		checked = true
		msg := strings.ToLower(f.Message)
		// The message must say the claim is WITHHELD — the vocabulary of
		// declining to promise, not of an authoring mistake.
		if !strings.Contains(msg, "withh") && !strings.Contains(msg, "declin") {
			t.Errorf("the %q message %q does not say the claim is withheld; "+
				"it is an honest withholding, not a model defect",
				f.Reason, f.Message)
		}
		for _, blame := range []string{"invalid", "illegal", "malformed", "must fix", "error in"} {
			if strings.Contains(msg, blame) {
				t.Errorf("the %q message %q reads as an authoring error "+
					"(%q); the model is fine and lint declines to promise",
					f.Reason, f.Message, blame)
			}
		}
	}
	if !checked {
		t.Fatalf("no %s finding carries reason %q, so the message assertion "+
			"is vacuous; report:%s", graphlint.CodeUnprovableCoverage,
			graphlint.ReasonRowCanRefuse, render(r))
	}
}

// REQ-83: "Graph lint MUST report every defect it can decide in one pass
// over a row group, not the first it encounters"
// REQ-84: "Withholding a group's exhaustiveness claim MUST NOT suppress
// overlap, coverage, or further withholding findings for that group."
// REQ-85: "Each unprovable dimension, each refusing row, each overlapping
// pair, and any coverage gap over a provable product is its own finding"
// SC-13.
// ADVERSARIAL
func TestReq83And84And85_EveryDecidableDefectIsEmittedInOnePass(t *testing.T) {
	// SC-13's fixture: one group carrying an overlapping ordinary pair AND
	// two rows that can refuse. A first-failure engine emits only the
	// first and fails this.
	const body = `
terminal = ["done"]

[initial]
status = "a"
opt = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "over-one"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "over-two"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "refuse-one"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
eq = "p"
[rule.write]
status = "b"

[[rule]]
id = "refuse-two"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
eq = "q"
[rule.write]
status = "b"
`
	r := lint(t, optionalGuardDecls, body)

	// One overlap finding naming BOTH rows, not one per row.
	requireCode(t, r, graphlint.CodeOverlap)
	// Withholding did not suppress the overlap findings for the group.
	requireCode(t, r, graphlint.CodeUnprovableCoverage)

	// REQ-85: one withheld-claim finding PER refusing row.
	var refusing []string
	for _, f := range withCode(r, graphlint.CodeUnprovableCoverage) {
		if f.Reason == graphlint.ReasonRowCanRefuse && f.Rule != "" &&
			!slices.Contains(refusing, f.Rule) {
			refusing = append(refusing, f.Rule)
		}
	}
	for _, want := range []string{"refuse-one", "refuse-two"} {
		if !slices.Contains(refusing, want) {
			t.Errorf("no withheld-claim finding names the refusing row %q; "+
				"each refusing row is its own finding; got %v; report:%s",
				want, refusing, render(r))
		}
	}
}

// REQ-30 (the pair half) / SC-13: "One overlapping pair MUST yield one
// finding naming both rows"
// BOUNDARY
func TestReq30_OneOverlappingPairYieldsOneFindingNamingBothRows(t *testing.T) {
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "left"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "right"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
	r := lint(t, statusOnlyDecls, body)

	if n := countCode(r, graphlint.CodeOverlap); n != 1 {
		t.Fatalf("%d %s findings for one overlapping pair; want exactly 1 "+
			"naming both rows, not one per row; report:%s",
			n, graphlint.CodeOverlap, render(r))
	}
	f := requireOneCode(t, r, graphlint.CodeOverlap)
	// Both rows must be recoverable from the one finding.
	joined := f.Rule + " " + f.Message + " " + f.Element + " " + f.Fingerprint
	for _, want := range []string{"left", "right"} {
		if !strings.Contains(joined, want) {
			t.Errorf("the single %s finding does not name the row %q: %+v",
				f.Code, want, f)
		}
	}
}

// REQ-86: "Graph lint findings MUST be emitted in deterministic order by
// finding identity: model id, invariant code, source rule/context id or
// graph element id, then normalized predicate/write fingerprint."
// REQ-114: the identity tuple is `(model id, invariant code, source
// rule/context id or graph element id, normalized predicate/write
// fingerprint)`.
// BOUNDARY
func TestReq86And114_FindingsAreOrderedByTheFindingIdentityTuple(t *testing.T) {
	r := lint(t, optionalGuardDecls, multiDefectBody)
	if len(r.Findings) < 2 {
		t.Fatalf("the fixture produced %d findings; the ordering assertion "+
			"needs at least 2; report:%s", len(r.Findings), render(r))
	}

	for i := 1; i < len(r.Findings); i++ {
		prev, cur := r.Findings[i-1], r.Findings[i]
		if identityKey(prev) > identityKey(cur) {
			t.Errorf("findings are not ordered by the identity tuple:\n"+
				"  [%d] %s\n  [%d] %s",
				i-1, identityKey(prev), i, identityKey(cur))
		}
	}
}

// REQ-85 (order-independence) / REQ-89: the emitted set never depends on
// row or dimension iteration order.
// ADVERSARIAL
func TestReq85_EmittedSetIsIndependentOfAuthoredRowOrder(t *testing.T) {
	r := lint(t, optionalGuardDecls, multiDefectBody)
	rev := lint(t, optionalGuardDecls, multiDefectBodyReversed)

	// Two empty lists are equal, so the comparison below only means
	// something once the fixture has decided something.
	if len(r.Findings) == 0 {
		t.Fatalf("the multi-defect fixture produced no findings, so the "+
			"order-independence comparison is vacuous; report:%s", render(r))
	}

	got, want := identityList(r), identityList(rev)
	if !slices.Equal(got, want) {
		t.Errorf("reversing the authored row order changed the emitted "+
			"finding list:\n  authored: %v\n  reversed: %v", want, got)
	}
}

// REQ-87: "The fingerprint MUST be a canonical sortable serialization,
// never a hash: RDR 0002's canonical atom sort … over the row's predicate
// atoms and next-state tags"
// ADVERSARIAL
func TestReq87_FingerprintIsACanonicalSortableSerializationNeverAHash(t *testing.T) {
	m := mustLoad(t, source(twoStateDecls, legalBody))
	on := rowByRuleID(t, m, "advance-on")
	off := rowByRuleID(t, m, "advance-off")

	fpOn, fpOff := graphlint.Fingerprint(on), graphlint.Fingerprint(off)
	if fpOn == "" || fpOff == "" {
		t.Fatalf("Fingerprint returned an empty string (on=%q off=%q)", fpOn, fpOff)
	}
	if fpOn == fpOff {
		t.Errorf("two rows differing in a guard literal share the "+
			"fingerprint %q; it must be total over the predicate set", fpOn)
	}

	// Never a hash: the atom content is READABLE in the serialization, so
	// a reviewer can tell two fingerprints apart by eye and a sort over
	// them is meaningful.
	for _, want := range []string{"flag", "status"} {
		if !strings.Contains(fpOn, want) {
			t.Errorf("the fingerprint %q does not carry the atom key %q; a "+
				"canonical sortable serialization is not a hash", fpOn, want)
		}
	}
	// Canonical: reordering a row's atoms leaves the fingerprint unchanged
	// — the normalized atom set is already sorted, so two authorings of
	// the same row agree.
	reordered := mustLoad(t, source(twoStateDecls, legalBodyReordered))
	if got := graphlint.Fingerprint(rowByRuleID(t, reordered, "advance-on")); got != fpOn {
		t.Errorf("reordering the authored atoms changed the fingerprint:\n"+
			"  %q\n  %q", fpOn, got)
	}
}

// REQ-88: "When one run mixes identity namespaces, a source rule/context
// id sorts before any graph element id."
// BOUNDARY
func TestReq88_RuleIDsSortBeforeGraphElementIDs(t *testing.T) {
	// The multi-defect fixture mixes namespaces: rule-scoped findings and
	// the model/traversal-scoped ones that carry an element id instead.
	r := lint(t, optionalGuardDecls, multiDefectBody)

	var sawRule, sawElement bool
	for _, f := range r.Findings {
		if f.Rule != "" {
			sawRule = true
		}
		if f.Rule == "" && f.Element != "" {
			sawElement = true
		}
	}
	if !sawRule || !sawElement {
		t.Fatalf("the run does not mix identity namespaces (saw rule id=%v, "+
			"saw element id=%v), so the ordering clause is untested. The "+
			"multi-defect fixture carries both rule-scoped findings and the "+
			"model/traversal-scoped ones that stand in an element id; "+
			"report:%s", sawRule, sawElement, render(r))
	}

	// Within one (model, code) bucket, every rule-id finding precedes
	// every element-id one.
	for i := 1; i < len(r.Findings); i++ {
		prev, cur := r.Findings[i-1], r.Findings[i]
		if prev.Model != cur.Model || prev.Code != cur.Code {
			continue
		}
		if prev.Rule == "" && prev.Element != "" && cur.Rule != "" {
			t.Errorf("a graph element id sorts before a source rule id "+
				"within the %s bucket:\n  [%d] element=%q\n  [%d] rule=%q",
				cur.Code, i-1, prev.Element, i, cur.Rule)
		}
	}
}

// --- helpers -------------------------------------------------------------

// identityKey renders the finding-identity tuple as one sortable string:
// model id, invariant code, source rule/context id OR graph element id,
// then the normalized predicate/write fingerprint. A source rule/context
// id sorts before any graph element id (REQ-88), so the namespace is
// carried as a leading discriminator.
func identityKey(f clierr.Finding) string {
	id, namespace := f.Rule, "0"
	if id == "" {
		namespace = "1"
		id = f.Element
		if id == "" {
			id = f.Span
		}
	}
	return strings.Join([]string{
		f.Model, f.Code, namespace, id, f.Fingerprint,
		f.Class, f.Reason, f.Dimension, f.Key, f.Operator, f.Literal, f.Block,
	}, "\x00")
}

// identityList renders a report's findings as their identity keys.
func identityList(r graphlint.Report) []string {
	out := make([]string, 0, len(r.Findings))
	for _, f := range r.Findings {
		out = append(out, identityKey(f))
	}
	return out
}
