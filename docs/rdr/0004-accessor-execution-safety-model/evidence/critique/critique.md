# Critique — Accessor Execution Safety Model

## Three likely implementation failures

1. **Write verification blesses collateral mutation.**
   - **Root cause in the RDR**: The read-back invariant checks that expected
     owned-tag values are present, but it does not require non-owned observed or
     recognized tags to remain unchanged.
   - **Enabling passage**: "compare expected owned-tag values" and
     "`write -> read = expected owned-tag value identity` for the written tag
     subset."
   - **User symptom**: A transition reports success after changing `status`, but
     the same write accessor also corrupts a caller-observed tag such as
     `profile`; the next guard or review run behaves differently even though the
     write path returned success.

2. **Timeout refusal is mistaken for rollback.**
   - **Root cause in the RDR**: Timeout is a distinct refusal class, but the RDR
     does not spell out what callers may assume about artifact state after a
     write accessor times out mid-mutation.
   - **Enabling passage**: "Every accessor invocation MUST have a bounded
     timeout. Timeout MUST be reported as its own refusal class."
   - **User symptom**: A skill retries a timed-out write against stale reads and
     either repeats the mutation or reports a confusing read-back mismatch.

3. **Typed capability names become a thin wrapper over arbitrary integration
   code.**
   - **Root cause in the RDR**: The RDR rejects shell strings and ambient
     artifact discovery, but it does not define how much the typed binding
     interface must expose for review.
   - **Enabling passage**: "It invokes typed bindings selected by accessor name
     and capability" and "External integrations remain possible, but only
     through typed bindings."
   - **User symptom**: Reviewers see `state.persist/write` in model data, but
     the actual binding still owns broad process authority and the first unsafe
     behavior appears only during an integration run.

## Section likely rewritten within six weeks

The **Round-Trip / Inverse Invariants** section will be rewritten first. The
current wording is the decisive safety claim for artifact mutation, and the
spike already demonstrates the hole: `write` can produce `tags={status=Final}`
while the observed artifact also contains another changed tag. The RDR needs to
say whether protected non-owned tags are part of success or deliberately outside
the contract.

## Assumption least likely to survive first contact

**A2** in its original shape: "Write accessors can verify their intended
owned-tag effect by re-reading the same artifact boundary after the write." That
survives for the owned tag, but it is insufficient as an artifact safety model
because a real write implementation can change adjacent observed or recognized
tags while still setting the expected owned tag.

## Premortem

Implementation shipped with a clean validator and a convincing fixture. The
accessor registry had `state.read`, `state.gate`, and `state.persist`; `write`
performed same-role read-back; `respond.Fail` mapped refusals through
`CLIError`. The first real user flow was a status transition on an RDR artifact.
`state.persist` wrote the planned owned `status=Final` value, then normalized a
nearby metadata field that the resolver treated as observed input. Read-back
checked only the planned owned tags, so `ExecuteWrite` returned success. The
next `Resolve` run read the mutated observed tag, a guard predicate evaluated
differently, and the user saw a transition that looked nondeterministic.

Debugging went in circles because every log line named the correct accessor
identity, capability, role, timeout, expected value, and observed owned value.
Nothing in the success result mentioned protected non-owned tags. QA reproduced
it only after adding a fixture where the write accessor sets `status=Final` and
also changes `profile=small` to `profile=large`. The implementation had obeyed
the literal RDR: expected owned tags matched. The RDR had failed to encode the
actual safety property users needed from authoritative artifact mutation.

## Acceptance tests that would have caught it

1. Given a caller-supplied `state` artifact with owned `status=Draft` and
   observed `profile=large`
2. And a successful transition plan expects owned `status=Final`
3. When `state.persist` writes `status=Final` and leaves `profile=large`
4. Then the write result is success

1. Given the same artifact and transition plan
2. When `state.persist` writes `status=Final` and changes observed
   `profile=small`
3. Then the write result is `read_back_mismatch`
4. And the mismatch details include the protected tag key, expected pre-write
   value, and observed post-write value

1. Given a write accessor times out after invocation starts
2. When a caller retries the transition
3. Then the caller first re-reads the authoritative artifact role before
   building a new transition plan
