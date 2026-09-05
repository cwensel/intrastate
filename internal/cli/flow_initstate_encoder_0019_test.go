package cli

// RDR 0019 — the SEED ENCODER: `0019:C1`'s "Seed values are taken from the
// LOADER-NORMALIZED `Model.Initial` assignments" paragraph, `0019:RT3`, and
// scenario S4's kind table.
//
// The suite's discriminating choice: every equality is asserted on ARTIFACT
// BYTES between two independently-built artifacts — one seeded by
// `init-state`, one written by `set-state --write`. RT3's claim is BYTE
// identity, so an oracle that compared `read-state` reports would pass
// against an encoder that persisted `["draft"]` and a reader that unwrapped
// it, which is precisely the third encoding C1 forbids.

import (
	"os"
	"strings"
	"testing"
)

// initKindModel is S4's fixture: one owned tag per value kind BOTH routes
// admit, all served by ONE file-backed writer, with `[initial]` assigning
// every one of them.
//
// The nine kinds `0019:S4` fixes: enum, scalar string, bool, int, float,
// set array, set with duplicates and HTML characters, empty set, and
// single-member set.
const initKindModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "initkinds"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.k_enum]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true

[tags.k_scalar]
provenance = "owned"
kind = "scalar"

[tags.k_bool]
provenance = "owned"
kind = "bool"
single_valued = true

[tags.k_int]
provenance = "owned"
kind = "int"
single_valued = true

[tags.k_float]
provenance = "owned"
kind = "float"
single_valued = true

[tags.k_set]
provenance = "owned"
kind = "set"
elements = ["a", "b", "c"]

[tags.k_set_html]
provenance = "owned"
kind = "set"
elements = ["a<b", "x&y", "plain"]

[tags.k_set_empty]
provenance = "owned"
kind = "set"
elements = ["a", "b"]

[tags.k_set_one]
provenance = "owned"
kind = "set"
elements = ["only", "other"]

[read.state]
role = "state"
path = "flow.state"
keys = ["k_enum", "k_scalar", "k_bool", "k_int", "k_float", "k_set", "k_set_html", "k_set_empty", "k_set_one"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["k_enum", "k_scalar", "k_bool", "k_int", "k_float", "k_set", "k_set_html", "k_set_empty", "k_set_one"]
timeout = "2s"
read_back = true

[initial]
k_enum = "draft"
k_scalar = "plain text"
k_bool = "true"
k_int = "5"
k_float = "1"
k_set = ["b", "a"]
k_set_html = ["x&y", "a<b", "a<b"]
k_set_empty = []
k_set_one = ["only"]

[context.done]
[context.done.match.k_enum]
eq = "final"

[[rule]]
id = "advance"
[rule.match.k_enum]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
k_enum = "final"
`

// initLoaderOnlyModel carries the THREE kinds only the LOADER admits — the
// admission-set divergence C1 turns on:
//
//   - `bare_for_set`   a BARE SCALAR assigned to a SET-valued tag;
//   - `array_for_one`  a SINGLE-MEMBER ARRAY LITERAL assigned to a SCALAR tag;
//   - `note`           an EMPTY SCALAR.
//
// All three load and normalize; all three are refused on the ARGV surface
// by `::canonicalValue` under JDR 0001 §D11. Transcribing the seed through
// argv would therefore refuse at seed time models that load clean, which is
// the first-run wall this verb exists to remove (A2).
const initLoaderOnlyModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "initloaderonly"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.bare_for_set]
provenance = "owned"
kind = "set"
elements = ["solo", "other"]

[tags.array_for_one]
provenance = "owned"
kind = "scalar"

[tags.note]
provenance = "owned"
kind = "scalar"

[tags.anchor]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true

[read.state]
role = "state"
path = "flow.state"
keys = ["bare_for_set", "array_for_one", "note", "anchor"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["bare_for_set", "array_for_one", "note", "anchor"]
timeout = "2s"
read_back = true

[initial]
bare_for_set = "solo"
array_for_one = ["wrapped"]
note = ""
anchor = "draft"

[context.done]
[context.done.match.anchor]
eq = "final"

[[rule]]
id = "advance"
[rule.match.anchor]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
anchor = "final"
`

// initFloatBoundaryModel is S4's ADDITIONAL boundary row: `[initial]
// threshold = 1.0`, whose loader-normalized seed is `"1"` while `--write
// threshold=1.0` writes `"1.0"`. RT3 does NOT claim these agree, and the
// divergence is pinned as DECIDED behaviour.
const initFloatBoundaryModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "initfloat"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.threshold]
provenance = "owned"
kind = "float"
single_valued = true

