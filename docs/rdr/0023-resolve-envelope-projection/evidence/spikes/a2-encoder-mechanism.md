# A2 spike — encoder mechanism for the C1 wire shape

Date: 2026-08-29

Claim tested: pointer-valued echo fields (`model`, `observed`, `owned`,
`readers`, `outcome`) with `,omitempty` — populated in default mode, nilled
at projection — reproduce the shipped `flow resolve` bytes exactly in
default mode and drop the echo keys cleanly under projection.

Verdict: **PASS** — the candidate mechanism satisfies every C1 clause.

## Encoder settings

The one wire encoder is `internal/cli/clierr/clierr.go` `WriteJSONLine`
(lines 174–178): `json.NewEncoder` + `enc.SetEscapeHTML(false)`, one NDJSON
line per call. `internal/cli/respond/respond.go` `writeJSONLine` (line 241)
routes the success envelope through it. The spike program uses identical
settings.

## Reference output (real CLI)

Built from repo HEAD (`go build ./cmd/intrastate`). Model:
`models/examples/pricing-decision-table.toml` (observed enums `tier`,
`region`, both `required = true`; outcome alphabet `["decide"]`).

```
intrastate flow resolve --model models/examples/pricing-decision-table.toml \
    --tag tier=paid --tag region=eu --outcome decide --as json
```

```json
{"type":"ok","data":{"model":"models/examples/pricing-decision-table.toml","revision":"","observed":{"region":"eu","tier":"paid"},"owned":{},"readers":[],"outcome":"decide","rule":"paid-eu","gates":[],"emit":{"dpa":"required","plan":"pro"},"next":{},"writes":{},"clear":[],"escaped":false}}
```

A tag-free success is not achievable with this model (`tier`/`region` are
required), but the same line already demonstrates the empty-`{}` case on
real output: a decision table declares no owned tags, so the shipped
`owned` is `{}` — present, not dropped.

Projected reference (jq key-deletion of the five echo keys):

```
jq -c 'del(.data.model,.data.observed,.data.owned,.data.readers,.data.outcome)' ref.json
```

```json
{"type":"ok","data":{"revision":"","rule":"paid-eu","gates":[],"emit":{"dpa":"required","plan":"pro"},"next":{},"writes":{},"clear":[],"escaped":false}}
```

## Spike program

Own module (`go mod init spike`), run as
`go run . ref.json ref-projected.json`. S2 is populated by DECODING the
reference line, so every compared value is the real call's value.

