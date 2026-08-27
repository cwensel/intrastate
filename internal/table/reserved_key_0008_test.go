package table_test

// RDR 0008 — ownership of the recognized-outcome tag key name, declaration
// half (blocks 2 and 3, plus the near-miss advisory).
//
// Two halves with different starting colours. RDR 0002 (BUILD-ORDER run 2)
// already carries `0008:C2`'s two naming rules under `CatReservedTagKey`, so
// those tests PIN shipped behavior against regression. `0008:C3`'s three-field
// failure payload and the Load-Bearing Decisions' near-miss advisory are
// net-new surface this RDR owes.
//
// Nothing here mocks the unit under test: every assertion hands real bytes to
// the real table.Load. The only stub is scenario 7's synthetic category
// consumer, which is a downstream collaborator, not the loader.

import (
	"errors"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// The reserved key, written out rather than read from `table.RecognizedTagKey`.
// Comparing the loader's behavior to the loader's own constant would verify
// the package against itself (premortem P-5, REQ-6).
const reservedKey = "recognized"

// The two direction-specific rule identifiers and the advisory's, quoted from
// `0008:C3` and the Identity decision. REQ-27 requires they be asserted
// byte-for-byte, so they are literals here and never derived.
const (
	ruleKernelOwned      = "reserved-tag-key/kernel-owned"
	ruleAuthorMustRename = "reserved-tag-key/author-must-rename"
	ruleNearMiss         = "reserved-tag-key/near-miss"
)

// --- Block 2: declaration naming (landed in RDR 0002) --------------------

// REQ-10: `0008:C2` "A tag declaration with provenance `recognized` MUST be
// named `recognized`."
// REQ-18: `0008:C2` "Any violation is a data-level validation failure in the
// `reserved_tag_key` category, reported at table load/lint before any
// resolution — never a kernel refusal."
// REQ-25: "a reserved-name owned/observed declaration, fails table load/lint
// in the `reserved_tag_key` category … before any resolution runs."
// HAPPY PATH (regression pin — landed in RDR 0002)
func TestReq10And18_RecognizedProvenanceUnderAnotherNameFailsReservedTagKey(t *testing.T) {
	if got := loadCategory(t, "neg/neg-recognized-misnamed.toml"); got != table.CatReservedTagKey {
		t.Errorf("category = %q; want %q", got, table.CatReservedTagKey)
	}
}

// REQ-11: `0008:C2` "A tag declaration with provenance `owned` or `observed`
// MUST NOT be named `recognized`."
// HAPPY PATH (regression pin — landed in RDR 0002)
func TestReq11_OwnedDeclarationNamedRecognizedFailsReservedTagKey(t *testing.T) {
	if got := loadCategory(t, "neg/neg-recognized-owned.toml"); got != table.CatReservedTagKey {
		t.Errorf("category = %q; want %q", got, table.CatReservedTagKey)
	}
}

// REQ-11: `0008:C2` "A tag declaration with provenance `owned` or `observed`
// MUST NOT be named `recognized`."
// DOMAIN EDGE (regression pin — landed in RDR 0002)
//
// The observed half of REQ-11's "owned or observed". The loader's guard tests
// provenance for INEQUALITY against `recognized`, so owned and observed are
// treated identically today; nothing but this fixture stops a later narrowing
// to `== ProvenanceOwned` from passing the whole suite.
func TestReq11_ObservedDeclarationNamedRecognizedFailsReservedTagKey(t *testing.T) {
	if got := loadCategory(t, "neg/neg-recognized-observed.toml"); got != table.CatReservedTagKey {
		t.Errorf("category = %q; want %q", got, table.CatReservedTagKey)
	}
}

// REQ-12: `0008:C2` "because `[tags.<tag>]` is keyed by tag name, one model
// admits at most one declaration named `recognized`, so at most one
// recognized-provenance declaration survives validation — cardinality is a
// consequence of the naming rule, not a separate check."
// REQ-82: TS-5 "the load is refused with **exactly one** failure, in the
// `reserved_tag_key` category, naming one of the two declarations as authored;
// which of the two is reported is unspecified and MUST NOT be asserted"
// REQ-83: TS-5 "The assertion is by category, not message text."
// DOMAIN EDGE (regression pin — landed in RDR 0002)
//
// A NEGATIVE requirement rides here too: no separate cardinality check may be
// added. Asserted by the failure being the naming category and being SINGULAR
// — a distinct cardinality rule would show up as a different category or as a
// list.
func TestReq12And82And83_TwoRecognizedDeclarationsYieldExactlyOneNamingFailure(t *testing.T) {
	_, err := table.Load(readFixture(t, "neg/neg-two-recognized-decls.toml"),
		"neg-two-recognized-decls.toml")
	if err == nil {
		t.Fatal("two recognized-provenance declarations loaded clean")
	}

	cat, ok := table.CategoryOf(err)
	if !ok {
		t.Fatalf("uncategorized refusal: %v", err)
	}
	if cat != table.CatReservedTagKey {
		t.Errorf("category = %q; want %q — cardinality is a consequence of the "+
			"naming rule, not a separate check", cat, table.CatReservedTagKey)
	}
	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		t.Errorf("the load reported %d failures; exactly one is required",
			len(joined.Unwrap()))
	}
	// Which of the two is named is UNSPECIFIED and deliberately unasserted.
}

