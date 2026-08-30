package flowbind_test

// RDR 0025 — carrier selection at the single production construction site
// (C1's residue rule, C6's gate and signature, the Selection LBD, S5b).
//
// The mutant this file kills is the one C1 spends a paragraph on: a
// registry whose selection is not TOTAL falls through to
// `Reader{Path: acc.Path}` with an EMPTY path, and `flowbind.go::load`
// maps a non-existent file to an empty store BY DESIGN — so a command
// entry that reached the file binding would read every declared key as
// ABSENT and confirm an unapplied write as verified. Silent state loss,
// never a crash. The residue arm is therefore asserted by the BEHAVIOUR of
// the constructed binding, not by its type.
//
// The second mutant: a gate implemented anywhere but the constructor
// leaves a construction path that yields an executing binding. So the gate
// is asserted on what `Registry` RETURNS, over the same model, with the
// flag on and off.

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/flowbind"
	"github.com/cwensel/intrastate/internal/table"
)

const cmdSelectionModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "selflow"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[read.state]
role = "state"
command = ["reader", "{artifact}"]
output = "raw"
keys = ["status"]
timeout = "2s"

[write.state]
role = "state"
command = ["writer", "{artifact}"]
keys = ["status"]
timeout = "2s"
read_back = true

[gate.approval]
role = "state"
command = ["gater", "{artifact}"]
exit_verdicts = { "0" = "allow", "1" = "deny" }
keys = ["status"]
timeout = "2s"

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "advance-draft"
gate = ["approval"]
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"
`

// pathBackedModel is the same shape with `path` carriers, so the two
// carriers can be compared side by side.
var pathBackedModel = strings.NewReplacer(
	`command = ["reader", "{artifact}"]
output = "raw"`, `path = "flow.status"`,
	`command = ["writer", "{artifact}"]`, `path = "flow.status"`,
	`command = ["gater", "{artifact}"]
exit_verdicts = { "0" = "allow", "1" = "deny" }`, `path = "flow.gate"`,
).Replace(cmdSelectionModel)

func loadModel(t *testing.T, src, id string) *table.Model {
	t.Helper()

	m, err := table.Load([]byte(src), id)
	if err != nil {
		t.Fatalf("%s must load clean: %v", id, err)
	}
	return m
}

func bindingFor(t *testing.T, reg accessor.Registry, name string, c accessor.Capability) accessor.Binding {
	t.Helper()

	for _, d := range reg.Definitions {
		if d.Identity.Name == name && d.Identity.Capability == c {
			return d.Binding
		}
	}
	t.Fatalf("no %s definition named %q in the registry", c, name)
	return nil
}

func ctx(t *testing.T) context.Context {
	t.Helper()

	c, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return c
}

// --- C6: the signature -----------------------------------------------------

// REQ-94: "signature: ... It gains BOTH new inputs and no others, in this
// order: `func Registry(m *table.Model, baseDir string, allowCommands bool)
// accessor.Registry`."
// REQ-95: "the return type is unchanged and `Registry` does not gain an
// error return — a carrier-less or gate-refused entry yields a refusing
// binding (C1), not a construction failure"
// BOUNDARY
//
// The compile-time half is the signature this file calls; the runtime half
// is that a model whose every entry is refused still CONSTRUCTS, with its
// definitions intact.
func TestReq94_RegistryTakesBaseDirAndAllowCommandsAndReturnsNoError(t *testing.T) {
	m := loadModel(t, cmdSelectionModel, "selection.toml")

	// Gate OFF over a wholly command-backed model: every entry is refused
	// at invocation, and construction still succeeds with every definition
	// present. A constructor that returned an error, or dropped the
	// refused entries, would break `unknown_accessor` diagnosis.
	reg := flowbind.Registry(m, "/models", false)

	if reg.Flow != m.ID {
		t.Errorf("Flow = %q; want %q", reg.Flow, m.ID)
	}
	if len(reg.Definitions) != 3 {
		t.Fatalf("Definitions = %d; want the three declared entries even with "+
			"the gate off — a refused entry is a refusing BINDING, not a "+
			"construction failure", len(reg.Definitions))
	}
	if !slices.Contains(reg.OwnedTags, "status") {
		t.Errorf("OwnedTags = %#v; want the model's owned tags", reg.OwnedTags)
	}
}

// --- Selection LBD / C1: the carrier discriminator ------------------------

