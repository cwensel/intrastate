package cli

// RDR 0011 — the `unknown` payload: its four sources, its closed reason
// vocabulary, its dedup rule, its sort, and its one stated fidelity limit.
//
// The rename `unresolved` → `unknown` is a payload change INDEPENDENT of
// the candidate predicate, and C3 says so outright: it "breaks assertions
// the predicate leaves untouched". The oracles here are about the payload,
// not about which rows survive.
//
// Two shapes are load-bearing throughout:
//
//   - an entry is a `{key, reason}` PAIR, never a bare key string, and the
//     reason names the caller's REMEDY — bind a reader or supply the tag
//     (`absent`), fix the value or the atom (`uncomparable`), pass
//     `--evaluate-gates` (`not-evaluated`).
//   - dedup is on the PAIR and runs LAST, after C2's `--all` filter has
//     been applied to the atom walk at emission.

import (
	"slices"
	"sort"
	"strings"
	"testing"
)

// REQ-22: "Every match atom over an absent key, every guard atom the kernel
// could not decide, every owned key no invoked reader established, and
// (absent --evaluate-gates) every gate id MUST appear in that candidate's
// `unknown` list as a `{key, reason}` pair."
// REQ-71: "the per-candidate field is `unknown`, a list of `{key, reason}`"
// A-2: the wire members are exactly `key` and `reason`, both always
// emitted (no `omitempty`).
// HAPPY PATH — the four sources on one candidate, which is also the shape
// the sort oracle below needs.
func TestReq22And71_AllFourSourcesLandOnUnknownAsKeyReasonPairs(t *testing.T) {
	model := writeFlowModel(t, flowSortOrderModel)
	// `mid` is an owned key the reader DECLARES but nothing establishes,
	// so the owned-key walk contributes it; `zebra` and `alpha` are
	// observed match keys nothing supplies, so the atom walk contributes
	// them; the gate `beta` is un-run because the flag is not passed.
	bind := seedMatchArtifact(t, model, "status=draft")

	data := runNext(t, model, []string{bind})
	pairs := candidateUnknown(t, requireCandidate(t, data, "sort-row"))

	for _, want := range []unknownPair{
		{Key: "alpha", Reason: "absent"},       // atom walk (match, absent key)
		{Key: "zebra", Reason: "absent"},       // atom walk (match, absent key)
		{Key: "mid", Reason: "absent"},         // owned-key walk (unestablished)
		{Key: "beta", Reason: "not-evaluated"}, // gate-id loop (un-run gate)
	} {
		if !hasUnknown(pairs, want.Key, want.Reason) {
			t.Errorf("`unknown` does not carry {%s, %s}: %#v\nAll four "+
				"sources C1 names contribute to ONE list, and every entry "+
				"is a PAIR — a bare key cannot carry the remedy, which is "+
				"the whole reason the flat `[]string` was rejected",
				want.Key, want.Reason, pairs)
		}
	}
}

