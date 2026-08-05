package resolver_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolver"
)

func legalInput() resolver.Input {
	return resolver.Input{
		FlowID:        "rdr",
		TableRevision: "2026-06-19",
		Owned: resolver.OwnedSnapshot{
			Available: true,
			Tags: resolver.TagSet{
				"status": "Draft",
			},
		},
		Observed: resolver.TagSet{
			"iteration": "2",
		},
		Recognized: resolver.Tag{
			Name:  "outcome",
			Value: "successful",
		},
		Table: []resolver.Edge{
			{
				Outcome: "successful",
				Guard: resolver.Guard{All: []resolver.TagPredicate{
					{
						Provenance: resolver.ProvenanceOwned,
						Name:       "status",
						Value:      "Draft",
					},
					{
						Provenance: resolver.ProvenanceObserved,
						Name:       "iteration",
						Value:      "2",
					},
					{
						Provenance: resolver.ProvenanceRecognized,
						Name:       "outcome",
						Value:      "successful",
					},
				}},
				NextTags: resolver.TagSet{
					"status": "Final",
				},
				Writes: []resolver.OwnedTagWrite{
					{
						Role:  "state",
						Name:  "status",
						Value: "Final",
					},
				},
			},
		},
	}
}

func noMatchInput() resolver.Input {
	input := legalInput()
	input.Table[0].Guard.All[0].Value = "Final"
	return input
}

func unmodeledOutcomeInput() resolver.Input {
	input := legalInput()
	input.Recognized.Value = "abandoned"
	return input
}

func ambiguousInput() resolver.Input {
	input := legalInput()
	second := input.Table[0]
	second.NextTags = resolver.TagSet{
		"status": "Implemented",
	}
	input.Table = append(input.Table, second)
	return input
}

func unevaluableInput() resolver.Input {
	input := legalInput()
	input.Table[0].Guard.Unevaluable = true
	return input
}

func assertRefusal(
	t *testing.T,
	input resolver.Input,
	want resolver.RefusalKind,
) resolver.Disposition {
	t.Helper()

	disposition, err := resolver.Resolve(input)
	if err != nil {
		t.Fatalf("Resolve returned Go error for modeled refusal %q: %v", want, err)
	}
	if disposition.Plan != nil {
		t.Fatalf("Resolve plan = %#v; want nil for refusal %q", disposition.Plan, want)
	}
	if disposition.Refusal == nil {
		t.Fatalf("Resolve refusal = nil; want %q", want)
	}
	if disposition.Refusal.Kind != want {
		t.Errorf("Resolve refusal kind = %q; want %q", disposition.Refusal.Kind, want)
	}
	return disposition
}

// REQ-1: "the same input must produce the same legal output"
// REQ-13: "Given the same flow identity, transition table revision, accessor-produced owned tag snapshot, caller-supplied observed tags, and freshly recognized outcome tag, resolve returns the same disposition: exactly one transition plan or exactly one typed refusal."
// REQ-24: "a resolution input is the tuple of flow identity, transition table revision, accessor-produced owned tag snapshot, observed tag-set, and freshly recognized outcome tag."
// REQ-25: "Replaying that tuple must replay the disposition."
// REQ-30: "Add focused kernel tests for deterministic replay and refusal classes, using fixtures that exercise owned, observed, and freshly recognized tags."
// REQ-33: "Both calls return value-identical transition plans, including next tags and accessor write descriptions."
// REQ-39: "No byte-stable hash or canonical serialization is introduced; determinism is value-level replay of the input tuple named in A1."
// HAPPY PATH
func TestResolve_ReplayReturnsValueIdenticalPlans(t *testing.T) {
	t.Parallel()

	input := legalInput()
	first, firstErr := resolver.Resolve(input)
	second, secondErr := resolver.Resolve(input)

	if firstErr != nil || secondErr != nil {
		t.Fatalf("Resolve errors = (%v, %v); want (nil, nil)", firstErr, secondErr)
	}
	if first.Plan == nil || first.Refusal != nil {
		t.Fatalf("first disposition = %#v; want exactly one plan", first)
	}
	if second.Plan == nil || second.Refusal != nil {
		t.Fatalf("second disposition = %#v; want exactly one plan", second)
	}
	if !reflect.DeepEqual(first, second) {
		t.Errorf("replayed dispositions differ:\nfirst:  %#v\nsecond: %#v", first, second)
	}
	if !reflect.DeepEqual(first.Plan.NextTags, resolver.TagSet{"status": "Final"}) {
		t.Errorf("next tags = %#v; want status=Final", first.Plan.NextTags)
	}
	wantWrites := []resolver.OwnedTagWrite{
		{
			Role:  "state",
			Name:  "status",
			Value: "Final",
		},
	}
	if !reflect.DeepEqual(first.Plan.Writes, wantWrites) {
		t.Errorf("writes = %#v; want %#v", first.Plan.Writes, wantWrites)
	}
}