// REQ-113: "the registry selects the binding constructor by carrier field:
// `path` → file binding (today's `flowbind`), `command` → command binding.
// C1's exactly-one rule makes the selection total; no precedence order
// exists to get wrong."
// REQ-114: "The file binding keeps its magic-suffix vocabulary
// (`verdictFor`, `unreachable`) **unchanged**"
// HAPPY PATH — REQ-114 is a NEGATIVE REQ.
//
// The oracle is BEHAVIOUR, not type identity: the command entry's read must
// attempt an invocation (and so refuse, with the gate on, because `reader`
// is not a real binary), while the path entry's read must reach the FILE
// and report established absence on a missing artifact — the deliberate
// empty-artifact semantics C1 promises not to disturb.
func TestReq113_TheRegistrySelectsTheBindingByCarrierField(t *testing.T) {
	t.Run("a command entry does not reach the file binding", func(t *testing.T) {
		m := loadModel(t, cmdSelectionModel, "selection.toml")
		reg := flowbind.Registry(m, "/models", true)

		b, ok := bindingFor(t, reg, "state", accessor.CapRead).(accessor.ReadBinding)
		if !ok {
			t.Fatal("the read binding does not implement ReadBinding")
		}
		art := accessor.Artifact{Role: "state", Path: "/nonexistent/state.json"}

		values, unreadable, err := b.Read(ctx(t), art, []string{"status"})

		if err == nil {
			t.Fatalf("a command entry over a missing artifact answered "+
				"values=%#v unreadable=%#v; the file binding maps a "+
				"non-existent file to an EMPTY store by design, so reaching it "+
				"would read every declared key as ABSENT and confirm an "+
				"unapplied write", values, unreadable)
		}
		for _, v := range values {
			if v.Absent {
				t.Errorf("%q came back established-absent; that is the file "+
					"binding's answer, not a command entry's", v.Key)
			}
		}
	})

	t.Run("a path entry still reaches the file binding unchanged", func(t *testing.T) {
		m := loadModel(t, pathBackedModel, "path-backed.toml")
		reg := flowbind.Registry(m, "/models", true)

		b, ok := bindingFor(t, reg, "state", accessor.CapRead).(accessor.ReadBinding)
		if !ok {
			t.Fatal("the read binding does not implement ReadBinding")
		}
		art := accessor.Artifact{Role: "state", Path: "/nonexistent/state.json"}

		values, _, err := b.Read(ctx(t), art, []string{"status"})
		if err != nil {
			t.Fatalf("a PATH entry over a missing artifact refused: %v — the "+
				"empty-artifact semantics let a flow's first write establish "+
				"state, and C1 must not disturb them", err)
		}
		if len(values) != 1 || !values[0].Absent {
			t.Errorf("values = %#v; want the file binding's established-absent "+
				"answer, unchanged", values)
		}

		// The DISCRIMINATING half. "The path entry still works" is satisfied
		// vacuously while the registry has no discriminator at all, so the
		// sibling must differ: the SAME artifact under a COMMAND entry must
		// NOT produce this answer.
		cm := loadModel(t, cmdSelectionModel, "selection.toml")
		cb := bindingFor(t, flowbind.Registry(cm, "/models", true),
			"state", accessor.CapRead).(accessor.ReadBinding)
		cvalues, _, cerr := cb.Read(ctx(t), art, []string{"status"})
		if cerr == nil && len(cvalues) == 1 && cvalues[0].Absent {
			t.Error("a COMMAND entry produced the file binding's " +
				"established-absent answer over the same missing artifact; " +
				"the two carriers must select different bindings")
		}
	})
}

