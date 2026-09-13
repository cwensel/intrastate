Model: claude-opus-5[1m]

# Whole-set Critique — {0021, 0029} at Stage 7.1

## Findings Ledger

| ID | RDR | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-----|-------------|--------------|-------------------|--------|
| C-1 | 0021 | `0021:C2` | 0029's settled joint decision obliges 0021 to assign tiers to `--emit`, `graph-export-too-large`, and C2's field spellings "in its own Normative Contracts." C2 assigns none. The tier words `frozen`/`append-only`/`growing` appear nowhere in 0021's normative text. | An agent reads `docs/cli-output-contract.md`, finds a tier for every vocabulary the CLI emits, finds none for the graph document's fields, and cannot tell whether `rows[].gate` may be renamed at 0.3.0. It guesses. | §1, §3, premortem, AT-1 |
| C-2 | 0021 | `0021:C1` | `--emit`'s value set (`json`, `dot`) is an emitted machine-readable vocabulary — an unknown value produces `flag-invalid-value` naming `emit` — and takes no tier. By `0029:C4`'s own rule ("an unassigned machine-readable surface is a defect") 0021 ships a defect on day one of its own implementation. | A consumer writing `--emit=mermaid` support guard code cannot tell whether the set is closed. Adding a third format later silently reclassifies its error handling. | §1, §3, AT-2 |
| C-3 | 0021 | `0021:§approach` | 0021 introduces a whole new machine-readable surface (an entire JSON document schema, ~11 top-level members plus nested row/atom/tag/reach objects) and treats "the document versions itself via `schema`" as discharging the stability question. Self-versioning answers *what changed*; it does not answer *what is allowed to change*, which is the tier. | An agent pins `intrastate.graph/1`, sees the marker unchanged across 0.2.0→0.3.0, and assumes the field list is stable — but C2's additive rule permits new fields within `/1` and says nothing about removal being forbidden beyond the bump. | §1, §3, premortem |
| C-4 | both | `0029:A4` | A4 claims the three tiers "partition every machine-readable surface this CLI emits" and is marked **Verified**. It was verified against `main`, where 0021's surface does not exist. The census is verified against a tree the cluster's own build order guarantees will be obsolete before 0021 lands. A4's own "If wrong" text names this exact recurrence and calls it "the residual risk this assumption carries." | Nothing, until an agent depends on an unpromised field. Then a silent break with no record to appeal to. | §3, premortem, AT-1 |
| C-5 | 0029 | `0029:C4` | C4 discharges the 0021 overlap by declaring the census "complete as of this record's implementation, not complete for all time" and delegating 0021's tiers to 0021. The delegation has no enforcement: A7's snapshot check reads *enumeration seams*, and an unassigned surface with no seam produces no seam to diff. The mechanism cannot detect the absence it was built to detect. | A consumer of the graph export operates on an untiered surface indefinitely; CI stays green through every release that widens or narrows it. | §1, §4, premortem, AT-3 |
| C-6 | 0029 | `0029:C3` | C3's mechanical backing (A7's snapshot) covers seams and severities only. 0021 introduces `graph-export-too-large` into the CLIError `code` vocabulary — the one `append-only` row C4 records as `seam: none (prose-only)` because A5 refuted the seam. So the single vocabulary 0021 provably grows is the single vocabulary 0029 provably cannot police. | A new refusal code appears with no snapshot movement, no release-note trigger, and no test failure. A consumer branching exhaustively on `code` takes its default branch. | §1, §3, premortem, AT-4 |
| C-7 | both | `0021:C5` | C5 embeds the graph document in `data` under `--as=json`. `0029:C1` adds `schema_version` at the envelope top level. Neither record states which marker an agent reads first, nor what an agent does when the envelope major is supported but `intrastate.graph/1` is not (or the inverse). The two-marker protocol is declared to coexist and never sequenced. | An agent that rejects on unsupported envelope major passes a graph document it cannot parse; an agent that checks only `schema` parses an envelope whose shape moved. | §1, §4, AT-5 |
| C-8 | 0021 | `0021:S2` | S2 pins golden fixtures for `intrastate.graph/1` byte-for-byte. `0029:S7` declares a repo-wide negative that byte-goldens over emitted output "would itself assert the cardinality C2 forbids consumers from asserting," and states "no golden/snapshot test of the JSON envelope exists in the repo today, and this RDR does not add one." Under `--as=json`, 0021's golden IS an envelope golden — and it will carry `schema_version`. | CI breaks on the release that moves `schema_version` `"0.1"`→`"0.2"`, for reasons unrelated to the graph document. A maintainer re-baselines and the additive tripwire C2 wanted is trained to be ignored. | §1, §4, AT-6 |
| C-9 | 0021 | `0021:§decision-rationale` | 0021's joint-check reads "clear (12 peers) … share no undecided contract," and its peer-citation paragraph names 0013, 0014, 0015, 0022, 0023, 0024 — but not 0029, the record that fired a joint decision AT it. The rationale records the wrong answer to the question the cluster exists to ask. | Invisible to the user; visible to the next implementer, who reads 0021 alone and never learns a peer constrained it. | §1, §2, AT-7 |
| C-10 | 0029 | `0029:C4` | C4's 0021 paragraph asserts 0021 "is an early Draft with every assumption still Pending, so it takes the citation in its own flow." 0021 is Final, dated 2026-08-28, with A1–A8 all Verified — and was finalized at Gate PASS before this cluster convened. The premise of the delegation is factually false about the record it delegates to. | None directly. But the delegation was accepted on a premise that would have forced a different remedy (an amendment, not a "flow") had it been checked. | §1, §2, AT-7 |
| C-11 | 0021 | `0021:A4` | A4 verifies that a new root verb needs no amendment to `0005:C1`'s envelope contract, reasoning from 0005's scope sentence. It is silent on 0029:C1, which adds a field to *every* terminal record at the gateway — including the `graph` verb's. 0021's envelope question was verified against the envelope as it stood, not as the cluster's higher-priority record leaves it. | The `graph` verb's `--as=json` envelope carries a field no passage in 0021 mentions. Any 0021 test asserting an exact envelope key set fails on integration. | §1, §4, AT-8 |
| C-12 | 0029 | `0029:§activation-step-3-promote-the-schema-to-1-0-at-1-0-0` | The 1.0.0 gate requires "the C4 census re-run against the then-current tree: every emitted machine-readable surface appears in the table." If 0021 has landed without tiering its surfaces, this gate fails at the worst possible moment — at the release where every tier becomes load-bearing at once — and the recorded remedy is "take the 1.0.0 slip." | A release slip attributed to a low-priority export verb shipped months earlier. | §2, §4, AT-1 |
| C-13 | 0021 | `0021:§approach` | Build order is 0029 then 0021, but 0029 is the record that *names* 0021's obligation and 0021 is the record that must *discharge* it. 0021 is Final and locked; discharging requires amending a Final record's Normative Contracts. The cluster's build order guarantees the obligation is created before the only record that can satisfy it is legally editable. | Structural: the set cannot be implemented in its stated order without a demotion. | §1, §2 |
| C-14 | 0029 | `0029:C2` | C2's `append-only` tier forbids consumers from asserting "the set's cardinality, a member's ordinal position, or a tail position." `0021:C2` fixes `rows[]` field order as "RDR 0002's dump field list in its canonical row order … a closed 11-member list with `emit` appended last per `0010:C3`" — a cardinality and a tail position, normatively, on a wire surface. | A consumer that reads 0021's C2 literally writes exactly the assertion 0029's C2 forbids, with each record citing the other as authority. | §1, §3, AT-9 |
| C-15 | 0029 | `0029:C2` | C2 prohibits the descriptive word `closed` "repo-wide, not one package," and S9 enumerates the sites. `0021:C2` uses "a closed 11-member list" to describe `dumpColumns`, and `0021:§briefly-rejected` / `0021:C1` reuse the framing. 0021's text is written today in the vocabulary 0029 retires tomorrow. | A reader hits both records and cannot tell whether "closed" in 0021 means `frozen` or means nothing. Exactly the ambiguity 0029 exists to kill, reintroduced by its cluster partner. | §1, §3, AT-10 |

