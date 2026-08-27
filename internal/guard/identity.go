package guard

import (
	"slices"
	"strconv"
	"strings"

	"github.com/cwensel/intrastate/internal/table"
)

// SemanticKind is one of this RDR's predicate semantic kinds — the
// rejection vocabulary this RDR owns. RDR 0006 mints no lint finding code
// for any of them; they are carried by RDR 0002's load categories.
type SemanticKind string

// The predicate semantic kinds this RDR owns.
const (
	SemanticKindUnknownOperator         SemanticKind = "unknown-operator"
	SemanticKindOperatorKindMismatch    SemanticKind = "operator-kind-mismatch"
	SemanticKindLiteralParseFailure     SemanticKind = "literal-parse-failure"
	SemanticKindLiteralOutsideDomain    SemanticKind = "literal-outside-declared-domain"
	SemanticKindDeclarationDisagreement SemanticKind = "declaration-kind-disagreement"
	SemanticKindUnknownTag              SemanticKind = "unknown-tag"
)

var semanticKinds = []SemanticKind{
	SemanticKindUnknownOperator,
	SemanticKindOperatorKindMismatch,
	SemanticKindLiteralParseFailure,
	SemanticKindLiteralOutsideDomain,
	SemanticKindDeclarationDisagreement,
	SemanticKindUnknownTag,
}

// SemanticKinds returns the closed predicate-semantic-kind set.
func SemanticKinds() []SemanticKind { return slices.Clone(semanticKinds) }

// LoadCategoryFor maps a predicate semantic kind onto RDR 0002's load
// category that carries it.
//
// The two the RDR names explicitly are the stage boundary: a malformed
// DECLARATION (a domain disagreeing with its kind) is rejected by the
// declaration loader before normalization completes, so it is a `malformed
// tag declaration`; a literal outside its tag's declared domain is rejected
// by guard parsing AFTER the declarations have loaded and before rows are
// yielded, so it is a `malformed predicate atom`. Neither is the other's,
// which is how a consumer tells which stage rejected.
func LoadCategoryFor(k SemanticKind) table.Category {
	if k == SemanticKindDeclarationDisagreement {
		return table.CatMalformedTagDeclaration
	}
	if k == SemanticKindUnknownTag {
		return table.CatUnknownTag
	}
	return table.CatMalformedPredicateAtom
}

// --- identity ------------------------------------------------------------

// Identity is one guard atom's identity tuple: the enclosing row's source
// identity joined to the atom's own four fields.
//
// Source identity is NOT an atom field — it is carried by the enclosing
// normalized row — which is why the tuple is the row's identity joined to
// the atom's four rather than an identity stored on the atom.
type Identity string

// String renders the identity tuple.
func (i Identity) String() string { return string(i) }

// AtomIdentity builds the six-field identity tuple `(RuleID,
// SourceLocator, key, block, operator, literal)`. An atom is never
// identified by its index within `all` or `unless`, so reordering a row's
// atoms leaves every identity unchanged, and the tuple is TOTAL: two atoms
// over one key in one block differing only in literal are two identities.
//
// Each field is length-prefixed so the encoding is injective for any field
// content — any delimiter is authorable inside a tag value or a rule id.
func AtomIdentity(row table.Row, atom table.Atom) Identity {
	var b strings.Builder
	writeField(&b, row.RuleID)
	writeField(&b, row.SourceLocator)
	writeField(&b, atom.Key)
	writeField(&b, string(atom.Block))
	writeField(&b, atom.Operator)
	for _, m := range atom.Literal {
		writeField(&b, m)
	}
	return Identity(b.String())
}

// SemanticKey renders an atom's semantic equality key: `(tag, operator,
// literal)`, with neither block nor source identity.
//
// Which equality each operation uses is fixed: diagnostics, deduplication,
// and any "same atom" claim use the identity tuple; only domain computation
// — deciding what subset of the product an atom denotes — uses this.
func SemanticKey(atom table.Atom) string {
	var b strings.Builder
	writeField(&b, atom.Key)
	writeField(&b, atom.Operator)
	for _, m := range atom.Literal {
		writeField(&b, m)
	}
	return b.String()
}

func writeField(b *strings.Builder, s string) {
	b.WriteString(strconv.Itoa(len(s)))
	b.WriteString(":")
	b.WriteString(s)
	b.WriteString("|")
}
