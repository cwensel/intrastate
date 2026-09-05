Model: claude-opus-5

# A8 — read-side carrier discrimination, source verification (Stage 6 reconcile)

Verdict: part (i) PASS as spelled; part (ii) REFUTED as stated, contract
amended after a strong consult returned PASS on the repair.

## Part (i) — the read-construction set. HOLDS.

`internal/cli/flowbind/registry.go`, read loop:

    var binding accessor.Binding = Reader{Path: acc.Path}
    if commandBacked(acc) {
        binding = cmdbind.Reader{Accessor: acc, Name: name, Config: cfg}
    }

Write loop, for contrast:

    var binding accessor.Binding = &Writer{Path: acc.Path}
    switch {
    case len(acc.Edit) != 0:   binding = NewEditWriter(acc, name)
    case commandBacked(acc):   binding = &cmdbind.Writer{...}
    }

Read side: an initializer plus a single `if` — TWO arms, no third, no shared
type, no `default`. Both VALUES against the write loop's POINTERS.
Receivers confirm it: `func (r Reader) Read(...)` in both flowbind and
cmdbind (VALUE); `(w *Writer) Apply`, `(w *EditWriter) Apply`,
`(w *cmdbind.Writer) Apply` (POINTER). A read-side case spelled
`*flowbind.Reader` matches nothing and would refuse EVERY model — which is
why C1 fixes the spelling normatively.

`EditWriter` appears only in the `m.Writers` loop, declares `CapWrite`, and
has no `Read` method: no third read carrier.

NOT a residue, unlike the write side: `commandBacked(acc)` is
`len(acc.Command) != 0 || acc.Path == ""`, so `flowbind.Reader` is reachable
on exactly a declared path with no command — no predicate the verb must call.
Admitting only `flowbind.Reader` is exactly the file-backed read set.

## Part (ii) — reachability of the resolver. REFUTED AS STATED.

    func (reg Registry) readerFor(role string) (Definition, bool)
    internal/accessor/model.go:249

UNEXPORTED. Only non-test caller: `internal/accessor/executor.go:343`,
inside the declaring package. `internal/cli` cannot call it.

A8 asserted this symbol is "reachable from internal/cli". It is not. This is
the SAME defect A7 found for `flowbind::commandBacked` — a contract naming a
check the verb cannot perform — recurring one package over.

## The repair (consulted, PASS)

`Registry.Definitions` is an exported field; `Definition.Identity`,
`.Accessor`, `.Binding`, `accessor.CapRead` and `table.Accessor.Role` are all
exported. The gate re-derives the role's reader in-package.

Equivalence is EXACT, not approximate. `readerFor`'s body:

    for _, d := range reg.Definitions {
        if d.Identity.Capability == CapRead && d.Accessor.Role == role {
            return d, true
        }
    }
    return Definition{}, false

- Bare FIRST-MATCH (`return` inside the loop), no accumulate-then-pick.
- No `::bound` check: a missing reader and a reader bound under another
  capability both yield the flat `false`.
- No ordering rule beyond `Definitions`' own slice order.
- No normalization: `Identity.Capability` compared to the `CapRead` constant,
  `Accessor.Role` compared with `==` against a value the loader stores
  verbatim (`Role: *a.Role`, internal/table/load.go:999) after only an
  absent-or-empty rejection.

FIRST MATCH IN REGISTRY ORDER is therefore normative, not incidental: read-side
role uniqueness is NOT enforced today — neither `internal/accessor/validate.go`
nor the loader's `::accessorTable` carries a uniqueness arm, and RDR 0016,
which would make role→reader a function, is Draft. Two same-role readers are
constructible, so a gate collecting matches and refusing on ambiguity would be
WIDER than `readerFor`, refusing models the executor serves. Taking the last
would likewise disagree.

## No-accessor-invoked property: preserved, and visibly stronger

`Binding`'s only method is `Capability()`; `Read` lives on the separate
`ReadBinding`. The filter touches `Identity.Capability` (a field on the
Definition, not the binding) and `Accessor.Role` (a table field). The
subsequent type switch inspects a dynamic type tag and calls nothing.
So the gate learns the carrier from struct fields and a type tag: zero
writes, zero binding invocation. This is what makes S10's assertion true by
construction rather than by timing.

## Option rejected

Exporting `readerFor` as `ReaderFor` would serve equally but is an
`internal/accessor` API change this RDR does not authorize — the same ground
on which C1 already declines to export `commandBacked`, and a strictly larger
blast radius (a permanent API member with one external caller, inviting future
callers to bypass the executor's read path) for zero additional capability.

## Blast radius of the repair

Dependents stay total: per needed role the gate is a trichotomy — reader found
and `flowbind.Reader` (admit), reader found and `cmdbind.Reader` (carrier
refusal naming capability=read), no reader (the unbound-needed-role arm C1
already defines, reached identically before and after).
No exported symbol added, no signature changed, no file outside docs/rdr/0019
edited. Peers: 0004 owns `internal/accessor` and is strictly better off
(`readerFor` stays unexported, C12/C13 read-back path unchanged); 0002's
`Accessor.Role` is read, never redefined; 0005 sees only the already-recorded
additive verb override; 0006 and 0010 are lint-side and never see the gate.
0016 is not contradicted: it is Draft, and if it lands, first-match and
unique-match coincide, so the amended C1 stays true under both worlds.
