Model: claude-opus-5[1m]

# Persona 1 — Product Manager (iter-2, delta)

Question: does this RDR actually deliver the user outcome?

Delta scope reviewed: `0025:§problem-statement` (residue paragraph),
`0025:MVV` (step 3 + End state), `0025:§decision-rationale` (QOC row +
deciding-rows paragraph), `0025:§consequences` (execution-gate bullet),
`0025:C6` and `0025:A10` (flag-only v1 + charted config successor).

Grounded in source (`/Users/cwensel/sandbox/newcoinc/intrastate`) before
reporting: `internal/cli/flow_state.go:311` is still the sole non-test
`exec.Write(...)` caller and its planned tags still come from `--write` flags
(`parseWrites`, line 279); `internal/cli/flow_resolve.go:139` still tells the
caller to "transcribe" the plan; `internal/cli/config/` does not exist; zero
`intrastate.toml` / `allow_commands` hits in any `*.go`; no `os.Getenv` in
`internal/` outside tests; `--as` (`internal/cli/root.go:186`) is the only
persistent flag in the binary.

---

## Prior findings — disposition

- **H1 closed.** `0025:MVV` step 3 no longer asserts the unbuilt
  resolve→apply path. It now names the real path — "the existing write path
  (`internal/cli/flow_state.go`'s set-state verb, whose planned tags come from
  its `--write` flags)" — and states in the same step that automatic
  resolve→apply "is not this RDR's scope and no Phase builds it." That matches
  the source exactly. `0025:G-scope` is now answerable.
- **H2 closed.** `0025:§problem-statement` now carries the residue explicitly:
  "The caller still carries the decision from resolve to apply by hand …
  automatic resolve→apply wiring is a successor's scope, not this record's."
  The `Charted.md` ledger records the successor. This is the same
  bounded-not-retired honesty the Selection LBD already had, which is what H2
  asked for.
- **H3 closed.** `0025:A10` is now **Verified** with a positive decision
  (flag alone), not Pending on a fork. Its absence evidence is exactly what I
  found in the repo. `0025:C6` now reads "**v1 ships the flag alone**" and
  drops the `allow_commands = true` config literal from the normative block.
  The Capability Dependencies row moved to "Not built". `0025:G-assumptions`
  is no longer blocked here.
- **M1 closed.** The MVV End state no longer says "declared commands only";
  it now distinguishes the read ("`git config --get` binds directly") from the
  write ("*the declared wrapper* … its body is script the model does not
  carry"). That is the disclosure M1 asked for, in the one sentence a reader
  quotes.
- **M2 closed.** The deciding-rows paragraph now concedes that the
  established-tool-write half "separates them not at all (A3, Approach)" and
  relocates O2's margin to "the positional read/gate case plus integration
  cost." The tension with the Approach's integration-honesty paragraph is
  gone.
- **L2 closed.** MVV step 3 now gates on `--allow-commands`, matching C6's
  flag-only v1. The MVV is executable as written.

---

## New findings

### Low

#### L1 — `0025:§consequences` execution-gate bullet claims the successor "converts it back to one-time", which `0025:A10` says the successor cannot do alone

The Consequences bullet closes: "the successor that adds the config file
converts it back to one-time." `0025:A10`'s **If wrong** clause specifies the
composition rule the successor inherits: "a later config file composes
disjunctively with the flag (either grants, neither revokes)". Under a
disjunctive rule the config file makes a one-time opt-in *available*; it does
not convert the flag path, and a scripted caller that already passes
`--allow-commands` keeps working unchanged. "Converts it back" reads as a
restoration of the original property; "makes the one-time opt-in available
again" is what A10 and C6 actually license.

This is a wording overreach in a bullet whose whole job is to state the cost
honestly, so it is worth fixing — but it does not change what ships and does
not contradict any normative clause.

**Decision blocked:** none that gates lock. It is a Contradiction-Check
(`0025:G-contradiction`) nit between the Consequence's promise and A10's
composition rule; a reviewer re-deriving the successor's obligation will hit
the mismatch.

---

## Outcome assessment

(a) **Residue disclosure** — honest and complete, and it neither under- nor
over-states. The Problem Statement now separates the *declaration* half this
RDR fixes from the *carry-the-decision* half it does not, names the code that
proves the residue, and routes the remainder to a successor. It does not
overclaim in the other direction either: it still asserts the real win ("what
changes is that the edit it drives is declared and bounded rather than ad
hoc"), which the C1/C2/C5 contracts support.

(b) **Flag-only v1 and the user outcome** — the outcome survives. The
delivered product is "a declared, linted, bounded command applies the state
change"; the gate is an admission condition on that, not part of it. The
friction is real but bounded to one repeated flag on the verbs that execute an
accessor, and A10's sufficiency argument is sound on the code: the gate is a
boolean, the flag sets it, and there is no other opt-in mechanism in the binary
to compare against — no env reads anywhere in `internal/`, one persistent flag
total. Building a config subsystem to avoid retyping a flag would be a larger
new user-facing surface than the contract it gates, which would have moved the
Profile. Deferring it is the proportionate call.

(c) **MVV still validates the promise** — yes. The Problem Statement promises
that a declared command can carry the apply; MVV steps 3–4 drive exactly that
through the write path that exists and verify by read-back through the declared
reader, with step 3's own sentence naming what it does and does not prove. S6
and S7 in the Testing Strategy match the rewritten MVV (S7 now ends "v1 has no
config file to malform (A10)"), so the test surface no longer splits on an
undecided A10.

(d) **New contradictions introduced** — one Low (L1 above). Nothing else in
the rewritten passages is contradicted elsewhere in the record; C6, A10, the
Capability Dependencies row, S7 and the Consequences bullet all now tell the
same flag-only story.
