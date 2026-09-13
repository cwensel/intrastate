Model: claude-opus-5

# Critique — RDR 0021, lint's normalized-graph export

## 6. Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0021:C4` | Mints refusal code `graph-export-too-large` for the SAME condition, same `reach()` completeness bit and same `nodeCeiling = 4096`, that already ships as `CodeProductTooLarge = "graph-product-too-large"`. The blocking tier is CLOSED and pinned by exact-set equality in `findings_0006_test.go`; `NodeCeiling()`'s own doc says "Exceeding it is `graph-product-too-large`". | One bound reports under two names depending on which verb you ran. A CI matcher for the lint code silently never fires on the export. | §1, §4, premortem, AT-1 |
| C-2 | `0021:A8` | Load-bearing assumption is `Pending` while `C4` states its consequence as normative MUST, `S7` tests it, the disposition table lists it as decided, and Phase 2 schedules it. The record's own gate text forbids exactly this. | Lock proceeds on an unverified surface; the ceiling arm is discovered unimplementable during Phase 2. | §1, §2, premortem, AT-2 |
| C-3 | `0021:C2` | `values` is fixed as the tag's declared value array with no mention of `OpaqueValue = "<opaque>"`, which `heldValues` substitutes for every tag with no finite declared domain. The sentinel reaches `Node.Values` and is pinned by three shipped tests. | A formal-model generator reads `<opaque>` as a real state value and derives a model over a state that does not exist. | §1, §3, §4, premortem, AT-3 |
| C-4 | `0021:A6` | Verified on a spike whose `Doc.Values` is `map[string]string` — one string per tag — while C2/A7 normatively fix `map[string][]string`. The only DOT renderer ever built for this document cannot represent a merged multi-valued node, the case the abstraction exists for. | The DOT diagram silently drops all but one value of every merged node; the diagram and the JSON disagree. | §1, §2, §4, AT-4 |
| C-5 | `0021:S6` | The oracle says "invert A6's recorded escaping ORDER", but that order exists only as `dotQuote` in a deleted `/tmp/a6-spike/main.go`. `C2` keeps DOT styling non-normative, so no clause pins a grammar to invert. | The escaping test cannot be written as specified; the mis-escaping bug A6 actually hit ships. | §1, §2, AT-5 |
| C-6 | `0021:C2` | Exports "the declared `[initial]` assignments" unfiltered, while `reach()` roots the graph only at `ProvenanceOwned` entries and passes them through `heldValues`. | The document's `initial` block names a tag that appears in no node; the diagram has no such start state. | §1, §4, AT-6 |
| C-7 | `0021:C5` | The `--as=json --emit dot` cell says `data` carries "a single string field" and never names the key; `C2`'s normative spellings cover the document, not the envelope wrapper. | `jq .data.???` — the documented one-step unwrap has no path. Two implementations both conform. | §1, §3, AT-7 |
| C-8 | `0021:A1` | Rides `TextLiner`, whose interface contract is "its whole content is one canonical LINE — a build identity, say", with a multi-kilobyte multi-line document, citing `version` as precedent. Widens a seam by usage, not amendment. | None at first. The seam's next maintainer optimizes for the documented single-line case and truncates the export. | §2, §3 |
| C-9 | `0021:§metadata` | `graph-export-too-large` is already published in RDR 0029, which is **Final** and no-amend, as one of 0021's vocabularies. Renaming the code to the shipped one now requires touching a locked record's cited literal. | Two Final records disagree with the binary about what the code is called. | §1, §4 |
| C-10 | `0021:RT1` | Claims `json-decode ∘ export = value identity on the exported projection` — trivially true, since the projection is defined as whatever was exported. No invariant ties the document to the traversal it claims to certify. | A document that faithfully round-trips through JSON while misrepresenting the graph passes every stated inverse. | §3, §4, AT-8 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The ceiling refusal ships under a code that does not exist, for a condition that already has one

