Model: claude-opus-5[1m]

# Persona 1 — Product Manager (iteration 2)

Delta-scope re-run over the rewritten draft. Scope read: `0019:§background`,
`0019:§problem-statement`, `0019:§approach`, `0019:MVV`,
`0019:§activation-step-1-contract-and-reference-docs`,
`0019:§decision-rationale`, `0019:§consequences`. Widened once — to
`0019:C1` (carrier-scope clause), `0019:S10`–`0019:S12`, `0019:BR4`/`BR5`,
`0019:F3`, `0019:A6`/`0019:A7` — because the delta passages cite them as the
owners of the boundary they now qualify, and the boundary-agreement question
cannot be answered without reading the owner text.

## Findings

**Zero findings.** The rewrites did not open a new gap in the PM lens. Every
delta passage grounds out against source, and the three iteration-1 PM
concerns are resolved rather than relocated.

## Verification performed

### `§background` Priority justification — grounded, not overclaimed

Each factual clause checked against `models/rdr.toml`:

- "no `class` key ⇒ `ClassStateMachine`" — confirmed. The parsed TOML has no
  `class` under `[model]` or top level, and
  `internal/table/model.go:53-56` documents the empty `Model.Class`
  reading as `ClassStateMachine`.
- "three `provenance = "owned"` tags" — confirmed: `[tags.stage]`,
  `[tags.status]`, `[tags.gate_passed]`, all `required = true`,
  `single_valued = true`.
- "an `[initial]` block seeding all three" — confirmed:
  `stage = "seeded"`, `status = "draft"`, `gate_passed = "false"`.
- "bound to a file-backed `[write.rdr-status]`" — confirmed by the
  discriminator that actually decides it. The accessor declares only
  `role`/`path`/`keys`/`timeout`/`read_back`; with no `edit` block and no
  `command`, `internal/cli/flowbind/registry.go:78-80` falls through both
  branches (`len(acc.Edit) != 0`, `commandBacked(acc)`) to the JSON
  `&Writer{}`. So it is file-backed by the same type-switch `0019:A7`
  specifies the verb will use — the Background and C1's carrier scope agree
  on this model by construction, not by assertion.

The claim is therefore not overclaimed in the other direction. I looked
specifically for the failure mode the brief names — a Background that now
overstates blocking-ness — and instead found the record *under*-claiming its
own evidence. See the next item.

### `Priority: Medium` is consistent, and the supporting evidence is stronger than the record says

The Background rests Medium on "the repo's own gate model is a live blocked
first run, not a hypothetical one." That is true, and there is an in-repo
witness the record does not cite:
`internal/cli/flow_mvv_0011_test.go::mvvRdrArtifact` drives
`models/rdr.toml` operationally through `flow set-state` + `flow next`, and
its seeding step is:

```
"flow", "set-state", "--model", model, "--artifact", bind,
"--write", "stage=resolved", "--write", "status=draft",
"--write", "gate_passed=false"
```

Two of those three values (`status=draft`, `gate_passed=false`) are
byte-identical to the model's `[initial]` block. That is the hand-typed
transcription `§problem-statement` describes, committed in shipped test code
against the shipped model. The wall is not merely reachable on the gate
model — it has already been walked around by hand, in-tree.

I checked the counter-hypothesis before accepting Medium: that `rdr.toml` is
a lint subject only. It is CI's lint subject
(`.github/workflows/ci.yml:97`, `intrastate lint --model models/rdr.toml`;
no `flow` verb appears in `ci.yml` or `Makefile`), and no `rdr.status`
artifact exists in the tree. But `0011`'s MVV harness is an operational flow
consumer, so the "nobody operates this model" objection does not hold.
Medium is warranted. This is an optional strengthening, not a finding — the
record's existing justification already carries the Priority.

No passage anywhere in the record still asserts the old "nobody is blocked"
premise. Swept all sections for `hypothetical`/`nobody`/`no operator`/
`not blocked`/`no live`: the only two `hypothetical` hits are
`§background` ("not a hypothetical one" — the new claim) and `§consequences`
("a live scope gap, not a hypothetical one" — the carrier gap, grounded
below). Both are in the Medium direction; no stale Low residue.

### MVV step 10 demonstrates the user outcome and is consistent with steps 1-9

Step 10 is executable as written and is not redundant with the fixture:

- It runs against a model that needs no additional artifact binding beyond
  role `rdr`. I checked whether `recognized` would force a second binding:
  it will not — `recognized` is argv-supplied
  (`internal/table/model.go:16-19`, `RecognizedTagKey`), with no accessor of
  its own. `[read.rdr-status]`/`[write.rdr-status]` share role `rdr` and
  path `rdr.status`, so one `--artifact rdr=<fresh>` is the whole binding.
