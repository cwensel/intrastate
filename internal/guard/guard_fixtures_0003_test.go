package guard_test

// Shared fixture helpers for the RDR 0003 guard/lint suite.
//
// Nothing here mocks the unit under test. Models are built by the real
// `internal/table` loader from authored TOML, and every lint assertion
// drives the real `internal/guard`. The loader is a collaborator that RDR
// 0002 owns; this RDR owns what the declarations MEAN.

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// --- model construction --------------------------------------------------

// modelHeader is the minimum JDR 0001 §D7 preamble every fixture needs: a
// model id, the recognized-outcome alphabet, and the reserved
// `[tags.recognized]` declaration. Tests supply their own tag declarations
// and rules on top of it.
const modelHeader = `
outcomes = ["go", "stop"]

[model]
id = "t"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true
`

// declBlock prefixes tag declarations with the model header and appends
// the accessor bindings RDR 0002 requires of them: every owned tag served
// by exactly one reader and one writer.
//
// The bindings are DERIVED from the declarations rather than authored per
// fixture, so an accessor-arity slip can never masquerade as an RDR 0003
// failure — every fixture below either loads clean or fails on a clause
// this RDR owns.
func declBlock(decls string) string {
	return modelHeader + decls + accessorsFor(decls)
}

// accessorsFor emits the reader/writer accessors covering every owned tag
// the declarations name.
func accessorsFor(decls string) string {
	var owned []string
	var key string
	for line := range strings.Lines(decls) {
		line = strings.TrimSpace(line)
		if name, ok := strings.CutPrefix(line, "[tags."); ok {
			key = strings.TrimSuffix(name, "]")
			continue
		}
		if line == `provenance = "owned"` && key != "" {
			owned = append(owned, key)
			key = ""
		}
	}
	if len(owned) == 0 {
		return ""
	}
	slices.Sort(owned)
	quoted := make([]string, 0, len(owned))
	for _, k := range owned {
		quoted = append(quoted, strconv.Quote(k))
	}
	keys := "[" + strings.Join(quoted, ", ") + "]"
	return `
[read.own]
role = "t"
path = "t.own"
keys = ` + keys + `
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ` + keys + `
timeout = "2s"
read_back = true
`
}

// mustLoadSource loads authored TOML that the spec says must load clean.
func mustLoadSource(t *testing.T, src string) *table.Model {
	t.Helper()

	m, err := table.Load([]byte(src), "fixture.toml")
	if err != nil {
		t.Fatalf("fixture must load clean; refused: %v\n--- source ---\n%s", err, src)
	}
	if m == nil {
		t.Fatal("loader returned a nil model with no error")
	}
	return m
}

// loadCategory loads decls+rule expected to refuse and returns the stable
// data-level category. The oracle is the category, never message text.
//
// The rule carries a write block because RDR 0002 checks rule SHAPE before
// it parses any atom: without one, every case here would refuse as
// `malformed rule shape` and the atom defect under test would never be
// reached. The write block is empty, so it constrains nothing this RDR
// owns — it only keeps the rule well-formed enough for its guard to be
// the defect the load reports.
func loadCategory(t *testing.T, decls, guardBlock string) table.Category {
	t.Helper()

	src := declBlock(decls) + `
[[rule]]
id = "r"
[rule.match.recognized]
eq = "go"
` + guardBlock + `[rule.write]
`

	_, err := table.Load([]byte(src), "fixture.toml")
	if err == nil {
		t.Fatalf("fixture loaded clean; want a refusal\n--- source ---\n%s", src)
	}
	cat, ok := table.CategoryOf(err)
	if !ok {
		t.Fatalf("fixture refused with a non-categorized error: %v", err)
	}
	return cat
}

// loadDeclCategory loads a tag-declaration-only fixture expected to refuse.
func loadDeclCategory(t *testing.T, decls string) table.Category {
	t.Helper()

	_, err := table.Load([]byte(declBlock(decls)), "fixture.toml")
	if err == nil {
		t.Fatalf("declaration loaded clean; want a refusal\n--- source ---\n%s", decls)
	}
	cat, ok := table.CategoryOf(err)
	if !ok {
		t.Fatalf("declaration refused with a non-categorized error: %v", err)
	}
	return cat
}

// --- lint result inspection ----------------------------------------------

// allFindings flattens every finding across every group report, for
// failure messages that must show what lint DID emit.
func allFindings(reports []guard.GroupReport) []guard.Finding {
	var out []guard.Finding
	for _, r := range reports {
		out = append(out, r.Findings...)
	}
	return out
}

// anyFinding reports whether some finding satisfies pred.
func anyFinding(reports []guard.GroupReport, pred func(guard.Finding) bool) bool {
	for _, f := range allFindings(reports) {
		if pred(f) {
			return true
		}
	}
	return false
}

// findingsWithCode returns every finding carrying code.
func findingsWithCode(reports []guard.GroupReport, code guard.Code) []guard.Finding {
	var out []guard.Finding
	for _, f := range allFindings(reports) {
		if f.Code == code {
			out = append(out, f)
		}
	}
	return out
}

// countCode counts findings carrying code.
func countCode(reports []guard.GroupReport, code guard.Code) int {
	return len(findingsWithCode(reports, code))
}

// greenGroups returns the reports whose verdict is a certified-exhaustive
// green — the only verdict this RDR lets stand as a proof.
func greenGroups(reports []guard.GroupReport) []guard.GroupReport {
	var out []guard.GroupReport
	for _, r := range reports {
		if r.Green {
			out = append(out, r)
		}
	}
	return out
}

// hasGreen reports whether any group certified green.
func hasGreen(reports []guard.GroupReport) bool {
	return len(greenGroups(reports)) > 0
}

// reportFor returns the single group report whose context names outcome
// and carries ruleID, so a multi-group fixture can be asserted per group.
func reportFor(t *testing.T, reports []guard.GroupReport, ruleID string) guard.GroupReport {
	t.Helper()

	var found []guard.GroupReport
	for _, r := range reports {
		if slices.Contains(r.RuleIDs, ruleID) {
			found = append(found, r)
		}
	}
	switch len(found) {
	case 1:
		return found[0]
	case 0:
		t.Fatalf("no group report carries rule %q; reports=%s", ruleID, renderReports(reports))
	default:
		t.Fatalf("rule %q appears in %d group reports; a row sits in exactly "+
			"one selection context", ruleID, len(found))
	}
	return guard.GroupReport{}
}

// renderReports formats reports for a failure message.
func renderReports(reports []guard.GroupReport) string {
	var b strings.Builder
	for _, r := range reports {
		fmt.Fprintf(&b, "\n  group %s verdict=%s green=%v rules=%v findings=%v",
			r.Context, r.Verdict, r.Green, r.RuleIDs, r.Findings)
	}
	if b.Len() == 0 {
		return " (none)"
	}
	return b.String()
}