// REQ-13: `0008:C2` "two declarations spelled `recognized` are a TOML
// duplicate-key error in the `malformed TOML` category before this rule is
// reached."
// REQ-84: TS-5 "The same-name variant is out of scope for this category:
// `[tags.recognized]` twice is a TOML duplicate-key error in `malformed TOML`,
// asserted as such."
// BOUNDARY (regression pin — landed in RDR 0002)
func TestReq13And84_SameNameTwiceIsMalformedTOMLNotReservedTagKey(t *testing.T) {
	got := loadCategory(t, "neg/neg-recognized-declared-twice.toml")
	if got == table.CatReservedTagKey {
		t.Fatalf("the same-name variant was folded into %q; it is a TOML "+
			"duplicate-key error and out of this category's scope", got)
	}
	if got != table.CatMalformedTOML {
		t.Errorf("category = %q; want %q", got, table.CatMalformedTOML)
	}
}

// REQ-15: `0008:C2` "The bound is an upper one only: this RDR's naming rule
// imposes **no lower bound** of its own."
// REQ-16: `0008:C2` "0002 decides that the declaration exists, this RDR
// decides what it is named."
// DOMAIN EDGE (regression pin — landed in RDR 0002)
//
// A NEGATIVE requirement: the missing-declaration lower bound is RDR 0002's
// `malformed_model_declaration`, and this RDR must NOT re-home it into
// `reserved_tag_key`. The pin fails if a future change moves it.
func TestReq15And16_TheMissingDeclarationLowerBoundStaysOutsideReservedTagKey(t *testing.T) {
	got := loadCategory(t, "neg/neg-no-recognized-decl.toml")
	if got == table.CatReservedTagKey {
		t.Fatalf("the absent-declaration lower bound was re-homed into %q; "+
			"REQ-15 says this RDR's naming rule imposes no lower bound", got)
	}
	if got != table.CatMalformedModelDeclaration {
		t.Errorf("category = %q; want %q — 0002's outcome-binding contract owns "+
			"the lower bound", got, table.CatMalformedModelDeclaration)
	}
}

// REQ-19: `0008:C2` "The reserved-key comparison is byte-exact on the
// **post-parse** key string: case-sensitive, no trimming, no folding (so
// `Recognized` is an ordinary, unreserved name, and `\" recognized\"` — a
// whitespace-bearing key — is likewise ordinary and unreserved)."
// REQ-56: "equality is byte-exact string match on the **post-parse** key
// string — no case folding, no trimming, no aliasing."
// REQ-57: "Near-spellings that parse to a different key string (`Recognized`,
// `RECOGNIZED`, a whitespace-bearing `\" recognized\"`) are by definition
// ordinary unreserved names."
// REQ-79: TS-4 "the byte-exact rule on the post-parse key string treats
// `Recognized`, `RECOGNIZED`, and `\" recognized\"` as ordinary unreserved
// names (no folding, no trimming)"
// REQ-100: Phase 3's "reserved-key normalization fixtures (case and whitespace
// variants stay unreserved…)"
// INPUT EDGE (regression pin — landed in RDR 0002)
func TestReq19And56And57And79_NearSpellingsAreOrdinaryUnreservedNames(t *testing.T) {
	unreserved := map[string]string{
		"Recognized (case-folded)":     "pos-near-miss-folded.toml",
		"RECOGNIZED (upper)":           "pos-near-miss-upper.toml",
		"\" recognized\" (whitespace)": "pos-near-miss-space.toml",
	}

	for name, fixture := range unreserved {
		t.Run(name, func(t *testing.T) {
			_, err := table.Load(readFixture(t, fixture), fixture)
			if err == nil {
				return
			}
			if cat, _ := table.CategoryOf(err); cat == table.CatReservedTagKey {
				t.Errorf("%s was treated as the reserved key; the comparison is "+
					"byte-exact with no folding and no trimming", name)
			} else {
				t.Errorf("%s failed load for an unrelated reason: %v", name, err)
			}
		})
	}
}

