Model: claude-fable-5

# Spike — A9 / A4 (Stage 4 re-entry, 2026-08-23): what Final 0002's normalizer does with the reserved key

Supersedes `a9-normalization.md` (2026-08-11), which re-ran RDR 0002's spike
as it stood before 0002's Stage 4 rebuilt it to the dump contract
(`8cf1dbc`, 2026-08-22) and before 0002 re-locked (`5bed32a`, 2026-08-23).
Every run below is against the spike **as committed at HEAD**
(`docs/rdr/0002-transition-table-as-reviewable-data/evidence/spikes/main.go`).

## Setup

```sh
cd docs/rdr/0002-transition-table-as-reviewable-data/evidence/spikes
export GOCACHE=<scratch>/gocache GOFLAGS=-mod=mod
```

## 1. Positive — 0002's two canonical fixtures (A4 first clause, A9 lifting)

```sh
go run . rdr-fixture.toml kata-fixture.toml > pos.out; diff output.txt pos.out
```

Result: exit 0, **byte-identical** to `output.txt` (0002's normative fixture
for its Testing Strategy scenario 2). Verbatim:

```
MODEL rdr rows=5 outcomes=round-clean,verdict-flapping,reconcile-block,finalized
rdr.continue-prelock kind=transition source=rdr:prelock outcome=round-clean atoms=[finalized_at.exists=false@all; iter.lt=3@all; profile.in=foundational large@match; profile.eq=small@unless; stage.eq=prelock@match; status.eq=Draft@match] next=[iter=2; stage=prelock] write=[iter=2; stage=prelock] requires_owned=[iter,stage] escape=[]
rdr.draft-no-match-escape kind=escape source=rdr:escape outcome=round-clean atoms=[status.eq=Draft@match] next=[] write=[] requires_owned=[] escape=[no_match]
rdr.reconcile-rewind kind=transition source=rdr:reconcile outcome=reconcile-block atoms=[status.eq=Draft@match] next=[prelock_lens=<clear>; rewind_scope=assumptions; stage=resolve; status=Draft] write=[prelock_lens=<clear>; rewind_scope=assumptions; stage=resolve; status=Draft] requires_owned=[prelock_lens,rewind_scope,stage,status] escape=[]
rdr.terminal-archive#finalized kind=transition source=rdr:terminal outcome=finalized atoms=[status.eq=Draft@match] next=[stage=archive] write=[stage=archive] requires_owned=[stage] escape=[]
rdr.terminal-archive#verdict-flapping kind=transition source=rdr:terminal outcome=verdict-flapping atoms=[status.eq=Draft@match] next=[stage=archive] write=[stage=archive] requires_owned=[stage] escape=[]
MODEL kata rows=2 outcomes=accepted,needs-work,closed
kata.review-accepted kind=transition source=kata:review outcome=accepted atoms=[owner.eq=current-session@all; phase.eq=review@match; status.eq=open@match; status.eq=closed@unless] next=[phase=ship; status=accepted] write=[phase=ship; status=accepted] requires_owned=[phase,status] escape=[]
kata.review-needs-work kind=transition source=kata:review outcome=needs-work atoms=[phase.eq=review@match; status.eq=open@match] next=[phase=resolve; status=open] write=[phase=resolve; status=open] requires_owned=[phase,status] escape=[]
```

Reading: every row carries `outcome=<literal>`; no `atoms=[…]` list contains a
`recognized` atom. `main.go::Row` now has `Outcome string`; `::normalize` calls
`::liftOutcome`, which deletes the `recognized` atom from the predicate map and
returns its literal(s). Both fixtures declare `[tags.recognized]` and every
predicate site is `[rule.match.recognized]` — the rename JDR 0001 §JD-10
ordered is on disk. **A9's 2026-08-11 verdict ("no outcome field; all
positions view-read") is refuted by the current artifact; A4's three-violation
census is stale.**

