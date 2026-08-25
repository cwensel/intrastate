# REQ List — RDR 0002 Transition Table As Reviewable Data

Phase 0 audit artifact. Every clause below is a testable obligation drawn from
`docs/rdr/0002-transition-table-as-reviewable-data.md`. Quotes are verbatim,
copied from the projector (`rdr inspect --select <id>`) or from the record
itself, never transcribed.

Element ids are carried where the REQ derives from a labelled contract
(`0002:C1` … `0002:C24`, `0002:MVV`). Source sections are abbreviated:

- `NC` = Proposed Solution / Technical Design / Normative Contracts
- `LBD` = Proposed Solution / Technical Design / Load-Bearing Decisions
- `TD` = Proposed Solution / Technical Design (prose, outside fences)
- `AP` = Proposed Solution / Approach
- `RT` = Proposed Solution / Technical Design / Round-Trip / Inverse Invariants
- `MC` = Proposed Solution / Technical Design / Conditional Mini-Checks
- `CA` = Research Findings / Critical Assumptions
- `RM` = Trade-offs / Risks and Mitigations
- `FM` = Trade-offs / Failure Modes
- `EIA` = Implementation Plan / Existing Infrastructure Audit
- `MVV` = Implementation Plan / Minimum Viable Validation
- `IP` = Implementation Plan (phases)
- `ND` = Implementation Plan / New Dependencies
- `TS` = Validation / Testing Strategy
- `PE` = Validation / Performance Expectations

**Two standing corrections apply to this REQ set and are marked at each REQ
they touch** (RDRs are never amended; the correction lives here and in
`deviations.md` D4/D5):

1. **§D13 (§JD-22) is a LANDING this RDR owes, and its clause is not in the
   record.** The record's Status line says "§D13 assigns this RDR the normative
   set-value encoding clause (a landing, not a citation)" (`0002:11-12`), but no
   `normative` fence in the body states it — `grep '§D13\|JSON array'` over the
   record hits only the Status line. JDR 0001 §D13 fixes the byte form:
   *"`Tag.Value` stays `string`; a set crosses as its canonical JSON array —
   members sorted, duplicate-free, compact encoding — and 0002 declares it."*
   REQ-93 states it. RDR 0007 already wrote its `contains` leg against §D13
   directly (`BUILD-ORDER.md`, "The one seam this order strains"), so this build
   must **match** that form, not redefine it. Recorded at `deviations.md` D5.
2. **The Prerequisites' "RDR 0007's reshape is unimplemented" is stale.** The
   record sequences Phases 2–3 behind a reshape it says has "no owner"
   (`0002:2073-2077`), and marks REQ-quoted assertions "unsatisfiable until it
   lands". **It has landed**: `internal/resolve/guard.go` now exports
   `GuardAtom{Key, Operator, Literal, Block}`, `OpExists`, `LiteralTrue`,
   `LiteralFalse`, `BlockAll`, `BlockUnless`, `BlockMatch`, and
   `resolve.go::Row.Guard` is `[]GuardAtom`. Per `BUILD-ORDER.md` 0007 is run 1
   of 8 and completed (`0007/artifacts/status.md`: state COMPLETE). Every REQ
   below marked *(was blocked)* is therefore satisfiable in this build.
   Recorded at `deviations.md` D4.

---

## A. Format and carrier

- [REQ-1] `0002:C1` "The transition model MUST be authored as sparse TOML data, not generated code and not a fully expanded Cartesian-product table." — (NC)

- [REQ-2] `0002:C2` "The source schema is the closed layout JDR 0001 §D7 fixes: root `outcomes`, root `terminal`, `[model]` (with the free-form sub-table `[model.metadata]`), `[initial]`, `[tags.<tag>]`, `[read.<id>]`, `[write.<id>]`, `[gate.<id>]`, `[context.<id>]`, `[[rule]]`, and `[dump]`." — (NC)

- [REQ-3] `0002:C2` "Context predicates live under `[context.<id>.match.<tag>]`; rule predicates live under `[rule.match.<tag>]`, `[rule.guard.all.<tag>]`, and `[rule.guard.unless.<tag>]`; writes live under `[rule.write]`; explicit clears live in a rule-level `clear` list; gate accessor references live in a rule-level `gate` list of `[gate.<id>]` ids; modeled escape rows live in a rule-level `escape` list." — (NC)

- [REQ-4] `0002:C2` "`[model]` additionally carries an optional human `description`, and each `[[rule]]` an optional `source` — a provenance *annotation* carried alongside the locator, which is derived from `(model id, rule id)` and never from `source` (below). Both are admitted keys" — (NC)

- [REQ-5] `0002:C2` "No other root key or table is admitted (strict decoding, below)." — (NC)

- [REQ-6] `0002:C3` "**Strict decoding is an obligation on this format, not a property of a library.** The decoder MUST reject unmapped keys so an unknown schema field is a stable refusal rather than a silent no-op." — (NC)

- [REQ-7] "Use `github.com/pelletier/go-toml/v2` as the TOML parser candidate." — (ND)

- [REQ-8] `0002:EIA` "The load entry therefore takes **already-read bytes plus a source id** for the locator, not a filesystem path — this package performs no file I/O and no path resolution." — (EIA, Config discovery)

## B. `[model.metadata]` extension namespace

- [REQ-9] `0002:C2` "**`[model.metadata]` is the one sanctioned extension namespace.** The loader MUST decode it as a free-form table, carry it through to the normalized model untouched, and MUST NOT interpret any key in it; strictness applies everywhere else." — (NC)

- [REQ-10] `0002:C2` "It is a **model-level field and reaches no candidate row**, so it is outside the dump's field list and outside the Round-Trip invariant" — (NC)

- [REQ-11] `0002:C2` "the normalized model MUST expose the decoded table as a field (Phase 2 deliverable), and scenario 1 MUST assert **value and nesting equality** against the authored table, not merely that its top-level key names survive" — (NC)