// REQ-20: `0008:C2` "TOML quoting is a surface artifact, not a name variant:
// `[tags.\"recognized\"]` and `[tags.recognized]` parse to the identical key
// string and are therefore both reserved."
// REQ-79 (second half): TS-4 "the quoted `\"recognized\"`, which parses to the
// identical key string as the bare form, **is** reserved."
// ADVERSARIAL (regression pin — landed in RDR 0002)
//
// The witness is the quoted form under a NON-recognized provenance: if quoting
// made a distinct name, this would be an ordinary owned declaration and would
// load. It must refuse `reserved_tag_key` exactly as the bare form does.
func TestReq20And79_QuotedFormParsesToTheSameKeyAndIsReserved(t *testing.T) {
	got := loadCategory(t, "neg/neg-recognized-quoted-owned.toml")
	if got != table.CatReservedTagKey {
		t.Errorf("category = %q; want %q — quoting is a surface artifact, not a "+
			"name variant", got, table.CatReservedTagKey)
	}
}

// REQ-21: `0008:C2` "The checked positions are the `[tags.<tag>]` declaration
// keys only."
// REQ-22: `0008:C2` "Predicate positions need no separate reserved-key check
// because RDR 0002's outcome-binding contract already decides every one of
// them"
// REQ-23: `0008:C2` "`[rule.write]` is covered by a *different* pre-existing
// rule and must not be folded into the sentence above … what rejects it is RDR
// 0002's `write to non-owned tag` category … This RDR adds no write-position
// check and depends on that category holding."
// DOMAIN EDGE (regression pin — landed in RDR 0002)
//
// NEGATIVE requirements. Positions other than the declaration key must NOT
// report `reserved_tag_key`: the guard position is 0002's, the accessor `keys`
// position is 0002's, and both keep their own categories. This test fails if a
// future change widens `reserved_tag_key` past the declaration keys.
func TestReq21And22And23_ReservedTagKeyIsScopedToDeclarationKeysOnly(t *testing.T) {
	otherPositions := map[string]string{
		"guard position":         "neg/neg-recognized-in-guard.toml",
		"accessor keys position": "neg/neg-recognized-in-keys.toml",
	}

	for name, fixture := range otherPositions {
		t.Run(name, func(t *testing.T) {
			got := loadCategory(t, fixture)
			if got == table.CatReservedTagKey {
				t.Errorf("the %s reported %q; the checked positions are the "+
					"[tags.<tag>] declaration keys only", name, got)
			}
		})
	}
}

// REQ-17: `0008:C2` "For a `resolve.Table` built by any producer other than
// 0002's loader no lower bound applies, and a row meant to match the
// recognized tag under an undeclared or innocent name refuses `no_match` at
// resolve time"
// REQ-14: `0008:C2` "The scope unit is one `[model]` (RDR 0002's schema unit);
// this RDR defines no cross-model or cross-file cardinality rule."
// DOMAIN EDGE (regression pin)
//
// The kernel-side half of REQ-17 lives in `internal/resolve`; here the
// table-side half is pinned: the loader's lower bound is a per-`[model]`
// obligation and the canonical fixtures each satisfy it on their own.
func TestReq14And17_TheScopeUnitIsOneModelWithNoCrossFileRule(t *testing.T) {
	for _, fixture := range []string{rdrFixture, kataFixture} {
		t.Run(fixture, func(t *testing.T) {
			m := mustLoad(t, fixture)
			if _, ok := m.Tags[reservedKey]; !ok {
				t.Errorf("%s declares no %q tag; the lower bound is per-[model]",
					fixture, reservedKey)
			}
		})
	}
}

