Model: claude-opus-5

# Stage 4 Resolve — verification pass for cli/0029

Accepted citations and rejected branches for the four Critical Assumptions.
Propose's passes are `prior-art.md` (external) and `in-repo-prior-art.md`
(at-HEAD vocabularies); this file is the verification pass and does not
restate them.

## A2 — no release exists yet (Source Search, VERIFIED)

- `git tag --list` empty; `git describe --tags` fails.
- `internal/version/version.go::resolve` defaults to
  `Info{Version: "dev", Commit: "none", Date: "unknown"}`; fills commit/date
  from VCS stamps, never moves `version` off `"dev"` without an `-ldflags`
  injection or module version.
- `.github/workflows/release.yml` — `on.push.tags: ["v*"]`, a tag that does
  not exist.

negative: no ldflags-independent path by which the reported version becomes
1.0.0 or higher.

## A3 — no strict envelope consumer (Source Search, VERIFIED)

- `internal/table/source.go::decodeStrict` — `DisallowUnknownFields`, TOML
  only, four call sites, all input, reached via `internal/table/load.go::Load`.
- `internal/cli/flow_input.go::planEnvelope` — the ONE real production
  consumer that re-parses the output envelope (`--plan` reading a
  `flow resolve --as json` record). Four-field struct, plain
  `json.Unmarshal`, so unknown top-level properties are ignored. Tolerant.
- Test-site envelope decodes are plain structs likewise.

negative: `schema_version` / `format_version` — no code hits, only this
record's prose. No `.claude/` skill, hook, Makefile or CI workflow parses the
envelope at all. No `DisallowUnknownFields` on any JSON path.

## A4 — three tiers partition the emitted surfaces (Source Search, VERIFIED as to the partition; C4 corrected)

Thirteen emitted vocabularies enumerated against C4's seven. The partition
holds — every unassigned set fits an existing tier — but C4's list was not
total. The five added, with defining sites:

| vocabulary | site | tier |
| --- | --- | --- |
| `level` (`note`/`warning`) | `internal/cli/respond/respond.go::Note`, `::Warn` | frozen |
| gate `verdict` | `internal/accessor/model.go::Verdicts`, emitted `internal/cli/flow_exec.go::gateResult` | frozen (also the accessor INPUT alphabet; `cmdbind` rejects a stranger) |
| `data.escape_class` | `internal/resolve/resolve.go::RefusalKinds` | frozen (enforced: `guard_atoms_test.go::TestReq79_NoSixthRefusalKindIsMinted`) |
| `flow next` unknown-`reason` | `internal/cli/flow_next.go::unknownFact.Reason` | append-only (DISTINCT from the `graph-unprovable-coverage` `reason` set) |
| `graph-lint-failed` | `internal/graphlint/taxonomy.go::AggregateCode` | frozen |

Correctly tier-less (model-authored open namespaces, `clierr` ascribes them
no meaning): `findings[].class`, `data.dispositions`.

Rejected branch: treating the graph-lint finding codes as one vocabulary.
Blocking (10) and advisory (4) take different tiers — a new BLOCKING code is
a refusal a previously-passing model now takes, which `append-only`
describes and `growing` understates.

negative: `level` is specified-but-unexercised — `Note`/`Warn` are reached by
no verb today, and no test asserts on the discriminator.

## A1 — the advisory closure (Peer RDR, PENDING — record, not code)

- Code permits a fifth freely: `internal/graphlint/taxonomy.go::severityFor`
  derives severity (`!IsBlocking` → `info`); no per-code severity field, no
  exhaustive switch, no cardinality constant; all eleven `AdvisoryCodes()`
  consumers iterate.
- One blocker, by value:
  `internal/graphlint/findings_0006_test.go::TestReq74_TheAdvisoryTierIsClosedAtExactlyFourMembers`
  — `slices.Equal` against a four-element literal at lines 288-293.
- `0006:C17` is Implemented and says "closed at" those four.

negative: no `len(AdvisoryCodes())` assertion, no severity lookup table keyed
by code, anywhere in `internal/graphlint/` or `internal/cli/`.

## Prior art accepted (C1's precedent)

- `opentofu/website/docs/internals/json-format.mdx:15-26` — the consumer
  rule verbatim, both halves. `internal/command/jsonplan/plan.go:31-58` —
  `format_version` beside `terraform_version`, not derived from it.
  Per-surface versions: `jsonstate/state.go:28` and
  `jsonprovider/provider.go:18` at `"1.0"` while jsonplan is `"1.2"`.
- `kubebuilder/pkg/plugin/external/types.go:26` — `apiVersion` on an
  inter-process JSON envelope.
- DevRef `go-modules-reference.pdf` p.4 — the 0.x instability rule,
  normative for a Go toolchain.
- DevRef `service_design_patterns.pdf` p.244 — Tolerant Reader / Postel's
  Law; p.245 for validated-on-send.
- DevRef `Designing Data Intensive Applications.pdf` p.121 — ignore the
  unrecognized field, forward compatibility.

Rejected / not found:
- SemVer 2.0.0 §4's own wording — not in any local corpus; not quoted from
  memory. The Go Modules Reference substitutes and is stronger here.
- The strict-input / tolerant-output INVERSION as received practice — the
  literature supports tolerant-on-receive and validated-on-send, not this
  specific audience split. Recorded in the RDR as this project's own
  decision.
- `helm/pkg/repo/v1/index.go:85` — `apiVersion`, but on a repository index
  FILE, not command output. Weaker; not cited in the record.
- Negative sweep, directly comparable CLIs with no envelope version field:
  gh, goreleaser, golangci-lint, hugo, consul, prometheus, etcd, semgrep.
