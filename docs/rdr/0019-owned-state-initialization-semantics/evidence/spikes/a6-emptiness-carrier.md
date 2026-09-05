Model: claude-opus-5

# A6 — emptiness carrier, source verification (Stage 6 reconcile)

Verdict: PASS. Candidate (b), an exported `flowbind` cardinality probe, survives.

## Implementer count: exactly TWO

`internal/cli/flowbind/flowbind.go::Reader.Read`
`internal/cli/cmdbind/cmdbind.go::Reader.Read`

`internal/accessor/executor.go::Executor.Read` is NOT a third implementer:
its signature `Read(ctx, name string) ReadResult` does not satisfy
`ReadBinding`; it is the driver that type-asserts `def.Binding.(ReadBinding)`
at `::Executor.invokeRead`.

`internal/cli/flowbind/edit.go::EditWriter` declares `CapWrite` and has no
`Read` method — not a read-side carrier.

## Seam reports no cardinality

`internal/accessor/binding.go::ReadBinding`:
  Read(ctx, art, requested []string) (values []KeyValue, unreadable []string, err error)
One KeyValue per REQUESTED key plus an unreadable list. No total, no count,
no enumeration of unrequested keys on Binding, ReadBinding, GateBinding or
WriteBinding.

## store is package-private

`type store map[string]string` (lowercase); `load`/`save` unexported.
Exported surface of package flowbind: Registry, OwnedTags, EditWriter/
NewEditWriter, Reader, Writer, Gate and their seam methods. Nothing
returns a key count.

## The load-bearing finding: WHERE the seal short-circuit lives

`::Reader.Read` short-circuits on the seal — after `load(art.Path)` it tests
`if _, sealed := s[sealedKey]; sealed` and returns
`return nil, slices.Clone(requested), nil`, every requested key unreadable,
before per-key resolution.

But that test is in `Read`'s BODY, not in `::load`. `load` returns the full
map INCLUDING `sealedKey`, which `::Writer.Apply` writes as an ordinary
entry. So on a read-back-sealed artifact `len(load(path))` is 1 — the real
count C1's sealed arm needs — obtained without entering the seal branch.

This is what keeps S9 a no-op success at exit 0 instead of degrading to an
exit-3 refusal.

## Per-candidate

(a) new `accessor.ReadBinding` capability — VIABLE ONLY as a DISTINCT method
    calling `load` directly. Reusing `Read` inherits the short-circuit and
    fails the sealed-store requirement outright. Also forces `cmdbind.Reader`
    and the accessor test doubles to grow a method answering a question the
    carrier gate guarantees is never asked of them.
(b) exported `flowbind` cardinality probe — VIABLE, cleanest. New exported
    method on `flowbind.Reader` calling the existing unexported `load`,
    returning `len(s)`. Answers 1 on a sealed store, never enters the seal
    branch, zero change to `ReadBinding`, zero change to `cmdbind.Reader`.
    THE SURVIVING CANDIDATE.
(c) `read-state`-family surface — no `read-state`/`readState`/`ReadState`
    symbol exists under `internal/`. A hypothetical new verb that would sit
    on (a) or (b) underneath. Not the carrier to pick.

## 0004:C3 satisfied

C3 bans discovering authoritative artifacts from ambient process state.
Candidate (b) takes the artifact from the per-invocation `accessor.Artifact`
(`art.Path`), exactly as `::Reader.Read` and `::Writer.Apply` already do. The
`os.ReadFile` inside `load` is the binding's own file access — what a
file-backed binding IS — not a caller bypassing the seam.
