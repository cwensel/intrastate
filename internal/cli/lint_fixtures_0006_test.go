package cli

// RDR 0006 — the MVV fixture corpus.
//
// Each model is a complete, self-contained conforming document: RDR 0002's
// loader must accept it, so the defect under test is LINT's and not a load
// failure. `[model].id` is `lintfix` throughout, which is what every
// finding's `model` field must carry (REQ-13).

// mvvHeader is the preamble every fixture shares, parameterised by its tag
// declarations and accessor bindings. Root scalar keys precede all tables,
// because TOML binds a bare key to the most recent table header.
const mvvHeader = `outcomes = ["go", "stop"]
`

const mvvModelTable = `
[model]
id = "lintfix"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true
`

// mvvStatus declares the workhorse owned enum.
const mvvStatus = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true
`

// mvvFlag declares a required owned bool — a finite guard dimension whose
// atoms cannot refuse.
const mvvFlag = `
[tags.flag]
provenance = "owned"
kind = "bool"
single_valued = true
required = true
`

// mvvOpt declares an OPTIONAL owned enum. An unmarked key defaults to
// optional (RDR 0003), so a value atom over it can refuse
// `guard_unevaluable`.
const mvvOpt = `
[tags.opt]
provenance = "owned"
kind = "enum"
domain = ["p", "q"]
single_valued = true
`

// mvvFree declares an owned `scalar` — no finite domain, so any dimension
// over it is not projectable.
const mvvFree = `
[tags.free]
provenance = "owned"
kind = "scalar"
required = true
`

// mvvAlways declares an owned key carrying the EXPLICIT always-present
// marker.
const mvvAlways = `
[tags.always]
provenance = "owned"
kind = "enum"
domain = ["p", "q"]
single_valued = true
required = true
`

// mvvWide declares an owned int whose finite domain far exceeds the
// published product bound.
const mvvWide = `
[tags.n]
provenance = "owned"
kind = "int"
min = 0
max = 1000000
single_valued = true
required = true
`

func mvvAccessors(keys string) string {
	return `
