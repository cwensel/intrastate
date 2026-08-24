Model: claude-fable-5

Model stamp: single-model pass (claude-fable-5); iteration-1 critique was
claude-opus-5 — this is the second model of the dual-model draw.

# Hostile Set Critique — RDR 0002–0009, iteration 2

Assignment: fresh, hostile premortem of the Final-and-unimplemented cluster
(0002, 0003, 0004, 0005, 0006, 0007, 0008, 0009) plus the home JDR 0001, read
in their 2026-08-23 text. 0001 is Implemented and is context only. Postmortem
files ignored. No hedging.

**Verdict up front.** Iteration 1 said the set would fail because obligations
were routed into a peer's *implement stage*, a destination that is not a
document. The cluster fixed that by routing them into a peer's *refine* instead
— and then locked every peer, on two consecutive days, without running the
refines. The result is a set of eight `Final` documents in which the newest
normative text (0002 and 0006, both Gate PASS 2026-08-23) contradicts itself on
the first thing a real table author will type, and in which every remaining
`Pending` record names as its venue a document that can no longer receive it.
The only Go that exists is a `version` subcommand and a kernel whose `Row.Guard`
is still a `string`; the corpus of design text has grown from ~29k to ~65k
lines while `internal/` has had no commit since 2026-08-09.

---

## Findings Ledger

