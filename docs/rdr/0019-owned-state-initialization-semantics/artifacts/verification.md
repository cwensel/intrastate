
## Phase 3b — adversarial review (`0019:F1`–`F5`)

Independent pass over the record's `Failure Modes`, written without reading
Phase 3a's findings. Tests live in `internal/cli/flow_initstate_adv_0019_test.go`.

### ADV-1 — F2, the torn state the VERB ITSELF produces — PASSES

**Failure mode.** `0019:F2` "Torn (multi-writer)": "no cross-writer atomicity
— a failed writer can leave the artifact seeded for some keys only. The store
is then non-empty, so a re-run is a NO-OP whose payload lists the missing
`[initial]` keys … but it does NOT repair them." (REQ-47, REQ-51, REQ-60,
REQ-67.)

**The gap in the shipped suite.** S12
(`TestReq51And93_0019_ATornStateIsReportedNotRepaired`) constructs its torn
state with a prior `flow set-state`. No shipped test drives `init-state`'s own
write loop to a partial commit, so F2's central claim was asserted only over a
state the verb did not create.

**The test.** `TestAdv1_0019_ATornSeedTheVerbItselfProducesIsReportedNotRepaired`.
A two-role model whose `aux` writer — first in the loop's sorted name order —
declares an unreachable read-back locator. The seeding arm commits `note` into
a sealed artifact, refuses at exit 3, and never reaches `state`. Asserts: the
partial commit is real and the unreached role's artifact does not exist; the
re-run is a no-op success leaving both artifacts byte-identical; the re-run's
payload names exactly `{stage}`.

**Result: PASSES.** The implementation defends F2's torn case. Kept for
regression value the shipped suite lacks — it is the only test whose torn
state the verb itself produced.

### ADV-2 — F2, the absent report is answered from the UNION of stores — FAILS

**Failure mode.** `0019:F2`: "The payload's absent-key list is what keeps the
manual route mechanical — it names precisely the keys to pass." REQ-47 fixes
that the no-op payload "reports which `[initial]` keys THE STORE does not
carry", and the record's repair route is "explicit `set-state` of the listed
keys" — which routes each key to its own writer's artifact.

**The defect.** `internal/cli/flow_initstate.go::initAbsentKeys` folds every
role's key set into one `present` map before testing membership:

```go
present := map[string]bool{}
for _, keys := range stores {
	for _, key := range keys {
		present[key] = true
	}
}
```

A key carried by ANY bound store suppresses its absent entry, regardless of
which writer serves it. Every other quantifier in the verb is role-scoped
(`initNeededRoles`, `initStoresEmpty` over per-role key sets); this one is not.

**The test.**
`TestAdv2_0019_TheAbsentReportIsAnsweredPerRoleNotFromTheUnionOfStores`.
Built through the CLI only, with no hand-staged artifact: both roles of
`initTwoRoleModel` are bound to one shared artifact and `note` is established
there by `set-state`; role B is then rebound to a fresh artifact (a
per-invocation binding change no clause forbids). Role A's store now carries
`note` — a key only role B's writer serves — while role B's store is empty.
The correct absent set is `{note, stage}`; the implementation reports
`{stage}`.

**Result: FAILS.** Observed payload: `absent: ["stage"]`.

**Harm, probed separately and confirmed.** An operator who runs exactly and
only what the payload lists — `flow set-state --write stage=seeded` — is left
with `note` still reading ABSENT from `flow read-state`. F2's claim that the
list "names precisely the keys to pass" is false on a multi-role model.

