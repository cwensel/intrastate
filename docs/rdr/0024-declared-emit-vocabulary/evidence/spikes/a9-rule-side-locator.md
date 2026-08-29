Model: claude-opus-5[1m]

# A9 rule-side leg — is an emit-refusal SOURCE LINE recoverable without new plumbing?

Stage 6 (Reconcile) source spike, RDR 0024. Read-only; no file was edited.

Operative claim under test: "Phase 1's rule-side `emitHeaderLine` is writable
WITHOUT new plumbing."

VERDICT: **CONFIRMS.** A rule-side source line is recoverable at HEAD from data
the loader already holds. It needs a NEW HELPER FUNCTION, but NO new data
threaded through the loader.

## 1. Does the loader hold the raw source where rules are processed?

Yes.

- Field: `internal/table/load.go:68` — `loader.src []byte`, documented as kept
  "for the ONE thing the decoded structs cannot supply: where in the source a
  declaration was written."
- Populated once at `internal/table/load.go:59` (`ld := &loader{doc: &doc, src: src, sourceID: sourceID}`).
- Rule processing is a `loader` METHOD, so `l.src` is in scope unchanged:
  - `internal/table/normalize.go:282` — `func (l *loader) normalizeRules() error`
  - `internal/table/normalize.go:335` — `func (l *loader) normalizeRule(rule *sourceRule, id string, setKeys []string) ([]Row, error)`
  - `internal/table/normalize.go:457` — `Emit: emitSequence(rule.Emit)` — the
    single site that reads `sourceRule.Emit`, inside `normalizeRule`, which is
    a `loader` method.

So the step that reads `sourceRule.Emit` already has both `l.src` and the rule's
identity in scope. Nothing new has to be threaded.

## 2. What identifies a rule at source level?

`internal/table/source.go:79-103` (`sourceRule`) carries `ID *string \`toml:"id"\``.
It is REQUIRED and UNIQUE, enforced before any emit work runs:

- `internal/table/normalize.go:290-292` — absent/empty id refuses `CatMalformedRuleShape`.
- `internal/table/normalize.go:303-306` — `seenID[id]` refuses `CatDuplicateRuleID`
  ("Compared by exact byte equality").

Two usable identities are therefore in hand at the emit site:
1. the rule id string (`id`, unique per document by construction), and
2. the array-of-tables index `i` from `internal/table/normalize.go:288` (`for i := range l.doc.Rule`).

`sourceRule` carries NO line/offset field, and none is needed (see §4/§6).

## 3. Is `[rule.emit]` the authored shape? (the crux)

Yes — and its header text is IDENTICAL for every rule, so a naive
`tagHeaderLine`-style scan for the header alone CANNOT work.

Authored shape is standard `[[rule]]` array-of-tables with a `[rule.emit]`
sub-table header on its own line, e.g. rdr-status.toml:342-352:

    [[rule]]
    id = "locate-demoted"
    [rule.match.recognized]
    eq = "locate"
    [rule.guard.all.status]
    eq = "Demoted"
    [rule.emit]
    next    = "none"
    ...

Not `[[rule.emit]]`, not an inline table.

Counts (commands in §7):
- `models/rdr-status.toml`: 63 occurrences of the literal line `[rule.emit]`;
  63 rule `id = ` lines; 0 duplicate id lines.
- `models/examples/pricing-decision-table.toml`: 4 rules, all `[[rule]]` +
  `[rule.emit]` (headers at lines 59, 71, 83, 95).

CONSEQUENCE: `tagHeaderLine`'s "exactly one match or return 0" rule
(`internal/table/load.go:181-204`) would return 0 for `[rule.emit]` in EVERY
real multi-rule model. A generic header scan is not the mechanism. The rule id
is.

## 4. Does pelletier/go-toml/v2 expose position information?

Only for its OWN failures, never for a successfully decoded document.

- `go.mod:6` — `github.com/pelletier/go-toml/v2 v2.2.4`.
- `errors.go:18-25` — `DecodeError{message, line, column, key, human}`; the
  position fields are UNEXPORTED.
- `errors.go:71` — `func (e *DecodeError) Position() (row, column int)`.
- `strict.go` — `StrictMissingError{Errors []DecodeError}`; grep for
  `Position|Row|Line` in strict.go returns nothing of its own.

There is no `Meta`, no per-value row table, and no decode-time position hook on
`toml.Decoder`. `internal/table/source.go:131-142` (`decodeStrict`) already
consumes the only two error shapes available.

DECISIVE: a rule-side emit refusal is decided AFTER a successful strict decode,
against `sourceRule.Emit` (a plain `map[string]string`). The decoder is silent
by then. This is exactly what `internal/table/category.go:103-115` documents:
"the decoder itself reports a position only for its own syntax failures."

So the decoder is not the route — but it does not need to be, because the raw
bytes are already retained on purpose (§1).

## 5. Existing precedent for line recovery

- `internal/table/load.go:181` — `tagHeaderLine(src []byte, key string) int`.
  Textual scan, strips a trailing `#` comment and surrounding space, exact-matches
  a bare or `strconv.Quote`d bracketed header, returns 0 unless exactly one line
  matches.
- `internal/table/category.go:151-170` — `atLine(err error, line int) error`,
  a post-hoc stamp onto `Failure.Line` (`category.go:115`). Its doc comment states
  the design rule directly: "the position is recoverable at ONE seam — the caller
  that still holds the source bytes."
- Wired at exactly ONE loader site: `internal/table/load.go:214`, inside `loadTags`.
- No other line-recovery helper exists in `internal/`. `dumpLineFor`
  (`internal/table/emit_0010_test.go:673`) and `helpLineFor`
  (`internal/cli/lint_0006_test.go:825`) are TEST helpers that scan rendered
  OUTPUT text, not source. No graphlint/resolve/cli analogue.

