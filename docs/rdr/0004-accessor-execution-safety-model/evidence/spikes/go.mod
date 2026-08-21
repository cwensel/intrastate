// Nested module so this Stage-4 spike stays out of the root module's
// build (`go build/vet/test ./...`) while remaining runnable as evidence:
// RDR 0004 cites main.go::<symbol> anchors across ~10 assumptions.
// Same pattern as RDR 0002's spike. Stdlib only — no require block.
module rdr0004spike

go 1.26.3
