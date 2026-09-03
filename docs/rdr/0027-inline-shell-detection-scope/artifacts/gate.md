# Finalization Gate — cli/0027 (inline-shell detection scope)

RDR: `0027-inline-shell-detection-scope` · Date: 2026-09-03 · **Verdict: READY**

Mechanical pre-sweep: `evidence/tooling-pass/tooling-pass.md` — **PASS**
(`rdr lint --locking 0027` exit 0, `blocking=0`; one C10 dangling peer-element
reference fixed in-pass; joint-decision fence cleared by firing arm 2 against
0022). Item 4 (Cross-Cutting Concerns) is authored in the record at
`cli/0027:G-cross-cutting`, not here.

## 1. Contradiction Check

No contradictions found between research findings, design principles, and
proposed solution. Three places were checked directly because they are where
this record's rounds moved:

- **Research Findings vs C1's predicate.** Key Discoveries states that freeing
  the interpreter's position "DOES add a new false-refusal class"; C1 ships that
  cost rather than denying it, §consequences carries the negative bullet by
  name, and the `D-selection-predicate` LBD records the tie-break judged against
  it. The finding and the solution agree that the class exists and is accepted.
- **"Closes the wrapper class" vs the two admitted forms.** §approach's summary
  sentence and C1's `out of scope, BY NAME` line are the same claim at two
  lengths: the promise is stated over argv WORDS, so a one-word shell string and
  a stdin-fed interpreter are outside the surface by construction, not by
  omission. §decision-rationale's premortem names precisely the reader who stops
  at the short form, and C1's `promise:` clause is the mitigation. No conflict.
- **0025's stated principle vs re-promoting the deny-list.** 0025 scopes the
  deny-list as "defense in depth … not the barrier itself." C1 keeps that frame:
  it widens the predicate's *position* while explicitly prohibiting a wrapper
  table ("no wrapper table exists and none may be added"), which is what
  re-promotion would require. ALT2 is rejected on exactly that ground. The
  planned feature does not contradict the inherited principle.
- **A2's scope note vs A7's finding.** These are the one pair that reads like a
  contradiction and is not: A2 says the widen refuses nothing *currently green
  in this repo's fixtures*; A7 says the widen *does* open a new refusal class in
  general. A2's own Status text states the limit ("it does NOT license the
  general claim … which A7 computes false") and §mini-check-trace step 1 repeats
  the boundary. Consistent by construction.

## 2. Assumption Verification

Seven Critical Assumptions, `ca_total=7 / ca_verified=6 / ca_pending=1 /
ca_unverified=0`. Every record is internally consistent — Status, Method and
Evidence agree, and every "If wrong" is non-empty and states an observable
consequence rather than a restatement.

- **Method vocabulary**: `ca_off_vocabulary=0`. A1 `Source Search`; A2–A7
  `Spike`. **No `Docs Only` record exists**, so nothing is blocked on that head.
- **Self-reference**: the single `Source Search` record (A1) anchors
  `internal/table/load.go::carrierDefect` with line cites into the source tree —
  no path resolves to this RDR or its artifact dir. Tooling CHECK 3 clean.
- **Symbol resolution**: `anchors_total=35`, `anchors_unresolved=0`,
  `anchors_unlooked=0` — every cited `path::Symbol` resolves on the repo,
  including the cross-repo prior-art anchors. `peer_evidence_unresolved=0`.
- **No `Verified` stamp proves only an adjacent claim.** The two scope-limited
  stamps carry their limit in the Status field itself rather than in prose
  elsewhere: A2 (scoped to this repo's fixtures and 0025's models) and A7
  (measured rate zero over a 264-vector sample, with the genre caveat stated).
  Both are honest narrowings, and both are the narrowing the round asked for.

**The one `Pending`: A5.** It stays Pending deliberately, downgraded at Stage 6
and accepted here rather than carried as a blocker, on three grounds:

1. *C1 does not rest on it.* A5's withholding premise was falsified as a
   present-tense fact — `cmd.Stdin` is assigned unconditionally at
   `internal/cli/cmdbind/cmdbind.go:207`, the sole assignment, with no branch or
   field read before it. What died is a durability claim about a *future*
   record. C1 names the stdin forms as admitted **by the predicate's own
   definition** (`-` and `-s` are not listed inline-code flags), and §approach
   already called the successor routing "an ownership claim, not a closure."
2. *No settled-fact prose depends on it.* Checked every dependent site.
   §mini-check-trace row 3c states outright that "A5 is no longer a witness for
   it … so this row says the form is ADMITTED, not that it is owned," and rests
   the green verdict on the predicate alone. §capability-dependencies marks the
   withholding capability **Deferred** and states "this record depends on none
   of it." §consequences names the stdin forms as admitted-and-documented, and
   flags that no kata yet tracks the successor. The Status-consistency rule is
   satisfied: nothing treats the withholding mechanism as settled.
