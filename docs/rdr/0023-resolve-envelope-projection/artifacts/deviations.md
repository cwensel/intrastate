# Deviations — RDR 0023 Resolve-envelope projection opt-out

Phase 2 implementation record. Every entry carries exactly one Type.

---

## D-1 — The MVV step-1 golden cannot carry the test's absolute model path

**Type:** TEST-FIXTURE
**Status:** resolved — fixture corrected, implementation continues
**REQs:** REQ-95, REQ-96, REQ-97, REQ-107, REQ-129

**What the Phase 1 test does.** `flow_mvv_0023_test.go::TestReq27And28And95And96And97And129_…`
compares the full emitted line of `mvvCall0023` against the checked-in
golden byte-for-byte. `mvvCall0023` builds its `--model` argument from
`pricingModelPath(t)`, which is `repoRootFor(t) + "/models/examples/pricing-decision-table.toml"`
— an ABSOLUTE path, because the package's tests run with `cwd = internal/cli`
and a relative model path would not resolve.

**The defect.** `resolvePayload.Model` carries `modelRef`, which is the
`--model` argument VERBATIM (`flow_input.go::selectModel` returns `path`
unchanged). A golden captured from that invocation therefore embeds the
capturing machine's checkout root, e.g.
`{"model":"<abs-checkout>/models/examples/pricing-decision-table.toml",…}`.
Such a golden:

1. fails on every other checkout, including CI, for a reason that has
   nothing to do with the clause it guards; and
2. commits an absolute local filesystem path into a checked-in artifact.

**Why this is TEST-FIXTURE and not SPEC-DEFECT.** The RDR fixes the
comparison, not the argv. `0023:MVV` step 1 and `0023:REQ-107` name the
invocation as "over the checked-in decision-table fixture
(`models/examples/pricing-decision-table.toml`)" — the RELATIVE spelling —
and the whole evidence base measures it that way:
`evidence/spikes/a1-byte-width.md` row S1 records
`--model models/examples/pricing-decision-table.toml` at 290 B, and
`evidence/spikes/a2-encoder-mechanism.md` §Reference output records the same
relative spelling in the reference bytes the RDR calls normative. The
absolute path is the Phase 1 harness's accommodation of the test cwd, not a
clause.

**Resolution (derived from the evidence base).** The golden is captured with
the RELATIVE model path the RDR and both spikes name, and the golden
comparison normalizes the checkout root out of the observed line before
comparing. Byte-identity is still asserted over every byte of the record,
including the `model` value in its RDR-normative relative spelling; only the
checkout-specific prefix is folded. The assertion is not weakened: a change
to any payload byte — the `model` key's presence, its relative value, or any
other field — still fails.

**Ordering discharged.** The golden was captured from the PRE-CHANGE binary
on a clean tree and committed BEFORE any edit to `resolvePayload`
(`0023:A-8`, REQ-96). See the first Phase 2 commit.