// REQ-24: `0008:C2` "a reference points at the wrong declared tag. Closing it
// needs a rule over declaration *intent* rather than declaration *name*, which
// is the same out-of-scope heuristic the Failure Modes section charts"
// DOMAIN EDGE (regression pin)
//
// A NEGATIVE requirement: the stray-predicate case is explicitly NOT closed
// here. A model that declares the reserved key correctly and whose rules point
// at an ORDINARY declared tag under a plausible-but-wrong name must load
// clean — closing it would need an intent heuristic this RDR disclaims. This
// pin fails if a future change starts refusing well-named references.
func TestReq24_TheStrayPredicateCaseIsNotClosedByThisRDR(t *testing.T) {
	// The canonical fixture declares `finalized_at` (observed) alongside the
	// reserved key, and its rules reference it by name. Nothing in this RDR
	// may make such a reference a failure on intent grounds.
	m := mustLoad(t, rdrFixture)
	if _, ok := m.Tags["finalized_at"]; !ok {
		t.Fatalf("%s no longer declares the ordinary observed tag this pin "+
			"rests on", rdrFixture)
	}
	if _, ok := m.Tags[reservedKey]; !ok {
		t.Fatalf("%s no longer declares %q", rdrFixture, reservedKey)
	}
	// The near-miss fixtures are the sharper witness: a declaration whose name
	// all but announces the wrong intent still loads clean, because the rule
	// is over the NAME, never the intent.
	for _, fixture := range []string{
		"pos-near-miss-folded.toml",
		"pos-near-miss-upper.toml",
		"pos-near-miss-space.toml",
	} {
		if _, err := table.Load(readFixture(t, fixture), fixture); err != nil {
			t.Errorf("%s was refused: %v — closing the wrong-intent case needs a "+
				"rule over declaration intent, which this RDR does not add",
				fixture, err)
		}
	}
}

// REQ-32: `0008:C3` "The category's stable data-level value is the token
// `reserved_tag_key`, matching the snake_case discriminator grammar RDR 0001
// uses for `RefusalKind` values"
// BOUNDARY (regression pin — landed in RDR 0002)
func TestReq32_TheCategoryTokenIsExactlyReservedTagKey(t *testing.T) {
	if string(table.CatReservedTagKey) != "reserved_tag_key" {
		t.Errorf("CatReservedTagKey = %q; want %q",
			table.CatReservedTagKey, "reserved_tag_key")
	}
}

// REQ-106: Deviation D2's check — "no `[read.<id>]`/`[write.<id>]`/
// `[gate.<id>]` `keys` entry may name `recognized` (a writer key must be
// owned; `recognized` is kernel-supplied)."
// ADVERSARIAL (regression pin — landed in RDR 0002)
//
// D2 asks for the CHECK's existence, not its category. It is 0002's
// `malformed_accessor_binding`, and this pin fails if the check disappears.
func TestReq106_NoAccessorKeysEntryMayNameTheReservedKey(t *testing.T) {
	got := loadCategory(t, "neg/neg-recognized-in-keys.toml")
	if got == "" {
		t.Fatal("a capability keys entry naming the reserved key loaded clean")
	}
	if got != table.CatMalformedAccessorBinding {
		t.Errorf("category = %q; want %q — the check exists under 0002's "+
			"category, not reserved_tag_key", got, table.CatMalformedAccessorBinding)
	}
}

// --- Block 3: the three-field failure payload (NET-NEW) ------------------

// REQ-26: `0008:C3` "Every `reserved_tag_key` failure … MUST carry, at the
// data level, three distinct machine-readable fields: the offending name as
// authored, a **remedy name**, and a stable rule identifier."
// REQ-29: `0008:C3` "A recognized-provenance declaration under a wrong name
// must be renamed **to** the reserved key. Remedy name: the literal
// `recognized`. Rule identifier: `reserved-tag-key/kernel-owned`."
// REQ-30: `0008:C3` "An owned or observed declaration named `recognized` must
// be renamed **away from** the reserved key … the field carries the empty
// string, and the rule identifier `reserved-tag-key/author-must-rename` is what
// tells a consumer the remedy is \"choose any other name\" rather than \"use
// this one\". A renderer MUST NOT present the reserved key as the required
// name in this direction."
// REQ-92: TS-7's byte-for-byte payload for both directions.
// REQ-101: Phase 3's golden failure-data check, both directions.
// ADVERSARIAL (net-new)
//
// Asserted byte-for-byte. The empty remedy in the author-must-rename direction
// is load-bearing: it is what keeps a renderer from presenting `recognized` as
// the required name for a declaration that must move AWAY from it.
func TestReq26And29And30And92And101_BothDirectionsCarryTheThreeFieldPayload(t *testing.T) {
	cases := map[string]struct {
		fixture   string
		offending string
		remedy    string
		rule      string
	}{
		"recognized-provenance under a wrong name": {
			fixture:   "neg/neg-recognized-misnamed.toml",
			offending: "outcome",
			remedy:    reservedKey,
			rule:      ruleKernelOwned,
		},
		"owned declaration named recognized": {
			fixture:   "neg/neg-recognized-owned.toml",
			offending: reservedKey,
			remedy:    "",
			rule:      ruleAuthorMustRename,
		},
		"observed declaration named recognized": {
			fixture:   "neg/neg-recognized-observed.toml",
			offending: reservedKey,
			remedy:    "",
			rule:      ruleAuthorMustRename,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			f := reservedTagKeyFailure(t, tc.fixture)

			if f.Offending != tc.offending {
				t.Errorf("offending name = %q; want %q (as authored)",
					f.Offending, tc.offending)
			}
			if f.Remedy != tc.remedy {
				t.Errorf("remedy name = %q; want %q", f.Remedy, tc.remedy)
			}
			if f.Rule != tc.rule {
				t.Errorf("rule identifier = %q; want %q", f.Rule, tc.rule)
			}
		})
	}
}

