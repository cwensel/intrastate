Model: claude-opus-5

# RDR 0019 — Critique Grounding Pass

Each claim below was tested against `main` by reading the named symbol, and
where the answer turned on runtime behaviour, by compiling a throwaway probe
against the real package (probes were removed; no source file was modified).

---

## G-1 — registry write-loop exhaustiveness: **PARTIAL**

`internal/cli/flowbind/registry.go::Registry` (write loop), and
`internal/cli/flowbind/registry.go::commandBacked`.

(a) **CONFIRMED.** The write loop is a `switch { }` with exactly TWO explicit
arms, and the default is the pre-switch initializer, not a `default:` clause:

    var binding accessor.Binding = &Writer{Path: acc.Path}
    switch {
    case len(acc.Edit) != 0:
        binding = NewEditWriter(acc, name)
    case commandBacked(acc):
        binding = &cmdbind.Writer{Accessor: acc, Name: name, Config: cfg}
    }

The `&Writer{Path: acc.Path}` fallthrough is real, so the substance of the
claim holds; the shape is an initializer-plus-two-arm switch rather than a
literal `default:` arm.

(b) **CONFIRMED.**

    func commandBacked(acc table.Accessor) bool {
        return len(acc.Command) != 0 || acc.Path == ""
    }

(c) **CONFIRMED** — unreachable. The `&Writer` initializer survives only when
both arms fail; `commandBacked` false requires `acc.Path != ""`. So a
`*flowbind.Writer` with an empty `Path` cannot be constructed here. This is
exactly the property the arm ordering is commented to protect: the `edit` arm
is tested FIRST and on carrier presence, because `commandBacked`'s
`acc.Path == ""` disjunct IS the residue path.

(d) **PARTIAL — there is a fourth, if the question is scoped to the whole
`Registry` function rather than the write loop.** The write loop constructs
exactly three: `*flowbind.Writer`, `accessor.WriteBinding` via
`flowbind.NewEditWriter` (concrete `*flowbind.EditWriter` —
`internal/cli/flowbind/edit.go:74` returns the interface type, not the
struct), and `*cmdbind.Writer`. There is no fourth WRITE binding. The registry
as a whole also builds `Reader`/`cmdbind.Reader` and `Gate`/`cmdbind.Gate`,
but none is a write binding.

---

## G-2 — read accessor is a separate resolution: **CONFIRMED**

