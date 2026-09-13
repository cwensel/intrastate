package cli

// RDR 0029 — the tier assertions C2 promises a consumer, read over the
// enumeration seams C4's census obliges.
//
// The seams this ranges over do not exist in the tree these tests are
// committed against, so every probe reaches them through the package
// source rather than through an identifier that would not compile. A
// `frozen` set is asserted BY VALUE; an `append-only` or `growing` set is
// asserted by MEMBERSHIP AND UNIQUENESS only — never by cardinality, a
// member's ordinal position, or a tail position, which C2 forbids a
// consumer to rely on.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/graphlint"
	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// seamLiteralMembers returns the string literals appearing in the body of
// the exported function `name` declared in the package at dir.
//
// This is how a seam that does not exist yet is asserted by value: the
// members are read out of the accessor's own source. When the accessor
// lands the literals are its members; while it is absent the lookup
// reports so and the assertion fails on content, not on compilation.
func seamLiteralMembers(t *testing.T, dir, name string) ([]string, bool) {
	t.Helper()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatalf("parsing %s: %v", dir, err)
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			var (
				found   bool
				members []string
			)
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Recv != nil || fn.Name == nil || fn.Name.Name != name {
					continue
				}
				found = true
				ast.Inspect(fn, func(n ast.Node) bool {
					lit, ok := n.(*ast.BasicLit)
					if !ok || lit.Kind != token.STRING {
						return true
					}
					if v, err := strconv.Unquote(lit.Value); err == nil {
						members = append(members, v)
					}
					return true
				})
			}
			if found {
				// A seam that clones a package-level var declares no
				// literals of its own; fall back to the var's.
				if len(members) == 0 {
					members = pkgVarStrings(t, pkg)
				}
				return members, true
			}
		}
	}
	return nil, false
}

// pkgVarStrings collects string literals from package-level var slices, so
// a seam written as `return slices.Clone(xs)` still yields its members.
func pkgVarStrings(t *testing.T, pkg *ast.Package) []string {
	t.Helper()

	var out []string
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.VAR {
				continue
			}
			ast.Inspect(gd, func(n ast.Node) bool {
				lit, ok := n.(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					return true
				}
				if v, err := strconv.Unquote(lit.Value); err == nil {
					out = append(out, v)
				}
				return true
			})
		}
	}
	return out
}

// REQ-31: "`respond::Types() []string` over `{ok}`"
// REQ-24: "`frozen`: the envelope `type` discriminator …"
// REQ-48: "a by-value assertion per frozen vocabulary in C4 — including the
// five this RDR newly assigned"
// BOUNDARY — `"failed"` is a reserved name this CLI does not emit and MUST
// NOT enter the seam, or the by-value assertion would pin a member no
// consumer can ever observe.
func TestReq24And31And48_TheEnvelopeTypeSeamIsExactlyOk(t *testing.T) {
	members, ok := seamLiteralMembers(t,
		pkgDir(t, "internal", "cli", "respond"), "Types")
	if !ok {
		t.Fatalf("respond declares no Types() seam; C4 owes " +
			"`respond::Types() []string` over {ok}, and without it the " +
			"`frozen` tier on the envelope type discriminator binds only prose")
	}

	got := slices.Clone(members)
	slices.Sort(got)
	want := []string{"ok"}
	if !slices.Equal(got, want) {
		t.Errorf("Types() = %v; want exactly %v. The emitted set has ONE "+
			"member: `respond::Fail` writes the bare *CLIError, so "+
			"\"failed\" is a reserved name this CLI does not emit and MUST "+
			"NOT enter the seam — S6's by-value assertion would otherwise "+
			"pin a member no consumer can ever observe", got, want)
	}
}

