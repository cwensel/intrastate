package cli

// RDR 0005 — `flow set-state` (REQ-60, REQ-64..REQ-68, REQ-78),
// `flow read-state` (REQ-58, REQ-77), the read-back-gated success rule
// (REQ-65, REQ-103..REQ-105), and the round-trip / inverse invariants
// (REQ-107..REQ-109).
//
// Every state oracle here goes through the CLI in both directions:
// `set-state` establishes, `read-state` observes. Nothing reads the
// artifact file, so these assertions hold for whatever on-disk format the
// implementation chooses — and they still discriminate, because a write
// that did not persist, or persisted wrong, changes what `read-state`
// reports.

import (
	"slices"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
)

// REQ-60: "flow set-state MUST invoke only declared write accessors over
// caller-supplied role=path artifact bindings and the planned owned-tag
// mutations given as --write name=value and --clear <key>."
// REQ-78: `set-state` data minimum — "`model`, `revision`, artifact role
// bindings, `writers[]`, requested `writes` and `clear[]`, and the
// read-back-confirmed owned-tag values."
// HAPPY PATH
func TestReq60And78_SetStatePayloadCarriesItsFullDataMinimum(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft", "stale=obsolete")

	data := flowData(t, requireSuccess(t, "flow", "set-state",
		"--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--write", "status=final",
		"--clear", "stale", "--as=json"))

	for _, key := range []string{
		"model", "revision", "artifacts", "writers", "writes", "clear", "owned",
	} {
		if _, ok := data[key]; !ok {
			t.Errorf("`set-state` payload carries no %q; the data minimum "+
				"fixes it. keys = %v", key, keysOf(data))
		}
	}

	// `writers[]` names the write accessor identities invoked — and ONLY
	// write accessors: a reader or a gate appearing here would mean
	// set-state invoked a capability it may not.
	if ids, ok := stringsAt(data, "writers"); ok {
		if !slices.Contains(ids, "state") {
			t.Errorf("`writers[]` = %v omits `state`", ids)
		}
		if slices.Contains(ids, "approval") {
			t.Error("`writers[]` names the GATE `approval`; set-state " +
				"invokes only declared WRITE accessors")
		}
	}

	// The read-back-confirmed owned values, not merely the request echoed.
	owned, ok := data["owned"].(map[string]any)
	if !ok {
		t.Fatalf("`owned` is not an object: %#v", data["owned"])
	}
	if owned["status"] != "final" {
		t.Errorf("owned[status] = %#v; read-back must confirm the planned "+
			"value %q", owned["status"], "final")
	}
	// REQ-66 / `0004:C11`: a CLEARED key reads back ABSENT — it is never
	// reported as holding the sentinel or an empty string.
	if v, held := owned["stale"]; held {
		t.Errorf("owned[stale] = %#v; a cleared key reads back ABSENT and a "+
			"verified removal is not a written value", v)
	}
}

