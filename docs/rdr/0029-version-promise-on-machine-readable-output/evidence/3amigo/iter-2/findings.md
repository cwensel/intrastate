Model: claude-opus-5

# Delta re-check, iter-2 — cli/0029

Scope: only the ids the rewrite touched. Three personas worn over each.
Four findings, all class 1 or class 3. Verified-clean list at the bottom.

---

## F1 — `0029:C1` — class 3 (FALSE against the source tree)

**New text.** C1's scoping paragraph now reads:

> These two rules bind an EXTERNAL consumer; this CLI emits no check for
> them, and the repo's own `--plan` reader is not one —
> `internal/cli/flow_input.go::planEnvelope` decodes an authored INPUT
> document, so `0002:C3`'s strict-input rule governs it, not this clause.

**Defect.** `planEnvelope` does not decode an authored input document. It
decodes THIS CLI's own emitted terminal envelope, piped. Its own doc
comment at `internal/cli/flow_input.go:310-317` says so:

> A plan arrives in one of two shapes and BOTH are accepted: the full
> `{"type":"ok","data":…}` envelope `respond.OK` emits, or the bare `data`
> object a caller extracted from it (`jq .data`, say).

`readPlan`'s refusal messages say the same thing twice — "`--plan` takes
one `flow resolve --as json` envelope or its `data` object"
(`flow_input.go:390`, `:398`) — and `flow_input.go:331-337` declares
`Code`/`Message` on the struct precisely to recognise the BARE `CLIError`
refusal record, i.e. the second of the two terminal records C1 governs.

So `planEnvelope` IS a consumer of the `--as=json` terminal envelope, and
it is in-repo. The new clause disclaims the one in-repo witness that
actually exists, on a factual premise the source contradicts. It also
misattributes the governing contract: `0002:C3` is the model-TOML
strict-input rule; a piped `flow resolve` envelope is not a model TOML.

Two live consequences, not merely a wording slip:

- `planEnvelope` is a struct decode with four named fields and no
  `DisallowUnknownFields`, so it is ALREADY a tolerant reader and will
  already survive the added `schema_version` key. That is a real,
  assertable satisfaction of C1's first consumer rule — and the rewrite
  now forbids anyone from citing it.
- C1 ends "A surface MUST NOT be given the input rule and the output rule
  at once." The new paragraph assigns `planEnvelope` the INPUT rule
  (`0002:C3`, strict) while the document it reads is an OUTPUT envelope.
  If that assignment were taken literally, a future maintainer would add
  `DisallowUnknownFields` to `planEnvelope` and break the `--plan` pipe on
  the very release that adds `schema_version` — the exact defect C1's
  tolerant-reader rule exists to prevent, produced by C1's own new text.

