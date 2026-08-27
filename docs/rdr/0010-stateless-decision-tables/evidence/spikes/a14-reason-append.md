Model: claude-opus-5[1m]

# A14 — appending `no-participating-dimension` to 0006's closed `reason` set

**Assumption**: "`no-participating-dimension` can be appended to 0006's closed
`reason` set without breaking `TestReq80`'s exhaustive assertion or any other
consumer of `graphlint.Reasons()`, and 0006 accepts the append."

**Verdict: PASS** (all three sub-questions verify).

Baseline before any change: `go test ./internal/graphlint/...` → `ok
github.com/newcoinc/intrastate/internal/graphlint 3.397s`.

Method note: committed source was NOT modified. The append was exercised on a
scratch copy of the tree at
`<scratchpad>/a14/repo`. `git status --porcelain` over
`internal/` in the real repo is empty throughout.

---

## 1. MECHANICAL — the set is declared closed AND append-only; no consumer breaks

### The declaration (verbatim)

`/Users/cwensel/sandbox/newcoinc/intrastate/internal/graphlint/taxonomy.go:56-63`

```go
// The closed, append-only `reason` set `graph-unprovable-coverage`
// carries (Technical Design, reason table). It is carried on that code
// alone.
const (
	ReasonDimensionNotFinite = "dimension-not-finite"
	ReasonTagNotSingleValued = "tag-not-single-valued"
	ReasonRowCanRefuse       = "row-can-refuse"
)
```

`taxonomy.go:87-92`:

```go
// reasons is the closed `reason` discriminator set.
var reasons = []string{
	ReasonDimensionNotFinite,
	ReasonTagNotSingleValued,
	ReasonRowCanRefuse,
}
```

`taxonomy.go:103-104`:

```go
// Reasons returns the closed `reason` discriminator set.
func Reasons() []string { return slices.Clone(reasons) }
```

The word **append-only** is present in the shipped source comment, not merely in
the RDR. That is the load-bearing phrase: the set is closed against *arbitrary*
values, and explicitly open to *appends*.

Note also `Reasons()` returns `slices.Clone(reasons)` — a fresh slice per call,
so no consumer can alias or mutate kernel state, and a longer slice cannot
corrupt a caller.

### Every consumer of `graphlint.Reasons()`

`grep -rn 'Reasons()' internal/` returns exactly two families. Only one is
graphlint's:

| Consumer | file:line | Breaks on APPEND? |
| --- | --- | --- |
| `Reasons()` definition | `internal/graphlint/taxonomy.go:104` | No — returns the longer set |
| `TestReq80_...ClosedSet` | `internal/graphlint/findings_0006_test.go:519` | **Yes, mechanically** — sorted literal at `:518`; see §2 |

**Not a consumer**: `internal/resolve/guard.go:81 func Reasons() []Reason` is a
*different, unrelated* closed set (`ReasonAbsent`/`ReasonUncomparable`, 0007:C8),
pinned at exactly two members by `internal/resolve/guard_atoms_test.go:804`.
It shares only the identifier name. `internal/graphlint` neither imports nor
references it, and 0011's Overrides explicitly keeps that seam at
`absent`/`uncomparable`. An append to graphlint's set cannot touch it.

### Every consumer of the individual reason constants

Production code (non-test), all in `internal/graphlint/coverage.go`:

| Site | Shape | Breaks on APPEND? |
| --- | --- | --- |
| `coverage.go:120` | `Reason: reason` (from `unprovableReason`) | No — assigns, does not enumerate |
| `coverage.go:162,166` | `ReasonTagNotSingleValued` literal | No |
| `coverage.go:180-211` `unprovableReason` | returns one of two constants | No — a producer, total over its own arms |
| `coverage.go:235-247` `unprovableMessage` | `switch reason { case ReasonTagNotSingleValued: ... default: ... }` | No — has a `default`, so a new value is *total*; it would phrase the dimension message, which is why C5's arm must supply its own message rather than route here |
| `coverage.go:277` | `Reason: ReasonRowCanRefuse` | No |
| `internal/graphlint/engine.go:109` | `f.Reason` in `identityKey` join | No — opaque string in a sort tuple |
| `internal/cli/clierr` `Finding.Reason` | serialized field | No — untyped string, `omitempty` |