// REQ-27: `0008:C3` "The rule identifier is a comparable token, not prose: a
// golden test asserts it byte-for-byte, and human-readable wording is the
// renderer's to choose."
// REQ-31: `0008:C3` "Both rule identifiers sit inside the one
// `reserved_tag_key` category"
// REQ-33: `0008:C3` "The two tokens are not alternatives and a consumer MUST
// NOT choose between them: `reserved_tag_key` is the **category**
// discriminator … The **rule identifier** … is carried inside the failure
// payload, identifying which rule within the category fired; it is for golden
// assertions and remediation lookup, never for category dispatch."
// BOUNDARY (net-new)
//
// One category, two rule identifiers, and the two axes are distinct: the
// category never varies with direction, and the rule identifier always does.
func TestReq27And31And33_OneCategoryCarriesTwoDistinctRuleIdentifiers(t *testing.T) {
	kernelOwned := reservedTagKeyFailure(t, "neg/neg-recognized-misnamed.toml")
	mustRename := reservedTagKeyFailure(t, "neg/neg-recognized-owned.toml")

	if kernelOwned.Category != table.CatReservedTagKey ||
		mustRename.Category != table.CatReservedTagKey {
		t.Fatalf("both directions must sit in %q; got %q and %q",
			table.CatReservedTagKey, kernelOwned.Category, mustRename.Category)
	}
	if kernelOwned.Rule == mustRename.Rule {
		t.Fatalf("both directions carry the rule identifier %q; the identifier "+
			"is what distinguishes which rule within the category fired",
			kernelOwned.Rule)
	}
	// The identifier is a token, not prose: byte-for-byte against the literals.
	if kernelOwned.Rule != ruleKernelOwned {
		t.Errorf("rule = %q; want %q", kernelOwned.Rule, ruleKernelOwned)
	}
	if mustRename.Rule != ruleAuthorMustRename {
		t.Errorf("rule = %q; want %q", mustRename.Rule, ruleAuthorMustRename)
	}
}

