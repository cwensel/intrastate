package resolve_test

// Fixture builders for RDR 0007 (guard predicate totality over an
// incomplete evaluation view). Nothing here mocks the unit under test:
// the kernel is always the real resolve.Resolve. What IS stubbed is the
// RDR 0003 value-evaluation seam, which this RDR deliberately narrows to
// "compare a present value against a literal" — a collaborator, not the
// unit.
//
// Every fixture drives verdicts THROUGH atoms and a value stub, never by
// injecting row verdicts, which would leave the kernel combinator off the
// tested path (RDR 0007 Testing Strategy preamble; premortem P-20).

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/resolve"
)

// --- atom builders -------------------------------------------------------

// allAtom builds a value-comparing atom in the `all` block.
func allAtom(key, op, literal string) resolve.GuardAtom {
	return resolve.GuardAtom{Key: key, Operator: op, Literal: literal, Block: resolve.BlockAll}
}

// unlessAtom builds a value-comparing atom in the `unless` block.
func unlessAtom(key, op, literal string) resolve.GuardAtom {
	return resolve.GuardAtom{Key: key, Operator: op, Literal: literal, Block: resolve.BlockUnless}
}

// existsAll builds an existence atom in the `all` block. present==true
// spells LiteralTrue, present==false spells LiteralFalse.
func existsAll(key string, present bool) resolve.GuardAtom {
	return resolve.GuardAtom{
		Key:      key,
		Operator: resolve.OpExists,
		Literal:  existsLiteral(present),
		Block:    resolve.BlockAll,
	}
}

// existsUnless builds an existence atom in the `unless` block.
func existsUnless(key string, present bool) resolve.GuardAtom {
	return resolve.GuardAtom{
		Key:      key,
		Operator: resolve.OpExists,
		Literal:  existsLiteral(present),
		Block:    resolve.BlockUnless,
	}
}

func existsLiteral(present bool) string {
	if present {
		return resolve.LiteralTrue
	}
	return resolve.LiteralFalse
}

// --- value seam stubs ----------------------------------------------------

// valueSeam is the RDR 0003 stand-in: a per-atom value evaluator that
// never sees the view (RDR 0007 `0007:C1`). It is keyed on the pair the
// seam is actually given — the atom's identity and the present value —
// so a test can distinguish "the seam decided" from "the kernel decided".
//
// Anything the table does not name is reported GuardUnevaluable, which is
// the A18 leg: a present value the seam cannot compare.
type valueSeam struct {
	// decided maps an atom-plus-value key to the seam's verdict.
	decided map[seamCall]resolve.GuardResult
	// calls records every invocation, so a test can assert the kernel did
	// NOT consult the seam for an absent key or an existence atom
	// (`0007:C1` states both as MUST NOTs on the kernel's call site, so
	// the call itself is spec-named behaviour, not an implementation
	// detail).
	calls *[]seamCall
}

// seamCall is one (atom, present value) pair handed across the seam.
type seamCall struct {
	Key      string
	Operator string
	Literal  string
	Block    resolve.Block
	Value    string
}

func newValueSeam() *valueSeam {
	return &valueSeam{decided: map[seamCall]resolve.GuardResult{}, calls: &[]seamCall{}}
}

// decide programs a verdict for one atom over one present value.
func (s *valueSeam) decide(atom resolve.GuardAtom, value string, verdict resolve.GuardResult) *valueSeam {
	s.decided[callOf(atom, value)] = verdict
	return s
}

func (s *valueSeam) Evaluate(atom resolve.GuardAtom, value string) resolve.GuardResult {
	c := callOf(atom, value)
	*s.calls = append(*s.calls, c)
	if v, ok := s.decided[c]; ok {
		return v
	}
	return resolve.GuardUnevaluable
}

// seen reports the calls the kernel made across the seam.
func (s *valueSeam) seen() []seamCall { return *s.calls }

func callOf(atom resolve.GuardAtom, value string) seamCall {
	return seamCall{
		Key:      atom.Key,
		Operator: atom.Operator,
		Literal:  atom.Literal,
		Block:    atom.Block,
		Value:    value,
	}
}

// eqSeam decides `eq` by exact byte equality of the present value against
// the atom's literal, and answers unevaluable for anything else. It is the
// smallest honest stand-in for RDR 0003's typed `eq` and is what the
// present-key legs of this suite compare against.
type eqSeam struct{}

func (eqSeam) Evaluate(atom resolve.GuardAtom, value string) resolve.GuardResult {
	if atom.Operator != opEq {
		return resolve.GuardUnevaluable
	}
	if value == atom.Literal {
		return resolve.GuardTrue
	}
	return resolve.GuardFalse
}

// Operator tokens from RDR 0003's operator/kind matrix. Only OpExists is
// this RDR's; the rest are opaque tokens the kernel hands to the seam.
const (
	opEq       = "eq"
	opIn       = "in"
	opGte      = "gte"
	opContains = "contains"
)

// --- table builders ------------------------------------------------------

