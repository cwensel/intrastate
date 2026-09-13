# REQ List — RDR 0021 Lint's normalized-graph export

Phase 0 audit artifact. Every clause below is a testable obligation drawn from
`docs/rdr/0021-lint-normalized-graph-export.md`. Quotes are verbatim — copied
from the projector (`rdr inspect --select <id>`) for fenced elements and read
from the record for testable prose outside the fences — never transcribed by
hand.

Element ids (`0021:C1` … `0021:C5`, `0021:MVV`, `0021:S1` … `0021:S9`,
`0021:D-*`, `0021:RT*`, `0021:F*`, `0021:G-cross-cutting`) are carried wherever
a REQ derives from a labelled element, so a later stage can trace the REQ back
to its contract.

`counts.elements.C = 5`, and the five fences carry far more than five
obligations: C1 is the whole verb surface (arm set, flag order, tier), C2 the
whole document (field list, orders, empty/absent rule, marker, stability), C4
the neutrality rule plus the ceiling refusal. Testable prose also lives outside
the fences — `Load-Bearing Decisions`, `Round-Trip / Inverse Invariants`, the
`Pre-Lock Mini-Checks` disposition/authority tables, `Failure Modes`, the
`Implementation Plan` phases, `Testing Strategy`, and `Cross-Cutting Concerns`
— and each is mined below.

Source sections are abbreviated:

- `NC` = Proposed Solution / Technical Design / Normative Contracts (fenced)
- `AP` = Proposed Solution / Approach
- `TD` = Proposed Solution / Technical Design (unfenced prose)
- `LBD` = Technical Design / Load-Bearing Decisions
- `RT` = Technical Design / Round-Trip / Inverse Invariants
- `MC` = Technical Design / Pre-Lock Mini-Checks (the five cue-fired tables)
- `CAP` = Proposed Solution / Capability Dependencies
- `EIA` = Proposed Solution / Existing Infrastructure Audit
- `CA` = Critical Assumptions
- `FM` = Trade-offs / Failure Modes
- `PRE` = Implementation Plan / Prerequisites
- `IP` = Implementation Plan (phases)
- `MVV` = Implementation Plan / Minimum Viable Validation
- `TS` = Validation / Testing Strategy (numbered scenarios)
- `XC` = Finalization Gate / Cross-Cutting Concerns

---

## Verb surface and registration

- [REQ-1] "A new root verb `graph` is registered beside `lint` — outside the `flow` group, under `0005:C1`'s carve-out for command groups \"owned by the RDR that names them\"." — (NC, `0021:C1`) — registration is at root via `NewRootCmd`'s `cmd.AddCommand`, the shipped siblings being `newVersionCmd`/`newLintCmd`/`newFlowCmd`/`newDocsCmd` (`internal/cli/root.go:192-200`).

- [REQ-2] "Model selection mirrors `lint`'s arm set — mirroring is scoped to the ARM SET and its CODE SPELLINGS, not to lint's help/usage text, which each verb words for itself — and codes verbatim (`internal/cli/lint.go::runLint`)" — (NC, `0021:C1`).

- [REQ-3] "exactly one of `--model <path>` / `--flow <id>`; both → `flag-mutually-exclusive`" — (NC, `0021:C1`) — the shipped spelling is at `internal/cli/lint.go:153`.

- [REQ-4] "`--flow` alone → `flag-invalid-value` (this build resolves no ids)" — (NC, `0021:C1`) — shipped spelling at `internal/cli/lint.go:164`.

- [REQ-5] "neither → `flag-required`" — (NC, `0021:C1`) — shipped spelling at `internal/cli/lint.go:173`.

- [REQ-6] "unreadable file → `model-unreadable`" — (NC, `0021:C1`) — shipped spelling at `internal/cli/lint.go:183`.

- [REQ-7] "load failure → `model-invalid` with one findings[] entry per load category." — (NC, `0021:C1`) — shipped spelling at `internal/cli/lint.go:202`; the entries come from the shared `internal/cli/flow_input.go::loadFindings` (`:223`), which populates each entry's `code` with the load-category slug.

- [REQ-8] "`--emit <format>` selects the document: `json` (default) or `dot`; any other value → `flag-invalid-value` naming `emit`." — (NC, `0021:C1`).

- [REQ-9] "That format set is an `append-only` vocabulary (`0029:C2`): a third format MAY be added in a minor, none removed or renamed within a major, so a consumer MUST NOT assert its cardinality or a member's position." — (NC, `0021:C1`).

- [REQ-10] "Its enumeration seam, which `0029:C4` obliges of every tier assignment, is a new exported `internal/cli` accessor over the two members — the package owns the set and nothing it imports imports it back, so the accessor builds with no new import and no cycle." — (NC, `0021:C1`).