[read.own]
role = "t"
path = "t.own"
keys = ` + keys + `
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ` + keys + `
timeout = "2s"
read_back = true
`
}

// --- the illegal matrix, one model per blocking invariant class ----------

// mvvNoRoot declares no initial owned state. Never "nothing reachable,
// therefore clean" — the disposition table routes it to
// `graph-dangling-edge` naming the missing declaration.
const mvvNoRoot = mvvHeader + `terminal = ["done"]
` + mvvModelTable + mvvStatus + `
[read.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"
read_back = true

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`

// mvvDeadEnd reaches `status = b`, which satisfies no declared terminal
// and is the source of no non-escape row.
const mvvDeadEnd = mvvHeader + `terminal = ["never"]
` + mvvModelTable + mvvStatus + `
[read.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"
read_back = true

[initial]
status = "a"

[context.never]
[context.never.match.status]
eq = "a"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`

// mvvOrdinaryOverlap carries two ordinary rows in one group that every
// assignment enables together.
const mvvOrdinaryOverlap = mvvHeader + `terminal = ["done"]
` + mvvModelTable + mvvStatus + `
[read.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"
read_back = true

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "over-one"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "over-two"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`

// mvvEscapeOverlap carries two ESCAPE rows sharing the `no_match` class
// whose assignments overlap — the escape population's own overlap, which
// is reported once per shared class.
const mvvEscapeOverlap = mvvHeader + `terminal = ["done"]
` + mvvModelTable + mvvStatus + `
[read.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"
read_back = true

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "esc-one"
escape = ["no_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"

[[rule]]
id = "esc-two"
escape = ["no_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
`

// mvvCoverageGap leaves the `flag = false` assignment uncovered over a
// fully provable product.
var mvvCoverageGap = mvvHeader + `terminal = ["done"]
` + mvvModelTable + mvvStatus + mvvFlag + mvvAccessors(`["flag", "status"]`) + `
[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "only-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"
`

// mvvNonFinite reads a `scalar` guard dimension: no finite declared
// domain, so the claim is withheld with reason `dimension-not-finite`.
var mvvNonFinite = mvvHeader + `terminal = ["done"]
` + mvvModelTable + mvvStatus + mvvFree + mvvAccessors(`["free", "status"]`) + `
[initial]
status = "a"
free = "anything"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "reads-free"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.free]
eq = "x"
[rule.write]
status = "b"
`

// mvvWithheld carries a value atom over an OPTIONAL key, so the row can
// refuse `guard_unevaluable` and the claim is withheld with reason
// `row-can-refuse`, naming the row and the atom.
var mvvWithheld = mvvHeader + `terminal = ["done"]
` + mvvModelTable + mvvStatus + mvvOpt + mvvAccessors(`["opt", "status"]`) + `
[initial]
status = "a"
opt = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "can-refuse"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
eq = "p"
[rule.write]
status = "b"
`

// mvvAlwaysPresentOwned declares an always-present OWNED key that
// `[initial]` omits, so it is absent at the root.
var mvvAlwaysPresentOwned = mvvHeader + `terminal = ["done"]
` + mvvModelTable + mvvStatus + mvvAlways + mvvAccessors(`["always", "status"]`) + `
[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`

// mvvOwnedBeforeWrite reads an owned tag by a GUARD atom that no row
// writes and `[initial]` does not establish. `RequiresOwned` names only
// `status`, so a lint reading that field never fires here.
var mvvOwnedBeforeWrite = mvvHeader + `terminal = ["done"]
` + mvvModelTable + mvvStatus + mvvAlways + mvvAccessors(`["always", "status"]`) + `
[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "reads-unwritten"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.always]
eq = "p"
[rule.write]
status = "b"
`

// mvvTerminalEscape declares NO terminal and ends at `status = b`,
// relying on lint to infer an implied terminal from the missing row.
const mvvTerminalEscape = mvvHeader + mvvModelTable + mvvStatus + `
[read.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"
read_back = true

[initial]
status = "a"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`

// mvvProductTooLarge declares a finite int domain whose product exceeds
// the published bound.
var mvvProductTooLarge = mvvHeader + `terminal = ["done"]
` + mvvModelTable + mvvStatus + mvvWide + mvvAccessors(`["n", "status"]`) + `
[initial]
status = "a"
n = 0

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "wide"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.n]
gte = 5
[rule.write]
status = "b"
`

// --- the multi-defect group ----------------------------------------------

// mvvMultiDefect carries several INDEPENDENT decidable defects so a
// first-failure engine cannot pass: an overlapping ordinary pair, two rows
// that can refuse `guard_unevaluable`, and a guard read of an owned key no
// row writes.
//
// All four rows bind the SAME recognized outcome, so they form ONE scoped
// row group per REQ-19 — the RDR names a multi-defect GROUP (0006:1450),
// and grouping is by the authored match pattern (gate.md §1). Splitting
// `over-*` onto `go` and `refuse-*` onto `stop` made two groups carrying
// one defect class each, which is not the scenario the record prescribes.
var mvvMultiDefect = mvvHeader + `terminal = ["done"]
` + mvvModelTable + mvvStatus + mvvOpt + mvvAlways +
	mvvAccessors(`["always", "opt", "status"]`) + `
[initial]
status = "a"
opt = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "over-one"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "over-two"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "refuse-one"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
eq = "p"
[rule.write]
status = "b"

[[rule]]
id = "refuse-two"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
eq = "q"
[rule.guard.all.always]
eq = "p"
[rule.write]
status = "b"
`

// --- the legal advisory matrix -------------------------------------------

// mvvLegalAdvisory is a CLEAN model carrying all three advisory codes the
// MVV names: a group closed by a bare escape row
// (`graph-coverage-closed-by-escape`), a redundant row whose accepted
// assignments are a proper subset of a sibling's (`graph-redundant-row`),
// and a rule no reachable owned-state node satisfies
// (`graph-unreachable-rule`).
var mvvLegalAdvisory = mvvHeader + `terminal = ["done"]
` + mvvModelTable + `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b", "orphan"]
single_valued = true
required = true
` + mvvFlag + mvvAccessors(`["flag", "status"]`) + `
[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "wide"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "narrow"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"

[[rule]]
id = "only-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"

[[rule]]
id = "bare-rescue"
escape = ["no_match", "ambiguous_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"

[[rule]]
id = "orphaned"
[rule.match.status]
eq = "orphan"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