// guardedRow builds one ordinary candidate row over the legalInput view
// (owned status=Draft, observed reviews=2, recognized). Writes and
// RequiresOwned are the row's post-guard transition dependency, which
// `0007:C9` narrows to exactly that.
func guardedRow(ruleID, locator string, atoms ...resolve.GuardAtom) resolve.Row {
	return resolve.Row{
		RuleID:        ruleID,
		SourceLocator: locator,
		Outcome:       "successful",
		Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
		RequiresOwned: []string{"status"},
		Guard:         atoms,
		NextTags:      []resolve.Tag{{Key: "status", Value: "Final"}},
		Writes:        []resolve.Tag{{Key: "status", Value: "Final"}},
	}
}

// conformingEscapeRow builds an escape row that conforms to `0007:C9`'s
// A21 composition: empty Writes, and therefore empty RequiresOwned. The
// shipped fixtures_test.go::escapeRow is non-conforming on BOTH fields
// (Testing Strategy row 15) and is deliberately not reused here.
func conformingEscapeRow(ruleID, locator string, class resolve.RefusalKind, atoms ...resolve.GuardAtom) resolve.Row {
	return resolve.Row{
		RuleID:        ruleID,
		SourceLocator: locator,
		Outcome:       "successful",
		Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
		RequiresOwned: nil,
		Guard:         atoms,
		NextTags:      []resolve.Tag{{Key: "status", Value: "Blocked"}},
		Writes:        nil,
		Escape:        []resolve.RefusalKind{class},
	}
}

// guardInput assembles an input over the canonical fixture view whose
// table holds exactly the given rows.
func guardInput(revision string, seam resolve.GuardEvaluator, rows ...resolve.Row) resolve.Input {
	return resolve.Input{
		Flow: "rdr",
		Table: resolve.Table{
			Revision: revision,
			Outcomes: []string{"successful", "failed"},
			Rows:     rows,
		},
		Owned:      []resolve.Tag{{Key: "status", Value: "Draft"}},
		Observed:   []resolve.Tag{{Key: "reviews", Value: "2"}},
		Recognized: "successful",
		Guards:     seam,
	}
}

// absentKey is a key no fixture view carries. `legalInput` supplies
// status (owned), reviews (observed), and recognized (recognized).
const absentKey = "iterations"

// --- payload helpers -----------------------------------------------------

// undecidedAtomsFor returns the payload atoms reported for one rule id.
func undecidedAtomsFor(r *resolve.Refusal, ruleID string) []resolve.UndecidedAtom {
	for _, row := range r.Undecided {
		if row.RuleID == ruleID {
			return row.Atoms
		}
	}
	return nil
}

// undecidedRuleIDs returns the rule ids named by the payload, in payload
// order (which `0007:C8` requires to be sorted, so order is meaningful).
func undecidedRuleIDs(r *resolve.Refusal) []string {
	out := make([]string, 0, len(r.Undecided))
	for _, row := range r.Undecided {
		out = append(out, row.RuleID)
	}
	return out
}

// reasonFor returns the reason reported for the atom on key in block, and
// whether the payload carried such an entry at all.
func reasonFor(atoms []resolve.UndecidedAtom, key string, block resolve.Block) (resolve.Reason, bool) {
	for _, a := range atoms {
		if a.Key == key && a.Block == block {
			return a.Reason, true
		}
	}
	return "", false
}

// mustRefuse resolves and fails unless the kernel refused with kind.
func mustRefuse(t *testing.T, in resolve.Input, kind resolve.RefusalKind) *resolve.Refusal {
	t.Helper()
	got := mustResolve(t, in)
	if !got.Refused() {
		t.Fatalf("kernel emitted a plan for rule %q (escaped=%v, writes=%v); want refusal %q",
			got.Plan.RuleID, got.Plan.Escaped, got.Plan.Writes, kind)
	}
	if got.Refusal.Kind != kind {
		t.Fatalf("refusal kind = %q; want %q", got.Refusal.Kind, kind)
	}
	return got.Refusal
}

// mustPlan resolves and fails unless the kernel produced a plan.
func mustPlan(t *testing.T, in resolve.Input) *resolve.Plan {
	t.Helper()
	got := mustResolve(t, in)
	if got.Refused() {
		t.Fatalf("kernel refused %q (undecided=%v, missingOwned=%v); want a plan",
			got.Refusal.Kind, got.Refusal.Undecided, got.Refusal.MissingOwned)
	}
	return got.Plan
}

// --- §D13 set encoding ---------------------------------------------------

// d13Set renders a set as the canonical carriage form JDR 0001 §D13
// fixes: "a set crosses as its canonical JSON array — members sorted,
// duplicate-free, compact encoding". RDR 0002 declares the clause; per
// docs/rdr/BUILD-ORDER.md ("The one seam this order strains") 0007 writes
// its `contains` leg against §D13 directly.
func d13Set(members ...string) string {
	sorted := slices.Clone(members)
	slices.Sort(sorted)
	sorted = slices.Compact(sorted)
	b, err := json.Marshal(sorted)
	if err != nil {
		panic(err)
	}
	return string(b)
}
