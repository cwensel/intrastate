package table_test

// Helpers for RDR 0008's declaration half. Every helper drives the real
// table.Load / table.LoadWithAdvisories over real fixture bytes; nothing here
// stubs the loader.

import (
	"errors"
	"reflect"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// reservedTagKeyFailure loads a fixture expected to refuse in the
// `reserved_tag_key` category and returns the categorized failure carrying
// block 3's three-field payload. The oracle is the category and the payload
// fields, never the message text (0002's failing-control rule, REQ-83).
func reservedTagKeyFailure(t *testing.T, fixture string) *table.Failure {
	t.Helper()

	_, err := table.Load(readFixture(t, fixture), fixture)
	if err == nil {
		t.Fatalf("fixture %s loaded clean; want a reserved_tag_key refusal", fixture)
	}
	var f *table.Failure
	if !errors.As(err, &f) {
		t.Fatalf("fixture %s refused with a non-categorized error: %v", fixture, err)
	}
	if f.Category != table.CatReservedTagKey {
		t.Fatalf("fixture %s refused %q; want %q",
			fixture, f.Category, table.CatReservedTagKey)
	}
	return f
}

// loadAdvisories loads a fixture through the advisory-carrying entry point and
// returns the advisory list. The near-miss advisory is non-blocking, so a
// refusal here is a test failure rather than an expected outcome.
func loadAdvisories(t *testing.T, fixture string) []table.Advisory {
	t.Helper()

	m, advisories, err := table.LoadWithAdvisories(readFixture(t, fixture), fixture)
	if err != nil {
		t.Fatalf("fixture %s must load clean; refused: %v", fixture, err)
	}
	if m == nil {
		t.Fatalf("fixture %s loaded a nil model with no error", fixture)
	}
	return advisories
}

// advisoryFor finds the advisory naming an authored spelling.
func advisoryFor(advisories []table.Advisory, authored string) (table.Advisory, bool) {
	for _, a := range advisories {
		if a.Authored == authored {
			return a, true
		}
	}
	return table.Advisory{}, false
}

// carriesCategory reports whether an advisory carries a category discriminator
// in any of its fields. REQ-63 forbids one: the advisory is not a validation
// failure and must not participate in category dispatch. The probe is over the
// value's own fields, because REQ-64 leaves the carrier's field names to the
// implementer — what is pinned is that no category token rides along.
func carriesCategory(a table.Advisory) bool {
	v := reflect.ValueOf(a)
	for i := range v.NumField() {
		f := v.Field(i)
		if f.Kind() != reflect.String {
			// A typed category field would not be a plain string; a
			// table.Category-typed field is caught by the type check below.
			if f.Type() == reflect.TypeOf(table.Category("")) {
				return true
			}
			continue
		}
		if f.Type() == reflect.TypeOf(table.Category("")) {
			return true
		}
		for _, cat := range table.Categories() {
			if f.String() == string(cat) {
				return true
			}
		}
	}
	return false
}
