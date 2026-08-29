# Resolve — literature and prior art for declared emit vocabulary

Model: claude-opus-5[1m]

Scope: literature and prior-art bearing on an OPTIONAL per-key `kind` + allowed-value `domain`
declaration on the `emit` block that, when present, makes the loader refuse undeclared keys and
out-of-domain values at authoring time, with evaluation semantics unchanged (values stay
uninterpreted, byte-compared).

Corpora searched: `DevRef` (148 SE reference books), `PapersFast` (~3623 CS papers),
plus the sibling analysis repo (`state-machines/`) and the peer-tool checkouts (`langref/`).
`StateMachineLit` and `StateMachineRes` were not queried — budget was spent where the first
two corpora were producing hits, and the code checkouts answered Q5 directly.

Path convention: paths below are repo-relative or checkout-relative. `DevRef` corpus items are
named by book title plus the corpus-internal filename, since the corpus is a PDF library.

## Query log

| # | Corpus / target | Mode | Query string | Outcome |
|---|---|---|---|---|
| 1 | PapersFast | semantic | `gradual typing soundness trade-offs optional type annotations partial adoption` | REJECTED — returned bibliography and table-of-contents fragments only |
| 2 | DevRef | semantic | `optional static type checking pay as you go opt-in strictness dynamic language` | HIT — Meyer OOSC §17.2 |
| 3 | DevRef | text | `a little bit typed` | HIT — confirms and expands #2 |
| 4 | DevRef | semantic | `schema validation unknown fields open closed content model additionalProperties tolerant reader` | HIT — Daigneau, Tolerant Reader |
| 5 | DevRef | semantic | `compiler error recovery report multiple errors one pass diagnostics fail on first error` | WEAK — returned debugging chapters, not diagnostics design |
| 6 | DevRef | semantic | `validating outputs against enumerated set of allowed values closed vocabulary controlled values` | PARTIAL — surfaced SQL Antipatterns Ch. 11 |
| 7 | DevRef | text | `error recovery parser continue after syntax error report as many errors` | REJECTED — matched Go `panic`/`recover` parser examples, not error-recovery design |
| 8 | DevRef | semantic | `restrict column to fixed set of values ENUM check constraint lookup table adding new value requires schema change` | HIT — SQL Antipatterns Ch. 11 "31 Flavors", full argument |

Two further DevRef queries (`fail fast crash early validate input at boundary...`, and an
identity-confirmation full-text query for the SQL Antipatterns chapter) were run to disambiguate
hits already in hand.

## Q1 — Gradual / optional typing and opt-in strictness — VERIFIED