- [REQ-11] "`--emit` validity is checked with the argument-shaped arms, AFTER the `--model`/`--flow` selection arms and BEFORE any file I/O. So `graph --model <unreadable> --emit=xml` refuses `flag-invalid-value` naming `emit`, never `model-unreadable`; the request is wrong independent of the environment." — (NC, `0021:C1`) — the ordering is normative and directly assertable.

- [REQ-12] "This is the one value-checked flag on the verb" — (NC, `0021:C1`) — no other flag on `graph` performs an enum-value check.

- [REQ-13] "RunE starts with `respond.ValidateMode`, success routes through `respond.OK`, failure through `respond.Fail`" — (NC, `0021:C1`) — the shipped seams are `internal/cli/respond/respond.go:142` (`ValidateMode`), `:162` (`OK`), `:212` (`Fail`).

- [REQ-14] "exit codes are the existing 0/2 mapping — no new exit group (`0005:C1`: exit 3 is environment-not-consulted only, and an unreadable model file is a wrong request, so it stays exit 2)." — (NC, `0021:C1`).

- [REQ-15] "verb `graph`, flag `--emit`." — (LBD, `0021:D-naming`) — the exact spellings are normative; `dump`, `export`, and overloading `--as` are the named rejected alternatives in the same decision.

## The exported document (JSON)

- [REQ-16] "The JSON document is a documentation artifact, versioned by a required leading `schema` field, initial value `intrastate.graph/1`" — (NC, `0021:C2`) — `schema` is both REQUIRED and LEADING (first on the wire), which `XC`'s "schema first" restates.

- [REQ-17] "evolution within `/1` is additive (a consumer ignoring unknown fields keeps working), and an incompatible change bumps the marker under the `0.x` promise RDR 0029 governs." — (NC, `0021:C2`).

- [REQ-18] "Exact field spellings are normative AS SPELLED HERE — the Illustrative Code is an exhibit that disclaims literal assertion and shows only some members, so it binds nothing" — (NC, `0021:C2`) — the field list in REQ-19…REQ-25 is the normative surface; the `Illustrative Code` JSON block is explicitly NOT assertable.

- [REQ-19] "`schema`; `model` (`table.Model.ID`, the AUTHORED `[model] id`, never the `--model <path>` argument or any path-derived string); `class`;" — (NC, `0021:C2`) — `internal/table/model.go:497` carries `ID`, `:506` carries `Class`.

- [REQ-20] "`tags[{name, provenance, kind, required, single_valued, domain}]` (`domain` present exactly when `guard.AssignmentCount` reports it finite)" — (NC, `0021:C2`) — `domain`'s presence is conditional and exactly determined, so both arms (finite → present, non-finite → absent) are assertable.

- [REQ-21] "`initial`; `terminal` (the declared predicate sets, carried as declared);" — (NC, `0021:C2`).

- [REQ-22] "`rows[{identity, source, kind, outcome, atoms, next, writes, requires_owned, gate, escape, emit}]` — RDR 0002's dump field list in the row order `internal/table/dump.go::dumpColumns` fixes, atoms (each `{key, operator, literal[], block}`) in the atom order 0002 fixes, lowercase snake_case" — (NC, `0021:C2`) — `dumpColumns` is the closed column vocabulary at `internal/table/dump.go:17`.

- [REQ-23] "the `emit` member is the one `0010:C3` adds, and this document carries the dump's field vocabulary rather than restating a cardinality" — (NC, `0021:C2`).

- [REQ-24] "`groups[{context, rules}]`;" — (NC, `0021:C2`).

- [REQ-25] "`reach{abstraction, nodes[{id, values}], edges[{from, to, rule}]}` — the merged fixpoint relation, `id` the node key `reach.go::(Node).key` fixes, `values` an OBJECT keyed by tag name whose every value is that tag's sorted, deduplicated value ARRAY (the shape `reach.go::Node`'s `Values map[string][]string` projects without invention)" — (NC, `0021:C2`; A7) — `values` is an object, NOT an array of joined `key=value` strings.

- [REQ-26] "nodes sorted by node key and edges by (from, to, rule) so construction order is unobservable." — (NC, `0021:C2`).

- [REQ-27] "A tag with no finite declared domain carries the single value `<opaque>` (`reach.go::OpaqueValue`) in its array, passed through verbatim; this record attaches no meaning to that spelling." — (NC, `0021:C2`) — `internal/graphlint/reach.go:28` declares `OpaqueValue = "<opaque>"`.

- [REQ-28] "Every declared collection renders as an empty JSON array `[]` (or object `{}`) when it has no members — never `null`; an optional member that does not apply is ABSENT, its key omitted, never `null`." — (NC, `0021:C2`; MC `disposition` rows 8–9) — three states are distinguishable on the wire and each is separately assertable.

- [REQ-29] "The `reach` block carries a REQUIRED marker, `abstraction` with the token value `declared-over-approximation`, stating the relation is the DECLARED over-approximation, not the runtime — merged nodes, guard/observed atoms unpruned" — (NC, `0021:C2`).