No production switch over the reason set is exhaustive-without-default, and Go
has no exhaustiveness check on string constants. No JSON golden/testdata file
pins the reason set (`testdata/` contains no `reason` matches).

Test-side literal references (`internal/cli/lint_fixtures_0006_test.go:295,320`;
`internal/graphlint/coverage_0006_test.go:805,999,1246,1302,1305,1310,1864`) all
assert that a *specific* finding carries a *specific* reason. Each is an
existential assertion about one arm; none enumerates the set. An append leaves
every one of them true.

**Conclusion §1**: exactly ONE consumer in the tree is sensitive to an append,
and it is the conformance test that is *designed* to be updated with one.

---

## 2. TEST — `TestReq80` asserts both halves; the update is mechanical and sufficient

`internal/graphlint/findings_0006_test.go:511-543` (verbatim):

```go
// REQ-80: "**`graph-unprovable-coverage` carries a `reason`
// discriminator.** … The finding MUST carry a stable `reason` from this
// closed, append-only set" — `dimension-not-finite`,
// `tag-not-single-valued`, `row-can-refuse`.
// ASSUMPTION-4: the discriminator is carried only on this code.
// BOUNDARY
func TestReq80_UnprovableCoverageCarriesAReasonFromTheClosedSet(t *testing.T) {
	want := []string{"dimension-not-finite", "row-can-refuse", "tag-not-single-valued"}
	got := slices.Clone(graphlint.Reasons())
	slices.Sort(got)

	if !slices.Equal(got, want) {
		t.Errorf("reason set = %v; want exactly %v", got, want)
	}

	// Every emitted finding of this code carries a reason from the set,
	// and no OTHER code carries one.
	r := lint(t, optionalGuardDecls, multiDefectBody)
	for _, f := range r.Findings {
		if f.Code == graphlint.CodeUnprovableCoverage {
			if !slices.Contains(want, f.Reason) {
				t.Errorf("%s carries reason %q, outside the closed set %v",
					f.Code, f.Reason, want)
			}
			continue
		}
		if f.Reason != "" {
			t.Errorf("%s carries the reason discriminator %q; the "+
				"discriminator is scoped to %s alone",
				f.Code, f.Reason, graphlint.CodeUnprovableCoverage)
		}
	}
}
```

Confirmed — it asserts **three** things, exactly as A14 states:

1. **Sorted-literal exhaustiveness** (`:518-524`): `slices.Sort(got)` then
   `slices.Equal(got, want)` against a sorted literal. Order-insensitive on the
   declaration side, so an append at any position in `taxonomy.go` works; the
   literal must be updated in **sorted** position.
2. **Every emission carries a member** (`:529-535`): walks all findings of
   `CodeUnprovableCoverage` and requires the reason be in the set. This is the
   half that will cover C5's new arm once it emits — and the half that fails
   outright on a `""` reason, exactly as `0010:C5` states.
3. **Scoping** (`:537-541`): no other code carries a reason. Unaffected by an
   append.

### The update is mechanical and sufficient — demonstrated

Applied to the scratch copy: append the constant + slice member in
`taxonomy.go`, and update the one literal at `findings_0006_test.go:518` to
sorted position (`no-participating-dimension` sorts second).

Diff (2 files, 6 insertions, 4 deletions — the extra lines are `gofmt`
re-aligning the const block):

```
 const (
-	ReasonDimensionNotFinite = "dimension-not-finite"
-	ReasonTagNotSingleValued = "tag-not-single-valued"
-	ReasonRowCanRefuse       = "row-can-refuse"
+	ReasonDimensionNotFinite       = "dimension-not-finite"
+	ReasonTagNotSingleValued       = "tag-not-single-valued"
+	ReasonRowCanRefuse             = "row-can-refuse"
+	ReasonNoParticipatingDimension = "no-participating-dimension"
 )

 var reasons = []string{
 	ReasonDimensionNotFinite,
 	ReasonTagNotSingleValued,
 	ReasonRowCanRefuse,
+	ReasonNoParticipatingDimension,
 }

-	want := []string{"dimension-not-finite", "row-can-refuse", "tag-not-single-valued"}
+	want := []string{"dimension-not-finite", "no-participating-dimension", "row-can-refuse", "tag-not-single-valued"}
```