**Bertrand Meyer, _Object-Oriented Software Construction_ (2nd ed.), §17.2 "Static typing: why and
how", subsection titled "A little bit typed"?, p. 646** (`DevRef`, `Books/Object Oriented Software
Construction-Meyer.pdf`, chunk 801, page 646).

Meyer states the partial-adoption problem directly, and it is the closest thing in the corpus to a
normative rule for this decision:

> "It was noted above that we should aim for a _strong_ form of static typing. This means that we
> should avoid any loopholes in the static requirements — or, if any such loopholes remain,
> **identify them clearly, if possible providing tools to flag any software using them.**"

He then argues that a checking regime with an unrestricted escape hatch does not deliver a
guarantee: "It seems difficult to accept claims of static typing if at any stage the developer can
eschew the type rules through casts. Accordingly, the rest of this chapter will assume that the
type system is strict and allows no casts."

Crucially, Meyer draws the distinction that licenses the proposal rather than condemning it. He
separates two things that look alike:

> "You may have noted that assignment attempts ... superficially resemble casts. But there is a
> fundamental difference: an assignment attempt does not blindly force a different type; it
> _tries_ a candidate type, and **enables the software to check whether the object actually matches
> that type.** This is safe, and indispensable in some circumstances."

Bearing on the decision. Meyer supports opt-in declaration on two conditions: (a) the unchecked
region must be *identifiable*, not merely tolerated, and tooling should be able to flag it; (b) an
escape hatch that *checks* is legitimate, an escape hatch that *asserts without checking* is not.
An emit key with a declared `domain` is Meyer's "assignment attempt" — it proposes a value set and
then verifies membership. An undeclared key, or a `kind` with no enumerable domain (the
`scalar`/`any` shape named in the task), is Meyer's loophole: acceptable only if it is nameable and
countable, so that authors can be told how much of their emit surface is unchecked. The design
consequence is that the loader should be able to report the declared-vs-undeclared ratio, not just
silently permit the undeclared remainder.

Rejected branch: `PapersFast` semantic search for gradual-typing theory (Siek/Taha soundness,
"pay-as-you-go") returned only bibliography lines and TOC tables across two attempts (queries 1 and
the diagnostics retry). The corpus contains the PDFs but its embeddings for these documents are
dominated by front matter; no primary gradual-typing paper was retrievable. `negative:` Q1 primary
gradual-typing literature (Siek & Taha, mypy/TypeScript/Dialyzer design rationale) — searched
`PapersFast`, not found in retrievable form; the Meyer citation carries this question instead.

## Q2 — Schema validation of outputs vs inputs — VERIFIED

**Robert Daigneau, _Service Design Patterns_, Tolerant Reader, pp. 272–274** (`DevRef`,
`Books2/rest/service_design_patterns.pdf`, chunks 257 and 259).

The pattern is normally cited for the *reader* side — "Design the client or service to extract only
what is needed, ignore unknown content, and expect variant data structures" (p. 272). The relevant
finding is the *sender*-side counterpart on p. 274, which is the direct answer to the
inputs-vs-outputs question:

> "all message senders should make every effort to conform to the agreed-upon protocols for message
> composition because a message sender that commits a gross violation of the 'message contract' can
> cause significant problems, even for a _Tolerant Reader._ One such example occurs when the
> message sender fails to submit a required element or **uses the wrong data type for some item.**
> Message senders can therefore facilitate effective communications by **using schema validation
> before sending a message.**"

Bearing on the decision. This is the established practice the proposal is asking for, stated as an
asymmetry rather than a symmetry: tolerance is a *reader* virtue and strictness is a *sender*
virtue. The rule engine is the sender of the emit payload; the downstream consumer that executes
the answer verbatim is the reader. Daigneau says the sender validates against the schema before
emitting, precisely because the reader's tolerance cannot rescue a misspelled value. That the
project already validates inputs but not outputs is, in this framing, the wrong half: input
tolerance and output strictness are the pattern, and validating only inputs inverts it.

**Robert Karwin, _SQL Antipatterns_ (Pragmatic Bookshelf, P1.0 May 2010), Chapter 11 "31 Flavors",
pp. 131–138** (`DevRef`, `Database-Books/db(33).pdf`). The chapter opens with exactly the emit
problem, in the objective statement on p. 132:

> "Ideally, we need the database to reject invalid data: `INSERT INTO Bugs (status) VALUES ('NEW');
> -- OK` / `INSERT INTO Bugs (status) VALUES ('BANANA'); -- Error!`"

This is prior art for the goal (a closed, enumerated value domain that rejects the typo) and,
below, the sharpest available evidence on its costs.

## Q3 — Whole-model vs per-key strictness, and the ergonomics of the all-or-nothing switch — VERIFIED

Three independent sources, one supporting the switch and two costing it.

**(a) The switch exists and is a single global boolean in a peer Go CLI.**
`goreleaser/internal/yaml/yaml.go:12-24` (checkout `langref/goreleaser`) wraps the YAML decoder in
exactly two modes and nothing between them:

```go
// UnmarshalStrict unmarshals a YAML document with strict behavior (only declared fields are tolerated).
func UnmarshalStrict(in []byte, out any) error {
	decoder := yaml.NewDecoder(bytes.NewReader(in))
	decoder.KnownFields(true)
	return handleErr(decoder.Decode(out))
}
```

`goreleaser/pkg/config/load.go:59` calls `yaml.UnmarshalStrict` for the user-authored config, and
`goreleaser/internal/artifact/artifact.go:317` calls `decoder.DisallowUnknownFields()` on the JSON
side. The comment "only declared fields are tolerated" is the whole-model rule stated verbatim:
strictness is a property of the decode, not of individual fields, and the declared struct *is* the
vocabulary. Helm does the same at `helm/pkg/repo/v1/index.go` and `helm/pkg/chart/v2/util/chartfile.go`.
Restate's protocol schema (`state-machines/repos/restate/service-protocol/endpoint_manifest_schema.json`)
sets `"additionalProperties": false` at the document root (line 301) and at every nested object
(lines 92, 130, 221, 292) while using `"enum"` for the value domains (lines 10, 29, 48, 70, 204,
273) — a fully closed content model with enumerated value sets, i.e. both halves of the proposal.

Note the shape of the analogy and its limit: in all of these the vocabulary is declared *in the
host language's type* (a Go struct, a JSON Schema), so "declared" and "undeclared" are never
ambiguous and the all-or-nothing switch costs nothing. The proposal's situation is harder because
the declaration is optional and partial, which is what Q1's Meyer citation and (b) below address.

**(b) The dominant cost of a closed enumerated domain is widening it.** Karwin, Ch. 11, pp. 134–135:

> "The most common alterations are to add or remove one of the permitted values. There's no syntax
> to add or remove a value from an ENUM or check constraint; **you can only redefine the column
> with a new set of values.**"

> "As a matter of policy, changing metadata ... should be infrequent and with attention to testing
> and quality assurance. If you need to change metadata to add or remove a value from an ENUM, then
> you either have to skip the appropriate testing or spend a lot of software engineering effort on
> short notice to make the change. Either way, these changes introduce risk and destabilize your
> project."

And the decision rule, p. 135: "before considering using ENUM, first ask yourself whether the set
of values are expected to change or even whether they might change. If so, it's probably not a good
time to employ an ENUM ... ENUM is most likely to succeed when it would make no sense to alter the
set of permitted values."

Karwin also names the failure mode of *not* declaring: "The list of values in the application code
got out of sync with the business rules in the database — again. This is a risk of maintaining
information in two different places." (p. 135). Both horns are real; the proposal is choosing the
declared horn.

His recommended mitigation (pp. 137–138) is a lookup table with an `active` flag, so a value can be
retired without being deleted:

> "you can add another attribute column to the lookup table to designate some values as obsolete.
> This allows you to maintain historical data ... while distinguishing between the obsolete values
> and values that are eligible to appear in your user interface."

Bearing on the decision. The proposal is materially cheaper than Karwin's antipattern on the axis
he cares about, and the RDR should say so rather than inherit his verdict: the emit domain lives in
the same authored, version-controlled, text file as the rules, so widening it is a one-line edit to
the same artifact the author is already editing — not an `ALTER TABLE` against populated data, not
an offline migration, and not a second location that can drift. Karwin's objection is aimed at
declarations that are *remote from the values they constrain*; this one is adjacent to them. What
does transfer is his "old flavors never die" point (p. 135): if emit values are ever persisted in
prior decision records, narrowing a domain later can orphan historical answers, so the RDR should
state whether domain *narrowing* is permitted or whether domains are append-mostly.

**(c) The all-or-nothing switch is ergonomically survivable when the refusal is a good message.**
See Q4 — golangci-lint refuses unknown names hard but accumulates them and names the discovery
command. The ergonomic cost of the switch is paid almost entirely at the moment of refusal.

`negative:` Q3 protobuf unknown-field preservation semantics and TypeScript excess-property checking
— searched `DevRef` (queries 4, 6, 8) and the `langref` checkouts; no citable primary discussion of
either was found. Rust `#[serde(deny_unknown_fields)]` was searched for directly in
`state-machines/repos/statewright/crates/` and is **not** used there, which is itself a finding
(see Q5).