- [REQ-12] `0002:C2` "Its internal shape is deliberately unconstrained (arbitrary keys, values, and nesting)" — (NC)

## C. Model header and version gate

- [REQ-13] `0002:C3` "`[model]` MUST contain `id` and `version`. Version `1` is the only version this RDR accepts; any other version MUST be refused before normalization." — (NC)

- [REQ-14] `0002:C3` "**The version check MUST run before strict field validation, not merely before normalization.** Loading MUST therefore proceed in two passes: read `[model]` permissively enough to obtain `version`, refuse on any value but `1`, and only then decode the document strictly." — (NC)

- [REQ-15] `0002:C3` "the ordering is `malformed TOML` → version gate → strict decoding → the remaining categories." — (NC)

- [REQ-16] `0002:C3` "**Load is fail-fast: the first category a document trips is the refusal.** Beyond the two fixed precedences above, the order in which independent defects are checked is deliberately **unspecified**" — (NC)

- [REQ-17] `0002:C3` "What is forbidden is the single `load` entry point returning a list instead of a refusal" — (NC)

- [REQ-18] `0002:C3` "Because this class trips no category, it is **asserted positively, not by refusal**: the oracle is a count over the normalized value — the atom count for the merge case (Testing Strategy 2's distinct-literal control) and the write value's member sequence for the write case — each asserted to survive normalization intact." — (NC)

## D. Accessor tables and bindings

- [REQ-19] `0002:C2` "Every entry carries `role`, `path`, `keys` (a non-empty list of declared tag keys), and `timeout` (a Go duration string). **Each of the four is required, and any of them absent, empty, or ill-formed is a `malformed accessor declaration`** — an absent or empty `role` or `path`, an absent or empty `keys` list, a `keys` member naming an undeclared tag (`unknown tag`), an absent `timeout`, a `timeout` that does not parse as a Go duration, or one that parses non-positive." — (NC)

- [REQ-20] `0002:C2` "A `[write.<id>]` entry additionally carries `read_back = true`; `read_back = false` or its absence on a write entry is the same category" — (NC)

- [REQ-21] `0002:C2` "The same id MAY appear under two capability tables — RDR 0004's identity is `(flow, name, capability)`." — (NC)

- [REQ-22] `0002:C2` "every key in a `keys` list MUST be a declared tag (`unknown tag`), and MUST NOT be `recognized`" — (NC)

- [REQ-23] `0002:C2` "every **owned** tag MUST be served by exactly one reader; an **observed** tag MAY be served by at most one reader, and zero is legal because an observed key may arrive from the caller (JDR 0001 §JD-9)" — (NC)

- [REQ-24] `0002:C2` "every key named by any rule's write block or clear list, and every key assigned in `[initial]`, MUST be served by exactly one writer (`malformed accessor binding`)" — (NC)

- [REQ-25] `0002:C2` "every key in a `[write.<id>]` entry's `keys` list MUST be an **owned** tag (`write to non-owned tag`)" — (NC)

- [REQ-26] `0002:C2` "every id in a rule's `gate` list MUST resolve to a `[gate.<id>]` entry (`unknown accessor`)" — (NC)

- [REQ-27] `0002:C2` "A violation of the second or third bullet is a `malformed accessor binding`; the fourth carries `write to non-owned tag`." — (NC)

- [REQ-28] `0002:C2` "A rule's `gate` list names the gate accessors whose `allow` the executor requires before applying that rule's plan; the list is carried on the normalized row and is part of its value." — (NC)

- [REQ-29] `0002:C2` "An escape rule MUST NOT carry a `gate` list (it yields no plan to gate) — `malformed escape declaration`." — (NC)

## E. Root, stop set, and terminal contexts

- [REQ-30] `0002:C2` "`[initial]` is a table of `tag = value` assignments; every key MUST be a declared owned tag and every value MUST be well-formed for that tag's declared kind and domain (`malformed initial declaration` otherwise; the `<clear>` sentinel is refused under the reserved-value rule)." — (NC)

- [REQ-31] `0002:C2` "Root `terminal` is a list of context ids, each of which MUST resolve (`unknown context`)" — (NC)

- [REQ-32] `0002:C2` "**A terminal context id is a reference, not a state name: normalization MUST dereference each to the context's explicit predicate set over owned tags**, in the same atom shape rules use, and the normalized value carries the predicate sets — never the bare ids." — (NC)

- [REQ-33] `0002:C2` "Whether a model *omits* `[initial]` or `terminal`, and whether a terminal context matches a non-owned tag, are RDR 0006's blocking findings, not load failures: the loader accepts the absence" — (NC)

## F. Rule shape, ids, and writes

- [REQ-34] `0002:C4` "Rule ids MUST be unique within a model, compared by exact byte equality; a duplicate rule id is a load failure." — (NC)

- [REQ-35] `0002:C4` "Each transition rule MUST contain a stable rule id, zero or more shared-context references, and a local match block. An ordinary transition rule MUST contain a write block, MAY contain a rule-level explicit clear list, and MAY contain a rule-level `gate` list. A write block MAY assign more than one tag." — (NC)

- [REQ-36] `0002:C4` "An escape rule MUST contain an `escape` list and MUST NOT contain a write block, clear list, or `gate` list, even an empty one (RDR 0009 binds the write-free obligation at the kernel boundary)." — (NC) — *see `deviations.md` D2: `write = []` is presence-keyed, not length-keyed.*

- [REQ-37] `0002:C4` "**A write replaces.** A write assigns a tag's whole value and supplants whatever was held; for a `set` kind the array literal is the whole new set. There is no accumulate form." — (NC)

- [REQ-38] `0002:C4` "a `set`-kind write value MUST normalize to an ordered, member-sorted sequence and MUST NOT be joined into a single delimiter-separated string, in the stored value or in any comparison derived from it." — (NC)

## G. Escape rows

- [REQ-39] `0002:C5` "An `escape` list MUST contain only resolver failure classes that RDR 0001 allows the table to model: `no_match` and `ambiguous_match`." — (NC)

