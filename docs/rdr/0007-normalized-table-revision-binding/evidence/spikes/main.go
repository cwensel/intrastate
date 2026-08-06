package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sort"
	"unicode/utf8"
)

const domain = "intrastate.transition-table-revision\x00v1"

type sourceLocator struct {
	Path   string
	Line   uint64
	Column uint64
}

type predicate struct {
	Key   string
	Value string
}

type nextTag struct {
	Tag   string
	Value string
}

type write struct {
	Tag       string
	Operation string
	Value     string
}

type row struct {
	RuleID          string
	ExpansionSuffix string
	Kind            string
	Predicates      []predicate
	EscapeClasses   []string
	NextTags        []nextTag
	Writes          []write
	SourceLocator   sourceLocator
}

type table struct {
	ModelID             string
	Outcomes            []string
	Rows                []row
	SourceSchemaVersion uint64
}

func main() {
	base := baseline()
	basePreimage := canonicalPreimage(base)
	baseRevision := revision(basePreimage)
	fmt.Printf("golden preimage_hex=%s\n", hex.EncodeToString(basePreimage))
	fmt.Printf("golden preimage_len=%d revision=%s\n", len(basePreimage), baseRevision)

	equivalent := []struct {
		name  string
		value table
	}{
		{"all_semantic_sets_reordered", reordered(base)},
		{"map_iteration_order", withMapConstructedNextTags(base)},
		{"source_locator_reformatted", withLocator(base, sourceLocator{Path: "reformatted/model.toml", Line: 900, Column: 41})},
		{"source_schema_version_excluded", withSourceSchemaVersion(base, 99)},
		{"nil_empty_semantic_sets", emptyTable(false)},
	}
	for _, tc := range equivalent {
		want := basePreimage
		if tc.name == "nil_empty_semantic_sets" {
			want = canonicalPreimage(emptyTable(true))
		}
		got := canonicalPreimage(tc.value)
		fmt.Printf("equivalent %-30s same_preimage=%t same_revision=%t\n",
			tc.name, bytes.Equal(want, got), revision(want) == revision(got))
	}

	changes := []struct {
		name   string
		mutate func(table) table
	}{
		{"model_id", func(t table) table { t.ModelID = "kata"; return t }},
		{"outcome", func(t table) table { t.Outcomes[0] = "rejected"; return t }},
		{"rule_id", func(t table) table { t.Rows[0].RuleID = "advance-renamed"; return t }},
		{"expansion_suffix", func(t table) table { t.Rows[0].ExpansionSuffix = "b"; return t }},
		{"kind", func(t table) table { t.Rows[0].Kind = "escape"; return t }},
		{"predicate_key", func(t table) table { t.Rows[0].Predicates[0].Key = "observed.status.eq"; return t }},
		{"predicate_value", func(t table) table { t.Rows[0].Predicates[0].Value = "Final"; return t }},
		{"escape_classes", func(t table) table { t.Rows[1].EscapeClasses = t.Rows[1].EscapeClasses[1:]; return t }},
		{"next_tag_name", func(t table) table { t.Rows[0].NextTags[0].Tag = "phase"; return t }},
		{"next_tag_value", func(t table) table { t.Rows[0].NextTags[0].Value = "reconcile"; return t }},
		{"write_tag", func(t table) table { t.Rows[0].Writes[0].Tag = "review"; return t }},
		{"write_operation", func(t table) table { t.Rows[0].Writes[0].Operation = "clear"; return t }},
		{"write_value", func(t table) table { t.Rows[0].Writes[0].Value = "open"; return t }},
	}
	for _, tc := range changes {
		got := canonicalPreimage(tc.mutate(clone(base)))
		fmt.Printf("changed    %-30s different_preimage=%t different_revision=%t\n",
			tc.name, !bytes.Equal(basePreimage, got), baseRevision != revision(got))
	}

	clear := clone(base)
	clear.Rows[0].Writes = []write{{Tag: "note", Operation: "clear"}}
	setEmpty := clone(base)
	setEmpty.Rows[0].Writes = []write{{Tag: "note", Operation: "set", Value: ""}}
	fmt.Printf("boundary   clear_vs_set_empty             different_preimage=%t different_revision=%t\n",
		!bytes.Equal(canonicalPreimage(clear), canonicalPreimage(setEmpty)),
		revision(canonicalPreimage(clear)) != revision(canonicalPreimage(setEmpty)))

	left := clone(base)
	left.Outcomes = []string{"ab", "c"}
	right := clone(base)
	right.Outcomes = []string{"a", "bc"}
	fmt.Printf("boundary   length_framing                 same_unframed_payload=%t different_preimage=%t different_revision=%t\n",
		bytes.Equal(unframed(left.Outcomes), unframed(right.Outcomes)),
		!bytes.Equal(canonicalPreimage(left), canonicalPreimage(right)),
		revision(canonicalPreimage(left)) != revision(canonicalPreimage(right)))
}

