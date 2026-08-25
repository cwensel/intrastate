// Package table owns the sparse transition-model source schema and its
// normalization into deterministic candidate rows.
//
// RDR 0002 locks this package's contract. The source is sparse TOML data
// authored for review; normalization expands it into the candidate-row
// value that lint, the resolver, diagnostics, and the expanded-table dump
// all read. The package performs no file I/O and no path resolution — the
// load entry takes already-read bytes plus a source id — and it never
// imports internal/cli: its validation failures carry stable data-level
// categories that the CLI maps, not CLI envelope shapes.
package table
