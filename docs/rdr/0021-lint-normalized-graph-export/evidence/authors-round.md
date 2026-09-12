Model: claude-opus-5[1m]

# Author's round — RDR 0021 (Stage 4 Resolve, 2026-09-12)

Stage 4's single consolidated author interaction. Nothing below has been
written into the record: an unapproved fixture is not Evidence, and a
pre-applied recommendation is a rewrite waiting on a different answer.
The record file is unmodified this pass.

Per-assumption verdicts (the basis for everything below):

| CA | Method | Status | Evidence |
| --- | --- | --- | --- |
| A1 | Spike | Verified | `evidence/spikes/a1-textliner.md` — `stdout == document + "\n"` byte-exact |
| A2 | Spike | Verified (conditional — see Q3) | `evidence/spikes/a2-edge-recovery.md` — recovery == in-traversal observer |
| A3 | Source Search | **FALSIFIED in part** — see Q3 | `evidence/research/a3-a4.md` |
| A4 | Peer RDR | Verified, citation re-anchor owed — see note | `0005:C1` scope sentence |
| A5 | Spike | Verified | `evidence/spikes/a5-encoder-stability.md` — 25-process byte identity |
| A6 | Spike | Verified | `evidence/spikes/a6-dot-derivability.md` — zero graphlint imports, `dot -Tsvg` exit 0 |

A4 needs no author ruling. `--outcome ground` returned `apply` (the question is
settled in source). A4's conclusion holds — `0005:C1`'s MUST-clauses are scoped
to the skill-integration `flow` group ("All four verbs MUST…", "flow next MUST…")
and impose nothing on a root `graph` verb — but its cited text is wrong: the
parenthetical "Other command groups (lint, dump, parse)" is a closed enumeration
that does not name `graph`. Re-anchoring A4's Evidence to C1's scope sentence is
a citation fix applied from source, not a fork. Peer 0005 is `status=Implemented`,
`gate_stale=false`.

---

## FIXTURES FOR APPROVAL

