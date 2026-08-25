package accessor

import (
	"time"

	"github.com/newcoinc/intrastate/internal/resolve"
)

// Disposition is the recorded, replay-stable outcome of one invocation.
// It records artifact role, accessor name, capability, timeout, and the
// returned tag values — the five the model must record for accessor
// execution to be deterministic enough for resolver replay (CA A4).
//
// Ordering is normalized: tags are rendered sorted, so two runs over the
// same model and fixture results compare EQUAL rather than merely
// error-free (FID `resolve -> replay`, ORA 4).
type Disposition struct {
	Accessor   string
	Capability Capability
	Role       string
	Timeout    time.Duration
	// Refusal is the refusal class, empty on success.
	Refusal RefusalClass
	// Tags is the returned tag values, rendered deterministically.
	Tags []resolve.Tag
	// Verdict is the gate verdict, empty for non-gate invocations.
	Verdict Verdict
}

func dispositionOf(def Definition) Disposition {
	timeout, _ := def.timeout()
	return Disposition{
		Accessor:   def.Identity.Name,
		Capability: def.Identity.Capability,
		Role:       def.Accessor.Role,
		Timeout:    timeout,
	}
}

// ReadDisposition renders a read result as a replay-stable disposition.
func ReadDisposition(def Definition, r ReadResult) Disposition {
	d := dispositionOf(def)
	if r.Refusal != nil {
		d.Refusal = r.Refusal.Class
		return d
	}
	d.Tags = r.OwnedSnapshot()
	return d
}

// GateDisposition renders a gate result as a replay-stable disposition.
func GateDisposition(def Definition, r GateResult) Disposition {
	d := dispositionOf(def)
	d.Verdict = r.Verdict
	if r.Refusal != nil {
		d.Refusal = r.Refusal.Class
	}
	return d
}

// WriteDisposition renders a write result as a replay-stable disposition.
func WriteDisposition(def Definition, r WriteResult) Disposition {
	d := dispositionOf(def)
	if r.Refusal != nil {
		d.Refusal = r.Refusal.Class
		return d
	}
	d.Tags = observedTags(tagsAsValues(r.Written))
	return d
}

func tagsAsValues(tags []resolve.Tag) []KeyValue {
	out := make([]KeyValue, 0, len(tags))
	for _, t := range tags {
		out = append(out, KeyValue{Key: t.Key, Value: t.Value})
	}
	return out
}
