// A1 evaluation harness — RDR 0003 Critical Assumption A1.
//
// Claim under test: the target RDR/kata flows fit a CLOSED TYPED predicate
// vocabulary (eq, in, lt, lte, gt, gte, exists, contains).
//
// Method: encode the Resolve-spike fixture's guard rows as normalized atoms,
// EVALUATE them against representative tag inputs, and assert each row's
// expected qualify/prune verdict. This is what the prior check.sh could not
// do (it scanned operator spellings only).
//
// Scope note: declared domains (RDR 0002 A11/A7/A9) are NOT consumed here.
// Evaluation needs values; domains are needed only for exhaustiveness (A2).
// The fixture's domain/min/max keys are deliberately ignored to avoid the
// circularity A2's Evidence line warns about.
package main

import (
	"fmt"
	"os"
	"sort"
)

// ---- Normalized atom shape (RDR 0003 / JDR 0001 §D1: key, operator, literal, block).

type Block int

const (
	BlockAll Block = iota
	BlockUnless
)

func (b Block) String() string {
	if b == BlockUnless {
		return "unless"
	}
	return "all"
}

type Op string

const (
	OpEq       Op = "eq"
	OpIn       Op = "in"
	OpLt       Op = "lt"
	OpLte      Op = "lte"
	OpGt       Op = "gt"
	OpGte      Op = "gte"
	OpExists   Op = "exists"
	OpContains Op = "contains"
)

// Kind is the declared tag value kind.
type Kind string

const (
	KindEnum Kind = "enum"
	KindInt  Kind = "int"
	KindBool Kind = "bool"
	KindSet  Kind = "set"
)

// Value is a typed tag value.
type Value struct {
	Kind Kind
	S    string
	I    int
	B    bool
	Set  []string
}

// Literal is an atom's right-hand side: a typed scalar or a typed set.
type Literal struct {
	Kind   Kind
	S      string
	I      int
	B      bool
	Set    []string // for `in` and `contains`
	IsSet  bool
}

type Atom struct {
	Key   string
	Op    Op
	Lit   Literal
	Block Block
}

type Row struct {
	RuleID string
	Match  map[string]Literal
	Atoms  []Atom
}

// ---- Guard verdict, mirroring internal/resolve.GuardResult.

type Verdict int

const (
	False Verdict = iota
	True
	Unevaluable
)

func (v Verdict) String() string {
	switch v {
	case True:
		return "TRUE"
	case False:
		return "FALSE"
	default:
		return "UNEVALUABLE"
	}
}

// ---- Operator/kind matrix — the closed typed vocabulary A1 is about.

var allowed = map[Op]map[Kind]bool{
	OpEq:       {KindEnum: true, KindInt: true, KindBool: true},
	OpIn:       {KindEnum: true, KindInt: true},
	OpLt:       {KindInt: true},
	OpLte:      {KindInt: true},
	OpGt:       {KindInt: true},
	OpGte:      {KindInt: true},
	OpExists:   {KindEnum: true, KindInt: true, KindBool: true, KindSet: true},
	OpContains: {KindSet: true},
}

type wellFormedErr struct{ msg string }

func (e wellFormedErr) Error() string { return e.msg }

// checkWellFormed enforces the operator/kind matrix against the DECLARED kind
// of the tag, before any evaluation. Unknown operator or unsupported pairing
// is rejected — the "closed and typed" half of A1.
func checkWellFormed(a Atom, declKind Kind) error {
	kinds, ok := allowed[a.Op]
	if !ok {
		return wellFormedErr{fmt.Sprintf("unknown operator %q on key %q", a.Op, a.Key)}
	}
	if !kinds[declKind] {
		return wellFormedErr{fmt.Sprintf("operator %q not allowed on kind %q (key %q)", a.Op, declKind, a.Key)}
	}
	// exists takes a bool literal regardless of the tag's kind.
	if a.Op == OpExists {
		if a.Lit.Kind != KindBool {
			return wellFormedErr{fmt.Sprintf("exists literal must be bool on key %q", a.Key)}
		}
		return nil
	}
	if a.Lit.IsSet != (a.Op == OpIn || a.Op == OpContains) {
		return wellFormedErr{fmt.Sprintf("operator %q literal set-ness mismatch on key %q", a.Op, a.Key)}
	}
	if !a.Lit.IsSet && a.Lit.Kind != declKind {
		return wellFormedErr{fmt.Sprintf("literal kind %q does not match declared kind %q (key %q)", a.Lit.Kind, declKind, a.Key)}
	}
	return nil
}

