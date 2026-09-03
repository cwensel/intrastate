Model: Claude Opus 5 (1M context)

# Repeatability-LITE diff — 0028

Single alternate-model reconstruction (`run-1.md`, Claude Sonnet 5, variant
`lite (profile: large)`) compared against `kind=="C"` (0028:C1, clauses
C1.1–C1.6), 0028:MVV, and the S1–S31 proof list, widened into
§load-bearing-decisions, §performance-expectations, §failure-modes,
§risks-and-mitigations, §consequences, §illustrative-code and the twelve
Critical Assumptions wherever a candidate silence had no covering span.

The reconstruction is unusually faithful. Every normative behaviour it names —
the three-way carrier rule, the two asymmetric template dialects, `parse once`,
regexp-quoted tag substitution, pre-edit index holding, the re-anchor IDENTITY
check, the no-op inode arm, mode preservation and symlink resolution, the
one-buffer/one-write per-entry atomicity boundary and its non-transactionality
across entries, and the full C1.4 category set in registration order — is
recoverable from C1 and matches it. The findings below are the residue: places
where the reconstruction had to *choose* an algorithmic behaviour and the RDR
does not decide it.

## Admitted findings

### 1. Apply-time refusal precedence is unordered — `0028:C1` (C1.3 `order:`)