---

## 1. The most likely INTER-RDR failure mode

**The defect: 0029 builds a census-plus-enforcement machine for machine-readable surfaces, then writes itself a delegation clause that exempts the only new surface in the cluster from both halves — and 0021, the delegee, never received the delegation.**

Read the two records alone and each is defensible. 0029 is a careful, heavily-spiked stability policy: three tiers, a `schema_version` on both terminal records, an enumeration seam per vocabulary, a CI snapshot that turns an undisclosed promotion red. 0021 is a clean, well-grounded export verb with byte-determinism proofs, a hostile-content DOT fixture, and a mechanism-independent neutrality oracle. Neither is sloppy. The defect is in the seam, and the seam is where both records stopped looking.

**Root cause.** `0029:C4` contains the passage that creates the defect:

> `cli/0021:C2`'s graph export document is exactly such a later surface, and it is named here because it is in flight against this record. … the vocabularies `cli/0021` emits under `data` — its `--emit` format set, its `graph-export-too-large` refusal code, and the field names its C2 fixes at Resolve — are machine-readable surfaces this census does not enumerate, because they do not exist on `main` yet. They take their tier assignments in `cli/0021` itself, in the change that adds them, per the rule above.

And `0029:§decision-rationale` settles it as a joint decision:

> 0021 then assigns tiers to the surfaces it introduces (`--emit`, `graph-export-too-large`, and C2's field spellings once Resolve fixes them) in its own Normative Contracts, per `0029:C4`. … 0021 is an early Draft with every assumption still Pending, so it takes the citation in its own flow.

Two things are wrong with this, and they compound.

First, **0021 is not an early Draft.** It is Final, dated 2026-08-28, all eight assumptions Verified, locked at a Gate PASS whose evidence sits in this repo's recent commits. 0029 is dated 2026-09-11. The delegation was written about a record that had already closed. The "flow" in which 0021 was to "take the citation" had ended two weeks earlier. This is C-10, and it is not a nitpick: it is the false premise on which the entire remedy rests. Had the author known 0021 was Final, the remedy would have had to be an amendment obligation, a demotion, or a hoist into a JDR — all of which are visible, tracked events. "It takes the citation in its own flow" is invisible and untracked, and so nothing happened.

Second, **the delegation is unenforceable by the very mechanism 0029 built to enforce it.** `0029:C3` backs the tier system with A7's snapshot check, and describes its reach precisely:

> the tiered-vocabulary seams C4 obliges are enumerable, so CI records each one's members … and fails when the current tree differs from it.

The check reads *seams*. A surface that was never assigned a tier has no seam obligation, therefore no seam, therefore nothing for the snapshot to record, therefore nothing to diff. `0029:A4` admits this in its own "If wrong" clause — "That check cannot discover a vocabulary nobody enumerated — the residual risk this assumption carries" — and then proceeds as though naming the risk discharged it. It does not. 0021 is not a hypothetical future surface that might slip past an attentive team; it is a *named, in-flight, Final* surface that the cluster's own build order schedules for implementation, and the mechanism is blind to it by construction. C-5.

**What 0021 actually says.** Nothing. I read every normative element of 0021 — C1, C2, C3, C4, C5 — plus its Load-Bearing Decisions, its Cross-Cutting Concerns rider, its Trade-offs, its Testing Strategy, its References, and its Decision Rationale. The words `frozen`, `append-only`, and `growing` do not appear in any of them. `0029:C4` is cited nowhere. `0029:C1` is cited nowhere. The only acknowledgement of 0029 anywhere in 0021 is three prose mentions of "RDR 0029's `0.x` promise" governing the `schema` marker's bump (`0021:C2`, `0021:§consequences`, `0021:G-cross-cutting`) — which is the *opposite* of what 0029 delegated. 0029 delegated tier assignment for three specific surfaces. 0021 acknowledged only that its document version bump rides under a promise 0029 owns. The two records agree on the one thing that was never in dispute and are silent on the thing that was. C-1, C-2, C-3.

Worse: `0021:§decision-rationale` records **"Joint-check: clear (12 peers) — open peers 0012–0020, 0022–0024 checked on both arms … share no undecided contract."** 0029 is not in that range and not in that list. 0021's joint-check literally could not have seen 0029 — 0029 did not exist when 0021 ran it, and 0021 never re-ran it. Meanwhile 0029's joint-check fired *at* 0021 "at the lock fence." The edge is declared mutual in the metadata and is, in the actual text, entirely one-directional. C-9.

**The symptom the agent sees.** An agent consuming `intrastate graph --model flow.toml --as=json`:

```json
{"type":"ok","schema_version":"0.3","data":{"schema":"intrastate.graph/1","model":"flow",...}}
```

It does exactly what `0029:C1` and the published register tell it: reads the envelope major (`0`, unstable, but it pins and proceeds), reads `data.schema` (`intrastate.graph/1`, recognised), then looks up the stability tier for the fields it wants to branch on — `rows[].outcome`, `reach.edges[].rule`, `tags[].domain`. It goes to `docs/cli-output-contract.md`, which `0029:§activation-step-1` established as the register of record carrying "the tier word beside every vocabulary C4 tiers." It finds fourteen vocabularies tiered. It finds nothing for the graph document. Not `frozen`. Not `append-only`. Not "untiered by design," which is the honest answer C4 gives for `findings[].class` and `data.dispositions` and would have given here.

The agent has two choices and both are wrong. Treat the absence as "no promise" and it declines to consume a surface that is, in practice, stable — the export exists precisely to be diffed in CI. Treat the absence as "tiered like its neighbours" and it writes an exhaustive switch over `--emit` values or asserts on `rows[]`'s 11 members, and breaks on the release that adds a twelfth column — which `0021:§consequences` explicitly anticipates ("a column added there (the `emit` precedent) has one obvious landing in the schema"). The register 0029 built to end exactly this guessing game has a hole in it shaped like the cluster's other record.

---

## 2. Will any RDR be rewritten within 6 weeks of shipping?

**Yes. 0021, and it will be reopened before it ships, not after.**

Not because its design is wrong — 0021's design is the strongest thing in this cluster. Because it is structurally impossible to implement the set as ordered without editing 0021's Normative Contracts.

The chain is short and has no escape. `0029:C4` is normative and says an unassigned machine-readable surface is a defect. It names 0021's three surfaces. It assigns their tiers to 0021's own Normative Contracts. 0021's Normative Contracts assign nothing. 0021 is Final. Build order is 0029 first. So the moment 0029 lands, 0021 is a Final record carrying a defect *by a peer's normative definition*, and the only place the defect can be repaired is inside 0021's `C2` and `C1` — locked text. C-13.

This is a 7.1 SPEC-DEFECT in the ordinary sense, and the honest disposition is a demotion of 0021 to `Draft [revised from Final … @refine]` with a scoped obligation: add the tier assignments and the two citations (`0029:C1` for the envelope marker, `0029:C4` for the boundary rule) that 0029's settled joint decision already obliges. That is a Refine-stage edit to two contracts, not a redesign. But it is a rewrite of locked normative text, and pretending otherwise is how the cluster gate gets skipped.

Three further forces make 0021's reopening non-optional rather than merely tidy:

1. **`0021:A4` is stale against the cluster.** It verified "a new root export verb requires no amendment to RDR 0005's envelope contract" against the envelope as it stood. 0029 adds a non-`omitempty` field to every terminal record at the gateway. 0021's envelope reasoning needs re-verification against the envelope 0029 leaves behind, not the one A4 examined. C-11.
2. **`0021:S2`'s golden fixtures collide with `0029:S7`'s repo-wide negative**, and 0021's `--as=json` goldens will carry `schema_version` — a field that moves on a schedule 0021 does not control. C-8.
3. **0021's own vocabulary is written in the term 0029 retires.** `0021:C2` describes `dumpColumns` as "a closed 11-member list." `0029:C2` prohibits `closed` as a tier descriptor repo-wide and `0029:S9` enumerates the sites to fix — a census taken against `main`, which does not contain 0021's text. C-15.

**0029 will not be rewritten within six weeks**, and I want to be explicit that this is a finding rather than a softening. Its exposure is deferred by design: the whole apparatus is declarations of intent through `0.x`, and `0029:C2` makes the qualifier ride the emitted major rather than prose, so a wrong tier is correctable for the entire `0.x` series at no consumer cost. The thing that would force a 0029 rewrite inside six weeks is a *fourth* tier proving necessary — and having read C4's fourteen assignments and the two deliberate non-assignments, I can't construct one. What *will* happen to 0029 inside six weeks is a census amendment, which C4 pre-authorized and which is not a rewrite. Its real reckoning is `0029:§activation-step-3`, at 1.0.0, where the census is re-run and a mismatch costs a release slip — and the most likely cause of that mismatch is 0021's untiered surfaces. C-12. That is months out, not weeks.

---

## 3. The one cross-cutting assumption that will not survive first contact

**`0029:A4` — "The three tier names partition every machine-readable surface this CLI emits — no surface needs a fourth tier." Status: Verified.**

The partition claim survives. The **Verified** status does not, and the status is what the downstream machinery consumes.

A4's evidence is a census: "An enumeration of every string vocabulary reachable on the `--as=json` wire found **sixteen** emitted sets." Every word of that is true of `main`. Not one word of it is true of the tree that exists after this cluster implements, because 0021 adds an entire document schema to the `--as=json` wire — 11 top-level members, nested objects for tags/rows/atoms/groups/reach, an `--emit` format set, and a refusal code. A4 was verified against a tree the cluster's own build order guarantees is obsolete. C-4.

The record knows this. A4's own history is the tell: the census "was wrong twice under authoring attention" — seven surfaces missed on the first pass, two more found inside `findings[]` by the grounding lens. A4 draws the right lesson ("because a census that was wrong twice under authoring attention will not stay right under none, completeness is enforced mechanically") and then reaches for a mechanism that cannot cover this case: A7's snapshot reads seams, and an unassigned surface has no seam. The residual is stated and then walked past.

What makes this specifically fatal rather than merely stale is the interaction with C4's delegation. C4 says the census is "complete as of this record's implementation, not complete for all time," and that later surfaces tier in the record that adds them. That framing is sound for *unknown* future surfaces. It is unsound for 0021, for one reason: **0021 is not a future surface. It is a Final, locked, in-flight peer whose text cannot take the assignment without being reopened.** C4 names 0021 explicitly, hands it an obligation, and misstates its status as "early Draft with every assumption still Pending" in the same paragraph. The delegation is made to a record that cannot accept it. C-10.

So A4's partition claim is fine and its Verified status is a claim about coverage that the cluster falsifies on the day it implements.

Two secondary cross-cutting failures ride the same seam:

- **`0021:C1`'s `--emit` set** is an emitted vocabulary by 0029's own definition — it produces a `flag-invalid-value` refusal naming `emit`, which is machine-observable behaviour keyed to set membership. Untiered. C-2.
- **`0021:C2`'s field-order clause** — "a closed 11-member list with `emit` appended last" — states a cardinality *and* a tail position as normative wire facts, which is precisely the assertion class `0029:C2` forbids over `append-only` vocabularies. Whether `rows[]` is `frozen` (in which case the cardinality is legitimate and must be declared) or `append-only` (in which case 0021's own wording is the anti-pattern) is unresolved, because nobody assigned the tier. C-14.