// REQ-26: "The gate ids are the ROW's own declared gates —
// `internal/table/model.go::Row.Gate`, the same field
// `internal/cli/flow_next.go` already passes to `runGates` when gates do run
// — not a model-level set and not `0005:C1`'s gate-accessor list"
// REQ-27: "A gate id un-run because --evaluate-gates was not passed is
// neither `absent` … nor `uncomparable` …, so it carries `not-evaluated`,
// which this RDR mints for the gate class ONLY".
// REQ-57: "The gate-id reads re-home to the `not-evaluated` reason C1 names;
// a re-homing that has to invent a reason token is not mechanical and is a
// defect."
// REQ-81: "gate id, `--evaluate-gates` not passed | yes | `{gate-id,
// not-evaluated}` | no | 0"
// A-3 / A-5: the token is the literal lowercase-HYPHEN `not-evaluated`, and
// the gate id itself goes in the `key` member.
// DOMAIN EDGE
func TestReq26And27And57And81_AnUnrunGateIdCarriesNotEvaluatedAndTheIdIsTheKey(t *testing.T) {
	model := writeFlowModel(t, flowGatedNextModel)
	art := seedArtifact(t, model, "status=draft", "flag=true")
	bind := artifactBinding(flowStateRole, art)

	t.Run("without-the-flag", func(t *testing.T) {
		data := runNext(t, model, []string{bind})
		pairs := candidateUnknown(t, requireCandidate(t, data, "gated-reported"))

		if !hasUnknown(pairs, "reported", "not-evaluated") {
			t.Errorf("`unknown` does not carry {reported, not-evaluated}: "+
				"%#v\nThe gate id goes in the KEY member and the reason is "+
				"the token this RDR mints for the gate class only. It is "+
				"NOT `absent` — the gate is declared and its id is reported "+
				"— and NOT `uncomparable` — a gate id is not a view key "+
				"with a value", pairs)
		}
		// The set is the ROW's own declared gates, not a model-level one:
		// `gate.excluded` is declared by the model and belongs to a
		// different row, so it must not appear here.
		if unknownKeys(pairs, "excluded") {
			t.Errorf("`unknown` names the gate `excluded`, which this ROW "+
				"does not declare: %#v\nA build reading a wider set reports "+
				"ids the row does not declare", pairs)
		}
	})

	t.Run("with-the-flag", func(t *testing.T) {
		data := runNext(t, model, []string{bind}, "--evaluate-gates")
		pairs := candidateUnknown(t, requireCandidate(t, data, "gated-reported"))

		if unknownKeys(pairs, "reported") {
			t.Errorf("the gate id is still listed under `unknown` with "+
				"--evaluate-gates passed: %#v\nThe entry exists BECAUSE the "+
				"caller did not ask; once they do, the gate RAN and its "+
				"result rides the candidate instead", pairs)
		}
	})
}

// REQ-28: "The gate-id entries are emitted on EVERY reported disposition,
// including `owned_state_unavailable`"
// ADVERSARIAL — the class means "the caller did not ask", which is true of
// a row whose evaluation stopped early exactly as it is of a plan row, and
// the fidelity limit below is scoped to the kernel's GUARD payload, which
// gate ids do not come from.
func TestReq28_GateIdsAreEmittedEvenOnAnOwnedStateUnavailableCandidate(t *testing.T) {
	model := writeFlowModel(t, flowSortOrderModel)
	// `mid` is an owned key the row WRITES — so it enters `RequiresOwned`
	// — and which nothing establishes, so the probe refuses
	// `owned_state_unavailable`. The row is still a candidate.
	bind := seedMatchArtifact(t, model, "status=draft")

	data := runNext(t, model, []string{bind})
	pairs := candidateUnknown(t, requireCandidate(t, data, "sort-row"))

	if !hasUnknown(pairs, "mid", "absent") {
		t.Fatalf("the row does not report its unestablished owned key "+
			"{mid, absent}: %#v\nWithout this the assertion below cannot "+
			"tell an `owned_state_unavailable` disposition from a plan",
			pairs)
	}
	if !hasUnknown(pairs, "beta", "not-evaluated") {
		t.Errorf("the gate id is missing from a candidate whose probe "+
			"refused `owned_state_unavailable`: %#v\nGate-id entries are "+
			"emitted on EVERY reported disposition; the fidelity limit is "+
			"scoped to the kernel's guard payload, which gate ids do not "+
			"come from", pairs)
	}
}

