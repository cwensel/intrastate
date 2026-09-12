Model: claude-opus-5

# Grounding checks — cli/0029 design-review code claims vs `main`

Source tree: /Users/cwensel/sandbox/newcoinc/intrastate. Verified against working tree on `main`.

---

## 1. `planEnvelope` is an INPUT struct, not an output envelope constructor — HOLDS

Anchor: `internal/cli/flow_input.go::planEnvelope` (struct at line 338; doc comment lines 310-337).

- No `func planEnvelope` anywhere in the tree; the only hit for the identifier is the struct declaration and its uses in `internal/cli/flow_input.go`.
- It decodes a `--plan` document read from stdin or a file (`planFlagUsage`, `planStdinSentinel = "-"` at the top of the same span) — INPUT.
- Doc comment states the discriminator explicitly: "The discriminator is the presence of a top-level `\"type\"` key" and "`Type` must be RAW, not `*string` … Raw bytes make presence a byte-level fact".
- Refusal handling: doc says "`Code` and `Message` are the REFUSAL envelope's discriminators"; `readPlan`'s doc (≈line 358-361) adds "A REFUSAL envelope on stdin refuses."
- Fields are exactly `Type json.RawMessage`, `Data json.RawMessage`, `Code string`, `Message string`. No `schema_version` field.

## 2. `respond` and `clierr` are separate leaf packages; no `SchemaVersion` — HOLDS

Anchors: `internal/cli/clierr/clierr.go` (package doc lines 1-11), `internal/cli/respond/respond.go`.

- Separate dirs: `internal/cli/respond/{respond.go,text.go}` and `internal/cli/clierr/{clierr.go,+3 tests}`.
- Package doc: "Package clierr is the leaf home of the CLI's structured-error type. It lives in its own package so other internal packages (config, …) can construct CLIErrors without importing internal/cli and forming an import cycle."
- Repo-wide `rg "SchemaVersion|schema_version" --type go` returns ZERO hits. No exported constant of that shape anywhere.

## 3. "closed" inside a wire-visible finding-code literal — HOLDS

Anchor: `internal/graphlint/taxonomy.go:40` — `CodeCoverageClosedByEscape = "graph-coverage-closed-by-escape"`.

The word appears inside the string literal, which is emitted on the wire as a finding `code` (see `internal/graphlint/coverage.go:438` `Code: CodeCoverageClosedByEscape`, and `internal/guard/lint.go:25` which carries the SAME literal as a `Code` type).

### `internal/graphlint/taxonomy.go` — 9 occurrences (case-insensitive)

| line | text | class |
|---|---|---|
| 37 | "The tier is CLOSED at these four" | (b) comment, fixed-set sense |
| 40 | `= "graph-coverage-closed-by-escape"` | **(a) wire-visible string literal** |
| 56 | "The closed, append-only `reason` set" | (b) comment |
| 68 | "The set is closed and APPEND-ONLY, and this is the append." | (b) comment |
| 86 | "advisoryCodes is the closed four-member advisory tier." | (b) comment |
| 88 | `CodeCoverageClosedByEscape,` (slice member) | (a) identifier referencing the wire literal |
| 94 | "reasons is the closed `reason` discriminator set." | (b) comment |
| 108 | "AdvisoryCodes returns the closed four-member advisory tier." | (b) comment |
| 111 | "Reasons returns the closed `reason` discriminator set." | (b) comment |

So: 1 true (a) string literal, 1 (a)-adjacent identifier, 7 (b) comments. Zero (c).

### `internal/table/category.go` — 3 occurrences

| line | text | class |
|---|---|---|
| 52 | "A constant is not in the closed set until it is appended to" | (b) append-only/compatible sense |
| 89 | "Categories returns the closed load-category set in declaration order." | (b) |
| 158 | "than none — a reviewer who believes the class is closed stops reading" | (c) unrelated prose (cautionary aside about reviewer behaviour, not a vocabulary claim) |

Zero (a) in `category.go`.

## 4. "closed" describing vocabularies in four more files — HOLDS (counts match)