| ID | RDR | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-----|-------------|--------------|-------------------|--------|
| C-1 | 0006 / 0002 | 0006 A6 "RDR 0002 is `Draft` with every assumption Verified, so this reopens no locked document"; 0006 Normative "A model that declares no initial owned state MUST be rejected with a blocking finding"; 0002 Normative "The source schema MUST use the Resolve spike field layout: root `outcomes`, `[model]`, `[tags.<tag>]`, `[accessors.<id>]`, `[context.<id>]`, `[[rule]]`, and `[dump]`" + "The decoder MUST reject unmapped keys" | Both locked 2026-08-23; 0002's closed, strictly-decoded layout has no `initial`/`terminal` and 0006 requires one. Every model 0002 can load is rejected by 0006; every model 0006 accepts is rejected by 0002 | `intrastate lint` fails every table with `graph-dangling-edge`; adding `[initial]` fails load with `unknown schema field` | §1, premortem, AT-1 |
| C-2 | 0006 / 0002 | 0006 A10 "Flips when RDR 0002 states write-replaces for single-valued tags. RDR 0002 is `Draft`, so this is a scheduled edit on an open peer" | 0002 locked with no replace-vs-accumulate rule; invariant 5's per-row reading is unsound by 0006's own "If wrong" | `graph-single-valued-state` never fires on a two-write path; multi-valued state reaches `assemble`, lint green | §1, AT-1 |
| C-3 | 0003 / 0006 | 0003 Status "A10, A12, A17 and A19 are on RDR 0006's refine"; 0003 A10 "RDR 0006 is now `Draft`"; 0006 A2 "they close on the peer, not here — a Draft cannot flip a locked peer's record … a route-back for Stage 6 to route" | Six `Pending` records in Final 0003 name Final 0006's refine as venue; 0006 in turn books a route-back on 0003's "next touch". RDRs are never amended; neither touch is scheduled. The records are Pending forever | Implementer of 0003 Phase 2 finds group construction gated on A10 "wait rather than build twice" with nothing to wait for | §1, AT-1 |
| C-4 | 0007 / 0009 | 0007 Prereq "[ ] **RDR 0009 re-locks against the loss of `Refusal.Guard`.** 0009 is `Final` and names `Refusal.Kind` / `Refusal.Guard` as the asserted properties"; JDR §JD-12 "corrected at its re-lock"; 0009 A3 "the asserted properties are `Refusal.Kind` / `Refusal.Guard`" | 0009 Gate PASS 2026-08-11, before §D1/§D4 existed; never re-locked; its 154-PASS frozen-suite evidence is against a kernel 0007 deletes | 0009 Phase 1 lands on a kernel whose `Refusal` no longer has the field its fixtures assert; whichever RDR lands second re-decides frozen tests under pressure (0007 A25) | §1, AT-2 |
| C-5 | 0007 / JDR | 0007 Prereq "[ ] **JDR 0001 §D4 is reopened and §JD-8 amended to grant ONE `omitempty` structured field on `CLIError` (A19).** … a REVERSAL of that ruling"; JDR §D4 "it reaches the CLI through §JD-8's `Detail`" (unchanged) | A Final RDR carries, as an unchecked Prerequisite, the reversal of a closed home decision; the home was not reopened | 0005's renderer implements `Detail` per the home; 0007's payload contract says structured field; JSON consumers parse a sorted two-level array out of a prose string | §1, AT-3 |
| C-6 | 0005 / 0006 / 0007 / 0009 | 0005 A5 "not new envelope fields or exit groups"; 0006 Normative "an append-only typed `findings` field owned by `clierr` … MUST NOT be `omitempty`"; 0009 Normative "a NEW omitempty field … a clierr-local representation … MUST NOT be `resolve.RowRef`"; 0007 A19 "ONE `omitempty` structured field"; JDR §JD-8 "*(blank, except the exit-3 call)*" | Four Final documents give four answers about one shipped 60-line struct (`clierr.CLIError`); the home entry that owns it is blank | Three uncoordinated fields land on the envelope, each with its own serialization rule; text mode enumerates one and drops the others | §1, AT-3 |
| C-7 | 0002 / 0003 | 0003 Illustrative Code "The `single_valued` spelling above is illustrative: the field's authoring location is RDR 0002's to fix (**A16**)"; 0002 Normative "RDR 0003 is the normative home of that model … MUST NOT be restated here"; 0002 "The decoder MUST reject unmapped keys" | The TOML keys for `kind`, `domain`, `min`/`max`, `optional`, `single_valued`, `elements` are spelled by no normative text — 0002 owns location, 0003 owns meaning, nobody owns the byte | Author writes `kind = "enum"`; strict decoder returns `unknown schema field`, or the implementer picks a spelling the spec cannot review | §3, AT-4 |
| C-8 | 0002 / 0004 | 0004 Technical Design "Each definition declares a stable name, capability, artifact role, expected tag keys, timeout policy, and whether read-back verification is required"; 0004 Normative "Missing or non-positive timeout metadata MUST fail validation before execution"; 0002 Normative "An `[accessors.<id>]` entry carries `mode` and `path`, as the spike fixtures author it" | No authored form exists for capability/role/keys/timeout/read-back; under strict decoding they cannot be added, and without them 0004's validator MUST reject every accessor | Every accessor definition fails validation (`missing_or_non_positive_timeout`) or fails load (`unknown schema field`) | §3, AT-4 |
| C-9 | 0005 / 0004 / 0002 / 0007 | 0005 Normative "`flow resolve` … MAY invoke declared gate accessors … when the matched candidate requires a gate fact; gate denied or indeterminate results MUST surface as stable CLIError refusals" and "MUST NOT … coerce gate … results into tag values"; 0004 gate capability; 0002 grammar has no rule-level gate reference; 0007 guards are atoms over tag keys | Gate accessors exist in 0004 and 0005 and are unrepresentable everywhere else: no authored binding, no kernel input, no guard atom can name one | `flow next`/`resolve` gate binding is dead surface; either it is silently dropped or an implementer invents a rule-level `gate` key and a kernel input | §3, AT-4 |
| C-10 | 0003 / 0007 / 0004 | 0007 Phase 3 "**Blocked on one RDR 0003 declaration.** `resolve.Tag.Value` is a bare `string`, so a set-valued tag reaches the seam as one opaque string with no stated element encoding"; 0003 Final set-literal clause covers the *literal* only; 0004 "typed tag values" | The runtime encoding of a set-valued tag *value* at the accessor→kernel→seam boundary is unowned; 0003 locked without it | `contains` cannot be implemented or contract-tested; a set tag read by an accessor arrives as an unparseable string and every `contains` guard is `uncomparable` | §3, AT-4 |
| C-11 | 0007 / 0005 / JDR | 0007 Testing Strategy row 16 "`--tag X=v` against `X eq v` ⇒ row selected and a plan"; 0007 Normative "PRESENCE IS PROVENANCE-BLIND"; JDR §JD-9 "Caller-supplied tags enter as `ProvenanceObserved` and never satisfy an owned-state dependency … Wire it to `Observed`" | JD-9's answer defends the owned path, which ADV-4/ADV-5 already defend; the guard path is provenance-blind by 0007's own MUST, so the home's decision does not touch the exposure it was opened for. The set now specifies a *test* that pins the bypass as correct | Operator hits `guard_unevaluable`, re-runs with `--tag`, gets a plan whose writes derive from typed-in state | §1, premortem, AT-5 |
| C-12 | 0008 / 0007 | 0008 Scenario 3 "through a **new** view-capturing guard seam … (`fixtureGuards.Evaluate` takes `_ resolve.TagSet`)"; 0007 SEAM "`Evaluate(atom, value) GuardResult` … It never sees the view" | 0008's only net-new coverage (A2, "the test gap") is unwritable against 0007's normative seam; both Gate PASS 2026-08-21 | 0008 Phase 3 cannot deliver scenario 3; A2's "unpinned" property stays unpinned | §1, AT-2 |
| C-13 | 0006 | Load-Bearing "**Join rule and termination** — two edges reaching the same successor produce one node whose per-tag value sets are the union of theirs (a lattice widening)" | A node *is* its per-tag value-set vector, so "the same successor" means identical vectors and the union is the identity — no widening ever occurs and the reachable set is the full product; or the merge key is something the RDR never states. Either way the node ceiling is the operative bound | First real model (status × profile × stage × iteration × …) trips `graph-product-too-large` on the traversal; lint declines to prove anything | §2, AT-6 |
| C-14 | 0006 | Consequences "if that count is non-zero on a model the maintainers accept, the response is a successor RDR on guard-aware pruning"; scenario 23 | The RDR pre-authorizes its own rewrite on the first false positive against the only model that will exist (A8) | Lint blocks the checked-in model; the fix is a new RDR, not a waiver | §2 |
| C-15 | 0003 / 0006 | 0003 Normative "a declaration carrying no single-valued marker is not single-valued, so a finite kind takes the `2^\|domain\|` row"; operator/kind clause "Such an atom MUST take the blocking inability-to-prove outcome"; 0003 desk trace step 4 "**UNPROVABLE** … The defect is in the *fixture*"; 0006 "There is no suppression or waiver mechanism" | Conservative defaults make every unmarked `eq`/`in`/comparison guard a *blocking* lint failure; 0003's own canonical fixture is unprovable at defaults; A21 predicts the default inversion and is Pending | First author's table fails `graph-unprovable-coverage` on every group until they learn markers whose spelling nobody owns (C-7) | §3, premortem, AT-4 |
| C-16 | 0007 / 0003 | 0007 Prereq "[ ] **RDR 0003 states the empty/omitted `unless` identity (A9).** 0003 mentions the case zero times"; 0003 Final `unless` clauses: "the single conjunctive assignment set matched by the row's full `unless` block" | 0003 locked without the sentence; read literally an empty conjunction matches the whole product and every `unless`-less row is subtracted to nothing in lint while the kernel selects it — 0007's own words: "P5's promise failing in the direction that makes lint worthless" | Lint reports gaps on every group whose rows omit `unless`; or the implementer picks the identity and lint and kernel silently agree by luck | §1, AT-2 |
| C-17 | 0003 / 0006 / JDR | 0003 `authority` table "Escape-row overlap population … **Contested — A17.** §JD-14 decided for one population"; 0006 invariant 3 "Escape rows are checked for overlap in their own populations"; JDR §JD-14 "escape rows are ordinary participants in both the union and the overlap check" | Both members implement the reading the home rejected; the home still states the opposite; correction deferred to "the next JDR 0001 touch" which does not exist | An implementer reading the home builds one overlap population; reading either member builds two; `graph-overlap` verdicts differ by document | §1, AT-3 |
| C-18 | 0003 / 0006 / 0007 | 0003 Normative "**Conformance currently has no enforcer, and that is a booked gap** … A18 books the producer"; 0006 Normative "discharges only the owned half of RDR 0003's conformance premise"; 0007 states no view-conformance obligation and is Final | An observed key declared always-present is enforced by nothing; every lint claim in 0003 is conditional on it | Green lint, `guard_unevaluable` at runtime — the exact P5 breach the narrowing exists to prevent, arriving through the declaration nobody checks | §3, AT-6 |
| C-19 | 0004 / 0002 / JDR | JDR "STILL OPEN … JD-3"; 0002 Normative "`Row.RequiresOwned` has no authored form. The normalizer MUST derive it"; 0004 Capability Dependencies "`Row.RequiresOwned` populated by some layer … **Open** … JD-3 records the field appears zero times in 0002"; 0007 Prereq "[ ] JDR 0001 §JD-3 landed in RDR 0002" | Three documents disagree on whether a decision 0002 has already taken normatively is closed; the home lags its members | Nothing user-visible yet; the home cannot be trusted as the artifact of record for what is decided | §1 |
| C-20 | 0002 / 0008 / JDR | 0002 Normative "The `recognized` tag is total over matching … every view that reaches a row binds `recognized` (JDR 0001 §JD-10)"; 0008 "**Boundary against JDR 0001 §JD-10 (open).** … Blocks 1 and 2 do **not** settle it"; JDR §JD-10 "Still open: whether a declared `recognized` tag is total" | 0002 cites as decided what the home lists as open and 0008 says unsettled; 0002 decided it unilaterally in its own voice | An `exists = false` atom on `recognized` is dead under 0002, a lint error or satisfiable under JD-10's unanswered arms | §1 |
| C-21 | 0004 | Disposition note "A row that consumes the key through `Row.Match` or a guard string alone is dropped by `TagSet.matches` (which is provenance-blind) before `missingOwned` runs, yielding `no_match`" | False for guards under §D4 (an absent guard key is `guard_unevaluable`, non-escapable); a Final document restates the pre-D4 control flow as current | Implementer of 0004 scenario 6 expects `no_match` on the guard leg; the kernel refuses `guard_unevaluable` | §1 |
| C-22 | 0004 | Prerequisites "- [ ] All Critical Assumptions verified … - [ ] RDR 0001 keeps … - [ ] RDR 0002 carries … - [ ] RDR 0003 consumes" — all four unchecked; Status `Final`; Gate PASS 2026-08-21 | Unchanged since iteration 1 (C-29 there); the gate passed a document whose own checklist is empty | Implementer cannot tell whether 0004's preconditions hold | AT-1 |
| C-23 | 0005 | Failure Modes code table: twelve `flow-*` codes; no code for `owned_state_unavailable`, `reserved_tag_key`, `escape-row-shape-breach`, `incomplete_read`, `read_back_incomplete`, `unknown_accessor`, `capability_mismatch`, `execution_failure`; 0004 "mapping accessor refusals onto CLI codes is RDR 0005's"; 0007 "`flow-guard-unevaluable` … `GroupUserEnv` … mis-signals as a caller-input problem" | 0005 is the only user-facing member and the only one never re-entered; every sibling minted failure classes it does not map; JD-8 blank | Skill branches on documented refusal classes and finds no code; missing artifact state exits 2 as "bad input" | §1, AT-3 |
| C-24 | 0007 | Background "RDR 0003 and RDR 0002 are both **Draft**"; Technical Environment "RDR 0006 (`Draft`, re-entry at refine)"; Capability Dependencies "RDR 0003 (Draft)", "RDR 0002 (Draft, §JD-3)"; Phase 4 "Hand over, with 0003 in Draft" | Every peer-status claim in Final 0007 is false; the handoffs it schedules name a receiving Draft that does not exist | 0007 Phase 4's authoring guidance (sentinel anti-pattern, placement rule, conjoined-row grammar) is handed to nobody | §1, AT-1 |
| C-25 | 0008 | Testing Strategy "this RDR can reach 'implemented' with its user-facing outcome still unshipped … hold this RDR at Final-unimplemented until 0002's work starts and the JDs this RDR depends on close" | 0008 locks on the record that its own implementation is inert and that the cluster verdict is NOT RECONCILED | A green 0008 row on the index changes nothing an author perceives | §2 |
| C-26 | 0009 | Normative "Resolve MUST return CheckValid's error VERBATIM, never fmt.Errorf-wrapped"; "`Unwrap() []error` MUST be EXACTLY ONE LEVEL deep"; "exactly ONE RowRef field, Ref, plus a Count int field" | Carried unchanged from iteration 1; zero consumers; the first `flow` verb that wants command context is forbidden by four MUSTs | Rewrite on first real caller | §2 |
| C-27 | 0009 | Technical Design "callers MUST check the error before reading the disposition; a caller that branches on `Refused()` first would read a success-shaped value with a nil `Plan`"; scenario 7b "Asserted directly, because it is the one newly-introduced hazard" | The trap is documented *and tested* rather than removed | `flow resolve` branches on `Refused()` first, dereferences nil `Plan` | §2 |
| C-28 | 0008 / 0009 / JDR | JDR §JD-5 open; 0008 block 4 "either error is conforming"; 0008 scenario 6 fourth variant "must NOT assert which of the two errors is reported"; 0009 "precedes every modeled disposition" | Two entry preconditions on one `error` return, order unpinned; the test JD-5 asks for ("pin it with a test") is forbidden from asserting order until JD-5 closes, and JD-5 waits for the test | Doubly-breaching table reports a different error per build | §1, AT-3 |
| C-29 | 0002 | Risks "This RDR is the cluster's only unlocked document … the producer locks *after* its consumers, inverting the dependency order … **Mitigation**: … `/rdr-cluster-reconcile`, which is the mechanism that owns cross-RDR drift" | 0002 locked knowing 0006 A6/A10 and 0004's accessor metadata had not landed and delegated the consequence to this gate | This critique | §1 |
| C-30 | 0006 / 0005 | 0006 Normative "Text mode MUST enumerate every finding's code and message, which requires extending `clierr.EmitText` (failure) and the `respond.OK` text branch (success …) — both are booked edits"; 0005 Prereq "[ ] Add text success payload rendering through `respond.OK` or a respond-owned helper" | Two Final RDRs independently specify edits to the same `respond.OK` text branch; neither cites the other | Text rendering implemented twice with two shapes | §1 |
| C-31 | set | 65,070 lines of `docs/rdr`+`docs/jdr` markdown vs 3,971 Go lines; zero commits under `internal/`/`cmd/` since 2026-08-09; shipped `Row.Guard string`, no `OpExists`, no `CheckValid`, no `Undecided` | Iteration 1's 7:1 is now 16:1; every kernel change every member describes as "landing" has not landed | Binary still answers `version` | §1, premortem |
| C-32 | 0003 | A21 "Unverified today because no fixture declares the marker at all (A16 owns the authoring location, still `Pending` on RDR 0002)" beside A16 "Status: Verified" and Prereq "[x] **A16 closed** (2026-08-22)" | Internal contradiction inside a Final record | Cosmetic; reader cannot tell which sentence is current | — |
| C-33 | 0006 | Performance "Lint MUST therefore publish a **model-independent node ceiling** … and MUST emit `graph-product-too-large` naming the traversal (not a group)" beside the same code for "Declared finite product exceeds the published proof bound" | One code, two remedies, while `graph-unprovable-coverage` got a `reason` discriminator for exactly this problem | Machine consumer cannot tell "declare a smaller domain" from "the model has too many owned tags" | §2 |