// --- row and resolution helpers ------------------------------------------

// rowByID returns the single normalized row carrying ruleID.
func rowByID(t *testing.T, m *table.Model, ruleID string) table.Row {
	t.Helper()

	for _, r := range m.Rows {
		if r.RuleID == ruleID {
			return r
		}
	}
	t.Fatalf("model carries no row with rule id %q", ruleID)
	return table.Row{}
}

// guardRow builds a bare normalized row carrying atoms, for the syntactic
// declaration-only predicates (CanRefuse) that read no model.
func guardRow(ruleID string, atoms ...table.Atom) table.Row {
	return table.Row{RuleID: ruleID, SourceLocator: "f:1", Atoms: atoms}
}

// resolveWith runs the real kernel over kt with view supplied under both
// tag provenances. It is how a lint claim is checked against the runtime it
// describes.
//
// One evaluation view is exactly what this RDR's `View` is: it carries no
// provenance, because conformance and the narrowing read the DECLARATION,
// never a runtime trace. The kernel, though, gates on provenance in two
// separate places — a row's `RequiresOwned` needs an OWNED-provenance
// entry, while a guard atom reads whatever the merged view holds — so
// supplying the view under one provenance alone would refuse
// `owned_state_unavailable` before any guard was reached, testing the
// accessor gate rather than the predicate semantics under test. Supplying
// both is value-identical: `assemble` resolves a key crossing provenances
// by precedence, never as a conflict.
func resolveWith(t *testing.T, kt resolve.Table, view guard.View) resolve.Result {
	t.Helper()

	tags := make([]resolve.Tag, 0, len(view))
	for _, key := range slices.Sorted(maps.Keys(view)) {
		tags = append(tags, resolve.Tag{Key: key, Value: view[key]})
	}
	res, err := resolve.Resolve(resolve.Input{
		Table:      kt,
		Recognized: recognizedFor(kt),
		Owned:      tags,
		Observed:   tags,
		Guards:     guard.Evaluator{},
	})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	return res
}

// recognizedFor picks the recognized outcome to resolve under: the
// fixtures' `go`, or — for a slice declaring its own alphabet, as the MVV's
// two do — that alphabet's first member. An outcome outside the table's
// declared alphabet refuses `unmodeled_outcome` before any row is
// considered, which would test RDR 0001's alphabet gate rather than this
// RDR's predicate semantics.
func recognizedFor(kt resolve.Table) string {
	if slices.Contains(kt.Outcomes, "go") {
		return "go"
	}
	if len(kt.Outcomes) == 0 {
		return "go"
	}
	return kt.Outcomes[0]
}

// describe renders a disposition for a failure message.
func describe(r resolve.Result) string {
	if r.Refused() {
		return "refusal:" + string(r.Refusal.Kind)
	}
	if r.Plan != nil {
		return "plan:" + r.Plan.RuleID
	}
	return "neither plan nor refusal"
}

// sameDisposition reports whether two results agree on plan-vs-refusal and
// on the identity the disposition names.
func sameDisposition(a, b resolve.Result) bool {
	return describe(a) == describe(b)
}

// verdictOf renders the set of group verdicts, so two lint runs can be
// compared for identity without depending on report order.
func verdictOf(reports []guard.GroupReport) string {
	var out []string
	for _, r := range reports {
		out = append(out, string(r.Verdict))
	}
	slices.Sort(out)
	return strings.Join(out, ",")
}

// --- kernel-atom helpers -------------------------------------------------

// kernelAtom builds a value-comparing kernel atom in the `all` block, the
// shape the evaluator seam is handed.
func kernelAtom(key, op, literal string) resolve.GuardAtom {
	return resolve.GuardAtom{Key: key, Operator: op, Literal: literal, Block: resolve.BlockAll}
}

// resolveTrue names the decided-true verdict without repeating the
// package-qualified constant in every assertion.
func resolveTrue() resolve.GuardResult { return resolve.GuardTrue }

// --- shared fixture sources ----------------------------------------------

// optionalKeyGroupSource builds a group that is domain-exhaustive over its
// declared domain but whose `gate` key is declared optional (or, when
// alwaysPresent is true, always-present). It is the narrowing's positive
// case and its paired negative control: the two differ in ONE marker.
func optionalKeyGroupSource(alwaysPresent bool) string {
	required := ""
	if alwaysPresent {
		required = "required = true\n"
	}
	return declBlock(`
[tags.gate]
provenance = "owned"
kind = "bool"
single_valued = true
`+required) + `
[[rule]]
id = "gate-open"
source = "t:open"
[rule.match.recognized]
eq = "go"
[rule.guard.all.gate]
eq = true
[rule.write]

[[rule]]
id = "gate-shut"
source = "t:shut"
[rule.match.recognized]
eq = "go"
[rule.guard.all.gate]
eq = false
[rule.write]
`
}

// twoRuleIdenticalGuardSource is two ordinary rows in ONE selection
// context carrying a byte-identical guard atom. Both dimensions are
// finitely declared, single-valued, and always-present, so the group is
// provable and the only defect is the overlap.
func twoRuleIdenticalGuardSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`) + `
[[rule]]
id = "dup-a"
source = "t:a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "large"
[rule.write]
profile = "large"

[[rule]]
id = "dup-b"
source = "t:b"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "large"
[rule.write]
profile = "small"
`
}

// completePartitionSource is one group over one single-valued,
// always-present enum whose two rows partition the domain exactly: the
// only case in this suite whose expected verdict is green.
func completePartitionSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`) + `
[[rule]]
id = "part-small"
source = "t:small"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"
[rule.write]
profile = "large"

[[rule]]
id = "part-large"
source = "t:large"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "large"
[rule.write]
profile = "small"
`
}

// mixedBlockSource is one row carrying atoms in BOTH guard blocks over
// finitely-declared, single-valued, always-present keys.
func mixedBlockSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true

[tags.iter]
provenance = "owned"
kind = "int"
min = 0
max = 3
single_valued = true
required = true
`) + `
[[rule]]
id = "mixed"
source = "t:mixed"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "large"
[rule.guard.unless.iter]
gte = 3
[rule.write]

[[rule]]
id = "mixed-peer"
source = "t:peer"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"
[rule.write]
`
}

// twoAtomUnlessSource is one row whose `unless` block carries TWO atoms,
// so per-atom negation and full-conjunction subtraction differ observably.
func twoAtomUnlessSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true

[tags.flag]
provenance = "owned"
kind = "bool"
single_valued = true
required = true
`) + `
[[rule]]
id = "two-atom-unless"
source = "t:unless"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
in = ["small", "large"]
[rule.guard.unless.profile]
eq = "large"
[rule.guard.unless.flag]
eq = true
[rule.write]
`
}