**Blocks.** The implementer cannot decide whether `planEnvelope` needs a
tolerance assertion, and `0029:S3` ("the existing in-repo envelope consumer
survives the added field") now has no stated subject — C1 has disclaimed
the only consumer S3 could mean.

**Also.** The same false premise was copied into `0029:MVV`'s `trace`
table, last row ("the repo's `planEnvelope` reads authored INPUT under
`0002:C3`'s strict rule, so it is not a witness for this row"). Same
defect, same fix.

---

## F2 — `0029:C4` — class 1 (internal contradiction: seam table vs. C4's own tier lists)

**New text.** C4's ENUMERATION SEAM clause: "A tier assignment obliges an
ENUMERATION SEAM… Five do, and adding their accessors is part of this
RDR's implementation", followed by a 5-row table.

**Defect.** C4's own tier lists name SIXTEEN tiered vocabularies. The seam
table accounts for five owed plus four claimed already-carried
(`RefusalKinds`, `Categories`, `Verdicts`, `Reasons`) plus one explicit
carve-out (`graph-lint-failed`) — ten. Six tiered entries are covered by
neither the "owed" table nor the "already carries the idiom" list, and two
of the six provably have NO accessor in the source tree:

| C4-tiered vocabulary | Tier | Seam in tree | Where C4 places it |
| --- | --- | --- | --- |
| envelope `type` (`ok`/`failed`) | `frozen` | **none** — bare literals, `respond.go:136` sets `s.Type = "ok"`; no exported accessor in `respond` (its exported funcs are `ModeOf`, `ValidateMode`, `OK`, `Fail`, `Note`, `Warn`) | neither table nor idiom list |
| `findings[].block` (`resolve::Block`) | `append-only` | **none** — `internal/resolve` exports only `RefusalKinds`, `Reasons`, `Resolve`, `CheckInput`, `TestGuardEvaluatorContract`; there is no `Blocks()` | neither table nor idiom list |
| severity (`blocking`/`info`) | `frozen` | `graphlint::Severities()` exists | unlisted (harmless, but uncounted) |
| `graph-unprovable-coverage` `reason` | `append-only` | `graphlint::Reasons()` | idiom list — but the idiom list credits it to the FOUR-package claim, double-counting with `resolve::Reasons()`, a DIFFERENT accessor in a different package |
| graph-lint ADVISORY codes | `growing` | `graphlint::AdvisoryCodes()` exists | unlisted |
| `version` payload field names | `frozen` | struct tags on `internal/version::Info`, not a set accessor | unlisted |

The two hard gaps are the first two rows. C4 asserts the seam obligation in
absolute terms — "Without one the tier is unassertable… an unassigned
machine-readable surface is a defect" — then enumerates a closed list of
five that does not include the envelope `type` discriminator, which is the
FIRST entry in its own `frozen` list, nor `findings[].block`, which is the
LAST entry in its own `append-only` list.

The count is also contradicted from a second direction.
`0029:§capability-dependencies` (also rewritten) now says "Nine tiered
vocabularies enumerate themselves today; five do not". Nine + five = 14,
but C4 tiers 16. The two rewritten passages disagree with each other and
both disagree with C4's lists.

**Blocks.** `0029:S6` and `0029:S7` were rewritten so that "the assertion
subject is C4's enumeration seam". For the envelope `type` set (S6,
`frozen`) and `findings[].block` (S7, `append-only`) there is no seam and
the table does not oblige one, so neither scenario has a subject and
neither is writable as a pass/fail test. `0029:§step-3` says "Build the
enumeration seams C4 obliges for the five vocabularies lacking one" — an
implementer following it literally ships without the two.

**Also.** `0029:A5` is scoped to "the five that lack one" and reasons about
only two hard cases (CLIError `code`, `flow next` reason). If the true
count is seven, A5's spike as written does not cover the addition, and A5's
"If wrong" fallback names only those two rows.

---

## F3 — `0029:S7` / `0029:S8` / `0029:§step-3` — class 1 (internal contradiction with C2's cardinality prohibition)

**New text.** S8 now instructs that
`internal/graphlint/findings_0006_test.go::TestReq74_TheAdvisoryTierIsClosedAtExactlyFourMembers`
be REPLACED, not updated, on the reasoning that "its name and its
four-element `want` literal both encode a cardinality that C2 forbids
asserting on a `growing` set, so re-pinning it at five reproduces the
defect at six."

**Defect.** The reasoning is right and the scope is wrong. C2's prohibition
binds `append-only` as well as `growing` — "A consumer… MUST NOT assert on
the set's cardinality, a member's ordinal position, or a tail position" is
stated in the `append-only` bullet, and `growing` inherits it by "as
`append-only`, and additionally…". At least two existing tests commit the
identical defect over vocabularies C4 tiers `append-only`, and no delta id
names either:

- `internal/graphlint/findings_0006_test.go:257`
  `TestReq73_TheBlockingCodeSetIsExactlyTheTenNamed` — an exhaustive
  `slices.Equal` over `graphlint.BlockingCodes()`, with the cardinality in
  the test NAME ("TheTen"), over the graph-lint BLOCKING codes, which C4
  tiers `append-only`. This is S8's defect verbatim, one function above
  S8's target, in the same file.
- `internal/graphlint/findings_0006_test.go:517`
  `TestReq80_UnprovableCoverageCarriesAReasonFromTheClosedSet` — exhaustive
  `slices.Equal` over `graphlint.Reasons()`, the `graph-unprovable-coverage`
  `reason` set, which C4 tiers `append-only`. Its in-test comment makes the
  contradiction explicit and self-aware: "The set is declared 'closed,
  **append-only**'… This stays an EXACT-set assertion — only the expected
  set grew — so a fifth member arriving unannounced still fails here."
  That is precisely the growth-blocking assertion S8 rejects for
  `TestReq74`, already adopted as a deliberate policy in the tree.

S7 tells the implementer to extend
`internal/table/command_carrier_0025_test.go::TestReq79_TheCategorySetIsAppendOnlyAndDuplicateFree`
as the pattern for `append-only`, but never says the two exhaustive tests
above are the ANTI-pattern that must go. So the RDR lands with
`append-only` asserted two contradictory ways in one package, and S7's
claim — "`append-only` vocabularies are never asserted by cardinality or
ordinal" — is FALSE in the tree the moment the RDR is declared done.

**Blocks.** S7 cannot be run as a pass/fail check: it asserts a repo-wide
negative that two live tests violate, and the RDR does not authorise
touching them. QA cannot tell whether a red `TestReq73` is a regression or
the intended consequence. Note S8 got exactly this right for its own
target, calling the red test "the intended visible failure, not a flake to
suppress" — the same sentence is owed to `TestReq73` and `TestReq80` and is
absent.

---

## F4 — `0029:S9` / `0029:§step-2` — class 1 (C2's exemption list does not cover what the widened grep hits)

**New text.** S9 widened the grep from two files to six packages —
`internal/graphlint`, `internal/table`, `internal/guard`,
`internal/accessor`, `internal/resolve`, `internal/cli` — "with an explicit
allowlist: the emitted identifier `graph-coverage-closed-by-escape` and the
lines quoting `0003:C7`, both exempted by C2." `§step-2` mirrors the six
packages. C2 scopes its prohibition to descriptive uses that reach "every
site that describes a tiered vocabulary".

**Defect.** The widened scope now sweeps a large body of `closed` text
those two exemptions do not reach and C2's descriptive rule does not
actually forbid, because the vocabularies described are NOT tiered by C4.
Non-test, descriptive, in-scope-by-grep, out-of-scope-by-C2:

- `internal/accessor/model.go:28,46` — "closed at three: read, gate, and
  write"; `Capabilities()`. Not in C4.
- `internal/accessor/model.go:89` — "the closed accessor refusal-class
  set". Not in C4 (C4 tiers `data.escape_class` = `resolve::RefusalKinds`,
  a DISJOINT set, as `model.go:52` itself notes).