An approved fixture is recorded as a normative fixture on its assumption's
Evidence line (and in the covering Testing Strategy scenario's Expected),
citing the spike artifact that produced it. Rejected → the assumption is
genuinely open and resolves by another Method.

### F1 — A1: text-mode passthrough byte shape

Grounding: `evidence/spikes/a1-textliner.md`, run through the real
`respond.OK` path (`internal/cli/respond/respond.go:166`, `Fprintln` of
`TextLiner.TextLine()`).

    document_len_bytes: 14494   (300-line payload)
    stdout_len_bytes:   14495   (document + exactly one "\n")
    stdout_sha256:      d2972f38be5499c1f1fd60ff5d022868d1b576da5d703571c85f284e564f325a
    stderr_len_bytes:   0       (empty-string sha256 e3b0c442...)

Holds identically with advisories provoked: stdout unchanged byte-for-byte,
`note:`/`warning:` land only on stderr. A decoy line `{"status":"ok"}` inside
the payload passed through unsniffed and unreordered.

Proposed normative fixture: **text-mode stdout is exactly `document + "\n"`,
and stderr is empty absent advisories.** Backs C5 and Testing Strategy
scenario 5.

APPROVE / REJECT: ________

### F2 — A5: canonical document line and empty-vs-nil rendering

Grounding: `evidence/spikes/a5-encoder-stability.md`, emitted through
`internal/cli/clierr/clierr.go:174 WriteJSONLine` (`SetEscapeHTML(false)`),
double-emitted and re-emitted in 25 fresh processes under map-seed variation —
all identical.

    sha256: dcb0259aa822615592f11b65d571ba8a8b33c03f24b98b3bd54a3a177af8049f
    (independently recomputed in Python over the exact line bytes + trailing newline)

The emitted line is recorded verbatim in the spike file (a C2-shaped document
carrying HTML-significant chars, quotes, newlines, non-ASCII). Confirmed on the
wire: `<`, `>`, `&` serialize as themselves; `schema` is first-emitted with
value `intrastate.graph/1`; output is compact with exactly one trailing `\n`.

Proposed normative fixture: **that exact line + sha256 as the canonical
`intrastate.graph/1` emission for the fixture document.** Backs C3 and Testing
Strategy scenario 2.

APPROVE / REJECT: ________

### F3 — A6: canonical DOT rendering

Grounding: `evidence/spikes/a6-dot-derivability.md`. Rendered from a decoded
document with zero `internal/graphlint` imports; `dot -Tsvg` exit 0 under
graphviz 12.2.1. Node/edge set and abstraction marker matched the document
exactly (marker in both header comment and graph `label`).

    // abstraction: declared over-approximation: merged nodes, guard/observed atoms unpruned
    digraph reach {
      label="abstraction: declared over-approximation: merged nodes, guard/observed atoms unpruned";
      "n0" [label="n0", shape=doublecircle, color=blue];
      "n1" [label="n1\\nstatus=held"];
      "n2" [label="n2\\nstatus=absent", peripheries=2, style=filled, fillcolor=lightgrey];
      "n0" -> "n1" [label="r-open"];
      "n1" -> "n2" [label="r-close"];
    }

Escaping rule established (order is load-bearing): backslash → `\\` FIRST,
then `"` → `\"`, then literal newline → the two-character DOT escape. Non-ASCII
passes through. Every node id and label is wrapped, so a hostile id cannot
break out of the node statement.

TWO CAVEATS THE AUTHOR SHOULD WEIGH BEFORE APPROVING:
1. The rendered label above shows `\\n` where the stated rule specifies the
   two-character escape `\n`. The rule text and the emitted bytes disagree;
   which one is normative needs your call (C2 makes the node/edge SET and the
   marker normative, styling not — so this may fall outside the fixture).
2. The spike's `Doc` type is a reconstruction of C2's field list, not a decode
   against a real exported value (the export does not exist yet). Derivability
   is validated at the level C2 fixes now; exact field spellings C2 defers to
   Resolve.

Proposed normative fixture: **the node/edge set and abstraction-marker
placement above** (not the styling). Backs A6 and Testing Strategy scenario 6.

APPROVE / REJECT: ________

---

## QUESTIONS ONLY THE AUTHOR CAN SETTLE

### Q1 — the contract count blocks the Profile latch (BLOCKING this stage)

`--outcome profile` returns `dispositions.op: stop`, `emit.op:
stopped:split-signal`: "Two or more durable contracts — the record spans more
than one seam." The accretion floor is `none` (`floor-below-two`), so no tier
raise is in play; `contracts=5`, `contracts_durable=2+`, `contracts_prose=true`.

The tension is real and internal to the record: the Normative Contracts
preamble asserts the five blocks "are facets of the one contract this RDR owns
— the export surface and its document — stated separately by concern", yet they
are fenced as five separately-labelled durable contracts, which is what the
model counts.

`emit.surface` names exactly two dispositions:
- **One seam** → collapse C1–C5 into clauses of a single `**C1**` fence
  (matching the preamble's own claim), then the Profile resolves.
- **Separate seams** → the RDR spans more than one seam; split it, back to
  Stage 2/3.

I cannot pick: this is a scope judgement about what this RDR owns, and Stage 4
must not choose a profile past a `stop`. Until it is answered the `Profile`
field cannot be written, and `/rdr-prelock grounding` (which the lens row would
otherwise name) has no latch to read.

RULING: ________

### Q2 — A5: does the wire commit to `[]` or `null` for "no items"?

Observed, directly: an explicit empty slice emits `[]`; a nil slice emits
`null`. Both are valid JSON and they are a distinguishable wire difference.
C2's strictly-additive evolution rule ("a consumer ignoring unknown fields
keeps working") depends on producers consistently choosing one for empty
collections — `atoms`, `writes`, `requires_owned`, `members`, `values`,
`edges`, `nodes`. A consumer written against `[]` breaks on a later `null`.

C2 does not currently state which. Recommended (not applied): commit to `[]`
for every declared collection, reserving `null` for genuinely absent optional
objects — the fixture already renders `[]` throughout because it uses
`[]T{}` rather than nil.

RULING: ________

### Q3 — A3's edges clause, and how durable A2's equality is

These are one question because the answer to A2's caveat determines A3's wording.

A3 claims the normalized model value "carries everything the document needs".
The field walk confirms every C2 field has a carrier EXCEPT the reachability
edges `{from, to, rule}`: `reach.go::Reach` returns only `[]Node`, and no
`Edge` type or edges function exists anywhere in the tree. The rule id needed
for the triple IS on `table.Row` (`RuleID`) and is in scope inside
`successorsOf`, so the relation is constructible from the model value with no
TOML re-parse — A3's "no re-parse" half survives; its "carries everything" half
does not, as written.

A2 then verified that post-hoc recovery equals the traversal's own edges
exactly (9 adversarial fixtures, 36 real models, with negative and positive
controls proving the comparator detects divergence). But the equality is
STRUCTURAL, not incidental: it holds because `indexOf` uses `subsumes`, which
makes `reach()`'s in-place widening arm dead code (exhaustive lemma; zero
widening events observed). The spike's positive control deliberately weakened
`indexOf` to merge on presence footprint — the widening arm went live and
recovery diverged, reproducing premortem P-1 exactly.

So: A2 holds only while `indexOf` uses `subsumes`.

Recommended (not applied), needing your ruling on each part:
(a) Narrow A3 to "…with no re-parse of the authored TOML; the edge relation is
    derived in the export path, which the model value fully determines."
(b) Record A2's Verified status as conditional on the `subsumes` property, and
    add a regression test pinning the widening branch unreachable — the spike's
    `next_action` asks for exactly this.
(c) Confirm the edges-carrier is new code in the export path (not an
    extension of `Reach`'s public signature), since C4 requires reaching the
    traversal through the same `graphlint` entry surface lint uses.

Any narrowing of a normative clause here triggers rdr-common §amendment-sweep
before lock.

RULING: ________

---

## Not asked (resolved without the author)

- Reuse audit: no existing graph/export/DOT capability. Root verbs are
  `version`, `lint`, `flow`, `docs`; no `--emit`, no `digraph`/graphviz, and no
  pre-existing JSON serialization of the model or graph (`rg -ln 'json:"'` over
  `internal/table/`, `internal/graphlint/`, `internal/guard/` returns nothing).
  The approach builds genuinely new capability — no return to Stage 2/3 on
  reuse grounds.
- `guard.AssignmentCount` exists at `internal/guard/declaration.go` with shape
  `(int, bool)`, matching A3's and C2's usage.
- Four call sites duplicate `WriteJSONLine`'s encoder settings. None is on this
  document's path (the export is unimplemented), so it is not an A5 finding —
  noted as latent drift against C3's one-encoder rule.
- Evidence body: Testing Strategy is authored (9 scenarios, byte/set oracles).
  Performance Expectations is absent by a prior stage's deliberate omission,
  not an unfilled placeholder — no `_Draft placeholder._` survives anywhere.
- Baseline lint: `blocking=0 resolution=0`. The five `placeholder:survived`
  hits and `gate:inline` all sit in the Finalization Gate (lines 814–934),
  Stage 7's sections — not ones this stage owed authored.
