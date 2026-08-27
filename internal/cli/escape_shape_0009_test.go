package cli

// RDR 0009 `0009:C7` — surfacing an escape-row shape breach through the CLI.
//
// The defect this file pins against lives at `flow_resolve.go`'s post-Resolve
// error branch, which today re-codes EVERY kernel error as
// `codeAccessorFailed` and folds the row identities into `err.Error()` prose.
// C7's obligation is to discriminate the breach out of that branch: a stable
// `Code`, `GroupInternal`, the identities on the SERIALIZED envelope, and a
// remedy `Hint`.
//
// How the breach is reached. The authored path cannot produce one — RDR
// 0002's normalizer rejects an escape rule carrying a write block, a clear
// list, or even an empty one (REQ-11/REQ-12), which is exactly the layering
// this RDR describes. A breach therefore reaches `Resolve` only from a
// NON-TOML producer, so the CLI-side obligation is exercised with the real
// kernel error, produced by the real `resolve.Resolve` over a hand-built
// table, handed to the verb's real classifier. No mock stands in for either.
//
// Every assertion below is on the typed `Code` and the structured `Findings`,
// never on the exit code alone: an unwrapped kernel error also exits 2
// (`ExecuteAndEmit` -> `cobraErrorToCLIError`, `command-error`), so the exit
// code cannot detect a missing wrap. The stable code can.
//
// The unit under test is `kernelResolveFailure(error) *clierr.CLIError`: the
// verb-layer classifier `flow_resolve.go`'s post-`Resolve` error branch must
// route through, sibling to the shipped `kernelFailure(resolve.Refusal)`. It
// is the smallest surface that carries the whole C7 obligation and is
// reachable from a test, since the authored path cannot construct a breach.
// Naming it here is the test's binding on the implementation, in the same way
// `0009:C4` binds the kernel's four spellings.

import (
	"encoding/json"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/resolve"
)

// The spellings `0009:C7` fixes verbatim.
const (
	// wantBreachCode is fixed by the contract's literal. Note it is NOT
	// `flow`-prefixed, unlike the neighbouring vocabulary in flow_input.go;
	// the contract's literal wins (see coverage.md, Q4).
	wantBreachCode = "escape-row-shape-breach"
	// wantBreachHint is the remedy TS scenario 11 fixes verbatim.
	wantBreachHint = "fix the table producer: an escape row must carry no writes"
)

// breachedKernelError produces a REAL escape-row shape breach by running the
// real kernel over a hand-built table carrying the given breaching
// identities. Nothing is mocked: the error under test is the one a non-TOML
// producer would provoke in production.
func breachedKernelError(t *testing.T, refs ...resolve.RowRef) error {
	t.Helper()

	tbl := resolve.Table{
		Revision: "escape-shape-0009",
		Outcomes: []string{"advance"},
	}
	for _, ref := range refs {
		tbl.Rows = append(tbl.Rows, resolve.Row{
			RuleID:        ref.RuleID,
			SourceLocator: ref.SourceLocator,
			Outcome:       "advance",
			Escape:        []resolve.RefusalKind{resolve.KindNoMatch},
			Writes:        []resolve.Tag{{Key: "status", Value: "Escaped"}},
		})
	}

	_, err := resolve.Resolve(resolve.Input{
		Flow:       "escapeshape",
		Table:      tbl,
		Recognized: "advance",
		Guards:     guardSeam(),
	})
	if err == nil {
		t.Fatalf("the kernel accepted a breaching table; there is no breach " +
			"to surface")
	}
	if !errors.Is(err, resolve.ErrEscapeShapeBreach) {
		t.Fatalf("the kernel errored for some other reason: %v", err)
	}
	return err
}

// theBreachRef is the identity used wherever one suffices.
var theBreachRef = resolve.RowRef{
	RuleID:        "rdr.escape.needsowned",
	SourceLocator: "flows/rdr.toml:90",
}

