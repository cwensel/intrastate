Model: claude-sonnet-5

# Hostile Premortem Critique — RDR 0025 (Command-Invoking Accessor Bindings)

## 6. Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0025:C6` | The "one new parameter on the registry constructor" framing is false: `flowbind.Registry` is a free function with exactly ONE call site (`flow_exec.go:830`) producing a single shared `accessor.Registry` consumed by both `flow_state.go`'s and `flow_exec.go`'s three `NewExecutor` sites, and today's `Reader{Path: acc.Path}`/`Writer{...}`/`Gate{...}` construction has zero notion of a flag or a carrier discriminator. Threading `--allow-commands` in means changing a free function's signature, the `flowRequest` struct, every CLI verb's flag registration, and inventing the carrier-discriminating branch C1's Selection LBD merely asserts exists "exactly there." | Author implements the flag as a bolt-on check inside the new command binding's constructor only, misses that `flowbind.Registry` still unconditionally builds a `path`-only `Reader`/`Writer`/`Gate` and never even reaches the command binding for entries the old code doesn't expect — or spends a full extra day rediscovering that "the registry constructor" the RDR names does not exist as a method, it is a package-level func called once, shared across verbs. | §1, premortem |
| C-2 | `0025:A11` | A11 is `Pending`, not `Verified`, yet C4's normative contract text presents `Detail string` on `Refusal`, the `*accessor.ExecError` typed-error carrier, and the `refusalOf` single-hook plan as settled fact throughout §Technical Design (`0025:C4`) and Testing Strategy (S3). If `refusalOf`'s call sites don't all cleanly supply-or-nil the error the way A11 hopes, C4's whole "0004 executor unchanged, one additive field" framing collapses and the interface change A11's own "If wrong" clause predicts becomes necessary. | Implementer discovers mid-build that one of `Executor.Gate`/`Executor.Write`'s inline error-handling paths can't route through `refusalOf` without restructuring, and now `Detail` either doesn't populate for a class of refusals the tests expect it for, or the binding interfaces must change after all — contradicting the "0004 executor unchanged" claim baked into the Decision Rationale's scoring. | §2, premortem |
| C-3 | `0025:A12` | A12 is `Pending`. The entire `argv0` path-separator resolution mechanism in C2 ("resolves against the model file's directory") depends on a loader change that has not been verified: neither `table.Model` nor `table.Accessor` carries a source path today (confirmed on `main`), and `flowbind/registry.go::Registry` receives only a `*table.Model`. C2's prose treats "the loader records the model file's directory" as decided, but A12 only proposes to verify this is *possible without disturbing in-memory model construction* — it has not checked whether every in-memory `table.Model` constructor (including test fixtures, other CLI paths, and any future embedder) now needs a new field wired through it. | A separator-bearing `argv0` (the `tools/flowstate-write` example in Illustrative Code itself) resolves correctly in the MVV's hand-built scenario but breaks or silently refuses for any caller that constructs a `table.Model` in memory (embedding, programmatic API, a future SDK) — because A12's fallback for that case is "empty field ⇒ load-time-shaped defect with no category," an ad hoc runtime refusal for a case C5's registered category list was never built to cover. | §2, §3, premortem |
| C-4 | `0025:C6`, `0025:§consequences` | The "one-time vs per-invocation opt-in" trade-off is waved off as an accepted cost, but the RDR never confronts what "per-invocation" means for a CI pipeline, a Makefile, or a wrapped script that calls `intrastate flow` dozens of times per build: every single call site across every CI job, cron, and script must be edited to add `--allow-commands`, forever, with no way to set it once for a session, a repo, or a user. The RDR's own reversal-ledger citations (Consul, Hugo, beads) all eventually moved to a *persistent* gate (a config file, an env var, a repo-scoped flag) — the RDR cites this pattern as validation of the "gate early" instinct while simultaneously choosing the one shape (env-var-less, config-less, per-call flag) none of those projects kept as the end state. | Every team adopting command-backed accessors either (a) papers over the gate with a shell alias/wrapper that always appends `--allow-commands` — silently defeating the "explicit per-invocation friction is the point" rationale — or (b) files a kata for the config file successor within the first sprint of real usage, because CI scripts calling `intrastate` from a dozen places is the median case, not the exception. | §2, §3, premortem |
| C-5 | `0025:A4`, `0025:§implementation-plan §Prerequisites` | RDR 0025 is proposed while its own named prerequisite, RDR 0016 (fail-closed `readerFor`), sits at Draft with `gate_stale=false, gate_written=false` — confirmed live on `main` (`readerFor` in `internal/accessor/model.go` is a bare first-match loop with no ambiguity check). The RDR's own Testing Strategy (S5b) *admits* this is untested/unpinned behavior ("pending 0016's fail-closed uniqueness... This pins the inherited behaviour so 0016's landing is a visible change, not a silent one") — i.e., it plans to ship a write path whose verification correctness depends on an accident of map/slice ordering that a **second, undelivered RDR** will later change underneath it. | A model author declares two readers for the same role (one path-backed, one command-backed) expecting explicit selection; `readerFor` silently picks whichever comes first in registration order, verifies a write against the *wrong* reader, and the write reports success while the actually-intended reader was never consulted — a correctness bug indistinguishable from a passing test until 0016 ships and the reader selection visibly changes underneath existing models. | §1, §3, premortem |
| C-6 | `0025:F6`, `0025:A3` | The stdin-sinking hazard (a write tool like `tee {artifact}` that reads the C3 JSON envelope as its own input and overwrites the artifact with `{"state.phase":"review"}` verbatim) is confirmed by the RDR's own spike and demoted to a "Silent risk" in Failure Modes, discovered only at the *next* read-back with a generic `bad config line 1` parse error — not at the point of corruption. The RDR's mitigation is "the wrapper contract... is the answer for every stdin-reading tool," i.e., a documentation convention with no lint check, no schema check, and no runtime guard preventing a model author from declaring exactly this footgun. | A model author, following the Illustrative Code's own write example almost literally but pointing `command` at a raw tool instead of a wrapper, silently corrupts a shared state artifact (a git config file, a JSON file another tool also reads) in a way that is invisible until the *next* invocation of an unrelated capability fails with a cryptic parse error days later — the diagnosis path in F6 says "surfaces only at read-back," meaning the blast radius (who else reads that artifact) is unbounded and untracked by the RDR. | §1, premortem |

