package cli

// RDR 0005 — the shared fixture corpus for the `flow` command group.
//
// Every model here is a complete, self-contained document RDR 0002's loader
// must ACCEPT, so a failure under test is the `flow` verb's and not a load
// refusal. The one deliberately-invalid model (`flowInvalidModel`) is the
// exception and is named as such.
//
// The fixture family is shaped by the MVV (`0005:MVV`), which fixes three
// discriminating structures the corpus must carry:
//
//  1. a reader NO candidate row needs, with its artifact role left unbound —
//     the only assertion that separates the narrowed invoked read-accessor
//     set (REQ-35) from "every declared reader";
//  2. one gated candidate PLUS one guard-excluded gated row, so
//     `next --evaluate-gates` can be shown to run exactly one gate (REQ-42);
//  3. an escape row, so an escaped plan is provably a success and not a
//     laundered refusal (REQ-52).

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- canonical set literal probes ---------------------------------------
//
// JDR 0001 §D13's canonical form with the two members the MVV fixes: one
// carrying `<` and one carrying `&`. These are the ONLY spellings the
// contract admits on the wire — the `<` / `&` forms below are the
// defect REQ-70 forbids, kept here so a test can assert their ABSENCE by
// name rather than by hand-rolled escape.
const (
	// canonicalSetLiteral is sorted, duplicate-free, compact, and rendered
	// with HTML escaping DISABLED (REQ-70, REQ-73).
	canonicalSetLiteral = `["a<b","x&y"]`
	// escapedSetLiteral is what bare `json.Marshal` produces. It MUST NOT
	// appear at any emit site (REQ-70, REQ-72).
	escapedSetLiteral = "[\"a\\u003cb\",\"x\\u0026y\"]"

	// setMemberLT and setMemberAmp are the two members asserted
	// byte-identical through plan -> request -> read-back (`0005:MVV`).
	setMemberLT  = "a<b"
	setMemberAmp = "x&y"
)

// --- fixture identities --------------------------------------------------

const (
	// flowStateRole is the artifact role the MVV model's PRIMARY reader,
	// writer, and gate all bind. It is the role a caller supplies with
	// `--artifact state=<path>`.
	flowStateRole = "state"
	// flowOrphanRole is the artifact role of the reader NO candidate row
	// needs. It is deliberately left UNBOUND in `next` / `resolve`
	// invocations (`0005:MVV`).
	flowOrphanRole = "orphan"
)

// flowMVVModel is the MVV fixture: one flow proving all four verbs.
//
// Structure, keyed to the clauses it must discriminate:
//
//   - `[tags.status]`   owned scalar   — the scalar write (`0005:MVV`)
//   - `[tags.labels]`   owned set      — the set write carrying `<` and `&`
//   - `[tags.stale]`    owned scalar   — the `--clear` target
//   - `[tags.profile]`  observed       — the `--tag` channel (REQ-26)
//   - `[tags.note]`     owned scalar   — served ONLY by `read.orphan`, and
//     required by NO rule, so no candidate row's `RequiresOwned` names it
//     and `read.orphan` must never run under `next` / `resolve` (REQ-35,
//     REQ-36) while `read-state` must run it (REQ-37).
//
// `[read.orphan]` binds role `orphan`, which the `next` / `resolve`
// invocations leave unbound.
const flowMVVModel = `outcomes = ["advance", "hold", "bail"]
terminal = ["done"]

[model]
id = "mvvflow"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[tags.labels]
provenance = "owned"
kind = "set"
elements = ["a<b", "x&y", "plain"]

[tags.stale]
provenance = "owned"
kind = "scalar"

[tags.note]
provenance = "owned"
kind = "scalar"

[tags.profile]
provenance = "observed"
kind = "enum"
domain = ["mid", "foundational"]
single_valued = true

[read.state]
role = "state"
path = "flow.state"
keys = ["status", "labels", "stale"]
timeout = "2s"

[read.orphan]
role = "orphan"
path = "flow.note"
keys = ["note"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status", "labels", "stale"]
timeout = "2s"
read_back = true

[gate.approval]
role = "state"
path = "flow.gate"
keys = ["status"]
timeout = "2s"

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "advance-draft"
gate = ["approval"]
clear = ["stale"]
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"

[[rule]]
id = "hold-draft"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "hold"
[rule.write]
status = "draft"
`

// flowEscapeModel adds a modeled escape row for `no_match`, so `bail`
// resolves through the escape edge rather than refusing (REQ-52).
const flowEscapeModel = flowMVVModel + `
[[rule]]
id = "bail-escape"
escape = ["no_match"]
[rule.match.recognized]
eq = "bail"
`

