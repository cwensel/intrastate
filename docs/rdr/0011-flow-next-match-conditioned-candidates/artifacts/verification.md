# Verification — RDR 0011 `flow next` match-conditioned candidates

Phase 3a Chain-of-Verification. Each entry below designs an input that would
make a correct implementation VISIBLY violate a REQ, then runs it against the
built CLI (`make build` → `./bin/intrastate`, HEAD `0c2e6e5`). Probes were
authored independently of the Phase 1 test files, which were not read.

Model fixtures were authored fresh under a scratch directory (not committed);
each is described inline by the row shape it exercises, so a later reader can
reconstruct it from the description.

---

## CoVe pass 1 — independent probe of the implementation (no test files read)

**Result: no violations found.** 41 probes across 13 purpose-built scratch
models plus `models/rdr.toml`. Every probe designed to expose a breach of the
candidate predicate, the `--all` restoration, the `unknown` payload, the
demand-set widening, the `disposition` grid, and the exit/output contract
returned the behaviour the record specifies.

### Probes that would have exposed a breach, and what they returned

**The candidate predicate and its boundaries (REQ-1, 3–5, 7–14, 68, 77–80, 83).**

- `models/rdr.toml` at owned `stage=resolved` → exactly `prelock`,
  `resolve-abandon`, `resolve-route-back`; `outcomes` the full three-member
  alphabet; `readers` `["rdr-status"]`. `--all` → 21 candidates, same
  `outcomes`, same `readers`. A stripped-match build reports 21 by default,
  so the count discriminates (REQ-87/88/95).
- A two-row model over an OBSERVED match key: no `--tag` → both rows
  candidates each carrying `{mode, absent}`; `--tag mode=fast` → only the
  present-and-EQUAL row, `unknown` `[]`; the present-and-UNEQUAL row is
  gone (REQ-11/12/77/78/79).
- **Empty-string presence (REQ-4, ASSUMPTION A-1).** A reader reporting
  `mode = ""` against a row matching `mode eq "fast"` → the row is
  EXCLUDED (zero candidates), i.e. `""` is PRESENT and compared, not
  treated as absent. The same model with the key genuinely omitted →
  candidate carrying `{mode, absent}`. A non-empty presence test would
  have inverted the first case; it does not.
- **Dead row (REQ-9/80).** A row authoring `eq = "x"` and `in = ["y"]` on one
  key: candidate carrying `{dk, absent}` while `dk` is absent; excluded at
  `dk=x`, at `dk=y`, AND at `dk=z` — a value on neither side of the pairing.
  The exclusion is therefore the kernel's conjunction, not a CLI literal
  comparison.
- **Multi-key mixed-state row (REQ-103).** One row matching `k1 eq x`,
  `k2 eq x`, `k3 eq x` under `--tag k1=x --tag k2=y` (present-equal,
  present-unequal, absent simultaneously) → EXCLUDED by default; under
  `--all` a candidate whose `unknown` carries no match entry. This is the
  per-key-independence case A12's spike recorded as unmeasured.
- **Set-kinded match key (REQ-6/102/127).** A `set`-kinded observed key
  matched by `eq = ["green"]`, supplied as `--tag 'labels=["green","green"]'`
  → the duplicate collapses to the canonical `["green"]` and the row
  MATCHES; `--tag 'labels=["blue"]'` selects the sibling row instead. A
  filter-before-`KernelRow()` build reimplementing `seamValue` would
  mis-compare here; it does not.
- **Guard decided false with an absent match key** → excluded in BOTH modes,
  with a guard-true sibling surviving in both as the negative control
  (REQ-83/86).
- Source check: the only `Value ==` comparison in `flow_next.go` is the
  pre-existing `ClearSentinel` write-preview split; no literal comparison
  exists in `internal/cli` for match (REQ-3/5/10).

**`--all` restores the PREDICATE, not the reader set (REQ-36–42, 46–47, 50, 86, 97, 115; Q1/Q2).**

- `--all` registered on `next` only. `flow resolve|read-state|set-state
  --all` and `flow --all` and root `--all` each → `command-error`, exit 2,
  one NDJSON envelope. `-a` shorthand is unregistered (`unknown shorthand
  flag: 'a'`), confirming ASSUMPTION A-9.