- [REQ-40] `0002:C5` "Normalization MUST render an escape rule as a candidate row carrying its normal predicate set, outcome, source rule id, source locator, and modeled failure class list." — (NC)

- [REQ-41] `0002:C5` "An escape row rescues only resolves carrying the outcome it binds" / "There is no table-wide catch-all: covering an alphabet of N outcomes requires N escape rows." — (NC)

- [REQ-42] `0002:C5` "An escape rule MAY bind its outcome with an `in` atom and expand like any other rule … each expansion is a separate row rescuing its own outcome, and the expansion suffix distinguishes them." — (NC)

- [REQ-43] `0002:C5` "Row kind (`transition` / `escape`) is a **derived view property, never a row field**: escape identity is discriminated solely by a non-empty escape class list" — (NC)

- [REQ-44] `0002:C5` "a render-time view struct MAY materialize the computed column … What it MUST NOT do is let that view feed back into the normalized value or the kernel row" — (NC)

## H. Contexts and merging

- [REQ-45] `0002:C6` "Shared contexts MAY inherit from other contexts, but inheritance MUST normalize to an explicit predicate set before lint or resolution." — (NC)

- [REQ-46] `0002:C6` "The combined predicate set is a **set over the full atom identity** — the tuple `(key, block, operator token, literal)`, the same tuple the dump sorts on. Merging MUST NOT key on any proper prefix of it: two atoms agreeing on `(key, block, operator)` but differing in literal are **distinct atoms** and both survive the merge." — (NC)

- [REQ-47] `0002:C6` "Merging is idempotent on identical atoms: the same atom contributed by a rule and by one or more inherited contexts collapses to one. Inheritance therefore never overrides — it only accumulates." — (NC)

- [REQ-48] `0002:C6` "authoring two atoms on one key that no view can satisfy together yields a dead rule, which is RDR 0006's unreachable-rule finding, not a load failure here." — (NC)

## I. Atoms: blocks, identity, literals

- [REQ-49] `0002:C7` "Guard predicates MUST be represented as positive `all` predicates and negative `unless` predicates. Normalization MUST combine the match atoms and both guard blocks into one candidate-row predicate set before ambiguity checks, and each atom in that set MUST retain the key, operator token, literal, and the block it was authored in — a three-valued domain, `match`, `all`, or `unless`." — (NC)

- [REQ-50] `0002:C7` "An implementer building the atom type from RDR 0007 alone gets two members and MUST widen to three here." — (NC) — **satisfied by the shipped kernel**: `resolve.BlockMatch` exists (JDR 0001 §D12, landed in 0007 Phase 1).

- [REQ-51] `0002:C7` "Normalization MUST NOT fold `unless` atoms into `all` or `match` atoms into either guard block." — (NC)

- [REQ-52] `0002:C7` "one atom authored in **both** `all` and `unless` is two distinct atoms and both survive normalization; this is not a load failure." / "It MUST NOT be silently pruned at load" — (NC)

- [REQ-53] `0002:C8` "For an existence atom the normalizer MUST emit the kernel's exported constants verbatim — operator token `OpExists` and literal `LiteralTrue` or `LiteralFalse` (RDR 0007, existence-operator clause; JDR 0001 §D4) — and MUST reject at load, as a malformed predicate atom, any existence literal that is not one of those two forms." — (NC) — *(was blocked; the constants now exist in `internal/resolve/guard.go`.)*

- [REQ-54] `0002:C9` "Tag-key identity is exact byte equality on the post-parse key string at every stage — declaration lookup, context and rule predicate references, write and clear targets, and the keys a normalized row carries. The normalizer performs no case folding, trimming, or namespace rewriting" — (NC)

- [REQ-55] `0002:C9` "The canonical spelling of a tag is its `[tags.<tag>]` declaration key; every reference resolves to a declaration by exact match or fails `unknown tag`" — (NC)

- [REQ-56] `0002:C10` "A set-valued literal MUST normalize to an ordered sequence of its members, each compared byte-exactly; it MUST NOT be rendered into a single delimiter-joined string as its normalized value." — (NC)

- [REQ-57] `0002:C10` "Members sort byte-lexicographically so two authored orderings of one set are one literal" — (NC)

- [REQ-58] `0002:C10` "Wherever atoms are compared, keyed, deduplicated, or sorted — in particular the **merge key** … the literal field MUST be compared as the member sequence, element by element. An implementation MUST NOT derive that key, or any other identity, by joining members into a string." — (NC)

- [REQ-59] `0002:C10` "Producers rendering into that field MUST spell a set visibly as a set **with each member delimited unambiguously** — quoting each member is the cheap way" — (NC)

- [REQ-60] `0002:C11` "**Rule ids, outcome literals, and the members of any match-block `in` literal MUST NOT contain the expansion-suffix separator `#`.** … Violations are a malformed rule id, a malformed recognized outcome alphabet, and a malformed predicate atom respectively." — (NC)

- [REQ-61] `0002:C11` "The string `<clear>` MUST be refused at load wherever a tag value is authored — a write-block value, an `[initial]` assignment, or a predicate literal — as a `reserved tag value`." — (NC)

## J. Tag declarations and provenance

- [REQ-62] `0002:C21` "The model MUST declare every tag it matches or writes, including each tag's provenance: owned, observed, or recognized." — (NC)

- [REQ-63] `0002:C22` "A tag declaration also carries its **type model**, spelled as the wire keys `kind`, `domain`, `min`, `max`, `elements`, `single_valued`, and `required` (JDR 0001 §D7(iii); `required` is the optionality marker and defaults to optional)." — (NC)

