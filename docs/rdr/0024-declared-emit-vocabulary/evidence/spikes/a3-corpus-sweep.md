Model: claude-opus-5[1m]

# A3 corpus sweep — `[emit.<key>]` declarations and `[rule.emit]` authorship

Repo HEAD at capture: `d771278` (branch `main`, clean tree).

## 1. Corpus sweep for `[emit.` declarations

```
$ rg -n --glob '*.toml' '^\s*\[emit[.\]]' .
```

```
(no output; rg exit status 1 = no matches)
```

```
$ rg -n --glob '*.toml' '^\s*\[emit' .
```

```
(no output; rg exit status 1 = no matches)
```

Hit count: **0** for both patterns, across every `.toml` in the repo (models, all
`internal/**/testdata`, `docs/`, and repo-root config TOML).

## 2. Corpus sweep for `[rule.emit]`

```
$ rg -n --glob '*.toml' '\[rule\.emit\]' .
```

```
./models/examples/pricing-decision-table.toml:18:#   - every ordinary row carries [rule.emit] and no write block
./models/examples/pricing-decision-table.toml:59:[rule.emit]
./models/examples/pricing-decision-table.toml:71:[rule.emit]
./models/examples/pricing-decision-table.toml:83:[rule.emit]
./models/examples/pricing-decision-table.toml:95:[rule.emit]
./docs/rdr/0010-stateless-decision-tables/evidence/spikes/a8-rule-emit-table.toml:53:[rule.emit]
```

```
$ rg -n --glob '*.toml' 'rule\.emit' .
```

```
./models/examples/pricing-decision-table.toml:18:#   - every ordinary row carries [rule.emit] and no write block
./models/examples/pricing-decision-table.toml:59:[rule.emit]
./models/examples/pricing-decision-table.toml:71:[rule.emit]
./models/examples/pricing-decision-table.toml:83:[rule.emit]
./models/examples/pricing-decision-table.toml:95:[rule.emit]
./docs/rdr/0010-stateless-decision-tables/evidence/spikes/a8-rule-emit-table.toml:53:[rule.emit]
```

(The dotted/inline form `rule.emit = ...` produces no additional hits; the two
patterns return identical result sets.)

### Full TOML inventory outside `docs/`

```
$ find . -name '*.toml' -not -path './docs/*' | sort
```