## 1. The three most likely ways implementation goes wrong

### (a) The gate site is not actually one seam, and the flag doesn't reach it cleanly

**Root cause.** C6 asserts: "The registry builds bindings once from the model (`flowbind/registry.go`), so a constructor that returns a refusing binding when the flag is unset gates every present and future caller by construction." This sentence describes a world where `flowbind.Registry` is a method with a constructor argument. It is not. Grounding against `main`:

- `internal/cli/flowbind/registry.go::Registry` is `func Registry(m *table.Model) accessor.Registry` — a free function, no receiver, no flag parameter today.
- It has exactly **one** non-test call site in the whole CLI: `internal/cli/flow_exec.go:830`, inside a `flowRequest{registry: flowbind.Registry(model), ...}` struct literal.
- That single `flowRequest.registry` value is then read by `flow_state.go`'s set-state path (`req.registry.Lookup`, `req.registry` consumed by `accessor.NewExecutor(req.registry, ...)` at line 306) *and* by `flow_exec.go`'s own two `NewExecutor` call sites (lines 256 and 644).

So the real shape is: one registry-building call site, feeding a shared struct, consumed by three different executor-construction call sites across two files and (implicitly) whatever CLI verbs route through `flowRequest`. C6's claim that "the flag therefore reaches the registry constructor, which is the one new parameter this clause adds" is not a one-line diff — it requires (1) adding a flag to `Registry()`'s signature, (2) plumbing that flag from CLI flag-parsing (a `cobra`/pflag `--allow-commands` bool that must be registered on however many verb definitions actually call into `flowRequest`) through to the `flowbind.Registry(model)` call, (3) inventing the carrier-discriminating branch inside `Registry()` that does not exist today (today it unconditionally does `Reader{Path: acc.Path}` for every entry — there is no `if acc.Command != nil` anywhere), and (4) deciding what the *command binding's* constructor does with a flag that was decided one layer up, in a completely different package, at model-load time, not at binding-construction time.

**Specific passage.** `0025:C6`: "The gate sites **in the command binding's constructor**, not in the executor... The flag therefore reaches the registry constructor, which is the one new parameter this clause adds."