// REQ-8: "at RUNTIME the constructor refuses a carrier-less entry with a
// refusing binding (`execution_failure`, Detail naming the malformed entry)
// and MUST NOT fall through to a `Path: \"\"` file binding, which would
// read every declared key as absent and confirm an unapplied write"
// REQ-11: "The selection in `flowbind.Registry` must be **exhaustive, and
// fail closed on the residue**." ... "the constructor selects `command`
// first, `path` second, and **panics-free refuses** the third case"
// REQ-131 (S7's C1 residue arm): "an entry with neither carrier builds a
// **refusing** binding, never a `Path: \"\"` file binding."
// ADVERSARIAL
//
// C1's load-time exactly-one rule makes this residue unreachable THROUGH
// THE LOADER; the constructor refuses it anyway, because a binding built
// from an IN-MEMORY model never passed the loader. So the fixture is
// deliberately an in-memory model — the only way this arm is reachable.
func TestReq8_ACarrierLessEntryBuildsARefusingBindingNotAnEmptyPathOne(t *testing.T) {
	// An in-memory model with an entry carrying NEITHER carrier. The loader
	// would refuse this (C5 `command_and_path_conflict`); nothing here goes
	// through the loader, which is exactly the case the clause covers.
	m := &table.Model{
		ID: "memflow",
		Tags: map[string]table.TagDecl{
			"status": {Provenance: table.ProvenanceOwned, Kind: "scalar"},
		},
		Readers: map[string]table.Accessor{
			"state": {Role: "state", Keys: []string{"status"}, Timeout: "2s"},
		},
		Writers: map[string]table.Accessor{
			"state": {
				Role: "state", Keys: []string{"status"}, Timeout: "2s", ReadBack: true,
			},
		},
		Gates: map[string]table.Accessor{
			"approval": {Role: "state", Keys: []string{"status"}, Timeout: "2s"},
		},
	}

	reg := flowbind.Registry(m, "/models", true)
	art := accessor.Artifact{Role: "state", Path: "/nonexistent/state.json"}

	t.Run("read", func(t *testing.T) {
		b := bindingFor(t, reg, "state", accessor.CapRead).(accessor.ReadBinding)

		values, _, err := b.Read(ctx(t), art, []string{"status"})
		if err == nil {
			t.Fatalf("a carrier-less read answered %#v; a `Path: \"\"` file "+
				"binding reads every declared key as ABSENT, which confirms an "+
				"unapplied write — silent state loss, never a crash", values)
		}
		requireNames(t, err, "state")
	})

	t.Run("write", func(t *testing.T) {
		b := bindingFor(t, reg, "state", accessor.CapWrite).(accessor.WriteBinding)

		err := b.Apply(ctx(t), art, nil)
		if err == nil {
			t.Fatal("a carrier-less write applied; the residue must refuse")
		}
		// The DISCRIMINATING half. A `Path: ""` file binding also fails on
		// Apply — it just fails for the wrong reason, and its READ silently
		// confirms the unapplied write. So the refusal must name the
		// malformed ENTRY, which is what distinguishes a refusing binding
		// from an empty-path file binding that happens to error.
		requireNames(t, err, "state")
	})

	t.Run("gate", func(t *testing.T) {
		b := bindingFor(t, reg, "approval", accessor.CapGate).(accessor.GateBinding)

		verdict, _, err := b.Gate(ctx(t), art)
		if err == nil {
			t.Fatalf("a carrier-less gate answered %q; the residue must refuse",
				verdict)
		}
	})
}

func requireNames(t *testing.T, err error, want string) {
	t.Helper()

	var ee *accessor.ExecError
	if !asExec(err, &ee) {
		t.Fatalf("the refusal is not an *accessor.ExecError: %v", err)
	}
	if !strings.Contains(ee.Detail, want) {
		t.Errorf("Detail = %q; want it to name the malformed entry %q",
			ee.Detail, want)
	}
}

// --- C6: the gate, at the construction site ------------------------------