---

## 4. The premortem — written as if it already happened

*Eight months on. Post-incident writeup, `intrastate` 1.0.0 slipped one release; two downstream agent harnesses broken.*

It shipped clean. 0029 landed first, exactly as planned. `clierr.SchemaVersion` went in as one exported constant; `respond.Success` and `clierr.CLIError` both marshalled `"schema_version":"0.1"`; `0023`'s golden at `docs/rdr/0023-resolve-envelope-projection/artifacts/mvv-step1-default-golden.json` was re-captured in the same commit and differed by exactly the one key, just as `0029:S3` required. `respond::Levels()`, `clierr::ExitCodes()`, `resolve::Blocks()`, `cli::UnknownReasons()` all built. The A7 snapshot test went into `internal/graphlint/graphlint_test` and went red on cue when someone moved `CodeVacuousAtom` across the severity partition. `docs/cli-output-contract.md` gained its named section with fourteen vocabularies and a tier beside each, and `llms.txt` pointed at it. The `closed` census in `0029:S9` was executed site by site. It was, by any internal measure, a model implementation. Nobody was careless.

0021 landed six weeks later. `runGraph` mirrored `runLint`'s arm set verbatim. `graphlint.ReachWithEdges` came in beside `Reach` — additive, `reach()` untouched, exactly as `0021:A8` promised. `renderDOT` imported stdlib only. The neutrality oracle held: `intrastate lint` byte-identical with and without the export code. `emitDocument` double-emitted identically under `GODEBUG=randmapiter=1` across all four `--as`×`--emit` cells. Every scenario in `0021:§testing-strategy` passed. `0021:S2`'s goldens for the state-machine and decision-table fixtures were checked in.

