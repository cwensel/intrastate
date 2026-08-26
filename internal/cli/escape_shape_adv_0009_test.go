package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
	"github.com/newcoinc/intrastate/internal/resolve"
)

// RDR 0009 Phase 3b — ADVERSARIAL failure-mode tests.
//
// These are written against the RDR's `## Trade-offs > ### Failure Modes`
// section, not against the Phase 1 suite, and were authored without reading
// it. Each test targets a mode the section names and asserts the property
// the section claims holds.

// adv0009BreachRow builds one breaching escape row: a non-empty Escape list
// together with a non-empty Writes slice.
func adv0009BreachRow(ruleID, locator, writeKey string) resolve.Row {
	return resolve.Row{
		RuleID:        ruleID,
		SourceLocator: locator,
		Outcome:       "adv-outcome",
		Escape:        []resolve.RefusalKind{resolve.KindNoMatch},
		Writes:        []resolve.Tag{{Key: writeKey, Value: "v"}},
	}
}

// adv0009KernelError runs the rows through Resolve so the error under test is
// the one the kernel actually returns on the live path, not a hand-built one.
func adv0009KernelError(t *testing.T, rows ...resolve.Row) error {
	t.Helper()
	_, err := resolve.Resolve(resolve.Input{
		Table:      resolve.Table{Outcomes: []string{"adv-outcome"}, Rows: rows},
		Recognized: "adv-outcome",
		Guards:     guardSeam(),
	})
	if err == nil {
		t.Fatal("Resolve returned nil error for a breaching table")
	}
	return err
}

// --- ADV-1 -------------------------------------------------------------
//
// FAILURE MODE (RDR 0009, Failure Modes, "Visible break (programmatic
// producer)" and "Diagnosis (revised at Resolve)"): the breach report must
// stay machine-readable and DIAGNOSTIC without reading error prose. The
// Normative Contracts make the per-identity count structural precisely
// because collapsing equal identities discards multiplicity, and require
// that for the degenerate identity "the error MUST remain diagnostic in
// that case by stating the breach count alongside the identities."
//
// The CLI's DEFAULT output mode is text: respond.Fail routes ModeJSON to
// EmitJSON and everything else to EmitText. clierr.Finding.Count reaches
// the JSON envelope, but clierr.EmitFindingsText renders a finding through
// Finding.identitySuffix(), whose field list was NOT widened alongside the
// new Count field. So in the operator's default view the multiplicity is
// gone — and for rows carrying no source identity (RowRef{"",""}) the text
// finding carries no identity AND no count, i.e. no diagnostic content at
// all, while the kernel's own Error() string says "3 breaching rows".
//
// This asserts the property the RDR claims: the breach count must be
// recoverable from BOTH serialized output modes.
func TestAdv0009_TheCollapsedBreachCountSurvivesTextRendering(t *testing.T) {
	// Three breaching rows sharing one identity — the RDR's own worked
	// example (three rows, one reported entry, Count == 3).
	err := adv0009KernelError(t,
		adv0009BreachRow("dup-rule", "flow.toml:10", "a"),
		adv0009BreachRow("dup-rule", "flow.toml:10", "b"),
		adv0009BreachRow("dup-rule", "flow.toml:10", "c"),
	)
	ce := kernelResolveFailure(err)

	if len(ce.Findings) != 1 {
		t.Fatalf("three rows sharing one identity must collapse to one finding, got %d",
			len(ce.Findings))
	}
	if ce.Findings[0].Count != 3 {
		t.Fatalf("Count = %d, want 3 (pre-collapse rows)", ce.Findings[0].Count)
	}

	var jsonOut bytes.Buffer
	clierr.EmitJSON(&jsonOut, ce)
	if !strings.Contains(jsonOut.String(), `"count":3`) {
		t.Fatalf("JSON envelope dropped the per-identity count:\n%s", jsonOut.String())
	}

	var textOut bytes.Buffer
	clierr.EmitText(&textOut, ce)
	// The count must be readable from text output too. Layout is free;
	// presence is not. Any rendering that carries the number 3 as the
	// collapsed-row count satisfies this — the assertion is deliberately
	// weak about FORM and strict about PRESENCE.
	if !strings.Contains(textOut.String(), "3") {
		t.Errorf("text rendering dropped the collapsed breach count of 3, so an "+
			"operator in the default output mode cannot tell that three rows "+
			"breached; JSON carries it and text does not:\n%s", textOut.String())
	}
}

