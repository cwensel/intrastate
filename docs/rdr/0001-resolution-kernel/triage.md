# RDR 0001 Roborev Triage

## Frozen World

- **BASE:** `991bd43`
- **HEAD:** `240c856`
- **Batch:** `batch:rdr-0001`
- **Triage basis:** nine review jobs produced 27 finding instances,
  consolidated into nine root causes and grounded against the RDR cluster and
  current implementation head.

## Finding Dispositions

| Root cause | Origin jobs | Route | Stable anchor | Grounding evidence and outcome |
| --- | --- | --- | --- | --- |
| Replay revision/table binding is not enforced | 4468, 4479 | **RDR-SEED `jv85`** | `internal/resolver/resolver.go::Input.TableRevision`, `Resolve` | `TableRevision` is carried in the replay tuple but is not bound to `Input.Table` during resolution. The architectural ownership question is preserved in `jv85`; no speculative kernel fix was made. |
| Recognized outcome has no declared alphabet boundary | 4468, 4471, 4474, 4475, 4477, 4478, 4479, 4483 | **KATA-BUG `a2m3`** | `internal/resolver/resolver.go::modelsOutcome` | The kernel treats an outcome as modeled when any supplied edge uses its value; it cannot validate the value against a declared outcome alphabet. Job 4483's duplicate outcome issue is the same root cause and is attached to the existing implementation bug `a2m3`. |
| Successful plans omit selected rule identity | 4474, 4477, 4478, 4479 | **RDR-SEED `zy0m`** | `internal/resolver/resolver.go::TransitionPlan`, `plan` | A plan contains next tags and inert writes but no stable row/source identity. The contract and ownership decision is preserved in `zy0m`; no ungrounded field was added. |
| Resolver lacks the RDR 0003 guard integration path | 4474, 4477, 4478, 4479, 4483 | **KATA-BUG `n07s`** | `docs/rdr/0003-guard-predicate-exhaustiveness.md::Approach`, `Normative Contracts`; `internal/resolver/resolver.go::Guard` | RDR 0003 owns equality, membership, bounded integer comparison, existence, and set-containment semantics, but the consumer integration path into RDR 0001 remains actionable. Job 4483 changed the earlier DROP adjudication to the implementation bug `n07s`. |
| Artifact-missing findings were stale | 4468, 4471, 4475 | **DROP — `superseded-at-HEAD`** | `docs/rdr/0001-resolution-kernel/{coverage.md,deviations.md,req-list.md,status.md,verification.md}` | The implementation artifacts exist at the grounded head. The missing-artifact claims no longer describe the current tree and require no additional write. |
| Returned-plan aliasing could corrupt replay | 4471 | **DROP — `superseded-at-HEAD`** | `internal/resolver/resolver_adversarial_test.go::TestResolveAdversarialReturnedPlanCannotMutateReplay` | Commit `393fda0` added a persistent adversarial replay test proving caller mutation cannot alter the supplied table or the next disposition. |
| No-match escape precedence lacked regression coverage | 4471 | **FIX-NOW `86b89d9`** | `internal/resolver/resolver_adversarial_test.go::TestResolveAdversarialNoMatchEscapes` | Commit `86b89d9` now covers one matching escape, duplicate escapes, and ordinary-match precedence. The focused resolver suite is green. |
| Verification evidence lacks a durable claim-to-test mapping | 4478, 4483 | **KATA-BUG `cvg8`** | `docs/rdr/0001-resolution-kernel/verification.md::Phase 3a — Chain-of-Verification` | Commit `9b3480e` replaced the deleted scratch command with `go test ./internal/resolver`, but that reproducibility edit was partial: the recorded aggregate probe result is not durably mapped to named persistent tests. Job 4483 routes the final evidence correction to `cvg8`. |
| Bounded fix review found a separate follow-up | 4481 | **KATA-BUG `z3ae`** | roborev job `4481` over fix commit `86b89d9` | The bounded review outcome is preserved as the independently actionable bug `z3ae`; it is not broadened into this documentation-only increment. |

## Open Questions Preserved in RDR Seeds

- **`jv85`:** Should a revision-to-table identity binding be enforced by the
  pure resolver, guaranteed by normalized table construction, or remain a
  caller precondition?
- **`zy0m`:** Must a successful transition plan expose stable selected
  rule/source identity, and which peer RDR owns that identity contract?

## Outcome

- **Root causes:** 9
- **Review jobs:** 9 (`4468`, `4471`, `4474`, `4475`, `4477`, `4478`,
  `4479`, `4481`, `4483`)
- **Raw finding instances:** 27
- **DROP:** 2 (stale artifacts, returned-plan aliasing)
- **FIX-NOW:** 1 (escape regression tests at `86b89d9`)
- **KATA-BUG:** 4 (`a2m3`, `z3ae`, `n07s`, `cvg8`)
- **RDR-SEED:** 2 (`jv85`, `zy0m`)
- **Batch label:** `batch:rdr-0001`
- **Next scoped flight:**
  `/kata-flight --label batch:rdr-0001 --drain`

All findings are routed. No unresolved finding requires an uncommitted code or
test change in this worktree.