## Q4 — Fail-fast vs accumulate — VERIFIED (from peer source, not from the book corpora)

The book corpora did not yield a principle. `negative:` Q4 compiler error-recovery and
error-accumulation literature — searched `DevRef` (semantic query 5, full-text query 7) and
`PapersFast` (semantic); DevRef returned debugging-technique chapters (Code Complete Ch. 23, The
Pragmatic Programmer Tip 25, Effective Debugging) about diagnosing *your own* bugs, and the
full-text query matched Go `panic`/`recover` parser examples. No treatment of diagnostic
accumulation as a design choice was retrievable. The one adjacent remark is Code Complete p. 587,
"Don't trust the compiler's second message" — which is an argument about cascading *spurious*
errors after a recovery, i.e. a caution about accumulation quality, not a rule.

The peer checkouts answer it cleanly, and the answer is a hybrid rather than either pole.

**golangci-lint: accumulate within a check, fail fast across checks.**
`golangci-lint/pkg/lint/lintersdb/validator.go` (checkout `langref/golangci-lint`).

Within `validateLintersNames` (lines 37–67) every unknown name is collected before any error is
returned, and they are reported together with a discovery hint:

```go
	if len(unknownNames) > 0 {
		return fmt.Errorf("unknown linters: '%v', run 'golangci-lint help linters' to see the list of supported linters",
			strings.Join(unknownNames, ","))
	}
```

