package cli

// Kata kf00, envelope half — a non-advancing row (`advance = false`)
// resolves with `writes: {}` and an EMPTY `next`.
//
// The kata's whole point is a state-machine row that refuses without
// applying, so the plan it hands `flow set-state --plan -` must name no
// write. The self-write workaround the kata rejected — assigning a tag the
// value it already holds — produced a REFUSAL carrying a non-empty write
// plan, which the executor then performed and reported as a successful
// apply. Nothing in the envelope may leave that door open through this
// shape either, so both fields are asserted directly rather than assumed
// from the missing write block.

import (
	"testing"
)

// kf00MixedModel is the mixed shape the kata's consumer needs: ONE
// state-machine model in which `apply` advances `status` and `refuse`
// emits its reason and stays put.
const kf00MixedModel = `outcomes = ["apply", "refuse"]
terminal = ["done"]

[model]
id = "kf00mixed"
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
path = "flow.state"
keys = ["status"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"
read_back = true

[context.done]
[context.done.match.status]
eq = "final"

[emit.op]
kind = "enum"
[emit.op.domain]
edit = ["applied"]
stop = ["stopped:no-gate"]

[[rule]]
id = "apply"
[rule.match.recognized]
eq = "apply"
[rule.match.status]
eq = "draft"
[rule.write]
status = "final"
[rule.emit]
op = "applied"

[[rule]]
id = "refuse"
advance = false
[rule.match.recognized]
eq = "refuse"
[rule.match.status]
eq = "draft"
[rule.emit]
op = "stopped:no-gate"
`

// resolveKf00 resolves one outcome over the mixed model against a seeded
// `status = draft` artifact and returns the decoded payload.
func resolveKf00(t *testing.T, outcome string) map[string]any {
	t.Helper()

	model := writeFlowModel(t, kf00MixedModel)
	art := seedArtifact(t, model, "status=draft")
	return flowData(t, requireSuccess(t,
		"flow", "resolve",
		"--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--outcome", outcome,
		"--as=json"))
}

// TestKataKf00_ANonAdvancingRowResolvesWithNoPlannedWrite is the GREEN
// assertion at the envelope: `writes` and `next` are both present and
// EMPTY, so `flow set-state --plan -` has nothing to apply.
func TestKataKf00_ANonAdvancingRowResolvesWithNoPlannedWrite(t *testing.T) {
	data := resolveKf00(t, "refuse")

	for _, key := range []string{"writes", "next"} {
		raw, present := data[key]
		if !present {
			t.Errorf("the payload carries no %q key; it keeps its `0005:C1` "+
				"object shape — never omitted. keys = %v", key, keysOf(data))
			continue
		}
		obj, ok := raw.(map[string]any)
		if !ok {
			t.Errorf("`%s` = %#v; it keeps its `0005:C1` object shape", key, raw)
			continue
		}
		if len(obj) != 0 {
			t.Errorf("`%s` = %#v; a non-advancing row plans NO write — the "+
				"self-write workaround this marker replaces is exactly the "+
				"refusal that carried one", key, obj)
		}
	}

	// The refusal is still an ANSWER: the emit and its declared disposition
	// are what the caller branches on, and they must survive the row
	// carrying no write.
	if got := dispositionsOf(t, data); got["op"] != "stop" {
		t.Errorf("dispositions = %v; want op=stop — the emit is what makes "+
			"declining to advance an outcome rather than a silent no-op", got)
	}
}

// TestKataKf00_TheSameModelStillPlansTheAdvancingRowsWrite is the control:
// the marker is PER ROW, so the model's advancing rule is untouched. One
// model both applies and refuses, which is the capability the kata names.
func TestKataKf00_TheSameModelStillPlansTheAdvancingRowsWrite(t *testing.T) {
	data := resolveKf00(t, "apply")

	writes, ok := data["writes"].(map[string]any)
	if !ok {
		t.Fatalf("`writes` = %#v; want an object", data["writes"])
	}
	if writes["status"] != "final" {
		t.Errorf("writes = %#v; want status=final — `advance = false` on a "+
			"SIBLING row must not disarm this one", writes)
	}
}
