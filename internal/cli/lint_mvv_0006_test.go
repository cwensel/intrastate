package cli

// RDR 0006 — the Minimum Viable Validation.
//
// REQ-MVV is one runnable end-to-end test. Per ASSUMPTION-7 the matrix
// drives the root Cobra command through ExecuteAndEmit — "the same
// command shape intended for CI" — using `--model <path>` (REQ-14), and
// reads the JSON envelope the gateway emits. Nothing here calls the
// engine API directly.

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
	"github.com/newcoinc/intrastate/internal/graphlint"
	"github.com/newcoinc/intrastate/internal/guard"
	"github.com/newcoinc/intrastate/internal/table"
)

// REQ-MVV: "Add a fixture-backed lint invocation that uses the same
// command shape intended for CI. It must pass one legal transition model
// and fail one illegal model for each blocking invariant class named in
// this RDR, asserting stable finding codes and source rule/context
// identity in JSON mode. … It must also include one **multi-defect group**
// asserting every expected code from a single run, so a first-failure
// engine cannot pass. The legal matrix must include a group closed by a
// bare escape row, a redundant row, and an unreachable rule, asserting
// `graph-coverage-closed-by-escape`, `graph-redundant-row`, and
// `graph-unreachable-rule` on success. Every blocking run must return
// aggregate `CLIError.Code = graph-lint-failed` with `GroupUserEnv` exit
// behavior and a machine-readable finding list; every run, clean or not,
// must emit the `findings` key even when empty."
// HAPPY PATH
func TestMVV_GraphLintAuthorityAndGuarantees(t *testing.T) {
	t.Run("legal model passes", mvvLegalModelPasses)
	t.Run("illegal matrix", mvvIllegalMatrix)
	t.Run("multi-defect group", mvvMultiDefectGroup)
	t.Run("legal advisory matrix", mvvLegalAdvisoryMatrix)
	t.Run("findings key always emitted", mvvFindingsKeyAlwaysEmitted)
}

// mvvLegalModelPasses: "It must pass one legal transition model".
func mvvLegalModelPasses(t *testing.T) {
	path := writeModel(t, legalModel)

	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err != nil {
		t.Fatalf("the legal transition model failed: %v\nstdout:\n%s", err, stdout)
	}
	env := parseSuccess(t, stdout)
	if env.Type != "ok" {
		t.Errorf("envelope type = %q; want %q", env.Type, "ok")
	}
	if env.Data == nil || env.Data.Findings == nil {
		t.Fatalf("the success envelope carries no `data.findings`:\n%s", stdout)
	}
	for _, f := range *env.Data.Findings {
		if graphlint.IsBlocking(f.Code) {
			t.Errorf("the legal model reported the blocking code %q", f.Code)
		}
	}

	// A pass is a PROOF, not the absence of a run. The receipt is the
	// `findings` key emitted with the empty list, and the raw JSON is the
	// oracle so an omitted key cannot decode to a satisfied nil slice.
	var raw map[string]json.RawMessage
	if jerr := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &raw); jerr != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", jerr, stdout)
	}
	dataRaw, ok := raw["data"]
	if !ok {
		t.Fatalf("the legal run emitted no `data` key:\n%s", stdout)
	}
	var data map[string]json.RawMessage
	if jerr := json.Unmarshal(dataRaw, &data); jerr != nil {
		t.Fatalf("`data` is not an object: %v\n%s", jerr, stdout)
	}
	findings, ok := data["findings"]
	if !ok {
		t.Fatalf("the legal run emitted no `data.findings` receipt:\n%s", stdout)
	}
	if got := strings.TrimSpace(string(findings)); got != "[]" {
		t.Errorf("data.findings = %s; want the empty-list receipt `[]`", got)
	}

	// And the proof must have been ATTEMPTED: the model's one group carries
	// two rows partitioning a finite guard dimension, so a lint that
	// examined it decides that group. The illegal counterpart below shares
	// this model's shape and DOES fail, which is what makes this pass a
	// verdict rather than a no-op.
	illegal := writeModel(t, illegalModel)
	if _, _, ierr := runCmd(t, "lint", "--model", illegal, "--as=json"); ierr == nil {
		t.Fatal("the illegal counterpart also passed, so the legal model's " +
			"green is a no-op rather than a verdict")
	}
}