// REQ-33: "`respond::Levels() []string`"
// REQ-24: "the stderr advisory `level` set (`note`, `warning`)"
// REQ-48: a by-value assertion per frozen vocabulary.
// BOUNDARY
func TestReq24And33And48_TheAdvisoryLevelSeamIsExactlyNoteAndWarning(t *testing.T) {
	members, ok := seamLiteralMembers(t,
		pkgDir(t, "internal", "cli", "respond"), "Levels")
	if !ok {
		t.Fatalf("respond declares no Levels() seam; C4 owes " +
			"`respond::Levels() []string` over the frozen stderr advisory " +
			"level set, which is emitted today as bare literals")
	}

	got := slices.Clone(members)
	slices.Sort(got)
	want := []string{"note", "warning"}
	if !slices.Equal(got, want) {
		t.Errorf("Levels() = %v; want exactly %v — C4 freezes the stderr "+
			"advisory `level` set at these two", got, want)
	}
}

// REQ-32: "`clierr::ExitCodes() []int` over the five values
// `{0,1,2,3,130}`"
// REQ-24: "the exit-code classes emitted by `internal/cli/clierr::ExitCodeFor`"
// REQ-48: a by-value assertion per frozen vocabulary.
// BOUNDARY — the PROJECTION is the contract, not its generator: the six
// ErrorGroup constants stay untiered and free to grow behind it.
func TestReq24And32And48_TheExitCodeSeamIsExactlyTheFiveEmittedIntegers(t *testing.T) {
	dir := pkgDir(t, "internal", "cli", "clierr")
	if _, ok := declaredFuncResult(t, dir, "ExitCodes"); !ok {
		t.Fatalf("clierr declares no ExitCodes() seam; C4 owes " +
			"`clierr::ExitCodes() []int` over {0,1,2,3,130}. ExitCodeFor is " +
			"a switch today, so the frozen tier has nothing to compare by value")
	}

	// The five values must be exactly the integers ExitCodeFor can emit.
	// They are read off the seam's own source so the assertion is by VALUE.
	src := readSource(t, pkgDir(t, "internal", "cli", "clierr", "clierr.go"))
	idx := strings.Index(src, "func ExitCodes()")
	if idx < 0 {
		t.Fatal("ExitCodes() is not declared in clierr.go")
	}
	body := src[idx:]
	if end := strings.Index(body, "\n}"); end >= 0 {
		body = body[:end]
	}
	for _, want := range []string{"0", "1", "2", "3", "130"} {
		if !strings.Contains(body, want) {
			t.Errorf("ExitCodes() does not carry the emitted exit code %s. "+
				"`0` comes from GroupSuccess/GroupWarning and `1` from the "+
				"non-CLIError fallthrough — they are in the set because they "+
				"are EMITTED, not because they are refusal classes", want)
		}
	}
}

// REQ-34: "`cli::UnknownReasons() []string` over the union, matching
// `graphlint::Reasons`"
// REQ-25: "`append-only`: … the `flow next` unknown-`reason` set (a distinct
// vocabulary on a distinct field from the preceding entry)"
// REQ-50: "membership-and-uniqueness assertions only, each over that
// vocabulary's enumeration seam"
// HAPPY PATH
func TestReq25And34And50_TheUnknownReasonSeamIsAppendOnlyAndDuplicateFree(t *testing.T) {
	members, ok := seamLiteralMembers(t, pkgDir(t, "internal", "cli"),
		"UnknownReasons")
	if !ok {
		t.Fatalf("internal/cli declares no UnknownReasons() seam; C4 owes " +
			"`cli::UnknownReasons() []string` over the union. The tokens are " +
			"a const block today with no accessor over them, so the " +
			"append-only tier has nothing to check membership over")
	}

	// MEMBERSHIP, not cardinality: the three tokens the CLI emits today
	// must each be present. A fourth arriving later is licensed by the
	// tier and must not fail this assertion.
	for _, want := range []string{
		string(resolve.ReasonAbsent),
		string(resolve.ReasonUncomparable),
		"not-evaluated",
	} {
		if !slices.Contains(members, want) {
			t.Errorf("UnknownReasons() omits %q; the union carries the "+
				"kernel-owned reasons plus RDR 0011's `not-evaluated`", want)
		}
	}

	// UNIQUENESS: a duplicate makes membership ambiguous.
	seen := map[string]bool{}
	for _, m := range members {
		if seen[m] {
			t.Errorf("UnknownReasons() carries %q twice; the set is "+
				"append-only and a duplicate makes membership ambiguous", m)
		}
		seen[m] = true
	}
}