- [REQ-30] "the schema docs state the soundness rule in one sentence: universal claims (\"no path does X\") proved over this relation hold at runtime; existence claims (\"some path reaches X\") may be spurious." — (NC, `0021:C2`) — a documentation obligation with an observable artifact (the schema docs), landed by IP Phase 4 (REQ-66).

- [REQ-31] "The document carries NO verdict or finding field — an export is never a lint pass — and NO per-node terminal marking" — (NC, `0021:C2`) — a negative with a direct observable: no such key on the wire, in any arm.

- [REQ-32] "the declared sets travel in the document for a consumer to evaluate." — (NC, `0021:C2`) — the `terminal` member carries the declared predicate sets even though no node is marked.

- [REQ-33] "Set-valued members are JSON arrays, closing `0002:§round-trip-inverse-invariants`'s lossy set-literal rendering for this document" — (NC, `0021:C2`).

- [REQ-34] "The field spellings above are an `append-only` vocabulary in `0029:C2`'s sense — a member MAY be added in a minor, none removed or renamed within a major" — (NC, `0021:C2`, STABILITY).

- [REQ-35] "A consumer MUST therefore tolerate an unrecognized field and MUST NOT assert on the field set's cardinality, a member's ordinal position, or a tail position; the `emit` member's arrival last in `dumpColumns` is 0002's row order, not a position this document promises." — (NC, `0021:C2`, STABILITY) — binds this record's OWN tests: they may not assert field-set size or member position.

- [REQ-36] "The enumeration seam `0029:C4` obliges is the Go struct's json tags, as it is for the `version` payload field names — the compiler is the seam, so the document's own decode is the by-value assertion." — (NC, `0021:C2`, STABILITY).

- [REQ-37] "Two markers coexist at different levels and neither substitutes for the other: this `schema` field versions the DOCUMENT, while `0029:C1`'s `schema_version` versions the ENVELOPE carrying it and is never projected into `data`." — (NC, `0021:C2`, STABILITY) — `schema_version` MUST NOT appear inside `data`.

- [REQ-38] "C2's schema-versioned JSON document; exact field spellings normative per C2's list; empty declared collections render `[]`/`{}` and an inapplicable optional member is absent, never `null`; DOT styling explicitly non-normative." — (LBD, `0021:D-wire-byte-format`).

## The exported document (DOT)

- [REQ-39] "The DOT document renders the same value: one node per reachability node, one edge per reachability edge labeled with its rule id, the initial node marked, and the abstraction marker rendered in the graph header comment/label so the diagram carries it too" — (NC, `0021:C2`).

- [REQ-40] "its node/edge SET and the marker are normative, its styling/attributes are not." — (NC, `0021:C2`) — fixes the assertion boundary: set equality and marker placement are assertable, styling is not.

- [REQ-41] "A pure renderer over the export value; equality test against the JSON arm's node/edge set." — (IP, Phase 3; A6) — the renderer consults no `graphlint`/analysis internals.

## Determinism

- [REQ-42] "For one model input and one build, emission is byte-for-byte identical across invocations, in every `--emit` and `--as` combination" — (NC, `0021:C3`) — all four cells, byte-level.

- [REQ-43] "every sequence on the wire is pre-sorted by C2's orders before marshaling, no Go map iteration reaches the wire, and JSON is emitted through the one shared non-HTML-escaping encoder (`clierr.WriteJSONLine`, per `0005:C1`'s one-encoder rule)." — (NC, `0021:C3`) — `internal/cli/clierr/clierr.go:220`; A5 records four ad hoc `SetEscapeHTML(false)` call sites that MUST NOT be joined by a fifth.

- [REQ-44] "The promise is scoped to (model, build), not across builds" — (NC, `0021:C3`) — the narrowing is normative; no cross-build byte promise may be asserted.

- [REQ-45] "a build bump may change the JSON only by C2's additive rule and may re-baseline DOT diffs freely, since DOT styling is non-normative." — (NC, `0021:C3`).

- [REQ-46] "Identity is therefore invocation-independent: the same model exported from two checkouts is byte-identical" — (LBD, `0021:D-identity`) — follows from `model` being the authored `[model] id`, never a path-derived string.

- [REQ-47] "a node is its canonical node key (`reach.go::(Node).key` — injective by escaping, ⇒ two exported nodes never collide); a row is RDR 0002's row identity; an edge is the `(from, to, rule)` triple" — (LBD, `0021:D-identity`).

## Neutrality

- [REQ-48] "The export runs load, normalization, grouping, and the reachability traversal only, reaching the traversal through the same `graphlint` entry surface lint uses" — (NC, `0021:C4`) — structural, not disciplinary: one traversal, two call sites.