func baseline() table {
	return table{
		ModelID:             "rdr",
		SourceSchemaVersion: 1,
		Outcomes:            []string{"accepted", "blocked"},
		Rows: []row{
			{
				RuleID:          "advance",
				ExpansionSuffix: "a",
				Kind:            "transition",
				Predicates: []predicate{
					{Key: "owned.profile.eq", Value: "mid"},
					{Key: "recognized.outcome.eq", Value: "accepted"},
				},
				NextTags: []nextTag{
					{Tag: "stage", Value: "prelock"},
					{Tag: "status", Value: "Draft"},
				},
				Writes: []write{
					{Tag: "note", Operation: "set", Value: ""},
					{Tag: "scope", Operation: "clear"},
				},
				SourceLocator: sourceLocator{Path: "model.toml", Line: 12, Column: 1},
			},
			{
				RuleID:          "escape",
				ExpansionSuffix: "a",
				Kind:            "escape",
				Predicates: []predicate{
					{Key: "owned.profile.eq", Value: "mid"},
					{Key: "recognized.outcome.eq", Value: "blocked"},
				},
				EscapeClasses: []string{"ambiguous_match", "no_match"},
				NextTags:      []nextTag{{Tag: "stage", Value: "reconcile"}},
				SourceLocator: sourceLocator{Path: "model.toml", Line: 28, Column: 1},
			},
		},
	}
}

func reordered(t table) table {
	t = clone(t)
	reverse(t.Outcomes)
	reverse(t.Rows)
	for i := range t.Rows {
		reverse(t.Rows[i].Predicates)
		reverse(t.Rows[i].EscapeClasses)
		reverse(t.Rows[i].NextTags)
		reverse(t.Rows[i].Writes)
	}
	return t
}

func withMapConstructedNextTags(t table) table {
	t = clone(t)
	tags := map[string]string{"status": "Draft", "stage": "prelock"}
	t.Rows[0].NextTags = t.Rows[0].NextTags[:0]
	for tag, value := range tags {
		t.Rows[0].NextTags = append(t.Rows[0].NextTags, nextTag{Tag: tag, Value: value})
	}
	return t
}

func withLocator(t table, locator sourceLocator) table {
	t = clone(t)
	for i := range t.Rows {
		t.Rows[i].SourceLocator = locator
	}
	return t
}

func withSourceSchemaVersion(t table, version uint64) table {
	t = clone(t)
	t.SourceSchemaVersion = version
	return t
}

func emptyTable(useNil bool) table {
	t := table{
		ModelID: "empty",
		Rows: []row{{
			RuleID: "empty-row",
			Kind:   "transition",
		}},
	}
	if !useNil {
		t.Outcomes = []string{}
		t.Rows[0].Predicates = []predicate{}
		t.Rows[0].EscapeClasses = []string{}
		t.Rows[0].NextTags = []nextTag{}
		t.Rows[0].Writes = []write{}
	}
	return t
}

func unframed(values []string) []byte {
	var result []byte
	for _, value := range values {
		result = append(result, value...)
	}
	return result
}

