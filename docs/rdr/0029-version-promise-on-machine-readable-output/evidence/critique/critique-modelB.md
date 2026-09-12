Model: claude-sonnet-5

## Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | 0029:C1 | The tolerant-reader promise is unenforceable prose; a strict-parsing agent breaks on the very release that announces the promise | Agent's parser hard-fails on `schema_version` appearing, on the release meant to prevent exactly that | §1, premortem, AT-1 |
| C-2 | 0029:C3 | Promotion-to-blocking disclosure is a release-note convention with zero tooling; nothing detects an undisclosed promotion | Consumer's CI goes red on a model nobody edited, with no diagnostic pointing at the cause | §1, §2, AT-2 |
| C-3 | 0029:A5 | Enumeration-seam spike is Pending at lock; two of sixteen tiered vocabularies (CLIError `code`, `flow next` reason) may be structurally unbuildable without inverting `clierr`'s leaf layering | Implementer discovers mid-Phase-1 that C4's promised tier for the largest vocabulary can't be asserted at all | §1, §3 |
| C-4 | 0029:A6 | `schema_version` constant's "one home" is Pending — asserted as settled fact in C1 while its own assumption record is still open | If A6 resolves differently than assumed, C1's normative text (already written as fact) is wrong at lock | §1 |
| C-5 | 0029:A4 | The tier census was already found incomplete once during review (seven of fourteen surfaces missed), and the fix mechanism (grep-based enumeration) has already proven fallible | A shipped surface has no stated tier; a consumer builds on an unstated assumption that later breaks | §1, §2, premortem |
| C-6 | 0029:§approach / 0029:C2 | The entire promise is distributed as prose in `docs/cli-output-contract.md` + `llms.txt`, with no machine-readable, runtime-queryable form | Agent author must read and correctly parse documentation to know what's safe — the same "read a file and hope" failure mode the RDR exists to eliminate for the envelope itself | §1, §3, premortem |
| C-7 | 0029:§decision-rationale (Premortem) / 0029:A2 | The 0.x "free window" has no forcing function; the project has a documented history (four `table.Categories()` growth events) of treating version boundaries informally, with no version signal ever consulted | `schema_version` lands after 1.0.0 by omission, forcing the exact breaking major-bump the field exists to prevent | §2, premortem |
| C-8 | 0029:C3 (re-attribution clause) | "MAY be introduced directly at `blocking` when it re-attributes an existing refusal" is a subjective judgment call with no test and no worked example distinguishing re-attribution from a genuinely new finding | Implementer or reviewer disagrees in good faith about whether a given new code counts as re-attribution, and a verdict-changing code ships undisclosed because it was misclassified as re-attribution | §1, §3 |
| C-9 | 0029:§validation (Testing Strategy, item 7) | S7 retires three existing exact-cardinality tests (`TestReq73`, `TestReq74`, `TestReq80`) as "the anti-pattern," but nothing prevents a *future* PR from re-adding an exact-set assertion over an append-only/growing vocabulary — there is no lint, no CI check, only a testing-strategy sentence | Six months post-ship, a new contributor adds a cardinality-asserting test on a `growing` vocabulary (the natural instinct when adding a test), silently reintroducing the exact defect this RDR retires | §2, premortem |
| C-10 | 0029:§problem-statement / 0029:A3 | The one real in-repo consumer this RDR found (`planEnvelope`) is a four-field tolerant struct that "just happens" to survive; the RDR's confidence that "no current consumer parses strictly" rests on an audit of *this repo only* — external/future consumers (the actual target audience: "an agent author pins a released version") are categorically unauditable | The RDR's central empirical claim (A3) can never be re-verified against the actual audience it's written for; it's Verified against a population of one internal caller and asserted true of an external population that cannot be enumerated | §1, §3 |

## 1. The three most likely ways implementation goes wrong

**Failure 1 — The tolerant-reader promise is a promise this project cannot keep or check.**

Root cause: C1 states "A consumer MUST ignore object properties with unrecognized names, and MUST reject an envelope reporting an unsupported major." This is a MUST directed at code the RDR authors do not control, cannot test, and cannot even inventory. The RDR's own premortem admits this directly: *"a consumer validating the envelope against a strict schema... sees an unexpected `schema_version` key and hard-fails. The promise's first act is a violation of itself."* The mitigating argument offered — that the 0.x series is a "free window" — is not a fix, it's a deferral. It says the failure mode is real but not yet load-bearing. That is not the same as designed-around.

The specific passage that enabled it: `0029:C1`'s tolerant-reader clause, combined with the premortem's own admission in `0029:§decision-rationale` that this is "the one clause here that does not wait for 1.0.0" for disclosure, but the tolerant-reader clause itself has zero enforcement — it's a MUST with no test, no schema artifact, nothing but a documentation sentence in `docs/cli-output-contract.md`.

