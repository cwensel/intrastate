# Flight — batch:rdr-0024

`/kata-flight --label batch:rdr-0024 --drain` · 1 wave · **7 shipped, 0 stopped, 0 skipped**

## Shipped

| # | kata | merged | what landed |
|---|---|---|---|
| 1 | `m6fz` | `ec7efdc` | `tableHeader` canonicalizes quoted rule headers (`[["rule"]]`, `[['rule']]`) so they join `emitRuleLine`'s census; also fixed the dotted-key guard running *before* unquoting |
| 2 | `wkqs` | `985bd42` | one `decodeScalarString` (decoder round-trip, not a hand-rolled unescaper) serves both `scalarAssignment` and `unquoteKey`; `stripComment` honours escaped quotes; ADV-1's no-id leg pins the exact header line |
| 3 | `5hy1` | `8ff701f` | four oracle legs converted from source-text proxies to structural observables; category floor gains its missing completeness direction |
| 4 | `jrp2` | `436795e` | member byte-exactness fixture + `Dispositions != nil` on the flat carrier path |
| 5 | `hpnv` | `0da7d97` | REQ-16 disposition-table arms for `int`/`scalar`; non-`verdict` arms made single-defect via `armSource()` |
| 6 | `4v92` | `3d507bb` | `TestReq22_0024`'s narrowing leg actually narrows (5→3), widened baseline rather than a narrower domain |
| 7 | `av9m` | `7e5c44d` | 0023 MVV dispositions fold made exact on the value, not just the position; populated values now report against REQ-97 |

Ship order was **dependency-aware, not priority-ordered**: `m6fz` is a production
defect `5hy1`'s weak oracle could not catch, so it shipped first and `5hy1`'s
strengthened oracle became its regression test. `wkqs` followed `m6fz` because
both alter what the ordinal census sees. The final two touch disjoint packages
and were resolved in parallel, then merged serially.

## Spin-offs

One `KATA_PUSH` during the wave: `m6fz`'s refine found `unquoteKey` does no TOML
escape decoding (`[["rule"]]` still misses the census). It was **not** fixed
in place — it is the header-key twin of `wkqs` Leg A, and both are the single
unescaper decision D12 defers, so fixing half the seam would have left the
matched half open. Pushed to `wkqs` (already in the wave), severity raised
low→medium on attach, and it shipped as part of that kata's `decodeScalarString`.

No spin-off was filed onto a standing selector, so the drain re-sweep terminated
empty after one wave, as designed.

## Scope review (ran before the wave, per the default gate)

All 7 IN-SCOPE. No merges, no closes, no RDR-seeds. Two escalation questions
both resolved *against* escalating, on evidence rather than counters:

- **Seam-accretion tripwire overruled.** `5hy1`'s seam carries EIGHT closed
  point-fixes, which by count routes to `kind:rdr-seed`. Grounding showed the
  eight are three unrelated failure modes (vacuous control flow / assertion
  strength / wrong observable); only the third is `5hy1`'s root cause, and it is
  its first crisp statement. Neither candidate unifying contract can be stated
  once. Eight point-fixes at a **category**, not a **seam**.
- **Scanner-class fork closed.** `wkqs`/`m6fz` both pre-committed to seeding an
  RDR on a third "hand-written TOML scanning misses a legal spelling". The count
  threshold was met (four spellings) but the shape threshold was not: the class
  is bounded by TOML's finite key/string syntax set, and the fork's alternative
  is already refuted — `Position()` exists only on `*DecodeError`
  (`pelletier/go-toml/v2@v2.2.4/errors.go:71`), which a cleanly-decoding
  document never produces. That is REQ-33, verified at the module source.

## Corrections the flight produced

Three plan errors were caught by verification rather than carried through:

1. **`hpnv` Leg B's RED recipe was unsound.** Reordering `run`'s step slice is
   inert — `checkRuleEmit` returns early on the empty `EmitDecls` map that
   `loadEmitDecls` populates, so reversing them makes the check a no-op rather
   than flipping the category. Replaced with a direct probe.
2. **`jrp2` Leg B's cite had rotted** — the real site is
   `emit_grammar_0024_test.go::TestReq4_0024`, not `emit_carry:151`.
3. **`jrp2` Leg A's mutation was ambiguous** — a pre-existing tag-side test also
   caught it through a shared helper, so the run was repeated with the new
   subtest stashed to prove the 0024 suite was genuinely blind to the gap.

## Verification discipline

Every kata was mutation-verified: change the artifact so the bytes stay and the
structure breaks, confirm the old oracle stays green, confirm the new one goes
red. `make check` (fmt, vet, golangci-lint, build, graph-lint, docs-check,
`go test -race`) exit 0 read from `$?` — never inferred from output text — on
every ship, and re-run on the **rebased** tip for `av9m`.