**The first break was a CI job.** Four weeks after 0021, someone added a fifth advisory finding code — the event `0029:A1` and `0029:S8` had prepared for, `0006:C17` duly amended to append-only, `TestReq74` replaced with a membership-and-uniqueness assertion. Correct in every respect. `clierr.SchemaVersion` moved `"0.1"` → `"0.2"` under `0029:C1`'s minor-increment rule. And a CI job nobody had thought about went red: `TestGraphExportGolden_StateMachine`, 0021's `--as=json` golden, which carried the envelope and therefore carried `schema_version`. The failure had nothing to do with the graph document. A maintainer re-baselined it in ninety seconds and moved on. Two releases later they re-baselined it again. By the fourth time, re-baselining that golden was muscle memory — which is how, on the fifth, a genuine field rename inside `reach.nodes[].values` got re-baselined too, silently, with no schema bump. `0021:C2`'s additive tripwire had been trained to be ignored. C-8.

**The second break was an agent.** A harness maintainer had written a consumer months earlier against `intrastate graph --as=json`. They had done the work: read `0029:C1`, checked the envelope major, tolerated unknown keys. Then they went looking for the tier on `rows[].outcome` and `reach.edges[].rule`, because the register told them a tier existed for every emitted vocabulary. It wasn't there. They reasoned — correctly, from the published material — that the graph document was `frozen`-adjacent: it carried its own `schema` marker, `0021:C3` promised byte determinism, `0021:S2` pinned goldens, and `0021:§consequences` called it "a first-class, diffable CI artifact." So they wrote a strict decoder over the 11-member `rows[]` list, because `0021:C2` told them in normative text that it was "a closed 11-member list with `emit` appended last."