// REQ-64: "It MUST never run gate accessors." — restated in AP: "It never
// runs gates".
// `advance-draft` row, and a `set-state` that ran gates would surface a
// gate verdict or a gate refusal. Neither may appear.
// ADVERSARIAL — the observable proof: the MVV model's gate is on the
func TestReq64_SetStateNeverRunsGateAccessors(t *testing.T) {
	// This model's ONLY gate cannot be consulted (REQ-47's fixture). If
	// `set-state` ran gates, it would refuse exit 3 instead of succeeding.
	model := writeFlowModel(t, flowGateFailModel)
	art := seedArtifact(t, model, "status=draft")

	stdout, _, err := runCmd(t, "flow", "set-state", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--write", "status=final", "--as=json")
	if err != nil {
		var ce *clierr.CLIError
		if asCLIError(err, &ce) && strings.HasPrefix(ce.Code, "flow-accessor-") {
			t.Fatalf("`set-state` refused %q, an ACCESSOR refusal raised by "+
				"the model's unconsultable GATE; set-state MUST NEVER run "+
				"gate accessors", ce.Code)
		}
		t.Fatalf("`set-state` failed: %v", err)
	}

	data := flowData(t, stdout)
	if _, present := data["gates"]; present {
		t.Error("`set-state` reported a `gates` field; it never runs gates")
	}
	if ids, ok := stringsAt(data, "writers"); ok {
		if slices.Contains(ids, "unreachable") {
			t.Error("`writers[]` names the gate accessor `unreachable`")
		}
	}
}

// REQ-65: "It MUST report success only after the accessor layer's read-back
// verification confirms the planned owned-tag values and that non-owned tags
// are unchanged."
// REQ-103: `flow-write-readback-mismatch` — "read-back disagrees with the
// plan, or a non-owned tag changed" — `GroupUserEnv` / 2 — "`findings[]`
// per key".
// SEPARATE read-state call rather than by trusting the write's own echo.
// DOMAIN EDGE — success implies the value is really there, proved by a
func TestReq65_SuccessImpliesTheValueIsActuallyReadableAfterwards(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")

	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--write", "status=final", "--as=json")

	// A fresh, independent read must observe the value. A `set-state` that
	// reported success without the read-back holding fails here.
	data := flowData(t, requireSuccess(t, "flow", "read-state",
		"--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--artifact", artifactBinding(flowOrphanRole, art), "--as=json"))

	if got := readerTagValue(t, data, "state", "status"); got != "final" {
		t.Errorf("after a successful `set-state`, `read-state` reports "+
			"status = %#v; success is reported ONLY after read-back "+
			"confirms the planned value", got)
	}
}

// REQ-66: "`[]` (empty set) and `--clear` (absent) stay distinct through
// read-back."
// REQ-107: "a cleared key reads back absent, an empty set reads back `[]`".
// spellings that a careless implementation collapses into one.
// REQ-127 / `0005:S5`: "`set-state` reports success only after read-back
// proves the scalar, the set (byte-equal canonical array), and the cleared
// key absent".
// BOUNDARY — the single sharpest boundary in the write grammar: two
func TestReq66And107_EmptySetAndClearedKeyStayDistinctThroughReadBack(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)

	t.Run("empty-set-reads-back-as-empty-array", func(t *testing.T) {
		art := seedArtifact(t, model, "status=draft", `labels=["plain"]`)

		requireSuccess(t, "flow", "set-state", "--model", model,
			"--artifact", artifactBinding(flowStateRole, art),
			"--write", "labels=[]", "--as=json")

		data := flowData(t, requireSuccess(t, "flow", "read-state",
			"--model", model,
			"--artifact", artifactBinding(flowStateRole, art),
			"--artifact", artifactBinding(flowOrphanRole, art), "--as=json"))

		raw, held := readerTag(t, data, "state", "labels")
		if !held {
			t.Fatal("`labels` is ABSENT after being written the empty set; " +
				"`[]` and a cleared key are DISTINCT — an empty set is a " +
				"present key holding no members")
		}
		if s, ok := raw.(string); ok {
			if s != "[]" {
				t.Errorf("`labels` reads back %q; an empty set reads back "+
					"the canonical empty array `[]`", s)
			}
			return
		}
		if arr, ok := raw.([]any); ok {
			if len(arr) != 0 {
				t.Errorf("`labels` reads back %#v; want an empty array", arr)
			}
			return
		}
		t.Errorf("`labels` reads back %#v; want the canonical `[]`", raw)
	})

	t.Run("cleared-key-reads-back-absent", func(t *testing.T) {
		art := seedArtifact(t, model, "status=draft", `labels=["plain"]`)

		requireSuccess(t, "flow", "set-state", "--model", model,
			"--artifact", artifactBinding(flowStateRole, art),
			"--clear", "labels", "--as=json")

		data := flowData(t, requireSuccess(t, "flow", "read-state",
			"--model", model,
			"--artifact", artifactBinding(flowStateRole, art),
			"--artifact", artifactBinding(flowOrphanRole, art), "--as=json"))

		if raw, held := readerTag(t, data, "state", "labels"); held {
			t.Errorf("`labels` reads back %#v after `--clear`; a cleared key "+
				"reads back ABSENT, which is what distinguishes it from the "+
				"empty set", raw)
		}
	})
}

