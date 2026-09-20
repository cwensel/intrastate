Model: claude-fable-5-1

# 0030 — prior-art record (Stage 2 propose, 2026-09-19)

Read BEFORE the approaches were enumerated. Budget: ≤3 corpus queries per
claim, ≤5 opened hits per claim. A re-run reads this file instead of
re-searching.

## Problem class

Computed (non-literal) assignment to an extended-state variable on a
transition, in a table-as-data model with no host expression language.

## Accepted citations (load-bearing, quoted)

### Class read — every peer with extended state attaches host-language code to the action

- **Ragel POC review** (`../state-machines/contrast/poc-rdr-ragel/REVIEW.md`,
  §Faithfulness): "A 'loop at most 3 times' is **not a regular language
  property** — counting to a bound is exactly what a DFA cannot do without N
  distinct states. I had two honest choices: (a) unroll the loop into three
  explicit `prelock1/2/3` states (pure-Ragel, but the cap becomes hard-wired
  topology), or (b) keep one `prelock` self-loop and count in an action
  variable. I chose (b) … **The bound is enforced in C, not in the regular
  grammar.**" and "Ragel handled edges and cycles for free; it handled guards,
  counters, bounds, and the coupled status register only as host-language
  action side-effects."
  ⇒ names the exact fork this record decides: unroll into N rows (what
  intrastate forces today) vs. count in an action (what every peer does).
  The POC also shows the cost of (b): "the .rl is no longer the single source
  of truth" — the drift intrastate's data-only table exists to prevent.
- **sismic** (`../state-machines/study/sismic/docs/code.rst`, §Context of the
  Python code evaluator): `on entry: x += 1` executed in a Python context;
  `{'x': 1, 'y': 0}` → `{'x': 2, 'y': 0}`. Cited by
  `../state-machines/MODEL-transition.md` §3a as the EFSM persist-unless-assigned
  convention. ⇒ the class answer is a general expression language over an
  execution context; intrastate has none by RDR 0002's charter.
- **SCXML** (`<assign>`, `<datamodel>`) — the standard's assignment element
  evaluates an `expr` in the datamodel's script language (uscxml ships
  ECMAScript/Jexl datamodels: `../state-machines/study/uscxml/README.md`;
  scxmlcc lists `<assign>` among supported tags, its manual entry is `todo`:
  `../state-machines/repos/scxmlcc/doc/user-manual.md` §Assign). ⇒ same class
  answer; no data-only assignment form in the standard.

### Instance read — what does each peer do for increment / enum-successor?

| Peer | Increment | Enum successor | Bound behaviour |
| --- | --- | --- | --- |
| sismic | `x += 1` (Python) | Python code | host language |
| SCXML/uscxml | `<assign expr="x+1">` | script | host language |
| Ragel POC | `ctx->prelock_iter++` in a C action | n/a | `if` in C, surfaced as a flag |
| stateless / transitions | callbacks | callbacks | host language |

No peer offers a data-only computed write; none decides the bound
declaratively. ⇒ the bound disposition has no external prior art and rests
on in-repo principle: `internal/table/load.go::valueMembers` comment quoting
`0002:C3`'s governing principle — "a malformed authoring is a stable refusal,
never a silent no-op".

### In-repo prior art the choice rests on

- `internal/table/normalize.go::expand` — the loader already mints one row
  per member for a match-block `in` atom, with the chosen member as the
  row-identity suffix (`0002:C13`). Its comment explains why guard `in`
  atoms are deliberately NOT expanded (an `unless` conjunction would become
  a disjunction; `#` is authorable in guard members).
- `docs/model-authoring.md` §Discriminate with guard atoms: "A **match** atom
  *scopes* the row group … and contributes **no dimension** to the product
  lint proves exhaustiveness over. Only `guard.all` and `guard.unless` atoms
  collect as dimensions." ⇒ an expansion that must make the bound cell
  provable has to emit guard atoms, not match atoms.
- `internal/table/load.go::conform` / `conformDomain` — a literal write past
  `max` or outside `domain` is refused at load under
  `malformed_tag_declaration` (renderWrites comment: "filed under the
  declaration category (deviations.md D3)"). ⇒ the bound refusal on an
  expanded literal is the same check, same category.
- `internal/table/load.go::tagDecl` carries `Domain: src.Domain` unsorted;
  `0021:C2`: "`domain` carries the declaration's AUTHORED members". ⇒
  authored enum order is available to carry step order.
- `internal/table/model.go::Identity`: "Every string that can reach here has
  `#` banned at load". ⇒ a domain member reaching the suffix needs the same
  guarantee (assumption A3).

## Queries run

- arc `StateMachineRes`: "transition action assigns a computed value:
  increment a counter or step to the next enum member instead of a literal"
  → xstate-like runtime docs, no assignment grammar. Rejected.
- arc `StateMachineRes`: "SCXML assign element location expr datamodel update
  variable on transition" → state-machine-cat SCXML docs (actions as opaque
  text). Rejected as non-load-bearing.
- arc `StateMachineRes`: "assign location expr: the assign element modifies
  the data model…" → scxmlcc user-manual (§Assign `todo`). Accepted as the
  "supported tag" cite only.
- arc `StateMachineLit`: "TLC computes successor states by evaluating the
  next-state action …" → TraceFix/StateFlow papers, score ≤0.55. No hit;
  ⚠ no literature coverage for analysis-time successor enumeration — the
  chosen approach does not rest on it (it rests on `expand`, in-repo).
- Sibling repo (scoped rg over `study/*/README.md`, `research/*.md`,
  `MODEL-transition.md`, `contrast/poc-rdr-ragel/REVIEW.md`): accepted the
  Ragel and sismic passages above.

## Rejected branches

- TLA+ next-state-relation analogy for a symbolic successor: not found in
  corpus at a load-bearing score; and the chosen approach needs no symbolic
  successor.
- Ada `'Succ` / Haskell `succ maxBound` as prior art for "error at the last
  member": not in any corpus; would be a model-prior citation, so it is NOT
  used — the bound disposition cites `0002:C3`'s principle instead.
