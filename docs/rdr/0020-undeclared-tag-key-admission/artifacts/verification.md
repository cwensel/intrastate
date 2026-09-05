# Verification — RDR 0020 Undeclared `--tag` key admission

## CoVe pass (Phase 3a) — no violation found

Method: for each REQ in `artifacts/req-list.md`, an adversarial input was
derived **from the record alone** (Phase 1 test files and `coverage.md`
deliberately unread), then executed against a build of `./cmd/intrastate`.
Probes ran against a fixture declaring scalar/enum `tier`, set `labels`,
scalar `note`, bool `flag`, int `count`, and nothing named `extra`/`extras`
— the shape normative fixture F1 specifies — plus
`models/examples/review-state-machine.toml` for the owned-key arm.

**Result: no REQ was observed to be violated. No FAIL-N entries.**

### Normative fixtures reproduced byte-for-byte

F1, F1-array-leg, F2, F3, F4, E and G all matched the record's quoted
strings exactly, including `flow-tag-invalid` message text and `param`.

### Verbatim fidelity (REQ-1/2/4/6/26/30/31/38/51)

Every probe crossed byte-for-byte into `observed`, confirmed by decoding the
payload and comparing to the argv value: unsorted array `["z","a"]` (not
sorted), duplicate members `["a","a"]` (not compacted), spaced array
`[ "a" , "b" ]` (not compacted), embedded quotes, `=` inside the value
(`a=b=c`, ASSUMPTION-6's post-`strings.Cut` remainder), unicode and CJK
(`héllo→世界`), a unicode KEY, leading/trailing whitespace, tab and newline,
HTML-escapable `<a>&"b"`, backslashes `a\b\\c`, a lone space (not empty),
nested JSON objects, arrays of non-strings `[1,2,3]`, a malformed array
`["a",`, and a bare `[`. None were canonicalised, folded, re-encoded, or
refused. Declared set `labels` kept its canonical array form (REQ-5).

### Ordering / precedence (REQ-7..11, 21..24, 29, 36, 37)

Every pairwise collision was probed. Grammar beats all (`noequals`, `=v`).
Reserved `recognized` beats empty value, array value and duplication. Owned
`status` beats empty value, array value and duplication. **Duplicate beats
the hoisted empty-value arm** — `--tag x=a --tag x=` refuses
`flow-tag-duplicate`, confirming the arm sits after the `seen[key]` mark and
the anti-inversion rule of REQ-23 holds. The converse order
(`--tag x= --tag x=a`) correctly refuses empty-value, since the first
iteration reaches the arm before a second occurrence exists.

### Guard correctness (REQ-15..20, 35, 55, 56)

The guard was probed across all three input classes and every declared kind:
declared SET `labels=` keeps the set-specific conformance message (fixture
G, so the unconditional `value == ""` form of REQ-20/56 is not present);
undeclared `extra=` and every declared NON-set kind (`note=`, `flag=`,
`count=`) emit the scalar-shaped "was given an empty value" message.

### Declared-only conformance arms (REQ-11, 14, 33, 34)

Kind/shape/domain arms fire only under a real declaration: `tier=["a","b"]`
and `note=["a"]` refuse "is not set-valued"; `count=99` and `flag=maybe`
refuse on domain. The undeclared near-misses `countt=99999` and
`flagg=maybe` pass and echo. Misspellings of declared keys (`labelz`,
`teir`, and case variants `TIER`/`Labels`) are admitted as carriers,
confirming REQ-3's no-case-folding and REQ-33's accepted open-world cost.

### Non-reach (REQ-13, 27, 28)

`--write` and `--clear` over an undeclared key still refuse
`flow-write-unbound` / `flow-clear-unbound` at the binding check, ahead of
`canonicalValue` — scalar, array and empty values alike. The model-side
closed world still refuses at load: guard atom, `[rule.write]` block,
`[initial]` assignment and accessor `keys` each answer `unknown_tag`; the
clear-list site remains at `internal/table/normalize.go:586`.

### Pins (REQ-12, 39, 40, 41, 44, 52, 57, 58)

`flow-tag-undeclared` appears nowhere in non-test source or docs. The
superseded guarded-lookup comment ("a different decision, under a different
code") is gone. `docs/model-authoring.md` §"Passing a set-valued tag" no
longer asserts the retired refusal and now documents the carrier rule.
Resolving with and without the carrier flags diffs empty across every
projected field — the sole payload delta is `observed` gaining the carried
keys. `go test ./...` is green repo-wide, so no pre-existing green test over
`--tag` admission turned red.

---

## Adversarial review (Phase 3b)

Findings derived independently from the record's Trade-offs / Failure Modes
section, not from the Testing Strategy. The reviewing posture was that
`0020:C1` WIDENS an admission seam, so the questions were what a caller can
now push through that admission previously refused, and whether the
compensating controls the Failure Modes section names actually exist on the
paths where an operator needs them.

Tests live in `internal/cli/flow_carrier_adv_0020_test.go`.

### ADV-1 — the verbatim boundary is bounded by the echo's UTF-8 handling

**Anchor**: `0020:G-cross-cutting` / REQ-6 — "A carried value is admitted
VERBATIM — byte-preserved, with no canonicalisation, folding, normalisation,
or re-encoding — and echoes in the resolve payload's `observed` field as
given."

**Failure mode**: the widening strips every DOMAIN from a carried value — a
declared `enum` or `set` holds its value to `domain`/`elements`, a carrier is
held to nothing — so ARBITRARY BYTES become reachable at this seam from
ordinary argv. "Byte-preserved" is therefore a materially stronger claim
after C1 than before it. The suite's existing REQ-6 table exercises only
printable ASCII (sorting, spacing, case, HTML metacharacters); the byte
classes that actually break a JSON echo were untested.

**Test**: `TestAdv1_0020_CarriedControlBytesAndInvalidUTF8AtTheVerbatimBoundary`
(eight legs: C0 controls, NUL, DEL, an NFD combining sequence, literal
`\uXXXX` text, and two ill-formed UTF-8 sequences).

**Result — PASSES; already handled, with one documented limit.** Control
bytes, NUL, DEL, the NFD sequence and the escape-looking text all round-trip
byte-for-byte: no folding, no normalisation, no re-encoding. Ill-formed UTF-8
does NOT round-trip — `encoding/json` replaces each bad byte with U+FFFD, so
`a\xffb` echoes as `a�b` and two distinct carried values can collapse to
the same echo.

**Scope — this is NOT a regression, and the test says so.**
`TestAdv1b_0020_TheEchosUTF8LimitIsProvenanceBlindAndPredatesTheCarrier`
proves the mangling is provenance-blind: a DECLARED bare `scalar` (which
carries no domain) is mangled identically at the same seam, and the two arms
agree byte-for-byte. C1 widens nothing here; it only makes a pre-existing
encoder limit reachable for a caller's whole undeclared vocabulary. The
finding is a **bound on the record's REQ-6 wording**, not a defect in the
seam C1 moved. Both tests are kept as clearly-labelled regression pins — ADV-1b
in particular guards against a "fix" applied to the carrier branch alone,
which would put the declared and carried arms out of agreement and breach
`0020:G-cross-cutting`'s actual invariant.

### ADV-2 — a surviving refusal must still leave no accessor side effect

**Anchor**: Failure Modes, "Visible" — "every refusal that survives is
unchanged and still fires at exit 2 BEFORE ANY ACCESSOR".

**Failure mode**: "before any accessor" is the clause with a SIDE EFFECT
behind it, and the widening moves the boundary it describes. An invocation
carrying an undeclared ARRAY used to exit 2 at admission with no accessor run
at all; under C1 the same argv is admitted and the verb proceeds to run
readers, gates and WRITERS. An exit-2 admission refusal that had already run
the write accessor would leave the artifact mutated while telling the caller
its input was rejected — the worst shape this widening could take, and one no
exit-code assertion can observe.

**Test**: `TestAdv2_0020_AnAdmissionRefusalLeavesNoAccessorSideEffect` — a
write-bearing verb (`flow set-state --write status=final`) over a seeded
artifact, with the artifact's raw bytes compared before and after, across
five surviving refusals: the hoisted empty-value arm, reserved, owned,
grammar, and duplicate.

**Result — PASSES; already handled.** All five refuse at exit 2 with the
artifact byte-identical. Kept as a regression pin: it is the only assertion
in the suite that would catch admission being reordered behind accessor
execution, which the exit code alone cannot see.

`TestAdv2b_0020_TheAdmittedArrayRunsAccessorsWithoutReachingOwnedState` pins
the accepted cost on the other side: the previously-refusing argv now
succeeds and performs its declared write, and the carried key appears in
neither the persisted owned state nor the write set — the runtime leg of
Decision Rationale (b). **PASSES.**

### ADV-3 — the diagnosis the record leans on is absent from the refusing path

**Anchor**: Failure Modes, the "Silent" and "Diagnosis" clauses, which are
one argument — "a misspelled key (declared or not) passes as a carrier and
the intended rule fails to match; resolution then refuses no-match or routes
to an escape row … Diagnosis: the resolve payload's `observed` field echoes
every carried key byte-for-byte — the stray spelling sits beside the declared
keys **in the same envelope the refusal rides**." The Premortem rests the
whole acceptance of the silent-typo cost on this echo.

**Failure mode**: the Silent clause and the Diagnosis clause describe
DIFFERENT EXITS. The Silent clause's outcome is a REFUSAL (no-match) or an
escape route; the `observed` echo is a field of the SUCCESS payload. The
claim that the stray spelling rides "the same envelope the refusal rides" is
a claim about the REFUSAL envelope — the artifact the operator is holding at
exactly the moment they are stuck.

**Test**: `TestAdv3_0020_TheStraySpellingDoesNotRideTheRefusalEnvelope` —
supplies `teir=free` (a misspelling of the declared, guarded `tier`) against
the fixture model, reproducing the Failure Modes narrative exactly, and
asserts the record's claim against the resulting refusal envelope.

**Result — FAILS against the implementation.** The refusal envelope's
top-level keys are `[code message findings]`. There is no `observed` field
and no path that could populate one: `clierr.CLIError` is
`{code, message, param, detail, hint, findings}` and the resolve payload's
`Observed` field lives only on the success envelope
(`internal/cli/flow_resolve.go:53`). The operator receives an envelope naming
the DECLARED key as absent ("rule `free`: the atom on `tier` could not be
decided (absent)") and carrying **no trace of the stray `teir` they actually
typed**.

**Disposition — record-level, not implementation-level.** This is NOT a
defect in `parseTags`: the admission seam does exactly what C1 says, and no
REQ in `req-list.md` obligates the refusal envelope to carry the echo (REQ-38
is mined from the Failure Modes prose and asserts only that `observed` echoes
every carried key). The gap is between the record's Failure Modes prose and
the shipped envelope shape. The remedy is an amended Failure Modes claim, or
a follow-on that puts the echo on the refusal envelope — not a change to the
carrier branch. Flagged for the orchestrator to route as a deviation or a
follow-on seed.

`TestAdv3b_0020_OnTheEscapeExitTheEchoIsPresentAndCarriesTheTypo` pins the
other exit and makes the finding legible: on the exit-0 escape route the
`observed` echo IS present and carries the typo verbatim, and the escape row
is confirmed selected (`rule: bail`, `escaped: true`). **PASSES.** The echo is
real; it is just on the wrong exit.

### Adversarial verdict

One finding requiring disposition (ADV-3, a record/implementation mismatch on
the Failure Modes diagnosis claim) and one documented limit (ADV-1, a bound
on REQ-6's "byte-preserved" wording that predates this change). The admission
seam itself — the carrier branch, the guarded hoisted empty-value arm, the
refusal precedence, and the accessor boundary — held under every attack
attempted.
