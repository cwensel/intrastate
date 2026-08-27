package accessor_test

// RDR 0004 — absence crosses the seam as OMISSION (REQ-30..REQ-34).
//
// This is the cross-layer clause. The accessor's own success branch
// carries absence as a value (REQ-21), but what reaches the RESOLVER must
// leave the key absent from the owned snapshot, or `TagSet.has` reads it
// as present and `owned_state_unavailable` silently retires for that key.
//
// ORA 6 fixes the oracle: the assertion is on the RESOLVER's disposition,
// not the accessor's branch — "a placeholder value makes `TagSet.has` true
// and the refusal never fires" — with a carried key as the control so the
// refusal cannot pass vacuously.
//
// Nothing in the type system prevents the defect (REQ-33):
// `resolve.Input.Owned` is `[]Tag` and `Tag{Key: "profile", Value:
// "<absent>"}` typechecks. The prohibition is carried by these tests.

import (
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/resolve"
)

// REQ-30: "A requested key the artifact genuinely does not carry MUST NOT
// be presented to the resolver as a present owned tag."
// REQ-31: "what crosses the seam MUST leave the key absent from the owned
// snapshot, so that a required absent key resolves as
// `owned_state_unavailable` rather than matching against a placeholder
// value."
// DOMAIN EDGE
func TestReq30_AbsentKeyIsOmittedFromTheOwnedSnapshot(t *testing.T) {
	// `profile` is requested and legitimately not carried.
	s := newStore(map[string]string{keyStatus: "Draft"})
	e, _ := readerExec(t, s, keyStatus, keyProfile)

	got := e.Read(ctxOf(t), readerName)

	mustReadSucceed(t, got)
	snapshot := got.OwnedSnapshot()

	if _, present := seamValueOf(snapshot, keyProfile); present {
		t.Errorf("the owned snapshot carries %q; a genuinely-absent key must be "+
			"OMITTED, never presented to the resolver as a present owned tag "+
			"(snapshot = %+v)", keyProfile, snapshot)
	}
	if got, want := tagKeysOf(snapshot), []string{keyStatus}; !slices.Equal(got, want) {
		t.Errorf("owned snapshot key set = %v; want %v — exactly the keys the "+
			"artifact carries, with the absent key omitted", got, want)
	}
	if v, _ := seamValueOf(snapshot, keyStatus); v != "Draft" {
		t.Errorf("owned snapshot value for %q = %q; want %q — the carried key "+
			"crosses normally", keyStatus, v, "Draft")
	}
}

// REQ-32: "A sentinel value would make an absent key read as *present*
// and silently retire `owned_state_unavailable` for it — the one encoding
// that turns this RDR's safety rule into a regression." — "The spike's
// `<absent>` string is fixture shorthand, not that seam value."
// ADVERSARIAL
//
// The named defect, asserted by name: no seam tag may carry the spike's
// `<absent>` shorthand or any other placeholder for a key the artifact
// does not hold.
func TestReq32_NoPlaceholderValueCrossesTheSeamForAnAbsentKey(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	e, _ := readerExec(t, s, keyStatus, keyProfile)

	snapshot := e.Read(ctxOf(t), readerName).OwnedSnapshot()

	// Positive precondition: the CARRIED key must cross. Without it, an
	// empty snapshot satisfies "no placeholder" vacuously.
	if v, ok := seamValueOf(snapshot, keyStatus); !ok || v != "Draft" {
		t.Fatalf("the snapshot does not carry %q = %q (snapshot = %+v); a "+
			"snapshot that omits EVERY key satisfies the no-placeholder "+
			"assertion vacuously", keyStatus, "Draft", snapshot)
	}

	for _, tag := range snapshot {
		for _, placeholder := range []string{"<absent>", "", "<none>", "<clear>", "null"} {
			if tag.Key == keyProfile && tag.Value == placeholder {
				t.Errorf("the seam carries %q = %q; a placeholder makes "+
					"`TagSet.has` TRUE and silently retires "+
					"owned_state_unavailable for that key", tag.Key, placeholder)
			}
		}
	}
}

