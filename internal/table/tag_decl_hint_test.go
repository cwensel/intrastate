package table

// A kind/field mismatch must name the field the kind DOES take,
// and a tag-declaration refusal must carry the line its `[tags.<key>]`
// header sits on.
//
// The three mismatch arms in `tagDecl` used to say only what the kind does
// NOT admit ("kind set admits no domain"). A first-time author reads that as
// "a set cannot be finite" and goes looking for the wrong remedy. The three
// arms share one shape, so the hint is mirrored across all three: an
// `elements`-only hint on the domain arm would be wrong advice for `int` and
// `scalar`.
//
// These tests exercise `tagDecl` directly rather than through `Load`. The
// arms are per-declaration and order-independent, so a document fixture
// would only add the loader's step order between the assertion and the
// message under test.

import (
	"errors"
	"strings"
	"testing"
)

// declErrorDetail runs one declaration through `tagDecl` and returns the
// refusal's `Detail`.
func declErrorDetail(t *testing.T, key string, src sourceTagDecl) string {
	t.Helper()
	_, err := tagDecl(key, src)
	if err == nil {
		t.Fatalf("declaration %q loaded; want a malformed_tag_declaration refusal", key)
	}
	var f *Failure
	if !errors.As(err, &f) {
		t.Fatalf("refusal is not a *Failure: %v", err)
	}
	if f.Category != CatMalformedTagDeclaration {
		t.Fatalf("category = %q; want %q", f.Category, CatMalformedTagDeclaration)
	}
	if !strings.HasPrefix(f.Detail, "tag "+key+": ") {
		t.Fatalf("detail %q drops the `tag <key>: ` prefix", f.Detail)
	}
	return f.Detail
}

// TestKindFieldMismatchNamesTheFieldTheKindTakes pins that each of the three
// kind/field mismatch arms names the field its kind DOES declare, not only
// the one it refuses.
//
// FAILURE MODE GUARDED: a message that reports only the absence. "kind set
// admits no domain" is true and useless — it withholds the fact the code
// comment above the check already states, that a set declares `elements`.
func TestKindFieldMismatchNamesTheFieldTheKindTakes(t *testing.T) {
	cases := []struct {
		name string
		key  string
		src  sourceTagDecl
		// want is the token naming the field the offending kind DOES take.
		want string
	}{
		{
			name: "a set carrying a domain is pointed at elements",
			key:  "flavors",
			src: sourceTagDecl{
				Provenance: "owned",
				Kind:       "set",
				Domain:     []string{"a", "b"},
			},
			want: "elements",
		},
		{
			name: "an int carrying elements is pointed at min/max",
			key:  "count",
			src: sourceTagDecl{
				Provenance: "owned",
				Kind:       "int",
				Elements:   []string{"a"},
			},
			want: "min",
		},
		{
			name: "an enum carrying bounds is pointed at domain",
			key:  "state",
			src: sourceTagDecl{
				Provenance: "owned",
				Kind:       "enum",
				Min:        ptrTo(1),
			},
			want: "domain",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			detail := declErrorDetail(t, c.key, c.src)
			if !strings.Contains(detail, c.want) {
				t.Fatalf("detail %q never names %q — the refusal reports "+
					"only what the kind refuses and withholds what it takes",
					detail, c.want)
			}
		})
	}
}

// TestKindFieldMismatchHintsAreSymmetric pins that the three arms carry the
// SAME hint. The arms are one shape, and a hint naming only the field of the
// kind that happened to be authored would be wrong advice for the other
// kinds sharing the arm: `admits no domain` fires for `set`, `int`, and
// `scalar` alike, so the hint must enumerate the mapping rather than answer
// for one kind.
//
// FAILURE MODE GUARDED: fixing the `set` case in isolation, leaving the
// sibling arms reporting only an absence and the trio asymmetric.
func TestKindFieldMismatchHintsAreSymmetric(t *testing.T) {
	// One declaration per arm, each authored with a DIFFERENT offending kind
	// so a hint hardcoded to one kind cannot pass all three.
	details := []string{
		declErrorDetail(t, "a", sourceTagDecl{
			Provenance: "owned", Kind: "set", Domain: []string{"x"}}),
		declErrorDetail(t, "b", sourceTagDecl{
			Provenance: "owned", Kind: "scalar", Min: ptrTo(0)}),
		declErrorDetail(t, "c", sourceTagDecl{
			Provenance: "owned", Kind: "enum", Elements: []string{"x"}}),
	}

	// Every arm names every kind's own field, so an author who reached the
	// wrong arm still reads the whole mapping.
	for _, want := range []string{"domain", "elements", "min"} {
		for i, detail := range details {
			if !strings.Contains(detail, want) {
				t.Fatalf("arm %d detail %q never names %q; the three "+
					"kind/field arms must carry the same mapping hint",
					i, detail, want)
			}
		}
	}
}

func ptrTo[T any](v T) *T { return &v }
