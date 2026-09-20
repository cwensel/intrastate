# RDR 0012 — author's round

Answers land in `rulings.md` beside this file, one
`- **<Qn|fixture>** — RULED: <verbatim>` line each under a
`## <date> — resolve` heading. The re-run reads them before the record.

## 2026-09-20 — resolve

- **Q1 A4 is refuted on its literal side — how should the claim narrow?** Typed `int` comparison parses BOTH sides (C2), but the literal is never rendered from the declaration: `valueSatisfies` joins the authored bytes and the loader's kind check is `strconv.Atoi` (`internal/table/load.go::conformKind`), which ACCEPTS `"00"`, `"01"`, `"+1"` (it rejects `" 1"` and `"many"`). So an authored `n eq "00"` flips lint from blocking `graph-coverage-gap` (exit 2) to clean (exit 0) — blocking→clean, the worst direction. Three options: **(a)** canonicalize `int` literals at load (reject `"00"`/`"+1"`) so the literal side cannot drift, which adds a load refusal to RDR 0003's conformKind; **(b)** restate A4 to admit authored non-canonical literals DO flip lint, and pin the flip as accepted behavior with a regression case; **(c)** take the recorded fallback in §Load-Bearing Decisions — validate-then-byte-compare (parse to prove comparability, compare raw bytes), which A4's If-wrong already names. — grounding: spike `evidence/spikes/a4-lint-diff.md` (fixture `cover07`, both binaries from `6c22a12`); held side CONFIRMED safe at `internal/guard/assignment.go` int arm (`strconv.Itoa`, no leading zero/`+`/whitespace) and `internal/guard/product.go::valueAssignments`; §Load-Bearing Decisions "Selection / predicate" already records (c) as the fallback if A4's lint-agreement verification fails — it now has.

- **Q2 The MVV's step 2 cannot be satisfied as written — which venue does the defect actually live in?** MVV step 2 expects `flow resolve --tag iter=many` to refuse `flow-guard-unevaluable` with reason `uncomparable`. Measured, it refuses `flow-tag-invalid` at CLI input instead: `internal/cli/flow_input.go` calls `table.ConformValue` on every DECLARED `--tag` before resolution, so the value never reaches the seam. The `GuardFalse → GuardUnevaluable` leg the RDR calls "the fix" is unreachable via `--tag` on a declared `int`. The unconformed ingress is the reader-supplied OWNED value (`internal/resolve/resolve.go` merges `ProvenanceOwned` unconformed; nothing in `internal/accessor` conforms). Either re-author the MVV onto an owned-value reader path, or re-scope the Problem Statement to name the owned ingress as the defect's home. — grounding: spike `evidence/spikes/a4-typed-compare.md` (all four `--tag` routes run against both binaries); JDR 0004 §D2's ingress census independently reaches the same table — "Reader-supplied owned values | **no** | nothing in `internal/accessor` conforms".

- **Q3 JDR 0004 binds this record and settled three clauses after its last write — apply them?** JDR 0004 (settled 2026-09-20; `Binds: 0012, 0030`) postdates the RDR's last refine (`6c22a12`) and the record cites JDR 0001 §JD-18 but carries no citation of 0004. Three landings are named in the JDR's own text: **JD-1** — "*Lands in 0012*: C1's closing sentence is scoped to the resolution path", i.e. C1's "the kernel (`internal/resolve`) continues to carry no declarations" must scope to `Resolve` and its callees, not the package as a namespace; **JD-3** — "**0012** gains the loader as a declared construction site", a FOURTH site (`internal/table/normalize.go::renderWrites`, a `*loader` method holding `l.model.Tags`) that C4 and A3 currently do not list; **JD-2** — lands in neither record, but "0012 may cite this entry where it states the widening's reach". Confirm all three, or name which to hold. — grounding: `docs/jdr/0004-declaration-carrying-package-and-undecidable-value-venue.md` §D1/§D2/§D3 and its Interface record; `internal/table/normalize.go::renderWrites` exists on `main` (it becomes a construction site only once 0030's shim extraction lands, so C4's list is correct for today's tree and incomplete against the settled forward obligation); `rg` confirms exactly three `Evaluator{}` sites on `main`.