// REQ-118 / MVV Scenario 6 (A9): "The read succeeds at the accessor
// boundary, the absent key is omitted from the owned snapshot rather than
// carried as a placeholder value, and the resolver refuses
// `owned_state_unavailable` naming that key."
// DOMAIN EDGE
//
// ORA 6's oracle exactly: the assertion is on the RESOLVER's disposition.
// The control is the second leg — a key the artifact DOES carry must
// resolve normally through the same path, so an implementation that omits
// every key cannot pass the refusal assertion vacuously.
func TestReq118_AbsentRequiredKeyReachesTheResolverAsOwnedStateUnavailable(t *testing.T) {
	t.Run("absent_required_key_refuses_at_the_resolver", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		e, _ := readerExec(t, s, keyStatus, keyProfile)

		read := e.Read(ctxOf(t), readerName)
		mustReadSucceed(t, read)

		result, err := resolve.Resolve(resolve.Input{
			Flow:       flowID,
			Table:      requiringTable(keyProfile),
			Owned:      read.OwnedSnapshot(),
			Recognized: recognizedOutcome,
		})
		if err != nil {
			t.Fatalf("the kernel used the Go error path: %v", err)
		}
		if result.Refusal == nil {
			t.Fatalf("the resolver produced a PLAN over an owned snapshot that "+
				"never carried %q; a placeholder value makes `TagSet.has` true "+
				"and the refusal never fires (plan = %+v)", keyProfile, result.Plan)
		}
		if result.Refusal.Kind != resolve.KindOwnedStateUnavailable {
			t.Fatalf("resolver refusal kind = %q; want %q",
				result.Refusal.Kind, resolve.KindOwnedStateUnavailable)
		}
		if !slices.Contains(result.Refusal.MissingOwned, keyProfile) {
			t.Errorf("MissingOwned = %v; want %q named",
				result.Refusal.MissingOwned, keyProfile)
		}
	})

	t.Run("control_a_carried_key_resolves_normally", func(t *testing.T) {
		// The SAME path, but the artifact carries the required key. An
		// implementation that omits every key fails here.
		s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
		e, _ := readerExec(t, s, keyStatus, keyProfile)

		read := e.Read(ctxOf(t), readerName)
		mustReadSucceed(t, read)

		result, err := resolve.Resolve(resolve.Input{
			Flow:       flowID,
			Table:      requiringTable(keyProfile),
			Owned:      read.OwnedSnapshot(),
			Recognized: recognizedOutcome,
		})
		if err != nil {
			t.Fatalf("the kernel used the Go error path: %v", err)
		}
		if result.Refusal != nil {
			t.Fatalf("resolver refused with kind %q over an owned snapshot that "+
				"DOES carry %q; the refusal above must not pass vacuously",
				result.Refusal.Kind, keyProfile)
		}
		if result.Plan == nil {
			t.Fatal("resolver produced neither a plan nor a refusal")
		}
	})
}

// REQ-33: "**Nothing in the type system prevents that** … The prohibition
// is carried by this contract and by MVV Scenario 6, not by a build
// error" — the obligation is a TEST, not a type.
// ADVERSARIAL
//
// The test states the counterfactual it guards: the defective encoding
// typechecks and resolves to a PLAN, which is exactly why the seam rule
// needs a runtime assertion rather than a compiler.
func TestReq33_ThePlaceholderEncodingTypechecksAndSilentlyRetiresTheRefusal(t *testing.T) {
	// This is the defect, spelled out. It compiles.
	defective := []resolve.Tag{
		{Key: keyStatus, Value: "Draft"},
		{Key: keyProfile, Value: "<absent>"},
	}

	result, err := resolve.Resolve(resolve.Input{
		Flow:       flowID,
		Table:      requiringTable(keyProfile),
		Owned:      defective,
		Recognized: recognizedOutcome,
	})
	if err != nil {
		t.Fatalf("the kernel used the Go error path: %v", err)
	}
	if result.Plan == nil {
		t.Fatalf("the placeholder encoding refused with %+v; this test's premise "+
			"is that it does NOT — `TagSet.has` decides on map presence, so a "+
			"placeholder makes the key read as present and the seam rule must "+
			"be enforced here rather than by the kernel", result.Refusal)
	}

	// The placeholder made the key read as PRESENT and the resolution
	// PLANNED — `owned_state_unavailable` silently retired for that key.
	// No build error stops it, so the executor's own snapshot must omit
	// the key.
	s := newStore(map[string]string{keyStatus: "Draft"})
	e, _ := readerExec(t, s, keyStatus, keyProfile)

	snapshot := e.Read(ctxOf(t), readerName).OwnedSnapshot()

	if v, ok := seamValueOf(snapshot, keyStatus); !ok || v != "Draft" {
		t.Fatalf("the snapshot does not carry %q = %q (snapshot = %+v); an "+
			"empty snapshot omits the absent key for the wrong reason",
			keyStatus, "Draft", snapshot)
	}
	if _, present := seamValueOf(snapshot, keyProfile); present {
		t.Errorf("the executor's own snapshot carries %q, reproducing the "+
			"placeholder encoding that just resolved to a plan instead of "+
			"%q; the prohibition is carried by this test, not by a type",
			keyProfile, resolve.KindOwnedStateUnavailable)
	}
}

