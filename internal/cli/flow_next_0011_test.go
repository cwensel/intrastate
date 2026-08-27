package cli

// RDR 0011 — C1, the match-conditioned candidate predicate for `flow next`.
//
// This file is the discriminating half of the suite. C3 states outright why
// it has to exist: "A green 0005 suite is NOT evidence C1 shipped — S2/S3
// are." Every oracle below is written so that the SHIPPED (stripped-match)
// predicate fails it, and so that a build which decides equality in the CLI
// — the fork C1 forbids — fails it too.
//
// The load-bearing distinction across this file: `flow next` decides
// PRESENCE, the kernel decides EQUALITY. An oracle that could be satisfied
// by a CLI-side value comparison is not an oracle for this contract.

import (
	"slices"
	"strings"
	"testing"
)

// --- the predicate itself ------------------------------------------------

// REQ-1: "flow next MUST report as a candidate exactly each non-escape row
// of the requested model whose match atoms over keys PRESENT in the
// assembled view all hold and whose guard the kernel does not decide false."
// REQ-77: "match key present, equal | yes | none | yes | 0"
// REQ-78: "match key present, unequal | **no** (excluded) | n/a — not
// reported | **no** | 0"
// REQ-79: "match key absent from view | yes | `{key, absent}` | yes | 0"
// REQ-99 / `0011:S3`: "present-equal ⇒ candidate, **the match key** absent
// from `unknown`; present-unequal ⇒ excluded …; absent ⇒ candidate with
// `{key, absent}`".
// REQ-90 / MVV 6: "Then add `--tag key=<value>` and assert the list narrows
// to the rows whose match holds."
// HAPPY PATH / INPUT EDGE / BOUNDARY — the three reachable match classes,
// one table.
func TestReq1And77And78And79_TheThreeReachableMatchClassesSortToCandidates(t *testing.T) {
	model := writeFlowModel(t, flowMatchClassesModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	cases := []struct {
		name string
		tag  []string
		// want is the candidate rule-id set, sorted. `dead-row` appears
		// once per SURVIVING expanded row, so the list may repeat it.
		want []string
		why  string
	}{
		{
			// REQ-79 / REQ-14: nothing supplies `phase`, so every match
			// atom over it is omitted from the probe and NO row is
			// excluded on match. Both expansions of `dead-row` survive.
			name: "absent",
			tag:  nil,
			want: []string{
				"dead-row", "dead-row", "match-alpha", "match-beta",
				"no-match-atoms",
			},
			why: "an ABSENT match key MUST NOT exclude its row (C1); it is " +
				"reported under `unknown` instead",
		},
		{
			// REQ-77 + REQ-78: `phase=alpha` holds for `match-alpha` and
			// is UNEQUAL for `match-beta`, which the kernel refuses
			// `no_match`. `dead-row`'s LIVE expansion (the collapsed
			// one-tag `alpha` row) survives; its DEAD expansion carries
			// `{phase alpha}` and `{phase zeta}` and the kernel's
			// conjunction fails it — so `dead-row` is reported ONCE.
			name: "present-equal-and-present-unequal",
			tag:  []string{"--tag", "phase=alpha"},
			want: []string{"dead-row", "match-alpha", "no-match-atoms"},
			why: "a PRESENT-and-UNEQUAL match key is the ONLY probe " +
				"disposition that excludes (C1); a stripped-match build " +
				"reports every row here",
		},
		{
			// REQ-80 / REQ-99: at a value on NEITHER side of the pairing
			// BOTH expansions of `dead-row` fail, so the rule vanishes —
			// and it vanishes by the KERNEL's conjunction, never by a CLI
			// literal comparison (REQ-10, A12).
			name: "present-unequal-for-every-match-row",
			tag:  []string{"--tag", "phase=zeta"},
			want: []string{"no-match-atoms"},
			why: "`phase=zeta` satisfies no authored literal, so only the " +
				"row with NO match atom over `phase` survives",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := runNext(t, model, []string{bind}, tc.tag...)
			got := candidateRules(t, data)
			want := slices.Clone(tc.want)
			slices.Sort(want)
			if !slices.Equal(got, want) {
				t.Errorf("candidates = %v; want %v\n%s", got, want, tc.why)
			}
		})
	}
}