func canonicalPreimage(t table) []byte {
	outcomes := make([][]byte, 0, len(t.Outcomes))
	for _, outcome := range t.Outcomes {
		outcomes = append(outcomes, stringValue(outcome))
	}
	rows := make([][]byte, 0, len(t.Rows))
	for _, candidate := range t.Rows {
		rows = append(rows, encodeRow(candidate))
	}
	root := recordValue(map[string][]byte{
		"model_id": stringValue(t.ModelID),
		"outcomes": setValue(outcomes),
		"rows":     setValue(rows),
	})
	return append([]byte(domain), root...)
}

func encodeRow(candidate row) []byte {
	predicates := make([][]byte, 0, len(candidate.Predicates))
	for _, item := range candidate.Predicates {
		predicates = append(predicates, recordValue(map[string][]byte{
			"key":   stringValue(item.Key),
			"value": stringValue(item.Value),
		}))
	}
	escapeClasses := make([][]byte, 0, len(candidate.EscapeClasses))
	for _, class := range candidate.EscapeClasses {
		escapeClasses = append(escapeClasses, stringValue(class))
	}
	nextTags := make([][]byte, 0, len(candidate.NextTags))
	for _, item := range candidate.NextTags {
		nextTags = append(nextTags, recordValue(map[string][]byte{
			"tag":   stringValue(item.Tag),
			"value": stringValue(item.Value),
		}))
	}
	writes := make([][]byte, 0, len(candidate.Writes))
	for _, item := range candidate.Writes {
		fields := map[string][]byte{
			"operation": stringValue(item.Operation),
			"tag":       stringValue(item.Tag),
		}
		if item.Operation == "set" {
			fields["value"] = stringValue(item.Value)
		}
		writes = append(writes, recordValue(fields))
	}
	return recordValue(map[string][]byte{
		"rule_id":          stringValue(candidate.RuleID),
		"expansion_suffix": stringValue(candidate.ExpansionSuffix),
		"kind":             stringValue(candidate.Kind),
		"predicates":       setValue(predicates),
		"escape_classes":   setValue(escapeClasses),
		"next_tags":        setValue(nextTags),
		"writes":           setValue(writes),
	})
}

func stringValue(value string) []byte {
	if !utf8.ValidString(value) {
		panic("invalid UTF-8")
	}
	return frame(0x01, []byte(value))
}

func setValue(values [][]byte) []byte {
	sort.Slice(values, func(i, j int) bool { return bytes.Compare(values[i], values[j]) < 0 })
	payload := make([]byte, 8)
	binary.BigEndian.PutUint64(payload, uint64(len(values)))
	for _, value := range values {
		payload = append(payload, value...)
	}
	return frame(0x03, payload)
}

func recordValue(fields map[string][]byte) []byte {
	names := make([]string, 0, len(fields))
	for name := range fields {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool { return bytes.Compare([]byte(names[i]), []byte(names[j])) < 0 })
	payload := make([]byte, 8)
	binary.BigEndian.PutUint64(payload, uint64(len(names)))
	for _, name := range names {
		payload = append(payload, stringValue(name)...)
		payload = append(payload, fields[name]...)
	}
	return frame(0x04, payload)
}

func frame(tag byte, payload []byte) []byte {
	result := make([]byte, 9, 9+len(payload))
	result[0] = tag
	binary.BigEndian.PutUint64(result[1:], uint64(len(payload)))
	return append(result, payload...)
}

func revision(preimage []byte) string {
	digest := sha256.Sum256(preimage)
	return "sha256:" + hex.EncodeToString(digest[:])
}

func clone(t table) table {
	result := t
	result.Outcomes = append([]string(nil), t.Outcomes...)
	result.Rows = append([]row(nil), t.Rows...)
	for i := range result.Rows {
		result.Rows[i].Predicates = append([]predicate(nil), t.Rows[i].Predicates...)
		result.Rows[i].EscapeClasses = append([]string(nil), t.Rows[i].EscapeClasses...)
		result.Rows[i].NextTags = append([]nextTag(nil), t.Rows[i].NextTags...)
		result.Rows[i].Writes = append([]write(nil), t.Rows[i].Writes...)
	}
	return result
}

func reverse[T any](values []T) {
	for left, right := 0, len(values)-1; left < right; left, right = left+1, right-1 {
		values[left], values[right] = values[right], values[left]
	}
}
