package cli

// A load finding must locate the defect at the line the loader
// attributed it to, not at line 1 unconditionally.
//
// `loadFindings` hardcoded `path + ":1"` because the loader carried no
// position. The loader now stamps a line on the refusals it can honestly
// ground (a `[tags.<key>]` header) and leaves the rest at zero, so the CLI
// renders the real line where there is one and keeps `:1` as the DOCUMENTED
// fallback everywhere else — a reader still needs the file to act on.

import (
	"strconv"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/table"
)

// locatorDoc declares `[tags.flavors]` as a `set` carrying a `domain`, a
// kind/field mismatch, with the header placed well below line 1.
const locatorDoc = `
outcomes = ["go"]

[model]
id = "locator-probe"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
domain = ["go"]
single_valued = true
required = true

[tags.flavors]
provenance = "owned"
kind = "set"
domain = ["a", "b"]
`

// TestLoadFindingLocatesTheDeclarationsOwnLine pins that the finding's
// locator carries the tag header's line.
//
// FAILURE MODE GUARDED: `Locator: path + ":1"`, which points every load
// finding at the first line of the file whatever the author wrote where.
func TestLoadFindingLocatesTheDeclarationsOwnLine(t *testing.T) {
	want := 0
	for i, line := range strings.Split(locatorDoc, "\n") {
		if strings.TrimSpace(line) == "[tags.flavors]" {
			want = i + 1
		}
	}
	if want <= 1 {
		t.Fatalf("fixture must place [tags.flavors] below line 1; found %d", want)
	}

	_, err := table.Load([]byte(locatorDoc), "flows/probe.toml")
	if err == nil {
		t.Fatal("the fixture loaded; want a malformed_tag_declaration refusal")
	}
	findings := loadFindings("flows/probe.toml", err)
	if len(findings) != 1 {
		t.Fatalf("findings = %d; want exactly 1", len(findings))
	}
	got := findings[0].Locator
	if wantLoc := "flows/probe.toml:" + strconv.Itoa(want); got != wantLoc {
		t.Fatalf("locator = %q; want %q — the finding must name the line "+
			"the loader attributed the defect to", got, wantLoc)
	}
}

// TestLoadFindingFallsBackToLineOne pins the fallback. A refusal the loader
// could not ground to a line still reports `<path>:1`: the file is what a
// reader acts on, and the finding record carries a locator either way.
//
// FAILURE MODE GUARDED: the line-aware branch inventing a position for a
// refusal that has none — a locator pointing at innocent source text is a
// worse answer than the honest, documented `:1`.
func TestLoadFindingFallsBackToLineOne(t *testing.T) {
	// A model header carrying no version: refused before any tag is read,
	// so no `[tags.<key>]` header identifies the defect.
	const noVersion = "outcomes = [\"go\"]\n\n[model]\nid = \"no-version\"\n"

	_, err := table.Load([]byte(noVersion), "flows/probe.toml")
	if err == nil {
		t.Fatal("a versionless model loaded; want a refusal")
	}
	findings := loadFindings("flows/probe.toml", err)
	if len(findings) != 1 {
		t.Fatalf("findings = %d; want exactly 1", len(findings))
	}
	if got := findings[0].Locator; got != "flows/probe.toml:1" {
		t.Fatalf("locator = %q; want %q — an ungrounded refusal keeps the "+
			"documented file-only fallback", got, "flows/probe.toml:1")
	}
}
