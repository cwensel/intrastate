# Deviations — RDR 0023 Resolve-envelope projection opt-out

Phase 2 implementation record. Every entry carries exactly one Type.

No entry is `needs author decision`. Every gap below is resolved by the
record's own evidence base, and each entry names the evidence it rests on.

| id | Type | Status | Subject |
| --- | --- | --- | --- |
| D-1 | TEST-FIXTURE | resolved | the MVV golden cannot carry the harness's absolute model path |
| D-2 | TEST-FIXTURE | mechanical translation | `payloadKeyOrder` reported string VALUES as keys |
| D-3 | TEST-FIXTURE | mechanical translation | the reader read-back left `read-state`'s second role unbound |
| D-4 | TEST-FIXTURE | mechanical translation | the flag-blindness `gates` arm bound no artifact |
| D-5 | TEST-FIXTURE | resolved | MVV step 6 ran a release-grammar shape the 40% bar was not measured on |
| D-6 | SPEC-DEFECT | resolved | `llms.txt` is not an artifact any flag's usage string reaches |
| D-7 | IMPL-DECISION | recorded | the reader seam ships as a nil-by-default package hook |

---

## D-1 — The MVV step-1 golden cannot carry the test's absolute model path

**Type:** TEST-FIXTURE
**Status:** resolved — fixture corrected, implementation continues
**REQs:** REQ-95, REQ-96, REQ-97, REQ-107, REQ-129

**What the Phase 1 test does.** `flow_mvv_0023_test.go::TestReq27And28And95And96And97And129_…`
compares the full emitted line of `mvvCall0023` against the checked-in
golden byte-for-byte. `mvvCall0023` builds its `--model` argument from
`pricingModelPath(t)`, which is
`repoRootFor(t) + "/models/examples/pricing-decision-table.toml"` — an
ABSOLUTE path, because the package's tests run with `cwd = internal/cli` and
a relative model path would not resolve.

**The defect.** `resolvePayload.Model` carries `modelRef`, which is the
`--model` argument VERBATIM (`flow_input.go::selectModel` returns `path`
unchanged). A golden captured from that invocation therefore embeds the
capturing machine's checkout root. Such a golden (1) fails on every other
checkout, including CI, for a reason that has nothing to do with the clause
it guards, and (2) commits an absolute local filesystem path into a
checked-in artifact.

**Why TEST-FIXTURE and not SPEC-DEFECT.** The RDR fixes the comparison, not
the argv. `0023:MVV` step 1 and REQ-107 name the invocation as "over the
checked-in decision-table fixture (`models/examples/pricing-decision-table.toml`)"
— the RELATIVE spelling — and the whole evidence base measures it that way:
`evidence/spikes/a1-byte-width.md` row S1 records
`--model models/examples/pricing-decision-table.toml` at 290 B, and
`evidence/spikes/a2-encoder-mechanism.md` §Reference output records the same
relative spelling in the bytes the RDR calls normative. The absolute path is
the Phase 1 harness's accommodation of the test cwd, not a clause.

**Resolution.** The golden is captured with the RELATIVE model path the RDR
and both spikes name, and the golden comparison folds the checkout root out
of the observed line before comparing. Byte-identity is still asserted over
every byte of the record, `model` included, in its RDR-normative relative
spelling; only the prefix that differs between two checkouts of the same
commit is folded. A change to any payload byte — the key's presence, its
relative value, or any other field — still fails.

**Ordering discharged.** The golden was captured from the PRE-CHANGE binary
on a clean tree and committed BEFORE any edit to `resolvePayload`
(`0023:A-8`, REQ-96), as the first Phase 2 commit.

---

## D-2 — `payloadKeyOrder` reported top-level string VALUES as keys

**Type:** TEST-FIXTURE
**Status:** mechanical translation
**REQs:** REQ-21, REQ-121, REQ-122

The helper walked `data`'s tokens tracking delimiter depth and appended
every `string` token seen at depth 0. A top-level string VALUE sits at the
same depth as a key, so `model`'s path, `revision`'s `""`, `rule`'s id and
`outcome`'s tag were all reported as key names. The default and projected
readings then differed in ways that had nothing to do with key order — the
very property the helper exists to expose. Observed reading:
`[model <path> revision "" observed owned readers outcome decide rule paid-eu …]`.

