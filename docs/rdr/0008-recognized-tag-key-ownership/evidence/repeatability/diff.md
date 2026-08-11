model: claude-opus-5[1m]

# Repeatability diff — RDR 0008

Diffed `run-1.md` (claude-opus-5[1m]), `run-2.md` (claude-fable-5),
`run-3.md` (glm-5.2:cloud). `variant: full (profile: foundational)` on all
three — coverage matches the profile. This session authored none of the runs.

Ordered by how many runs disagree (3-way first), then by whether the contract
is load-bearing. Each finding names the RDR section it lands in.

---

## 1. Disagreements

### D-1 — The advisory's carrier and payload (3-way divergence; **[boundary split]** on required-ness)

**Contract**: the non-blocking near-miss advisory (Load-Bearing Decisions /
Identity, asserted by scenario 4). What data it carries, what channel it
travels, and whether anything makes it comparable.

| Run | Rendering |
| --- | --- |
| run-1 | "separate advisory/diagnostic channel on the normalizer result, distinct from the failure list"; carries authored + reserved spelling. Explicitly flagged: "**the largest under-specified surface in the RDR**" — no category token, no rule id, no field list, no carrier. |
| run-2 | GUESS: "its own type, since it is not a validation failure and must not alter the verdict"; carries authored + reserved spelling. Also GUESSes the fold: "not pinned to a Go API — guessed `strings.ToLower`-equivalent simple fold". |
| run-3 | `type NearMissAdvisory struct { AuthoredName, ReservedName string }`; "carried as a separate advisory list alongside the failures" — GUESS. Additionally GUESSes the trigger's conjunction semantics: "case-folded-only or trimmed-only is `GUESS` — the RDR says folding *and* trimming, conjoined; I read that as both required." |

**RDR passage that let them diverge**: block 2's failure payload is pinned to
the byte — three named fields, two literal token values, a golden test. The
advisory beside it (Load-Bearing Decisions / Identity, lines 1289–1296) is
normative and tested but names *no* fields, *no* token, and *no* carrier: "MUST
raise a **non-blocking advisory** naming both spellings." Scenario 4 asserts
"the advisory is pinned as targeted rather than blanket" and that it "MUST NOT
change the load/lint verdict" — both properties of the *verdict*, neither a
property of the *payload*. The asymmetry with block 3 is what all three runs
noticed independently.

**Boundary split**: run-3 (alt-model) is the only run that read the fold+trim
conjunction as a live ambiguity ("case-folded-only or trimmed-only is GUESS").
run-1 and run-2 both silently implemented `fold(trim(name)) == "recognized"` —
the conjunction of the *transforms*, applied together. Under run-3's reading,
`Recognized` (case-only) and `" recognized"` (whitespace-only) may not qualify;
scenario 4 requires that all three near-misses raise the advisory, so run-1/2's
reading is the one that satisfies the scenario. The RDR's prose ("under Unicode
-simple case folding *and* trimming") is genuinely readable both ways, and the
alt-model draw is where that surfaced.

**Lands in**: Normative Contracts (a new advisory clause, or block 3 widened) —
the advisory needs at minimum a rule identifier token and a named field pair,
on the same footing as `reserved-tag-key/kernel-owned`, so scenario 4's
assertion has something byte-comparable to assert. Plus a one-line
disambiguation of the fold+trim predicate.

---

### D-2 — The exported predicate's error arity: first breach vs. all breaches (3-way GUESS, same resolution)

**Contract**: what the exported `Input` predicate returns when more than one
reserved-key breach is present.

| Run | Rendering |
| --- | --- |
| run-1 | "GUESS — whether `ValidateInput` reports one breach or all… This run returns the first breach found." Notes block 4's "detected whenever present" constrains *this* predicate vs. 0009's, not multiplicity within. |
| run-2 | "first hit returns the breach error" — stated without a GUESS marker in the helper description, but the GUESS summary covers "the breach error's type and text". |
| run-3 | "`GUESS`: returns a single aggregate error for the first breach; the RDR does not specify whether one error or a list is returned." Separately GUESSes intra-predicate ordering: "reports the `Input`-key breach before the `RequiresOwned` breach, but the RDR does not say." |