- **fixture F1 — runtime non-canonical held value flips a plan into a refusal.** Approve as the normative fixture for C2's parsed-comparison leg.
  Input: model `mvv-int` (decision-table; `iter` int min 0 max 9 observed single-valued required; row `guarded` guarded `iter eq 7`; unguarded row `fallback`), invoked
  `intrastate flow resolve --model mvv-int.toml --outcome step --tag iter=07 --plan-only --as json`.
  Expected TODAY (raw-string seam): exit 0,
  `{"type":"ok","schema_version":"0.1","data":{"revision":"","rule":"fallback","gates":[],"emit":{"next":"fallback"},"dispositions":{"next":"go"},"next":{},"writes":{},"clear":[],"escaped":false}}`
  Expected AFTER typed comparison: exit 2,
  `{"code":"flow-ambiguous-match","message":"more than one rule matches the recognized outcome `step`; resolve does not choose among them","schema_version":"0.1","findings":[{"code":"flow-ambiguous-match","message":"rule `fallback`: this rule matches and conflicts with the others named here","locator":"mvv-int:fallback","rule":"fallback"},{"code":"flow-ambiguous-match","message":"rule `guarded`: this rule matches and conflicts with the others named here","locator":"mvv-int:guarded","rule":"guarded"}]}`
  — grounding: spike `evidence/spikes/a4-typed-compare.md`; values read from the run, not invented. Confirms A4's predicted `GuardFalse → GuardTrue` at the CLI.

- **fixture F2 — the MVV's malformed case is refused UPSTREAM, not at the seam.** Approve as the normative fixture pinning where `--tag` conformance fires (it is the Expected for Q2's re-authoring, whichever venue wins).
  Input: same `mvv-int` model; `intrastate flow resolve --model mvv-int.toml --outcome step --tag iter=many --plan-only --as json`.
  Expected, BOTH before and after the change: exit 2,
  `{"code":"flow-tag-invalid","message":"the value for `iter` does not conform to its declaration: \"many\" is not an int","schema_version":"0.1","param":"iter"}`
  — grounding: spike `evidence/spikes/a4-typed-compare.md`. This is what the record's MVV step 2 currently mis-predicts as `flow-guard-unevaluable`/`uncomparable`.

- **fixture F3 — authored non-canonical int literal flips lint from blocking to clean.** Approve as the normative fixture for Q1's regression case, under whichever narrowing is ruled.
  Input: model `cover07` (decision-table; `[tags.n]` int min 0 max 1 observed single-valued required; `r-zero` guarded `n eq "00"` emitting `zero`; `r-one` guarded `n eq 1` emitting `one`), invoked `intrastate lint --model cover07.toml --as json`.
  Expected TODAY (raw-string seam): exit 2,
  `{"code":"graph-lint-failed","message":"the model carries blocking graph-lint findings","schema_version":"0.1","findings":[{"code":"graph-coverage-gap","message":"group cover07/step over rows [r-one r-zero] leaves 1 of 2 assignments in its scoped product uncovered for the no_match arm; the coverage union must equal the scoped product","model":"cover07","severity":"blocking","rule":"r-one","element":"cover07/step","dimension":"n","class":"no_match"}]}`
  Expected AFTER typed comparison: exit 0, `{"type":"ok","schema_version":"0.1","data":{"findings":[]}}`
  — grounding: spike `evidence/spikes/a4-lint-diff.md`; both binaries built from `6c22a12` in a throwaway worktree. This is the blocking→clean flip that refutes A4 as written.

### Context the round does not ask about

The corpus lint diff came back clean — 123 models (`models/` plus every
lintable `*.toml` under `internal/*/testdata`), no verdict flip, `go test
./...` green against the typed scaffold. That result does NOT support A4:
the corpus contains zero `eq`/`in` atoms over an `int` tag (the single
`kind = "int"` declaration, `models/examples/release-grammar.toml`
`[tags.risk]`, is guarded only by `gte`/`lt`/`gt`/`lte`, which already
parse both sides). It is a coverage statement, not a safety one, and is
recorded that way in §Validation scenario 5.
