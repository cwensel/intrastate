package resolve

// RDR 0008 `0008:C4` / `0008:C5` — the reserved-key producer precondition.
//
// PHASE 1 DECLARATION ONLY. The predicate below is UNIMPLEMENTED: it accepts
// every Input. The conformance suite in reserved_key_0008_test.go is red
// against it by design, and Phase 2 fills in the scan. The declaration lands
// separately so the repository's pre-commit gate (`go vet ./...`) stays green
// while the tests that drive the implementation are already in place.

// CheckInput reports whether in breaches RDR 0008's reserved-key producer
// obligation, returning a non-nil error on the first breach it finds.
//
// The reserved key `recognized` enters the assembled evaluation view only
// through Input.Recognized. A producer supplying it any other way is making a
// programmer mistake, which travels this Go error path rather than the modeled
// refusal path RDR 0001 reserves for value-level dispositions.
//
// Its read-domain is exactly three sequences — Input.Owned, Input.Observed,
// and each row's RequiresOwned reached through Input.Table.Rows — and no other
// Input field. The check is unconditional on Input.Recognized: a reserved-keyed
// owned or observed tag is a breach whether or not the resolve carries an
// outcome.
//
// Resolve applies this same predicate at entry, so the reserved name has
// exactly one enforcement point across every channel it can arrive through.
func CheckInput(in Input) error {
	// TODO(rdr-0008 Phase 2): scan Owned, Observed, and Table.Rows'
	// RequiresOwned for recognizedTagKey; return the first breach as a single
	// non-nil error. Unimplemented: every conformance assertion over this
	// function currently fails.
	_ = in
	return nil
}