Then RDR 0010's successor added a twelfth dump column. `0021:§consequences` had said it would: "a column added there (the `emit` precedent) has one obvious landing in the schema." It landed additively within `/1`, exactly per `0021:C2`. `schema` stayed `intrastate.graph/1`. `schema_version` moved a minor. The A7 snapshot showed nothing, because no *seam* changed — the graph document never had one. And the harness's strict decoder hard-failed on an unrecognised member of a list a normative contract had called closed. C-1, C-3, C-14.

The maintainer filed an issue. It is the best-documented thing in the post-mortem, because it is quoted in `0029:§activation-step-1` as the exact trigger for machine-readable tier exposure: *"a consumer asserting cardinality on an `append-only` set."* The register's own escalation condition fired — on a surface the register never covered.

**The third break was quiet and is still shipping.** `graph-export-too-large` entered the CLIError `code` vocabulary. That vocabulary is the one `append-only` row `0029:A5` refuted and `0029:C4` records as `seam: none (prose-only)` — no accessor, no snapshot coverage, no assertion of any kind. So the code appeared with no seam diff, no release-note trigger, and no test failure. A consumer branching exhaustively on `code` took its default branch and reported "unknown internal error" for a model that was simply too large, with the narrow-a-domain remedy sitting right there in the `detail` field, unread. Nobody has noticed. C-6.

