package table

import "slices"

// CloneRowSetKeysForTest detaches the unexported slice that external tests
// cannot clone directly.
func CloneRowSetKeysForTest(in []Row) []Row {
	out := slices.Clone(in)
	for i := range out {
		out[i].setKeys = slices.Clone(out[i].setKeys)
	}
	return out
}