// wrapped runs the verb layer's classifier over a real kernel breach and
// requires it to have produced a CLIError.
func wrapped(t *testing.T, err error) *clierr.CLIError {
	t.Helper()

	ce := kernelResolveFailure(err)
	if ce == nil {
		t.Fatalf("the verb layer did not wrap the kernel error into a "+
			"*clierr.CLIError: %v", err)
	}
	return ce
}

// REQ-58: "When a breach reaches the CLI, the verb MUST wrap it into a
// *clierr.CLIError carrying a stable Code, Group GroupInternal (exit 2), the
// offending row identity, and a Hint stating the remedy"
// HAPPY PATH
//
// All four halves at once, asserted on the Code — never on the exit code
// alone, which an unwrapped kernel error would also satisfy.
func TestReq58_ABreachReachingTheCLIIsWrappedWithCodeGroupIdentityAndHint(t *testing.T) {
	ce := wrapped(t, breachedKernelError(t, theBreachRef))

	if ce.Code != wantBreachCode {
		t.Errorf("code = %q; want %q — the breach must be discriminated out "+
			"of the generic %q branch", ce.Code, wantBreachCode,
			codeAccessorFailed)
	}
	if ce.Group != clierr.GroupInternal {
		t.Errorf("group = %v; want GroupInternal", ce.Group)
	}
	if clierr.ExitCodeFor(ce) != 2 {
		t.Errorf("exit = %d; want 2", clierr.ExitCodeFor(ce))
	}
	if ce.Hint == "" {
		t.Errorf("the wrap carries no Hint stating the remedy")
	}

	// The offending row identity, carried structurally.
	if len(ce.Findings) != 1 {
		t.Fatalf("Findings carries %d entries; want the one offending row",
			len(ce.Findings))
	}
	f := ce.Findings[0]
	if f.Locator != theBreachRef.SourceLocator {
		t.Errorf("finding locator = %q; want the source locator %q",
			f.Locator, theBreachRef.SourceLocator)
	}
	if f.Rule != theBreachRef.RuleID && f.Param != theBreachRef.RuleID {
		t.Errorf("neither finding.rule (%q) nor finding.param (%q) carries "+
			"the offending rule id %q", f.Rule, f.Param, theBreachRef.RuleID)
	}
}

// REQ-59: "Because clierr.CLIError.Cause is json:"-", the Go error chain that
// carries the RowRef values is NOT wire-visible: the offending row identities
// MUST reach the envelope through a serialized field."
// ADVERSARIAL
//
// The discriminating assertion is on the SERIALIZED bytes: a wrap that
// carried the identities only on `Cause` would satisfy every Go-side
// assertion above and still lose them at the wire.
func TestReq59_TheRowIdentitiesReachTheSerializedEnvelope(t *testing.T) {
	ce := wrapped(t, breachedKernelError(t, theBreachRef))

	var buf strings.Builder
	clierr.EmitJSON(&buf, ce)
	line := buf.String()

	var env map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &env); err != nil {
		t.Fatalf("the envelope is not one JSON object: %v\n%s", err, line)
	}
	if _, present := env["cause"]; present {
		t.Errorf("the envelope serialized `cause`; it is json:\"-\"")
	}

	if !strings.Contains(line, theBreachRef.SourceLocator) {
		t.Errorf("the serialized envelope does not carry the source locator "+
			"%q:\n%s", theBreachRef.SourceLocator, line)
	}
	if !strings.Contains(line, theBreachRef.RuleID) {
		t.Errorf("the serialized envelope does not carry the rule id %q:\n%s",
			theBreachRef.RuleID, line)
	}

	// And the identities are in a STRUCTURED field, not folded into prose.
	raw, ok := env["findings"]
	if !ok {
		t.Fatalf("the envelope carries no `findings` array:\n%s", line)
	}
	var findings []map[string]any
	if err := json.Unmarshal(raw, &findings); err != nil {
		t.Fatalf("`findings` is not an array of objects: %v", err)
	}
	if len(findings) != 1 {
		t.Errorf("`findings` carries %d entries; want 1", len(findings))
	}
}