**The reckoning came at 1.0.0.** `0029:§activation-step-3` gates promotion to schema `"1.0"` on re-running the C4 census against the then-current tree: *"every emitted machine-readable surface appears in the table."* It didn't. The graph document's field names, the `--emit` set, and `graph-export-too-large` were all emitted and none were in the table. The recorded remedy — *"If the census moves at that point, that is the signal the promise was already drifting — take the 1.0.0 slip rather than promoting a table the tree does not match"* — was followed. 1.0.0 slipped one release while tiers were assigned retroactively to a surface that had been in consumers' hands for eight months, which meant the assignment was no longer a free declaration of intent. It was a ratification of whatever consumers had already guessed. C-12, C-5.

**The root cause, one sentence:** `0029:C4` delegated three tier assignments to a record it described as "an early Draft with every assumption still Pending," and that record had been Final for two weeks — so the obligation was created in a record that could not enforce it and delegated to a record that could not accept it, and the mechanism built to catch exactly this reads seams that an unassigned surface never has.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

**AT-1 — The census covers the cluster, not `main`** *(catches C-1, C-4, C-12)*

```gherkin
Scenario: Every machine-readable surface in the cluster has a tier
  Given RDR 0029:C4 declares "an unassigned machine-readable surface is a defect"
  And RDR 0021 is in the same reconcile cluster and introduces surfaces under --as=json
  When the C4 census is evaluated against the union of main and every Final
       unimplemented record in the cluster
  Then every field name, format token, and refusal code that 0021:C2, 0021:C1,
       and 0021:C5 place on the wire appears in a tier table
  And the table is reachable from the record that emits the surface
  But the census evaluated only against main omits all of them
  Therefore 0029:A4 MUST NOT be marked Verified
```
Plain: re-run A4's enumeration with 0021's C2 field list appended. It returns "unassigned" on every row. A4 flips to Pending-on-0021.