// flowAmbiguousModel carries two rows matching the SAME outcome and state
// with no discriminating guard, so the kernel refuses `ambiguous_match`
// (REQ-49, REQ-94).
const flowAmbiguousModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "ambig"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[read.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"
read_back = true

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "first"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"

[[rule]]
id = "second"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"
`

// flowGatedNextModel carries ONE gated candidate reachable under
// `status = draft` and ONE gated row EXCLUDED by its guard on the same
// outcome. `next --evaluate-gates` must run the reported candidate's gate
// and NOT the excluded row's (REQ-42, `0005:MVV`).
const flowGatedNextModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "gatednext"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[tags.flag]
provenance = "owned"
kind = "bool"
single_valued = true
required = true

[read.state]
role = "state"
path = "flow.state"
keys = ["status", "flag"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status", "flag"]
timeout = "2s"
read_back = true

[gate.reported]
role = "state"
path = "flow.gate.reported"
keys = ["status"]
timeout = "2s"

[gate.excluded]
role = "state"
path = "flow.gate.excluded"
keys = ["status"]
timeout = "2s"

[initial]
status = "draft"
flag = "true"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "gated-reported"
gate = ["reported"]
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "final"

[[rule]]
id = "gated-excluded"
gate = ["excluded"]
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.guard.all.flag]
eq = "false"
[rule.write]
status = "final"
`

// flowInvalidModel is DELIBERATELY unloadable: `[tags.status]` declares an
// undeclared kind, so RDR 0002's loader refuses it with a stable category
// the CLI maps to `flow-model-invalid` (REQ-24, REQ-91).
const flowInvalidModel = `outcomes = ["advance"]

[model]
id = "bad"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "not-a-declared-kind"
single_valued = true
required = true
`

// --- filesystem helpers --------------------------------------------------

// writeFlowModel writes src to a fresh temp file and returns its path. It
// is the `--model <path>` argument; the CLI performs the file I/O and hands
// RDR 0002's loader BYTES plus a source id (REQ-23).
func writeFlowModel(t *testing.T, src string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "flow.toml")
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatalf("write fixture model: %v", err)
	}
	return path
}

// newFlowArtifact returns a path in a fresh temp dir for the caller-owned
// artifact a `--artifact role=path` binding names. The artifact's on-disk
// FORMAT is the implementation's choice and is never asserted here: every
// round-trip oracle in this suite goes `set-state` -> `read-state` through
// the CLI, so it holds for any format the implementation picks (REQ-107).
func newFlowArtifact(t *testing.T, name string) string {
	t.Helper()

	return filepath.Join(t.TempDir(), name)
}

// artifactBinding renders the `role=path` pair `--artifact` takes.
func artifactBinding(role, path string) string { return role + "=" + path }

// flowGateDenyModel carries three gated rows on distinct outcomes so the
// gate precedence rule can be exercised separately (REQ-51, REQ-97,
// REQ-98):
//
//   - `deny-path`          one gate that DENIES
//   - `indeterminate-path` one gate that is INDETERMINATE, none denying
//   - `both-path`          one denying AND one indeterminate gate, so deny
//     must override indeterminate
//
// Which verdict a fixture gate returns is the implementation's binding
// concern; the gate ACCESSOR IDS below are what the contract fixes, and the
// suite asserts on the resulting refusal identity.
const flowGateDenyModel = `outcomes = ["deny-path", "indeterminate-path", "both-path"]
terminal = ["done"]

[model]
id = "gatedeny"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[read.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"
read_back = true

[gate.deny]
role = "state"
path = "flow.gate.deny"
keys = ["status"]
timeout = "2s"

[gate.indeterminate]
role = "state"
path = "flow.gate.indeterminate"
keys = ["status"]
timeout = "2s"

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "deny-row"
gate = ["deny"]
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "deny-path"
[rule.write]
status = "final"

[[rule]]
id = "indeterminate-row"
gate = ["indeterminate"]
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "indeterminate-path"
[rule.write]
status = "final"

[[rule]]
id = "both-row"
gate = ["deny", "indeterminate"]
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "both-path"
[rule.write]
status = "final"
`

