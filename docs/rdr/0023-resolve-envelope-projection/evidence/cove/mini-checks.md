Model: claude-opus-5[1m]

# Stage-5 mini-check cue read — cli/0023

The RDR's first lens pass owes the cue read (`$RDR_HOME/stages/05-prelock.md`).
Four of five cues fire; each fired cue's table is written into the RDR itself.

| Mini-check | Cue | Fires | Table location in the RDR |
| --- | --- | --- | --- |
| source-authority census | derived/propagated output — `readers` and `owned` were adjudicated ECHO at A7 *because* they are derivable, which is the ambiguity the census names | yes | Technical Design § Source-authority census |
| test-discriminability | S3/S4 are structural-absence oracles (S4 asserts a non-registration and the RDR itself names a vacuity guard); S1 passes by non-difference | yes | Technical Design § Oracle discriminability |
| round-trip / fidelity | no import/export, parse/deparse, compose/decompose, migrate/rollback. Projection is key DELETION, explicitly not re-encoding, and has no inverse (F7 "recovery" is re-running without the flag, not an inverse). Byte-identity here is equality-with-default, which S1 owns, not round-trip fidelity | no | — (correctly absent) |
| disposition | drops an input class (the echo group) and sets exit outcomes (`command-error`, exit 2 on the other verbs) | yes | Technical Design § Disposition |
| desk trace | C1's key-set/order/byte-identity clauses, C2's always-keep core and completeness rule, and S1–S5's Expected lines all bear on one output surface; the MVV walks it stepwise and a normative fixture exists | yes | Validation § Desk trace |

## Desk-trace outcome

No CONTRADICTION row. Every MVV step has a concrete witness from the
normative fixture (the pricing 2×2 call in
`evidence/spikes/a2-encoder-mechanism.md`) or from a live run on
repo HEAD.

Step 3's text witness was taken live rather than reasoned:

```
intrastate flow resolve --model models/examples/pricing-decision-table.toml \
    --tag tier=paid --tag region=eu --outcome decide --as text
```

returns 15 alphabetically-sorted lines; projection removes exactly the 6
echo lines (`model`, `observed.region`, `observed.tier`, `outcome`,
`owned`, `readers`), leaving 9. Strict subset, byte-identical per line.

## Findings the trace surfaced

**M-1 (0023:C2, class c — internal contradiction against shipped source).**
C2's always-keep clause read "no mode may OMIT" `rule, escaped,
escape_class, revision`, and C1 restated `escape_class`'s presence rule as
"omitempty-on-unescaped". Both are false against the producer:
`internal/cli/flow_resolve.go::escapeClassOf` returns `""` — rendering the
key ABSENT — when an escaped row declares several classes and the probe
cannot name one it rescues. Its own doc comment states this is deliberate:
"a row declaring several and rescuing an unprobeable kind reports nothing
rather than guessing." So an escaped plan can legitimately carry no
`escape_class` today, and an implementer following C1 would write an oracle
asserting presence whenever `escaped` is true, which fails on shipped
behavior.

Grounding: JDR 0002 §D1's own rationale for always-keep is that a projection
"can never launder a rescued plan into an ordinary one" — a
projection-INVARIANCE claim, not an unconditional-presence claim. The field
that actually carries the rescued/ordinary distinction is `escaped`, which
is non-`omitempty` and unconditionally present. So the fix is a wording
correction inside the settled decision, not a reopening of it.

Disposition: **fixed**. C1 now defers to the producer's own presence rule and
requires presence-rule fields be asserted by comparison against the same
run's default output rather than an unconditional literal; C2 restates
always-keep as projection-invariance and names `escaped` as the guard's real
basis. Amendment sweep over `escape_class` (10 sites): S2 Expected and MVV
step 2 re-worded to "iff the default carried it"; the Decision Rationale
safety row changed "never omissible" → "never PROJECTABLE". Sites at 108,
328, 371, 879, 941 are group-membership lists or already projection-phrased
and need no change.

Not this RDR's to fix, noted for the record: the shipped 0005 oracles assert
`escape_class` present on an escaped plan but only exercise single-class
escape rows (`flowEscapeModel` models `no_match`), so the producer's
documented multi-class-unprobeable path is untested. That is a 0005 gap, not
a 0023 one — this RDR neither widens nor narrows it.

**M-2 (0023:S5, class d — silence, resolved in place).** The draft states
the wire keeps declaration order (C1) and that text lines are a subset, but
never says the two orderings differ. Text is alphabetically sorted by the
shipped `flatten` (`respond/text.go`, `sort.Strings`). The subset property
holds under either ordering because projection only deletes keys, so this
was a legibility gap rather than a defect. Disposition: **fixed** — the desk
trace's closing paragraph names both orderings and why no clause requires
them to agree.