// The degenerate case the Normative Contracts single out: rows carrying no
// source identity collapse to RowRef{"",""}, and the report "MUST remain
// diagnostic in that case by stating the breach count alongside the
// identities." In text mode the identity fields are all empty, so
// identitySuffix() renders nothing but the hint and the count is absent —
// the finding line is indistinguishable from a single-row breach.
func TestAdv0009_ADegenerateIdentityBreachIsStillDiagnosticInText(t *testing.T) {
	err := adv0009KernelError(t,
		adv0009BreachRow("", "", "a"),
		adv0009BreachRow("", "", "b"),
		adv0009BreachRow("", "", "c"),
	)
	ce := kernelResolveFailure(err)

	var textOut bytes.Buffer
	clierr.EmitText(&textOut, ce)

	// Isolate the finding lines; the envelope's own message/hint lines are
	// constant prose and carry no per-breach information.
	var findingLines []string
	for _, line := range strings.Split(textOut.String(), "\n") {
		if strings.Contains(line, codeEscapeRowShapeBreach+":") &&
			!strings.HasPrefix(line, "error:") {
			findingLines = append(findingLines, line)
		}
	}
	if len(findingLines) != 1 {
		t.Fatalf("want one finding line, got %d:\n%s", len(findingLines), textOut.String())
	}
	if !strings.Contains(findingLines[0], "3") {
		t.Errorf("a three-row breach with no source identity renders a text finding "+
			"carrying neither an identity nor a count, so it is indistinguishable "+
			"from a single-row breach and names nothing an operator can act on; "+
			"the kernel's own Error() says %q but the wire says %q",
			err.Error(), strings.TrimSpace(findingLines[0]))
	}
}

// --- ADV-2 -------------------------------------------------------------
//
// FAILURE MODE (RDR 0009, Failure Modes, "Visible break (programmatic
// producer)"): "A single errors.As/AsType call reports only the FIRST
// breach in the chain, so it is the wrong instrument for reading a
// multi-breach report — a test asserting on all offending rows must
// traverse the aggregate."
//
// escapeShapeBreaches is the verb's traversal, and its own comment claims
// it "traverses it rather than calling errors.As once, which would report
// the first breach alone." But it reaches the aggregate with a BARE TYPE
// ASSERTION — err.(interface{ Unwrap() []error }) — rather than errors.As.
// A bare assertion only matches the OUTERMOST error. Under any wrap the
// assertion fails and control falls through to exactly the single
// errors.As the RDR calls the wrong instrument, silently reducing an
// N-breach report to its first breach. Nothing signals the loss:
// errors.Is still classifies, so the Code, Group, exit code and Hint all
// stay correct while the identity report goes incomplete.
//
// errors.As demonstrably finds the join through a wrap (it finds it even
// through *clierr.CLIError, whose Unwrap returns a single error), so the
// robust instrument is available and unused.
func TestAdv0009_MultiBreachTraversalSurvivesAWrappedKernelError(t *testing.T) {
	err := adv0009KernelError(t,
		adv0009BreachRow("rule-a", "flow.toml:1", "x"),
		adv0009BreachRow("rule-b", "flow.toml:2", "y"),
		adv0009BreachRow("rule-c", "flow.toml:3", "z"),
	)

	direct := escapeShapeBreaches(err)
	if len(direct) != 3 {
		t.Fatalf("unwrapped: got %d breaches, want 3", len(direct))
	}

	wrapped := fmt.Errorf("resolving the transition table: %w", err)

	// The classification survives the wrap, which is what makes the loss
	// silent: the envelope still looks right.
	if !errors.Is(wrapped, resolve.ErrEscapeShapeBreach) {
		t.Fatal("errors.Is stopped classifying the wrapped breach")
	}

	got := escapeShapeBreaches(wrapped)
	if len(got) != len(direct) {
		t.Errorf("one %%w wrap reduced the breach report from %d identities to %d; "+
			"escapeShapeBreaches reaches the aggregate with a bare type assertion "+
			"instead of errors.As, so it degrades to the single-errors.As read the "+
			"RDR names as the wrong instrument for a multi-breach report",
			len(direct), len(got))
	}

	ce := kernelResolveFailure(wrapped)
	if ce.Code != codeEscapeRowShapeBreach {
		t.Fatalf("Code = %q, want %q", ce.Code, codeEscapeRowShapeBreach)
	}
	if len(ce.Findings) != 3 {
		t.Errorf("the wrapped breach reached the wire with %d findings, want 3; "+
			"the exit code and the stable Code are both still correct, so nothing "+
			"reveals that two offending rows were dropped from the envelope",
			len(ce.Findings))
	}
}