**Suggested repair (a later step's job; no production code touched here).**
Answer `initAbsentKeys` per role: for each `[initial]` assignment, test its key
against the key set of the store its own writer's role names, rather than
against the union.

### Discarded probe — F1 exit-group confusion — PASSED, redundant

"An exit-3 seed leaves an APPLIED store the re-run will not reseed" (REQ-67,
F1's "refused before writing" vs "applied but unverified"). Written, PASSED,
and discarded rather than kept: it is
`TestReq67_0019_AnIncompleteReadBackExitsThreeAndLeavesANonEmptyStore` under a
different name. The implementation defends this mode; a redundant copy adds no
regression value. ADV-1 asserts the same disposition incidentally, on the
harder multi-writer construction.

---

## Phase 3a — CoVe verification (execution against the shipped binary)

Independent pass over `<art>/req-list.md`'s 103 REQs + REQ-MVV, written
without reading the Phase 1 test files. Method: for each REQ, name an input
that would make a correct implementation visibly violate it, then RUN that
input against `bin/intrastate` over hand-built fixtures. Findings are recorded
only where a violation was ACTUALLY REPRODUCED; a suspicion is not a FAIL.

### FAIL-1 — REQ-15, the verb cardinal in `docs/model-authoring.md` is stale

**The obligation.** REQ-15: "Six code/doc sites hard-code the cardinal and
change with it: `internal/cli/flow.go`'s package comment, `::newFlowCmd`'s
`Long` body, `::flowExtendedDesc`, `internal/cli/flow_exec.go`'s header
comment, `docs/cli-reference.md` (two strings), and `docs/model-authoring.md`."
(IP, Phase 1 Step 1.)

**The failing input.**

```
grep -n 'four `flow` verbs' docs/model-authoring.md
```

**Observed.** `docs/model-authoring.md:15`:

> analysis, and both are driven by the same four `flow` verbs. The class is
> declared in `[model]`, and it decides which invariants apply.

**What the REQ requires.** The cardinal changes with the verb set.
`init-state` is registered — `intrastate flow --help` lists five
`Available Commands` — so this site must read five, as the other five cardinal
sites now do:

- `internal/cli/flow.go:5` — "One group, five verbs" — UPDATED
- `internal/cli/flow.go:44` — "The five verbs are the whole..." — UPDATED
- `internal/cli/flow.go:111` — "The five verbs divide one job..." — UPDATED
- `internal/cli/flow_exec.go:3` — reworded to carry no cardinal — SATISFIED
- `docs/cli-reference.md:181,220` — both read "The five verbs" — UPDATED
- `docs/model-authoring.md:15` — reads "four" — **NOT UPDATED**

**Corroboration that this is a defect and not a permitted reading.** The same
document contradicts itself. `docs/model-authoring.md:380-384` enumerates the
verbs and names five, including the one line 15 excludes:

> **The same verbs.** `flow next`, `flow resolve`, `flow read-state`,
> `flow set-state`, and `flow init-state` drive both classes — though a
> decision table has no state to read, set, or initialize, so only `next` and
> `resolve` are meaningful for it.

So the Phase 2 activation edits (REQ-95/96) DID land in this file — the
start-state sentence gained its runtime carrier at line 752, the file-backed
scope is stated at line 768, and `flow next`'s `unknown[].reason: absent`
report is named at line 762 — while the Phase 1 Step 1 cardinal edit at line
15 was missed. The reader who meets the cardinal first is told there are four.

**Suggested repair (a later step's job; nothing edited here).** Change
`docs/model-authoring.md:15` "the same four `flow` verbs" to name five.

### FAIL-2 — REQ-47/REQ-51, the absent report is answered from the UNION of stores

Reproduced INDEPENDENTLY of Phase 3b's ADV-2, through the CLI rather than the
test harness. The finding is ADV-2's; this is corroboration on a second route,
not a second defect.

**The failing input.** Two roles, `ra` (whose writer serves `alpha`) and `rb`
(whose writer serves `beta`), both keys in `[initial]`. Role A's store is made
to carry `beta` — a key only role B's writer serves — while role B's store is
empty:

```
$ printf '{"beta":"B"}' > /tmp/U1.json          # role ra's artifact
$ rm -f /tmp/U2.json                            # role rb's artifact: empty
$ intrastate flow init-state --model two.toml \
    --artifact ra=/tmp/U1.json --artifact rb=/tmp/U2.json --as json
```

**Observed.**

```json
{"type":"ok","data":{"...":"...","seeded":[],"absent":["alpha"]}}
```

**What the REQ requires.** REQ-47: the payload "reports which `[initial]` keys
the store does not carry". `beta`'s store is role `rb`'s, which is EMPTY, so
`beta` is absent and must be listed. The correct report is `{alpha, beta}`.
`::initAbsentKeys` folds every role's key set into one `present` map, so role
A's unrelated `beta` suppresses role B's genuine absence. REQ-51's repair
route — "explicit `set-state` of the listed keys" — is therefore incomplete:
an operator who writes exactly what the payload lists leaves `beta` absent.

---

### Probes that HELD (no finding)

Run against `bin/intrastate` over fixtures in a scratch dir. Each names the
input that would have exposed a violation.

**Empty-store predicate (REQ-24/25/26/27/28/29/30/31).**

- ALL quantifier, cross-artifact: role A carrying a key with role B empty is a
  NO-OP at exit 0 with ZERO writes to BOTH — role A's bytes unchanged, role
  B's artifact never created. A per-artifact or ANY reading would have seeded
  into B. HELD.
- SEALED store: an artifact whose only key is `flowbind::sealedKey` (the
  NUL-prefixed `flow.readback-unreachable` marker, byte-verified with `od -c`
  before the run) is a NO-OP SUCCESS at exit 0, bytes unchanged — not the
  exit-3 refusal `::Reader.Read`'s unreadable short-circuit would produce. The
  count is over STORE keys: `absent` correctly reports `["note","stage"]`
  while the store is non-empty. `flowbind.StoreKeys` reads `len(load(path))`
  and does not enter the seal branch, as A6 requires. HELD.
- Substituting "every `[initial]` key reads absent" (REQ-30) would have seeded
  over the seal; it does not. HELD.
- An unbound needed role is a REFUSAL, never an artifact treated as empty
  (REQ-31): `flow-artifact-missing` at exit 2 even when the OTHER bound store
  is non-empty, so the no-op arm is unreachable from there. HELD.

**Cleared-key tombstones (REQ-48/49/74, MVV.6/7/9).**

- Seed `{note,stage}`, `set-state --clear note`, re-run `init-state`: exit 0,
  ZERO writes (md5 unchanged), `absent:["note"]`, and `read-state` still
  reports `note` absent. The key does NOT resurrect. HELD.
- Clear the LAST key so the store is `{}`, then re-run: RESEEDS both keys.
  REQ-49's accepted residual, asserted in the failing direction. HELD.
- Cross-artifact: clear `alpha` in role A while role B still carries `beta` —
  the re-run is a no-op and `alpha` stays absent. This is REQ-26's per-model
  guarantee, the property the ALL quantifier exists to protect. HELD.

**Carrier discrimination, BOTH sides (REQ-38/39/40/41/42/90/91).**

- Write side, edit-carried (`[write.state.edit.stage]`): exit 2,
  `flow-init-carrier-unsupported`, message names the WRITE capability and the
  accessor, artifact never created. HELD.
- Write side, command-backed: same code with `--allow-commands` UNSET **and**
  SET — never the allow-commands refusal. REQ-36's preemption. HELD.
- Read side, file-backed `[write.state]` paired with command-backed
  `[read.state]` on the same role: exit 2 with the carrier code naming the
  READ capability. A write-side-only type-switch would have ADMITTED this and
  seeded. HELD.
- POINTER/VALUE spelling (REQ-41): the file-backed control seeds at exit 0, so
  the read case is not written `*flowbind.Reader` — that spelling would match
  nothing and refuse every model, this control included. HELD.

**Registry first-match reader selection (REQ-43/44/46).**

- Two readers on one role serving different keys. Named `aaa`(file) /
  `zzz`(command): ADMITS and seeds. Names swapped so `aaa` is the
  command-backed one: REFUSES with the carrier code naming `aaa`. The
  discrimination flips with registry slice order exactly as `::readerFor`
  selects, so nothing sorts or filters before scanning. HELD.
- A role with no reader at all reaches the unbound-needed-role arm
  (`flow-artifact-missing`), not a separate missing-reader code (REQ-44). HELD.

**Decision-table class exclusion (REQ-4/5/6/7/35, MVV.8, S7).**

- A loading decision-table model (`class = "decision-table"`, no `[initial]`,
  no owned tags): exit 2, `flow-init-class-unsupported`, artifact never
  created. HELD.
- With NO `--artifact` bound at all: still the CLASS code, so class precedes
  role binding. HELD.
- With a command-backed read accessor declared: still the CLASS code, never
  the carrier code — REQ-35's ordering, answerable from the loaded model
  before the registry is consulted. HELD.

**Gate ordering (REQ-32/33/34/37).**

- Non-file-backed model with NO `--artifact`: the carrier code, not the
  artifact-binding family — carrier precedes role binding (REQ-33). Verified
  on all three carrier shapes (edit write, command write, command read). HELD.
- Writer arity is settled at LOAD (REQ-34/61/84): a model routing one
  `[initial]` key to two writers refuses `flow-model-invalid` with NO
  artifacts created — never the writer-routing family. HELD.

**Seed encoding (REQ-17..23, RT3, S4).**

- Nine kinds both routes admit (enum, scalar, bool, int, decimal-as-scalar,
  sorted set, set with duplicates and `<`/`&`, empty set, single-member set)
  seeded into a fresh artifact vs. the `set-state --write` transcription into
  a second fresh artifact: **byte-identical** (`cmp` clean). Sets are sorted,
  deduplicated, compact, HTML escaping off; scalars are `members[0]` verbatim
  — never `["draft"]`. HELD.
- Three loader-only spellings (REQ-23/53/72/83): `labels = "bare"` (bare
  scalar for a set tag) seeds `["bare"]`; `status = ["single"]` (array for a
  scalar tag) seeds `single`; `note = ""` seeds `{"note":""}` and reads back
  PRESENT with an empty value. All three seed at exit 0 — no re-conform pass
  refuses at seed time what the loader admits. HELD.
- The empty scalar is distinct from a cleared key: after `--clear note` the
  key leaves the object entirely. HELD.
- REQ-55's asymmetry is preserved untouched: `set-state --write note=` still
  refuses `flow-write-invalid`, "the tag `note` was given an empty value". HELD.
- Numeric boundary (REQ-73/82): `[initial] threshold = 1.0` seeds `"1"` while
  `--write threshold=1.0` writes `"1.0"`; `--write n=+5` writes `"+5"` where
  the loader seeds `"5"`. Both read back their own value. The divergence is
  pinned as decided behaviour, not repaired. HELD.

**Read purity / never-fill (REQ-2/3, S13).**

- A bound artifact with `stage` present and `note` never seeded: `read-state`
  reports only `stage`; `flow next` reports `owned:{stage}` with no synthesis;
  `resolve` refuses `flow-no-match` over the assembled state. A fully empty
  artifact makes `next` report `unknown:[{"key":"stage","reason":"absent"}]`.
  No read path returns the `[initial]` value. Control: the same fixture with
  `note` PRESENT reads back its stored value. HELD.

**Payload content (REQ-56/57/58).**

- Seeding arm: `seeded:["note","stage"]`, `absent:[]`. No-op arm: `seeded:[]`,
  `absent:["note"]`. Keys by NAME only — no value is echoed on either arm, so
  `read-state` stays the single surface reporting values. The absent list is
  empty on the seeded arm as REQ-58 requires. HELD.

**Read-back and exit groups (REQ-63/67, F1).**

- A write accessor declaring an unreachable read-back locator seeds, then
  refuses `flow-write-readback-incomplete` at **exit 3**, leaving a store
  carrying the seal plus the written keys — non-empty, so the re-run is the
  no-op that does not repair. Exit 3 for could-not-complete vs. exit 2 for
  completed-and-disagreed is the group split REQ-67 fixes. HELD.

**Refusal codes and exit groups (REQ-64/65/66).**

- `flow-init-class-unsupported` and `flow-init-carrier-unsupported` are
  DISTINCT codes, both at exit 2, and distinct from the shared classes
  (`flow-artifact-missing`, `flow-model-invalid`, `flow-write-unbound`), which
  keep their shipped groups. Asserted on group and distinctness only — the
  spellings are non-normative per REQ-65 / ASSUMPTION-3 and nothing here pins
  them. HELD.

**Verb surface (REQ-8/9/10/11/13/14).**

- `intrastate flow --help` lists exactly `{init-state, next, read-state,
  resolve, set-state}` — set equality with `init-state` added. HELD.
- No write grammar: `--write stage=x` is `command-error: unknown flag`. HELD.
- `--allow-commands` is INHERITED, not declared: accepted on this verb and
  inert (the run succeeds unchanged). HELD.
- `respond.ValidateMode` runs first: `--as bogus` refuses `flag-invalid-value`
  before any model work. HELD.

**MVV spine (REQ-MVV.1-.10).** All ten steps run green end to end over a
two-owned-key fixture, including MVV.10 against the SHIPPED `models/rdr.toml`:
bound to a fresh artifact it seeds exactly `{gate_passed, stage, status}` at
their declared values
(`{"gate_passed":"false","stage":"seeded","status":"draft"}`), `flow next`
then reports NO `unknown[].reason: absent` for any of them, and the re-run is
byte-identical (RT2). HELD.

**ASSUMPTION-4 ("needed" role scoping).** A command-backed write accessor
serving no `[initial]` key does NOT trigger the carrier gate, and its role may
go unbound without refusal — the reading the req-list records. Probed for harm
against REQ-26: a bound-but-not-needed artifact carrying a key does not block
the seed, but every case reachable that way is one REQ-49 already licenses
(the needed role's own store was emptied), so no cleared-key guarantee is
lost. Recorded as consistent with the audit, not as a finding.

### REQs not verifiable by execution

Recorded rather than converted into FAILs.

- **REQ-5's "UNCONSTRUCTIBLE" state** (a decision-table declaring `[initial]`)
  cannot be built — the loader refuses it — which is the claim itself. The
  reachable half (the refusal keys on the class, not on `len(owned)`) is
  asserted by the state-machine-with-empty-`[initial]` control, which lint
  blocks first (REQ-6, 0006:C18).
- **REQ-42** (the `table.Accessor` re-derivation ban) and **REQ-45** (the gate
  invokes no `Binding` method) are source-shape obligations with no
  behavioural discriminator distinguishable from the type-switch that ships.
  Read in `flow_initstate.go::initWriteCarrier` / `::initReadCarrier` /
  `::initCarrierGate`: both hold as written.
- **REQ-16** (`llms.txt` regenerated) verified by inspection, not execution:
  `llms.txt:24` carries `intrastate flow init-state`.
- **REQ-59/68/98/99/100/101/102/103** are non-goals, prerequisites, or
  taxonomy-registration obligations with no runtime surface.
- **REQ-87 (S8)** — a read accessor bound to a DIFFERENT artifact path than
  its writer — is not constructible through the CLI on the file-backed
  carrier: `--artifact` binds per ROLE, and the read-back reader is selected
  by the writer's own role, so both necessarily resolve to the same path. A
  role with no reader reaches `flow-artifact-missing` (REQ-44), a different
  arm. The mismatch class is reachable only from the test harness.