// REQ-28: `0008:C3` "A category consumer may map the failure, but the guidance
// travels in the failure data, not the renderer."
// REQ-91: TS-7 "**both directions** of the naming rule, each passed to a
// synthetic consumer that maps only the categories RDR 0002 enumerates today
// (i.e. treats this one as unknown and falls through to a generic branch).
// **Expected**: the failure **data** still carries all three fields, asserted
// byte-for-byte and independent of anything the consumer renders."
// REQ-93: TS-7 "The consumer is a test stub, not RDR 0005's exit-code map:
// this RDR asserts the payload contract, and 0005 owns whatever mapping it
// later adds."
// ADVERSARIAL (net-new)
//
// The consumer is a stub that deliberately does NOT know `reserved_tag_key`
// and renders a generic string for it. The three fields must survive that
// ignorance intact, which is what "the guidance travels in the failure data"
// means operationally.
func TestReq28And91And93_PayloadSurvivesAConsumerThatDoesNotKnowTheCategory(t *testing.T) {
	// A synthetic consumer over the categories that predate this one. It is a
	// downstream collaborator, not the loader.
	known := map[table.Category]string{
		table.CatMalformedTOML:             "the file is not valid TOML",
		table.CatUnknownTag:                "a rule names an undeclared tag",
		table.CatWriteToNonOwnedTag:        "a rule writes a non-owned tag",
		table.CatMalformedModelDeclaration: "the [model] table is malformed",
		table.CatMalformedAccessorBinding:  "an accessor binding is malformed",
	}
	render := func(cat table.Category) string {
		if msg, ok := known[cat]; ok {
			return msg
		}
		return "the model failed to load"
	}

	for _, fixture := range []string{
		"neg/neg-recognized-misnamed.toml",
		"neg/neg-recognized-owned.toml",
	} {
		t.Run(fixture, func(t *testing.T) {
			f := reservedTagKeyFailure(t, fixture)

			if got := render(f.Category); got != "the model failed to load" {
				t.Fatalf("the stub consumer knows %q; it must fall through to the "+
					"generic branch for this test to mean anything", f.Category)
			}
			if f.Offending == "" {
				t.Error("the offending name did not survive the generic branch")
			}
			if f.Rule == "" {
				t.Error("the rule identifier did not survive the generic branch")
			}
			// The remedy is legitimately empty in the author-must-rename
			// direction, so its presence is asserted through the rule
			// identifier's direction rather than by non-emptiness.
			switch f.Rule {
			case ruleKernelOwned:
				if f.Remedy != reservedKey {
					t.Errorf("remedy = %q; want %q", f.Remedy, reservedKey)
				}
			case ruleAuthorMustRename:
				if f.Remedy != "" {
					t.Errorf("remedy = %q; the author-must-rename direction "+
						"carries the empty string so no renderer presents the "+
						"reserved key as the required name", f.Remedy)
				}
			default:
				t.Errorf("unknown rule identifier %q", f.Rule)
			}
		})
	}
}

// REQ-34: `0008:C3` "The Go type, package, and field names carrying these
// values are RDR 0002's to choose — this RDR constrains the values and their
// distinctness, not the struct."
// REQ-35: `0008:C3` "Where the offending key is the reserved name, this payload
// requirement extends RDR 0002's pre-existing `unknown tag` failure; that
// extension is authored by this RDR and lands in 0002's implementation."
// DOMAIN EDGE (net-new)
//
// The constraint is DISTINCTNESS of the three values, not the field spelling.
// Per deviation D6 the payload binds `reserved_tag_key` failures only: 0002's
// implemented mandatory-declaration check fires before any `unknown tag`
// failure can carry the reserved key as its offending key, so the `unknown
// tag` arm is unreachable and is deliberately not asserted.
func TestReq34And35_TheThreeValuesAreDistinctFieldsNotOneConflatedString(t *testing.T) {
	f := reservedTagKeyFailure(t, "neg/neg-recognized-misnamed.toml")

	if f.Offending == f.Rule {
		t.Error("the offending name and the rule identifier are the same value; " +
			"block 3 requires three DISTINCT machine-readable fields")
	}
	if f.Remedy == f.Rule {
		t.Error("the remedy name and the rule identifier are the same value")
	}
	if f.Offending == f.Remedy {
		t.Errorf("offending and remedy are both %q; the wrong-named direction "+
			"must name a DIFFERENT remedy", f.Offending)
	}
}

// --- The near-miss advisory (NET-NEW) ------------------------------------

// REQ-58: "**Advisory warning — normative, not deferred.**"
// REQ-59: "a declaration whose post-parse key is not `recognized` but becomes
// `recognized` under **either** Unicode-simple case folding **or** trimming of
// leading/trailing whitespace, or both applied together, MUST raise a
// **non-blocking advisory** naming both spellings."
// REQ-60: "The trigger is deliberately disjunctive: `Recognized` (folding
// alone) and `\" recognized\"` (trimming alone) are each near-misses on their
// own, and scenario 4 requires the advisory on both, so a conjunctive reading
// would fire on neither."
// REQ-80: TS-4 "each of the three unreserved near-misses raises the
// **non-blocking advisory** while `[tags.result]` — an ordinary name that is
// not a near-miss — raises none, so the advisory is pinned as targeted rather
// than blanket."
// ADVERSARIAL (net-new)
//
// The disjunction is the load-bearing part: `Recognized` needs folding alone
// and `" recognized"` needs trimming alone, so a conjunctive implementation
// fires on neither. `result` is the control that keeps the advisory targeted.
func TestReq58And59And60And80_DisjunctiveNearMissTriggerAndTheTargetedControl(t *testing.T) {
	cases := map[string]struct {
		fixture  string
		authored string
		want     bool
	}{
		"folding alone (Recognized)": {
			fixture: "pos-near-miss-folded.toml", authored: "Recognized", want: true,
		},
		"folding alone (RECOGNIZED)": {
			fixture: "pos-near-miss-upper.toml", authored: "RECOGNIZED", want: true,
		},
		"trimming alone (\" recognized\")": {
			fixture: "pos-near-miss-space.toml", authored: " recognized", want: true,
		},
		"both together (\" Recognized\")": {
			fixture: "pos-near-miss-both.toml", authored: " Recognized", want: true,
		},
		"neither (result) — the targeted control": {
			fixture: "pos-not-near-miss.toml", authored: "result", want: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			advisories := loadAdvisories(t, tc.fixture)

			got, found := advisoryFor(advisories, tc.authored)
			switch {
			case tc.want && !found:
				t.Fatalf("no advisory for the near-miss %q; advisories: %+v",
					tc.authored, advisories)
			case !tc.want && found:
				t.Fatalf("an ordinary name %q raised advisory %+v; the advisory "+
					"is targeted, never blanket", tc.authored, got)
			case !tc.want:
				return
			}

			if got.Authored != tc.authored {
				t.Errorf("authored spelling = %q; want %q", got.Authored, tc.authored)
			}
			if got.Reserved != reservedKey {
				t.Errorf("reserved spelling = %q; want the literal %q",
					got.Reserved, reservedKey)
			}
		})
	}
}