[tags.anchor]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true

[read.state]
role = "state"
path = "flow.state"
keys = ["threshold", "anchor"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["threshold", "anchor"]
timeout = "2s"
read_back = true

[initial]
threshold = 1.0
anchor = "draft"

[context.done]
[context.done.match.anchor]
eq = "final"

[[rule]]
id = "advance"
[rule.match.anchor]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
anchor = "final"
`

// --- the RT3 comparison harness -----------------------------------------

// seedByInit runs `init-state` into a FRESH artifact and returns its bytes.
func seedByInit(t *testing.T, model string) []byte {
	t.Helper()

	art := newFlowArtifact(t, "init.artifact")
	requireSuccess(t, initStateArgs(model, artifactBinding(initRoleA, art))...)

	b, err := os.ReadFile(art)
	if err != nil {
		t.Fatalf("init-state left no artifact at %s: %v", art, err)
	}
	return b
}

// seedBySetState writes the same assignments into a SECOND fresh artifact
// through `set-state --write` and returns its bytes.
func seedBySetState(t *testing.T, model string, writes ...string) []byte {
	t.Helper()

	art := newFlowArtifact(t, "set.artifact")
	args := []string{"flow", "set-state", "--model", model,
		"--artifact", artifactBinding(initRoleA, art)}
	for _, w := range writes {
		args = append(args, "--write", w)
	}
	requireSuccess(t, append(args, "--as=json")...)

	b, err := os.ReadFile(art)
	if err != nil {
		t.Fatalf("set-state left no artifact at %s: %v", art, err)
	}
	return b
}

// REQ-17: "Seed values are taken from the LOADER-NORMALIZED `Model.Initial`
// assignments, NOT by transcribing them to their `--write` argv spelling"
// REQ-18: "each value is rendered in JDR 0001 §D13's canonical form,
// DISPATCHED ON THE DECLARED KIND, which is what makes read-back equality
// byte equality and leaves no third encoding."
// REQ-71 / `0019:RT3`: "`init-state` into a fresh artifact and the
// `set-state --write` transcription of the same `[initial]` assignments
// into a second fresh artifact produce byte-identical artifacts, for every
// value kind BOTH routes admit AND every VALUE whose two routes carry the
// same spelling"
// REQ-80 / `0019:S4`: "table-driven over every value kind BOTH routes admit
// (9 kinds: enum, scalar string, bool, int, float, set array, set with
// duplicates and HTML characters, empty set, single-member set) …
// **Expected**: the two artifacts are byte-identical (RT3; this is A2's
// spike promoted to a standing test)."
// REQ-81 / SC-4 argv spelling: "Each row's `--write` argv MUST use the
// spelling the loader normalizes to (`1`, not `1.0`; `5`, not `+5`) — the
// numeric rows otherwise diff on a transcription difference RT3 does not
// claim, and a test written the other way fails on day one for the wrong
// reason (A2)."
// HAPPY PATH — the whole nine-kind table in ONE comparison, since RT3's
// claim is over the two ARTIFACTS, not per key.
func TestReq17And18And71And80And81_0019_InitAndSetStateProduceByteIdenticalArtifacts(t *testing.T) {
	model := writeFlowModel(t, initKindModel)

	byInit := seedByInit(t, model)

	// Every argv spelling below is the one the LOADER NORMALIZES TO
	// (REQ-81): `5` not `+5`, `1` not `1.0`. A row written the other way
	// would diff on a transcription difference RT3 does not claim.
	bySet := seedBySetState(t, model,
		"k_enum=draft",
		"k_scalar=plain text",
		"k_bool=true",
		"k_int=5",
		"k_float=1",
		`k_set=["b","a"]`,
		`k_set_html=["x&y","a<b","a<b"]`,
		`k_set_empty=[]`,
		`k_set_one=["only"]`,
	)

	if string(byInit) != string(bySet) {
		t.Errorf("the two artifacts DIFFER; RT3 fixes BYTE identity across "+
			"the two routes\ninit-state:  %s\nset-state:   %s",
			string(byInit), string(bySet))
	}
}

// REQ-19: "the seed encoder is exactly the SET arm of `set-state`'s: for
// `decl.Kind == \"set\"` it is `internal/cli/flow_input.go::canonicalSet`
// (sort, compact, JSON array); for every scalar kind it is the single
// member verbatim, `members[0]`."
// REQ-20: "`::canonicalSet` alone is NOT the whole encoder — it takes
// `[]string` and renders a JSON array, so applying it to a scalar seed
// would persist `[\"draft\"]` where `set-state --write status=draft`
// persists `draft`"
// ADVERSARIAL — the named defect, asserted directly on the seeded artifact.
// A scalar seed persisted as a one-element JSON array is what a whole-table
// byte comparison would ALSO catch, but this row names the defect so a
// failure reads as the mistake C1 warns about.
func TestReq19And20_0019_ScalarSeedsArePersistedVerbatimNotWrappedInAnArray(t *testing.T) {
	model := writeFlowModel(t, initKindModel)
	art := newFlowArtifact(t, "state.artifact")

	requireSuccess(t, initStateArgs(model, artifactBinding(initRoleA, art))...)

	body := string(mustReadArtifact(t, art))
	for _, wrapped := range []string{
		`["draft"]`, `["plain text"]`, `["true"]`, `["5"]`, `["1"]`,
	} {
		if strings.Contains(body, wrapped) {
			t.Errorf("the artifact carries %s — `::canonicalSet` was applied "+
				"to a SCALAR seed. For every scalar kind the encoder is the "+
				"single member VERBATIM, `members[0]`\nartifact: %s",
				wrapped, body)
		}
	}
}

// REQ-19: the SET arm is `::canonicalSet` — sort, compact, JSON array.
// REQ-18: JDR 0001 §D13's canonical form.
// DOMAIN EDGE — the three properties of the canonical set form, each given
// a seed a non-conforming encoder would visibly change: the `[initial]`
// members are declared UNSORTED and with a DUPLICATE, and two carry HTML
// metacharacters that bare `json.Marshal` would escape.
func TestReq18And19_0019_SetSeedsTakeTheCanonicalSortedCompactUnescapedForm(t *testing.T) {
	model := writeFlowModel(t, initKindModel)
	art := newFlowArtifact(t, "state.artifact")

	requireSuccess(t, initStateArgs(model, artifactBinding(initRoleA, art))...)
	owned := ownedFromReadState(t, requireSuccess(t,
		readStateArgs(model, artifactBinding(initRoleA, art))...))

	// `["b","a"]` declared; sorted-compact is `["a","b"]`.
	requireReadsBack(t, owned, "k_set", `["a","b"]`)
	// `["x&y","a<b","a<b"]` declared; sorted, DEDUPLICATED, and rendered
	// with HTML escaping DISABLED.
	requireReadsBack(t, owned, "k_set_html", canonicalSetLiteral)
	requireReadsBack(t, owned, "k_set_empty", `[]`)
	requireReadsBack(t, owned, "k_set_one", `["only"]`)

	if body := string(mustReadArtifact(t, art)); strings.Contains(body, `<`) ||
		strings.Contains(body, `&`) {
		t.Errorf("the artifact carries HTML-ESCAPED members; the canonical "+
			"form disables HTML escaping, and %s is the defect spelling\n"+
			"artifact: %s", escapedSetLiteral, body)
	}
}

// REQ-21: "The argv encoder `::canonicalValue` is the correct *reference*
// for both arms' output but is NOT the call: it re-runs the argv admission
// checks this clause rules out below."
// REQ-22: "Conformance to the declaration is ALREADY HELD by the loader and
// is not re-established here … A re-conform pass on a loader-normalized
// value cannot fail and buys nothing"
// REQ-23: "the loader's `[initial]` admission set is a proper SUPERSET of
// the argv route's: a bare scalar for a set-valued tag and a SINGLE-MEMBER
// array literal for a scalar tag both load and normalize … and are refused
// only on the argv surface … Transcribing through argv would therefore
// refuse at seed time models that load clean, re-creating the first-run
// wall this verb exists to remove (A2)."
// REQ-83 / SC-4 loader-only kinds: "the 2 kinds only the loader admits
// (bare scalar for a set tag, array for a scalar tag) seed successfully via
// the normalized path and read back value-for-value"
// ADVERSARIAL — the whole point of NOT calling `::canonicalValue`. Each of
// these three seeds REFUSES on the argv surface, so an implementation that
// routed the seed through argv fails this test at exit 2 rather than
// seeding, which is the first-run wall re-created.
func TestReq21And22And23And83_0019_LoaderOnlyKindsSeedThroughTheNormalizedPath(t *testing.T) {
	model := writeFlowModel(t, initLoaderOnlyModel)
	art := newFlowArtifact(t, "state.artifact")

	requireSuccess(t, initStateArgs(model, artifactBinding(initRoleA, art))...)
	owned := ownedFromReadState(t, requireSuccess(t,
		readStateArgs(model, artifactBinding(initRoleA, art))...))

	cases := []struct {
		req, key, want, why string
	}{
		{
			req:  "REQ-23 / REQ-83",
			key:  "bare_for_set",
			want: `["solo"]`,
			why: "a BARE SCALAR assigned to a SET-valued tag; the loader " +
				"normalizes it to a one-member sequence, and the set arm " +
				"renders that as the canonical array. `::canonicalValue` " +
				"REFUSES this spelling on argv",
		},
		{
			req:  "REQ-23 / REQ-83",
			key:  "array_for_one",
			want: "wrapped",
			why: "a SINGLE-MEMBER ARRAY LITERAL assigned to a SCALAR tag; " +
				"the loader admits it at arity 1 and normalizes to one " +
				"member, which the scalar arm renders VERBATIM. " +
				"`::canonicalValue` REFUSES this spelling on argv",
		},
		{
			req:  "REQ-53 / REQ-54 / REQ-83",
			key:  "note",
			want: "",
			why: "an EMPTY SCALAR; the scalar arm renders `members[0]` and " +
				"the key reads back PRESENT. `::canonicalValue` REFUSES " +
				"`--write note=` on argv",
		},
	}

	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			requireReadsBack(t, owned, tc.key, tc.want)
			t.Logf("%s: %s", tc.req, tc.why)
		})
	}
}

