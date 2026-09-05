# REQ List — RDR 0020 Undeclared `--tag` key admission (what the zero TagDecl means)

Phase 0 audit artifact. Every clause below is a testable obligation drawn from
`docs/rdr/0020-undeclared-tag-key-admission.md`. Quotes are verbatim — copied
from the projector (`rdr inspect --select <id>`) for the fenced element and read
from the record for testable prose outside the fences — never transcribed by
hand.

Element ids (`0020:C1`, `0020:MVV`, `0020:D-*`, `0020:S1`…`0020:S5`,
`0020:G-cross-cutting`) are carried wherever a REQ derives from a labelled
element, so a later stage can trace the REQ back to its contract.

Source sections are abbreviated:

- `NC` = Proposed Solution / Technical Design / Normative Contracts (fenced, `0020:C1`)
- `AP` = Proposed Solution / Approach
- `TD` = Proposed Solution / Technical Design (unfenced prose)
- `LBD` = Technical Design / Load-Bearing Decisions
- `DISP` = Technical Design / Disposition (table)
- `FM` = Trade-offs / Failure Modes
- `MVV` = Implementation Plan / Minimum Viable Validation
- `IP` = Implementation Plan / Phase 1 (Steps 1–3)
- `TS` = Validation / Testing Strategy (scenarios `0020:S1`…`0020:S5`)
- `XC` = Finalization Gate / Cross-Cutting Concerns (`0020:G-cross-cutting`)

## Standing notes for this REQ set

**S1 — `counts.elements.C = 1`, but C1 is a multi-clause fence.** The record has
exactly one `normative` fence (`0020:C1`, 3217 bytes) carrying the carrier rule,
a four-item closed refusal list, a hoist obligation with a guard, and a set of
"unchanged" pins. Each of those is an independent testable obligation and is
mined separately below. Further testable prose sits OUTSIDE the fence:
`Technical Design`, `Load-Bearing Decisions`, the `Disposition` table (a
seven-row input-class → exit mapping), `Failure Modes`, the MVV, Phase 1's three
Steps, and the five Testing Strategy scenarios. All are mined.

**S2 — this is a one-seam change with a widening direction.** C1 widens what
`internal/cli/flow_input.go::parseTags` accepts and narrows nothing except by
relocation (the empty-value arm). The verification surface is therefore
asymmetric: the NEW behaviour (REQ-1..8) needs red tests, while most of the
remaining REQs are **pin-the-unchanged** obligations whose failure mode is
regression, not absence. Each REQ below is marked *(new)*, *(pin)*, or
*(move — behaviour-preserving)* where the distinction is load-bearing.

**S3 — the docs obligation names a file that does not exist.** Phase 1 Step 3
(IP) obligates "the `docs/model-schema.md` line for authors composing a
producer's whole tag output"; `docs/model-schema.md` is absent at HEAD
(`ls docs/*.md` → `cli-output-contract.md`, `cli-reference.md`,
`model-authoring.md`, `README.md`). The prose C1 falsifies lives in
`docs/model-authoring.md` §"Passing a set-valued tag" (lines 925–943), which
states verbatim at HEAD: "**An undeclared tag whose value is an array is
refused.**" See QUESTIONS Q1 and REQ-40.

---

## A. The carrier rule (`0020:C1`, NC)

- [REQ-1] `0020:C1` *(new)* "PURE CARRIER. At `--tag` admission, a key ABSENT from the loaded model's normalized tag table (the two-value `m.Tags[key]` lookup) is a pure carrier: the CLI admits its value VERBATIM — byte-preserved, JSON array literals included — and never canonicalises, conforms, kind-checks, or compares it." — (NC)

- [REQ-2] `0020:C1` *(new)* "The carried value rides the kernel's Observed view and echoes in the resolve payload's `observed` field byte-for-byte as given" — (NC)