```
./.kata.toml
./.roborev.toml
./internal/table/testdata/delim/merge-delim-comma.toml
./internal/table/testdata/delim/merge-delim-empty.toml
./internal/table/testdata/delim/merge-delim-pipe.toml
./internal/table/testdata/delim/merge-delim-semi.toml
./internal/table/testdata/delim/merge-delim-space.toml
./internal/table/testdata/delim/write-delim-comma-a.toml
./internal/table/testdata/delim/write-delim-comma-b.toml
./internal/table/testdata/delim/write-delim-empty-a.toml
./internal/table/testdata/delim/write-delim-empty-b.toml
./internal/table/testdata/delim/write-delim-pipe-a.toml
./internal/table/testdata/delim/write-delim-pipe-b.toml
./internal/table/testdata/delim/write-delim-semi-a.toml
./internal/table/testdata/delim/write-delim-semi-b.toml
./internal/table/testdata/delim/write-delim-space-a.toml
./internal/table/testdata/delim/write-delim-space-b.toml
./internal/table/testdata/dup/dup-model-a.toml
./internal/table/testdata/dup/dup-model-b.toml
./internal/table/testdata/dup/dup-model-distinct.toml
./internal/table/testdata/kata-fixture.toml
./internal/table/testdata/merge-delim-atom.toml
./internal/table/testdata/merge-distinct-rev.toml
./internal/table/testdata/merge-distinct.toml
./internal/table/testdata/merge-idempotent.toml
./internal/table/testdata/neg/neg-accessor-no-timeout.toml
./internal/table/testdata/neg/neg-alphabet-hash.toml
./internal/table/testdata/neg/neg-bad-exists.toml
./internal/table/testdata/neg/neg-case-folded-tag.toml
./internal/table/testdata/neg/neg-clear-as-write.toml
./internal/table/testdata/neg/neg-clear-in-initial.toml
./internal/table/testdata/neg/neg-cyclic-context.toml
./internal/table/testdata/neg/neg-dump-omitted-column.toml
./internal/table/testdata/neg/neg-dump-repeated-column.toml
./internal/table/testdata/neg/neg-dump-unknown-column.toml
./internal/table/testdata/neg/neg-dup-alphabet-member.toml
./internal/table/testdata/neg/neg-duplicate-rule-id.toml
./internal/table/testdata/neg/neg-empty-alphabet-member.toml
./internal/table/testdata/neg/neg-escape-no-match-block.toml
./internal/table/testdata/neg/neg-escape-with-clear.toml
./internal/table/testdata/neg/neg-escape-with-empty-write.toml
./internal/table/testdata/neg/neg-escape-with-gate.toml
./internal/table/testdata/neg/neg-escape-with-write.toml
./internal/table/testdata/neg/neg-gate-is-reader.toml
./internal/table/testdata/neg/neg-guard-all-bad-operator.toml
./internal/table/testdata/neg/neg-guard-all-clear.toml
./internal/table/testdata/neg/neg-guard-unless-bad-operator.toml
./internal/table/testdata/neg/neg-guard-unless-clear.toml
./internal/table/testdata/neg/neg-guard-unless-out-of-domain.toml
./internal/table/testdata/neg/neg-guard-unless-unknown-tag.toml
./internal/table/testdata/neg/neg-in-member-hash.toml
./internal/table/testdata/neg/neg-initial-bad-value.toml
./internal/table/testdata/neg/neg-initial-int-out-of-range.toml
./internal/table/testdata/neg/neg-initial-unknown.toml
./internal/table/testdata/neg/neg-literal-outside-domain.toml
./internal/table/testdata/neg/neg-malformed-toml.toml
./internal/table/testdata/neg/neg-match-exists.toml
./internal/table/testdata/neg/neg-match-lt.toml
./internal/table/testdata/neg/neg-no-alphabet.toml
./internal/table/testdata/neg/neg-no-model-id.toml
./internal/table/testdata/neg/neg-no-model-table.toml
./internal/table/testdata/neg/neg-no-recognized-decl.toml
./internal/table/testdata/neg/neg-no-version.toml
./internal/table/testdata/neg/neg-no-write-block.toml
./internal/table/testdata/neg/neg-outcome-outside.toml
./internal/table/testdata/neg/neg-owned-no-reader.toml
./internal/table/testdata/neg/neg-owned-two-readers.toml
./internal/table/testdata/neg/neg-recognized-declared-twice.toml
./internal/table/testdata/neg/neg-recognized-in-guard.toml
./internal/table/testdata/neg/neg-recognized-in-keys.toml
./internal/table/testdata/neg/neg-recognized-misnamed.toml
./internal/table/testdata/neg/neg-recognized-observed.toml
./internal/table/testdata/neg/neg-recognized-owned.toml
./internal/table/testdata/neg/neg-recognized-quoted-owned.toml
./internal/table/testdata/neg/neg-rule-empty-match-block.toml
./internal/table/testdata/neg/neg-rule-no-match-block.toml
./internal/table/testdata/neg/neg-ruleid-hash.toml
./internal/table/testdata/neg/neg-tagdecl-domain-on-bool.toml
./internal/table/testdata/neg/neg-tagdecl-single-valued-on-scalar.toml
./internal/table/testdata/neg/neg-tagdecl-single-valued-on-set.toml
./internal/table/testdata/neg/neg-terminal-unknown.toml
./internal/table/testdata/neg/neg-two-recognized-decls.toml
./internal/table/testdata/neg/neg-unknown-context.toml
./internal/table/testdata/neg/neg-unknown-field.toml
./internal/table/testdata/neg/neg-unknown-tag-match.toml
./internal/table/testdata/neg/neg-v2-shaped.toml
./internal/table/testdata/neg/neg-version.toml
./internal/table/testdata/neg/neg-write-observed.toml
./internal/table/testdata/neg/neg-write-value-outside-domain.toml
./internal/table/testdata/neg/neg-write-value-wrong-kind.toml
./internal/table/testdata/neg/probe-a-initial-observed-no-writer.toml
./internal/table/testdata/neg/probe-b-initial-observed-with-writer.toml
./internal/table/testdata/perm/rdr-eq-as-in.toml
./internal/table/testdata/perm/rdr-keyorder.toml
./internal/table/testdata/perm/rdr-ruleorder.toml
./internal/table/testdata/pos-near-miss-both.toml
./internal/table/testdata/pos-near-miss-folded.toml
./internal/table/testdata/pos-near-miss-space.toml
./internal/table/testdata/pos-no-initial.toml
./internal/table/testdata/pos-no-terminal.toml
./internal/table/testdata/pos-not-near-miss.toml
./internal/table/testdata/rdr-fixture.toml
./internal/table/testdata/write-delim-a.toml
./internal/table/testdata/write-delim-b.toml
./models/examples/pricing-decision-table.toml
./models/examples/release-grammar.toml
./models/examples/review-state-machine.toml
./models/rdr.toml
```

**No checked-in `testdata/*.toml` fixture authors `[rule.emit]`.** Emit fixtures
are constructed in Go source as inline TOML strings, not as on-disk fixtures.

### Split by location