// mvvIllegalMatrix: "fail one illegal model for each blocking invariant
// class named in this RDR, asserting stable finding codes and source
// rule/context identity in JSON mode."
//
// The named matrix is `graph-dangling-edge`, `graph-dead-end`,
// `graph-overlap` (both an ordinary-row pair and an escape-row pair
// sharing a failure class), `graph-coverage-gap`,
// `graph-unprovable-coverage` (both a non-finite dimension and a withheld
// claim naming row and atom), `graph-single-valued-state`,
// `graph-always-present-owned`, `graph-owned-before-write`,
// `graph-terminal-escape`, and `graph-product-too-large`.
func mvvIllegalMatrix(t *testing.T) {
	for _, tc := range []struct {
		name string
		src  string
		want string
		// rule, when set, is the source rule/context id the finding must
		// carry. Where the defect is a model-level one there is no rule to
		// name and the identity assertion falls back to REQ-127's span or
		// graph element id.
		rule string
		// atom asserts the four atom fields on the withheld-claim arm.
		atom bool
		// class asserts the failure class on an escape-scoped finding.
		class string
	}{
		{name: "dangling-edge/missing-root", src: mvvNoRoot,
			want: graphlint.CodeDanglingEdge},
		{name: "dead-end", src: mvvDeadEnd,
			want: graphlint.CodeDeadEnd},
		{name: "overlap/ordinary-pair", src: mvvOrdinaryOverlap,
			want: graphlint.CodeOverlap, rule: "over-one"},
		{name: "overlap/escape-pair-shared-class", src: mvvEscapeOverlap,
			want: graphlint.CodeOverlap, class: "no_match"},
		{name: "coverage-gap", src: mvvCoverageGap,
			want: graphlint.CodeCoverageGap, rule: "only-on"},
		{name: "unprovable/non-finite-dimension", src: mvvNonFinite,
			want: graphlint.CodeUnprovableCoverage},
		{name: "unprovable/withheld-claim", src: mvvWithheld,
			want: graphlint.CodeUnprovableCoverage, rule: "can-refuse", atom: true},
		{name: "always-present-owned", src: mvvAlwaysPresentOwned,
			want: graphlint.CodeAlwaysPresentOwned},
		{name: "owned-before-write", src: mvvOwnedBeforeWrite,
			want: graphlint.CodeOwnedBeforeWrite, rule: "reads-unwritten"},
		{name: "terminal-escape", src: mvvTerminalEscape,
			want: graphlint.CodeTerminalEscape},
		{name: "product-too-large", src: mvvProductTooLarge,
			want: graphlint.CodeProductTooLarge},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := writeModel(t, tc.src)

			stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
			if err == nil {
				t.Fatalf("the illegal model succeeded; stdout:\n%s", stdout)
			}

			// "Every blocking run must return aggregate
			// `CLIError.Code = graph-lint-failed` with `GroupUserEnv` exit
			// behavior".
			if code := clierr.ErrorCode(err); code != graphlint.AggregateCode {
				t.Errorf("aggregate code = %q; want %q", code, graphlint.AggregateCode)
			}
			var ce *clierr.CLIError
			if !asCLIError(err, &ce) {
				t.Fatalf("the failure is not a CLIError: %v", err)
			}
			if ce.Group != clierr.GroupUserEnv {
				t.Errorf("group = %v; want GroupUserEnv", ce.Group)
			}
			if got := clierr.ExitCodeFor(err); got != 2 {
				t.Errorf("exit = %d; want 2", got)
			}

			// "and a machine-readable finding list".
			env := parseFailure(t, stdout)
			if env.Findings == nil {
				t.Fatalf("the failure envelope carries no `findings` key:\n%s", stdout)
			}
			findings := *env.Findings
			if len(findings) == 0 {
				t.Fatalf("the failure envelope carries an empty finding list "+
					"alongside a blocking failure:\n%s", stdout)
			}

			// "asserting stable finding codes".
			var matched []clierr.Finding
			for _, f := range findings {
				if f.Code == tc.want {
					matched = append(matched, f)
				}
			}
			if len(matched) == 0 {
				var got []string
				for _, f := range findings {
					got = append(got, f.Code)
				}
				t.Fatalf("no %s finding; got codes %v\nstdout:\n%s",
					tc.want, got, stdout)
			}
			for _, f := range matched {
				if f.Severity != graphlint.SeverityBlocking {
					t.Errorf("%s severity = %q; want %q",
						f.Code, f.Severity, graphlint.SeverityBlocking)
				}
				if f.Model != "lintfix" {
					t.Errorf("%s model = %q; want the `[model].id` %q",
						f.Code, f.Model, "lintfix")
				}
				if strings.TrimSpace(f.Message) == "" {
					t.Errorf("%s carries no human message", f.Code)
				}
			}

			// "and source rule/context identity in JSON mode" — REQ-127's
			// fallback applies where the defect names no rule.
			var identified bool
			for _, f := range matched {
				if f.Rule != "" || f.Span != "" || f.Element != "" {
					identified = true
				}
			}
			if !identified {
				t.Errorf("no %s finding carries a source rule/context id, a "+
					"source span, or a graph element id:\n%s", tc.want, stdout)
			}
			if tc.rule != "" {
				var named bool
				for _, f := range matched {
					if f.Rule == tc.rule ||
						strings.Contains(f.Message, tc.rule) {
						named = true
					}
				}
				if !named {
					t.Errorf("no %s finding names the source rule %q:\n%s",
						tc.want, tc.rule, stdout)
				}
			}

			// The withheld-claim arm names row AND atom.
			if tc.atom {
				var ok bool
				for _, f := range matched {
					if f.Reason == graphlint.ReasonRowCanRefuse &&
						f.Key != "" && f.Operator != "" &&
						f.Literal != "" && f.Block != "" {
						ok = true
					}
				}
				if !ok {
					t.Errorf("no %s finding carries reason %q with all four "+
						"atom fields:\n%s", tc.want,
						graphlint.ReasonRowCanRefuse, stdout)
				}
			}

			// The escape-scoped arm names the shared failure class.
			if tc.class != "" {
				var ok bool
				for _, f := range matched {
					if f.Class == tc.class {
						ok = true
					}
				}
				if !ok {
					t.Errorf("no %s finding carries the shared failure class "+
						"%q:\n%s", tc.want, tc.class, stdout)
				}
			}
		})
	}
}

