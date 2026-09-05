package cli

// RDR 0019 — the shared fixture corpus and oracles for `flow init-state`.
//
// Three rules hold across the whole 0019 suite, and each is load-bearing
// for whether these tests discriminate the contract rather than merely
// observe an implementation:
//
//  1. STATE IS OBSERVED THROUGH THE CLI, never by parsing the artifact,
//     except where a clause is stated ON ARTIFACT BYTES. C1 scopes its
//     byte claims to the file-backed carrier, and S2/S9/S11/S12 say
//     "asserted on artifact bytes, not exit code alone" — so those tests
//     read bytes deliberately, and every other oracle goes through
//     `flow read-state`.
//
//  2. NO TEST PINS EITHER NEW REFUSAL CODE'S SPELLING. REQ-65 fixes the
//     exit GROUP (2) and DISTINCTNESS only, on the `0028:C1.3`
//     `codeWriteEditRefused` carve-out, and says "no test may pin the
//     string". The oracles below therefore assert group + distinctness,
//     and capture the observed spelling only to compare it against ANOTHER
//     observed spelling.
//
//  3. NO PAYLOAD FIELD NAME C1 DOES NOT FIX IS PINNED. `0019:D-wire-byte-format`
//     defers NAMES and fixes CONTENT, so the payload oracles search the
//     decoded `data` object for a key-name list carrying the expected
//     content rather than demanding a particular field.

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
)

// --- fixture identities --------------------------------------------------

const (
	// initRoleA is the primary artifact role every single-role fixture
	// binds. It reuses the shipped `state` spelling so an invocation reads
	// like the shipped verbs'.
	initRoleA = "state"
	// initRoleB is the SECOND role, bound only by the two-role fixture.
	// S11's ALL-quantifier test needs two artifacts that can differ in
	// emptiness independently, which one role cannot express.
	initRoleB = "aux"
)

// --- the S1-shaped fixture ----------------------------------------------

// initMVVModel is the `0019:MVV` fixture: "a fixture state-machine model
// declaring `[initial]` with one always-present owned key and one plain
// owned key, both writer-served".
//
//   - `[tags.stage]`  owned enum, required (always-present) — the MVV's
//     always-present key;
//   - `[tags.note]`   owned scalar, NOT required — the MVV's plain key,
//     and the `--clear` target of MVV steps 6 and 9;
//   - one `[write.state]`/`[read.state]` pair, both file-backed on role
//     `state`, both `read_back = true`.
//
// It is deliberately NOT a decision table and its `[initial]` names both
// owned keys, so the seeded arm has two keys to report and the no-op arm
// has a non-empty absent list available.
//
// TEST-FIXTURE (Phase 2): the rule carries NO guard on `note`. `0019:MVV`
// requires only "one always-present owned key and one plain owned key, both
// writer-served", and MVV step 1 requires the fixture to LINT CLEAN with
// 0006's arms all green. A guard atom over `note` — an OPTIONAL scalar with
// no finite declared domain — is a blocking `graph-unprovable-coverage`
// under two of 0006's arms at once (`row-can-refuse` and
// `dimension-not-finite`), so a fixture carrying one can never satisfy step 1.
// Dropping the atom costs the MVV nothing: `note` is the CLEAR target of
// steps 6 and 9, and neither reads it through a guard.
const initMVVModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "initmvv"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.stage]
provenance = "owned"
kind = "enum"
domain = ["seeded", "final"]
single_valued = true
required = true

[tags.note]
provenance = "owned"
kind = "scalar"