- `internal/guard/grammar.go` — **6** occurrences, claim of ~6 correct. Lines 47 ("operators is the closed typed vocabulary in matrix order"), 56 ("Operators returns this RDR's closed operator vocabulary"), 59-60 ("KnownOperator reports whether token is in the closed vocabulary. The set is closed, so a free-form expression string is simply not an operator" — two hits), 67 ("operator outside the closed vocabulary is accepted by no kind"), 145 ("the vocabulary is closed, so a free-form expression string is not an").
- `internal/accessor/model.go:308` — exact phrase: "Verdicts returns the closed three-member gate verdict vocabulary." Claim HOLDS (file has 6 `closed` hits total: 28, 46, 52, 89, 131, 308).
- `internal/resolve/resolve.go` — lines 45 and 66, as claimed. L45: "closed: RDR 0001's Normative Contracts say the kind set \"is exactly\""; L66: "RefusalKinds returns the closed kernel-owned refusal kind set."
- `internal/cli/flow_input.go:349` — exact phrase: "PUBLISHED code deliberately: RDR 0005's refusal table is closed (see `writerFor`'s doc comment)". Claim HOLDS.

**Repo-wide total, non-test Go files (excluding `docs/`): 110 occurrences across 24 files.**
Top contributors: `internal/guard/lint.go` 18, `internal/graphlint/coverage.go` 13, `internal/cli/cmdbind/cmdbind.go` 9, `internal/table/edit.go` 9, `internal/graphlint/taxonomy.go` 9, `internal/table/load.go` 8, `internal/accessor/model.go` 6, `internal/guard/grammar.go` 6, `internal/resolve/guard.go` 6, `internal/table/model.go` 5, `internal/table/category.go` 3.

## 5. Two exported `Operators()` — HOLDS, member sets IDENTICAL, both comment on not widening

- `internal/guard/grammar.go::Operators` (line 57) → `slices.Clone(operators)`, backed by `var operators` at line 50.
- `internal/table/model.go::Operators` (line 96) → `slices.Clone(operators)`, backed by `var operators` at line 93.
- Both vars are the byte-identical literal: `[]string{"eq", "in", "lt", "lte", "gt", "gte", "exists", "contains"}` — 8 members, same order. Sets are IDENTICAL.
- `internal/table/model.go:91-92` carries the not-widening comment verbatim: "operators is RDR 0003's closed operator set. This RDR does not mint operators and does not widen the set (`0002:C16`)." Claim HOLDS.
- Note (not claimed, but adjacent duplication of the same shape): `internal/guard/grammar.go:54` and `internal/table/model.go:84` both declare the same five-token kind list.

## 6. `clierr::ExitCodeFor` — PARTIAL

Anchor: `internal/cli/clierr/clierr.go::ExitCodeFor` (lines 131-149), `::ErrorGroup` (lines 27-44).

Signature is exactly `func ExitCodeFor(err error) int`, switching on `ce.Group` and ending in `return 1`. What the claim gets right and wrong:

- (a) **ErrorGroup constants: SIX**, `iota`-based on `type ErrorGroup int`: `GroupSuccess`, `GroupWarning`, `GroupUserEnv`, `GroupEnvUnavailable`, `GroupInternal`, `GroupSignalCancel`.
- (b) **`Group ErrorGroup \`json:"-"\`\`** (line 77) — NOT serialized. Struct doc line 47 says "Group drives the exit code and is not serialized." Claim HOLDS.
- (c) **There IS a comment inviting the taxonomy to grow**, at lines 25-26: "Add groups as the exit taxonomy grows — keep ExitCodeFor in sync." (There is a second growth invitation on the envelope itself at lines 48-49: "Extend with new optional fields as needed — keep them `omitempty` so the envelope stays append-only and stable for tools.")
- (d) **Exit codes the function can return: `{0, 1, 2, 3, 130}`** — five distinct values. `0` also on `err == nil` (line 133); `1` is the non-CLIError / unmatched-group default.
- (e) **PARTIAL against the claim's framing**: `clierr` exports NO enumerator over groups or codes. Exported functions in `clierr.go` are exactly: `ErrorCode`, `ExitCodeFor`, `EmitJSON`, `WriteJSONLine`, `EmitText`, `EmitFindingsText`, plus methods `(*CLIError).Error` and `(*CLIError).Unwrap`. No `Groups()`, `ExitCodes()`, or `Codes()`.

Why PARTIAL and not HOLDS: the claim's bullet-(c) framing ("whether there is any comment inviting the taxonomy to grow") is answered YES — a reviewer reading the claim as "there is no such comment" would be wrong. Everything else in the claim is accurate. The switch has no `default:` clause; the `return 1` is a bare post-switch statement, which is behaviourally the same but not literally a `default` arm.

## 7. No enumeration of the CLIError `code` vocabulary — HOLDS

- No `func Codes()` and no `codes = []string{…}` in the CLI-error path. The ONLY `Codes()` in the tree is `internal/guard/lint.go:41` — `func Codes() []Code { return slices.Clone(codes) }` over `var codes []Code` (line 30) — which enumerates the GRAPH-LINT finding codes, a different vocabulary from CLIError's `code`.
- How codes are actually produced: at raise sites, as a `Code:` field on a composite literal. Repo-wide non-test, 71 `Code:` assignments — **10 are bare string literals** (e.g. `internal/cli/root.go:324` `Code: "command-error"`) and **61 reference package-level constants** (`graphlint.AggregateCode`, `guard.Code*`, `accessor` validation codes, etc.). So the vocabulary is a union of per-package constant blocks plus a tail of inline literals, with no single accessor over it and no single declaration site.

## 8. `respond.Success` marshals four fields — HOLDS

Anchor: `internal/cli/respond/respond.go::Success` (lines 51-56). Exact field list and tags:

```
Type     string     `json:"type"`
Notes    []Advisory `json:"notes,omitempty"`
Warnings []Warning  `json:"warnings,omitempty"`
Data     any        `json:"data,omitempty"`
```

`Type` is non-omitempty; the other three are `omitempty`. Exactly as claimed.

## 9. `CLIError.Hint` is `omitempty` — HOLDS

Anchor: `internal/cli/clierr/clierr.go::CLIError` (lines 50-88). Full marshalled field set:

```
Code     string    `json:"code"`
Message  string    `json:"message"`
Param    string    `json:"param,omitempty"`
Detail   string    `json:"detail,omitempty"`
Hint     string    `json:"hint,omitempty"`
Findings []Finding `json:"findings,omitempty"`
```

Plus two NON-marshalled fields: `Group ErrorGroup \`json:"-"\`` and `Cause error \`json:"-"\``. `Hint` is present with `omitempty` (line 60), claim HOLDS.

