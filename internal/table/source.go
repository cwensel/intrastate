package table

import (
	"bytes"

	toml "github.com/pelletier/go-toml/v2"
)

// sourceDoc is the closed source layout JDR 0001 §D7 fixes: root
// `outcomes`, root `terminal`, `[model]`, `[initial]`, `[tags.<tag>]`,
// `[read.<id>]`, `[write.<id>]`, `[gate.<id>]`, `[context.<id>]`,
// `[emit.<key>]` (`0024:C1`), `[[rule]]`, and `[dump]`. No other root key
// or table is admitted (`0002:C2`).
type sourceDoc struct {
	Outcomes *[]string                 `toml:"outcomes"`
	Terminal []string                  `toml:"terminal"`
	Model    *sourceModel              `toml:"model"`
	Initial  map[string]any            `toml:"initial"`
	Tags     map[string]sourceTagDecl  `toml:"tags"`
	Read     map[string]sourceAcc      `toml:"read"`
	Write    map[string]sourceAcc      `toml:"write"`
	Gate     map[string]sourceAcc      `toml:"gate"`
	Context  map[string]sourceContext  `toml:"context"`
	Emit     map[string]sourceEmitDecl `toml:"emit"`
	Rule     []sourceRule              `toml:"rule"`
	Dump     *sourceDump               `toml:"dump"`
}

// sourceEmitDecl is one `[emit.<key>]` declaration (`0024:C1`). It carries
// exactly two keys — `kind` and `domain` — and no tag-declaration key: an
// emit declaration is NOT a tag declaration, so `provenance`, `min`, `max`,
// `elements`, `single_valued` and `required` are refused by strict decoding
// rather than by a hand-written arm.
type sourceEmitDecl struct {
	Kind string `toml:"kind"`
	// Domain is declared `any` because it is TWO-SHAPED: a flat member
	// array (`domain = [...]`) or a sub-table of disposition tokens
	// (`[emit.<key>.domain]`) are two spellings of ONE key, and no concrete
	// Go type accepts both. It mirrors sourceModel.Metadata.
	//
	// The carve-out is this FIELD TYPE, not a decoder option:
	// `pelletier/go-toml/v2` has no per-subtree strictness switch and
	// DisallowUnknownFields stays set for the whole document. What the
	// `any` costs is that malformed shapes INSIDE the domain — a
	// non-array, a non-string member, nesting below the disposition level
	// — reach the loader instead of the decoder, which is why C1 carries
	// explicit arms for them.
	Domain any `toml:"domain"`
}

type sourceModel struct {
	ID      *string `toml:"id"`
	Version *int    `toml:"version"`
	// Class is the model class (`0010:C1`). It is a POINTER so absence is
	// distinguishable from `class = ""`: an absent key reads as
	// `state-machine`, while the empty string is a value outside the
	// admitted set and refuses `malformed model declaration`.
	Class *string `toml:"class"`
	// Description is an admitted optional human annotation (`0002:C2`).
	Description string `toml:"description"`
	// Metadata is the one sanctioned extension namespace. It decodes as a
	// free-form map so strict decoding descends no further: strictness
	// applies everywhere else, and no key inside is interpreted
	// (`0002:C2`).
	Metadata map[string]any `toml:"metadata"`
}

type sourceTagDecl struct {
	Provenance   string   `toml:"provenance"`
	Kind         string   `toml:"kind"`
	Domain       []string `toml:"domain"`
	Min          *int     `toml:"min"`
	Max          *int     `toml:"max"`
	Elements     []string `toml:"elements"`
	SingleValued *bool    `toml:"single_valued"`
	Required     bool     `toml:"required"`
}

// sourceAcc is one accessor entry. Every field is a pointer or a slice so
// ABSENCE is distinguishable from an empty or false value: `role = ""` and
// an absent `role` are the same category, but an absent `read_back` and
// `read_back = false` must both refuse while `read_back = true` passes
// (`0002:C2`).
type sourceAcc struct {
	Role     *string   `toml:"role"`
	Path     *string   `toml:"path"`
	Keys     *[]string `toml:"keys"`
	Timeout  *string   `toml:"timeout"`
	ReadBack *bool     `toml:"read_back"`
}

type sourceContext struct {
	Inherits string                    `toml:"inherits"`
	Match    map[string]map[string]any `toml:"match"`
}

// sourceRule is one `[[rule]]`. The write block is a pointer so its
// PRESENCE is keyed rather than its length: an escape rule carrying
// `write = []` is refused just as one carrying a populated block is
// (`0002:C4`, deviations.md D2). The same holds for `clear`, `gate`, and
// `escape`.
type sourceRule struct {
	ID     *string   `toml:"id"`
	Use    []string  `toml:"use"`
	Source string    `toml:"source"`
	Clear  *[]string `toml:"clear"`
	Gate   *[]string `toml:"gate"`
	Escape *[]string `toml:"escape"`

	Match map[string]map[string]any `toml:"match"`
	Guard *sourceGuard              `toml:"guard"`
	Write *map[string]any           `toml:"write"`
	// Emit is the `[rule.emit]` answer block (`0010:C3`): a FLAT table
	// whose values MUST be strings.
	//
	// The map's value type is what enforces that. A non-string value — a
	// nested `[rule.emit.sub]` included — is a decoder TYPE error, which
	// decodeStrict maps to `malformed TOML`; no hand-written type check is
	// added, and duplicate keys are refused by the TOML decoder before any
	// of this code runs, which is why normalization carries no dedup pass.
	//
	// Unlike `write`/`clear`/`gate` this is NOT a pointer: `emit` keys on
	// LENGTH rather than on key presence, so an absent block and a
	// present-but-empty one are the same empty sequence.
	Emit map[string]string `toml:"emit"`
}

type sourceGuard struct {
	All    map[string]map[string]any `toml:"all"`
	Unless map[string]map[string]any `toml:"unless"`
}

// sourceDump carries exactly one key, `order` (`0002:C19`).
type sourceDump struct {
	Order []string `toml:"order"`
}

// versionProbe is the permissive first pass. The version check MUST run
// before strict field validation, not merely before normalization, so
// loading proceeds in two passes: read `[model]` permissively enough to
// obtain `version`, refuse on any value but 1, and only then decode
// strictly (`0002:C3`).
type versionProbe struct {
	Model *struct {
		ID      *string `toml:"id"`
		Version *int    `toml:"version"`
	} `toml:"model"`
}

// decodeStrict decodes src into dst, rejecting any key the layout does not
// map. Strict decoding is an obligation on this format, not a property of
// a library: an unknown schema field is a stable refusal, never a silent
// no-op (`0002:C3`).
func decodeStrict(src []byte, dst any) error {
	dec := toml.NewDecoder(bytes.NewReader(src))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		var strict *toml.StrictMissingError
		if asStrictError(err, &strict) {
			return fail(CatUnknownSchemaField, strict.Error())
		}
		return fail(CatMalformedTOML, err.Error())
	}
	return nil
}

func asStrictError(err error, target **toml.StrictMissingError) bool {
	e, ok := err.(*toml.StrictMissingError)
	if ok {
		*target = e
	}
	return ok
}
