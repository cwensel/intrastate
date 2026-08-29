Model: claude-fable-5

# Premortem critique — RDR 0021 lint normalized-graph export (`graph` verb, `--emit json|dot`)

Scope discipline: this critique is written from the brief alone. No repository file, design record,
or contract text was opened. Where a finding says "quote the clause," that is the mitigation — the
brief itself does not supply the clause, which is exactly the defect class being hunted.

## Findings ledger

| ID | passage/claim | Failure mode | Symptom user sees | Origin |
|----|---------------|--------------|-------------------|--------|
| P-1 | Claim 3: edge list "recovered AFTER the fixpoint (by re-running the successor function over final nodes with subsumption mapping)" | Post-hoc successor re-run over *merged* nodes is not the traversal lint ran; a merged node carries wider state than any pre-merge node, so re-running successors can fire rows the certified traversal never fired (and subsumption mapping can drop edges the traversal did take). Exported graph ≠ certified graph. Seed-2 class: false premise about the normalized value's shape, and the mechanism inverts the goal (the export exists to show what lint certified; it shows something else). | Formal checker derived from the export proves/refutes properties the linted model does not have; CI diff shows edge churn with no authored change | Element 4 / 2 |
| P-2 | Claim 1: envelope contract "carves non-flow command groups out ('lint, dump, parse ...')" — ellipsis | Selective citation (seed-1/3 class). The elided text and the surrounding clause are exactly where a closed-enumeration or "new verbs require amendment" veto would live. `graph` is not in the quoted list. | Envelope-contract owner rejects the verb at review or, worse, after ship a contract-conformance sweep flags `graph` as an unregistered root verb | Element 4 / 2 |
| P-3 | Claim 5: byte-identical "for one (model, build) pair" vs. invariant (b) "byte-for-byte deterministic for the same model" | Quiet narrowing of the stated invariant: (model) became (model, build). A tool upgrade in CI reorders or reformats the document; every open PR's graph diff goes red simultaneously. | CI graph-diff gate fails on every branch the day the tool is bumped, with zero model changes | Element 4 / 2 |
| P-4 | Claim 5: "every sequence pre-sorted by canonical orders and no Go map on the wire" | Marshaling determinism does not give value determinism. If the fixpoint's worklist/merge order depends on map iteration, the *node identities* (merged-node names, subsumption winners) differ run to run; sorting a nondeterministic set of names is still nondeterministic content. Sort keys with ties are a second hole. | Same commit, two CI runs, two different byte streams; "replay-safe" invariant falsified | Element 2 |
| P-5 | Claim 2: TextLiner "prints a payload's single string verbatim ... no stray bytes" | Line-oriented seams normalize: trailing-newline appended or stripped, CRLF handling, per-line processing of a multi-line payload; and text mode may interleave warnings/progress on stdout. One extra `\n` breaks `sha256`-based diff gates and DOT parsers are the lucky case — silent diff noise is the unlucky one. | `intrastate graph --emit dot | dot -Tsvg` works but the CI byte-diff never converges; or a warning line lands inside the JSON document | Element 2 |
| P-6 | Claim 6: skipping lint invariants "serves the outcome" | Two failure arms. (a) Some lint invariants may be *well-formedness preconditions* for the traversal (declared initial state exists, tags declared, domains finite); skipping them means the traversal runs on garbage — crash, non-termination up to the ceiling, or a confidently wrong graph. (b) Social failure: a clean, exported, versioned artifact reads as certification. Reviewer approves the diagram; the lint failure sits in a different CI job. | Panic/ceiling-blowout on a model lint would have refused with a clear finding; or a PR merged on the strength of a graph its own linter rejected | Element 2 |
| P-7 | Claim 9: abstraction marker + prose "suffice" | Markers do not travel. The DOT rendering — "the diagram is half the stated outcome" — has no natural cell for `abstraction: over-approx`, and an SVG pasted into a review carries nothing. Formal consumers: over-approximation makes safety checks sound but liveness/reachability-of-terminal results on the export meaningless; the brief's own sibling concern (terminal-reachability liveness, cli/0022) is the first thing someone will check against this graph. | External checker "verifies" liveness the runtime violates; team cites the export as proof | Element 2 |
| P-8 | Claim 4: normalized value "carries everything ... no re-parse" | Seed-2-adjacent: the normalizer is lossy by design (one authored key → two atoms is already on record). Authored identity (rule names/order, selection-context grouping as authored, domain declarations vs. inferred domains) may be gone or transformed. The document then can't be traced back to the TOML a human reviews, halving the "review it as a diagram" outcome. | Reviewer cannot map a graph edge to the authored rule that produced it | Element 4 / 2 |
| P-9 | Claim 8: peer viz CLI `--output-type svg|dot|json|scxml` as prior art | Borrowed authority (seed-3 class): no transcript of the peer's actual behavior is cited. Peers in this family commonly write files or non-deterministic JSON; and the peer's flag conflates format with mode — the very sin R5 rejects — so the citation cuts against the design where it's inconvenient and for it where it's not. | Prior-art warrant collapses at review; flag shape re-litigated post-ship | Element 4 |
| P-10 | R5 rejection ("conflates envelope mode with document format") vs. the chosen matrix | The chosen design's `--as text --emit json` cell emits a *bare* JSON document on stdout of a tool whose json-looking stdout otherwise means "envelope." A generic wrapper that sniffs `{` mis-parses it; `--as json --emit dot` buries DOT as an escaped string a human must unwrap. The rejection reason applies, diluted, to the winner — inverted-rejection escape class. | Org wrapper script crashes on `graph` output; users pipe envelope-wrapped DOT into `dot` and get a parse error | Element 2 / 4 |
| P-11 | R1 rejection: "error envelope has no graph carrier" | Asserts the error envelope is closed without quoting the closure clause (seed-1 class). If the envelope has an extensible data field, `lint --emit` was viable and cheaper — and it would have structurally guaranteed the export IS the certified graph (same run, same traversal), killing P-1 at the root. The separate verb creates two load+normalize+traverse call paths that can drift — the precise drift the RDR was chartered to eliminate. | Graph verb and lint disagree about the same model after a refactor touches one path | Element 2 / 4 |
| P-12 | R4 rejection: "non-recoverable by a prior locked decision" | Seed-3 pattern verbatim: a rejection resting on what a prior record allegedly says, unquoted. If the locked decision says "we won't guarantee round-trip," not "cannot be recovered," the cheapest alternative was dismissed on a misread. | Reviewer opens the cited decision, finds it doesn't say that; alternatives table invalid | Element 4 |
| P-13 | Claim 7: refuse outright at node ceiling; exit codes 2 vs 3 unassigned | The model most in need of visual inspection is a blown-up one; total refusal removes the tool exactly when wanted. And the brief never assigns the refusal an exit code: ceiling breach is arguably neither "bad model" (2) nor "environment" (3). CI scripts branch on the wrong code. | Engineer debugging state-space explosion gets "refused: too big" and nothing else; CI retries an error that will never succeed | Element 2 |
| P-14 | Claim 6 first half: "running no lint invariants ... guarantees verdict-neutrality structurally" | True but vacuous as stated — separate process, separate verb, nothing shared to corrupt. The real neutrality risk is unaddressed: if edge extraction (P-1) requires instrumenting the shared traversal, that instrumentation runs inside lint's code path too, and a bug there CAN alter lint. The structural guarantee holds only while the P-1 mechanism (post-hoc re-run) holds — and P-1 says it can't. | A later fix for P-1 quietly voids the neutrality argument nobody re-checked | Element 2 |
| P-15 | R2 rejection: "a shell redirect already covers the need" | Redirect purity depends on P-5 stdout purity. If text mode ever interleaves any non-document byte, redirect captures the pollution and the "covered" need isn't. The rejection is only as strong as claim 2, which is unproven. | Redirected `graph.json` fails schema validation because a deprecation warning is line 1 | Element 2 |
| P-16 | R3 rejection: "at the cost of one pure renderer" | Underpriced: DOT is a second canonical surface — identifier escaping (quotes/newlines in tag values), stable node ids across diffs, label truncation. Each is a determinism and correctness surface with its own golden tests; "one pure renderer" becomes a growing presentation dialect inside a validator CLI. | DOT output unparseable for a model whose tag value contains a quote; diagram diffs churn on label formatting fixes | Element 2 |

