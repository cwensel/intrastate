package table_test

// RDR 0002's Minimum Viable Validation and the resolver handshake it
// depends on.

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
	"github.com/newcoinc/intrastate/internal/table"
)

// fixtureGuards is a value stand-in for RDR 0003's guard evaluator seam.
//
// The kernel decides PRESENCE itself and never hands the seam an absent
// key (`0007:C8` — an absent key is ReasonAbsent, decided by the kernel),
// so this seam is only ever asked to compare a PRESENT value against a
// literal. It implements just enough of RDR 0003's typed semantics for
// the fixture's own atoms: `lt` over an int, `eq` over a scalar, and
// `exists` never reaches it because the kernel owns that operator.
//
// This mocks the SEAM, never the unit under test: the loader and
// normalizer under test are the real ones, and the kernel is the real
// resolve.Resolve.
type fixtureGuards struct{}

func (fixtureGuards) Evaluate(atom resolve.GuardAtom, value string) resolve.GuardResult {
	switch atom.Operator {
	case "lt":
		got, err1 := strconv.Atoi(value)
		want, err2 := strconv.Atoi(atom.Literal)
		if err1 != nil || err2 != nil {
			return resolve.GuardUnevaluable
		}
		if got < want {
			return resolve.GuardTrue
		}
		return resolve.GuardFalse
	case "eq":
		if value == atom.Literal {
			return resolve.GuardTrue
		}
		return resolve.GuardFalse
	}
	return resolve.GuardUnevaluable
}

// ownedTags builds the accessor-produced owned snapshot from a map.
func ownedTags(kv map[string]string) []resolve.Tag {
	out := make([]resolve.Tag, 0, len(kv))
	for k, v := range kv {
		out = append(out, resolve.Tag{Key: k, Value: v})
	}
	return out
}

// REQ-137: "Run `internal/resolve::Resolve` over one matching ordinary
// tag-set …, one tag-set with no ordinary match but one matching
// `no_match` escape row, and one tag-set in which a sibling candidate's
// guard is unevaluable … **Expected**: … the unevaluable-sibling tag-set
// refuses `guard_unevaluable` even though a decidable sibling and a
// `no_match` escape row exist"
// ADVERSARIAL (TS 4)
//
// All three tag-sets are drawn from the fixture's two same-outcome sibling
// rows — without them the gate's ordering is untestable, since a
// one-candidate-per-outcome fixture can never exercise the third leg.
func TestReq137_ResolveOverTheNormalizedFixtureRows(t *testing.T) {
	m := mustLoad(t, rdrFixture)
	tbl := m.KernelTable()

	t.Run("one matching ordinary tag-set resolves to one transition row", func(t *testing.T) {
		// Every sibling candidate's guard is decidable: `cluster_ready` is
		// present (so continue-prelock-cluster's guard decides FALSE) and
		// `finalized_at` is absent (so the exists=false atom decides TRUE).
		owned := map[string]string{
			"status": "Draft", "stage": "prelock", "profile": "large",
			"iter": "1", "cluster_ready": "false",
		}
		res, err := resolve.Resolve(resolve.Input{
			Flow:       "rdr",
			Table:      tbl,
			Owned:      ownedTags(owned),
			Recognized: "round-clean",
			Guards:     fixtureGuards{},
		})
		if err != nil {
			t.Fatalf("the Go error path was used: %v", err)
		}
		if res.Refusal != nil {
			t.Fatalf("refused %q; want one transition row", res.Refusal.Kind)
		}
		if res.Plan == nil {
			t.Fatal("no plan and no refusal")
		}
		if res.Plan.RuleID != "continue-prelock" {
			t.Errorf("plan names rule %q; want continue-prelock", res.Plan.RuleID)
		}
		if res.Plan.Escaped {
			t.Error("the ordinary match came back as an escape")
		}
	})

	t.Run("no ordinary match resolves to the modeled escape row", func(t *testing.T) {
		// `stage` is not prelock, so neither sibling matches; the escape
		// row's only atom is status.eq=Draft, which holds.
		owned := map[string]string{
			"status": "Draft", "stage": "propose", "profile": "large",
			"iter": "1", "cluster_ready": "false",
		}
		res, err := resolve.Resolve(resolve.Input{
			Flow:       "rdr",
			Table:      tbl,
			Owned:      ownedTags(owned),
			Recognized: "round-clean",
			Guards:     fixtureGuards{},
		})
		if err != nil {
			t.Fatalf("the Go error path was used: %v", err)
		}
		if res.Refusal != nil {
			t.Fatalf("refused %q; want the modeled escape disposition", res.Refusal.Kind)
		}
		if res.Plan == nil || !res.Plan.Escaped {
			t.Fatalf("plan = %+v; want an escaped disposition", res.Plan)
		}
		if res.Plan.RuleID != "draft-no-match-escape" {
			t.Errorf("escape row = %q; want draft-no-match-escape", res.Plan.RuleID)
		}
	})

	t.Run("an unevaluable sibling refuses guard_unevaluable", func(t *testing.T) {
		// `cluster_ready` is ABSENT from the view, so
		// continue-prelock-cluster's `guard.all` eq atom over that optional
		// owned key is unevaluable — even though continue-prelock is
		// decidable and a no_match escape row exists.
		owned := map[string]string{
			"status": "Draft", "stage": "prelock", "profile": "large", "iter": "1",
		}
		res, err := resolve.Resolve(resolve.Input{
			Flow:       "rdr",
			Table:      tbl,
			Owned:      ownedTags(owned),
			Recognized: "round-clean",
			Guards:     fixtureGuards{},
		})
		if err != nil {
			t.Fatalf("the Go error path was used: %v", err)
		}
		if res.Refusal == nil {
			t.Fatalf("resolved to plan %+v; want a guard_unevaluable refusal — "+
				"an unevaluable survivor refuses even though a decidable sibling "+
				"and a no_match escape row both exist", res.Plan)
		}
		if res.Refusal.Kind != resolve.KindGuardUnevaluable {
			t.Errorf("refusal kind = %q; want %q",
				res.Refusal.Kind, resolve.KindGuardUnevaluable)
		}
	})
}