## 10. `respond` package doc vs `Fail`'s actual behavior — HOLDS (doc/behavior divergence confirmed)

Anchors: `internal/cli/respond/respond.go` package doc (lines 1-25), `::Success` doc (lines 48-50), `::Fail` (lines 181-192), `internal/cli/clierr/clierr.go::EmitJSON` (lines 151-157).

Doc text, verbatim:
- Lines 12-13: the discriminated NDJSON block, including `{"type": "failed", ...}   # graceful failure`.
- Lines 23-24: "The terminal \"ok\"/\"failed\" type names are reserved."
- Lines 48-49 (`Success` doc): "Type is always \"ok\" (Fail emits \"failed\")."

Actual behavior: `Fail` under `ModeJSON` calls `clierr.EmitJSON(cmd.OutOrStdout(), ce)`, which is `WriteJSONLine(out, e)` over the **bare `*CLIError`**. `CLIError` has no `Type` field and no `type` json tag (see claim 9) — so the emitted failure line carries `code`/`message`/… with **NO `type` key at all**. The doc's `{"type": "failed", ...}` shape and the "Fail emits \"failed\"" parenthetical are both WRONG about the wire.

Corroborating: `internal/cli/flow_input.go`'s own `planEnvelope` doc (lines 331-337) states the correct behavior — "A failed run under `--as=json` emits the bare `CLIError` — `{\"code\":…,\"message\":…}` … so a refusal carries NO `type` key".

**Second half of the claim also HOLDS**: `readPlan`'s doc comment (`internal/cli/flow_input.go`, ≈lines 357-360) says "a pipe that lost its plan carries `{\"type\":\"failed\",…}` instead" — the same stale `type: failed` shape, in a SECOND file. So the wrong shape is documented in two places while the code in the same file documents the right one 20 lines above.

## 11. Advisory `level` values are bare literals, no `Levels()` — HOLDS

Anchor: `internal/cli/respond/respond.go` lines 207 and 217.

- L207 (`Note`): `_ = writeJSONLine(stderr, map[string]string{"level": "note", "message": message})`
- L217 (`Warn`): `payload := map[string]string{"level": "warning", "message": w.Message}`
- Both `"note"` and `"warning"` are inline string literals inside `map[string]string{…}`. No `Levels()` accessor and no `levels` var anywhere in the package. The only other mention is the package doc line 17: `advisories only, discriminated by "level" ("note" | "warning")`.
- `internal/cli/respond/text.go` carries no `level` handling.

## 12. `graphlint.AggregateCode` and the tier accessors — HOLDS

Anchor: `internal/graphlint/taxonomy.go`.