// flowGateFailModel carries one row whose gate cannot be CONSULTED — its
// artifact path names a locator the binding cannot reach. A gate that could
// not be consulted is exit 3, never a deny (REQ-47, REQ-55).
const flowGateFailModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "gatefail"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[read.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"
read_back = true

[gate.unreachable]
role = "state"
path = "flow.gate.unreachable"
keys = ["status"]
timeout = "2s"

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "gated-advance"
gate = ["unreachable"]
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"
`

// flowGuardUnevaluableModel guards a row on an OPTIONAL owned key. An
// unmarked key is optional (RDR 0003), so a value atom over it refuses
// `guard_unevaluable` when nothing established it (REQ-96).
const flowGuardUnevaluableModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "guardunev"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[tags.opt]
provenance = "owned"
kind = "enum"
domain = ["p", "q"]
single_valued = true

[read.state]
role = "state"
path = "flow.state"
keys = ["status", "opt"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status", "opt"]
timeout = "2s"
read_back = true

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "guarded-advance"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.guard.all.opt]
eq = "p"
[rule.write]
status = "final"
`

// flowReadBackFailModel declares a writer whose post-mutation read-back
// cannot COMPLETE — its read-back locator names something the binding
// cannot re-read. A read-back that did not run is exit 3 and carries a
// "may have been applied" detail (REQ-104, REQ-105); it is NOT a mismatch,
// which would assert the artifact is positively wrong.
const flowReadBackFailModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "readbackfail"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[read.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state.readback-unreachable"
keys = ["status"]
timeout = "2s"
read_back = true

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "advance"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"
`

// flowSetObservedModel declares a SET-kind OBSERVED tag, so a set literal
// can be shown to cross the CLI in `--tag` under the same canonical form it
// takes in `--write` and in any payload (REQ-29, REQ-70).
const flowSetObservedModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "setobserved"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[tags.marks]
provenance = "observed"
kind = "set"
elements = ["a<b", "x&y", "plain"]

[read.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"
read_back = true

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "advance"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"
`

// flowSetPlanModel's one rule PLANS a set-valued write whose members carry
// `<` and `&`. It is the fixture behind the MVV's three-hop byte-identity
// assertion: the plan emits the set, the request carries it back verbatim,
// and the read-back returns it (`0005:MVV`, REQ-70, REQ-108).
const flowSetPlanModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "setplan"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[tags.labels]
provenance = "owned"
kind = "set"
elements = ["a<b", "x&y", "plain"]

[read.state]
role = "state"
path = "flow.state"
keys = ["status", "labels"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status", "labels"]
timeout = "2s"
read_back = true

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "advance"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"
labels = ["a<b", "x&y"]
`

// flowEscapeOtherOutcomeModel carries an escape row bound to a DIFFERENT
// outcome than the one under test, whose guard reads an owned key served
// only by a reader the caller never binds. The kernel skips a differing-
// outcome escape row before it can rescue anything
// (`internal/resolve.escapeOrRefuse`), so that row's reader is not part of
// the requested outcome's demand set and its unbound role must not refuse
// the request (REQ-35, REQ-36; DEV-8).
const flowEscapeOtherOutcomeModel = flowMVVModel + `
[[rule]]
id = "bail-escape-guarded"
escape = ["no_match"]
[rule.match.recognized]
eq = "bail"
[rule.guard.all.sidenote]
eq = "ready"

[tags.sidenote]
provenance = "owned"
kind = "scalar"

[read.sidecar]
role = "sidecar"
path = "flow.sidecar"
keys = ["sidenote"]
timeout = "2s"
`

// flowDomainModel declares one tag of each CONSTRAINED kind, all served by
// the same writer, so a `--write` value can be held to its declaration
// independently of any accessor or rule:
//
//   - `[tags.status]`  enum with a `domain`
//   - `[tags.iter]`    int with `min` / `max`
//   - `[tags.ready]`   bool
//   - `[tags.labels]`  set with `elements`
//   - `[tags.free]`    scalar — the UNCONSTRAINED control, which must keep
//     accepting anything, so a domain refusal elsewhere is provably the
//     DECLARATION biting and not a blanket tightening of `--write`
//   - `[tags.hint]`    observed enum with a `domain`, the `--tag` channel
//
// `[write.state]` names every owned key, so none of them can be refused
// `flow-write-unbound` — a refusal here is the value's, not the binding's.
const flowDomainModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "domain"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[tags.iter]
provenance = "owned"
kind = "int"
min = 0
max = 9

[tags.ready]
provenance = "owned"
kind = "bool"

[tags.labels]
provenance = "owned"
kind = "set"
elements = ["alpha", "beta"]

[tags.free]
provenance = "owned"
kind = "scalar"

[tags.hint]
provenance = "observed"
kind = "enum"
domain = ["low", "high"]
single_valued = true

[read.state]
role = "state"
path = "flow.state"
keys = ["status", "iter", "ready", "labels", "free"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status", "iter", "ready", "labels", "free"]
timeout = "2s"
read_back = true

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "advance"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"
`

// flowExpiredTimeoutModel is `flowMVVModel` with every declared accessor
// timeout narrowed from `2s` to `1ns`.
//
// `1ns` is a LEGAL declared value, not a test backdoor: the loader accepts
// any positive Go duration (`internal/table/load.go`), so this model is one
// a caller could write. It is the observable surface through which the two
// TIMEOUT classes are reachable — `flow-accessor-timeout` on a read and
// `flow-write-readback-timeout` on the post-mutation read-back — and
// without it those two of the five exit-3 codes have no CLI-reachable
// provocation at all.
//
// It is deterministic rather than racy: `context.WithTimeout(ctx, 1ns)` is
// already expired before the binding returns, and the executor checks the
// deadline BEFORE it inspects the binding's result
// (`internal/accessor/executor.go`, `invokeRead`), so `timeout` is reported
// on every run rather than winning a race some of the time.
var flowExpiredTimeoutModel = strings.ReplaceAll(flowMVVModel, `timeout = "2s"`, `timeout = "1ns"`)