// evalAtom decides VALUE semantics over a present value only (the RDR's seam
// scope). Presence and `exists` belong to the caller/kernel.
func evalAtom(a Atom, v Value) Verdict {
	switch a.Op {
	case OpEq:
		switch v.Kind {
		case KindEnum:
			return b2v(v.S == a.Lit.S)
		case KindInt:
			return b2v(v.I == a.Lit.I)
		case KindBool:
			return b2v(v.B == a.Lit.B)
		}
	case OpIn:
		switch v.Kind {
		case KindEnum:
			for _, e := range a.Lit.Set {
				if v.S == e {
					return True
				}
			}
			return False
		}
	case OpLt:
		return b2v(v.I < a.Lit.I)
	case OpLte:
		return b2v(v.I <= a.Lit.I)
	case OpGt:
		return b2v(v.I > a.Lit.I)
	case OpGte:
		return b2v(v.I >= a.Lit.I)
	case OpContains:
		// set-valued tag contains every element of the literal set
		have := map[string]bool{}
		for _, e := range v.Set {
			have[e] = true
		}
		for _, want := range a.Lit.Set {
			if !have[want] {
				return False
			}
		}
		return True
	}
	return Unevaluable
}

func b2v(b bool) Verdict {
	if b {
		return True
	}
	return False
}

// evalRow combines atoms: `all` is conjunctive (strong Kleene), `unless` is a
// conjunctive exclusion block subtracted as one unit. An absent key under a
// value atom is UNEVALUABLE, not false (RDR 0007's veto).
func evalRow(r Row, view map[string]Value, decl map[string]Kind) (Verdict, error) {
	allV := True
	unlessAtoms := 0
	unlessV := True

	for _, a := range r.Atoms {
		k, ok := decl[a.Key]
		if !ok {
			return Unevaluable, wellFormedErr{fmt.Sprintf("undeclared tag %q", a.Key)}
		}
		if err := checkWellFormed(a, k); err != nil {
			return Unevaluable, err
		}
		v, present := view[a.Key]

		var got Verdict
		if a.Op == OpExists {
			// TOTAL operator: decides on presence alone, never unevaluable.
			got = b2v(present == a.Lit.B)
		} else if !present {
			got = Unevaluable
		} else {
			got = evalAtom(a, v)
		}

		if a.Block == BlockUnless {
			unlessAtoms++
			unlessV = kleeneAnd(unlessV, got)
		} else {
			allV = kleeneAnd(allV, got)
		}
	}

	if allV == False {
		return False, nil
	}
	if unlessAtoms > 0 {
		// The row is disabled when the full `unless` block matches.
		if unlessV == True {
			return False, nil
		}
		if unlessV == Unevaluable && allV != False {
			return Unevaluable, nil
		}
	}
	return allV, nil
}

func kleeneAnd(a, b Verdict) Verdict {
	if a == False || b == False {
		return False
	}
	if a == Unevaluable || b == Unevaluable {
		return Unevaluable
	}
	return True
}

// matches applies the row's `match` pattern (source-state selection).
func matches(r Row, view map[string]Value) bool {
	for k, lit := range r.Match {
		v, ok := view[k]
		if !ok || v.S != lit.S {
			return false
		}
	}
	return true
}

// ---- The fixture, encoded as normalized rows.

func decls() map[string]Kind {
	return map[string]Kind{
		"status":             KindEnum,
		"profile":            KindEnum,
		"stage":              KindEnum,
		"lens":               KindEnum,
		"prelock_iterations": KindInt,
		"cluster_eligible":   KindBool,
		"rewind_target":      KindEnum,
		// Added for the `contains` cell — a set-valued tag. NOTE: RDR 0002
		// cannot declare an element universe today (A9); this harness needs
		// only the KIND to evaluate, not the universe, so the cell is
		// exercisable for A1 (vocabulary) while A9 remains open for A2.
		"lenses_done": KindSet,
	}
}

