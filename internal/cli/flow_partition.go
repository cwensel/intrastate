package cli

// RDR 0023 `0023:C2` — the DECLARED ECHO/PLAN partition of `flow resolve`'s
// success payload, and nothing else.
//
// This file is a DECLARATION SITE, deliberately separate from the file that
// applies the projection. `0023:C2` fixes both halves of that separation:
//
//	"The carrier is CONSTRAINED, not free: it MUST be a standalone table
//	 keyed by field name, in its own declaration site, NOT a per-field
//	 marker on `resolvePayload`."
//
//	"Independence has to be structural — a separate site a projection edit
//	 does not open — or S3 asserts only that the implementation agrees with
//	 itself."
//
// So the assignment lives here and the projection lives in
// flow_projection.go. A marker carried on `resolvePayload` itself would be
// the carrier the clause forbids for a concrete reason: `0023:A2` retypes
// those very fields, so a marker there is edited in the same commit that
// changes the projection — moving the assertion and its subject together
// and leaving the reflective oracle asserting only that the implementation
// agrees with itself.
//
// This file names no flag and reads no flag. It states WHICH SIDE each
// field is on; whether a given run narrows to one side is the projection's
// question, asked once, elsewhere.

import "reflect"

// payloadGroup is the two-valued axis JDR 0002 §D1 fixes for a projectable
// verb payload. There is no third value: `0023:C2`'s reflective oracle
// asserts every field is assigned to EXACTLY ONE group, so an unassigned
// new field is a test failure rather than a silent default.
type payloadGroup string

const (
	// groupEcho is the request the caller already holds, read back: "the
	// model reference, the observed tags (--tag, echoed unchanged), the
	// assembled owned view and the invoked reader identities, and the
	// requested outcome" (`0023:C2`).
	groupEcho payloadGroup = "echo"

	// groupPlan is what the run DECIDED: "rule identity, gate results,
	// authored answers and their interpretations, planned
	// next/writes/clear, the escape disposition, and revision"
	// (`0023:C2`).
	groupPlan payloadGroup = "plan"
)

// resolvePayloadGroups is the declared assignment: one entry per field of
// `resolvePayload`, keyed by the Go FIELD NAME, carrying `echo` or `plan`.
//
// The per-field assignment is the Source-authority census's
// (`0023:CEN`), transcribed once:
//
//	model, observed, owned, readers, outcome  ECHO
//	revision                                  PLAN  (identity slot)
//	rule, gates, emit, next, writes, clear,
//	escaped, escape_class                     PLAN
//
// Two entries carry more weight than the rest and are called out because a
// later reader will be tempted to move them:
//
//   - `Revision` rides the PLAN side deliberately. It is produced by the
//     loader, never restated from the request, so it is not an echo of
//     anything the caller sent — and carrying it keeps the wire shape
//     stable for the day RDR 0002 admits the key. It is constant-empty
//     today; that is a fact about 0002, not about this partition.
//
//   - `EscapeClass` is PLAN and is always-keep core, even though it is a
//     plain `string` with `,omitempty` rather than a pointer. Group
//     membership is THIS TABLE's, never a Go type's: an implementation that
//     read "pointer + omitempty" as the mark of a projectable field would
//     sweep `escape_class` away with the echo group and drop the very
//     disposition that says a plan was RESCUED.
//
// This verb's ALWAYS-KEEP core (JDR 0002 §D1) is `rule`, `escaped`,
// `escape_class`, `revision`: if the axis ever generalizes past one
// boolean, no mode may project those away. Today every mode keeps every
// plan field, so the core needs no separate enforcement — but the
// membership above is where a successor reads it from.
var resolvePayloadGroups = map[string]payloadGroup{
	"Model":    groupEcho,
	"Observed": groupEcho,
	"Owned":    groupEcho,
	"Readers":  groupEcho,
	"Outcome":  groupEcho,

	"Revision":    groupPlan,
	"Rule":        groupPlan,
	"Gates":       groupPlan,
	"Emit":        groupPlan,
	"Next":        groupPlan,
	"Writes":      groupPlan,
	"Clear":       groupPlan,
	"Escaped":     groupPlan,
	"EscapeClass": groupPlan,
}

// echoFieldNames returns the `resolvePayload` field names this table
// assigns to the ECHO group, in the struct's own declaration order.
//
// Declaration order is what makes the caller's deletion loop stable and
// keeps the surviving keys in "their default-mode relative order": the
// projection is the default payload MINUS these names, never a re-sort.
//
// It is derived from the table by REFLECTION over the struct rather than
// written out a second time, so the two cannot drift: a field renamed
// without its entry renamed simply has no entry, and the reflective
// completeness oracle fails at its first assertion rather than here.
func echoFieldNames() []string {
	rt := reflect.TypeOf(resolvePayload{})
	out := make([]string, 0, rt.NumField())
	for i := range rt.NumField() {
		name := rt.Field(i).Name
		if resolvePayloadGroups[name] == groupEcho {
			out = append(out, name)
		}
	}
	return out
}
