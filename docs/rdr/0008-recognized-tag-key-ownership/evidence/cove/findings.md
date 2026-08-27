Model: claude-opus-5[1m]

# CoVe pre-lock lens — RDR 0008 (recognized-tag-key ownership)

Grounded on `main` at `03e6629`. Every codebase claim below was checked by
reading the named source file, not by trusting the RDR's prose.

## Step 0 — Grounding sweep

### Codebase claims (`path::Symbol`)

| # | RDR claim | Verdict | Cite / what I found |
| --- | --- | --- | --- |
| C-1 | `internal/resolve/resolve.go::recognizedTagKey` is the unexported constant `"recognized"` | CONFIRMED | `internal/resolve/resolve.go:110` — `const recognizedTagKey = "recognized"` |
| C-2 | `internal/resolve/resolve.go::assemble` injects it with `ProvenanceRecognized` only when `in.Recognized` is non-empty | CONFIRMED | `resolve.go:151-156` — `if in.Recognized != "" { view.tags[recognizedTagKey] = taggedValue{value: in.Recognized, provenance: ProvenanceRecognized} }` |
| C-3 | Outcome gate is `row.Outcome != in.Recognized` plus the alphabet check; never reads the key | CONFIRMED | `resolve.go:321` (`!in.Table.models(in.Recognized)`), `:332`, `:479` |
| C-4 (A2) | Exactly four non-test `Input.Recognized` reads, exactly one reaching tag vocabulary | CONFIRMED with a correction | `grep -n '\.Recognized' internal/resolve/resolve.go` -> lines 151, 153, 321, 332, 479, 518. That is **five** distinct read statements (151/153 are one site). The RDR counts "four non-test sites" then parenthesizes `::refuse` as "a fourth" — the arithmetic is internally muddled (F-6) but the substance (one tag-vocabulary site) is confirmed. |
| C-5 (A2) | `::Resolve` assembles once; `view := assemble(in)` is the only `assemble` call in the package | CONFIRMED | `resolve.go:319` is the sole call; `grep -rn 'assemble(' internal/` returns only `:148` (definition), `:319`, and a prose comment in `adversarial_test.go:398` |
| C-6 (A2) | The one view threads to matcher `view.matches(row.Match)` and guards `gate(candidates, in.Guards, view)` -> `::evaluateGuard` -> `seam.Evaluate(guard, view)`; `::escapeOrRefuse` reuses it | CONFIRMED | `resolve.go:335`, `:341`, `:376`, `:380`, `:425-432`, `:473-485` |
| C-7 (A2) | `TagSet` wraps a map so pass-by-value shares one backing map | CONFIRMED | `resolve.go:98-100` — `type TagSet struct { tags map[string]taggedValue }` |
| C-8 (A2 test gap) | `internal/resolve/resolve_test.go::TestReq17_OwnedObservedAndRecognizedTagsAllReachSelection` pins the recognized tag reaching *match* selection | CONFIRMED | `resolve_test.go:615` (func), subtest at `:633` asserting `p.RuleID == "rdr.draft.on-recognized-tag"` |
| C-9 (A2 test gap) | `internal/resolve/fixtures_test.go::fixtureGuards.Evaluate` has signature `Evaluate(guard string, _ resolve.TagSet)` — the suite's *only* guard seam, so it discards the view | CONFIRMED | `fixtures_test.go:20` — `func (g fixtureGuards) Evaluate(guard string, _ resolve.TagSet) resolve.GuardResult`; `rg 'func .*Evaluate\(' internal/` returns exactly one hit repo-wide |
| C-10 (A4) | In Go table data the literal key `recognized` appears exactly once, as a read-only **match pattern** in `internal/resolve/fixtures_test.go::recognizedTagSensitiveTable` | CONFIRMED | `fixtures_test.go:237` — `{Key: "recognized", Value: "successful"}` inside `Match`; `RequiresOwned: []string{"status"}` (`:239`); `NextTags`/`Writes` both `{Key: "status", ...}` (`:240-241`). No owned/observed/write occupancy. |
| C-11 (A4) | The three committed TOML fixtures use `recognized` only as a provenance value under an author-named key: `[tags.outcome]` in the two 0002 spikes, `[tags.rewind_target]` in the 0003 spike | CONFIRMED | `docs/rdr/0002-.../evidence/spikes/rdr-fixture.toml:25-26` and `kata-fixture.toml:22-23`; `docs/rdr/0003-.../evidence/spikes/guard-fixture.toml:36-37` |
| C-12 (A5) | The `Observed` then `Owned` loops in `::assemble` implement last-writer-wins over the recognized entry | CONFIRMED | `resolve.go:157-162` — unconditional map assignment, recognized written first at `:152` |
| C-13 (A6) | `internal/resolve/resolve.go::Resolve` has signature `(Result, error)` and its doc comment says verbatim "the error return is reserved for programmer mistakes, not for modeled refusals" | CONFIRMED | `resolve.go:318` signature; `resolve.go:295-296` doc comment, exact wording matches |
| C-14 (A6) | `Resolve`'s body has **no** non-nil error path today | CONFIRMED | Every `return` in `Resolve` (`resolve.go:322, 344, 348, 350, 352`) ends `, nil`; `grep -n 'errors\.\|fmt.Errorf' internal/resolve/resolve.go` -> zero hits; the file imports only `slices` and `strings` (`:13-16`) |
| C-15 (D3) | Precedence `owned > observed > recognized` in `assemble` | CONFIRMED | `resolve.go:142-147` doc comment plus the write order at `:151-162` |
| C-16 (Sibling-path check) | "`recognizedTagKey` is the only reserved key in the kernel (whole-package sweep for key constants and 'reserved' yields nothing else)" | CONFIRMED for the kernel package; the sweep is package-scoped and misses a repo sibling | Within `internal/resolve` the only string-key constant is `recognizedTagKey:110`; "reserved" appears only at `:295`. **However** `internal/cli/respond/respond.go:24` states "The terminal \"ok\"/\"failed\" type names are **reserved**." — a shipped in-repo reserved-token decision the check never surfaces (F-2). |
| C-17 (Prior art) | XState `repos/xstate/packages/core/src/guards.ts::GuardArgs`, guards receive `{ context, event }` | CONFIRMED as content, REFUTED as path | Symbol exists at `../state-machines/repos/xstate/packages/core/src/guards.ts:49` (`export interface GuardArgs<`) and `:132` (`{ context, event }: GuardArgs<any, any>`). `repos/xstate/...` does **not** resolve from the intrastate repo root (`ls repos` -> no such dir); the research cache writes `../state-machines/repos/xstate/...`. See F-4. |
| C-18 (Prior art) | Inngest "Read-Only Fields ... system-managed and cannot be modified through write operations" at `repos/inngest/pkg/api/v2/README.md` | CONFIRMED as content, REFUTED as path | `../state-machines/repos/inngest/pkg/api/v2/README.md:255` — `## Read-Only Fields`. Same missing `../state-machines/` prefix. See F-4. |

### Peer-RDR and fixture claims

| # | RDR claim | Verdict | Cite |
| --- | --- | --- | --- |
| P-1 (A6) | "The string `Input` appears **zero** times in RDR 0009 (case-sensitive, whole file)" | CONFIRMED | `grep -c "Input" docs/rdr/0009-escape-row-shape-conformance-ownership.md` -> `0` |
| P-2 (A6) | RDR 0009 scopes the obligation to "a PRODUCER obligation on every constructor of resolve.Row values" | CONFIRMED verbatim | `0009...md:476-477` |
| P-3 (A6) | 0009's predicate is "exactly `len(row.Escape) != 0 && len(row.Writes) != 0`" applied to every row "at `Resolve` entry" | CONFIRMED verbatim | `0009...md:531`; `:532-533` — "applied to every row of the supplied table at `Resolve` / entry" |
| P-4 (A6) | 0009's rationale: its breach is a property of the table value "independent of the input tuple" | CONFIRMED verbatim | `0009...md:503` |
| P-5 (A6) | 0009's predicate form is "sharpened at Pre-Lock" | CONFIRMED verbatim | `0009...md:513`; also `:543` |
| P-6 (A6) | 0009 is `Draft` with "all seven of its assumptions `Pending`" | CONFIRMED | `0009...md:9`; seven `**Status**: Pending` lines at `:305, 319, 333, 355, 376, 387, 402` |
| P-7 (A6) | "**0009's own triage** records 'Disjoint answer spaces, no ordering dependency' against this RDR's kata" | REFUTED as to locus | Sentence exists at `0009...md:169` but in 0009's **Background prose**. `ls -R docs/rdr/0009-escape-row-shape-conformance-ownership/` shows only `evidence/propose-premortem/critic.md` and `evidence/research/propose-prior-art.md` — **no `triage.md` exists** (contrast `docs/rdr/0001-resolution-kernel/artifacts/triage.md`). See F-3. |
| P-8 (A1) | RDR 0002: "Validation failures MUST retain stable data-level categories before CLI mapping, **including at minimum** malformed TOML, ... and ambiguous overlap." | CONFIRMED verbatim | `0002...md:342-348` (normative block) |
| P-9 (A1) | RDR 0002 Technical Design calls the outcome alphabet "the **closed set**" | CONFIRMED | `0002...md:217` |
| P-10 (A1) | RDR 0001 Normative Contracts: the refusal-kind set "is **exactly**" | CONFIRMED | `0001...md:214`; also `:106` |
| P-11 (A1) | RDR 0002 Status is `Final` | CONFIRMED | `0002...md:9` |
| P-12 (A1) | "RDR 0005 never enumerates RDR 0002's data-level categories at all (**zero** occurrences)" | CONFIRMED | `grep -c "data-level" docs/rdr/0005-skill-integration-cli-contract.md` -> `0` |
| P-13 (A1) | RDR 0005's nearest table is a "**Minimum** stable code strings" list of CLI codes | CONFIRMED | `0005...md:599`, table at `:601-614` |
| P-14 (A1) | RDR 0005: "Error envelopes stay append-only through `CLIError` fields" | CONFIRMED verbatim | `0005...md:313` |
| P-15 (A1d) | "both canonical fixtures name the recognized-provenance tag `outcome`" | CONFIRMED | `rdr-fixture.toml:25` and `kata-fixture.toml:22`, both `[tags.outcome]` |
| P-16 (A1d) | "`recognized` appears in RDR 0002 exclusively as a provenance value, never as a tag key" | CONFIRMED | 10 hits in `0002...md` (`:45, 108, 195, 215, 333, 344, 453, 490, 657, 717`); each a provenance value, alphabet reference, lint-category name, search-command line, or prose. Zero uses as a tag key. |
| P-17 | "no tag key named `outcome` appears in the RDR 0002 body" | CONFIRMED, with a caveat | True of the body; `[tags.outcome]` appears in 0002's own canonical evidence spikes. Slightly self-serving framing. |
| P-18 (A5 P1) | RDR 0002: "The model MUST declare every tag it matches or writes, including each tag's provenance: owned, observed, or recognized" | CONFIRMED verbatim | `0002...md:332-334` (normative block) |
| P-19 (A5 P2) | RDR 0002 rejects "rules that match on missing tag declarations" and "writes to non-owned tags" | CONFIRMED verbatim | `0002...md:233-234` |
| P-20 (A5 P2) | RDR 0002 Failure Modes: "a lint failure before the model is accepted" | CONFIRMED verbatim | `0002...md:636-637` |
| P-21 (A5 P2) | The categories are named `unknown tag` and `write to non-owned tag` | CONFIRMED | `0002...md:346` normative list |
| P-22 | RDR 0002 reserves the `<clear>` write sentinel | CONFIRMED | `0002...md:337-339` normative block |
| P-23 (D2) | RDR 0001 D2 records the key as implementation latitude, "adds no public surface" | CONFIRMED | `docs/rdr/0001-resolution-kernel/artifacts/deviations.md:69-86`; `:82`; `:84-86` |
| P-24 (D3) | RDR 0001 D3 owns precedence `owned > observed > recognized` | CONFIRMED | `deviations.md:88-105`; `:100` |
| P-25 | RDR 0001 `verification.md`: "deterministic and refusal-shaped, not a breach" | CONFIRMED verbatim | `verification.md:526-528` |
| P-26 | REQ-17 is "Merge owned, observed, and freshly recognized tags into the evaluation view" | CONFIRMED verbatim | `docs/rdr/0001-resolution-kernel/artifacts/req-list.md:60` |
| P-27 (Alt 2 con) | "RDR 0007's chosen guard domain evaluates predicates over the assembled view, so guards could no longer reference the recognized value" | CONFIRMED, and stronger than stated | `0007...md:1258-1268` normative block — "PRESENCE IS PROVENANCE-BLIND ... `exists` MUST decide TRUE for an observed- or recognized-supplied key exactly as for an owned one." RDR 0007 is `Final` (`:9`), so withdrawal would breach a locked normative clause. |
| P-28 | "Joint-check: clear (7 peers)" | REFUTED / unsupported | RDR 0007 (`Final`) at `:953-956` **quotes RDR 0008's normative sentence** and leans A13's Verified conclusion on it. RDR 0008 records no inbound dependency. See F-1. |
| P-29 | Sibling-path check enumerates `recognizedTagKey` and 0002's `<clear>` as the only adjacent reserved-token decisions | REFUTED (incomplete) | `grep -c "0006" docs/rdr/0008-...md` -> `0`. RDR **0006** is `Final` (`0006...md:9`) and normatively owns blocking static graph-lint authority over the normalized model, with declared-tags-with-provenance in its minimum input contract (`:251`). Never reconciled. See F-5. |

**Step 0 tally: 3 REFUTED (P-7 locus, P-28 joint-check, P-29 sibling sweep), 2 partially REFUTED on path/arithmetic (C-17/C-18 paths, C-4 count), 0 NOT-FOUND.** Every `path::Symbol` in the RDR resolves on `main`.

## Step 1-2 — Verification questions

Each answer written independently against source or the named peer file.

**Q1 (source). Does `internal/resolve/resolve.go::Resolve` already have a non-nil error path that an `Input` precondition could join, or would candidate (b) introduce the first one?**
It would introduce the first one. Every `return` in `Resolve` (`resolve.go:322, 344, 348, 350, 352`) ends `, nil`. The file imports only `slices` and `strings` (`:13-16`) — no `errors`, no `fmt` — and `grep -n 'errors\.\|fmt.Errorf' internal/resolve/resolve.go` returns zero hits. The `error` in the signature at `:318` is dead surface today. The RDR's phrasing "no non-nil error path in its body today" is exactly right.

**Q2 (source). Is the guard seam's view parameter used by any implementation in the repo, or is "only guard seam discards the view" understated?**
There is exactly one `Evaluate` implementation repo-wide — `rg -n 'func .*Evaluate\(' internal/` returns only `internal/resolve/fixtures_test.go:20`, whose signature is `Evaluate(guard string, _ resolve.TagSet) resolve.GuardResult` (blank-named, hence unusable). The interface `GuardEvaluator.Evaluate(guard string, view TagSet)` at `resolve.go:91-94` is never exercised with a view-reading implementation anywhere — production or test. The RDR's claim is accurate and if anything conservative.

**Q3 (source). Does `::assemble` distinguish an absent recognized outcome from an empty-string one, and does anything downstream depend on the key's presence?**
`resolve.go:151` gates injection on `in.Recognized != ""`, so an empty recognized outcome yields a view with **no** `recognized` key. Downstream `resolve.go:321` calls `in.Table.models("")` = `slices.Contains(t.Outcomes, "")` — false for any normal alphabet — so `Resolve` refuses `KindUnmodeledOutcome` before selection. The absent-key state is unreachable past the alphabet gate. RDR 0008's normative block 1 states the binding unconditionally; the shipped code is conditional. See F-9.

**Q4 (source). Could a row's `RequiresOwned` name `recognized` and produce a confusing `owned_state_unavailable`?**
Yes, and the RDR is silent. `resolve.go:444-458` iterates `row.RequiresOwned` and tests `view.has(key, ProvenanceOwned)` — a **provenance-specific** test (`resolve.go:125-128`). `RequiresOwned: []string{"recognized"}` finds the key present under `ProvenanceRecognized`, so `has` is false and the kernel refuses `owned_state_unavailable` naming `recognized`. RDR 0008 constrains declarations and `Input` producers, never `Row.RequiresOwned`. See F-8.

**Q5 (source). Is there a sibling reserved-token decision in shipped Go code the sibling-path check missed?**
Yes: `internal/cli/respond/respond.go:22-25` — "The terminal \"ok\"/\"failed\" type names are **reserved**." A shipped, carrier-owned reserved-name decision on a wire-format discriminator with no enforcement (`OK` at `:112` unconditionally sets `s.Type = "ok"`; nothing validates a future `Stream` emitter). RDR 0008's sweep was scoped to `internal/resolve`, so it could not surface this. See F-2.

**Q6 (source). Does `recognizedTagSensitiveTable` place the literal `recognized` in any forbidden position?**
No. `fixtures_test.go:226-245`: the literal appears only at `:237` inside `Match`. `RequiresOwned` is `[]string{"status"}` (`:239`); `NextTags` and `Writes` are both `{Key: "status", ...}` (`:240-241`). A4 is accurate on this fixture.

**Q7 (peer RDR). Does RDR 0009 anywhere reach the `Input` value or the tag channel?**
No. `grep -c "Input" docs/rdr/0009-...md` -> `0` (case-sensitive, whole file). Its normative surface is `resolve.Row` construction (`:476-477`), a whole-table entry precondition (`:492-503`), and an exported table predicate (`:509-517`). Its rationale grounds the breach as "independent of the input tuple" (`:503`) — the structurally opposite trigger. A6's negative conclusion holds.

**Q8 (peer RDR). Is RDR 0008 a dependency of any *Final* peer?**
Yes. RDR 0007 is `Final` (`0007...md:9`) and its A13 Evidence quotes RDR 0008's normative producer sentence verbatim (`0007...md:953-956`) as the load-bearing premise that "The single input-side producer obligation in the whole RDR set constrains one reserved key *name* ... proof the authors knew how to write an input-boundary obligation and wrote only that one." A `Final` peer has consumed a `Draft` RDR's not-yet-locked text. See F-1.

**Q9 (RDR internal). Does the RDR state a spelling rule that survives RDR 0002's actual `[tags.<tag>]` TOML schema?**
Partially. `0002...md:200-207` fixes the layout as `[tags.<tag>]`. RDR 0008's normative block 2 says the comparison is "byte-exact on the parsed TOML key value: case-sensitive, no trimming, no folding". Testing Strategy scenario 4 (`0008...md:1044-1052`) then resolves quoting and whitespace in opposite directions inside one Expected clause, and contradicts Load-Bearing Decisions / Identity (`:645-652`). See F-7.

**Q10 (RDR internal). What happens when a table declares a recognized-provenance tag named `recognized` but the flow binds no recognized outcome?**
RDR is silent. Normative block 2 sets naming and an upper cardinality bound ("at most one such declaration") but no lower bound and no reachability obligation. Failure Modes (`:928-935`) covers only the inverse case (an observed tag intended as recognizer semantics). Combined with Q3, a declared-but-never-bound reserved key is authorable and undiagnosed.

**Q11 (RDR internal). Does the RDR reconcile its enforcement locus with RDR 0006, the Final holder of blocking-lint authority?**
No — RDR 0006 is not mentioned once (`grep -c "0006"` -> `0`). RDR 0006 (`Final`, `:9`) normatively holds blocking acceptance-gate authority over the normalized model, requires blocking findings to carry "a stable code, model identity, severity, human-readable message, and the source rule/context id or source span", makes `intrastate lint` the authoritative surface, and states "resolver-local validation flags MAY call that engine, but MUST NOT define different acceptance rules". Its minimum lint input contract includes "declared tags with provenance" (`:251`). See F-5.

**Q12 (RDR internal + source). Is A2's `Input.Recognized` sweep arithmetic self-consistent?**
No. A2 says "four non-test sites", names one + three = four, then adds "(A fourth, `::refuse` ...)" — making five while calling it a fourth. Source shows five distinct read statements: `resolve.go:151/153`, `:321`, `:332`, `:479`, `:518`. See F-6.

## Step 3 — Findings

### F-1 — (c) internal contradiction / undisclosed inbound dependency — "Joint-check: clear (7 peers)" is false: a *Final* peer already depends on this Draft's normative text

**RDR passage** (`0008...md:758`): `Joint-check: clear (7 peers)`.

**Source cite** (`docs/rdr/0007-guard-predicate-totality.md:9, 949-956`): RDR 0007 `Status: Final`, and its A13 Evidence reads:

> The single input-side producer obligation in the whole RDR set constrains one reserved key *name*, not provenance origin: RDR 0008's "Producers of kernel `Input` MUST NOT supply an owned or observed tag keyed `recognized`" — proof the authors knew how to write an input-boundary obligation and wrote only that one.

**Why it is a defect.** RDR 0008 asserts no live joint dependency while a locked, `Final` peer has already consumed and quoted verbatim RDR 0008's normative block 4 as the load-bearing premise of a `Verified` assumption. This inverts the dependency direction the RDR documents: Predecessors lists only 0001 and 0002, and Overrides says it "Narrows nothing in either peer" — but it now has a downstream consumer that breaks if block 4 is reworded, retargeted (the Enforcement-locus decision is explicitly still **open**, A6 Residual), or if candidate (a) is chosen and the obligation becomes unchecked prose. Worse, RDR 0007's A13 reasoning is negative-existential ("wrote only that one"), so it is falsified not by 0008 changing wording but by any peer adding a second input-side obligation. RDR 0008 cannot lock a "clear" joint-check while a Final peer's Verified assumption rests on the exact sentence whose enforcement locus 0008 leaves unsettled. The joint-check line must record the 0007->0008 inbound dependency, and the RDR must state that block 4's text is now stability-bearing for RDR 0007 A13.

### F-2 — (b) a new rule with an existing sibling the RDR's own sibling-path check missed

**RDR passage** (`0008...md:741-752`): "Sibling-path check: the chosen approach adds one identity rule (a reserved tag key). Searched for an adjacent path already making that decision: `internal/resolve/resolve.go::recognizedTagKey` is the only reserved key in the kernel (**whole-package sweep** for key constants and \"reserved\" yields nothing else), and RDR 0002's `<clear>` sentinel is the nearest sibling reserved-token decision ... No parallel signal is invented."

**Source cite** (`internal/cli/respond/respond.go:22-25`):

```
// As the CLI grows, add an intermediate-record emitter (Stream) for
// verbs that produce zero-or-more records before the terminal line. The
// terminal "ok"/"failed" type names are reserved.
```

**Why it is a defect.** The sweep was scoped to `internal/resolve` ("whole-**package** sweep"), so it structurally could not find a sibling outside the kernel — yet the RDR reports the result as if repo-wide ("the only reserved key", "No parallel signal is invented"). A shipped in-repo sibling exists: `respond` reserves the wire-format discriminator values `"ok"`/`"failed"` against future emitters. That sibling is directly informative for the one decision RDR 0008 leaves open, because it is the **candidate (a)** shape — a documented, unenforced reservation with zero detection (`respond.OK` at `:112` unconditionally sets `s.Type = "ok"`; nothing validates that a future `Stream` emitter avoids the reserved names). A reviewer weighing (a) vs (b) should have that precedent in front of them. Under the RDR's own rules an adjacent path already making this kind of decision must be named; this one is not, and the sweep's stated scope is narrower than the conclusion it licenses.

### F-3 — (a) peer claim REFUTED: the cited RDR 0009 "triage" artifact does not exist

**RDR passage** (`0008...md:472-474`): "No ordering dependency is created — **0009's own triage** records \"Disjoint answer spaces, no ordering dependency\" against this RDR's kata."

**What I found.** The sentence exists, but in RDR 0009's **Background prose**: `docs/rdr/0009-escape-row-shape-conformance-ownership.md:166-171` — "Triage also examined and rejected merging this with kata `z53t` (RDR 0008) ... Disjoint answer spaces, no ordering dependency." And `ls -R docs/rdr/0009-escape-row-shape-conformance-ownership/` shows only `evidence/propose-premortem/critic.md` and `evidence/research/propose-prior-art.md` — **there is no `triage.md`** (contrast `docs/rdr/0001-resolution-kernel/artifacts/triage.md`, which does exist).

**Why it is a defect.** A6 is `Verified` by `Method: Peer RDR` and this is one of its supporting cites. The substance survives (0009 does say it), but the evidence points at a non-existent artifact, so the citation cannot be followed and the load-bearing "no ordering dependency" conclusion rests on a Draft peer's self-assessment of a sibling — not on an independent triage record, which is what "0009's own triage" implies. Given A6 already carries a "Dependency qualification (RDR 0009 is `Draft`, not Final)" caveat, mis-attributing the source of the no-ordering claim inflates its authority precisely where the RDR is trying to be careful. The cite should read `0009-escape-row-shape-conformance-ownership.md` Background, marked as 0009's self-assessment.

### F-4 — (a) codebase claims with unresolvable paths: both prior-art `path::Symbol` cites are wrong relative to the repo root

**RDR passages** (`0008...md:232-235`): "XState: guards receive `{ context, event }` (`repos/xstate/packages/core/src/guards.ts::GuardArgs`) ... Inngest: \"Read-Only Fields ... system-managed and cannot be modified through write operations\" (`repos/inngest/pkg/api/v2/README.md`)." Repeated in References (`:1197-1200`).

**What I found.** The content is accurate — `../state-machines/repos/xstate/packages/core/src/guards.ts:49` declares `export interface GuardArgs<`, `:132` destructures `{ context, event }: GuardArgs<any, any>`; `../state-machines/repos/inngest/pkg/api/v2/README.md:255` is `## Read-Only Fields`. But `repos/` does not exist under the intrastate repo root (`ls repos _repos` -> nothing), so neither path resolves from the repo the RDR lives in. The RDR's own research cache writes them correctly: `docs/rdr/0008-.../evidence/research/prior-art.md:44` — `../state-machines/repos/xstate/packages/core/src/guards.ts::GuardArgs` — and `:48` for Inngest. The RDR body stripped the `../state-machines/` prefix.

**Why it is a defect.** The Finalization Gate's Assumption Verification requires "that each cited `path::Symbol` **resolves on `main`**". These two do not. They are also the sole external prior-art anchors for the Decision Rationale's deciding "prior-art alignment" row, so a reviewer who cannot follow them cannot audit the row that (with blast radius) decides the fork. Mechanical fix — restore the `../state-machines/` prefix used by the cache, or mark the citations external-checkout-relative — but it must be fixed before lock, since the Gate checks exactly this.

### F-5 — (d) RDR is silent, revealing a missing requirement: RDR 0006 holds Final blocking-lint authority and is never mentioned

**RDR passage.** `grep -c "0006" docs/rdr/0008-recognized-tag-key-ownership.md` -> `0`. Phase 2 (`:976-981`) places the check thus: "Add the reserved-key checks to RDR 0002's load/lint path", and Prerequisites (`:951-953`) gate on "RDR 0002 implementation underway (the enforcement point is its normalizer's load/lint path)". References (`:1182-1200`) list 0001, 0002, 0007, 0009 — not 0006.

**Source cite** (`docs/rdr/0006-graph-lint-authority-and-guarantees.md:9` and its normative blocks): `Status: Final`; "Graph lint MUST be a blocking acceptance gate over the normalized transition model. A model with any blocking lint finding MUST NOT be accepted for resolver use or CI success."; "Graph lint MUST check **at least** these blocking invariant classes: dangling edge, dead end, determinism/overlap, guard exhaustiveness/gap, single-valued state, owned-set-before-match, and declared terminal/escape handling."; "Every blocking finding MUST carry a stable code, model identity, severity, human-readable message, and the source rule/context id or source span when the normalized model can provide one."; "The authoritative CLI surface for graph acceptance MUST be the root command `intrastate lint` or a same-engine CI invocation of that command. Pre-commit hooks, aliases, and **resolver-local validation flags MAY call that engine, but MUST NOT define different acceptance rules.**" Its minimum lint input contract explicitly includes "declared tags with provenance, value kind, finite-domain metadata ..." (`0006...md:251`).

**Why it is a defect.** RDR 0008 introduces a new blocking acceptance rule over declared tags and homes it in 0002's normalizer without ever addressing the Final RDR that claims blocking-acceptance authority and whose declared input contract includes declared tags with provenance. Three concrete gaps follow. (i) It is unstated whether a `reserved tag key` failure is a *0002 validation category* (pre-normalization, as written) or a *0006 blocking lint finding* — 0006 requires a strictly richer payload (stable code, model identity, severity, message, source rule id/span) than RDR 0008's normative block 3 mandates (offending name, required name, rule). (ii) RDR 0008's Risks section anticipates only "RDR 0005's exit-code map, refusal renderers, remediation docs" as category-enumerating consumers (`:890-893`) and misses 0006's `graph-lint-failed` aggregate envelope entirely. (iii) 0006's "MUST NOT define different acceptance rules" creates a real risk that a 0002-local reserved-key check *is* a second acceptance rule. RDR 0008's Cross-Cutting Concerns and Proportionality gate items are unanswered templates, so nothing else in the document would catch this. Either RDR 0006 belongs in Predecessors/References with an explicit locus reconciliation, or the RDR must state why declaration-name validation is pre-normalization and therefore outside 0006's normalized-graph scope.

### F-6 — (c) internal contradiction: A2's site count is self-inconsistent and off by one against source

**RDR passage** (`0008...md:296-306`): "Package-wide sweep of every `Input.Recognized` read in `internal/resolve/resolve.go` finds **four** non-test sites, exactly one of which reaches tag vocabulary: `::assemble` ... is the sole injection point. **The other three** compare the raw string and never route through `TagSet` — `::Resolve` `in.Table.models(in.Recognized)`, `::Resolve` and `::escapeOrRefuse` `row.Outcome != in.Recognized` ... (**A fourth**, `::refuse` echoing `Recognized` onto the `Refusal`, is a diagnosis field ...)"

**Source cite** (`grep -n '\.Recognized' internal/resolve/resolve.go`): lines `151`, `153` (assemble, one site), `321` (`models`), `332` (Resolve gate), `479` (escapeOrRefuse gate), `518` (`refuse` echo) — five distinct read sites.

**Why it is a defect.** The passage says "four", enumerates one + three = four, then introduces "a fourth" which is actually the fifth — so the stated total is wrong and the ordinal contradicts the count in the same sentence. A2 is `Verified` by `Method: Source Search`, and the only thing a Source Search assumption asserts is that the sweep was exhaustive. A sweep whose own arithmetic does not close cannot be relied on to be exhaustive; a reader re-running it gets five and cannot tell whether the RDR missed one or miscounted. The substantive conclusion (exactly one tag-vocabulary site) is correct and independently confirmed at `resolve.go:151-156`, so this is a precision defect in the evidence record, not a false conclusion — but it is exactly what the Gate's "no `Verified` stamp ... proves only an adjacent claim" check exists to catch.

### F-7 — (c) internal contradiction: Testing Strategy scenario 4 resolves quoting/whitespace two opposite ways in one Expected clause

**RDR passage** (`0008...md:1044-1052`), scenario 4:

> **Scenario** (identity/exactness): near-spellings of the reserved word as declared tag names — `Recognized`, `RECOGNIZED`, a quoted TOML key `"recognized"`, and **a whitespace-bearing variant**.
> **Expected**: the byte-exact rule on the parsed TOML key value treats **every case-variant** as an ordinary unreserved name (**no folding, no trimming**); the quoted form, which parses to the same key value, **is** reserved.

**Cross-reference** (`0008...md:601-605`, normative block 2): "The reserved-key comparison is byte-exact on the parsed TOML key value: case-sensitive, no trimming, no folding (so `Recognized` is an ordinary, unreserved name)." And Load-Bearing Decisions / Identity (`:645-652`): "Near-spellings (`Recognized`, a quoted or whitespace-bearing variant) are **by definition ordinary unreserved names**".

**Why it is a defect.** The scenario lists four inputs but its Expected clause covers two dispositions and assigns them by the wrong discriminator. "Every case-variant" disposes of `Recognized`/`RECOGNIZED`; "the quoted form ... **is** reserved" disposes of `"recognized"` — correctly, since TOML `["recognized"]` and `[recognized]` parse to the identical key. That leaves the whitespace-bearing variant with **no** Expected outcome, and it is the ambiguous one: `[tags." recognized"]` parses to the key value `" recognized"` (byte-distinct -> unreserved under "no trimming"), but Load-Bearing Decisions / Identity flatly classifies "a quoted **or** whitespace-bearing variant" as unreserved — which directly contradicts the scenario's ruling that the quoted form IS reserved. The two passages disagree about the quoted case, and neither settles the whitespace case. Since this RDR's single contract *is* the identity of a name and the comparison rule is normative, an unsettled parse-level edge is a hole in the one thing the document owns. Fix: state the rule over the post-parse key string only, drop "quoted" from the Identity near-spellings list (quoting is a TOML surface artifact, not a name variant), and give the whitespace variant an explicit Expected.

### F-8 — (d) RDR is silent, revealing a missing requirement: `Row.RequiresOwned` is a third un-reserved channel for the key

**RDR passage.** Normative block 4 (`0008...md:618-628`) reserves the key on `Input` producers. Block 2 (`:595-606`) reserves it on tag declarations. The Normative Contracts preamble (`:570-585`) enumerates the channels as exactly two: "the two channels that key can arrive through: declaration (blocks 1-3) and resolve-time data (block 4)". `RequiresOwned` appears nowhere in the RDR.

**Source cite** (`internal/resolve/resolve.go:180-182, 444-458, 125-128`):

```
	// RequiresOwned names owned tag keys the row's evaluation needs.
	// A key absent from the owned snapshot yields owned_state_unavailable.
	RequiresOwned []string
```

and `missingOwned` tests `view.has(key, ProvenanceOwned)`, where `has` is `tv.provenance == prov` — a provenance-**specific** test.

**Why it is a defect.** A normalized row may carry `RequiresOwned: []string{"recognized"}`. At resolve time the key *is* present in the assembled view but under `ProvenanceRecognized`, so `has(key, ProvenanceOwned)` is false and the kernel emits `owned_state_unavailable` naming `recognized` — the precise failure shape RDR 0008 exists to eliminate: a refusal whose text names the reserved key while nothing tells the author the two concepts could not agree. This is neither the declaration channel (a `RequiresOwned` entry references a declared key, and under this RDR a `recognized` declaration is legal and required) nor the `Input` data channel (block 4 constrains producers of `Input.Owned`/`Input.Observed`, not `Row`). So the "two channels" enumeration is incomplete — and the completeness argument defending a single-contract Profile ("They are not separable seams — reserving only the declaration channel leaves the silent-shadowing hole this RDR exists to close") is built on that enumeration. Note RDR 0007 (`Final`) already re-defined `RequiresOwned` as "post-guard write-dependency keys" (`0007...md:23-25`) and separately mandated that guard presence is provenance-blind while `missingOwned` stays provenance-specific (`0007...md:1270-1277`: "The two predicates answer different questions and MUST NOT be conflated") — making this asymmetry a live, peer-acknowledged seam. RDR 0008 needs either a normative clause forbidding `recognized` in `RequiresOwned` (natural home: the 0002 normalizer, alongside block 2) or an explicit written exclusion saying why the case is out of scope.

### F-9 — (d) RDR is silent: normative block 1 is unconditional where `::assemble` is conditional

**RDR passage** (`0008...md:587-593`, normative block 1): "The assembled evaluation view MUST bind the freshly recognized outcome under exactly the tag key `recognized` (`internal/resolve/resolve.go::recognizedTagKey`)."

**Source cite** (`internal/resolve/resolve.go:151`): `if in.Recognized != "" {` — the binding is gated on a non-empty outcome, so with `Recognized: ""` the view carries **no** `recognized` key. (Key Discoveries at `0008...md:191-198` does note the injection happens "only when `in.Recognized` is non-empty", so the RDR knows this — the normative block simply does not carry the condition.)

**Why it is a defect.** The RDR's central normative sentence states an unconditional MUST that the shipped, "already conforming, so no kernel code changes" (`:499-504`) implementation does not satisfy for the empty-outcome input. It is not reachable as a bug — `resolve.go:321` refuses `KindUnmodeledOutcome` for `""` against any normal alphabet before selection — but two things follow. First, Phase 1 ratifies this exact sentence as contract and pins it with a behavioral conformance test (`:960-966`, `:1012-1023`); a literal reading of the MUST is falsifiable by an empty-outcome input, so the test's input domain needs stating. Second, a table whose `Outcomes` alphabet contains `""` would pass the alphabet gate with no reserved key in the view — an interaction with RDR 0002's "closed set of outcome tags" (`0002...md:217`) that neither RDR excludes. The block needs the same "when a recognized outcome is present" scope its own Key Discoveries paragraph already carries.
