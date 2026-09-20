package cli

// A reader's returned values enter the resolver's OWNED snapshot only when
// the MODEL declares the key owned.
//
// RDR 0016 legalizes the premise this file probes: reader cardinality is
// fixed per KEY — "exactly one reader per owned key, at most one per
// observed key" — so ONE reader's declared `keys` may legally span an owned
// key and an observed key it also happens to serve. Nothing in RDR 0016
// makes a reader's key set provenance-homogeneous.
//
// The consequence, before the filter below, was a precedence inversion the
// caller had no remedy for. `runReaders` folded EVERY value a reader
// returned into the `Owned` slice; `resolve.assemble` stamps everything in
// `Input.Owned` as `ProvenanceOwned` regardless of the model's own
// `[tags.<key>].Provenance`; and owned beats observed in the kernel's
// precedence (pinned by `resolve/adversarial_test.go`,
// `TestAdv4_ObservedTagsMustNotShadowTheOwnedSnapshot`). So a reader-served
// observed key silently overrode the value the caller supplied with
// `--tag`, and selected a different transition.
//
// `--tag` could not defend itself either: `parseTags` refuses a key only
// when `flowbind.OwnedTags` lists it as OWNED, so it accepts the tag — and
// is then outranked — exactly in this case. RDR 0011's cross-provenance
// analysis is scoped to the inverse direction and does not adjudicate this
// one.
//
// Each oracle below fails for a DIFFERENT defect:
//
//   - the precedence arm fails if any non-owned reader value reaches the
//     resolver, and names the SELECTED RULE rather than a payload echo, so
//     it is a claim about the transition the caller gets;
//   - the echo arm fails if the filter is applied somewhere that also
//     rewrites the `owned` echo's membership incorrectly;
//   - the owned-key arm fails if the filter over-reaches and starves the
//     RDR 0016 single-reader path of a genuinely owned key;
//   - the diagnostic arm fails if the filter leaks into `readerOutput.Tags`,
//     which must keep reporting every key the reader answered, filtered
//     only for established absence.

import (
	"testing"
)

// --- fixture -------------------------------------------------------------

// wmv3SeedModel declares `profile` OWNED so `flow set-state` can establish
// a value for it on the artifact. It exists only to seed: the assertions
// all run against `wmv3SpanModel`, where the same key is OBSERVED.
//
// Seeding this way keeps the artifact's on-disk FORMAT unasserted, the rule
// the whole `flow` suite holds (`flow_harness_0005_test.go`): the CLI's own
// write path establishes the state, and no test here claims to know what a
// tag looks like at rest.
const wmv3SeedModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "wmv3seed"
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