// --- ADV-3 (companion guards, expected to PASS) ------------------------
//
// FAILURE MODE (RDR 0009, Failure Modes, "Silent failure guarded against"
// and Consequences: "zero disposition change for conforming tables"): the
// new whole-table entry precondition must not perturb any conforming
// table, and the collapsed report must be a function of the input tuple
// rather than of row order. These are regression guards on properties the
// implementation currently HOLDS.
func TestAdv0009_TheCLIBreachEnvelopeIsAFunctionOfTheTableNotRowOrder(t *testing.T) {
	rows := []resolve.Row{
		adv0009BreachRow("rule-b", "flow.toml:2", "x"),
		adv0009BreachRow("rule-a", "flow.toml:1", "y"),
		adv0009BreachRow("rule-a", "flow.toml:1", "z"),
		adv0009BreachRow("rule-a", "flow.toml:3", "q"),
	}
	permutations := [][]int{
		{0, 1, 2, 3},
		{3, 2, 1, 0},
		{2, 0, 3, 1},
		{1, 3, 0, 2},
	}

	var want string
	for i, perm := range permutations {
		permuted := make([]resolve.Row, 0, len(perm))
		for _, idx := range perm {
			permuted = append(permuted, rows[idx])
		}
		ce := kernelResolveFailure(adv0009KernelError(t, permuted...))
		var out bytes.Buffer
		clierr.EmitJSON(&out, ce)
		if i == 0 {
			want = out.String()
			continue
		}
		if out.String() != want {
			t.Errorf("permutation %v produced a different envelope:\n got %s\nwant %s",
				perm, out.String(), want)
		}
	}
}

// A non-breach kernel error must NOT be recoded to this RDR's Code. RDR
// 0008's reserved-key breach travels the same Go error channel, and a
// blanket recode would mislabel it.
func TestAdv0009_ANonBreachKernelErrorKeepsTheGenericInternalCode(t *testing.T) {
	_, err := resolve.Resolve(resolve.Input{
		Table:      resolve.Table{Outcomes: []string{"adv-outcome"}},
		Recognized: "adv-outcome",
		Owned:      []resolve.Tag{{Key: "recognized", Value: "adv-outcome"}},
		Guards:     guardSeam(),
	})
	if err == nil {
		t.Fatal("want the RDR 0008 reserved-key breach, got nil")
	}
	ce := kernelResolveFailure(err)
	if ce.Code == codeEscapeRowShapeBreach {
		t.Fatalf("RDR 0008's reserved-key breach was recoded to %q", ce.Code)
	}
	if ce.Code != codeAccessorFailed {
		t.Fatalf("Code = %q, want %q", ce.Code, codeAccessorFailed)
	}
	if clierr.ExitCodeFor(ce) != 2 {
		t.Fatalf("exit = %d, want 2", clierr.ExitCodeFor(ce))
	}
}

// The Count field must round-trip through the wire unchanged: a consumer
// reads it back off JSON rather than re-parsing prose.
func TestAdv0009_ThePerIdentityCountRoundTripsThroughJSON(t *testing.T) {
	err := adv0009KernelError(t,
		adv0009BreachRow("rule-a", "flow.toml:1", "x"),
		adv0009BreachRow("rule-a", "flow.toml:1", "y"),
		adv0009BreachRow("rule-b", "flow.toml:2", "z"),
	)
	ce := kernelResolveFailure(err)

	var out bytes.Buffer
	clierr.EmitJSON(&out, ce)

	var back clierr.CLIError
	if e := json.Unmarshal(out.Bytes(), &back); e != nil {
		t.Fatalf("envelope did not round-trip: %v", e)
	}
	if len(back.Findings) != 2 {
		t.Fatalf("round-tripped findings = %d, want 2", len(back.Findings))
	}
	want := map[string]int{"rule-a": 2, "rule-b": 1}
	for _, f := range back.Findings {
		if f.Count != want[f.Rule] {
			t.Errorf("rule %q round-tripped Count = %d, want %d",
				f.Rule, f.Count, want[f.Rule])
		}
	}
}