// REQ-24: "Undecided guard atoms carry the reason the kernel reports —
// `absent` or `uncomparable`, `0007:C8`'s closed set — read off the probe's
// `guard_unevaluable` refusal payload (`Refusal.Undecided`); a guard atom
// over an absent key therefore appears once, the walk and the payload
// agreeing on `{key, absent}`."
// REQ-75 / `0011:D-undecided-vocabulary`: "reuses `0007:C8`'s closed reason
// set (`absent`, `uncomparable`) for every fact that comes off the kernel's
// guard payload, rather than minting one there; the CLI's own walk emits
// only `absent`."
// REQ-84: "guard atom over present key, seam unevaluable | yes | `{key,
// uncomparable}` (payload only) | yes | 0"
// REQ-112 / `0011:S7`: "a guard atom over a present key the seam cannot
// decide appears as `{key, uncomparable}` … the same atom over an absent
// key appears once as `{key, absent}` (walk and payload agree, dedup on the
// pair)".
// REQ-66: the fixture is A13's ONE reachable producer — a present value the
// operator cannot parse.
// DOMAIN EDGE — `uncomparable` is reachable ONLY through the payload read
// C1 adds; the CLI's own view walk cannot see it, because the key is
// PRESENT.
func TestReq24And66And75And84And112_AnUnevaluableGuardOverAPresentKeyIsUncomparable(t *testing.T) {
	// The value is ESTABLISHED under a model where `size` is a `scalar`
	// and `notanint` conforms, and READ under one where the same key is an
	// `int` the row guards with `gt` — the only CLI-reachable route to a
	// present value the operator cannot parse. Both hops go through the
	// CLI; nothing here asserts the artifact's on-disk format.
	writer := writeFlowModel(t, flowLooseWriterModel)
	art := newFlowArtifact(t, "loose.artifact")
	bind := artifactBinding(flowStateRole, art)
	requireSuccess(t, "flow", "set-state", "--model", writer,
		"--artifact", bind, "--write", "status=draft",
		"--write", "size=notanint", "--as=json")

	reader := writeFlowModel(t, flowUncomparableGuardModel)
	data := runNext(t, reader, []string{bind})

	t.Run("present-key-unevaluable", func(t *testing.T) {
		pairs := candidateUnknown(t, requireCandidate(t, data, "uncomparable-row"))
		if !hasUnknown(pairs, "size", "uncomparable") {
			t.Errorf("`unknown` does not carry {size, uncomparable}: %#v\n"+
				"`size` is PRESENT at a value `gt` cannot parse, so the "+
				"seam answers unevaluable and the kernel reports the reason "+
				"on the probe's `guard_unevaluable` payload. The CLI's own "+
				"view walk CANNOT see this fact — the key is present — so "+
				"it is reachable only through the payload read C1 adds, "+
				"and it appears on `next` for the FIRST time here", pairs)
		}
		if hasUnknown(pairs, "size", "absent") {
			t.Errorf("`unknown` reports {size, absent} for a PRESENT key: "+
				"%#v\n`absent` and `uncomparable` have DIFFERENT remedies "+
				"— bind a reader or supply the tag; fix the value or the "+
				"atom's literal — and the payload names which", pairs)
		}
	})

	t.Run("absent-key-appears-once", func(t *testing.T) {
		pairs := candidateUnknown(t, requireCandidate(t, data, "absent-guard-row"))
		var n int
		for _, p := range pairs {
			if p.Key == "gone" {
				n++
			}
		}
		if n != 1 {
			t.Errorf("the guard key `gone` appears %d times in `unknown`: "+
				"%#v; want exactly 1\nThe WALK and the PAYLOAD both produce "+
				"{gone, absent} for a guard atom over an absent key, and "+
				"dedup on the PAIR collapses them to one entry — this is "+
				"the reachable half of the dedup rule", n, pairs)
		}
		if !hasUnknown(pairs, "gone", "absent") {
			t.Errorf("`unknown` does not carry {gone, absent}: %#v", pairs)
		}
	})
}