// REQ-61: "Advisory, not a failure, because the name is legal and this RDR
// must not reject a tag it does not own."
// REQ-63: "It travels on an advisory channel distinct from the
// validation-failure list — it carries **no** `reserved_tag_key` category
// discriminator, because it is not a validation failure and MUST NOT
// participate in category dispatch or alter the load/lint verdict."
// REQ-81: TS-4 "The advisory's payload is asserted byte-for-byte like the
// failure payload: authored spelling, reserved spelling `recognized`, rule
// identifier `reserved-tag-key/near-miss`, and **no** `reserved_tag_key`
// category discriminator. The advisory MUST NOT change the load/lint verdict
// for any of them."
// REQ-62: "the advisory carries the same machine-readable discipline as the
// failure payload (block 3) … two distinct fields, the **authored spelling**
// and the **reserved spelling** (the literal `recognized`), plus a stable rule
// identifier whose value is the literal `reserved-tag-key/near-miss`."
// REQ-100: Phase 3's near-miss advisory assertions.
// BOUNDARY (net-new)
//
// Three obligations at once: the verdict is unchanged (the model still loads
// clean and non-nil), the advisory carries no category discriminator, and the
// payload is byte-for-byte.
func TestReq61And62And63And81_AdvisoryIsNonBlockingCategorylessAndByteExact(t *testing.T) {
	fixtures := map[string]string{
		"pos-near-miss-folded.toml": "Recognized",
		"pos-near-miss-upper.toml":  "RECOGNIZED",
		"pos-near-miss-space.toml":  " recognized",
		"pos-near-miss-both.toml":   " Recognized",
	}

	for fixture, authored := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			m, advisories, err := table.LoadWithAdvisories(readFixture(t, fixture), fixture)

			// The verdict is unchanged: the name is legal.
			if err != nil {
				t.Fatalf("the near-miss altered the load verdict: %v", err)
			}
			if m == nil {
				t.Fatal("clean load returned a nil model")
			}

			got, ok := advisoryFor(advisories, authored)
			if !ok {
				t.Fatalf("no advisory for %q", authored)
			}
			if got.Authored != authored {
				t.Errorf("authored = %q; want %q", got.Authored, authored)
			}
			if got.Reserved != reservedKey {
				t.Errorf("reserved = %q; want %q", got.Reserved, reservedKey)
			}
			if got.Rule != ruleNearMiss {
				t.Errorf("rule identifier = %q; want the literal %q",
					got.Rule, ruleNearMiss)
			}
			if got.Authored == got.Reserved {
				t.Error("authored and reserved are the same value; block 3's " +
					"discipline requires two DISTINCT fields")
			}
			if carriesCategory(got) {
				t.Errorf("the advisory carries a category discriminator (%+v); it "+
					"is not a validation failure and must not participate in "+
					"category dispatch", got)
			}
		})
	}
}

