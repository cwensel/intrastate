Model: claude-opus-5[1m]

# Stage 6 Reconcile — RDR 0008 (recognized-tag key ownership)

Verdict: **RECONCILED**. Thirteen assumptions terminal, no BLOCKER. One
assumption (A12) refuted as stated and downgraded with its repair landed in
the RDR as an explicit Prerequisite.

## Stage 5 preflight

- `Profile: foundational` → required lenses `cove 3amigo critique
  repeatability`. All four evidence dirs present and resolved.
- Determinacy / repeatability: `full` variant as the profile requires (not
  `lite`), three runs with three **distinct** model stamps —
  `claude-opus-5[1m]`, `claude-fable-5`, `glm-5.2:cloud` — plus `diff.md` and
  `resolve.md`. No variant mismatch.
- **Caveat (not a gate failure)**: the `critique` lens ran single-model
  (Pass A + Pass B, both `claude-opus-5[1m]`, two fresh contexts) as a
  documented fallback — `resources.md` records "Alt-model roster: omitted
  (single model)". `foundational` prefers dual-model convergence; the
  recorded fallback stands in. Repeatability got the genuine multi-model
  draw the critique lens could not.

## Open set — four sources

1. **Pre-Lock needs-verification lists** (3amigo, critique, repeatability,
   cove): A7, A9, A10, A11, A12; plus A8/A13 booked Verified at authoring
   and flagged for confirmation here.
2. **Still-Pending assumptions**: the same five — A7, A9, A10, A11, A12.
3. **Named-but-unrun spikes**: A9's plan named "a normalization spike over
   the committed fixtures" as an alternative to a peer read; `{SPIKE_DIR}`
   did not exist. **Run at this stage** — see `spikes/a9-normalization.md`.
4. **Exactness-word delta**: the round-introduced exactness claims are the
   byte-exact identity rule, the direction-specific payload tokens, the
   disjunctive advisory trigger, and the first-breach arity. Each is a rule
   this RDR *authors* rather than a fact about an external system, and each
   is pinned by a named scenario (4, 6, 7) — the falsifiability this RDR
   requires of its own MUSTs. No unbacked exactness claim survives.
   A13's `[]Tag` shape (the one round-introduced *source* claim) re-confirmed
   at `internal/resolve/resolve.go::Input` / `::Tag`: `Owned []Tag`,
   `Observed []Tag`, `Tag{Key, Value string}` — unmoved.

Absorption audit over all four lens dirs: **zero unabsorbed defect residue**.
Every finding is fixed, dismissed-with-cite, or charted. Surviving residue is
forward-looking only (the five assumptions above, two charted successors,
two cluster-reconcile watch items).

## Dispositions

| Item | Source | Disposition | Evidence / plan |
|---|---|---|---|
| **A7** — `RequiresOwned` name reservation vs RDR 0007's ownership | 1,2 | **VERIFIED** | Ownership: RDR 0007 (`Final`) scopes "single normative home" to the *guard-decidability domain rule*; its `RequiresOwned` block is wholly semantic. Decisive — 0007 *itself* cites this RDR's rule as a peer's **name** constraint (`0007…md:951-955`), drawing the name/meaning line and placing this clause on the name side. Reachability: `::missingOwned` → `::TagSet.has` provenance-specific test → `KindOwnedStateUnavailable`. Still pinned by scenario 9. |
| **A9** — which normalized field a tag predicate produces | 1,2,3 | **VERIFIED** (spike) | Re-ran RDR 0002's own normalizer over all fixtures: `evidence/spikes/a9-normalization.md`. Normalized `Row` has **no outcome field**; `[rule.match.<tag>]` lands unprefixed in the predicate set; `rule.Guard.All`/`Unless` merge into the same set with `all:`/`unless:` prefixes. First reading holds → Phase 2 rename is mechanical and ungated. |
| **A10** — accessor-derived key spelling `recognized` | 1,2 | **VERIFIED** | Accessor binds to a *declared tag*, never names a key: 0002's reference is one-way (tag → accessor), the spike's `Accessor` type is `{Mode, Path}` with no key field, and RDR 0004 (`Final`) declares "expected tag keys" with identity `(flow, name, capability)`, a read accessor returning "typed tag values". Data channel carries no key-spelling authority → block 4's programmer-mistake channel is correct. Scope: contractual (accessor layer unimplemented at HEAD). |
| **A11** — composition with RDR 0009's entry precondition | 1,2 | **VERIFIED** (independence) | Disjoint fields: this RDR reads `Input.Owned`/`Observed` + each row's `RequiresOwned`; 0009's predicate is exactly `len(row.Escape) != 0 && len(row.Writes) != 0`. Both pure reads, no shared field → neither can suppress the other's verdict. `Input` occurs **zero** times in 0009 (re-counted). Report order is a cross-RDR *decision*, not a lookup; block 4 carries a conforming interim rule. Stage 7.1 item by design. |
| **A12** — 0002 implementer encounters the constraint | 1,2 | **DOWNGRADED** (refuted as stated) | Mechanism does not exist: `launch.md:14-15` implements one RDR and puts "Cross-RDR orchestration … out of scope"; `Predecessors` points backward only; `grep -rn "Overrides"` over the engine returns **no matches**; RDR 0002 carries **zero** `0008` occurrences. Repair adopted and landed: an explicit **Prerequisite** requiring the carried half be handed to 0002 (named Stage 8 input, or a pointer edited into 0002). Form is a Stage 7.1 call; the obligation is settled here. |
| **A13** — `Input.Owned`/`Observed` are slices, not maps | 1 | **VERIFIED** (re-confirmed) | `::Input` `Owned []Tag` / `Observed []Tag`, `::Tag {Key, Value string}`; `::assemble` consumes them as sequences. Shape has not moved under RDR 0001. |
| **A8** — first non-nil error path | 1 | **VERIFIED** (stands) | Confirmed unchanged; A10 closes the scope limit A8 declared for the data channel. |
| Charted C-1 (intent heuristic), C-2 (user-facing surface) | audit | **out of scope, charted** | Successor RDR / Stage 7.1 watch item. Boundary stated in Failure Modes and Trade-offs, not silently omitted. |
| Watch: 0007 A13 negative-existential | audit | **cluster-reconcile** | Falsifiable by a *future* peer, not by this RDR. |

## Hard rules

- **Refutation → BLOCKER?** No. A12 is the only refuted item, and it asserted
  a *discoverability mechanism*, not a design premise. No normative contract
  changes: the rule, category, payload, and predicate all stand as written.
  What it changes is enforcement *reach*, and that is repaired in-document by
  a Prerequisite rather than papered over. No return stage is owed.
- **MVV-critical deferral?** None. The kernel half runs during this RDR's
  implementation. The normalizer half is *scoped to the peer owning the code*
  — explicitly not deferred by choice — and A12's disposition **strengthens**
  that gate by making the handoff a written prerequisite instead of an
  assumed traversal.

## Completeness

- No `_Draft placeholder._` and no seed-skeleton header (grep: 0).
- `## References` is filled with real citations; no template brackets.
- Remaining bracketed text is the **Finalization Gate**'s own response
  prompts — Stage 7's to answer, correctly unfilled at Stage 6.

## RDR edits made by this stage

A7/A9/A10/A11 flipped `Pending` → `Verified` with evidence written in place;
A12 recorded as refuted-as-stated with its repair; and every dependent
passage reconciled so no settled-fact prose rests on a Pending assumption —
Approach item 5, Capability Dependencies, A1's scope limit, A4's
sub-question, Risks/Mitigations, Failure Modes, Phase 2, scenario 9, and
Prerequisites (new gate item).