[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["mid", "foundational"]
single_valued = true

[read.state]
role = "state"
path = "flow.state"
keys = ["status", "profile"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status", "profile"]
timeout = "2s"
read_back = true

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "seed-advance"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"
`

// wmv3SpanModel is the fixture the contract turns on: ONE reader whose
// declared `keys` SPAN an owned key (`status`) and an observed key
// (`profile`). RDR 0016 permits exactly this — one reader per owned key,
// at most one per observed key — and `read.state` is the unique server of
// both.
//
// Two rules discriminate on `profile` alone, so the selected rule IS the
// answer to "whose `profile` won": `mid-path` if the caller's `--tag`
// survived, `foundational-path` if the reader's artifact value overrode it.
// The two rules differ on their `[rule.write]` too, so a payload consumer
// reading only the plan sees the divergence as well.
const wmv3SpanModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "wmv3span"
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

[tags.profile]
provenance = "observed"
kind = "enum"
domain = ["mid", "foundational"]
single_valued = true

[read.state]
role = "state"
path = "flow.state"
keys = ["status", "profile"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"
read_back = true

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "mid-path"
[rule.match.status]
eq = "draft"
[rule.match.profile]
eq = "mid"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"

[[rule]]
id = "foundational-path"
[rule.match.status]
eq = "draft"
[rule.match.profile]
eq = "foundational"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "draft"
`

// wmv3Artifact seeds an artifact carrying `status=draft` and
// `profile=foundational`, then returns its `role=path` binding for the SPAN
// model. The seed runs under `wmv3SeedModel`, where `profile` is owned and
// therefore writable; the reads run under `wmv3SpanModel`, where it is not.
func wmv3Artifact(t *testing.T) string {
	t.Helper()

	seed := writeFlowModel(t, wmv3SeedModel)
	art := seedArtifact(t, seed, "status=draft", "profile=foundational")
	return artifactBinding(flowStateRole, art)
}

// --- oracles -------------------------------------------------------------

// The caller's `--tag` on a model-declared OBSERVED key decides the
// transition, even when the invoked reader also answers that key.
func TestOwnedFilter_AReaderServedObservedKeyDoesNotOutrankTheCallersTag(t *testing.T) {
	model := writeFlowModel(t, wmv3SpanModel)
	bind := wmv3Artifact(t)

	stdout := requireSuccess(t, "flow", "resolve", "--model", model,
		"--artifact", bind, "--outcome", "advance",
		"--tag", "profile=mid", "--as=json")
	data := flowData(t, stdout)

	if got, _ := data["rule"].(string); got != "mid-path" {
		t.Fatalf("selected rule = %q; want %q — the reader also answers the "+
			"model-declared OBSERVED key `profile`, and its artifact value "+
			"(`foundational`) must not be promoted to owned and outrank the "+
			"caller's --tag", got, "mid-path")
	}

	// The `owned` echo is the resolver's own view, so the promoted key must
	// be absent from it too. Asserting only on `rule` would pass against an
	// implementation that filtered at the kernel call and still reported a
	// view the caller cannot reconcile with the plan.
	owned, ok := data["owned"].(map[string]any)
	if !ok {
		t.Fatalf("payload carries no `owned` object; keys: %v", keysOf(data))
	}
	if v, present := owned["profile"]; present {
		t.Errorf("`owned` carries profile = %v; an observed-declared key is "+
			"never part of the owned snapshot", v)
	}
	if got, _ := owned["status"].(string); got != "draft" {
		t.Errorf("`owned` carries status = %q; want %q — the filter must not "+
			"starve the genuinely owned key its reader serves", got, "draft")
	}
}

// The same reader's genuinely OWNED key still reaches the resolver: the
// filter narrows by provenance, it does not disable the read path RDR 0016
// fixes.
func TestOwnedFilter_AnOwnedKeyServedByTheSameReaderStillResolves(t *testing.T) {
	model := writeFlowModel(t, wmv3SpanModel)
	bind := wmv3Artifact(t)

	// No `--tag` at all: `profile` is then unsupplied, so the ONLY way a
	// rule can match is through the reader's promoted value. The request
	// must refuse rather than silently resolve `foundational-path`.
	_, _, err := runCmd(t, "flow", "resolve", "--model", model,
		"--artifact", bind, "--outcome", "advance", "--as=json")
	if err == nil {
		t.Fatal("resolve SUCCEEDED with no --tag; the only candidate " +
			"discriminator is `profile`, which no caller supplied, so a " +
			"success means the reader's value was promoted to owned")
	}

	// And with the tag supplied the owned key must still be carried: a
	// `status` that failed to arrive would refuse instead of selecting.
	stdout := requireSuccess(t, "flow", "resolve", "--model", model,
		"--artifact", bind, "--outcome", "advance",
		"--tag", "profile=foundational", "--as=json")
	data := flowData(t, stdout)
	if got, _ := data["rule"].(string); got != "foundational-path" {
		t.Fatalf("selected rule = %q; want %q — the owned `status` read must "+
			"still decide the match arm it serves", got, "foundational-path")
	}
}

// The DIAGNOSTIC path is unfiltered: `flow read-state` reports every key the
// reader answered, whatever the model calls it. The filter belongs to the
// resolver seam alone.
func TestOwnedFilter_ReadStateStillReportsEveryKeyTheReaderAnswered(t *testing.T) {
	model := writeFlowModel(t, wmv3SpanModel)
	bind := wmv3Artifact(t)

	stdout := requireSuccess(t, "flow", "read-state", "--model", model,
		"--artifact", bind, "--as=json")
	data := flowData(t, stdout)

	readers, ok := objectsAt(data, "readers")
	if !ok || len(readers) == 0 {
		t.Fatalf("payload carries no `readers` array; keys: %v", keysOf(data))
	}
	var found bool
	for _, r := range readers {
		if id, _ := r["id"].(string); id != "state" {
			continue
		}
		found = true
		tags, ok := r["tags"].(map[string]any)
		if !ok {
			t.Fatalf("reader `state` carries no `tags` object: %v", r)
		}
		if got, _ := tags["profile"].(string); got != "foundational" {
			t.Errorf("reader `state` reports profile = %q; want %q — the "+
				"diagnostic output is UNFILTERED, so an observed-declared "+
				"key the reader answered is still reported", got,
				"foundational")
		}
		if got, _ := tags["status"].(string); got != "draft" {
			t.Errorf("reader `state` reports status = %q; want %q",
				got, "draft")
		}
	}
	if !found {
		t.Fatalf("no reader `state` in the read-state payload: %v", readers)
	}
}