Symptom the user sees: An agent author who followed reasonable practice — validating machine input with a JSON Schema or a strict Go struct — has their harness break on the exact release meant to prevent breaking harnesses. They open an issue saying "your compatibility promise broke my pipeline," and the maintainers' only answer is "you should have used a tolerant parser," which is true but useless after the fact, and nowhere stated as a *hard requirement* an agent author would go looking for before writing their parser.

**Failure 2 — Silent promotion, because the disclosure mechanism is a sentence, not a system.**

Root cause: C3's entire enforcement mechanism for the single most consumer-damaging event this RDR defines (`info`→`blocking` promotion, "the verdict-changing event") is: *"MUST be disclosed in the release notes for the release that carries it, naming the code."* There is no test for this. The RDR says so itself in Risks and Mitigations: *"The promotion disclosure in C3 is process, not a test, so a promotion ships unannounced... Accepted as a known limit."* Accepting a known limit on the load-bearing clause of the entire RDR — the one event severe enough to flip a consumer's exit code from 0 to 2 — is accepting the central risk the RDR was written to retire.

The specific passage: `0029:C3` and the accompanying Risk/Mitigation in `0029:§trade-offs` ("Accepted as a known limit"), and the Capability Dependencies row stating the disclosure discipline is "asserted by review not by a test."

Symptom: A model that lints clean today starts failing tomorrow with no code change on the consumer's side. The consumer's CI goes red. They diff their own repo, find nothing, and only after wasting an afternoon do they think to check the CLI's release notes — if they even know to look there, since nothing prompts them to.

**Failure 3 — The tier census keeps being wrong, and there's no mechanism that catches the next miss.**

Root cause: A4 was "Verified as to the partition" only after correction — the RDR's own text admits the first pass of C4 named seven of fourteen surfaces and missed seven, two of which ("the two model-authored open namespaces... which the first pass of this census missed by enumerating only the top level and `data`") were caught only because a review lens went looking inside `findings[]`. The RDR's fix for this recurring miss is more prose discipline: *"A surface added later takes a tier assignment in the same document as part of the change that adds it; an unassigned machine-readable surface is a defect."* That's a policy statement, not a gate. Nothing in CI checks that every emitted JSON key has a corresponding tier row.

The specific passage: `0029:A4`'s own "If wrong" clause, which reads as already partially triggered: *"C4 is incomplete and an unassigned surface has no stated promise — the exact gap this RDR exists to close, reappearing inside the fix."*

Symptom: Some future emitted field — a new verb's payload, a new nested struct — ships with no tier, and a consumer has no way to know whether it's safe to switch on exhaustively. They find out when it changes shape and their code silently mis-parses (the exact failure mode from the Problem Statement: "its parse silently takes the wrong branch").

## 2. The section that will be rewritten within 6 weeks

`docs/cli-output-contract.md`'s new tier table and prose, i.e. Activation Step 1 (`0029:§phase-2-operational-activation`).

This is the section most exposed to reality because it's the only part of the RDR that is *pure prose with no code backing it whatsoever* — no test asserts its content matches the actual C4 census, no generator keeps it in sync with the enumeration seams, and it's explicitly the piece the RDR itself flags as sitting in the part of the read path the project "deprioritizes" (`llms.txt` tells agents to prefer asking the binary over reading any file). The RDR even names its own successor: *"Exposing the tier table from the binary is the successor if the prose register proves unread; it is out of scope here and not a fallback, since it answers a different question."* That's an admission, in the same document, that the chosen mechanism is provisional and a better one is already known. The first time an agent author reports "I didn't know the tier was X" — which will happen, because nothing forces them to read this file rather than just calling the binary as `llms.txt` tells them to — this section gets rewritten into a `--as=json`-emitted tier manifest, i.e., exactly the mechanism this RDR chose not to build.

## 3. The assumption that will not survive first contact with a real user

A3: *"No current consumer parses the envelope strictly enough that one added field breaks it."*

This is Verified, but only against a population of one — the repo's own internal `planEnvelope` reader — plus the absence of any strict parser in CI, tests, or `.claude/` skills. The RDR is explicit that the actual audience is external: *"An agent author pins a released `intrastate` version in a skill, a CI job, or a harness, and writes code that parses `--as=json`"* — someone the RDR authors have never seen and cannot audit. The verification method is "Source Search" over a codebase that, by the Problem Statement's own framing, is not where the risk lives. The risk lives in every external agent skill, CI harness, and one-off script written against `--as=json` that the authors have zero visibility into and that plausibly *does* use `DisallowUnknownFields` or a JSON Schema `additionalProperties: false`, because that's the default posture any careful engineer takes when parsing machine input, and the RDR's own Problem Statement describes exactly this population ("if it lints models — expects a model that passes today to pass tomorrow"). The first real bug report will be from exactly this consumer, and the RDR's own premortem already scripted this outcome before it happened.

## 4. Premortem

It is six weeks after `intrastate` v0.3.0 shipped `schema_version`. Two things have gone wrong, and they are the two things the RDR flagged as accepted risk rather than designed away.