// optionalAtomInBlockSource places the SAME value atom over the SAME
// optional key in the named block, holding every other input fixed. It is
// the `all` / `unless` symmetry witness.
func optionalAtomInBlockSource(block table.Block) string {
	blockName := "all"
	if block == table.BlockUnless {
		blockName = "unless"
	}
	return declBlock(`
[tags.gate]
provenance = "owned"
kind = "bool"
single_valued = true
`) + `
[[rule]]
id = "gate-open"
source = "t:open"
[rule.match.recognized]
eq = "go"
[rule.guard.` + blockName + `.gate]
eq = true
[rule.write]

[[rule]]
id = "gate-shut"
source = "t:shut"
[rule.match.recognized]
eq = "go"
[rule.guard.` + blockName + `.gate]
eq = false
[rule.write]
`
}

// withheldPlusOverlapSource is one group carrying BOTH two refusing rows
// (over an optional key) and two decidable rows that overlap. The
// withholding must not suppress the overlap, and no coverage gap may be
// emitted for the group.
func withheldPlusOverlapSource() string {
	return declBlock(`
[tags.gate]
provenance = "owned"
kind = "bool"
single_valued = true

[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`) + `
[[rule]]
id = "refuser-a"
source = "t:ra"
[rule.match.recognized]
eq = "go"
[rule.guard.all.gate]
eq = true
[rule.write]

[[rule]]
id = "refuser-b"
source = "t:rb"
[rule.match.recognized]
eq = "go"
[rule.guard.unless.gate]
eq = false
[rule.write]

[[rule]]
id = "decidable-a"
source = "t:da"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "large"
[rule.write]

[[rule]]
id = "decidable-b"
source = "t:db"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
in = ["small", "large"]
[rule.write]
`
}

// twoRefusingRowsSource is a group with exactly two refusing rows, so
// "one finding per refusing row, never just the first" is observable.
func twoRefusingRowsSource() string {
	return declBlock(`
[tags.gate]
provenance = "owned"
kind = "bool"
single_valued = true
`) + `
[[rule]]
id = "refuser-a"
source = "t:ra"
[rule.match.recognized]
eq = "go"
[rule.guard.all.gate]
eq = true
[rule.write]

[[rule]]
id = "refuser-b"
source = "t:rb"
[rule.match.recognized]
eq = "go"
[rule.guard.all.gate]
eq = false
[rule.write]
`
}

// guardedEscapeRefuserSource is a group whose ONLY possibly-refusing row
// is a guarded ESCAPE row. The ordinary rows partition their always-present
// dimension, so a narrowing reading the ordinary-row population would
// certify the group green.
func guardedEscapeRefuserSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true

[tags.gate]
provenance = "owned"
kind = "bool"
single_valued = true
`) + `
[[rule]]
id = "ordinary"
source = "t:ord"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
in = ["small", "large"]
[rule.write]

[[rule]]
id = "guarded-escape"
source = "t:esc"
escape = ["no_match"]
[rule.match.recognized]
eq = "go"
[rule.guard.all.gate]
eq = true
`
}

// existsPartitionSource partitions the PRESENCE dimension of one optional
// key with two `exists` atoms, and carries no value atom at all — so no
// row can refuse and the group proves over the presence dimension alone.
func existsPartitionSource() string {
	return declBlock(`
[tags.gate]
provenance = "owned"
kind = "bool"
single_valued = true
`) + `
[[rule]]
id = "gate-present"
source = "t:present"
[rule.match.recognized]
eq = "go"
[rule.guard.all.gate]
exists = true
[rule.write]

[[rule]]
id = "gate-absent"
source = "t:absent"
[rule.match.recognized]
eq = "go"
[rule.guard.all.gate]
exists = false
[rule.write]
`
}

// --- group helpers -------------------------------------------------------

// groupOf returns the single scoped row group containing ruleID.
func groupOf(t *testing.T, m *table.Model, ruleID string) guard.Group {
	t.Helper()

	var found []guard.Group
	for _, g := range guard.Groups(m) {
		if slices.Contains(ruleIDsOf(g), ruleID) {
			found = append(found, g)
		}
	}
	switch len(found) {
	case 1:
		return found[0]
	case 0:
		t.Fatalf("no scoped row group contains rule %q", ruleID)
	default:
		t.Fatalf("rule %q sits in %d groups; a row shares exactly one "+
			"selection context", ruleID, len(found))
	}
	return guard.Group{}
}

// ruleIDsOf lists the rule ids a group carries.
func ruleIDsOf(g guard.Group) []string {
	out := make([]string, 0, len(g.Rows))
	for _, r := range g.Rows {
		out = append(out, r.RuleID)
	}
	slices.Sort(out)
	return out
}

// --- product fixture sources ---------------------------------------------

// productOnlyGapSource covers each of two dimensions completely when they
// are checked independently, but leaves the assignment (large, false)
// uncovered in the 2-D product.
func productOnlyGapSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true

[tags.flag]
provenance = "owned"
kind = "bool"
single_valued = true
required = true
`) + `
[[rule]]
id = "gap-a"
source = "t:ga"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"
[rule.write]

[[rule]]
id = "gap-b"
source = "t:gb"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "large"
[rule.guard.all.flag]
eq = true
[rule.write]
`
}

// twoContextSource carries two selection contexts over ONE outcome,
// distinguished by their match pattern (the source state).
func twoContextSource() string {
	return declBlock(`
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["Draft", "Final"]
single_valued = true
required = true

[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`) + `
[[rule]]
id = "draft-a"
source = "t:da"
[rule.match.status]
eq = "Draft"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"
[rule.write]

[[rule]]
id = "draft-b"
source = "t:db"
[rule.match.status]
eq = "Draft"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "large"
[rule.write]

[[rule]]
id = "final-a"
source = "t:fa"
[rule.match.status]
eq = "Final"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
in = ["small", "large"]
[rule.write]
`
}

// twoOutcomeSource carries one source state and TWO recognized outcomes,
// so the outcome half of the selection context splits the group.
func twoOutcomeSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`) + `
[[rule]]
id = "on-go"
source = "t:go"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
in = ["small", "large"]
[rule.write]

[[rule]]
id = "on-stop"
source = "t:stop"
[rule.match.recognized]
eq = "stop"
[rule.guard.all.profile]
in = ["small", "large"]
[rule.write]
`
}

// participationSource carries two keys each guarded by exactly ONE row —
// one under `all` with `eq`, one under `unless`. Neither discriminates
// between the rows, so a rows-must-differ reading would drop both.
func participationSource() string {
	return declBlock(`