```go
// Spike for RDR 0023 A2: can pointer-valued echo fields with `,omitempty`
// reproduce the shipped `flow resolve` wire bytes in default mode and drop
// the echo keys cleanly under projection?
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type gateResult struct {
	ID     string `json:"id"`
	Result string `json:"result"`
	Reason string `json:"reason,omitempty"`
}

// envelope mirrors respond.Success (internal/cli/respond/respond.go).
type envelope struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

// s1Payload is the shipped resolvePayload verbatim
// (internal/cli/flow_resolve.go): field order, types, and tags.
type s1Payload struct {
	Model       string            `json:"model"`
	Revision    string            `json:"revision"`
	Observed    map[string]string `json:"observed"`
	Owned       map[string]string `json:"owned"`
	Readers     []string          `json:"readers"`
	Outcome     string            `json:"outcome"`
	Rule        string            `json:"rule"`
	Gates       []gateResult      `json:"gates"`
	Emit        map[string]string `json:"emit"`
	Next        map[string]string `json:"next"`
	Writes      map[string]string `json:"writes"`
	Clear       []string          `json:"clear"`
	Escaped     bool              `json:"escaped"`
	EscapeClass string            `json:"escape_class,omitempty"`
}

// s2Payload is the candidate mechanism: the five echo fields (model,
// observed, owned, readers, outcome) become pointer-valued with
// `,omitempty`; every other field is untouched. Declaration order is
// identical to s1Payload.
type s2Payload struct {
	Model       *string            `json:"model,omitempty"`
	Revision    string             `json:"revision"`
	Observed    *map[string]string `json:"observed,omitempty"`
	Owned       *map[string]string `json:"owned,omitempty"`
	Readers     *[]string          `json:"readers,omitempty"`
	Outcome     *string            `json:"outcome,omitempty"`
	Rule        string             `json:"rule"`
	Gates       []gateResult       `json:"gates"`
	Emit        map[string]string  `json:"emit"`
	Next        map[string]string  `json:"next"`
	Writes      map[string]string  `json:"writes"`
	Clear       []string           `json:"clear"`
	Escaped     bool               `json:"escaped"`
	EscapeClass string             `json:"escape_class,omitempty"`
}

// s3Payload is the hazard the RDR names: bare NON-pointer maps under
// `,omitempty` — an empty (or nil) map is dropped from the wire.
type s3Payload struct {
	Model       string            `json:"model,omitempty"`
	Revision    string            `json:"revision"`
	Observed    map[string]string `json:"observed,omitempty"`
	Owned       map[string]string `json:"owned,omitempty"`
	Readers     []string          `json:"readers,omitempty"`
	Outcome     string            `json:"outcome,omitempty"`
	Rule        string            `json:"rule"`
	Gates       []gateResult      `json:"gates"`
	Emit        map[string]string `json:"emit"`
	Next        map[string]string `json:"next"`
	Writes      map[string]string `json:"writes"`
	Clear       []string          `json:"clear"`
	Escaped     bool              `json:"escaped"`
	EscapeClass string            `json:"escape_class,omitempty"`
}

// encode uses the CLI's exact wire settings: json.Encoder with
// SetEscapeHTML(false) (internal/cli/clierr/clierr.go WriteJSONLine).
func encode(v any) []byte {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		panic(err)
	}
	return bytes.TrimRight(buf.Bytes(), "\n")
}

func mustRead(path string) []byte {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes.TrimRight(b, "\n")
}

func report(name string, got, want []byte) bool {
	same := bytes.Equal(got, want)
	fmt.Printf("%-34s %v\n", name+":", verdict(same))
	if !same {
		fmt.Printf("  got:  %s\n  want: %s\n", got, want)
	}
	return same
}

func verdict(ok bool) string {
	if ok {
		return "BYTE-IDENTICAL"
	}
	return "MISMATCH"
}

func main() {
	ref := mustRead(os.Args[1])          // shipped CLI output
	refProjected := mustRead(os.Args[2]) // jq del() of the five echo keys

	// Decode the reference into the candidate struct so every value is the
	// real call's value, not a hand-typed copy.
	var env struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(ref, &env); err != nil {
		panic(err)
	}
	var s2 s2Payload
	if err := json.Unmarshal(env.Data, &s2); err != nil {
		panic(err)
	}

	ok := true

	// (a) default mode: all echo pointers populated.
	ok = report("S2 default vs shipped reference",
		encode(envelope{Type: "ok", Data: s2}), ref) && ok

	// (b) projection: echo pointers nilled -> keys absent, order kept.
	proj := s2
	proj.Model, proj.Observed, proj.Owned, proj.Readers, proj.Outcome =
		nil, nil, nil, nil, nil
	projBytes := encode(envelope{Type: "ok", Data: proj})
	ok = report("S2 projected vs jq key-deletion", projBytes, refProjected) && ok
	fmt.Printf("projected line: %s\n", projBytes)
	for _, k := range []string{`"model"`, `"observed"`, `"owned"`, `"readers"`, `"outcome"`} {
		if strings.Contains(string(projBytes), k) {
			fmt.Printf("  FAIL: projected line still carries %s\n", k)
			ok = false
		}
	}

	// {} preservation: a non-nil pointer to an EMPTY map under `,omitempty`
	// still renders "observed":{} — the mechanism distinguishes empty from
	// projected (nil pointer -> key absent).
	empty := map[string]string{}
	type tiny struct {
		Observed *map[string]string `json:"observed,omitempty"`
	}
	presentEmpty := encode(tiny{Observed: &empty})
	absent := encode(tiny{Observed: nil})
	fmt.Printf("pointer to empty map:              %s\n", presentEmpty)
	fmt.Printf("nil pointer:                       %s\n", absent)
	if string(presentEmpty) != `{"observed":{}}` || string(absent) != `{}` {
		fmt.Println("  FAIL: {}-preservation property does not hold")
		ok = false
	}

	// Hazard: the same payload as a NON-pointer struct with `,omitempty` —
	// the empty owned/next/... maps are DROPPED, diverging from shipped bytes.
	var s3 s3Payload
	if err := json.Unmarshal(env.Data, &s3); err != nil {
		panic(err)
	}
	s3Bytes := encode(envelope{Type: "ok", Data: s3})
	fmt.Printf("S3 (non-pointer omitempty) line:   %s\n", s3Bytes)
	if bytes.Equal(s3Bytes, ref) {
		fmt.Println("  UNEXPECTED: hazard did not reproduce (S3 == reference)")
		ok = false
	} else {
		fmt.Println("hazard confirmed: non-pointer omitempty drops empty {} (byte diff vs shipped)")
	}

	if !ok {
		os.Exit(1)
	}
	fmt.Println("ALL ASSERTIONS PASS")
}
```