| Location | Files authoring `[rule.emit]` |
| --- | --- |
| `models/` | `models/examples/pricing-decision-table.toml` (4 blocks) |
| `internal/**/testdata` (on-disk fixtures) | none |
| `docs/` (prior-RDR evidence artifact, not a lintable corpus model) | `docs/rdr/0010-stateless-decision-tables/evidence/spikes/a8-rule-emit-table.toml` (1 block) |

### Emit authored inline in Go source

```
$ rg -l 'rule\.emit' --glob '*.go' . | sort
```

```
./internal/cli/decision_table_0010_test.go
./internal/cli/failure_modes_0010_test.go
./internal/cli/flow_resolve.go
./internal/cli/mvv_0010_test.go
./internal/graphlint/adversarial_0010_test.go
./internal/graphlint/class_0010_test.go
./internal/table/class_0010_test.go
./internal/table/emit_0010_test.go
./internal/table/emit_shape_0010_test.go
./internal/table/fixtures_0010_test.go
./internal/table/model.go
./internal/table/normalize.go
./internal/table/source.go
```

```
$ rg -l 'rule\.emit' --glob '*.go' --glob '!*_test.go' . | sort
```

```
./internal/cli/flow_resolve.go
./internal/table/model.go
./internal/table/normalize.go
./internal/table/source.go
```

The four non-test files reference `[rule.emit]` only in doc comments and field
documentation (they do not author model text). The ten `_test.go` files author
`[rule.emit]` blocks as inline TOML string literals — these are the real
`[rule.emit]`-authoring corpus outside `models/`. Notable examples (verbatim
excerpts of the authoring sites, not full files):

```
./internal/table/emit_0010_test.go:264:		"integer": "[rule.emit]\nverdict = 42\n",
./internal/table/emit_0010_test.go:267:		"nested table": "[rule.emit]\nverdict = \"alpha\"\n" +
./internal/table/emit_0010_test.go:268:			"[rule.emit.sub]\ninner = \"x\"\n",
./internal/table/fixtures_0010_test.go:74:		sb.WriteString("[rule.emit]\n")
./internal/graphlint/class_0010_test.go:64:		"[rule.emit]\nverdict = \"" + id + "\"\n"
./internal/cli/mvv_0010_test.go:622:	if strings.Contains(string(src), "[rule.emit]") {
./internal/cli/mvv_0010_test.go:623:		t.Error("models/rdr.toml authors `[rule.emit]`; it is untouched by " +
```

Note `internal/table/emit_0010_test.go:268` is the one place in the whole repo
where the token `[rule.emit.sub]` appears — a deliberate negative case
(nested-table-under-emit is a decoder type error), not an `[emit.<key>]`
top-level declaration.

## 3. Baseline suite run

```
$ go build ./... && go test ./... 2>&1 | tail -40
```

```
?   	github.com/cwensel/intrastate/cmd/intrastate	[no test files]
ok  	github.com/cwensel/intrastate/docs/rdr/0010-stateless-decision-tables/evidence/spikes/a13-zerodim	0.323s
ok  	github.com/cwensel/intrastate/internal/accessor	1.034s
ok  	github.com/cwensel/intrastate/internal/cli	0.797s
ok  	github.com/cwensel/intrastate/internal/cli/clierr	0.248s
ok  	github.com/cwensel/intrastate/internal/cli/flowbind	0.702s
?   	github.com/cwensel/intrastate/internal/cli/respond	[no test files]
ok  	github.com/cwensel/intrastate/internal/graphlint	4.142s
ok  	github.com/cwensel/intrastate/internal/guard	1.496s
ok  	github.com/cwensel/intrastate/internal/resolve	1.595s
ok  	github.com/cwensel/intrastate/internal/table	1.925s
ok  	github.com/cwensel/intrastate/internal/version	1.117s
```

Build succeeded; every package `ok`. This is today's green baseline.

**Deferred leg.** A3 also asks for a full-suite run "with the loader change in
place and no model changed". That loader change does not exist at HEAD `d771278`
— RDR 0024 is at Draft and nothing implementing the `[emit.<key>]` declaration or
its zero-declaration opt-out has been written. That leg is therefore unrunnable
now and is deferred to implementation; it is not simulated or approximated here.

## 4. Baseline lint over the model corpus

```
$ ls models/ models/examples/
```

```
models/:
examples
rdr.toml

models/examples/:
pricing-decision-table.toml
release-grammar.toml
review-state-machine.toml
```

```
$ go run ./cmd/intrastate lint --help
```

```
Check a transition model against the mandatory graph invariants.

A model with any blocking finding is refused: the command returns the
aggregate error graph-lint-failed at exit 2 and carries every
blocking finding in the machine-readable findings list.

Bounds enforced by this build (model-independent implementation
constants, not per-model inputs):

  product bound   2048
  node ceiling    4096

Usage:
  intrastate lint [flags]

Flags:
      --flow string    flow id to lint (reserved; this build resolves none — use --model)
  -h, --help           help for lint
      --help-all       show extended help (vocabulary, wire shapes, exit codes)
      --model string   path to the transition model to lint

Global Flags:
      --as string   output mode: text | json (default "text")
```