All three landed on first-breach, so this is not a *behavioral* three-way
split — it is a three-way agreement that the RDR is silent. run-1 and run-3
both went further and noted that block 4's "detected whenever present" phrasing
reads as if it settles this, and does not.

**RDR passage**: block 4 — "What this RDR requires of the implementation is
therefore only that its own breach be **detected whenever present** — never
skipped because another precondition also fired." The clause is scoped to
cross-predicate interaction (vs. RDR 0009), and a reader can mistake it for a
completeness requirement *within* the predicate. Block 5 adds a second
obligation to the same predicate without saying whether a doubly-breaching
`Input` reports one or both.

**Lands in**: Normative Contracts block 4 (one clause: single error, first
breach, intra-predicate order unspecified — or the opposite, if aggregation is
wanted). Note this is adjacent to A11 but *not* covered by it: A11 is the
0008↔0009 order, this is 0008-internal.

---

### D-3 — `Input`'s table field: named once, never introduced (3-way GUESS)

**Contract**: how the predicate reaches the rows it must check.

| Run | Rendering |
| --- | --- |
| run-1 | "**GUESS — the RDR never states that `Input` carries the table.** Block 5 says the predicate checks 'the rows of the supplied table'… The field name `Table` is inferred from that phrase." Lists it #3 in its ranked determinacy notes. |
| run-2 | `in.Table.Rows` used throughout without a marker; the GUESS summary does not list it. |
| run-3 | `Table Table // carries Rows` in the struct, and "It reads … `in.Table.Rows[*].RequiresOwned` (block 5 widens it onto the rows, so it is not decidable without the table)". |

All three converged on `in.Table.Rows`, so the *inference* is stable — but it
is an inference from a single incidental phrase, not from a stated contract.

**RDR passage**: block 5 — "this widens the `Input` predicate to read
`in.Table.Rows`". That is the only place the RDR asserts `Input` carries the
table, and it appears in a *consequence* clause about A11, not in the contract
that establishes the predicate's domain (block 4). Block 4 says only "a
construction-time predicate over `Input`".

**Lands in**: Normative Contracts block 4 — state the predicate's read-domain
once, where the predicate is defined, rather than leaving it to be recovered
from block 5's A11 aside. Low risk (all three recovered it) but it is the kind
of silence that a fourth reader recovers differently.

---

### D-4 — The `required name` field in the owned/observed direction (1-run finding, load-bearing)

**Contract**: block 3's mandated three-field payload, applied to the
*second* naming rule.

Only **run-1** surfaced this, and it is the sharpest finding in the set:

> "GUESS — the owned/observed-named-`recognized` case reuses the same
> `required` field value (`"recognized"`), which reads oddly since the author
> must RENAME AWAY from it there. The RDR mandates three fields for 'every
> reserved_tag_key failure' without splitting by direction."