// REQ-4: "The kernel must treat state as a tag-set, with guards expressed as predicates over tags and tag provenance distinguished as owned, observed, or freshly recognized."
// REQ-6: "the caller supplies the transition table, the recognized typed outcome, all non-owned context tags, and the owned tag snapshot already read from caller-provided artifacts by the accessor layer."
// REQ-9: "Inputs are: flow identity, transition table revision, owned-state snapshot values produced by the accessor layer, observed tags supplied by the caller, the freshly recognized outcome tag, and the reviewable transition table."
// REQ-10: "Evaluation builds a single tag-set view, selects matching candidate edges, refuses zero or multiple matches unless the table contract explicitly models an escape edge, and emits a transition plan."
// REQ-26: "the only successful selection is exactly one matching edge after guard evaluation."
// REQ-29: "Implement tag-set assembly, guard evaluation delegation, exact-one edge selection, and typed refusal behavior."
// REQ-38: "Resolution is bounded by the supplied transition table and tag snapshots."
// DOMAIN EDGE
func TestResolve_GuardsDistinguishTagProvenance(t *testing.T) {
	t.Parallel()

	input := legalInput()
	disposition, err := resolver.Resolve(input)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if disposition.Plan == nil || disposition.Refusal != nil {
		t.Fatalf("disposition = %#v; want exactly one plan", disposition)
	}

	input.Observed["status"] = "Draft"
	input.Owned.Tags["status"] = "Implemented"
	assertRefusal(t, input, resolver.RefusalNoMatch)
}

// REQ-5: "it must stay stateless and non-orchestrating, consume owned-state snapshots and write targets produced by the accessor layer, and refuse illegal or incomplete transition inputs rather than initiating work."
// REQ-7: "A successful resolution returns the next state tags and the owned-tag writes that the accessor layer applies back to the same caller-provided artifact boundary"
// REQ-11: "Artifact selection and accessor execution stay outside this RDR's contract."
// REQ-22: "The kernel MUST NOT initiate work"
// REQ-23: "The kernel MUST NOT execute persistence side effects directly."
// REQ-28: "Define the internal resolver package boundary, result taxonomy, and pure resolution entry point without CLI output or persistence side effects."
// REQ-31: "Expose only the kernel values needed by RDR 0005; do not add command output or state mutation semantics in this RDR."
// HAPPY PATH
func TestResolve_ReturnsInertOwnedTagWrites(t *testing.T) {
	t.Parallel()

	input := legalInput()
	before := resolver.TagSet{"status": input.Owned.Tags["status"]}

	disposition, err := resolver.Resolve(input)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if disposition.Plan == nil || disposition.Refusal != nil {
		t.Fatalf("disposition = %#v; want exactly one plan", disposition)
	}
	if !reflect.DeepEqual(input.Owned.Tags, before) {
		t.Errorf("owned snapshot mutated to %#v; want %#v", input.Owned.Tags, before)
	}
	want := input.Table[0].Writes
	if !reflect.DeepEqual(disposition.Plan.Writes, want) {
		t.Errorf("planned writes = %#v; want inert descriptions %#v", disposition.Plan.Writes, want)
	}
}

