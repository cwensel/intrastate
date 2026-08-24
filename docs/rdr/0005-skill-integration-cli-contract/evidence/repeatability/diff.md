model: claude-opus-5
lens: repeatability-lite (diff pass)
variant: lite (profile: mid) — diffs `run-1.md` (claude-sonnet-5) against the RDR

# RDR 0005 — Repeatability-lite diff

Neutral pass: this context authored no run. Lite admits only **concrete
contract silences** — a step whose order the RDR doesn't fix, a field whose
owner it doesn't name, a transform whose tie-break it doesn't state. Naming,
wording, and helper-decomposition differences are excluded by the lens and are
not reported below.

`run-1.md` reproduced the RDR's normative surface with unusual fidelity: all
four verbs, their flag sets, the full 25-row code table with correct exit
groups, gate ordering (post-selection, deny-overrides, exit-3 for gate
timeout), the flat `Finding` record with the exact field set, the
`SetEscapeHTML(false)` canonical-set rule, and the three round-trip
invariants. Those are agreement, not findings (§Agreement below). The
divergences that remain are the RDR's silences.

## Findings

### R-1 — `revision` has no defined source or format [field owner unnamed]

The RDR names `revision` as a required member of all four verbs' success
`data` payloads (§Technical Design *Envelope*, four bullets) and makes "the
same model **revision**" load-bearing in §Load-Bearing Decisions *Identity*
("The same request over the same model revision and the same artifact contents
must produce the same success or refusal") and again in the Finalization Gate
(line 1128). It never says **who produces it, what it is, or how it is
derived** — RDR 0002's loader, a content hash of the model bytes, a
`[model.metadata]` field, or a file mtime are all consistent with the text.

`run-1.md` carried `revision` through every payload without ever describing
it, and did **not** mark it GUESS — it reproduced the token while inheriting
the silence. That is the leak: a determinism claim ("same model revision ⇒
same result") whose subject is undefined cannot be implemented or tested, and
two implementers will pick different derivations.

Lands in: §Technical Design *Envelope* (definition) and/or §Normative
Contracts. Load-bearing: the Identity clause and the MVV's determinism
assertion both rest on it.

**Disposition: pin.**

### R-2 — `--as` default and scope unstated at this contract [transform tie-break unstated]

The RDR calls `--as text|json` "persistent" (line 48, quoting
`docs/cli-output-contract.md`) and its Normative Contracts fix behavior "under
--as=json" and "under --as=text", but never states the **default** when the
flag is omitted. `run-1.md` marked exactly this as **GUESS**: "default
presumed `text` — **GUESS**, RDR only says the flag 'selects the wire
format'."

This is a cited-not-restated boundary: the default plausibly belongs to
`docs/cli-output-contract.md` (existing surface), not to this RDR. The
admissible finding is not "pick a default" but that the RDR leaves a reader
unable to tell **whether it is inheriting one or leaving it open**. A
one-clause citation ("the default is the root contract's; this RDR does not
change it") discharges it without over-specifying.

Lands in: §Technical Design *Envelope* or §Capability Dependencies (the
output-gateway row).

**Disposition: cut/cite** — state the inheritance, do not re-specify the
default here.

### R-3 — reader-selection scope for `next` / `resolve` is unfixed [step order / set membership unstated]

§Normative Contracts says `flow next` and `flow resolve` "MAY run declared
read accessors over explicit `--artifact role=path` bindings and MUST assemble
owned state only from them", and `flow-artifact-missing` fires for "a role **an
invoked accessor needs** that no binding supplies" (line 474, 910). Which
accessors are *invoked* is never fixed: **all** declared readers, or only
those whose keys the candidate rows actually require.

`run-1.md` resolved this silently and differently in two places — its
pseudo-code loops "for each declared read accessor **the model needs**" while
its `assembleOwnedState` helper says it "runs **each declared** read accessor".
The run contradicted itself inside one document, which is the signature of a
contract the source left open rather than a transcription slip.

This is observable, not cosmetic: it decides whether an unbound artifact role
irrelevant to the requested outcome refuses with `flow-artifact-missing` or
succeeds. It also decides how many accessor invocations (each with a timeout
and exit-3 surface) a single `next` costs.

Lands in: §Normative Contracts (the `flow next` / `flow resolve` paragraph).

**Disposition: pin.**

### R-4 — no stated relationship between `flow next` candidates and `--evaluate-gates` gate invocation scope [step order unstated]

§Normative Contracts requires `flow next` to "run gate accessors only when
--evaluate-gates is given; otherwise it MUST list gate ids as unresolved
facts." It does not say **which** gates run under the flag: every gate on
every candidate row, or only gates on rows still viable after guard
evaluation. Unlike `resolve` — where §D9's "after exact-one selection" fixes
the site exactly — `next` has no selection event to anchor to, so the
ordering clause that makes `resolve` determinate has no analogue here.

`run-1.md` did not surface this (it wrote "gates run only under
`--evaluate-gates`" and stopped), so it is a silence the run inherited rather
than diverged on. It is admissible because it is concrete and consequential:
under a model with many candidate rows, the two readings differ in effect
count and in whether a deny on a non-viable row appears in `candidates[]`.

Lands in: §Normative Contracts (the `flow next` paragraph).

**Disposition: pin.**

## Non-findings (excluded by the lens)

- **Go handler signatures** (`runFlowNext(cmd, req) (Result, *clierr.CLIError)`)
  — `run-1.md` marked these GUESS. Correctly so: this RDR fixes a CLI verb
  contract, not a Go package API, and §Approach explicitly leaves internal
  shape to Phases 1–3. Underspecifying this is **by design** (the over-spec
  trap); leave non-normative.
- **Internal helper names** (`assembleOwnedState`, `runSelectedRowGates`,
  `mapRefusal`) — helper decomposition, excluded by the lens.
- **Model file schema** — `run-1.md` marked GUESS and scoped it out itself;
  RDR 0002 owns it and this RDR cites-not-restates by design. Leave
  non-normative.

## Agreement (determinate — no action)

One line, per the lens: `run-1.md` and the RDR agree byte-for-contract on the
four-verb surface and flag grammar; the 25-row failure-code table including
exit groups; exit 3 = environment-unavailable / exit 2 = everything else with
no new group; the flat `Finding` record and its full field list with
`omitempty` and no nested `atom`; `findings` as the single new `CLIError`
field; gate placement after exact-one selection with deny-overrides and
gate-timeout-as-exit-3; escaped plans as successes carrying `escaped` +
`escape_class`; the `<clear>` sentinel being unauthorable and `[]`-vs-absent
staying distinct; and all three round-trip invariants including
`SetEscapeHTML(false)` byte equality. No `GUESS` marker landed on any
normative clause — the GUESS markers cluster exactly on the three
by-design-open areas listed above, which is the healthy signal for this lens.

## Escalation assessment

Escalation to the full ×3 lens is **not** triggered. The run was not
"suspiciously clean with no GUESS markers" — it produced four explicit GUESS
markers, and they localize onto genuinely out-of-scope areas rather than onto
normative clauses. R-1 is load-bearing but is a *definition gap on a token the
RDR already owns*, not a cross-RDR contract a one-run sample under-powers;
pinning it needs a decision, not more draws. Variant stays `lite`.
