# Charted — net-new scope out of cli/0021 critique

Charted: `graphlint.NodeCeiling`'s doc comment attributes the node-ceiling
refusal to `graph-product-too-large` (`internal/graphlint/taxonomy.go:150-152`),
but that code is the guard-product bound (`taxonomy.go:34`, `guard/lint.go:22`);
the node ceiling is a distinct constant (`nodeCeiling = 4096`). Out of scope —
a stale comment in shipped source, not an RDR defect: 0021:C4 already states the
two bounds are distinct and mints its own `graph-export-too-large`. Suggested
successor: a `type:bug` kata against `internal/graphlint`, alongside kata `yybx`
(which already ships the stale `reach.go` comment fix). Likely the source of
critique pass A's C-1/C-9 (both dismissed as unfounded).