`internal/accessor/model.go::Registry.readerFor`, used at
`internal/accessor/executor.go:311`.

    func (reg Registry) readerFor(role string) (Definition, bool) {
        for _, d := range reg.Definitions {
            if d.Identity.Capability == CapRead && d.Accessor.Role == role {
                return d, true
            }
        }

The read-back re-read resolves by ROLE over `CapRead` definitions only — it
never consults the write definition's binding. The two construction sites in
`flowbind.Registry` are separate loops over separate maps (`m.Readers`,
`m.Writers`), each running its own `commandBacked` test on its own
`table.Accessor`. So yes: a model can bind a file-backed `[write.x]`
(`path` carrier) and a command-backed `[read.y]` (`command` carrier) on the
same role, and the two carriers are independently determined. The executor
comment states the coupling explicitly: "The read-back re-read goes through
the READ path over the SAME caller-supplied role the write binding names."

---

## G-3 — numeric spelling round-trip: **CONFIRMED (byte divergence is real)**

`internal/table/load.go::valueMembers` and
`internal/cli/flow_input.go::canonicalValue` /
`internal/table/load.go::conformKind`.

`valueMembers` renders via `strconv` exactly as claimed:

    case int64:   return []string{strconv.FormatInt(t, 10)}, nil
    case float64: return []string{strconv.FormatFloat(t, 'g', -1, 64)}, nil

`canonicalValue` stores the argv string VERBATIM for every non-set kind — it
returns `value` unchanged after conformance, with no re-encoding:

    if err := table.ConformValue(decl, value); err != nil { ... }
    return value, nil

`conformKind` has arms for `int` and `bool` ONLY — there is **no `float`
arm**, and `float` is not a declared kind at all: `internal/table/model.go`
fixes the closed five-token vocabulary as
`{"enum", "bool", "int", "set", "scalar"}`. A `float`-kind declaration would
be refused upstream as a malformed kind, so the "does conformKind have a float
arm" question is moot in the direction that matters: unmatched kinds
(including `scalar`) fall through `conformKind` unchecked.

`conformKind` DOES accept `+5` for `int` via `strconv.Atoi` (Atoi parses a
leading `+`), so `+5` is admitted and stored verbatim as `"+5"`.

Probe results against the real package:

    valueMembers(float64(1.0))       = []string{"1"}
    valueMembers(int64(5))           = []string{"5"}
    conformKind({Kind:"int"}, "+5")  = nil   (admitted)

**Decisive:** for `[initial] threshold = 1.0` the init route stores `"1"`
(TOML decodes to `float64(1)`; `FormatFloat(1, 'g', -1, 64)` is `"1"`), while
`set-state --write threshold=1.0` stores `"1.0"` verbatim. The byte divergence
is CONFIRMED. The same divergence holds for `+5` vs `5` on an `int` tag.

---

## G-4 — refusal code spellings normative: **PARTIAL**

`internal/cli/flow_input.go` — the code table header, and
`internal/cli/flow_input.go::codeWriteEditRefused`.

The header comment says exactly what the claim says, verbatim:

    // --- the stable code table (`0005:FM`) -----------------------------------
    //
    // The spellings are normative and the group fixes the exit. They are
    // declared as constants in one place so a producer cannot spell one
    // slightly differently at a second site.

BUT the claim's second half is **REFUTED for `codeWriteEditRefused`
specifically**. It is a fixed literal `const` (`= "flow-write-edit-refused"`),
so it is a constant — yet its own doc comment carves it OUT of the header's
normativity:

    // codeWriteEditRefused is RDR 0028 `0028:C1.3` EXIT GROUP:'s distinct
    // code for a declared-line-edit refusal decided before mutation. Its
    // SPELLING is this stage's to choose and is explicitly non-normative
    // — no test may pin the string; what the contract fixes is the exit
    // GROUP and the `findings[]` carriage.

Its sibling `codeRequestRefused` carries the same carve-out ("Its SPELLING is
this stage's to choose on the same terms as its sibling's"). So: header says
normative; two constants in that same table say explicitly non-normative. Any
0019 argument leaning on "all spellings in this table are normative" is
unsound for these two.

---

## G-5 — empty scalar admission: **CONFIRMED (reachable today)**

`internal/table/load.go::loader.loadInitial`,
`internal/table/model.go::IsDecisionTable`,
`internal/cli/flowbind/flowbind.go::Writer.Apply`.

`loadInitial` ADMITS an empty string. Its guards are: undeclared tag,
`valueMembers` error, the `<clear>` reserved sentinel, `conform`, and an ARITY
check (`decl.Kind != "set" && len(members) != 1`). An empty string is a
one-member sequence `[""]`, so arity passes, and `conform` -> `conformKind`
has no `scalar` arm, so it returns nil. There is no emptiness guard anywhere
on the path.

`IsDecisionTable` does NOT exclude it — it reads one field and nothing else:

    func IsDecisionTable(m *Model) bool {
        return m != nil && m.Class == ClassDecisionTable
    }

`Writer.Apply` stores it without complaint: the only special case is
`accessor.IsClear(t.Value)` (removal); everything else is `s[t.Key] = t.Value`.

Probe against the real loader, with a valid model declaring
`[tags.note] kind="scalar" provenance="owned"` and `[initial] note = ""`:

    LOADED OK; Initial = []table.TagValue{{Key:"note", Value:[]string{""}}}

And end-to-end through the binding:

    Apply err = <nil>
    artifact bytes = {"note":""}
    Read vals = [{Key:"note", Value:"", Absent:false}]

**Decisive: an empty-scalar `[initial]` value is reachable at a seed path
TODAY.** It is not rejected at load, it round-trips to the artifact, and it
reads back PRESENT (`Absent:false`) — distinct from a cleared key, which
leaves the object entirely.

Note the one adjacent asymmetry: the CLI argv path REFUSES what the model path
admits. `canonicalValue` rejects `value == ""` for non-set kinds with "the tag
`<key>` was given an empty value". So `[initial] note = ""` is admitted while
`set-state --write note=` is refused.

---

## G-6 — allow-commands refusal site: **CONFIRMED (fires at invocation)**

`internal/cli/cmdbind/cmdbind.go::spawn` (the `!cfg.AllowCommands` arm,
cmdbind.go:245) and `internal/accessor/model.go::Registry.AllowCommands`.

The refusal is the FIRST rung of the pre-spawn ladder inside `spawn`:

    if !cfg.AllowCommands {
        return invocation{}, refuse("command execution requires the " +
            "`allow_commands` opt-in (--allow-commands); the accessor `" +
            name + "` declares a command and none was given")
    }

`flowbind.Registry` PASSES the gate (into `cmdbind.Config` and onto
`accessor.Registry.AllowCommands`) but never refuses on it at construction —
construction always succeeds and yields a binding. The `Registry.AllowCommands`
field comment says it is carried "here as STATE" precisely because `cmdbind`
imports `accessor` and the reverse would be an import cycle.

**Is a refusal reachable WITHOUT invoking an accessor? YES — one.**
`internal/accessor/executor.go:311` uses `readerFor` for a C1.3 read-back
PRE-CHECK that refuses before `Apply`: when an `edit` write's role reader is
command-backed and the gate is off, the write refuses without spawning
anything. Its comment: "Refusing here — before `Apply`, unconditionally — is
what stops a forgotten `--allow-commands` from mutating the artifact and then
reporting `read_back_incomplete`". So the gate has TWO refusal sites: the
invocation-time ladder in `spawn`, and this pre-mutation executor check. Both
are past registry construction; neither is at binding-construction time.

---

## G-7 — ReadBinding cardinality: **CONFIRMED**

`internal/accessor/binding.go::ReadBinding.Read`.

    Read(ctx context.Context, art Artifact, requested []string) (
        values []KeyValue, unreadable []string, err error)

Key-scoped exactly as claimed: one `KeyValue` per requested key plus an
`unreadable` list. It reports NO cardinality — there is no "does this store
hold any key" answer anywhere on the interface. `KeyValue.Absent` is per-key
and is documented as a VALUE ("the binding established the key is not
carried"), not a store-level count.

Exactly TWO implementers, no third:
`internal/cli/flowbind/flowbind.go::Reader.Read` (flowbind.go:189) and
`internal/cli/cmdbind/cmdbind.go::Reader.Read` (cmdbind.go:1123). A sweep for
non-test `Read(` methods matching the seam signature returns only those two
plus `accessor/executor.go::Executor.Read`, which is the executor's own
caller-facing method, not a `ReadBinding`.

`flowbind::store` IS package-private (`type store map[string]string`) with no
exported surface returning a key count. The package's exported surface is
`Reader`, `Writer`, `Gate`, `EditWriter`, `NewEditWriter`, `Registry`,
`OwnedTags` — none exposes cardinality.

---

## G-8 — verb cardinal sites: **CONFIRMED (six code/doc sites)**

Hard-coded "four verbs" or equivalent cardinal, outside RDR prose:

1. `internal/cli/flow.go:5` — package comment, "One group, four verbs".
2. `internal/cli/flow.go:42` — `newFlowCmd` Long body, "The four verbs are the
   whole skill-integration surface".
3. `internal/cli/flow.go:107` — `flowExtendedDesc`, "The four verbs divide one
   job".
4. `internal/cli/flow_exec.go:3` — "the shared execution path the four verbs
   sit on".
5. `docs/cli-reference.md:181` and `:218` — both strings, mirroring flow.go.
6. `docs/model-authoring.md:369` — "**The same four verbs.**".

Also `internal/cli/flow_state.go:405` ("The four the contract...") is a
cardinal about REFUSALS, not verbs — not counted.

**`llms.txt` does NOT hard-code the cardinal.** It enumerates the verbs as
bullet lines (`intrastate flow read-state`, `intrastate flow set-state`, ...)
with no "four". **`docs/cli-output-contract.md` does NOT hard-code it either**
— it says "each verb" generically. README has no such cardinal. So two of the
six sites the claim asks about are clean.

**Test corpus asserts the verb set verb-by-verb: YES.**
`internal/cli/flow_surface_0005_test.go:21` fixes
`var flowVerbs = []string{"next", "resolve", "read-state", "set-state"}`, and
`TestReq1And3_FlowGroupExposesExactlyTheFourNormativeVerbs` asserts SET
EQUALITY against `flowGroupNames`. A fifth registered verb fails this test,
independently of any doc string. Registration is
`internal/cli/flow.go::newFlowCmd`'s single `cmd.AddCommand(...)` of four
constructors.

---

## G-9 — sealed-artifact internals: **CONFIRMED**

All in `internal/cli/flowbind/flowbind.go`.

`unreachable` IS a pure suffix test on the DECLARED path — no filesystem
access, no ambient state:

    func unreachable(path string) bool {
        return strings.HasSuffix(path, "."+unreachableSuffix) ||
            strings.HasSuffix(path, "-"+unreachableSuffix)
    }

`Writer.Apply` seals and mutates in ONE `save`. The mutation is applied into
`s`, then the seal is set or deleted on the same map, and there is a single
terminal `return save(art.Path, s)`. The comment fixes the atomicity as
load-bearing: "The mutation lands — `s` already carries it — AND the artifact
is marked unverifiable, in ONE write... a seal that discarded the mutation
would make the refusal's own detail false." The else-branch `delete(s,
sealedKey)` makes the seal NON-monotonic — it names the LAST write's locator.

`Reader.Read` DOES short-circuit on a sealed store, returning every requested
key as unreadable and no values:

    if _, sealed := s[sealedKey]; sealed {
        return nil, slices.Clone(requested), nil
    }

**Contract vs private implementation detail:**

- CONTRACT (exported / documented): `Reader`, `Writer`, `Gate` and their
  `Path` fields; the seam methods `Read` / `Apply` / `Gate` / `Capability` /
  `Invocations`; and the two package-comment format properties (key PRESENCE
  distinguishes empty-set from cleared; values stored VERBATIM).
- PRIVATE IMPLEMENTATION DETAIL: `unreachable`, `unreachableSuffix`,
  `gatePrefix`, `verdictFor`, `store`, `load`, `save`, `sealedKey`
  (`"\x00flow.readback-unreachable"`), and `sealedMarker`. All lowercase, none
  exported. The package comment explicitly labels the artifact format an
  "IMPL-DECISION recorded in the deviations artifact", noting "RDR 0005 fixes
  the CLI contract, not a file schema, and no REQ constrains this."

A 0019 S9 fixture depending on the sealed-key spelling, the suffix vocabulary,
or the one-save ordering is depending on private internals of another RDR's
implementation, not on contract.

---

## Cross-cutting note

Two claims came back PARTIAL for the same reason: a general statement in the
source is true, but a specific instance carves itself out. G-4's code table
declares its spellings normative while two constants inside it declare
themselves explicitly non-normative. G-1(d)'s write-binding enumeration is
complete for the write loop but the registry builds six binding types overall.
Neither carve-out is visible from the general statement alone.