// REQ-60: "The chosen carrier is a NEW omitempty field added under the type's
// own "Extend with new optional fields as needed" allowance — NOT the existing
// Detail"
// BOUNDARY
//
// The carrier already ships as `CLIError.Findings` (RDR 0008). The live
// obligation is to USE it rather than `Detail`, and to mint no second
// parallel identity field beside it.
func TestReq60_TheCarrierIsFindingsAndNotDetail(t *testing.T) {
	ce := wrapped(t, breachedKernelError(t, theBreachRef))

	if len(ce.Findings) == 0 {
		t.Errorf("the identities were not carried on `Findings`")
	}
	if strings.Contains(ce.Detail, theBreachRef.SourceLocator) ||
		strings.Contains(ce.Detail, theBreachRef.RuleID) {
		t.Errorf("the identities were folded into `Detail` (%q); reading "+
			"them back out would require re-parsing rendered prose", ce.Detail)
	}

	// The carrier stays `omitempty` and no SECOND identity list was minted
	// beside it.
	rt := reflect.TypeOf(clierr.CLIError{})
	f, ok := rt.FieldByName("Findings")
	if !ok {
		t.Fatalf("clierr.CLIError has no Findings field")
	}
	if !strings.Contains(string(f.Tag), "omitempty") {
		t.Errorf("Findings is not omitempty: %q", f.Tag)
	}
	for _, banned := range []string{"Rows", "RowRefs", "Identities", "Refs"} {
		if _, found := rt.FieldByName(banned); found {
			t.Errorf("clierr.CLIError gained a second identity carrier %q "+
				"beside Findings", banned)
		}
	}
}

// REQ-61: "The clause binds the carrier CLASS only; the field's NAME stays the
// implementer's."
// BOUNDARY
//
// Resolved by HEAD: the name is `Findings`. What the clause still binds is the
// CLASS — a serialized, structured, repeated identity carrier — so the test
// asserts the class properties rather than re-litigating the name.
func TestReq61_TheCarrierIsAStructuredRepeatedSerializedField(t *testing.T) {
	f, ok := reflect.TypeOf(clierr.CLIError{}).FieldByName("Findings")
	if !ok {
		t.Fatalf("clierr.CLIError has no Findings field")
	}
	if f.Type.Kind() != reflect.Slice {
		t.Errorf("the carrier is %v; the class is a REPEATED field so a "+
			"multi-breach report needs no second envelope", f.Type)
	}
	if f.Type.Elem().Kind() != reflect.Struct {
		t.Errorf("the carrier's element is %v; the class is STRUCTURED, "+
			"never rendered prose", f.Type.Elem())
	}
	if strings.HasPrefix(string(f.Tag.Get("json")), "-") {
		t.Errorf("the carrier is not serialized: %q", f.Tag)
	}
}

// REQ-62: "the field's type MUST be a clierr-local representation (plain
// strings or a small row-identity struct declared in clierr) and MUST NOT be
// resolve.RowRef, keeping clierr a leaf that does not import internal/resolve
// — the CLI verb layer, which already imports both, does the conversion"
// ADVERSARIAL
//
// Two obligations: the carrier's identity fields are plain strings (never
// `resolve.RowRef`), and `clierr` gains no `internal/resolve` import. The
// import guard is the load-bearing half — a `RowRef`-typed field would compile
// and pass every value assertion while breaking the leaf property.
func TestReq62_ClierrStaysALeafAndTheConversionIsTheVerbLayers(t *testing.T) {
	rt := reflect.TypeOf(clierr.Finding{})
	rowRef := reflect.TypeOf(resolve.RowRef{})
	for i := range rt.NumField() {
		if rt.Field(i).Type == rowRef {
			t.Errorf("clierr.Finding field %q is resolve.RowRef; the carrier "+
				"must be a clierr-local representation", rt.Field(i).Name)
		}
	}
	for _, name := range []string{"Rule", "Param", "Locator"} {
		f, ok := rt.FieldByName(name)
		if !ok {
			t.Fatalf("clierr.Finding has no %s field", name)
		}
		if f.Type.Kind() != reflect.String {
			t.Errorf("clierr.Finding.%s is %v; the identity fields are "+
				"plain strings", name, f.Type)
		}
	}

	// The leaf property, read off the shipped source rather than inferred.
	src := readRepoFile(t, "internal/cli/clierr/clierr.go")
	if strings.Contains(src, "internal/resolve") {
		t.Errorf("internal/cli/clierr imports internal/resolve; clierr must " +
			"stay a leaf and the verb layer does the conversion")
	}

	// And the conversion really did happen: the verb layer produced plain
	// strings from the kernel's RowRef values.
	ce := wrapped(t, breachedKernelError(t, theBreachRef))
	if len(ce.Findings) == 0 {
		t.Fatalf("no findings to inspect")
	}
	if ce.Findings[0].Locator != theBreachRef.SourceLocator {
		t.Errorf("the RowRef.SourceLocator was not converted into the "+
			"finding's Locator: %q vs %q",
			ce.Findings[0].Locator, theBreachRef.SourceLocator)
	}
}

