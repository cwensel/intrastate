Model: claude-opus-5[1m] (dispatcher consolidation)

# 3amigo consolidation — cli/0020

Three isolated persona passes (no cross-persona visibility). Model stamps:
persona-1-pm claude-sonnet-4-5, persona-2-implementer claude-sonnet-5,
persona-3-qa claude-sonnet-5.

## Hotspots (≥2 personas anchoring a finding on one id)

Mechanical count over `rdr anchors` across the three files, restricted to
FINDING anchors (the raw count also picks up widening notes and non-findings):

- **0020:C1** — 2 findings (implementer P2-M2 duplicated empty-check literal,
  implementer P2-L1 adjacent stale comment) plus QA's grounding of its ordering
  claim. The one true hotspot: the empty-value hoist is where all three
  personas' attention landed.

`0020:§failure-modes` and `0020:§risks-and-mitigations` each surface twice in
the raw count, but from PM's widening prose plus one QA finding — not two
independent findings. Not hotspots.

## Ledger

| id | Persona | Sev | Anchor | Finding |
| --- | --- | --- | --- | --- |
| P1 | PM | — | — | No anchored findings. Outcome clear and singular; widened to consequences/risks/failure-modes/normative-contracts and found no gap. |
| P2-M1 | Implementer | Medium | 0020:D-selection-predicate | Exact splice point of the hoisted empty-value check is prose-only; illustrative code omits the reserved/owned/duplicate switch, so the splice could invert duplicate-vs-empty precedence. |
| P2-M2 | Implementer | Medium | 0020:C1 | The empty-check message literal will exist twice (parseTags + canonicalValue); RDR does not say whether to share a helper. |
| P2-L1 | Implementer | Low | 0020:C1 | `codeInvalidFor`'s doc comment claims `--outcome` shares this path; it does not. Pre-existing, adjacent to Step 1's diff. |
| P2-L2 | Implementer | Low | 0020:§implementation-plan | MVV fixture's fourth key's declared kind unstated — self-answered from the spike fixture. |
| P3-L1 | QA | Low | 0020:S1 | Array-carrier echo has no witnessed fixture for its exact wire form; inferable from `Observed map[string]string`, not witnessed. |
| P3-L2 | QA | Low | 0020:§failure-modes | No scenario exercises the documented "Silent" misspelled-key failure mode. |
| P3-I1 | QA | INFO | 0020:S2 | "outcome" in S2/MVV-3 is unambiguous only after cross-referencing the resolve payload JSON keys. |

## Dispositions

- **P1** — no findings. Nothing to disposition.
- **P2-M1** — **fixed**. Grounded: duplicate (:663) precedes the empty site
  (:723) at HEAD, and `seen[key]=true` is at :667 (ground-source.md Q1), so
  duplicate already wins over empty-value for a repeated key. The RDR's stated
  order (D-selection-predicate) IS that precedence, so the splice is bounded,
  not free — but the illustrative code showed the empty check with no duplicate
  check in view, inviting the inverting splice. Fixed by making the binding
  explicit in C1 and D-selection-predicate and by showing the switch in the
  illustrative snippet.
- **P2-M2** — **dismissed-with-cite**. Implementation latitude, already
  reserved: §technical-design grants "implementation latitude, bounded by C1's
  refusal list: whichever shape is taken, the empty-value arm must still fire
  for a carrier". Whether one literal or a shared helper produces that message
  is not a normative property — C1 pins the message BYTES (fixture F4), which is
  the drift guard the finding asks for. Specifying the helper would be the
  over-spec trap on an identity/format RDR.
- **P2-L1** — **charted-to-successor**. Real (grounded: ground-source.md Q3
  confirms `--outcome` refuses inline at flow_resolve.go:239 and never enters
  `codeInvalidFor`), but pre-existing and outside 0020's seam — 0020 makes no
  claim about `--outcome`. See Charted below.
- **P2-L2** — **dismissed-with-cite**. Self-answered by the persona from
  `evidence/spikes/carrier-table.toml`; the desk trace row 1 already names that
  fixture as MVV step 1's witness. No gap.
- **P3-L1** — **fixed** (narrowly). The unwitnessed array leg is already
  disclosed twice (S1's own text, desk-trace row 2). What was missing is the
  wire form the red test must assert. Fixed by naming it in S1 from the
  `Observed map[string]string` field type, which makes the assertion writable
  without inventing a witness the spike never produced.
- **P3-L2** — **dismissed-with-cite**. Re-raise of a decided call: the silent
  misspelling arm is an ACCEPTED cost, adjudicated in the Disposition table
  ("**silent** — accepted open-world cost"), in §failure-modes' Diagnosis /
  Recovery, and survived the propose premortem (§decision-rationale). Its
  downstream behaviour (no-match refusal / escape row) is the resolution
  kernel's contract, not this admission seam's — testing it here would pin
  another RDR's surface.
- **P3-I1** — **dismissed-with-cite**. INFO, self-resolved on inspection; S2
  already enumerates the eight payload fields and names fixture F2 with the
  literal JSON. The ambiguity exists only for a reader who does not open the
  cited fixture, which the scenario requires.

## Charted

- **Charted (P2-L1):** `internal/cli/flow_input.go::codeInvalidFor`'s doc
  comment (:754-756) says `--tag` and `--outcome` take `flow-tag-invalid`,
  implying `--outcome` routes through this helper; it does not — it refuses
  inline at `flow_resolve.go:239`. The code VALUE it names is correct, so this
  is a comment-accuracy defect, not a behaviour defect. Out of scope for 0020,
  whose seam is `--tag` admission and which makes no `--outcome` claim.
  Suggested successor: a kata (comment fix), not an RDR — no decision is at
  stake. Step 2 of 0020 rewrites a DIFFERENT comment in this file, so the two
  should not be conflated in review.