- [REQ-64] `0002:C22` "This RDR owns where a declaration is authored — under `[tags.<tag>]`, beside `provenance`, with no accessor reference — and the two load categories that carry RDR 0003's rejection rules: `malformed tag declaration` … and `malformed predicate atom` (a literal outside the tag's declared domain …)." — (NC) — *see `deviations.md` D1 (kind vocabulary is RDR 0003's five tokens, not `string`) and D3 (write-value domain conformance).*

- [REQ-65] `0002:C22` "Normalization MUST carry every declared field through to the normalized model without loss" — (NC)

- [REQ-66] `0002:C12` "A tag declaration with provenance `recognized` MUST be named `recognized`, and no owned or observed declaration may take that name; violations fail in RDR 0008's `reserved_tag_key` category" — (NC)

- [REQ-67] `0002:C12` "the `outcomes` alphabet MUST be non-empty, duplicate-free, and MUST NOT contain the empty string." — (NC)

## K. Outcome binding and expansion

- [REQ-68] `0002:C13` "Every rule — ordinary or escape — MUST bind exactly one outcome. **Outcome binding reads the match blocks only** — the rule's local `match` block plus the `match` blocks of its inherited contexts. That set MUST contain exactly one atom on `recognized`, using `eq` or `in`, whose literal(s) are members of the `outcomes` alphabet." — (NC)

- [REQ-69] `0002:C13` "Normalization lifts that atom out of the predicate set into the row's outcome field" — (NC)

- [REQ-70] `0002:C13` "A rule binding zero outcomes, more than one `recognized` atom, or a literal outside the alphabet is a load failure." — (NC)

- [REQ-71] `0002:C13` "A `recognized` atom authored under `guard.all` or `guard.unless` MUST be refused at load as a malformed outcome binding — never lifted." — (NC)

- [REQ-72] `0002:C13` "Every `in` atom in a rule's match blocks — local `match` and inherited context `match`, on `recognized` or on any other declared tag — expands into one candidate row per member, and the rule's rows are the Cartesian product of its expanding atoms." — (NC)

- [REQ-73] `0002:C13` "In each expanded row the `in` atom becomes an `eq` atom on the chosen member, still carrying block `match`" — (NC)

- [REQ-74] `0002:C13` "The **expansion suffix** is the sequence of chosen members, one per atom with **more than one member**, taken in the atoms' sort order `(key, block, operator token, literal)`; it is empty when no such atom exists." — (NC)

- [REQ-75] `0002:C13` "A product row whose expanded atoms cannot be satisfied together (two `eq` on one key, different literals) is a dead row for RDR 0006, not a load failure." — (NC)

- [REQ-76] `0002:C13` "A single-member `in` expands to one row whose suffix is empty, so `eq = "x"` and `in = ["x"]` are one spelling of one edge and cannot mint two identities. A suffix is non-empty exactly when the rule produced more than one row." — (NC)

## L. Derived row fields

- [REQ-77] `0002:C14` "`Row.RequiresOwned` has no authored form. The normalizer MUST derive it for every row as the sorted, duplicate-free set of tag keys named by the rule's write block and clear list" — (NC)

- [REQ-78] `0002:C14` "An escape rule carries neither a write block nor a clear list, so a normalized escape row MUST carry an empty set (RDR 0007 A21)." — (NC)

- [REQ-79] `0002:C14` "this RDR is its producer (JDR 0001 §JD-3) and does not add guard-read keys to it." — (NC)

- [REQ-80] `0002:C15` "Normalization MUST populate the next-state tags with the tag values the rule's write block and clear list produce, and the writes with the owned-tag writes the accessor layer applies — the same rendered set, including `<clear>` entries." — (NC)

- [REQ-81] `0002:C15` "An escape row carries neither a write block nor a clear list, so it normalizes to a row with both empty" — (NC)

- [REQ-82] `0002:C15` "Normalization MUST populate both explicitly from the rule and MUST NOT populate one by aliasing the other" — (NC) — asserted by the **mutation test** in TS scenario 2, not by value comparison.

- [REQ-83] `0002:C23` "Clearing a tag MUST be represented by an explicit rule-level `clear` entry that normalization renders as a `<clear>` write. Absence from both the write block and the clear list MUST NOT imply deletion." — (NC)

## M. Kernel handoff and match/guard routing

- [REQ-84] `0002:C16` "The normalized candidate row carries the atoms as **one unified set with each atom's authored block retained**; that single field is what the dump contract lists, what the Round-Trip invariant compares, and what RDR 0003 reads downstream. The split is applied when a `resolve.Row` is constructed for the kernel." — (NC) — *(was blocked; `Row.Guard []GuardAtom` now exists.)*

- [REQ-85] `0002:C16` "the handoff MUST route each atom to exactly one of them by its block: every atom carrying block `match` populates `Match`; every atom carrying block `all` or `unless` belongs to the guard predicate, **regardless of operator** — an `eq` atom under `guard.all` is a guard atom. The routing is exhaustive and disjoint by construction" — (NC)

- [REQ-86] `0002:C16` "a **match block admits only `eq` and `in`** (`in` by expansion into per-member `eq` rows): a comparison, existence, `contains`, or any other operator under `[rule.match]` or `[context.<id>.match]` is refused at load as a malformed predicate atom" — (NC)

- [REQ-87] `0002:C16` "**The admitted operator set, guard blocks included, is RDR 0003's closed set `{eq, in, lt, lte, gt, gte, exists, contains}`.** This RDR does not mint operators and does not widen that set … A token under `[rule.guard.all.<tag>]` or `[rule.guard.unless.<tag>]` outside the set is a `malformed predicate atom` (`unknown operator`) refused at load" — (NC)

- [REQ-88] `0002:C17` "**Atom-level validation is block-agnostic.** Every rule this RDR states about an authored atom — operator membership (above), the `<clear>` reserved value, domain and kind conformance of a literal, `#` reservation in members, and tag-key declaration — MUST be enforced identically in all three atom blocks: `match`, `guard.all`, and `guard.unless`." — (NC)

- [REQ-89] `0002:C17` "**Conformance is per operator, not per literal.** … `eq`, `in`, and `contains` take members and are domain-checked; `exists` takes a bool literal (`<true>`/`<false>`), never a member; `lt`/`lte`/`gt`/`gte` take an ordered **bound**, which is kind-checked but not domain-checked" — (NC)

- [REQ-90] `0002:C17` "The `<clear>` ban and the tag-key declaration rule bind every atom regardless of operator. Two rules are legitimately match-only and are **not** block-agnostic: the `eq`/`in` operator restriction (§D6 routing) and the `#` reservation in `in` members, since only match blocks expand." — (NC)

## N. Normalization output and ordering

- [REQ-91] `0002:C18` "The tool MUST normalize the sparse source into deterministic candidate rows for lint, resolver lookup, diagnostics, and table dumps. Each candidate row MUST retain its source rule id and source locator." — (NC)

- [REQ-92] `0002:TD` "The locator must identify at least the model id and rule id, and it is **derived, not authored**: the loader composes it from `(model id, rule id)` … An authored `[[rule]].source` is a provenance *annotation* carried alongside it, not the locator itself" — (TD)

- [REQ-93] **§D13 LANDING (owed; not in the record).** JDR 0001 §D13: "`Tag.Value` stays `string`; a set crosses as its canonical JSON array — members sorted, duplicate-free, compact encoding — and 0002 declares it." A set-valued tag value crossing the kernel seam (`resolve.Tag.Value`, `Row.NextTags`, `Row.Writes`, and any set-valued atom literal handed to the guard seam) MUST be encoded as that canonical JSON array; read-back equality (RDR 0004) is byte equality over it. — (JDR 0001 §D13 / `deviations.md` D5) — *This is the only REQ whose text is not in the RDR record; RDR 0007 already wrote its `contains` leg against it.*

- [REQ-94] `0002:C19` "The expanded table dump MUST be derived from the normalized candidate-row value and MUST carry every field of it: row identity, source locator, outcome, predicate atoms (with block), gate ids, next-state tags, writes including `<clear>` entries, required-owned keys, and escape failure classes — plus the derived row kind column." — (NC)

- [REQ-95] `0002:C19` "`[dump]` carries exactly one key, `order`: a list of column identifiers." with the closed vocabulary "identity, source, kind, outcome, atoms, next, writes, requires_owned, gate, escape" — (NC)

- [REQ-96] `0002:C19` "`[dump]` settings MAY reorder the rendered columns; they MUST NOT omit a field. An `order` naming an unknown identifier, repeating one, or omitting one is a `malformed dump declaration` refused at load — not a silently truncated dump." — (NC)

- [REQ-97] `0002:C19` "`[dump]` is **presentation, and MUST NOT reach the normalized value** … Normalization MUST ignore it entirely, so two models differing only in `[dump]` normalize to identical candidate-row sets" — (NC)

- [REQ-98] `0002:C19` "Dump ordering MUST be deterministic across source key order. Rows sort by row identity, compared field by field as the tuple `(model id, rule id, expansion suffix)`; `model id` and `rule id` compare byte-lexicographically on the post-parse string, and the expansion suffix compares as a sequence — element by element, byte-lexicographically, a shorter sequence that is a prefix of a longer one sorting first, so the empty suffix sorts before any non-empty one." — (NC)

- [REQ-99] `0002:C19` "Within each row, atoms sort by (key, block, operator token, literal), and next-state tags, writes, required-owned keys, gate ids, and escape classes sort by key." — (NC)

- [REQ-100] `0002:C19` "The **source locator MUST NOT participate in row ordering**." — (NC)

- [REQ-101] `0002:C19` "**One source document carries exactly one model** … `duplicate model id` is never decidable within a single document. It is scoped to a caller that loads **several documents into one dump or lint invocation** and MUST be checked there, over the set of loaded models, before rows are merged for rendering. A loader handed one document MUST NOT report it." — (NC)

- [REQ-102] `0002:C19` "The dump MUST be emitted from a pre-sorted sequence of normalized rows, never by iterating a map and never by delegating key order to an encoder." — (NC)

- [REQ-103] `0002:C20` "Source order and rendered-row order MUST NOT decide a successful transition. … the normalized row set is unordered as far as selection is concerned … and normalization MUST NOT emit any positional field a consumer could tiebreak on." — (NC)

- [REQ-104] `0002:C20` "The selection procedure itself — gate-then-count, and which refusals are escapable — is **JDR 0001 §D2's and the kernel's, and MUST NOT be restated here**" — (NC)

- [REQ-105] `0002:C19` "A dump spanning more than one model MUST carry a model-unique `model id`" — (NC)

## O. Load category taxonomy

- [REQ-106] `0002:C24` "Load-time validation failures MUST retain stable data-level categories before CLI mapping, including at minimum malformed TOML, unknown schema field, missing recognized outcome alphabet, malformed recognized outcome alphabet (an empty alphabet, a duplicate member, the empty string as a member, or a member containing the expansion-suffix separator — each distinguishable), unknown tag, unknown context, cyclic context inheritance, write to non-owned tag, unknown accessor, malformed accessor declaration, malformed accessor binding, malformed tag declaration, malformed initial declaration, reserved tag value, unsupported version, malformed predicate atom (unknown operator; literal ill-formed for its operator; literal outside the declared domain; operator not admitted under a match block; `#` in a match-block `in` member — each distinguishable), malformed escape declaration (a write block, clear list, or `gate` list on an escape rule), **malformed model declaration** (an absent `[model]`, an absent `id`, an absent `version`, or an absent `[tags.recognized]` declaration …), **malformed dump declaration** …, malformed outcome binding (including a `recognized` atom authored in a guard block), **malformed rule shape** (an ordinary rule carrying no write block …), malformed rule id (containing the expansion-suffix separator), duplicate rule id, duplicate model id, and `reserved_tag_key` (RDR 0008)." — (NC)

- [REQ-107] `0002:C24` "an absent `version` MUST refuse here and MUST NOT be folded into `unsupported version`, which would name a version the author never wrote" — (NC)

- [REQ-108] `0002:C24` "**"Load" names the whole source-to-candidate-rows pipeline, not one callable.** Every category above MUST be refused by that pipeline before it yields candidate rows, whatever internal decomposition an implementation chooses." — (NC)

- [REQ-109] `0002:C24` "The categories are **data-level and MUST NOT depend on the CLI envelope**: this package MUST NOT import `internal/cli`" — (NC)

- [REQ-110] `0002:C24` "Category identifiers are **snake_case renderings of the prose names above**, anchored on `reserved_tag_key` … The Testing Strategy asserts on the category, so these identifiers are an API surface, not message text." — (NC)

- [REQ-111] `0002:C24` "Cross-row findings — overlap, gap, dead row, read-before-write — carry RDR 0006's lint categories, not these." — (NC)

- [REQ-112] `0002:TD` "a check decidable from one rule plus the model's declarations is load-time and this RDR's; a check that must compare normalized rows against each other — overlap, gap, dead row, read-before-write — is graph lint and RDR 0006's. Ambiguous overlap between candidate rows is therefore a lint finding, not a load failure." — (TD)

## P. Round-trip invariant

- [REQ-113] `0002:RT` "`parse ∘ normalize = expanded-table value identity` on valid transition model fixtures: two authorings of one model — differing in TOML key order, rule declaration order, or `eq`-versus-single-member-`in` spelling — must normalize to the same candidate-row set, compared as values over the full field list the dump contract carries … with the source locator's optional line/column detail excluded from the comparison." — (RT)

- [REQ-114] `0002:RT` "The locator's **presence and its rule-identifying part are compared**, even though its line/column detail is not" — (RT)

- [REQ-115] `0002:RT` "a set-valued field MUST render its members delimited in a way that shows the field is a set (for example bracketed), so a reviewer reading two rows can see that a difference *may* be hiding" — (RT)

- [REQ-116] `0002:RT` "Source rewrite is likewise out of scope for this RDR" — (RT)

## Q. MVV and testing obligations

- [REQ-MVV] `0002:MVV` "Parse two hand-authored sparse TOML fixtures, one for a representative RDR flow slice and one for a representative kata flow slice, into typed source data; normalize them into candidate rows; dump the expanded table; validate tag declarations, context references, predicate references, recognized outcomes, supported model version, multi-tag writes, outcome binding, and escape rows; then prove by unit test that one sample tag-set resolves through `internal/resolve::Resolve` to exactly one ordinary row, one unmatched sample tag-set resolves to exactly one modeled escape row when the table declares one, one unsupported-version variant is refused before normalization, and one deliberately overlapping variant is reported by lint as an ambiguous overlap." — (MVV)

- [REQ-117] `0002:MVV` "The overlap item is the one MVV assertion this RDR cannot discharge alone … It is **deferred to the Phase 5 lint handshake rather than counted as satisfied here**" / "Marking the overlap assertion green before RDR 0006 lands would be asserting on a stub." — (MVV)

- [REQ-118] `0002:MVV` "`evidence/spikes/iter-2/` is the promoted set" / "a fixture may not be narrowed on promotion, only extended." — (MVV)

- [REQ-119] `0002:MVV` "Each load-time validation category is asserted by **category**, not by message text, and each has one mutated fixture that trips it and no other" — (MVV)

- [REQ-120] `0002:MVV` "The existence-atom test exercises `exists = true` **and** `exists = false`" — (MVV)

- [REQ-121] `0002:MVV` "The byte-equality contract is asserted by a negative control — a reference spelling a declared `status` as `Status` must fail `unknown tag`" — (MVV)

- [REQ-122] `0002:MVV` "assert that the normalized value is unchanged when the same model is authored with its TOML keys in a different order, **and when its rules are declared in a different order**, and that repeated dumps of one model are byte-identical." — (MVV)

- [REQ-123] `0002:TS` "**`[model.metadata]` is asserted by deep equality against the authored table** — every value and every nesting level" and "The `source` and `description` keys decode onto the rule and model structs rather than refusing as unknown schema fields." — (TS 1)

- [REQ-124] `0002:TS` "This SHA MUST NOT be asserted as a golden hash by any implementation test" / "Assert over the normalized value" — (TS 2)

- [REQ-125] `0002:TS` "**The control is a mutation test**: normalize a row, mutate one field in place, and assert the other is unchanged." — (TS 2, no-alias)

- [REQ-126] `0002:TS` "Two contexts contributing `in = ["a,b", "c"]` and `in = ["a", "b,c"]` on one key and block MUST yield **two** atoms in the normalized set" and "the control instead **parameterizes over the plausible delimiters** (`,` `;` `|` space, and the empty string)" — (TS 2)

- [REQ-127] `0002:TS` "A `set`-kind write authoring `["a,b", "c"]` and one authoring `["a", "b,c"]` MUST normalize to two distinct write values, parameterized over the same delimiter set." / "The assertion must read the normalized value, not the rendered one" — (TS 2)

- [REQ-128] `0002:TS` "a rule and an inherited context contributing a **byte-identical** `(block, key, operator, literal)` atom collapse to **exactly one** atom, asserted by count on that key." — (TS 2)

- [REQ-129] `0002:TS` "the fixture's `draft-no-match-escape` binds its outcome with `in` over two alphabet members, and normalization must yield **two** escape rows carrying distinct expansion suffixes, each with an empty write set and each retaining the modeled failure-class list." — (TS 2, A11)

- [REQ-130] `0002:TS` "constructing the `resolve.Row` for each normalized RDR-fixture row, the union of `Match` and the guard's atoms equals the normalized atom set and their intersection is empty, and every `Match` atom carries block `match`." — (TS 2, A12) — *(was blocked; satisfiable now.)*

- [REQ-131] `0002:TS` "Each variant is refused with the **one** category its mutation targets and no other — the assertion is on the category, not the message text" — (TS 3)

- [REQ-132] `0002:TS` "The **seven** still owed at implementation are `malformed TOML`, `missing recognized outcome alphabet`, `cyclic context inheritance`, `malformed tag declaration` …, `duplicate model id` …, `malformed model declaration` … and `malformed dump declaration`" — (TS 3)

- [REQ-133] `0002:TS` "**Guard-block controls are owed for every atom-level rule (A15).** … Each of operator membership, `<clear>`, domain/kind conformance, `#` reservation, and tag-key declaration owes a `guard.all` and a `guard.unless` variant." — (TS 3) — Stage 6 closed six of these in the spike; scenario 3 "continues to owe one negative control per rule **per block** for any rule added later" (`0002:C17`).

- [REQ-134] `0002:TS` "**Fixture conformance is a control, not a review step.** Both normative fixtures MUST load clean under the finished loader, asserted as part of scenario 1." — (TS 3)

- [REQ-135] `0002:TS` "The precedence control is a **v2-shaped** document — `version = 2` plus a v2-only key the strict decoder would reject as an unknown schema field — asserted to refuse `unsupported version` and **not** `unknown schema field`." — (TS 3, A14)

- [REQ-136] `0002:TS` "Positive controls in the same fixture set: an observed key served by no reader loads (JD-9), and a model with no `[initial]` loads and is left to lint." — (TS 3)

- [REQ-137] `0002:TS` "Run `internal/resolve::Resolve` over one matching ordinary tag-set …, one tag-set with no ordinary match but one matching `no_match` escape row, and one tag-set in which a sibling candidate's guard is unevaluable … **Expected**: … the unevaluable-sibling tag-set refuses `guard_unevaluable` even though a decidable sibling and a `no_match` escape row exist" — (TS 4) — *(was blocked; satisfiable now.)*

- [REQ-138] `0002:TS` "The non-boolean literal is refused at load as a malformed predicate atom; the mis-cased reference fails `unknown tag` rather than folding." — (TS 6)

- [REQ-139] `0002:CA` (A10) "the scenario 3 fixture for this category is a **pair** of documents sharing a `model id`, and a single-document load of either one must not report it." — (CA A10 / TS 3)

## R. Package placement and phases

- [REQ-140] `0002:EIA` "New internal package can own sparse source structs and normalization." — (EIA, Normalizer/table package)

- [REQ-141] `0002:IP` Phase 2 "Introduce typed source structures, normalization to `internal/resolve::Row` values (outcome lifted, per-atom block retained, `RequiresOwned` derived, existence constants emitted), and an expanded-table dump with deterministic ordering and source rule ids." — (IP Phase 2)

- [REQ-142] `0002:IP` Phase 3 "Add validation rules for tag and accessor declarations, the provenance-scoped `keys` bindings, the reserved `recognized` name and `<clear>` value, rule ids, context references, `[initial]` assignments and `terminal` ids, predicate atoms (including the match-block operator restriction and domain membership), outcome bindings, gate references, write targets, explicit clears, escape declarations, and expansion-count diagnostics." — (IP Phase 3)

- [REQ-143] `0002:IP` Phase 4 "Connect the normalized rows to RDR 0001's gate-then-count exact-one contract without adding runtime ordering or host-code predicate callbacks." — (IP Phase 4)

- [REQ-144] `0002:IP` Phase 5 "Expose the parsed representation needed by RDR 0006 for graph determinism, reachability, and read-before-write checks." — (IP Phase 5)

- [REQ-145] `0002:MC` "expansion counts per rule | dump (derived: rows per source rule id) | this RDR | reviewer-visible, non-blocking; no threshold here." — (MC `disposition`) — Phase 3's "expansion-count diagnostics" are **derived from the dump**, not a lint finding and not a refusal (`0002:FM`).

- [REQ-146] `0002:PE` "source key order is neutralized by sorting every emitted sequence; Go's map iteration is never the emission order (rows, atoms, next tags, writes, required-owned keys, and escape classes are all sorted slices before rendering)" — (PE)

---

## ASSUMPTIONS

- ASSUMPTION: **RDR 0007's kernel reshape is treated as landed, and Phases 2–3
  open without a local mirror of the constants.** The record's Prerequisites
  bullet (`0002:2049-2077`) states the reshape is unimplemented, has no owner,
  and that Phases 2–3 "sequence behind" it, forbidding a local mirror. Grounded
  against source: `internal/resolve/guard.go` exports `GuardAtom`, `OpExists`,
  `LiteralTrue`, `LiteralFalse`, `BlockAll`, `BlockUnless`, `BlockMatch`, and
  `resolve.go::Row.Guard` is `[]GuardAtom`; `docs/rdr/BUILD-ORDER.md` puts 0007
  at run 1 of 8 and `0007/artifacts/status.md` records it COMPLETE. The bullet
  is stale-by-design (a Final RDR describing a `main` that has since moved), so
  the sequencing block is discharged rather than deviated. Recorded at
  `deviations.md` D4.

- ASSUMPTION: **The §D13 set-value encoding clause is written as REQ-93 from
  JDR 0001 §D13 verbatim, not re-derived.** Grounded: `deviations.md` D5 calls
  it "a *landing*, not a citation, and it is the only one in this sweep";
  `BUILD-ORDER.md` §"The one seam this order strains" says RDR 0007 already
  wrote its `contains` leg and TS scenario 8 "against §D13 directly" and that
  "0002's later run must match it, not redefine it". So this build's job is to
  implement the JSON-array form, not to choose one.

- ASSUMPTION: **REQ-38 (member-sequence write value) and REQ-93 (canonical JSON
  array at the seam) are two layers, not a contradiction.** The normalized value
  holds a `[]string` member sequence (REQ-38, REQ-56, REQ-58); the *kernel seam*
  encoding — what lands in `resolve.Tag.Value`, a `string` — is the canonical
  JSON array (REQ-93). §D13 says so explicitly: "The member-sequence fence is
  satisfied: a JSON array is a sequence, not a joined string." So the handoff
  serializes; the normalized value never does, and no identity is derived from
  the serialized form.

- ASSUMPTION: **`kind = "string"` in the promoted fixtures is rewritten to RDR
  0003's five-token vocabulary before promotion.** Grounded: `deviations.md` D1
  check 1 requires `grep 'kind = "string"'` → 0 and names `scalar` as 0003's
  token (`0003:791-793`). The fixture is evidence, not a fence, so this is a
  fixture rewrite, not a spec change. REQ-64 carries the pointer.

- ASSUMPTION: **`write = []` on an escape rule is refused by key *presence*, not
  by length.** REQ-36's fence says "even an empty one"; `deviations.md` D2
  records the spike checking `len(rule.Write) > 0` and requires a
  `neg-escape-with-empty-write` fixture. The loader keys on presence.

- ASSUMPTION: **`duplicate model id` (REQ-101, REQ-139) is implemented as a
  multi-document entry point, and the single-document loader does not carry the
  check at all.** Grounded: `0002:C19` "A loader handed one document MUST NOT
  report it", and A10's decidability half is source-verified (`[model]` is a
  singular struct field). The Phase 2/3 API therefore needs a second, set-taking
  call — a check over loaded models, "before rows are merged for rendering".