But `Validate` (lines 22–35) runs its validators in sequence and returns on the first one that
fails:

```go
	for _, v := range validators {
		if err := v(&cfg.Linters); err != nil {
			return err
		}
	}
```

Two further details are worth copying. First, the message names the remedy (`run 'golangci-lint
help linters'`) — with a closed vocabulary the tool knows the full legal set, so the refusal can
always tell the author what the legal values are. Second, golangci-lint grades severity by
authority of the source: an unknown linter in the *config* is a hard error, while an unknown linter
in an inline `//nolint` directive is only a warning — `pkg/result/processors/nolint_filter.go:104`,
`p.log.Warnf("Found unknown linters in //nolint directives: %s", ...)`.

**kubebuilder / Kubernetes API machinery: accumulate into a typed list.**
`kubebuilder/docs/book/src/cronjob-tutorial/testdata/project/internal/webhook/v1/cronjob_webhook.go:189-205`
collects `field.ErrorList` across independent checks and returns them as one aggregate
`apierrors.NewInvalid(...)`, with each error carrying a `field.NewPath("spec").Child("schedule")`
path to the offending location. The comment at line 214 states the intent: "The field helpers from
the kubernetes API machinery help us return nicely structured validation errors."

Bearing on the decision. The established shape is: accumulate all violations that are *independent*
and *locatable* (each undeclared key, each out-of-domain value, each with a path to its rule), and
fail fast only across *phases* where a later phase's results would be meaningless or spurious given
an earlier failure — parse before validate, resolve the declaration block before checking emit keys
against it. For this RDR that means one pass over all rules reporting every bad key and value with
its rule location and the legal set, not the first one.

## Q5 — Peer engines declaring output vocabularies — VERIFIED

**statewright (`state-machines/repos/statewright`) is the closest analogue found and is a partial
counter-example worth recording.** It is a state machine that drives an LLM, so like the rule
engine it consumes an emitted symbol from a component that can misspell. Its machine definition
declares the per-state event vocabulary as the keys of `on`
(`crates/engine/src/types.rs:34-36`, `pub on: BTreeMap<String, TransitionDef>`), and the engine
refuses a symbol outside it — `crates/engine/src/transition.rs:439-444`,
`fn no_safe_next_still_errors_on_unrecognized()` asserts that `"GIBBERISH"` from state `draft` is an
error.

The counter-example is *how* it decides membership. Before refusing, `resolve_transition` attempts
recovery against the declared set (`transition.rs:35-76`): a state-name-vs-event disambiguation that
lists the events reaching that target, then a case-insensitive **substring match** ("PASS" →
"TESTS_PASS"), and only then an author-declared `safe_next` fallback state, documented at
`types.rs:51-55`:

```rust
    /// Fallback target when the model emits an unrecognized transition event.
    /// Only fires on unknown events — valid transitions and FAIL are unaffected.
    /// Declared by the state machine author, not inferred.
```

