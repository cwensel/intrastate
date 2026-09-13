# deviations — 0029 version promise on machine-readable output

## D1 — the fifth advisory code's identity is implementation latitude

Type: IMPL-DECISION
Status: mechanical translation

The MVV's step 2 commands "Add a new graph-lint finding code at severity
`info` that fires on the model from step 1", and `0029:A1` licenses "a fifth
advisory code" against the amended `0006:C17`. Neither the record nor any
REQ-N names WHICH code. The implementation mints
`graph-idempotent-write` — a row whose write block sets a key back to the
value its own match pattern pinned by equality, reported and never
rejected, exactly the disposition `graph-vacuous-atom` already carries.

It fires on `models/rdr.toml` (rule `seed-revise` matches `stage eq
"seeded"` and writes `stage = "seeded"`), which is the MVV's named clean
model, so step 2's "fires on the model from step 1" is satisfied by the
checked-in fixture rather than by a scratch model.

Recorded because it adds a member to a `growing` vocabulary a consumer
reads, which affects future interpretation of the tier's census. It is not
a new SURFACE: `graphlint::AdvisoryCodes()` is the seam C4 already obliges,
the severity is DERIVED (`severityFor` returns `info` for anything
`IsBlocking` rejects), and no REQ-N is bypassed.

## D2 — the MVV's minor-version movement cannot be observed in one build