// REQ-2: "unmodeled matches must be refused instead of guessed"
// REQ-8: "an illegal, ambiguous, incomplete, or unmodeled input returns an explicit refusal class."
// REQ-14: "The kernel MUST refuse instead of guessing when no edge matches, more than one edge matches, required owned state is unavailable, a guard cannot be evaluated, or the recognized outcome is not modeled by the table."
// REQ-36: "The kernel returns the typed unmodeled-outcome refusal and performs no persistence side effect."
// INPUT EDGE
func TestResolve_RefusesUnmodeledOutcome(t *testing.T) {
	t.Parallel()
	assertRefusal(t, unmodeledOutcomeInput(), resolver.RefusalUnmodeledOutcome)
}

// REQ-14: "The kernel MUST refuse instead of guessing when no edge matches, more than one edge matches, required owned state is unavailable, a guard cannot be evaluated, or the recognized outcome is not modeled by the table."
// REQ-34: "The kernel returns the typed no-match refusal and performs no persistence side effect."
// INPUT EDGE
func TestResolve_RefusesNoMatch(t *testing.T) {
	t.Parallel()
	assertRefusal(t, noMatchInput(), resolver.RefusalNoMatch)
}

// REQ-3: "escape behavior must be explicit table data rather than a confident wrong edge"
// REQ-14: "The kernel MUST refuse instead of guessing when no edge matches, more than one edge matches, required owned state is unavailable, a guard cannot be evaluated, or the recognized outcome is not modeled by the table."
// REQ-27: "Zero, multiple, unavailable, or unevaluable candidates are refusals unless the table contains a modeled escape edge that itself matches exactly once."
// REQ-35: "The kernel returns the typed ambiguous-match refusal unless one explicit escape edge matches exactly once."
// ADVERSARIAL
func TestResolve_AmbiguityRequiresOneExplicitEscapeEdge(t *testing.T) {
	t.Parallel()

	input := ambiguousInput()
	assertRefusal(t, input, resolver.RefusalAmbiguousMatch)

	input.Table = append(input.Table, resolver.Edge{
		Outcome:   "successful",
		EscapeFor: resolver.RefusalAmbiguousMatch,
		Guard: resolver.Guard{All: []resolver.TagPredicate{
			{
				Provenance: resolver.ProvenanceOwned,
				Name:       "status",
				Value:      "Draft",
			},
		}},
		NextTags: resolver.TagSet{
			"status": "NeedsReview",
		},
		Writes: []resolver.OwnedTagWrite{
			{
				Role:  "state",
				Name:  "status",
				Value: "NeedsReview",
			},
		},
	})

	disposition, err := resolver.Resolve(input)
	if err != nil {
		t.Fatalf("Resolve with explicit escape: %v", err)
	}
	if disposition.Plan == nil || disposition.Refusal != nil {
		t.Fatalf("disposition = %#v; want exactly one escape plan", disposition)
	}
	if got := disposition.Plan.NextTags["status"]; got != "NeedsReview" {
		t.Errorf("escape next status = %q; want %q", got, "NeedsReview")
	}
}

// REQ-14: "The kernel MUST refuse instead of guessing when no edge matches, more than one edge matches, required owned state is unavailable, a guard cannot be evaluated, or the recognized outcome is not modeled by the table."
// REQ-15: "Modeled refusal is a value-level resolver disposition, not a CLI error and not the Go error path for parser bugs, IO failures, or programmer mistakes."
// REQ-37: "The kernel returns the corresponding value-level typed refusal and does not fall back to ambient discovery or the CLI/Go error path."
// BOUNDARY
func TestResolve_RefusesUnavailableOwnedStateAsValue(t *testing.T) {
	t.Parallel()

	input := legalInput()
	input.Owned.Available = false
	assertRefusal(t, input, resolver.RefusalOwnedStateUnavailable)
}

// REQ-14: "The kernel MUST refuse instead of guessing when no edge matches, more than one edge matches, required owned state is unavailable, a guard cannot be evaluated, or the recognized outcome is not modeled by the table."
// REQ-15: "Modeled refusal is a value-level resolver disposition, not a CLI error and not the Go error path for parser bugs, IO failures, or programmer mistakes."
// REQ-37: "The kernel returns the corresponding value-level typed refusal and does not fall back to ambient discovery or the CLI/Go error path."
// DOMAIN EDGE
func TestResolve_RefusesUnevaluableGuardAsValue(t *testing.T) {
	t.Parallel()
	assertRefusal(t, unevaluableInput(), resolver.RefusalGuardUnevaluable)
}

