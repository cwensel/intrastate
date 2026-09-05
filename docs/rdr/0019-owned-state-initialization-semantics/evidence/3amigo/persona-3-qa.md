Model: claude-opus-5[1m]

# Persona 3 — QA / Tester

Owned set read in full: `0019:C1`, `0019:MVV`, `0019:S1`–`0019:S10`.

**Widening.** Four things sent me outside the owned set. (1) Two scenarios
(`S9`, `S10`) name fixture machinery — `flowbind::sealedKey`, the edit and
command carriers — whose constructibility is decided by shipped code, not by
the scenario, so I grounded each against `internal/cli/flowbind/`. (2) `C1`
declares an ALL quantifier "load-bearing rather than stylistic" but no
scenario mentions a second artifact, so I read `0019:D-selection-predicate`
and `0019:§mini-checks` to see whether an oracle lived there instead. (3) The
scenarios cite invariants by number only, so I read `0019:RT1`–`RT4` to get
the actual equalities under test. (4) A missing scenario has no line range, so
I read `0019:F1`–`F5` to enumerate the failure modes the record names in prose
and diffed that list against `S1`–`S10`.

The record's `oracle`, `fidelity`, and `disposition` mini-check tables in
`0019:§mini-checks` are unusually strong — most scenarios already carry a
negative control and a stated discriminating assertion. The findings below are
what survives that.

---

## High

### 1. `0019:S9` — the setup is not constructible as written

**Claim.** S9's fixture ("write through a writer with an unreachable read-back
locator … then clear every owned key") cannot be built, because the seal is a
property of the writer's declared locator, not of the write, and clearing
re-runs that same writer.

**Evidence.** `internal/cli/flowbind/flowbind.go::Writer.Apply` sets
`s[sealedKey] = sealedMarker` under `if unreachable(w.Path)` and `delete`s it
in the `else` arm — the seal names the LAST write's locator, as the inline
comment says. `::unreachable` (same file, line 99) is a pure suffix test on the
accessor's declared `Path` string, so a given writer is either always-sealing
or never-sealing across every invocation. S9's clears must route through a
writer for those owned keys, and `0019:A3` fixes exactly one writer per key.
So the clears go through the same unreachable writer: the seal is re-set (fine)
but each clear also exits non-zero, since `Reader.Read` short-circuits on
`sealedKey` and reports every requested key UNREADABLE
(`flowbind.go:200`), which the executor turns into `read_back_incomplete`
(`internal/accessor/executor.go:527`). The alternative reading — clear through
a second, reachable writer — is barred by A3 and would `delete(s, sealedKey)`
anyway, destroying the very state S9 exists to construct.

**Test it prevents.** The one test that pins the emptiness predicate's count to
STORE keys rather than owned keys. `0019:§mini-checks` `oracle` row 9 names its
negative control ("count owned keys only → seeds into sealed store"), and the
`disposition` table's "Sealed (one-key) store" row is the only input class
whose entire distinguishing property is that the store's sole key is not owned.
If the fixture cannot be built, that row has no test, and a future
implementation that counts owned keys ships green.

**What would fix it.** State the arrangement that actually yields a
one-key sealed store — e.g. the clears run first through a reachable writer and
the sealing write lands last (which requires the sealing write to target a key
outside the cleared set, or a second role), or the fixture is built by direct
artifact construction rather than through the CLI. Either is a statement the
scenario currently does not make.

### 2. `0019:C1`, `0019:D-selection-predicate` — the ALL quantifier has no scenario

**Claim.** C1 declares the ALL-over-bound-artifacts quantifier load-bearing and
names its exact hazard, but no scenario, MVV step, or mini-check row exercises
a model with more than one bound artifact.

**Evidence.** C1: "The quantifier is ALL, not ANY and not per-artifact, and it
is load-bearing rather than stylistic: a model may declare N write accessors
over N roles … so a per-artifact or ANY reading would seed into a store that
still carries a key whenever a SIBLING artifact happened to be empty." The
construction is real — `internal/table/model.go:521` declares
`Writers map[string]Accessor`, and `internal/cli/flow_state.go:357` already
resolves `req.artifacts[def.Accessor.Role]` per writer. But every scenario
fixture is single-artifact: `S1`/`MVV` bind `--artifact state=<fresh path>`;
`S4` uses two artifacts only as a diff target between two separate runs, not as
two roles of one model; grep of the whole validation section finds no "sibling",
"second role", or "per-artifact" fixture. The `oracle` mini-check table has
rows for the predicate's STORE-key count (row 9) and its carrier scope (row 10)
but none for its quantifier.

**Test it prevents.** Bind two roles, seed role A, leave role B empty, run
`init-state`: it must be a NO-OP over both. Under an ANY or per-artifact
reading it seeds role B — resurrecting into a store the model still holds keys
in, which is precisely the hazard `0019:A5` exists to exclude. Nothing in the
record currently fails on that change. The cleared-key guarantee C1 states
"holds per model rather than merely per artifact" is exactly the property with
no test.

## Medium

### 3. `0019:S8` — no constructible mismatch on the only carrier the RDR admits

**Claim.** S8 needs "a writer whose read-back disagrees with the value
written," but on the file-backed carrier — the only one `C1` admits — writer
and reader share one store, so a disagreement cannot be produced.

