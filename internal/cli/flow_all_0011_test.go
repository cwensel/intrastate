package cli

// RDR 0011 — C2, the `--all` flag: the same contract's second surface.
//
// `--all` restores 0005's ENUMERATION — the match pattern taking no part in
// the verdict — for the caller who wants the guard-conditioned alphabet
// rather than the state-conditioned one. Two things about it are easy to
// get subtly wrong, and both have their own oracle here:
//
//  1. `--all` is NOT output-identical to 0005. Under `--all` a match atom
//     is neither an exclusion nor an entry in `unknown` — a departure from
//     shipped behaviour, since `summarize` today walks `Row.Atoms` with no
//     `Block` test and a match key absent from the view lands in the list.
//  2. The filter is scoped to the ATOM WALK's contribution ONLY, at
//     EMISSION. A post-merge filter keyed on "is this key a match-block
//     atom's key" deletes the OWNED-key walk's identical `{key, absent}`
//     pair too, because C1 dedups on the pair and the two have no surviving
//     provenance to discriminate on.

import (
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
)

// REQ-36: "flow next MUST accept a boolean flag --all, default false."
// REQ-50 / `0011:D-naming`: "the flag is `--all`, boolean." — the exact
// spelling is normative; `--enumerate`, `--ignore-match`, `--no-match`, a
// sense flip, and an enum (`--state=all|applicable`) are named REJECTED.
// A-9: the flag has NO shorthand — `0011:A14` records that `-a` produces
// different pflag text, and registering it would widen the surface REQ-47
// pins.
// HAPPY PATH
func TestReq36And50_NextRegistersABooleanAllFlagDefaultFalseWithNoShorthand(t *testing.T) {
	// Vacuously true of an unregistered verb, so `next` must exist before
	// anything about its flags means anything.
	if !flowSubcommand(t, "next") {
		t.Fatal("`flow next` is not registered; the flag surface this " +
			"contract fixes is not yet assertable")
	}

	var checked bool
	for _, c := range NewRootCmd().Commands() {
		if c.Name() != "flow" {
			continue
		}
		for _, sub := range c.Commands() {
			if sub.Name() != "next" {
				continue
			}
			checked = true

			f := sub.Flags().Lookup("all")
			if f == nil {
				t.Fatalf("`flow next` does not register --all; the spelling "+
					"is normative (`0011:D-naming`), and %v are named "+
					"REJECTED alternatives",
					[]string{"--enumerate", "--ignore-match", "--no-match"})
			}
			if f.Value.Type() != "bool" {
				t.Errorf("--all is of type %q; want `bool` — the axis has "+
					"exactly two members, which is why an enum "+
					"(`--state=all|applicable`) was rejected", f.Value.Type())
			}
			if f.DefValue != "false" {
				t.Errorf("--all defaults to %q; want \"false\" — the DEFECT "+
					"is the default, so flipping the sense (`--select` on a "+
					"0005 default) is the rejected alternative", f.DefValue)
			}
			if f.Shorthand != "" {
				t.Errorf("--all registers the shorthand %q; A14 observed "+
					"`-a` yielding different pflag text and A-9 reads the "+
					"spelling as the long form only", f.Shorthand)
			}

			// The rejected spellings must NOT ship as aliases.
			for _, rejected := range []string{
				"enumerate", "ignore-match", "no-match", "select", "state",
			} {
				if sub.Flags().Lookup(rejected) != nil {
					t.Errorf("`flow next` registers --%s; it is a NAMED "+
						"rejected spelling, not an alias", rejected)
				}
			}
		}
	}
	if !checked {
		t.Fatal("the `flow` group registers no `next` verb")
	}
}