Type: SPEC-UNDER
Status: needs author decision → RESOLVED (0029:S7, the record's own Testing Strategy matrix — "the matrix the verified assumptions imply" — contains NO minor-movement row and NO empty-findings row: row 1 asserts `"schema_version":"0.1"` by value and is explicitly additive; row 4 asserts exit 0, `type == "ok"` and `"severity":"info"`. Those are the single-build-witnessable assertions. The MVV oracle at :1546 assigns release-classification to a human and routes the mechanical half to REQ-21/A7's committed CI snapshot (:956); the oracle's step-1 control at :1542-1545 makes the empty-findings baseline a property of the PRE-change build. req-list.md:130 already carries the release-classification clause as EXCLUDED on the same warrant. Re-scope the in-process runner to those assertions — faithful implementation, not a scope change. Phase 3c applies it.)

`internal/cli/schema_mvv_0029_test.go::TestMVV0029_AnAgentPinsAVersionAndSurvivesANewFindingCode`
asserts, in ONE process against ONE build, both:

- step 1: the emitted `schema_version` is exactly `"0.1"`; and
- step 4: the minor component MOVED between the step-1 and step-3 runs
  (`movedMinor == baseMinor` is an error).

`clierr.SchemaVersion` is a single constant, so one build emits one value
on both runs. The two assertions are jointly unsatisfiable in-process.

The record's own text is a TWO-BUILD procedure: step 1 says "Build at the
current HEAD with `schema_version` implemented", step 2 "Add a new
graph-lint finding code", step 3 "Re-run the same command". The `"0.1"` →
`"0.2"` movement C1 describes is therefore a release-classification act
across two builds, and `req-list.md` already carries that half as EXCLUDED
("release-classification rule binding a future release … nothing in-repo
defining a release to range over"). The MVV oracle concedes the same:
"Nothing in-repo defines a release for a test to range over, so the
release-classification half stays human."

The same test also asserts, at step 1, `data.findings == []` on
`models/rdr.toml` while requiring at step 2 that a new advisory code fire
on that same unmodified model — a second in-process contradiction with the
same root cause: one build cannot be both the pre-change and the
post-change binary.

REQ-8 and REQ-42 fix the constant at `"0.1"` and are asserted by three
other tests (`TestReq1And8And42`, `TestReq1And43`, `TestReq1And58`), so
raising the constant to `"0.2"` to satisfy the movement row would turn
those red and contradict the record's quoted value. The implementation
therefore holds `SchemaVersion = "0.1"` and leaves this test's
movement/empty-findings rows red.

This is not resolvable from the record's evidence: it needs the author to
decide whether the in-process MVV runner is re-scoped to the assertions a
single build can witness (the field's presence, value, `type`, exit code,
and the new finding's `info` severity), with the minor-movement and
pre-change-baseline rows recorded as the human release-classification step
the oracle already assigns to a human.

Evidence:
- `0029:C1` — "the minor component still increments on every change"
- `0029:MVV` steps 1-3 and the `oracle` table's step-4 row
- `req-list.md` EXCLUDED entry 2 (release-classification)
- `internal/cli/schema_mvv_0029_test.go:180-182` (`"0.1"` baseline),
  `:190-193` (empty findings), `:287-291` (minor must move)

## D3 — two Phase 1 tests bind the same bytes in `internal/accessor/model.go`

Type: TEST-FIXTURE
Status: needs author decision → RESOLVED (0029:C2 scopes the prohibition to a tiered vocabulary; 0029:S9 pairs internal/accessor/model.go with one site, :308, the verdict set, and names the accessor capability set among sites where `closed` stays; 0029:C4 tiers neither the capability nor the validation-code set. Source is correct; the defect is the file-scoped census sweep at internal/cli/schema_closed_wording_0029_test.go:93-106, which must honour the per-site vocabulary subject the census table carries — per REQ-57, "an unscoped grep for the word is not" the assertion. Fix in that assertion, not model.go.)

Two of this RDR's own Phase 1 tests make incompatible demands of
`internal/accessor/model.go`:

- `schema_closed_wording_0029_test.go::TestReq15And54And55And57_TheClosedWordingIsRetiredAcrossTheEnumeratedCensus`
  reports EVERY non-exempt line containing `closed` in each census file. The
  census pairs `internal/accessor/model.go` with the subject "the verdict
  set (frozen)", but the assertion is file-scoped, not vocabulary-scoped.
- `schema_closed_wording_0029_test.go::TestReq17And63_UntieredVocabulariesAndNonTierSensesKeepTheWord`
  requires two snippets in that same file to SURVIVE verbatim:
  `closed three-member capability vocabulary` (line 46) and
  `closed eight-member validation-code set` (line 131).

Neither snippet matches `exemptLine`, so the first test flags exactly the
two lines the second test requires. No wording satisfies both.

The record is on the second test's side, and says so twice. REQ-63 /
Phase 1 Step 2: "Sites describing an untiered vocabulary, and the non-tier
senses of the word, are out of scope and stay". C2 scopes the prohibition
to "a tiered vocabulary" and S9's census names the SUBJECT per site — for
this file, the verdict set alone. The capability set and the
validation-code set are untiered by C4, so their `closed` stays.

The implementation therefore retired `closed` on every line in that file
that describes a TIERED vocabulary (the kernel's five-kind refusal set at
line 52, the capability count at line 28, the refusal-class set at line 89)
and left lines 46 and 131 as REQ-63 requires. The census test stays red on
those two lines.

The fix is to the census assertion, not to the source: it must honour the
per-site SUBJECT the census table already carries rather than sweeping the
whole file. That is a change to a Phase 1 test, which this phase may not
make, so it is referred.

Evidence:
- `0029:C2` — the prohibition's subject is "a tiered vocabulary"
- `0029:REQ-63` / Phase 1 Step 2 — untiered sites "are out of scope and stay"
- `0029:S9` — the census is a site list paired with the vocabulary that put
  each site in scope; "an unscoped grep for the word is not" the assertion
- `internal/cli/schema_closed_wording_0029_test.go:29-42` (the census and
  its `vocabulary` column), `:93-106` (the file-scoped sweep), `:160-168`
  (the two snippets required verbatim)

## D4 — the added field costs `0023`'s 40% projection-saving bar its margin

Type: DEPENDENCY-LIMIT
Status: needs author decision → RESOLVED (0023:C1 is normative and fixes the measurand as the FULL EMITTED LINE, "not on the `.data` payload alone", calling it the stricter reading because a constant wrapper makes any payload-level saving a smaller proportion of the line. Reading 1 is therefore foreclosed without amending C1; Reading 2 stands — 0023 re-measures its fixtures against the post-0024/post-0029 envelope and records the new floor. Cause correction: of the 41 B per-side growth, only 23 B is `schema_version`; 18 B is pre-existing drift from 0024:C4's non-omitempty `dispositions` (commit e2624d1), landed after A1's 2026-08-29 spike. Main already measured 40.5% — a 0.5-point margin, not 2.2 — so the bar was near-red independent of 0029. Sibling fixtures now sit at 41.7% and 41.3%: any further non-omitempty envelope or plan-group field pushes them red too.)

`internal/cli/flow_mvv_0023_test.go::TestMVV0023_ResolveEnvelopeProjectionEndToEnd`
step 6 requires each checked-in fixture's `--plan-only` projection to save
at least 40% of the full emitted line. With `schema_version` on the
envelope, `release-grammar/begin` measures 38.4% (453 B → 279 B) and the
row goes red.

This is not a stale literal. `schema_version":"0.1"` adds the same ~26
bytes to BOTH the default and the projected line, so it shrinks the RATIO
without changing what either side carries: the numerator (bytes saved) is
untouched while the denominator grows. `0023:A1` measured 42.2-51.3% on
these fixtures pre-change, and `release-grammar/begin` had the thinnest
margin at 42.2% — 2.2 points, less than this field costs it.

The interaction is structural and belongs to both records:

- `0029` puts a non-`omitempty` field on every terminal record by contract
  (C1), and it rides the projected and unprojected lines alike.
- `0023`'s bar is stated on "the full emitted line", envelope included, so
  the envelope's own growth counts against a verb's payload saving.

Nothing in `0029` licenses moving the bar, and nothing in this phase may
edit `0023`'s threshold: the 40% figure is `0023:A1`'s route-back
condition, and the test says so — a landing below 40% "ROUTES BACK rather
than recording a number".

Referred to the author with two candidate readings:

1. `0023`'s bar measures the VERB PAYLOAD (the `data` object), not the
   envelope that carries it. Re-measuring on `data` alone restores the
   pre-change ratio exactly, since the field is outside `data` by
   `0029:C1`, and keeps A1's numbers comparable across the change.
2. The bar stands as written and `0023` re-measures its fixtures against
   the post-`0029` envelope, recording the new floor.

Reading 1 is the better fit for what A1 set out to measure — whether the
plan group carries materially more than the projection — but it restates a
peer record's assertion subject, which this phase may not decide.

Evidence:
- `0029:C1` — the field is on both terminal records and is not `omitempty`
- `0029:G-cross-cutting` — "Field ORDER is not asserted here"; the field is
  additive, not a payload change
- `internal/cli/flow_mvv_0023_test.go:441-457` (the 40% bar, measured on
  `foldCheckoutRoot(emittedLine(...))` — the full line, envelope included)
- observed: `release-grammar/begin` 453 B → 279 B, 38.4% saved