## 1. Prospective hindsight — the failure, as it happened

The verb shipped on schedule. `intrastate graph --emit json` produced a tidy
`intrastate.graph/1` document, the golden tests were green, and the CI graph-diff gate went in a
week later.

The first incident was quiet. A platform team derived a TLC model from the export — the whole
point of the RDR — and reported that terminal states were reachable from every node. They closed
the liveness ticket. Six weeks later a runtime trace showed an unreachable terminal. Post-incident
analysis found two compounding causes. First, the abstraction marker was present in the JSON but
the team had worked from the DOT rendering, which carried no marker at all; over-approximated
edges made liveness "hold." Second — and this was the one that ended the design's credibility —
the exported edge set was not the edge set lint had traversed. The edges were reconstructed
post-fixpoint by re-running the successor function over the final *merged* nodes. Merged nodes are
wider than any node the traversal actually visited, so the re-run fired rule rows the certified
traversal never fired, and the subsumption mapping silently dropped an edge the traversal did
take. The artifact everyone had been diffing "as what lint certifies" had never been that.

The second incident was loud. The org's generic CLI wrapper, which parses any intrastate stdout
that starts with `{` as a terminal envelope, was pointed at `graph` in its default
(`--as text --emit json`) configuration. The bare document parsed as an envelope with no `ok`
field; the wrapper classified every export as a malformed error and paged the on-call. The same
week, a PR was approved on the strength of a clean exported diagram; the model's lint job — a
separate CI stage — had exited 2. The graph verb, by design, "succeeds for any model that
loads," and the reviewer took a versioned, schema-stamped artifact as certification.

