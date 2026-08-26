package table

// RDR 0008 Load-Bearing Decisions / Identity — the near-miss advisory.
//
// PHASE 1 DECLARATION ONLY. LoadWithAdvisories below returns no advisories.
// The conformance suite in reserved_key_0008_test.go is red against it by
// design, and Phase 2 fills in the near-miss scan.

// AdvisoryNearMiss is the stable rule identifier a near-miss advisory carries.
// It is a comparable token asserted byte-for-byte, never prose.
const AdvisoryNearMiss = "reserved-tag-key/near-miss"

// Advisory is one non-blocking load-time advisory. It travels a channel
// distinct from the validation-failure list: it carries NO category
// discriminator, does not participate in category dispatch, and never alters
// the load/lint verdict (RDR 0008, Identity).
type Advisory struct {
	// Authored is the declaration key exactly as the author spelled it.
	Authored string
	// Reserved is the reserved spelling the authored one nearly matches —
	// for the near-miss rule, the literal `recognized`.
	Reserved string
	// Rule is the stable rule identifier, e.g. AdvisoryNearMiss.
	Rule string
}

// LoadWithAdvisories is Load plus the non-blocking advisory channel.
//
// A declaration whose post-parse key is not `recognized` but becomes
// `recognized` under EITHER Unicode-simple case folding OR trimming of
// leading/trailing whitespace — or both together — raises a near-miss
// advisory naming both spellings. The trigger is disjunctive: `Recognized`
// is a near-miss on folding alone and `" recognized"` on trimming alone.
//
// The advisory is not a failure: the name is legal and RDR 0008 must not
// reject a tag it does not own. The returned error and model are exactly
// what Load returns for the same bytes.
func LoadWithAdvisories(src []byte, sourceID string) (*Model, []Advisory, error) {
	// TODO(rdr-0008 Phase 2): collect near-miss advisories over the
	// [tags.<tag>] declaration keys. Unimplemented: every advisory assertion
	// currently fails.
	m, err := Load(src, sourceID)
	return m, nil, err
}
