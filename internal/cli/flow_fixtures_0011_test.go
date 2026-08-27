package cli

// RDR 0011 — the discriminating fixture corpus for the match-conditioned
// `flow next` candidate predicate.
//
// C3 states outright why this file has to exist: "Because no shipped 0005
// fixture contains a row whose match key is present and unequal, the 0005
// suite cannot by itself distinguish this predicate from the stripped-match
// one it replaces." Every model below is authored to separate the two
// predicates, and each names the clause it discriminates.
//
// Two authoring constraints bind every fixture here and are recorded once:
//
//  1. An ordinary rule needs a write block (`malformed_rule_shape`), so a
//     row that must match on a key WITHOUT that key entering
//     `RequiresOwned` has to write some OTHER key (REQ-62). Where the
//     fixture ignores that, C1's owned-key walk emits the same
//     `{key, absent}` pair the `--all` filter leaves alone and the fixture
//     discriminates nothing.
//  2. An AUTHORED empty `[rule.match]` is refused at load (REQ-105), so
//     C1's empty-match-pattern clause is reached only through a row all of
//     whose match atoms are over ABSENT keys.

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
)

// --- fixture role identities ---------------------------------------------

const (
	// flowMatchRole is the artifact role the 0011 fixtures' primary reader
	// binds.
	flowMatchRole = "state"
	// flowSideRole is the role of the reader serving a MATCH-ONLY owned key
	// (S8). It is bound, unbound, or pointed at a refusing locator by the
	// demand-set oracles.
	flowSideRole = "side"
)

// flowMatchClassesModel authors the three reachable match input classes
// (`disposition`, S3, MVV 6) over ONE observed key, plus the `0002:C13`
// dead row, plus a control row that matches on nothing observable.
//
// `phase` is OBSERVED, so the caller supplies it with `--tag` and the three
// classes are reachable from one model by varying the flag:
//
//   - no `--tag`            ⇒ `phase` ABSENT from the view for every row
//   - `--tag phase=alpha`   ⇒ present-and-EQUAL for `match-alpha`,
//     present-and-UNEQUAL for `match-beta`
//   - `--tag phase=zeta`    ⇒ present-and-unequal for BOTH
//
// Every match row writes `status`, an owned key the reader establishes, so
// the row's `RequiresOwned` never names `phase` — REQ-62's requirement, and
// what makes the `--all` filter oracle (REQ-39/REQ-42) discriminate here.
//
// `dead-row` authors `eq` + `in` on one key (`0002:C13`). A12 fixes the
// expansion's shape precisely, and the fixture is authored to it: the `in`
// member that COINCIDES with the `eq` literal (`alpha`) renders
// byte-identically and collapses to ONE tag, so that expanded row is LIVE;
// the member that DIFFERS (`zeta`) leaves a row carrying TWO tags on one
// key with differing literals — `0002:C13`'s dead row, which no value
// satisfies. Both expanded rows share the rule id `dead-row`.
//
// REQ-100 turns on that asymmetry: the two-match-tags assertion is written
// over the DEAD row only, because "the pairing's other expanded row
// collapses to one tag and is live, so asserting two there would assert
// something false."
//
// C1's disposition for the dead row: a candidate while `phase` is absent —
// carrying `{phase, absent}` exactly as any undecided row does — and
// excluded by the KERNEL's conjunction once `phase` is present at ANY
// value, including a value on neither side of the pairing. Never by a CLI
// literal comparison, which REQ-10 forbids outright.
const flowMatchClassesModel = `outcomes = ["alpha", "beta", "gamma", "dead"]
terminal = ["done"]

[model]
id = "matchclasses"
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

[tags.phase]
provenance = "observed"
kind = "enum"
domain = ["alpha", "beta", "zeta"]
single_valued = true

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
id = "match-alpha"
gate = ["approval"]
[rule.match.phase]
eq = "alpha"
[rule.match.recognized]
eq = "alpha"
[rule.write]
status = "final"

[[rule]]
id = "match-beta"
gate = ["approval"]
[rule.match.phase]
eq = "beta"
[rule.match.recognized]
eq = "beta"
[rule.write]
status = "final"

[[rule]]
id = "no-match-atoms"
[rule.match.recognized]
eq = "gamma"
[rule.write]
status = "final"

[[rule]]
id = "dead-row"
[rule.match.phase]
eq = "alpha"
in = ["alpha", "zeta"]
[rule.match.recognized]
eq = "dead"
[rule.write]
status = "final"
`

