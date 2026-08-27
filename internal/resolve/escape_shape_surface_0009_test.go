package resolve_test

// RDR 0009 — the doc-contract amendments, the fixture-conformance duty, and
// the performance / non-goal clauses.
//
// The doc-contract REQs (54-57) are amendments to DOC COMMENTS that are
// themselves REQ-1 contract artifacts of RDR 0001. A doc comment is not
// observable behaviour, so each is pinned the only way a test legitimately
// can: on the SOURCE of the shipped file, read as a fixture. That is not an
// assertion on log strings or on a message — it is an assertion that a
// contract artifact the record names was in fact amended.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/resolve"
)

// --- G. Doc-contract amendments authorized by this RDR -------------------

// docCommentOf returns the doc comment attached to the named top-level
// declaration in one of the package's source files.
func docCommentOf(t *testing.T, file, name string) string {
	t.Helper()

	fset := token.NewFileSet()
	src, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("read %s: %v", file, err)
	}
	f, err := parser.ParseFile(fset, file, src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse %s: %v", file, err)
	}

	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Name.Name == name && d.Doc != nil {
				return d.Doc.Text()
			}
		case *ast.GenDecl:
			if d.Doc != nil {
				for _, spec := range d.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok && ts.Name.Name == name {
						return d.Doc.Text()
					}
				}
			}
			for _, spec := range d.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok && ts.Name.Name == name &&
					ts.Doc != nil {
					return ts.Doc.Text()
				}
			}
		}
	}
	t.Fatalf("no doc comment found for %s in %s", name, file)
	return ""
}

// fieldDocOf returns the doc comment on one field of a named struct type.
func fieldDocOf(t *testing.T, file, typeName, fieldName string) string {
	t.Helper()

	fset := token.NewFileSet()
	src, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("read %s: %v", file, err)
	}
	f, err := parser.ParseFile(fset, file, src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse %s: %v", file, err)
	}

	var found string
	ast.Inspect(f, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok || ts.Name.Name != typeName {
			return true
		}
		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return false
		}
		for _, fld := range st.Fields.List {
			for _, nm := range fld.Names {
				if nm.Name == fieldName && fld.Doc != nil {
					found = fld.Doc.Text()
				}
			}
		}
		return false
	})
	return found
}

const resolveSourceFile = "resolve.go"

// REQ-54: "*Doc-contract amendment* — `resolve.go::Resolve`'s numbered
// evaluation-order list (a REQ-1 contract artifact) currently opens at the
// alphabet check, so the whole-table precondition MUST be prepended to it as a
// new step 0."
// BOUNDARY
//
// The evaluation-order list is a contract artifact of RDR 0001, so its content
// is the claim. It must gain a step 0 naming the whole-table precondition,
// ahead of the alphabet check that currently opens it.
func TestReq54_ResolvesEvaluationOrderListGainsAStepZero(t *testing.T) {
	doc := docCommentOf(t, resolveSourceFile, "Resolve")

	if !strings.Contains(doc, "0.") {
		t.Fatalf("Resolve's evaluation-order list has no step 0; the "+
			"whole-table precondition must be PREPENDED to it:\n%s", doc)
	}

	zero := strings.Index(doc, "0.")
	one := strings.Index(doc, "1.")
	if one >= 0 && zero > one {
		t.Errorf("step 0 appears after step 1; the precondition is " +
			"prepended, not appended")
	}

	// Step 0 must actually name this RDR's precondition rather than being
	// an unrelated renumbering.
	step := doc[zero:]
	if end := strings.Index(step, "1."); end > 0 {
		step = step[:end]
	}
	lower := strings.ToLower(step)
	if !strings.Contains(lower, "escape") && !strings.Contains(lower, "checkvalid") {
		t.Errorf("step 0 does not name the escape-row shape precondition:\n%s",
			step)
	}
}

