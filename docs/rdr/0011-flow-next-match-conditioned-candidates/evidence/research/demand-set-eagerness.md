Model: claude-opus-5

# cli/0011 — Stage 4 prior-art ledger: demand-set eagerness at a reader seam

Stage: 4-resolve. Date: 2026-08-26.

Question. The demand set — which readers to invoke — is computed statically as a
UNION over all CANDIDATE rows for the requested outcome, before evaluation. So a
reader runs even when the row ultimately selected does not reference its key, and
an unbound artifact behind that reader fails the whole invocation (exit 2) even
though the selected row needed nothing from it.

- Option A (eager/total): union demand set. Behaviour is a function of the RULE
  TABLE plus the requested outcome, not of which row happens to win.
- Option B (lazy/scoped): fetch per selected branch. Never fails on a reader the
  winner did not need.

## Accepted citations

### Angle 1 — access control: Cedar, and the split it makes

- **Cedar** (Cutler et al., "Cedar: A New Language for Expressive, Fast, Safe,
  and Analyzable Authorization (Extended Version)", arXiv 2024; PapersFast,
  `arXiv-2024-Cedar-…_1.pdf`). Cedar makes the split this RDR is arguing about,
  and makes it EXPLICITLY, in two places.

  §1 p.2 — the posture, stated as a design goal:

  > "While deciding a request considers **all active policies**, Cedar is able
  > to quickly **slice** the full policy set to evaluate only those relevant to
  > the request."

  §2.4 p.5 — the correctness condition that licenses the slice:

  > "…the same outcome as if it had considered all available policies."

  §3.3 p.8/10 — the slicing mechanism, which is a STATIC index on the policy
  SCOPE (principal/resource), not on the policy conditions:

  > "Let pof(c) be the entity E::s named in the Principal portion of policy c, or
  > a special identifier Any if no entity is named… Define the key for c to be
  > ⟨pof(c), rof(c)⟩… The slice of relevant policies is then ⋃_{k∈K} Σ(k). Other
  > schemes are possible too, **e.g., which take into account the policy
  > conditions**."

  §3.5 p.11 states the proved obligation: `Sound slicing: ∀c∈C. c∉slice(C,σ)
  implies c cannot satisfy request σ.` — an OVER-approximation is sound; an
  under-approximation is not.

  ⇒ **Supports A, with a named refinement.** Cedar narrows by SCOPE (the
  selection seam — this codebase's `match`) and explicitly declines to narrow by
  CONDITION (the guard seam). It fetches for the whole scope-slice and evaluates
  every member. The union is over CANDIDATES, and candidacy is decided by the
  match keys only — which is exactly `0011`'s conditioned-candidate set. Cedar
  ratifies the shape rather than the alternative.

- **Cedar totality / errors are not decisions** (same paper, §2.5 p.5, §3.2 p.7,
  and n.2 p.10). Cedar's answer to "the data needed isn't there" is neither a
  lazy skip nor a whole-request abort: a missing entity or attribute makes
  evaluation *stuck* ("Stuck states in the semantics correspond to raised errors
  in the implementation", §3.2), the offending policy is dropped from the
  decision, and the failure is surfaced as DIAGNOSTICS alongside the decision —
  "the Allow/Deny decision is coupled with extra diagnostics, which include…
  any policies that exhibited errors when evaluated" (n.2). Errors are moved OUT
  of the decision channel and into a separate reporting channel, and the missing
  data is instead ruled out ahead of time by the schema validator (§2.5):
  "ensuring that Cedar policy evaluation will not result in an error, e.g., if
  the policy attempts to access a non-existent attribute."

  ⇒ **Third posture, neither A nor B.** Cedar prevents the situation statically
  (validator + schema) rather than choosing between failing eagerly and skipping
  lazily. The nearest transfer to `0011`: the unbound-artifact case is arguably a
  CONFIGURATION defect to catch at load, not an evaluation-time outcome.

- **XACML** — the classic PDP/PIP attribute-retrieval axis. **Not evidenced
  here.** See Rejected branches; the sibling ledger `three-valued-match.md`
  recorded the same absence.

### Angle 2 — eager vs lazy in general engineering

- **Fail Fast** — Wolff, *Microservices: Flexible Software Architecture* (DevRef,
  `microservices-flexible-software-architecture.pdf` p.240):

  > "Each system should recognize errors as quickly as possible and indicate them
  > immediately. When a call requires a certain service and that service is
  > unavailable at the moment, the call can be directly answered with an error
  > message. **The same is true when other resources are not available at the
  > time. Also, a call should be validated right at the start. When it contains
  > errors, there is nothing to be gained by processing** [it further]."

  ⇒ **Supports A**, and is the only named pattern found that argues the eager
  side directly. Note its scope: it justifies failing early on a resource the
  call *requires* — it does not speak to a resource the call turns out not to
  need. That is the precise gap in this RDR.

- **Non-strict / short-circuit evaluation** — the counter-argument, well attested
  and unusually explicit that the point is AVOIDING AN ERROR, not saving work:
  - Friedman & Wand, *Essentials of Programming Languages* 3e (DevRef, p.159 /
    `2003 Daniel P Friedman_1.pdf` p.137): "Sometimes in a given call a procedure
    never refers to one or more of its formal parameters. In this case time
    devoted to evaluating the corresponding operands is wasted. **It may even be
    that evaluation of such an operand would result in an error or never
    terminate.** For example, were it not for such problems, `if` could be a
    procedure, instead of having to be a syntactic form."
  - Meyer, *Object-Oriented Software Construction* (DevRef, p.477): the
    `and then` / `or else` non-strict operators "may yield a value (false) in
    cases when the `and` form does not… the non-commutative operators may be said
    to be **'more defined than or equal to'** their respective counterparts."
  - McConnell, *Code Complete* (DevRef, pp.475-476): short-circuiting is what
    makes `while (i < MAX) if (item[i] <> 0)` correct at all.

  ⇒ **Supports B**, and supplies the sharpest framing available: lazy evaluation
  is *strictly more defined* than eager. Option A is not merely slower — it is
  partial where B is total, on the same inputs. This inverts the intuition that
  the eager path is the "total" one.

- **negative:** No source was found arguing that a static, over-approximating
  demand set is preferable BECAUSE it makes behaviour a function of the query
  rather than of the data. DDIA's declarative-query material (Kleppmann, DevRef
  pp.42, 449) says only that declarative languages let the optimizer choose the
  *how* — which is an argument for optimizer FREEDOM (hence for pruning), not
  against it. PostgreSQL's planner docs (DevRef, `postgresql-1{3..18}-US.pdf`
  "Overview of PostgreSQL Internals") describe projection/selection pruning as
  standard, i.e. planners over-approximate only as far as the analysis is
  imprecise, and prune wherever they can prove they may. **The determinism
  argument for A is this codebase's own; the literature does not supply it.**

