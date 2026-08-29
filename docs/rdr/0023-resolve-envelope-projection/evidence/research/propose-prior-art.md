# cli/0023 — Stage 2 prior-art pass (Propose)

Bounded read, not a spike. Budget: ≤3 queries per claim, ≤5 opened hits per
claim. Sources: Domain priors (`docs/cli-output-contract.md`, governing RDRs
0005/0011 via the projector), semble over `/Users/cwensel/sandbox/langref`
(the Go/CLI fallback the resources index names), arc `StateMachineRes`.

## Accepted citations

### P1 — gh ships caller-controlled output reduction on read verbs through one funnel

- Query (semble, `langref/gh-cli`): "AddJSONFlags json fields flag exporting
  selected payload fields" → `pkg/cmdutil/json_flags.go`.
- Quote: `f.StringSlice("json", nil, "Output JSON with the specified
  `fields`")` (`pkg/cmdutil/json_flags.go::addJsonFlag`), beside
  `addJqFlag` ("Filter JSON output using a jq `expression`") and
  `addTemplateFlag`, registered together by
  `AddJSONFlagsWithoutShorthand` / `AddJSONFlags`
  (`pkg/cmd/workflow/list/list.go:65` shows the per-command wiring with a
  command-owned allowed-field list).
- ⇒ The general form is NAMED-FIELD selection with a per-command declared
  field universe — the caller picks fields, the command declares which are
  pickable. Note `gh`'s `--json` doubles as the output-mode switch;
  intrastate already has `--as`, so importing the general form imports only
  its selection axis, plus the obligation to police a field universe.

### P2 — the boolean fixed-projection form (`<result>-only`) has direct peer precedent

- Query (semble, `langref/gh-cli`): "name-only flag print only filenames
  diff" → `pkg/cmd/pr/diff/diff.go`.
- Quote: `cmd.Flags().BoolVar(&opts.NameOnly, "name-only", false, "Display
  only names of changed files")` (`gh pr diff --name-only`, mirroring
  `git diff --name-only`).
- ⇒ A boolean flag that projects the result to one named, command-defined
  part — no caller-supplied field list, no field universe to police — is an
  established peer form for "same computation, smaller report".

### P3 — the governing fences (in-repo, projector reads)

- `0005:C1` fence: text mode "MUST emit human output derived from the same
  verb-specific result" — so any projection must be applied to the
  verb-specific result BEFORE `respond.OK`; a JSON-only omission is
  structurally impossible through `internal/cli/respond/text.go::flatten`
  (one marshal feeds both modes).
- `0005:§technical-design` "Envelope": the minimum `data` shape for
  `resolve` lists `model`, `revision`, `observed`, `owned`, `readers[]`,
  `outcome`, `rule`, `gates[]`, `next`, `writes`, `clear[]`, `escaped`,
  `escape_class` — pinned by `0005:A6` ("minimum success payload fields").
  ⇒ an opt-in absence of some of these must ride an Overrides entry
  narrowing A6's minimum to the default mode.
- `0011:C2` fence: `--all` "MUST NOT be accepted by flow resolve, flow
  read-state, or flow set-state — satisfied by NON-REGISTRATION",
  structurally pinned (`internal/cli/flow_all_0011_test.go::TestReq46And47And65And98_AllIsAbsentFromTheOtherThreeVerbsAndTheFlowGroup`).
  ⇒ the fence names `--all`, not "any flag on resolve"; a new flag on
  `resolve` trips no shipped oracle. The reusable idiom (Lookup != nil,
  vacuity-guarded) is the shipped absence-oracle house form.
- `0011:D-naming`: rejected `--select` / `--match` spellings "because the
  defect is the default and every skill would carry the flag" — a ground
  that does NOT bind here (the default echo is agreed correct; only chained
  calls carry the flag). Reopening recorded deliberately per `0011:BR3`'s
  own instruction ("deliberate rather than rediscovered").
- `internal/cli/flow_resolve.go` extended help "Reading a successful plan"
  lists exactly `rule`, `next{}`, `writes{}`, `clear[]`, `emit{}`,
  `escaped`, `escape_class`, `gates[]` — the verb's own vocabulary already
  partitions the payload into "the plan" and the request echo.
- `0024:C4` (PROPOSED peer, same envelope): appends `dispositions` last,
  after `emit`, "never omitted"; joined by `Plan.RuleID` like `emit`.
  ⇒ answer-group member for any projection partition.

## Rejected / no-coverage branches

- **DMN-engine instance claim** ("a decision-engine evaluation result
  carries output entries only; inputs are not echoed") — ⚠ no prior-art
  coverage. Two arc queries on `StateMachineRes` ("DMN decision evaluation
  result contains output entries only, inputs not echoed back"; "camunda
  evaluate decision REST API response result variables") returned no
  on-topic hit (top hits: awf-cli error-code docs, inngest API README).
  The claim is NOT leaned on; the choice rests on P1–P3. If wanted later it
  is a Resolve read, not a Propose ground.
- Sibling-path check (step 5, rg over `internal/cli`, non-test): the only
  conditional success-payload fields are producer presence rules —
  `flow_resolve.go::resolvePayload.EscapeClass` (`escape_class,omitempty`,
  present only when an escape row fired) and
  `flow_next.go` candidate `Gates` (`gates,omitempty`, present only when
  `--evaluate-gates` ran them). No caller-controlled projection of computed
  content exists anywhere in the CLI; nothing to reuse, no parallel signal
  invented.