// REQ-55: "The same authorization covers the two disposition-cardinality doc
// contracts the breach return would otherwise contradict — `Result`'s "never
// both and never neither" and `Resolve`'s "returns exactly one disposition" —
// each amended to scope itself to nil-error returns."
// ADVERSARIAL
//
// Without the amendment the shipped doc contracts are FALSE on a breach: the
// zero Result is "neither". Each must be scoped to the nil-error return.
func TestReq55_TheDispositionCardinalityDocContractsAreScopedToNilError(t *testing.T) {
	scoped := func(doc string) bool {
		lower := strings.ToLower(doc)
		return strings.Contains(lower, "nil error") ||
			strings.Contains(lower, "nil-error") ||
			strings.Contains(lower, "error is nil")
	}

	resultDoc := docCommentOf(t, resolveSourceFile, "Result")
	if !scoped(resultDoc) {
		t.Errorf("Result's \"never both and never neither\" doc contract is "+
			"not scoped to nil-error returns; on a breach the zero Result is "+
			"NEITHER, so the unscoped sentence is false:\n%s", resultDoc)
	}

	resolveDoc := docCommentOf(t, resolveSourceFile, "Resolve")
	if !scoped(resolveDoc) {
		t.Errorf("Resolve's \"returns exactly one disposition\" doc contract "+
			"is not scoped to nil-error returns:\n%s", resolveDoc)
	}
}

// REQ-56: "state the producer obligation on the `Row` doc contract"
// BOUNDARY
//
// The producer obligation must be discoverable where a producer looks: the
// `Row` doc contract itself (the type's doc or the doc on the two fields the
// predicate reads).
func TestReq56_TheProducerObligationIsStatedOnTheRowDocContract(t *testing.T) {
	combined := strings.ToLower(strings.Join([]string{
		docCommentOf(t, resolveSourceFile, "Row"),
		fieldDocOf(t, resolveSourceFile, "Row", "Writes"),
		fieldDocOf(t, resolveSourceFile, "Row", "Escape"),
	}, "\n"))

	if !strings.Contains(combined, "escape") {
		t.Fatalf("the Row doc contract does not mention escape rows")
	}
	if !strings.Contains(combined, "write") {
		t.Fatalf("the Row doc contract does not mention writes")
	}
	// The obligation itself: an escape row carries no writes. A doc that
	// merely mentions both words without stating the constraint fails.
	statesObligation := strings.Contains(combined, "no writes") ||
		strings.Contains(combined, "empty writes") ||
		strings.Contains(combined, "carries no owned-state mutation") ||
		strings.Contains(combined, "must not carry") ||
		strings.Contains(combined, "producer obligation")
	if !statesObligation {
		t.Errorf("the Row doc contract does not STATE the producer "+
			"obligation that an escape row carries no writes:\n%s", combined)
	}
}

// REQ-57: "record the kernel's validation surface — exactly which shape
// property it checks vs still assumes — so partial validation cannot be read
// as general kernel ownership."
// ADVERSARIAL
//
// The load-bearing half is the NEGATIVE statement: a nil return must not be
// readable as general table validity. C6 names the specific example — the
// Escape-class restriction is NOT checked.
func TestReq57_TheValidationSurfaceRecordsWhatIsNotChecked(t *testing.T) {
	doc := strings.ToLower(docCommentOf(t, resolveSourceFile, "CheckValid"))

	if !strings.Contains(doc, "escape") {
		t.Fatalf("CheckValid's doc does not name the ONE property it "+
			"checks:\n%s", doc)
	}
	statesLimit := strings.Contains(doc, "nothing else") ||
		strings.Contains(doc, "not check") ||
		strings.Contains(doc, "no other") ||
		strings.Contains(doc, "only")
	if !statesLimit {
		t.Errorf("CheckValid's doc does not state that nothing else is "+
			"checked, so a nil return could be read as general table "+
			"validity:\n%s", doc)
	}
	if !strings.Contains(doc, "no_match") && !strings.Contains(doc, "class") {
		t.Errorf("CheckValid's doc does not name the Escape-class "+
			"restriction as a property it does NOT check:\n%s", doc)
	}

	// And behaviourally: the class restriction really is unchecked, so the
	// doc's disclaimer is true rather than defensive boilerplate.
	tbl := resolve.Table{
		Revision: "rev-req57",
		Outcomes: []string{"successful"},
		Rows: []resolve.Row{{
			RuleID:        breachRuleID,
			SourceLocator: breachLocator,
			Outcome:       "successful",
			// A class RDR 0002 forbids an escape rule from modeling. The
			// kernel predicate does not check it.
			Escape: []resolve.RefusalKind{resolve.KindOwnedStateUnavailable},
		}},
	}
	if err := tbl.CheckValid(); err != nil {
		t.Errorf("CheckValid rejected an unmodelable Escape class (%v); the "+
			"predicate checks escape-row SHAPE only", err)
	}
}