- `internal/accessor/model.go:131` — "closed eight-member validation-code
  set". Not in C4.
- `internal/guard/identity.go:35` — "the closed predicate-semantic-kind
  set". Not in C4.
- `internal/guard/doc.go:1`, `internal/guard/grammar.go:47,56,59,60,67,145`
  — "closed typed vocabulary" over `guard::Operators`. This one IS tiered
  (C4 `frozen`), so it IS in scope — but it is not a `0003:C7` quotation,
  it is `guard`'s own prose, so the quotation exemption does not save it
  and it must be rewritten. Same for `internal/table/model.go:91,95,122`
  and `internal/table/normalize.go:70`.
- `internal/table/dump.go:9,22` — "closed column vocabulary"
  (`DumpColumns`). Not in C4.
- `internal/table/source.go:9`, `load.go:55,280,1365,1377,1434`,
  `edit.go:62,404`, `category.go:52` — closed layouts, emit-kinds, the
  `clear` set, the verdict set at load. Not in C4 (C4 tiers
  `table::Categories()` only).
- `internal/cli/cmdbind/cmdbind.go:49,726,920` — "closed placeholder
  vocabulary". Not in C4. (Note `cmdbind` is under `internal/cli`, so S9's
  package-level scope pulls it in.)

Plus a large mass of `closed` used in senses no exemption names and no
reader would call a tier description: `fail closed` / `fails closed`
(`internal/resolve/guard_adversarial_0007_test.go` ~14 sites,
`internal/cli/flowbind/registry.go:22`), `closed-world` matching, pipe/file
"write end is closed" (`cmdbind.go:331,464,556,866`), brace "never closed"
(`edit.go:100,122,255`), and `closedBy`/`ClosedByEscape` identifiers in
`guard/lint.go` and `graphlint/coverage.go` (~20 sites) — which are the
INTERNAL spelling of the emitted identifier C2 exempts, but C2 exempts the
emitted string `graph-coverage-closed-by-escape`, not the Go identifiers
`ClosedByEscape` / `closedBy` that produce it.

`rg -c closed` over the six packages returns matches in 91 files. C2's
two-item exemption list plus S9's restatement of it is not a specification
an implementer can run, and it is not a grep QA can write: as literally
stated ("no occurrence of `closed`… across every package that owns one",
allowlisting exactly two things), the assertion fails on dozens of sites
the RDR does not intend to change.

**Blocks.** `0029:§step-2` is unexecutable as scoped — the implementer
cannot tell which of ~91 files to edit — and S9 is not writable as the
"grep assertion [that] is sufficient" it claims to be. The narrow version
(sites describing a C4-TIERED vocabulary) is implementable; the widened
text does not say that, because C2's "every site that describes a tiered
vocabulary" is contradicted by S9's "no occurrence of `closed`… across
every package that owns one" — a package-scoped, not vocabulary-scoped,
predicate.

---

## Checked and clean

Verified against the source tree, no defect found:

- **`0029:C4` two-`Operators()` paragraph** — TRUE. `internal/guard/grammar.go:50`
  and `internal/table/model.go:93` both declare
  `{"eq","in","lt","lte","gt","gte","exists","contains"}`, same eight, same
  order; `table/model.go:91` carries the cited "does not mint operators and
  does not widen the set (`0002:C16`)" comment verbatim.
- **`0029:C4` aggregate-code carve-out** — TRUE. `graphlint::AggregateCode`
  is a single `const` (`taxonomy.go:48`) with one emit site
  (`internal/cli/lint.go:216`); the "closure of the EMIT SITES" claim is
  assertable as written.
- **`0029:C4` exit-code row** — TRUE. `clierr::ExitCodeFor`
  (`clierr.go:131-149`) is a `switch` yielding exactly `{0,1,2,3,130}`, no
  accessor. Row's "Owed" is correct.
- **`0029:C4` stderr `level` row** — TRUE. `respond.go:207,217` carry bare
  `"note"`/`"warning"` literals; no accessor.
- **`0029:C4` `flow next` row** — TRUE. `internal/cli/flow_next.go:91-95` is
  a `const` block (`reasonAbsent`, `reasonUncomparable`, `reasonNotEvaluated`)
  with no accessor over the union.
- **`0029:C4` four-package idiom claim** — TRUE. `resolve::RefusalKinds`
  (`resolve.go:67`), `table::Categories` (`category.go:90`),
  `accessor::Verdicts` (`model.go:309`), `graphlint::Reasons`
  (`taxonomy.go:112`) all exist and all return `slices.Clone`.
- **`0029:C1` "one home in `clierr`"** — survives the layering C1 describes.
  `respond` already imports `internal/cli/clierr` (`respond.go:32`), and
  `clierr` imports neither, so an exported `clierr.SchemaVersion` read by
  `respond` forms no cycle. `0029:A6`'s one open item ("Unconfirmed: that
  `respond` may import `clierr`") is in fact already settled by the source;
  leaving A6 Pending is conservative, not wrong. No `SchemaVersion` symbol
  exists anywhere today, as A6 states.
- **`0029:A5`'s premise** — consistent with the tree as far as checked;
  `clierr` has no import of `graphlint`/`guard`/`accessor`, so the
  layering-inversion risk A5 flags is real and correctly left Pending.
- **`0029:S1`/`0029:S2` additive criterion** — sound. `Success` does carry
  `omitempty` `notes`/`warnings` and `CLIError` does carry `omitempty`
  `hint`/`findings`, so the exact-key-set assertion the criterion replaces
  would indeed have been flaky. The named baselines (`data,type` and
  `code,detail,message,param`) are consistent with `0005:C1`'s bare-record
  shape, which `flow_input.go:331-337` independently corroborates.
- **`0029:S8`'s target** — exists at
  `internal/graphlint/findings_0006_test.go:287` with exactly the
  four-element `want` and the `slices.Equal` S8 describes; `AdvisoryCodes()`
  exists at `taxonomy.go:109`, so the replacement S8 specifies is writable.
  (The defect in F3 is scope, not this target.)
- **`0029:§prerequisites` 0023 golden item** — TRUE.
  `docs/rdr/0023-resolve-envelope-projection/artifacts/mvv-step1-default-golden.json`
  exists and `internal/cli/flow_mvv_0023_test.go` is present, so the
  re-capture obligation lands on a real artifact.
- **`0029:§prerequisites` 0006:C17 item** — `0006:C17` does read "The
  advisory tier is closed at" its four members, so the amendment
  prerequisite is correctly stated.
- **`0029:D-identity` 0.x qualifier** — internally consistent with C1's 0.x
  clause and with the MVV `trace` table's "version moves" row; no
  contradiction found.
- **`0029:§problem-statement` 0.x reconciliation paragraph** — consistent
  with C2's declaration-of-intent clause and C3's "binds during `0.x`"
  carve-out. No contradiction found.
- **`0029:§activation-step-1`** — the four acceptance criteria are each
  checkable, and the `llms.txt` tension note is disclosure, not a new claim.
  Criterion 3 (branch on presence of `code`, not on `type`) is corroborated
  by `flow_input.go:331-337`, which already implements exactly that
  discriminator.
- **`0029:MVV` `oracle` step-1 row and version-movement row** — the
  amendments are consistent with S1's additive criterion and with A5's
  "no mechanical link" caveat. Only the `trace` table's last row is
  defective (F1).
- **`0029:§capability-dependencies` release-notes note** — an accurate
  statement of C3's unbacked-by-test status; no defect.