- [REQ-49] "It MUST NOT run the lint invariants, MUST NOT emit findings, and MUST NOT alter any input it shares with lint: `intrastate lint`'s verdict, finding set, and bytes are identical with and without the export code present" — (NC, `0021:C4`) — the byte-compare oracle, against a pre-change capture.

- [REQ-50] "the oracle is MECHANISM-INDEPENDENT — it binds equally if Resolve picks A2's in-traversal edge observer" — (NC, `0021:C4`) — the assertion may not be written so that it only holds for the post-hoc-recovery mechanism.

- [REQ-51] "The verb succeeds for ANY model that loads, including a model lint refuses; the model-loads-but-lint-refuses case is a dedicated fixture with an asserted, defined document — never whatever the traversal happens to do" — (NC, `0021:C4`) — the document content must be ASSERTED, not merely non-empty.

- [REQ-52] "When the traversal is incomplete under the published node ceiling (`graphlint.NodeCeiling()`), the verb refuses with the scalar code `graph-export-too-large` (GroupUserEnv, exit 2) naming the ceiling and the narrow-a-domain remedy — never a partial document, because a partial graph diffs as a graph change." — (NC, `0021:C4`) — `internal/graphlint/taxonomy.go:174` publishes `NodeCeiling()`; `:158` fixes `nodeCeiling = 4096`.

- [REQ-53] "Completeness MUST reach the verb through an EXPORTED `graphlint` surface: today's `Reach` discards the private `reach`'s `complete` bool (`nodes, _ := reach(m)`), so the edges-carrying function Q3(c) already adds beside it returns completeness with the nodes and edges." — (NC, `0021:C4`; A8) — the new function carries nodes, edges, AND completeness; `Reach` (`reach.go:73`) and `reach` (`:90`) are untouched.

- [REQ-54] "The ceiling this refusal names is `reach`'s node-count completeness — the SAME bound `analysis.go`'s `checkNodeCeiling` already fires on (both read one `reach()` and one `nodeCeiling` constant, published as `graphlint.NodeCeiling()`), observed at two call sites" — (NC, `0021:C4`).

- [REQ-55] "The distinct bound this refusal does NOT involve is the guard-product one, `graphlint.ProductBound()` (`guard.Bound()`); no guard-product input reaches the export's refusal path." — (NC, `0021:C4`) — `internal/graphlint/taxonomy.go:169`; a negative with an observable (no guard-product input on that path).

- [REQ-56] "Add a sibling exported function (`Reach` and `reach()` both untouched)" — (EIA, *Reachability* row) — the extension is additive; neither existing function's signature or behaviour changes.

## Mode coexistence (`--as` × `--emit`)

- [REQ-57] "Under `--as=text` (the default) stdout carries the selected document verbatim and nothing else — the payload satisfies the gateway's `TextLiner` seam, so the DOT stream pipes to `dot` and the JSON document diffs raw in CI; advisories stay on stderr." — (NC, `0021:C5`) — `internal/cli/respond/respond.go:92` declares `TextLiner`, `:196` is the `OK` branch that prints it.

- [REQ-58] "Under `--as=json` stdout carries exactly one terminal `ok` envelope whose `data` embeds the same document: the document object for `--emit json`, and for `--emit dot` an object carrying the DOT text as the single REQUIRED string member `dot` — spelled normatively here, so the documented unwrap is exactly `jq -r .data.dot`." — (NC, `0021:C5`) — `dot` is the normative member spelling.

- [REQ-59] "Both modes derive from one export value (0005's two-modes agreement)." — (NC, `0021:C5`).

- [REQ-60] "All four `--as`×`--emit` cells are defined — none refused, none dead" — (NC, `0021:C5`) — each of the four cells is separately assertable.

- [REQ-61] "a consumer that parses stdout as an envelope MUST use `--as=json`: the text-mode stream is the bare document by contract, exactly as `version`'s `TextLine` is its identity string" — (NC, `0021:C5`).

- [REQ-62] "The verb offers NO caller-controlled projection of its success payload; if a later revision adds one it MUST conform to JDR 0002 §D1 (projection before respond.OK, echo-group-only, enforced partition, always-keep core)." — (NC, `0021:C5`) — the present-tense half is a testable negative: no projection flag on the verb.

- [REQ-63] "`--emit` selects the DOCUMENT, `--as` selects the ENVELOPE; the 2×2 composes with no refused or dead cell" — (LBD, `0021:D-selection-predicate`).

- [REQ-64] "Reuse (`TextLiner` + `OK`)" — (EIA, *Output gateway* row) — no gateway exception and no second stream is introduced.

## Round-trip / inverse invariants

- [REQ-65] "`json-decode ∘ export = value identity on every field C2 lists`: decoding the JSON document reconstructs, value-for-value, each of those fields — including exact set members, the recoverability 0002's text dump deliberately declined" — (RT, `0021:RT1`) — the quantifier is C2's field list, not "whatever was exported".

