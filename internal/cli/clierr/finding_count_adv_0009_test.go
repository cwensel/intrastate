package clierr_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
)

// RDR 0009 Phase 3b — ADVERSARIAL: the shared-record widening.
//
// clierr.Finding.Count is NEW PUBLIC SURFACE on a record RDRs 0005, 0006
// and 0008 co-own (RDR 0009 deviations D3). The widening has two
// obligations that pull in opposite directions:
//
//  1. it must be INVISIBLE to every co-owner that does not set it, so no
//     shipped envelope, golden, or round-trip changes; and
//  2. it must be CARRIED by every shared renderer, so a producer that does
//     set it is not silently truncated on one output path.
//
// Obligation 1 holds. Obligation 2 does not: EmitFindingsText renders a
// finding through Finding.identitySuffix(), whose field list was not
// widened, so Count reaches JSON and is dropped from text — for THIS RDR
// and for any co-owner that adopts the field later. `0006:C14` fixes text
// mode as the completeness-preserving renderer ("a renderer that drops one
// has dropped a defect the author needs to see"); the widening opened a
// hole in exactly that renderer.

// --- ADV-3: obligation 2 — the shared text renderer must carry Count ----
func TestAdv0009_SharedTextRendererCarriesTheWidenedCountField(t *testing.T) {
	ce := &clierr.CLIError{
		Code:    "some-aggregate-failure",
		Message: "several subjects collapsed",
		Group:   clierr.GroupInternal,
		Findings: []clierr.Finding{
			{Code: "subject-defect", Message: "this subject is broken",
				Rule: "r1", Locator: "f.toml:1", Count: 7},
		},
	}

	var jsonOut bytes.Buffer
	clierr.EmitJSON(&jsonOut, ce)
	if !strings.Contains(jsonOut.String(), `"count":7`) {
		t.Fatalf("JSON dropped Count:\n%s", jsonOut.String())
	}

	var textOut bytes.Buffer
	clierr.EmitText(&textOut, ce)
	if !strings.Contains(textOut.String(), "7") {
		t.Errorf("Finding.Count is serialized to JSON but dropped by the shared "+
			"text renderer: Finding.identitySuffix()'s field list was not widened "+
			"alongside the new field, so a collapsed report loses its multiplicity "+
			"in the CLI's DEFAULT output mode. Every RDR co-owning this record "+
			"inherits the hole.\nJSON: %s\nTEXT: %s",
			strings.TrimSpace(jsonOut.String()), strings.TrimSpace(textOut.String()))
	}
}

// --- obligation 1 — the widening must not perturb the co-owners ---------
//
// These are regression guards on properties the implementation HOLDS. A
// producer that never sets Count must serialize byte-identically to what
// it did before the field existed, because the field is `omitempty` and
// its zero value is 0.
func TestAdv0009_TheCountWideningIsInvisibleToProducersThatDoNotSetIt(t *testing.T) {
	cases := []struct {
		name    string
		finding clierr.Finding
		want    string
	}{
		{
			name:    "RDR 0005 model-load finding",
			finding: clierr.Finding{Code: "model-load-failed", Message: "bad model", Locator: "flow.toml:12"},
			want:    `{"code":"model-load-failed","message":"bad model","locator":"flow.toml:12"}`,
		},
		{
			name: "RDR 0006 graph-lint finding",
			finding: clierr.Finding{
				Code: "graph-unprovable-coverage", Message: "cannot prove coverage",
				Model: "rdr", Severity: "error", Rule: "r1",
				Reason: "guard-undecidable", Fingerprint: "a|b",
			},
			want: `{"code":"graph-unprovable-coverage","message":"cannot prove coverage",` +
				`"model":"rdr","severity":"error","rule":"r1","reason":"guard-undecidable",` +
				`"fingerprint":"a|b"}`,
		},
		{
			name: "RDR 0008 reserved-key finding",
			finding: clierr.Finding{
				Code: "model-load-failed", Message: "reserved key", Key: "recognized",
				Block: "match", Class: "no_match",
			},
			want: `{"code":"model-load-failed","message":"reserved key",` +
				`"key":"recognized","block":"match","class":"no_match"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			if err := clierr.WriteJSONLine(&out, tc.finding); err != nil {
				t.Fatal(err)
			}
			if got := strings.TrimSuffix(out.String(), "\n"); got != tc.want {
				t.Errorf("the Count widening changed a co-owner's envelope:\n"+
					" got %s\nwant %s", got, tc.want)
			}
		})
	}
}

// A finding decoded from a pre-widening payload — one carrying no `count`
// key at all — must decode with Count == 0 and re-encode identically.
func TestAdv0009_APreWideningPayloadRoundTripsUnchanged(t *testing.T) {
	const payload = `{"code":"model-load-failed","message":"bad model","locator":"flow.toml:12"}`

	var f clierr.Finding
	if err := json.Unmarshal([]byte(payload), &f); err != nil {
		t.Fatal(err)
	}
	if f.Count != 0 {
		t.Fatalf("Count = %d on a payload carrying no count key, want 0", f.Count)
	}

	var out bytes.Buffer
	if err := clierr.WriteJSONLine(&out, f); err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSuffix(out.String(), "\n"); got != payload {
		t.Errorf("re-encode differs:\n got %s\nwant %s", got, payload)
	}
}