---

## 1. The inter-RDR failure mode

### The failure

**The set will fail because it relocated every unresolved cross-document
question from "the peer's implement stage" (iteration 1) to "the peer's
refine", locked the peer without the refine, and left the relocation text in
place as if it still described a live venue.** Every member now carries
`Final`, and every member still carries Pending records, unchecked
Prerequisites, and peer-status sentences that name a Draft that no longer
exists. Under the project's own rule — RDRs are never amended — a Pending
record homed at a Final peer's refine is Pending forever. The queue nobody owned
in iteration 1 still exists; it has a new address.

The mechanism is visible in the two documents that locked last, on the same
day, each believing the other was still open:

- **0006** (Gate PASS 2026-08-23) needs a declared initial owned state and
  terminal predicates. It says so normatively — "A model that declares no
  initial owned state MUST be rejected with a blocking finding" — and books the
  producer as A6: "**RDR 0002's authoring schema** … RDR 0002 is `Draft` with
  every assumption Verified, so this reopens no locked document." Its Status
  line repeats it: "A6 and A10 are scheduled edits on RDR 0002 (`Draft`)".
- **0002** (Gate PASS 2026-08-23) locked with the layout enumeration unchanged:
  "root `outcomes`, `[model]`, `[tags.<tag>]`, `[accessors.<id>]`,
  `[context.<id>]`, `[[rule]]`, and `[dump]`", plus "`[model]` MUST contain
  `id` and `version`", plus "The decoder MUST reject unmapped keys so an unknown
  schema field is a stable refusal rather than a silent no-op."

