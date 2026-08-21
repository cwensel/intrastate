Model: claude-opus-5[1m]

# Grounding Findings — iter-2 (re-entry, delta-scoped to A8)

Scope: claims added or edited after the Final commit `c149031` — A8, the
read-completeness normative clause, the Read completeness invariant, MVV
Scenario 2, and the Failure Modes / Assumption Verification edits that carry
them. A1-A7 were swept in iter-1 and are carried forward.

## Findings

- **NOT-FOUND — A8 Evidence cites three spike symbols that do not resolve.**
  The A8 Evidence line cites
  `.../evidence/spikes/main.go::completeRead`,
  `.../evidence/spikes/main.go::truncatedRead`, and
  `.../evidence/spikes/main.go::absentKeyRead`.
  None of the three exists in `main.go`. The file declares no `completeRead`,
  `truncatedRead`, or `absentKeyRead` symbol at any scope. The three read
  dispositions are real, but they are produced by differently-named fixture
  builders: `main.go::newArtifacts` (complete read),
  `main.go::newPartialArtifacts` (unreadable `profile`, drives
  `incomplete_read`), and `main.go::newSparseArtifacts` (artifact genuinely
  lacks `profile`, resolves to `<absent>`). The cited-symbol form is what the
  `Method: Spike` / RDR anchor doctrine requires be greppable, and these are
  not.

## Confirmed codebase claims (delta scope)

- `internal/resolve/resolve.go::missingOwned` — CONFIRMED. Present at
  `resolve.go`; it walks each candidate's `RequiresOwned` and accumulates a key
  when `view.has(key, ProvenanceOwned)` is false.
- `internal/resolve/resolve.go::TagSet.has` — CONFIRMED. Declared as
  `func (s TagSet) has(key string, prov Provenance) bool`, whose body is
  `tv, ok := s.tags[key]; return ok && tv.provenance == prov` — map presence
  plus provenance, with no error or unreadability channel.
- **The downstream-irrecoverability claim** ("an absent and an unread owned key
  both become `owned_state_unavailable`") — CONFIRMED at the call flow, not
  merely at the named symbol. `TagSet` carries `map[string]taggedValue` holding
  only `value` + `provenance`; `Input.Owned` supplies it through `assemble`
  with no error channel, so a key omitted by a truncated read and a key the
  artifact lacks are byte-identical inputs to `missingOwned`. The refusal name
  is confirmed at `internal/resolve/fixtures_test.go:312`, which maps
  `"owned_state_unavailable"` to `missingOwnedInput`.
- `.../evidence/spikes/main.go::read` — CONFIRMED. Takes
  `requested ...string`, and for each requested key returns
  `refusalIncompleteRead` when `art.unreadable[key]`, else substitutes
  `absentValue` when the artifact does not carry it — the branch split A8
  asserts.
- `.../evidence/spikes/main.go::absentValue` — CONFIRMED. `const absentValue =
  "<absent>"`.
- Transcript claims `output.txt:11-13` — CONFIRMED. Line 11
  `tags={profile=large,status=Draft}` (complete), line 12
  `refusal=incomplete_read`, line 13 `tags={profile=<absent>,status=Draft}`.
  The renumbering of the pre-existing A4 (`output.txt:14-15`) and A6
  (`output.txt:16-23`) cites is also CONFIRMED against the current transcript.
- **JDR 0001 §D3** — CONFIRMED. `docs/jdr/0001-resolve-kernel-seam.md:141`
  is `## D3 — Must a read accessor return a *complete* tag set?`, resolving
  **(b)**: "A read accessor MUST return the complete tag set for the keys it
  was asked for, or take the refusal branch." The RDR's normative clause states
  the same rule, and §D3's landing instruction ("one normative clause plus an
  MVV scenario asserting a truncated read refuses") is satisfied by the added
  clause and MVV Scenario 2.

## Inverse search — does a sibling path already make this decision?

Searched `internal/` and `cmd/` for `incomplete`, `truncat`, `partial`,
`unreadable`, and `accessor`/`Accessor` across production `.go` files. **None
exists.** Every `accessor` hit in production code is a doc comment in
`internal/resolve/resolve.go` describing the snapshot's provenance; there is no
accessor executor, no read-completeness predicate, and no partial-read
discriminator anywhere outside the spike. The new rule does not duplicate an
existing sibling decision — it binds at a boundary production code has not yet
built.
