Model: claude-opus-5[1m]

# A13 / A15 closure — Stage 6, RDR 0002

Both assumptions were `Pending` on entry, both had a wide arm that **changes
match semantics while tripping zero load categories**, and both were closed
this pass by fixing the spike and minting the controls they were owed — rather
than by deferring the controls to implementation.

Why closed rather than downgraded: when a defect trips no load category, the
owed fixture is the *only* detector. Deferring the fixture defers the only
detector, and the fixture set is **promoted, not re-authored** at
implementation ("a fixture may not be narrowed on promotion, only extended"),
so the promoted set would have been the one that cannot discriminate.

## A13 — a member sequence is never compared as a joined string

**Defect, as found.** `Atom.identity` keyed the merge on
`literalString()` = `strings.Join(a.Literal, ",")`, and `renderWrites` stored
`Write.Value` as `renderValue(...)`, also comma-joined.

**Fix.**

| Site | Before | After |
| --- | --- | --- |
| `Atom.identity` | `… + literalString()` | `… + memberKey(a.Literal)` — length-prefixed (`%d:%s`), injective for any member content |
| `Write.Value` | `string` (joined) | `[]string` member sequence |
| `Write.identity` | *(did not exist)* | `Key + \x00 + memberKey(Value)` |
| `Write.render` | `strings.Join(Value, "|")` | `renderMembers` — each member `%q`-quoted |
| `literalString` | used for identity | display only, with a comment forbidding identity use |

**Controls** (`gen-cases.py` → `delim/`, 15 fixtures). Both arms parameterized
over `,` `;` `|` space and the empty string, per Testing Strategy 2.

- Atom arm: two contexts contribute `in = ["a<d>b","c"]` and `in = ["a","b<d>c"]`
  on one key+block → **two atoms** for every delimiter; the rule stays dead.
- Write arm: one set write spelled `["x<d>y","z"]` vs `["x","y<d>z"]` →
  **different values** for every delimiter.

**The parameterization earned its keep twice.**

1. The pre-fix spike joined on `,` and collapsed **only** the comma case. All
   four other delimiters yielded two atoms *while the defect was live*. A
   single-delimiter control would have passed it — which is precisely the
   argument Testing Strategy 2 makes for parameterizing.
2. The `|` case then caught a residual in the *render*: an unquoted `|`
   separator spelled `["x|y","z"]` and `["x","y|z"]` alike. The identity was
   already correct, but any oracle asserting on rendered text would have passed
   the collision. Fixed by quoting members; the display clause now requires it.

## A15 — atom-level validation is block-agnostic

**Defect, as found.** The guard walk checked only that the key was declared and
was not `recognized`; `checkMatchBlock` held every other atom rule. So
`frobnicate = "small"`, `eq = "<clear>"` and `eq = "NOT_A_PROFILE"` all loaded
clean under `[rule.guard.unless.profile]` and emitted a normalized atom.
Reproduction preserved at `a15-guard-block-repro.txt`.

Invisible by construction: `gen-cases.py` mutated match blocks only, so no
transcript line could witness it — which is why two lenses passed over it.

**Fix.** `checkMatchBlock` → `checkAtomBlock(m, owner, blockName, block)`; the
guard walk routes through it. Two rules stay legitimately match-only, gated on
`isMatch`: the `eq`/`in` restriction (§D6 routing) and `#` reservation in `in`
members (only match blocks expand). Guard blocks admit RDR 0003's closed set,
cited as `guardOperators`.

**A design refinement this forced.** "Domain/kind conformance" reads as
unqualified, but applying it literally to every operator breaks the RDR
fixture's own `iter.lt = 3`: a comparison operator takes an ordered **bound**,
not a domain member. Conformance is therefore per-operator —
`eq`/`in`/`contains` domain-checked, `exists` a bool literal,
`lt`/`lte`/`gt`/`gte` kind-checked only. The clause now says so.

**Controls.** Six guard-block negatives (37 → 43), one per atom rule per block;
each refuses exactly the category it names.

## Regression — nothing else moved

| Check | Result |
| --- | --- |
| Baseline dump digest | `6ccfe901…` — **unchanged**, 3 consecutive runs |
| Row census | 9 RDR + 2 kata — unchanged |
| Negative fixtures | 43/43 refuse one category each (41 category controls + 2 probes) |
| Positive fixtures | all load; both `-strict-only` checks `STRICT-OK` |
| Permutations | key-order, rule-order, `eq`-as-`in` all digest identically |
| Merge idempotence | one atom |
| Merge distinct-literal | two atoms; reversed `use` order digests identically |

The digest is unchanged because neither fixture carries a set literal or write
value whose members contain a delimiter — which is exactly why they could not
discriminate the defect, and exactly why the new controls had to be minted.

**Consequence for the SHA passage.** The RDR previously argued the recorded SHA
must not be a golden because "a *correct* implementation necessarily changes
these bytes." That premise is now false — the correct implementation reproduced
them. The ban stands, but on the oracle's character (a rendered-text hash
witnesses nothing about the identity rules, and breaks on legitimate rendering
changes), not on a predicted byte change. Corrected in Testing Strategy 2.
