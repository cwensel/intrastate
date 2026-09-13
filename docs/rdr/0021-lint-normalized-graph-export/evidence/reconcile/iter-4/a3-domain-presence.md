Model: claude-sonnet-5

Refs read: main (5fe32cf) for internal/guard/declaration.go (identical on both refs,
diff confirmed no delta). internal/cli/graph_document.go, graph_document_0021_test.go,
and testdata/graph_0021_state_machine.json exist ONLY on branch worktree-rdr-0021
(143d919) — absent on main (`git show main:<path>` -> fatal: path does not exist).

Query log: git rev-parse main / worktree-rdr-0021; git show <ref>:<path> | diff;
grep -n on each shown file; git ls-tree -r worktree-rdr-0021 --name-only | grep _test.go.

Q1 — internal/guard/declaration.go (main+branch identical), domainSize, lines 112-150:
  case "enum": L116-18 `if len(d.Domain) == 0 { return 0, false }` L119 `return spread(len(d.Domain), d.SingleValued), true` -> finite only w/ non-empty Domain; authors members.
  case "bool": L123 `return spread(2, d.SingleValued), true` -- finite, no d.Domain set/read.
  case "int": L128 `return spread(width, d.SingleValued), true` (L136 ceiling arm too) -- finite, no d.Domain.
  case "set": L145 `return spread(len(d.Elements), false), true` -- finite, uses d.Elements not d.Domain.
  default (scalar): L149 `return 0, false`.
  Confirms table exactly: bool/int/set finite-true with no Domain populated; only non-empty-domain enum authors members.

Q2 — internal/cli/graph_document.go (branch only), L185:
  `if _, finite := guard.AssignmentCount(decl); finite && len(decl.Domain) > 0 {`
  Gate is `finite && len(decl.Domain) > 0` -- requires BOTH finite AND non-empty Domain.
  This implements the AUTHORED-MEMBERS reading, not bare FINITE (bool/int/set are finite but len(Domain)==0, so excluded).

Q3 — doc comment, graph_document.go L51-54 (branch, CURRENT text):
  "// graphTagDoc is one declared tag. `domain` is present exactly when
  // `guard.AssignmentCount` reports the declaration's domain finite, and
  // ABSENT — key omitted — otherwise (REQ-20); `omitempty` is what delivers
  // the absence, and the finite arm always carries at least one member."
  NOT corrected. Still states the false FINITE-only claim ("always carries at
  least one member" is false for bool/int/set per Q1). Contradicts the L185 gate
  it documents. Second gate-adjacent comment at L180-181 repeats the same stale
  "present exactly when...carries a finite domain" framing.

Q4 — internal/cli/graph_document_0021_test.go (branch only) L140-262,
  TestReq20And28_TagDomainIsPresentExactlyWhenFinite:
  Fixture declares tags.recognized (enum, no domain, L153-157), tags.status
  (enum, domain=["a","b"], L159-164), tags.free (scalar, L~167). Assertions:
  L229-232 `finite, ok := seen["status"]` ... L233 `if _, has := finite["domain"]; !has { t.Errorf(...) }` (asserts PRESENT for enum-with-members).
  L239-242 `opaque, ok := seen["free"]` ... L243 `if raw, has := opaque["domain"]; has { t.Errorf(...) }` (asserts ABSENT for scalar).
  `recognized` is decoded into `seen` but NEVER referenced again -- no domain assertion.
  No bool or bounded-int tag declared anywhere in this test's fixture.
  CONFIRMED gap: only enum-with-members and scalar arms asserted; bool, bounded
  int, and member-less enum are UNASSERTED.

  Fixture testdata/graph_0021_state_machine.json (branch only), tags array:
  {"name":"flag",...,"kind":"bool",...} -- NO domain key.
  {"name":"recognized",...,"kind":"enum",...} -- NO domain key (member-less).
  {"name":"status",...,"kind":"enum",...,"domain":["a","b"]} -- domain PRESENT.
  Fixture data is consistent with the authored-members reading, but nothing in
  the test suite asserts flag/recognized's absence specifically.