[read.state]
role = "state"
path = "flow.state"
keys = ["stage", "note"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["stage", "note"]
timeout = "2s"
read_back = true

[initial]
stage = "seeded"
note = "hello"

[context.done]
[context.done.match.stage]
eq = "final"

[[rule]]
id = "advance"
[rule.match.stage]
eq = "seeded"
[rule.match.recognized]
eq = "advance"
[rule.write]
stage = "final"
`

// initTwoRoleModel declares TWO write accessors over TWO roles, each
// serving its own `[initial]` key. It is the ONLY shape that can
// discriminate C1's ALL quantifier from an ANY or per-artifact reading
// (S11): role A's artifact can carry a key while role B's is empty.
const initTwoRoleModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "inittworole"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.stage]
provenance = "owned"
kind = "enum"
domain = ["seeded", "final"]
single_valued = true
required = true

[tags.note]
provenance = "owned"
kind = "scalar"

[read.state]
role = "state"
path = "flow.state"
keys = ["stage"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["stage"]
timeout = "2s"
read_back = true

[read.aux]
role = "aux"
path = "flow.aux"
keys = ["note"]
timeout = "2s"

[write.aux]
role = "aux"
path = "flow.aux"
keys = ["note"]
timeout = "2s"
read_back = true

[initial]
stage = "seeded"
note = "hello"

[context.done]
[context.done.match.stage]
eq = "final"

[[rule]]
id = "advance"
[rule.match.stage]
eq = "seeded"
[rule.match.recognized]
eq = "advance"
[rule.write]
stage = "final"
`

// initDecisionTableModel is a decision-table model the loader ADMITS: no
// `[initial]`, no owned tags, no accessors (S7). That combination is what
// makes S7's discriminating control possible — the class refusal cannot be
// asserted by "a writer that fails if invoked", because such a model can
// bind no owned-tag writer at all, so the ordering is asserted by the
// ABSENCE of any accessor invocation.
const initDecisionTableModel = `outcomes = ["decide"]

[model]
id = "initdt"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.a]
provenance = "observed"
kind = "enum"
domain = ["x", "y"]
single_valued = true
required = true

[[rule]]
id = "cell-x"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "x"

[[rule]]
id = "cell-y"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "y"
`

// initEditCarriedModel declares an EDIT-carried write accessor on the role
// its `[initial]` key needs. Its binding is `*flowbind.EditWriter`, which
// C1's write-side type switch must REFUSE — it has no JSON store, no
// `sealedKey`, and no artifact bytes, so "the store carries no key" is
// UNDEFINED for it (S10 write row 1).
//
// For this carrier the refusal is OVER-DETERMINED — an edit writer with no
// anchor match would also fail on its own terms — so S10's discriminating
// control is that the test must still see the CARRIER code, and that
// removing the carrier check makes the edit case ATTEMPT A WRITE.
const initEditCarriedModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "initedit"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.stage]
provenance = "owned"
kind = "scalar"

[read.state]
role = "state"
path = "flow.state"
keys = ["stage"]
timeout = "2s"

[write.state]
role = "state"
keys = ["stage"]
timeout = "2s"
read_back = true

[write.state.edit.stage]
anchor  = "^- \\*\\*Stage\\*\\*: (.+)$"
replace = "- **Stage**: {stage}"

[initial]
stage = "seeded"

[context.done]
[context.done.match.stage]
eq = "final"

[[rule]]
id = "advance"
[rule.match.stage]
eq = "seeded"
[rule.match.recognized]
eq = "advance"
[rule.write]
stage = "final"
`

// initCommandWriteModel declares a COMMAND-backed write accessor. Its
// binding is `*cmdbind.Writer`, which the write-side type switch must
// refuse (S10 write row 2).
//
// The command is `true`, which SUCCEEDS if ever spawned — deliberately, so
// that a carrier check moved after registry invocation would reach
// `cmdbind::spawn` and the row would fail on the allow-commands refusal
// instead of the carrier code. C1 states the preemption is "a real
// ordering, not a tautology".
const initCommandWriteModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "initcmdwrite"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.stage]
provenance = "owned"
kind = "scalar"

[read.state]
role = "state"
path = "flow.state"
keys = ["stage"]
timeout = "2s"

[write.state]
role = "state"
command = ["true", "{artifact}"]
keys = ["stage"]
timeout = "10s"
read_back = true

[initial]
stage = "seeded"

[context.done]
[context.done.match.stage]
eq = "final"

[[rule]]
id = "advance"
[rule.match.stage]
eq = "seeded"
[rule.match.recognized]
eq = "advance"
[rule.write]
stage = "final"
`