**AT-2 — A refusing flag value is a vocabulary** *(catches C-2)*

```gherkin
Scenario: --emit's value set is tiered
  Given 0021:C1 refuses an unrecognised --emit value with flag-invalid-value naming emit
  Then the accepted set {json, dot} is a machine-readable vocabulary under 0029:C2
  And it carries exactly one of frozen | append-only | growing
  And that tier is stated in 0021's Normative Contracts per 0029:C4
  When 0021:C1 is read for any of those three tokens
  Then none is found — the test fails at review
```

**AT-3 — The enforcement mechanism can see the thing it enforces** *(catches C-5)*

```gherkin
Scenario: A7's snapshot detects an unassigned surface
  Given 0029:C3 backs the tier system with a committed seam/severity snapshot
  And 0029:C4 declares an unassigned machine-readable surface a defect
  When a surface ships with no tier and therefore no enumeration seam
  Then the snapshot records no row for it and the diff is empty
  And CI is green
  Therefore C3's mechanism cannot detect C4's defect class
  And C4's delegation to a peer record MUST carry a tracked obligation
       (an amendment, a demotion, or a JDR hoist) rather than prose
```

**AT-4 — The growing vocabulary is the unpoliceable one** *(catches C-6)*

```gherkin
Scenario: 0021's new refusal code lands in the seamless vocabulary
  Given 0029:A5 refuted the CLIError code registry seam
  And 0029:C4 records that row as "seam: none (prose-only)"
  And 0021:C4 introduces graph-export-too-large into that same vocabulary
  Then the one vocabulary the cluster provably grows is the one the snapshot
       cannot police
  And 0029:C3's coverage statement MUST name this cluster-specific gap,
       not only the general prose-only limit
```