## Results

Program output, verbatim (exit 0):

```
S2 default vs shipped reference:   BYTE-IDENTICAL
S2 projected vs jq key-deletion:   BYTE-IDENTICAL
projected line: {"type":"ok","data":{"revision":"","rule":"paid-eu","gates":[],"emit":{"dpa":"required","plan":"pro"},"next":{},"writes":{},"clear":[],"escaped":false}}
pointer to empty map:              {"observed":{}}
nil pointer:                       {}
S3 (non-pointer omitempty) line:   {"type":"ok","data":{"model":"models/examples/pricing-decision-table.toml","revision":"","observed":{"region":"eu","tier":"paid"},"outcome":"decide","rule":"paid-eu","gates":[],"emit":{"dpa":"required","plan":"pro"},"next":{},"writes":{},"clear":[],"escaped":false}}
hazard confirmed: non-pointer omitempty drops empty {} (byte diff vs shipped)
ALL ASSERTIONS PASS
```

Clause by clause:

- **(a) default-mode byte identity**: S2 with all echo pointers populated
  marshals byte-identical to the shipped reference line, including
  `"owned":{}` (non-nil pointer to an empty map survives `,omitempty` —
  omitempty tests pointer nilness, not pointee emptiness).
- **(b) projection**: nilling the five echo pointers yields exactly the
  reference minus the five keys — absent, not null, not empty; surviving
  keys keep declaration order and are byte-identical to the jq
  key-deletion of the reference.
- **Empty vs projected are distinguishable**: `&map[string]string{}` under
  `,omitempty` renders `"observed":{}` (present); `nil` renders the key
  absent.
- **Named hazard reproduced**: S3 (non-pointer maps with `,omitempty`)
  drops the empty `"owned":{}` and `"readers":[]` from the shipped line —
  a byte diff against today's output, so bare `omitempty` on the existing
  field types cannot be the mechanism.

## Text-mode determinism (c)

Same successful call with `--as=text`, 10 runs, sha256 of each full
output, `sort -u`:

```
49ebd9c187c14f2a0e82531089d74e3bf6646dac701a82b777623bada55bfdc8
```

One unique hash — all 10 runs byte-identical. Determinism comes from
`internal/cli/respond/text.go` `flatten`, which sorts each decoded
object's keys before emitting leaf lines (`sort.Strings(keys)`,
text.go:75), and from the shared non-HTML-escaping encoder
(`marshalCanonical`, text.go:131–139).
