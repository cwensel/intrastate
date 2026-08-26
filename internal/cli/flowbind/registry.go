package flowbind

import (
	"slices"

	"github.com/newcoinc/intrastate/internal/accessor"
	"github.com/newcoinc/intrastate/internal/table"
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
func Registry(m *table.Model) accessor.Registry {
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
