package flowbind

import (
	"slices"

	"github.com/cwensel/intrastate/internal/accessor"
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
//
// PHASE 1 DECLARATION ONLY for the two new parameters: the signature is
// the one `0025:C6` fixes so the conformance suite compiles, and Phase 2
// lands the carrier discrimination and the gate.
func Registry(m *table.Model, baseDir string, allowCommands bool) accessor.Registry {
	_, _ = baseDir, allowCommands

	reg := accessor.Registry{
		Flow:      m.ID,
		OwnedTags: OwnedTags(m),
	}

	for _, name := range slices.Sorted(keys(m.Readers)) {
		acc := m.Readers[name]
		reg.Definitions = append(reg.Definitions, accessor.Definition{
			Identity: accessor.Identity{
				Flow: m.ID, Name: name, Capability: accessor.CapRead,
			},
			Accessor: acc,
			Binding:  Reader{Path: acc.Path},
		})
	}
	for _, name := range slices.Sorted(keys(m.Writers)) {
		acc := m.Writers[name]
		reg.Definitions = append(reg.Definitions, accessor.Definition{
			Identity: accessor.Identity{
				Flow: m.ID, Name: name, Capability: accessor.CapWrite,
			},
			Accessor: acc,
			Binding:  &Writer{Path: acc.Path},
		})
	}
	for _, name := range slices.Sorted(keys(m.Gates)) {
		acc := m.Gates[name]
		reg.Definitions = append(reg.Definitions, accessor.Definition{
			Identity: accessor.Identity{
				Flow: m.ID, Name: name, Capability: accessor.CapGate,
			},
			Accessor: acc,
			Binding:  Gate{Path: acc.Path},
		})
	}
	return reg
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
