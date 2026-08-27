package table

// A tag-declaration refusal must carry the line its
// `[tags.<key>]` header sits on, so a CLI finding can locate the defect
// instead of pointing every refusal at line 1.
//
// SCOPE OF THE POSITION. go-toml v2 surfaces a position only for its OWN
// decode failures; the semantic checks in `tagDecl` run against already
// decoded structs, where nothing of the source survives. Rather than plumb
// a position through every `fail` site, the loader recovers the ONE
// position it can honestly recover — the `[tags.<key>]` header — by scanning
// the source bytes it already holds. Every other refusal keeps `Line` zero,
// and zero means unknown, never line 1.
//
// The scan is deliberately conservative. An inline (`tags = { … }`) or
// dotted-key authoring has no `[tags.<key>]` header to find, and a document
// could in principle present a key more than once; in both cases the scan
// declines rather than guessing, because a confidently WRONG line is worse
// for an author than an honest absence.

import (
	"errors"
	"strings"
	"testing"
)

// tagLineDoc is a valid model but for `[tags.flavors]`, which declares a
// `domain` its `set` kind does not take. The header sits well below line 1
// so a hardcoded fallback cannot pass by coincidence.
const tagLineDoc = `
outcomes = ["go"]

[model]
id = "line-probe"
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

// failureOf returns the *Failure a load refusal carries.
func failureOf(t *testing.T, err error) *Failure {
	t.Helper()
	if err == nil {
		t.Fatal("the document loaded; want a refusal")
	}
	var f *Failure
	if !errors.As(err, &f) {
		t.Fatalf("refusal is not a *Failure: %v", err)
	}
	return f
}

// TestTagDeclRefusalCarriesItsHeaderLine pins that a malformed
// `[tags.<key>]` declaration reports the line of its own header.
//
// FAILURE MODE GUARDED: a refusal that carries no position at all, forcing
// every consumer to attribute the defect to line 1 regardless of where the
// author wrote it.
func TestTagDeclRefusalCarriesItsHeaderLine(t *testing.T) {
	// The header's line is derived from the fixture rather than written as a
	// literal, so editing the fixture cannot silently decouple the two.
	want := 0
	for i, line := range strings.Split(tagLineDoc, "\n") {
		if strings.TrimSpace(line) == "[tags.flavors]" {
			want = i + 1
		}
	}
	if want <= 1 {
		t.Fatalf("fixture must place [tags.flavors] below line 1; found %d", want)
	}

	_, err := Load([]byte(tagLineDoc), "line-probe.toml")
	f := failureOf(t, err)
	if f.Category != CatMalformedTagDeclaration {
		t.Fatalf("category = %q; want %q", f.Category, CatMalformedTagDeclaration)
	}
	if f.Line != want {
		t.Fatalf("Line = %d; want %d — the refusal must name the line its "+
			"[tags.flavors] header sits on", f.Line, want)
	}
}

// TestRefusalsWithNoTagHeaderCarryNoLine pins the honest-absence half of the
// contract: a refusal the header scan cannot ground leaves `Line` zero.
//
// FAILURE MODE GUARDED: a scan that reaches for the nearest plausible header
// and stamps a line onto a refusal that has nothing to do with it. Zero is
// the only correct answer where no `[tags.<key>]` header identifies the
// defect; a fabricated line sends the author to innocent source text.
func TestRefusalsWithNoTagHeaderCarryNoLine(t *testing.T) {
	cases := []struct {
		name string
		src  string
	}{
		{
			name: "malformed TOML never reaches the tag table",
			src:  "outcomes = [\n",
		},
		{
			name: "a model-header refusal names no tag",
			src:  "outcomes = [\"go\"]\n\n[model]\nid = \"no-version\"\n",
		},
		{
			name: "an inline tags authoring has no [tags.<key>] header",
			src: "outcomes = [\"go\"]\n\n[model]\nid = \"inline\"\nversion = 1\n\n" +
				"[tags]\nrecognized = { provenance = \"recognized\", kind = \"enum\", " +
				"domain = [\"go\"], single_valued = true, required = true }\n" +
				"flavors = { provenance = \"owned\", kind = \"set\", domain = [\"a\"] }\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := Load([]byte(c.src), "probe.toml")
			f := failureOf(t, err)
			if f.Line != 0 {
				t.Fatalf("Line = %d; want 0 — this refusal has no "+
					"[tags.<key>] header to ground a line on, and a "+
					"fabricated line is worse than none", f.Line)
			}
		})
	}
}