The cleanup uncovered the paper trail. The envelope-contract clause justifying "no amendment
needed" turned out to be an ellipsis: the elided sentence enumerated the exempt verbs closed and
required registration for new roots. The prior-art peer cited for the `--emit` shape wrote its
DOT to files, not stdout. And the locked decision cited to kill the dump-based alternative said
"round-trip not guaranteed," not "non-recoverable." Three citations, zero opened. The verb was
frozen behind a feature flag pending a re-propose in which the edge list is recorded *during* the
shared traversal — the mechanism R1's rejection had taken off the table.

## 2. Obstacle negation — claims

**C1 negated.** The carve-out is quoted with an ellipsis over exactly the span where a veto would
live. Assume the contract enumerates non-flow verbs closed, or requires every new root verb to
register its envelope behavior. Then `graph` ships in breach, and the `--as text` bare-document
behavior — the most contract-adjacent thing in the design — was never adjudicated by the contract
that owns `--as`. Scenario: contract conformance sweep (or the contract RDR's owner) flags the
verb; the "no amendment needed" premise was the RDR's cheapest sentence and its most expensive.

**C2 negated.** "Prints a payload's single string verbatim" was specced for one-line human
payloads. Assume the seam appends a terminal newline (most line-writers do), or normalizes the
payload's own trailing newline, or shares a writer with warning output in text mode. For a
byte-diffed artifact, one nondeterministic trailing byte is a broken invariant; for DOT, an
interleaved warning is a parse failure. The claim is stated as a property of the seam but was
never tested with a multi-kilobyte multi-line payload.

**C3 negated (primary).** The premise is that (successor function ∘ final merged nodes ∘
subsumption map) reproduces the traversal's edge relation. It does not in general: subsumption
merging is precisely the operation that makes a node's state a superset, and this project's own
ledger (seed 2) records that "supplying more state" changes which candidates fire — in that
episode it *removed* candidates against intuition. Either direction is fatal here: the re-run
emits edges the traversal never certified, or the mapping collapses distinct traversal edges into
one, or drops them. The export's charter is "the graph the linter actually certifies"; this
mechanism exports a related but different graph, and no consumer can tell.

**C4 negated.** Normalization is lossy and non-injective (one authored key → two atoms, already
on record). Assume any of: authored rule identity is not retained on normalized rows;
selection-context groups are flattened during normalization; tag domains are interned or widened;
terminal predicates are compiled to row form. Then the schema's promised fields ("tag declarations
with finite domains," "selection-context groups") either cannot be populated from the normalized
value or are populated with the normalized *shape*, which a human cannot map back to the TOML they
authored — and the "review it as a diagram" outcome degrades to "review an artifact only the
normalizer understands."

**C5 negated.** Two holes. (a) The determinism argument covers serialization, not the value:
if merged-node identity or merge order derives from map-iteration or worklist order inside the
fixpoint, canonical sorting launders nondeterministic content into stable-looking bytes that still
differ run to run. Ties in the canonical sort orders are the same hole smaller. (b) "(model,
build)" is not the invariant; invariant (b) says "same model." Cross-build instability means
every tool upgrade invalidates every stored graph diff baseline at once — a CI outage by design.

**C6 negated.** The structural half is vacuous (nothing shared, nothing corrupted) and becomes
false the moment C3's failure forces edge recording into the shared traversal (see P-14). The
outcome half inverts twice: models that lint refuses may violate preconditions the traversal
assumes (undeclared tag in an initial assignment, unbounded domain) — so "succeeds for any model
that loads" is either untrue or true only by emitting garbage; and a successful export from a
failing model is a certification-shaped artifact that reviewers will treat as certification.

