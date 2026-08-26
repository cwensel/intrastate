package table

import (
	"sort"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
)

// RDR 0008 Load-Bearing Decisions / Identity — the near-miss advisory.

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
	m, err := Load(src, sourceID)
	return m, nearMissAdvisories(src), err
}

// nearMissAdvisories scans the [tags.<tag>] declaration keys for near-misses
// on the reserved key. It reads the document's own key set rather than a
// loaded model, so the advisory channel stays independent of the load verdict
// it must not alter: a document that refuses still reports its near-misses,
// and a document whose only issue is a near-miss loads clean.
func nearMissAdvisories(src []byte) []Advisory {
	var doc struct {
		Tags map[string]any `toml:"tags"`
	}
	if err := toml.Unmarshal(src, &doc); err != nil {
		// Unparseable bytes carry no declaration keys to advise on; the
		// malformed_toml refusal is Load's to report.
		return nil
	}

	keys := make([]string, 0, len(doc.Tags))
	for key := range doc.Tags {
		if isNearMiss(key) {
			keys = append(keys, key)
		}
	}
	if len(keys) == 0 {
		return nil
	}
	// Map iteration is unordered; a stable list keeps the channel
	// reproducible for a golden consumer.
	sort.Strings(keys)

	advisories := make([]Advisory, 0, len(keys))
	for _, key := range keys {
		advisories = append(advisories, Advisory{
			Authored: key,
			Reserved: RecognizedTagKey,
			Rule:     AdvisoryNearMiss,
		})
	}
	return advisories
}

// isNearMiss reports whether key is not the reserved key but becomes it under
// EITHER simple case folding OR trimming of leading/trailing whitespace, or
// both together. The trigger is disjunctive on purpose: `Recognized` needs
// folding alone and `" recognized"` needs trimming alone, so a conjunctive
// reading would fire on neither.
func isNearMiss(key string) bool {
	if key == RecognizedTagKey {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(key), RecognizedTagKey)
}