[tags.only_in_all]
provenance = "owned"
kind = "bool"
single_valued = true
required = true

[tags.only_in_unless]
provenance = "owned"
kind = "bool"
single_valued = true
required = true
`) + `
[[rule]]
id = "part-a"
source = "t:pa"
[rule.match.recognized]
eq = "go"
[rule.guard.all.only_in_all]
eq = true
[rule.write]

[[rule]]
id = "part-b"
source = "t:pb"
[rule.match.recognized]
eq = "go"
[rule.guard.unless.only_in_unless]
eq = true
[rule.write]
`
}

// identicalConstraintSource guards one key IDENTICALLY in every row, so it
// does not discriminate — and must still bound the product.
func identicalConstraintSource() string {
	return declBlock(`
[tags.shared]
provenance = "owned"
kind = "bool"
single_valued = true
required = true

[tags.other]
provenance = "owned"
kind = "bool"
single_valued = true
required = true
`) + `
[[rule]]
id = "same-a"
source = "t:sa"
[rule.match.recognized]
eq = "go"
[rule.guard.all.shared]
eq = true
[rule.guard.all.other]
eq = true
[rule.write]

[[rule]]
id = "same-b"
source = "t:sb"
[rule.match.recognized]
eq = "go"
[rule.guard.all.shared]
eq = true
[rule.guard.all.other]
eq = false
[rule.write]
`
}

// identicalConstraintUnboundedSource is identicalConstraintSource with the
// shared key declaring no finite domain.
func identicalConstraintUnboundedSource() string {
	return declBlock(`
[tags.shared]
provenance = "observed"
kind = "scalar"
required = true

[tags.other]
provenance = "owned"
kind = "bool"
single_valued = true
required = true
`) + `
[[rule]]
id = "same-a"
source = "t:sa"
[rule.match.recognized]
eq = "go"
[rule.guard.all.shared]
eq = "x"
[rule.guard.all.other]
eq = true
[rule.write]

[[rule]]
id = "same-b"
source = "t:sb"
[rule.match.recognized]
eq = "go"
[rule.guard.all.shared]
eq = "x"
[rule.guard.all.other]
eq = false
[rule.write]
`
}

// dualUseKeySource uses ONE key as a guard key in one group and as a match
// key in another, so the per-atom-per-group split is observable.
func dualUseKeySource() string {
	return declBlock(`
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["Draft", "Final"]
single_valued = true
required = true

[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`) + `
[[rule]]
id = "guarded-use"
source = "t:gu"
[rule.match.status]
eq = "Draft"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
in = ["small", "large"]
[rule.write]

[[rule]]
id = "matched-use"
source = "t:mu"
[rule.match.status]
eq = "Final"
[rule.match.profile]
eq = "large"
[rule.match.recognized]
eq = "go"
[rule.write]
`
}

// existsOverKindSource declares `gate` single-valued or not, guarding it
// with an `exists` atom alone — the operator unaffected by the marker.
func existsOverKindSource(singleValued bool) string {
	marker := ""
	if singleValued {
		marker = "single_valued = true\n"
	}
	return declBlock(`
[tags.gate]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
`+marker) + `
[[rule]]
id = "gate-present"
source = "t:present"
[rule.match.recognized]
eq = "go"
[rule.guard.all.gate]
exists = true
[rule.write]

[[rule]]
id = "gate-absent"
source = "t:absent"
[rule.match.recognized]
eq = "go"
[rule.guard.all.gate]
exists = false
[rule.write]
`
}

// existsIgnoredWouldGoGreenSource covers both VALUES of an optional key but
// only the {present} half of its presence dimension. Dropping the `exists`
// atom from the product would certify it green.
func existsIgnoredWouldGoGreenSource() string {
	return declBlock(`
[tags.gate]
provenance = "owned"
kind = "bool"
single_valued = true
`) + `
[[rule]]
id = "exists-only-present"
source = "t:p"
[rule.match.recognized]
eq = "go"
[rule.guard.all.gate]
exists = true
[rule.write]
`
}

// setAndEnumSource declares one single-valued enum and one set tag with an
// element universe, for the two projection directions.
func setAndEnumSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "mid", "large"]
single_valued = true
required = true

[tags.labels]
provenance = "owned"
kind = "set"
elements = ["bug", "chore"]
required = true
`) + `
[[rule]]
id = "set-and-enum"
source = "t:se"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
in = ["small", "mid"]
[rule.guard.all.labels]
contains = ["bug"]
[rule.write]
`
}

// unmarkedEnumGuardSource guards a finite enum that declares NO
// single-valued marker with single-value operators.
func unmarkedEnumGuardSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
required = true
`) + `
[[rule]]
id = "unmarked-a"
source = "t:ua"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "large"
[rule.write]

[[rule]]
id = "unmarked-b"
source = "t:ub"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"
[rule.write]
`
}

// unboundedIntGuardSource guards an int declaring no {min..max} bound.
func unboundedIntGuardSource() string {
	return declBlock(`
[tags.iter]
provenance = "owned"
kind = "int"
single_valued = true
required = true
`) + `
[[rule]]
id = "unbounded-a"
source = "t:ia"
[rule.match.recognized]
eq = "go"
[rule.guard.all.iter]
lt = 3
[rule.write]

[[rule]]
id = "unbounded-b"
source = "t:ib"
[rule.match.recognized]
eq = "go"
[rule.guard.all.iter]
gte = 3
[rule.write]
`
}

// undeclaredDomainExamplesSource guards a scalar whose covered examples
// enumerate every value the rules mention — a lint treating examples as
// complete would certify it.
func undeclaredDomainExamplesSource() string {
	return declBlock(`
[tags.loose]
provenance = "observed"
kind = "scalar"
required = true
`) + `
[[rule]]
id = "loose-a"
source = "t:la"
[rule.match.recognized]
eq = "go"
[rule.guard.all.loose]
eq = "one"
[rule.write]

[[rule]]
id = "loose-b"
source = "t:lb"
[rule.match.recognized]
eq = "go"
[rule.guard.all.loose]
eq = "two"
[rule.write]
`
}

// overLargeProductSource declares enough finitely-declared dimensions that
// the scoped product exceeds any sane published bound while every
// dimension stays provable.
func overLargeProductSource() string {
	var decls, guards string
	for _, k := range []string{"a", "b", "c", "d", "e", "f", "g", "h"} {
		decls += `
[tags.set_` + k + `]
provenance = "owned"
kind = "set"
elements = ["e0", "e1", "e2", "e3", "e4", "e5", "e6", "e7", "e8", "e9"]
required = true
`
		guards += `[rule.guard.all.set_` + k + `]
