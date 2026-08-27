package table_test

// RDR 0002 Phase 3c — regression tests for the CoVe (Phase 3a) findings.
//
// FAIL-1 and FAIL-2 were found by probe rather than by a shipped test, so
// each gets an oracle here that fails against the pre-fix implementation.
// The Phase 3b adversarial tests already guard ADV-1/2/3 and are untouched.
//
// These are ADDED tests. Nothing existing is weakened.

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// ---------------------------------------------------------------- FAIL-1

// FAIL-1 — a multi-member literal was silently truncated at the kernel
// seam.
//
// REQ-93 (§D13 landing): "`Tag.Value` stays `string`; a set crosses as its
// canonical JSON array — members sorted, duplicate-free, compact encoding",
// covering "any set-valued atom literal handed to the guard seam". REQ-56
// forbids rendering a set-valued literal into a single string; here the
// surplus members were not merely joined, they were DROPPED.
//
// `model.go::seamValue` gated the array encoding on the tag's DECLARED
// kind, so on any non-`set` tag a multi-member literal crossed as
// `members[0]` alone. The resolution splits the question in two, because
// the seam asks two different ones:
//
//   - A tag VALUE (`Tag.Value` on a match tag, a next-state tag, a write)
//     is a set by its declared KIND — unchanged, and what §D13 governs.
//   - A guard atom LITERAL is a set by its OPERATOR. RDR 0003 fences `eq`
//     and the integer comparisons as SINGLE-VALUE operators ("the tag's
//     single held value IS the literal"), leaving `in` and `contains` as
//     the two whose right-hand side is a member set — exactly the two the
//     shipped cross-RDR contract test drives as JSON arrays
//     (`internal/resolve/guardcontract.go`).
//
// Both arms of the fix are asserted: the single-value operators refuse a
// multi-member literal at LOAD (`0002:C3` — a stable refusal, never a
// silent reinterpretation), and `in`/`contains` cross the seam whole.
func TestFail1_MultiMemberLiteralIsNeverTruncatedAtTheKernelSeam(t *testing.T) {
	// `status` is `scalar`-kind, so nothing carried on it is a declared
	// `set`: that is what shows carriage is not decided by the kind alone.
	// `labels` is the declared `set` the `contains` arms need, since RDR
	// 0003's operator/kind matrix admits `contains` over `set` alone — the
	// carriage question the arms ask is unchanged by which tag carries it,
	// and the `in` arm still asks it on the `scalar`.
	build := func(guard, write string) string {
		return `outcomes = ["done"]

[model]
id = "m"
version = 1

[tags.status]
provenance = "owned"
kind = "scalar"

[tags.labels]
provenance = "owned"
kind = "set"

[tags.recognized]
provenance = "recognized"
kind = "enum"
required = true

[read.r]
role = "m"
path = "m"
keys = ["labels", "status"]
timeout = "2s"

[write.w]
role = "m"
path = "m"
keys = ["labels", "status"]
timeout = "2s"
read_back = true

[[rule]]
id = "advance"
[rule.match.recognized]
eq = "done"
` + guard + `
[rule.write]
status = ` + write + `
`
	}

	guardLiteralOn := func(t *testing.T, m *table.Model, key, operator string) string {
		t.Helper()
		kr := m.Rows[0].KernelRow()
		for _, a := range kr.Guard {
			if a.Key == key && a.Operator == operator {
				return a.Literal
			}
		}
		t.Fatalf("no `%s` %s guard atom reached the kernel row: %+v",
			key, operator, kr.Guard)
		return ""
	}
	guardLiteral := func(t *testing.T, m *table.Model, operator string) string {
		t.Helper()
		return guardLiteralOn(t, m, "status", operator)
	}

	// --- the two set-valued operators cross whole ---

	t.Run("a guard `contains` literal crosses as the §D13 array", func(t *testing.T) {
		m, err := table.Load([]byte(build(
			"[rule.guard.all.labels]\ncontains = [\"done\", \"open\"]", `"done"`)),
			"fail1-contains.toml")
		if err != nil {
			t.Fatalf("refused: %v", err)
		}
		got := guardLiteralOn(t, m, "labels", "contains")
		if got != `["done","open"]` {
			t.Errorf("GuardAtom.Literal = %q; want %q — REQ-93 sends any "+
				"set-valued atom literal across as its canonical JSON array, "+
				"and `guardcontract.go` drives `contains` against exactly "+
				"that byte form", got, `["done","open"]`)
		}
	})

	t.Run("a guard `in` literal crosses as the §D13 array", func(t *testing.T) {
		m, err := table.Load([]byte(build(
			"[rule.guard.all.status]\nin = [\"alpha\", \"beta\"]", `"done"`)),
			"fail1-in.toml")
		if err != nil {
			t.Fatalf("refused: %v", err)
		}
		got := guardLiteral(t, m, "in")
		if got != `["alpha","beta"]` {
			t.Errorf("GuardAtom.Literal = %q; want %q. This is the sharpest "+
				"arm: RDR 0007's evaluator parses an `in` literal as a JSON "+
				"array, so a bare string must answer GuardUnevaluable and the "+
				"unevaluable-candidate veto then DEADLOCKS the flow rather "+
				"than refusing", got, `["alpha","beta"]`)
		}

		// Decoding it back must yield every authored member. A truncation
		// is invisible to a byte comparison alone if the encoding ever
		// changes shape, so assert the members too.
		var members []string
		if err := json.Unmarshal([]byte(got), &members); err != nil {
			t.Fatalf("literal %q is not a JSON array: %v", got, err)
		}
		if len(members) != 2 {
			t.Errorf("literal carries %d members; the author wrote 2 — no "+
				"member may be dropped at the seam", len(members))
		}
	})

	t.Run("the encoding is canonical: sorted, duplicate-free, compact", func(t *testing.T) {
		m, err := table.Load([]byte(build(
			"[rule.guard.all.labels]\ncontains = [\"zulu\", \"alpha\"]", `"done"`)),
			"fail1-canonical.toml")
		if err != nil {
			t.Fatalf("refused: %v", err)
		}
		if got := guardLiteralOn(t, m, "labels", "contains"); got != `["alpha","zulu"]` {
			t.Errorf("GuardAtom.Literal = %q; want the sorted, compact %q",
				got, `["alpha","zulu"]`)
		}
	})

	t.Run("two literals sharing a member stay distinguishable", func(t *testing.T) {
		// The consequence FAIL-1 demonstrated: `[AAA open]` and
		// `[ZZZ open]` crossed as "AAA" and "ZZZ", both dropping the
		// member the author shared between them.
		load := func(members string) string {
			m, err := table.Load([]byte(build(
				"[rule.guard.all.labels]\ncontains = "+members, `"done"`)),
				"fail1-share.toml")
			if err != nil {
				t.Fatalf("refused: %v", err)
			}
			return guardLiteralOn(t, m, "labels", "contains")
		}
		a := load(`["AAA", "open"]`)
		b := load(`["ZZZ", "open"]`)
		if !strings.Contains(a, "open") || !strings.Contains(b, "open") {
			t.Errorf("the shared member `open` did not reach the seam: %q / %q", a, b)
		}
		if a == b {
			t.Errorf("two distinct literals crossed the seam as the same bytes %q", a)
		}
	})

	// --- the single-value operators refuse instead ---

	t.Run("a multi-member `eq` literal is refused at load", func(t *testing.T) {
		for _, block := range []string{
			"[rule.guard.all.status]",
			"[rule.guard.unless.status]",
			"[rule.match.status]",
		} {
			_, err := table.Load([]byte(build(
				block+"\neq = [\"done\", \"open\"]", `"done"`)),
				"fail1-eq.toml")
			if err == nil {
				t.Errorf("%s: a multi-member `eq` literal loaded clean; RDR "+
					"0003 fences `eq` as a single-value operator, so it has "+
					"no denotation and `0002:C3` requires a stable refusal "+
					"rather than a silent truncation at the seam", block)
				continue
			}
			if cat, _ := table.CategoryOf(err); cat != table.CatMalformedPredicateAtom {
				t.Errorf("%s: category = %q; want %q (got: %v)",
					block, cat, table.CatMalformedPredicateAtom, err)
			}
		}
	})

	t.Run("a one-member `eq` array stays admitted", func(t *testing.T) {
		// `0002:C13` fixes `eq = "x"` and the one-member spelling as one
		// spelling of one edge, so the refusal must key on arity above
		// one, never on array-ness.
		if _, err := table.Load([]byte(build(
			"[rule.guard.all.status]\neq = [\"done\"]", `"done"`)),
			"fail1-eq-one.toml"); err != nil {
			t.Errorf("a one-member `eq` array refused: %v", err)
		}
	})

	t.Run("a multi-member write on a non-set kind is refused at load", func(t *testing.T) {
		_, err := table.Load([]byte(build(
			"[rule.guard.all.status]\nexists = true", `["done", "open"]`)),
			"fail1-write.toml")
		if err == nil {
			t.Fatal("a multi-member write on a `scalar` tag loaded clean; " +
				"`0002:C4` gives the array literal a meaning for a `set` kind " +
				"only, and RDR 0004's read-back compares the seam value for " +
				"equality")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatMalformedTagDeclaration {
			t.Errorf("category = %q; want %q (got: %v)",
				cat, table.CatMalformedTagDeclaration, err)
		}
	})

	t.Run("a multi-member [initial] value on a non-set kind is refused", func(t *testing.T) {
		src := strings.Replace(build(
			"[rule.guard.all.status]\nexists = true", `"done"`),
			"[tags.status]", "[initial]\nstatus = [\"done\", \"open\"]\n\n[tags.status]", 1)
		_, err := table.Load([]byte(src), "fail1-initial.toml")
		if err == nil {
			t.Fatal("a multi-member `[initial]` value on a `scalar` tag " +
				"loaded clean; REQ-30 requires every value to be well-formed " +
				"for the declared kind, and arity is part of that")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatMalformedInitialDeclaration {
			t.Errorf("category = %q; want %q (got: %v)",
				cat, table.CatMalformedInitialDeclaration, err)
		}
	})
}

// TestFail1b_GuardLiteralCarriageMatchesTheShippedSeamContract pins the
// resolution against the cross-RDR contract itself rather than against a
// literal this package chose. `internal/resolve/guardcontract.go` is the
// contract RDR 0007 shipped and RDR 0003's evaluator is instantiated
// against; its `in` and `contains` cases fix the literal's byte form. This
// asserts that a literal authored here, on a tag whose declared kind is
// NOT `set`, reaches the seam in a form that contract accepts.
func TestFail1b_GuardLiteralCarriageMatchesTheShippedSeamContract(t *testing.T) {
	src := `outcomes = ["done"]

[model]
id = "m"
version = 1

[tags.subject]
provenance = "owned"
kind = "scalar"

[tags.recognized]
provenance = "recognized"
kind = "enum"
required = true

[read.r]
role = "m"
path = "m"
keys = ["subject"]
timeout = "2s"

[write.w]
role = "m"
path = "m"
keys = ["subject"]
timeout = "2s"
read_back = true

[[rule]]
id = "advance"
[rule.match.recognized]
eq = "done"
[rule.guard.all.subject]
in = ["alpha", "beta"]
[rule.write]
subject = "done"
`
	m, err := table.Load([]byte(src), "fail1b.toml")
	if err != nil {
		t.Fatalf("refused: %v", err)
	}

	var atom resolve.GuardAtom
	for _, a := range m.Rows[0].KernelRow().Guard {
		if a.Key == "subject" {
			atom = a
		}
	}
	if atom.Key == "" {
		t.Fatal("no `subject` guard atom reached the kernel row")
	}

	// The contract's `in/member` case: literal `["alpha","beta"]`, a
	// present value `alpha`, verdict GuardTrue. The literal this package
	// emits must be byte-identical to the one the contract drives, or
	// RDR 0003's evaluator cannot decide the atom and must answer
	// GuardUnevaluable — which the unevaluable-candidate veto turns into a
	// deadlock rather than a refusal.
	const contractLiteral = `["alpha","beta"]`
	if atom.Literal != contractLiteral {
		t.Errorf("GuardAtom.Literal = %q; the shipped contract "+
			"(`internal/resolve/guardcontract.go`, case \"in/member\") "+
			"drives %q. A literal this evaluator cannot parse is "+
			"GuardUnevaluable, never false, so a mismatch here deadlocks "+
			"the flow behind RDR 0007's veto", atom.Literal, contractLiteral)
	}
}

// ---------------------------------------------------------------- FAIL-2

// FAIL-2 — a nested one-element array literal was silently flattened.
//
// REQ-6 (`0002:C3`): "**Strict decoding is an obligation on this format,
// not a property of a library.** The decoder MUST reject unmapped keys so
// an unknown schema field is a stable refusal rather than a silent no-op."
// The governing principle — a malformed authoring is a stable refusal,
// never a silent reinterpretation — was defeated at the VALUE layer:
// `load.go::valueMembers` recursed on `[]any` and rejected an element only
// when it yielded more than one member, so `[["a"], "b"]` was unwrapped
// into `["a", "b"]` while `[["a","b"], "c"]` refused.
//
// The refusal therefore depended on the inner array's ARITY rather than on
// its SHAPE, which is what makes it a defect rather than a leniency: the
// same authoring error refused or succeeded depending on how many elements
// the nested array happened to hold.
func TestFail2_NestedArrayLiteralsAreRefusedOnShapeNotArity(t *testing.T) {
	build := func(matchIn, guard, write, initial string) string {
		return `outcomes = ["done"]

[model]
id = "m"
version = 1
` + initial + `
[tags.status]
provenance = "owned"
kind = "scalar"

[tags.labels]
provenance = "owned"
kind = "set"

[tags.recognized]
provenance = "recognized"
kind = "enum"
required = true

[read.r]
role = "m"
path = "m"
keys = ["status", "labels"]
timeout = "2s"

[write.w]
role = "m"
path = "m"
keys = ["status", "labels"]
timeout = "2s"
read_back = true

[[rule]]
id = "advance"
[rule.match.recognized]
eq = "done"
[rule.match.status]
in = ` + matchIn + `
` + guard + `
[rule.write]
labels = ` + write + `
`
	}

	// Each site is exercised at BOTH inner arities. The one-element case is
	// the regression; the multi-element case is the control that already
	// refused, and asserting both is what pins the refusal to shape.
	sites := []struct {
		name  string
		build func(nested string) string
		want  table.Category
	}{
		{
			name: "a match-block `in` member",
			build: func(nested string) string {
				return build(`[`+nested+`, "b"]`, "", `["x"]`, "")
			},
			want: table.CatMalformedPredicateAtom,
		},
		{
			name: "a guard `contains` member",
			build: func(nested string) string {
				return build(`["a", "b"]`,
					"[rule.guard.all.labels]\ncontains = ["+nested+", \"y\"]",
					`["x"]`, "")
			},
			want: table.CatMalformedPredicateAtom,
		},
		{
			name: "a write-block value",
			build: func(nested string) string {
				return build(`["a", "b"]`, "", `[`+nested+`, "y"]`, "")
			},
			want: table.CatMalformedTagDeclaration,
		},
		{
			name: "an [initial] assignment",
			build: func(nested string) string {
				return build(`["a", "b"]`, "", `["x"]`,
					"\n[initial]\nlabels = ["+nested+", \"y\"]\n")
			},
			want: table.CatMalformedInitialDeclaration,
		},
	}

	arities := map[string]string{
		// The regression: a one-element nested array satisfied the old
		// arity test and was unwrapped.
		"one-element nested array": `["x"]`,
		// The control: this shape already refused, and must keep doing so.
		"multi-element nested array": `["x", "z"]`,
		// Depth is not a defence either — the recursion flattened
		// arbitrarily deep one-element nesting.
		"deeply nested one-element array": `[[["x"]]]`,
	}

	for _, site := range sites {
		for arity, nested := range arities {
			t.Run(site.name+"/"+arity, func(t *testing.T) {
				_, err := table.Load([]byte(site.build(nested)), "fail2.toml")
				if err == nil {
					t.Fatalf("a nested array literal loaded clean. A nested "+
						"array is not a member sequence and the format admits "+
						"no such spelling, so `0002:C3`'s stable-refusal "+
						"obligation binds it — and keying the refusal on the "+
						"inner array's arity made the same authoring error "+
						"refuse or succeed by accident (nested: %s)", nested)
				}
				if cat, _ := table.CategoryOf(err); cat != site.want {
					t.Errorf("category = %q; want %q (got: %v)", cat, site.want, err)
				}
			})
		}
	}

	// The control that the fix is not over-broad: a FLAT array at each of
	// the same sites still loads.
	t.Run("a flat array at every site still loads", func(t *testing.T) {
		src := build(`["a", "b"]`,
			"[rule.guard.all.labels]\ncontains = [\"x\", \"y\"]",
			`["x", "y"]`,
			"\n[initial]\nlabels = [\"x\", \"y\"]\n")
		if _, err := table.Load([]byte(src), "fail2-control.toml"); err != nil {
			t.Errorf("a flat array literal refused: %v", err)
		}
	})
}

// ------------------------------------------------------------------- R2

// R2 — a guard literal's carriage form is keyed on the OPERATOR ALONE.
//
// FAIL-1 above established that a guard atom's literal is a set by its
// OPERATOR (`in`/`contains` → the §D13 canonical JSON array), while a tag
// VALUE is a set by its declared KIND. The fix ADDED the operator trigger
// but left the kind trigger standing beside it, so a guard literal on a
// `set`-kind tag still array-encoded whatever the operator was.
//
// That is a defect for every operator RDR 0003 fences as single-valued.
// `exists` is the sharpest case because the kernel owns it outright:
// `internal/resolve/guard.go::evaluateAtom` compares the literal VERBATIM
// against `resolve.LiteralTrue`/`LiteralFalse` and, per `0007:C3`, a
// literal matching neither is never decided from presence — it is
// `GuardUnevaluable` with reason `uncomparable`. An `exists` literal
// carried as `["true"]` therefore matches neither constant, and RDR 0007's
// unevaluable-candidate veto turns that undecidable row into a whole-
// outcome `guard_unevaluable` refusal that poisons its decidable siblings.
//
// The atom is legal by construction: `load.go::conform` gives `exists` its
// own arm that returns before `conformKind`/`conformDomain`, so `exists`
// on a `set`-kind tag is authorable by design and must resolve.
func TestR2_GuardLiteralCarriageIsKeyedOnTheOperatorAlone(t *testing.T) {
	// `labels` is a declared `set`. The guard is `exists`, whose literal is
	// a bare bool constant regardless of the tag's kind.
	build := func(literal string) string {
		return `outcomes = ["done"]

[model]
id = "m"
version = 1

[tags.labels]
provenance = "owned"
kind = "set"

[tags.recognized]
provenance = "recognized"
kind = "enum"
required = true

[read.r]
role = "m"
path = "m"
keys = ["labels"]
timeout = "2s"

[write.w]
role = "m"
path = "m"
keys = ["labels"]
timeout = "2s"
read_back = true

[[rule]]
id = "advance"
[rule.match.recognized]
eq = "done"
[rule.guard.all.labels]
exists = ` + literal + `
[rule.write]
labels = ["alpha"]
`
	}

	// The carriage form itself, at the seam.
	t.Run("the literal crosses as the bare constant", func(t *testing.T) {
		for _, literal := range []string{"true", "false"} {
			m, err := table.Load([]byte(build(literal)), "r2-carriage.toml")
			if err != nil {
				t.Fatalf("exists = %s refused: %v", literal, err)
			}
			guard := m.Rows[0].KernelRow().Guard
			if len(guard) != 1 {
				t.Fatalf("%d guard atoms; want 1: %+v", len(guard), guard)
			}
			if guard[0].Literal != literal {
				t.Errorf("GuardAtom.Literal = %q; want the bare %q. The "+
					"kernel compares an `exists` literal verbatim against "+
					"resolve.LiteralTrue/LiteralFalse (`0007:C3`), so an "+
					"array-encoded literal matches neither constant",
					guard[0].Literal, literal)
			}
		}
	})

	// And end to end: the row must reach a DECIDED verdict through the
	// real kernel, not the veto.
	t.Run("the row resolves rather than refusing guard_unevaluable", func(t *testing.T) {
		cases := []struct {
			literal string
			owned   map[string]string
			plan    bool
		}{
			// `exists = true` over a present `labels`: the guard is TRUE
			// and the row is the plan.
			{literal: "true", owned: map[string]string{"labels": `["alpha"]`}, plan: true},
			// `exists = false` over a present `labels`: the guard is FALSE
			// and the row is pruned, so the outcome is a clean `no_match`
			// — a decided verdict, never `guard_unevaluable`.
			{literal: "false", owned: map[string]string{"labels": `["alpha"]`}, plan: false},
		}
		for _, c := range cases {
			m, err := table.Load([]byte(build(c.literal)), "r2-resolve.toml")
			if err != nil {
				t.Fatalf("exists = %s refused: %v", c.literal, err)
			}
			res, err := resolve.Resolve(resolve.Input{
				Flow:       "f",
				Table:      m.KernelTable(),
				Owned:      ownedTags(c.owned),
				Recognized: "done",
				Guards:     fixtureGuards{},
			})
			if err != nil {
				t.Fatalf("exists = %s: Resolve errored: %v", c.literal, err)
			}
			if res.Refused() && res.Refusal.Kind == resolve.KindGuardUnevaluable {
				t.Errorf("exists = %s refused guard_unevaluable (%+v); the "+
					"atom is legal — `conform` gives `exists` its own arm "+
					"— and the kernel decides presence itself, so the only "+
					"way it is undecidable is a literal this seam mangled",
					c.literal, res.Refusal.Undecided)
				continue
			}
			if c.plan && res.Plan == nil {
				t.Errorf("exists = true over a present tag did not plan: %+v",
					res.Refusal)
			}
			if !c.plan && res.Plan != nil {
				t.Errorf("exists = false over a present tag planned %+v; the "+
					"guard is FALSE and the row must be pruned", res.Plan)
			}
		}
	})
}