### Angle 3 — CLI prior art: OpenTofu splits the two postures

OpenTofu is the closest peer instance, and it lands on BOTH sides — deliberately,
at two different layers. This is the most useful finding in the ledger.

- **Eager/total at the PLUGIN-PRESENCE layer.**
  `../langref/opentofu/internal/tofu/context.go:327-362`,
  `checkConfigDependencies`, iterates `config.ProviderRequirements()` — the
  requirements of the WHOLE configuration — and errors per missing plugin:

      "This configuration requires provider %s, but that provider isn't
      available. You may be able to install it automatically by running:
        tofu init"

  Its doc comment names the motive as error quality, not correctness: "can avoid
  some potentially-more-confusing errors from later operations."

  It is called from `context_plan.go:163`, and the caller bails on the spot:

      // If required dependencies are not available then we'll bail early since
      // otherwise we're likely to just see a bunch of other errors related to
      // incompatibilities, which could be overwhelming for the user.

  Crucially this runs BEFORE the graph is built. `PlanOpts.Targets` is read only
  later, by `TargetingTransformer` at
  `internal/tofu/graph_builder_plan.go:260`. ⇒ **`tofu plan -target=X` fails on a
  provider that X does not use.** This is a direct, exact instance of Option A's
  cost: a hard failure on a resource the selected work did not need.