// REQ-16: "The kernel-owned refusal kind set is exactly: `no_match`, `ambiguous_match`, `owned_state_unavailable`, `guard_unevaluable`, `unmodeled_outcome`"
// REQ-17: "Each refusal kind must be stable enough for RDR 0005 to map to a CLI error code without inspecting error strings."
// DOMAIN EDGE
func TestResolve_RefusalKindsAreStableValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input resolver.Input
		want  resolver.RefusalKind
		wire  string
	}{
		{
			name:  "no match",
			input: noMatchInput(),
			want:  resolver.RefusalNoMatch,
			wire:  "no_match",
		},
		{
			name:  "ambiguous match",
			input: ambiguousInput(),
			want:  resolver.RefusalAmbiguousMatch,
			wire:  "ambiguous_match",
		},
		{
			name: "owned state unavailable",
			input: func() resolver.Input {
				input := legalInput()
				input.Owned.Available = false
				return input
			}(),
			want: resolver.RefusalOwnedStateUnavailable,
			wire: "owned_state_unavailable",
		},
		{
			name:  "guard unevaluable",
			input: unevaluableInput(),
			want:  resolver.RefusalGuardUnevaluable,
			wire:  "guard_unevaluable",
		},
		{
			name:  "unmodeled outcome",
			input: unmodeledOutcomeInput(),
			want:  resolver.RefusalUnmodeledOutcome,
			wire:  "unmodeled_outcome",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			disposition := assertRefusal(t, tt.input, tt.want)
			if got := string(disposition.Refusal.Kind); got != tt.wire {
				t.Errorf("refusal wire value = %q; want %q", got, tt.wire)
			}
		})
	}
}

// REQ-12: "the kernel only exposes structured success/refusal values that the CLI can map."
// REQ-13: "Given the same flow identity, transition table revision, accessor-produced owned tag snapshot, caller-supplied observed tags, and freshly recognized outcome tag, resolve returns the same disposition: exactly one transition plan or exactly one typed refusal."
// REQ-15: "Modeled refusal is a value-level resolver disposition, not a CLI error and not the Go error path for parser bugs, IO failures, or programmer mistakes."
// BOUNDARY
func TestResolve_ReturnsExactlyOneStructuredDisposition(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    resolver.Input
		wantPlan bool
	}{
		{
			name:     "legal plan",
			input:    legalInput(),
			wantPlan: true,
		},
		{
			name:  "modeled refusal",
			input: noMatchInput(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			disposition, err := resolver.Resolve(tt.input)
			if err != nil {
				t.Fatalf("Resolve returned Go error: %v", err)
			}
			hasPlan := disposition.Plan != nil
			hasRefusal := disposition.Refusal != nil
			if hasPlan == hasRefusal {
				t.Fatalf("disposition = %#v; want exactly one plan or refusal", disposition)
			}
			if hasPlan != tt.wantPlan {
				t.Errorf("has plan = %t; want %t", hasPlan, tt.wantPlan)
			}
		})
	}
}

