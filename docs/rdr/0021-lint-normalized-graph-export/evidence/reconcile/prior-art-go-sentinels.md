Model: claude-sonnet-5

# Prior art: reserved/sentinel string protection in Go codebases

Corpus searched: /Users/cwensel/sandbox/langref (~35 checked-out Go repos: gh-cli, helm,
goreleaser, kubebuilder, roborev, beads, prometheus, loki, consul, opentofu, gonum,
golangci-lint, eslint, sqlfluff, and others), via semble semantic code search
(MCP tool, first call indexed and cached the tree).

---

## 1. Reserved/sentinel strings REFUSED at ingress (not just documented)

### Prometheus: `__`-prefixed label namespace

- Repo: prometheus (and vendored copy in loki)
- File: `prometheus/vendor` root form at `github.com/prometheus/common/model/labels.go:43-83`
  (also `loki/vendor/github.com/prometheus/common/model/labels.go`)
- Mechanism: labels whose name starts and ends with exactly `__` are the
  reserved namespace (`isReservedLabel`, `prometheus/vendor/.../otlptranslator/strconv.go:73-91`
  parses this pattern generically). Within it, sub-prefixes are further
  reserved and documented as owned by specific pipeline stages:
  `MetaLabelPrefix = "__meta_"`, `TmpLabelPrefix = "__tmp_"` ("reserved for
  use ... by users" — i.e. the *inner* `__tmp_` band is the one users may
  write to; the rest of `__*__` is not).
- Enforcement site (ingress): `prometheus/scrape/target.go:626-694`,
  function `PopulateLabels`. This is the actual admission boundary for a
  scrape target's label set: it runs relabeling, then at line 681-687:
  ```go
  // Meta labels are deleted after relabelling. Other internal labels propagate to
  // the target which decides whether they will be part of their label set.
  lb.Range(func(l labels.Label) {
      if strings.HasPrefix(l.Name, model.MetaLabelPrefix) {
          lb.Del(l.Name)
      }
  })
  ```
  Reserved-but-not-meta `__` labels that survive relabeling are validated
  against the general label-name grammar (`LabelNameRE`,
  `prometheus/vendor/.../common/model/labels.go:97-111`) and duplicate/empty
  checks further down `PopulateLabels`; a user-authored series label
  colliding with a protocol-reserved name is caught here, at parse/ingress
  time, before the target is scraped — not at render/query time.
- Verdict: **ingress-time deletion + validation**, not merely documented.

### beads: reserved KV-key namespace, and reserved config keys

- Repo: beads
- File: `beads/cmd/bd/kv.go:20-48` (also mirrored at `internal`-relative
  path shown as `cmd/bd/kv.go:21-36` in a second index pass — same file)
- Mechanism: `validateKVKey` explicitly rejects:
  - empty/whitespace keys,
  - keys that would nest under the store's own `kv.` prefix,
  - keys under `kvkeys.MemoryPrefix` (reserved for `bd remember`/`bd forget`),
  - keys under `sync.`, `conflict.`, `federation.`, `jira.`, `linear.`,
    `export.`, `import.` (reserved internal-config prefixes).
  Each returns a `fmt.Errorf` — a hard refusal, not a warning.
- File: `beads/cmd/bd/config.go:1029-1042`, `rejectProtectedConfigKey`.
  Mechanism: a switch on protected key names (`issue_prefix`/`issue-prefix`)
  that returns a user-facing error plus the correct alternate command,
  rather than letting `bd config set` silently corrupt derived state.
- Enforcement site: both are called from the CLI command handler before
  the key is persisted — i.e. **ingress**, at the `bd kv set` / `bd config
  set` command boundary, not at read/render time.

---

## 2. Provenance: distinguishing computed/synthesized vs user-authored data in the same field

### Kubernetes: `ObjectMeta.ManagedFields` — parallel provenance map, not inline flags

- Repo: loki (vendored k8s.io/apimachinery), file:
  `loki/vendor/k8s.io/apimachinery/pkg/apis/meta/v1/types.go:285-296`