// REQ-63: "the per-identity Count MUST serialize alongside the identities"
// BOUNDARY
//
// The count must reach the WIRE, structurally. Exercised on the degenerate
// producer — three rows sharing one identity — where the count is the only
// diagnostic content the report has left (REQ-43). Asserted on the serialized
// bytes, so a count carried only in Go memory fails.
func TestReq63_ThePerIdentityCountSerializesAlongsideTheIdentities(t *testing.T) {
	zero := resolve.RowRef{}
	ce := wrapped(t, breachedKernelError(t, zero, zero, zero))

	if len(ce.Findings) != 1 {
		t.Fatalf("Findings carries %d entries; three rows sharing one "+
			"identity collapse to ONE", len(ce.Findings))
	}

	var buf strings.Builder
	clierr.EmitJSON(&buf, ce)

	var env struct {
		Findings []map[string]any `json:"findings"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())),
		&env); err != nil {
		t.Fatalf("the envelope is not one JSON object: %v", err)
	}
	if len(env.Findings) != 1 {
		t.Fatalf("the serialized envelope carries %d findings; want 1",
			len(env.Findings))
	}

	raw, present := env.Findings[0]["count"]
	if !present {
		t.Fatalf("the serialized finding carries no `count`; the "+
			"per-identity count must reach the wire alongside the "+
			"identities, never only as prose: %v", env.Findings[0])
	}
	n, ok := raw.(float64)
	if !ok {
		t.Fatalf("`count` is %T; want a JSON number", raw)
	}
	if int(n) != 3 {
		t.Errorf("`count` = %d; want 3 pre-collapse rows", int(n))
	}
}

// REQ-64: "The stable Code is "escape-row-shape-breach"."
// BOUNDARY
//
// Byte-for-byte, and NOT `flow`-prefixed. The contract's literal wins over
// the neighbouring vocabulary's convention.
func TestReq64_TheStableCodeIsSpelledExactly(t *testing.T) {
	ce := wrapped(t, breachedKernelError(t, theBreachRef))

	if ce.Code != "escape-row-shape-breach" {
		t.Errorf("code = %q; the contract fixes the literal "+
			"\"escape-row-shape-breach\"", ce.Code)
	}
	if strings.HasPrefix(ce.Code, "flow-") {
		t.Errorf("code = %q carries the `flow-` prefix; the contract's "+
			"literal is unprefixed", ce.Code)
	}
}

// REQ-65: "the Cause comment's now-stale second clause MUST be amended in the
// same change (a Prerequisite; amending a stale code comment is not reopening
// RDR 0005's envelope contract)."
// BOUNDARY
//
// The shipped comment claims "the wire-visible cause surface is Detail",
// which `Findings` already makes false at HEAD and which this RDR's own
// identity carrier makes decisively false. Asserted on the shipped source.
func TestReq65_TheStaleCauseCommentIsAmended(t *testing.T) {
	doc := causeFieldComment(t)
	if doc == "" {
		t.Fatalf("CLIError.Cause carries no doc comment")
	}

	lower := strings.ToLower(doc)
	if strings.Contains(lower, "the wire-visible cause surface is detail") {
		t.Errorf("CLIError.Cause's doc still claims Detail is THE "+
			"wire-visible cause surface; `findings` carries the row "+
			"identities and the sentence is false:\n%s", doc)
	}
	if !strings.Contains(lower, "findings") {
		t.Errorf("the amended Cause comment does not name `findings` as a "+
			"wire-visible surface:\n%s", doc)
	}
}

// REQ-66: "The same change that adds the `omitempty` identity field to
// `clierr.CLIError` amends the now-stale second clause of
// `internal/cli/clierr/clierr.go::CLIError.Cause`'s doc comment ("Not
// serialized — the wire-visible cause surface is Detail"), which the addition
// makes false, and updates `docs/cli-output-contract.md`'s error-envelope
// field list alongside."
// BOUNDARY
//
// The doc half: the output contract's error-envelope field list must name
// `findings` AND the count key the identities travel with, so a consumer
// reading the contract can decode what the CLI emits.
func TestReq66_TheOutputContractDocumentsTheIdentityCarrier(t *testing.T) {
	doc := readRepoFile(t, "docs/cli-output-contract.md")

	if !strings.Contains(doc, "findings") {
		t.Errorf("docs/cli-output-contract.md does not document `findings` " +
			"on the error envelope")
	}
	if !strings.Contains(doc, "count") {
		t.Errorf("docs/cli-output-contract.md does not document the " +
			"per-identity `count` the row identities serialize alongside")
	}
}

// REQ-67: "conformance MUST be asserted on the Code, never on the exit code
// alone."
// ADVERSARIAL
//
// The claim is that the exit code cannot discriminate a missing wrap. Pinned
// by showing the two are OBSERVABLY different on the Code and IDENTICAL on
// the exit: an unwrapped kernel error also reaches exit 2.
func TestReq67_TheExitCodeCannotDiscriminateAMissingWrap(t *testing.T) {
	raw := breachedKernelError(t, theBreachRef)
	ce := wrapped(t, raw)

	unwrappedExit := clierr.ExitCodeFor(cobraErrorToCLIError(raw))
	wrappedExit := clierr.ExitCodeFor(ce)
	if unwrappedExit != wrappedExit {
		t.Logf("the unwrapped kernel error exits %d and the wrap exits %d; "+
			"the exit code happens to discriminate here, but the Code is "+
			"still the contract's oracle", unwrappedExit, wrappedExit)
	}

	// The Code always discriminates.
	if got := clierr.ErrorCode(cobraErrorToCLIError(raw)); got == wantBreachCode {
		t.Errorf("an UNWRAPPED kernel error already carries the breach "+
			"code %q; the code would then not discriminate the wrap", got)
	}
	if ce.Code != wantBreachCode {
		t.Errorf("the wrapped error's code = %q; want %q", ce.Code,
			wantBreachCode)
	}
}

// REQ-68: "Adding this Code does NOT open RDR 0001's refusal taxonomy: the
// five refusal kinds stay closed and gain no member"
// ADVERSARIAL
func TestReq68_TheRefusalTaxonomyStaysClosed(t *testing.T) {
	kinds := resolve.RefusalKinds()
	if len(kinds) != 5 {
		t.Fatalf("the taxonomy carries %d kinds (%v); the five stay closed",
			len(kinds), kinds)
	}
	for _, k := range kinds {
		if string(k) == wantBreachCode {
			t.Errorf("the CLI code %q was minted as a refusal kind; clierr "+
				"codes are a separate, extensible surface", k)
		}
	}

	// And the breach's own CLIError is not a refusal: it is GroupInternal,
	// the CLI's own-invariant bucket, never a modeled refusal's group.
	ce := wrapped(t, breachedKernelError(t, theBreachRef))
	if ce.Group != clierr.GroupInternal {
		t.Errorf("group = %v; a breach is the producer's programmer mistake, "+
			"not a modeled refusal", ce.Group)
	}
}

// REQ-69: "This is an additive envelope change, not a change to RDR 0005's
// refusal-to-exit-code mapping"
// BOUNDARY
//
// The mapping is unchanged: every shipped group still maps to the exit it
// mapped to before, and the new code rides GroupInternal's existing exit 2
// rather than introducing an exit of its own.
func TestReq69_TheRefusalToExitCodeMappingIsUnchanged(t *testing.T) {
	mapping := map[clierr.ErrorGroup]int{
		clierr.GroupSuccess:        0,
		clierr.GroupWarning:        0,
		clierr.GroupUserEnv:        2,
		clierr.GroupInternal:       2,
		clierr.GroupEnvUnavailable: 3,
		clierr.GroupSignalCancel:   130,
	}
	for group, want := range mapping {
		got := clierr.ExitCodeFor(&clierr.CLIError{Code: "x", Group: group})
		if got != want {
			t.Errorf("group %v maps to exit %d; RDR 0005 fixes %d",
				group, got, want)
		}
	}

	ce := wrapped(t, breachedKernelError(t, theBreachRef))
	if clierr.ExitCodeFor(ce) != 2 {
		t.Errorf("the breach exits %d; it rides GroupInternal's existing "+
			"exit 2 and introduces none of its own", clierr.ExitCodeFor(ce))
	}
}

// REQ-109: "A breach surfaced through the CLI. **Expected**: exit 2 via
// `CLIError{Group: GroupInternal}` carrying `Code: "escape-row-shape-breach"`,
// the offending row identity, and the remedy `Hint` "fix the table producer:
// an escape row must carry no writes" — asserted on the `Code`"
// HAPPY PATH
//
// The end-to-end C7 assertion. Every element of the expectation is pinned,
// including the Hint VERBATIM — the record fixes its wording, so unlike the
// kernel's message prose it is a contract string.
func TestReq109_ABreachSurfacedThroughTheCLICarriesCodeIdentityAndHint(t *testing.T) {
	ce := wrapped(t, breachedKernelError(t, theBreachRef))

	if ce.Code != wantBreachCode {
		t.Fatalf("code = %q; want %q — the conformance oracle", ce.Code,
			wantBreachCode)
	}
	if ce.Group != clierr.GroupInternal {
		t.Errorf("group = %v; want GroupInternal", ce.Group)
	}
	if clierr.ExitCodeFor(ce) != 2 {
		t.Errorf("exit = %d; want 2", clierr.ExitCodeFor(ce))
	}
	if ce.Hint != wantBreachHint {
		t.Errorf("hint = %q; the record fixes the remedy verbatim as %q",
			ce.Hint, wantBreachHint)
	}

	// The offending row identity, read off the SERIALIZED envelope — the Go
	// chain that carries the RowRef values is json:"-".
	var buf strings.Builder
	clierr.EmitJSON(&buf, ce)
	var env struct {
		Code     string           `json:"code"`
		Hint     string           `json:"hint"`
		Findings []map[string]any `json:"findings"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())),
		&env); err != nil {
		t.Fatalf("the envelope is not one JSON object: %v", err)
	}
	if env.Code != wantBreachCode {
		t.Errorf("serialized code = %q; want %q", env.Code, wantBreachCode)
	}
	if env.Hint != wantBreachHint {
		t.Errorf("serialized hint = %q; want %q", env.Hint, wantBreachHint)
	}
	if len(env.Findings) != 1 {
		t.Fatalf("serialized findings carries %d entries; want the one "+
			"offending row", len(env.Findings))
	}
	body, _ := json.Marshal(env.Findings[0])
	for _, want := range []string{
		theBreachRef.RuleID, theBreachRef.SourceLocator,
	} {
		if !strings.Contains(string(body), want) {
			t.Errorf("the serialized finding does not carry %q: %s",
				want, body)
		}
	}
}