3. *There is a named plan, and it is "no code."* This record ships nothing for
   the stdin axis. The obligation it does carry is MVV step 3 / scenario S4 —
   `["sh","-s"]` and `["python","-"]` must lint GREEN and load with argv
   unchanged — which is a predicate assertion, verified by the predicate. The
   withholding point is the successor's to add, at the single call site named
   above, characterized by the Stage 6 spike as a one-branch, one-field local
   change.

The residual cost, stated plainly: until that successor ships, a stdin-fed
interpreter is admitted and a reviewer must read argv for it. That is exactly
what C1's `promise:` clause tells a reviewer to expect, so the gap is documented
rather than hidden — which is this RDR's whole thesis about what a check may
promise. Accepted, not deferred silently.

## 3. Scope Verification

The Minimum Viable Validation is **in scope and executed during
implementation**, not deferred. Its five steps map onto the three phases with
nothing left over, and Testing Strategy scenarios S1–S7 are the executable form.

The specific tests:

- **The predicate half** — S1 (one mutant per named wrapper refuses with
  `command_shell_interpreter`), S2 (`["nice","sh","-c","cat {artifact}"]` reports
  the interpreter defect and **not** `command_unknown_placeholder` — asserted on
  the category *string*), S3 (the nine `env`-option forms q2q1 enumerates
  refuse, proving the deleted walk is subsumed), S4 (admitted forms lint green
  **and** load with argv unchanged), S5 (`["ruby","tool.rb","-e","prod"]` refuses
  naming `ruby -e`, asserted on the detail *text*), S6 (`table.Categories()`
  membership and order unchanged, as a golden assertion over the full ordered
  slice).
- **The promise half** — **S7 is the named proof and the one that matters
  here**: the description Phase 3 ships is read from the `--help-all` extended-help
  body with no model loaded and no defect present, and must contain **both**
  out-of-scope forms in C1's words. This is the MVV's step 5.

S7 is the scope answer rather than a formality. The `oracle` mini-check states
in writing that steps 1–4 "all pass with no description at all" — so without
S7 the `promise:` clause would ship untested, and the record's central claim
(what a reviewer may rely on) would be unverified. Its control starts red by
construction: no description surface exists today
(`internal/table/category.go::Categories` returns identifiers only). A6 is
Verified, so the landing site is settled and S7 is writable.

The `oracle` mini-check also disarms the weak row honestly: S4's "lint passes"
is an absence-of-error oracle a no-op satisfies, and it is retained only as a
*pair* with S1/S3, with `["python","-"]` and `["sh","-es"]` named as the
discriminating members. That is the correct treatment, not a deferral.

## 5. Proportionality

Right-sized. Nothing to trim before locking.

- **Contract count — the split test.** This RDR is the sole author of **exactly
  one** independent load-bearing contract: `C1`, the successor text for 0025:C5's
  `command_shell_interpreter` line and its `interpreter set` line
  (`contracts=1`, `contracts_durable=1`, `contracts_transient=0`). It owns one
  seam — clause 4's predicate and the claim wording over it. Everything adjacent
  is cited, not restated: 0025:C5 keeps the clause map, within-entry precedence,
  registration and the other five categories; cli/0028:C1 owns the `edit` arm
  and the tail categories. No second seam to split off.
- **Profile re-validation.** The field reads `large — C1, the successor
  predicate and claim wording for 0025:C5's `command_shell_interpreter`
  deny-list; user-facing yes; locks format`. Re-derived against the contracts
  just counted, it holds: one contract, but a **user-facing surface** (the
  refusal a model author sees, plus the description Phase 3 ships) and it
  **locks a format** (the reported form `argv[i] + " " + argv[j]`, the detail
  text, `Categories()` order). Those two dispositions are what carry it above
  `small`/`mid`. The form is correct — value plus one clause naming the
  contract, no matrix or provenance prose left from the template. The accretion
  floor does not apply (`seam_lineage_count=0`, rule `floor-below-two`), so the
  recorded value is at or above it.
- **The lenses that ran agree with it.** `--outcome lens` returns
  `emit.row: none` (rule `lens-large-row-complete`): grounding → 3amigo →
  critique all present, with `critique_models=differ`, so the cross-model
  critique is a real second pass and not a single-model fallback. Repeatability
  is an add-on gated on the written `Determinacy:` line, which reads `n/a` with
  its reason (no new signature, type or API surface — C1 replaces
  `interpreterForm`'s body at its existing name and signature, with the
  quantifier, basename rule, tie-break, reported form and read-scope all pinned
  in the clause). `--outcome repeatability` → `resolve:determinacy` → `none`.
  Nothing on the row is owed. No under-sized Profile routed past a lens.
- **Document size.** 601 lines for a `large` record that carries seven Critical
  Assumptions with spike evidence, three scored alternatives, four mini-checks
  and a seven-scenario testing strategy. The two long sections — the CA
  evidence fields and `§mini-check-trace` — are load-bearing verification
  content the grounding sweep reads, not narration; lint raises no
  `evidence:over-budget`. §research-findings correctly records the negative
  prior-art result (no peer CLI detects shell-ness from argv) rather than
  padding around it. No section flagged for trimming.
