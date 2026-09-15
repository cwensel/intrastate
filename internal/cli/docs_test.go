package cli

// The generated-docs contract.
//
// `make docs-check` compares regenerated output against what is
// committed, so the gate is only meaningful if generation is a pure
// function of the command tree. These tests pin the properties that
// makes true — anything that varies per build, per checkout, or per run
// would turn the gate into a random CI failure rather than a drift
// signal.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// renderBoth returns the two generated files' contents for a fresh tree.
func renderBoth(t *testing.T) (reference, llms string) {
	t.Helper()
	root := NewRootCmd()
	var ref, idx strings.Builder
	writeCLIReference(&ref, root)
	writeLLMsTxt(&idx, root)
	return ref.String(), idx.String()
}

// TestDocs_GenerationIsDeterministic is the property the staleness gate
// rests on. Two renders of the same tree must be byte-identical; a map
// iterated in range order or a timestamp would break this and make
// `make check` fail at random.
func TestDocs_GenerationIsDeterministic(t *testing.T) {
	for i := range 8 {
		ref1, llms1 := renderBoth(t)
		ref2, llms2 := renderBoth(t)
		if ref1 != ref2 {
			t.Fatalf("cli-reference render %d is not deterministic", i)
		}
		if llms1 != llms2 {
			t.Fatalf("llms.txt render %d is not deterministic", i)
		}
	}
}

// TestDocs_CarryNoBuildStamp pins the one value that genuinely varies
// between builds. version.Get() is ldflags-stamped, so if it reached the
// generated output every release build would dirty the committed files
// and the gate would cry wolf.
func TestDocs_CarryNoBuildStamp(t *testing.T) {
	ref, llms := renderBoth(t)
	// The stamp is "dev" in an unstamped test binary, which is too common
	// a word to grep for. Assert on the shape the stamp renders as
	// instead: cobra's version line and the identity string's parens.
	for name, body := range map[string]string{
		"docs/cli-reference.md": ref,
		"llms.txt":              llms,
	} {
		if strings.Contains(body, "version for intrastate") {
			t.Errorf("%s carries cobra's --version line; it is ldflags-stamped "+
				"and would dirty the file on every release build", name)
		}
		if strings.Contains(body, "(commit ") {
			t.Errorf("%s carries the build-identity string; it is "+
				"ldflags-stamped and must not reach generated output", name)
		}
	}
}

// TestDocs_ReferenceCoversEveryVisibleCommand pins coverage: a verb added
// to the tree appears in the reference without anyone remembering to
// list it. The hidden `docs` command itself is excluded — it addresses
// maintainers, not callers.
func TestDocs_ReferenceCoversEveryVisibleCommand(t *testing.T) {
	ref, llms := renderBoth(t)
	walkCommandTree(NewRootCmd(), func(c *cobra.Command) {
		if c.Hidden {
			if strings.Contains(ref, "## "+c.CommandPath()) {
				t.Errorf("reference includes hidden command %q", c.CommandPath())
			}
			return
		}
		if !strings.Contains(ref, "## "+c.CommandPath()) {
			t.Errorf("reference omits %q", c.CommandPath())
		}
		if c.Parent() == nil {
			return
		}
		if !strings.Contains(llms, "`"+c.CommandPath()+"`") {
			t.Errorf("llms.txt omits %q", c.CommandPath())
		}
	})
}

// TestDocs_ReferenceCarriesExtendedBodies pins the reason this generator
// exists rather than cobra/doc's GenMarkdownTree: the extended help is
// the substance, and a reference without it would just be a flag list.
func TestDocs_ReferenceCarriesExtendedBodies(t *testing.T) {
	ref, _ := renderBoth(t)
	walkCommandTree(NewRootCmd(), func(c *cobra.Command) {
		if c.Hidden {
			return
		}
		ext := extendedHelpFor(c)
		if ext == "" {
			return
		}
		// Compare on the body's first line: the whole body is fenced, and
		// a substring match on it would be brittle against wrapping.
		first := strings.TrimSpace(strings.SplitN(ext, "\n", 2)[0])
		if first != "" && !strings.Contains(ref, first) {
			t.Errorf("reference omits %q's extended body (first line: %q)",
				c.CommandPath(), first)
		}
	})
}

// TestDocs_GeneratedFilesCarryTheBanner pins the "do not edit" marker and
// the regeneration command in it, so a reader who opens the file is told
// what to run instead of making an edit that will not survive.
func TestDocs_GeneratedFilesCarryTheBanner(t *testing.T) {
	ref, _ := renderBoth(t)
	if !strings.HasPrefix(ref, "<!-- AUTO-GENERATED") {
		t.Errorf("cli-reference.md does not open with the generated banner")
	}
	if !strings.Contains(ref, "make docs") {
		t.Errorf("the banner does not name the regeneration command")
	}
}

// TestDocs_WritesBothFilesUnderDir exercises the verb end to end, which
// is what `make docs` actually runs.
func TestDocs_WritesBothFilesUnderDir(t *testing.T) {
	dir := t.TempDir()
	if _, _, err := runCmd(t, "docs", "--dir", dir); err != nil {
		t.Fatalf("docs: %v", err)
	}
	for _, rel := range []string{filepath.Join("docs", "cli-reference.md"), "llms.txt"} {
		body, err := os.ReadFile(filepath.Join(dir, rel))
		if err != nil {
			t.Fatalf("%s was not written: %v", rel, err)
		}
		if len(body) == 0 {
			t.Errorf("%s is empty", rel)
		}
	}
}

// TestDocs_LLMsTxtListItemsUnderH2AreLinks pins the llmstxt.org rule that
// file-list items appearing under an H2 must be hyperlinks, not bare
// text or backticked names. Once a "## " heading has been seen, every
// "- " list item's content must open with a markdown link.
func TestDocs_LLMsTxtListItemsUnderH2AreLinks(t *testing.T) {
	_, llms := renderBoth(t)
	seenH2 := false
	for _, ln := range strings.Split(llms, "\n") {
		if strings.HasPrefix(ln, "## ") {
			seenH2 = true
			continue
		}
		if !seenH2 || !strings.HasPrefix(ln, "- ") {
			continue
		}
		item := strings.TrimPrefix(ln, "- ")
		if !strings.HasPrefix(item, "[") || !strings.Contains(item, "](") {
			t.Errorf("list item under an H2 is not a markdown link: %q", ln)
		}
	}
}

// TestDocs_IsHiddenFromTheCommandList pins that the generator does not
// clutter the surface a user reads to learn the tool.
func TestDocs_IsHiddenFromTheCommandList(t *testing.T) {
	stdout, _ := runHelp(t, "--help")
	if strings.Contains(stdout, "docs ") {
		t.Errorf("`docs` appears in the root command list; it addresses "+
			"maintainers and should stay hidden:\n%s", stdout)
	}
	// Hidden, but still reachable and documented for whoever runs it.
	if _, _, err := runCmd(t, "docs", "--help"); err != nil {
		t.Errorf("`docs --help` is not reachable: %v", err)
	}
}