// REQ-53 / F5: "An `[initial]` value that is an EMPTY SCALAR (`note =
// \"\"`) seeds like any other: the scalar arm renders `members[0]`, the
// artifact carries `{\"note\":\"\"}`, and the key reads back PRESENT —
// distinct from a cleared key, which leaves the object entirely."
// REQ-83: "`[initial] note = \"\"` seeds, the artifact carries
// `{\"note\":\"\"}`, and the key reads back PRESENT with an empty value —
// asserted distinct from a cleared key".
// BOUNDARY — the whole claim is a DISTINCTION, so the test must exhibit
// BOTH states: the seeded empty scalar (present, empty) and the same key
// CLEARED (absent from the object entirely). An assertion on the seeded
// state alone cannot discriminate a reader that collapses "" to absent.
func TestReq53And83_0019_EmptyScalarSeedsPresentAndIsDistinctFromCleared(t *testing.T) {
	model := writeFlowModel(t, initLoaderOnlyModel)
	art := newFlowArtifact(t, "state.artifact")
	bind := artifactBinding(initRoleA, art)

	requireSuccess(t, initStateArgs(model, bind)...)

	// The seeded state: PRESENT with an empty value, and the artifact
	// carries the key.
	seeded := ownedFromReadState(t, requireSuccess(t, readStateArgs(model, bind)...))
	requireReadsBack(t, seeded, "note", "")

	body := string(mustReadArtifact(t, art))
	if !strings.Contains(body, `"note":""`) {
		t.Errorf("the artifact does not carry `\"note\":\"\"`; C1 fixes that "+
			"an empty scalar seed PERSISTS as a present key with an empty "+
			"value\nartifact: %s", body)
	}

	// The cleared state: ABSENT, and the key leaves the object entirely.
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", bind, "--clear", "note", "--as=json")

	cleared := ownedFromReadState(t, requireSuccess(t, readStateArgs(model, bind)...))
	requireReadsAbsent(t, cleared, "note", "after --clear note")

	if after := string(mustReadArtifact(t, art)); strings.Contains(after, `"note"`) {
		t.Errorf("the artifact still carries the `note` KEY after --clear; a "+
			"clear is a key REMOVAL that leaves the object entirely, which is "+
			"what makes it DISTINCT from a seeded empty scalar\nartifact: %s",
			after)
	}
}