- [REQ-66] "`export ∘ export = byte identity on any loadable model` (C3's replay form; the MVV asserts it as byte equality, not exit-code green)." — (RT, `0021:RT2`).

- [REQ-67] "No `load ∘ export` inverse is claimed: the document is derived output, never a model source." — (RT, `0021:RT3`) — a declared exemption; no test may assert an export→load round trip.

## Minimum Viable Validation

- [REQ-68] "Author a small state-machine fixture (two owned states, one terminal, one escape row) and a decision-table fixture." — (MVV, `0021:MVV` step 1).

- [REQ-69] "`intrastate graph --model <fixture>` twice → the two stdouts are BYTE-identical, parse as JSON, and carry every C2 field (schema, model identity and class, tags, initial, terminal, rows, groups, reach.nodes, reach.edges, and the `reach` abstraction marker)." — (MVV, `0021:MVV` step 2).

- [REQ-70] "`intrastate graph --model <fixture> --emit dot | dot -Tsvg` renders; the DOT node/edge id set equals step 2's `reach` block." — (MVV, `0021:MVV` step 3).

- [REQ-71] "`intrastate graph --model <fixture> --as=json | jq .data` equals step 2's document, value-for-value." — (MVV, `0021:MVV` step 4).

- [REQ-72] "`intrastate lint` over a fixture WITH blocking findings still refuses identically (byte-compared against a pre-change capture), while `intrastate graph` over the same model succeeds with an ASSERTED document (defined content, not incidental) — the neutrality oracle, mechanism-independent per C4." — (MVV, `0021:MVV` step 5).

## Testing Strategy scenarios (obligations beyond the MVV)

- [REQ-73] "Determinism under map-seed variation — emit the same fixture repeatedly across `--emit`×`--as` cells, with map order PROVOKED under `GODEBUG=randmapiter=1` … **Expected**: byte-identical stdout per cell (C3); no cell depends on Go map iteration order." — (TS, `0021:S1`) — `randmapiter=1` is the required mechanism, not an option.

- [REQ-74] "Golden fixtures pin `intrastate.graph/1` for a state-machine model and a decision-table model. The goldens capture the DOCUMENT, not the envelope — `--as=text`, or `jq .data` off `--as=json` — so `0029:C1`'s `schema_version` … never enters the pinned bytes." — (TS, `0021:S2`).

- [REQ-75] "**Expected**: the goldens hold; a field added without a schema decision fails the pin (C2's additive rule tripwire), and an envelope version bump does not." — (TS, `0021:S2`) — the golden must be discriminating in one direction and insensitive in the other.

- [REQ-76] "Lint neutrality — run `intrastate lint` over blocking and clean fixtures with the export code present, byte-compared against a pre-change capture. **Expected**: identical verdict, finding set, and bytes (C4)" — (TS, `0021:S3`) — both a blocking and a clean fixture.

- [REQ-77] "Model that loads but lint refuses. **Expected**: `graph` succeeds with an ASSERTED document (defined content, not incidental), carrying no verdict or finding field (C2, C4)." — (TS, `0021:S4`).

- [REQ-78] "Mode agreement — `--as=json | jq .data` against the `--as=text` document, for both `--emit` values. **Expected**: value-for-value equality for `json`; for `dot`, the `jq -r .data.dot` unwrap reproduces the DOT text byte-for-byte AFTER accounting for the one trailing newline F1 pins on text-mode stdout (`document + \"\\n\"`) — the `dot` string member carries the document without that gateway newline, so the comparison is against the document, not the stream (C5, A1)." — (TS, `0021:S5`) — the trailing-newline accounting is normative, not incidental.

- [REQ-79] "DOT arm equality and hostile content — tag values carrying quotes, newlines, and non-ASCII. **Expected**: DOT node/edge id set equals the JSON document's `reach` block, the abstraction marker is present in the header, and every fixture survives `dot -Tsvg`" — (TS, `0021:S6`).

- [REQ-80] "Exit 0 is NOT sufficient: each hostile node and edge identifier's emitted quoted string is unescaped by inverting A6's recorded escaping ORDER (undo the `\\n` line-break escape, then `\\\"`, then `\\\\`) and asserted EQUAL to its source id, so a well-formed but mis-escaped identifier fails — the bug class A6's spike hit. Identifiers only; DOT label STYLING stays non-normative (F3)." — (TS, `0021:S6`) — the un-escape-and-compare is the required oracle; a `dot -Tsvg` exit code alone fails this REQ.

- [REQ-81] "Traversal incomplete under the published node ceiling. **Expected**: refusal with `graph-export-too-large` (GroupUserEnv, exit 2) naming the ceiling and remedy; no document on stdout (C4)." — (TS, `0021:S7`).

