Model: claude-opus-5[1m]

# Persona 3 — QA / Tester (iteration 2)

Delta-scope re-run. Checked items 1–8 of the brief against
`internal/cli/flowbind/flowbind.go`, `internal/cli/flowbind/registry.go`,
`internal/accessor/executor.go`, `internal/accessor/model.go`,
`internal/table/load.go`, `internal/cli/flow_state.go`,
`internal/cli/flow_input.go`, and `models/rdr.toml`.

Three findings. Two are new gaps the rewrite opened (S8, S9); one is a
pre-existing mismatch that the rewrite's own new disposition row now makes
self-contradictory (S5).

---

## HIGH — `0019:S8`: the new fixture is unconstructible; reader and writer cannot be bound to different artifact paths

**Test it prevents:** the read-back-mismatch scenario — the only test
distinguishing `ClassReadBackMismatch` from `ClassReadBackIncomplete`, and the
only test of C1's "a re-run does NOT repair a present-and-wrong key".

The rewrite replaced the previously-unconstructible fixture with: "the fixture
binds the READ accessor to a different artifact path than its writer,
pre-seeded with a conflicting value for the same key." That is also
unconstructible, for two independent reasons in `internal/accessor/executor.go`:

1. The read-back reader is selected BY ROLE, not by accessor:
   `reader, hasReader := e.Registry.readerFor(def.Accessor.Role)`
   (`executor.go:311`), and `readerFor` matches `d.Accessor.Role == role`
   (`internal/accessor/model.go:249-256`). A read accessor on a *different*
   role is simply never consulted for this writer's read-back.
2. Even a same-role reader gets the writer's artifact, not its own:
   `def, art, timeout, refusal := e.selects(name, CapWrite)` resolves
   `art = e.Artifacts[def.Accessor.Role]` once (`executor.go:59`), and that
   same `art` is handed to the read-back:
   `raw := e.invokeRead(ctx, reader, art, readTimeout, compared)`
   (`executor.go:517`). Artifact paths are bound per ROLE
   (`--artifact role=path`), so there is exactly one path per role.

The accessor's own declared `path` cannot redirect it either: `flowbind.Reader`
uses `r.Path` only for the `unreachable(r.Path)` test and reads the store from
`load(art.Path)` (`flowbind.go:186-193`). So writer and read-back reader
provably share one store, which is exactly the situation the scenario's own
first sentence says makes a disagreement unconstructible.

The scenario's parenthetical fallback — a command-backed READER paired with a
file-backed writer — IS constructible and IS grounded (`registry.go:49-51`
builds `cmdbind.Reader` on `commandBacked(acc)`; C1's carrier refusal is scoped
to write accessors, confirmed in C1's CARRIER SCOPE paragraph). But it is
written as an aside, not as the fixture, and it is the only route the code
admits.

**Fix direction:** promote the command-backed-reader construction to the
scenario's stated fixture (with its `--allow-commands` and helper-binary cost
acknowledged), or drop the different-path claim. As written, the primary
fixture cannot be built, so no test can be written from S8.

---

## HIGH — `0019:S9`: the "second role" arrangement does not produce a store whose sole key is the seal, and the scenario loses its discriminating control

**Test it prevents:** the pin that the emptiness predicate counts STORE keys
rather than OWNED keys — S9's stated purpose ("a future change that counted
owned keys instead would seed beneath an unverifiable state and fail here").

The rewrite's fixture: clear every owned key through the ordinary reachable
writer FIRST (role A's store → `{}`), then land the sealing write LAST on a
SECOND ROLE "whose own key is outside the cleared set", claimed to yield "a
store whose sole key is the seal."

It does not. Each role has its own artifact path, and `Writer.Apply` operates on
`art.Path` (`flowbind.go:252`), so role B is a *different store*. `Apply`'s
unreachable arm sets `s[sealedKey] = sealedMarker` only after applying the
planned tags (`flowbind.go:250-270`), and a write with an empty plan never
reaches `Apply` at all — `Executor.Write` short-circuits on
`len(plan.Writes) == 0` (`executor.go:277-281`). So the sealing write on role B
must carry at least one planned tag, and if that tag is B's own key (outside the
cleared set, i.e. assigned rather than cleared), B's store ends as
`{Bkey: v, sealedKey: "1"}` — TWO keys, one of them owned.