First: a user built a GitHub Action around `intrastate lint --as=json`, piping the output into a small Python script that does `json.loads()` then validates against a hand-written `TypedDict`-derived JSON Schema with `additionalProperties: false`, because that's the natural thing to do when you want your CI to fail loudly on unexpected drift rather than silently ignore it — the exact caution the Problem Statement's own opening paragraph describes ("if it lints models — expects a model that passes today to pass tomorrow"). The v0.3.0 release adds `schema_version` as a non-`omitempty` top-level key. The schema validation step fails with `Additional properties are not allowed ('schema_version' was unexpected)`. The user's build goes red. They open an issue titled "0.3.0 broke my CI with no code change on my end." The maintainer's answer — "the compatibility doc says to ignore unrecognized keys" — is correct and useless: the user never read `docs/cli-output-contract.md`, because `llms.txt` (which they did read, because it's the thing the project tells agents to consult) told them to prefer asking the binary, and their skill was written before the tier table existed.

Second, in parallel: a contributor adds a new graph-lint advisory finding, `graph-redundant-escape`, at severity `info`. It's a good addition and passes review — nobody flags it, because C3 says introduction at `info` needs no disclosure beyond release notes, and the PR author, reasonably, doesn't loop in anyone downstream. It ships in v0.3.2, a patch release, without release notes calling it out beyond a changelog line reading "add graph-redundant-escape lint finding." Three weeks later, the same contributor, working a different card, notices this new advisory code actually always co-occurs with an existing blocking condition and decides to promote it to `blocking` in v0.4.0 to consolidate the two. This IS the disclosed, verdict-changing event C3 requires — and the release notes do call it out by name, exactly as C3 demands. But a downstream user's CI, which had been silently accumulating this advisory finding on models that were otherwise clean, goes red on the v0.4.0 upgrade. The user does not read release notes before upgrading (nobody does, for a patch-adjacent minor). They spend forty minutes diffing their own model file against yesterday's version, find nothing changed, and only think to check `intrastate`'s changelog after a teammate suggests it. The disclosure obligation was met to the letter of C3 and the user was still surprised, because "disclosed in release notes" and "discovered before the pipeline goes red" are not the same thing, and the RDR conflates them.

Both failures were named, nearly verbatim, in the RDR's own Risks and Mitigations and Premortem sections, and both were accepted rather than designed around.

## 5. Acceptance tests that would have caught each failure at RDR-review time

```gherkin
Feature: schema_version tolerant-reader claim is load-bearing, not aspirational

  Scenario: A representative strict external consumer is modeled, not just the in-repo one
    Given a JSON Schema with "additionalProperties: false" generated from the
      pre-schema_version envelope shape (the shape any careful external consumer
      would plausibly have authored before this RDR)
    When the post-change envelope (carrying schema_version) is validated against it
    Then the RDR's own text must show this failure explicitly, not only in prose,
      and must record an explicit decision on whether the project ships a
      deprecation/announcement window before schema_version reaches consumers who
      are NOT proven to be tolerant (i.e., everyone outside this repo)

  Scenario: Promotion disclosure has a machine-checkable gate, not only a release-note convention
    Given a finding code that moves from severity "info" to severity "blocking"
      between two tagged releases
    When the release is prepared
    Then CI fails the release build unless the diff between the two releases'
      severity tables is named in the release notes
    # As written, 0029:C3 has no such gate; "Accepted as a known limit" in
    # 0029:§trade-offs is the ledger entry this AT would have forced the RDR
    # to close instead.

  Scenario: Every emitted machine-readable key has a tier assignment, enforced by a program, not a sentence
    Given the full, current --as=json output surface across every verb
    When a new field or vocabulary member is introduced anywhere in the tree
    Then a CI check (not a documentation clause) fails the PR unless the new
      surface has a corresponding tier row in docs/cli-output-contract.md
    # 0029:C4's "an unassigned machine-readable surface is a defect" is
    # currently enforced by nothing; this AT is what would have caught the
    # A4 miss (seven surfaces on the first pass) automatically, rather than
    # via a review lens noticing by hand.

  Scenario: The enumeration-seam spike (A5) is resolved before lock, not deferred into implementation
    Given the CLIError "code" vocabulary has 61 raise sites across three
      packages with no single declaration site
    When a clierr-side registry is attempted per C4's seam obligation
    Then either the accessor builds without an import cycle, or the RDR
      records — before lock, not as a Phase-1 discovery — that this
      vocabulary's tier is prose-only and cannot be asserted
    # As written this is a Pending assumption at lock time; an AT run
    # pre-lock would force the spike to actually run rather than ride on
    # Method: Spike, Status: Pending into Finalization.

  Scenario: Re-attribution vs new-finding classification is not a judgment call
    Given a new finding code introduced directly at severity "blocking"
    When it is reviewed under C3's re-attribution exception
    Then the PR must cite the specific pre-existing refusal code it
      re-attributes, and a test must assert the two codes are mutually
      exclusive on the same input
    # 0029:C3's re-attribution clause has no worked example and no test;
    # this AT is what would have surfaced that gap before lock.
```