Put the two together. A table with no `initial` is rejected by 0006. A table
with an `initial` is rejected by 0002. **Under the current Final text there is
no transition model that both documents accept.** This is not silence, which is
what iteration 1 diagnosed; it is a contradiction between two normative blocks
in two Final documents, produced by the refine-relocation mechanism on the very
day both were locked (C-1). The same mechanism leaves 0006's invariant 5
resting on a write-replaces rule that 0002 never states (C-2).

### Root cause across the RDRs

Three enablers, all new since iteration 1, all consequences of the fix applied
to iteration 1's diagnosis.

**Enabler 1 — the "scheduled edit on an open peer" idiom.** It appears in
0006 (A6, A10: "scheduled edit on an open peer, not a route-back"), in 0003
(A10, A12, A17, A18, A19, A20: "homed at RDR 0006's refine"), in 0007 (Phase 4
"with 0003 in Draft"; Prerequisites "Two duties are added to RDR 0002's
Refinement Context Direction list"), and in 0008 (A12 "a pointer lands in
0002's own surface as a cross-RDR edit"). The idiom is honest at the moment of
writing and false at the moment of locking, and the lock does not re-run the
check. 0006's own A2 states the trap in one sentence and walks into it: "they
close on the peer, not here — a Draft cannot flip a locked peer's record." When
0006 wrote that, 0003 was locked and 0006 was Draft. Then 0006 locked. Now
neither can flip the other's record, and 0003's Status still says its records
are "on RDR 0006's refine" (C-3).

**Enabler 2 — stale peer-status is not a lock condition.** 0007 is `Final`
(Gate PASS 2026-08-21) and describes 0002, 0003 and 0006 as `Draft` in its
Background, Technical Environment, Capability Dependencies and Phase 4 (C-24).
0006's Status line calls 0002 `Draft` on the day 0002 locked. 0003's body calls
0006 `Draft` and A21 calls A16 `Pending` two paragraphs after A16 is marked
Verified (C-32). 0004's Capability Dependencies says `RequiresOwned` "appears
zero times in 0002" while 0002 now derives it normatively (C-19). The README
index is correct; the members are not; and the members are where an implementer
reads. Nothing in the gate diffs a member's peer-status sentences against the
index, which is the same failure iteration 1 found on 0004's Prerequisites
(still all unchecked, still Final — C-22), one level up.

**Enabler 3 — the home lags the members and the members lag the home, in
both directions.** JDR 0001 is supposed to be the artifact of record for
joint decisions. Today:

- The home says JD-3 is open; 0002 has decided it normatively; 0004 and 0007
  still record it as open (C-19).
- The home says JD-10's totality question is open; 0002 asserts totality
  normatively "(JDR 0001 §JD-10)"; 0008 says blocks 1–2 "do **not** settle it"
  (C-20).
- The home says JD-14 puts escape rows in one overlap population; 0003 and
  0006 both implement two populations and 0003 marks the home's reading
  "Contested" with the correction deferred to "the next JDR 0001 touch" (C-17).
- The home's §D4 routes the refusal payload through `Detail`; 0007's
  Prerequisites require §D4 reopened and reversed, unchecked (C-5).
- The home's JD-8 entry is blank on the one shipped surface four Final RDRs
  now write contradictory MUSTs against: 0005 "not new envelope fields", 0006
  "typed `findings` field … MUST NOT be `omitempty`", 0009 "a NEW omitempty
  field … clierr-local", 0007 "ONE `omitempty` structured field" (C-6). 0006
  and 0005 additionally both specify edits to `respond.OK`'s text branch
  without citing each other (C-30).

"Cite, never restate" only works if the citation target is current. It is not.
The set has replaced restatement drift with pointer drift, and pointer drift is
worse, because a stale restatement is at least visible in the document you are
reading.

### The symptom the user sees

The user is a table author with the first real RDR transition model. They
write `[tags.status]` with `provenance = "owned"`, `kind = "enum"`,
`domain = [...]` — the spelling they copy from 0003's illustrative block. Load
fails: `unknown schema field` (C-7; 0002 does not spell the keys and its decoder
is strict). They strip the type model. Load passes. Lint fails: every group is
`graph-unprovable-coverage`, `reason: tag-not-single-valued`, because 0003's
defaults treat every unmarked enum as a power set and 0006 makes that blocking
with no waiver (C-15). They cannot add the marker, because they just learned
its key is rejected. They give up on lint, run `flow resolve`, and get the JD-9
speed bump: `guard_unevaluable`, then `--tag`, then a plan (C-11).

That is if `flow resolve` exists. It does not. `internal/cli/root.go`
registers `version`. 0005, the only user-facing member, is the only member
never re-entered since 2026-06-19 and has no code for any failure class its
seven siblings minted since (C-23). The cluster wrote 36,000 more lines of
design text since iteration 1 and zero lines of `internal/` (C-31).

---

## 2. The RDR that will be rewritten within six weeks of shipping

**RDR 0006 — Graph Lint Authority And Guarantees.**

Iteration 1 named 0009, and 0009 is still the runner-up for exactly the reasons
given then (C-26, C-27: four MUSTs that forbid the natural implementation, a
documented-and-tested nil-`Plan` trap, zero consumers). But 0009 will be
*patched* — one exported-API revision — whereas 0006 will be *rewritten*,
because its central algorithm is underspecified in a way no patch fixes, its
inputs do not exist, and it says in its own text that its first contact with
the real model triggers a successor RDR.

**Its inputs do not exist.** The reachability relation is "rooted at the
declared initial owned state (A6)" and invariant 2 reads "declared terminals".
0002 locked without either and cannot be amended (C-1). 0006 rejected a sidecar
("would split model authoring across two files"). So the first implementation
either reopens 0002 — a Final, `foundational` producer — or ships the sidecar
0006 rejected, or infers a root from the model, which invariant 7 forbids. All
three change 0006's normative text.

**Its algorithm is not specified.** The Load-Bearing "Join rule" says "two
edges reaching the same successor produce one node whose per-tag value sets are
the union of theirs (a lattice widening)". But a node is defined as the
per-tag value-set vector; "the same successor" can only mean the same vector;
two identical vectors union to themselves. Either no widening ever happens and
the reachable set is the full product of subsets — exponential in owned tags
times domain size, which the RDR itself concedes in Performance Expectations —
or the merge key is something the RDR never names (C-13). The RDR's answer to
the exponential case is a "model-independent node ceiling" that emits
`graph-product-too-large` naming the traversal. The RDR model has status,
profile, stage, prelock iteration, cluster eligibility, rewind target and lens
as owned or observed dimensions; that ceiling fires on it. Then invariant 6
(owned-before-write), invariant 2 (dead end) and `graph-always-present-owned`
are all unprovable, and the RDR's blocking authority is a blocking refusal to
lint.

**It pre-authorizes its own rewrite.** Consequences: "if that count is non-zero
on a model the maintainers accept, the response is a successor RDR on
guard-aware pruning — still not a suppression flag." Scenario 23 measures it.
Given 0003's defaults (C-15) and the guard-blind traversal, a non-zero count on
the first real model is not a risk, it is the expected value. The RDR names its
replacement before it has been implemented.

**Its escape-row and overlap semantics contradict the home** (C-17), its
findings field contradicts 0005 (C-6), its text-mode rendering edits overlap
0005's (C-30), its single-valued invariant rests on a rule 0002 never states
(C-2), and its A5/A8 CI gate has no subject (still Pending, still no checked-in
model). Twenty-three scenarios, none runnable, several — 11, 17, 20, 21 —
designed as regression guards for design choices no code has made.

The one thing that would have to be true for 0006 *not* to be rewritten: 0002
is reopened first and absorbs `initial`, `terminal`, write-replaces, the
type-model keys and the accessor metadata in one refine before any
implementation starts. That is the same event that would make Step 3's
assumption survive, and nothing schedules it.

*(Runner-up, for the record: 0002 — not because it is wrong, but because it is
the producer that locked last with every consumer's un-landed request pointed at
it, and it is where C-1, C-2, C-7, C-8 all have to land. It will be reopened
first; 0006 will be rewritten.)*

---

## 3. The cross-cutting assumption that will not survive first contact

**The assumption: that splitting a contract into "authoring location" (0002)
and "meaning" (a peer) is a complete assignment — that once meaning has an
owner, the wire spelling will exist.** It will not, because 0002's decoder is
strict by MUST, and a strict decoder turns an unspelled key into a rejected
key. The doctrine "one home per contract; cite, never restate" (JDR P6) has been
applied so thoroughly that the bytes an author types have no home at all.

This assumption is not written down anywhere, which is why it is load-bearing
everywhere. Its instances:

- **The tag type model.** 0002: "the type model is RDR 0003's and is cited, not
  restated." 0003: "The `single_valued` spelling above is illustrative: the
  field's authoring location is RDR 0002's to fix (**A16**)." A16 is Verified
  because 0002 *mentions* five fields by prose name; neither document spells
  `kind`, `domain`, `min`, `max`, `optional`, `single_valued`, `elements` as
  TOML keys. 0003's own 3amigo pass blocked lock on exactly this (A16, iteration
  5) and un-blocked on a prose sentence (C-7).
