Model: claude-sonnet-4-5

No findings.

## Coverage

### Widened spans

- Widened past the pre-scoped anchor list to read the full post-propose diff
  (`/tmp/0028-post-propose.diff`, 885 lines) end to end, since the C1-C6→C1.1-C1.6
  merge touches nearly every clause and the diff itself is the authoritative
  post-propose delta.
- Widened into `docs/rdr/0028-declared-line-edit-writer/evidence/spikes/*` (A3,
  A4, A5 spike outputs) because the RDR's post-propose A3/A4/A5 evidence blocks
  cite these files as the source of the "Verified" status and the wrapped-
  qualifier claim requires reading the spike's raw finding to check whether the
  resolved RDR text's "SAFE, not truncating" framing is a real correction of
  the spike's own "genuine boundary finding" language, or a walk-back that
  contradicts it.
- Widened into `docs/rdr/0028-declared-line-edit-writer/evidence/research/prior-art.md`
  to determine whether `pkg/plugin/util/util.go::InsertCode` (in the
  pre-computed resolved-anchor list) is an in-repo intrastate symbol or
  external prior art, since it does not exist anywhere under `intrastate/`.
- Widened into the sibling `rdr` engine repo
  (`/Users/cwensel/sandbox/newcoinc/rdr/tools/rdr/`) and its records
  (`/Users/cwensel/sandbox/newcoinc/intrastate/docs/rdr/README.md`,
  `/Users/cwensel/sandbox/newcoinc/rdr/models/rdr-write.toml`,
  `/Users/cwensel/sandbox/newcoinc/rdr/tools/rdr/internal/scan/fields.go`,
  `/Users/cwensel/sandbox/newcoinc/rdr/tools/rdr/internal/model/template.go`)
  because A3/A4's claims are about "the consumer" — the RDR record-keeping
  system itself, whose reader/model code lives in that repo, not in
  `intrastate`'s Go source tree. The task's source-repo binding names
  `intrastate` as the tree every `path::Symbol` anchors into, but several
  post-propose claims are explicitly about consumer-side behavior
  (`rdr status -json -filter status`, `model.ValueContinues`,
  `internal/scan/fields.go`) that only exist in the sibling repo; grounding
  them required reading there.

### Unanchored prose claims grounded — CONFIRMED

- "No accessor dump/marshal/normalize round-trip exists anywhere in
  `internal/`" (A6) → `internal/table/dump.go::Dump`/`DumpAll` render candidate
  ROWS (identity/source/kind/atoms/writes/…), never an `Accessor`/carrier
  field; `internal/table/model.go::seamValue` encodes tag VALUES only. No
  dump/marshal path touches `Path`/`Command`/`Edit`.
- "`checkAccessorBindings` is carrier-agnostic (arity only)" (A6) →
  `internal/table/load.go:1265` walks `l.model.Readers`/`Writers` by `.Keys`
  count only; no carrier-type branch.
- The loader's `sourceAcc` field copy is "ADDITIVE per-field copying, not an
  exclusionary branch" (A6) → `internal/table/load.go:1004-1009`,
  `if a.Path != nil {...}` then `if a.Command != nil {...}`, confirmed.
- "No other change to 0025:C1–C1.6" / "searched the repo for an existing
  line-selection signal to reuse: none exists (no non-test `regexp` import
  under `internal/`)" (Load-Bearing Decisions) → `grep` for `"regexp"` imports
  under `internal/` returns only four `_test.go` files.
- `flowbind.go::save` fixes mode 0600 and is the only stage-and-rename site in
  the repo (A5 / Existing Infrastructure Audit) → confirmed at
  `internal/cli/flowbind/flowbind.go:139-166`; repo-wide grep for
  `os.Rename`/`Chmod` under `internal/`,`pkg/` returns only this file.
- `cmdbind.go::substitute` compares whole argv elements only (A9) →
  `internal/cli/cmdbind/cmdbind.go:400-408`, `if el != ArtifactPlaceholder`;
  it is the only substitution/placeholder site in the repo (grep for
  `func substitute`/`Placeholder`/`Template` under `internal/` returns only
  this one).
- "`Executor.Write` resolves the role's reader... BEFORE it calls `Apply`"
  (A2) → `internal/accessor/executor.go`: `readerFor` at line 302,
  `binding.Apply` at line 355, both inside one function body, confirmed order.
- "today's gate check lives inside `internal/cli/cmdbind::spawn` and runs only
  during execution" (A2) → `internal/cli/cmdbind/cmdbind.go:170`,
  `if !cfg.AllowCommands { return invocation{}, refuse(...) }` inside `spawn`,
  which only runs at command-execution time, not at reader resolution.