// REQ-109 (multi-breach half): the CLI carries EVERY offending identity, in
// the kernel's sorted order, not just the first.
// ADVERSARIAL
//
// A single `errors.As` at the verb layer would report one identity out of
// several (REQ-38). This pins the verb layer traversing the aggregate.
func TestReq109_TheCLICarriesEveryOffendingIdentityInKernelOrder(t *testing.T) {
	a := resolve.RowRef{RuleID: "rule.a", SourceLocator: "flows/a.toml:1"}
	b := resolve.RowRef{RuleID: "rule.b", SourceLocator: "flows/b.toml:2"}
	c := resolve.RowRef{RuleID: "rule.c", SourceLocator: "flows/c.toml:3"}

	// Supplied out of identity order, so a position-preserving verb layer
	// would report them out of order.
	ce := wrapped(t, breachedKernelError(t, c, a, b))

	if len(ce.Findings) != 3 {
		t.Fatalf("Findings carries %d entries; the verb must traverse the "+
			"aggregate rather than reporting the first breach alone",
			len(ce.Findings))
	}
	want := []resolve.RowRef{a, b, c}
	for i, f := range ce.Findings {
		if f.Locator != want[i].SourceLocator {
			t.Errorf("finding %d locator = %q; want %q — the kernel's "+
				"compareRefs order is carried through", i, f.Locator,
				want[i].SourceLocator)
		}
		if f.Code != wantBreachCode {
			t.Errorf("finding %d code = %q; want %q", i, f.Code,
				wantBreachCode)
		}
	}
}