func fixtureRows() []Row {
	s := func(v string) Literal { return Literal{Kind: KindEnum, S: v} }
	i := func(v int) Literal { return Literal{Kind: KindInt, I: v} }
	bl := func(v bool) Literal { return Literal{Kind: KindBool, B: v} }
	set := func(vs ...string) Literal { return Literal{Kind: KindEnum, Set: vs, IsSet: true} }

	return []Row{
		{
			RuleID: "profile-to-grounding",
			Match:  map[string]Literal{"status": s("Draft")},
			Atoms: []Atom{
				{Key: "profile", Op: OpIn, Lit: set("mid", "large"), Block: BlockAll},
				{Key: "prelock_iterations", Op: OpGte, Lit: i(3), Block: BlockUnless},
			},
		},
		{
			RuleID: "foundational-to-cove",
			Match:  map[string]Literal{"status": s("Draft")},
			Atoms: []Atom{
				{Key: "profile", Op: OpEq, Lit: s("foundational"), Block: BlockAll},
				{Key: "cluster_eligible", Op: OpExists, Lit: bl(true), Block: BlockAll},
			},
		},
		{
			RuleID: "continue-prelock-lenses",
			Match:  map[string]Literal{"stage": s("prelock")},
			Atoms: []Atom{
				{Key: "lens", Op: OpIn, Lit: set("grounding", "3amigo", "critique", "repeatability", "cove"), Block: BlockAll},
				{Key: "prelock_iterations", Op: OpLt, Lit: i(3), Block: BlockAll},
			},
		},
		{
			RuleID: "reconcile-rewind-legality",
			Match:  map[string]Literal{"stage": s("reconcile")},
			Atoms: []Atom{
				{Key: "rewind_target", Op: OpIn, Lit: set("resolve", "refine", "propose"), Block: BlockAll},
				{Key: "status", Op: OpEq, Lit: s("Implemented"), Block: BlockUnless},
			},
		},
		// `contains` cell — the operator Phase 3 gates. Encoded over a
		// set-valued tag so the vocabulary claim covers all 8 operators.
		{
			RuleID: "lens-row-complete",
			Match:  map[string]Literal{"stage": s("prelock")},
			Atoms: []Atom{
				{Key: "lenses_done", Op: OpContains, Lit: set("grounding", "3amigo", "critique"), Block: BlockAll},
			},
		},
		// `lte` / `gt` cell — the two operators check.sh's hardcoded report
		// loop could never print. Exercised so the closed set is fully covered.
		{
			RuleID: "iteration-band",
			Match:  map[string]Literal{"stage": s("prelock")},
			Atoms: []Atom{
				{Key: "prelock_iterations", Op: OpGt, Lit: i(0), Block: BlockAll},
				{Key: "prelock_iterations", Op: OpLte, Lit: i(2), Block: BlockAll},
			},
		},
	}
}

// ---- Test cases: representative tag inputs with EXPECTED verdicts.

type tc struct {
	name   string
	view   map[string]Value
	expect map[string]Verdict // ruleID -> expected verdict for MATCHING rows
	skip   map[string]bool    // rows whose match pattern excludes them
}

func enum(v string) Value { return Value{Kind: KindEnum, S: v} }
func inum(v int) Value    { return Value{Kind: KindInt, I: v} }
func bval(v bool) Value   { return Value{Kind: KindBool, B: v} }
func sval(vs ...string) Value {
	return Value{Kind: KindSet, Set: vs}
}