- Line 48: `const AggregateCode = "graph-lint-failed"` — a single untyped const, not a set. Doc (lines 46-47): "AggregateCode is the one aggregate CLIError code a blocking run returns (`0006:C14`)."
- Line 106: `func BlockingCodes() []string { return slices.Clone(blockingCodes) }`
- Line 109: `func AdvisoryCodes() []string { return slices.Clone(advisoryCodes) }`
- Both return `slices.Clone` of unexported package vars (`blockingCodes`, `advisoryCodes`). Two more of the same shape exist: `Reasons()` (line 112) and `Severities()` (line 115).
- `severityFor` (line 123) does NOT call `slices.Contains` directly — it delegates: `if IsBlocking(code) { return SeverityBlocking }; return SeverityInfo`. `IsBlocking` (line 120) is `slices.Contains(blockingCodes, code)`. **The derivation is via `slices.Contains(blockingCodes, code)` as claimed, but through one level of indirection.** Substance HOLDS.

## 13. Test symbols and slice accessors — HOLDS

- `internal/graphlint/findings_0006_test.go:287` — `func TestReq74_TheAdvisoryTierIsClosedAtExactlyFourMembers`. The four-element `want` literal is at lines **288-293** (claim said ~287-293), holding `"graph-coverage-closed-by-escape"`, `"graph-redundant-row"`, `"graph-unreachable-rule"`, `"graph-vacuous-atom"`, compared with `slices.Equal` against a sorted `AdvisoryCodes()`.
- `internal/resolve/guard_atoms_test.go:894` — `func TestReq79_NoSixthRefusalKindIsMinted`.
- `internal/table/command_carrier_0025_test.go:924` — `func TestReq79_TheCategorySetIsAppendOnlyAndDuplicateFree`.
- `internal/resolve/resolve.go:67` — `func RefusalKinds() []RefusalKind` returning a composite literal (`KindNoMatch`, `KindAmbiguousMatch`, …).
- `internal/table/category.go:90` — `func Categories() []Category` returning a composite literal (`CatMalformedTOML`, `CatUnknownSchemaField`, …).
- `internal/accessor/model.go:309` — `func Verdicts() []Verdict { return slices.Clone(verdicts) }`.

Note the two accessor IDIOMS in play across the repo: `slices.Clone(pkgVar)` (graphlint, guard, table/model, accessor) versus a freshly-built composite literal with no backing var (resolve, table/category). Both are safe against caller mutation; they are not the same shape.

## 14. MVV step-1 golden — HOLDS, and the assertion IS byte-identity (modulo two documented folds)

- File exists: `docs/rdr/0023-resolve-envelope-projection/artifacts/mvv-step1-default-golden.json`, 291 bytes.
- Path constant: `internal/cli/flow_mvv_0023_test.go:43` — `const goldenPath0023 = "docs/rdr/0023-resolve-envelope-projection/artifacts/" + "mvv-step1-default-golden.json"`.
- Read at line **82** (claim said ~44; line 44 is the second half of the path constant, the READ is at 82): `path := filepath.Join(repoRootFor(t), goldenPath0023); want, err := os.ReadFile(path)`, inside `TestReq27And28And95And96And97And129_DefaultModeIsByteIdenticalToTheCapturedPreChangeGolden` (line 81).
- Assertion (line ~119): `if got != strings.TrimRight(string(want), "\n")` — a whole-string equality, i.e. **byte-identity**, not a field-subset or unmarshal-and-compare.
- Two transforms are applied to `got` before the compare, both documented in-place as normalization rather than narrowing: `foldCheckoutRoot` (folds the absolute checkout prefix out of the `model` value so the golden can carry the relative spelling) and `foldEmitDispositions0024`. The test comment states explicitly "This narrows nothing. Byte-identity is still asserted over the whole record, `model` included."

## 15. `docs/cli-output-contract.md` — HOLDS on the three words; NO release/versioning language; does NOT document a `type` key on failure

File exists, 26325 bytes.

- `frozen`: **0**. `append-only`: **0**. `growing`: **0**. All three counts zero, claim HOLDS.
- Release/upgrade/breaking/versioning language: `upgrade` 0, `breaking` 0, `release` 0, `compat` 0. `version` has **1** hit, line 320, and it is NOT about the wire: "A declaration is author-owned and unversioned: widen or narrow a domain by …". So the document carries **no release, upgrade, breaking-change, or versioning language at all**.
- Failure envelope: the doc does NOT document a `type` key on failure. It describes `findings` as "a top-level array on the `CLIError` envelope — a **sibling** of `code`, not nested under an `error` wrapper" (lines 21-23) — i.e. the bare-`CLIError` shape the code actually emits, matching claim 10's finding and contradicting `respond.go`'s package doc.
- The doc explicitly defers to the binary: "**The binary states the contract itself.** `intrastate --help-all` … This document does not repeat that reference" (lines 5-10).