- **Lazy/scoped at the PROVIDER-CONFIGURATION layer.** In the same graph builder,
  `&PruneProviderTransformer{}` at `graph_builder_plan.go:229` ("Remove unused
  providers and proxies") and `&pruneUnusedNodesTransformer{}` at :255 both run
  before `TargetingTransformer` at :260. A provider whose resources are all
  pruned or untargeted is never *configured* — never handed its credentials,
  never asked to connect. ⇒ For the expensive, failure-prone step, OpenTofu is
  Option B.

  ⇒ **The line OpenTofu draws is between static availability and dynamic
  acquisition.** Cheap, local, statically-known facts (is the plugin on disk) are
  demanded eagerly over the whole config. Expensive, remote, credential-bearing
  acquisition is demanded lazily per surviving node. Read across to `0011`: the
  question is not "eager or lazy" in the abstract but which side of that line a
  READER sits on.

- **helm** — checked, rejected as a weaker analogue: dependency and values
  resolution is chart-wide by construction (a chart's `Chart.lock` is resolved
  before any template renders), so there is no per-template branch on which the
  posture could differ. It has nothing to say about the winner-didn't-need-it
  case.

## Verdict

**Genuinely split — but not evenly, and the split is informative rather than a
coin-flip.**

1. Cedar (the closest structural analogue) supports the union-over-candidates
   shape, and independently confirms the level at which to draw the candidate
   set: narrow by SELECTION scope (match keys), do not attempt to narrow by
   CONDITION. `0011`'s conditioned-candidate set is exactly Cedar's slice.
2. Language semantics supply the strongest argument against A, and it is sharper
   than "lazy is faster": non-strict evaluation is *more defined* than strict
   (Meyer). Option A is partial where Option B is total.
3. OpenTofu shows a mature peer CLI doing BOTH, split by cost and locality of the
   fetch, and paying exactly this RDR's cost at the eager layer:
   `plan -target` fails on an uninstalled provider the target does not use. That
   this is tolerated in a tool of that scale is real evidence that the cost is
   acceptable — but note OpenTofu's own comments justify it by ERROR QUALITY, not
   by determinism or caller-independence.
4. **negative:** the determinism/caller-independence argument for the eager union
   — "behaviour should be a function of the query, not the data" — has no
   support in these corpora. If `0011` rests on it, it rests on its own
   reasoning, and should say so rather than cite prior art.

Suggested reading for the tie-breaker: prior art does not settle A vs B, but it
does say the axis is COST AND LOCALITY of the fetch, not eagerness as a
principle. A reader that is cheap, local, and statically checkable belongs in the
eager union (OpenTofu's plugin check, Cedar's validator). A reader that performs
remote or credentialed acquisition is where every surveyed system goes lazy. If
`0011` keeps the eager union, the OpenTofu precedent is the honest citation, and
Cedar's diagnostics channel (n.2 p.10) is the honest mitigation — report the
unbound reader without necessarily making it exit 2.

## Rejected branches

- **XACML spec, PIP attribute retrieval, `Indeterminate` scoping.** The OASIS
  XACML specification is absent from all four corpora — confirmed again here
  (queries return SAML profile documents that merely *reference* XACML, plus
  unrelated policy-architecture papers). Every XACML claim reachable from these
  corpora is second-hand through Cedar §6 or SAML. The question "is a missing PIP
  a whole-request Indeterminate or scoped to the rule that needed it" — the
  single most on-point question in the brief — **could not be answered.** Verify
  against the spec directly before citing XACML either way.
- **PTaCL / ATRAP** (PapersFast, `9A5BA134-…pdf`): re-surfaced, but its content
  is about three-valued *target evaluation* and attribute-hiding attacks, not
  about when attributes are fetched. Already fully mined by
  `three-valued-match.md`; nothing further for this question.
- **OPA/Rego partial evaluation.** Cedar §5.2 p.19 discusses Rego's scaling but
  not its partial-evaluation feature; no Rego source in these corpora. Rego's
  `partial` API is very likely relevant (it is the named instance of evaluating
  with incomplete input) but is **unexplored here, not refuted.**
- **GraphQL / DataLoader / N+1.** No retrievable material. DevRef's nearest is
  Newman, *Building Microservices* p.177 on batch APIs, which is about round-trip
  cost, not about failing on data the resolved branch didn't need. Not cited.
- **Strictness / demand analysis (compilers).** Searched PapersFast; retrieval
  returned pure noise (Bloom filters, RAPPOR, conjunctive-query parallelism). The
  literature certainly exists — a demand set IS a strictness analysis, and the
  over-approximation direction there is the safe one — but it is not retrievable
  from this corpus. Unexplored, not refuted.
- **helm, gh-cli, kubebuilder, goreleaser.** helm characterized and rejected
  above. The others were not pursued once OpenTofu produced an exact instance;
  additional weak analogues would not move the verdict.

## Queries run

| Corpus / tool | Query | Outcome |
| --- | --- | --- |
| PapersFast (semantic) | `XACML policy decision point attribute retrieval missing attribute Indeterminate` | rejected — noise (XML Schema, XQuery) |
| PapersFast (text) | `PIP attribute finder Indeterminate missing attribute policy evaluation` | rejected — PDP/PIP architecture only, no retrieval semantics |
| PapersFast (text) | `XACML MustBePresent attribute missing Indeterminate rule combining algorithm deny-overrides` | rejected — SAML profile refs; spec absent |
| PapersFast (text) | `Cedar policy language total function no partial evaluation error entity does not exist` | accepted — Cedar §1, §5.3, §6 |
| PapersFast (text) | `slice of policies must contain all policies that could be satisfied scope-based indexing same outcome as if it had considered all available policies` | accepted — decisive, Cedar §3.3 p.8/10 + §1 p.2 |
| PapersFast (text) | `if evaluation of a policy results in an error that policy is ignored diagnostics authorization response` | accepted — Cedar §2.5, §3.2, n.2; also EOPL lazy-eval |
| DevRef (semantic) | `declarative query language specifies what data is required not how to retrieve it optimizer freedom predictable` | rejected for A — DDIA argues optimizer freedom, i.e. pruning |
| DevRef (semantic) | `lazy evaluation short circuit avoids evaluating operands that would raise an error strict versus non-strict` | accepted — Meyer, McConnell, EOPL |
| DevRef (text) | `query optimizer determines which columns and tables are required before execution column pruning projection pushdown` | partial — PostgreSQL internals; supports pruning, not over-approximation |
| DevRef (text) | `N+1 query problem eager loading versus lazy loading fetch all at once round trip` | rejected — Encyclopedia-of-Algorithms noise |
| DevRef (semantic) | `batch fetch all needed data in one round trip instead of many small requests network latency` | rejected — batch APIs, wrong axis |
| DevRef (text) | `fail fast validate all preconditions up front rather than failing partway through predictable errors` | accepted — Wolff, Fail Fast |
| PapersFast (semantic) | `strictness analysis demand analysis safe approximation evaluate eagerly only if definitely needed` | rejected — pure noise |
| semble (opentofu) | `check that all required providers are installed and error if provider not initialized run terraform init` | accepted — decisive, `context.go:327` |
| grep (opentofu) | `checkConfigDependencies` callers; `Targets` in `graph_builder_plan.go` | accepted — establishes ordering :163 before :260 |