func cases() []tc {
	return []tc{
		{
			name: "large draft, first prelock iteration -> grounding qualifies",
			view: map[string]Value{
				"status": enum("Draft"), "profile": enum("large"),
				"prelock_iterations": inum(0),
			},
			expect: map[string]Verdict{
				"profile-to-grounding": True,
				"foundational-to-cove": False, // profile != foundational
			},
		},
		{
			name: "large draft at iteration cap -> unless block prunes the row",
			view: map[string]Value{
				"status": enum("Draft"), "profile": enum("large"),
				"prelock_iterations": inum(3),
			},
			expect: map[string]Verdict{
				"profile-to-grounding": False,
				"foundational-to-cove": False,
			},
		},
		{
			name: "foundational draft with cluster_eligible present -> cove qualifies",
			view: map[string]Value{
				"status": enum("Draft"), "profile": enum("foundational"),
				"cluster_eligible":   bval(true),
				"prelock_iterations": inum(1),
			},
			expect: map[string]Verdict{
				"profile-to-grounding": False, // foundational not in {mid,large}
				"foundational-to-cove": True,
			},
		},
		{
			name: "foundational draft, cluster_eligible ABSENT -> exists decides FALSE (never unevaluable)",
			view: map[string]Value{
				"status": enum("Draft"), "profile": enum("foundational"),
				"prelock_iterations": inum(1),
			},
			expect: map[string]Verdict{
				"profile-to-grounding": False,
				"foundational-to-cove": False, // exists=true on an absent key -> FALSE, decided
			},
		},
		{
			name: "prelock, lens set membership + bounded int -> qualifies",
			view: map[string]Value{
				"stage": enum("prelock"), "lens": enum("critique"),
				"prelock_iterations": inum(2),
				"lenses_done":        sval("grounding", "3amigo", "critique", "repeatability"),
			},
			expect: map[string]Verdict{
				"continue-prelock-lenses": True,
				"lens-row-complete":       True, // contains: superset satisfies
				"iteration-band":          True, // gt 0 AND lte 2
			},
		},
		{
			name: "prelock at iteration cap -> lt prunes; gt/lte band prunes",
			view: map[string]Value{
				"stage": enum("prelock"), "lens": enum("cove"),
				"prelock_iterations": inum(3),
				"lenses_done":        sval("grounding", "3amigo"),
			},
			expect: map[string]Verdict{
				"continue-prelock-lenses": False, // lt 3 fails
				"lens-row-complete":       False, // contains: critique missing
				"iteration-band":          False, // lte 2 fails
			},
		},
		{
			name: "prelock, iteration 0 -> gt 0 prunes the band row",
			view: map[string]Value{
				"stage": enum("prelock"), "lens": enum("grounding"),
				"prelock_iterations": inum(0),
				"lenses_done":        sval("grounding", "3amigo", "critique"),
			},
			expect: map[string]Verdict{
				"continue-prelock-lenses": True,
				"lens-row-complete":       True, // exact set satisfies contains
				"iteration-band":          False,
			},
		},
		{
			name: "reconcile with legal rewind target -> qualifies",
			view: map[string]Value{
				"stage": enum("reconcile"), "status": enum("Draft"),
				"rewind_target": enum("refine"),
			},
			expect: map[string]Verdict{
				"reconcile-rewind-legality": True,
			},
		},
		{
			name: "reconcile on an Implemented RDR -> unless prunes rewind",
			view: map[string]Value{
				"stage": enum("reconcile"), "status": enum("Implemented"),
				"rewind_target": enum("refine"),
			},
			expect: map[string]Verdict{
				"reconcile-rewind-legality": False,
			},
		},
		{
			name: "reconcile with illegal rewind target -> all block prunes",
			view: map[string]Value{
				"stage": enum("reconcile"), "status": enum("Draft"),
				"rewind_target": enum("finalize"),
			},
			expect: map[string]Verdict{
				"reconcile-rewind-legality": False,
			},
		},
		{
			name: "prelock, lens key ABSENT -> value atom is UNEVALUABLE, not false (0007 veto)",
			view: map[string]Value{
				"stage": enum("prelock"), "prelock_iterations": inum(1),
				"lenses_done": sval("grounding"),
			},
			expect: map[string]Verdict{
				"continue-prelock-lenses": Unevaluable,
				"lens-row-complete":       False,
				"iteration-band":          True,
			},
		},
	}
}

// ---- Negative controls: well-formedness rejections.

type badCase struct {
	name string
	atom Atom
	key  Kind
}