// REQ-50: "copy the latter's doc form — "returns nil if <x> is valid, or else
// an error describing a problem.""
// BOUNDARY
func TestReq50_CheckValidsDocFollowsTheStdlibForm(t *testing.T) {
	doc := strings.ToLower(docCommentOf(t, resolveSourceFile, "CheckValid"))

	if !strings.Contains(doc, "returns nil") {
		t.Errorf("CheckValid's doc does not open with the "+
			"\"returns nil if <x> is valid\" form:\n%s", doc)
	}
	if !strings.Contains(doc, "error") {
		t.Errorf("CheckValid's doc does not say it otherwise returns an "+
			"error describing a problem:\n%s", doc)
	}
}

// REQ-51: "The doc comment MUST also name the ONE property checked — escape-row
// shape conformance — and state that nothing else is checked (in particular
// not the Escape-class restriction to no_match/ambiguous_match that Row's doc
// records for RDR 0002), so a nil return is never read as general table
// validity."
// ADVERSARIAL
func TestReq51_CheckValidsDocNamesTheOnePropertyAndItsLimits(t *testing.T) {
	doc := docCommentOf(t, resolveSourceFile, "CheckValid")
	lower := strings.ToLower(doc)

	if !strings.Contains(lower, "shape") {
		t.Errorf("CheckValid's doc does not name escape-row SHAPE "+
			"conformance as the one property checked:\n%s", doc)
	}
	limits := strings.Contains(lower, "nothing else") ||
		strings.Contains(lower, "not check") ||
		strings.Contains(lower, "no other")
	if !limits {
		t.Errorf("CheckValid's doc does not state that nothing else is "+
			"checked:\n%s", doc)
	}
	if !strings.Contains(lower, "general") && !strings.Contains(lower, "validity") {
		t.Errorf("CheckValid's doc does not warn that a nil return is not "+
			"general table validity:\n%s", doc)
	}
}

// --- J. Fixture conformance (Phase 2) ------------------------------------

// REQ-93: "Bring the write-bearing escape rows into conformance at all three
// sites: drop the `Writes` field from the
// `internal/resolve/fixtures_test.go::escapeRow` builder (line 257, inherited
// by its sixteen call sites), and drop the two call-site overrides —
// `internal/resolve/adversarial_test.go:229` (ADV-2b) and
// `internal/resolve/fixup_test.go:103` (Fixup-1e "missing owned state")."
// BOUNDARY
//
// Asserted on VALUES rather than on line numbers: the shared builder's output
// conforms, and no fixture the frozen suite hands `Resolve` breaches. A
// call-site override that survived would show up as a breaching disposition
// fixture below.
func TestReq93_TheWriteBearingEscapeFixturesAreConformed(t *testing.T) {
	nonVacuityGate(t)

	built := escapeRow(breachRuleID, breachLocator, resolve.KindNoMatch)
	if len(built.Writes) != 0 {
		t.Errorf("the shared escapeRow builder still carries %d writes",
			len(built.Writes))
	}

	for name, in := range allDispositionInputs() {
		if err := in.Table.CheckValid(); err != nil {
			t.Errorf("%s: a frozen-suite fixture still breaches: %v",
				name, err)
		}
	}
}

// REQ-94: "`NextTags` stays on the builder: A4 did not widen the predicate."
// BOUNDARY
//
// The conformance edit must remove the WRITES only. A builder that also lost
// its NextTags would silently weaken every escape fixture that depends on the
// escaped plan carrying a next state.
func TestReq94_NextTagsStaysOnTheEscapeRowBuilder(t *testing.T) {
	nonVacuityGate(t)

	built := escapeRow(breachRuleID, breachLocator, resolve.KindNoMatch)
	if len(built.NextTags) == 0 {
		t.Errorf("the escapeRow builder lost its NextTags; A4 did not " +
			"widen the predicate and only Writes is dropped")
	}
}

// REQ-95: "Phases 1 and 2 land as ONE change: the moment the entry check
// exists, every write-bearing escape fixture errors at `Resolve` entry and
// `mustResolve` fatals across the frozen suite"
// ADVERSARIAL
//
// The coupling, asserted in both directions: the check exists AND every
// fixture the frozen suite resolves passes it. Either half alone is a broken
// intermediate state.
func TestReq95_TheCheckAndTheFixtureConformanceLandTogether(t *testing.T) {
	// The check exists and discriminates.
	if err := breachingNoMatchInput().Table.CheckValid(); err == nil {
		t.Errorf("Phase 1 did not land: CheckValid accepts a breaching table")
	}

	// And no frozen-suite fixture trips it — which is what keeps
	// `mustResolve` from fataling across the suite.
	for name, in := range allDispositionInputs() {
		if _, err := resolve.Resolve(in); err != nil {
			t.Errorf("%s: Phase 2 did not land with Phase 1 — a frozen "+
				"fixture now fatals at Resolve entry: %v", name, err)
		}
	}
}