// REQ-37: "Under --all the candidate predicate MUST be exactly the one
// 0005:C1 specified: every non-escape row whose guard the kernel does not
// decide false, with the row's match pattern taking no part in the verdict."
// REQ-38: "A match atom is then neither an exclusion nor an entry in
// `unknown`, whatever its key's presence, achieved by the probe omitting the
// match pattern so the kernel never sees those atoms."
// REQ-86: "any of the above, under `--all` | match takes no part; guard
// rules alone decide | no match entry, whatever the key's presence | per
// guard | 0"
// HAPPY PATH — the same three match classes C1 sorts, under `--all`, where
// none of them sorts anything.
func TestReq37And38And86_UnderAllTheMatchPatternTakesNoPartInTheVerdict(t *testing.T) {
	model := writeFlowModel(t, flowMatchClassesModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	// Under `--all` the reported set is the SAME whatever `phase` holds,
	// because match takes no part: every non-escape row survives, the
	// `eq`+`in` pairing's two expansions included.
	want := []string{
		"dead-row", "dead-row", "match-alpha", "match-beta", "no-match-atoms",
	}
	slices.Sort(want)

	for _, tag := range [][]string{
		nil,
		{"--tag", "phase=alpha"},
		{"--tag", "phase=zeta"},
	} {
		name := "absent"
		if tag != nil {
			name = tag[1]
		}
		t.Run(name, func(t *testing.T) {
			data := runNext(t, model, []string{bind},
				append(slices.Clone(tag), "--all")...)

			if got := candidateRules(t, data); !slices.Equal(got, want) {
				t.Errorf("candidates under --all = %v; want %v\nUnder --all "+
					"the predicate is EXACTLY 0005's: every non-escape row "+
					"whose guard the kernel does not decide false, with the "+
					"match pattern taking no part — so the set cannot vary "+
					"with `phase`", got, want)
			}
		})
	}
}

// REQ-39: "Under --all that entry MUST disappear — reporting a fact the
// mode ignores would be reporting on a predicate it does not apply — so the
// --all branch MUST filter BlockMatch atoms out of `unknown`, and an oracle
// MUST assert their absence."
// REQ-42: "The oracle asserting the absence MUST therefore be written over
// a match key that is NOT in the row's `RequiresOwned`" — "asserted over a
// matched-and-written key it fails on the very model MVV 4 runs, and
// forcing it to pass would delete a shipped 0005 fact."
// REQ-62: "That absent-key row MUST match on a key it does NOT write or
// clear, writing some OTHER key instead" — "a matched-AND-written key enters
// `RequiresOwned`, where C1's owned-key walk produces the same `{key,
// absent}` pair that C2's --all filter does not touch, so the row would
// report the pair in both modes and discriminate nothing."
// A-4: the filter keys on the ATOM's `Block` field, not on a key-set
// difference.
// ADVERSARIAL — this is the ONE respect in which `--all` is not
// output-identical to 0005, and it lies outside A4's predicate-scoped
// classification, so a green 0005 suite says nothing about it.
func TestReq39And42_UnderAllAMatchAtomsUnknownEntryDisappears(t *testing.T) {
	model := writeFlowModel(t, flowAbsentMatchKeyModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	// `absent-row` matches on `wanted` and WRITES `status`, so `wanted` is
	// NOT in its `RequiresOwned` — REQ-42's requirement, and what keeps
	// this oracle from being confounded by the owned-key walk.
	def := runNext(t, model, []string{bind})
	defPairs := candidateUnknown(t, requireCandidate(t, def, "absent-row"))
	if !hasUnknown(defPairs, "wanted", "absent") {
		t.Fatalf("by DEFAULT the row does not report {wanted, absent}: "+
			"%#v\nWithout this control the --all assertion below passes "+
			"against a build that never reports the fact in either mode",
			defPairs)
	}

	all := runNext(t, model, []string{bind}, "--all")
	allPairs := candidateUnknown(t, requireCandidate(t, all, "absent-row"))
	if unknownKeys(allPairs, "wanted") {
		t.Errorf("under --all the row still reports `wanted` under "+
			"`unknown`: %#v\nThe mode IGNORES the match predicate, so "+
			"reporting a match fact would be reporting on a predicate it "+
			"does not apply. This is a DEPARTURE from shipped behaviour, "+
			"not a preservation of it: `summarize` today walks `Row.Atoms` "+
			"with no `Block` test", allPairs)
	}
}

// REQ-40: "The filter is scoped to the ATOM WALK's contribution ONLY. It
// MUST NOT suppress a `{key, absent}` pair the OWNED-key walk independently
// produces for the same key"
// REQ-41: "That scope is only achievable at EMISSION: the filter MUST be
// applied to the atom walk as it contributes, before the sources merge and
// before C1's dedup, never as a predicate over the merged list."
// REQ-30: "Deduplication is the LAST step: each source contributes its
// entries, C2's --all filter is applied to the atom walk's contribution at
// emission, and only then are the merged entries deduped and sorted."
// REQ-97 / `0011:S2`: "over a matched-AND-written key whose reader
// established nothing … `--all` MUST still report `{stage, absent}` from the
// owned-key walk, since only the atom walk's contribution is filtered (C2)
// and dedup runs last (C1). A build that dedups before filtering deletes
// that pair and fails this assertion".
// ADVERSARIAL — the filter/dedup ORDER, which no other oracle in the suite
// discriminates. `models/rdr.toml` is the model where the distinction is
// load-bearing: `stage` is BOTH a `[rule.match]` key and a `[rule.write]`
// key, so it enters `RequiresOwned` and both walks name it.
func TestReq30And40And41And97_TheAllFilterIsScopedToTheAtomWalkAndDedupRunsLast(t *testing.T) {
	model := repoModelPath(t)
	// The artifact is seeded but the reader establishes NOTHING for
	// `stage`: a fresh artifact with no write leaves every owned key
	// absent, so the owned-key walk contributes `{stage, absent}` while
	// the atom walk would too.
	art := newFlowArtifact(t, "unestablished.artifact")
	bind := artifactBinding("rdr", art)

	all := runNext(t, model, []string{bind}, "--all")
	candidates := nextCandidates(t, all)
	if len(candidates) == 0 {
		t.Fatal("`--all` over the checked-in model reported no candidate; " +
			"the guard-conditioned predicate reports every row the guards " +
			"do not exclude")
	}

	var checked bool
	for _, c := range candidates {
		required, ok := stringsAt(c, "required")
		if !ok || !slices.Contains(required, "stage") {
			continue
		}
		checked = true
		pairs := candidateUnknown(t, c)
		if !hasUnknown(pairs, "stage", "absent") {
			t.Errorf("candidate %v names `stage` in `required` — so the "+
				"OWNED-key walk independently produces {stage, absent} — "+
				"but its `unknown` under --all does not carry it: %#v\n"+
				"Only the ATOM walk's contribution is filtered, and dedup "+
				"runs LAST. A build that dedups first cannot honour C2's "+
				"scope: once merged, the two walks' identical pairs are "+
				"indistinguishable and a filter applied afterwards deletes "+
				"BOTH — including the shipped 0005 fact C2 forbids deleting",
				c["rule"], pairs)
		}
	}
	if !checked {
		t.Fatal("no candidate names `stage` in `required`; every rule of " +
			"the checked-in model both matches on and writes `stage`, " +
			"which is what makes this the model where the distinction is " +
			"load-bearing")
	}
}

// REQ-43: "Everything else 0005:C1 requires of flow next — the data
// minimum, gate accessors only under --evaluate-gates and only for reported
// candidates, a deny reported on its candidate with exit 0, no invented
// guard facts, the narrowed invoked reader set — MUST hold identically in
// both modes." — [0005-carried].
// REQ-128: "Every other `0005:C1` obligation on `flow next` is unchanged."
// DOMAIN EDGE — the boundary of the override: only the candidate clause
// moves.
func TestReq43And128_EveryOther0005ObligationHoldsIdenticallyInBothModes(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	for _, mode := range []struct {
		name string
		args []string
	}{
		{"default", nil},
		{"all", []string{"--all"}},
	} {
		t.Run(mode.name, func(t *testing.T) {
			data := runNext(t, model, []string{bind},
				append(slices.Clone(mode.args), "--tag", "profile=mid")...)

			// The data minimum (`0005:REQ-74`), unchanged.
			for _, key := range []string{
				"model", "revision", "observed", "owned", "readers",
				"outcomes", "candidates",
			} {
				if _, ok := data[key]; !ok {
					t.Errorf("payload carries no %q; the data minimum is a "+
						"0005 obligation this contract does not override. "+
						"keys = %v", key, keysOf(data))
				}
			}

			// The NARROWED invoked reader set: `read.orphan` serves only
			// `note`, which no row requires, matches on, or guards — so
			// it must not run in EITHER mode, and its unbound role must
			// not refuse.
			if readers := readersOf(t, data); slices.Contains(readers, "orphan") {
				t.Errorf("readers = %v under %s; `read.orphan` serves only "+
					"`note`, which no row requires, MATCHES on, or guards — "+
					"C1's demand-set term adds match-owned keys and adds "+
					"nothing here", readers, mode.name)
			}

			// No invented guard facts: `stale` is owned and unseeded.
			if owned, ok := data["owned"].(map[string]any); ok {
				if v, held := owned["stale"]; held {
					t.Errorf("`owned[stale]` = %#v under %s though no "+
						"reader established it; `next` invents no fact in "+
						"either mode", v, mode.name)
				}
			}
		})
	}
}

// REQ-44: "The success payload shape MUST be identical in both modes:
// `unknown` is present in both, as `[]` rather than omitted when it carries
// nothing, so one consumer struct parses either mode."
// REQ-71 / `0011:D-undecided-reporting-shape`: "the per-candidate field is
// `unknown`, a list of `{key, reason}`".
// A-2: the wire members are exactly `key` and `reason`, both always
// emitted, on a list named `unknown`.
// BOUNDARY — the empty list is the shape that separates "present as []"
// from "omitted", and one consumer struct must parse either mode.
func TestReq44And71_TheUnknownFieldIsPresentAsAnEmptyListInBothModes(t *testing.T) {
	model := writeFlowModel(t, flowFullyResolvedModel)
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
			data := runNext(t, model, []string{bind},
				append(slices.Clone(mode.args), "--evaluate-gates")...)

			c := requireCandidate(t, data, "resolved-row")

			// The field must be PRESENT — an omitted field and an empty
			// one are different wire shapes, and only one of them lets a
			// single consumer struct parse both modes.
			raw, held := c["unknown"]
			if !held {
				t.Fatalf("candidate carries no `unknown` field at all under "+
					"%s; it is emitted as `[]` rather than OMITTED when it "+
					"carries nothing. candidate keys = %v",
					mode.name, keysOf(c))
			}
			if raw == nil {
				t.Fatalf("`unknown` is JSON null under %s; the shape is an "+
					"empty ARRAY", mode.name)
			}

			pairs := candidateUnknown(t, c)
			if len(pairs) != 0 {
				t.Errorf("`unknown` under %s = %#v; this fixture's row has "+
					"its match and guard keys established by the invoked "+
					"reader, writes owned keys that reader serves, and "+
					"declares no gate — so nothing about it is undecided",
					mode.name, pairs)
			}
		})
	}
}