```
$ go build -o <tmp>/intrastate-a3 ./cmd/intrastate
$ for m in models/rdr.toml models/examples/pricing-decision-table.toml \
           models/examples/release-grammar.toml models/examples/review-state-machine.toml; do
      echo "=== intrastate lint --model $m ==="
      <tmp>/intrastate-a3 lint --model "$m" 2>&1
      echo "exit=$?"
  done
```

```
=== intrastate lint --model models/rdr.toml ===
exit=0
=== intrastate lint --model models/examples/pricing-decision-table.toml ===
exit=0
=== intrastate lint --model models/examples/release-grammar.toml ===
  graph-coverage-closed-by-escape: the coverage of group release/build is closed by the bare escape row "begin-otherwise" rather than proved over its declared domains (rule="begin-otherwise" element="release/build")
  graph-coverage-closed-by-escape: the coverage of group release/hold is closed by the bare escape row "hold-otherwise" rather than proved over its declared domains (rule="hold-otherwise" element="release/hold")
  graph-coverage-closed-by-escape: the coverage of group release/ship is closed by the bare escape row "ship-otherwise" rather than proved over its declared domains (rule="ship-otherwise" element="release/ship")
exit=0
=== intrastate lint --model models/examples/review-state-machine.toml ===
  graph-coverage-closed-by-escape: the coverage of group review/approve is closed by the bare escape row "approve-otherwise" rather than proved over its declared domains (rule="approve-otherwise" element="review/approve")
  graph-coverage-closed-by-escape: the coverage of group review/reject is closed by the bare escape row "reject-otherwise" rather than proved over its declared domains (rule="reject-otherwise" element="review/reject")
  graph-coverage-closed-by-escape: the coverage of group review/submit is closed by the bare escape row "submit-otherwise" rather than proved over its declared domains (rule="submit-otherwise" element="review/submit")
```

All four models exit `0`. The `graph-coverage-closed-by-escape` lines on the two
grammar/state-machine models are advisory (non-blocking) findings, present today
and unrelated to emit. Today's model corpus lints clean.

## Findings

- `[emit.` / `[emit` top-level declaration hits in `.toml`: **0** repo-wide.
- `[rule.emit]` authoring, split by location:
  - `models/`: `models/examples/pricing-decision-table.toml` only (4 blocks). The
    A3 claim holds exactly for the `models/` corpus.
  - `internal/**/testdata`: **none** — no on-disk fixture authors `[rule.emit]`.
  - `docs/`: `docs/rdr/0010-stateless-decision-tables/evidence/spikes/a8-rule-emit-table.toml`
    (1 block) — a prior-RDR evidence artifact, not a lintable corpus model, but it
    is a checked-in file authoring `[rule.emit]`.
  - Inline TOML in Go test source: ten `_test.go` files across
    `internal/table`, `internal/cli`, and `internal/graphlint` author `[rule.emit]`
    blocks as string literals. These are the bulk of `[rule.emit]` authorship in
    the repo and they exercise the loader directly.
- Baseline suite: `go build ./...` clean, `go test ./...` all `ok`.
- Baseline lint: all four `models/*.toml` exit `0`.
- Post-change leg of A3 (full suite with the loader change in place): **not run**
  — the loader change does not exist at Draft.

### Verdict on A3-as-written

A3's *first* clause — "No checked-in model or fixture authors an `[emit.<key>]`
declaration" — is **accurate**: zero hits repo-wide.

A3's *second* clause — "only `models/examples/pricing-decision-table.toml` authors
`[rule.emit]`" — is **accurate only under a `models/`-scoped reading** and is
**false read corpus-wide**. Two other categories of checked-in file author
`[rule.emit]`: one prior-RDR evidence TOML under `docs/`, and ten Go test files
authoring inline TOML. The Go tests matter for the opt-out claim precisely because
they feed the loader; every one of those inline models declares zero
`[emit.<key>]`, so the zero-declaration opt-out must keep them passing unchanged —
which is exactly what the deferred post-change suite run has to prove.

Recommended narrowing: "only `models/examples/pricing-decision-table.toml` authors
`[rule.emit]` **under `models/`**; the remaining `[rule.emit]` authorship is inline
TOML in `internal/{table,cli,graphlint}` test sources, all of which likewise
declare zero `[emit.<key>]`."

The underlying *substance* of A3 — every checked-in emit author declares zero
`[emit.<key>]`, so the zero-declaration opt-out leaves the whole corpus linting as
today — is **supported by the evidence**. Only the enumeration in the wording is
too narrow.