**AT-5 — The two-marker read order is specified** *(catches C-7)*

```gherkin
Scenario Outline: An agent reads both version markers
  Given the envelope carries schema_version "<env>" per 0029:C1
  And data carries schema "<doc>" per 0021:C2
  When the agent decides whether it can parse
  Then a normative passage states which marker is checked first
  And states the behaviour for <outcome>

  Examples:
    | env | doc                  | outcome                          |
    | 1.0 | intrastate.graph/1   | parse                            |
    | 1.0 | intrastate.graph/2   | UNSPECIFIED — no passage governs |
    | 2.0 | intrastate.graph/1   | UNSPECIFIED — no passage governs |
```
Two of three rows resolve to nothing in either record. The test fails at review and forces one sentence into `0021:C5` or `0029:C1`.

**AT-6 — Goldens do not cross the envelope boundary** *(catches C-8)*

```gherkin
Scenario: 0021's golden fixtures survive a schema_version bump
  Given 0029:S7 states no golden/snapshot test of the JSON envelope exists
        and this RDR does not add one
  And 0021:S2 pins golden fixtures for intrastate.graph/1
  When a 0021 golden is captured under --as=json
  Then it contains the envelope, and therefore schema_version
  And it breaks on the next minor bump for reasons unrelated to the graph document
  Therefore 0021:S2 MUST scope its goldens to the DOCUMENT (--as=text, or jq .data),
       never the envelope
```

**AT-7 — The joint-check edge is bidirectional in the text** *(catches C-9, C-10)*

```gherkin
Scenario: A fired joint decision is visible from both ends
  Given 0029's joint-check fired at 0021 and settled C4 as the home
  When 0021:§decision-rationale is read
  Then it records "Joint-check: clear (12 peers) — open peers 0012–0020, 0022–0024"
  And 0029 appears in neither the range nor the peer-citation paragraph
  And 0029:C4 describes 0021 as "an early Draft with every assumption still Pending"
  While 0021:§metadata reads Status: Final, dated 2026-08-28, A1–A8 Verified
  Therefore both ends of the edge are wrong about the other and the gate MUST
       re-run 0021's joint-check against 0029 before either implements
```

**AT-8 — Envelope assumptions are re-verified against the cluster** *(catches C-11)*

```gherkin
Scenario: 0021:A4 holds after 0029 lands
  Given 0021:A4 verifies a new root verb needs no envelope-contract amendment
  And 0029:C1 adds a non-omitempty schema_version to every terminal record
  When the graph verb's --as=json envelope is emitted post-0029
  Then it carries a field no passage in 0021 mentions
  And any 0021 assertion over an exact envelope key set fails
  Therefore 0021:A4 MUST be re-verified against the post-0029 envelope
```

**AT-9 — No record states an assertion its peer forbids** *(catches C-14)*

```gherkin
Scenario: rows[] cardinality is consistent across the cluster
  Given 0029:C2 forbids consumers asserting cardinality, ordinal, or tail position
        on an append-only vocabulary
  And 0021:C2 states rows[] is "a closed 11-member list with emit appended last"
  When rows[]'s tier is looked up
  Then no tier is assigned
  And the cluster cannot say whether 0021:C2's wording is a legitimate frozen
       declaration or the exact anti-pattern 0029:C2 retires
  Therefore the tier MUST be assigned before either record implements
```

**AT-10 — The retired word is retired across the cluster** *(catches C-15)*

```gherkin
Scenario: "closed" does not describe a tiered vocabulary anywhere in the cluster
  Given 0029:C2 prohibits the descriptive use of "closed" repo-wide
  And 0029:S9 enumerates the in-scope sites from a census taken against main
  When 0021's normative text is added to that census
  Then 0021:C2's "a closed 11-member list" describing dumpColumns is in scope
  And 0029:S9's site list does not include it
  Therefore S9's census MUST extend to cluster records, or 0021's wording MUST
       be changed in the same reconcile
```