- [REQ-3] `0020:D-identity` *(new)* "\"declared\" means the key is present, byte-exact, in the normalized model's tag table (`m.Tags`, two-value lookup); the same presence signal `internal/table/load.go::accessorTable` already branches on (`_, ok := l.model.Tags[key]`). No parallel registry, no case folding, no zero-decl sentinel." — (LBD)

- [REQ-4] `0020:D-wire-byte-format` *(new)* "an undeclared value crosses verbatim, never re-canonicalised: canonicalising would interpret as a set a value whose declaration never asserted set-ness." — (LBD)

- [REQ-5] `0020:C1` *(pin)* "declared set values keep the canonical-array form of `docs/cli-output-contract.md` §Set values on the wire" — (NC; also LBD `0020:D-wire-byte-format`: "Declared set values keep the §Set-values-on-the-wire canonical array unchanged.")

- [REQ-6] `0020:G-cross-cutting` *(new)* "A carried value is admitted VERBATIM — byte-preserved, with no canonicalisation, folding, normalisation, or re-encoding — and echoes in the resolve payload's `observed` field as given." — (XC, Character encoding)

## B. The closed refusal list for an undeclared key (`0020:C1`, NC)

- [REQ-7] `0020:C1` *(pin)* "The only refusals reachable for an undeclared key are, unchanged and still preceding any accessor (REQ-27):" — (NC) — i.e. the list REQ-8..11 is EXHAUSTIVE for an undeclared key; any other refusal firing over an undeclared key is a violation.

- [REQ-8] `0020:C1` *(pin)* "the flag grammar (`--tag` takes `name=value`), `flow-tag-invalid`;" — (NC)

- [REQ-9] `0020:C1` *(pin)* "`flow-tag-reserved` (REQ-26), `flow-tag-owned` (REQ-27), `flow-tag-duplicate` (REQ-28);" — (NC)

- [REQ-10] `0020:C1` *(move — behaviour-preserving)* "the empty-value arm, `flow-tag-invalid`" — (NC)

- [REQ-11] `0020:C1` *(new)* "`flow-tag-invalid`'s kind, shape, and domain arms — including \"is not set-valued\" — are reachable only for a DECLARED key, whose loaded declaration (kind required at load, one of RDR 0003's five tokens) is what the message then truthfully reports." — (NC)

- [REQ-12] `0020:C1` *(pin)* "No new refusal code is minted; `flow-tag-undeclared` does not exist." — (NC; also `0020:D-naming`: "no new code; the rejected alternative's `flow-tag-undeclared` is deliberately not minted")

- [REQ-13] `0020:C1` *(pin)* "The model-side closed world is untouched: a rule atom, accessor key, write target, clear target, or `[initial]` assignment naming an undeclared tag still refuses at load (`unknown_tag`)." — (NC)

- [REQ-14] `0020:C1` *(new)* "Declaring a key later is the opt-in tightening: admission then enforces that declaration's kind and domain." — (NC)

## C. The hoisted empty-value arm and its guard (`0020:C1`, NC)

- [REQ-15] `0020:C1` *(move)* "an empty observed value is indistinguishable from unset, which is a grammar-level fact about the value that holds with or without a declaration. It therefore binds the carrier too, and MOVES to the admission path ahead of the carrier branch; today it sits inside `canonicalValue`'s `!isSet` arm, which the carrier no longer enters." — (NC)

- [REQ-16] `0020:C1` *(move — behaviour-preserving)* "The hoisted arm carries the SCALAR-shaped message for a carrier — ``the tag `X` was given an empty value`` — which is byte-identical to what an undeclared key already receives at HEAD (normative fixture **F4**, Testing Strategy scenario 4, from `evidence/spikes/d-undeclared-empty.txt`), so the hoist is behaviour-preserving." — (NC)

- [REQ-17] `0020:C1` *(pin)* "The set-specific empty-value message (\"is set-valued and takes a JSON array literal; got \") is a CONFORMANCE message and stays declared-only, inside the arms below: it presupposes a declared kind, which a carrier by definition has none of." — (NC)