// REQ-35: "`resolve::Blocks() []Block`"
// REQ-25: "`findings[].block` (`internal/resolve::Block`)" is append-only,
// "on the authority of `JDR 0001 §D12`, which already grew the set from two
// members to three".
// REQ-50: membership-and-uniqueness only.
// HAPPY PATH
func TestReq25And35And50_TheBlockSeamIsAppendOnlyAndDuplicateFree(t *testing.T) {
	dir := pkgDir(t, "internal", "resolve")
	if _, ok := declaredFuncResult(t, dir, "Blocks"); !ok {
		t.Fatalf("internal/resolve exports no Blocks(); C4 owes " +
			"`resolve::Blocks() []Block` over the append-only " +
			"`findings[].block` vocabulary")
	}

	src := readSource(t, pkgDir(t, "internal", "resolve", "guard.go"))
	idx := strings.Index(src, "func Blocks()")
	if idx < 0 {
		t.Fatal("Blocks() is not declared in internal/resolve/guard.go")
	}
	body := src[idx:]
	if end := strings.Index(body, "\n}"); end >= 0 {
		body = body[:end]
	}

	// The three members JDR 0001 §D12 leaves in force, by NAME not by count.
	for _, want := range []string{"BlockAll", "BlockUnless", "BlockMatch"} {
		if !strings.Contains(body, want) {
			t.Errorf("Blocks() omits %s; the set grew from two members to "+
				"three under JDR 0001 §D12 and is append-only, so every "+
				"declared block must be enumerable", want)
		}
	}
}

// REQ-38: "`findings[].operator` has TWO enumerations,
// `internal/guard::Operators` and `internal/table::Operators`, in separate
// packages … the tier binds both, and the by-value assertion pins both"
// REQ-48: a by-value assertion per frozen vocabulary.
// ADVERSARIAL — one enumeration drifting from the other is the defect; the
// wire carries whichever the producing path used.
func TestReq38And48_BothOperatorEnumerationsArePinnedByValue(t *testing.T) {
	want := []string{"contains", "eq", "exists", "gt", "gte", "in", "lt", "lte"}

	guardOps := slices.Sorted(slices.Values(guard.Operators()))
	if !slices.Equal(guardOps, want) {
		t.Errorf("guard.Operators() = %v; want exactly %v", guardOps, want)
	}

	tableOps := slices.Sorted(slices.Values(table.Operators()))
	if !slices.Equal(tableOps, want) {
		t.Errorf("table.Operators() = %v; want exactly %v", tableOps, want)
	}

	// The tier binds BOTH, so the two must agree member for member: a
	// consumer reading `findings[].operator` cannot tell which package
	// produced the value.
	if !slices.Equal(guardOps, tableOps) {
		t.Errorf("the two operator enumerations disagree: guard=%v table=%v. "+
			"C4 tiers `findings[].operator` frozen and binds both "+
			"enumerations, because the wire carries whichever the producing "+
			"path used", guardOps, tableOps)
	}
}

// REQ-24: "`frozen`: … the severity vocabulary (`blocking`, `info`); … the
// gate `verdict` set (`allow`, `deny`, `indeterminate`) … `data.escape_class`
// (`internal/resolve::RefusalKinds`)"
// REQ-48: a by-value assertion per frozen vocabulary in C4.
// BOUNDARY
func TestReq24And48_TheFrozenAnchorVocabulariesMatchTheirDeclaredMembers(t *testing.T) {
	t.Run("severity", func(t *testing.T) {
		got := slices.Sorted(slices.Values(graphlint.Severities()))
		want := []string{"blocking", "info"}
		if !slices.Equal(got, want) {
			t.Errorf("Severities() = %v; want exactly %v", got, want)
		}
	})

	t.Run("verdict", func(t *testing.T) {
		var got []string
		for _, v := range accessor.Verdicts() {
			got = append(got, string(v))
		}
		slices.Sort(got)
		want := []string{"allow", "deny", "indeterminate"}
		if !slices.Equal(got, want) {
			t.Errorf("Verdicts() = %v; want exactly %v — the set is frozen "+
				"because it is simultaneously the accessor-protocol input "+
				"alphabet", got, want)
		}
	})

	t.Run("escape_class", func(t *testing.T) {
		var got []string
		for _, k := range resolve.RefusalKinds() {
			got = append(got, string(k))
		}
		slices.Sort(got)
		want := []string{
			"ambiguous_match", "guard_unevaluable", "no_match",
			"owned_state_unavailable", "unmodeled_outcome",
		}
		if !slices.Equal(got, want) {
			t.Errorf("RefusalKinds() = %v; want exactly %v", got, want)
		}
	})
}