Consequences for the test:

- The stated end state ("a store whose sole key is the seal") is not reached by
  the stated construction. Nothing in the fixture description is assertable
  against a real store.
- The oracle survives only vacuously. The expected NO-OP holds because role B is
  non-empty — but it is non-empty because it carries an OWNED key, not because
  of the seal. So the discriminating control fails: a hypothetical predicate
  that counted owned keys would ALSO decline here, and this scenario would pass
  under the very bug it exists to catch. `0019:§mini-checks` `oracle` row
  "9 sealed artifact declined … Negative control: count owned keys only → seeds
  into sealed store" is therefore false against this fixture.

Note the route the rewrite RULED OUT does work. S9 rejects clearing through the
sealing writer on the ground that it "both re-seals and exits non-zero." The
non-zero exit is a read-back-incomplete refusal, but the mutation still lands —
`Apply` deletes the key and sets the seal in one atomic `save` (`flowbind.go:255-270`),
and `Reader.Read`'s short-circuit (`flowbind.go:200-208`) only affects the
verifying re-read, not the persisted bytes. A `set-state --clear <last owned key>`
through a `-unreachable`-suffixed writer therefore yields exactly `{sealedKey}`,
which IS the state S9 needs. A non-zero exit on a fixture SETUP step is not a
disqualifier.