// REQ-25: "That payload read MUST be guarded on the probe's disposition
// being `guard_unevaluable`, the only refusal that carries
// `Refusal.Undecided` (A13): `summarize` MUST NOT read the field off
// whatever `Result` it is handed and filter afterwards"
// ADVERSARIAL — a guard-decided-FALSE row surfaces as `no_match` on the
// probe and is not reported at all, so a build that read the field off
// whatever Result it was handed would have to invent facts for a row that
// is not on the wire. The observable form: an excluded row contributes
// nothing anywhere on the payload.
func TestReq25_TheGuardPayloadIsReadOnlyOnAGuardUnevaluableDisposition(t *testing.T) {
	model := writeFlowModel(t, flowGatedNextModel)
	art := seedArtifact(t, model, "status=draft", "flag=true")
	bind := artifactBinding(flowStateRole, art)

	data := runNext(t, model, []string{bind})

	// `gated-excluded`'s guard is decided FALSE. It is not a candidate, so
	// nothing about it — including any undecided-atom fact — may appear.
	if c := candidateNamed(t, data, "gated-excluded"); c != nil {
		t.Fatalf("`gated-excluded` is reported: %#v", c)
	}
	for _, c := range nextCandidates(t, data) {
		for _, p := range candidateUnknown(t, c) {
			if p.Key == "flag" {
				t.Errorf("candidate %v carries an `unknown` entry for "+
					"`flag`, the key the EXCLUDED row's guard decided "+
					"false: %#v\nThe payload read is guarded on the probe's "+
					"disposition being `guard_unevaluable` — the only "+
					"refusal that carries `Refusal.Undecided` — never read "+
					"off whatever `Result` is in hand and filtered after",
					c["rule"], p)
			}
		}
	}
}

// REQ-31: "when the probe refuses `owned_state_unavailable`, `gate` returns
// before it collects the guard payload, so an `uncomparable` guard atom on
// that row is not on the wire and is not reported; every `absent` fact on
// the row still is"
// REQ-85: "owned key no reader established | yes | `{key, absent}`; an
// `uncomparable` guard atom on the same row is not reported (precedence
// limit) | yes | 0"
// REQ-112: "asserted so the limit is pinned rather than discovered."
// DOMAIN EDGE — the one stated fidelity limit, asserted rather than left to
// be found by a future caller.
func TestReq31And85And112_TheOwnedStateUnavailablePrecedenceLimitIsPinned(t *testing.T) {
	writer := writeFlowModel(t, flowLooseWriterModel)
	art := newFlowArtifact(t, "loose.artifact")
	bind := artifactBinding(flowStateRole, art)
	// `size` is established at a non-integer; `missingowned` is NOT
	// established, so the probe refuses `owned_state_unavailable` before
	// `gate` ever collects the guard payload.
	requireSuccess(t, "flow", "set-state", "--model", writer,
		"--artifact", bind, "--write", "status=draft",
		"--write", "size=notanint", "--as=json")

	reader := writeFlowModel(t, flowOwnedUnavailableModel)
	data := runNext(t, reader, []string{bind})

	c := requireCandidate(t, data, "unavail-row")
	pairs := candidateUnknown(t, c)

	// The row is STILL a candidate, and every `absent` fact on it is
	// still reported — from the walk, which does not depend on the
	// payload.
	if !hasUnknown(pairs, "missingowned", "absent") {
		t.Errorf("the missing owned key is not reported {missingowned, "+
			"absent}: %#v\n`owned_state_unavailable` leaves the row a "+
			"candidate with its facts reported, and the missing owned key "+
			"is the remedy the caller acts on FIRST", pairs)
	}
	// And the `uncomparable` guard atom on the SAME row is NOT reported —
	// `gate` returned before collecting the payload, so the fact is not
	// on the wire. The limit is stated, so it is asserted.
	if hasUnknown(pairs, "size", "uncomparable") {
		t.Errorf("`unknown` carries {size, uncomparable} on a row whose "+
			"probe refused `owned_state_unavailable`: %#v\nC1 states the "+
			"limit rather than hiding it — `gate`'s `missingOwned` check "+
			"returns BEFORE the undecided-collection loop begins, so that "+
			"fact never reaches the wire. A build reporting it invented it",
			pairs)
	}
}