// mvvMultiDefectGroup: "one **multi-defect group** asserting every
// expected code from a single run, so a first-failure engine cannot pass."
func mvvMultiDefectGroup(t *testing.T) {
	path := writeModel(t, mvvMultiDefect)

	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err == nil {
		t.Fatalf("the multi-defect model succeeded; stdout:\n%s", stdout)
	}
	if code := clierr.ErrorCode(err); code != graphlint.AggregateCode {
		t.Errorf("aggregate code = %q; want %q", code, graphlint.AggregateCode)
	}

	env := parseFailure(t, stdout)
	if env.Findings == nil {
		t.Fatalf("the failure envelope carries no `findings` key:\n%s", stdout)
	}

	// The scenario the RDR prescribes is a multi-defect GROUP: all four
	// rows bind the same match pattern and outcome, so RDR 0003's grouping
	// puts them in ONE group carrying every defect class. Splitting them
	// across two single-class groups would still satisfy a code-set
	// assertion while no longer being the prescribed scenario, so the
	// grouping is pinned directly against the fixture the run reads.
	m, lerr := table.Load([]byte(mvvMultiDefect), path)
	if lerr != nil {
		t.Fatalf("the multi-defect fixture does not normalize: %v", lerr)
	}
	groups := guard.Groups(m)
	if len(groups) != 1 {
		var ctxs []string
		for _, g := range groups {
			ctxs = append(ctxs, g.Context.String())
		}
		t.Fatalf("the multi-defect fixture forms %d scoped row groups (%v); "+
			"want exactly 1 — the RDR names a multi-defect GROUP, and rows "+
			"split across per-class groups are a different scenario",
			len(groups), ctxs)
	}
	var groupRows []string
	for _, row := range groups[0].Rows {
		groupRows = append(groupRows, row.RuleID)
	}
	slices.Sort(groupRows)
	wantRows := []string{"over-one", "over-two", "refuse-one", "refuse-two"}
	if !slices.Equal(groupRows, wantRows) {
		t.Errorf("the single group %q carries rows %v; want all four "+
			"defect-bearing rows %v in ONE group",
			groups[0].Context.String(), groupRows, wantRows)
	}

	// Every expected code from a SINGLE run, asserted as the exact
	// MULTISET the run reports. A deduped code SET leaves per-row
	// multiplicity unproven: the two rows that can refuse are two
	// withheld-claim findings, and an engine collapsing them to one
	// would report the same set.
	var got []string
	for _, f := range *env.Findings {
		got = append(got, f.Code)
	}
	slices.Sort(got)

	want := []string{
		graphlint.CodeAlwaysPresentOwned,
		graphlint.CodeOverlap,
		graphlint.CodeOwnedBeforeWrite,
		graphlint.CodeUnprovableCoverage,
		graphlint.CodeUnprovableCoverage,
	}
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Errorf("the multi-defect run reported the finding codes %v; want "+
			"exactly %v. A first-failure engine reports a prefix; an engine "+
			"deduping per-row findings reports one %q instead of two.\n"+
			"stdout:\n%s", got, want,
			graphlint.CodeUnprovableCoverage, stdout)
	}

	// One overlapping pair is ONE finding naming both rows, not one per row.
	var overlaps int
	for _, f := range *env.Findings {
		if f.Code == graphlint.CodeOverlap {
			overlaps++
		}
	}
	if overlaps != 1 {
		t.Errorf("%d %s findings for one overlapping pair; want exactly 1 "+
			"naming both rows:\n%s", overlaps, graphlint.CodeOverlap, stdout)
	}

	// And the two withheld claims name the two DISTINCT refusing rows, so
	// the multiplicity above is per-row and not one row reported twice.
	var refusing []string
	for _, f := range *env.Findings {
		if f.Code == graphlint.CodeUnprovableCoverage {
			refusing = append(refusing, f.Rule)
		}
	}
	slices.Sort(refusing)
	if wantRefusing := []string{"refuse-one", "refuse-two"}; !slices.Equal(
		refusing, wantRefusing) {
		t.Errorf("the withheld-claim findings name %v; want one per "+
			"refusing row, %v:\n%s", refusing, wantRefusing, stdout)
	}
}