- ASSUMPTION: **Load-time domain/kind conformance of *write-block values* is
  filed under `malformed tag declaration` / literal-outside-domain, per
  `deviations.md` D3**, since REQ-106's list names no dedicated category and D3's
  check names that one. If a fenced category must widen, D3 says escalate.

- ASSUMPTION: **Snake_case category identifiers are derived mechanically from
  REQ-106's prose names** (`malformed_toml`, `unknown_schema_field`,
  `malformed_predicate_atom`, …), anchored on `reserved_tag_key`, per REQ-110.
  The one place the record forbids re-derivation is the `[dump]` column
  vocabulary (REQ-95), which is fixed verbatim to the fixtures' spelling.

- ASSUMPTION: **The Round-Trip invariant (REQ-113) is asserted over the
  normalized value, never over dump text.** Stated by REQ-113 itself and
  reinforced by REQ-97/REQ-124; `dump` has no specified inverse (`0002:RT`).

- ASSUMPTION: **The `[model.metadata]` field on the normalized model is
  `map[string]any`.** Grounded: A1's evidence names exactly this — "the strict
  check descends only into struct branches, leaving a `map[string]any` field
  untouched, which is what makes `[model.metadata]` free-form while the rest
  stays strict" — and REQ-12 forbids constraining its shape.