// initCommandReadModel is S10's READ-SIDE row, and it is the row the
// write-side type switch ALONE would ADMIT: `[write.state]` is file-backed
// (`*flowbind.Writer`, admitted) while `[read.state]` on the SAME role is
// command-backed (`cmdbind.Reader`, refused).
//
// It doubles as the detector for the POINTER/VALUE spelling C1 makes
// normative: the registry constructs readers as VALUES (`Reader{...}`), so
// a read-side type-switch case written `*flowbind.Reader` matches NOTHING
// and would refuse every model, file-backed ones included — a total
// failure this row's file-backed siblings would expose.
const initCommandReadModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "initcmdread"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.stage]
provenance = "owned"
kind = "scalar"

[read.state]
role = "state"
command = ["cat", "{artifact}"]
output = "raw"
keys = ["stage"]
timeout = "10s"

[write.state]
role = "state"
path = "flow.state"
keys = ["stage"]
timeout = "2s"
read_back = true

[initial]
stage = "seeded"

[context.done]
[context.done.match.stage]
eq = "final"

[[rule]]
id = "advance"
[rule.match.stage]
eq = "seeded"
[rule.match.recognized]
eq = "advance"
[rule.write]
stage = "final"
`

// initSealModel's write accessor declares an UNREACHABLE read-back
// locator, so a write through it leaves the artifact SEALED — carrying
// `flowbind::sealedKey` and nothing else once its owned keys are cleared.
// It is S9's setup path: the sealed store is a ONE-key, NON-EMPTY store
// that init-state must decline to seed, at exit 0 no-op success.
const initSealModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "initseal"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.stage]
provenance = "owned"
kind = "scalar"

[read.state]
role = "state"
path = "flow.state"
keys = ["stage"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state.readback-unreachable"
keys = ["stage"]
timeout = "2s"
read_back = true

[initial]
stage = "seeded"

[context.done]
[context.done.match.stage]
eq = "final"

[[rule]]
id = "advance"
[rule.match.stage]
eq = "seeded"
[rule.match.recognized]
eq = "advance"
[rule.write]
stage = "final"
`

// initMismatchModel binds the READ accessor to a DIFFERENT artifact path
// than its writer, which is S8's construction: the read-back reads a store
// the write never touched, so a value pre-seeded there conflicts and the
// read-back COMPLETES AND DISAGREES.
//
// The two roles are what makes the divergence bindable from argv: the
// writer takes `state`, the reader takes `mirror`, and the invocation
// points them at different files.
const initMismatchModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "initmismatch"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.stage]
provenance = "owned"
kind = "scalar"

