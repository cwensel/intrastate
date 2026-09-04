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

	// The three RDR 0024 emit-vocabulary categories (`0024:D-naming`). The
	// wire slug and the registered category are ONE decision, so each
	// constant here also takes an entry in Categories() below.
	CatMalformedEmitDeclaration Category = "malformed_emit_declaration"
	CatUnknownEmitKey           Category = "unknown_emit_key"
	CatEmitValueOutOfDomain     Category = "emit_value_out_of_domain"

	// The six RDR 0025 command-carrier categories (`0025:C5`), in clause
	// order. The wire STRINGS here are the contract; these identifiers are
	// not. A constant is not in the closed set until it is appended to
	// Categories() below — an unregistered constant refuses correctly at
	// the call site while staying invisible to every consumer that
	// enumerates the set, which is why the registration is asserted
	// separately from the refusal (`0025:C5`, S1).
	CatCommandAndPathConflict    Category = "command_and_path_conflict"
	CatCommandEmpty              Category = "command_empty"
	CatCommandUnknownPlaceholder Category = "command_unknown_placeholder"
	CatCommandShellInterpreter   Category = "command_shell_interpreter"
	CatCommandOutputShape        Category = "command_output_shape"
	CatCommandEnvConflict        Category = "command_env_conflict"
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
		CatMalformedEmitDeclaration,
		CatUnknownEmitKey,
		CatEmitValueOutOfDomain,

		// The six RDR 0025 command-carrier categories, appended at the
		// TAIL in C5 clause order so a consumer enumerating the list sees
		// additions only there (`0025:C5`).
		CatCommandAndPathConflict,
		CatCommandEmpty,
		CatCommandUnknownPlaceholder,
		CatCommandShellInterpreter,
		CatCommandOutputShape,
		CatCommandEnvConflict,
	}
}

// categoryDescriptions carries the reviewer-facing text a category ships
// on the `--help-all` surface. It is PER-CATEGORY OPT-IN, not total over
// Categories(): a category absent here carries no text and renders as its
// identifier alone.
//
// Opt-in rather than total because `Categories()` is append-only and its
// total size is not a contract at any point (`0025:REQ-79`), so requiring
// text for every member would couple this map to every future append —
// and `0027:C1` ships text for exactly one category (`0027:D2`).
//
// The text states what a reviewer may RELY on. For a check that is
// deliberately not a barrier, that means naming what it does NOT cover:
// the promise is the predicate line plus the forms admitted by name, and a
// description that claimed more than the predicate delivers would be worse
// than none — a reviewer who believes the class is closed stops reading
// argv (`0027:C1` promise:).
var categoryDescriptions = map[Category]string{
	CatCommandShellInterpreter: `Refused: a listed interpreter word followed, at ANY later argv
position, by one of that interpreter's own inline-code flags — under any
prefix (env and its options, nice, timeout, xargs, doas, and wrappers
nobody enumerated). Nothing before the interpreter word is read, so no
wrapper table exists and none is consulted.

The check reads argv WORDS only. It never splits a word on whitespace,
and never reads stdin, files, PATH, or the resolved binary. The
interpreter set is an OPEN deny-list, so an unlisted spelling (python3,
nodejs, busybox) is admitted.

Out of scope, BY NAME — admitted by lint, and an interpreter may still
run:

  a shell string carried in ONE word, such as env -S "sh -c …", or a
  single "sh -c …" element handed to a tool that re-splits it;

  an interpreter that reads its script from STDIN — sh -s, bare sh,
  sh -es, python -, node -. The channel is the scope: any listed
  interpreter taking its code on stdin rather than as a later argv word
  is admitted, however spelled.

sh script.sh is the sanctioned wrapper-file form and never a defect.`,
}

// CategoryDescription returns the reviewer-facing text a category ships on
// the `--help-all` surface, and whether the category carries any. The
// surface is per-category opt-in, so a category with no text reports
// false rather than an empty string — a caller can then render the
// identifier alone instead of a blank line (`0027:C1` promise:, REQ-37).
func CategoryDescription(c Category) (string, bool) {
	d, ok := categoryDescriptions[c]
	return d, ok
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