// REQ-55 / F5: "The ARGV route's contrary refusal (`::canonicalValue`
// rejects `--write note=`) is a loader/write-surface asymmetry predating
// this RDR and is NOT resolved here in either direction (F5, routed to RDR
// 0002)" — out of scope; no change to `::canonicalValue`.
// ADVERSARIAL — the non-goal. The tempting "fix" is to make the argv route
// admit the empty scalar for symmetry; F5 forbids resolving the asymmetry
// in EITHER direction, so `--write note=` must still refuse.
func TestReq55_0019_ArgvStillRefusesTheEmptyScalarWrite(t *testing.T) {
	// F5 is an ASYMMETRY between the loader route and the argv route, and
	// this record resolves it in NEITHER direction. Both sides must exist
	// before "not resolved" can be discriminated: against a tree with no
	// `init-state`, the argv refusal is simply shipped behaviour and this
	// test would restate it rather than pin the non-goal.
	if !flowSubcommand(t, initStateVerb) {
		t.Fatalf("the `flow` group registers no `%s`; F5's asymmetry needs "+
			"BOTH routes to exist before \"not resolved in either "+
			"direction\" can be discriminated", initStateVerb)
	}
	model := writeFlowModel(t, initLoaderOnlyModel)
	art := newFlowArtifact(t, "state.artifact")

	_, _, err := runCmd(t, "flow", "set-state", "--model", model,
		"--artifact", artifactBinding(initRoleA, art),
		"--write", "note=", "--as=json")
	if err == nil {
		t.Errorf("`set-state --write note=` was ADMITTED; the argv route's " +
			"contrary refusal is a predating asymmetry this record does NOT " +
			"resolve in either direction (F5, routed to RDR 0002)")
	}
}

