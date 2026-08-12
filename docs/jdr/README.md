# Joint Decision Registries (JDRs)

A **JDR** is the single normative home for decisions that two or more sibling
RDRs *jointly* own — decisions no one RDR solely owns, which therefore cannot
live inside any one of them.

Prototyped here in `intrastate`; intended to be pushed down into the RDR engine
(`$RDR_HOME`) as a skill once the shape settles.

## Why this exists

The cluster gate (`/rdr-cluster-reconcile`) types a cross-RDR finding one of three ways. The
third, **JOINT-DECISION**, is explicitly *not* a spec defect: the siblings
contradict over a decision that is jointly theirs. The prescribed cure is to
**hoist** it to a single normative home that the siblings **cite** rather than
restate, and let them proceed under recorded tolerance.

Restating a shared decision in each sibling is what makes them drift. The RDR
engine states the general law:

> Two copies of one contract drift into self-contradiction — the most common
> internal-defect class — so the cure is deletion, not a lint to keep them in
> sync.
> — `$RDR_HOME/stages/README.md`, *Single-source each contract*

The failure mode is concrete: two RDRs each carry a copy of one shared rule,
their re-locks land a day apart, and the copies contradict. Every gate passes,
because each document is internally consistent. An implementer follows the wrong
copy and rebuilds the defect with a green conscience.

## What a JDR is *not*

- **Not an RDR.** An RDR is one decision, with alternatives, assumptions,
  lenses, and a finalize gate; it locks. A JDR is a **registry of many entries**
  with independent per-entry lifecycle. Seeding an umbrella *RDR* to fix
  cross-RDR drift reproduces the drift one level up — that is precisely what
  RDRs 0007/0008/0009 did, each declaring itself "the single normative home" of
  a seam rule, and the 0002-0009 gate returned 24 blocking findings.
- **Not a decision log.** In ADR practice "decision log" already means the whole
  directory of records. A JDR is a curated subset with normative force.

## Identity — topic ∧ cluster, never the gate run

**A JDR is one question that a set of RDRs jointly answer.** Two things fix its
identity, and both are load-bearing:

- **The cluster** — *who* must agree. The set of RDRs that share the decisions.
- **The topic** — *what* they are agreeing about. A JDR's title is a question;
  every entry is a fragment of its answer.

Neither alone is sufficient, and the failure mode of dropping either is
concrete.

**Drop the topic and you get a ledger** — a junk drawer that accretes every
cross-RDR item touching a popular RDR, with no principle for resolving any of
them. Entries stop informing each other.

**Key on "the identical RDR set" and you get one JDR per entry.** This is
tempting — it sounds like the crispest possible rule — but it collapses to
ADRs. Check it against real data: in this repo's `JDR 0001` the ten entries span
ten *different* subsets (`JD-1` is 0007×0003, `JD-10` is 0008×0002×0009, `JD-4`
reaches five members). **Entries in a cohesive registry share a topic, not a
membership list.** The cluster is the outer boundary of who may be bound;
individual entries bind subsets of it.

So: entries are topically cohesive *because* they answer one question — and the
cluster is the union of everyone that question binds, not a set each entry must
match exactly.

### Routing a new entry

Ask the two questions in order:

1. **Is it the same question?** If yes, it belongs to that JDR — even if it
   binds a subset of the members, and even if it drags in a member the cluster
   did not previously list (widen the `cluster` frontmatter and say so).
2. **Is it a different question over the same RDRs?** Then it is a *second*
   JDR. Same cast, different play. Two JDRs may share members; that is normal
   and costs only a cross-reference.

Consequences:

1. **Re-running the gate on the same set appends; it does not create.** Entries
   accrete across iterations; mark later ones with the iteration that found them.
2. **A subset invocation still lands in the same JDR.** Running the gate over
   `7 8 9` when the cluster is `0002…0009` finds the same shared decisions;
   they have one home.
3. **Two genuinely disjoint clusters get two JDRs.** If a later set of RDRs
   shares decisions among themselves and not with this set, it gets its own.