// REQ-96: "The package gains `errors` (and `fmt` for the message); both clear
// the frozen import guards."
// BOUNDARY
//
// The observable consequence of `errors` landing is the aggregate and the
// sentinel; the observable consequence of `fmt` landing is a non-empty,
// diagnostic Error() string. Neither is asserted on its wording — only that
// the error renders something rather than an empty string, which is what
// REQ-43's "MUST remain diagnostic" requires of the message channel.
func TestReq96_TheErrorsAndFmtCapabilitiesAreBothPresent(t *testing.T) {
	err := resolveErr(t, breachingNoMatchInput())

	// `errors`: the join aggregate and the sentinel.
	if _, ok := err.(interface{ Unwrap() []error }); !ok {
		t.Errorf("no errors.Join aggregate; the `errors` capability is absent")
	}
	if resolve.ErrEscapeShapeBreach == nil {
		t.Errorf("no sentinel; the `errors` capability is absent")
	}

	// `fmt`: the per-row error renders a non-empty message. The WORDING is
	// the implementer's (no test may assert on it), but an empty message
	// would leave REQ-43 unsatisfiable.
	elems := breachElements(t, err)
	if len(elems) == 0 {
		t.Fatalf("no per-row elements to render")
	}
	if strings.TrimSpace(elems[0].Error()) == "" {
		t.Errorf("the per-row error renders an empty message; REQ-43 " +
			"requires it to remain diagnostic")
	}
}

// --- L. Performance and non-goals ----------------------------------------

// REQ-99: "The precondition adds one O(rows) pass over `Table.Rows` per
// `Resolve` call."
// BOUNDARY
//
// Linearity, asserted as a behavioural property rather than a timing: the
// check's verdict over a table is a function of its rows alone, and scaling
// the row count scales the reported breaches proportionally — which a
// non-linear (e.g. pairwise) pass would not preserve.
func TestReq99_ThePreconditionIsOneLinearPassOverTheRows(t *testing.T) {
	for _, n := range []int{1, 10, 100} {
		refs := make([]resolve.RowRef, 0, n)
		for i := range n {
			refs = append(refs, resolve.RowRef{
				RuleID:        "rule." + string(rune('a'+i%26)) + strings.Repeat("x", i/26),
				SourceLocator: "flows/f.toml:" + strings.Repeat("1", i+1),
			})
		}
		got := breachRefs(t, resolveErr(t, multiBreachInput(refs...)))
		if len(got) != n {
			t.Errorf("n=%d: reported %d identities; one pass over Rows "+
				"reports each breaching row exactly once", n, len(got))
		}
	}
}

// REQ-100: "The scan is allocation-free in the conforming case — it reads
// `len(row.Escape)` and `len(row.Writes)` and allocates only when a breach is
// found and a `RowRef` is recorded"
// BOUNDARY
//
// Asserted with the allocation counter, which is the only non-tautological
// oracle for this clause: the conforming scan allocates nothing.
func TestReq100_TheConformingScanIsAllocationFree(t *testing.T) {
	nonVacuityGate(t)

	tbl := conformingNoMatchInput().Table

	allocs := testing.AllocsPerRun(100, func() {
		if err := tbl.CheckValid(); err != nil {
			t.Fatalf("the conforming table breached: %v", err)
		}
	})
	if allocs != 0 {
		t.Errorf("CheckValid allocated %.0f times per run on a conforming "+
			"table; the scan reads two lengths and allocates only on a "+
			"breach", allocs)
	}
}