**C7 negated.** Refusal is defensible for the diff gate but indefensible for the debugging
journey: the over-ceiling model is the one whose graph a human most needs to see. A partial
document with a required `truncated: true` marker serves that user; total refusal serves nobody
and the brief doesn't even assign the refusal an exit code (2? 3? neither fits cleanly).

**C8 negated.** No transcript, no invocation, no version of the peer is cited — the warrant is
the peer's flag *name*. Assume the peer writes `--output-type` results to files (common in that
tool family), or its JSON is unversioned and nondeterministic: then the prior art supports R2's
rejected alternative better than the chosen one. Worse, the peer's single flag conflating format
and destination is the shape R5 rejects — the citation is used selectively, load-bearing where
convenient, ignored where not.

**C9 negated.** The marker lives in the JSON; the outcomes live in the DOT/SVG and in the formal
model derived downstream, neither of which carries it. Over-approximation semantics are
property-class-dependent (safety: sound; liveness: void), and the brief's own ecosystem has a
live terminal-reachability-liveness concern — the first query a consumer will run is the one the
export cannot answer. "Doc prose" is the weakest control in the hierarchy and the only one
proposed.

## 2b. Obstacle negation — rejection reasons

**R1 inverted.** "The error envelope has no graph carrier" is asserted, not quoted (P-11). If the
envelope admits structured data on exit 2 — or if a lint run may emit an `ok` envelope containing
findings *and* a graph under a single-terminal-envelope reading — then `lint --emit` was viable.
And R1's rejection buys a real cost: a second load→normalize→traverse call path that can drift
from lint's, re-creating the "second artifact that drifts" problem this RDR exists to kill, one
layer down. The strongest version of this design shares one traversal with an edge observer;
R1's rejection is what forbade thinking about it.

**R2 inverted.** "Shell redirect covers it" is only as true as C2's stdout purity, which is
untested (P-15). It also silently narrows the audience to POSIX-shell contexts; a `--out` flag is
trivially deterministic and side-steps the entire TextLiner question. "Side-effecting write to a
read-only surface" proves too much: emitting bytes to stdout is also a write; the surface's
read-only-ness is about the *model*, which `--out` does not touch.

**R3 inverted.** "One pure renderer" is the famous last words of every embedded presentation
format (P-16). DOT brings escaping, stable node identity, label policy — each a determinism
surface with golden tests, each a magnet for "just add colors" feature requests inside a
validator CLI. JSON-only plus a blessed converter script keeps the certified surface minimal;
the rejection prices the renderer at its day-one cost.