**Root cause.** `C4` requires the export to refuse with `graph-export-too-large` when the traversal is incomplete under the node ceiling. That code appears nowhere in the tree. The condition it names — `reach()` returning `complete == false` because `len(nodes) > nodeCeiling` — is not a new condition. It is the exact bit `analysis.go::checkNodeCeiling` fires on, and lint already reports it as `CodeProductTooLarge = "graph-product-too-large"` (`taxonomy.go:34`). `NodeCeiling()`'s own doc comment states it outright: "Exceeding it is `graph-product-too-large`".

**The passage that enabled it.** `C4` spent its revision budget on the wrong half of the sentence. The 3amigo implementer persona caught that C4 had misidentified `checkNodeCeiling` as "a distinct guard-product bound", and the fix corrected the *prose* — "the SAME bound … observed at two call sites, which is the neutrality rule above rather than an exception to it". Nobody then asked the obvious follow-up: if it is the same bound at two call sites, why does it get a second name? The corrected sentence makes the defect *worse*, because it now explicitly argues these are one bound while the clause continues to mint a second code for it.

There is a second-order problem. The blocking tier is CLOSED (`0006:C17`) and pinned by exact-set equality:

```go
got := slices.Clone(graphlint.BlockingCodes())
slices.Sort(got)
if !slices.Equal(got, want) { ... }
```

So the implementer has two doors, both bad. Add the code to `blockingCodes` and a shipped RDR 0006 test fails, dragging a locked record's closed-tier guarantee into this RDR's scope. Or keep it out of the taxonomy as a bare `CLIError` code, and the binary now has a graph-family refusal that `graphlint.BlockingCodes()` does not enumerate — which is precisely the "unassigned machine-readable surface" that RDR 0029:C4 calls a defect.

**Symptom the user sees.** A CI pipeline greps for `graph-product-too-large` to catch models that outgrew the ceiling. It matches on `intrastate lint` and never matches on `intrastate graph`, because the identical condition on the identical constant reports under a different name. The model that is too large to lint is also too large to export, and the pipeline notices only one of them.

### 1.2 The opaque sentinel escapes onto the wire as a real value

**Root cause.** `reach.go` abstracts any tag with no finite declared domain to a single representative:

```go
const OpaqueValue = "<opaque>"

func heldValues(m *table.Model, key string, value []string) []string {
	if _, finite := guard.AssignmentCount(m.Tags[key]); !finite {
		return []string{OpaqueValue}
	}
	return canonicalValues(value)
}
```

This is not an edge case. It is called at both places a node acquires values — root seeding (`reach.go:120`) and successor production (`reach.go:317`) — and three shipped tests in `opaque_0006_test.go` pin it.

**The passage that enabled it.** `C2` fixes `values` as "an OBJECT keyed by tag name whose every value is that tag's value ARRAY, sorted and deduplicated — the shape `reach.go::Node`'s `Values map[string][]string` projects without invention". "Without invention" is doing enormous work here, and it is what caused the miss. `A7` verified the *container* shape — map to object, slice to array — and never audited the *inhabitants*. The word "projects" made the members look like someone else's problem.

The sentinel is spelled `<opaque>`. `C3` routes the document through the non-HTML-escaping encoder, so `<` and `>` serialize as themselves. On the wire it is an ordinary string, indistinguishable from a declared domain value that happens to be spelled `<opaque>`.

**Symptom the user sees.** The RDR's headline user — "let any formal model be *derived* from the code instead" — writes a generator over the JSON. It reads `{"stage": ["<opaque>"]}` and emits a TLA+ model with a state constant named `<opaque>`. The model checks green. The property it verified is about a state the runtime never enters, and every real value of that tag was never modeled at all. This is a false green in a tool whose stated purpose is to stop drift between the spec and the source of truth.

### 1.3 The DOT arm cannot represent a merged node

**Root cause.** `A6` is `Verified` by a spike whose document type declares:

```go
type Node struct {
	ID     string            `json:"id"`
	Values map[string]string `json:"values"`
}
```

One string per tag. But a node in this abstraction is a *set* — that is the entire content of the merge rule, and `C2`/`A7` normatively fix `map[string][]string`. The spike rendered fixtures that could not exercise a multi-valued node, then reported an "exact set match" against them.

