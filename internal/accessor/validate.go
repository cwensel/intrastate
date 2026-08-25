package accessor

import "slices"

// Validate runs the eight validation arms over the registry and the
// identities the model references. It returns every finding it can
// decide, each carrying its own named code. A valid registry returns the
// EMPTY set — the negative control that fails a validator which rejects
// everything (ORA 1).
//
// The eight arms (`0004:TD`, TS 1): missing accessor, multiply-bound
// accessor identity, capability mismatch, missing/non-positive timeout,
// missing write read-back, ambient artifact discovery, write to a
// non-owned tag, and a missing or empty read requested-key set.
//
// RDR 0002's loader already refuses several of these shapes at load
// (`internal/table/load.go::accessorTable`). This RDR owns EXECUTION, so
// the arms here decide the same defects at the execution boundary over
// the values the executor is actually handed; nothing is re-parsed.
func Validate(reg Registry, referenced []Identity) []Finding {
	var findings []Finding

	// --- per-definition arms -------------------------------------------
	seen := map[Identity]int{}
	for _, d := range reg.Definitions {
		seen[d.Identity]++

		if d.AmbientDiscovery {
			findings = append(findings, Finding{
				Code:       CodeAmbientArtifactDiscovery,
				Accessor:   d.Identity.Name,
				Capability: d.Identity.Capability,
				Detail:     d.Accessor.Role,
			})
		}
		if _, ok := d.timeout(); !ok {
			findings = append(findings, Finding{
				Code:       CodeMissingOrNonPositiveTimeout,
				Accessor:   d.Identity.Name,
				Capability: d.Identity.Capability,
				Detail:     d.Accessor.Timeout,
			})
		}
		switch d.Identity.Capability {
		case CapRead:
			// The requested key set is the definition's declared
			// metadata; a reader without one has nothing to be complete
			// over (`0004:C5`, the eighth arm).
			if len(d.Accessor.Keys) == 0 {
				findings = append(findings, Finding{
					Code:       CodeMissingRequestedKeySet,
					Accessor:   d.Identity.Name,
					Capability: d.Identity.Capability,
				})
			}
		case CapWrite:
			if !d.Accessor.ReadBack {
				findings = append(findings, Finding{
					Code:       CodeMissingWriteReadBack,
					Accessor:   d.Identity.Name,
					Capability: d.Identity.Capability,
				})
			}
			for _, key := range d.Accessor.Keys {
				if !slices.Contains(reg.OwnedTags, key) {
					findings = append(findings, Finding{
						Code:       CodeWriteNonOwnedTag,
						Accessor:   d.Identity.Name,
						Capability: d.Identity.Capability,
						Detail:     key,
					})
				}
			}
		}
	}

	// --- multiply-bound identities -------------------------------------
	//
	// Only a second binding of the SAME triple is multiply-bound:
	// capability is part of the identity, so `[read.x]` and `[write.x]`
	// are two identities, not a rebinding (LBD, Identity).
	for _, d := range reg.Definitions {
		if seen[d.Identity] > 1 {
			findings = append(findings, Finding{
				Code:       CodeMultiplyBoundAccessor,
				Accessor:   d.Identity.Name,
				Capability: d.Identity.Capability,
			})
			seen[d.Identity] = 0
		}
	}

	// --- referenced identities -----------------------------------------
	for _, ref := range referenced {
		if _, ok := reg.Lookup(ref.Name, ref.Capability); ok {
			continue
		}
		code := CodeMissingAccessor
		if reg.bound(ref.Name) {
			// The id exists, but under another capability table. That is
			// a different defect from "not bound at all".
			code = CodeCapabilityMismatch
		}
		findings = append(findings, Finding{
			Code:       code,
			Accessor:   ref.Name,
			Capability: ref.Capability,
		})
	}

	return findings
}