contains = ["e0"]
`
	}
	return declBlock(decls) + `
[[rule]]
id = "too-large"
source = "t:tl"
[rule.match.recognized]
eq = "go"
` + guards + `[rule.write]
`
}

// --- escape and set fixture sources --------------------------------------

// bareEscapeClosesGapSource is a group whose ordinary row leaves a gap and
// whose BARE escape row (declaring only no_match) closes it.
func bareEscapeClosesGapSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`) + `
[[rule]]
id = "partial"
source = "t:partial"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"
[rule.write]

[[rule]]
id = "bare-escape"
source = "t:bare"
escape = ["no_match"]
[rule.match.recognized]
eq = "go"
`
}

// overlappingOrdinaryBareEscapeSource is REQ-77's witness: its ordinary
// population OVERLAPS — so the `ambiguous_match` arm is reachable per RDR
// 0006's reachability precondition — while its bare escape row declares only
// `no_match`. The two arms therefore compute different unions AND the
// undeclared arm draws its gap finding, which is what REQ-77 claims. An
// overlap-free population cannot witness the clause: 0006 treats the
// `ambiguous_match` arm as vacuously closed there.
func overlappingOrdinaryBareEscapeSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`) + `
[[rule]]
id = "ov-a"
source = "t:oa"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"
[rule.write]

[[rule]]
id = "ov-b"
source = "t:ob"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"
[rule.write]

[[rule]]
id = "bare-escape"
source = "t:bare"
escape = ["no_match"]
[rule.match.recognized]
eq = "go"
`
}

// overlappingOrdinaryCoversProductSource is the union-source witness: its
// ordinary population OVERLAPS on `small` — so the `ambiguous_match` arm is
// reachable — and jointly COVERS the whole product {small, large}, while no
// row declares `ambiguous_match` at all. Nothing rescues the ambiguity, so
// the arm must draw its gap; a union that counted the overlapping ordinary
// rows would let the very overlap that mints the refusal certify it
// rescued.
func overlappingOrdinaryCoversProductSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`) + `
[[rule]]
id = "ov-a"
source = "t:oa"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"
[rule.write]

[[rule]]
id = "ov-b"
source = "t:ob"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
in = ["small", "large"]
[rule.write]
`
}

// overlappingOrdinaryAmbiguousEscapeSource is the same overlapping,
// fully-covering ordinary population plus a BARE escape row declaring
// `ambiguous_match`. That row denotes the whole scoped product, so the arm
// really is rescued and draws no gap.
func overlappingOrdinaryAmbiguousEscapeSource() string {
	return overlappingOrdinaryCoversProductSource() + `
[[rule]]
id = "ambig-escape"
source = "t:ambig"
escape = ["ambiguous_match"]
[rule.match.recognized]
eq = "go"
`
}

// overlappingOrdinaryGapAmbiguousEscapeSource overlaps on `small` — making
// the `ambiguous_match` arm reachable — while leaving `large` uncovered, so
// the ordinary population does NOT close alone. Its bare escape row
// declares BOTH rescuable classes and is therefore what closes each arm,
// which is reported rather than taken silently.
func overlappingOrdinaryGapAmbiguousEscapeSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`) + `
[[rule]]
id = "ov-a"
source = "t:oa"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"
[rule.write]

[[rule]]
id = "ov-b"
source = "t:ob"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"
[rule.write]

[[rule]]
id = "ambig-escape"
source = "t:ambig"
escape = ["no_match", "ambiguous_match"]
[rule.match.recognized]
eq = "go"
`
}

// guardedEscapePartialSource carries a GUARDED escape row, whose accepted
// assignments are its own rather than the whole product.
func guardedEscapePartialSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`) + `
[[rule]]
id = "ordinary"
source = "t:ord"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"
[rule.write]

[[rule]]
id = "guarded-escape"
source = "t:esc"
escape = ["no_match"]
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "large"
`
}

// optionalKeyPlusBareEscapeSource is optionalKeyGroupSource(false) with a
// bare escape row added — the false-green A19 guards.
func optionalKeyPlusBareEscapeSource() string {
	return optionalKeyGroupSource(false) + `
[[rule]]
id = "bare-escape"
source = "t:bare"
escape = ["no_match"]
[rule.match.recognized]
eq = "go"
`
}

// twoPopulationSource carries two overlapping ordinary rows AND two
// overlapping escape rows for one class, so both populations are exercised
// and the cross-population pairing can be shown absent.
//
// The domain carries a third value the ORDINARY rows do not reach, so the
// escape rows contribute an assignment their guarded peers do not. Without
// it the ordinary rows close coverage on their own and "the escape rows
// were excluded from coverage" becomes untestable — the two unions agree
// whether the escape rows were counted or dropped.
func twoPopulationSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "mid", "large"]
single_valued = true
required = true
`) + `
[[rule]]
id = "ord-a"
source = "t:oa"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"
[rule.write]

[[rule]]
id = "ord-b"
source = "t:ob"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
in = ["small", "mid"]
[rule.write]

[[rule]]
id = "esc-a"
source = "t:ea"
escape = ["no_match"]
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "large"

[[rule]]
id = "esc-b"
source = "t:eb"
escape = ["no_match"]
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
in = ["small", "mid", "large"]
`
}

// multiClassEscapeSource carries one escape row declaring BOTH classes and
// two declaring one each, so the per-class partition is observable.
func multiClassEscapeSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`) + `
[[rule]]
id = "ordinary"
source = "t:ord"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
in = ["small", "large"]
[rule.write]

[[rule]]
id = "both"
source = "t:both"
escape = ["no_match", "ambiguous_match"]
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"

[[rule]]
id = "only-no-match"
source = "t:onm"
escape = ["no_match"]
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"

[[rule]]
id = "only-ambiguous"
source = "t:oam"
escape = ["ambiguous_match"]
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "large"
`
}

// twoClassOverlapPairSource carries two escape rows that BOTH declare both
// classes and overlap, so the pair yields one finding per class.
func twoClassOverlapPairSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`) + `
[[rule]]
id = "ordinary"
source = "t:ord"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
in = ["small", "large"]
[rule.write]

[[rule]]
id = "pair-a"
source = "t:pa"
escape = ["no_match", "ambiguous_match"]
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"

[[rule]]
id = "pair-b"
source = "t:pb"
escape = ["no_match", "ambiguous_match"]
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
in = ["small", "large"]
`
}

// setLiteralOrderSource authors one `contains` literal in the given
// element order.
func setLiteralOrderSource(literal string) string {
	return declBlock(`