**The passage that enabled it.** `A6`'s Evidence anticipates a naming mismatch and waves it off: "The spike's field spellings are provisional placeholders (`abstraction_marker`, `terminal_satisfying`); C2's spellings govern." That sentence trained every subsequent reader to treat divergence between the spike and C2 as cosmetic. The divergence in `Values` is not a spelling — it is a different type with different cardinality, and it sits three lines from the spellings the note excuses.

`joinNodes` unions value sets precisely so a node reached by two paths writing different values carries both. That is the common case in any real state machine, and it is the case the spike's `map[string]string` cannot hold.

**Symptom the user sees.** A reviewer runs `intrastate graph --model flow.toml --emit dot | dot -Tsvg`. Every merged node renders with one value where the JSON document carries several. The diagram, which the Problem Statement names as half the user outcome, disagrees with the document it is supposedly a projection of — and `S6`'s oracle compares node/edge *id sets*, which match, so the test passes.

---

## 2. The one section rewritten within 6 weeks

**`0021:§pre-lock-mini-checks`** (lines 577–639, 6.6K — the largest non-contract block in the record).

It will be rewritten because it is the only section whose every row is a claim about an artifact that does not exist yet, written in the register of a test report that already ran. Its own preamble admits the problem and then proceeds anyway: "Witnesses are spike output, not exporter output — the exporter is unbuilt."

Five tables — `authority`, `oracle`, `fidelity`, `disposition`, `trace` — restate C1–C5 in a second normative voice. That is a duplicate contract, and duplicates drift. The drift has already started *before implementation*: the 3amigo QA persona found the `oracle` and `trace` rows for step 3 still asserting the exit-0-plus-set-match bar that `S6` had just declared insufficient. Two rows went stale from one scenario edit, and were caught only because a persona widened scope beyond its brief. That is the failure mode of the whole section, observed once already, with four more tables and no mechanism preventing the next occurrence.

Then Phase 1 runs and the goldens land. Every `fidelity` and `oracle` row currently cites a spike; each must be re-pointed at a real fixture. `C-1` forces a rewrite of the `disposition` table's ceiling row. `C-3` and `C-4` force rewrites of `oracle` row 3 and `trace` row 3 — the same two rows that already went stale. `C-5` forces a rewrite of the escaping oracle, since the cited escaping order lives in a deleted `/tmp` file.

The record says these tables were "cue-fired" at Stage 5 — generated to satisfy a lens rather than because anyone needed them. Six weeks after shipping, the goldens are the oracle, the tables are a stale second copy, and the cheapest correct action is deletion.

---

## 3. The assumption that will not survive first contact

**`A3` — that the normalized model value carries everything the document needs.**

Not the narrowed version. `A3` was already narrowed once at Resolve, when the edge relation turned out not to be carried. The ruling patched that one hole and left the underlying claim — that the *model value* is the right source for the *document* — standing. It is that claim that breaks.

The first real user is the one the Problem Statement names: someone deriving a formal model. They discover three things the model value cannot tell them, and every one is a schema change inside a locked wire format.

**First, they cannot distinguish an abstraction from a value.** `<opaque>` arrives as a string (C-3). The user needs to know "this tag was abstracted because its domain is not finite" — a fact `guard.AssignmentCount` computes and the document throws away. The fix is a new field marking abstracted tags.

**Second, they cannot tell which `initial` entries root the graph.** `C2` exports the declared `[initial]` assignments; `reach()` roots only the `ProvenanceOwned` subset (C-6). A generator that seeds its model from `initial` seeds states the traversal never used.

**Third, the `abstraction` marker is one token for the whole document.** `C2` fixes the value `declared-over-approximation` and states the soundness rule: universal claims hold, existence claims may be spurious. But over-approximation is *per-edge*. `matchSatisfiable` prunes owned-tag atoms and deliberately does not prune observed/recognized match atoms or any guard atom — "every such edge is taken as traversable". So some edges are exact and some are speculative, and the document flattens that to one flag. A user proving a universal property is fine; a user asking "which of these 400 edges are real?" — the actual question when a diagram looks wrong — gets no answer, and must re-derive it from the model.