- [REQ-82] "Selection-arm refusals — both/neither of `--model`/`--flow`, `--flow` alone, unreadable file, load failure, and an unknown `--emit` value. **Expected**: lint's code spellings verbatim plus `flag-invalid-value` naming `emit`; exit codes stay the existing 0/2 mapping (C1)" — (TS, `0021:S8`).

- [REQ-83] "stdout carries NO document on every refusing arm (C4's never-a-partial-document rule, asserted here and not only in the mini-check disposition table)." — (TS, `0021:S8`) — every refusing arm, not just the ceiling arm.

- [REQ-84] "JSON round-trip — decode the exported document. **Expected**: value identity on every C2 field including exact set members (RT1); no `load ∘ export` inverse is exercised (RT3)." — (TS, `0021:S9`).

- [REQ-85] "The MVV is the acceptance floor; these scenarios are what \"done\" adds beyond it. Every oracle is byte- or set-equality — none asserts an exit code alone." — (TS, preamble) — a standing constraint on every test this record authorizes.

## Disposition table (input class × exit · code · artifact)

- [REQ-86] "Model loads, lint clean | 0 | — | document on stdout | loud" — (MC, `disposition`).

- [REQ-87] "Model loads, lint would refuse | 0 | — | ASSERTED document (defined content) | loud" — (MC, `disposition`).

- [REQ-88] "Traversal incomplete at node ceiling (`graphlint.NodeCeiling()`, via the exported completeness surface C4 requires) | 2 | `graph-export-too-large` (GroupUserEnv), names ceiling + narrow-a-domain remedy | **none** — never a partial document | loud (C4)" — (MC, `disposition`).

- [REQ-89] "Both/neither `--model`/`--flow` | 2 | `flag-mutually-exclusive` / `flag-required` | none | loud (C1)" — (MC, `disposition`).

- [REQ-90] "`--flow` alone | 2 | `flag-invalid-value` (no ids resolve this build) | none | loud" — (MC, `disposition`).

- [REQ-91] "Unreadable file / load failure | 2 | `model-unreadable` / `model-invalid` (one findings[] entry per load category) | none | loud" — (MC, `disposition`).

- [REQ-92] "Unknown `--emit` value | 2 | `flag-invalid-value` naming `emit` | none | loud" — (MC, `disposition`).

## Authority (writer / reader / call-site fences)

- [REQ-93] "Reachability relation | `reach.go::reach` (private) | lint via `graphlint.Run`; export via the new sibling exporter | `analysis.go:46` (lint); the new exporter | … | `reach()` — both arms reach the one traversal (C4)" — (MC, `authority`) — lint's call site stays `analysis.go:46`, unchanged.

- [REQ-94] "Empty collection on the wire | producer code choosing `[]T{}`, never nil | … | `null` (rejected) | `[]`/`{}`; optional member ABSENT" — (MC, `authority`) — the producer constructs empty slices, never nil, so `null` cannot reach the wire.

- [REQ-95] "Selection-arm refusal codes | `lint.go::runLint` | `graph` verb mirrors the arm set and code spellings verbatim; help text is the verb's own (C1) | … | lint — C1 mirrors, never redefines" — (MC, `authority`).

- [REQ-96] "Document format selection | `--emit` (this verb only) | … | `--as` (envelope mode, `respond.go::FlagName`) | `--emit` selects DOCUMENT, `--as` selects ENVELOPE" — (MC, `authority`) — `respond.FlagName` is `"as"` (`internal/cli/respond/respond.go:52`); `--emit` is scoped to this verb alone.

## Failure modes

- [REQ-97] "refusals reuse lint's arm codes (C1) and the new `graph-export-too-large` refusal names the ceiling and remedy … a killed process leaves no terminal envelope (existing contract)." — (FM, `0021:F1`).

- [REQ-98] "an envelope-sniffing wrapper pointed at text-mode output mis-parses the bare document — by contract it must use `--as=json` (C5); the four-cell behavior is documented and tested." — (FM, `0021:F2`).

- [REQ-99] "nondeterministic bytes across runs — caught by the MVV's double-emit byte compare, never shipped silently; a partial graph after an incomplete traversal — structurally impossible, C4 refuses instead." — (FM, `0021:F4`).

## Implementation phases and prerequisites

- [REQ-100] "A8 verified — the exported completeness-carrying surface exists. It gates Phase 2: C4's `graph-export-too-large` MUST, S7, and the `disposition` table's ceiling row all depend on it, so Phase 2 cannot ship while it is Pending." — (PRE).

- [REQ-101] "Build the document-assembly component over `*table.Model` + `Reach`-with-edges; goldens pin `intrastate.graph/1`; the A2 differential test (observer vs recovery, subsumption-merge fixtures) decides the recovery mechanism before the schema freezes." — (IP, Phase 1).

- [REQ-102] "Register root `graph` with C1's arm set, `--emit`, the `TextLiner` ride, and the `graph-export-too-large` refusal." — (IP, Phase 2).