// REQ-109 (non-breach half): a kernel error that is NOT a shape breach still
// falls through to the generic branch.
// ADVERSARIAL
//
// The classifier must DISCRIMINATE, not blanket-recode. A verb that answered
// `escape-row-shape-breach` to every kernel error would pass every assertion
// above and mislabel RDR 0008's reserved-key breach.
func TestReq109_ANonBreachKernelErrorStillFallsThroughToTheGenericBranch(t *testing.T) {
	// A real, different kernel programmer-mistake error: RDR 0008's
	// reserved-key producer breach.
	_, err := resolve.Resolve(resolve.Input{
		Flow:       "escapeshape",
		Table:      resolve.Table{Revision: "r", Outcomes: []string{"advance"}},
		Owned:      []resolve.Tag{{Key: "recognized", Value: "advance"}},
		Recognized: "advance",
		Guards:     guardSeam(),
	})
	if err == nil {
		t.Fatalf("the reserved-key producer breach no longer errors")
	}
	if errors.Is(err, resolve.ErrEscapeShapeBreach) {
		t.Fatalf("the reserved-key breach is classified as a shape breach")
	}

	ce := kernelResolveFailure(err)
	if ce == nil {
		t.Fatalf("the verb layer produced no CLIError for a kernel error")
	}
	if ce.Code == wantBreachCode {
		t.Errorf("a non-breach kernel error was coded %q; the classifier "+
			"must discriminate rather than blanket-recode", ce.Code)
	}
	if ce.Code != codeAccessorFailed {
		t.Errorf("code = %q; a non-breach kernel error still falls through "+
			"to %q", ce.Code, codeAccessorFailed)
	}
	if len(ce.Findings) != 0 {
		t.Errorf("the generic branch carries %d findings; the identity "+
			"carrier belongs to the breach branch alone", len(ce.Findings))
	}
}

// causeFieldComment returns the doc comment on `clierr.CLIError.Cause`, read
// off the shipped source. The comment is a contract artifact RDR 0005 owns
// and this RDR's Prerequisite amends, so its text is the assertion subject —
// this is not an assertion on a log or a message string.
func causeFieldComment(t *testing.T) string {
	t.Helper()

	path := filepath.Join(repoRootFor(t), "internal/cli/clierr/clierr.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse clierr.go: %v", err)
	}

	var doc string
	ast.Inspect(f, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok || ts.Name.Name != "CLIError" {
			return true
		}
		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return false
		}
		for _, fld := range st.Fields.List {
			for _, nm := range fld.Names {
				if nm.Name == "Cause" && fld.Doc != nil {
					doc = fld.Doc.Text()
				}
			}
		}
		return false
	})
	return doc
}