Bearing on the decision. Three things transfer and one must be resisted. Transfers: (1) the
declared vocabulary lives on the *state*, i.e. per-context rather than globally, so "legal" is
scoped; (2) the fallback is opt-in and author-declared, and its doc comment insists it is "declared
by the state machine author, not inferred" — the same instinct as making the domain declaration
optional but explicit; (3) even with a fallback, an unrecognized symbol in a state with no fallback
is still a hard error, so the escape hatch does not dissolve the check. Must be resisted: the
substring coercion. Statewright silently rewrites the emitted symbol to a declared one, which is
exactly what the proposal forbids by keeping values uninterpreted and byte-compared. A near-miss
"PASS"/"TESTS_PASS" coercion is how a typo becomes a *silently different answer* rather than a
refusal. If the RDR wants near-miss help it should be a **suggestion in the error message** ("did
you mean ...?", as golangci-lint points to `help linters`), never a substitution.

**The sibling analysis repo already states the requirement.**
`state-machines/attic/VALIDATOR-REQUIREMENTS.md` line 39, requirement **R6**: "**Queryable current
state** + **enumerate legal next states / legal events from a state**", justified as "The LLM must
be *told* the legal outcome alphabet". Line 38, **R5**, cites "kata verdict reason enum" as an
existing enumerated output domain in this project's own design.

`state-machines/ANALYSIS-kernel-vocabulary.md`, finding **F1**, gives the refusal taxonomy from a
production FSM (`looplab-fsm/errors.go`) and draws the distinction this RDR needs between two
different failures:

- `InvalidEventError{Event, State}` — "event cannot be called in the current state" (`errors.go:21-30`).
  The symbol is in the alphabet but illegal from here.
- `UnknownEventError{Event}` — "event is not defined" → "event `<e>` does not exist" (`errors.go:32-39`).
  The symbol is **not in the alphabet at all**.

Bearing on the decision. This is the vocabulary for the two refusals the loader must emit and they
should not share a message: an undeclared emit *key* is `UnknownEventError`-shaped ("no such key"),
while an out-of-domain *value* on a declared key is `InvalidEventError`-shaped ("not legal here",
and the legal set is enumerable and should be printed).

`negative:` Q5 in other engines — searched `state-machines/repos/` (ms-conductor, inngest, awf-cli,
AgentSpec, StateSmith, scxmlcc, fsm, qmuntal-stateless) for declared output/action vocabularies and
unknown-symbol handling. ms-conductor returned no task-type enum or allowed-value registry. inngest's
matches were OpenAPI/protobuf artifacts, not an authored event-vocabulary declaration. No SCXML
tooling in the checkout was found to constrain an output/action vocabulary. Only statewright and
restate carry the pattern. DMN was not re-searched, per the standing finding that it is absent from
all corpora.

## Summary of what the citations settle

1. Opt-in checking is defensible if the unchecked remainder is *identifiable and flaggable* rather
   than merely permitted (Meyer, OOSC §17.2 p. 646), and a declaration that *verifies* membership is
   categorically different from one that *asserts* a type — the design should count and report the
   undeclared surface.
2. Validating outputs is the established asymmetry, not an innovation: tolerant reader, strict
   sender, "using schema validation before sending a message" (Daigneau, pp. 272–274).
3. Whole-model strictness is what real loaders do (goreleaser `KnownFields(true)`, helm,
   restate's `additionalProperties: false` throughout), and its documented cost is the expense of
   widening the domain (Karwin Ch. 11) — a cost this design largely escapes because the declaration
   is co-located with the rules it constrains. The open question it raises is domain *narrowing*
   against historical answers.
4. Accumulate independent, locatable violations into one report; fail fast only across phases
   (golangci-lint `validator.go`, kubebuilder `field.ErrorList`). Always print the legal set or the
   command that reveals it.
5. The nearest peer (statewright) confirms the declared-vocabulary-plus-optional-fallback shape and
   supplies the anti-pattern to avoid: never coerce a near-miss symbol to a declared one. Suggest,
   do not substitute. Name the two refusals separately (`UnknownEventError` vs `InvalidEventError`,
   per `ANALYSIS-kernel-vocabulary.md` F1).
