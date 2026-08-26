package flowbind

// RDR 0005 — the mutation-persistence guarantee behind REQ-104/REQ-105.
//
// `flow set-state` refuses at exit 3 with "the mutation may have been
// applied; inspect the artifact" when a post-mutation read-back cannot
// complete. That detail is a CLAIM ABOUT THE ARTIFACT, so something must
// pin it: an implementation that discarded the mutation on the unreachable
// path would make the refusal's own detail false while still emitting the
// right code, exit, and detail.
//
// The CLI-level oracle cannot cover this. A sealed artifact reports every
// key UNREADABLE by design, so `flow read-state` cannot observe the
// persisted value, and `flow_setstate_0005_test.go` deliberately reads no
// artifact file so its assertions hold for any on-disk format. The
// guarantee is therefore only observable at this package's own seam, which
// is where this test lives.

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/newcoinc/intrastate/internal/accessor"
	"github.com/newcoinc/intrastate/internal/resolve"
)

// REQ-104: `flow-write-readback-incomplete` — "read-back returned an
// incomplete key set" — `GroupEnvUnavailable` / 3 — "`detail`: may have
// been applied".
// DOMAIN EDGE — the mutation must SURVIVE the unreachable read-back. This
// is the regression that guards the truncating implementation which once
// replaced the artifact's contents with invalid JSON.
func TestReq104_AnUnreachableReadBackKeepsTheMutation(t *testing.T) {
	dir := t.TempDir()
	art := accessor.Artifact{Role: "state", Path: filepath.Join(dir, "a.json")}

	// Establish a prior value through a REACHABLE writer, so the test can
	// tell "kept the mutation" apart from "kept the old artifact".
	reachable := &Writer{Path: "flow.state"}
	if err := reachable.Apply(context.Background(), art, []resolve.Tag{
		{Key: "status", Value: "draft"},
		{Key: "note", Value: "keep-me"},
	}); err != nil {
		t.Fatalf("seeding through a reachable writer failed: %v", err)
	}

	// Now write through a writer whose declared read-back locator cannot
	// be reached. The mutation must land anyway.
	sealed := &Writer{Path: "flow.state.readback-unreachable"}
	if err := sealed.Apply(context.Background(), art, []resolve.Tag{
		{Key: "status", Value: "final"},
	}); err != nil {
		t.Fatalf("Apply returned %v; the write itself must still be "+
			"attempted — the unreachable leg is the VERIFICATION, not the "+
			"mutation (REQ-104)", err)
	}

	s, err := load(art.Path)
	if err != nil {
		t.Fatalf("the artifact is unreadable after the write: %v; the "+
			"refusal's detail claims it MAY HAVE BEEN APPLIED, which is a "+
			"claim about an artifact a caller can still inspect", err)
	}
	if got := s["status"]; got != "final" {
		t.Errorf("status = %q, want %q; the planned write MUST persist "+
			"even though its read-back could not complete, or the "+
			"\"may have been applied\" detail is false (REQ-104)",
			got, "final")
	}
	if got := s["note"]; got != "keep-me" {
		t.Errorf("note = %q, want %q; an unreachable read-back MUST NOT "+
			"discard values the write never planned to touch",
			got, "keep-me")
	}
	if _, sealedNow := s[sealedKey]; !sealedNow {
		t.Error("the artifact carries no seal; the verifying re-read " +
			"refuses on it, which is what makes the write " +
			"applied-but-unverified rather than a mismatch")
	}
}