// REQ-93: "gate site: ONE — `buildRequest` reads the flag and passes it to
// `flowbind.Registry`, the single production construction site
// (`flow_exec.go:830`), which builds refusing command bindings when it is
// unset."
// REQ-96: "The gate sites **in the command binding's constructor**, not in
// the executor"
// REQ-117: "Phase 3: Selection — Discriminate the constructor by carrier
// field at `internal/cli/flowbind/registry.go::Registry` ... Register
// `--allow-commands` ... and thread it (C6) from `buildRequest` ... in the
// SAME phase: selection is the step that first makes a command reachable,
// so the gate must not lag it by even one commit."
// ADVERSARIAL
//
// The discriminating pair: the SAME model, the same entries, the flag off
// and on. With it off every command binding refuses naming the gate; with
// it on the refusal is about the missing tool, not about the gate — which
// is what proves the gate is the thing that changed.
func TestReq93_TheGateIsAppliedAtTheSingleConstructionSite(t *testing.T) {
	m := loadModel(t, cmdSelectionModel, "selection.toml")
	art := accessor.Artifact{Role: "state", Path: "/tmp/state.json"}

	off := flowbind.Registry(m, "/models", false)
	on := flowbind.Registry(m, "/models", true)

	t.Run("gate off: every capability refuses naming allow_commands", func(t *testing.T) {
		r := bindingFor(t, off, "state", accessor.CapRead).(accessor.ReadBinding)
		_, _, rerr := r.Read(ctx(t), art, []string{"status"})
		requireGateNamed(t, rerr)

		w := bindingFor(t, off, "state", accessor.CapWrite).(accessor.WriteBinding)
		requireGateNamed(t, w.Apply(ctx(t), art, nil))

		g := bindingFor(t, off, "approval", accessor.CapGate).(accessor.GateBinding)
		_, _, gerr := g.Gate(ctx(t), art)
		requireGateNamed(t, gerr)
	})

	t.Run("gate on: the refusal is no longer about the gate", func(t *testing.T) {
		// `reader` is not a real binary, so the read still fails — but the
		// cause must be the spawn, not the gate. Without this arm, "the
		// gate refuses" is satisfied by a registry that refuses everything.
		r := bindingFor(t, on, "state", accessor.CapRead).(accessor.ReadBinding)
		_, _, err := r.Read(ctx(t), art, []string{"status"})
		if err == nil {
			t.Fatal("the read succeeded against a non-existent argv0")
		}
		var ee *accessor.ExecError
		if asExec(err, &ee) && strings.Contains(ee.Detail, "allow_commands") {
			t.Errorf("Detail = %q with --allow-commands SET; the gate must be "+
				"the only thing the flag changes", ee.Detail)
		}
	})

	t.Run("a PATH-backed entry is ungated", func(t *testing.T) {
		// The gate is on command execution, not on the seam. A path-backed
		// model must behave identically with the flag off — otherwise every
		// existing model breaks (REQ-139).
		pm := loadModel(t, pathBackedModel, "path-backed.toml")
		reg := flowbind.Registry(pm, "", false)

		b := bindingFor(t, reg, "state", accessor.CapRead).(accessor.ReadBinding)
		values, _, err := b.Read(ctx(t), accessor.Artifact{
			Role: "state", Path: "/nonexistent/state.json",
		}, []string{"status"})
		if err != nil {
			t.Fatalf("a path-backed read refused with the gate off: %v — the "+
				"gate bounds COMMAND execution and nothing else", err)
		}
		if len(values) != 1 {
			t.Errorf("values = %#v; want the file binding's unchanged answer", values)
		}
	})
}

func requireGateNamed(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("a command binding ran with --allow-commands unset")
	}
	var ee *accessor.ExecError
	if !asExec(err, &ee) {
		t.Fatalf("the refusal is not an *accessor.ExecError: %v", err)
	}
	if !strings.Contains(ee.Detail, "allow_commands") {
		t.Errorf("Detail = %q; want it to name `allow_commands`", ee.Detail)
	}
}