- **Accessor definitions.** 0004 requires name, capability, artifact role,
  expected tag keys, timeout, read-back flag, and refuses a definition missing
  any of them. 0002 gives `[accessors.<id>]` exactly `mode` and `path` "as the
  spike fixtures author it" and validates nothing else. Every accessor 0004 can
  execute is a load failure in 0002, and every accessor 0002 loads is a
  validation failure in 0004 (C-8).
- **Gate accessors.** 0004 declares them; 0005 invokes them from `flow next` /
  `resolve` "when the matched candidate requires a gate fact" and forbids
  coercing the result to a tag. 0002's grammar has no rule-level gate
  reference; 0001's `Input` has no gate channel; 0007's guard atoms range over
  tag keys. A "gate fact" a candidate "requires" is expressible in no document
  (C-9).
- **Set-valued tag values.** `resolve.Tag.Value` is a `string`. 0003 canonicalizes
  the *literal*; nobody says how a set-valued *value* read by an accessor is
  encoded when it crosses to the kernel and on to the seam. 0007 Phase 3 says
  its `contains` contract test "cannot be written from any current document"
  and hands the request to a 0003 that has since locked without answering
  (C-10).
- **Initial and terminal declarations.** 0006 needs them, 0002 has no key for
  them, and 0002's decoder rejects the key an author would add (C-1).

