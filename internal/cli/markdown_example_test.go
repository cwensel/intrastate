package cli

// Keep this example in the regular CLI suite so make check exercises the
// shipped model and document through the same path as the walkthrough.

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
)

const markdownReviewModel = "../../models/examples/markdown-review.toml"

func TestMarkdownExamplePlansAndVerifiesOneLineEdit(t *testing.T) {
	t.Parallel()
	target, original := markdownReviewFixture(t)
	plan := markdownReviewPlan(t, target)
	assertMarkdownDocument(t, target, original)
	data := flowData(t, requireSuccess(t, markdownReviewArgs(target,
		"set-state", "--plan", plan)...))
	owned, ok := data["owned"].(map[string]any)
	if !ok || owned["status"] != "approved" {
		t.Fatalf("read-back status = %#v; want approved", data["owned"])
	}
	expected := bytes.Replace(original, []byte("Status: draft"), []byte("Status: approved"), 1)
	assertMarkdownDocument(t, target, expected)

	data = flowData(t, requireSuccess(t, markdownReviewArgs(target,
		"resolve", "--outcome", "approve")...))
	if data["rule"] != "already-approved" {
		t.Fatalf("repeated approval selected %#v; want already-approved", data["rule"])
	}
	writes, ok := data["writes"].(map[string]any)
	if !ok || len(writes) != 0 {
		t.Fatalf("repeated approval writes = %#v; want empty", data["writes"])
	}
	assertMarkdownDocument(t, target, expected)
}

func TestMarkdownExampleRefusesChangedAnchorsWithoutWriting(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name        string
		replacement string
		token       string
	}{
		{
			name:        "missing status line",
			replacement: "Review pending",
			token:       "edit_anchor_unmatched",
		},
		{
			name:        "duplicate status line",
			replacement: "Status: draft\nStatus: draft",
			token:       "edit_anchor_ambiguous",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			target, original := markdownReviewFixture(t)
			plan := markdownReviewPlan(t, target)
			changed := bytes.Replace(original, []byte("Status: draft"), []byte(tc.replacement), 1)
			if err := os.WriteFile(target, changed, 0o600); err != nil {
				t.Fatal(err)
			}

			_, _, err := runCmd(t, markdownReviewArgs(target,
				"set-state", "--plan", plan)...)
			var ce *clierr.CLIError
			if !asCLIError(err, &ce) {
				t.Fatalf("expected a structured refusal, got %v", err)
			}
			if clierr.ExitCodeFor(err) != 2 || !strings.Contains(ce.Detail, tc.token) {
				t.Fatalf("expected exit 2 with %s, got %+v", tc.token, ce)
			}
			assertMarkdownDocument(t, target, changed)
		})
	}
}

func markdownReviewFixture(t *testing.T) (string, []byte) {
	t.Helper()
	if _, err := exec.LookPath("sed"); err != nil {
		t.Skip("the Markdown example requires sed on PATH")
	}
	original, err := os.ReadFile("../../models/examples/markdown-review.md")
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(t.TempDir(), "release checklist.md")
	if err := os.WriteFile(target, original, 0o600); err != nil {
		t.Fatal(err)
	}
	return target, original
}

func markdownReviewArgs(target, verb string, extra ...string) []string {
	args := []string{
		"flow", verb, "--model", markdownReviewModel,
		"--artifact", "document=" + target, "--allow-commands", "--as=json",
	}
	return append(args, extra...)
}

func markdownReviewPlan(t *testing.T, target string) string {
	t.Helper()
	stdout := requireSuccess(t, markdownReviewArgs(target,
		"resolve", "--outcome", "approve")...)
	data := flowData(t, stdout)
	if data["rule"] != "approve-draft" {
		t.Fatalf("selected %#v; want approve-draft", data["rule"])
	}
	plan := filepath.Join(filepath.Dir(target), "plan.json")
	if err := os.WriteFile(plan, []byte(stdout), 0o600); err != nil {
		t.Fatal(err)
	}
	return plan
}

func assertMarkdownDocument(t *testing.T, target string, expected []byte) {
	t.Helper()
	actual, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Fatalf("document bytes differ:\ngot:\n%s\nwant:\n%s", actual, expected)
	}
}