Result — `go build ./...` OK, and `go test ./...` **fully green**:

```
ok  	.../docs/rdr/0010-stateless-decision-tables/evidence/spikes/a13-zerodim	0.308s
ok  	.../internal/accessor	1.632s
ok  	.../internal/cli	1.306s
ok  	.../internal/cli/clierr	0.984s
ok  	.../internal/cli/flowbind	0.422s
ok  	.../internal/graphlint	3.688s
ok  	.../internal/guard	0.866s
ok  	.../internal/resolve	1.305s
ok  	.../internal/table	1.523s
```

Targeted: `TestReq80_UnprovableCoverageCarriesAReasonFromTheClosedSet` PASS,
`TestReq81_RowCanRefuseMessageReadsAsAWithholdingNotAnError` PASS.

No other assertion in the file or the suite fails. The two-line append plus the
one-literal edit is **mechanical and sufficient**.

> Transient note for the record: an earlier run of the scratch copy showed a
> failure in the `a13-zerodim` spike package. It reproduced with the patch
> **reverted**, and was traced to the spike directory being untracked
> (concurrent sibling-agent work), so `git checkout` in the copy could not
> restore it — a stale-copy artifact, not a regression. Refreshing that
> directory cleared it; the green run above is post-refresh.

---

## 3. OWNER ASSENT — structural; 0006's own text anticipates appends

### 0006's status

`rdr inspect --select 0006:§metadata 0006`:

- **Status**: `Implemented` (not merely Final — the code exists and ships)
- **Overrides**: `None`
- **Seam Lineage**: `no prior accretion`

So 0006 is locked. The question is whether an append requires *reopening* it.

### Where the reason set is declared normative

The reason set is **not** in any numbered `C`-contract. I checked the plausible
candidates: `0006:C8` (the narrowing that mints
`graph-unprovable-coverage`, no `reason` mention), `0006:C13` (the finding
payload contract — lists code/model/severity/message/ids and the atom fields,
but **not** `reason`), `0006:C17` (closes the *advisory tier*, unrelated), and
`0006:BR5` (rejects a second code / warning tier, unrelated).

The declaration lives in the **Technical Design**, `§technical-design`, the
paragraph introducing the reason table
(`docs/rdr/0006-graph-lint-authority-and-guarantees.md`, TD reason table;
projector-relative TD lines 247-251), and is carried as **REQ-80**:

> **`graph-unprovable-coverage` carries a `reason` discriminator.** One code, three
> triggers, and each has a *different* remedy, so a machine consumer branching on
> `code` alone cannot tell the author what to do. The finding MUST carry a stable
> `reason` from this closed, append-only set, naming the key, dimension, or row at
> fault:

followed by the three-row table (`dimension-not-finite` / `tag-not-single-valued`
/ `row-can-refuse`, each with its remedy column), and:

> `row-can-refuse` is *not* a model defect — it is an honest withholding — and the
> message MUST say so rather than reading as an authoring error.

Restated in `docs/rdr/0006-.../artifacts/req-list.md:217`:

> - [REQ-80] "**`graph-unprovable-coverage` carries a `reason` discriminator.** …
>   The finding MUST carry a stable `reason` from this closed, append-only set,
>   naming the key, dimension, or row at fault" — the set is `dimension-not-finite`,
>   `tag-not-single-valued`, `row-can-refuse`. — (REASON)

### Does that authorize an append without amending 0006?

**Yes. The assent is structural.**

The operative phrase is **"this closed, append-only set"**. Read against 0006's
own vocabulary elsewhere, the two adjectives do distinct work, and the contrast
is decisive:

- Where 0006 means a set that may **not** grow, it says so and names the members
  as the closure: `0006:C17` — "The advisory tier is **closed at**
  `graph-coverage-closed-by-escape`, `graph-redundant-row`,
  `graph-unreachable-rule`, and `graph-vacuous-atom`." No "append-only". Same in
  the source: `taxonomy.go:37-38` — "The tier is CLOSED at these four
  (`0006:C17`)."