// REQ-29: "Entries are deduplicated by `{key, reason}` pair, not by key,
// and sorted by `(key, reason)`."
// REQ-113 / `0011:S7`: "`unknown` is SORTED by `(key, reason)` — asserted on
// a fixture whose facts are minted out of that order, so an unsorted build
// fails."
// BOUNDARY — the fixture mints `zebra`, `alpha` (atom walk, authored
// order), then `mid` (owned walk), then `beta` (gate loop). Sorted the list
// is `alpha, beta, mid, zebra` — an order NO single source produces, so a
// build that merely concatenates its sources fails here.
func TestReq29And113_UnknownIsSortedByKeyThenReasonAndDedupedOnThePair(t *testing.T) {
	model := writeFlowModel(t, flowSortOrderModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	data := runNext(t, model, []string{bind})
	pairs := candidateUnknown(t, requireCandidate(t, data, "sort-row"))

	got := make([]string, 0, len(pairs))
	for _, p := range pairs {
		got = append(got, p.Key+"\x00"+p.Reason)
	}
	want := slices.Clone(got)
	sort.Strings(want)
	if !slices.Equal(got, want) {
		t.Errorf("`unknown` is not sorted by (key, reason): %#v\nThe "+
			"fixture mints its facts OUT of that order — the atom walk "+
			"contributes `zebra` then `alpha`, the owned walk `mid`, the "+
			"gate loop `beta` — so a build that concatenates its sources "+
			"without sorting reports them in emission order and fails "+
			"here. The sort exists so a consumer diffing a golden payload "+
			"sees no churn", pairs)
	}

	// Dedup is on the PAIR: no pair repeats.
	seen := map[unknownPair]bool{}
	for _, p := range pairs {
		if seen[p] {
			t.Errorf("the pair %#v appears twice in `unknown`: %#v\n"+
				"Entries are deduplicated by `{key, reason}` pair, and "+
				"dedup is the LAST step", p, pairs)
		}
		seen[p] = true
	}
}

// REQ-29: "Entries are deduplicated by `{key, reason}` pair, not by key,
// and sorted by `(key, reason)`."
// REQ-113 / `0011:S7`: "`unknown` is SORTED by `(key, reason)` — asserted on
// a fixture whose facts are minted out of that order, so an unsorted build
// fails."
// BOUNDARY — the PAIR half of both clauses. `flowSortOrderModel` above
// declares four DISTINCT keys, so it cannot separate pair dedup from key
// dedup, nor the `(key, reason)` sort from a key-only sort: both broken
// builds stay green over it. This fixture supplies the discriminating
// input — ONE key minted twice with different reasons, from two different
// sources — which `flowUnknownPairModel` gets from a deliberate tag-key /
// gate-id namespace collision.
func TestReq29_DedupAndSortDiscriminateOnTheReasonNotOnlyTheKey(t *testing.T) {
	// `size` is established at a PRESENT non-integer through the writer
	// model, where it is a `scalar`; the pair model below declares it an
	// `int` and guards it with `gt`, so the operator cannot parse it and
	// the kernel reports `{size, uncomparable}`.
	writer := writeFlowModel(t, flowUnknownPairWriterModel)
	art := newFlowArtifact(t, "pair.artifact")
	bind := artifactBinding(flowMatchRole, art)
	requireSuccess(t, "flow", "set-state", "--model", writer,
		"--artifact", bind, "--write", "status=draft",
		"--write", "size=notanint", "--as=json")

	model := writeFlowModel(t, flowUnknownPairModel)
	data := runNext(t, model, []string{bind})
	pairs := candidateUnknown(t, requireCandidate(t, data, "pair-row"))

	// (i) BOTH same-key entries survive. A build keying `seen` on `f.Key`
	// instead of the whole pair collapses these two into one.
	for _, want := range []unknownPair{
		{Key: "beta", Reason: "absent"},        // atom walk (absent match key)
		{Key: "beta", Reason: "not-evaluated"}, // gate loop (un-run gate id)
	} {
		if !hasUnknown(pairs, want.Key, want.Reason) {
			t.Errorf("`unknown` is missing %#v: %#v\nThe key `beta` is BOTH "+
				"an observed match key nothing supplies and an un-run gate "+
				"id, so it is minted twice with DIFFERENT reasons. Dedup is "+
				"on the `{key, reason}` PAIR, not on the key — `absent` and "+
				"`not-evaluated` name different remedies, so collapsing "+
				"them loses one of them (REQ-29, `0011:C1`)", want, pairs)
		}
	}

	// (ii) exactly two `beta` entries — the pair dedup still DEDUPES, it
	// just does not over-collapse.
	betas := 0
	for _, p := range pairs {
		if p.Key == "beta" {
			betas++
		}
	}
	if betas != 2 {
		t.Errorf("`unknown` carries %d entries for `beta`; want exactly 2: "+
			"%#v\nOne per reason — dedup on the pair removes repeats of the "+
			"SAME pair without merging distinct ones", betas, pairs)
	}

	// (iii) the REASON is the sort tiebreak, and the ordering it produces
	// is TOTAL. Over two entries sharing a key a key-only comparison
	// returns 0, leaving their relative order to `slices.SortFunc`, which
	// is NOT stable — so a build without the tiebreak has no defined
	// answer here. Assert the full expected sequence rather than a
	// sortedness predicate, which a key-only comparator also satisfies.
	want := []unknownPair{
		{Key: "alpha", Reason: "absent"},
		{Key: "beta", Reason: "absent"},
		{Key: "beta", Reason: "not-evaluated"},
		{Key: "size", Reason: "not-evaluated"},
		{Key: "size", Reason: "uncomparable"},
	}
	if !slices.Equal(pairs, want) {
		t.Errorf("`unknown` = %#v; want %#v\nSorted by `(key, reason)` the "+
			"list is total: `alpha` before both `beta` entries (the KEY "+
			"comparison), and within a key the REASON breaks the tie. The "+
			"`size` pair is the DISCRIMINATING one — it is minted "+
			"`uncomparable` by the undecided-atom read (step 3) and "+
			"`not-evaluated` by the gate loop (step 4), and since "+
			"`not-evaluated` < `uncomparable` the LATER source contributes "+
			"the smaller reason, so emission order is the REVERSE of "+
			"sorted order. A build comparing only the key cannot restore "+
			"it: over two entries sharing a key a key-only comparison "+
			"returns 0 and `slices.SortFunc` is NOT stable, so the answer "+
			"is undefined (REQ-29, `0011:C1`)", pairs, want)
	}

	// (iv) the reason tiebreak is load-bearing, not incidental: on `size`
	// the required order is the OPPOSITE of the order the sources mint it
	// in, so a key-only comparator cannot produce this sequence by
	// preserving emission order the way it can for `beta`.
	iNot := slices.Index(pairs, unknownPair{Key: "size", Reason: "not-evaluated"})
	iUnc := slices.Index(pairs, unknownPair{Key: "size", Reason: "uncomparable"})
	if iNot < 0 || iUnc < 0 || iNot > iUnc {
		t.Errorf("`size` entries are not ordered `not-evaluated` before "+
			"`uncomparable`: %#v\nBoth must be present — the pair dedup "+
			"keeps them apart — and the REASON orders them against their "+
			"emission order, which mints `uncomparable` first", pairs)
	}
}

// REQ-124 / `0011:F3`: "Only guard atoms can carry this reason on `next`."
// REQ-129: "adds the reason token `not-evaluated` for un-run gate ids on
// that CLI list only — `0007:C8`'s seam vocabulary stays at
// `absent`/`uncomparable` and its payload is not widened."
// A-3: the reason values are the literal lowercase-hyphen tokens `absent`,
// `uncomparable`, `not-evaluated`.
// ADVERSARIAL — the vocabulary is closed and each token is scoped to a
// class: `uncomparable` never appears on a match, owned-key, or gate entry,
// and `not-evaluated` never appears on a match, guard, or owned-key entry.
func TestReq124And129_TheReasonVocabularyIsClosedAndEachTokenIsScoped(t *testing.T) {
	closed := []string{"absent", "uncomparable", "not-evaluated"}

	// Every reported candidate across the corpus, so an invented token
	// anywhere fails here.
	type run struct {
		name  string
		model string
		seed  []string
		extra []string
	}
	for _, r := range []run{
		{"match-classes", flowMatchClassesModel, []string{"status=draft"}, nil},
		{"absent-match", flowAbsentMatchKeyModel, []string{"status=draft"}, nil},
		{"sort-order", flowSortOrderModel, []string{"status=draft"}, nil},
		{"gated", flowGatedNextModel, []string{"status=draft", "flag=true"}, nil},
		{"gated-evaluated", flowGatedNextModel,
			[]string{"status=draft", "flag=true"}, []string{"--evaluate-gates"}},
	} {
		t.Run(r.name, func(t *testing.T) {
			model := writeFlowModel(t, r.model)
			bind := seedMatchArtifact(t, model, r.seed...)
			data := runNext(t, model, []string{bind}, r.extra...)

			for _, c := range nextCandidates(t, data) {
				for _, p := range candidateUnknown(t, c) {
					if !slices.Contains(closed, p.Reason) {
						t.Errorf("candidate %v carries the reason %q, which "+
							"is outside the closed set %v — `0007:C8`'s two "+
							"members plus the ONE token this RDR mints for "+
							"the gate class", c["rule"], p.Reason, closed)
					}
					if p.Reason == "not_evaluated" || p.Reason == "notEvaluated" {
						t.Errorf("the gate reason is spelled %q; the token "+
							"is the lowercase-HYPHEN `not-evaluated`",
							p.Reason)
					}
				}
			}
		})
	}
}

// REQ-124 / `0011:F3`: "Only guard atoms can carry this reason on `next`."
// ADVERSARIAL — `uncomparable` requires a PRESENT key the seam could not
// decide, which only a guard atom can produce: the CLI's own walk emits
// only `absent`, and a gate id is not a view key with a value.
func TestReq124_UncomparableNeverRidesAMatchOwnedKeyOrGateEntry(t *testing.T) {
	writer := writeFlowModel(t, flowLooseWriterModel)
	art := newFlowArtifact(t, "loose.artifact")
	bind := artifactBinding(flowStateRole, art)
	requireSuccess(t, "flow", "set-state", "--model", writer,
		"--artifact", bind, "--write", "status=draft",
		"--write", "size=notanint", "--as=json")

	reader := writeFlowModel(t, flowUncomparableGuardModel)
	data := runNext(t, reader, []string{bind})

	// The control: `uncomparable` IS reachable on this fixture, so the
	// scoping assertion below is not vacuous.
	control := candidateUnknown(t, requireCandidate(t, data, "uncomparable-row"))
	if !hasUnknown(control, "size", "uncomparable") {
		t.Fatalf("the fixture does not produce `uncomparable` at all: %#v\n"+
			"Without it the scoping assertion below passes vacuously",
			control)
	}

	// `size` is the ONLY key a guard atom names here. Any other key
	// carrying `uncomparable` would be a match, owned-key, or gate entry.
	for _, cand := range nextCandidates(t, data) {
		for _, p := range candidateUnknown(t, cand) {
			if p.Reason == "uncomparable" && p.Key != "size" {
				t.Errorf("candidate %v carries {%s, uncomparable}; only "+
					"GUARD atoms can carry this reason on `next`, and "+
					"`size` is the only guard key on this model that the "+
					"seam cannot decide", cand["rule"], p.Key)
			}
		}
	}
}

// REQ-72 / `0011:D-undecided-reporting-shape`: "Text mode (`--as=text`)
// carries the same pairs — both members of each, so the text caller gets
// the diagnosis and not just the symptom"
// REQ-73: "It does NOT get a bespoke `key (reason)` form: `flow next`
// declares no `FindingCarrier` and `respond.OK` therefore renders its
// payload through the shared generic flattener … which emits path-qualified
// leaf lines — `candidates[0].unknown[0].key: stage`,
// `candidates[0].unknown[0].reason: absent`."
// DOMAIN EDGE — the RDR says "no scenario asserts a text-mode string", so
// this oracle asserts the INFORMATION (both members reach text mode) and
// the NEGATIVE (no per-verb template was minted), never an exact line.
func TestReq72And73_TextModeCarriesBothMembersThroughTheSharedFlattener(t *testing.T) {
	model := writeFlowModel(t, flowAbsentMatchKeyModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	text := requireSuccess(t, "flow", "next", "--model", model,
		"--artifact", bind, "--as=text")

	// Both MEMBERS of the pair must reach the text caller, so the
	// diagnosis travels with the symptom. Assert the PATH-QUALIFIED leaf,
	// not a bare substring: this fixture's rule ids (`absent-row`,
	// `all-absent`) carry the word "absent" on their own, so a build that
	// dropped every `reason` leaf would satisfy a member scan while
	// reporting a symptom with no diagnosis attached.
	for _, leaf := range []string{"key: wanted", "reason: absent"} {
		if !strings.Contains(text, leaf) {
			t.Errorf("text mode does not carry the %q leaf:\n%s\nThe text "+
				"caller gets BOTH members of each pair — the diagnosis and "+
				"not just the symptom — which is the property that matters",
				leaf, text)
		}
	}

	// And it is the SHARED generic flattener, not a per-verb template:
	// the leaf lines are path-qualified. Minting a bespoke `key (reason)`
	// form to get `wanted (absent)` would edit a predecessor's surface for
	// cosmetics — the flattener exists precisely so "the two modes cannot
	// disagree about what the run reported".
	if strings.Contains(text, "wanted (absent)") {
		t.Errorf("text mode renders the bespoke `key (reason)` form:\n%s\n"+
			"`flow next` declares no `FindingCarrier`, so `respond.OK` "+
			"renders through the shared flattener; the rendering is the "+
			"gateway's to change, uniformly, and this contract owns the "+
			"payload's INFORMATION, not its text layout", text)
	}
	if !strings.Contains(text, "unknown") {
		t.Errorf("text mode names no `unknown` path at all:\n%s\nThe "+
			"flattener emits path-qualified leaf lines such as "+
			"`candidates[0].unknown[0].key`", text)
	}
}

// REQ-125 / `0011:G-cross-cutting`: "Version marker: none minted — the
// payload field rename `unresolved` → `unknown` IS the break"
// ADVERSARIAL — no schema or version field may be added to the payload to
// soften the rename, and the OLD field must be gone: a build emitting both
// has not made the break, it has widened the surface.
func TestReq125_NoVersionMarkerIsMintedAndTheOldFieldIsGone(t *testing.T) {
	model := writeFlowModel(t, flowAbsentMatchKeyModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	data := runNext(t, model, []string{bind})

	for _, field := range []string{"schema", "version", "payload_version", "$schema"} {
		if _, held := data[field]; held {
			t.Errorf("the payload carries a %q field; NO version marker is "+
				"minted — the field rename IS the break, and callers "+
				"migrate on it", field)
		}
	}

	for _, c := range nextCandidates(t, data) {
		if _, held := c["unresolved"]; held {
			t.Errorf("candidate %v still carries the OLD `unresolved` "+
				"field alongside `unknown`: %v\nThe rename is the break; "+
				"emitting both widens the surface instead of making it",
				c["rule"], keysOf(c))
		}
		// And the new field is a list of OBJECTS, not of strings — the
		// element type change is half of what the five re-homings are for.
		//
		// An EMPTY JSON array carries no element type: `stringsAt` returns
		// `([], true)` over `[]` in every possible build, so the check is
		// scoped to a candidate that actually carries an entry. `absent-row`
		// does (`{wanted, absent}`), which is what makes the assertion
		// discriminating rather than vacuous. Recorded as DEV-7.
		if pairs := candidateUnknown(t, c); len(pairs) == 0 {
			continue
		}
		if _, ok := stringsAt(c, "unknown"); ok {
			t.Errorf("candidate %v's `unknown` decodes as a []string; its "+
				"elements are `{key, reason}` OBJECTS. A flat list cannot "+
				"carry the remedy, which is why it was rejected",
				c["rule"])
		}
	}
}