// REQ-72 / `0019:RT3`: "the three kinds only the loader admits (bare scalar
// for a set tag, array for a scalar tag, empty scalar) have no `--write`
// transcription to compare against and are outside this invariant"
// DOMAIN EDGE — the scope carve-out, asserted in the direction that
// matters: each of the three has NO admissible argv transcription, so no
// byte comparison is even constructible for them. A suite that tried to
// include them in the RT3 table would fail at setup.
func TestReq72_0019_LoaderOnlyKindsHaveNoArgvTranscriptionToCompare(t *testing.T) {
	// RT3's carve-out is about the SEED route this record adds: these three
	// kinds are outside the invariant because they have no `--write`
	// transcription to compare a SEEDED artifact against. Without the verb
	// there is no invariant to be outside of.
	if !flowSubcommand(t, initStateVerb) {
		t.Fatalf("the `flow` group registers no `%s`; RT3's scope carve-out "+
			"is about the seed route this verb adds", initStateVerb)
	}
	model := writeFlowModel(t, initLoaderOnlyModel)

	// The argv spelling each loader-only kind WOULD need, and which
	// `::canonicalValue` refuses. That refusal is what puts them outside
	// RT3 rather than making RT3 false.
	for _, write := range []string{
		"bare_for_set=solo",         // bare scalar for a SET tag
		`array_for_one=["wrapped"]`, // array literal for a SCALAR tag
		"note=",                     // empty scalar
	} {
		t.Run(write, func(t *testing.T) {
			art := newFlowArtifact(t, "state.artifact")
			_, _, err := runCmd(t, "flow", "set-state", "--model", model,
				"--artifact", artifactBinding(initRoleA, art),
				"--write", write, "--as=json")
			if err == nil {
				t.Errorf("`--write %s` was ADMITTED on argv; RT3 excludes "+
					"this kind on the ground that it has NO `--write` "+
					"transcription to compare against", write)
			}
		})
	}
}