func asExec(err error, target **accessor.ExecError) bool {
	for err != nil {
		if ee, ok := err.(*accessor.ExecError); ok {
			*target = ee
			return true
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = u.Unwrap()
	}
	return false
}

// --- S5b: read-back reader selection --------------------------------------

// REQ-128 (S5b arm b): "The disjoint-keys arm: pending 0016's fail-closed
// uniqueness, the selection is `registry.go`'s name-sorted first match —
// asserted **by the selected reader's identity**, never by \"read-back
// succeeded\" ... arm (b) asserts **both halves**: that the registry emits
// readers name-sorted, and that `readerFor` returns the first."
// REQ-120: "0025 lands after 0016 or inherits first-match until it does"
// DOMAIN EDGE
//
// The property is EMERGENT across two packages, so a test on the composed
// outcome alone passes while either half moves. Both halves are asserted
// here: the registry's emission ORDER, and — as the sibling arm below —
// that `readerFor` takes the first. The fixture is the LIVE migration
// shape: two readers on one role with DISJOINT key sets, which loads clean
// because `checkAccessorBindings` counts readers per TAG KEY, not per role.
func TestReq128_TheRegistryEmitsReadersNameSortedForFirstMatchSelection(t *testing.T) {
	m := loadModel(t, twoReaderModel, "two-readers.toml")
	reg := flowbind.Registry(m, "/models", true)

	// Half one: the emission order. `readerFor` takes the FIRST match by
	// slice order, so the slice order is half the selection.
	var readers []string
	for _, d := range reg.Definitions {
		if d.Identity.Capability == accessor.CapRead {
			readers = append(readers, d.Identity.Name)
		}
	}
	if !slices.IsSorted(readers) {
		t.Errorf("the registry emitted readers %#v; want them NAME-SORTED — "+
			"first-match selection is only predictable if the order is",
			readers)
	}
	want := []string{"a-legacy", "b-command"}
	if !slices.Equal(readers, want) {
		t.Errorf("readers = %#v; want %#v", readers, want)
	}
}

// REQ-128 (S5b arm a): "The owned-key arm: a second reader serving a
// command write's owned key is REFUSED at load
// (`malformed_accessor_binding`) — so read-back's reader is unique by
// construction, not by selection order."
// ADVERSARIAL
func TestReq128_ASecondReaderOnAWritesOwnedKeyIsRefusedAtLoad(t *testing.T) {
	// Both readers serve `status`, which the write owns. The loader counts
	// readers per TAG KEY and requires exactly one for an owned tag.
	src := strings.Replace(twoReaderModel, `keys = ["stage"]`, `keys = ["status"]`, 1)
	if src == twoReaderModel {
		t.Fatal("the mutant substitution did not apply")
	}

	_, err := table.Load([]byte(src), "two-readers-owned.toml")
	if err == nil {
		t.Fatal("two readers on a write's owned key loaded clean; read-back's " +
			"reader must be unique by CONSTRUCTION, not by selection order")
	}
	cat, ok := table.CategoryOf(err)
	if !ok {
		t.Fatalf("the refusal carries no category: %v", err)
	}
	if cat != table.CatMalformedAccessorBinding {
		t.Errorf("category = %q; want %q", cat, table.CatMalformedAccessorBinding)
	}
}

// twoReaderModel is the LIVE migration shape S5b names: a path-backed
// reader kept for rollback beside a new command-backed one, on the same
// role, with DISJOINT key sets — which loads clean and where `readerFor`
// silently takes the first.
const twoReaderModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "migrflow"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[tags.stage]
provenance = "observed"
kind = "scalar"

[read.a-legacy]
role = "state"
path = "flow.stage"
keys = ["stage"]
timeout = "2s"

[read.b-command]
role = "state"
command = ["reader", "{artifact}"]
output = "raw"
keys = ["status"]
timeout = "2s"

[write.state]
role = "state"
command = ["writer", "{artifact}"]
keys = ["status"]
timeout = "2s"
read_back = true

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "advance-draft"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"
`

// REQ-109: "the accessor identity stays RDR 0004's `(flow, name,
// capability)` triple; whether an entry is path- or command-backed does not
// enter identity."
// BOUNDARY — NEGATIVE REQ.
//
// The discriminating pair: two models identical except for the carrier
// produce identical IDENTITY sets. A registry that folded the carrier into
// identity would produce different triples and break every lookup.
func TestReq109_TheCarrierDoesNotEnterAccessorIdentity(t *testing.T) {
	cmd := flowbind.Registry(loadModel(t, cmdSelectionModel, "cmd.toml"), "/m", true)
	pth := flowbind.Registry(loadModel(t, pathBackedModel, "path.toml"), "/m", true)

	ids := func(reg accessor.Registry) []accessor.Identity {
		var out []accessor.Identity
		for _, d := range reg.Definitions {
			out = append(out, d.Identity)
		}
		return out
	}

	got, want := ids(cmd), ids(pth)
	if len(got) != len(want) {
		t.Fatalf("the two carriers produced %d and %d definitions; the carrier "+
			"does not enter identity", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("identity[%d] = %+v under `command`; want %+v, the same "+
				"triple the `path` carrier produces", i, got[i], want[i])
		}
	}

	// And the discriminating half: the two registries must nonetheless bind
	// DIFFERENT bindings, or the identity match above is satisfied by a
	// registry that ignored the carrier entirely.
	art := accessor.Artifact{Role: "state", Path: "/nonexistent/state.json"}
	cb := bindingFor(t, cmd, "state", accessor.CapRead).(accessor.ReadBinding)
	pb := bindingFor(t, pth, "state", accessor.CapRead).(accessor.ReadBinding)

	_, _, cerr := cb.Read(ctx(t), art, []string{"status"})
	_, _, perr := pb.Read(ctx(t), art, []string{"status"})
	if (cerr == nil) == (perr == nil) {
		t.Errorf("both carriers answered alike (command err=%v, path err=%v); "+
			"identity is shared but the BINDING is selected by carrier",
			cerr, perr)
	}
}