[tags.labels]
provenance = "owned"
kind = "set"
elements = ["bug", "chore"]
required = true
`) + `
[[rule]]
id = "set-lit"
source = "t:sl"
[rule.match.recognized]
eq = "go"
[rule.guard.all.labels]
contains = ` + literal + `
[rule.write]
`
}

// setOverlapSource is two rows whose `contains` literals overlap, with the
// first row's literal authored in the given element order.
func setOverlapSource(literal string) string {
	return declBlock(`
[tags.labels]
provenance = "owned"
kind = "set"
elements = ["bug", "chore"]
required = true
`) + `
[[rule]]
id = "set-a"
source = "t:sa"
[rule.match.recognized]
eq = "go"
[rule.guard.all.labels]
contains = ` + literal + `
[rule.write]

[[rule]]
id = "set-b"
source = "t:sb"
[rule.match.recognized]
eq = "go"
[rule.guard.all.labels]
contains = ["bug"]
[rule.write]
`
}

// unprovablePlusLargeSource carries BOTH an unprovable dimension and
// enough finite dimensions that the product would exceed the bound if it
// were computable. No cardinality may be reported.
func unprovablePlusLargeSource() string {
	var decls, guards string
	for _, k := range []string{"a", "b", "c", "d", "e", "f", "g", "h"} {
		decls += `
[tags.set_` + k + `]
provenance = "owned"
kind = "set"
elements = ["e0", "e1", "e2", "e3", "e4", "e5", "e6", "e7", "e8", "e9"]
required = true
`
		guards += `[rule.guard.all.set_` + k + `]
contains = ["e0"]
`
	}
	decls += `
[tags.loose]
provenance = "observed"
kind = "scalar"
required = true
`
	return declBlock(decls) + `
[[rule]]
id = "mixed-a"
source = "t:ma"
[rule.match.recognized]
eq = "go"
` + guards + `[rule.guard.all.loose]
eq = "x"
[rule.write]
`
}

// overBoundStructuralPlusOverlapSource drives the over-bound,
// structurally-unprovable branch while keeping an overlap decidable.
//
// `loose` is an enum with a finite declared domain — so `Cardinality` is
// computable — but it carries no single-valued marker, which makes its
// spread 2^12 and pushes the product past the bound, AND makes the `eq`
// atom over it structurally unprojectable. `flag` is a plain required bool,
// so the two rows guarding only `flag` have an accepted-assignment set over
// the decidable sub-product and their overlap is decidable without ever
// touching `loose`.
func overBoundStructuralPlusOverlapSource() string {
	return declBlock(`
[tags.loose]
provenance = "owned"
kind = "enum"
domain = ["v0", "v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8", "v9", "v10", "v11"]
required = true

[tags.flag]
provenance = "owned"
kind = "bool"
single_valued = true
required = true
`) + `
[[rule]]
id = "unprovable-row"
source = "t:up"
[rule.match.recognized]
eq = "go"
[rule.guard.all.loose]
eq = "v0"
[rule.write]

[[rule]]
id = "decidable-a"
source = "t:da"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = true
[rule.write]

[[rule]]
id = "decidable-b"
source = "t:db"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = true
[rule.write]
`
}

// overBoundStructuralDisjointSource is the same shape with the two
// decidable rows made mutually exclusive, so the overlap the branch now
// runs reports nothing rather than reporting spuriously.
func overBoundStructuralDisjointSource() string {
	return declBlock(`
[tags.loose]
provenance = "owned"
kind = "enum"
domain = ["v0", "v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8", "v9", "v10", "v11"]
required = true

[tags.flag]
provenance = "owned"
kind = "bool"
single_valued = true
required = true
`) + `
[[rule]]
id = "unprovable-row"
source = "t:up"
[rule.match.recognized]
eq = "go"
[rule.guard.all.loose]
eq = "v0"
[rule.write]

[[rule]]
id = "decidable-a"
source = "t:da"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = true
[rule.write]

[[rule]]
id = "decidable-b"
source = "t:db"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = false
[rule.write]
`
}

// overBoundStructuralAllUnprovableSource is the negative control: EVERY row
// in the group guards the unprovable dimension, so no row has an accepted
// assignment set at all and the overlap check correctly reports nothing.
func overBoundStructuralAllUnprovableSource() string {
	return declBlock(`
[tags.loose]
provenance = "owned"
kind = "enum"
domain = ["v0", "v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8", "v9", "v10", "v11"]
required = true

[tags.flag]
provenance = "owned"
kind = "bool"
single_valued = true
required = true
`) + `
[[rule]]
id = "unprovable-a"
source = "t:ua"
[rule.match.recognized]
eq = "go"
[rule.guard.all.loose]
eq = "v0"
[rule.guard.all.flag]
eq = true
[rule.write]

[[rule]]
id = "unprovable-b"
source = "t:ub"
[rule.match.recognized]
eq = "go"
[rule.guard.all.loose]
eq = "v0"
[rule.guard.all.flag]
eq = true
[rule.write]
`
}

// --- diagnostics fixture sources -----------------------------------------

// gapAndOverlapSource carries BOTH an intentional gap and an intentional
// overlap in one group over a fully-provable 2-D product.
func gapAndOverlapSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true

[tags.flag]
provenance = "owned"
kind = "bool"
single_valued = true
required = true
`) + `
[[rule]]
id = "ov-a"
source = "t:oa"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"
[rule.write]

[[rule]]
id = "ov-b"
source = "t:ob"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"
[rule.guard.all.flag]
eq = true
[rule.write]

[[rule]]
id = "gap-large-true"
source = "t:gl"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "large"
[rule.guard.all.flag]
eq = true
[rule.write]
`
}

// twoUnprovableDimensionsSource carries TWO separately unprovable
// dimensions in one group, so "never stop at the first" is observable.
func twoUnprovableDimensionsSource() string {
	return declBlock(`
[tags.loose_a]
provenance = "observed"
kind = "scalar"
required = true

[tags.loose_b]
provenance = "observed"
kind = "scalar"
required = true
`) + `
[[rule]]
id = "two-unprovable"
source = "t:tu"
[rule.match.recognized]
eq = "go"
[rule.guard.all.loose_a]
eq = "one"
[rule.guard.all.loose_b]
eq = "one"
[rule.write]
`
}

// reversedTwoRefusingRowsSource is twoRefusingRowsSource with the two rows
// authored in the opposite order.
func reversedTwoRefusingRowsSource() string {
	return declBlock(`
[tags.gate]
provenance = "owned"
kind = "bool"
single_valued = true
`) + `
[[rule]]
id = "refuser-b"
source = "t:rb"
[rule.match.recognized]
eq = "go"
[rule.guard.all.gate]
eq = false
[rule.write]

[[rule]]
id = "refuser-a"
source = "t:ra"
[rule.match.recognized]
eq = "go"
[rule.guard.all.gate]
eq = true
[rule.write]
`
}