// flowSetKindedMatchModel is S3(a) — the ONE fixture that separates a
// filter-applied-over-`KernelRow()` build from a filter-applied-over-
// `row.Atoms` build (REQ-6, REQ-102).
//
// `marks` is a SET-kinded observed tag. `Row.KernelRow` renders each match
// atom's literal through `internal/table/model.go::seamValue`, which
// canonicalizes a set literal — sorted, compacted JSON, HTML escaping
// disabled — against the row's own set-key list. A build that rebuilds the
// match tags from `row.Atoms` in `internal/cli` reimplements that encoding
// in a second place and mis-compares the literal against the kernel's
// canonical form. Over every other fixture in this file the two builds are
// observationally identical.
//
// The authored literal carries `&`, one of the two members RDR 0005 pins
// for escaping. `eq` is a single-value operator, so a set-kinded match key
// takes a ONE-member literal — but set-ness here is the DECLARED KIND, not
// the member count, so `seamValue` still renders it as the §D13 canonical
// JSON array `["x&y"]` with `SetEscapeHTML(false)`. A build that rebuilds
// the tag from `row.Atoms` hands the kernel the bare member `x&y`, or the
// HTML-escaped `["x\u0026y"]`, either of which mis-compares against the
// caller's canonically-encoded `--tag marks=["x&y"]` and reports a spurious
// `no_match`. `scalar-control` is the negative control the oracle pairs it
// with: over a scalar match key the two builds agree.
const flowSetKindedMatchModel = `outcomes = ["advance", "hold"]
terminal = ["done"]

[model]
id = "setmatch"
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
id = "set-match"
[rule.match.marks]
eq = ["x&y"]
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"

[[rule]]
id = "scalar-control"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "hold"
[rule.write]
status = "final"
`

// flowMixedStateRowModel is S3(b) — one row carrying THREE match keys in
// mixed states at once (REQ-103): one present-and-equal, one
// present-and-unequal, one absent.
//
// C1's exclusion is the kernel's conjunction, per key: the unequal key
// decides and the row is excluded by default, even though another key on
// the same row holds and a third is absent. Under `--all` the same row is a
// candidate whose `unknown` carries NO match entry (C2), whatever any of
// the three keys' presence.
//
// All three keys are OBSERVED and none is written, so none enters
// `RequiresOwned` (REQ-62) and the `--all` assertion is not confounded by
// the owned-key walk.
const flowMixedStateRowModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "mixedstate"
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

[tags.holds]
provenance = "observed"
kind = "enum"
domain = ["yes", "no"]
single_valued = true

[tags.fails]
provenance = "observed"
kind = "enum"
domain = ["yes", "no"]
single_valued = true

[tags.missing]
provenance = "observed"
kind = "enum"
domain = ["yes", "no"]
single_valued = true

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
id = "mixed-row"
[rule.match.holds]
eq = "yes"
[rule.match.fails]
eq = "yes"
[rule.match.missing]
eq = "yes"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"
`

// flowAbsentMatchKeyModel is C3's absent-key discriminator and S4's
// all-absent row (REQ-61, REQ-62, REQ-104).
//
//   - `absent-row` matches on `wanted`, an observed key nothing supplies,
//     and WRITES `status` — a different key — so `wanted` never enters
//     `RequiresOwned`. By default it is a candidate carrying
//     `{wanted, absent}`; under `--all` that entry MUST disappear.
//   - `recognized-only` is A5's fixture B: its ONLY match atom is
//     `recognized`, lifted into `Row.Outcome` at normalize, so it is a
//     candidate with NO match fact in `unknown` (REQ-82, REQ-104).
//   - `all-absent` matches on TWO observed keys nothing supplies, the
//     reachable form of C1's empty-probe clause: its probe carries an empty
//     match pattern, which matches unconditionally, and every omitted key
//     is listed `absent`, so its `unknown` is NOT empty (REQ-14, REQ-104).
const flowAbsentMatchKeyModel = `outcomes = ["advance", "hold", "bail"]
terminal = ["done"]