**Fix direction:** either restore the clear-through-the-sealing-writer setup
(stating that the setup step's exit 3 is expected and the mutation lands), or
state the second-role construction as clearing B's own key through B's sealing
writer so B's store is `{seal}` alone. Either way the scenario must end with a
store carrying NO owned key, or the owned-key-count negative control is not run.

---

## MEDIUM — `0019:S5` vs `0019:C1` vs the new `disposition` writer-arity row: S5's stated refusal family is wrong, and the three no longer agree

**Test it prevents:** a correct writer-arity assertion. As written, S5 asserts
on the wrong refusal family, so the test passes for the wrong reason and cannot
detect a regression in the arm it names.

The rewrite changed the `disposition` table's writer-arity row to
"Role unbound (writer-arity is load-enforced, unreachable here)", matching C1,
which now says the writer-arity half "is already enforced at LOAD —
`internal/table/load.go::checkAccessorBindings` folds every `[initial]` key into
its `written` set" and that "the verb states the requirement … but owes no
runtime check and no fixture for it."

Grounded and correct: `checkAccessorBindings` does
`for _, t := range l.model.Initial { written[t.Key] = true }` (`load.go:1539-1541`)
and refuses `writerCount[key] != 1` as `CatMalformedAccessorBinding`. At the CLI
that surfaces as `codeModelInvalid = "flow-model-invalid"` via `loadFailure`
(`internal/cli/flow_input.go:47,194-203`).

But `0019:S5` was NOT updated and still reads: "`[initial]` names a key with
zero declared writers, and separately one with more than one. **Expected**:
refusal from the existing writer-routing family with ZERO writes committed."
The writer-routing family is `flow-write-unbound` / `flow-clear-unbound`
(`flow_input.go:42-43`), which this input never reaches — the model never loads.
The new both-arms assertion ("bytes/mtime unchanged, or the artifact path still
absent") is satisfied trivially for the same reason.

So the record now says in three places: C1 — no fixture is owed; the disposition
table — unreachable here; S5 — a scenario asserting a family the input cannot
reach. A tester writing S5 from the text writes a passing test that proves
nothing about the verb.

**Fix direction:** either delete S5's writer-arity halves (C1 says no fixture is
owed) or restate them as a LOAD refusal (`flow-model-invalid`, before the verb
runs), and reconcile with the disposition row's parenthetical.

---

## Items checked, no gap found

- **`0019:S11`** — constructible and correctly grounded. `Model.Writers` is
  `map[string]Accessor` (`internal/table/model.go:521`) and `runFlowSetState`
  iterates per writer, checking `req.artifacts[def.Accessor.Role]` for each
  (`internal/cli/flow_state.go:355-361`), so two write accessors over two roles
  is a real model shape. The oracle (both artifacts' bytes) and the negative
  control (switch to ANY/per-artifact → role B gains its key while every
  single-artifact row still passes) are both real and discriminating. This is
  the strongest new scenario.
- **`0019:S12`** — the F2 property is worth testing and the oracle (A's
  bytes/mtime unchanged; payload names exactly B's un-seeded keys) is real. Not
  raised as a finding because the un-stated setup mechanism is a fixture detail
  an implementer can resolve — the record does NOT say how "role B's write
  fails", and the obvious declarative lever does not work (a `-unreachable`
  writer APPLIES then seals, `flowbind.go:250-270`, so B would be seeded, not
  un-seeded). Constructing the torn state by writing role A's artifact directly
  during setup (never through the verb) is legal for a fixture and sufficient.
  Worth a sentence when someone next touches the scenario; not blocking.
- **`0019:S5`/`0019:S6` both-arms phrasing** — sound as an assertion form.
  S6's role-binding arm is genuinely reachable and correctly attributed: roles
  are bound per invocation and `runFlowSetState` refuses with
  `codeArtifactMissing` (`flow_state.go:355-361`). The finding above is against
  S5's family, not against the both-arms phrasing.
- **`0019:C1` non-normative code spellings (item 6)** — sound. A stated oracle
  exists: `clierr.Finding.Code` / the error envelope's `code` field are on the
  wire (`internal/cli/clierr/clierr.go:51,294`;
  `internal/cli/respond/respond.go:219`), so S7, S10 and MVV steps 7–8 can
  assert exit 2 plus pairwise inequality of the two new codes and their
  inequality from the shared classes, without pinning a literal. S10's
  command-backed row can specifically assert code ≠
  `codeRequestRefused = "flow-accessor-request-refused"`
  (`flow_input.go:96`), which is the ordering C1 fixes. The precedent cited
  (`0028:C1.3`'s `codeWriteEditRefused`) is real (`flow_input.go:81`).
- **`0019:MVV` step 10 (item 7)** — runnable and assertable, grounded against
  `models/rdr.toml`: no `[model] class` key (so `IsDecisionTable` is false),
  three owned tags `stage`/`status`/`gate_passed`, an `[initial]` block seeding
  exactly those three, and `[write.rdr-status]` carrying `path`/`keys`/
  `read_back` with no `command` and no `edit` — file-backed by
  `registry.go`'s fallthrough. The `flow next` assertion is real: the `unknown[]`
  reason vocabulary includes `absent` (`internal/cli/flow_next.go:86,338`).
  The claim that the three keys span two value kinds (enum, bool) is correct.
- **`0019:RT3` reconciliation with `0019:S4` (item 8)** — consistent across
  RT3, S4, and the `fidelity` table's third row: three loader-only kinds, two
  testable by read-back equality, the empty scalar unwritable (F5) and untested.
  Grounded: `loadInitial` admits a bare scalar for a set tag and an arity-1
  array for a scalar tag (`load.go:1688-1713`, the
  `decl.Kind != "set" && len(members) != 1` refusal), and `canonicalValue`
  refuses both on argv (`flow_input.go:718-738`).
- **`0019:§mini-checks` new `oracle` rows (item 5)** — the S11 row agrees with
  S11 as written. The S12 row agrees with S12 as written. The S9 row does NOT
  agree with the rewritten S9 fixture; that is folded into the S9 finding above
  rather than reported separately.
