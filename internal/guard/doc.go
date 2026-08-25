// Package guard owns RDR 0003's guard predicate model: the closed typed
// operator vocabulary, the tag declaration model every guard atom is
// written against, the value-comparison evaluator that satisfies the
// kernel's `resolve.GuardEvaluator` seam, and the finite-domain lint that
// proves coverage and overlap over a scoped row group.
//
// RDR 0002 owns where a declaration is authored and how it is carried
// through normalization; this package owns what a declaration MEANS. RDR
// 0007 owns the kernel seam, key presence, existence-atom verdicts, and
// the combination of per-atom verdicts; this package's evaluator decides
// value semantics over a PRESENT VALUE only and never reads the tag view.
//
// The package performs no I/O and imports no third-party code: it takes an
// already-loaded `table.Model` and returns values.
//
// Phase 1 status: the RDR 0003 test suite in this directory is the spec's
// enforcement surface and is written against the API this package will
// export. No implementation has landed yet, so every symbol the tests
// reference is undefined and each undefined-symbol error names the REQ its
// test header quotes.
package guard