- ASSUMPTION: **Phase 5's lint handshake (REQ-144) mints an exported accessor
  over the normalized model, not an RDR 0006 implementation.** RDR 0006 is
  unimplemented and run 4 of 8 (`BUILD-ORDER.md`); REQ-117 forbids marking the
  overlap assertion green against a stub. The deliverable is the exposed shape.

- ASSUMPTION: **REQ-145's "expansion-count diagnostics" in Phase 3 is a derived
  dump property, not a validation rule.** Grounded: the `disposition` table
  routes it to "dump (derived: rows per source rule id)", Failure Modes says it
  is "not a lint diagnostic and not a refusal", and Risks says it is "derivable
  from the dump itself … needs no lint output". So Phase 3 mints no check for it.

- ASSUMPTION: **The kernel `Row` fields REQ-80/REQ-82 name are `NextTags` and
  `Writes`, both `[]Tag`** — confirmed at source (`internal/resolve/resolve.go`),
  which is what makes the no-alias mutation test (REQ-125) assertable, exactly as
  TS scenario 2 argues.

---

## QUESTIONS

- **Q1 — Does `Match []Tag` survive the 0007 reshape, or do match atoms travel
  in `Row.Guard` as `BlockMatch` atoms?** Two readings, materially different
  for REQ-85 (the handoff routing) and REQ-130 (the totality/disjointness
  assertion).
  (a) `Row.Match []Tag` stays the equality pattern and `Row.Guard []GuardAtom`
  takes only `all`/`unless` atoms; `BlockMatch` exists on the shared `Block`
  type but no match atom is ever placed in `Guard`. REQ-85 then routes to two
  different Go fields of different types, and REQ-130's "union … equals the
  normalized atom set" is asserted across `Match`-as-`[]Tag` and
  `Guard`-as-`[]GuardAtom`.
  (b) Match atoms travel in `Row.Guard` discriminated by `Block == BlockMatch`,
  making the union a single-slice partition.
  **Proceeding under (a).** The evidence is one-sided but not conclusive:
  `internal/resolve/resolve.go` still declares `Match []Tag` alongside
  `Guard []GuardAtom`; `guard.go::isGuardBlock` returns false for `BlockMatch`
  and its comment says the kernel's guard evaluation excludes it; `0002:C16`
  says "every atom carrying block `match` populates `Match`", naming `Match` as
  a distinct destination; and `0007/artifacts/req-list.md` Q1 asks the mirror
  question and also proceeds under (a). This is recorded rather than resolved
  only because 0007's own auditor left it open and this RDR is the first
  producer that actually fills both fields — if a later phase finds `Match`
  cannot carry a `BlockMatch`-tagged atom's operator token, the routing contract
  is unchanged but the assertion shape in REQ-130 is.

- **Q2 — Which category does an `[initial]` value that is well-formed for its
  kind but outside its declared domain carry?** REQ-30 folds "kind and domain"
  into `malformed initial declaration`, while REQ-106 lists "literal outside the
  declared domain" as a distinguishable arm of `malformed predicate atom`, and
  TS scenario 3 names `neg-initial-bad-value` / `neg-initial-int-out-of-range`
  under the `malformed initial declaration` arm.
  **Proceeding on `malformed initial declaration`** for `[initial]` values and
  `malformed predicate atom` for predicate literals — the two clauses scope by
  *site*, not by defect, and the fixture names in TS 3 confirm the `[initial]`
  site. Recorded because REQ-106's "literal outside the declared domain" arm
  reads unqualified, and a category-not-message oracle (REQ-131) will fail on
  the wrong choice. No predecessor settles it: RDR 0003 owns the domain rules
  but not this RDR's category assignment.