// ownedBeforeWriteSource matches an OWNED tag no reachable predecessor
// writes: no rule's write block ever assigns `never_written`.
func ownedBeforeWriteSource() string {
	return declBlock(`
[tags.never_written]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[tags.other]
provenance = "owned"
kind = "bool"
single_valued = true
required = true
`) + `
[[rule]]
id = "reads-unwritten"
source = "t:ru"
[rule.match.never_written]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.other]
in = ["true", "false"]
[rule.write]
other = true
`
}

// observedMatchSource matches an OBSERVED tag the same way, which carries
// no predecessor-write obligation.
func observedMatchSource() string {
	return declBlock(`
[tags.watched]
provenance = "observed"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[tags.other]
provenance = "owned"
kind = "bool"
single_valued = true
required = true
`) + `
[[rule]]
id = "reads-observed"
source = "t:ro"
[rule.match.watched]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.other]
in = ["true", "false"]
[rule.write]
other = true
`
}

// optionalButAlwaysWrittenSource declares a key OPTIONAL while every rule
// writes it — so a reachability-driven reading would find it always
// present, and the declaration-driven one must still say can-refuse.
func optionalButAlwaysWrittenSource() string {
	return declBlock(`
[tags.gate]
provenance = "owned"
kind = "bool"
single_valued = true
`) + `
[[rule]]
id = "writer"
source = "t:w"
[rule.match.recognized]
eq = "stop"
[rule.write]
gate = true

[[rule]]
id = "reader"
source = "t:r"
[rule.match.recognized]
eq = "go"
[rule.guard.all.gate]
eq = true
[rule.write]
gate = false

[[rule]]
id = "reader-peer"
source = "t:rp"
[rule.match.recognized]
eq = "go"
[rule.guard.all.gate]
eq = false
[rule.write]
gate = true
`
}

// --- MVV and testing-strategy fixture sources ----------------------------

// mvvRDRSlice is the MVV's RDR flow slice. It carries three scoped row
// groups over one outcome, distinguished by their match pattern:
//
//   - `stage eq propose` — the COMPLETE PARTITION, the MVV's only passing
//     exhaustiveness verdict (obligation 2);
//   - `stage eq prelock` — the intentional GAP, visible only in the
//     2-D product (obligation 3);
//   - `stage eq resolve` — the intentional OVERLAP, likewise 2-D only.
//
// Every target-flow dimension declares `single_valued`, per A21/REQ-127,
// except `lens`, which is a `set` under `contains` — the recorded
// alternative for a co-occurring dimension.
func mvvRDRSlice() string {
	return `
outcomes = ["round-clean", "stop"]

[model]
id = "rdr"
version = 1
description = "MVV RDR flow slice."

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.stage]
provenance = "owned"
kind = "enum"
domain = ["propose", "prelock", "resolve"]
single_valued = true
required = true

[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true

[tags.cluster_eligible]
provenance = "owned"
kind = "bool"
single_valued = true
required = true

[tags.prelock_iterations]
provenance = "owned"
kind = "int"
min = 0
max = 3
single_valued = true
required = true

[tags.rewind_target]
provenance = "owned"
kind = "enum"
domain = ["assumptions", "approach"]
single_valued = true
required = true

[tags.lens]
provenance = "owned"
kind = "set"
elements = ["grounding", "critique"]
required = true

[read.own]
role = "rdr"
path = "rdr.own"
keys = ["cluster_eligible", "lens", "prelock_iterations", "profile", "rewind_target", "stage"]
timeout = "2s"

[write.own]
role = "rdr"
path = "rdr.own"
keys = ["cluster_eligible", "lens", "prelock_iterations", "profile", "rewind_target", "stage"]
timeout = "2s"
read_back = true

# --- group A: the complete partition (green) -----------------------------
# Two rows partitioning profile x cluster_eligible's 2 x 2 product with
# no unless block anywhere - deviations D1's check.

[[rule]]
id = "partition-small"
source = "rdr:ps"
[rule.match.stage]
eq = "propose"
[rule.match.recognized]
eq = "round-clean"
[rule.guard.all.profile]
eq = "small"
[rule.write]
profile = "large"

[[rule]]
id = "partition-large"
source = "rdr:pl"
[rule.match.stage]
eq = "propose"
[rule.match.recognized]
eq = "round-clean"
[rule.guard.all.profile]
eq = "large"
[rule.write]
profile = "small"

# --- group B: the intentional 2-D gap ------------------------------------
# profile is covered by one row and prelock_iterations by the other, so
# each dimension is complete on its own; the assignment
# (profile = large, prelock_iterations = 3) is uncovered.

[[rule]]
id = "gap-a"
source = "rdr:ga"
[rule.match.stage]
eq = "prelock"
[rule.match.recognized]
eq = "round-clean"
[rule.guard.all.profile]
eq = "small"
[rule.guard.unless.prelock_iterations]
gte = 3
[rule.write]
stage = "resolve"

[[rule]]
id = "gap-b"
source = "rdr:gb"
[rule.match.stage]
eq = "prelock"
[rule.match.recognized]
eq = "round-clean"
[rule.guard.all.profile]
eq = "large"
[rule.guard.all.prelock_iterations]
lt = 3
[rule.write]
stage = "resolve"

# --- group C: the intentional 2-D overlap --------------------------------
# Neither rewind_target alone nor lens alone reveals it: the two rows
# intersect only at (rewind_target = assumptions, lens contains grounding).

[[rule]]
id = "ov-a"
source = "rdr:oa"
[rule.match.stage]
eq = "resolve"
[rule.match.recognized]
eq = "round-clean"
[rule.guard.all.rewind_target]
in = ["assumptions", "approach"]
[rule.guard.all.lens]
contains = ["grounding"]
[rule.write]
stage = "propose"

[[rule]]
id = "ov-b"
source = "rdr:ob"
[rule.match.stage]
eq = "resolve"
[rule.match.recognized]
eq = "round-clean"
[rule.guard.all.rewind_target]
eq = "assumptions"
[rule.guard.all.cluster_eligible]
exists = true
[rule.write]
stage = "propose"
`
}

