package table

import "errors"

// Category is a stable data-level load-failure discriminator. The Testing
// Strategy asserts on the category, never on message text, so these
// identifiers are an API surface (`0002:C24`).
//
// Identifiers are snake_case renderings of the prose category names the
// record enumerates, anchored on `reserved_tag_key`, which RDR 0008 fixes.
type Category string

// The load category floor `0002:C24` enumerates. Cross-row findings —
// overlap, gap, dead row, read-before-write — are RDR 0006's lint
// categories and deliberately absent here.
const (
	CatMalformedTOML                      Category = "malformed_toml"
	CatUnknownSchemaField                 Category = "unknown_schema_field"
	CatMissingRecognizedOutcomeAlphabet   Category = "missing_recognized_outcome_alphabet"
	CatMalformedRecognizedOutcomeAlphabet Category = "malformed_recognized_outcome_alphabet"
	CatUnknownTag                         Category = "unknown_tag"
	CatUnknownContext                     Category = "unknown_context"
	CatCyclicContextInheritance           Category = "cyclic_context_inheritance"
	CatWriteToNonOwnedTag                 Category = "write_to_non_owned_tag"
	CatUnknownAccessor                    Category = "unknown_accessor"
	CatMalformedAccessorDeclaration       Category = "malformed_accessor_declaration"
	CatMalformedAccessorBinding           Category = "malformed_accessor_binding"
	CatMalformedTagDeclaration            Category = "malformed_tag_declaration"
	CatMalformedInitialDeclaration        Category = "malformed_initial_declaration"
	CatReservedTagValue                   Category = "reserved_tag_value"
	CatUnsupportedVersion                 Category = "unsupported_version"
	CatMalformedPredicateAtom             Category = "malformed_predicate_atom"
	CatMalformedEscapeDeclaration         Category = "malformed_escape_declaration"
	CatMalformedModelDeclaration          Category = "malformed_model_declaration"
	CatMalformedDumpDeclaration           Category = "malformed_dump_declaration"
	CatMalformedOutcomeBinding            Category = "malformed_outcome_binding"
	CatMalformedRuleShape                 Category = "malformed_rule_shape"
	CatMalformedRuleID                    Category = "malformed_rule_id"
	CatDuplicateRuleID                    Category = "duplicate_rule_id"
	CatDuplicateModelID                   Category = "duplicate_model_id"
	CatReservedTagKey                     Category = "reserved_tag_key"
)

// Categories returns the closed load-category set in declaration order.
func Categories() []Category {
	return []Category{
		CatMalformedTOML,
		CatUnknownSchemaField,
		CatMissingRecognizedOutcomeAlphabet,
		CatMalformedRecognizedOutcomeAlphabet,
		CatUnknownTag,
		CatUnknownContext,
		CatCyclicContextInheritance,
		CatWriteToNonOwnedTag,
		CatUnknownAccessor,
		CatMalformedAccessorDeclaration,
		CatMalformedAccessorBinding,
		CatMalformedTagDeclaration,
		CatMalformedInitialDeclaration,
		CatReservedTagValue,
		CatUnsupportedVersion,
		CatMalformedPredicateAtom,
		CatMalformedEscapeDeclaration,
		CatMalformedModelDeclaration,
		CatMalformedDumpDeclaration,
		CatMalformedOutcomeBinding,
		CatMalformedRuleShape,
		CatMalformedRuleID,
		CatDuplicateRuleID,
		CatDuplicateModelID,
		CatReservedTagKey,
	}
}

// Failure is a categorized load refusal. Load is fail-fast, so a document
// yields exactly one of these and never a list (`0002:C3`).
type Failure struct {
	Category Category
	Detail   string
}

func (f *Failure) Error() string {
	if f.Detail == "" {
		return string(f.Category)
	}
	return string(f.Category) + ": " + f.Detail
}

// CategoryOf reports the stable category a refusal carries.
func CategoryOf(err error) (Category, bool) {
	var f *Failure
	if errors.As(err, &f) {
		return f.Category, true
	}
	return "", false
}

// fail builds a categorized refusal.
func fail(cat Category, detail string) error {
	return &Failure{Category: cat, Detail: detail}
}