// REQ-MVV: "Parse two hand-authored sparse TOML fixtures, one for a
// representative RDR flow slice and one for a representative kata flow
// slice, into typed source data; normalize them into candidate rows; dump
// the expanded table; validate tag declarations, context references,
// predicate references, recognized outcomes, supported model version,
// multi-tag writes, outcome binding, and escape rows; then prove by unit
// test that one sample tag-set resolves through `internal/resolve::Resolve`
// to exactly one ordinary row, one unmatched sample tag-set resolves to
// exactly one modeled escape row when the table declares one, one
// unsupported-version variant is refused before normalization, and one
// deliberately overlapping variant is reported by lint as an ambiguous
// overlap."
// HAPPY PATH
//
// The gating Minimum Viable Validation for RDR 0002, as one runnable
// end-to-end test.
//
// The Round-Trip / Inverse Invariant this RDR declares is
// `parse ∘ normalize = expanded-table value identity`, so leg 7 compares
// the reconstructed value VALUE-FOR-VALUE against the original across all
// three semantics-preserving permutations, and leg 3 compares repeated
// dumps BYTE-FOR-BYTE. A green exit or "did not error" is explicitly not
// sufficient anywhere below.
//
// The overlap item (leg 9) is the one MVV assertion this RDR cannot
// discharge alone: the ambiguous-overlap check is cross-row and therefore
// RDR 0006's by the arity split, and RDR 0006 is unimplemented. It is
// DEFERRED to the Phase 5 lint handshake rather than counted as satisfied
// here — marking it green against a stub would be asserting on nothing.
// What leg 9 asserts instead is what this RDR actually owes before that
// handshake: that the overlapping rows carry enough structure for the
// check.
func TestMVV_ParseNormalizeDumpValidateAndResolve(t *testing.T) {
	t.Run("1_parse_two_fixtures_into_typed_source_data", func(t *testing.T) {
		rdr := mustLoad(t, rdrFixture)
		kata := mustLoad(t, kataFixture)

		if rdr.ID != "rdr" || kata.ID != "kata" {
			t.Fatalf("model ids = %q / %q; want rdr / kata", rdr.ID, kata.ID)
		}
		// Typed source data, not a bag of strings: the type model, the
		// accessor entries, and the metadata all decode into fields.
		if rdr.Tags["iter"].Kind != "int" {
			t.Errorf("iter kind = %q; want int", rdr.Tags["iter"].Kind)
		}
		if !rdr.Writers["rdr-status"].ReadBack {
			t.Error("the writer entry did not decode read_back")
		}
		if !reflect.DeepEqual(kata.Metadata, map[string]any{"owner": "kata-team"}) {
			t.Errorf("kata metadata = %#v; want the authored table", kata.Metadata)
		}
	})

	t.Run("2_normalize_into_candidate_rows", func(t *testing.T) {
		rdr := mustLoad(t, rdrFixture)
		kata := mustLoad(t, kataFixture)

		if got := len(rdr.Rows); got != 9 {
			t.Errorf("RDR fixture normalized to %d rows; want 9 (%v)",
				got, rowIdentities(rdr))
		}
		if got := len(kata.Rows); got != 2 {
			t.Errorf("kata fixture normalized to %d rows; want 2 (%v)",
				got, rowIdentities(kata))
		}

		// The approved row identities, value-for-value.
		want := []string{
			"rdr.continue-prelock#foundational",
			"rdr.continue-prelock#large",
			"rdr.continue-prelock-cluster#foundational",
			"rdr.continue-prelock-cluster#large",
			"rdr.draft-no-match-escape#reconcile-block",
			"rdr.draft-no-match-escape#round-clean",
			"rdr.reconcile-rewind",
			"rdr.terminal-archive#finalized",
			"rdr.terminal-archive#verdict-flapping",
		}
		if !reflect.DeepEqual(rowIdentities(rdr), want) {
			t.Errorf("row identities =\n %v\nwant\n %v", rowIdentities(rdr), want)
		}
	})

	t.Run("3_dump_the_expanded_table", func(t *testing.T) {
		rdr := mustLoad(t, rdrFixture)
		out := table.Dump(rdr)

		if strings.TrimSpace(out) == "" {
			t.Fatal("the dump is empty")
		}
		for _, r := range rdr.Rows {
			if !strings.Contains(out, r.Identity()) {
				t.Errorf("the dump omits row %s", r.Identity())
			}
		}
		// Byte-for-byte across repeated dumps, not "did not error".
		for i := range 8 {
			if got := table.Dump(mustLoad(t, rdrFixture)); got != out {
				t.Fatalf("dump %d is not byte-identical to the first", i)
			}
		}
	})

	t.Run("4_validate_the_eight_named_families", func(t *testing.T) {
		families := map[string]struct {
			rel  string
			want table.Category
		}{
			"tag declarations":     {"neg/neg-tagdecl-domain-on-bool.toml", table.CatMalformedTagDeclaration},
			"context references":   {"neg/neg-unknown-context.toml", table.CatUnknownContext},
			"predicate references": {"neg/neg-unknown-tag-match.toml", table.CatUnknownTag},
			"recognized outcomes":  {"neg/neg-dup-alphabet-member.toml", table.CatMalformedRecognizedOutcomeAlphabet},
			"supported version":    {"neg/neg-version.toml", table.CatUnsupportedVersion},
			"outcome binding":      {"neg/neg-outcome-outside.toml", table.CatMalformedOutcomeBinding},
			"escape rows":          {"neg/neg-escape-with-write.toml", table.CatMalformedEscapeDeclaration},
		}
		for name, spec := range families {
			if got := loadCategory(t, spec.rel); got != spec.want {
				t.Errorf("%s: category = %q; want %q", name, got, spec.want)
			}
		}

		// Multi-tag writes are validated positively: they must SURVIVE.
		row := rowByID(t, mustLoad(t, rdrFixture), "rdr.reconcile-rewind")
		if len(row.Writes) != 4 {
			t.Errorf("the multi-tag write normalized to %d writes; want 4 "+
				"(three assignments plus one rendered clear): %+v",
				len(row.Writes), row.Writes)
		}
	})

	t.Run("5_one_sample_tag_set_resolves_to_exactly_one_ordinary_row", func(t *testing.T) {
		m := mustLoad(t, rdrFixture)
		owned := map[string]string{
			"status": "Draft", "stage": "prelock", "profile": "large",
			"iter": "1", "cluster_ready": "false",
		}
		res, err := resolve.Resolve(resolve.Input{
			Flow:       "rdr",
			Table:      m.KernelTable(),
			Owned:      ownedTags(owned),
			Recognized: "round-clean",
			Guards:     fixtureGuards{},
		})
		if err != nil {
			t.Fatalf("the Go error path was used: %v", err)
		}
		if res.Refusal != nil {
			t.Fatalf("refused %q; want exactly one ordinary row", res.Refusal.Kind)
		}
		if res.Plan == nil {
			t.Fatal("neither plan nor refusal")
		}
		if res.Plan.Escaped {
			t.Error("the disposition is an escape; want an ordinary row")
		}
		if res.Plan.RuleID != "continue-prelock" {
			t.Errorf("plan rule = %q; want continue-prelock", res.Plan.RuleID)
		}
		// Value-for-value on the plan's substance, not just its identity.
		wantNext := []resolve.Tag{{Key: "iter", Value: "2"}, {Key: "stage", Value: "prelock"}}
		if !reflect.DeepEqual(res.Plan.NextTags, wantNext) {
			t.Errorf("plan next tags = %+v; want %+v", res.Plan.NextTags, wantNext)
		}
		if !reflect.DeepEqual(res.Plan.Writes, wantNext) {
			t.Errorf("plan writes = %+v; want %+v", res.Plan.Writes, wantNext)
		}
	})

	t.Run("6_one_unmatched_tag_set_resolves_to_one_modeled_escape_row", func(t *testing.T) {
		m := mustLoad(t, rdrFixture)
		owned := map[string]string{
			"status": "Draft", "stage": "propose", "profile": "large",
			"iter": "1", "cluster_ready": "false",
		}
		res, err := resolve.Resolve(resolve.Input{
			Flow:       "rdr",
			Table:      m.KernelTable(),
			Owned:      ownedTags(owned),
			Recognized: "round-clean",
			Guards:     fixtureGuards{},
		})
		if err != nil {
			t.Fatalf("the Go error path was used: %v", err)
		}
		if res.Refusal != nil {
			t.Fatalf("refused %q; the table declares a no_match escape row for "+
				"this outcome", res.Refusal.Kind)
		}
		if res.Plan == nil || !res.Plan.Escaped {
			t.Fatalf("plan = %+v; want exactly one modeled escape row", res.Plan)
		}
		if res.Plan.RuleID != "draft-no-match-escape" {
			t.Errorf("escape row = %q; want draft-no-match-escape", res.Plan.RuleID)
		}
		// An escape row yields no writes.
		if len(res.Plan.Writes) != 0 || len(res.Plan.NextTags) != 0 {
			t.Errorf("the escape plan carries next=%+v writes=%+v; both must be "+
				"empty", res.Plan.NextTags, res.Plan.Writes)
		}
	})

	t.Run("7_round_trip_value_identity", func(t *testing.T) {
		// `parse ∘ normalize = expanded-table value identity`. The
		// reconstructed value is compared VALUE-FOR-VALUE against the
		// original over the full field list, across all three
		// semantics-preserving authorings.
		base := mustLoad(t, rdrFixture)
		for _, rel := range []string{
			"perm/rdr-keyorder.toml",
			"perm/rdr-ruleorder.toml",
			"perm/rdr-eq-as-in.toml",
		} {
			got := mustLoad(t, rel)
			if !reflect.DeepEqual(got.Rows, base.Rows) {
				t.Errorf("%s did not reconstruct the original candidate-row set "+
					"value-for-value:\n got = %+v\nwant = %+v", rel, got.Rows, base.Rows)
			}
			// And byte-for-byte through the dump, which is derived from
			// that value.
			if table.Dump(got) != table.Dump(base) {
				t.Errorf("%s dumped differently from the baseline", rel)
			}
		}
	})

	t.Run("8_unsupported_version_is_refused_before_normalization", func(t *testing.T) {
		m, err := table.Load(readFixture(t, "neg/neg-version.toml"), "neg-version.toml")
		if err == nil {
			t.Fatal("an unsupported version loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatUnsupportedVersion {
			t.Errorf("category = %q; want %q", cat, table.CatUnsupportedVersion)
		}
		// "before normalization": no candidate rows are produced.
		if m != nil && len(m.Rows) > 0 {
			t.Errorf("the refused document yielded %d candidate rows", len(m.Rows))
		}

		// And the precedence control: the version gate runs before strict
		// field validation (A14).
		if got := loadCategory(t, "neg/neg-v2-shaped.toml"); got != table.CatUnsupportedVersion {
			t.Errorf("the v2-shaped precedence control refused %q; want %q",
				got, table.CatUnsupportedVersion)
		}
	})

	t.Run("9_overlap_is_deferred_to_the_phase_5_lint_handshake", func(t *testing.T) {
		// DEFERRED, not satisfied: RDR 0006 is unimplemented, and this RDR
		// refuses to count an assertion against a stub. What is asserted
		// here is what this RDR owes before the handshake — the overlapping
		// rows are distinguishable by row identity and comparable by
		// predicate set.
		base := string(readFixture(t, rdrFixture))
		src := strings.Replace(base, "[rule.guard.all.cluster_ready]\neq = true\n", "", 1)
		if src == base {
			t.Fatal("overlap substitution did not apply")
		}
		m, err := table.Load([]byte(src), "overlapping.toml")
		if err != nil {
			t.Fatalf("the overlapping variant refused at load with %v; ambiguous "+
				"overlap is cross-row and therefore RDR 0006's", err)
		}

		a := rowByID(t, m, "rdr.continue-prelock#large")
		b := rowByID(t, m, "rdr.continue-prelock-cluster#large")
		if a.Identity() == b.Identity() {
			t.Error("the overlapping rows are not distinguishable by row identity")
		}
		if a.Outcome != b.Outcome {
			t.Fatalf("the rows bind %q and %q; the variant is not an overlap",
				a.Outcome, b.Outcome)
		}
		if len(a.Atoms) == 0 || len(b.Atoms) == 0 {
			t.Error("the overlapping rows are not comparable by predicate set")
		}
		// The lint handshake surface exists and exposes both rows.
		if len(m.KernelTable().Rows) != len(m.Rows) {
			t.Error("the exposed representation loses rows")
		}
	})
}