// REQ-64: "The carrier's Go type and field names are RDR 0002's to choose, on
// the same footing as the failure payload; this RDR constrains the values and
// the channel's separateness."
// DOMAIN EDGE (net-new)
//
// The constraint asserted here is SEPARATENESS: the advisory channel is not
// the validation-failure channel. A failing load and a near-miss are
// independent — a model can carry both, and the advisory does not become a
// failure nor the failure an advisory.
func TestReq64_TheAdvisoryChannelIsSeparateFromTheFailureChannel(t *testing.T) {
	// A clean load with a near-miss: advisories present, error nil.
	m, advisories, err := table.LoadWithAdvisories(
		readFixture(t, "pos-near-miss-folded.toml"), "pos-near-miss-folded.toml")
	if err != nil || m == nil {
		t.Fatalf("a near-miss-only model must load clean; got model=%v err=%v", m != nil, err)
	}
	if len(advisories) == 0 {
		t.Fatal("the advisory channel is empty on a near-miss model")
	}

	// A failing load: the error is a categorized failure, and no advisory has
	// been promoted into it.
	_, _, ferr := table.LoadWithAdvisories(
		readFixture(t, "neg/neg-recognized-owned.toml"), "neg-recognized-owned.toml")
	if ferr == nil {
		t.Fatal("the owned-named-recognized fixture loaded clean")
	}
	var f *table.Failure
	if !errors.As(ferr, &f) {
		t.Fatalf("the failure is not a categorized table.Failure: %v", ferr)
	}
	if f.Rule == ruleNearMiss {
		t.Error("an advisory rule identifier reached the validation-failure " +
			"channel; the two channels are distinct")
	}
}

// REQ-65: "**Naming** — canonical name `recognized`, matching the shipped
// `recognizedTagKey` constant and RDR 0001's frozen fixture. Rejected: a
// sigil-guarded name"
// BOUNDARY (regression pin — landed in RDR 0002)
//
// A NEGATIVE requirement: the canonical name is the bare literal, not a
// sigil-guarded variant. Pinned on the exported spelling, which is the only
// place a sigil could be introduced.
func TestReq65_CanonicalNameIsTheBareLiteralWithNoSigil(t *testing.T) {
	if table.RecognizedTagKey != reservedKey {
		t.Fatalf("RecognizedTagKey = %q; want the bare literal %q",
			table.RecognizedTagKey, reservedKey)
	}
	for _, sigil := range []string{"@", "$", "_", "!", ":", "%"} {
		if len(table.RecognizedTagKey) > 0 &&
			string(table.RecognizedTagKey[0]) == sigil {
			t.Errorf("the canonical name is sigil-guarded with %q", sigil)
		}
	}
}

// REQ-73: `0008:MVV` "**Normalizer half — executed inside RDR 0002's
// implementation.** A table declaring a recognized-provenance tag named
// `recognized` loads and lints clean; the same table with the declaration
// renamed (and separately, with an owned tag named `recognized`) fails
// load/lint in the `reserved_tag_key` category whose failure data carries the
// direction-appropriate remedy name and rule identifier — not a silent
// no-match at resolve time."
// REQ-75: TS-2 "**Expected**: both fail load/lint in the `reserved_tag_key`
// data-level category before any resolution — not a silent no-match at resolve
// time."
// HAPPY PATH (mixed: the categories landed in 0002, the payload is net-new)
func TestReq73And75_NormalizerHalfCleanLoadPlusBothFailuresWithPayload(t *testing.T) {
	// Clean: the canonical fixture declares the recognized tag under the
	// reserved name.
	m := mustLoad(t, rdrFixture)
	if decl, ok := m.Tags[reservedKey]; !ok {
		t.Fatalf("%s carries no %q declaration", rdrFixture, reservedKey)
	} else if decl.Provenance != table.ProvenanceRecognized {
		t.Errorf("%s's %q declaration has provenance %q; want recognized",
			rdrFixture, reservedKey, decl.Provenance)
	}

	// Both failure directions, before any resolution, with their payloads.
	directions := map[string]struct {
		fixture string
		remedy  string
		rule    string
	}{
		"declaration renamed":           {"neg/neg-recognized-misnamed.toml", reservedKey, ruleKernelOwned},
		"owned tag named recognized":    {"neg/neg-recognized-owned.toml", "", ruleAuthorMustRename},
		"observed tag named recognized": {"neg/neg-recognized-observed.toml", "", ruleAuthorMustRename},
	}
	for name, d := range directions {
		t.Run(name, func(t *testing.T) {
			f := reservedTagKeyFailure(t, d.fixture)
			if f.Remedy != d.remedy {
				t.Errorf("remedy = %q; want %q", f.Remedy, d.remedy)
			}
			if f.Rule != d.rule {
				t.Errorf("rule identifier = %q; want %q", f.Rule, d.rule)
			}
		})
	}
}