Each is additive under `C2`'s `/1` rule, so nothing breaks loudly. Instead `intrastate.graph/1` accretes three compensating fields in its first quarter, and the fourth consumer reads a schema whose original fields are all subtly wrong to use alone.

---

## 4. Premortem

*Written from six months after ship.*

We shipped `intrastate graph` in week three. The MVV passed: four invocations, byte-identical output, `jq .data` agreed with text mode, `intrastate lint` was byte-identical before and after. C4's neutrality oracle held exactly as designed. Nothing in the acceptance floor touched what actually broke.

**Week 5 — the CI matcher that never fired.** The platform team had a pipeline step matching `graph-product-too-large` to catch models outgrowing the ceiling. When they added `intrastate graph --model flow.toml > graph.json` to the same job, a 5000-node model refused with `graph-export-too-large` and the matcher missed it. `runGraph`'s ceiling arm and `checkNodeCeiling` were reading the same `complete` bool off the same `reach()` behind the same `nodeCeiling = 4096`, and reporting it under two names. The fix was one string. It could not be made in one place: RDR 0029 is Final and had already published `graph-export-too-large` in its census, so correcting the binary put a locked record in conflict with the tree. We shipped the rename and carried an erratum.

**Week 7 — the formal model that verified nothing.** A team derived a TLA+ spec from `graph.json` for a flow whose `assignee` tag is a free-form string with no declared domain. `heldValues` had abstracted it to `["<opaque>"]`, `buildNodeValues` copied `Node.Values` into the document without inspection, and the generator emitted `assignee \in {"<opaque>"}`. Two hundred states of real behavior collapsed to one fictional constant. TLC checked it green in four seconds. They shipped on that green and hit the bug in production three weeks later — a state pair the abstraction had erased. When we traced it back, nothing had malfunctioned: `A7` verified the container shape, `C2` said "projects without invention", and the sentinel was never in anyone's field of view.

**Week 9 — the diagram that lied.** A reviewer rendered `--emit dot` for a model where two rows write different values to `stage`. `joinNodes` had unioned them into one node holding `["review","approved"]`. `renderDOT`, written against A6's spike, took `values[tag]` as a single string and printed `stage=review`. The JSON carried both. `S6` compared node and edge *id* sets — identical, since ids come from `(Node).key` — and passed. We found it because a reviewer approved a PR believing a transition to `approved` did not exist. The renderer had been wrong since the first commit, built faithfully from a spike whose `Doc.Values` was `map[string]string` while C2 said `map[string][]string`, under an Evidence note telling us the spike's divergences were provisional spellings.

**Week 10 — the test we could not write.** Fixing the renderer meant writing the escaping test `S6` specified: invert A6's recorded escaping order and assert equality with the source id. A6's order lived in `dotQuote` in `/tmp/a6-spike/main.go`, deleted with the tmpdir months earlier. `C2` puts DOT styling outside the normative boundary, so no clause pinned a grammar to invert. We wrote an approximate inverse from the Evidence prose — "backslash, then quote, then newline" — and it disagreed with our renderer on a tag value containing a literal backslash-n. Neither artifact was authoritative. We invented a grammar, pinned it in code, and the RDR has never described what ships.

**What the record got right.** C4's neutrality held perfectly — `intrastate lint` never moved a byte. C3's determinism held. The A2 edge recovery, the one thing the premortem and three review rounds hammered, was exactly correct: every edge equaled the traversal's own.

**What the record got wrong.** Every failure above is a *content* defect — what the fields mean — while C1–C5, the mini-checks, and nine test scenarios are almost entirely about *form*: byte order, sort order, empty-vs-null, field spellings, envelope cells. The record verified the document's shape to four decimal places and never asked whether the values inside it meant what a consumer would take them to mean. `RT1` made that gap invisible: `json-decode ∘ export = value identity on the exported projection` is a tautology over whatever we chose to export. There was no invariant tying the document to the graph.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