- Its `flow next` assertion is correctly key-scoped, and this is the detail
  that keeps it consistent with `0019:F3`. Every rule in `rdr.toml` carries
  a `[rule.match.recognized]` atom, and `recognized` is not in `[initial]`,
  so after a successful seed `flow next` will *still* emit
  `{recognized, absent}` on every candidate via the atom walk
  (`internal/cli/flow_next.go:339-353`). Step 10 says "no
  `unknown[].reason: absent` **for them**" — scoped to the three seeded
  keys. Fixture step 4 uses the same construction ("no
  `unknown[].reason: absent` entry **for either key**"). The two steps state
  the assertion the same way, and both agree with F3's disclosed scope gap
  (init's claim covers exactly the `[initial]` key set). Had either dropped
  the qualifier it would assert something false; neither does.
- It adds coverage the fixture cannot: the fixture's two keys do not span
  the `bool` kind, and `rdr.toml`'s `gate_passed` does. The record states
  this reason explicitly rather than presenting step 10 as a repeat.

Steps 1-9 prove mechanism on a fixture; step 10 proves outcome on the
shipped model. That division is the right one for the PM question and it is
the division the record claims.

### The scope qualifiers agree; they are not two statements of one boundary

The brief asks whether `§problem-statement`, Phase 2, and `0019:C1` state
the boundary differently. They do not — and the two quantifiers that look
like a conflict are over different domains:

- **Seeding predicate** (`§approach`, `C1`, RT table): ALL bound
  *artifacts* must be empty. Stated identically in all three
  ("EVERY bound artifact carries NO key", "the quantifier is ALL, not
  per-artifact").
- **Carrier refusal** (`C1`): ANY bound *write accessor* being
  non-file-backed refuses.

Both point the conservative direction — seed only when everything is clean,
refuse when anything is unsupported — so there is no arm where the two
disagree about what the verb does. `0019:S11` exists specifically to pin the
ALL quantifier with a discriminating control, and `0019:S12` pins the torn
case; the multi-accessor domain is tested, not just asserted.

`§problem-statement` and `§consequences` phrase the carrier boundary with a
singular accessor ("flows whose bound write accessor is edit-carried…",
"A model binding an edit-carried… write accessor"), where C1 uses "any of
whose bound write accessors". I considered flagging this and did not: every
model in the tree (`models/rdr.toml`, `models/examples/release-grammar.toml`,
`models/examples/review-state-machine.toml`) binds exactly one write
accessor, C1 is the named owner of the boundary and both prose passages cite
it, and S11 covers the plural case. This is prose number agreement, not a
boundary stated two ways, and amending prose to match would be editorial.

### The qualifier is not missed anywhere that still promises the unqualified outcome

Swept the full record for unqualified outcome promises
(`every state-machine flow`, `every flow`, `all flows`, `the first run of`).
Two hits, both correct:

- `§problem-statement`: "The first run of every state-machine flow needs a
  manual `set-state`" — this describes the **status quo wall**, which *is*
  universal today. Qualifying it would make it false.
- `§consequences` positive bullet: "the first run of a state-machine flow
  becomes one explicit command" — indefinite article, not universal, and the
  fourth negative bullet in the same list carries the full FILE-BACKED
  disclosure with the 0025/0028 citations and the successor pointer.

`0019:ALT1`'s con ("every state-machine flow's first run stays a hand-typed
duplication") is likewise a status-quo description under the rejected
alternative, not a promise made by the chosen approach.

I verified the `§consequences` scope-gap bullet's own claim that the two
excluded carriers are "shipped and Implemented": `0025` and `0028` both
carry `Status: Implemented`, and both writers exist in code
(`internal/cli/flowbind/edit.go`, `internal/cli/flowbind/registry.go:78-80`
routing to `NewEditWriter` and `cmdbind.Writer`). The bullet claims the
*carriers* are shipped — which is true — not that models use them. No
in-tree model does, which makes the gap disclosed-and-currently-unpopulated.
The record does not overstate this either way.

### Phase 2 discovery-affordance constraint is implementable against real surface

The constraint requires naming the verb from `flow next`'s
`unknown[].reason: absent` report, on the ground that `0019:BR4` and
`0019:BR5` close both automatic-invocation routes. Both rejections confirmed
as read (external wrapper → "ambient by construction"; auto-init on first
`set-state` → "implicit again, just relocated"), so the premise holds: the
verb is typed by hand, and discoverability is load-bearing for the user
outcome rather than a nicety. The insertion point exists —
`internal/cli/flow_next.go:100-106` already carries the sentence "it leaves
the row a candidate with the key listed under unknown and its reason named,
so you can see what to supply," which is exactly where "what to supply" for
`[initial]` keys is `flow init-state`. The constraint is actionable, not
aspirational.

## Iteration-1 PM concerns — disposition

All three resolved by the rewrite, none relocated:

1. **Priority understated / Background asserted no live blocker.** Fixed.
   `Priority` is Medium, and `§background` now grounds it on the shipped
   gate model with the four model facts that make it a member of the served
   population. Verified above against `models/rdr.toml`.
2. **MVV proved mechanism but not the user outcome.** Fixed by step 10,
   which runs the shipped model and asserts the wall is gone on it. Verified
   executable and correctly key-scoped above.
3. **Outcome promised unqualified while C1 scoped it to one carrier.** Fixed
   in both directions — `§problem-statement` carries the carrier-scope
   qualifier, `§consequences` carries the disclosed scope-gap bullet with a
   successor pointer gated on `0019:A6`, and Phase 2 constrains the shipped
   doc text to state FILE-BACKED scope rather than "every state-machine
   flow."

## Verdict

PASS, non-blocking. The record delivers the user outcome it claims, states
its boundary once and consistently, and proves the outcome on a shipped
model rather than only a fixture. No decision is blocked on the PM lens.

The two `Pending` assumptions in the delta's orbit (`0019:A6` emptiness
carrier, `0019:A7` binding type-switch exhaustiveness) are Stage 4 work and
are correctly scoped as approach-level and contract-level respectively; they
are not PM-lens blockers and I do not re-report them here.