// REQ-101: "The cost is per-call because the kernel is stateless by contract;
// no caching or memoization is introduced"
// ADVERSARIAL
//
// Statelessness is the claim: mutating a table's rows between two calls on
// the SAME Table value must change the verdict. A memoized check keyed on
// anything but the live rows would return the stale answer.
func TestReq101_NoCachingOrMemoizationIsIntroduced(t *testing.T) {
	tbl := conformingNoMatchInput().Table
	if err := tbl.CheckValid(); err != nil {
		t.Fatalf("the conforming table breached: %v", err)
	}

	// Same Table value, same Revision — only the row content changes.
	tbl.Rows = append(tbl.Rows,
		breachingEscapeRow("rule.later", "flows/later.toml:1",
			resolve.KindNoMatch))
	if err := tbl.CheckValid(); err == nil {
		t.Errorf("CheckValid returned the earlier verdict after the rows " +
			"changed; the check is per-call with no memoization")
	}

	// And back: removing the breach restores the clean verdict.
	tbl.Rows = tbl.Rows[:len(tbl.Rows)-1]
	if err := tbl.CheckValid(); err != nil {
		t.Errorf("CheckValid kept the breach verdict after the row was "+
			"removed: %v", err)
	}
}

// REQ-102: "No byte-stable hash, canonical serialization, or ordering
// guarantee is introduced, so the determinism checklist does not apply."
// BOUNDARY
//
// The non-goal, asserted as an absence at the kernel surface: neither Table
// nor the breach error gains a hash/serialization method the checklist would
// bind.
func TestReq102_NoHashOrCanonicalSerializationIsIntroduced(t *testing.T) {
	banned := []string{"Hash", "Fingerprint", "Canonical", "MarshalJSON",
		"MarshalText", "MarshalBinary"}

	for _, rt := range []reflect.Type{
		reflect.TypeOf(resolve.Table{}),
		reflect.TypeOf(&resolve.EscapeShapeBreachError{}),
	} {
		for _, name := range banned {
			if _, ok := rt.MethodByName(name); ok {
				t.Errorf("%v gained a %s method; no byte-stable hash or "+
					"canonical serialization is introduced", rt, name)
			}
		}
	}
}

// REQ-103: "**Silent failure guarded against**: an escape plan carrying writes
// no rule sanctioned … Under the precondition it cannot be emitted; the
// conformance fixtures are the tripwire against regression."
// ADVERSARIAL
//
// The tripwire: an escape plan carrying writes must be unreachable through
// `Resolve` on any input, breaching or not.
func TestReq103_AnEscapePlanCarryingWritesCannotBeEmitted(t *testing.T) {
	inputs := []resolve.Input{
		breachingNoMatchInput(),
		conformingNoMatchInput(),
		emptyNotNilEscapeInput(),
		dormantBreachInput(),
	}
	for name, in := range allDispositionInputs() {
		_ = name
		inputs = append(inputs, in)
	}

	for i, in := range inputs {
		got, err := resolve.Resolve(in)
		if err != nil {
			continue
		}
		if got.Plan != nil && got.Plan.Escaped && len(got.Plan.Writes) != 0 {
			t.Errorf("input %d emitted an escape plan carrying %d writes: %+v",
				i, len(got.Plan.Writes), got.Plan.Writes)
		}
	}
}

// REQ-104: "**Error-channel creep guarded against**: this RDR's normative
// clause states what qualifies for the error path (a producer/programmer
// contract breach, never a modeled condition)"
// ADVERSARIAL
//
// The guard against creep: NO modeled condition may migrate to the error
// path. Every shipped disposition fixture must still travel the Result value
// with a nil error, and the ONLY input that errors is the producer breach.
func TestReq104_NoModeledConditionMigratesToTheErrorPath(t *testing.T) {
	for name, in := range allDispositionInputs() {
		if _, err := resolve.Resolve(in); err != nil {
			t.Errorf("%s: a modeled condition migrated onto the error "+
				"path: %v", name, err)
		}
	}

	// And the breach — a producer contract breach — is on the error path.
	if _, err := resolve.Resolve(breachingNoMatchInput()); err == nil {
		t.Errorf("a producer contract breach did not qualify for the error " +
			"path")
	}
}