Every one of these is a case where two Final documents each say "the other
owns that part", and the part they both disclaim is the part the author types.
Iteration 1's floor was "tags reaching the kernel are a faithful picture of the
artifact"; that floor was fixed (§D3). This floor is one layer up: "the author
can write what the set specifies." They cannot.

### How it breaks

**First contact, the author.** They write the model with the fields the
documents describe. Load rejects every declaration key beyond `provenance` and
`accessor`. They strip to what loads. 0004's validator rejects every accessor
for missing timeout. They remove accessors. 0006 rejects the model for no
initial state. They add `[initial]`; load rejects it. There is no fixed point
in which a model with a guard, an accessor and a root passes both 0002 and
0006 as written.

**First contact, the implementer.** Faced with the above, the 0002 implementer
does the reasonable thing: invents the keys. `kind`, `domain`, `single_valued`,
`optional`, `elements`, `initial`, `terminal`, and under `[accessors.<id>]`
`capability`, `role`, `keys`, `timeout`, `read_back`. None is reviewable
against a normative sentence, so review cannot reject any spelling, so the
first implementation's guess becomes the wire format of a `foundational` RDR by
default. 0002's own strict-decoding clause was written to prevent "a parser swap
that silently lost it would retire the `unknown schema field` category without
any contract appearing to change" — and the same clause now guarantees that the
first implementer defines the schema without any contract appearing to change.

**First contact, the CLI.** The one place code exists, `clierr.CLIError`, is
the one place four Final RDRs disagree (C-6). The implementer of whichever
lands first adds their field; the second finds a struct with a field they were
told would not exist and a text renderer that enumerates the first field only.
The home entry that owns it is blank.

---

## 4. Premortem

*Written from nine months after the cluster was declared Final.*

We implemented in dependency order, as the index said: 0002, then 0003 and
0007 together, then 0004, 0006, 0008, 0009, and finally 0005. Every RDR reached
`Implemented`. The product does not accept the RDR flow's own transition model,
and we have three open successor RDRs to explain why.

**Weeks 1–3, `internal/model`.** The 0002 implementer built `Load` against the
spike fixtures, exactly as the "canonical examples" clause requires. Strict
decoding, two-pass version gate, twenty-two validation categories, byte-exact
keys, the atom-set merge, `<clear>` writes, `RequiresOwned` derivation — all
green against `rdr-fixture.toml` and `kata-fixture.toml`. Then they tried to
declare an enum domain so 0003's proof could run. There was no key. They read
0003: "illustrative … RDR 0002's to fix." They read 0002: "cited, not
restated." They picked `kind`, `domain`, `optional`, `single_valued`,
`elements`, `min`, `max`, wrote them into `internal/model/schema.go`, and
opened a note. Then 0004's accessor definition: `capability`, `role`, `keys`,
`timeout`, `read_back` under `[accessors.<id>]`. Same note. Then 0006's
`initial` and `terminal`. Same note. The note became RDR 0010, *Transition
model wire keys*, seeded from the implementation and locked in a week, because
it had to describe what already shipped.

**Weeks 4–7, the kernel reshape.** 0007 Phase 1 replaced `Row.Guard string`
with `[]GuardAtom`, narrowed the seam to `Evaluate(atom, value)`, and
re-encoded 154 frozen tests. 0009's `Table.CheckValid` was scheduled for the
same pass because both rewrite `fixtures_test.go::escapeRow` (0007 A25). The
0009 implementer opened 0009 and found every assertion in its A3 evidence
anchored on `Refusal.Guard`, a field that no longer existed, and a scenario 3
that "asserts on `Refusal.Kind` / `Refusal.Guard`". 0009 had never been
re-locked after §D1/§D4; 0007's Prerequisite said it would be; nobody had. They
re-ran the frozen suite against the reshaped kernel, found four dispositions
that "re-decided" rather than "re-encoded", and — under schedule pressure,
exactly as 0007 A25 predicted — pinned them to whatever the reshaped kernel
returned. 0008's scenario 3 was dropped outright: its "view-capturing guard
seam" cannot capture a view the seam is normatively forbidden to see.

**Week 8, the envelope.** 0006 landed `Findings []Finding` on `CLIError`,
non-`omitempty`, and extended `EmitText`. 0009 landed `Rows []RowIdentity`,
`omitempty`. 0007's `Undecided` payload was rendered into `Detail` because the
JDR's §D4 said so and nobody had reopened it; the implementer serialized a
sorted two-level array into a prose string. 0005 was implemented last, against
a `CLIError` with two new fields its normative text says do not exist and an
`EmitText` that enumerates one of them. A skill parsing `--as=json` found
`findings` present-and-empty on every failure and `rows` absent on most, and
`detail` occasionally containing JSON.

**Week 10, first real table.** The RDR-flow model, authored by the person the
cluster was written for. `intrastate lint --model rdr.toml`. Result:
`graph-lint-failed`, forty-one findings. Nineteen `graph-unprovable-coverage`
with `reason: tag-not-single-valued` — every `eq`/`in` guard over `status`,
`profile`, `stage`, `lens`, none of which had been marked, because the marker
key was RDR 0010's and the author had read 0003. They marked them. Eleven
`graph-unprovable-coverage` with `reason: row-can-refuse` — every guard over an
optional key, which by 0003's default was every key. They marked them
always-present. Then `graph-always-present-owned` on the root, because
`[initial]` did not list them all. Then, with the model finally provable, one
`graph-product-too-large` naming the traversal: seven owned tags, domains of 4
to 6, and the "lattice widening" in `reach.go` that never merged anything
because the merge key was the value-set vector itself and no two paths produced
an identical one. The node ceiling was 50,000. The model had 2.4 million. Lint
declined to prove anything. Scenario 23 recorded a non-zero count. Per
Consequences, the response was a successor RDR: 0011, *Guard-aware reachability
pruning*. 0006 was `Implemented` and unusable on the one model it exists for.

**Week 12, the operator.** `flow resolve` shipped. First `guard_unevaluable`
in anger: `reason: absent`, key `cluster_eligible`, on a kata whose
`.kata.toml` the accessor could not fully read. The refusal named the key
plainly, as 0007 promised. The on-call did what the key told them to:
`--tag cluster_eligible=true`. Plan. Writes applied. 0007 Testing Strategy row
16 had shipped a test asserting that this is the correct disposition. JD-9's
"wire it to Observed" had been done to the letter and made no difference,
because guard presence is provenance-blind by 0007's own MUST. The runbook
gained the entry iteration 1 predicted, word for word.