// REQ-73 / `0019:RT3`: "on the numeric kinds the two routes agree on the
// value and may DISAGREE on its spelling, so the invariant is scoped to the
// spelling the loader normalizes to … `[initial] threshold = 1.0` seeds
// `\"1\"` where `--write threshold=1.0` writes `\"1.0\"`"
// REQ-82 / SC-4 boundary row: "One ADDITIONAL row asserts the boundary
// rather than the invariant: `[initial] threshold = 1.0` seeded, against
// `--write threshold=1.0`, produces artifacts that DIFFER (`\"1\"` vs
// `\"1.0\"`) and both read back their own value — the divergence is pinned
// as decided behavior".
// BOUNDARY — the invariant's EDGE, asserted in the FAILING direction. This
// is the one row where byte identity must NOT hold; a test that only
// asserted identity everywhere would let an implementation "fix" the
// divergence and silently change decided behaviour.
func TestReq73And82_0019_TheFloatSpellingDivergenceIsDecidedBehaviour(t *testing.T) {
	model := writeFlowModel(t, initFloatBoundaryModel)

	initArt := newFlowArtifact(t, "init.artifact")
	requireSuccess(t, initStateArgs(model, artifactBinding(initRoleA, initArt))...)

	setArt := newFlowArtifact(t, "set.artifact")
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", artifactBinding(initRoleA, setArt),
		"--write", "threshold=1.0", "--write", "anchor=draft", "--as=json")

	initBytes := string(mustReadArtifact(t, initArt))
	setBytes := string(mustReadArtifact(t, setArt))

	if initBytes == setBytes {
		t.Errorf("the two artifacts are IDENTICAL; the boundary row fixes "+
			"that they DIFFER — the loader normalizes `1.0` to `1` while "+
			"`--write threshold=1.0` writes `1.0`, and RT3 does not claim "+
			"agreement on spelling\nboth: %s", initBytes)
	}

	// Both read back THEIR OWN value — the divergence is a spelling
	// difference, not a lost or corrupted value.
	initOwned := ownedFromReadState(t, requireSuccess(t,
		readStateArgs(model, artifactBinding(initRoleA, initArt))...))
	requireReadsBack(t, initOwned, "threshold", "1")

	setOwned := ownedFromReadState(t, requireSuccess(t,
		readStateArgs(model, artifactBinding(initRoleA, setArt))...))
	requireReadsBack(t, setOwned, "threshold", "1.0")
}

// mustReadArtifact reads an artifact that MUST exist. A missing artifact
// where the contract fixes a write is a failure, never a skip.
func mustReadArtifact(t *testing.T, path string) []byte {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("the artifact at %s is missing or unreadable: %v", path, err)
	}
	return b
}

// REQ-54: "It creates no third encoding and is not refused at seed time,
// because refusing it would require re-conforming a loader-normalized
// value, which this clause rules out above."
// ADVERSARIAL — two claims, and the second is the reason for the first. The
// empty scalar is the value a re-conform pass would REJECT, so an
// implementation that re-conformed at seed time refuses here. And "no third
// encoding" is asserted by byte-comparing the seeded artifact against the
// two encodings C1 admits: the scalar arm's `members[0]` verbatim, never a
// wrapped array and never an omitted key.
func TestReq54_0019_TheEmptyScalarIsNotRefusedAtSeedTimeAndAddsNoThirdEncoding(t *testing.T) {
	model := writeFlowModel(t, initLoaderOnlyModel)
	art := newFlowArtifact(t, "state.artifact")
	bind := artifactBinding(initRoleA, art)

	// NOT REFUSED at seed time. A re-conform pass over the loader-normalized
	// value — which this clause rules out — is the only thing that refuses
	// here, so a refusal names that defect.
	stdout, _, err := runCmd(t, initStateArgs(model, bind)...)
	if err != nil {
		t.Fatalf("the seed REFUSED (%v); an `[initial]` empty scalar seeds "+
			"like any other. Refusing it would require RE-CONFORMING a "+
			"loader-normalized value, which C1 rules out: conformance is "+
			"already held by the loader and is not re-established here\n"+
			"stdout: %s", err, stdout)
	}

	// NO THIRD ENCODING: the scalar arm renders `members[0]` verbatim. The
	// two defect spellings are a wrapped array (the set arm misapplied) and
	// an omitted key (the empty value collapsed to absence).
	body := string(mustReadArtifact(t, art))
	if strings.Contains(body, `"note":[]`) || strings.Contains(body, `"note":[""]`) {
		t.Errorf("the artifact encodes the empty scalar as an ARRAY; that is "+
			"a THIRD encoding — the scalar arm renders `members[0]` "+
			"verbatim\nartifact: %s", body)
	}
	if !strings.Contains(body, `"note":""`) {
		t.Errorf("the artifact does not carry `\"note\":\"\"`; the empty "+
			"scalar creates NO third encoding and persists as a present key "+
			"with an empty value\nartifact: %s", body)
	}
}