## 16. `llms.txt` — HOLDS

File exists at repo root, 2170 bytes.

- Opening instruction, verbatim (lines 11-13): "The binary is authoritative and self-describing. Prefer asking it over reading any file here — a file is a snapshot, the binary is the build you actually have:" followed by a fenced `sh` block: `intrastate --help-all`, `intrastate <cmd> --help-all`, `intrastate <cmd> --as=json  # one terminal JSON envelope per run`.
- Link to the contract, line 39-40, one-line description verbatim: "- [Output contract](docs/cli-output-contract.md): worked JSON payloads and the rationale behind the envelope shape."
- Also lists `intrastate version`: "Print build version, commit, and date" (Commands section).

## 17. `version` verb emits build identity under the `ok` envelope — HOLDS

Anchor: `internal/cli/version.go::newVersionCmd`.

- Gateway: `return respond.OK(cmd, respond.Success{Data: version.Get()})` — BOTH modes route through `respond`, no direct printing. The doc comment (lines 9-12) says "in json mode it emits the structured Info under the terminal \"ok\" envelope"; the Long help repeats it.
- Payload: `internal/version/version.go::Info` — `Version string \`json:"version"\``, `Commit string \`json:"commit"\``, `Date string \`json:"date"\``.
- **The collision the claim points at is real**: on `intrastate version --as=json` the wire is `{"type":"ok","data":{"version":"…","commit":"…","date":"…"}}`. A new TOP-LEVEL `version` key on the envelope would sit as a sibling of `type`/`data` while `data.version` already means "the build's release identity" — two different senses of the same word one nesting level apart, on the very verb a caller reaches for to pin a build.
- Supporting: `versionExtendedDesc` states "the finding codes, the refusal vocabulary, and the lint bounds are all properties of a particular build, and `--help-all` on this binary is the authoritative statement of what THIS build does."

## 18. `flow_next.go` unknown-`reason` union, no accessor — HOLDS

Anchor: `internal/cli/flow_next.go` const block at lines **91-95** (claim said ~92-93; the block spans 91-95).

```
const (
	reasonAbsent       = string(resolve.ReasonAbsent)
	reasonUncomparable = string(resolve.ReasonUncomparable)
	reasonNotEvaluated = "not-evaluated"
)
```

- Two members are re-spellings of `resolve.Reason*` kernel constants; `reasonNotEvaluated` is a locally-declared bare string literal, described by the preceding comment (lines 86-90) as "`not-evaluated` is the ONE token RDR 0011 adds, scoped to un-run gate ids on this CLI list".
- The comment calls the union "The closed reason vocabulary" but there is **no accessor over the union** — no `Reasons()`, no backing slice, in `internal/cli`. The three constants are consumed only at three raise sites (lines 352, 360, 377) as `unknownFact{Key: …, Reason: …}` values. A consumer cannot enumerate the union from the binary through this package.
- Contrast: `internal/graphlint` DOES export `Reasons()` over its own `reason` set (taxonomy.go:112) — so the repo has the accessor idiom for exactly this shape, just not applied here.

---

## Tally

- **HOLDS: 17** — claims 1, 2, 3, 4, 5, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18.
- **PARTIAL: 1** — claim 6 (`ExitCodeFor`). Right: signature, six `ErrorGroup` constants, `json:"-"`, exit-code set `{0,1,2,3,130}`, no enumerating accessor. Wrong/incomplete: there IS a growth-inviting comment ("Add groups as the exit taxonomy grows — keep ExitCodeFor in sync.", clierr.go:25-26), and the `return 1` is a post-switch statement rather than a `default:` arm.
- **FAILS: 0.**

Minor line-number drift, non-substantive: claim 13's `want` literal is 288-293 not 287-293; claim 14's golden read is line 82 not ~44 (44 is the path constant's continuation); claim 18's const block is 91-95 not 92-93.

Substantive finding for the RDR's benefit (claim 10): the `{"type":"failed",…}` failure shape is documented in TWO places — `internal/cli/respond/respond.go` package doc (lines 12, 23-24, 48-49) and `internal/cli/flow_input.go`'s `readPlan` doc — and is emitted in NEITHER. `docs/cli-output-contract.md` and `planEnvelope`'s own doc both describe the correct bare-`CLIError` shape.