- [REQ-103] "`docs/cli-output-contract.md` gains the document section; `--help-all` gains the verb; the schema docs carry C2's `abstraction` marker and its soundness sentence." — (IP, Phase 4) — three separate documentation observables.

- [REQ-104] "The verb runs load + normalize + reachability only: it never runs the lint invariants, emits no findings, and exports any model that loads" — (AP) — the approach-level restatement of C4, carrying the same observables.

- [REQ-105] "selection flags → load (`table.LoadWithAdvisories`'s underlying load path — advisories are ignored here; they are lint's advisory channel, not graph data) → assemble → render (`--emit`) → respond gateway." — (TD) — advisories are ignored by the export path and never become document content.

## Cross-cutting (projected policy, `0021:G-cross-cutting`)

- [REQ-106] "This record claims no content-addressed identity and no replay-stable hash; the promise is byte identity of the emitted stream (C3, RT2), not a digest." — (XC) — a negative fence: no hash field, no digest claim.

- [REQ-107] "Every sequence is pre-sorted by C2's declared orders before marshaling: `rows` in RDR 0002's canonical row order, atoms in 0002's canonical atom order, `reach.nodes` by node key, `reach.edges` by (from, to, rule), and each tag's `values` array sorted and deduplicated — so construction order is unobservable." — (XC).

- [REQ-108] "field names are lowercase snake_case against `internal/table/dump.go::dumpColumns` … and no case folding is applied to values." — (XC).

- [REQ-109] "The shared encoder is non-HTML-escaping, so tag values and rule ids reach the wire unescaped" — (XC) — `<`, `>`, `&` serialize as themselves in the document.

- [REQ-110] "The export is a new root verb beside `lint`, adding no flag to `lint` and altering no existing output." — (XC) — directly assertable against `lint`'s flag set and output bytes.

---

## EXCLUDED

Clauses the auditor judges NOT testable in this record's implementation. Each
is a peer-owned contract, a scope this record defers, a code-siting claim, or a
negative with no observable of its own.

- EXCLUDED: "The new `graph-export-too-large` refusal code below mints no vocabulary of its own: it is a member of the CLIError `code` vocabulary `0029:C4` already tiers `append-only`, and it inherits that row's `seam: none (prose-only)` disposition — the tier stands, unasserted and legibly so, for the reason recorded there." — (NC, `0021:C1`) — the record states in terms that the tier is UNASSERTED and the seam is prose-only; there is no by-value assertion to write. The code's behaviour is REQ-52/REQ-81/REQ-88.

- EXCLUDED: "which merged nodes satisfy a `terminal` predicate set is the dead-end quantifier RDR 0015 owns (JDR 0001 §JD-23), and this record declares no evaluator of its own" — (NC, `0021:C2`) — peer-owned (RDR 0015). The testable residue is REQ-31 (no per-node terminal marking) and REQ-32 (declared sets travel).

- EXCLUDED: "the document is NOT a model source and no export→load inverse is claimed." — (NC, `0021:C2`) — restates RT3, already carried as REQ-67; not a second obligation.

- EXCLUDED: "what `intrastate lint` promises as a repository gate is RDR 0014's contract, which this record only leaves unchanged" — (NC, `0021:C4`) — peer-owned (RDR 0014); this record's own obligation is the neutrality byte-compare, REQ-49.

- EXCLUDED: "if a later revision adds one it MUST conform to JDR 0002 §D1 (projection before respond.OK, echo-group-only, enforced partition, always-keep core)." — (NC, `0021:C5`) — a conditional binding a FUTURE revision; unobservable in this implementation. The present-tense half is REQ-62.

- EXCLUDED: "One new document-assembly component (package boundary settled at Resolve; it sits beside the relation it exports)" — (TD) — a code-siting claim; the record explicitly leaves the package boundary to Resolve, so no byte or behaviour follows from it.

- EXCLUDED: "Extension point: the document is schema-versioned, so later producers (a declared-emit vocabulary, new invariant metadata) join additively without a second export surface." — (TD) — a claim about future producers; the assertable content is REQ-17/REQ-34.

- EXCLUDED: "The two `guard` functions are read as-is: this record changes nothing in that package, whose observation surface is RDR 0013's to evolve." — (CA, A3) — peer-owned (RDR 0013); a negative whose only observable is "no diff in `internal/guard`".

- EXCLUDED: "The property is pinned in CODE, not restated here — kata `yybx` ships the behaviour regression test and the stale-comment fix in isolation from this RDR; once landed, that test is the property's home." — (CA, A2) — the record explicitly sites this obligation in a separate kata, outside this implementation's scope.

- EXCLUDED: "Latent risk, not a falsifier: four call sites build ad hoc `SetEscapeHTML(false)` encoders instead of routing through `WriteJSONLine`, so the implementation must reuse the shared encoder rather than add a fifth." — (CA, A5) — the positive obligation is REQ-43 (emit through `clierr.WriteJSONLine`); "do not add a fifth ad hoc encoder" has no observable beyond it.

