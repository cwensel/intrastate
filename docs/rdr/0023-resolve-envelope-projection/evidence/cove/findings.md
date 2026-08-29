Model: claude-opus-5[1m]

# CoVe pre-lock lens — cli/0023

## Step 0 — Codebase claim sweep

| # | Claim (element) | Verdict | Source cite |
|---|---|---|---|
| 1 | A2: `respond/text.go::flatten` sorts keys from the decoded map | CONFIRMED | `internal/cli/respond/text.go::flatten` — `sort.Strings(keys)` in the `map[string]any` arm |
| 2 | A2 / Perf: one wire encoder, `SetEscapeHTML(false)`, compact NDJSON | CONFIRMED | `internal/cli/clierr/clierr.go::WriteJSONLine` (l.174-178); `respond.writeJSONLine` routes to it |
| 3 | A2: non-pointer `omitempty` drops empty `{}`/`[]` | CONFIRMED | Go encoding/json semantics; reproduced in `evidence/spikes/a2-encoder-mechanism.md` |
| 4 | A2: non-nil pointer to empty map survives `omitempty` | CONFIRMED | A2 spike program `*map[string]string` / `*[]string` arms; matches Go's rule (pointer emptiness is nil-ness) |
| 5 | A2: shipped default line has `"owned":{}`, `"readers":[]` | CONFIRMED | live run `flow resolve --model models/examples/pricing-decision-table.toml --outcome decide --tag tier=free --tag region=eu --as=json` |
| 6 | A3: existing resolve-payload oracles all run default mode | CONFIRMED | `rg plan-only internal/` → zero hits; `flow_resolve_0005_test.go`, `decision_table_0010_test.go` invoke with no projection flag |
| 7 | A3: `decision_table_0010_test.go` pins `resolvePayload` at exactly 14 fields | CONFIRMED | `internal/cli/decision_table_0010_test.go:434-438` — `rt.NumField() != 14` |
| 8 | A3: 14 is the current count | CONFIRMED | `internal/cli/flow_resolve.go::resolvePayload` — Model, Revision, Observed, Owned, Readers, Outcome, Rule, Gates, Emit, Next, Writes, Clear, Escaped, EscapeClass |
| 9 | A3: "swept all 16 test files invoking `flow resolve`" | IMPRECISE (not load-bearing) | 21 files in `internal/cli/` match; substantive claim CONFIRMED |
| 10 | A5: the `--all` probes are exact-name `Lookup("all")` | CONFIRMED | `internal/cli/flow_all_0011_test.go:55,431,442,456` |
| 11 | A5: root builds via plain `AddCommand` | CONFIRMED | `internal/cli/root.go:190,193,196,198` |
| 12 | A5: "the shipped tests already walk `Commands()` recursively" | REFUTED | every shipped test sweep is a fixed 2-level nested loop (`flow_all_0011_test.go:45/49`, `:423/438`; `flow_surface_0005_test.go:27/31`; `flow_input_0005_test.go:460/464`). No test recurses. |
| 13 | A5: "the hidden `docs` command still enumerates, so the whole-tree walk misses nothing" | REFUTED as stated | the only recursive walker `internal/cli/help_all.go::walkCommandTree` (l.343-351) skips `help` and `completion` by name; its `docs.go` callers further gate `if c.Hidden` (`docs.go:184`, `:245`). `docs` enumerates because the caller allows hidden, not because the walker is total. |
| 14 | Inverse: does a sibling path already register a flag named `all` outside the fenced set? | CONFIRMED (sibling exists) | `internal/cli/help_all.go::wireHelpSubcommandAll` l.231-236 — `helpCmd.Flags().Bool("all", ...)` |
| 15 | A5: `registerSelectionFlags` registers only `model`/`flow`/`artifact` | CONFIRMED | `internal/cli/flow.go:162-167` |
| 16 | A6: refusal envelope is `code,message,param,detail,hint,findings`, no echo member | CONFIRMED | `internal/cli/clierr/clierr.go:50-77` (`Group`, `Cause` are `json:"-"`); live refusal `{"code":"flow-unmodeled-outcome","message":...,"param":"begin"}` |
| 17 | A6: `respond.Fail` marshals the `*CLIError` directly | CONFIRMED | `respond.Fail` → `clierr.EmitJSON` → `WriteJSONLine(out, e)` (`clierr.go:152-157`) |
| 18 | A7: `parseTags` refuses rather than coerces | CONFIRMED | `internal/cli/flow_input.go::parseTags` l.264-306 — reserved/owned/duplicate/malformed all return `userErr` |
| 19 | A7: `canonicalSet` re-encodes value-preservingly | CONFIRMED | `internal/cli/flow_input.go::canonicalSet` l.392-405 — `slices.Compact(slices.Sorted(...))` + non-HTML-escaping encode |
| 20 | A7: `runFlowResolve` copies `--model`/`--outcome` verbatim | CONFIRMED | `flow_resolve.go:174,224,229`; `observedTagMap`/`tagMap` (`flow_exec.go:739-755`) copy k/v unchanged |
| 21 | A7: `invokedReaders` is a pure function of model + outcome | CONFIRMED | `internal/cli/flow_exec.go::invokedReaders(m *table.Model, outcome string)` l.86 |
| 22 | C2: `internal/resolve/resolve.go::Resolve` takes `Owned` and `Observed` as Input | CONFIRMED | `type Input struct { Flow; Table; Owned []Tag; Observed []Tag; Recognized; Guards }` l.324-333 |
| 23 | C2 / D-identity: `revision` is loader-produced provenance, "the plan's only provenance" | **REFUTED** | `internal/cli/flow_exec.go:49` — `func (flowRequest) revision() string { return "" }`; doc comment l.38-40: "NO model can declare one … empty (DEV-7)". Live wire and every golden carry `"revision":""`. |
| 24 | Approach's plan-group list vs `flowResolveExtendedDesc` | CONFIRMED (exact match) | `internal/cli/flow_resolve.go:86-102` lists the same eight, no echo fields |
| 25 | C1: `emit` stays present as `{}` (0010:C4) | CONFIRMED | `flow_resolve.go:48` no omitempty; `emitMap` l.285-291 allocates unconditionally |
| 26 | C1: `escape_class` omitempty-on-unescaped | CONFIRMED | `flow_resolve.go:54` `json:"escape_class,omitempty"`; set only under `if plan.Escaped` l.264-271 |
| 27 | A4: `0024:C4` joins `dispositions` by `Plan.RuleID`, never observed/owned | CONFIRMED | `rdr inspect --select 0024:C4 0024` — "the same single `Plan.RuleID` join path as `emit`" |
| 28 | "no caller-controlled projection exists to reuse" | CONFIRMED | only producer presence rules: `EscapeClass omitempty`; `flow_next.go:71` `gates,omitempty` under `--evaluate-gates` (`flow_next.go:225`) |
| 29 | A1: byte table reproduces on live source | CONFIRMED | `evidence/spikes/a1-byte-width.md` S1 default line matches a live run; jq deletion path is the five echo keys |
| 30 | C1 text-subset clause holds on a real shape | CONFIRMED (this shape) | live `--as=text` sorted leaf lines; plan-group subset is byte-identical subset of default lines |

