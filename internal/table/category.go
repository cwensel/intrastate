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

	// RDR 0008 `0008:C3` — the three-field payload every `reserved_tag_key`
	// failure carries at the data level. The guidance travels here, in the
	// failure data, not in the renderer: a consumer that does not know the
	// category still reads all three.
	//
	// PHASE 1 DECLARATION ONLY. Nothing populates these yet; the RDR 0008
	// conformance suite is red against them by design and Phase 2 fills them
	// in at the two `fail(CatReservedTagKey, …)` sites in load.go.

	// Offending is the offending declaration name exactly as authored.
	Offending string
	// Remedy is the name the author must use. For a recognized-provenance
	// declaration under a wrong name it is the literal `recognized`; for an
	// owned or observed declaration named `recognized` it is the EMPTY
	// STRING, because the remedy there is "choose any other name" and a
	// renderer must not present the reserved key as the required name.
	Remedy string
	// Rule is the stable, comparable rule identifier naming which rule within
	// the category fired — RuleKernelOwned or RuleAuthorMustRename. It is for
	// golden assertions and remediation lookup, never for category dispatch.
	Rule string

	// Line is the 1-based source line the refusal is attributed to, or ZERO
	// where the loader cannot ground one. Zero is the common case and it
	// means UNKNOWN, never line 1: most categories are decided against
	// already-decoded structs, from which no position survives, and the
	// decoder itself reports a position only for its own syntax failures.
	//
	// A renderer chooses its own fallback for zero. Line is diagnostic, not
	// identity: RDR 0002's round-trip invariant compares a locator's
	// presence and its rule-identifying part but EXCLUDES the optional
	// line/column detail, precisely because that detail moves when
	// unrelated source text is edited. Adding a real line is therefore a
	// diagnostic improvement, not a contract change.
	Line int
}

// The two direction-specific rule identifiers `0008:C3` fixes. Both sit inside
// the one `reserved_tag_key` category; they are not alternatives to the
// category discriminator and a consumer must not choose between them.
const (
	// RuleKernelOwned: a recognized-provenance declaration must be renamed TO
	// the reserved key.
	RuleKernelOwned = "reserved-tag-key/kernel-owned"
	// RuleAuthorMustRename: an owned or observed declaration must be renamed
	// AWAY FROM the reserved key.
	RuleAuthorMustRename = "reserved-tag-key/author-must-rename"
)

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

// fail builds a categorized refusal carrying no source position.
func fail(cat Category, detail string) error {
	return &Failure{Category: cat, Detail: detail}
}

// atLine stamps a source line onto a refusal that already carries its
// category and detail, and returns it unchanged when the line is not
// grounded or the error is not a *Failure.
//
// It is a post-hoc stamp rather than a `fail` parameter on purpose: the
// position is recoverable at ONE seam — the caller that still holds the
// source bytes — while the refusals themselves are minted deep in checks
// that ran against decoded structs. Threading a position every `fail` site
// cannot supply would buy nothing but a wider signature.
func atLine(err error, line int) error {
	if line <= 0 {
		return err
	}
	var f *Failure
	if !errors.As(err, &f) {
		return err
	}
	f.Line = line
	return err
}