[model]
id = "absentmatch"
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

[tags.wanted]
provenance = "observed"
kind = "enum"
domain = ["yes", "no"]
single_valued = true

[tags.other]
provenance = "observed"
kind = "enum"
domain = ["yes", "no"]
single_valued = true

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
id = "absent-row"
[rule.match.wanted]
eq = "yes"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"

[[rule]]
id = "recognized-only"
[rule.match.recognized]
eq = "hold"
[rule.write]
status = "final"

[[rule]]
id = "all-absent"
[rule.match.wanted]
eq = "yes"
[rule.match.other]
eq = "yes"
[rule.match.recognized]
eq = "bail"
[rule.write]
status = "final"
`

// flowFullyResolvedModel is MVV 7's fixture (REQ-91): a candidate with
// NOTHING undecided, so `unknown` can be asserted present as `[]` rather
// than omitted in BOTH modes.
//
// MVV 7 fixes the shape and says why no shipped 0005 fixture serves: the
// row's match and guard keys are all established by the invoked reader, the
// reader also serves the owned keys it writes, and the run carries
// `--evaluate-gates` — without which any declared gate id lands in
// `unknown` as `not-evaluated` and the list is non-empty by construction.
// This model therefore declares NO gate at all, so the oracle holds with or
// without the flag.
const flowFullyResolvedModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "fullyresolved"
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

[initial]
status = "draft"
flag = "true"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "resolved-row"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "final"
`

// flowLooseWriterModel and flowUncomparableGuardModel are S7's
// `uncomparable` pair (REQ-66, REQ-84, REQ-112).
//
// A13 records exactly ONE reachable producer of `ReasonUncomparable` from
// `flow next`: the seam answering unevaluable for a PRESENT value the
// operator cannot parse. The other three producers are unreachable from the
// CLI and MUST NOT be fixtured, "since a fixture that forces one tests
// `internal/resolve`, which this RDR does not change" (C3).
//
// Reaching that producer through the CLI takes TWO models over ONE
// artifact, and the reason is recorded here rather than discovered: every
// value channel into the view is validated against the DECLARATION it
// crosses — `--tag` refuses `flow-tag-invalid` and `--write` refuses
// `flow-write-invalid` for a value that does not conform to its kind — so
// no single model can both hold a non-integer and guard it with `gt`
// (`gt` is refused at load on any kind but `int`).
//
// The artifact, however, is CALLER-OWNED and is not a per-model store. So
// the value is ESTABLISHED through the production write path under
// `flowLooseWriterModel`, where `size` is a `scalar` and `notanint`
// conforms, and READ under `flowUncomparableGuardModel`, where the same key
// is an `int` the row guards with `gt`. The value is present in the view
// and the operator cannot parse it: `internal/guard/grammar.go`'s
// `strconv.Atoi(value)` fails and the seam answers `GuardUnevaluable`,
// which the kernel reports as `{size, uncomparable}` on the probe's
// `guard_unevaluable` refusal payload. That fact is one the CLI's own view
// walk CANNOT see, because the key is PRESENT — which is precisely why the
// payload read C1 adds is the only way to surface it.
//
// No test here asserts on the artifact's on-disk FORMAT: both hops go
// through the CLI, exactly as the 0005 harness requires.
const flowLooseWriterModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "loosewriter"
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

[tags.size]
provenance = "owned"
kind = "scalar"

[tags.gone]
provenance = "owned"
kind = "scalar"