// mvvKataSlice is the MVV's kata flow slice. It carries the narrowing's
// POSITIVE case — a domain-exhaustive group over a possibly-absent guard
// key — and its always-present NEGATIVE CONTROL.
func mvvKataSlice() string {
	return `
outcomes = ["accepted", "stop"]

[model]
id = "kata"
version = 1
description = "MVV kata flow slice."

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.phase]
provenance = "owned"
kind = "enum"
domain = ["review", "ship"]
single_valued = true
required = true

# The possibly-absent guard key: no required marker, so it is optional.
[tags.assignee]
provenance = "observed"
kind = "bool"
single_valued = true

# The control's key: the SAME shape, declared always-present.
[tags.triaged]
provenance = "owned"
kind = "bool"
single_valued = true
required = true

[tags.labels]
provenance = "owned"
kind = "set"
elements = ["bug", "chore"]
required = true

[read.own]
role = "kata"
path = "kata.own"
keys = ["labels", "phase", "triaged"]
timeout = "2s"

[write.own]
role = "kata"
path = "kata.own"
keys = ["labels", "phase", "triaged"]
timeout = "2s"
read_back = true

# --- the narrowing's POSITIVE case ---------------------------------------
# Domain-exhaustive over assignee's two values, yet assignee is declared
# optional, so both rows can refuse guard_unevaluable at runtime.

[[rule]]
id = "kata-absent-key"
source = "kata:aa"
[rule.match.phase]
eq = "review"
[rule.match.recognized]
eq = "accepted"
[rule.guard.all.assignee]
eq = true
[rule.write]
phase = "ship"

[[rule]]
id = "kata-absent-key-peer"
source = "kata:ap"
[rule.match.phase]
eq = "review"
[rule.match.recognized]
eq = "accepted"
[rule.guard.all.assignee]
eq = false
[rule.write]
phase = "ship"

# --- the NEGATIVE CONTROL ------------------------------------------------
# Otherwise identical, over a key declared always-present. It certifies
# green, proving the narrowing is tight rather than blanket.

[[rule]]
id = "kata-control-a"
source = "kata:ca"
[rule.match.phase]
eq = "ship"
[rule.match.recognized]
eq = "accepted"
[rule.guard.all.triaged]
eq = true
[rule.write]
phase = "review"

[[rule]]
id = "kata-control-b"
source = "kata:cb"
[rule.match.phase]
eq = "ship"
[rule.match.recognized]
eq = "accepted"
[rule.guard.all.triaged]
eq = false
[rule.write]
phase = "review"
`
}

// mvvLegalView is a conforming view over the RDR slice's prelock group that
// selects exactly one row.
func mvvLegalView() guard.View {
	return guard.View{
		"stage":              "prelock",
		"profile":            "small",
		"cluster_eligible":   "true",
		"prelock_iterations": "1",
		"rewind_target":      "assumptions",
		"lens":               `["grounding"]`,
	}
}

// mvvNoRowView is a conforming view over the prelock group that no row
// accepts — the 2-D hole.
func mvvNoRowView() guard.View {
	v := mvvLegalView()
	v["profile"] = "large"
	v["prelock_iterations"] = "3"
	return v
}

// shapePairSource builds one group whose scoped product has a cardinality
// fixed relative to the published bound B: over it when over is true, just
// under it otherwise. `wide` picks few wide dimensions; otherwise many
// narrow boolean ones. Both shapes reach the same cardinality.
func shapePairSource(wide, over bool) string {
	target := guard.Bound() * 2
	if !over {
		target = guard.Bound() / 2
	}
	if target < 2 {
		target = 2
	}

	var decls, guards string
	add := func(name, decl, atom string) {
		decls += "\n[tags." + name + "]\nprovenance = \"owned\"\n" + decl
		guards += "[rule.guard.all." + name + "]\n" + atom
	}

	if wide {
		// Two wide single-valued int dimensions of ~sqrt(target) values.
		side := 2
		for side*side < target {
			side++
		}
		for i, name := range []string{"wide_a", "wide_b"} {
			_ = i
			add(name,
				"kind = \"int\"\nmin = 1\nmax = "+itoa(side)+"\nsingle_valued = true\nrequired = true\n",
				"gte = 1\n")
		}
	} else {
		// Many narrow single-valued bool dimensions: 2^n.
		n := 1
		for pow2(n) < target {
			n++
		}
		for i := range n {
			add("narrow_"+itoa(i),
				"kind = \"bool\"\nsingle_valued = true\nrequired = true\n",
				"in = [\"true\", \"false\"]\n")
		}
	}

	return declBlock(decls) + `
[[rule]]
id = "shaped"
source = "t:shaped"
[rule.match.recognized]
eq = "go"
` + guards + `[rule.write]
`
}

func pow2(n int) int {
	out := 1
	for range n {
		out *= 2
	}
	return out
}

func itoa(n int) string { return strconv.Itoa(n) }

// markerAcceptanceSource holds every input fixed but the single-valued
// marker on `marker_tag`. Under the marker the two rows partition the
// domain's two values exactly; without it the dimension is a 4-assignment
// powerset the rows cannot close.
func markerAcceptanceSource(marked bool) string {
	marker := ""
	if marked {
		marker = "single_valued = true\n"
	}
	return declBlock(`
[tags.marker_tag]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
required = true
`+marker) + `
[[rule]]
id = "mark-a"
source = "t:ma"
[rule.match.recognized]
eq = "go"
[rule.guard.all.marker_tag]
eq = "a"
[rule.write]

[[rule]]
id = "mark-b"
source = "t:mb"
[rule.match.recognized]
eq = "go"
[rule.guard.all.marker_tag]
eq = "b"
[rule.write]
`
}

// reorderSource authors the same two rows, and the same `in` set literal,
// in one of two orders. The semantics are identical.
func reorderSource(reordered bool) string {
	decls := `
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true

[tags.labels]
provenance = "owned"
kind = "set"
elements = ["bug", "chore"]
required = true
`
	first := `
[[rule]]
id = "ro-a"
source = "t:ra"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
in = ["small", "large"]
[rule.guard.all.labels]
contains = ["bug"]
[rule.write]
`
	second := `
[[rule]]
id = "ro-b"
source = "t:rb"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "small"
[rule.write]
`
	if reordered {
		first = `
[[rule]]
id = "ro-a"
source = "t:ra"
[rule.match.recognized]
eq = "go"
[rule.guard.all.labels]
contains = ["bug"]
[rule.guard.all.profile]
in = ["large", "small"]
[rule.write]
`
		return declBlock(decls) + second + first
	}
	return declBlock(decls) + first + second
}

// sameKeyBothBlocksSource carries atoms over ONE key in BOTH guard blocks.
func sameKeyBothBlocksSource() string {
	return declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "mid", "large"]
single_valued = true
required = true
`) + `
[[rule]]
id = "both-blocks"
source = "t:bb"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
in = ["small", "mid", "large"]
[rule.guard.unless.profile]
eq = "large"
[rule.write]
`
}
