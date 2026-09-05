# Repeatability-LITE diff: run-1 vs RDR 0020

Model: run-1 is an alternate-model (claude-fable-5-1) reconstruction of
0020's API, helpers, data model, and pseudo-code from C1/MVV/S1-S4 plus
widening into D-identity, D-wire-byte-format, D-naming,
D-selection-predicate, the technical-design/proposed-solution sections,
testing-strategy, and A1-A5. This is a single-run diff against the RDR
text, not a disagreement count.

## Admitted findings

### Finding 1 — declared-SET key given a literal empty string: which
refusal arm fires is internally unsettled (0020:C1 vs disposition table
vs illustrative code)

**RDR elements**: 0020:C1 (empty-value arm text), the Disposition table
rows "undeclared key, empty value" and "declared set key, empty value"
(§technical-design), and the Illustrative Code block (§technical-design).

**The under-determination**: C1's prose says the empty-value arm "MOVES
to the admission path ahead of the carrier branch" and applies "declared
or not" — i.e., it fires on `value == ""` unconditionally, before any
`m.Tags[key]` lookup. The Illustrative Code makes this literal:

```go
if value == "" { // empty is indistinguishable from unset, declared or not
    return nil, userErr(codeTagInvalid, key, "…was given an empty value")
}
decl, declared := m.Tags[key]
```

Read this way, a declared SET key given the raw empty string (`--tag
labels=`) would hit the hoisted arm and get the SCALAR-shaped message
("was given an empty value"), because the code never reaches
`m.Tags[key]` — let alone `canonicalValue` — before returning.

But the Disposition table lists "declared set key, empty value" as a
**separate row** with a **different** message: `flow-tag-invalid` "is
set-valued and takes a JSON array literal; got ". That message is
sourced from fixture `evidence/spikes/g-declared-set-empty.txt`, which
is HEAD (pre-change) behavior for exactly the input `--tag labels=`
(declared set `labels`, raw empty string) — and the RDR's own Desk Trace
never walks this fixture through the post-hoist code path to show it is
preserved. Testing Strategy scenario 4's aside says only that the two
messages "differ" and that the set one "stays inside the conformance
arms, declared-only" — asserting the outcome without reconciling it
against the illustrative code's unconditional `value == ""` check that
appears to make that arm unreachable for a byte-empty argument.

The RDR text therefore permits two readings that give different exit
messages for the identical concrete input `--tag labels=` against a
model that declares `labels` as a set:
1. The hoisted empty-value arm is truly unconditional (per C1 prose and
   the illustrative code as written) → declared-set-empty now gets the
   scalar message, changing HEAD behavior for that case (an unflagged,
   undiscussed behavior change).
2. The hoisted arm is implicitly scoped to only the cases where
   "empty" and "undeclared" coincide, or "empty value" in the
   disposition table secretly means something narrower than
   `value == ""` (e.g., an empty array literal after JSON parsing) →
   preserves fixture g, but then the illustrative code's guard is
   mis-drawn/misleading as written, and the RDR never states the
   narrower predicate.

**How run-1 resolved it**: run-1 flagged this explicitly as a GUESS and
picked reading 2 — "step 5 fires on the raw empty string for every key;
step 7's set-specific message remains for a declared set whose array
literal canonicalises to nothing (e.g. `[]`)." That is a plausible
peacemaking guess, but the RDR text does not state a `[]`-after-parse
predicate anywhere; run-1 invented it to reconcile two RDR passages that,
read as written, do not reconcile on their own. The RDR does not force
this resolution — a literal reading of the illustrative code forces
reading 1 instead, which is inconsistent with the Disposition table and
fixture g.

**Why load-bearing**: this is the one input class (declared set +
byte-empty value) where the RDR's own evidence (fixture g) and its own
illustrative code point at different wire messages, i.e. a real behavior
fork the RDR does not adjudicate, not a decomposition or naming choice.

## What the RDR DID determine unambiguously

- Identity/discriminator for declared vs. undeclared: byte-exact
  presence in `m.Tags`, two-value lookup, no folding, no sentinel
  (0020:D-identity) — run-1 reproduced this exactly with no guess.
- Wire behavior for the carrier itself: verbatim, byte-preserved,
  never canonicalized or kind/shape/domain-checked (0020:C1,
  D-wire-byte-format) — reproduced exactly.
- Precedence of the provenance guards (grammar → reserved → owned →
  duplicate) ahead of the empty-value/carrier splice, and that the
  empty-value arm sits after duplicate, not before
  (0020:D-selection-predicate) — reproduced exactly, including the
  "why" (would invert duplicate-vs-empty precedence otherwise).
- No new refusal code is minted; `flow-tag-undeclared` does not exist
  (0020:D-naming) — reproduced exactly.
- Disposition for a misspelled declared key: silent accept, diagnosis
  by echo only, an accepted open-world cost — reproduced exactly, not
  flagged as a gap.
- That declaredness need not be recorded on the admitted-tag data
  structure is explicit implementation latitude ("or passes the
  declaredness bit into it — implementation latitude"), not a silence;
  run-1's GUESS here is correctly scoped as omitted future-work, not
  reported as a finding.
- The `name=value` split grammar and the grammar-arm's `param` value
  are unchanged, pre-existing HEAD behavior explicitly out of this
  RDR's scope (C1 lists "the flag grammar" as unchanged); run-1's GUESS
  markers on these are false positives — not contract silences, just
  code this RDR does not touch.

## Verdict on determinism

The run reads as substantially determinate, not shared-model overfit:
every major contract element (identity rule, wire format, precedence
order, naming restraint, misspelling disposition) was reproduced exactly
without needing a guess, and most of run-1's own GUESS markers turned out
to be false positives once checked against sections the run had to widen
into (D-identity, D-selection-predicate, the disposition table). The one
admitted finding is a genuine internal tension in the RDR's own text
(illustrative code vs. disposition table vs. grounding fixture) rather
than an artifact of the reconstructing model filling in unconstrained
space — the RDR itself supplies conflicting cues here, and run-1 quietly
resolved the conflict instead of surfacing it.