**Week 14, the gate accessor.** A flow needed "proceed only if CI is green" —
the gate capability 0004 defined and 0005 exposed. There was no place in the
table to attach it to a rule and no way for the kernel to consume the result.
The `flow next` gate-binding code path had shipped as dead code behind a flag
nothing set. RDR 0012, *Gate accessor binding*, seeded.

**What actually killed us.** Not a wrong decision — §D1, §D2, §D3, §D4 were
all right, and the atom reshape was the best kernel change we made. What
killed us was that we treated `Final` as a property of a document and never as
a property of the set: 0002 and 0006 locked on the same day each believing the
other was Draft, and after that there was no document that could receive the
five wire-key questions, the write-replaces rule, the empty-`unless` identity,
the set-value encoding, the envelope field, or the JD-14 correction. We had a
home for joint decisions and let it fall behind its members in both directions
until the members cited it for things it said the opposite of. And we grew the
design corpus by 36,000 lines between the two cluster gates while the code the
corpus describes as "landing" stayed at `Row.Guard string`. Everything above
would have surfaced in the first week of implementing 0002 against a real
table. We spent that week, and eleven more, locking.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

RDR-review-time gates, executable against the document set as it stands.

### AT-1 — A record homed at a peer's refine blocks lock while the peer is Final
*Catches C-1, C-2, C-3, C-22, C-24*

```gherkin
Given an RDR carries a Pending assumption, an unchecked Prerequisite, or a
  "scheduled edit" whose named venue is another RDR's refine or "next touch"
When the Finalization Gate runs for the carrying RDR
Then the named venue RDR MUST be Draft in the README index at gate time
And if it is Final, the gate MUST FAIL, because the venue cannot receive the
  edit and the record is Pending forever
And every sentence in the RDR asserting a peer's status MUST match the index
  (a mechanical grep for "Draft"/"Final" beside "RDR 000N")
```
Run today: 0006 (A6, A10 → 0002 Final), 0003 (A10, A12, A17, A18, A19, A20 →
0006 Final), 0007 (Phase 4, Prerequisites → 0002/0003 Final), 0004
(Prerequisites all unchecked) → four failures.

### AT-2 — A Final RDR that a later decision made stale must be re-locked before the deciding RDR locks
*Catches C-4, C-12, C-16*

```gherkin
Given RDR X's normative text deletes or reshapes a surface a Final peer Y
  cites as an asserted property (0007 deletes Refusal.Guard; 0007's seam
  never sees the view; 0007 requires 0003 to state the empty-unless identity)
When the Finalization Gate runs for X
Then Y MUST already carry a re-lock note citing X's decision, or the gate
  MUST record Y as Draft
And "rides to Y's re-lock" MUST NOT satisfy this while Y's re-lock is not
  scheduled
```

### AT-3 — Every joint decision the members cite must read identically at the home
*Catches C-5, C-6, C-17, C-19, C-20, C-23, C-28, C-30*

```gherkin
Given a member cites JDR 0001 §D-n or §JD-n as decided, or carries a
  Prerequisite that reopens one
When the cluster reconcile pass runs
Then for each cited entry the home's text and every member's fenced normative
  text MUST agree on the answer, or the pass MUST emit ANSWER-CONTRADICTS
And a home entry marked "(blank)" MUST FAIL the pass if two or more members
  write a MUST against the surface it names
And no member may cite an entry as decided that the home lists as open
```
Run today: JD-3 (home open, 0002 decided), JD-10 (home open, 0002 decided,
0008 says unsettled), JD-14 (home says one population, both members say two),
§D4 (home says Detail, 0007 says structured field), JD-8 (blank, four MUSTs) →
five failures.

### AT-4 — Every field a consumer requires must have a spelled wire key in the producer
*Catches C-7, C-8, C-9, C-10, C-15*

```gherkin
Given 0002 owns the authored TOML layout and MUST reject unmapped keys
When any peer's normative text requires a declaration, definition, or
  reference the author must write (type model fields, accessor metadata,
  gate references, initial/terminal, set-value encoding)
Then 0002's layout enumeration MUST name the exact key, or the peer MUST
  name it and 0002 MUST cite that spelling
And "authoring location is 0002's, meaning is the peer's" with no key spelled
  by either MUST FAIL the gate for both documents
And every field a peer's validator MUST reject when missing MUST be
  authorable under 0002's layout
```

### AT-5 — A guarantee must be evaluated against every producer, including the one the home routed
*Catches C-11*

```gherkin
Given 0007 claims missing artifact state cannot be masked, and JD-9 answers
  the --tag channel with "wire it to Observed"
When the Gate evaluates the guarantee
Then it MUST trace the channel through the guard path's presence rule
  (provenance-blind), not only the owned path
And a test that pins the bypass as the expected disposition (0007 row 16)
  MUST be recorded as an accepted exposure in the Problem Statement, not as
  conformance
```

### AT-6 — An algorithm stated as load-bearing must name its keys and be run on the target model
*Catches C-13, C-14, C-18*

```gherkin
Given 0006's reachability relation is the input to invariants 2, 5, 6, 7
When the Gate runs
Then the merge key of the fixpoint MUST be stated (what makes two successors
  "the same")
And the node count for the representative RDR model MUST be computed from
  its declared owned tags and domains and compared to the published ceiling
And an RDR whose own text names a successor RDR as the response to its
  first false positive MUST NOT lock until scenario 23 has a subject
And every premise a lint claim is conditional on (0003's conforming view)
  MUST name an enforcer for each provenance, or the claim MUST be scoped to
  the provenances that have one
```

### AT-7 — The set-level delivery gate, restated
*Catches C-25, C-29, C-31*

```gherkin
Given a cluster of Final, unimplemented RDRs
When any member would be locked or re-locked
Then the pass MUST report lines of RDR/JDR text against lines under internal/
  and the date of the last internal/ commit
And a producer RDR MUST NOT lock after its consumers while a consumer's
  request against it is unlanded (0002 Risks names this and locked anyway)
And a member whose Testing Strategy states its own implementation delivers
  no user-perceptible outcome (0008) MUST NOT count toward "Final" for
  ordering purposes
```

---

## 6. Delta vs iteration 1