- **Q1 discriminator (ASSUMPTION A-12).** A model whose key `k` is BOTH a
  match key and a guard key and is NOT in `RequiresOwned`: under `--all`,
  `{k, absent}` STILL appears — from the guard atom. A per-KEY filter would
  have deleted it. The filter is per-ATOM on `Block`, as C2 requires.
- **Filter/dedup ORDER (REQ-40/41/97).** Two independent fixtures: (i)
  `models/rdr.toml` over an empty artifact, where `stage` is both a
  `[rule.match]` and a `[rule.write]` key — `--all` still reports
  `{stage, absent}` from the owned-key walk; (ii) a purpose-built model
  whose owned key `mw` is matched AND written — `--all` still reports
  `{mw, absent}`. A build that dedups before filtering deletes both pairs;
  neither is deleted.
- **Q2 discriminator (ASSUMPTION A-13, REQ-18/115).** Over the match-only
  owned-key model, `next` and `next --all` report the IDENTICAL `readers`
  set `["md","st"]`. `--all` does not revert the demand set.

**The `unknown` payload: shape, vocabulary, dedup, ordering (REQ-22–31, 44, 71–73, 81–82, 84–85, 104, 113, 124–125).**

- Wire members are exactly `{"key","reason"}`, always emitted, on a list
  named `unknown` present as `[]` rather than omitted — confirmed by
  `keys` on the candidate object with and without `--evaluate-gates`
  (REQ-44, A-2). `gates` remains `omitempty` (0005-carried, unchanged).
- **`uncomparable`, the payload-only reason (REQ-84/112).** An `int`-kinded
  OWNED key delivered by a reader at the non-integer value `"abc"`, guarded
  `gt = 3` → `{num, uncomparable}`. (`--tag` cannot produce this: it
  conforms the value at parse and refuses `flow-tag-invalid`, so the
  reader route is the one reachable producer, exactly as A13 states.)
- **The `owned_state_unavailable` fidelity limit (REQ-31/85/112).** A row
  carrying BOTH that `uncomparable` guard atom AND an owned key its bound
  reader reports absent → candidate carrying `{extra, absent}` from the
  walk, and the `uncomparable` fact is NOT reported. The sibling row
  without the missing owned key reports `{num, uncomparable}` in the same
  run. This is the precedence limit pinned rather than a defect.
- **Dedup on the pair (REQ-24/29).** A guard atom over an ABSENT key
  produces `{num, absent}` exactly ONCE — the walk and the kernel payload
  agree and collapse.
- **Sort order (REQ-29/113).** A row minting facts deliberately out of
  order — match atoms authored `zz`, `mm`, `aa`; gate ids `zgate`,
  `agate`; plus an unestablished owned key `extra` — returns
  `aa/absent, agate/not-evaluated, extra/absent, mm/absent,
  zgate/not-evaluated, zz/absent`: all four sources interleaved and sorted
  by `(key, reason)`. An unsorted build fails here.
- **Gate ids (REQ-26/27/28/81).** Ids come from the ROW's own `gate` list,
  not a model-level set: two rows declaring different gates report their own
  id only. `not-evaluated` appears on no non-gate entry, and `uncomparable`
  on no match/owned/gate entry across every fixture (REQ-124). Gate ids are
  emitted on `owned_state_unavailable` too (REQ-28).
- **Set-kinded property.** The `unknown` list is a set of pairs: identical
  pairs from different sources collapse; distinct pairs on one key are not
  reachable (no key can carry two reasons), matching S7's stated limit.
- Text mode renders `candidates[0].unknown[0].key: dk` /
  `...reason: absent` through the shared flattener, both members present,
  no bespoke template (REQ-72/73).
- No schema or version field is emitted (REQ-125).

**Gates run only for REPORTED candidates (REQ-12/70).** The sharpest available
discriminator: a model whose match-EQUAL row declares gate `gfast` and whose
match-UNEQUAL row declares gate `gslow`, with the `gate` artifact role
deliberately UNBOUND. Under `--tag mode=fast --evaluate-gates` the run fails on
`gfast` only; under `--tag mode=slow --evaluate-gates` it fails on `gslow` only.
A build that ran an excluded row's gates would have failed on `gslow` in the
first case. It does not.

**The demand-set widening and its effect on `flow resolve`, in BOTH directions
(REQ-16–21, 92–93, 114–118, 120, 126).** A model declaring an owned key `mode`
that one row MATCHES on but no row writes, clears, or guards, served by exactly
one reader:

- `next` invokes `md`, carries `mode` in the view, and match-DECIDES the row
  (no `{mode, absent}`) — the discriminating oracle for the term (REQ-114).
  Same `readers` under `--all` (REQ-115).
- `flow resolve --outcome go`, all three named arms: reader bound and
  answering → PLAN, exit 0; role unbound → `flow-artifact-missing`, exit 2;
  reader refusing (malformed artifact) → `flow-accessor-failed`, exit 3.
  `flow-no-match` in NONE of the three (REQ-19/93).
- **The BREAKING arm reproduced (REQ-20/116).** With a second ordinary row
  for the same outcome that does NOT match on `mode`, and `md` unbound: the
  run that would have returned `plain-row`'s plan instead exits 2
  `flow-artifact-missing`. The union-over-rows demand pulls the reader in
  above the kernel. This is the class C1/Consequences NAME and ACCEPT; it
  fires exactly here and nowhere observed beyond it.
- **Neither over- nor under-fires.** Over a model whose match key is
  OBSERVED rather than owned, the term is empty and `readers` is unchanged
  between `next` and `resolve` (REQ-21). Over `models/rdr.toml`, `flow
  resolve` at `advance`/`revise`/`abandon` still plans and an unmodeled
  outcome still refuses `flow-unmodeled-outcome` (REQ-118).
- **Both verbs agree (REQ-117).** `next` and `resolve --outcome go` over one
  model return byte-equal `readers`.
- **Escape-row scoping.** `resolve --outcome go` does NOT demand the reader
  serving an escape row bound to `stop`; `resolve --outcome stop` DOES. `next`
  takes the union over all rows and demands both. Consistent with the demand
  set being a union over the OUTCOME's rows and `next`'s DEV-1 reading (a);
  not a violation.

**The three pins licensing the presence test (REQ-63/106).** A model declaring
one owned key served by two readers is refused at load
(`malformed_accessor_binding: owned tag mode is served by 2 readers; want
exactly one`); `--tag mode=fast --tag mode=slow` → `flow-tag-duplicate`, exit 2;
`--tag phase=b` on an owned key → `flow-tag-owned`, exit 2. All three precede
any probe.

**Structural REQs checked at source rather than at runtime.** `internal/resolve`
carries an EMPTY diff across the branch (`git diff --name-only main...HEAD`),
and the only production files changed are `internal/cli/flow_next.go` and
`internal/cli/flow_exec.go` — REQ-34/107/120 hold in the one mechanically
checkable form MVV 9 names. No occurrence of `unresolved`/`Unresolved` remains
in production or in the three 0005 test files' prose (REQ-55/60); the
`flow_next_0005_test.go` header names RDR 0011 as the source of the default
(REQ-53). `docs/cli-output-contract.md` documents the payload, the `unknown`
reason table, and the `--all` invocation, and `README.md:43` no longer states
the overridden default (REQ-56/122). `go test ./...` green.

### Coverage and what was sampled rather than probed

Probed directly, with a designed-to-break input: REQ-1..21, 22..31, 32..44,
46..50, 63..64, 66..73, 74..86, 87..93, 95..106, 112..118, 120, 122..127. That
is the candidate predicate, the demand set, the `unknown` payload, `--all`, the
`disposition` grid in full, the MVV, and the exit/output contract.

Sampled or checked structurally rather than by probe (they are review-time or
diff-time obligations no runtime input can express, and the record says so):
REQ-45 (request identity), REQ-51..62 (help wording and the test re-homing
census — the header, the production rename, and the prose sweep were checked;
the five per-line test reads were NOT, since the files are Phase 1 artifacts
this pass is forbidden to read), REQ-65, REQ-94, REQ-107..111, REQ-119,
REQ-121, REQ-128..130.

One class could not be probed from the CLI and is not a gap: the
present-but-CONFLICTED match key (REQ-64, A3), which the three pins above
refuse before a probe exists — verified negatively by confirming all three
refusals fire.

### Verdict

No `FAIL-N` entry is owed from this pass. Every probe designed to distinguish
the specified behaviour from its nearest wrong neighbour — a filter applied
per-key instead of per-atom, a dedup that runs before the `--all` filter, a
demand set that reverts under `--all` or unions the wrong rows, a presence test
that treats `""` as absent, a payload read off the wrong refusal kind, gates run
for an excluded row — returned the specified behaviour.