- [REQ-18] `0020:C1` *(new — the guard)* "The hoisted arm is therefore NOT unconditional on `value == \"\"`: it fires for a key that is undeclared, or declared with a NON-set kind, and it must not intercept a declared SET key, whose empty value keeps the set-specific conformance message (normative fixture **G**, `evidence/spikes/g-declared-set-empty.txt`: `--tag labels=` ⇒ ``the tag `labels` is set-valued and takes a JSON array literal; got ``)." — (NC)

- [REQ-19] `0020:C1` *(pin)* "A declared scalar's empty value keeps the scalar message it has at HEAD (`evidence/spikes/e-declared-empty.txt`), which is the same string the hoisted arm emits, so routing it through either site is byte-identical." — (NC)

- [REQ-20] `0020:C1` *(new — negative)* "Hoisting on the bare value test alone would move a declared set key off its conformance message — a declared-key change this RDR does not make." — (NC) — a test asserting the unconditional form FAILS.

## D. Refusal precedence at admission (`0020:D-selection-predicate`, LBD)

- [REQ-21] `0020:D-selection-predicate` *(pin + move)* "refusal precedence at admission is unchanged in order: grammar → reserved → owned → duplicate → empty-value (undeclared or non-set-declared keys) → (declared keys only) kind/shape/domain conformance;" — (LBD)

- [REQ-22] `0020:D-selection-predicate` *(new)* "the carrier branch sits where the conformance arms would have run." — (LBD)

- [REQ-23] `0020:D-selection-predicate` *(pin — the anti-inversion rule)* "It lands AFTER the duplicate arm and the `seen[key]` mark, not before: at HEAD the duplicate check already precedes the empty-value site (`flow_input.go:663` vs the `:723` arm inside `canonicalValue`), so a repeated key refuses `flow-tag-duplicate` today and must still do so — hoisting the arm ahead of the duplicate case would invert that precedence and change behaviour this RDR does not touch." — (LBD)

- [REQ-24] `0020:D-selection-predicate` *(latitude bound)* "The order above is the whole constraint on the splice point; where the arm sits within it is implementation latitude." — (LBD)

## E. Seam scope — one function moves (TD)

- [REQ-25] *(new)* "One seam moves: the admission path in `internal/cli/flow_input.go`. `parseTags` switches its guarded lookup to the two-value form and routes an undeclared key around `canonicalValue` (or passes the declaredness bit into it — implementation latitude, bounded by C1's refusal list: whichever shape is taken, the empty-value arm must still fire for a carrier and must still leave a declared set on its conformance message), so the kind/shape/domain arms run only under a real declaration." — (TD)

- [REQ-26] *(pin)* "The carried value flows exactly where an undeclared scalar already flows today: into the kernel's Observed view (`resolve.Input.Observed`, where only atoms — all declaration-checked at load — can read keys) and out through the resolve payload's `observed` echo, byte-for-byte as given." — (TD)

- [REQ-27] *(pin — out of scope)* "`--write`/`--clear` are out of scope: `internal/cli/flow_state.go::parseWrites` proves a writer binding (and therefore a declaration) before its `canonicalValue` call, so the zero-decl arm is reachable from `parseTags` alone." — (TD)

- [REQ-28] *(pin)* "`canonicalValue` keeps its copy for the `--write` carrier, which enters by its own path." — (IP, Phase 1 Step 1)

- [REQ-29] *(pin)* "The provenance guards are untouched and still precede admission: reserved key, owned key, duplicate (REQ-26/27/28)." — (AP)

## F. Disposition table — input class → exit (DISP)

Seven rows; each is a testable input-class obligation. Quoted cell-wise.

- [REQ-30] *(pin)* "undeclared key, scalar value | 0 | — | `observed` echo, verbatim | loud (echoed)" — (DISP)

- [REQ-31] *(new — the change)* "undeclared key, array literal | 0 (**changed**; refuses at HEAD) | — | `observed` echo, verbatim | loud (echoed)" — (DISP)

- [REQ-32] *(move)* "undeclared key, empty value | 2 | `flow-tag-invalid` \"was given an empty value\" | none | loud" — (DISP)

