Model: claude-sonnet-5

## 6. Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | 0028:A3 | Verified assumption is a fixture-fit spike, not a generality claim; RE2-suffices and every-anchor-selects-exactly-one is verified against exactly one consumer's 33 records and 26-row README, not against the class of "any text artifact" this RDR frames itself as solving | Every NEW consumer of `edit` re-runs the anchor-viability spike by hand against their own corpus in production, because lint cannot check it (C1.4 explicitly does not prove match-cardinality) — the first drifted anchor in a second adopter's fleet is discovered as a live `edit_anchor_unmatched` outage, not a design-time finding | §1, premortem |
| C-2 | 0028:A10, 0028:A11 | Two Critical Assumptions the gate itself requires to be resolved before lock (G-assumptions: "any that remain Pending... with a plan to verify before implementation begins") are still `Pending` at Draft, and C1.3's own prose treats their answers as already settled ("the gate is carried to the accessor as state, never read across the seam (A10)"; "That arm is 0004's, unchanged here and out of this RDR's scope (A11)") | Implementation starts, discovers A10's back-import cycle risk or A11's deadline-arm interaction is not as assumed, and the normative contract C1.3 has to be rewritten mid-implementation — the exact "Status consistency" violation the Finalization Gate template warns against | §1, premortem, G-assumptions |
| C-3 | 0028:C1.3 (`re-anchor:`) | The re-anchor pass re-runs every rule's anchor over the POST-EDIT buffer, but only checks that each rule's OWN anchor still selects its own rewritten line (or zero, if deleted) — it does not re-check that a sibling rule's anchor still selects exactly one line after this rule's rewrite | A multi-key entry where rewriting key A's line causes key B's anchor to now match zero or two lines (e.g., a templated line whose content changes shape after substitution) ships a partially-consistent write: A's contract is satisfied, B's is silently unchecked until B's own next write, when it surfaces as a confusing `edit_anchor_unmatched` on a rule that was never touched this invocation | premortem, §1 |
| C-4 | 0028:F4, 0028:F5, 0028:S31 | The RDR's OWN stated safety net for "anchor selects the wrong single line" is post-mutation read-back through the SAME role's reader that the model author also authored/configured — there is no independent verification path, so an author whose reader has the same blind spot as their anchor (e.g., both keyed to the first bullet-shaped line) gets silent corruption with a false "consistent" read-back | The artifact is corrupted (wrong line rewritten) but the CLI reports success, because the reader that verifies happens to read a decoy that also validates — an outcome the RDR acknowledges exists (F5) but treats as adequately mitigated by "read-back is the commit-time check," when read-back is authored by the same fallible party as the anchor | §1, premortem |
| C-5 | 0028:§decision-rationale (Alternatives matrix), 0028:ALT3 | Alternative 3 (adapter registry) is rejected largely on "document knowledge enters the tool," but the RDR's own MVV and Testing Strategy encode deep document knowledge of ONE specific consumer's document family (RDR flow's Status bullet and README row shape) into the acceptance criteria and even the illustrative code — the "generic line editor" framing is aspirational, the actual proof surface is one document family | A second, differently-shaped consumer (e.g., YAML front-matter, a different indentation convention, a table with escaped pipes) discovers the anchor grammar and testing muscle memory built during this RDR's implementation and review does not generalize, and the "generic as the flat-JSON path writer" claim (Background) is falsified on contact | §1, premortem, §3 |
| C-6 | 0028:C1.5 (`one-way:`) | `clear = "line"` deletes the line with no way to re-create it via `edit` (creation is fenced out entirely); the RDR states the mitigation is "declare `clear` only where the LINE is the key," but nothing in lint (C1.4) or in the runtime contract stops an author from declaring `clear = "line"` on a key whose anchor also matches a shared/reusable row structure, and nothing prevents the sequence BR2-adjacent mistake: clear, then discover you need the field back | An author clears a field, then needs to reinstate it (a Status flips back from a terminal state, or a corrected record needs a re-added bullet), and the tool refuses to help at all — the recovery path is "edit the file by hand outside intrastate," which quietly reopens exactly the wrapper/manual-edit drift this RDR exists to close, just relocated to the recovery path | premortem, §2 |
| C-7 | 0028:§approach, 0028:C1 (Normative Contracts header) | Profile is declared `large` on the strength of "one contract" (C1), but C1 has six clauses (C1.1–C1.6) spanning three separate seams: the loader/lint seam (C1.1, C1.4), the accessor Write seam plus a NEW cross-cutting context-tag channel (C1.2, C1.3, A1), and the UNRELATED command-reader placeholder vocabulary (C1.6, extending 0025:C2, a different predecessor's contract) — C1.6 is explicitly "scored separately... because the matrix above is about the write carrier and C1.6 is not one," which is the RDR admitting mid-document that it is authoring a second, independent load-bearing contract inside a "one contract" Profile claim | A reviewer or future maintainer treating this as one indivisible change unit cannot land or revert C1.6 (the read-side placeholder extension) independently of C1.1–C1.5 (the write carrier), even though the Decision Rationale itself separates their justification, their alternatives, and their scoring — a partial revert or partial rollout is architecturally blocked by a Profile decision that undercounted contracts | premortem, §2, G-proportionality |
| C-8 | 0028:A6, 0028:§capability-dependencies | A6 confirms "no accessor dump/marshal/normalize round-trip exists anywhere in internal/," so its own stated failure mode ("a dump/normalize round-trip drops or mangles edit silently") has, by the assumption's own text, "no live target" — meaning this assumption is Verified against nothing that exercises the new code, deferring the actual dump/normalize risk to whenever such a round-trip is added later, with no test or lint gate in THIS RDR to catch it then | A future RDR adds a table dump/normalize path (plausible — `intrastate lint --model` output, a table migration tool), forgets the `edit` carrier because A6 said "not exclusionary, just add a parallel `if`," and `edit` entries silently vanish from a round-tripped model — the exact silent-corruption class this whole RDR was written to eliminate, recreated one layer up | §1, premortem |
| C-9 | 0028:§prerequisites | The RDR explicitly ships with the consumer's OWN acceptance corpus in a known-broken state for two rows (0001, 0006, the unlinked `\| NNNN \|` README form) and calls this "partial adoption" that "blocks adoption... not blocking implementation" — but the MVV fixture is defined to exclude exactly these two rows, so the acceptance test that is supposed to prove "the composed consumer scenario works" is proving it against a corpus curated to avoid the one class of pre-existing drift the real corpus has | The feature ships, passes its MVV, and the FIRST thing the actual consumer does — flip record 0001 or 0006's state — refuses with `edit_anchor_unmatched`, which is contractually correct but was foreseeable and undemonstrated at RDR-review time; the "delivered" claim in the Problem Statement is true only for 24 of 26 rows on day one | §2, §3, premortem |
| C-10 | 0028:C1.3 (`write:`, no-op clause) | "A post-edit buffer equal to the input is not written at all" is witnessed ONLY by inode-stability (S20) because content can't distinguish it — this is a reasonable engineering call, but it means a caller cannot distinguish, from CLI output alone, "your edit was a no-op because the value already matched" from "your edit applied normally"; both report success with no visible signal, and no scenario or contract clause requires the CLI to surface which happened | An agent-loop caller (the RDR's own stated audience, per A8's framing) driving repeated `resolve \| set-state` cannot tell from output whether a state transition actually occurred this run or was already applied by a previous run, which matters for idempotency reasoning and audit trails, and has to shell out to `git diff` or inode-watch to find out — undocumented as a caller-visible gap | §1, premortem |

## 1. The three most likely ways implementation goes wrong

### Failure 1: The "generic line editor" is actually a one-consumer special case, and the second consumer breaks it

**Root cause.** The RDR's rhetorical frame (Problem Statement, Background, Decision Rationale) is "a first-class declared, in-process write carrier" as general infrastructure — explicitly compared to "the flat-JSON `path` writer" for genericity. But every piece of verification evidence, every MVV step, and most of the Testing Strategy (S8, S12, S13, S19, S26) is grounded in exactly one document family: the RDR flow's own `- **Status**:` bullet and README index row. A3 and A4 — the two spikes that actually prove RE2 suffices and read-back round-trips — are run against 33 files and one reader, both belonging to the tool's own self-hosted consumer.

**Enabling passage.** A3's assumption text: "Go's RE2 `regexp` is expressive enough for THE CONSUMER'S two anchors... each anchor selects exactly one line in every EXISTING CONSUMER RECORD and README" (0028:A3). The "if wrong" clause even concedes the fallback is "the anchor grammar needs a second dialect... or the consumer's records refuse `edit_anchor_ambiguous`" — the assumption's own escape hatch is "make the CONSUMER conform," not "the tool generalizes." Meanwhile Decision Rationale claims Alternative D (adapter registry) is rejected because "document knowledge enters the tool, which the kata fences out" — yet the MVV and half the S-scenarios encode exactly one document's structural knowledge (bracket-qualifier wrapping, table-row column counting) into the acceptance bar.

**Symptom.** A second consumer with a differently-shaped artifact (YAML front matter, an indented list, a table using escaped pipes, multi-byte content near an anchor boundary) discovers that "generic" meant "generic within RE2's syntax," not "field-tested outside one shape." Their first attempt to write a working anchor either fails cryptically at lint (a legitimate but under-explained `edit_anchor_invalid`) or, worse, passes lint and then hits `edit_anchor_ambiguous` or `edit_anchor_collision` against real data that the tool author never tested a comparable case for, because A3's spike evidence never had to.

### Failure 2: Two Critical Assumptions are Pending at Draft, but the normative contract already assumes their answers

**Root cause.** RDR process requires (per this record's own Finalization Gate template, G-assumptions) that no assumption remain `Pending` "with[out] a plan to verify before implementation begins," and explicitly bars "Status consistency" violations — "no assumption marked Pending... may have settled-fact prose elsewhere in the RDR depending on it." A10 and A11 are both `Pending` (method: Source Search, "to verify"). Yet C1.3's `read-back:` clause states as settled fact: "the gate is carried to the accessor as state, never read across the seam (A10)" and "That arm is 0004's, unchanged here and out of this RDR's scope (A11)" — both written in declarative, not conditional, voice, inside the NORMATIVE contract clause a lock would freeze.

**Enabling passage.** 0028:A10's own "If wrong" text: "C1.3's pre-check falls back to A2's named alternative — hoist it to `internal/cli/flow_state.go` before `exec.Write` — and the refusal is then minted by the CLI rather than the accessor, which MOVES WHERE THE GATE-NAMING DETAIL IS COMPOSED." That is not a footnote-level risk — it changes which layer owns a normative Detail-composition responsibility that C1.3 states as fact. A11's "If wrong": "this RDR must amend the deadline arm for non-spawning carriers, which widens it into `internal/accessor/executor.go` beyond A10's one field" — i.e., the blast radius of A11 resolving "wrong" is a scope increase into shared executor code this RDR currently claims is out of scope.

**Symptom.** Implementation (Stage 8) starts against a contract whose C1.3 prose already presumes outcomes for two Source-Search checks that haven't been run. If either resolves against the RDR's stated expectation, C1.3 needs a normative rewrite mid-implementation — not a clarification, an amendment to the clause locked at Finalize. This is the single most mechanically certain "will be rewritten within 6 weeks" candidate in the document (see §2 below).

### Failure 3: Silent wrong-line corruption is real, acknowledged, and given a mitigation that only works when the mitigation's author didn't make the same mistake as the anchor's author

**Root cause.** F4, F5, and S31 all describe the same hazard: an anchor whose exactly-one match selects the WRONG line (a decoy — quoted example, template comment, a stale copy) while the real target line has drifted elsewhere. The RDR's stated defense is uniform across all three: "the read-back through the role's reader... refuses `read_back_mismatch`... which is the reason read-back is the commit-time check." This is true only if the reader's own selection logic is independent of the anchor's selection logic. Nothing in C1 requires or verifies that independence — it's a property of the ambient reader implementation, asserted true for the RDR flow consumer (because A3/A4 spiked it) and NOT a contract this RDR imposes on future readers.

**Enabling passage.** S31 states the failure mode explicitly and then treats it as covered: "the write applies, and the read-back... refuses `read_back_mismatch` with `Applied()` TRUE — the one arm where the applied sense is set on a failure... and the post-mutation net F4 and F5 name as the SOLE DEFENCE against silent wrong-line corruption." Calling it "the sole defence" while also stating it depends on an ambient property of an unrelated reader implementation is not a defense — it's a documented single point of failure, and the RDR names it as such without hardening it further ("a per-rule value guard is a successor's," Briefly Rejected).

**Symptom.** A consumer whose reader was written casually — say, a `grep`-first-match reader that would itself be fooled by the same decoy the anchor matched — gets silent data corruption reported as success. This is the worst failure class the RDR itself identifies (irreversible-looking, no error surfaced) and the mitigation is "have a good reader," which is not enforced, tested, or even lint-checkable by this RDR's own admission (C1.4: "what lint does NOT prove: that an anchor matches exactly one line of a particular file").

## 2. The one section that will be rewritten within 6 weeks of shipping

**C1.3's `read-back:` clause** (the gate pre-check, and by extension the site-selection prose throughout C1.3 that leans on A10 and A11).

This is not a guess about general RDR fragility — it is the direct, mechanical consequence of shipping a locked normative clause whose content depends on two assumptions the record itself marks `Pending` at Draft. The Finalization Gate template this very RDR carries (G-assumptions) exists specifically to prevent locking with unresolved Pending assumptions feeding settled prose — and this RDR's Prerequisites section (unchecked checkbox: "All Critical Assumptions verified") explicitly flags that this hasn't happened yet. Two outcomes are possible when A10/A11 actually get verified during Stage 8:

- A10 resolves against the stated expectation → the gate-refusal Detail-composition site moves from the accessor layer to `internal/cli/flow_state.go`, which changes WHO composes the Detail message C1.3 currently attributes to the accessor's `Registry` field. That's a rewrite of C1.3's `read-back: SITE:` paragraph, not a comment update.
- A11 resolves against the stated expectation → the deadline arm needs amendment reaching into `internal/accessor/executor.go`, which C1.3 currently declares "unchanged here and out of this RDR's scope." That's a scope change discovered after lock, requiring either an amendment RDR or a reopened C1.3.

Either path rewrites prose inside the locked contract within the implementation window — i.e., within weeks of "shipping" the Draft-to-Final transition, likely during Stage 8 itself, which the RDR's own status timeline puts inside the 6-week frame trivially.

## 3. The one assumption that will not survive first contact with a real user

**A3/A4's implicit universal claim: "RE2 with no lookaround is expressive enough, and each anchor selects exactly one line" generalizes beyond the RDR flow's own artifact corpus.**

Both A3 and A4 are marked `Verified` — correctly, by their own narrow terms: they verified their claim against 33 records and one README belonging to the tool's own self-hosted consumer, using a reader (`rdr status -json`) written by the same team building the write carrier. That's not a criticism of the spike methodology; it's a correctly-scoped spike for the stated MVV. The failure is rhetorical scope creep in the surrounding document: the Problem Statement and Approach sections describe `edit` as solving "text-artifact writes" generally ("A model author whose state lives in text artifacts... wants..."), and Decision Rationale claims prior-art alignment with Ansible `lineinfile` and Puppet `file_line` — tools used across an enormous diversity of unrelated file shapes in the wild, with none of that diversity represented in this RDR's verification evidence.

The very first real user outside the RDR-flow self-hosting loop who tries to declare an `edit` rule against a genuinely different document shape (a YAML list item, a JSON-in-a-comment marker, a line that legitimately contains a literal `{` in RE2-significant position, a file using significant leading whitespace where "the whole replacement line" semantics silently eats indentation the author didn't intend to touch) will hit either an unhelpful lint failure with no comparable prior test to reference, or — worse per Failure 3 above — a corrupted write that a weaker reader doesn't catch. The RDR's "generic as the flat-JSON `path` writer" framing (Background) will not survive that contact; `path` truly is format-agnostic (it writes whatever bytes it's given to a whole file), while `edit`'s correctness is contingent on RE2 anchor-craft skill that this RDR has only demonstrated once, on friendly, self-authored data.

## 4. Premortem

*It is eight weeks post-lock. `edit` shipped, Stage 8 landed, the RDR-flow consumer migrated `rdr-write.toml` off its `sed` strings as planned. Here is what actually happened.*

During Stage 8, resolving A10 turned up exactly the branch its own "If wrong" clause predicted: a second, test-only construction site for `accessor.Registry` outside `flowbind.go::Registry` existed in an integration-test helper that did not thread `allowCommands`, so the new field's zero-value silently meant "gate off" for every test built through that helper. Nobody had audited that site because A10's evidence (marked Pending, deferred to implementation) never listed it. The fix required touching `internal/accessor` test scaffolding that A10 explicitly said would NOT need to change ("no new import, no Binding interface method"). C1.3's `read-back: SITE:` paragraph was edited three times during Stage 8 review before the wording stabilized — exactly the "rewritten within 6 weeks" prediction landing on schedule.

Separately, the consumer's own migration surfaced the Prerequisites-section warning nobody read carefully: records 0001 and 0006, in the unlinked `| NNNN |` README form, could not be flipped through the new machinery. The team had assumed "we'll fix the two rows first" was a five-minute chore; it turned out the unlinked form was load-bearing for an unrelated tool (a changelog generator) that keyed off the exact string shape, so normalizing it broke that tool's regex, which had never been reviewed against this RDR because it was "the consumer's problem, not this RDR's." Two weeks were spent untangling a dependency this RDR's Prerequisites section flagged existed but assumed would be trivial to clear.

The first non-RDR-flow adopter — a different internal tool tracking deployment approvals in a markdown file with nested checklists — tried to declare an `edit` rule against a line inside a nested list item. Their anchor, `^  - \[ \] Approved by: (.+)$`, worked fine in isolation but their file had a SECOND checklist item elsewhere in the file using the same "Approved by:" label inside a code-fenced example block documenting the file's own format (a decoy exactly matching F5's described hazard). Their reader — a quick, three-week-old command reader nobody had stress-tested — happened to also parse the first "Approved by:" line it found in the file, which was the decoy, not the real one, because their reader's own scan short-circuited on first match. The write applied to the decoy. Read-back "verified" against the decoy. `Applied()` came back true, no error surfaced. The real approval line was untouched; the deployment gate downstream (reading the real line via a different, correct parser) still showed "unapproved" for three days while the team's dashboard, keyed to the write-side confirmation, showed the deployment as cleared. Nobody noticed until the actual gate blocked a release everyone believed was already approved.

The postmortem for THAT incident cites `0028:F5` and `0028:S31` almost verbatim — the RDR had named this exact failure mode and mitigation gap in its own text, at draft time, and shipped anyway because the "sole defence" (a correctly-independent reader) was never a contract this RDR enforced, only an ambient property true of the one reader it tested against.

Meanwhile, a second migration attempt inside the RDR-flow consumer itself hit C1.5's one-way `clear` semantics the hard way: an author cleared a Status bullet meaning to immediately re-set it in the same batch operation, discovered `edit` refuses `edit_anchor_unmatched` on the second write in the same run (per S29's documented behavior), and had to hand-edit the file outside `intrastate` to recover — reopening, for that one incident, precisely the "an edit no model declared and no lint covered" drift class this RDR's Problem Statement names as the reason for its own existence.

## 5. Acceptance tests that would have caught each failure at RDR-review time

**For Failure 1 (single-consumer genericity):**
```
Scenario: Anchor viability across a second, unrelated document family
  Given a document family NOT authored by the RDR-flow consumer
    (e.g., YAML front matter, or a markdown file with nested list indentation)
  And a `- **Status**:`-shaped bullet is NOT present
  When an author declares an `edit` rule with a hand-written anchor
    for that family's natural state-holding line
  Then `intrastate lint --model` must not require special-casing
    beyond what is documented as generally true of RE2 anchors
  And at least one corpus outside `rdr/tools/rdr/testdata` must be
    walked by an A3-style spike before A3 is marked Verified
    for anything broader than "the RDR-flow consumer's records"
```

**For Failure 2 (Pending assumptions feeding settled prose):**
```
Scenario: Gate check — no Pending assumption may be cited as settled fact
  Given A10 is Status: Pending
  And A11 is Status: Pending
  When the Finalization Gate's Assumption Verification response is drafted
  Then it must list A10 and A11 under "remain Pending... with a plan to
    verify before implementation begins"
  And it must flag every normative-contract sentence that reads A10 or A11
    as already-decided fact (C1.3's "read-back: SITE:" and its A11
    deadline-arm paragraph) as provisional, not lockable, until resolved
  And the RDR must not proceed to Final while C1.3 contains declarative
    prose depending on a Pending assumption
```

**For Failure 3 (silent wrong-line corruption, decoy + weak reader):**
```
Scenario: Read-back independence is a contract, not an ambient property
  Given an anchor whose exactly-one match is a decoy line
  And a reader whose own selection logic ALSO selects the same decoy
    (i.e., anchor and reader share the same blind spot)
  When the write applies and read-back runs
  Then the RDR must specify what happens — today it explicitly does not:
    F5/S31 assume the reader is independent and correct, with no lint,
    runtime check, or authoring guidance forcing that independence
  And the RDR must either (a) require an anchor/reader-independence
    property as part of C1's contract, (b) provide a lint-time or
    apply-time cross-check between anchor selection and reader-visible
    line, or (c) explicitly document this as an accepted, UNMITIGATED
    risk in Risks and Mitigations rather than implying F4/F5 constitute
    a working defense
```

**For the C1.5 one-way `clear` trap (surfaced in premortem, folds into C-6):**
```
Scenario: Clear-then-restore in one logical operation
  Given a record whose Status bullet is cleared via `clear = "line"`
  When the SAME batch/session immediately needs to re-establish that
    line with a new value
  Then the RDR must document the required recovery path explicitly
    (not merely "declare clear only where the line is the key") —
    e.g., "author two separate models" or "this is an intentionally
    unsupported workflow, users must not use `clear` for any field
    that might need reinstatement in the same operational flow"
  And this constraint must appear in Risks and Mitigations, not only
    inferred from C1.5's "one-way" prose and S29's test scenario
```

**For the Profile/contract-count mismatch (C-7):**
```
Scenario: Profile re-validation counts C1.6 as a second contract
  Given C1.6 is explicitly "scored separately... because the matrix
    above is about the write carrier and C1.6 is not one" (Decision
    Rationale)
  When the Proportionality gate response re-validates Profile against
    contract count
  Then it must count C1.6 (extending 0025:C2, a different predecessor)
    as an independent load-bearing contract from C1.1-C1.5
  And either split C1.6 into its own RDR/contract id, or explicitly
    justify why one Normative Contracts section may carry two
    independently-scored, independently-alternatived decisions
    under a single "C1" id and a Profile of "one contract"
```

**For the MVV-excludes-known-drift gap (C-9):**
```
Scenario: MVV fixture must include at least one known-drifted real record
  Given the consumer's live README has two rows (0001, 0006) in a form
    the new `edit` rules cannot address
  When the Minimum Viable Validation fixture is built
  Then it must include at least one of these rows (or an equivalent
    fixture reproducing their shape) as a NEGATIVE scenario demonstrating
    the documented refusal, rather than curating the fixture to exclude
    all known-drifted data
  And the Scope Verification gate response must confirm the MVV proves
    the "delivered" claim against the REAL corpus's actual shape
    distribution, not a cleaned subset of it
```