[read.state]
role = "state"
path = "flow.state"
keys = ["status", "size", "gone"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status", "size", "gone"]
timeout = "2s"
read_back = true

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "seed-row"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"
`

// flowUncomparableGuardModel reads the artifact the loose writer seeded.
//
//   - `uncomparable-row` guards `gt = "3"` over the PRESENT non-integer, so
//     the seam answers unevaluable and the kernel reports
//     `{size, uncomparable}` (REQ-84).
//   - `absent-guard-row` is the paired negative control: the same shape of
//     guard atom over an ABSENT key, which the walk and the payload BOTH
//     produce as `{gone, absent}` and which MUST dedup to exactly ONE entry
//     (REQ-24, REQ-112).
const flowUncomparableGuardModel = `outcomes = ["advance", "hold"]
terminal = ["done"]

[model]
id = "uncomparable"
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

[tags.size]
provenance = "owned"
kind = "int"

[tags.gone]
provenance = "owned"
kind = "scalar"

[read.state]
role = "state"
path = "flow.state"
keys = ["status", "size", "gone"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status", "size", "gone"]
timeout = "2s"
read_back = true

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "uncomparable-row"
[rule.match.recognized]
eq = "advance"
[rule.guard.all.size]
gt = "3"
[rule.write]
status = "final"

[[rule]]
id = "absent-guard-row"
[rule.match.recognized]
eq = "hold"
[rule.guard.all.gone]
eq = "ready"
[rule.write]
status = "final"
`

// flowOwnedUnavailableModel is S7's precedence-limit fixture (REQ-31,
// REQ-85, REQ-112): a row that BOTH lacks an established owned key and
// carries an `uncomparable` guard atom.
//
// C1 states the limit rather than hiding it: when the probe refuses
// `owned_state_unavailable`, `gate` returns BEFORE it collects the guard
// payload, so the `uncomparable` guard atom is not on the wire and is not
// reported — while every `absent` fact on the row still is, from the walk.
//
// It reads the SAME artifact the loose writer seeded, for the same reason
// the model above does. `missingowned` is an owned key this model's row
// WRITES — so it enters `RequiresOwned` — and which the loose writer never
// established, so `missingOwned` fires; `size` is the present non-integer
// the `gt` operator cannot parse.
const flowOwnedUnavailableModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "ownedunavail"
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

[tags.missingowned]
provenance = "owned"
kind = "scalar"

[tags.size]
provenance = "owned"
kind = "int"

[read.state]
role = "state"
path = "flow.state"
keys = ["status", "missingowned", "size"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status", "missingowned", "size"]
timeout = "2s"
read_back = true

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "unavail-row"
[rule.match.recognized]
eq = "advance"
[rule.guard.all.size]
gt = "3"
[rule.write]
missingowned = "x"
status = "final"
`

// flowSortOrderModel is S7's sort fixture (REQ-29, REQ-113): the facts are
// MINTED out of `(key, reason)` order, so an unsorted build fails.
//
// The row's own emission order is: the atom walk contributes `zebra`
// (absent match key) and `alpha` (absent match key) in authored order,
// then the owned-key walk contributes `mid` (an owned key no reader
// established), then the gate-id loop contributes the gate id `beta` with
// reason `not-evaluated`. Sorted by `(key, reason)` the list is
// `alpha`, `beta`, `mid`, `zebra` — an order NO single source produces.
const flowSortOrderModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "sortorder"
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

[tags.mid]
provenance = "owned"
kind = "scalar"

[tags.zebra]
provenance = "observed"
kind = "scalar"

[tags.alpha]
provenance = "observed"
kind = "scalar"

[read.state]
role = "state"
path = "flow.state"
keys = ["status", "mid"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status", "mid"]
timeout = "2s"
read_back = true

[gate.beta]
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
id = "sort-row"
gate = ["beta"]
[rule.match.zebra]
eq = "z"
[rule.match.alpha]
eq = "a"
[rule.match.recognized]
eq = "advance"
[rule.write]
mid = "m"
status = "final"
`

// flowMatchOnlyOwnedModel is S8's demand-set fixture (REQ-16..21, REQ-114,
// REQ-116, REQ-117): a model declaring an owned key that some row MATCHES
// on but NO row writes, clears, or guards, served by exactly one reader.
//
// `mode` is that key. On the UN-EXTENDED demand set — `RequiresOwned` ∪
// guard-owned — `read.side` is never invoked, `mode` is absent from the
// view, and EVERY row comes back a candidate carrying `{mode, absent}`:
// the default silently degraded to `--all`. With C1's match-owned term the
// reader IS invoked, `mode` IS in the view, and `mode-row` is
// match-DECIDED.
//
// `plain-row` is S8's BREAKING arm (REQ-116, REQ-20): a SECOND ordinary row
// for the SAME outcome (`go`) that does NOT match on `mode`. Because the
// demand set is a union over the outcome's rows, `flow resolve --outcome go`
// invokes `read.side` even though the row it would select does not need it
// — so over an UNBOUND or refusing reader a run that returned that row's
// plan now exits 2/3. That is the class C1 names and accepts.
//
// `read.side` serves TWO keys (`mode` and `extra`) so REQ-117's oracle —
// both verbs' `readers` equal over a reader serving several keys — has its
// subject.
const flowMatchOnlyOwnedModel = `outcomes = ["go", "stop"]
terminal = ["done"]

[model]
id = "matchonlyowned"
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

[tags.mode]
provenance = "owned"
kind = "enum"
domain = ["fast", "slow"]
single_valued = true

[tags.extra]
provenance = "owned"
kind = "scalar"

[read.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"

[read.side]
role = "side"
path = "flow.side"
keys = ["mode", "extra"]
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
id = "mode-row"
[rule.match.mode]
eq = "fast"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "final"

[[rule]]
id = "plain-row"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "final"

[[rule]]
id = "stop-row"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "final"
`

// flowZeroOwnedMatchModel is S8's negative control (REQ-21, REQ-115): a
// model where NO row match-owns a key, so C1's added demand term finds
// nothing and the invoked reader set is unchanged.
//
// A genuine zero-owned-TAG model is unloadable on this build (an ordinary
// rule with no write block is refused `malformed_rule_shape`, and a write
// to an observed tag is refused `write_to_non_owned_tag`), which the A5
// spike records as a limit — so this fixture carries one owned tag that
// every row WRITES and none match-owns, exercising the property the term
// needs.
const flowZeroOwnedMatchModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "zeroownedmatch"
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

[tags.hint]
provenance = "observed"
kind = "enum"
domain = ["x", "y"]
single_valued = true

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
id = "hint-row"
[rule.match.hint]
eq = "x"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"
`

// flowTwoReadersOneKeyModel is the FIRST of C3's three pins (REQ-63,
// REQ-106): a model declaring ONE owned key served by TWO readers, which
// `internal/table/load.go::checkAccessorBindings` refuses at LOAD — before
// any probe is built.
//
// That refusal is what makes C1's presence test EXACT rather than
// approximate: a present-but-conflicted key, the one case the kernel's
// `matches` folds into `false` without comparing, cannot reach a CLI-built
// probe (A3).
const flowTwoReadersOneKeyModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "tworeaders"
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

[read.second]
role = "state"
path = "flow.second"
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

// --- 0011 payload readers ------------------------------------------------
//
// C3 forbids repointing `flow_harness_0005_test.go::stringsAt` at the new
// element type: "it is shared harness used well beyond these five sites,
// and changing its return shape would touch readers this contract has not
// counted. The re-homing adds a SIBLING helper beside it (a
// `[]{key, reason}` reader over the same payload)."
//
// unknownPair is that element as A-2 fixes the wire members: exactly `key`
// and `reason`, both always emitted, on a list named `unknown`.

type unknownPair struct {
	Key    string `json:"key"`
	Reason string `json:"reason"`
}

// unknownAt is the SIBLING helper C3 mandates. It reports whether
// `unknown` was present AND array-of-`{key, reason}`-shaped — the ok bool
// the five re-homed reads MUST assert, because three of them discard
// `stringsAt`'s comma-ok today and one is a NEGATIVE assertion that passes
// vacuously on the empty slice (REQ-58).
func unknownAt(m map[string]any, key string) ([]unknownPair, bool) {
	raw, ok := m[key]
	if !ok {
		return nil, false
	}
	arr, ok := raw.([]any)
	if !ok {
		return nil, false
	}
	out := make([]unknownPair, 0, len(arr))
	for _, e := range arr {
		obj, ok := e.(map[string]any)
		if !ok {
			return nil, false
		}
		k, kok := obj["key"].(string)
		r, rok := obj["reason"].(string)
		if !kok || !rok {
			return nil, false
		}
		out = append(out, unknownPair{Key: k, Reason: r})
	}
	return out, true
}

// candidateUnknown reads one candidate's `unknown` list, FAILING the test
// when it is missing or mis-shaped. C2 requires the field present in both
// modes, as `[]` rather than omitted, so an absent field is never a
// tolerable shape here.
func candidateUnknown(t *testing.T, c map[string]any) []unknownPair {
	t.Helper()

	pairs, ok := unknownAt(c, "unknown")
	if !ok {
		t.Fatalf("candidate %v carries no `unknown` list of {key, reason} "+
			"pairs; C1 renames `unresolved` to `unknown` and types its "+
			"entries, and C2 requires the field present in BOTH modes as `[]` "+
			"rather than omitted. got %#v (candidate keys %v)",
			c["rule"], c["unknown"], keysOf(c))
	}
	return pairs
}

// hasUnknown reports whether the pairs carry exactly {key, reason}.
func hasUnknown(pairs []unknownPair, key, reason string) bool {
	return slices.Contains(pairs, unknownPair{Key: key, Reason: reason})
}

// unknownKeys reports whether ANY pair names key, whatever its reason.
func unknownKeys(pairs []unknownPair, key string) bool {
	for _, p := range pairs {
		if p.Key == key {
			return true
		}
	}
	return false
}

// --- candidate navigation ------------------------------------------------

// nextCandidates decodes the `candidates` array off a `flow next` success
// payload, FAILING when it is not an array of objects.
func nextCandidates(t *testing.T, data map[string]any) []map[string]any {
	t.Helper()

	candidates, ok := objectsAt(data, "candidates")
	if !ok {
		t.Fatalf("`candidates` is not an array of objects: %#v",
			data["candidates"])
	}
	return candidates
}

// candidateRules returns the sorted rule ids of a candidate list, which is
// the form every set assertion in this suite compares.
func candidateRules(t *testing.T, data map[string]any) []string {
	t.Helper()

	var out []string
	for _, c := range nextCandidates(t, data) {
		id, _ := c["rule"].(string)
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

// candidateNamed returns the candidate for rule id, or nil.
func candidateNamed(t *testing.T, data map[string]any, id string) map[string]any {
	t.Helper()

	for _, c := range nextCandidates(t, data) {
		if c["rule"] == id {
			return c
		}
	}
	return nil
}

// requireCandidate returns the candidate for rule id, FAILING when the row
// was not reported.
func requireCandidate(t *testing.T, data map[string]any, id string) map[string]any {
	t.Helper()

	c := candidateNamed(t, data, id)
	if c == nil {
		t.Fatalf("row %q was not reported as a candidate; reported = %v",
			id, candidateRules(t, data))
	}
	return c
}

// --- invocation helpers --------------------------------------------------

// nextArgs builds a `flow next --as=json` argv over model and the given
// artifact bindings, appending extra.
func nextArgs(model string, bindings []string, extra ...string) []string {
	args := []string{"flow", "next", "--model", model}
	for _, b := range bindings {
		args = append(args, "--artifact", b)
	}
	args = append(args, extra...)
	return append(args, "--as=json")
}

// runNext drives `flow next --as=json` and returns the decoded payload.
func runNext(t *testing.T, model string, bindings []string, extra ...string) map[string]any {
	t.Helper()

	return flowData(t, requireSuccess(t, nextArgs(model, bindings, extra...)...))
}

// seedMatchArtifact seeds the `state` role for a 0011 fixture and returns
// its binding.
func seedMatchArtifact(t *testing.T, model string, writes ...string) string {
	t.Helper()

	return artifactBinding(flowMatchRole, seedArtifact(t, model, writes...))
}

// readersOf reads the payload's `readers` list, sorted.
func readersOf(t *testing.T, data map[string]any) []string {
	t.Helper()

	readers, ok := stringsAt(data, "readers")
	if !ok {
		t.Fatalf("`readers` is not an array of ids: %#v", data["readers"])
	}
	sort.Strings(readers)
	return readers
}

// viewKeysOf returns the key set of the payload's assembled view — the
// union of `owned` and `observed` — which is what `assembledView` holds.
func viewKeysOf(t *testing.T, data map[string]any) []string {
	t.Helper()

	seen := map[string]bool{}
	for _, field := range []string{"owned", "observed"} {
		m, ok := data[field].(map[string]any)
		if !ok {
			t.Fatalf("`%s` is not an object: %#v", field, data[field])
		}
		for k := range m {
			seen[k] = true
		}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// jsonRoundTrip re-encodes a decoded payload to canonical JSON, so two
// payloads can be compared for byte-identity independently of map order.
func jsonRoundTrip(t *testing.T, v any) string {
	t.Helper()

	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("re-encoding the payload failed: %v", err)
	}
	return string(b)
}

// --- kernel-untouched helpers --------------------------------------------

// resolveRefusalKindStrings returns the kernel's closed refusal-kind set as
// strings. `RefusalKinds()` is the kernel's own enumerator, so reading it
// here asserts the closure without reaching into `internal/resolve` or
// changing anything under it.
func resolveRefusalKindStrings() []string {
	kinds := resolve.RefusalKinds()
	out := make([]string, 0, len(kinds))
	for _, k := range kinds {
		out = append(out, string(k))
	}
	return out
}

// gitDiffStat returns `git diff --stat <base> -- <path>` against the branch
// point, which is the ONE mechanically checkable form of "the kernel is
// untouched" (`0011:S5`, MVV 9). A green kernel suite is its corroboration,
// not a second criterion.
//
// It resolves the base as the merge base with the default branch and skips
// rather than fails when git is unavailable or the tree is not a checkout:
// the claim is about a DIFF, and a harness that cannot see one has nothing
// to report either way.
func gitDiffStat(t *testing.T, path string) string {
	t.Helper()

	root := repoRootFor(t)
	base, err := exec.Command("git", "-C", root, "merge-base", "HEAD", "main").Output()
	if err != nil {
		t.Skipf("cannot resolve the merge base with `main`: %v — the "+
			"kernel-diff claim needs a branch point to compare against", err)
	}
	out, err := exec.Command("git", "-C", root, "diff", "--stat",
		strings.TrimSpace(string(base)), "--", path).Output()
	if err != nil {
		t.Skipf("`git diff --stat` failed: %v", err)
	}
	return string(out)
}

// repoModelPath returns the absolute path of the transition model checked
// into this repo (`models/rdr.toml`), which S1 and MVV 2-4 run against.
//
// It is the model whose narrowing the RDR's problem statement is about: 22
// rules, 21 reported by the stripped-match predicate at `stage=resolved`,
// 3 by C1's. It is also the model where C2's filter/dedup ORDER is
// load-bearing, because `stage` is BOTH a `[rule.match]` key and a
// `[rule.write]` key and so enters `RequiresOwned`.
func repoModelPath(t *testing.T) string {
	t.Helper()

	return filepath.Join(repoRootFor(t), checkedInModelPath)
}
