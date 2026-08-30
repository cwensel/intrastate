package cli

// RDR 0023 `0023:C1` — the `--plan-only` projection: the one flag-aware
// step on `flow resolve`'s success path.
//
// The contract this file implements is narrow and its shape is normative,
// not an implementation-plan detail:
//
//   - The projection SITE is "applied to the verb-specific result BEFORE
//     respond.OK, never in the respond gateway and never per output mode",
//     and "it sits on the SUCCESS path only, AFTER the last respond.Fail
//     return". So `runFlowResolve` calls this immediately before handing
//     the result to the gateway, and nothing here knows about output modes.
//
//   - The flag is read at EXACTLY ONE lexical site — the branch in
//     `projectResolvePayload` below. It is never passed into payload
//     assembly, and never consulted in gate evaluation, rule selection, or
//     refusal construction.
//
//   - There is deliberately NO "if refusing, skip projection" guard. C1
//     forbids one: refusal flag-blindness has to hold STRUCTURALLY, and it
//     does, because every refusal returns from `runFlowResolve` before this
//     code is reachable. A defensive guard here would replace a structural
//     property with a runtime one and quietly permit the site to move.
//
//   - Projection is KEY DELETION, never re-encoding. It takes the fully
//     assembled payload and returns it minus the echo group, so every
//     surviving field carries its default-mode bytes through the same
//     encoder in the same declaration order. It cannot compute a value the
//     default width would not have produced, because it computes nothing.
//
// The ECHO set is not decided here. It is read from the DECLARED partition
// in flow_partition.go — a separate site this file's edits do not open.
// That separation is what makes the reflective oracle able to fail: an
// omit-list authored here would be restated by any oracle that read it, and
// would agree with itself no matter which fields it dropped.

import (
	"reflect"

	"github.com/spf13/cobra"
)

// planOnlyFlagName is the flag's normative spelling (`0023:D-naming`):
// boolean, long form only, default false. `--select`, `--fields`,
// `--quiet`, `--no-echo`, `--short`/`--brief` and an inverted `--explain`
// are all named REJECTED spellings for this axis.
const planOnlyFlagName = "plan-only"

// planOnlyFlagUsage is the flag's authored usage string.
//
// It is authored ONCE, here, because it ships into the generated
// `docs/cli-reference.md` and `llms.txt` rather than staying an
// implementation detail — and it NAMES THE OMITTED GROUP, which is the
// mitigation for the one real risk the narrowing carries: `model` is
// exactly the field a projected caller loses, so a caller must be able to
// learn that from the flag itself rather than by diffing two runs.
const planOnlyFlagUsage = "omit the request echo from a successful plan " +
	"(model, observed, owned, readers, outcome); the plan itself is unchanged"

// projectResolvePayload returns the payload `respond.OK` should render.
//
// This is the ONLY site that reads the flag. When it is unset the payload
// is returned untouched, which is why default-mode output is byte-identical
// to the pre-change binary. When it is set, the ECHO group's fields are
// nilled — and because `0023:A2` retyped exactly those five to `*T` with
// `,omitempty`, a nilled field DROPS ITS KEY rather than rendering `null`
// or an empty placeholder. Absent means absent.
//
// The function takes the FULLY ASSEMBLED payload and returns the projected
// one. Fusing it into assembly — building a narrower payload under the flag
// — is forbidden, and for a checkable reason: it would make "every carried
// field is byte-identical to its default rendering" untestable, since there
// would no longer be a default rendering to compare against.
func projectResolvePayload(cmd *cobra.Command, payload resolvePayload) resolvePayload {
	if planOnly, _ := cmd.Flags().GetBool(planOnlyFlagName); !planOnly {
		return payload
	}
	return projectAwayEchoGroup(payload)
}

// projectAwayEchoGroup deletes the declared ECHO group from an assembled
// payload, leaving every PLAN field exactly as assembled.
//
// It is driven by `echoFieldNames()` — the DECLARATION in
// flow_partition.go — rather than by a list written here. That is the whole
// point of the split: a field added to the struct and assigned `echo` there
// is projected away without this file changing, and a field added with NO
// entry is caught by the reflective completeness oracle instead of
// defaulting silently into either width.
//
// Deletion is by setting the pointer to nil, which is key deletion under
// `,omitempty`. It rewrites no carried value, so key ORDER, encoder, and
// every surviving byte are the default's.
func projectAwayEchoGroup(payload resolvePayload) resolvePayload {
	projected := reflect.ValueOf(&payload).Elem()
	for _, name := range echoFieldNames() {
		field := projected.FieldByName(name)
		if !field.IsValid() || field.Kind() != reflect.Pointer {
			// A field the partition calls `echo` that is not
			// pointer-valued cannot be dropped by this mechanism, and
			// zeroing it instead would emit a stand-in rather than an
			// absence. Leaving it carried is the honest failure: the
			// projected key set then differs from the declaration's plan
			// side and the reflective oracle says so, loudly, instead of
			// this code inventing an empty placeholder the contract
			// forbids.
			continue
		}
		field.SetZero()
	}
	return payload
}