// REQ-34: the absent-key row's loudness is conditional and that condition
// is **not** this RDR's to enforce: "this RDR's seam guarantee delivers
// `owned_state_unavailable` **only when the consuming row declares the
// key in `RequiresOwned`**" — a key consumed only through `Row.Match` or a
// guard is OUTSIDE this RDR's surface and MUST NOT be "fixed" here.
// DOMAIN EDGE
//
// The residual is asserted as a residual: with the key consumed only
// through `Row.Match`, the omission yields `no_match`, not
// `owned_state_unavailable` — the kernel's specified behaviour. This test
// pins that the accessor layer does NOT compensate by inventing a
// placeholder or by widening `RequiresOwned` on its own.
func TestReq34_MatchOnlyKeyIsOutsideThisRDRsSurface(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	e, _ := readerExec(t, s, keyStatus, keyProfile)

	read := e.Read(ctxOf(t), readerName)
	mustReadSucceed(t, read)
	snapshot := read.OwnedSnapshot()

	if v, ok := seamValueOf(snapshot, keyStatus); !ok || v != "Draft" {
		t.Fatalf("the snapshot does not carry %q = %q (snapshot = %+v); the "+
			"residual below is only meaningful over a snapshot that carries "+
			"what the artifact HAS", keyStatus, "Draft", snapshot)
	}
	if _, present := seamValueOf(snapshot, keyProfile); present {
		t.Fatalf("the snapshot carries %q; the seam rule omits it regardless of "+
			"how the consuming row reads it", keyProfile)
	}

	result, err := resolve.Resolve(resolve.Input{
		Flow:       flowID,
		Table:      matchOnlyTable(keyProfile),
		Owned:      snapshot,
		Recognized: recognizedOutcome,
	})
	if err != nil {
		t.Fatalf("the kernel used the Go error path: %v", err)
	}
	if result.Refusal == nil {
		t.Fatalf("resolver produced a plan %+v; the row's match pattern names "+
			"%q, which the snapshot omits", result.Plan, keyProfile)
	}
	if result.Refusal.Kind != resolve.KindNoMatch {
		t.Errorf("resolver refusal kind = %q; want %q — a key consumed only "+
			"through Row.Match is dropped by match filtering before "+
			"missingOwned runs. That residual is peer-owned scope and MUST NOT "+
			"be closed by widening the accessor's seam encoding",
			result.Refusal.Kind, resolve.KindNoMatch)
	}
}

// --- kernel tables for the cross-layer assertions ------------------------

const recognizedOutcome = "Recognized"

// requiringTable is a one-row table whose row DECLARES the key in
// RequiresOwned — the set JDR 0001 §JD-3 assigns to RDR 0002's
// normalizer, derived as the rule's write block plus clear list.
func requiringTable(required string) resolve.Table {
	return resolve.Table{
		Revision: "rev-0004-seam",
		Outcomes: []string{recognizedOutcome},
		Rows: []resolve.Row{{
			RuleID:        "rdr.requires",
			SourceLocator: statePath + ":20",
			Outcome:       recognizedOutcome,
			RequiresOwned: []string{required},
			Writes:        []resolve.Tag{{Key: required, Value: "large"}},
			NextTags:      []resolve.Tag{{Key: required, Value: "large"}},
		}},
	}
}

// matchOnlyTable consumes the key through Row.Match alone, so it never
// enters RequiresOwned and the omission yields no_match.
func matchOnlyTable(consumed string) resolve.Table {
	return resolve.Table{
		Revision: "rev-0004-seam",
		Outcomes: []string{recognizedOutcome},
		Rows: []resolve.Row{{
			RuleID:        "rdr.match.only",
			SourceLocator: statePath + ":21",
			Outcome:       recognizedOutcome,
			Match:         []resolve.Tag{{Key: consumed, Value: "large"}},
			NextTags:      []resolve.Tag{{Key: keyStatus, Value: "Final"}},
		}},
	}
}

var _ = accessor.CapRead
