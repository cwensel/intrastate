package cli

// Kata ywwq — `flow read-state` accepts `--tag`, so a command reader whose
// argv names `{tag.<key>}` (RDR 0028 `0028:C1.6`) can be read through it.
//
// read-state runs EVERY declared reader, so before this a model carrying a
// `{tag.<key>}` reader could never be read here at all: the flag did not
// exist, and without it the placeholder was always unbound. The tags take
// the same path `resolve`, `next` and `set-state` take — `parseTags`, then
// the request's `artifactMap` — so their refusals are the same ones. A
// `--tag` refusal lands before any reader runs; an unbound or flag-shaped
// placeholder refuses per use site, before the affected reader spawns, so
// a reader ahead of it in the run may already have run.
//
// The fixture is q14r's fake `kata`, whose reader argv is
// `[kata, "show", "{tag.id}", "--json"]` and whose log records every argv
// the tool receives. It declares that one reader only, so here every
// refusal leaves the log empty.

import (
	"slices"
	"testing"
)

// readStateArgs is a read-state invocation against the fake tool with the
// given extra flags appended.
func (f fakeKata) readStateArgs(extra ...string) []string {
	return append([]string{"flow", "read-state", "--model", f.model,
		"--artifact", artifactBinding("kata", f.dir), "--as=json"}, extra...)
}

func TestReadStateYwwq_TagBindsACommandReadersArgv(t *testing.T) {
	f := newFakeKata(t, "lifecycle:queued")

	stdout := requireSuccess(t, f.readStateArgs("--allow-commands", "--tag", "id=ywwq")...)

	if got := f.read(t, "log"); !slices.Equal(got, []string{"show ywwq --json"}) {
		t.Fatalf("the tool saw %q; want the bound id substituted into the reader's argv", got)
	}

	readers, _ := flowData(t, stdout)["readers"].([]any)
	if len(readers) != 1 {
		t.Fatalf("readers = %v; want the one declared reader", readers)
	}
	reader := readers[0].(map[string]any)
	if reader["id"] != "kata" {
		t.Errorf("reader id = %v; want kata", reader["id"])
	}
	if keys, _ := reader["keys"].([]any); !slices.Equal(keys, []any{"lifecycle"}) {
		t.Errorf("reader keys = %v; want the declared [lifecycle]", keys)
	}
	tags, _ := reader["tags"].(map[string]any)
	if tags["lifecycle"] != "queued" {
		t.Errorf("reader tags = %v; want lifecycle = queued from the tool", tags)
	}
	if _, echoed := tags["id"]; echoed {
		t.Errorf("reader tags = %v; --tag is context for the argv and is never "+
			"reported as read state", tags)
	}
}

func TestReadStateYwwq_AnUndeclaredCarrierChangesNothing(t *testing.T) {
	f := newFakeKata(t, "lifecycle:queued")
	base := requireSuccess(t, f.readStateArgs("--allow-commands", "--tag", "id=ywwq")...)
	with := requireSuccess(t, f.readStateArgs("--allow-commands", "--tag", "id=ywwq",
		"--tag", "carrier=anything")...)

	if base != with {
		t.Errorf("an undeclared carrier changed the report:\n%s\nvs\n%s", base, with)
	}
}

// Every refusal is about the request, exits 2, and lands before the tool
// runs — asserted by the tool's empty log, not by the refusal's text.
func TestReadStateYwwq_TagRefusalsPrecedeAnySpawn(t *testing.T) {
	for _, tc := range []struct {
		name  string
		extra []string
		code  string
	}{
		{"no --tag leaves the placeholder unbound",
			[]string{"--allow-commands"}, codeRequestRefused},
		{"a flag-shaped bound value",
			[]string{"--allow-commands", "--tag", "id=-x"}, codeRequestRefused},
		{"an owned key",
			[]string{"--allow-commands", "--tag", "id=ywwq", "--tag", "lifecycle=queued"}, codeTagOwned},
		{"the reserved key",
			[]string{"--allow-commands", "--tag", "id=ywwq", "--tag", "recognized=claim"}, codeTagReserved},
		{"a repeated key",
			[]string{"--allow-commands", "--tag", "id=a", "--tag", "id=b"}, codeTagDuplicate},
		{"a malformed entry",
			[]string{"--allow-commands", "--tag", "noequals"}, codeTagInvalid},
		{"no --allow-commands",
			[]string{"--tag", "id=ywwq"}, codeRequestRefused},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFakeKata(t, "lifecycle:queued")

			requireRefusal(t, tc.code, 2, f.readStateArgs(tc.extra...)...)
			if got := f.read(t, "log"); len(got) != 0 {
				t.Errorf("the tool ran %q; every refusal precedes the spawn", got)
			}
		})
	}
}