- [REQ-33] *(new — accepted cost)* "undeclared key that is a MISSPELLING of a declared one | 0 | — | `observed` echo, verbatim | **silent** — accepted open-world cost; diagnosis is the echo (Failure Modes)" — (DISP)

- [REQ-34] *(pin)* "declared scalar/enum key, array literal | 2 | `flow-tag-invalid` \"is not set-valued\" — now truthful | none | loud" — (DISP)

- [REQ-35] *(pin — the guard's witness)* "declared set key, empty value | 2 | `flow-tag-invalid` \"is set-valued and takes a JSON array literal; got \" | none | loud" — (DISP)

- [REQ-36] *(pin)* "reserved / owned / duplicate key, any value | 2 | `flow-tag-reserved` / `-owned` / `-duplicate` | none | loud (unchanged, still first)" — (DISP)

- [REQ-37] *(pin)* "Visible: every refusal that survives is unchanged and still fires at exit 2 before any accessor — reserved, owned, duplicate and grammar for ANY key, and empty-value for any key that is not a declared set (a declared set's empty value refuses on its conformance message instead); wrong kind, wrong shape and out-of-domain for a DECLARED key only (per C1's list)." — (FM)

- [REQ-38] *(new)* "Diagnosis: the resolve payload's `observed` field echoes every carried key byte-for-byte — the stray spelling sits beside the declared keys in the same envelope the refusal rides." — (FM)

## G. Comments and docs (IP Phase 1 Steps 2–3)

- [REQ-39] *(new)* "Rewrite the guarded-lookup comment to cite this RDR's carrier contract instead of promising a different decision elsewhere; `ConformValue`'s zero-decl doc stays true as written." — (IP, Step 2) — the comment at `internal/cli/flow_input.go:669-675` ("Refusing an undeclared key is a different decision, under a different code, and this is not the place that makes it") is the target; `internal/table/load.go::ConformValue`'s doc is NOT edited.

- [REQ-40] *(new)* "the `docs/model-schema.md` line for authors composing a producer's whole tag output (Background's obligation)" — (IP, Step 3). See ASSUMPTION-1 and QUESTIONS Q1: the named file does not exist; the falsified prose is `docs/model-authoring.md` §"Passing a set-valued tag" (HEAD lines 925–943), whose lead sentence is "**An undeclared tag whose value is an array is refused.**"

- [REQ-41] *(new)* "under arm (1) or (3), a `docs/model-schema.md` line for authors composing a producer's whole tag output." — (Context / Background)

## H. Minimum Viable Validation (`0020:MVV`)

- [REQ-MVV] `0020:MVV` — the five steps below are the MVV; done is all five demonstrably satisfied.

- [REQ-42] `0020:MVV` "Author a fixture model declaring one scalar tag and one set tag, with no declaration for `extra` or `extras`." — (MVV 1)

- [REQ-43] `0020:MVV` "Run `flow resolve` with `--tag extra=plain` and `--tag extras=[\"a\",\"b\"]` beside the declared tags: the invocation exits 0 and the payload's `observed` field carries both values byte-for-byte as given." — (MVV 2)

- [REQ-44] `0020:MVV` "Run the same invocation without the two carrier flags: the selected rule and outcome are identical — the carried keys influenced nothing (A2's runtime leg)." — (MVV 3)

- [REQ-45] `0020:MVV` "Hand the DECLARED scalar tag an array literal: still refused `flow-tag-invalid` \"is not set-valued\" — now truthfully, and the red test pins this beside step 2 so the asymmetry is authored, not incidental." — (MVV 4)

- [REQ-46] `0020:MVV` "Run `--tag extra=` (undeclared, empty): still refused `flow-tag-invalid` \"was given an empty value\" — the one arm the carrier does not escape, pinned so the hoist out of `canonicalValue` cannot silently drop it." — (MVV 5)

- [REQ-47] *(new)* "Land the red test of MVV steps 2–5 (undeclared scalar AND array pass, declared-scalar-given-array still refuses, undeclared-empty still refuses)" — (IP, Step 3)

## I. Testing Strategy scenarios and normative fixtures (`0020:S1`…`0020:S5`)

- [REQ-48] *(structural)* "Coverage goal — done is all six green, with the undeclared/declared pair adjacent in one test so the asymmetry C1 fixes is pinned by construction rather than by two tests that could drift" — (TS). See ASSUMPTION-2 on "all six".

- [REQ-49] `0020:S1` "**Scenario**: undeclared scalar and undeclared array literal admitted (MVV 2). **Expected**: exit 0; `observed` echoes both byte-for-byte." — (TS)

- [REQ-50] `0020:S1` — normative fixture **F1**: "for a model declaring scalar `tier` and set `labels` and nothing named `extra`: `--tag tier=free --tag labels=[\"security\"] --tag extra=plain` ⇒ `\"observed\":{\"extra\":\"plain\",\"labels\":\"[\\\"security\\\"]\",\"tier\":\"free\"}`." — (TS)

- [REQ-51] `0020:S1` — the array leg's asserted wire form: "the asserted wire form follows from `resolve.Input.Observed`'s `map[string]string` type — the raw argument text as a string value, `\"extras\":\"[\\\"a\\\",\\\"b\\\"]\"`, the same shape F1 already witnesses for the declared set key `labels`." — (TS)

- [REQ-52] `0020:S2` "**Scenario**: the same resolve with and without the carrier flags (MVV 3). **Expected**: identical selected rule and outcome." Normative fixture **F2**: "`\"rule\":\"free\"`, `\"emit\":{\"plan\":\"basic\"}`, with `gates`, `next`, `writes`, `clear` empty and `escaped:false` in both runs — every projected field diffs empty and the sole payload delta is `observed` gaining the carried key." — (TS)

- [REQ-53] `0020:S3` "**Scenario**: declared scalar given an array literal (MVV 4). **Expected**: refused `flow-tag-invalid`, \"is not set-valued\"." Normative fixture **F3**: "`--tag tier=[\"a\",\"b\"]` ⇒ exit 2, ``{\"code\":\"flow-tag-invalid\",\"message\":\"the tag `tier` is not set-valued; got the array literal [\\\"a\\\",\\\"b\\\"]\",\"param\":\"tier\"}``." — (TS)

- [REQ-54] `0020:S4` "**Scenario**: undeclared key given an empty value (MVV 5). **Expected**: refused `flow-tag-invalid`, \"was given an empty value\" — the arm the carrier does not escape." Normative fixture **F4**: "`--tag extra=` ⇒ exit 2, ``{\"code\":\"flow-tag-invalid\",\"message\":\"the tag `extra` was given an empty value\",\"param\":\"extra\"}`` — byte-identical to HEAD". — (TS)

- [REQ-55] `0020:S5` "**Scenario**: declared SET key given an empty value — the hoist's one regression risk, pinned beside scenario 4 so the guard on the hoisted arm cannot be dropped silently. **Expected**: refused `flow-tag-invalid` on the SET-specific conformance message, not the scalar-shaped one." Normative fixture **G**: "`--tag labels=` ⇒ exit 2, ``{\"code\":\"flow-tag-invalid\",\"message\":\"the tag `labels` is set-valued and takes a JSON array literal; got \",\"param\":\"labels\"}`` — byte-identical to HEAD." — (TS)

- [REQ-56] `0020:S5` *(negative test obligation)* "An unconditional `value == \"\"` test ahead of the declared lookup fails this scenario, which is what makes it the guard's witness." — (TS)

## J. Cross-cutting invariants (`0020:G-cross-cutting`, XC)

- [REQ-57] `0020:G-cross-cutting` *(pin)* "C1 widens what admission accepts and narrows nothing, so no shipped invocation changes meaning: undeclared scalars passed before and pass now, and the only behaviour change is that an undeclared array is admitted where it was refused." — (XC) — i.e. NO migration surface; any pre-existing green test over `--tag` admission that turns red is a regression, not an expected update, EXCEPT a test asserting the undeclared-array refusal (A1 verified there is none).

- [REQ-58] `0020:G-cross-cutting` *(pin)* "The one arm that does not widen, the empty-value refusal, is hoisted rather than changed and is byte-identical for both affected input classes (fixtures F4, E, G), so it needs no migration either." — (XC)

- [REQ-59] `0020:G-cross-cutting` *(scoping)* "This RDR makes no byte-identical output, content-addressed identity, or replay-stable hash claim, so the determinism checklist does not apply. The two \"byte-identical\" uses in C1 are narrower: they assert message-string equality between two refusal sites across the hoist" — (XC) — "byte-identical" in this record means MESSAGE-STRING equality, not payload-hash stability; do not build a hash/golden-digest test on it.

---

## ASSUMPTIONS

- **ASSUMPTION-1 — the docs REQ (REQ-40/41) lands in `docs/model-authoring.md`,
  and includes correcting the now-false paragraph there.** IP Step 3 and
  Background both name `docs/model-schema.md`, which does not exist at HEAD.
  The record's own Technical Environment section names the surfaces it read and
  does not include `model-schema.md`; the doc that actually carries the tag
  declaration/authoring material is `docs/model-authoring.md`, whose
  §"Passing a set-valued tag" states at HEAD "**An undeclared tag whose value is
  an array is refused.**" and "**every set-valued key the caller passes must be
  declared, even one no row guards on.**" — both directly falsified by C1
  (REQ-1, REQ-31). Reading `model-schema.md` as a stale/aspirational filename
  for the model-authoring doc is the single defensible reading: the alternative
  (create a new `docs/model-schema.md` and leave `model-authoring.md` asserting
  the retired behaviour) would ship documentation that contradicts the shipped
  contract, which the record's own Problem Statement calls out as the defect it
  is fixing ("a message pointing at a declaration that does not exist"). See Q1.

- **ASSUMPTION-2 — TS's "all six green" counts the five numbered scenarios plus
  the scalar/array pair inside scenario 1 counted separately.** The Testing
  Strategy heading says "done is all six green" but enumerates five numbered
  scenarios (`0020:S1`…`0020:S5`). Scenario 1 explicitly carries two legs ("the
  scalar leg" witnessed at HEAD by F1, and "the array leg … is the post-change
  extension of F1 and is what the red test adds"). The Desk trace table
  independently lists six behavioural rows (MVV 1–5 plus row "5b — declared SET
  key, empty value (not an MVV step; the hoist's blast radius)"). Both readings
  land on the same six ASSERTIONS, so the count is not load-bearing for
  coverage: REQ-49..56 enumerate every one. Recorded because a coverage matrix
  keyed on "six scenarios" must not drop the S1 array leg or S5.

- **ASSUMPTION-3 — "the two-value `m.Tags[key]` lookup" (REQ-1, REQ-3) is a
  presence test on the NORMALIZED model, matching the sibling call sites.**
  `0020:D-identity` names `internal/table/load.go::accessorTable`'s
  `_, ok := l.model.Tags[key]` as "the same presence signal" and the Authority
  census row names four sibling arms all using `_, ok := Tags[key]`. Reading
  this as a plain Go two-value map index (not a helper, not a new predicate
  method) is the only reading consistent with "No parallel registry … no
  zero-decl sentinel". Implementation latitude covers whether the bit is
  branched on in `parseTags` or threaded into `canonicalValue` (REQ-25).

- **ASSUMPTION-4 — REQ-7's closed refusal list binds the ADMISSION path only,
  not the whole invocation.** C1 says "The only refusals reachable for an
  undeclared key are … still preceding any accessor (REQ-27)". A carrier key
  obviously does not exempt the invocation from every downstream refusal
  (missing `--outcome`, no-match resolution, accessor failures). Reading the
  list as exhaustive over `parseTags`'s own refusals for that key is the only
  coherent reading; the "still preceding any accessor" clause is what fixes the
  scope boundary.

- **ASSUMPTION-5 — "declared with a NON-set kind" (REQ-18) means
  `decl.Kind != "set"` over a key PRESENT in `m.Tags`.** C1's guard has three
  input classes (undeclared; declared non-set; declared set) and the record's
  Illustrative Code renders it as `if value == "" && (!declared || !decl.IsSet())`.
  The Illustrative Code is explicitly "shape only, not load-bearing", but the
  set-vs-non-set discriminator it uses matches `canonicalValue`'s existing
  `isSet := decl.Kind == "set"` at HEAD (`flow_input.go:715`), so no new kind
  taxonomy is introduced. `0020:D-identity`'s "no case folding" applies to the
  KEY, not the kind token.

- **ASSUMPTION-6 — REQ-2/REQ-26's "byte-for-byte as given" means the raw
  post-`strings.Cut` value.** `--tag k=v` is split on the first `=`
  (`flow_input.go:650`), so a value containing `=` carries its remainder
  verbatim. "Verbatim" is scoped to what admission received after the grammar
  split, not to the whole argv token; REQ-8 keeps the grammar arm unchanged, and
  no clause in the record proposes changing the split.

---

## QUESTIONS

- **Q1 — Does REQ-40's docs obligation create `docs/model-schema.md`, or amend
  `docs/model-authoring.md`?** Two readings, materially different in output:
  (a) Create a new `docs/model-schema.md` carrying the composition guidance, and
  leave `docs/model-authoring.md` §"Passing a set-valued tag" as written.
  (b) Treat `model-schema.md` as a stale name for the authoring doc: put the
  line in `docs/model-authoring.md` AND correct that section, whose HEAD text
  ("**An undeclared tag whose value is an array is refused.**"; "**every
  set-valued key the caller passes must be declared, even one no row guards
  on.**") C1 falsifies outright.
  **Proceeding under (b)** (ASSUMPTION-1). Grounding: the resources index
  (`.rdr/resources.md` §Design docs) names only `docs/cli-output-contract.md`
  and `CLAUDE.md` as authoritative contracts — no `model-schema.md` — and
  `docs/README.md`-listed authoring material lives in `model-authoring.md`.
  Predecessor precedent: 0005's and 0008's REQ sets both treat docs clauses as
  obligations against the doc that actually carries the surface, not against a
  filename literal. Reading (a) would leave a shipped doc asserting the exact
  behaviour the RDR retires — the failure mode the record's Problem Statement
  names ("undiscoverable from the error"). Recorded because (a) is a strictly
  additive, cheaper output and a later stage may prefer it; if so, the
  `model-authoring.md` correction is still owed under REQ-57's "no shipped
  invocation changes meaning" and REQ-31's changed disposition, and should be
  raised as a deviation rather than dropped.

- **Q2 — Does REQ-33 (misspelled key passes silently) carry any implementation
  obligation, or is it purely a recorded consequence?** Two readings:
  (a) It is a disposition row like the others and needs a pinned test asserting
  a misspelling of a declared set key exits 0 and echoes.
  (b) It is a documented accepted cost with no separate obligation — its
  behaviour is already fully determined by REQ-1/REQ-31 (a misspelling IS an
  undeclared key), and the MVV does not enumerate it.
  **Proceeding under (b) for the MVV bar, (a) as optional coverage.** Grounding:
  the MVV's five steps and TS's five scenarios never mention a misspelling; the
  Premortem calls it "today's shipped behavior for scalars" and answers it with
  the `observed` echo (REQ-38), not with a new check; and C1's normative text
  contains no misspelling clause. The row's "Silent or loud" column marks it
  **silent** — the only such row — which is a property of the disposition, not
  an assertion to author. Recorded because a coverage matrix built strictly
  row-by-row off the Disposition table would demand a test the MVV does not,
  and because reading (a) is harmless (the test would pass under REQ-31's
  implementation) but must not be counted as an MVV gate.