// REQ-45: "--all is part of request identity." — with `0011:D-identity`.
// REQ-76: "the same request in the same mode over the same model revision
// and artifact contents reports the same candidate set."
// BOUNDARY — identity means the flag CHANGES the answer where the modes
// differ, and repeats it where they do not.
func TestReq45And76_AllIsPartOfRequestIdentity(t *testing.T) {
	model := writeFlowModel(t, flowMatchClassesModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	defArgs := nextArgs(model, []string{bind}, "--tag", "phase=alpha")
	allArgs := nextArgs(model, []string{bind}, "--tag", "phase=alpha", "--all")

	def, all := requireSuccess(t, defArgs...), requireSuccess(t, allArgs...)
	if def == all {
		t.Errorf("the default and --all payloads are byte-identical over a "+
			"model with a present-and-unequal match key; --all is part of "+
			"request identity and this fixture is where the two predicates "+
			"disagree\npayload:\n%s", def)
	}

	// Same request, same mode ⇒ same answer, in BOTH modes.
	for name, args := range map[string][]string{
		"default": defArgs, "all": allArgs,
	} {
		first := requireSuccess(t, args...)
		for i := range 3 {
			if again := requireSuccess(t, args...); again != first {
				t.Fatalf("%s run %d differs from run 0\nfirst:\n%s\n"+
					"again:\n%s", name, i+1, first, again)
			}
		}
	}
}

// REQ-46: "--all MUST NOT be accepted by flow resolve, flow read-state, or
// flow set-state — satisfied by NON-REGISTRATION (registered on next only,
// not on the shared registerSelectionFlags)."
// REQ-47: "The build MUST pin that negative STRUCTURALLY, asserting --all is
// absent from those three verbs' flag sets and from the flow group's own,
// guarded against vacuity by requiring the four verbs to exist first"
// REQ-65: "These are TWO oracles, not one: `TestReq3` is presence-only … so
// it homes the positive (`--all` on `next`) and cannot carry C2's
// structural negative."
// REQ-98 / `0011:S2`: "C2's negative is asserted STRUCTURALLY … not on
// pflag's error text".
// ADVERSARIAL — the structural absence oracle, in the house idiom
// (`TestReq69_NoPlanFlagShipsOnAnyVerb`,
// `TestReq12_FlowGroupDoesNotRedefineTheRootAsFlag`). It holds whatever
// pflag's wording is and it fails if a future persistent `--all` on `flow`
// or root silently widens the surface.
func TestReq46And47And65And98_AllIsAbsentFromTheOtherThreeVerbsAndTheFlowGroup(t *testing.T) {
	// Vacuously true of an unregistered group, so all four verbs must
	// exist before the ABSENCE of `--all` on three of them means anything.
	if names := flowGroupNames(t); len(names) != len(flowVerbs) {
		t.Fatalf("the `flow` group registers %v; all four verbs must exist "+
			"before the ABSENCE of --all on three of them is assertable",
			names)
	}

	var group bool
	for _, c := range NewRootCmd().Commands() {
		if c.Name() != "flow" {
			continue
		}
		group = true

		// The GROUP's own flag set: a persistent --all on `flow` would
		// silently widen the surface on all four verbs at once.
		if c.Flags().Lookup("all") != nil || c.PersistentFlags().Lookup("all") != nil {
			t.Error("the `flow` group registers --all of its own; C2 " +
				"places the flag on `next` ONLY, so a group-level or " +
				"persistent registration widens the surface this clause " +
				"fences")
		}

		for _, sub := range c.Commands() {
			if sub.Name() == "next" {
				continue
			}
			if sub.Flags().Lookup("all") != nil {
				t.Errorf("`flow %s` registers --all; the negative is "+
					"satisfied by NON-REGISTRATION — the flag rides `next` "+
					"only, never the shared selection registrar",
					sub.Name())
			}
		}
	}
	if !group {
		t.Fatal("no `flow` command group is registered")
	}

	// The ROOT's persistent set: `--as` is the sole persistent flag, and a
	// persistent --all there would reach every verb.
	if NewRootCmd().PersistentFlags().Lookup("all") != nil {
		t.Error("the root command registers a PERSISTENT --all; it would " +
			"reach `resolve`, `read-state`, and `set-state` and defeat the " +
			"non-registration this clause rests on")
	}
}

// REQ-48: "`command-error` + exit 2 is the SHARED usage bucket for every
// unknown flag, unknown subcommand, and Args violation, so it MUST NOT be
// asserted as if it named this flag"
// REQ-49: "This contract mints no new refusal code for it."
// REQ-98: "the behavioural run (exit 2, `command-error`, envelope on stdout
// under `--as=json`) corroborates it".
// ADVERSARIAL — the corroborating run, written so it cannot be mistaken for
// the assertion. It asserts the CLASS and the exit, never pflag's free-text
// message, which A14 records as text this repo does not own.
func TestReq48And49_TheUnknownFlagRefusalIsTheSharedUsageClassAndMintsNoNewCode(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	for _, verb := range []struct {
		name string
		args []string
	}{
		{"resolve", []string{"flow", "resolve", "--model", model,
			"--artifact", bind, "--outcome", "hold", "--all", "--as=json"}},
		{"read-state", []string{"flow", "read-state", "--model", model,
			"--artifact", bind, "--all", "--as=json"}},
		{"set-state", []string{"flow", "set-state", "--model", model,
			"--artifact", bind, "--write", "status=draft", "--all",
			"--as=json"}},
	} {
		t.Run(verb.name, func(t *testing.T) {
			stdout, _, err := runCmd(t, verb.args...)
			if err == nil {
				t.Fatalf("`flow %s --all` SUCCEEDED; the flag is registered "+
					"on `next` only\nstdout:\n%s", verb.name, stdout)
			}
			if code := clierr.ErrorCode(err); code != "command-error" {
				t.Errorf("code = %q; want %q — pflag fails the parse before "+
					"any RunE, so no `flow-*` code is reachable and this "+
					"contract mints none", code, "command-error")
			}
			if exit := clierr.ExitCodeFor(err); exit != 2 {
				t.Errorf("exit = %d; want 2", exit)
			}
			// The envelope rides STDOUT under `--as=json`, which
			// `primeAsFlag` is what makes true despite the failed parse.
			if !strings.Contains(stdout, `"command-error"`) {
				t.Errorf("no structured envelope on stdout under --as=json:"+
					"\n%s\nThe emission path is the CLI's OWN gateway, not "+
					"cobra's printer", stdout)
			}
			// The code alone cannot attribute the refusal to THIS flag:
			// it is the shared usage bucket. So no oracle here matches on
			// pflag's text, and this run is corroboration for the
			// structural absence oracle above — never a substitute.
			if strings.Contains(stdout, "flow-all") {
				t.Errorf("a bespoke `flow-all*` code appears on the wire:"+
					"\n%s\nThis contract mints no new refusal code for the "+
					"flag", stdout)
			}
		})
	}
}

// REQ-65: "register --all on next in the per-verb flag map the surface
// suite already carries
// (`flow_surface_0005_test.go::TestReq3_EachVerbRegistersItsNormativeFlagSpellings`).
// These are TWO oracles, not one"
// HAPPY PATH — the POSITIVE half, homed where the surface suite homes every
// other normative flag spelling. `TestReq3` is presence-only and cannot
// carry the structural negative, which is why the two are separate.
func TestReq65_TheSurfaceFlagMapNamesAllOnNext(t *testing.T) {
	// The per-verb map the surface suite carries, with `--all` added to
	// `next`. Asserted here rather than by editing the 0005 map, because
	// C3 requires the 0005 file's assertions to keep passing UNCHANGED
	// where they assert the surface — this is an ADDITION, not a rewrite.
	want := map[string][]string{
		"next":       {"flow", "model", "tag", "artifact", "evaluate-gates", "all"},
		"resolve":    {"flow", "model", "tag", "artifact", "outcome"},
		"read-state": {"flow", "model", "artifact"},
		"set-state":  {"flow", "model", "artifact", "write", "clear"},
	}

	if names := flowGroupNames(t); len(names) != len(flowVerbs) {
		t.Fatalf("the `flow` group registers %v; the four verbs must exist "+
			"before their flag surface is assertable", names)
	}

	for _, c := range NewRootCmd().Commands() {
		if c.Name() != "flow" {
			continue
		}
		for _, sub := range c.Commands() {
			flags, ok := want[sub.Name()]
			if !ok {
				continue
			}
			for _, f := range flags {
				if sub.Flags().Lookup(f) == nil {
					t.Errorf("`flow %s` does not register --%s; the flag "+
						"spellings are normative", sub.Name(), f)
				}
			}
		}
	}
}

// REQ-123 / `0011:F7`: "`next` is effect-free in both modes" —
// [0005-carried], and it holds under `--all` too.
// REQ-126 / `0011:G-cross-cutting`: "No new execution path, no new
// concurrency surface".
// ADVERSARIAL — non-mutation is observable as the artifact's read-back
// being unchanged across a `next` in either mode.
func TestReq123And126_NextIsEffectFreeInBothModes(t *testing.T) {
	model := writeFlowModel(t, flowMatchClassesModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	readBack := func() string {
		t.Helper()
		return requireSuccess(t, "flow", "read-state", "--model", model,
			"--artifact", bind, "--as=json")
	}

	before := readBack()
	for _, args := range [][]string{
		{"--tag", "phase=alpha"},
		{"--tag", "phase=alpha", "--all"},
		{"--tag", "phase=alpha", "--evaluate-gates"},
		{"--all", "--evaluate-gates"},
	} {
		runNext(t, model, []string{bind}, args...)
	}
	if after := readBack(); after != before {
		t.Errorf("the artifact's observed state changed across `flow next` "+
			"runs; `next` is effect-free in BOTH modes and the recovery "+
			"note is \"none needed\"\nbefore:\n%s\nafter:\n%s",
			before, after)
	}
}