Precedent is therefore a pattern to follow (scan `l.src`, stamp via `atLine`),
not a reusable function for the nested case.

## 6. The concrete mechanism

A rule-side `emitHeaderLine(src []byte, ruleID string) int` is writable using
only `l.src` and `id`, both already in scope at `normalize.go:335`:

1. ANCHOR on the rule id line, not the emit header. The anchor text is
   `id = "<ruleID>"` (normalized for whitespace/quoting the way `tagHeaderLine`
   already normalizes: strip trailing `#` comment, `TrimSpace`). The id is
   unique per document by `CatDuplicateRuleID` (`normalize.go:303`), so this
   anchor matches at most once among rules.
   Caveat verified: `[model]` also carries an `id` key
   (pricing-decision-table.toml:26-27), so the anchor must match the id's VALUE,
   not the bare `id = ` key. Matching `id = "<ruleID>"` is inherently
   value-scoped and so is already safe; retaining `tagHeaderLine`'s
   "exactly one match or return 0" discipline makes a collision with a
   same-valued `[model] id` decline rather than mispoint.
2. SCAN FORWARD from the anchor for the first line equal to `[rule.emit]`,
   stopping at the next top-level `[[rule]]` / `[` header so the scan cannot
   run into the following rule's block.
3. Return 0 on anything ambiguous, matching `tagHeaderLine`'s existing contract
   and `Failure.Line`'s documented "zero means UNKNOWN" (`category.go:103-115`).
4. Stamp with the EXISTING `atLine` at the emit refusal site inside
   `normalizeRule`.

Ordering precondition VERIFIED mechanically: in both real models, every rule's
`id` line precedes that rule's `[rule.emit]` header — 63/63 in rdr-status.toml,
4/4 in pricing-decision-table.toml (§7). Where an author writes emit before id,
step 2 finds no match within the block and returns 0, which is the sanctioned
decline, not a wrong pointer.

To locate the offending KEY inside the block rather than the block header, the
same scan continues from the emit header to the first line whose key matches the
offending emit key — same inputs, no extra data.

## 7. Commands run

    rg -n "rule.emit|\[\[rule\]\]|^\[rule" models/examples/pricing-decision-table.toml
    rg -n "rule.emit|\[\[rule\]\]" <rdr-repo>/models/rdr-status.toml
    rg -n "normalizeRules" internal/ -l
    rg -n "Emit|sourceRule|func \(l \*loader\)|range l.doc.Rule" internal/table/normalize.go
    rg -n "emit" internal/table/normalize.go
    rg -ln "HeaderLine|headerLine|lineOf|sourceLine|LineFor" internal/
    rg -n "pelletier/go-toml" go.mod
    sed -n '1,110p' $(go env GOMODCACHE)/github.com/pelletier/go-toml/v2@v2.2.4/errors.go
    rg -n "type StrictMissingError" .../strict.go ; grep -n "Position\|Row\|Line" .../strict.go
    grep -c '^\[rule\.emit\]$' models/rdr-status.toml            # 63
    grep -c '^id = ' models/rdr-status.toml                      # 63
    grep '^id = ' models/rdr-status.toml | sort | uniq -d | wc -l # 0
    awk '/^\[\[rule\]\]/{r++; idl=0} /^id = /{if(!idl) idl=NR} /^\[rule\.emit\]$/{...}' \
      models/rdr-status.toml                                     # rules: 63  id-before-emit ok: 63
    awk '...' models/examples/pricing-decision-table.toml        # rules: 4   id-before-emit ok: 4
    grep -n '^id = \|^\[model\]\|^\[\[rule\]\]' models/examples/pricing-decision-table.toml

## 8. Anchors

- `internal/table/load.go:59` — loader constructed with `src`
- `internal/table/load.go:68` — `loader.src []byte`
- `internal/table/load.go:181` — `tagHeaderLine`, the pattern to follow
- `internal/table/load.go:214` — the ONE existing `atLine` wiring
- `internal/table/category.go:115` — `Failure.Line`
- `internal/table/category.go:160` — `atLine`, reusable as-is
- `internal/table/source.go:80` — `sourceRule.ID`
- `internal/table/source.go:102` — `sourceRule.Emit map[string]string`
- `internal/table/source.go:131` — `decodeStrict`
- `internal/table/normalize.go:288` — `for i := range l.doc.Rule`
- `internal/table/normalize.go:290-306` — id required, unique
- `internal/table/normalize.go:335` — `normalizeRule`, a loader method
- `internal/table/normalize.go:457` — the sole `sourceRule.Emit` read
- `<rdr-repo>/models/rdr-status.toml:342-352` — authored `[[rule]]` + `[rule.emit]`
- `models/examples/pricing-decision-table.toml:51-61` — same shape

## 9. Verdict on the operative claim

CONFIRMS. `emitHeaderLine` is writable without new plumbing.

The distinction that matters: this needs a NEW HELPER FUNCTION (a rule-id-anchored
scan, because `tagHeaderLine`'s single-match-on-header-text contract cannot
disambiguate 63 identical `[rule.emit]` lines). It does NOT need new DATA threaded
through the loader — no source index, no decode-captured line, no change to
`sourceRule`, no change to any signature. Both inputs (`l.src`, `id`) are already
present and in scope at `normalize.go:335`, and the stamping mechanism (`atLine`)
already exists and is reusable unchanged.

Residual (non-blocking, design note for Phase 1): the helper must anchor on the
rule id and bound its forward scan to the rule's own block, and must keep
`tagHeaderLine`'s decline-on-ambiguity discipline. An author who writes
`[rule.emit]` before `id` in the same block gets line 0 (file-only locator),
which the existing `Failure.Line` contract already sanctions.