## 2. Mutants — guard-position `recognized` atom (0002 fenced: MUST be refused at load)

```sh
sed 's/^\[rule\.match\.recognized\]$/[rule.guard.all.recognized]/'    rdr-fixture.toml > mut-guard-all.toml
sed 's/^\[rule\.match\.recognized\]$/[rule.guard.unless.recognized]/' rdr-fixture.toml > mut-guard-unless.toml
go run . mut-guard-all.toml; go run . mut-guard-unless.toml
```

Result: **both exit 0** and emit the same five rows as §1 (each rule's
`recognized` atom lifted into `outcome=`, absent from `atoms=[…]`; the
`@all` / `@unless` block tag never appears because the atom was removed
before rendering).

Reading: the spike witnesses the *lift*, not the guard-position *refusal*.
0002's fenced text — "A `recognized` atom authored under `guard.all` or
`guard.unless` MUST be refused at load as a malformed outcome binding — never
lifted" — is the contract; `main.go::liftOutcome` scans the merged atom map
without consulting `Atom.Block`, so it lifts from any block. **0002 already
books this**: its Testing Strategy scenario 3 lists "a `recognized` atom
authored under `guard.all` and one under `guard.unless`" among the mutants
owed, and records that the spike witnesses nine categories, "the remainder
… owed at implementation." Not a defect to route; 0008 restates the fenced
rule (block 2), not the spike's behavior.

## 3. Mutant — two recognized-provenance declarations (0008 scenario 5 fixture)

```sh
sed 's/^\[tags\.recognized\]$/[tags.outcome]/' rdr-fixture.toml \
  | sed 's/^\[tags\.finalized_at\]$/[tags.result]\nprovenance = "recognized"\nkind = "enum"\n\n[tags.finalized_at]/' \
  > mut-two-decls.toml
grep -c 'provenance = "recognized"' mut-two-decls.toml   # → 2
go run . mut-two-decls.toml
```

Result (verbatim, exit 1):

```
refused: reserved_tag_key: recognized declaration named "outcome"
```

Reading: **exactly one** refusal, category `reserved_tag_key`, naming one of
the two declarations. Which of `outcome` / `result` is named depends on Go map
iteration order in `main.go::validate` — unspecified, as 0002's fail-fast
clause says ("a document tripping two categories MAY be refused with
either"). This is the normative fixture for 0008 scenario 5 as restated:
one refusal by category, never two.

## 4. Mutant — two independent categories in one document (fail-fast)

```sh
sed 's/^\[tags\.recognized\]$/[tags.outcome]/; s/^id = "terminal-archive"$/id = "continue-prelock"/' \
  rdr-fixture.toml > mut-two-defects.toml
go run . mut-two-defects.toml
```

Result (verbatim, exit 1):

```
refused: reserved_tag_key: recognized declaration named "outcome"
```

Reading: a misnamed recognized declaration **and** a duplicate rule id yield
one refusal. Corroborates "Load is fail-fast: the first category a document
trips is the refusal" and "an accumulating loader that returns a list is a
different contract than this one."

## 5. RDR 0003's guard fixture (A4 remaining declaration)

```sh
go run . ../../../0003-guard-predicate-exhaustiveness/evidence/spikes/guard-fixture.toml
```

Result (verbatim, exit 1):

```
refused: missing recognized outcome alphabet
```

Reading: `guard-fixture.toml` declares `[tags.rewind_target]` with
`provenance = "recognized"` (`:36-37`) and references it at
`[rule.guard.all.rewind_target]` (`:72`), but carries no root `outcomes` and
no `[rule.match.recognized]` on any rule. It is not a 0002 model; it is
refused on the first category tripped, and would additionally trip
`reserved_tag_key` (wrong name), zero-outcome binding, and — under 0002's
fenced text — the guard-position refusal. Under Final 0002 it is a load
failure on several grounds, **not a rename**, and it is RDR 0003's artifact.
No Phase 2 rename work remains for 0008.