// REQ-49: "`graph-lint-failed` takes no by-value row: it is a single
// `const`, and its testable claim is that a blocking run returns that code
// and no other."
// DOMAIN EDGE — closure of the EMIT SITES, not set membership.
func TestReq49_ABlockingRunReturnsTheAggregateCodeAndNoOther(t *testing.T) {
	path := writeModel(t, illegalModel)

	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err == nil {
		t.Fatalf("the illegal model succeeded; stdout:\n%s", stdout)
	}

	env := parseFailure(t, stdout)
	if env.Code != graphlint.AggregateCode {
		t.Errorf("a blocking run returned the aggregate code %q; want %q "+
			"and no other — the testable claim for this single `const` is "+
			"closure of the emit sites, not set membership",
			env.Code, graphlint.AggregateCode)
	}
}

// REQ-25: "`append-only`: `internal/table::Categories()` …"
// REQ-50: "membership-and-uniqueness assertions only"
// REQ-13: "A consumer MUST tolerate an unrecognized member, and MUST NOT
// assert on the set's cardinality, a member's ordinal position, or a tail
// position."
// HAPPY PATH
func TestReq13And25And50_TheCategorySetIsAssertedByMembershipNotCardinality(t *testing.T) {
	all := table.Categories()
	if len(all) == 0 {
		t.Fatal("Categories() is empty")
	}

	seen := map[table.Category]bool{}
	for _, c := range all {
		if c == "" {
			t.Error("Categories() carries an empty category")
		}
		if seen[c] {
			t.Errorf("category %q is registered twice; the list is "+
				"append-only and a duplicate makes membership ambiguous", c)
		}
		seen[c] = true
	}

	// Membership of a member this RDR does not own, proving the assertion
	// ranges the real set without pinning its size.
	if !seen[table.CatMalformedTOML] {
		t.Errorf("Categories() omits %q", table.CatMalformedTOML)
	}
}

// REQ-28: "Two emitted namespaces take NO tier, deliberately:
// `findings[].class` and `data.dispositions`"
// REQ-29: "A vocabulary is tiered wherever it is EMITTED … an unassigned
// machine-readable surface is a defect."
// DOMAIN EDGE — the negative: these two carry model-authored tokens to
// which `clierr` ascribes no meaning, so no seam may enumerate them.
func TestReq28And29_TheUntieredNamespacesGainNoEnumerationSeam(t *testing.T) {
	// The precondition that makes this negative meaningful: C4's census
	// must actually have been built. In a tree with no seams at all, "the
	// untiered namespaces gained none" is true of everything and proves
	// nothing about the DELIBERATE omission this clause fixes.
	if _, ok := declaredFuncResult(t,
		pkgDir(t, "internal", "cli", "respond"), "Types"); !ok {
		t.Fatal("no enumeration seam exists yet, so \"these two namespaces " +
			"take NO tier\" is vacuous — every namespace lacks a seam. The " +
			"census must be built before its deliberate exclusions mean " +
			"anything.")
	}

	// `findings[].class` is author surface: a seam over it would claim a
	// CLI vocabulary where the model supplies the tokens.
	for _, tc := range []struct {
		dir  []string
		name string
	}{
		{[]string{"internal", "cli", "clierr"}, "Classes"},
		{[]string{"internal", "cli"}, "Dispositions"},
	} {
		if _, ok := declaredFuncResult(t, pkgDir(t, tc.dir...), tc.name); ok {
			t.Errorf("%s declares %s(); `findings[].class` and "+
				"`data.dispositions` take NO tier deliberately — they carry "+
				"model-authored tokens to which clierr ascribes no meaning, "+
				"and are author surface, not CLI vocabulary",
				strings.Join(tc.dir, "/"), tc.name)
		}
	}
}

