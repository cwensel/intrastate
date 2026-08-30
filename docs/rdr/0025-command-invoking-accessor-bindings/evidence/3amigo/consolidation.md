Model: claude-opus-5[1m]

# 3amigo consolidation — 0025

Mechanical merge of three isolated persona passes (`persona-1-pm.md`,
`persona-2-implementer.md`, `persona-3-qa.md`). Hotspots are the set
intersection of the element ids each persona anchored to — the personas never
saw each other, so agreement is evidence, not conformity. Overlap marks a
hotspot passage, not a validated finding; a single-persona finding is not
thereby weaker.

## Hotspots

### 3/3 personas

- **0025:A10** — PM: Pending on a *decision*, so C6's "set it once" adoption
  promise may degrade to `--allow-commands` per invocation; blocks
  `G-assumptions` (lock requires Verified) and Profile re-derivation.
  Implementer: leaves flag-vs-file open, so Phase 4's scope (does it build a
  config reader at all?) is undetermined. QA: C6 is normative while A10 is
  Pending; the flag-vs-config precedence rule and the malformed-config refusal
  class are both untestable.
- **0025:C6** — Implementer: the execution gate names no enforcement site among
  `registry.go::Registry`, `NewExecutor`, and the binding, and must reach three
  live `NewExecutor` call sites. QA: S7's "a malformed config must refuse,
  never fail open" is a *normative* requirement stated only inside a scenario
  Expected — C6 itself covers only absent-or-false. PM: the one-time-opt-in
  user outcome C6 promises is contingent on the unresolved A10.

### 2/3 personas

- **0025:C3** — Implementer: `output`/`exit_absent`/`exit_verdicts`/`env`/
  `env_pass` are normative but absent from C1's carrier block, and
  `source.go::decodeStrict`'s `DisallowUnknownFields` makes any of them
  `unknown_schema_field` over the whole document rather than the C5 category S1
  asserts. QA: "an omitted key is not-carried" is ambiguous between
  established-absent (`KeyValue.Absent`) and unreadable — `executor.go::classify`
  gives these different refusal outcomes — and the default `output = "json"`
  read mode has no scenario in S1–S7.
- **0025:C4** — Implementer: the new `Refusal.Detail` has no delivery path;
  `binding.go` interfaces return bare `error`, `executor.go::invokeRead`/`Gate`/
  `Write` discard it, and `flow_exec.go::accessorFailure` never reads `Detail`
  while already using `CLIError.Detail` for `detailMayHaveApplied`. Also: `LC_*`
  is the sole glob in an explicitly no-glob env policy, and `HOME` is
  allowlisted despite A7 naming `~/.gitconfig` as the ambient-widening vector.
  QA: the four-part env composition is called "total and ordered" with no
  collision precedence stated; S4b tests only the `GIT_CONFIG_COUNT` negative.
- **0025:C5** — Implementer: "all five are appended" contradicts A6's "the four
  defects", and constant naming/ordering into `table.Categories()` is unstated
  against order-sensitive golden maps. QA: the interpreter deny-list is
  normatively OPEN but S1's single `sh -c` mutant reads as "inline shell is
  blocked"; no scenario asserts that an unlisted interpreter is admitted.

## Single-persona findings

### PM (persona 1)

- High | 0025:MVV — step 3's "the resolved decision invokes the declared write
  command" has no code path: `flow_state.go:311` is the sole `exec.Write` caller
  and takes caller-supplied `--write` flags; no Phase 1–4 wires resolve→apply.
  Blocks Scope Verification (`G-scope`).
- High | 0025:§problem-statement — the stated harm conjoins "edit not declared"
  and "caller carries the decision across"; only the first is fixed, and no
  section records the second as surviving residue (unlike the Selection LBD's
  honest note). Blocks `G-proportionality` and the successor-seed decision.
- Medium | 0025:MVV — end state claims "through declared commands **only**" while
  step 1 binds the write to a wrapper script the model does not carry.
- Medium | 0025:§decision-rationale — the deciding QOC row "Binds established
  tools without wrappers" scores O2 = 4, but A3/§approach establish that
  established-tool *writes* always need a wrapper. Choice still holds (24 vs 22)
  but the stated reason misdescribes the outcome. Blocks `G-contradiction`.
- Low | 0025:§approach — cites "(Seam Lineage)" but no such section exists in the
  outline; the substance lives in §background.
- Low | 0025:MVV — step 3 says `allow_commands = true` "in the invoking
  configuration" while C6/A10 permit a flag-only v1; S7 handles both, the MVV
  does not.

### Implementer (persona 2)

- Medium | 0025:C2 — argv0 resolution "against the model file's directory" is
  unbuildable: `table.Model` and `table.Accessor` carry no source path, and
  `Registry` receives only `*table.Model`.
- Medium | 0025:C2 — the absolute-path precondition names no absolutizer;
  `flow_input.go::parseArtifacts` stores `--artifact` paths verbatim, so existing
  relative bindings would refuse for command entries, and siting the check "in the
  executor" contradicts C4's classes-unchanged claim.
- Medium | 0025:C1 (with C5) — the neither-declared case re-categorizes an
  existing `CatMalformedAccessorDeclaration` refusal asserted by live fixtures,
  and `accessorTable`'s fail-fast single-failure switch leaves C5 defect
  precedence unfixed.
- Low | 0025:§failure-modes — the Windows before-spawn refusal lives only in prose
  with no contract home, no defect spelling, and no build-tag-vs-runtime-check
  decision (determines whether `Setpgid`/`Cancel` code compiles unguarded).

### QA (persona 3)

- High | 0025:S3 — the cited normative fixture FX-exit-codes (E1a–E1h, R6) is
  entirely empty-stdout, so it supplies no oracle for the three non-empty-stdout
  ordering cases the scenario exists to prove, and the positional "respectively"
  list does not align to it. Prevents the ordering test; also leaves open whether
  C4's parse-before-classify applies to a stdout truncated at the 1 MiB cap.
- Medium | 0025:F7 — Windows refusal plus platform-neutral lint is a two-part
  normative behaviour with no covering scenario, no contract clause, and no CI
  runner (`ci.yml` is `ubuntu-latest` at all six jobs).
- Medium | 0025:S5b — the "name-sorted first match" oracle is emergent across two
  packages (`readerFor` scans slice order in `accessor/model.go`; `flowbind/
  registry.go` does the sorting), so S5b passes while the property it pins can
  move silently.
- Low | 0025:S5 — neither C3 nor any scenario addresses a command *read* returning
  the literal `<clear>`, which the inherited `executor.go:96` rule turns into
  `incomplete_read` — newly reachable via raw-mode reads of real tools.
- Low | 0025:S4 — the ablation arms assert against an "ablated build" that is not a
  defined construct and implies a test-injection seam for the
  `Setpgid`/`Cancel`/`WaitDelay` triple that C4 states as fixed behaviour.