[read.state]
role = "mirror"
path = "flow.state"
keys = ["stage"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["stage"]
timeout = "2s"
read_back = true

[initial]
stage = "seeded"

[context.done]
[context.done.match.stage]
eq = "final"

[[rule]]
id = "advance"
[rule.match.stage]
eq = "seeded"
[rule.match.recognized]
eq = "advance"
[rule.write]
stage = "final"
`

// --- invocation builders -------------------------------------------------

// initStateArgs renders one `flow init-state` argv over model with the
// given `role=path` bindings. The verb takes the shared selection flags and
// explicit `--artifact` bindings and NO write grammar (REQ-10).
func initStateArgs(model string, bindings ...string) []string {
	args := []string{"flow", "init-state", "--model", model}
	for _, b := range bindings {
		args = append(args, "--artifact", b)
	}
	return append(args, "--as=json")
}

// readStateArgs renders the `flow read-state` argv the round-trip oracles
// use. `read-state` runs EVERY declared reader, so every declared role must
// be bound here even when init-state needed fewer.
func readStateArgs(model string, bindings ...string) []string {
	args := []string{"flow", "read-state", "--model", model}
	for _, b := range bindings {
		args = append(args, "--artifact", b)
	}
	return append(args, "--as=json")
}

// --- oracles -------------------------------------------------------------

// codeCobraUnregistered is what cobra raises for a command it does not
// know. Every refusal oracle in this suite REJECTS it by name.
//
// This is the tautology gate. `init-state` is the verb this record adds, so
// before it is registered EVERY invocation below refuses with this code —
// and an oracle that asked only "did it refuse?" would go green against a
// tree that has no verb at all. The shipped 0005 harness states the same
// rule for the same reason: "an unregistered command yields cobra's
// `command-error`, which is exactly the tautology the red gate must
// exclude."
const codeCobraUnregistered = "command-error"

// requireVerbRegistered fails when err is cobra's unregistered-command
// refusal. Every oracle that accepts a refusal calls it first, so no
// assertion in this suite can be satisfied by the verb's ABSENCE.
//
// TEST-FIXTURE (Phase 2): the CODE alone does not discriminate.
// `::cobraErrorToCLIError` maps EVERY cobra/pflag error to `command-error`,
// so an unknown FLAG on a registered verb carries the same code as an
// unregistered COMMAND — and REQ-10 and REQ-102 both assert precisely that
// unknown-flag refusal, which a code-only guard rejects as the tautology.
// The two are separable by cobra's own message: an unregistered command says
// `unknown command "<name>" for "<parent>"`. Keying on the message is what
// makes the guard exclude the tautology it names WITHOUT excluding the
// refusals the contract fixes.
func requireVerbRegistered(t *testing.T, ce *clierr.CLIError) {
	t.Helper()

	if ce.Code == codeCobraUnregistered &&
		strings.Contains(ce.Message, "unknown command") {
		t.Fatalf("the `flow %s` verb is not registered — cobra refused the "+
			"COMMAND (%s: %s), so this assertion says nothing about the "+
			"contract. A refusal that only proves the verb is missing is "+
			"the tautology this suite excludes",
			initStateVerb, ce.Code, ce.Message)
	}
}

// initRefusal drives args, asserts it REFUSED at the given exit, and
// returns the structured error.
//
// It deliberately does NOT take a code: the two codes new to this verb are
// non-normative in spelling (REQ-65), so a shared oracle that demanded one
// would be the very string pin the contract forbids. Callers that need a
// SHARED code (`flow-artifact-missing`, `flow-model-invalid`) use the
// shipped `requireRefusal` instead, which pins the code exactly.
func initRefusal(t *testing.T, exit int, args ...string) *clierr.CLIError {
	t.Helper()

	stdout, _, err := runCmd(t, args...)
	if err == nil {
		t.Fatalf("invocation succeeded; want a refusal at exit %d\nstdout:\n%s",
			exit, stdout)
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("refusal is not a structured CLIError: %v", err)
	}
	requireVerbRegistered(t, ce)
	if got := clierr.ExitCodeFor(err); got != exit {
		t.Errorf("exit code = %d; want %d — C1 fixes the exit GROUP even "+
			"where it leaves the code spelling non-normative", got, exit)
	}
	return ce
}

// requireInitRefusedOnItsOwnTerms asserts args refused with a STRUCTURED
// refusal the VERB raised — never cobra's unregistered-command error, and
// never an unknown-flag error.
//
// It is the oracle for every clause whose claim is "this input is NOT
// admitted". Such a clause is satisfied vacuously by a missing verb or a
// missing flag, so the oracle names both defects out.
func requireInitRefusedOnItsOwnTerms(t *testing.T, what string, args ...string) *clierr.CLIError {
	t.Helper()

	stdout, _, err := runCmd(t, args...)
	if err == nil {
		t.Fatalf("%s was ACCEPTED; the contract fixes that it is not\n"+
			"stdout:\n%s", what, stdout)
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("%s: refusal is not a structured CLIError: %v", what, err)
	}
	requireVerbRegistered(t, ce)
	return ce
}

// artifactBytes returns the artifact's bytes and whether it exists. An
// absent file is the "never written" state, which is distinct from an
// empty one only at the filesystem level — every content-level claim in C1
// treats the two the same.
func artifactBytes(t *testing.T, path string) ([]byte, bool) {
	t.Helper()

	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, false
	}
	if err != nil {
		t.Fatalf("read artifact %s: %v", path, err)
	}
	return b, true
}

// requireBytesUnchanged asserts the artifact is byte-identical to the
// snapshot taken before the invocation, INCLUDING the absent-vs-present
// distinction. It is the "zero writes" oracle every no-op and refusal
// clause states on artifact bytes.
func requireBytesUnchanged(t *testing.T, path string, before []byte, existed bool, what string) {
	t.Helper()

	after, exists := artifactBytes(t, path)
	if exists != existed {
		t.Errorf("%s: artifact %s existence changed (before=%v after=%v); "+
			"the contract fixes ZERO writes on this arm", what, path, existed, exists)
		return
	}
	if string(after) != string(before) {
		t.Errorf("%s: artifact %s bytes changed\nbefore: %q\nafter:  %q\n"+
			"the contract fixes ZERO writes on this arm",
			what, path, string(before), string(after))
	}
}

// snapshot captures an artifact's bytes-and-existence for a later
// requireBytesUnchanged.
func snapshot(t *testing.T, path string) ([]byte, bool) {
	t.Helper()
	return artifactBytes(t, path)
}

// --- payload readers -----------------------------------------------------
//
// `0019:D-wire-byte-format` defers payload field NAMES to implementation
// while fixing its CONTENT, so these helpers locate content BY SHAPE:
// they scan the decoded `data` object for a string-array field and return
// the candidates. A test then asserts on the CONTENT (which key names
// appear), never on which field carried it.

// stringArrayFields returns every field of data whose value is a JSON array
// of strings, keyed by field name.
func stringArrayFields(data map[string]any) map[string][]string {
	out := map[string][]string{}
	for name, raw := range data {
		arr, ok := raw.([]any)
		if !ok {
			continue
		}
		vals := make([]string, 0, len(arr))
		allStrings := true
		for _, item := range arr {
			s, ok := item.(string)
			if !ok {
				allStrings = false
				break
			}
			vals = append(vals, s)
		}
		if allStrings {
			sort.Strings(vals)
			out[name] = vals
		}
	}
	return out
}

// payloadCarriesKeySet reports whether SOME string-array field of the
// payload carries exactly the given key set. It is how a test asserts the
// payload's CONTENT without pinning the field NAME the contract defers.
func payloadCarriesKeySet(data map[string]any, want []string) (string, bool) {
	sorted := append([]string(nil), want...)
	sort.Strings(sorted)
	for name, got := range stringArrayFields(data) {
		if len(got) != len(sorted) {
			continue
		}
		match := true
		for i := range got {
			if got[i] != sorted[i] {
				match = false
				break
			}
		}
		if match {
			return name, true
		}
	}
	return "", false
}

// payloadMentionsKey reports whether any string-array field of the payload
// contains key. Used where a clause fixes that a key IS named without
// fixing the whole set.
func payloadMentionsKey(data map[string]any, key string) bool {
	for _, got := range stringArrayFields(data) {
		for _, v := range got {
			if v == key {
				return true
			}
		}
	}
	return false
}

// payloadValueStrings returns every string leaf anywhere in the payload,
// so a test can assert a seeded VALUE does NOT appear (REQ-57).
func payloadValueStrings(t *testing.T, stdout string) []string {
	t.Helper()

	var env struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); err != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", err, stdout)
	}
	var any0 any
	if err := json.Unmarshal(env.Data, &any0); err != nil {
		t.Fatalf("`data` is not JSON: %v", err)
	}
	var out []string
	var walk func(v any)
	walk = func(v any) {
		switch t := v.(type) {
		case string:
			out = append(out, t)
		case []any:
			for _, item := range t {
				walk(item)
			}
		case map[string]any:
			for _, item := range t {
				walk(item)
			}
		}
	}
	walk(any0)
	return out
}

// --- read-state oracles --------------------------------------------------

// ownedFromReadState collapses `read-state`'s per-reader reports into one
// key->value map of the tags the readers actually returned. A key ABSENT
// from the artifact is absent from this map: `readerOutput` carries the
// DECLARED key set in `keys` beside the tags returned in `tags`, and a key
// in `keys` with no `tags` entry is established-absent.
func ownedFromReadState(t *testing.T, stdout string) map[string]string {
	t.Helper()

	data := flowData(t, stdout)
	readers, ok := data["readers"].([]any)
	if !ok {
		t.Fatalf("read-state payload has no `readers` array: %v", keysOf(data))
	}
	out := map[string]string{}
	for _, r := range readers {
		rm, ok := r.(map[string]any)
		if !ok {
			continue
		}
		tags, ok := rm["tags"].(map[string]any)
		if !ok {
			continue
		}
		for k, v := range tags {
			s, ok := v.(string)
			if !ok {
				continue
			}
			out[k] = s
		}
	}
	return out
}

// requireReadsBack asserts `read-state` reports key at exactly want — the
// REQ-107 value-for-value equality RT1 and `0019:D-identity` fix. Exit 0 is
// NOT sufficient and no oracle here accepts it as such.
func requireReadsBack(t *testing.T, owned map[string]string, key, want string) {
	t.Helper()

	got, present := owned[key]
	if !present {
		t.Errorf("read-state reports %q ABSENT; want it PRESENT at %q — "+
			"RT1 fixes value-for-value read-back, not merely exit 0", key, want)
		return
	}
	if got != want {
		t.Errorf("read-state reports %q = %q; want %q — RT1's equality is "+
			"REQ-107's, value-for-value in canonical form", key, got, want)
	}
}

// requireReadsAbsent asserts `read-state` reports key ABSENT. It is the
// cleared-key oracle: REQ-107 fixes that a cleared key reads back absent
// for every reader, and C1 forbids any read path synthesizing `[initial]`.
func requireReadsAbsent(t *testing.T, owned map[string]string, key, what string) {
	t.Helper()

	if got, present := owned[key]; present {
		t.Errorf("%s: read-state reports %q = %q; want it ABSENT — no read "+
			"path may synthesize, default, or fall back to `[initial]`",
			what, key, got)
	}
}

// --- `next` oracles ------------------------------------------------------

// nextAbsentKeys returns every key `flow next` reported under
// `unknown[].reason: absent`, across every candidate. It is the surface
// where an operator meets the first-run wall the verb exists to remove.
func nextAbsentKeys(t *testing.T, stdout string) map[string]bool {
	t.Helper()

	data := flowData(t, stdout)
	out := map[string]bool{}
	cands, ok := data["candidates"].([]any)
	if !ok {
		return out
	}
	for _, c := range cands {
		cm, ok := c.(map[string]any)
		if !ok {
			continue
		}
		unknown, ok := cm["unknown"].([]any)
		if !ok {
			continue
		}
		for _, u := range unknown {
			um, ok := u.(map[string]any)
			if !ok {
				continue
			}
			if reason, _ := um["reason"].(string); reason != "absent" {
				continue
			}
			if key, ok := um["key"].(string); ok {
				out[key] = true
			}
		}
	}
	return out
}