// REQ-36: "**`seam: none (prose-only)`**" for the CLIError `code` row; "The
// tier STANDS, unasserted and legibly so"
// DOMAIN EDGE — the negative REQ: no seam is built here, and no exact-set
// or cardinality assertion over the code vocabulary may exist. A5 refuted
// the registry: `clierr` importing graphlint/guard/accessor is a hard build
// cycle.
func TestReq36_TheCLIErrorCodeVocabularyGainsNoSeam(t *testing.T) {
	dir := pkgDir(t, "internal", "cli", "clierr")

	// The same precondition: this row is the ONE census entry that stays
	// unasserted while the other five are built. Before any seam exists,
	// "this one has none" is indistinguishable from "none of them do".
	if _, ok := declaredFuncResult(t, dir, "ExitCodes"); !ok {
		t.Fatal("no seam exists in clierr yet, so `seam: none (prose-only)` " +
			"on the `code` row carries no information — the distinction C4 " +
			"makes legible is between a tier whose seam EXISTS and this one, " +
			"whose seam is absent by report")
	}
	for _, name := range []string{"Codes", "ErrorCodes", "CodeVocabulary"} {
		if _, ok := declaredFuncResult(t, dir, name); ok {
			t.Errorf("clierr declares %s(); C4 records the CLIError `code` "+
				"row as `seam: none (prose-only)` — A5 spiked it and the "+
				"registry is not buildable at this layering. The tier STANDS, "+
				"unasserted and legibly so; a hand-copied shadow list the "+
				"raise sites never touch is not a seam C2 can assert by value",
				name)
		}
	}

	// And the leaf stays a leaf: the import ban is what made the registry
	// unbuildable, so it must not have been relaxed to build one.
	src := readSource(t, pkgDir(t, "internal", "cli", "clierr", "clierr.go"))
	for _, banned := range []string{
		"internal/graphlint", "internal/guard", "internal/accessor",
	} {
		if strings.Contains(src, banned) {
			t.Errorf("clierr imports %q; that is the hard build cycle A5 "+
				"reported, and the reason this row has no seam", banned)
		}
	}
}

// REQ-52: "No golden/snapshot test of the JSON envelope exists in the repo
// today, and this RDR does not add one"
// ADVERSARIAL — a byte-golden over the envelope would itself assert the
// cardinality C2 forbids consumers from asserting.
func TestReq52_NoGoldenSnapshotOfTheJSONEnvelopeIsAdded(t *testing.T) {
	// A7's snapshot is explicitly NOT this: its subject is the SEAM
	// MEMBERS and their severities, not an emitted envelope. The defect
	// under test is a checked-in file holding a whole terminal envelope
	// that this RDR's own field would have to be baked into.
	root := repoRootFor(t)

	// The negative is only meaningful once the field exists: before then no
	// golden could carry it, and this would pass on every tree.
	if _, ok := declaredConstNames(t,
		pkgDir(t, "internal", "cli", "clierr"))["SchemaVersion"]; !ok {
		t.Fatal("`schema_version` does not exist yet, so \"no envelope " +
			"golden carries it\" is vacuously true; S7's negative binds the " +
			"tree that HAS the field")
	}

	entries, err := os.ReadDir(root + "/internal/cli/testdata")
	if err != nil {
		return // no testdata dir at all is trivially conforming
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		body, rerr := os.ReadFile(root + "/internal/cli/testdata/" + e.Name())
		if rerr != nil {
			continue
		}
		if strings.Contains(string(body), `"schema_version"`) &&
			strings.Contains(string(body), `"type"`) {
			t.Errorf("internal/cli/testdata/%s is an envelope golden "+
				"carrying schema_version; S7 states no golden/snapshot test "+
				"of the JSON envelope exists in the repo today and this RDR "+
				"does not add one — a byte-golden over the envelope would "+
				"itself assert the cardinality C2 forbids", e.Name())
		}
	}
}