// REQ-67: "Each writer applies its own keys and reads back; no cross-writer
// atomicity is promised."
// REQ-68: "`set-state` trusts the caller to transcribe the plan; nothing
// links a `set-state` request to a prior `resolve`."
// request with no prior `resolve` in the session at all.
// BOUNDARY — a stated NON-guarantee: `set-state` accepts a hand-written
func TestReq68_SetStateAcceptsAnUnlinkedRequestWithNoPriorResolve(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")

	// No `resolve` runs first, and no plan token is carried. The request
	// stands on its own grammar.
	data := flowData(t, requireSuccess(t, "flow", "set-state",
		"--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--write", "status=final", "--as=json"))

	if owned, ok := data["owned"].(map[string]any); ok {
		if owned["status"] != "final" {
			t.Errorf("owned[status] = %#v; want %q — the request is accepted "+
				"on its own terms", owned["status"], "final")
		}
	}
}

// REQ-58: "flow read-state MUST invoke only declared read accessors over
// caller-supplied role=path artifact bindings and return, per reader, its
// declared keys and the tag-set read, or a stable CLIError failure."
// REQ-77: `read-state` data minimum — "`readers[]` each with its declared
// `keys` (the requested key set) and the tags it returned — so \"absent from
// the artifact\" and \"not requested\" are distinguishable from the payload
// alone."
// DOMAIN EDGE — the distinguishability property is the whole point.
func TestReq58And77_ReadStateReportsDeclaredKeysBesideTheTagsReturned(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	// `status` is established; `stale` is NOT, so it is a declared key that
	// is ABSENT from the artifact. `note` is declared on a DIFFERENT reader.
	art := seedArtifact(t, model, "status=draft")

	data := flowData(t, requireSuccess(t, "flow", "read-state",
		"--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--artifact", artifactBinding(flowOrphanRole, art), "--as=json"))

	readers, ok := objectsAt(data, "readers")
	if !ok {
		t.Fatalf("`readers` is not an array of objects: %#v", data["readers"])
	}

	var checked bool
	for _, r := range readers {
		if r["id"] != "state" {
			continue
		}
		checked = true
		keys, ok := stringsAt(r, "keys")
		if !ok {
			t.Fatalf("reader `state` carries no `keys`; the declared key set "+
				"rides beside the tags. got %#v", r["keys"])
		}
		// The DECLARED key set, verbatim from `[read.state].keys`.
		for _, want := range []string{"status", "labels", "stale"} {
			if !slices.Contains(keys, want) {
				t.Errorf("reader `state`.keys = %v omits the declared key "+
					"%q", keys, want)
			}
		}
		// `note` belongs to the OTHER reader: "not requested" here.
		if slices.Contains(keys, "note") {
			t.Errorf("reader `state`.keys = %v names `note`, which "+
				"`[read.orphan]` declares — the key set is per-reader", keys)
		}

		tags, ok := r["tags"].(map[string]any)
		if !ok {
			t.Fatalf("reader `state` carries no `tags` object: %#v", r["tags"])
		}
		// `status` was requested AND present.
		if tags["status"] != "draft" {
			t.Errorf("reader `state`.tags[status] = %#v; want %q",
				tags["status"], "draft")
		}
		// `stale` was requested and is ABSENT from the artifact — which the
		// payload distinguishes from "not requested" because `stale`
		// appears in `keys` while carrying no tag value.
		if _, held := tags["stale"]; held {
			t.Errorf("reader `state`.tags reports a value for `stale`, which " +
				"nothing established; absent-from-the-artifact is not a value")
		}
	}
	if !checked {
		t.Fatal("`read-state` reported no reader `state`")
	}
}

