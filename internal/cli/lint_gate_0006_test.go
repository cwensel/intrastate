package cli

// RDR 0006 — the repository acceptance gate (REQ-119..REQ-122).
//
// These assert on repository ARTIFACTS — the CI workflow, the Makefile,
// and the checked-in transition model A8 requires — because the gate this
// RDR fixes is a repository boundary, not a library call. Scenario 7 is
// mechanically checkable exactly as the RDR states: both jobs exist in
// that file, the new one names that command, and the model path it passes
// is the checked-in model.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/graphlint"
	"github.com/cwensel/intrastate/internal/table"
)

// checkedInModelPath is the home A8's transition model takes. The gate
// lints THIS file, never a fixture-only corpus.
const checkedInModelPath = "models/rdr.toml"

func readRepoFile(t *testing.T, rel string) string {
	t.Helper()

	body, err := os.ReadFile(filepath.Join(repoRootFor(t), rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	return string(body)
}

// REQ-121 / A8: "a transition model is checked into this repo for that
// gate to lint." Scenario 7 has no subject until one is authored and
// homed.
// HAPPY PATH
func TestReq121_ATransitionModelIsCheckedIntoTheRepo(t *testing.T) {
	path := filepath.Join(repoRootFor(t), checkedInModelPath)

	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("no transition model is checked in at %s: %v. Scenario 7 "+
			"has no subject until one is authored and homed (A8).",
			checkedInModelPath, err)
	}

	// It must be a model this repo's own loader accepts — a gate whose
	// subject does not normalize never reaches lint.
	m, lerr := table.Load(src, checkedInModelPath)
	if lerr != nil {
		t.Fatalf("the checked-in model does not conform to the transition "+
			"model schema: %v", lerr)
	}
	if m.ID == "" {
		t.Error("the checked-in model declares no `[model].id`")
	}
	if len(m.Rows) == 0 {
		t.Error("the checked-in model carries no candidate rows, so linting " +
			"it proves nothing")
	}
	if len(m.Initial) == 0 {
		t.Error("the checked-in model declares no initial owned state")
	}
}

// REQ-119 / SC-7: "`.github/workflows/ci.yml` contains a **`graph-lint`**
// job — not `lint`, which is the shipped golangci-lint job and must
// survive untouched — that runs `make build` and then
// `./bin/intrastate lint --as=json` over the checked-in transition model,
// with the assertion reading the JSON `code` field, not the exit integer
// alone."
// BOUNDARY
func TestReq119_CICarriesAGraphLintJobOverTheCheckedInModel(t *testing.T) {
	ci := readRepoFile(t, ".github/workflows/ci.yml")

	// The shipped golangci-lint job must survive untouched.
	if !strings.Contains(ci, "golangci-lint") {
		t.Error("the shipped golangci-lint job is gone from ci.yml; it must " +
			"survive untouched")
	}
	if !strings.Contains(ci, "\n  lint:") {
		t.Error("the shipped `lint:` job key is gone from ci.yml")
	}

	// The NEW job is keyed `graph-lint`, not `lint`.
	if !strings.Contains(ci, "graph-lint:") {
		t.Fatalf("ci.yml carries no `graph-lint` job; the shipped `lint` job " +
			"is the golangci-lint gate and is not reusable")
	}

	// It runs `make build` and then the built binary over the checked-in
	// model — not a hook wrapper and not a unit-test-only engine path.
	for _, want := range []string{
		"make build",
		"./bin/intrastate lint",
		"--as=json",
		checkedInModelPath,
	} {
		if !strings.Contains(ci, want) {
			t.Errorf("the graph-lint job does not name %q; it must run the "+
				"production command over the checked-in model", want)
		}
	}

	// The assertion reads the JSON `code` field, not the exit integer
	// alone.
	if !strings.Contains(ci, "code") {
		t.Errorf("the graph-lint job's assertion does not read the JSON " +
			"`code` field; the exit integer alone is not the oracle")
	}
}

// REQ-119 (second half) / IP Phase 3: "`make check` is additionally wired
// to the same command for local parity, which requires adding the `build`
// edge `check` currently lacks."
// BOUNDARY
func TestReq119_MakeCheckIsWiredToTheSameCommandWithABuildEdge(t *testing.T) {
	mk := readRepoFile(t, "Makefile")

	var checkLine string
	for _, line := range strings.Split(mk, "\n") {
		if strings.HasPrefix(line, "check:") {
			checkLine = line
			break
		}
	}
	if checkLine == "" {
		t.Fatal("the Makefile declares no `check` target")
	}
	if !strings.Contains(checkLine, "build") {
		t.Errorf("`check` does not depend on `build`: %q. Local parity runs "+
			"the built command, which needs the `build` edge `check` "+
			"currently lacks.", checkLine)
	}
	if !strings.Contains(mk, "graph-lint") {
		t.Error("the Makefile carries no graph-lint target wired to the same " +
			"command for local parity")
	}
}

// REQ-120 / SC-7: "A deliberately illegal commit to the checked-in model
// must fail that job."
// ADVERSARIAL
func TestReq120_ADeliberatelyIllegalEditToTheCheckedInModelFailsTheGate(t *testing.T) {
	src := readRepoFile(t, checkedInModelPath)

	// The gate's subject is the checked-in model, so the proof that the
	// gate BITES is that a deliberate defect in that same document fails
	// the same command. The defect is minimal and local: drop the
	// `[initial]` table, which the disposition table routes to a blocking
	// `graph-dangling-edge` and never to a clean empty reachable set.
	broken := dropInitialTable(src)
	if broken == src {
		t.Fatalf("the checked-in model carries no `[initial]` table to "+
			"remove, so this test cannot author a deliberate defect in it:\n%s",
			checkedInModelPath)
	}

	path := writeModel(t, broken)
	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err == nil {
		t.Fatalf("a deliberately illegal edit to the checked-in model passed "+
			"the gate; stdout:\n%s", stdout)
	}
	if code := clierr.ErrorCode(err); code != graphlint.AggregateCode {
		t.Errorf("the illegal edit returned %q; want %q", code,
			graphlint.AggregateCode)
	}
}

// REQ-122 / SC-23: "lint the checked-in transition model (A8) in its
// accepted state. … the blocking finding count is recorded. Zero is the
// pass condition; a non-zero count on a model its maintainers accept is
// the mechanical trigger for the guard-aware-pruning successor RDR named
// in `Consequences`, and MUST be reported rather than waived."
// HAPPY PATH
func TestReq122_FalsePositiveCensusOverTheCheckedInModelIsZero(t *testing.T) {
	src := readRepoFile(t, checkedInModelPath)
	path := writeModel(t, src)

	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err == nil {
		// Zero blocking findings is the pass condition.
		return
	}

	env := parseFailure(t, stdout)
	var count int
	var codes []string
	if env.Findings != nil {
		for _, f := range *env.Findings {
			if graphlint.IsBlocking(f.Code) {
				count++
				codes = append(codes, f.Code+"@"+f.Rule)
			}
		}
	}
	// Reported, never waived: the count and the codes are the record the
	// successor RDR is triggered from.
	t.Errorf("false-positive census over the checked-in model in its "+
		"ACCEPTED state: %d blocking findings %v. Zero is the pass "+
		"condition; a non-zero count on a model its maintainers accept is "+
		"the mechanical trigger for the guard-aware-pruning successor RDR "+
		"and MUST be reported rather than waived.", count, codes)
}

// dropInitialTable removes the `[initial]` table and its assignments,
// leaving the rest of the document byte-identical.
func dropInitialTable(src string) string {
	var out []string
	var inInitial bool
	for _, line := range strings.Split(src, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "[initial]" {
			inInitial = true
			continue
		}
		if inInitial {
			// The table ends at the next table header.
			if strings.HasPrefix(trimmed, "[") {
				inInitial = false
			} else {
				continue
			}
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}