**AT-1 — refusal code is not a synonym for a shipped code (C-1)**

```gherkin
Given the set of finding codes graphlint publishes via BlockingCodes()
And a proposed new refusal code from any clause of this RDR
When the condition the new code names is compared to the condition
     each published code names
Then no two codes may name the same predicate over the same constant
And if the predicate is reach()'s completeness bit against nodeCeiling,
     the code MUST be the one already published for it
```
At review: `C4` names `reach()` completeness under `graphlint.NodeCeiling()`. `checkNodeCeiling` fires on that bit under `graph-product-too-large`. Same predicate, same constant, two codes → BLOCK.

**AT-2 — no Pending assumption carries normative weight (C-2)**

```gherkin
Given an assumption whose Status is Pending or Unverified
When the record is searched for clauses depending on it
Then no MUST-clause, test scenario, disposition row, or phase may
     state its consequence as settled
```
At review: `A8` is Pending; `C4` states "Completeness MUST reach the verb through an EXPORTED `graphlint` surface", `S7` tests the arm, the disposition table lists the refusal, Phase 2 schedules it → BLOCK. This is the record's own gate text, unenforced.

**AT-3 — every inhabitant of an exported field is enumerated (C-3)**

```gherkin
Given a document field projected from a Go type
When the set of values that type can hold is enumerated from source
Then every sentinel or abstract value MUST be named in the clause
     fixing the field, and distinguishable on the wire from a
     user-authored value of the same spelling
```
At review: `Node.Values` can hold `OpaqueValue = "<opaque>"` via `heldValues`. `C2` names it nowhere; nothing on the wire distinguishes it from a declared value spelled `<opaque>` → BLOCK.

**AT-4 — a spike's types match the clause it verifies (C-4)**

```gherkin
Given a spike cited as Evidence for a Verified assumption
When each field of the spike's document type is compared to the
     normative clause
Then any difference in TYPE or CARDINALITY invalidates the Verified
     stamp, and only field-NAME differences may be waived as
     provisional spellings
```
At review: A6's `Doc.Values` is `map[string]string`; `C2`/`A7` fix `map[string][]string`. Cardinality mismatch → the DOT arm was never exercised on a merged node → BLOCK.

**AT-5 — a test oracle cites something a tester can open (C-5)**

```gherkin
Given a test scenario whose oracle names an inverse or decoder
When the cited definition is resolved
Then it MUST resolve to a normative clause or a file under version
     control — never to spike source in a temporary directory
```
At review: `S6` says "inverting A6's recorded escaping ORDER"; the order is `dotQuote` in `/tmp/a6-spike/main.go`, and `C2` leaves DOT styling non-normative → BLOCK.

**AT-6 — the document's root equals the traversal's root (C-6)**

```gherkin
Given a model whose [initial] assigns a tag of non-owned provenance
When the document is exported
Then either every initial entry appears in some node, or the clause
     states which entries are declared-but-not-rooted
```
At review: `C2` exports declared `[initial]`; `reach()` filters to `ProvenanceOwned` and applies `heldValues` → the blocks disagree, unstated → BLOCK.

**AT-7 — every envelope cell names its key (C-7)**

```gherkin
Given the --as × --emit matrix
When each cell's stdout is described
Then a consumer MUST be able to write the exact accessor path from
     the record alone
```
At review: `--as=json --emit dot` says "a single string field carrying the DOT text" with no key; the "one `jq -r` step" has no path → BLOCK.

**AT-8 — an inverse invariant is falsifiable (C-10)**

```gherkin
Given a claimed round-trip or inverse invariant
When a document is constructed that satisfies it while
     misrepresenting the source relation
Then if such a document exists, the invariant is a tautology and
     does not constrain the implementation
```
At review: `RT1` quantifies over "the exported projection" — whatever was exported. A document with `<opaque>` for every value, a single-valued DOT projection, and a non-rooted `initial` satisfies RT1, RT2 and RT3 completely. No invariant ties the document to the traversal → the inverse section constrains nothing.