- Mechanism: rather than tagging each field value with a "who set this"
  flag, Kubernetes keeps a **separate, parallel structure**
  (`ManagedFields []ManagedFieldsEntry`) that maps `(workflow-id, version)`
  to the *set of field paths* that workflow last wrote. The object's actual
  data fields (spec/status/metadata) contain only the merged, current
  values — never a marker of origin. Provenance is reconstructed by
  looking up the field's JSON path in `ManagedFields`, not by inspecting
  the field's own value.
- This is a "wrapper/parallel structure" pattern: the same serialized
  object carries both the payload and an out-of-band ownership index,
  and the two are kept structurally separate so a user-authored value and
  a controller-synthesized value are never ambiguous at the storage layer
  — the ambiguity is resolved by consulting the side-table, not by reserving
  a value in-band.
- This is populated by the API server during Server-Side Apply (i.e. at the
  admission/ingress boundary for every write), not computed lazily at
  render time.

### Kubernetes: annotations/labels stay a flat `map[string]string` — no provenance split at all

- Same vendor tree: `ObjectMeta.Annotations`/`Labels` remain plain
  `map[string]string` (`loki/vendor/.../types.go`, same struct). Kubernetes
  chose NOT to split "system-set" vs "user-set" annotations into different
  typed fields or maps. Instead it relies on **namespace-prefix convention**
  (`kubernetes.io/`, `k8s.io/`) plus admission-controller-level policy
  (not generic apimachinery validation) to keep controllers from clobbering
  user annotations and vice versa. This is documented convention backed by
  narrow validation (DNS-subdomain syntax via `IsQualifiedName`/
  `IsDNS1123Subdomain`, `loki/vendor/.../util/validation/validation.go:163-182`)
  rather than a structural provenance mechanism — i.e. for
  labels/annotations specifically, Kubernetes is closer to the "counter-example"
  camp (see section 4) than to the ManagedFields camp. The split only
  exists for whole-field ownership (ManagedFields), not for arbitrary
  string-keyed maps like annotations.

---

## 3. Retrofit: protection added after an unprotected sentinel caused real corruption

### beads: `kv.memory.*` namespace reservation retrofitted after a data-loss bug (GH#2474)