run-2 and run-3 both emitted the same payload for both directions without
noticing the semantics inverted — run-2's pseudo-code literally writes
`failures += reserved_tag_key(offending=key, ...)` with the same `required` for
the owned/observed branch; run-3's `ReservedTagKeyFailure.RequiredName //
literal "recognized"` is unconditional.

**RDR passage**: block 3 — "Every `reserved_tag_key` failure … MUST carry … the
offending name as authored, **the required name (the literal `recognized`)**,
and a stable rule identifier." Block 2 defines two rules with opposite
remedies: a recognized-provenance declaration must be *renamed to*
`recognized`; an owned/observed declaration must be *renamed away from* it. The
payload contract is written for the first rule and applied to both. Under the
second, "required name: recognized" is not merely uninformative — it is
actively wrong guidance, and the Trade-offs section leans on exactly this
payload as the mitigation for the Negative consequence ("authors lose naming
freedom … Mitigated by the failure *data* — offending name, required name, rule
identifier — which is normative and tested (scenario 7)").

This is a determinacy silence with a user-outcome consequence, not a style
difference: the RDR's own stated mitigation misfires for one of its two rules.

**Lands in**: Normative Contracts block 3 (split the payload by direction, or
give the two rules distinct rule identifiers, or restate what `required name`
means when the remedy is a rename-away). Scenario 7 and scenario 2's
owned-tag-named-`recognized` half both assert against it.

---

### D-5 — All three runs modeled `Input.Owned`/`Observed` as maps; on `main` they are `[]Tag` slices (3-way shared wrong guess, surfaced at grounding)

**Not a run disagreement** — a three-way *agreement* that grounding refutes,
recorded here because the lens's job is to find where the RDR reads
differently than the code it constrains.

| Run | Rendering |
| --- | --- |
| run-1 | `Observed map[string]string` / `Owned map[string]string`, "GUESS — shape. RDR says only 'tag keys'; map assumed". |
| run-2 | Does not restate the type; pseudo-code iterates `keys(in.Owned) + keys(in.Observed)` — a map idiom. |
| run-3 | `Owned map[string]TaggedValue` / `Observed map[string]TaggedValue`, GUESS. |

On `main`: `internal/resolve/resolve.go::Input` carries `Owned []Tag` and
`Observed []Tag`, where `Tag` is `struct { Key, Value string }`.
`::assemble` iterates them as slices (`for _, t := range in.Observed`), letting
the last tag win within a provenance.

**Why this is a finding and not trivia**: a slice admits **duplicate keys**,
a map does not. The RDR's block 4 says producers "MUST NOT supply an owned or
observed tag keyed `recognized`" and the predicate is described as checking
keys — under the runs' map model that is one lookup; under the real slice model
it is a scan that must decide whether a *second* `recognized`-keyed entry is a
second breach (feeding D-2's first-vs-all question with a concrete case the RDR
never contemplated). Scenario 6 constructs an `Input` "carrying an owned or
observed tag keyed `recognized`" — singular — and does not pin the duplicate
case either.

The RDR does describe `Input.Owned` correctly elsewhere as "an
accessor-produced owned tag snapshot" (A8's scope limit, quoting the package
doc), so this is silence rather than contradiction: nothing in the RDR states
the collection shape, and every run invented the same wrong one.

**Lands in**: Normative Contracts block 4 — one clause stating the predicate
scans the owned/observed tag *sequences* (they are slices, duplicates
admissible), which resolves alongside D-2's arity question rather than as a
separate edit.

---

### D-6 — Whether `recognizedTagKey` is exported (2-run divergence, both dispositions conforming)

| Run | Rendering |
| --- | --- |
| run-1 | Shows the export commented out: "the RDR does NOT require this to be exported… a conforming implementation may omit it." |
| run-2 | "The RDR explicitly does NOT require exporting a spelling constant" — omits it. |
| run-3 | "`GUESS`: left unexported, matching D2's 'no public surface' posture; the normalizer hardcodes the literal under the behavioral test." |

Not a defect — the RDR states this is implementation latitude (Technical
Design). Recorded so resolve does not count it as a silence. The one thing
worth noting: run-3 inferred *how* the normalizer then gets the spelling
("hardcodes the literal"), which the RDR also leaves open, and which the
behavioral-conformance requirement (premortem P-5) makes safe either way.

---

## 2. GUESS clusters (contracts ≥2 runs marked GUESS)

| Cluster | Runs | Load-bearing? |
| --- | --- | --- |
| **Advisory data shape / carrier / trigger conjunction** | 1, 2, 3 | **Yes** — normative + tested (scenario 4), and the only normative artifact in the RDR with no comparable token. → D-1 |
| **Exported predicate's name** | 1, 2, 3 | No — explicit latitude, and the RDR forbids binding it to an 0009 symbol. All three picked `ValidateInput` and all three marked it GUESS. Determinate *as latitude*. |
| **Breach error's type and text** | 1, 2, 3 | No — the RDR fixes the channel (Go error path) and the trigger; the error value is implementation latitude by the same posture as the predicate name. All three assumed `fmt.Errorf`-shaped, none assumed an `errors.Is` sentinel. |
| **`Input.Owned` / `Input.Observed` collection shape** | 1, 2, 3 | **Yes, but not for the reason the runs gave.** All three dismissed it as irrelevant ("only the key matters") and all three guessed *map*; `main` has `[]Tag`. The duplicate-key case a slice admits is what makes it load-bearing. → D-5 |
| **Failure/advisory struct shapes** | 1, 2, 3 | No — the RDR affirmatively assigns these to RDR 0002 ("The Go type, package, and field names … are RDR 0002's to choose"). Determinate by delegation. |
| **First-breach vs. all-breaches, and intra-predicate order** | 1, 3 (2 implicitly) | **Yes** — see D-2. |
| **`in.Table.Rows` as the predicate's read-domain** | 1, 3 | Marginal — see D-3. |
| **0008 vs 0009 entry-precondition order** | 1, 2, 3 | No — the RDR explicitly leaves it open (A11 / Stage 7.1) and *says so*. run-1 recorded it defensively: "Not a defect; recorded so the diff does not count it as one." Correct. |
| **Ordering of the reserved-key check among 0002's existing lint categories** | 3 only | No — 0002's validation order is 0002's, and block 2 already fixes the one ordering that matters (TOML duplicate-key precedes the naming rule). |

---

## 3. Agreement

The runs rendered these identically — the RDR is determinate here, and these
are not findings:

- **The reserved value and its byte-exact identity rule.** All three: `recognized`,
  post-parse key string, case-sensitive, no trim, no fold; `[tags."recognized"]`
  ≡ `[tags.recognized]` (both reserved); `Recognized` / `RECOGNIZED` /
  `" recognized"` ordinary and unreserved.
- **The two declaration rules and their direction** (recognized-provenance MUST
  be named `recognized`; owned/observed MUST NOT be), and that the checked
  positions are `[tags.<tag>]` keys **only**.
- **The fall-through for predicate positions** (`unknown tag`) and for
  `[rule.write]` (`write to non-owned tag`) — all three reproduced the
  distinction block 2 draws, including *why* `unknown tag` stops firing once
  `[tags.recognized]` is legally declared.
- **Both literal tokens and their non-interchangeability**: category
  `reserved_tag_key`, rule id `reserved-tag-key/kernel-owned`, "not
  alternatives, a consumer MUST NOT choose between them", category for RDR
  0005's dispatch, rule id for golden assertions.
- **`Resolve(in Input) (Result, error)` unchanged in signature**, breach →
  non-nil error + zero-valued `Result` (no `Plan`, no `Refusal`), never a new
  `RefusalKind`, nil error preserved for conforming input including the empty
  `Input{}` at `resolve_test.go:748`.
- **"One predicate, two call sites"**, the predicate carrying *both* block 4's
  and block 5's obligations, and unconditional on `in.Recognized`.
- **`assemble` unchanged**, sole view constructor, injects only for non-empty
  `in.Recognized`, D3 precedence (`owned` > `observed` > `recognized`)
  untouched.
- **The bypass residual**: a package-internal caller with
  `RequiresOwned: ["recognized"]` still refuses `owned_state_unavailable` with
  `MissingOwned` containing the reserved key.
- **Cardinality as a consequence, not a check**: ≤1 recognized-provenance
  declaration falls out of name-keying; two spelled `recognized` are
  `malformed TOML`; no lower bound.
- **The two-part Done split** (kernel half at HEAD; normalizer half carried into
  RDR 0002's implementation).

## Health

**Healthy.** The diff localizes to five contracts (D-1 through D-5), GUESS
markers cluster on the same under-specified surfaces across runs, and the
clusters that are *not* findings are ones the RDR affirmatively delegates
(0002's structs) or explicitly leaves open (A11). The three runs disagree where
the RDR is silent and agree where it is pinned, which is the signal this lens
exists to produce. The alt-model draw earned its keep once (D-1's fold+trim
conjunction) and the strongest run-sourced finding (D-4) came from a semantic
read no other run performed.

One qualifier on "healthy": D-5 is a **shared** wrong guess — three runs
independently modeled `Input.Owned`/`Observed` as maps where `main` has slices,
and each filed the question as irrelevant rather than as a GUESS worth ranking.
That is the identical-but-confidently-wrong shape in miniature, caught by
grounding rather than by run disagreement. It does not invalidate the pass (it
is one contract, and the runs *did* mark the shape as invented), but it is why
the diff was grounded against the kernel rather than resolved from the three
runs alone.