**Symptom.** Implementer either (a) burns a day discovering "the registry constructor" is a stateless free function called once and shared, forcing a redesign of how the flag is threaded that the RDR gave zero guidance on, or (b) ships a version where the flag is checked in one code path (say, `flow_state.go`) but not another (`flow_exec.go`'s two paths), because the RDR's own words ("no execution path can bypass the gate by reaching the executor another way") assert a guarantee the actual call graph does not structurally provide without work the RDR never scoped.

### (b) Two Pending assumptions (A11, A12) are load-bearing for contracts written as if already true

**Root cause.** RDR authoring convention distinguishes `Verified`, `Pending`, and `Unverified` for exactly this reason: a `Pending` assumption is one whose truth has not been established. Yet C4 (Detail field, ExecError, refusalOf hook) and C2 (argv0 path-separator resolution against "the model file's directory") are both written in fully normative, declarative voice — no hedging, no "if A11 holds" qualification — even though A11 and A12 are both stamped `Pending` in the Critical Assumptions section, with explicit "If wrong" clauses that each describe a real possibility of blowing up the design's central claims ("the stderr tail needs an interface change... contradicts C4's '0004 executor unchanged' framing"; "argv0 resolution must fall back to a constructor parameter... or separator-bearing argv0 is dropped from v1").

The RDR's own Finalization Gate template explicitly calls this out as a lock-blocking condition: "List any record whose Method is `Docs Only`... and any that remain `Pending` or `Unverified` with a plan to verify before implementation begins... no assumption marked `Pending` or `Unverified` may have settled-fact prose elsewhere in the RDR depending on it." A11 and A12 are Pending, and C2/C4 are exactly the settled-fact prose the gate rule forbids depending on unresolved assumptions. This is not a hypothetical risk; it is the RDR failing its own template's stated lock condition, visible on a single read of the Critical Assumptions section against the Normative Contracts.

**Specific passage.** `0025:A11` ("Status: Pending... Verify by reading `refusalOf` and every one of its call sites") and `0025:A12` ("Status: Pending... Verify by reading the load site and enumerating in-memory `table.Model` constructions"), both contradicted by the confident normative voice of `0025:C4` and `0025:C2`.

**Symptom.** At implementation time, either verification confirms the happy path and nothing changes (best case, but then why was C4/C2 allowed to lock as if this were already known?), or verification fails and the implementer discovers — after the contract has already been treated as fixed, tests written against it, and review sign-off obtained on the *assumption* that these were settled — that the "0004 executor unchanged" claim central to the Decision Rationale's scoring (row: "Reversibility... can add placeholders later") is false, forcing a re-litigation of the whole binding-interface shape after code is already written against the old contract.

### (c) The stdin-corruption hazard (F6) is real, spike-confirmed, and mitigated only by documentation

**Root cause.** A3's own spike (`evidence/spikes/a3-a7-a8-real-tools/`) explicitly ran the `tee {artifact}` scenario (R4b) and confirmed it overwrites the artifact with the raw JSON envelope, silently, with exit 0 — the failure only surfaces at the *next* read-back attempt, with a generic parse error (`bad config line 1`) that gives the operator zero clue what actually happened. The RDR is honest about this in Failure Modes (F6: "Silent risk"), but the "mitigation" offered anywhere in the record is prose: "The wrapper contract (envelope in, tool-native invocation out) is the answer for every stdin-reading tool" — a documented convention model authors are expected to follow, backed by **zero** load-time or lint-time enforcement. C5's six defect categories (`command_and_path_conflict`, `command_empty`, `command_unknown_placeholder`, `command_shell_interpreter`, `command_output_shape`, `command_env_conflict`) do not include anything that could catch "this command reads stdin and also happens to write the artifact it was told to write to" — there is no static way to know that `tee {artifact}` sinks stdin into the named path, because `{artifact}` is a positional CLI argument, not a stdin schema declaration, and the tool's I/O behavior is not visible from its argv at all.

**Specific passage.** `0025:F6`: "a write tool that sinks its stdin into the artifact ingests the C3 envelope as content — the A3 spike's `tee {artifact}` overwrote the artifact with `{"state.phase":"review"}` and exit 0; it surfaces only at read-back... The wrapper contract... is the answer for every stdin-reading tool."

**Symptom.** A model author who reads the Illustrative Code's write example (`tools/flowstate-write` — a wrapper, correctly) and then, for a *different* role's write entry, reaches for a raw established tool instead of writing a wrapper (because the RDR's Approach section says reads and gates can bind directly and doesn't loudly warn that writes categorically cannot for any tool that touches stdin) silently corrupts a real artifact — possibly one another accessor or another tool outside intrastate also reads — with no signal until an unrelated future read fails with a message that does not mention "command," "stdin," or "wrapper" anywhere.

## 2. The section that will be rewritten within 6 weeks of shipping

**`§Normative Contracts / C6` (the execution gate).** This is the highest-confidence prediction in this critique. The RDR itself supplies the evidence against its own design: the "reversal ledger" it cites approvingly (Consul script checks "flipped off-by-default and re-gated twice more," Hugo's "deny-by-default `security.exec.allow`," beads "refusing repo-persisted commands pending a trust gate") is a list of projects that **all eventually needed a persistent, non-per-invocation gate** — an allowlist file, a config surface, an environment variable — because a flag that must be retyped on every single invocation is not a security control real users tolerate; it is a control they routinely defeat with a wrapper script, a shell alias, or a `Makefile` target that always passes it. C6 even names this exact successor need explicitly: "the successor that adds the file inherits an already-gated seam." The RDR is not wrong that v1 needs *some* gate before a config surface exists — A10's absence-of-config-surface finding is real and well-grounded — but the specific shape chosen (flag only, no env var, no repo-scoped opt-in, no session state) guarantees that the first real CI adoption story generates a support kata within days, and that kata's fix will be "add a config file or an env var," i.e., a rewrite of C6, not an amendment. The RDR's Consequences section even names the cost honestly ("per-invocation friction... a scripted caller repeats the flag") without registering that this is not a tolerable steady state, only a v1 stopgap dressed as a design decision.

## 3. The assumption that will not survive first contact with a real user

**A3's framing that "established-tool writes... are served by a thin declared wrapper... and that wrapper class is small and mechanical, not bespoke per integration."** The spike backing this (`evidence/spikes/a3-a7-a8-real-tools/`) tests exactly one established tool: `git config`. The wrapper it produces (`wrapper-write.sh`, "8 lines") is indeed small for `git config`'s specific shape: stdin JSON, loop over keys, shell out to `git config --file {artifact} <key> <value>` per key. But `git config`'s write semantics (idempotent per-key CLI invocation, no read-modify-write race, plaintext key=value target) are close to the best possible case for "wrap an established CLI into a JSON-in/JSON-out envelope." The RDR generalizes from this single, favorable data point to "the wrapper class is small and mechanical, not bespoke per integration" as if this were a proven property of the *design*, when it is an artifact of choosing a well-behaved example tool. A real user's first non-`git-config` established tool — a database CLI needing a connection string and a transaction, a Kubernetes `kubectl patch` needing a JSON-merge-patch body, a Terraform-adjacent tool needing state locking, or literally any tool whose write is not "one flag per key" — will need a wrapper that is not 8 lines, that has its own error handling, its own partial-failure semantics, and its own idempotency story that the C3 envelope (flat string map in, nothing more) cannot express. The RDR's own Approach section admits "its body is script the model does not carry" as an accepted Consequence, but frames it as a bounded, minor cost rather than confronting that for anything beyond the single spiked tool, "small and mechanical" is an assertion with an n of one.

## 4. The Premortem

*Six weeks after `intrastate` 0025 ships, filed as a retrospective.*

The team adopted command-backed write accessors to finally close the "declared decision but hand-applied edit" gap RDR 0025 was built to fix. The first production model declared a `write.branch_state` entry delegating to a small internal CLI (`ourtool state set`) that reads its target path from a positional argument and its new value from a second positional argument — not from stdin. Because C3 mandates stdin-JSON as the *only* channel for per-invocation write data beyond `{artifact}`, and `ourtool` doesn't read stdin at all, the model author did what A3's spike explicitly demonstrated was the pattern: wrote a five-line shell wrapper (`tools/state-wrapper.sh`) that reads the JSON off stdin, `jq`-extracts the one key intrastate happens to write for that role, and calls `ourtool state set {artifact} <value>`.

Three weeks in, someone added a *second* owned tag to the same write role, expecting the wrapper to "just work" because the model's `keys` list grew. It didn't: the wrapper's `jq` extraction only ever looked at the first key, silently dropping the second tag's planned value. Read-back (`internal/accessor/executor.go::Write`, verified through `readerFor`) came back `read_back_mismatch` for the dropped key — correctly caught, per C4's design — but the on-call engineer's first instinct was to blame the *executor*, because nothing about `0025:C3`'s stdin contract or the wrapper's own five lines were reviewed as part of the model-review process; the model diff that added the second `keys` entry looked identical to every other reviewed model change, and the wrapper script lived in a different repo path that the model-review checklist didn't touch. The fix took two hours once found, but the finding took a day, because the failure mode — "a shell wrapper hand-written to match today's key set silently ignores new keys" — is not a class C5 can catch (no lint rule inspects wrapper *bodies*, by design: `0025:§Approach` explicitly says "its body is script the model does not carry"), and Failure Modes never named this specific silent-drop shape, only the `tee {artifact}` stdin-sink case.

Separately, and more seriously: a different team, migrating an existing path-backed reader/writer pair to a command-backed pair for a `repo` role, discovered mid-migration that they now had *two* readers registered for the `repo` role — the old path-backed one (not yet deleted, kept as a rollback fallback) and the new command-backed one. `readerFor` (`internal/accessor/model.go`, unchanged by 0025, confirmed first-match-only on `main`) silently picked whichever came first by map/slice iteration order under the hood of `flowbind.Registry`'s `slices.Sorted(keys(...))` — which is deterministic by *name*, not by intent — and their write's read-back verified against the **old, stale, path-backed reader** for two release cycles before anyone noticed the command-backed write's actual effects weren't being checked at all. The RDR's own S5b explicitly predicted this exact scenario and left it "unpinned until 0016" (`0025:A4`), shipping anyway on the theory that 0016 would arrive and make the behavior a "visible change." 0016 was still Draft when 0025 shipped, exactly as `0025:§Implementation-Plan §Prerequisites` recorded it might be, and nobody tracked the successor RDR closely enough to notice its landing (or non-landing) mattered to production correctness of an already-shipped feature.

Finally, the CI team, needing `--allow-commands` on every one of forty-odd `intrastate flow` invocations across their pipeline's YAML, did exactly what the reversal-ledger precedents the RDR itself cites did eventually: they wrote a wrapper Make target, `state-apply: intrastate flow set-state --allow-commands $(ARGS)`, that always passes the flag. The "explicit per-invocation friction" C6 relies on as its safety property evaporated on day one of real adoption, exactly as Consul's and Hugo's own histories — cited approvingly in the RDR's Decision Rationale as validation — predicted it would.

## 5. Acceptance tests that would have caught each failure at RDR-review time

```gherkin
Feature: RDR 0025 review-time acceptance checks

  Scenario: Registry threading for --allow-commands is concretely specified
    Given C6 claims "the flag reaches the registry constructor"
    When the reviewer asks "which function signature changes, and how many
         call sites of flowbind.Registry and NewExecutor must be edited"
    Then the RDR must name flowbind.Registry's actual signature
         (a free function with one non-test call site at flow_exec.go:830)
    And must show the flag threading from CLI flag parse through
         flowRequest through Registry() through binding construction
    And must not describe a struct-constructor pattern that does not exist
         in internal/cli/flowbind/registry.go

  Scenario: No normative contract depends on a Pending assumption
    Given the Finalization Gate's own Assumption Verification rule states
         "no assumption marked Pending... may have settled-fact prose
         elsewhere in the RDR depending on it"
    When A11 (ExecError/Detail carrier) is Status: Pending
    And C4 states the Detail-field mechanism in fully normative voice
    Then lock must be blocked until A11 resolves to Verified
    And the same check must apply to A12 (model file directory) vs C2

  Scenario: Stdin-sink write hazard has a structural, not documentary, guard
    Given the A3 spike already demonstrated `tee {artifact}` silently
         corrupts the artifact when a write command reads stdin as content
    When a model author declares a write `command` entry
    Then intrastate lint must be able to flag or the RDR must explain
         why it categorically cannot flag "a raw non-wrapper tool bound
         to a write capability" as at minimum an advisory
    And the record must not rely solely on author-followed convention
         ("the wrapper contract... is the answer") for a hazard its own
         spike reproduced

  Scenario: Two-reader ambiguity is not shipped as "unpinned until a
             successor RDR lands"
    Given readerFor is confirmed first-match-only on main with no
         ambiguity detection
    And RDR 0016 (the fix) is Draft, not locked, at 0025's proposal time
    When a model declares two readers for one role (one path-backed, one
         command-backed) during a migration
    Then 0025 must either block command-write read-back until 0016 lands,
         or add its own interim ambiguity refusal, rather than shipping
         a documented-but-unenforced "S5b pins first-match" test that
         passes today and silently changes meaning under 0016

  Scenario: --allow-commands does not get wrapped away on first CI contact
    Given the reversal ledger cited by the RDR shows every peer eventually
         needed a persistent (non-per-call) opt-in
    When a CI pipeline needs to invoke command-backed accessors more than
         once
    Then the RDR must state explicitly why a wrapper Make target or shell
         alias that always passes --allow-commands does not defeat C6's
         stated security property
    And if it cannot, C6 must be scoped as an admitted stopgap with a
         tracked successor kata, not presented as a settled design decision

  Scenario: Wrapper "small and mechanical" claim is tested against more
             than one tool
    Given A3's spike verifies exactly one established tool (git config)
    When the Approach section asserts wrapper scripts are "small and
         mechanical... not bespoke per integration" as a general property
    Then at least one additional established tool with non-trivial write
         semantics (multi-key transactional write, non-idempotent CLI,
         or non-positional target) must be spiked
    Or the claim must be scoped explicitly to "verified for git-config-like
         tools only"
```