## Steps 1-2 — Verification questions and independent answers

**Q1 (CODEBASE).** Does `revision` on a shipped `flow resolve` payload ever carry a non-empty value?
**A1.** No. `internal/cli/flow_exec.go:49` is `func (flowRequest) revision() string { return "" }` — constant return. Doc comment states no model can declare a revision (RDR 0002's `[model]` block has no such field) and `KernelTable().Revision` is deliberately not forwarded. Live run confirms `"revision":""`.

**Q2 (CODEBASE).** Is there a shipped recursive walk of the command tree a `plan-only` whole-tree oracle could reuse unchanged?
**A2.** Exactly one, `internal/cli/help_all.go::walkCommandTree` (l.343-351), and it is not total: it skips any child named `help` or `completion` before recursing. Every walk in shipped tests is a fixed two-level `for _, c := range NewRootCmd().Commands() { for _, sub := range c.Commands() }` (e.g. `flow_all_0011_test.go:45/49`).

**Q3 (CODEBASE).** Does any command outside the 0011-fenced set already register a flag named `all`?
**A3.** Yes — `internal/cli/help_all.go::wireHelpSubcommandAll` l.231-236 does `helpCmd.Flags().Bool("all", false, ...)`. The 0011 oracle does not see it because it only descends the `flow` group. A genuinely whole-tree sweep for `all` would fail today; the `plan-only` sweep is safe only because the name differs.

**Q4 (CODEBASE).** Does `respond.Fail` ever touch the verb payload struct?
**A4.** No. It branches on mode and calls `clierr.EmitJSON(cmd.OutOrStdout(), ce)` / `clierr.EmitText(...)`; `EmitJSON` is `_ = WriteJSONLine(out, e)` over the `*CLIError` (`clierr.go:152-157`).

**Q5 (CODEBASE).** Is the field pin still 14, and does `resolvePayload` declare 14?
**A5.** Both 14. `internal/cli/decision_table_0010_test.go:435` `if rt.NumField() != 14`; struct `flow_resolve.go:31-55`.

**Q6 (CODEBASE).** Does `flatten` distinguish a nil value from an absent key, and does an absent key produce no line?
**A6.** `flatten` walks the decoded map, so an omitted key contributes nothing. It renders `nil` and empty containers alike as `<path>: (none)` (`case nil` and the `len(t)==0` arms). The text-subset clause holds provided the projection is true key deletion; a `null` stand-in would emit `(none)` and break the subset. C1 already forbids `null`, so this is consistent.

**Q7 (CODEBASE).** Does `registerSelectionFlags` reach `resolve`, and could a `plan-only` added there leak?
**A7.** Yes it reaches resolve (`flow_resolve.go:79`), next (`flow_next.go:129`), and the state verbs; it registers exactly `model`, `flow`, `artifact` (`flow.go:162-167`). Verb-local registration on `newFlowResolveCmd` is the correct site.

**Q8 (CODEBASE).** Does a sibling path already implement caller-controlled payload narrowing?
**A8.** No. Only producer presence rules: `EscapeClass omitempty` (set only under `plan.Escaped`) and `flow_next.go:71` `gates,omitempty` gated on `--evaluate-gates` — a *work* flag deciding whether gates run (`flow_next.go:225`), not a report-width flag. `rg plan-only internal/` returns zero hits.

**Q9 (RDR-internal).** Does the Approach's plan-group list match `flowResolveExtendedDesc` exactly?
**A9.** Yes as a set and order. Approach l.328-330 names `rule, next{}, writes{}, clear[], emit{}, escaped, escape_class, gates[]`; help `flow_resolve.go:86-102` lists the same eight. Neither lists `revision`; the Approach adds it explicitly.

**Q10 (RDR-internal).** Does any clause require the projected width to be strictly smaller than the default?
**A10.** RDR is silent. C1 fixes the key set, D-wire-byte-format fixes key deletion, A1 measures bytes on four fixtures, but no clause makes "projected < default" normative. Risks P-12 names the hazard ("a large derived field could make `--plan-only` output *bigger* … with every test green") and mitigates only with the partition oracle and key-set literal — neither a width assertion. No S-item asserts width.

**Q11 (RDR-internal).** Does C2's always-keep core get a stated purpose that survives its members being constant?
**A11.** No. C2 l.457-461 grounds the core on "never … detach a plan from the model revision that produced it"; JDR 0002 §D1 repeats the wording. That purpose presupposes `revision` discriminates. See Q1.

**Q12 (RDR-internal).** Does the RDR specify what the whole-tree oracle does with cobra's auto-generated `help`/`completion` subcommands?
**A12.** RDR is silent. C1 says "walks the WHOLE command tree from the root"; A5 asserts the shipped recursion "misses nothing". Neither names `help`/`completion` — exactly where the one shipped recursive walker stops (Q2) and where a competing `all` already lives (Q3).

## Step 3 — Findings

**F-1** — anchor `0023:C2` (also `0023:D-identity`, `0023:MVV`, JDR 0002 §D1) — class **(a) codebase claim REFUTED**.
C2 places `revision` on the PLAN side on the ground that "it is produced by the loader, not restated from the request, and it is the plan's only provenance — a chained caller cannot otherwise detect that the model changed under the same path between calls." Source refutes both halves: `revision` is a hardcoded constant `""` on every payload this CLI can emit, so it is neither loader-produced nor capable of detecting any model change.
Cite: `internal/cli/flow_exec.go:49` — `func (flowRequest) revision() string { return "" }`; doc comment l.38-40 "RDR 0002's `[model]` block admits `id`, `version`, `description`, and `metadata` and no revision field, so NO model can declare one … empty (DEV-7)"; live wire `{"type":"ok","data":{…,"revision":"",…}}`.

**F-2** — anchor `0023:C2` — class **(c) internal contradiction** (consequence of F-1).
C2's ALWAYS-KEEP core is `rule, escaped, escape_class, revision`, justified so a projected payload "can never … detach a plan from the model revision that produced it." With `revision` constant-empty, always-keep membership binds every future projection mode to carry 14 bytes (`"revision":"",`) that encode nothing — in an RDR whose purpose is byte reduction. The RDR nowhere records that `revision` is currently vacuous.
Cite: RDR l.457-461 vs `internal/cli/flow_exec.go:49`; JDR 0002 §D1 l.56-60 carries the identical unsupported rationale.

**F-3** — anchor `0023:A5` — class **(a) codebase claim REFUTED**.
A5's Evidence states "the shipped tests already walk `Commands()` recursively (the hidden `docs` command still enumerates, so the whole-tree walk misses nothing)." No shipped test recurses — every sweep is a fixed two-level loop.
Cite: `internal/cli/flow_all_0011_test.go:45,49` and `:423,438`; `internal/cli/flow_surface_0005_test.go:27,31`; `internal/cli/flow_input_0005_test.go:460,464`.

**F-4** — anchor `0023:A5` (obligation lands on `0023:C1`/`0023:S4`) — class **(a) codebase claim REFUTED** + **(d) missing requirement**.
The one recursive walker `internal/cli/help_all.go::walkCommandTree` (l.343-351) skips any child named `help` or `completion` before recursing, and its `docs.go` callers further gate on `c.Hidden`. So `docs` enumerates because the caller allows hidden commands, not because the walker is total — and a `plan-only` oracle reusing `walkCommandTree` would assert a strictly weaker negative than C1's "WHOLE command tree". C1/S4 do not say whether `help`/`completion` are in scope.
Cite: `internal/cli/help_all.go:343-351`; `internal/cli/docs.go:184`, `:245`.

**F-5** — anchor `0023:A5` — class **(b) new rule with an existing sibling**.
A5 concludes "nothing … is name-generic in a way the new sweep would disturb," but the inverse check finds a live second registration of the very flag name 0011 fences: `intrastate help --all`. The 0011 oracle passes only because it descends the `flow` group; a truly whole-tree sweep for `all` would fail today. This does not break `plan-only` (different name), but it falsifies the generalization A5 rests on and shows the whole-tree structural-negative idiom C1 mints has a pre-existing counterexample in the same repo.
Cite: `internal/cli/help_all.go::wireHelpSubcommandAll` l.231-236 — `helpCmd.Flags().Bool("all", false, "print extended help …")`.

**F-6** — anchor `0023:C1` (and `0023:S1`/`0023:S2`) — class **(d) RDR is silent → missing requirement**.
C1's differential oracle asserts "a strict key-subset … with byte-identical values on every carried key" but never asserts the projected byte width is smaller than the default. Risks names the P-12 hazard (a later field making `--plan-only` output *bigger* "with every test green") and answers it only with the partition-completeness and key-set-literal oracles, neither of which is a width assertion — so the exact green-tests failure P-12 describes remains unguarded by any S-item. A "projected bytes < default bytes on the MVV fixtures" oracle would close it.
Cite: RDR l.382-391 (C1 differential clause), l.861-871 (P-12 risk and mitigation), §Testing Strategy S1-S5.

**F-7** — anchor `0023:A3` — class **(a) minor imprecision, not load-bearing**.
A3 states "swept all 16 test files invoking `flow resolve`"; the current tree has 21 `internal/cli` test files matching. The substantive claim — none passes a projection flag, every success-payload assertion runs default mode — is CONFIRMED, so this is count drift in the Evidence prose, not a defect in the assumption.
Cite: `rg -c 'flow resolve|"resolve"' internal/cli/*_test.go` → 21 files; `rg plan-only internal/` → no matches.