- EXCLUDED: "a build bump re-baselines DOT diffs and may add JSON fields (C3's stated narrowing); the cross-build golden makes the moment visible." — (FM, `0021:F3`) — explicitly marked "Silent (accepted, marked)"; the cross-build case is outside C3's (model, build) scope by REQ-44.

- EXCLUDED: "Diagnosis: byte-diff two runs (determinism), `jq .data` vs text output (mode agreement), `intrastate lint` before/after (neutrality oracle)." — (FM, `0021:F5`) — operator guidance; the three oracles it names are REQ-42, REQ-71, REQ-49.

- EXCLUDED: "Omitted as not applicable: build tool compatibility, licensing, deployment model, IDE compatibility, secret/credential lifecycle, memory management, concurrency model." — (XC) — a declared non-applicability list; nothing to assert.

- EXCLUDED: "All Critical Assumptions verified (A1/A2/A5/A6 are spikes; A2 — edge recovery equals the traversal's edges — gates Phase 1)" — (PRE) — a lifecycle gate discharged at lock (all eight assumptions are Verified), not an implementation obligation. A8's Phase-2 gate is carried separately as REQ-100 because it names a surface the implementation must build.

---

## ASSUMPTIONS

Implicit choices made where wording was imprecise but a single reading is
defensible.

- **A-1 — `--emit`'s default is applied, not merely documented.** C1 says
  "`--emit <format>` selects the document: `json` (default) or `dot`".
  **Reading taken:** invoking `graph` with no `--emit` produces the JSON
  document, byte-identical to an explicit `--emit=json`, so the default is a
  registered flag default rather than a doc-only statement. This makes REQ-8's
  default arm assertable and matches C5's "Under `--as=text` (the default)"
  phrasing for the sibling flag.

- **A-2 — "`schema` field … required leading" means first key on the wire.**
  C2 requires "a required leading `schema` field" and `XC`'s determinism list
  records "`schema` first". **Reading taken:** `schema` is the first key in the
  emitted JSON object, which a byte-level golden pins automatically; this is a
  property of the encoder's struct field order (A5), not a separate sort.

- **A-3 — the neutrality byte-compare covers `lint`'s stdout AND stderr.**
  C4 says lint's "verdict, finding set, and bytes are identical". **Reading
  taken:** "bytes" means the full observable output of the `lint` invocation —
  stdout and stderr both — plus the exit code, since a perturbation that moved
  an advisory to stderr would otherwise pass a stdout-only compare. A1 already
  fixes stderr as the advisory channel, so both streams are capturable.

- **A-4 — "one findings[] entry per load category" is mirrored, not
  re-derived.** C1 mirrors lint's `model-invalid` arm "with one findings[]
  entry per load category". **Reading taken:** the verb calls the same shared
  `internal/cli/flow_input.go::loadFindings` (`:223`) that `lint.go`'s arm
  calls, so the entry set and each entry's `code` are whatever that shipped
  loader produces — the record's Joint-check states exactly this ("satisfied by
  the shipped loader this record reuses, not by anything it declares").

- **A-5 — the `--emit dot` envelope member is an object, not a bare string.**
  C5 says `data` embeds "for `--emit dot` an object carrying the DOT text as
  the single REQUIRED string member `dot`". **Reading taken:** `data` is a JSON
  object `{"dot": "<text>"}`, so `jq -r .data.dot` is the documented unwrap and
  `data` is never a bare JSON string. "Single" constrains this object to that
  one member.

- **A-6 — `graph-export-too-large` refuses before any document bytes are
  written.** C4 requires "never a partial document". **Reading taken:**
  completeness is checked after the traversal and before any rendering or
  writing, so stdout carries zero document bytes on the refusing arm (REQ-83's
  "stdout carries NO document on every refusing arm"), not a truncated stream.

- **A-7 — the exported completeness surface is one new function, not two.**
  C4 says "the edges-carrying function Q3(c) already adds beside it returns
  completeness with the nodes and edges", and A8 describes "a NEW exported
  function beside `Reach` carrying nodes, edges, and completeness". **Reading
  taken:** a single new exported `graphlint` function returns all three, rather
  than a separate completeness accessor; `Reach` and `reach` keep their present
  signatures (`reach.go:73`, `:90`).

---

## QUESTIONS

None. Every clause this auditor read admits one defensible reading, or is
recorded above as an ASSUMPTION where the wording was imprecise but the
alternative readings do not produce materially different behaviour. The one
clause that previously carried a genuine ambiguity — the shape of a `reach`
node's `values` member — was resolved before lock: the `trace` mini-check
records it as "RESOLVED (was CONTRADICTION)", C2 now fixes the tag-keyed object
shape, and A7 is Verified. It is carried here as REQ-25, not as a question.
