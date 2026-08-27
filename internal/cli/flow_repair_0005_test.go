package cli

// RDR 0005 — the REPAIR path after a read-back that could not complete.
//
// REQ-104 lets a write refuse at exit 3 with "the mutation may have been
// applied; inspect the artifact". That disposition is about the invocation
// that refused, and REQ-110 fixes request identity over the model revision
// and the CURRENT artifact contents — so nothing makes the refusal STICKY.
// A caller who corrects the model and writes again through a reachable
// accessor must CONVERGE on an artifact `read-state` can report, or the
// CLI has a reachable state with no way out of it: RDR 0002's per-accessor
// `path` is AUTHORED, so a mistyped write locator is a user input, not a
// test-only fixture.
//
// This oracle goes through the CLI in both directions, like every other
// state oracle in this suite: nothing here reads the artifact file, so it
// holds for whatever on-disk format the implementation chooses.

import (
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
)

// REQ-104 / REQ-110 — DOMAIN EDGE: a write whose read-back could not
// complete leaves the artifact REPAIRABLE. Reachable writes converge on a
// verified artifact; they do not refuse forever.
func TestReq104_ACorrectedWriterRepairsAnUnverifiableArtifact(t *testing.T) {
	// The two models differ ONLY in the write accessor's declared locator:
	// `flowReadBackFailModel` names one whose read-back cannot complete,
	// `flowMVVModel` names a reachable one. That is exactly the correction
	// a model author makes after taking the REQ-104 refusal.
	broken := writeFlowModel(t, flowReadBackFailModel)
	fixed := writeFlowModel(t, flowMVVModel)

	art := newFlowArtifact(t, "state.artifact")
	bind := artifactBinding(flowStateRole, art)
	// The orphan reader is bound to its OWN artifact so that role's
	// unrelated key set cannot color this oracle's read (`0005:MVV`).
	orphanBind := artifactBinding(flowOrphanRole,
		newFlowArtifact(t, "orphan.artifact"))

	// 1. The mistyped locator. This MUST refuse at exit 3 with the
	//    read-back-incomplete code specifically — a different failure here
	//    would mean the rest of this test never reaches the repair it is
	//    about.
	_, _, err := runCmd(t, "flow", "set-state", "--model", broken,
		"--artifact", bind, "--write", "status=final", "--as=json")
	if err == nil {
		t.Fatal("the write through an unreachable read-back locator " +
			"SUCCEEDED; REQ-104 makes it exit 3, and without that refusal " +
			"there is nothing to repair")
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) || ce.Code != codeReadBackIncomplete {
		t.Fatalf("code = %v; want %q — this oracle needs the "+
			"read-back-incomplete refusal specifically, not any failure",
			err, codeReadBackIncomplete)
	}

	// 2. The author corrects the model and writes again, REPEATEDLY. Each
	//    such write is reachable, so the artifact must converge on verified
	//    within a bounded number of them. A monotonic seal never converges:
	//    every re-read is refused because a PAST write — not the last one —
	//    declared an unreachable read-back, so this loop would exhaust.
	//
	//    Convergence, not one-shot success, is the claim. The FIRST
	//    corrected write can still refuse honestly: its own PRE-write
	//    baseline of the reader's protected keys was taken while the
	//    artifact was still unverifiable, and `0004:C13` makes an
	//    unestablished baseline `read_back_incomplete` too. That refusal is
	//    about that invocation's own snapshot and does not survive it.
	const attempts = 3
	var repaired bool
	for i := range attempts {
		_, _, err := runCmd(t, "flow", "set-state", "--model", fixed,
			"--artifact", bind, "--write", "status=final", "--as=json")
		if err == nil {
			repaired = true
			break
		}
		if !asCLIError(err, &ce) || ce.Code != codeReadBackIncomplete {
			t.Fatalf("corrected write %d refused %v; a REACHABLE write may "+
				"only ever refuse for its own unverifiable read-back, "+
				"never for some other reason", i+1, err)
		}
	}
	if !repaired {
		t.Fatalf("%d successive writes through a REACHABLE accessor all "+
			"refused; correcting the locator must eventually restore a "+
			"verifiable artifact, or the earlier REQ-104 refusal has "+
			"poisoned it permanently with no repair path through the CLI",
			attempts)
	}

	// 3. And the value is now observable. `read-state` refusing
	//    `flow-read-incomplete` here would be the monotonic-seal defect:
	//    every key unreadable forever, however many correct writes follow.
	data := flowData(t, requireSuccess(t, "flow", "read-state",
		"--model", fixed, "--artifact", bind, "--artifact", orphanBind,
		"--as=json"))
	if got := readerTagValue(t, data, flowStateRole, "status"); got != "final" {
		t.Errorf("status = %#v, want %q; the repairing write's value must "+
			"be what `read-state` reports", got, "final")
	}
}