4. **A new RDR landing on an existing seam joins that seam's JDR** — if it joins
   the same question. It does not get a JDR of its own just for being new, nor
   does it join merely for touching a listed member.

**Compound titles are a smell worth checking.** If the topic question needs an
"and" to state, test whether the halves are causally chained (one JDR) or merely
adjacent (two). `JDR 0001` here is compound — trust *and* surfacing — and stays
one registry because the halves are one pipeline: what the resolver cannot trust
determines what it refuses, which determines what the user sees. Adjacent topics
would not survive that test.

### Why identity is not the "connected component"

A tempting rule is "one JDR per connected component of the shares-a-decision
graph." Reject it: connected components **merge** as edges are added, so a later
decision spanning two JDRs would retroactively fuse them — and merging documents
breaks citations. Citation stability is the overriding constraint: once an RDR
records `→ JDR 0001 §JD-5`, that anchor must never move.

Therefore JDRs are **append-only and never merge or split after seeding**. Seed
conservatively — one per cluster. If a later decision genuinely spans two JDRs,
it lands in the JDR owning the seam it sits on, and the other JDR
cross-references it. Two JDRs cross-referencing is cheap; a moved anchor is not.

### Naming

`docs/jdr/NNNN-<seam-slug>.md`, zero-padded, sequential. The slug names the
**seam or theme**, never the member list — membership grows, filenames must not
rot. Members are declared in frontmatter.

Must live **outside** `$RDR_RECORDS` (`docs/rdr`): `§rdr-resolve` globs
`$RDR_RECORDS/NNNN-*.md` non-recursively, so a numbered file placed there would
be silently resolvable as an RDR. `docs/jdr/` is invisible to that glob.

## Citation

Siblings cite `JDR NNNN §JD-n` — anchor only, never restated mechanism prose.
The cluster gate treats a restatement as the defect. Sub-clause and prose
anchors follow the same form: `JDR 0001 §D2`, `JDR 0001 §Principles`.

An RDR tolerated under a joint decision carries the TEMPLATE qualifier on its
live Status value:

```
**Status**: Final [joint decision → JDR 0001 §JD-5]
```

Per `$RDR_HOME/TEMPLATE.md` this stays `Final` for every binary gate and **does
not self-clear** — the home owns the decision, so the qualifier is permanent
until the home says otherwise.

**Stamp the qualifier when the entry resolves, not when it is filed.** The
qualifier asserts a sibling may proceed under tolerance; asserting that over an
entry still `open` would license implementation across an undecided contract.

## Entry lifecycle

Per entry, not per document:

| Status | Meaning |
| --- | --- |
| `open` | Identified, not yet decided. Blocks the siblings it spans. |
| `constraint` | The answer is already determined by shipped code or a stated principle. Nothing to negotiate; record it so nobody implements against it. Does not block. |
| `blank` | A detail the implementer fills (a code string, a test's home). Names an owner, not a negotiation. Does not block. |
| `deferred → <phase>` | Cannot honestly be settled on paper; the named implementation phase settles it. Does not block seeding that phase; does block claiming the seam is closed. |
| `decided` | Resolved here. Siblings may carry the tolerance qualifier. |
| `withdrawn` | Re-triaged as a single-RDR defect or a non-finding; kept with its reason so the anchor never dangles. |

Never delete an entry — withdraw it. Deleting breaks citations.

### Grade entries, or the registry over-gates

The cluster gate grades findings by severity (`blocks-impl` / `risks-impl` /
`cosmetic`). That axis asks *must an implementer confront this?* — true of a
thousand micro-decisions, so `blocks-impl` inflates. A registry that copies
severity straight into `open` will hold a cluster at NOT RECONCILED over
bookkeeping.

Grade every entry on a second axis before filing it:

- **Fork** → `open`. A genuine either/or with divergent consequences; wrong
  choice costs rework or ships a defect. **These gate implementation.**
- **Constraint** → `constraint`. Shipped code or a stated principle already
  determines it. Write it down; do not negotiate it.
- **Blank** → `blank`. The implementer fills it. Name the owner.

Check candidate entries against **shipped code, not sibling prose** — a gap
between two RDRs often closes the moment you read the package they both
describe. In `JDR 0001`, three of ten entries survived that check as forks.

The inverse error is real too: grade by *consequence*, never by how hard the
finding was to spot or how cheap the fix looks. An entry that is one line to fix
and silently unsound if skipped is a fork.

## Document state

`open` — entries still being decided · `settled` — every entry `decided` or
`deferred` · `superseded` — replaced (name the successor).

## Structure

1. YAML frontmatter — `authors`, `state`, `cluster` (member RDRs), `labels`
2. `# JDR NNNN <Title>` — a **question**, not a noun phrase; plus a
   non-normative note
3. `## Problem statement` — what the members jointly decide
4. `## Principles` — **required.** See below.
5. `## D1…Dn` — the genuine forks, as `(a)/(b)/(c)` with trade-offs and a
   recommendation; each closes with a bolded **Resolved:** naming **which RDRs
   the decision lands in and what changes there**
6. `## Interface record` — the `JD-n` entries, each stating its answer, the
   user-visible stake, and a short provenance parenthetical
7. `## What this does not decide` — local items that stay with their RDRs

**No work plan.** A JDR records decisions, not the sequence for executing them.
A "next steps" section is stale the moment the first RDR re-locks, nothing
updates it, and it duplicates what the RDR `Status` lines and `/rdr-status`
already say authoritatively — the same two-copies-drift failure this document
type exists to prevent, turned on the document itself. A section that must be
*emptied* before the JDR can reach `settled` is one nobody remembers to empty,
and half-true is worse than absent. Each **Resolved:** names its landing RDRs;
that is the durable half, and it stays true after the work is done.

**Keep instances short.** Rationale and decisions only. Convention lives here in
this README and is never repeated in an instance; process history lives in git
and is never narrated in the document. The most consequential questions get the
most prominent placement — a real fork buried in a flat list while a small
question holds a top-level section is the defect that makes a registry unusable.

### Guiding principles — what makes a JDR steer instead of merely record

A registry of independently-argued entries drifts exactly like the sibling RDRs
it replaced: ten locally-reasonable resolutions that do not add up to one
coherent behavior. The cure is a small set of standing commitments, stated
before the entries, that every resolution must satisfy.

Rules for them:

- **Derive, never invent.** Each principle must be traceable to normative text
  the members already lock, or to a decision resolved in this JDR. Quote the
  source. A principle with no anchor is the author legislating.
- **Make them decisive.** A principle earns its place by ruling options *out*.
  If every option on a fork satisfies it, it is decoration.
- **One line each.** A principle that needs a paragraph is an argument, and
  arguments belong in the decision that uses it.
- **A resolved fork may become a principle.** Resolutions compound rather than
  accumulating: once a fork fixes a promise at a strength, later entries bound
  themselves by it and cite it rather than re-arguing.

This is the section that answers "is a JDR a ledger?" — no. A ledger records; a
JDR records *and* constrains how its open entries may be closed.

## Ancestry

No canonical name exists for this artifact; it combines two long-standing
traditions.

- **Interface Control Documents** and their joint-ownership working groups
  (MIL-STD-490A; NASA SP-2016-6105) — a boundary jointly owned by parties whose
  own specs must not restate it.
- **IANA registries** (RFC 8126) — numbered entries, per-entry lifecycle, and
  specifications that cite entries rather than restating them.
- **arc42 §8 Cross-cutting Concepts** — rules stated once so sibling sections
  need not repeat them.
- **ViewPoints inter-viewpoint rules** (Finkelstein/Nuseibeh 1994; Easterbrook
  1996) and **boundary objects** (Star & Griesemer 1989) — the stance that
  cross-view inconsistency is *managed* at chosen checkpoints rather than
  eliminated. Stage 7.1 already rests on this literature.

Rejected names: `ICD`/`IRD` (aerospace-only recognition; `.icd` is an OpenCL
driver config), `ADR` (one-per-file, and collides with RDR), `SDR` (System
Design Review), `decision log` (inverts the established ADR meaning).