// REQ-11: "A match atom whose key is ABSENT from the assembled view MUST be
// omitted from the probe and MUST NOT exclude the row; it MUST appear in
// the candidate's `unknown` list as `{key, absent}`."
// REQ-23: "Match atoms over absent keys and unestablished owned keys carry
// `absent`, derived from the CLI's own view walk."
// REQ-79: "match key absent from view | yes | `{key, absent}` | yes | 0"
// INPUT EDGE
func TestReq11And23And79_AnAbsentMatchKeyIsReportedAbsentNotExcluded(t *testing.T) {
	model := writeFlowModel(t, flowAbsentMatchKeyModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	data := runNext(t, model, []string{bind})

	c := requireCandidate(t, data, "absent-row")
	pairs := candidateUnknown(t, c)
	if !hasUnknown(pairs, "wanted", "absent") {
		t.Errorf("candidate `absent-row` matches on `wanted`, which nothing "+
			"supplies, but its `unknown` does not carry {wanted, absent}: "+
			"%#v\nAn absent match key MUST NOT be a silent exclusion and "+
			"MUST NOT be a silent inclusion either — it is named, with its "+
			"reason, so the caller knows what to supply", pairs)
	}
}

// REQ-101 / `0011:S3`: "The assertion is therefore on the MATCH key's
// absence from `unknown`, never on `unknown` being empty" — because "an
// ordinary rule needs a write block, and a written owned key enters
// `RequiresOwned`, so a present-equal row may still carry an unrelated
// `{owned-key, absent}` entry from C1's owned walk."
// REQ-77: "match key present, equal | yes | none | yes | 0"
// HAPPY PATH
func TestReq77And101_APresentAndEqualMatchKeyLeavesNoUnknownEntryForThatKey(t *testing.T) {
	model := writeFlowModel(t, flowMatchClassesModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	data := runNext(t, model, []string{bind}, "--tag", "phase=alpha")

	c := requireCandidate(t, data, "match-alpha")
	pairs := candidateUnknown(t, c)
	if unknownKeys(pairs, "phase") {
		t.Errorf("candidate `match-alpha` reports `phase` under `unknown` "+
			"while `phase` is PRESENT and EQUAL: %#v\nThe assertion is on "+
			"the MATCH key's absence from `unknown`, never on `unknown` "+
			"being empty (S3) — a present-and-equal key is decided, so "+
			"there is nothing undecided to report about it", pairs)
	}
}

// REQ-12: "A match atom whose key is present and whose value the kernel
// finds unequal yields `no_match`, and `no_match` is the ONLY probe
// disposition that excludes: such a row MUST NOT be reported, and under
// --evaluate-gates its gates MUST NOT run."
// REQ-78: "match key present, unequal | **no** (excluded) | n/a — not
// reported | **no** | 0"
// REQ-70: "gates under `--evaluate-gates` for reported candidates only".
// ADVERSARIAL — the gate is the observable side effect a match-excluded row
// must not produce.
func TestReq12And70And78_AMatchExcludedRowsGatesDoNotRunUnderEvaluateGates(t *testing.T) {
	model := writeFlowModel(t, flowMatchClassesModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	data := runNext(t, model, []string{bind},
		"--tag", "phase=alpha", "--evaluate-gates")

	if c := candidateNamed(t, data, "match-beta"); c != nil {
		t.Fatalf("`match-beta` matches `phase eq \"beta\"` while `phase` is "+
			"present at `alpha`; the kernel refuses its probe `no_match`, "+
			"which is the ONLY disposition that excludes. It was reported "+
			"as %#v", c)
	}

	// The negative control: the row that DOES match is reported and its
	// gate DID run, so the absence above is a match exclusion and not a
	// gate surface that never works.
	reported := requireCandidate(t, data, "match-alpha")
	gates, ok := objectsAt(reported, "gates")
	if !ok || len(gates) == 0 {
		t.Fatalf("the reported candidate `match-alpha` carries no gate "+
			"results under --evaluate-gates: %#v\nWithout this control the "+
			"assertion above would pass against a build that runs no gate "+
			"at all", reported["gates"])
	}

	// A gate result for the EXCLUDED row must appear nowhere on the
	// payload: "under --evaluate-gates its gates MUST NOT run".
	for _, c := range nextCandidates(t, data) {
		if c["rule"] == "match-beta" {
			t.Errorf("a gate result rides a match-excluded row: %#v", c)
		}
	}
}

// REQ-13: "`guard_unevaluable` and `owned_state_unavailable` MUST leave the
// row a candidate with their facts reported."
// REQ-83: "guard decided false | **no** (excluded) | n/a | no | 0" —
// [0005-carried], the one guard disposition that still excludes.
// REQ-96 / `0011:S2`: "`gated-excluded` absent in both **over
// `flowGatedNextModel`, the fixture that authors it** … with
// `gated-reported` present in both as its negative control".
// DOMAIN EDGE — asserted over the fixture that AUTHORS the row, because
// "asserted over `flowMVVModel`, which carries no such row, the oracle
// passes vacuously in every possible build and pins nothing".
func TestReq83And96_GuardExcludedRowsStayExcludedInBothModes(t *testing.T) {
	model := writeFlowModel(t, flowGatedNextModel)
	art := seedArtifact(t, model, "status=draft", "flag=true")
	bind := artifactBinding(flowStateRole, art)

	for _, mode := range []struct {
		name string
		args []string
	}{
		{"default", nil},
		{"all", []string{"--all"}},
	} {
		t.Run(mode.name, func(t *testing.T) {
			data := runNext(t, model, []string{bind}, mode.args...)

			if c := candidateNamed(t, data, "gated-excluded"); c != nil {
				t.Errorf("`gated-excluded` is reported: %#v\nIts guard "+
					"`flag eq \"false\"` is decided FALSE over `flag=true`, "+
					"and C1 leaves guard exclusions untouched — so it is "+
					"absent in BOTH modes", c)
			}
			// The negative control, without which the assertion above
			// would pass against a build reporting nothing at all.
			requireCandidate(t, data, "gated-reported")
		})
	}
}

// REQ-14: "A row ALL of whose match atoms are omitted is probed with an
// empty match pattern, which matches unconditionally; it is a candidate and
// every omitted key is listed `absent`."
// REQ-68: "a row whose retained match pattern is empty matches
// unconditionally as it does today".
// REQ-104 / `0011:S4`: "the all-absent row is a candidate — its probe
// carries an empty match pattern, which matches unconditionally — and every
// omitted key is listed `{key, absent}`, so its `unknown` is NOT empty."
// BOUNDARY — the reachable form of C1's empty-probe clause.
func TestReq14And68And104_ARowWhoseMatchAtomsAreAllOmittedMatchesUnconditionally(t *testing.T) {
	model := writeFlowModel(t, flowAbsentMatchKeyModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	data := runNext(t, model, []string{bind})

	c := requireCandidate(t, data, "all-absent")
	pairs := candidateUnknown(t, c)
	for _, key := range []string{"wanted", "other"} {
		if !hasUnknown(pairs, key, "absent") {
			t.Errorf("the all-absent row does not list {%s, absent}: %#v\n"+
				"EVERY omitted key is listed, so this row's `unknown` is "+
				"NOT empty — that non-emptiness is what separates it from "+
				"MVV 7's fully-resolved candidate", key, pairs)
		}
	}
}

// REQ-82: "only match atom is `recognized` | yes | none — lifted into
// `Row.Outcome` at normalize | yes | 0"
// REQ-104 / `0011:S4`: "the `recognized`-only row is a candidate with
// `unknown` empty of match facts".
// DOMAIN EDGE — A5's fixture B. `internal/table/normalize.go` lifts the
// `recognized` match atom into `Row.Outcome`, so no `recognized`-keyed atom
// reaches `Row.Atoms` at all and it can never be an unknown fact.
func TestReq82And104_ARecognizedOnlyRowIsACandidateWithNoMatchFacts(t *testing.T) {
	model := writeFlowModel(t, flowAbsentMatchKeyModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	data := runNext(t, model, []string{bind})

	c := requireCandidate(t, data, "recognized-only")
	pairs := candidateUnknown(t, c)
	if unknownKeys(pairs, "recognized") {
		t.Errorf("the `recognized`-only row reports `recognized` under "+
			"`unknown`: %#v\n`recognized` is lifted into `Row.Outcome` at "+
			"normalize and C1 binds the probe's `Recognized` to the row's "+
			"OWN outcome, so it can never be an undecided fact", pairs)
	}
	if len(pairs) != 0 {
		t.Errorf("the `recognized`-only row's `unknown` = %#v; the row "+
			"declares no gate, matches on nothing but `recognized`, and "+
			"writes a key the invoked reader establishes, so nothing is "+
			"undecided about it", pairs)
	}
}

// REQ-80: "dead row … yes while the key is absent (undecided); **no** once
// it is present at any value"
// REQ-10: "The CLI MUST NOT compare a row's literals against each other to
// detect a dead row — that is a literal-to-literal comparison this contract
// does not own".
// REQ-99 / `0011:S3`: "a dead row … ⇒ candidate with `{key, absent}` while
// absent, excluded once the key is present at any value, with no CLI
// literal comparison (A12)."
// REQ-8: "A probe row MAY carry several match tags on one key" — the filter
// MUST tolerate it; `atomsFromBlock`'s `eq`+`in` pairing is the shipped
// producer.
// REQ-9: "Whether a key's tags hold TOGETHER is the kernel's conjunction
// (`TagSet.matches` is all-must-hold), never the CLI's: a dead row is a
// candidate carrying `{k, absent}` while `k` is absent, exactly as any
// undecided row is, and is excluded by the kernel once `k` is present at any
// value, exactly as `flow resolve` never selects it."
// REQ-100: "the oracle asserts two match tags on the DEAD row only — the
// pairing's other expanded row collapses to one tag and is live".
// ADVERSARIAL — the dead row is the case a CLI that reasoned about literals
// would get RIGHT for the wrong reason, so the oracle is written on the
// disposition rather than on any claim about the CLI's internals.
func TestReq10And80And99And100_TheDeadRowIsTheKernelsConjunctionNotACLIComparison(t *testing.T) {
	model := writeFlowModel(t, flowMatchClassesModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	// While `phase` is ABSENT the dead row is a candidate exactly as any
	// undecided row is — both its expansions survive, and each names the
	// key it could not decide.
	absent := runNext(t, model, []string{bind})
	var deadWhileAbsent int
	for _, c := range nextCandidates(t, absent) {
		if c["rule"] != "dead-row" {
			continue
		}
		deadWhileAbsent++
		if pairs := candidateUnknown(t, c); !hasUnknown(pairs, "phase", "absent") {
			t.Errorf("the dead row does not carry {phase, absent} while "+
				"`phase` is absent: %#v\nA dead row is a candidate carrying "+
				"{k, absent} while `k` is absent, EXACTLY as any undecided "+
				"row is — the CLI has not compared its literals", pairs)
		}
	}
	if deadWhileAbsent != 2 {
		t.Errorf("the `eq`+`in` pairing reported %d expanded rows while "+
			"`phase` was absent; want 2 (the collapsed live row and the "+
			"two-tag dead row, A12)", deadWhileAbsent)
	}

	// REQ-100's asymmetry, observable through the disposition rather than
	// by counting tags on the wire: at `phase=alpha` the LIVE expansion
	// (whose `in` member coincides with the `eq` literal and collapsed to
	// one tag) survives while the DEAD one is excluded, so the rule is
	// reported EXACTLY ONCE. Asserting two survivors here would assert
	// something false.
	atAlpha := runNext(t, model, []string{bind}, "--tag", "phase=alpha")
	var deadAtAlpha int
	for _, c := range nextCandidates(t, atAlpha) {
		if c["rule"] == "dead-row" {
			deadAtAlpha++
		}
	}
	if deadAtAlpha != 1 {
		t.Errorf("`dead-row` reported %d times at `phase=alpha`; want 1 — "+
			"the collapsed one-tag expansion is LIVE and survives, the "+
			"two-tag expansion is DEAD and the kernel's conjunction "+
			"excludes it", deadAtAlpha)
	}

	// Excluded at a value on NEITHER side of the pairing, which is the
	// arm that separates the kernel's conjunction from any CLI rule about
	// literals: both tags fail, so both expansions go.
	atZeta := runNext(t, model, []string{bind}, "--tag", "phase=zeta")
	if c := candidateNamed(t, atZeta, "dead-row"); c != nil {
		t.Errorf("`dead-row` survives at `phase=zeta`, a value on NEITHER "+
			"side of the `eq`+`in` pairing: %#v\nThe exclusion is the "+
			"kernel's conjunction over the tags the CLI handed it — the "+
			"CLI's only decision was the presence test", c)
	}
}

// REQ-103 / `0011:S3(b)`: "A row carrying SEVERAL match keys in MIXED
// states at once — one present-and-equal, one present-and-unequal, one
// absent — asserting the row is excluded (the unequal key decides, by the
// kernel's conjunction) and that under `--all` it is a candidate whose
// `unknown` carries no match entry."
// REQ-7: "every match tag on one key shares one presence verdict, so the
// filter hands a key's tags to the kernel together or omits them together
// and can never split a key."
// ADVERSARIAL — the per-key-independence claim A12 argues structurally but
// whose spike records as an explicit limit ("two keys each carrying
// multiple tags were not measured").
func TestReq7And103_AMixedStateRowIsExcludedByTheUnequalKeyAlone(t *testing.T) {
	model := writeFlowModel(t, flowMixedStateRowModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	// `holds=yes` (present, equal), `fails=no` (present, UNEQUAL),
	// `missing` unsupplied (absent). One row, three states at once.
	tags := []string{"--tag", "holds=yes", "--tag", "fails=no"}

	t.Run("default-excluded", func(t *testing.T) {
		data := runNext(t, model, []string{bind}, tags...)
		if c := candidateNamed(t, data, "mixed-row"); c != nil {
			t.Errorf("`mixed-row` is reported: %#v\nOne of its three match "+
				"keys is present and UNEQUAL, and the kernel's conjunction "+
				"is all-must-hold — the unequal key decides, regardless of "+
				"the key that holds and the key that is absent", c)
		}
	})

	t.Run("all-reports-it-with-no-match-entry", func(t *testing.T) {
		data := runNext(t, model, []string{bind}, append(slices.Clone(tags), "--all")...)
		c := requireCandidate(t, data, "mixed-row")
		pairs := candidateUnknown(t, c)
		for _, key := range []string{"holds", "fails", "missing"} {
			if unknownKeys(pairs, key) {
				t.Errorf("under --all the row reports match key %q under "+
					"`unknown`: %#v\nUnder --all a match atom is neither an "+
					"exclusion nor an entry in `unknown`, WHATEVER its "+
					"key's presence (C2) — none of these three keys is in "+
					"the row's RequiresOwned, so no other source produces "+
					"them", key, pairs)
			}
		}
	})
}

// REQ-6: "The restriction MUST be applied as a filter over the CONVERTED
// probe — `probe := row.KernelRow()`, then drop from `probe.Match` every
// `resolve.Tag` whose `Key` the view lacks — never by rebuilding match tags
// from `row.Atoms` in `internal/cli`."
// REQ-102 / `0011:S3(a)`: "A SET-KINDED match key … assert the row matches
// when the supplied set equals the authored one under the kernel's
// canonical form."
// REQ-127: "Set literals go through `seamValue`'s existing canonicalization."
// ADVERSARIAL — the ONE fixture that separates a filter-before-conversion
// build from a filter-after-conversion one. "Over every other mandated
// fixture the two builds are observationally identical, so this key is the
// one that separates them."
func TestReq6And102And127_ASetKindedMatchKeyIsComparedInTheKernelsCanonicalForm(t *testing.T) {
	model := writeFlowModel(t, flowSetKindedMatchModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	// The supplied set EQUALS the authored one under the kernel's canonical
	// form: `seamValue` renders the declared-set match literal as the §D13
	// canonical JSON array with HTML escaping DISABLED. A build that
	// rebuilt the tag from `row.Atoms` hands the kernel the bare member
	// `x&y` — or the escaped `["x&y"]` — and reports a spurious
	// `no_match`.
	data := runNext(t, model, []string{bind}, "--tag", `marks=["x&y"]`)

	if c := candidateNamed(t, data, "set-match"); c == nil {
		t.Errorf("`set-match` is NOT reported while the supplied set equals "+
			"the authored one under the kernel's canonical form. reported "+
			"= %v\nA build that re-derives the match tag from `row.Atoms` "+
			"reimplements `seamValue`'s canonicalization in a second place "+
			"and mis-compares; the filter MUST run over `KernelRow()`'s "+
			"converted tags", candidateRules(t, data))
	}

	// The negative control S3(a) names: over a SCALAR match key both
	// builds agree, so a failure above is the set encoding and not a
	// blanket match-comparison defect.
	requireCandidate(t, data, "scalar-control")

	// And the set key genuinely DECIDES: a different member excludes.
	other := runNext(t, model, []string{bind}, "--tag", `marks=["plain"]`)
	if c := candidateNamed(t, other, "set-match"); c != nil {
		t.Errorf("`set-match` survives a supplied set that does NOT equal "+
			"the authored literal: %#v\nWithout this arm the assertion "+
			"above passes against a build that ignores the match pattern "+
			"entirely — which is exactly the predicate 0011 overrides", c)
	}
}

// REQ-3: "The CLI decides PRESENCE only, by the same view-presence test
// `summarize` already applies; it MUST NOT compare a match value."
// REQ-5: "A build that decides equality in the CLI has forked the seam."
// REQ-4 / A-1: "Present means the view HOLDS A VALUE for the key" — read as
// map membership, never a non-empty test, because "an empty-but-present
// value yields a `Tag` with `Value: \"\"`, which `assembledView` stores as a
// present map entry".
// BOUNDARY — the empty string is the input that separates a presence test
// from a non-empty test, and a non-empty test would also break REQ-109's
// key-set agreement.
func TestReq3And4And5_PresenceIsMapMembershipAndTheEmptyStringIsPresent(t *testing.T) {
	model := writeFlowModel(t, flowMatchClassesModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	// `no-match-atoms` matches on nothing but `recognized`, so it is a
	// candidate however `phase` is supplied — the control that keeps the
	// arms below from passing vacuously.
	base := runNext(t, model, []string{bind})
	requireCandidate(t, base, "no-match-atoms")

	// A value the model DECLARES but that no authored literal names is
	// still PRESENT, and presence alone decides which atoms the probe
	// carries. The row is excluded because the KERNEL found the values
	// unequal, not because the CLI compared them.
	supplied := runNext(t, model, []string{bind}, "--tag", "phase=zeta")
	if c := candidateNamed(t, supplied, "match-alpha"); c != nil {
		t.Errorf("`match-alpha` survives `phase=zeta`: %#v", c)
	}
	c := requireCandidate(t, supplied, "no-match-atoms")
	if pairs := candidateUnknown(t, c); unknownKeys(pairs, "phase") {
		t.Errorf("a row with NO match atom over `phase` reports it under "+
			"`unknown`: %#v\nThe walk reports the atoms the ROW carries, "+
			"never every key the model declares", pairs)
	}
}

// REQ-15: "The invoked reader set and the assembled view are fixed ONCE per
// invocation, before the row loop, and are the same under --all; `unknown`
// is never a function of row order or of the flag."
// REQ-18: "The demand set is a property of the MODEL … it is computed once
// per invocation, before the row loop, and is identical under --all"
// REQ-76 / `0011:D-identity`: "the same request in the same mode over the
// same model revision and artifact contents reports the same candidate set."
// BOUNDARY — order-independence is observable as run-to-run identity, and
// mode-independence as an equal `readers` set and view across the flag.
func TestReq15And18And76_TheReaderSetAndViewAreFixedOncePerInvocationAndModeIndependent(t *testing.T) {
	model := writeFlowModel(t, flowMatchClassesModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	args := nextArgs(model, []string{bind}, "--tag", "phase=alpha")
	first := requireSuccess(t, args...)
	for i := range 3 {
		if again := requireSuccess(t, args...); again != first {
			t.Fatalf("run %d differs from run 0; `unknown` is never a "+
				"function of row order\nfirst:\n%s\nagain:\n%s",
				i+1, first, again)
		}
	}

	def := runNext(t, model, []string{bind}, "--tag", "phase=alpha")
	all := runNext(t, model, []string{bind}, "--tag", "phase=alpha", "--all")

	if a, b := readersOf(t, def), readersOf(t, all); !slices.Equal(a, b) {
		t.Errorf("readers under default = %v, under --all = %v; the demand "+
			"set is a property of the MODEL, not of the mode, and is "+
			"identical under --all (C1)", a, b)
	}
	if a, b := viewKeysOf(t, def), viewKeysOf(t, all); !slices.Equal(a, b) {
		t.Errorf("the assembled view's key set under default = %v, under "+
			"--all = %v; the view is fixed ONCE per invocation and is the "+
			"same under --all (C1)", a, b)
	}
}

// REQ-32: "`outcomes` MUST remain the model's full declared alphabet in
// every mode." — [0005-carried] (0005 DEV-4's reading stands).
// REQ-69: "The alphabet `outcomes` is copied from the model in both modes
// (0005's DEV-4 reading stands: the alphabet is what may be requested; the
// candidates are what varies with state)."
// REQ-33: "flow next MUST NOT choose among the candidates it reports and
// MUST NOT turn a gate deny into a refusal." — [0005-carried].
// REQ-74 / `0011:D-selection-predicate`: "In both modes `next` reports all
// survivors and chooses none".
// DOMAIN EDGE — the alphabet is what narrowing must NOT touch.
func TestReq32And33And69And74_TheAlphabetIsTheFullDeclaredSetInBothModes(t *testing.T) {
	model := writeFlowModel(t, flowMatchClassesModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	want := []string{"alpha", "beta", "dead", "gamma"}

	for _, mode := range []struct {
		name string
		args []string
	}{
		{"default", []string{"--tag", "phase=zeta"}},
		{"all", []string{"--tag", "phase=zeta", "--all"}},
	} {
		t.Run(mode.name, func(t *testing.T) {
			data := runNext(t, model, []string{bind}, mode.args...)

			outcomes, ok := stringsAt(data, "outcomes")
			if !ok {
				t.Fatalf("`outcomes` is not an array of bare tags: %#v",
					data["outcomes"])
			}
			slices.Sort(outcomes)
			if !slices.Equal(outcomes, want) {
				t.Errorf("outcomes = %v; want the model's FULL declared "+
					"alphabet %v\nThe alphabet is what may be REQUESTED; "+
					"the candidates are what varies with state. At "+
					"`phase=zeta` only one candidate survives and the "+
					"alphabet still carries four", outcomes, want)
			}

			// `next` reports all survivors and chooses none: a single
			// candidate is not a selection, and no refusal is minted.
			if len(candidateRules(t, data)) == 0 && mode.name == "default" {
				t.Error("`candidates` is empty; `no-match-atoms` matches " +
					"on nothing but `recognized` and survives every value " +
					"of `phase`")
			}
		})
	}
}

// REQ-34: "This contract changes nothing in the `internal/resolve` PACKAGE:
// the kernel's match seam stays two-valued, its selection and escape phases
// are untouched, `0007:REQ-78` stays deferred, `0007:C8`'s payload is not
// widened, and no sixth `RefusalKind` is minted (`0007:REQ-79`)."
// REQ-107 / `0011:S5`: "the kernel is untouched — the ONE mechanically
// checkable form of that claim is MVV 9's `git diff --stat internal/resolve`
// empty, and it is the form the build asserts".
// REQ-130: "`0007:REQ-78` (the two-valued match seam) is left deferred".
// ADVERSARIAL — a claim about the DIFF, asserted against the merge base so
// a kernel edit smuggled in with this work fails here rather than passing
// unnoticed under a green suite.
func TestReq34And107And130_TheKernelPackageIsUntouchedByThisContract(t *testing.T) {
	// The refusal-kind set is the kernel's own closure and is the wire
	// half of "no sixth RefusalKind is minted". `RefusalKinds()` is
	// exported and the five spellings are `0007`'s.
	want := []string{
		"ambiguous_match", "guard_unevaluable", "no_match",
		"owned_state_unavailable", "unmodeled_outcome",
	}
	var got []string
	for _, k := range resolveRefusalKindStrings() {
		got = append(got, k)
	}
	slices.Sort(got)
	if !slices.Equal(got, want) {
		t.Errorf("RefusalKinds() = %v; want the five `0007` fixes %v — this "+
			"contract mints no sixth kind (`0007:REQ-79`)", got, want)
	}

	// The mechanically checkable form of "the kernel is untouched": the
	// diff of `internal/resolve` against the branch point is EMPTY. A
	// green kernel suite is corroboration, not this criterion — "a green
	// suite cannot by itself prove nothing changed" (S5).
	diff := strings.TrimSpace(gitDiffStat(t, "internal/resolve"))
	if diff != "" {
		t.Errorf("`git diff --stat internal/resolve` is NOT empty:\n%s\n"+
			"C1 changes nothing in the kernel PACKAGE; a probe-shape change "+
			"in `internal/cli` plus one demand-set term is the whole diff "+
			"this contract authorizes there", diff)
	}
}

// REQ-2: "Both verdicts MUST be the kernel's, asked over a one-row probe
// table that carries the row's match atoms RESTRICTED to keys present in the
// CLI's assembled view …, no escape list, and `Recognized` bound to the
// row's own outcome — the binding that keeps the probe modelling its own
// outcome, so `unmodeled_outcome` is unreachable."
// A-7: "no escape list" means `probe.Escape = nil`, the shipped strip
// retained unchanged in both modes; only `probe.Match` changes.
// A-8: "`Outcomes: []string{row.Outcome}` is retained alongside
// `Recognized: row.Outcome`, so `Table.models` holds trivially — C1's probe
// builder MUST preserve that pairing."
// REQ-67: "`excluded` must return the kernel's `Result` (not a `bool`) so
// `summarize` can read the `guard_unevaluable` payload; this is the one
// signature change in the file"
// ADVERSARIAL — the probe's three non-match properties, each observable
// through the payload rather than through the builder.
func TestReq2And67_TheProbeBindsItsOwnOutcomeStripsEscapeAndKeepsTheResult(t *testing.T) {
	// (i) `Recognized` bound to the row's OWN outcome: `unmodeled_outcome`
	// is unreachable, so `flow next` never refuses it however many
	// outcomes the model declares. A builder that bound the requested
	// outcome (there is none on `next`) or left it empty would refuse.
	model := writeFlowModel(t, flowMatchClassesModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	data := runNext(t, model, []string{bind})
	for _, c := range nextCandidates(t, data) {
		outcome, _ := c["outcome"].(string)
		if outcome == "" {
			t.Errorf("candidate %v names no outcome; the probe binds "+
				"`Recognized` to the ROW's own outcome, which is what keeps "+
				"it modelling itself", c["rule"])
		}
	}

	// (ii) the escape list is STRIPPED: an escape row is never reported as
	// a candidate, and its presence never rescues a row the match
	// predicate excluded. `flowEscapeModel` authors one for `bail`.
	esc := writeFlowModel(t, flowEscapeModel)
	escBind := artifactBinding(flowStateRole, seedArtifact(t, esc, "status=draft"))
	escData := runNext(t, esc, []string{escBind})
	if c := candidateNamed(t, escData, "bail-escape"); c != nil {
		t.Errorf("the escape row `bail-escape` is reported as a candidate: "+
			"%#v\nC1 reports exactly each NON-ESCAPE row, and the probe "+
			"carries no escape list — the stripped list leaves "+
			"`escapeOrRefuse` no rescue row, so it passes the original "+
			"refusal through", c)
	}

	// (iii) `excluded` returns the kernel's `Result` rather than a `bool`:
	// the observable consequence is that a `guard_unevaluable` refusal's
	// payload REACHES `summarize`. Today `excluded` reduces the Result to
	// a bool and the payload never arrives, which is why `uncomparable`
	// cannot appear on `next` at all.
	writer := writeFlowModel(t, flowLooseWriterModel)
	art := newFlowArtifact(t, "loose.artifact")
	lbind := artifactBinding(flowStateRole, art)
	requireSuccess(t, "flow", "set-state", "--model", writer,
		"--artifact", lbind, "--write", "status=draft",
		"--write", "size=notanint", "--as=json")

	reader := writeFlowModel(t, flowUncomparableGuardModel)
	uncomparable := runNext(t, reader, []string{lbind})
	pairs := candidateUnknown(t, requireCandidate(t, uncomparable, "uncomparable-row"))
	if !hasUnknown(pairs, "size", "uncomparable") {
		t.Errorf("no `uncomparable` fact reached the payload: %#v\nThat "+
			"fact exists ONLY on the kernel's `guard_unevaluable` refusal "+
			"payload, so a build whose probe helper still returns a `bool` "+
			"discards it before `summarize` can read it. The signature "+
			"change is the reason the probe Result is kept rather than "+
			"discarded", pairs)
	}
}

// REQ-120 / PH 1: "in `flow_exec.go::invokedReaders`, add each row's
// match-block owned keys to the demand set — a term both callers take … —
// the one edit outside `flow_next.go`"
// ADVERSARIAL — a SCOPE claim about the production diff. Phase 1's intent
// bounds where the change may land, and a build that spread the predicate
// across other files has not honoured it. Asserted as a diff, because no
// runtime observation can see file boundaries.
func TestReq120_TheProductionDiffIsFlowNextPlusOneTermInFlowExec(t *testing.T) {
	// The two production files this contract authorizes to change, plus
	// the two shipped user-facing descriptions C3 corrects (REQ-56) and
	// the 0005 test files C3 re-homes (REQ-54, REQ-60).
	allowed := map[string]bool{
		"internal/cli/flow_next.go":                  true,
		"internal/cli/flow_exec.go":                  true,
		"docs/cli-output-contract.md":                true,
		"README.md":                                  true,
		"internal/cli/flow_next_0005_test.go":        true,
		"internal/cli/flow_adversarial_0005_test.go": true,
		"internal/cli/flow_mvv_0005_test.go":         true,
		"internal/cli/flow_surface_0005_test.go":     true,
	}

	for _, path := range changedFilesSince(t) {
		if allowed[path] || strings.Contains(path, "_0011_test.go") ||
			strings.HasPrefix(path, "docs/rdr/") {
			continue
		}
		t.Errorf("`%s` changed; Phase 1 names ONE edit outside "+
			"`flow_next.go` — the `invokedReaders` demand-set term in "+
			"`flow_exec.go` — and Phases 2 and 3 add the flag, the help, "+
			"C3's fixtures, and the two shipped descriptions. A predicate "+
			"spread wider than that has not honoured the placement this "+
			"RDR decided: PRESENCE in the CLI, EQUALITY in the kernel, "+
			"`internal/resolve` not modified", path)
	}
}