**R4 inverted.** The rejection's first leg is a citation to a "prior locked decision" whose text
is not quoted (P-12) — the seed-3 pattern exactly. Its second leg ("the dump carries no
reachability relation") is real but repairable (extend the dump), so the whole rejection stands
on the unopened citation.

**R5 partially self-applied.** The chosen matrix reproduces the rejected conflation in two cells
(P-10): bare JSON on a "text" stream that elsewhere signals envelope-JSON by its first byte, and
DOT-as-escaped-string inside an envelope. The rejection reason was correct; the winner does not
fully escape it and the RDR does not say so.

## 3. Consumer artifacts — what would have caught each finding at review time

- **P-1/C3 — `TestGraphEdgesEqualTraversalEdges` (differential, corpus-driven).** Instrument the
  actual lint traversal with a recording observer; run the post-hoc recovery; assert edge-set
  equality over a fixture corpus that *must* include models with subsumption merges and models
  where a merged node enables a row no pre-merge node enabled. This single test decides the RDR's
  central mechanism. Its absence from the plan is the tell.
- **P-2/C1 — evidence file quoting the carve-out clause in full**, ellipsis-free, with the
  sentence before and after, plus a one-line conformance check: is the non-flow set closed? Review
  gate: no elided citation in a load-bearing claim.
- **P-3/P-4/C5 — `TestGraphByteStableAcrossRuns`** (`-count=100`, hash-compare, map-seed
  shuffling via `GODEBUG`/`hashrandom`) and **`TestGraphByteStableAcrossBuilds`** (golden file
  checked against the previous released binary in CI). The second test forces the (model) vs
  (model, build) decision into the open.
- **P-5/C2 — `TestTextModeStdoutIsExactlyTheDocument`**: run `graph --emit dot` and
  `--emit json` under conditions that provoke warnings; assert `stdout == document` byte-for-byte
  (including trailing-byte policy, stated in the schema doc) and all diagnostics on stderr.
- **P-6/C6 — user journey "PR with failing lint, passing graph"**: a written CI walkthrough
  showing both jobs' outputs side by side and what the reviewer sees; plus fixture
  `model_loads_but_lint_refuses/` (e.g., initial assignment to an undeclared tag) with an asserted,
  documented `graph` outcome — success with defined content, or a defined refusal — not whatever
  the traversal happens to do.
- **P-7/C9 — journey "derive a TLC model from the export"**: a table in the RDR stating which
  property classes are sound over the over-approximation; schema test asserting the abstraction
  marker is a *required* field; DOT golden asserting the marker is rendered in the graph header
  comment/label so the diagram carries it too.
- **P-8/C4 — `TestSchemaFieldCoverage`**: one fixture per schema field, populated end-to-end from
  the normalized value with no TOML re-parse — mandatory fixtures: a two-atom authored key, a
  selection-context group, a declared-but-unused tag domain. Fails at review if any field has no
  populating fixture.
- **P-9/C8 — prior-art transcript**: actual captured invocations of the cited peer (command,
  stdout/file behavior, determinism check of its JSON) committed as evidence. Review gate: no
  prior-art claim without a transcript.
- **P-10/R5 — journey "generic envelope wrapper meets graph"**: run the org-style wrapper
  (parse-stdout-as-envelope) against all four `--as × --emit` cells; assert defined behavior in
  each. Also `TestEnvelopeDotUnwrap`: a documented one-liner (`jq -r`) from envelope to `dot`
  input, tested.
- **P-11/R1 — evidence quoting the error-envelope schema** showing it has (or lacks) an
  extensible data carrier; plus `TestGraphAndLintShareNormalization` — a seam/call-graph test
  asserting both verbs reach the traversal through one function, so drift is structural, not
  disciplinary.
- **P-12/R4 — the locked decision quoted verbatim** in the alternatives table, with the sentence
  containing "non-recoverable" (or its absence) visible.
- **P-13/C7 — `TestCeilingRefusalContract`**: over-ceiling fixture; assert exit code (chosen and
  documented as 2 or 3 with rationale), stderr message naming the ceiling and the count reached;
  journey "debug a state-space explosion" showing what the user gets instead of a graph.
- **P-15/R2 — same artifact as P-5**; the R2 rejection is valid iff that test exists and passes.
- **P-16/R3 — `TestDotEscaping`**: fixture with tag values containing `"`, newline, non-ASCII;
  assert `dot -Tsvg` parses the output; golden-stability test for node identifiers across an
  unrelated model edit.

## 4. Refutation targets — seed-class hunt results

**Selective citation of a contract (seed 1):** three instances. (i) C1's carve-out quote with a
load-bearing ellipsis — favorable fragment quoted, potential vetoing sentence elided, the exact
five-to-zero pattern from the ledger (P-2). (ii) R1's "error envelope has no graph carrier" —
a contract property asserted without the schema (P-11). (iii) C5's silent rewrite of invariant
(b) from "same model" to "(model, build)" — citing one's own invariant selectively (P-3).

**False premise about the normalized value's shape (seed 2):** the RDR's central mechanism.
C3 assumes successor-over-merged-nodes reproduces the traversal's edges; the project's own
escaped defect ("one key, two atoms; more state removed candidates") is a standing counterexample
to exactly this style of assumption about the normalizer/traversal, and here the mechanism again
inverts the proposal's goal — the export of "what lint certified" certifies something else (P-1).
C4 repeats the class at the schema level: fields promised from a value that may not carry them
(P-8).

**Borrowed authority (seed 3):** two instances. C8 cites a peer CLI by flag spelling with no
transcript — every reviewer can nod without anyone having run it (P-9). R4 rests a rejection on
what a "prior locked decision" says, unquoted; if the decision says round-trip-not-guaranteed
rather than non-recoverable, the cheapest alternative died on a misread (P-12).

**Inverted-rejection escape:** R5's rejection reason ("conflates envelope mode with document
format") is correct and partially indicts the chosen matrix's own `--as text --emit json` and
`--as json --emit dot` cells (P-10); R1's rejection forecloses the one mechanism (shared-traversal
edge observer) that would dissolve the design's central defect (P-11 ↔ P-1).

## Verdict rationale

The approach's *surface* (a root `graph` verb, versioned JSON document, DOT sibling, ceiling
refusal) survives. Its *mechanism* does not survive P-1 as specced: post-fixpoint edge recovery
cannot be assumed to reproduce the certified relation, and the fix (record edges during the
shared traversal) is available inside the same approach but voids the current verdict-neutrality
argument, so both must be re-specced together. The three citation defects (P-2, P-9, P-12) are
review-time repairable by opening and quoting the sources. None of this forces a different
approach; all of it forces changes before lock.
