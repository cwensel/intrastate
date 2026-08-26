Model: claude-opus-5

# cli/0011 — Stage 4 prior-art ledger: three-valued match at a selection seam

Stage: 4-resolve. Date: 2026-08-26.

Question. A match block is a set of key/value equalities that SELECTS which
rows apply, evaluated before the guard. With partial caller state a match key
may be (a) present and equal, (b) present and unequal, (c) absent, or (d)
present but supplied twice with differing values, so no single value can be
compared. Does prior art distinguish (c)/(d) from (b), at a SELECTION seam
rather than a condition seam?

The Stage-2 ledger (`prior-art.md`) recorded no coverage: every state-machine
peer assumes a current state is always present. This pass searched the
access-control and query literature, where partial input is the normal case.

## Accepted citations

- **PTaCL** (Policy Target Algebra for Access Control), read via Griesmayer &
  Morisset, "ATRAP: Automated Traceability of Access Request Policies"
  (PapersFast, `9A5BA134-9334-433E-89CD-71C072EB469E.pdf`) §2.1 3-valued Logic,
  §2.2 Target and Policy pp.552-553 incl. Table 1.

  A policy has a *target* — "specifies the requests to which the policy is
  applicable" — over a request that is a set of attribute name/value pairs.
  Structurally this codebase's match block over the assembled view.

      ⟦(n, v)⟧(q) = 1T  if (n, v′) ∈ q and v = v′
                    ⊥T  if (n, v′) ∉ q
                    0T  otherwise

  Absence is tested BEFORE inequality. And stated outright: "the target
  `(nat, FR)` evaluates to 1T (match) if the request contains `(nat, FR)`, to
  0T (no-match) if it contains `(nat, AT)`, but not `(nat, FR)`, and to ⊥T
  (indeterminate) if it does not contain any value for the attribute `nat`.
  **In other words, PTaCL can distinguish between a non-matching value for an
  attribute and a missing attribute.**"

  PTaCL also types the two domains apart — `{1T, 0T, ⊥T}` for target
  evaluation, `{1P, 0P, ⊥P}` for policy composition — "in order to avoid any
  confusion". ⇒ a selection-seam three-valued verdict kept distinct from the
  guard's is an established shape, not an invention.

  Note on case (d): PTaCL's atomic target is EXISTENTIAL over the request
  multiset, so `{(nat,FR),(nat,AT)}` matches both `(nat,FR)` and `(nat,AT)`
  (Table 1 row 4) rather than going ⊥T. That disposition is rejected here: it
  lets one request match mutually exclusive rows, which a transition table
  wants surfaced, not swallowed. It is cited as the one system that does NOT
  fold (d) into no-match either — it widens instead.

  PTaCL further supplies `optT`, "the optional target (i.e., transform ⊥T into
  0T)" — a deliberate, opt-in collapse of undecided into excluded. ⇒ precedent
  for a caller-side flag rather than a silent default collapse.

- **SAML 2.0 Core** (PapersFast, `2004-15.pdf`) §2.5.1 Conditions p.21 — the
  uncomparable case: "If an element contained within a `<Conditions>` element
  is encountered that is not understood, the status of the condition cannot be
  evaluated and the validity status of the assertion **MUST be considered to be
  Indeterminate**." ⇒ cannot-evaluate yields a third value, explicitly not
  false. This is the closest spec-level analogue to (d).

- **Cedar** (Cutler et al., arXiv 2024, PapersFast) §Related Work p.21 —
  XACML's Indeterminate as the ancestor vocabulary, and XACML targets as a
  policy-INDEXING seam ("labeled with an explicit set of target conditions that
  form the basis of policy indexing"). ⇒ corroborates that a target is a
  selection seam.

- **SQL three-valued logic** — Elmasri & Navathe, *Fundamentals of Database
  Systems* (DevRef, `db(11).pdf`) §4.1 p.121, §5.1 p.145 Table 5.1: "When a
  NULL is involved in a comparison operation, the result is considered to be
  UNKNOWN … SQL uses a three-valued logic with values TRUE, FALSE, and
  UNKNOWN." Celko, *SQL Database Programmers Handbook* (DevRef, `db(40).pdf`)
  §"The Null of It All" p.47: "A NULL is not a value; it is a marker for a
  value that is missing … even (NULL = NULL) is UNKNOWN." ⇒ the mass-adoption
  instance of (c), and the reason not to conflate unknown with false
  internally even where a boundary later filters on it.

## Verdict

No surveyed system folds present-but-uncomparable into a silent non-match.
The attested dispositions for (d) are Indeterminate (XACML/SAML lineage) or
existential widening (PTaCL's atomic target). Silent exclusion — the
pre-Stage-4 draft of `0011:C1` — has no prior art behind it.

Recommendation adopted in C1: `match` / `no-match` / `indeterminate` at the
match seam, `no-match` reserved strictly for present-and-unequal, with (c) and
(d) both indeterminate and distinguished by the reason `absent` /
`uncomparable` — the closed set `0007:C8` already owns, so no vocabulary is
minted.

## Rejected branches

- DMN hit policies / the `-` irrelevant marker: unevidenced here. The acronym
  collides with "default mode network" across the neuroscience papers and
  dominated retrieval in PapersFast. Unexplored, not refuted; PTaCL and the
  XACML lineage answer the question without it.
- XACML status-code URIs (`…status:missing-attribute`): the XACML spec itself
  is absent from these corpora. Every XACML claim above is sourced through
  PTaCL, Cedar, or SAML. Verify against the spec before putting a status-code
  URI in user-facing output — C1 does not.
- may/must analysis, abstract-interpretation lattices, content-based pub/sub
  matching: searched, no retrievable material on incomplete-input semantics.
- StateMachineLit / StateMachineRes on undecidable guards: guards are
  two-valued across the peer libraries; no unknown handling. Consistent with
  the Stage-2 finding.

## Queries run

| Corpus | Query | Outcome |
| --- | --- | --- |
| PapersFast (text) | `Indeterminate NotApplicable XACML` | accepted — surfaced SAML 2.0 Indeterminate |
| PapersFast (text) | `XACML policy target match Indeterminate attribute not found MustBePresent` | accepted — Cedar + PTaCL/ATRAP |
| PapersFast (text) | `PTaCL target evaluation three-valued logic match no-match indeterminate attribute missing` | accepted — decisive, §2.1/§2.2 |
| PapersFast (text) | `atomic target semantics attribute name value request contains multiple values indeterminate 1T 0T` | accepted — Table 1, the conflicting-duplicate row |
| PapersFast (text) | `the semantics of an atomic target for a request q is defined as follows 1T if 0T if bottom otherwise` | accepted — formal semantics + `optT` |
| DevRef (text) | `SQL NULL three-valued logic UNKNOWN comparison WHERE predicate neither true nor false` | accepted — Celko + Elmasri/Navathe |
| PapersFast (text) | `DMN decision table hit policy incomplete missing input value irrelevant` | rejected — acronym collision |
| PapersFast (text) | `may analysis must analysis unknown abstract value partial evaluation lattice three-valued` | rejected — noise |
| PapersFast (text) | `content-based subscription matching event missing attribute predicate undefined partial match` | rejected — pub/sub load balancing |
| StateMachineLit (text) | `enabled transitions guard cannot be evaluated unknown partial state three-valued` | rejected — no partial-state semantics |
| StateMachineRes (text) | `which transitions are possible from current state guard undecidable unknown variable` | rejected — two-valued guards |
| DevRef (text) | `decision table condition entry dash don't care irrelevant rule ambiguity completeness` | rejected — indexing/UML noise |