- File: `beads/cmd/bd/kv.go:31-39` (comment inline with the code)
- History reconstructed from the code comment itself (no git-log tool was
  used; this is the committed rationale, not commit history):
  > "Reserve the persistent-memory namespace: a generic memory.* key would
  > store to kv.memory.*, indistinguishable from a `bd remember` memory, and
  > the merge resolver auto-resolves kv.memory.* conflicts with --theirs
  > (GH#2474). Without this guard a user's deliberate kv value could be
  > silently overridden by a remote on pull."
- What happened: before this guard, `kv.memory.*` was an ordinary,
  unprotected KV namespace. The merge/conflict resolver had separately
  started auto-resolving that namespace with a `--theirs` strategy (because
  it assumed only `bd remember` wrote there). Once a user could also write
  plain KV data into the same namespace, the two writers' values became
  indistinguishable in storage, and the auto-resolver would silently
  overwrite a user's manually-set value on the next sync — silent data
  loss, not a crash.
- Migration shape: **not a breaking change to the wire format or schema** —
  it is a pure CLI-side validation addition (`validateKVKey` gained one more
  `strings.HasPrefix` branch). Any *existing* user data already living
  under `kv.memory.*` before the fix remains in place and ambiguous; the
  fix only prevents *new* collisions going forward. This is refusal at the
  write/ingress boundary, added after the fact, treated as a bug fix rather
  than a breaking/versioned change (no deprecation cycle, no migration
  script for pre-existing colliding keys was found in this corpus).

### beads: deterministic dependency-edge IDs replacing DB-random UUIDs (GH#4259) — adjacent pattern

- File: `beads/internal/storage/depid/depid.go` (whole file; header comment
  lines 1-21, separator rationale lines 33-36)
- Not a "reserved value refused on ingress" case, but the same family of
  problem (collision-proofing a value that must be unambiguous) solved by
  **construction instead of validation**: the primary key was originally a
  DB-side random `UUID()` default, which minted a different value per Dolt
  clone for the *same logical edge*, causing merge failures/duplicate rows
  across clones. The fix derives the id deterministically via
  `uuid.NewSHA1(Namespace, issueID + sep + target)` where `sep` is ASCII
  Unit Separator (`\x1f`), chosen specifically because it "cannot occur in
  an issue id or a dependency target, so the encoding is unambiguous."
  Rather than validating/rejecting inputs that might collide, the design
  picks a delimiter from outside the legal character set of user data,
  making collision structurally impossible. Treated as an internal
  migration (a backfill function is mentioned, "every insert path and the
  upgrade backfill call New"), not a public breaking change — the primary
  key's *value* changes but its type/shape (CHAR(36) UUID) does not.

---

## 4. Counter-examples: collision allowed, ambiguity documented instead of prevented

- Kubernetes labels/annotations (see section 2, second entry) are the
  clearest counter-example found in this corpus: `Annotations`/`Labels`
  stay a flat, unstructured `map[string]string`. The reserved-prefix
  convention (`kubernetes.io/`, `k8s.io/`) is enforced only insofar as
  generic syntax validation applies to *all* keys (DNS-subdomain rules);
  apimachinery itself does not refuse a user who writes
  `kubernetes.io/foo: bar` on their own object — that policing, where it
  exists, lives in admission webhooks/controllers outside this vendored
  tree, not in the core type or its validation package. The core rationale
  implied by the code structure (no rejection path in apimachinery) is that
  ownership is a policy concern layered on top of a deliberately generic
  map, not a property the storage type itself can or should enforce.
- `google.golang.org/protobuf/reflect/protoregistry/registry.go:40-46`
  (vendored in loki and service) has an explicit, named, *configurable*
  ambiguity policy for a different but related problem (duplicate proto
  registrations): `var conflictPolicy = "panic" | "warn" | "ignore"`, with
  the comment "Neither of the above are covered by the compatibility
  promise and may be removed." This is a documented, intentional
  non-enforcement default (not per-value sentinel collision, but the same
  shape of tradeoff: correctness vs. permissiveness, resolved by policy
  rather than by refusal) — included here as the closest analogue to a
  "we document the ambiguity" stance found in this corpus. No repo in this
  set was found explicitly rationalizing an *unprevented* reserved-sentinel
  collision in user data; the closest matches either prevent it (Prometheus,
  beads) or leave it to out-of-band policy (Kubernetes annotations).

---

## Summary of mechanism vs. enforcement point

| Repo | Sentinel/problem | Mechanism | Enforcement point |
|---|---|---|---|
| prometheus | `__`-prefixed labels | prefix check + deletion of `__meta_*` | ingress (`PopulateLabels`, scrape-time) |
| beads | `kv.`, `memory.`, internal-config prefixes | `strings.HasPrefix` refusal, hard error | ingress (`bd kv set` CLI handler) |
| beads | protected config keys | switch/refusal with remediation message | ingress (`bd config set` CLI handler) |
| beads | dependency edge PK | out-of-band delimiter outside legal input charset | construction-time (id derivation), not validation |
| kubernetes (apimachinery) | field ownership (spec/status) | parallel `ManagedFields` side-table | ingress (server-side apply, API server write path) |
| kubernetes (apimachinery) | annotation/label namespace convention | flat map, syntax-only validation, no ownership refusal | none in core type; pushed to external policy |
| protobuf registry | duplicate registration | named, documented, configurable policy (panic/warn/ignore) | registration time, but deliberately non-refusing by default in some modes |

Dominant pattern across this corpus: **ingress-time refusal by prefix/namespace
check** (Prometheus, beads) is the majority pattern for "must never collide
with user data" sentinels, with an out-of-band/construction-time alternative
(beads dependency id, using a delimiter outside the legal input alphabet)
as a second, stronger variant that removes the need for validation entirely.
Where a project instead relies purely on documented convention with no
structural or ingress enforcement (Kubernetes annotations), the ambiguity is
real and explicitly left to layered, external policy (admission
controllers) rather than the core type system.