The helper now reads a key, then CONSUMES its value with
`json.RawMessage`. The projected order was already correct before the fix
(`revision, rule, gates, emit, next, writes, clear, escaped`), which is
`0023:A-1`'s recorded witness and the desk trace's step-2 bytes. No
assertion was weakened.

---

## D-3 — The reader differential's read-back left `read-state`'s second role unbound

**Type:** TEST-FIXTURE
**Status:** mechanical translation
**REQs:** REQ-35, REQ-37, REQ-40–REQ-44, REQ-46

`TestReq35And37And40…`'s `readBack` closure invoked `flow read-state` with
only the `state` role bound, and the call refused
`flow-artifact-missing` on role `orphan`.

That refusal is correct behaviour, not a projection defect: `0005:REQ-37`
makes `read-state` the deliberate exception that "runs every declared
reader, because a diagnostic read has no candidate set to narrow by", so the
role `resolve` correctly leaves unbound is one `read-state` requires. The
0005 MVV (`flow_mvv_0005_test.go`) binds both roles to the same artifact for
this reason, and `flow_next_0005_test.go::TestReq37_…` pins the refusal as
the intended discriminating pair.

The read-back now binds `orphan` as well. The assertion it carries — the
artifact is unchanged across a projected `flow resolve` — is untouched.

---

## D-4 — The flag-blindness oracle's `gates` arm bound no artifact

**Type:** TEST-FIXTURE
**Status:** mechanical translation
**REQs:** REQ-47, REQ-48, REQ-49, REQ-50, REQ-53

`TestReq47And48And49And50And53_…`'s `gates` arm ran `dtGateModel0010` with
an empty artifact binding. Both of that fixture's gates declare
`role = "state"`, so the run refused `flow-artifact-missing` ABOVE the gate
run the arm exists to witness.

Bound to a fresh artifact, as the shipped 0010 oracle binds it
(`decision_table_0010_test.go:526`) — the gate's verdict is its declared
path's suffix (`flowbind.go::verdictFor`), not the artifact's content, so a
fresh artifact is sufficient and is the established pattern.

---

## D-5 — MVV step 6 ran a release-grammar shape the 40% bar was not measured on

**Type:** TEST-FIXTURE
**Status:** resolved — fixture corrected against the recorded baseline
**REQs:** REQ-103, REQ-104, REQ-105, REQ-109

**What happened.** Step 6's `release-grammar/begin` shape ran against an
UNSEEDED artifact and with no observed `--tag`s. It refused
`flow-owned-state-unavailable` — the `begin` row guards `phase = idle` and
writes `build-id`, so both owned keys must be established before the row is
decidable and there is no plan to measure at all.

**Why the bar was at risk even after seeding.** Seeding alone yields
388 B → 238 B = 38.7%, BELOW the 40% bar — not because the projection saves
too little, but because the measured shape is not the one the bar was set
against. `evidence/spikes/a1-byte-width.md` §S2b records the baseline
invocation verbatim:

```
intrastate flow resolve --model models/examples/release-grammar.toml \
    --artifact release=release-idle.artifact \
    --outcome build --tag risk=0 --tag 'checks=[]' --as json
```

with "owned state seeded to `phase=idle`, `build-id=none`" and the row
recorded at 412 B → 238 B, 42.2%. The two observed tags are echo-group
content; dropping them shrinks the DEFAULT side only and depresses the ratio
on a shape the RDR never measured.

**Why this is not REQ-105's route-back.** REQ-105 routes back when "the plan
group is carrying materially more than A1 measured". The plan group here is
238 B — the spike's recorded projected width to the byte. Running the
spike's own invocation reproduces `412 B → 238 B, 42.2%` EXACTLY against the
relative-path spelling. Nothing about the plan group moved; the fixture was
measuring a different call.