- Where 0006 means a set that may grow, it adds **append-only** — and it uses
  that word in exactly this sense elsewhere for structures explicitly designed to
  grow: `0006:C14` ("an **append-only** typed `findings` field owned by
  `clierr`") and `0006:D-wire-byte-format` ("JSON mode carries findings as
  **append-only** structured data").

`closed` here means *no value outside the set is legal* (a finding may not invent
a reason — which is what TestReq80's second half enforces). `append-only` means
*the set may be extended, never reordered, narrowed, or have a member's meaning
changed*. An append is therefore the **sanctioned** evolution, and a
change/removal would be the amendment. 0006 anticipated growth and pre-authorized
its only safe direction.

Two corroborating signals:

- The shipped source carries 0006's intent forward verbatim in the comment at
  `taxonomy.go:56` — "The closed, **append-only** `reason` set". The implementer
  of 0006 read it the same way.
- The reason table's own framing — "One code, three triggers, and each has a
  *different* remedy" — grounds the set in *remedy distinctness*, not in the
  number three. `no-participating-dimension` supplies a fourth trigger with a
  genuinely distinct remedy ("author the discriminators as guard atoms"), which
  is precisely the axis the set was built to expose. Borrowing
  `dimension-not-finite` instead would state the remedy "declare the domain" for
  a dimension that does not exist — the failure mode `0010:C5` and the critique
  pass both call out.

**This is not an amendment.** No 0006 sentence becomes false: the three existing
members keep their spellings, triggers, remedies, and order; REQ-81's
`row-can-refuse` clause is untouched; `reason` stays scoped to
`graph-unprovable-coverage` alone. 0006 need not be reopened. **No BLOCKER.**

The one honest caveat: 0006's TD reason table will, after the append, enumerate
three of four members. That is inherent to any append-only set owned by a locked
record and is exactly what the per-project rule "we never amend rdr content …
code is the source of truth" contemplates — `taxonomy.go` is the source of
truth for the set; 0006 gives the history and the growth rule. 0010's Overrides
metadata is the durable record of the fourth member.

### Precedent

`0010`'s own Overrides already books this correctly and explicitly
("**appended, not conditioned**"). For *peer* precedent of a downstream RDR
extending a locked upstream vocabulary without reopening it:

- **`0008`** (Status: Implemented) — Overrides: "**extends** RDR 0002's
  tag-declaration surface with a name constraint on recognized-provenance
  declarations and **one additional data-level validation category**
  (`reserved tag key`, within its 'including at minimum' extensible list)."
  This is the closest analogue: a downstream Implemented RDR adding a member to
  an upstream category list on the strength of the upstream's own
  extensibility marker.
- **`0011`** (Status: Final) — Overrides: "**adds the reason token**
  `not-evaluated` for un-run gate ids on that CLI list only — `0007:C8`'s seam
  vocabulary stays at `absent`/`uncomparable`." A downstream RDR adding a reason
  token, while explicitly *declining* to widen the closed kernel set next to it —
  showing the project already distinguishes an authorized append from a
  prohibited one, and books each accordingly.

Both confirm the house pattern: extend where the owner marked the set
extensible, record it in Overrides, do not reopen the owner.

---

## Verdict

| Sub-question | Verdict |
| --- | --- |
| 1. MECHANICAL — closed+append-only declared; consumers survive an append | **PASS** |
| 2. TEST — TestReq80 asserts both halves; update mechanical and sufficient | **PASS** |
| 3. OWNER ASSENT — 0006's "append-only" pre-authorizes; no reopen needed | **PASS** |

**Overall: PASS.** A14 verifies. `no-participating-dimension` can be appended to
`graphlint.Reasons()` with a two-line source change and a one-literal test
update; the full suite stays green; and 0006's assent is structural rather than
something that must be sought, because 0006 declared the set append-only by
design and reserved the append as its sanctioned growth path.

Residual work belongs to `0010:C5`'s emission arm (site + message), covered by
A13 — not to A14. In particular the new arm MUST supply its own message rather
than routing through `coverage.go::unprovableMessage`, whose `default` would
otherwise phrase it as "declare the domain".
