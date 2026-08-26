package resolve

import (
	"errors"
	"slices"
)

// RDR 0008 `0008:C4` / `0008:C5` — the reserved-key producer precondition.

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
	for _, tag := range in.Owned {
		if tag.Key == recognizedTagKey {
			return errReservedOwnedTag
		}
	}
	for _, tag := range in.Observed {
		if tag.Key == recognizedTagKey {
			return errReservedObservedTag
		}
	}
	for _, row := range in.Table.Rows {
		if slices.Contains(row.RequiresOwned, recognizedTagKey) {
			return errReservedRequiresOwned
		}
	}
	return nil
}

// The three breach channels, one sentinel each. They are plain errors: a
// producer breach is a programmer mistake, not table data, so it carries no
// load category and spells no RefusalKind (`0008:C4`, REQ-37/REQ-49).
var (
	errReservedOwnedTag = errors.New(
		"resolve: owned tag on the reserved key " + recognizedTagKey +
			"; the recognized outcome enters only through Input.Recognized")
	errReservedObservedTag = errors.New(
		"resolve: observed tag on the reserved key " + recognizedTagKey +
			"; the recognized outcome enters only through Input.Recognized")
	errReservedRequiresOwned = errors.New(
		"resolve: row RequiresOwned names the reserved key " + recognizedTagKey +
			"; the key is kernel-supplied and never owned state")
)