- `internal/accessor/model.go::Artifact` is `{Role, Path}` with no tag/context
  field (A1) → confirmed at `internal/accessor/model.go:249-252`.
- `WriteBinding.Apply`/`ReadBinding.Read` both take `(ctx, art Artifact, ...)`
  — same `art` type reaches both (A1) → confirmed at
  `internal/accessor/binding.go:40,57`.
- `registerTagFlag` registers `--tag` as "observed tag" and `parseTags`
  refuses an owned key (A7) → confirmed at `internal/cli/flow.go:182-185` and
  `internal/cli/flow_input.go:606-630`.
- JDR 0003 §D1 says "three rules, no new field, no new class" and names
  `0028:A8` by id (A8) → confirmed by reading
  `docs/jdr/0003-accessor-binding-seam.md` §D1 through the projector.
- The two README index rows for 0001 and 0006 are carried in the unlinked
  `| NNNN |` form while every other row (0002-0028) is linked
  `[NNNN](NNNN-slug.md)` (new Failure Modes / Load-Bearing Decisions bullet)
  → confirmed directly in `intrastate/docs/rdr/README.md:12` (`| 0001 | ... |`)
  and `:17` (`| 0006 | ... |`) against the linked form at every other row.
- `rdr-write.toml`'s `readme-add` rule emits the linked
  `| [NNNN](NNNN-slug.md) | <Title> | <Status> | <Priority> |` form (MVV
  fixture note) → confirmed at `rdr/models/rdr-write.toml:459` (`readme-add`
  rule's `edit` field, verbatim).
- The wrapped-qualifier fixture claim (A3) — a Status qualifier that continues
  onto a second physical line, and that `${2}`-preserving `replace` is safe
  rather than truncating — traced against the raw spike output
  (`evidence/spikes/a3-anchor-regex.md`) and the fixture
  `rdr/tools/rdr/testdata/status/records/0021-cache-warmup-order.md:9-10`.
  The spike's own prose calls this a "genuine boundary finding" using
  "truncat[ion]" language, but its own demonstrated mechanism (`${2}` captures
  and re-emits exactly the pre-edit line-9 bytes verbatim; line 10 is never
  touched) is the SAME mechanism the resolved A3 text describes as "SAFE, not
  truncating." The apparent tension is a framing correction between spike and
  resolve, not a factual contradiction: nothing the edit does drops or
  re-derives content — `${N}` is a literal capture-group copy, not a semantic
  qualifier reconstruction, so the resolved framing is the more accurate one.
  Not filed as a finding.
- `pkg/plugin/util/util.go::InsertCode`/`::ReplaceInFile` do not exist
  anywhere under `intrastate/`; traced to `evidence/research/prior-art.md`
  query #5, "semble `langref/kubebuilder`" — explicit EXTERNAL prior art
  (kubebuilder), and the RDR text (Background/References sections) correctly
  labels it "kubebuilder `pkg/plugin/util/util.go::InsertCode`" both places
  it is cited. Not an in-repo sibling misattribution; not filed as a finding.

### Inverse check — new discriminators/rules vs. existing siblings

1. **Carrier-selection/discrimination site** — `internal/table/load.go::carrierDefect`
   (the `hasPath`/`hasCommand` switch, confirmed at line 1058) and
   `internal/cli/flowbind/registry.go::commandBacked` (boolean discriminator,
   confirmed at line 100) ARE the existing carrier-decision sites, and the RDR
   explicitly extends them (its own stated design) rather than inventing a
   parallel decision site. No undisclosed sibling found.
2. **Line-oriented in-place text editing with anchor matching** — searched
   `internal/` and `pkg/` for any existing regexp-based line editor: none
   exists in-repo (only test-file `regexp` imports, see above).
   `pkg/plugin/util/util.go::InsertCode` is external kubebuilder prior art
   (confirmed above), correctly disclosed as such, not an intrastate sibling.
3. **Atomic stage-and-rename with mode preservation** — `internal/cli/flowbind/flowbind.go::save`
   is the only stage-and-rename site in the repo (grep confirmed), and it is
   fixed-0600, not mode-preserving. No sibling already does mode-preserving
   atomic writes; the RDR correctly frames this as new.
4. **`{placeholder}` substitution with escape rules** — `internal/cli/cmdbind/cmdbind.go::substitute`
   is the only placeholder-substitution site in the repo (grep for
   `func substitute`/placeholder/template under `internal/` returns only this
   one); `internal/table/load.go::commandPlaceholders` (a var, the closed
   list `[]string{"{artifact}"}`, confirmed at line 1029) is the lint-time
   companion the RDR extends. No second substitution mechanism found
   elsewhere.