func badCases() []badCase {
	return []badCase{
		{"unknown operator", Atom{Key: "profile", Op: Op("matches"), Lit: Literal{Kind: KindEnum, S: "x"}}, KindEnum},
		{"lt on enum", Atom{Key: "profile", Op: OpLt, Lit: Literal{Kind: KindEnum, S: "mid"}}, KindEnum},
		{"contains on enum", Atom{Key: "profile", Op: OpContains, Lit: Literal{Kind: KindEnum, Set: []string{"a"}, IsSet: true}}, KindEnum},
		{"literal kind mismatch", Atom{Key: "prelock_iterations", Op: OpEq, Lit: Literal{Kind: KindEnum, S: "three"}}, KindInt},
		{"in with scalar literal", Atom{Key: "profile", Op: OpIn, Lit: Literal{Kind: KindEnum, S: "mid"}}, KindEnum},
	}
}

func main() {
	decl := decls()
	rows := fixtureRows()
	byID := map[string]Row{}
	for _, r := range rows {
		byID[r.RuleID] = r
	}

	fmt.Println("A1 EVALUATION HARNESS — RDR 0003")
	fmt.Println("claim: target-flow guard rows are expressible AND decidable in the closed typed vocabulary")
	fmt.Printf("fixture: evidence/spikes/guard-fixture.toml (4 rules) + 2 rows added for contains/lte/gt\n")
	fmt.Println()

	pass, fail := 0, 0

	fmt.Println("== Evaluation cases ==")
	for _, c := range cases() {
		ids := make([]string, 0, len(c.expect))
		for id := range c.expect {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		for _, id := range ids {
			r := byID[id]
			if !matches(r, c.view) {
				fmt.Printf("  FAIL  %s / %s: match pattern did not select the row\n", c.name, id)
				fail++
				continue
			}
			got, err := evalRow(r, c.view, decl)
			if err != nil {
				fmt.Printf("  FAIL  %s / %s: unexpected well-formedness error: %v\n", c.name, id, err)
				fail++
				continue
			}
			want := c.expect[id]
			if got != want {
				fmt.Printf("  FAIL  %s / %s: want %s, got %s\n", c.name, id, want, got)
				fail++
				continue
			}
			fmt.Printf("  ok    %-28s %-26s -> %s\n", id, trunc(c.name), got)
			pass++
		}
	}

	fmt.Println()
	fmt.Println("== Well-formedness negative controls (closed+typed) ==")
	for _, b := range badCases() {
		err := checkWellFormed(b.atom, b.key)
		if err == nil {
			fmt.Printf("  FAIL  %s: expected rejection, got accept\n", b.name)
			fail++
			continue
		}
		fmt.Printf("  ok    rejected: %-24s (%v)\n", b.name, err)
		pass++
	}

	fmt.Println()
	fmt.Println("== Operator coverage (measured, not asserted) ==")
	seen := map[Op]int{}
	for _, r := range rows {
		for _, a := range r.Atoms {
			seen[a.Op]++
		}
	}
	closed := []Op{OpEq, OpIn, OpLt, OpLte, OpGt, OpGte, OpExists, OpContains}
	uncovered := []string{}
	for _, op := range closed {
		n := seen[op]
		mark := "covered"
		if n == 0 {
			mark = "UNCOVERED"
			uncovered = append(uncovered, string(op))
		}
		fmt.Printf("  %-9s atoms=%d  %s\n", op, n, mark)
	}

	fmt.Println()
	fmt.Println("== Operators REQUIRED but absent from the closed set ==")
	fmt.Println("  none — every target-flow guard row above was encoded without")
	fmt.Println("  an operator outside {eq,in,lt,lte,gt,gte,exists,contains}.")

	fmt.Println()
	fmt.Printf("RESULT: pass=%d fail=%d uncovered_operators=%v\n", pass, fail, uncovered)
	if fail > 0 || len(uncovered) > 0 {
		fmt.Println("VERDICT: A1 NOT ESTABLISHED")
		os.Exit(1)
	}
	fmt.Println("VERDICT: A1 SUPPORTED — vocabulary sufficient and decidable over the target-flow slice")
}

func trunc(s string) string {
	if len(s) > 26 {
		return s[:23] + "..."
	}
	return s
}
