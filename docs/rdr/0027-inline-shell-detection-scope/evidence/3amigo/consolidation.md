# 3amigo consolidation — RDR 0027

Three isolated persona passes (PM, Implementer, QA), each run as a separate
sub-agent that saw the record and its own persona block only. Model stamps are
in the per-persona files.

## Hotspots (anchored by >=2 personas)

Mechanical count over `rdr anchors --record 0027 persona-*.md`. Overlap marks a
hotspot passage, not a validated finding; a single-persona finding is not weaker.

| id | personas | what they were each looking at |
|---|---|---|
| `0027:C1` | 3 (PM, Impl, QA) | PM: the `promise:` clause ships text nothing verifies. Impl: does the predicate match the code. QA: the `report:` detail bar is applied unevenly across scenarios. |
| `0027:MVV` | 2 (PM, QA) | PM: no step asserts the promise text. QA: the `oracle` mini-check's strengthened control lives here, not in the scenario it fixes. |
| `0027:§testing-strategy` | 2 (PM, QA) | Both: scenarios assert category/detail, never the shipped description; S4/S3 oracles under-specified relative to the record's own bars. |
| `0027:§scope-verification` | 2 (PM, Impl) | Both: unfilled template placeholder. Stage 7's section — noted, not owed here. |

## Merged findings

Severity as reported by the raising persona. "Personas" names who raised it;
a hotspot id is marked.

### H1 — promise-wording text is the deliverable and nothing verifies it ships
Personas: PM (moderate-high). Anchors: `0027:C1` (hotspot), `0027:§phase-3-promise-wording`,
`0027:§testing-strategy` (hotspot), `0027:MVV` (hotspot).

C1's `promise:` clause says the lint's user-facing description states the
out-of-scope forms in C1's words, and Phase 3 creates that description (none
exists in shipped code). Every scenario and MVV step asserts the category string
or the refusal detail; none asserts the description text. The record could ship a
correct predicate with a missing or stale description and pass every listed test
— which is precisely the promise-wording failure §decision-rationale's premortem
names and claims mitigated.

Blocks: whether the validation plan is sufficient to call the user outcome
delivered at Stage 8.

### H2 — S4's stated oracle is weaker than the record's own prescribed control
Personas: QA (medium). Anchors: `0027:S4`, `0027:MVV` (hotspot).

S4's Expected is "every one lints GREEN" — an absence-of-error assertion a no-op
implementation satisfies. The MVV `oracle` mini-check names that exact row as
weak and prescribes the fix (assert green AND the binding loads with argv
unchanged), but fires from a different top-level section and is never folded back
into S4's own Expected line, where a test author reads.

Blocks: writing S4 as a test that discriminates a correct predicate from a stub.

### H3 — S3 omits the detail-text bar that C1's `report:` clause and S1 both set
Personas: QA (low). Anchors: `0027:S3`, `0027:C1` (hotspot), `0027:S1`.

S3's Expected is "all refuse `command_shell_interpreter`" — category only. C1's
`report:` clause requires the detail to name the two matched words, and S1
operationalizes it. The record's own `oracle` mini-check flags a bare category
assertion as insufficient for this contract.

Blocks: writing S3 with a pass/fail bar consistent with the contract it exercises.

### H4 — stdin-fed successor has no forcing function
Personas: PM (moderate). Anchors: `0027:A5`, `0027:F2`, `0027:§background`.

`env -i sh -c` was verified to actually execute a shell. C1 closes the wrapper
axis and admits the stdin-fed axis to a charted successor; A5's own durability
caveat records that no kata tracks it — it is a `Charted.md` line. Until it ships
the reviewer reads argv, which is the Problem Statement's opening complaint.

Blocks: whether landing C1 alone is an accepted interim state or should be paired
with tracking the successor.

### H5 — tie-break message attribution unexercised for two chained interpreters
Personas: Implementer (minor). Anchors: `0027:D-selection-predicate`, `0027:C1` (hotspot).

`["python","sh","-c","echo"]`: lowest-i-then-lowest-j reports `python -c`, where a
human reads `sh -c` as operative. Refusal correctness is unaffected — both words
are listed, the binding refuses either way — but no scenario or spike exercises a
two-distinct-interpreter argv, so the diagnostic wording has no worked example.

Blocks: nothing load-bearing; leaves the detail text for this shape unguided.

### H6 — deny-list membership has no forward pointer
Personas: QA (low/informational). Anchors: `0027:C1` (hotspot), `0027:S4`, `0027:S5`.

0027 scopes itself to argv position, not list membership (owned by 0025:C5) —
correct. But S4's `trace` row and S5 lean on "python is listed" / "ruby is listed"
as load-bearing prose premises, and 0027 carries no pointer to where the list is
read.

Blocks: nothing; a test author hunts for the list.

### N1 — Finalization Gate sections are unfilled template text
Personas: PM (low), Implementer (informational). Anchor: `0027:§scope-verification` (hotspot),
plus §contradiction-check, §assumption-verification, §proportionality.

Both personas flagged it and both correctly guessed it is sequencing: these are
Stage 7's sections, written at finalize. Baseline lint reports them as
`placeholder:survived`. Not owed by this stage — recorded so the hotspot count is
honest.

### N2 — implementer's empirical confirmation (not a defect)
Personas: Implementer. Anchors: `0027:C1` (hotspot), `0027:§existing-infrastructure-audit`.

Probe tests against the shipped `interpreterForm` confirm `nice`/`timeout`/`xargs`
prefixes currently lint green — i.e. the predicate is not yet position-free. The
record says exactly this ("replace the body, same name, same signature") and is
`Draft`. Confirmatory: the contract is actionable as written.