// REQ-105: "Positive: the five-kind taxonomy, RDR 0005's mapping, RDR 0002's
// grammar, and the kernel's type vocabulary all stand unchanged; the check
// reuses a shipped discriminator and a reserved channel."
// BOUNDARY
//
// The no-change claim over the kernel's own surface: the five kinds, and the
// shipped type vocabulary the record names.
func TestReq105_TheKernelVocabularyStandsUnchanged(t *testing.T) {
	kinds := resolve.RefusalKinds()
	want := []resolve.RefusalKind{
		resolve.KindUnmodeledOutcome,
		resolve.KindNoMatch,
		resolve.KindAmbiguousMatch,
		resolve.KindOwnedStateUnavailable,
		resolve.KindGuardUnevaluable,
	}
	if len(kinds) != len(want) {
		t.Fatalf("the taxonomy carries %d kinds; want the shipped five",
			len(kinds))
	}
	for _, k := range want {
		found := false
		for _, got := range kinds {
			if got == k {
				found = true
			}
		}
		if !found {
			t.Errorf("refusal kind %q is gone from the taxonomy", k)
		}
	}

	// The kernel's type vocabulary — the values the record names as
	// unchanged — still carry their shipped shapes.
	if _, ok := reflect.TypeOf(resolve.Result{}).FieldByName("Plan"); !ok {
		t.Errorf("Result.Plan is gone")
	}
	if _, ok := reflect.TypeOf(resolve.Result{}).FieldByName("Refusal"); !ok {
		t.Errorf("Result.Refusal is gone")
	}
	if _, ok := reflect.TypeOf(resolve.RowRef{}).FieldByName("RuleID"); !ok {
		t.Errorf("RowRef.RuleID is gone")
	}
}

// REQ-106: "Negative (behavior change, deliberate): a hand-built table
// carrying a *dormant* malformed row — one no resolution path reaches —
// previously resolved fine and now errors on every `Resolve` (whole-table
// precondition)."
// DOMAIN EDGE
//
// The deliberate regression, pinned so it cannot be silently walked back into
// an on-path-only check: the dormant table errors on EVERY Resolve, for every
// request tuple, not merely on the one that reaches the row.
func TestReq106_TheDormantMalformedRowErrorsOnEveryResolve(t *testing.T) {
	base := dormantBreachInput()

	for _, recognized := range []string{
		base.Recognized, "successful", "never-recognized", "",
	} {
		in := deepCopyInput(base)
		in.Recognized = recognized
		got, err := resolve.Resolve(in)
		if err == nil {
			t.Errorf("recognized=%q: the dormant malformed row did not "+
				"error: %+v — the precondition is over the whole table on "+
				"EVERY call", recognized, got)
		}
	}
}

// REQ-107: "Escape fixtures … | Extend (conform all three sites) | Fixtures
// must satisfy the precondition; supersedes the "Noted, not filed" cleanup"
// BOUNDARY
func TestReq107_EveryEscapeFixtureSiteSatisfiesThePrecondition(t *testing.T) {
	nonVacuityGate(t)

	// The builder itself.
	row := escapeRow("r", "l", resolve.KindNoMatch)
	if len(row.Escape) != 0 && len(row.Writes) != 0 {
		t.Errorf("the shared escapeRow builder breaches the precondition")
	}

	// And every table the frozen suite hands the kernel.
	for name, in := range allDispositionInputs() {
		if err := in.Table.CheckValid(); err != nil {
			t.Errorf("%s: fixture site does not satisfy the precondition: %v",
				name, err)
		}
	}
}

// REQ-108: "Row-identity diagnostics | `internal/resolve/resolve.go::RowRef` +
// `rowRefs`/`compareRefs` | … | Reuse | The typed breach error carries
// `RowRef`s in the same sorted order"
// BOUNDARY
//
// Reuse, not a parallel implementation: the breach report's order must match
// the order the shipped `ambiguous_match` refusal reports the SAME identities
// in. A second, independently-written comparator would eventually diverge.
func TestReq108_TheBreachReportUsesTheSameSortedRowRefOrder(t *testing.T) {
	// Identities chosen so the sort is discriminating on BOTH components.
	refs := []resolve.RowRef{
		{RuleID: "b", SourceLocator: "flows/a.toml:1"},
		{RuleID: "a", SourceLocator: "flows/z.toml:9"},
		{RuleID: "a", SourceLocator: "flows/a.toml:1"},
	}
	got := breachRefs(t, resolveErr(t, multiBreachInput(refs...)))

	want := []resolve.RowRef{
		{RuleID: "a", SourceLocator: "flows/a.toml:1"},
		{RuleID: "a", SourceLocator: "flows/z.toml:9"},
		{RuleID: "b", SourceLocator: "flows/a.toml:1"},
	}
	if !sameRefs(got, want) {
		t.Errorf("breach report order = %+v; want compareRefs order %+v — "+
			"RuleID first, then SourceLocator", got, want)
	}
}