// REQ-107: "`set-state` followed by `read-state` must return the planned
// owned-tag values for the artifact role that was written. The equality is
// value-for-value over the owned tags the write planned to mutate — byte
// equality for set values in canonical form."
// HAPPY PATH — the inverse invariant, across all three value shapes.
func TestReq107_SetStateThenReadStateReturnsThePlannedValues(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft", "stale=obsolete")

	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--write", "status=final",
		"--write", "labels="+canonicalSetLiteral,
		"--clear", "stale", "--as=json")

	data := flowData(t, requireSuccess(t, "flow", "read-state",
		"--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--artifact", artifactBinding(flowOrphanRole, art), "--as=json"))

	if got := readerTagValue(t, data, "state", "status"); got != "final" {
		t.Errorf("status reads back %#v; want %q", got, "final")
	}
	assertCanonicalSetValue(t, "labels", readerTagRaw(t, data, "state", "labels"))
	if _, held := readerTag(t, data, "state", "stale"); held {
		t.Error("`stale` reads back present after `--clear`; a cleared key " +
			"reads back absent")
	}
}

// REQ-108: "`resolve` → skill → `set-state` is copy-through: every `writes`
// value and `clear[]` key in a plan is accepted verbatim by `set-state`'s
// `--write` / `--clear` grammar, so a set write survives the round trip
// byte-identical."
// DOMAIN EDGE — the plan's OWN emitted values are fed straight back in.
func TestReq108_PlanValuesAreAcceptedVerbatimBySetStatesGrammar(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft", "stale=obsolete")

	plan := flowData(t, requireSuccess(t, "flow", "resolve", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--outcome", "advance", "--as=json"))

	writes, ok := plan["writes"].(map[string]any)
	if !ok {
		t.Fatalf("the plan's `writes` is not an object: %#v", plan["writes"])
	}
	clears, _ := stringsAt(plan, "clear")

	// Transcribe the plan into a `set-state` request VERBATIM — exactly what
	// a skill does between the two calls.
	args := []string{"flow", "set-state", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art)}
	for key, raw := range writes {
		args = append(args, "--write", key+"="+planValueLiteral(t, raw))
	}
	for _, key := range clears {
		args = append(args, "--clear", key)
	}

	if _, _, err := runCmd(t, append(args, "--as=json")...); err != nil {
		t.Fatalf("the plan's own values were REFUSED by `set-state`'s "+
			"grammar: %v\nCopy-through means every `writes` value and "+
			"`clear[]` key a plan emits is accepted verbatim", err)
	}
}

// REQ-109: "`read-state` → `--tag` re-pairs only for observed tags; an owned
// tag read by `read-state` cannot be handed back through `--tag` (refused
// `flow-tag-owned`) because `resolve` reads it itself."
// ADVERSARIAL — the asymmetry is deliberate and must be enforced.
func TestReq109_AnOwnedTagReadByReadStateCannotBeHandedBackThroughTag(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")

	// `read-state` reports `status`…
	data := flowData(t, requireSuccess(t, "flow", "read-state",
		"--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--artifact", artifactBinding(flowOrphanRole, art), "--as=json"))
	if got := readerTagValue(t, data, "state", "status"); got != "draft" {
		t.Fatalf("read-state reports status = %#v; want %q", got, "draft")
	}

	// …and handing it straight back through `--tag` is refused.
	requireRefusal(t, "flow-tag-owned", 2,
		"flow", "resolve", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--outcome", "hold", "--tag", "status=draft", "--as=json")

	// An OBSERVED tag re-pairs fine — so the refusal is about provenance,
	// not about `--tag` being generally unusable after a read.
	requireSuccess(t, "flow", "resolve", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--outcome", "hold", "--tag", "profile=mid", "--as=json")
}

// REQ-104: `flow-write-readback-incomplete` — "read-back returned an
// incomplete key set" — `GroupEnvUnavailable` / 3 — "`detail`: may have been
// applied".
// REQ-105: `flow-write-readback-timeout` — "post-mutation read-back timed
// out" — `GroupEnvUnavailable` / 3 — same detail.
// REQ-20: the exit-3 population is exactly five classes, two of which are
// these.
// refusal actionable, so its ABSENCE is a defect.
// REQ-128 / `0005:S6`: "the first three exit 3 with `GroupEnvUnavailable`
// and a \"may have been applied\" detail on the write case; the mismatch
// exits 2 with `findings[]` per key; none is a successful transition or a
// coerced `read-state` payload".
// DOMAIN EDGE — the "may have been applied" detail is what makes the
func TestReq104And105_ReadBackFailuresExit3AndSayTheWriteMayHaveApplied(t *testing.T) {
	model := writeFlowModel(t, flowReadBackFailModel)
	// The artifact is created EMPTY rather than seeded through
	// `seedArtifact`, which drives `flow set-state` against this same model
	// and would `t.Fatal` on the very refusal under test. REQ-110 fixes a
	// refusal identity by its inputs, so a seeding call over the same model
	// and the same writer refuses identically — the seed and the assertion
	// cannot both hold, and the seed fired first (DEV-5).
	//
	// The seed was never load-bearing: this writer's read-back is
	// unreachable by its declared locator, so the refusal follows the
	// MUTATION regardless of what the artifact held beforehand. Dropping it
	// is what lets the REQ-104/REQ-105 oracle below actually be reached.
	art := newFlowArtifact(t, "state.artifact")

	_, _, err := runCmd(t, "flow", "set-state", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--write", "status=final", "--as=json")
	if err == nil {
		t.Fatal("a write whose read-back could not complete reported " +
			"SUCCESS; success is reported ONLY after read-back confirms")
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("not a CLIError: %v", err)
	}
	if !slices.Contains([]string{
		"flow-write-readback-incomplete", "flow-write-readback-timeout",
	}, ce.Code) {
		t.Fatalf("code = %q; a read-back that did not COMPLETE is "+
			"`flow-write-readback-incomplete` or "+
			"`flow-write-readback-timeout` — it is NOT "+
			"`flow-write-readback-mismatch`, which asserts the artifact is "+
			"wrong", ce.Code)
	}
	if got := clierr.ExitCodeFor(err); got != 3 {
		t.Errorf("exit = %d; want 3 — the environment could not be consulted",
			got)
	}
	if ce.Detail == "" {
		t.Error("the refusal carries no `detail`; the code table fixes " +
			"\"may have been applied\" as its carrier, and without it the " +
			"caller cannot tell whether to retry or to inspect")
	}
}

// --- payload readers -----------------------------------------------------
//
// `read-state` reports per-reader tags. These helpers read them without
// presuming a rendering beyond the data minimum REQ-77 fixes.

func readerTag(t *testing.T, data map[string]any, reader, key string) (any, bool) {
	t.Helper()

	readers, ok := objectsAt(data, "readers")
	if !ok {
		t.Fatalf("`readers` is not an array of objects: %#v", data["readers"])
	}
	for _, r := range readers {
		if r["id"] != reader {
			continue
		}
		tags, ok := r["tags"].(map[string]any)
		if !ok {
			t.Fatalf("reader %q carries no `tags` object: %#v", reader,
				r["tags"])
		}
		v, held := tags[key]
		return v, held
	}
	t.Fatalf("`read-state` reported no reader %q", reader)
	return nil, false
}

func readerTagRaw(t *testing.T, data map[string]any, reader, key string) any {
	t.Helper()

	v, held := readerTag(t, data, reader, key)
	if !held {
		t.Fatalf("reader %q reports no value for %q", reader, key)
	}
	return v
}

func readerTagValue(t *testing.T, data map[string]any, reader, key string) any {
	t.Helper()

	v, _ := readerTag(t, data, reader, key)
	return v
}

// planValueLiteral renders a plan's emitted write value back into the
// `--write name=value` grammar. A set value is already the canonical JSON
// array literal on the wire, so it is passed through UNCHANGED — that
// pass-through is what "accepted verbatim" means (REQ-108).
func planValueLiteral(t *testing.T, raw any) string {
	t.Helper()

	switch v := raw.(type) {
	case string:
		return v
	case []any:
		members := make([]string, 0, len(v))
		for _, m := range v {
			s, ok := m.(string)
			if !ok {
				t.Fatalf("a set member is not a string: %#v", m)
			}
			members = append(members, `"`+s+`"`)
		}
		return "[" + strings.Join(members, ",") + "]"
	default:
		t.Fatalf("a plan write value is neither a scalar nor a set: %#v", raw)
		return ""
	}
}