Comparison against `cluster-reconcile/0002-0009/critique-set.md`
(claude-opus-5, 34 rows). Status per row: **CLOSED** (the home answered and
current member text is consistent), **STILL OPEN** (no change, or the same
defect in current text), **SUPERSEDED** (the finding no longer exists in the
current text, usually because a JDR decision dissolved it). Where a row has
mutated rather than closed, the successor row in this file's ledger is named.

| Iter-1 ID | Status | Note |
|---|---|---|
| C-1 (0007 A10 → 0003 Phase 1) | SUPERSEDED | §D1: row carries atoms; no guard encoding to specify. |
| C-2 (0007 A12 verified-by-silence) | CLOSED | 0003 A7 states the presence-dimension projection normatively; 0006 adopts row groups by citation. 0007 A12 is now honestly DOWNGRADED. |
| C-3 (0007 A15 conjoined-row grammar) | SUPERSEDED | 0002's atom-identity clause admits two atoms per key; A15 no longer exists. 0007 Phase 4 still "hands over" authoring grammar to a Draft 0003 that is Final — folded into C-24 here. |
| C-4 (0007 harness UNMITIGATED) | SUPERSEDED | §D4: kernel enforces the domain rule; the harness is gone. |
| C-5 (0003 Phase 1 too small) | SUPERSEDED | The routed obligations dissolved with §D1/§D4. A new receiving-hole exists (C-3, C-16 here). |
| C-6 (0003 Prerequisites checked and stale) | CLOSED | 0003's Prerequisites are now itemized per assumption with honest `[ ]`. The unchecked ones are the new problem (C-3). |
| C-7 (obligations accumulate in an unowned queue) | STILL OPEN | Mutated: the queue moved from "implement stage" to "refine of a Final peer". C-1, C-2, C-3, C-4, C-5, C-16, C-24 here. |
| C-8 (0007 mandated panic) | CLOSED | §D1/§JD-6. |
| C-9 (0005 envelope vs panic) | CLOSED | §D1/§JD-6. |
| C-10 (0007 A6b read completeness Pending) | CLOSED | §D3; 0007 A6b Verified by citation; 0004 carries the clause. |
| C-11 (0004 no completeness requirement) | CLOSED | 0004 normative completeness clause + seam-omission clause + scenarios 6/7. |
| C-12 (0007 A13 --tag bypass) | STILL OPEN | JD-9 open; JD-9's recommended answer does not reach the guard path; 0007 row 16 now tests the bypass as correct. C-11 here. |
| C-13 (0005 --tag unconstrained) | STILL OPEN | 0005 text unchanged. C-11 here. |
| C-14 (0008 producer obligation, CLI producer) | STILL OPEN (narrowed) | Detection now exists (0008's `Resolve`-entry predicate); classification of a user-typed reserved key is JD-9, open. |
| C-15 (0008 breaks all three fixtures) | CLOSED | JD-10 rules the rename; 0008 A4/A9/Phase 2 carry the inventory. |
| C-16 (0008 intent gap lints clean) | STILL OPEN | 0008 block 2 states it as out of scope. |
| C-17 (0009 frozen API) | STILL OPEN | Unchanged. C-26 here. |
| C-18 (0009 nil-Plan trap) | STILL OPEN | Now also tested (7b). C-27 here. |
| C-19 (0009 RowRef{"",""} collapse) | STILL OPEN (mitigated) | `Count` field added; distinct rows still unlocatable. |
| C-20 (0006 row-group membership undefined) | CLOSED | 0003 defines the group; 0006 "Source state is the authored match pattern". |
| C-21 (0002 ambiguous overlap no valve) | CLOSED | 0002 moved overlap to lint by arity; 0003's presence dimension makes the two-row pattern provable. |
| C-22 (Implemented while xg7p open) | SUPERSEDED | §D4 changes the kernel; 0007 Consequences now say the probe inverts once Phases 1–2 land. |
| C-23 (reconcile gate all-clear on silence) | SUPERSEDED | A JDR home exists and iteration 1 returned NOT RECONCILED. The home now lags (C-19, C-20, C-17, C-5 here). |
| C-24 (spec:code ratio) | STILL OPEN, worse | 16:1 by markdown lines; zero `internal/` commits since 2026-08-09. C-31 here. |
| C-25 (0007 MVV non-validating) | CLOSED | MVV runs against the real kernel with a real atom. |
| C-26 (view-reading evaluator does not exist) | SUPERSEDED | §D4: domain-rule vectors are kernel tests. |
| C-27 (0008 conformance test needs 0002's normalizer) | STILL OPEN | 0008 now states the split explicitly and holds itself at Final-unimplemented. C-25 here. |
| C-28 (0005 last in dependency order) | STILL OPEN | 0005 unchanged; now also unmapped for every sibling's failure class. C-23 here. |
| C-29 (0004 Prerequisites unchecked vs Final) | STILL OPEN | Verbatim. C-22 here. |
| C-30 (0004 self-reference claim false) | NOT RE-CHECKED | Gate responses moved to `artifacts/gate.md` (not read here); A1/A2/A6/A7 still cite 0004's own spike directory. |
| C-31 (0006 A5 CI gate Pending) | STILL OPEN | A5 and A8 Pending; no checked-in model. |
| C-32 (0005 no-new-fields vs 0006 findings) | STILL OPEN, widened | Now four documents (0005, 0006, 0007, 0009) on one struct; JD-8 blank. C-6, C-30 here. |
| C-33 (read-before-write, three owners) | CLOSED | 0002's arity split assigns it to 0006; 0003 keeps the normative clause but quantifies over the relation 0006 defines and 0006 cites 0003's clause. One relation, one code. |
| C-34 (owned_state_unavailable has no CLI code) | STILL OPEN | JD-8 open; 0007 disposition table says "no code row yet". C-23 here. |

Tally: 13 CLOSED, 7 SUPERSEDED, 13 STILL OPEN (three widened), 1 not
re-checked. What closed was almost entirely what the JDR's §D1–§D4 decided.
What stayed open is everything routed to a still-blank or still-open home
entry (JD-5, JD-8, JD-9, JD-10) plus 0005 and 0009, the two members never
re-entered. The net-new failures in this file (C-1, C-2, C-7, C-8, C-9, C-10,
C-13, C-16, C-17, C-20) are all products of the two lockings on 2026-08-22/23.

---

*Ledger rows: 33. Acceptance tests: 7. Sections 1–5 are the argument; the
ledger is the routable record. Any defect not carried in the ledger is
unreported by construction.*