// mvvLegalAdvisoryMatrix: "The legal matrix must include a group closed by
// a bare escape row, a redundant row, and an unreachable rule, asserting
// `graph-coverage-closed-by-escape`, `graph-redundant-row`, and
// `graph-unreachable-rule` on success."
func mvvLegalAdvisoryMatrix(t *testing.T) {
	path := writeModel(t, mvvLegalAdvisory)

	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err != nil {
		t.Fatalf("the legal advisory model failed: %v\nstdout:\n%s", err, stdout)
	}

	env := parseSuccess(t, stdout)
	if env.Data == nil || env.Data.Findings == nil {
		t.Fatalf("the success envelope carries no `data.findings`:\n%s", stdout)
	}

	var got []string
	for _, f := range *env.Data.Findings {
		if !slices.Contains(got, f.Code) {
			got = append(got, f.Code)
		}
		if graphlint.IsBlocking(f.Code) {
			t.Errorf("the success payload carries the blocking code %q", f.Code)
		}
	}
	slices.Sort(got)

	for _, c := range []string{
		graphlint.CodeCoverageClosedByEscape,
		graphlint.CodeRedundantRow,
		graphlint.CodeUnreachableRule,
	} {
		if !slices.Contains(got, c) {
			t.Errorf("the legal matrix did not report %q on success; got %v. "+
				"A lint that silently drops the advisory tier fails.\n"+
				"stdout:\n%s", c, got, stdout)
		}
	}
}

// mvvFindingsKeyAlwaysEmitted: "every run, clean or not, must emit the
// `findings` key even when empty."
func mvvFindingsKeyAlwaysEmitted(t *testing.T) {
	// Clean: `data.findings` present with value `[]`. The oracle is the
	// RAW JSON so an omitted key fails rather than decoding to nil.
	clean := writeModel(t, legalModel)
	stdout, _, err := runCmd(t, "lint", "--model", clean, "--as=json")
	if err != nil {
		t.Fatalf("the legal model failed: %v\nstdout:\n%s", err, stdout)
	}

	var raw map[string]json.RawMessage
	if jerr := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &raw); jerr != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", jerr, stdout)
	}
	dataRaw, ok := raw["data"]
	if !ok {
		t.Fatalf("the clean run emitted no `data` key; `Data` is itself "+
			"`omitempty`, so the success path must assign it a non-nil "+
			"struct:\n%s", stdout)
	}
	var data map[string]json.RawMessage
	if jerr := json.Unmarshal(dataRaw, &data); jerr != nil {
		t.Fatalf("`data` is not an object: %v\n%s", jerr, stdout)
	}
	findings, ok := data["findings"]
	if !ok {
		t.Fatalf("the clean run emitted no `data.findings` key; the empty "+
			"list is the proof's receipt:\n%s", stdout)
	}
	if got := strings.TrimSpace(string(findings)); got != "[]" {
		t.Errorf("clean run data.findings = %s; want `[]`", got)
	}

	// Not clean: the `findings` key is present on the failure envelope.
	dirty := writeModel(t, illegalModel)
	failOut, _, ferr := runCmd(t, "lint", "--model", dirty, "--as=json")
	if ferr == nil {
		t.Fatalf("the illegal model succeeded:\n%s", failOut)
	}
	var failRaw map[string]json.RawMessage
	if jerr := json.Unmarshal([]byte(strings.TrimSpace(failOut)), &failRaw); jerr != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", jerr, failOut)
	}
	if _, ok := failRaw["findings"]; !ok {
		t.Errorf("the failing run emitted no top-level `findings` key:\n%s",
			failOut)
	}
}