// REQ-18: "The kernel MUST NOT print output"
// REQ-19: "The kernel MUST NOT inspect CLI flags"
// REQ-20: "The kernel MUST NOT discover ambient state"
// REQ-21: "The kernel MUST NOT choose artifacts on behalf of the caller"
// REQ-22: "The kernel MUST NOT initiate work"
// REQ-23: "The kernel MUST NOT execute persistence side effects directly."
// ADVERSARIAL
func TestResolve_IgnoresAmbientProcessStateAndProducesNoSideEffects(t *testing.T) {
	dir := t.TempDir()
	sentinel := filepath.Join(dir, "state.md")
	wantFile := []byte("status: Draft\n")
	if err := os.WriteFile(sentinel, wantFile, 0o600); err != nil {
		t.Fatalf("write sentinel: %v", err)
	}

	input := legalInput()
	input.Table[0].Writes[0].Role = sentinel

	originalArgs := os.Args
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	os.Args = []string{"intrastate", "resolve", "--force-wrong-edge"}
	t.Setenv("INTRASTATE_AMBIENT_STATUS", "Implemented")

	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	stderrReader, stderrWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("stderr pipe: %v", err)
	}
	os.Stdout = stdoutWriter
	os.Stderr = stderrWriter
	t.Cleanup(func() {
		os.Args = originalArgs
		os.Stdout = originalStdout
		os.Stderr = originalStderr
	})

	disposition, resolveErr := resolver.Resolve(input)

	os.Stdout = originalStdout
	os.Stderr = originalStderr
	if err := stdoutWriter.Close(); err != nil {
		t.Fatalf("close stdout writer: %v", err)
	}
	if err := stderrWriter.Close(); err != nil {
		t.Fatalf("close stderr writer: %v", err)
	}
	stdout, err := io.ReadAll(stdoutReader)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	stderr, err := io.ReadAll(stderrReader)
	if err != nil {
		t.Fatalf("read stderr: %v", err)
	}
	if err := stdoutReader.Close(); err != nil {
		t.Fatalf("close stdout reader: %v", err)
	}
	if err := stderrReader.Close(); err != nil {
		t.Fatalf("close stderr reader: %v", err)
	}

	if resolveErr != nil {
		t.Fatalf("Resolve: %v", resolveErr)
	}
	if disposition.Plan == nil || disposition.Refusal != nil {
		t.Fatalf("disposition = %#v; want exactly one plan", disposition)
	}
	if !bytes.Equal(stdout, nil) || !bytes.Equal(stderr, nil) {
		t.Errorf("Resolve output stdout=%q stderr=%q; want both empty", stdout, stderr)
	}
	gotFile, err := os.ReadFile(sentinel)
	if err != nil {
		t.Fatalf("read sentinel: %v", err)
	}
	if !bytes.Equal(gotFile, wantFile) {
		t.Errorf("sentinel content = %q; want unchanged %q", gotFile, wantFile)
	}
}

// REQ-32: "No third-party dependency is proposed at this stage."
// BOUNDARY
func TestResolver_PublicContractUsesOnlyStandardGoValues(t *testing.T) {
	t.Parallel()

	input := legalInput()
	disposition, err := resolver.Resolve(input)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if disposition.Plan == nil {
		t.Fatal("Resolve plan = nil; want a transition plan composed of Go values")
	}
}

// REQ-MVV: "Resolve must name and implementation must add a replay test that feeds the same table, owned snapshot, observed tags, and recognized outcome to the kernel twice and asserts value-identical dispositions. The same validation must include at least one value-level refusal each for `no_match`, `ambiguous_match`, `owned_state_unavailable`, `guard_unevaluable`, and `unmodeled_outcome`, and must assert those modeled refusals do not use the CLI or Go error path."
// HAPPY PATH
func TestResolve_MinimumViableValidation(t *testing.T) {
	t.Parallel()

	input := legalInput()
	first, firstErr := resolver.Resolve(input)
	second, secondErr := resolver.Resolve(input)
	if firstErr != nil || secondErr != nil {
		t.Fatalf("replay errors = (%v, %v); want (nil, nil)", firstErr, secondErr)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("replay lost value fidelity:\nfirst:  %#v\nsecond: %#v", first, second)
	}
	if first.Plan == nil || first.Refusal != nil {
		t.Fatalf("replay disposition = %#v; want exactly one plan", first)
	}

	unavailable := legalInput()
	unavailable.Owned.Available = false
	tests := []struct {
		name  string
		input resolver.Input
		want  resolver.RefusalKind
	}{
		{
			name:  "no_match",
			input: noMatchInput(),
			want:  resolver.RefusalNoMatch,
		},
		{
			name:  "ambiguous_match",
			input: ambiguousInput(),
			want:  resolver.RefusalAmbiguousMatch,
		},
		{
			name:  "owned_state_unavailable",
			input: unavailable,
			want:  resolver.RefusalOwnedStateUnavailable,
		},
		{
			name:  "guard_unevaluable",
			input: unevaluableInput(),
			want:  resolver.RefusalGuardUnevaluable,
		},
		{
			name:  "unmodeled_outcome",
			input: unmodeledOutcomeInput(),
			want:  resolver.RefusalUnmodeledOutcome,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertRefusal(t, tt.input, tt.want)
		})
	}
}