**Resolution.** Step 6 now seeds `phase=idle, build-id=none` and passes the
spike's two observed tags, per MVV step 6's own instruction to "compare
against A1's baseline in `evidence/spikes/a1-byte-width.md`".

---

## D-6 — `llms.txt` is not an artifact any flag's usage string reaches

**Type:** SPEC-DEFECT
**Status:** resolved — the obligation REQ-117 states is asserted instead
**REQs:** REQ-116, REQ-117, REQ-118

**The clause.** REQ-116 states that "`flowResolveExtendedDesc` and the
flag's own usage string are inputs to the GENERATED artifacts
`docs/cli-reference.md` and `llms.txt`". The Phase 1 oracle reads that
literally and requires the substring `plan-only` to appear in both files.

**The verified fact.** It cannot appear in `llms.txt`, and no flag does.
`internal/cli/docs.go::writeLLMsTxt` emits a fixed prose header, then one
`- \`CommandPath\`: Short` line per command, then two static link lists. It
renders NO command's flag set and consults neither `Usage` nor the extended
help body. Confirmed by sweep over the committed file: `--all`, `--outcome`,
`--artifact`, `--model` and `--evaluate-gates` each occur zero times; the
only `--as` and `--help-all` occurrences are in the hand-written prose
block. `docs/cli-reference.md` is the artifact that renders flags, and it
carries `--plan-only` with its authored usage string after `make docs`.

**Why the premise is the defect and not the implementation.** Satisfying the
literal reading would require either changing `resolve`'s `Short` to mention
one flag — misdescribing the verb in a command index — or teaching
`writeLLMsTxt` to render flags, which is a NEW public surface on a generated
artifact that the RDR's Normative Contracts do not name and that would put
one verb's flag in an index carrying no others. Both are worse than the
clause is worth, and REQ-94's "reuse unchanged" posture toward shipped
rendering infrastructure argues against the second.

**Resolution, from the RDR's own statement of the obligation.** REQ-117
states the real requirement: the registering commit "MUST also carry the
regenerated `docs/cli-reference.md` and `llms.txt`, or `make check` fails on
`docs-check` at that commit". `docs-check` compares the committed files
against a fresh render. The oracle now asserts exactly that — it re-renders
both artifacts through the shipped `docs` command and requires byte equality
with the committed copies — plus the substring assertion on
`docs/cli-reference.md`, the file that does render flags. That is strictly
stronger than the substring probe for that file, and it is satisfiable and
meaningful for `llms.txt`, which the substring probe was not.

`make check` (including `docs-check`) is green at the documenting commit.

---

## D-7 — The reader-execution seam ships as a nil-by-default package hook

**Type:** IMPL-DECISION
**Status:** recorded — affects how a successor reads the seam
**REQs:** REQ-41, REQ-42, `0023:A-3`

C1 fixes the seam's POSITION and leaves its FORM to the implementer:
"recording execution at the reader invocation site (a counting or recording
seam around the reader pass, the run's own side effects, never a re-call of
the planning function)".

The form chosen is a package-level `readerExecutionHook func(readerID string)`,
nil in production, called through `recordReaderExecution(name)` immediately
before the per-reader `exec.Read` inside `(flowRequest).runReaders` — the one
pass that yields `readers` and `owned` together.

Two properties are deliberate and a successor should preserve them:

- **It is flag-blind by signature.** It takes the reader id and nothing
  else — no flag, no request, no payload. REQ-49 puts the projection flag at
  exactly one lexical site and this is not it; a seam able to observe the
  flag would make the very execution it records flag-dependent, which is the
  vacuous form `0023:A-3` rules out.
- **It records EXECUTION, not intent.** It sits at the `exec.Read` call, not
  in `invokedReaders`. A seam at the name computation would record what was
  planned to run, and comparing that across two runs is green on every
  implementation — including one that runs no reader at all, which is the
  form C1 rejects by name.

Recorded as IMPL-DECISION rather than SPEC-UNDER because C1 explicitly
delegates the form; the seam itself is a named construction obligation
(REQ-84), not an unnamed additive surface. It adds no CLI surface: no flag,
no output field, no error code, and the production call path is one nil
check.
