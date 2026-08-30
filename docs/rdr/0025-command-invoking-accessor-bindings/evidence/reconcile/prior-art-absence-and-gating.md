# Stage-6 prior art: absence representation, flag gating

## 1. Explicit-flag absence is the established pattern (not a workaround)

- **Go stdlib** `database/sql`: `sql.Null[T] struct { V T; Valid bool }`, doc
  "Valid is true if String is not NULL". A present struct with a flag, chosen
  over nil-pointer or column omission — for a domain where omission was available.
- **OpenTofu/cty** — closest structural analogue, models READ/ABSENT/UNREADABLE
  as three states of a PRESENT value via two predicates
  (`internal/configs/configschema/coerce_value.go:36-64`): `in.IsNull()` = absent,
  `!in.IsKnown()` = unknown. `internal/plans/dynamic_value.go:22-35` warns Go-level
  nil is a DIFFERENT thing from `cty.NullVal` — conflating container-level
  "not there" with domain-level null is a known bug source.
- **Consul KV** (`api/kv.go:69-89`): absent = `(nil pair, nil err)` with 404
  whitelisted as non-error; unreadable = `(nil pair, non-nil err)`. Two channels,
  three outcomes — same shape as `(values, unreadable, err)`.
- **gorm** (`migrator/column_type.go:91-117`): `Nullable() (nullable, ok bool)` —
  `ok=false` = driver cannot tell, distinct from `nullable=false`.
- **DevRef** SQL Antipatterns ch.14 "Fear of the Unknown": overloading an
  ordinary value for "missing" is the antipattern; the solution is a dedicated
  unambiguous marker. "A NULL is not a value; it is a marker for a value that is
  missing" — markers must be PRESENT to be read.

Counter-evidence, honestly reported: **etcd** `RangeResult` uses omission
(`kvstore_txn.go:90-109`). Not a counterexample — `Range` scans a key interval
with no caller-named requested set, and has NO third outcome (gRPC errors the
whole call), so the tri-state never arises.

## 2. Documented failure mode of omission-as-absence

- **proto3 field presence** (protobuf.dev/programming-guides/field_presence/) —
  the canonical case. Implicit presence collapses three facts: "explicitly set to
  default", "notionally cleared", "never set". Consequence: "an update patch alone
  cannot represent an update to the default value". Data round-tripping through an
  implicit-presence system is LOSSY. Google reversed this: `optional` returned in
  v3.12/v3.15, explicit presence is the default in editions. A decade-long
  industry reversal of exactly this choice.
- **RFC 7386 JSON Merge Patch** — null overloaded as "remove" makes a region of
  the value space inexpressible: "suitable for ... documents that ... do not make
  use of explicit null values". RFC 6902 avoids it with an explicit `{"op":"remove"}`
  field — the same move as an explicit `Absent` flag.
- **Live Go instance**: `zitadel/internal/crypto/key.go:43-70` — `key, ok :=
  readKeys[id]; if !ok { log; continue }` cannot distinguish a key missing from
  partial backend failure from one genuinely not configured.

## 3. The fail-closed default (omission => UNREADABLE) is well-precedented

`opentofu/internal/getproviders/errors.go:232-251` states the distinction and
picks the same default:
  "ErrIsNotExist returns true if and only if the given error is ... an
   AFFIRMATIVE RESPONSE that a requested object does not exist. This is as
   opposed to errors indicating that the source is unavailable ... where we
   therefore CANNOT SAY FOR CERTAIN whether the requested object exists."

Two principles, both satisfied by `classify`: absence requires an affirmative
response, never inferred from silence; the enumeration is closed and centralized,
everything else falls through to unknown. `io/fs` (fs.go:151-159) encodes the same
asymmetry — ENOENT is absent, EACCES is unreadable; `os.IsNotExist` documents
"whether its argument is KNOWN to report that a file ... does not exist".

=> `KeyValue{Absent: true}` is not a workaround forced by an arbitrary choice.
It is the REQUIRED consequence of a defensible, well-precedented default.

## 4. Mixed-granularity: whole-invocation exit code vs per-key envelope

Nearest prior art — GitHub Code Search as consumed by gh
(`pkg/cmd/skills/search/search.go:143-150`): `IncompleteResults bool
\`json:"incomplete_results"\`` sits INSIDE the envelope beside the items, not in
an HTTP status. When a result set's completeness is in question, the completeness
flag belongs at the same granularity as the results.

Hazard: an exit code is a small integer already claimed by the operational
vocabulary (crash, signal, cancellation, auth) — `gh`'s own exit-codes topic
(`help_topic.go:306-325`) warns a command "may have more exit codes". Putting a
SEMANTIC fact in that channel reintroduces the merge-patch overload one level up.
Second hazard: if the envelope can already express "all keys absent" (N records
with Absent:true), the exit code is a second source of truth that can DISAGREE.

Test to apply: a whole-invocation signal is sound only when it carries what the
per-key channel structurally CANNOT (e.g. "I could not produce a body at all").
"Every requested key is absent" fails that test — it is derivable from the envelope.

NOTE for 0025: `exit_absent` exists precisely for the case where there IS no
envelope (empty stdout). That is the "could not produce a body" case, so it PASSES
the test — but only because C3 restricts exit maps to empty stdout. That
restriction is load-bearing and should be read as such.

## 5. Flag gating (A13)

Eleven mature Go CLIs swept (helm, gh, goreleaser, kubebuilder, consul, etcd,
opentofu, golangci-lint, loki, prometheus, hugo). **Zero** rely on an
unregistered-flag lookup returning the zero value as a security decision.

pflag's own error text is the tell: `flag.go:376-378` returns
"flag accessed but not defined: %s" — the library classifies the miss as a
PROGRAMMER ERROR, not a query about user intent.

- **etcd** is the direct precedent — same problem shape (security boolean read by
  shared helpers). Registers on root `PersistentFlags` (`ctlv3/ctl.go:67,69`) AND
  checks the error at the shared read site, exiting on miss
  (`ctlv3/command/global.go:246-259`): registration discipline is what EARNS the
  right to check the error.
- **helm** funnels per-command registration through one shared helper
  (`pkg/cmd/helpers.go:29-41`); its read idiom `cmd.Flag(name).Value`
  (`helpers.go:45-47`) NIL-PANICS on an unregistered command.
- **hugo** — the exec-allowlist analogue: `security.exec.allow` is a config
  Whitelist, and the zero value denies BY EXPLICIT CONSTRUCTION
  (`whitelist.go:52-58` sets `acceptNone: true`), not by absence. Policy is a
  required constructor parameter (`hexec.New(cfg security.Config, ...)`), so no
  call site is reachable without a policy in hand — no lookup, no miss.

Hazard class: the safety property lives in a coincidence between two independent
decisions (no verb registers it; this call site discards the error), neither
documenting its dependence on the other. Three plausible edits break it, and the
dangerous one fails OPEN: a parent registering it persistently, or a verb
registering a different default, silently converts the "miss" into a "hit".

In-repo shape: `internal/cli/flow.go:84` — `newFlowCmd()`'s subtree is EXACTLY the
four verbs (next, resolve, read-state, set-state). `internal/cli/root.go:186`
already sets the precedent (`cmd.PersistentFlags()` for the shared output-mode
flag), read via `Lookup` + nil check (`root.go:243`), never a discarded GetBool.
`lint` sits at ROOT, outside the flow group — matching C6's "lint does not carry it".