**Evidence.** `flowbind::Writer.Apply` sets `s[t.Key] = t.Value` and saves;
`flowbind::Reader.Read` (`flowbind.go:189`) loads that same file and returns
each requested key's stored value. There is no divergence knob. The nearby
constructions produce a different class: an unreachable *writer* locator seals
and yields `read_back_incomplete` (`executor.go:527`), and an unreachable
*reader* locator errors out to the same incomplete class — never
`ClassReadBackMismatch`, which `executor.go:579` reaches only when
`verifyReadBack` finds an evaluated inequality. Grep of `internal/cli/*_test.go`
finds no existing `read_back_mismatch` fixture at the CLI layer to copy. The
constructions that WOULD work are all unstated: a read accessor bound to a
different artifact path than its writer, or a command-backed read accessor
(legal — C1's carrier refusal is scoped to WRITE accessors only, so a
command-backed reader passes the carrier gate).

**Test it prevents.** The PRESENT-AND-UNVERIFIED terminal arm and, more
importantly, its second half — "a subsequent re-run is a no-op that does NOT
repair it." That clause is the one C1 says "MUST NOT be documented as its
recovery," and `0019:F2` restates it for the present-and-wrong key. Without a
stated fixture, the arm ships untested and the anti-repair property is asserted
in prose only.

### 4. `0019:F2` — the torn multi-writer failure mode has no scenario

**Claim.** F2 describes a distinct post-condition — partial seed, non-empty
store, re-run lists the missing keys without repairing them — and no scenario
covers it.

**Evidence.** `0019:F2`: "a failed writer can leave the artifact seeded for
some keys only. The store is then non-empty, so a re-run is a NO-OP whose
payload lists the missing `[initial]` keys — visible rather than silently
skipped — but it does NOT repair them." C1's no-op clause names "torn" as its
own input class alongside "partially seeded, post-clear, or fully seeded." But
`S3` covers only post-clear on a single-writer fixture, and the `disposition`
table collapses torn into the generic "Any bound artifact non-empty" row. No
scenario constructs a partial-write state, and note that a torn state requires
multi-writer — the same multi-artifact fixture finding 2 says is absent.

**Test it prevents.** Seed partially (one writer fails mid-plan), re-run, and
assert the payload names exactly the un-seeded `[initial]` keys AND that the
seeded ones are not rewritten. Missing this, an implementation that "helpfully"
completes a torn seed on re-run — the single most natural thing to write, and
the direct violation of the cleared-key guarantee — passes every stated test,
because on a single-writer fixture torn and post-clear are indistinguishable.

### 5. `0019:S7`, `0019:S10`, `0019:MVV` step 8 — the refusal-code oracles have no spelling

**Claim.** Three scenarios assert on "the C1 class code" and "the carrier
refusal code," and `C1` explicitly defers both spellings, so the assertion
cannot be written as an equality.

**Evidence.** C1: "It carries a dedicated code in the `flow-*` family (spelling
sharpened pre-lock)". `0019:§technical-design` restates the deferral: "Exact
payload field names, the refusal-code spelling, and help text are deliberately
deferred to Resolve/Pre-Lock." The comparable shipped codes ARE fixed
constants — `internal/cli/flow_input.go:36-63` (`codeWriteUnbound =
"flow-write-unbound"` and 20 siblings) — and the record correctly says the
shared classes reuse them unchanged. The gap is only the two new ones.

**Test it prevents.** A test that fails when the refusal is raised under the
WRONG code — e.g. the decision-table refusal surfacing as `flow-model-invalid`
rather than its own class code, which would collapse a class refusal into a
selection refusal and silently defeat the "CLASS refusal, not an emptiness one"
distinction C1 spends a paragraph establishing. The implementer will invent a
spelling and assert against their own invention, so the test cannot fail. This
is a disclosed deferral rather than an omission, hence Medium; it becomes moot
the moment the two spellings land.

## Low

### 6. `0019:RT3` vs `0019:S4` — disagree on how many loader-only kinds are testable

**Claim.** RT3 exempts three kinds from the byte-identity invariant; S4 says
two of the three are testable by a different oracle. The two passages are
reconcilable but state different things about the same set.

**Evidence.** RT3: "the three kinds only the loader admits (bare scalar for a
set tag, array for a scalar tag, empty scalar) have no `--write` transcription
to compare against and are outside this invariant." S4: "the 2 kinds only the
loader admits AND that are writable … seed successfully via the normalized path
and read back value-for-value … the assertion is read-back equality, not
cross-route byte identity. The empty scalar … is unwritable by either route
(F5), so it is not a row here." The `fidelity` mini-check row already carries
the reconciliation ("the 3 kinds … 2 are testable; the empty scalar is
unwritable, F5"), and `0019:F5` grounds the empty-scalar exclusion in
`internal/cli/flow_input.go::canonicalValue`. So the record does resolve it —
just not at RT3, which is where a test author reading the invariant would stop.

**Test it prevents.** None outright; S4 is the operative passage and it is
precise. The risk is a test author working from RT3 alone dropping the two
loader-only writable kinds entirely, losing the read-back-equality rows that
cover the normalized seed path — the exact path C1 chose over argv
transcription.

### 7. `0019:S5`, `0019:S6` — "asserted on artifact bytes" is undefined when no artifact exists

**Claim.** Both scenarios specify a bytes-level oracle for a refusal that may
occur before any artifact file exists.

**Evidence.** S5: "ZERO writes committed — asserted on artifact bytes, not just
exit code." S6: "zero writes committed." Both are plan-validation refusals, and
C1 places plan validation before any write. If the fixture binds a fresh path,
there are no bytes to compare — the assertion has to be "the path is still
absent." S7 gets this right and states both arms explicitly ("bytes/mtime
unchanged, or the artifact path still absent"); S5 and S6 do not.

**Test it prevents.** Nothing structural — the test author will do the obvious
thing. Flagged only because S7 establishes the precise phrasing two scenarios
earlier and S5/S6 do not inherit it, which is the kind of asymmetry that reads
as intentional when it is not.
