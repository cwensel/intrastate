# Two-boundary absence: binding return vs resolver seam (A15 adjudication)

A Stage-6 consult disagreement was settled against source. One agent, reading
0004's RECORD TEXT (C6/C8), concluded absence must cross as omission from both
slices at the binding return. Two agents reading the EXECUTOR CODE concluded it
must cross as `KeyValue{Absent: true}` in `values`. Source shows the design has
TWO boundaries and 0004:C8 governs only the outer one.

## Boundary 1 — the binding return (`ReadBinding.Read`)

`internal/accessor/binding.go::ReadBinding.Read` doc:
  "It returns one KeyValue per key it could read (present or established-absent)
   and names in unreadable every requested key it could NOT read."

`internal/accessor/binding.go::KeyValue.Absent` doc:
  "Absent reports that the binding established the key is not carried.
   This is the accessor layer's INTERNAL representation; what crosses
   the seam omits the key entirely (`0004:C8`)."

So at boundary 1 absence is a PRESENT record with an explicit flag.
Omission from both slices is UNREADABLE:
`internal/accessor/executor.go::readOutcome.classify` — "A key the binding
classified as neither read nor unreadable defaults to UNREADABLE: guessing
absence is the failure this contract exists to prevent (`0004:C6`)."
The `!ok` branch appends to `unread`; `Executor.Read` then refuses
`ClassIncompleteRead`. Pinned by
`internal/accessor/read_completeness_0004_test.go::TestReq29_UnclassifiedKeyDefaultsToUnreadableNotAbsent`.

An Absent KeyValue is NOT reclassified: the `clearIsUnreadable` branch is
guarded by `!v.Absent` (executor.go:143), and executor.go:147 propagates
`Absent: v.Absent` verbatim.

## Boundary 2 — the resolver seam (what 0004:C8 governs)

`internal/accessor/model.go::ReadResult.OwnedSnapshot` (model.go:344-350)
converts flag -> omission:
  "if v.Absent { // Omission, never a placeholder: a sentinel would make
   `TagSet.has` true and silently retire the refusal. continue }"

`internal/cli/flow_exec.go:288-297` does the same for reported tags:
  "An established-absent key is OMITTED from the reported tags while STAYING
   in `keys` ... it is the same omission `0004:C8` requires at the resolver
   seam, where a placeholder would silently retire an
   `owned_state_unavailable` refusal."

## Consequence for 0025

C7's "`exit_absent` establishes absence by returning the key in NEITHER slice"
and A15's arm 3 describe boundary-2 behavior at boundary 1. Shipped as worded,
every `exit_absent` command read lands in classify's `!ok` branch and refuses
`incomplete_read` — the opposite of the intent.

The repair routes absence through `KeyValue{Key: k, Absent: true}` in `values`,
the same channel `internal/cli/flowbind/flowbind.go::Reader.Read` (flowbind.go:213)
already uses for the path-backed binding: `KeyValue{Key: key, Value: v, Absent: !held}`.
0004:C8 continues to hold unchanged at boundary 2, reached via `OwnedSnapshot()`.

No interface change. No path-backed binding change. No 0004 test change —
TestReq29 keeps pinning the defensive default, which the repair stops relying on.
