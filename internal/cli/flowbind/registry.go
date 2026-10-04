package flowbind

import (
	"slices"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/cmdbind"
	"github.com/cwensel/intrastate/internal/table"
)

// Registry builds the validated accessor registry for one loaded model,
// binding each declared accessor to its capability's binding.
//
// Identity is the triple `(flow, name, capability)` RDR 0004 fixes, so the
// same id under `[read.x]` and `[write.x]` yields TWO definitions — two
// identities, not a rebinding (JDR 0001 §D7(ii)).
//
// OwnedTags is populated from the model's tag declarations, because
// `0004:C10` bounds what a write may apply by what the MODEL calls owned,
// not merely by what a writer's `keys` list names.
//
// Selection is by CARRIER and must be TOTAL, failing closed on the residue
// (`0025:C1`): a `command`-carrying entry builds a command binding, a
// `path`-carrying one builds today's file binding, and an entry with
// NEITHER builds a REFUSING binding — never a `Path: ""` file binding,
// which would read every declared key as absent and confirm an unapplied
// write.
//
// baseDir is the model file's directory, the root a separator-bearing
// argv0 resolves against (`0025:C2`); allowCommands is the
// `--allow-commands` gate, checked here because this is the single
// production construction site every executor's registry comes from
// (`0025:C6`).
func Registry(m *table.Model, baseDir string, allowCommands bool) accessor.Registry {
	cfg := cmdbind.Config{BaseDir: baseDir, AllowCommands: allowCommands}

	reg := accessor.Registry{
		Flow:      m.ID,
		OwnedTags: OwnedTags(m),
		// RDR 0028 `0028:C1.3` SITE: — the gate reaches the accessor as
		// STATE from this single production construction site, never read
		// across the `cmdbind` seam, which would be an import cycle.
		AllowCommands: allowCommands,
	}

	for _, name := range slices.Sorted(keys(m.Readers)) {
		acc := m.Readers[name]
		var binding accessor.Binding = Reader{Path: acc.Path}
		if commandBacked(acc) {
			binding = cmdbind.Reader{Accessor: acc, Name: name, Config: cfg}
		}
		reg.Definitions = append(reg.Definitions, accessor.Definition{
			Identity: accessor.Identity{
				Flow: m.ID, Name: name, Capability: accessor.CapRead,
			},
			Accessor: acc,
			Binding:  binding,
		})
	}
	for _, name := range slices.Sorted(keys(m.Writers)) {
		acc := m.Writers[name]
		var binding accessor.Binding = &Writer{Path: acc.Path}
		switch {
		// The `edit` arm is tested FIRST and on the carrier's own
		// PRESENCE, never through `commandBacked`. That is deliberate:
		// `commandBacked`'s `acc.Path == ""` arm IS the residue path, so
		// folding `edit` into it would leave a carrier-less entry
		// indistinguishable from an `edit`-carried one — and rerouting
		// residue to a `Path: ""` file binding is the silent state loss
		// C1's runtime arm exists to prevent (`0028:C1.1` in-memory:,
		// REQ-118).
		//
		// `edit` is a WRITE carrier only, which is why this arm appears
		// in this loop alone: a read or gate entry carrying `edit` is
		// `edit_carrier_conflict` at load, and one built in memory falls
		// through to the residue rather than becoming a second line
		// parser.
		case len(acc.Edit) != 0:
			binding = NewEditWriter(acc, name)
		// The `steps` arm follows `edit`'s discipline for the same reason:
		// tested on its own presence and ahead of `commandBacked`, whose
		// `Path == ""` arm would otherwise take a `steps` entry as residue
		// (kata q14r). It is a command-spawning binding, so it takes the
		// same `cmdbind.Config` and with it the `--allow-commands` gate.
		case len(acc.Steps) != 0:
			binding = &cmdbind.StepsWriter{Accessor: acc, Name: name, Config: cfg}
		case commandBacked(acc):
			binding = &cmdbind.Writer{Accessor: acc, Name: name, Config: cfg}
		}
		reg.Definitions = append(reg.Definitions, accessor.Definition{
			Identity: accessor.Identity{
				Flow: m.ID, Name: name, Capability: accessor.CapWrite,
			},
			Accessor: acc,
			Binding:  binding,
		})
	}
	for _, name := range slices.Sorted(keys(m.Gates)) {
		acc := m.Gates[name]
		var binding accessor.Binding = Gate{Path: acc.Path}
		if commandBacked(acc) {
			binding = cmdbind.Gate{Accessor: acc, Name: name, Config: cfg}
		}
		reg.Definitions = append(reg.Definitions, accessor.Definition{
			Identity: accessor.Identity{
				Flow: m.ID, Name: name, Capability: accessor.CapGate,
			},
			Accessor: acc,
			Binding:  binding,
		})
	}
	return reg
}

// commandBacked is the carrier discriminator, and it is what makes the
// selection TOTAL: `command` first, `path` second, and the residue — an
// entry carrying NEITHER — takes the command arm too, where it builds a
// REFUSING binding rather than a `Path: ""` file binding (`0025:C1`).
//
// That third case is the one C1 spends a paragraph on: the file binding
// FAILS OPEN on an empty path, mapping a non-existent file to an empty
// store by design, so a carrier-less entry reaching it would read every
// declared key as ABSENT and confirm an unapplied write as verified —
// silent state loss, never a crash. C1's load-time exactly-one rule makes
// the residue unreachable through the loader; a binding built from an
// in-memory model never passed it, which is why the constructor refuses
// it anyway.
func commandBacked(acc table.Accessor) bool {
	return len(acc.Command) != 0 || acc.Path == ""
}

// OwnedTags returns the model's owned tag keys, sorted.
func OwnedTags(m *table.Model) []string {
	var out []string
	for key, decl := range m.Tags {
		if decl.Provenance == table.ProvenanceOwned {
			out = append(out, key)
		}
	}
	slices.Sort(out)
	return out
}

func keys[V any](m map[string]V) func(func(string) bool) {
	return func(yield func(string) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}