`C1.4 precedence:` gives a total fail-fast order for the five LOAD-time
categories, and scopes itself explicitly to lint ("evaluated after 0025:C5's
clauses 1–6"). No corresponding clause orders the APPLY-time refusals against
each other. C1.3 `order:` establishes only that all of them are decided before
any byte is written — a *timing* guarantee, not a *precedence* guarantee.

run-1's pseudo-code had to invent one: it places the `edit_value_multiline`
scan (its step 3) ahead of anchor selection (step 4), so an entry whose planned
value carries `\n` *and* whose anchor matches zero lines reports
`edit_value_multiline`. The opposite order reports `edit_anchor_unmatched`.
Both satisfy every sentence in C1. Nothing in §failure-modes,
§risks-and-mitigations or §performance-expectations narrows it —
§performance-expectations' determinism checklist covers map order for the
output buffer and for refusal reporting *across sibling tables*, and inherits
0025:C5's unspecified-order rule for that case only; it says nothing about the
order between distinct apply-time categories.

This is contract-visible: the Detail token a caller (an agent loop, per
`0028:A8`) branches on is different under the two orders. The same
under-determination applies between `edit_anchor_unmatched`,
`edit_anchor_ambiguous` and `edit_anchor_collision` across sibling rules — C1.3
`select:` names all three but does not say whether rule-by-rule cardinality is
resolved to completion before the cross-rule collision sweep runs, so an entry
with one unmatched rule and two colliding rules has two admissible answers.
run-1 chose per-rule-cardinality-first, then a collision sweep. S9's expectation
lists all three outcomes but pairs each with a single-defect fixture, so it does
not pin the composite. Note the asymmetry: `0028:S28` exists *precisely* because
"fail-fast is observable ONLY through a multi-defect input" — that reasoning is
applied to load time and not to apply time.

Either decide the apply-time order (and give it an S-scenario in the shape of
S28), or state explicitly that apply-time precedence is unspecified and that no
test may assert it, as C1.4 does for map-ranged siblings.

### 2. `<clear>` detection vs. the multi-line value scan — `0028:C1` (C1.2 `value shape:` / C1.5)

`<clear>` is a sentinel planned value, and C1.5 makes its handling depend on
whether the rule declared `clear = "line"`: declared ⇒ delete the line; absent ⇒
refuse `edit_clear_undeclared` pre-mutation. C1.2 `value shape:` independently
says a planned value containing `\n` or `\r` refuses `edit_value_multiline`
before mutation. Neither clause says whether the `<clear>` sentinel is
recognised before or after the value-shape scan, nor whether the sentinel is
recognised before or after `replace` template expansion.

The second half matters and is the sharper edge: C1.5 closes with "the literal
string `<clear>` is never substituted into `replace`", and `0028:S18` exercises
"the literal string `<clear>` reaching a `replace`" — but neither states the
mechanism. Two implementations are admissible. (a) The sentinel is tested
against the raw planned value *before* expansion, so a rule with
`clear = "line"` deletes the line and a rule without one refuses
`edit_clear_undeclared`. (b) The sentinel is only meaningful as a whole planned
value and a *literal* `<clear>` authored into a `replace` template's literal
segment passes through as bytes. run-1 chose (a) implicitly and placed the
`clear` branch inside the per-rule apply loop (its step 5), *after* selection
and after the value scan — which means under its reconstruction a planned
`<clear>` on a rule with no `clear` declaration reaches
`edit_clear_undeclared` only after anchor resolution has already succeeded, so
an unmatched anchor on such a rule would report `edit_anchor_unmatched`
instead. C1.5's zero-match-is-SUCCESS arm is stated for the `clear = "line"`
case only, and says nothing about the `clear`-absent case's interaction with
selection. This is finding 1's precedence question at a specific, tested
collision point (S17 fixtures all three arms singly, never composed).

### 3. Which side of the entry the value-shape scan covers is unstated for bound tags — `0028:C1` (C1.2 `value shape:`)

C1.2 says "a planned value or bound tag value containing `\n` or `\r` refuses
BEFORE mutation" and `0028:S10` proves both. But a bound tag value is used in
two distinct places: as `{tag.<key>}` inside an `anchor` (C1.2 `anchor admits:`,
regexp-quoted) and as a whole argv element in a `command` (C1.6). The RDR does
not say whether the multiline refusal applies to the tag value *per invocation*
(so a `\n`-bearing tag refuses even when no `edit` entry references it) or
*per use site* (so it refuses only when it reaches an anchor). run-1 chose
per-invocation, scanning `plannedValuesAndBoundTags(planned, art.Tags)`
wholesale at step 3, which refuses on tags the entry never uses. The narrower
reading is equally consistent: a regexp-quoted `\n` inside an anchor is
harmless to *line data* — it simply cannot match, since C1.3 `terminators:`
makes `\n` a terminator and never line content — and C1.6's argv path has no
line-data hazard at all, which is the stated ground for the refusal
("the single structural hazard of interpolation into line data"). Under the
per-use-site reading the refusal arguably should not fire for an anchor tag at
all. This lands on C1.2's own rationale, so it is a contract question, not
implementation taste. `0028:A1` settles that the map is invocation-wide, which
makes the scope question sharper rather than answering it.

## GUESSes dropped as false positives or as already-charted

Nine of run-1's marked GUESSes are not admissible. Recorded so the neutrality of
this pass is auditable:

1. **`Artifact.Tags` field name/type** — inadmissible as *naming*. `0028:A1`
   (Verified) settles the substance: one context map rides `Artifact`, it is
   invocation-wide not per-role, and `artifactMap()` is named as the sole
   construction site to change. The Go spelling is not a contract.
2. **`ReadBinding.Read` return shape** — false positive from a truncated read.
   `0028:A1` quotes the existing signature verbatim
   (`Read(ctx, art Artifact, requested []string)`) and states it is unchanged;
   run-1 wrote `(..., error)` because it widened into A1 but did not carry the
   return type back.
3. **Anchor-resolver helper name and file** — decomposition taste.
   `0028:D-selection-predicate` names the one dispatch point that *is* normative
   (`registry.go::commandBacked` becoming three-way); everything below it is
   free.
4. **Template parser/expander helper name** — same. C1.2 `parse once:` fixes the
   *semantics* (parse at LOAD into literal|group|placeholder segments; never
   re-scan emitted bytes; never call `regexp.Expand`, with the reason given).
   The name is not a contract.
5. **Line-buffer Go type (`[]Line` / `[][]byte`)** — inadmissible. The
   byte-level contract is fully pinned by `0028:D-wire-byte-format` and C1.3
   `terminators:`; the in-memory representation is unobservable.
6. **`Apply` function name / `edit_binding.go` path** — naming taste.
7. **Whether the gate pre-check lives in `Apply` or in `Executor.Write`** —
   ALREADY CHARTED, and run-1 says so itself. `0028:A10` is Pending and owns
   exactly this ("no new import, no `Binding` interface method"), with a named
   fallback in its If-wrong; `0028:A2` is Pending and owns the ordering and the
   `protectedKeys` baseline interaction. C1.3 `read-back:` marks the SITE
   conditional on A10 in the record's own text. Not a fresh silence.
8. **Whether the executor's deadline arm can mint a false applied sense over a
   pre-write refusal** — ALREADY CHARTED as `0028:A11` (Pending), and C1.3
   `order:` scopes its guarantee to refusals the binding itself mints, in
   terms, naming A11.
9. **Re-anchor index arithmetic** — run-1 reconstructed
   `reanchorMatchesOwnLine` without marking it, and it is correct: this is
   `0028:A12` (Pending), which states the pre-edit-index-minus-preceding-
   deletions derivation and its If-wrong (compare content instead). Already
   charted.

## Escalation note

The reconstruction agreed with the RDR closely enough to be an escalation cue in
its own right, and this should be read against the LITE variant's limits rather
than as a clean bill. A single alternate-model run gives no disagreement count,
so a shared blind spot between reconstructor and record is indistinguishable
from a determinate contract. Two structural facts argue the agreement is real
rather than an artefact: run-1 documents its widening (into A1,
§illustrative-code, §load-bearing-decisions) and read C1.1–C1.6 and the MVV in
full without truncation, and the three findings above are all in the same narrow
band — *ordering between independently-stated refusals* — which is precisely the
band C1 addresses at load time (C1.4 `precedence:`) and leaves open at apply
time. That is a coherent single gap surfacing three times, not three unrelated
misses. None is a blocker: each is resolvable by one added sentence in C1.3
plus, for finding 1, one S-scenario in S28's shape.
