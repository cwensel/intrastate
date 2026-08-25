package accessor_test

// RDR 0004 — read branch discipline and the completeness guarantee
// (REQ-13..REQ-29).
//
// The coverage floor this file exists to hold: an implementation that
// THINS the value set instead of refusing produces the same shape as a
// genuine absence, and every downstream presence predicate then decides
// on state that was never read. The discriminating pair is
// TestReq21_… — one artifact truncates, one legitimately lacks the key,
// and they must take DIFFERENT BRANCHES.
//
// The second mutant this file kills is the derived key set: an
// implementation computing the requested keys from what came back makes
// completeness self-fulfilling. Every fixture here pins `keys` in the
// definition and lets the artifact contents disagree with it.

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/newcoinc/intrastate/internal/accessor"
	"github.com/newcoinc/intrastate/internal/table"
)

// REQ-13: "A read accessor MUST return typed tag values or a typed
// refusal."
// HAPPY PATH
func TestReq13_ReadReturnsTypedValuesOrATypedRefusal(t *testing.T) {
	t.Run("values_branch", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
		e, _ := readerExec(t, s, keyStatus, keyProfile)

		got := e.Read(ctxOf(t), readerName)

		values := mustReadSucceed(t, got)
		if v, ok := valueOf(values, keyStatus); !ok || v.Value != "Draft" {
			t.Errorf("value for %q = %+v (found=%v); want the typed value %q",
				keyStatus, v, ok, "Draft")
		}
		if got.Refusal != nil {
			t.Error("the values branch carries a refusal; the disjunction is exclusive")
		}
	})

	t.Run("refusal_branch", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		s.unreadable[keyStatus] = true
		e, _ := readerExec(t, s, keyStatus)

		got := e.Read(ctxOf(t), readerName)

		mustReadRefuse(t, got, accessor.ClassIncompleteRead)
		if len(got.Values) != 0 {
			t.Errorf("the refusal branch carries values %+v; the disjunction is "+
				"exclusive and a refusal carries no values", got.Values)
		}
	})
}

// REQ-14: "It MUST NOT mutate authoritative artifacts."
// ADVERSARIAL
//
// The oracle is the store's content before and after. A read that writes
// back a normalized value would pass every branch assertion above.
func TestReq14_ReadDoesNotMutateTheAuthoritativeArtifact(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
	before := map[string]string{}
	for k, v := range s.tags {
		before[k] = v
	}
	e, b := readerExec(t, s, keyStatus, keyProfile)

	values := mustReadSucceed(t, e.Read(ctxOf(t), readerName))

	// Positive precondition: the read must have actually reached the
	// artifact and returned its values. Otherwise "nothing changed" is
	// satisfied by a read that never ran.
	if b.reads != 1 {
		t.Fatalf("the read binding ran %d times; want exactly 1 — a read that "+
			"never reaches the artifact cannot witness non-mutation", b.reads)
	}
	if v, ok := valueOf(values, keyStatus); !ok || v.Value != "Draft" {
		t.Fatalf("the read returned %+v; want %q = %q read from the artifact",
			values, keyStatus, "Draft")
	}

	if len(s.tags) != len(before) {
		t.Fatalf("artifact key count moved from %d to %d across a READ",
			len(before), len(s.tags))
	}
	for k, want := range before {
		if got := s.tags[k]; got != want {
			t.Errorf("artifact key %q moved from %q to %q across a READ; a read "+
				"accessor must not mutate authoritative artifacts", k, want, got)
		}
	}
}

// REQ-19: "The typed-tag-values branch carries a completeness guarantee:
// a read accessor MUST return the tag set for exactly the keys it was
// asked for — no requested key missing, no unrequested key added — or
// take the refusal branch."
// REQ-24: "Consumers may therefore treat a returned tag set as exactly
// the requested keys."
// HAPPY PATH
//
// The "no unrequested key added" half is the one a store-dumping
// implementation fails: the store carries three keys and the definition
// requests two.
func TestReq19_SuccessBranchCarriesExactlyTheRequestedKeys(t *testing.T) {
	s := newStore(map[string]string{
		keyStatus:  "Draft",
		keyProfile: "large",
		keyLabels:  `["alpha"]`,
	})
	requested := []string{keyStatus, keyProfile}
	e, _ := readerExec(t, s, requested...)

	values := mustReadSucceed(t, e.Read(ctxOf(t), readerName))

	want := slices.Clone(requested)
	slices.Sort(want)
	if got := keysOf(values); !slices.Equal(got, want) {
		t.Errorf("returned key set = %v; want exactly the requested %v — no "+
			"requested key missing, no unrequested key added (the artifact also "+
			"carries %q, which was NOT requested)", got, want, keyLabels)
	}
}

// REQ-15: "Every read accessor definition MUST declare the requested key
// set as validated metadata." [0002-delivered] — `keys` is required on
// every capability-table entry.
// HAPPY PATH
//
// This RDR consumes RDR 0002's carried `keys` rather than spelling a
// second field. The test cites the shipped surface: `table.Accessor.Keys`
// is what `RequestedKeys` reads.
func TestReq15_RequestedKeySetIsTheDefinitionsCarriedKeys(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	def := readerDef(&readBinding{store: s}, keyStatus, keyProfile)

	if got, want := def.Accessor.Keys, []string{keyStatus, keyProfile}; !slices.Equal(got, want) {
		t.Fatalf("fixture carries keys %v; want %v — RDR 0002's `table.Accessor` "+
			"is the carrier this RDR consumes", got, want)
	}
	got := def.RequestedKeys()
	want := []string{keyStatus, keyProfile}
	slices.Sort(got)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Errorf("RequestedKeys() = %v; want the definition's declared %v",
			got, want)
	}
}

// REQ-16: "The set MUST NOT be derived from the keys a read actually
// resolved."
// REQ-18: "a read that lost keys still reports success over a smaller
// set" is the failure the pinned set prevents; the implementation MUST
// NOT port the spike's `expectedTagKeys`.
// ADVERSARIAL
//
// ORA 2's second control, stated exactly: the definition requests two
// keys, the binding can resolve only one, and a derived-key-set
// implementation returns a ONE-KEY SUCCESS. The assertion is that the
// truncation refuses instead.
func TestReq16_RequestedSetIsPinnedNotDerivedFromWhatResolved(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
	s.unreadable[keyProfile] = true
	e, _ := readerExec(t, s, keyStatus, keyProfile)

	got := e.Read(ctxOf(t), readerName)

	if !got.Refused() {
		t.Fatalf("read returned success over values %+v; the requested key set is "+
			"the definition's %v, so losing %q is a REFUSAL — deriving the set "+
			"from what resolved makes completeness self-fulfilling",
			got.Values, []string{keyStatus, keyProfile}, keyProfile)
	}
	mustReadRefuse(t, got, accessor.ClassIncompleteRead)
}

// REQ-17: "A missing or empty requested key set MUST fail validation
// before execution." — the EIGHTH validation arm, explicitly unwitnessed
// by the spike (A9).
// INPUT EDGE
func TestReq17_MissingOrEmptyRequestedKeySetFailsValidation(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})

	for name, keys := range map[string][]string{
		"empty": {},
		"nil":   nil,
	} {
		t.Run(name, func(t *testing.T) {
			def := readerDef(&readBinding{store: s})
			def.Accessor.Keys = keys

			fs := accessor.Validate(registryOf(def), []accessor.Identity{def.Identity})

			if !hasCode(fs, accessor.CodeMissingRequestedKeySet) {
				t.Errorf("validation codes = %v; want %q — a reader with a %s "+
					"requested key set must fail BEFORE execution",
					findingCodes(fs), accessor.CodeMissingRequestedKeySet, name)
			}
		})
	}
}

// REQ-20: "A partial or truncated read is a refusal, not a value."
// REQ-22: "One unreadable requested key MUST refuse the whole read; the
// executor MUST NOT return the keys that did resolve."
// REQ-23: "the success predicate for a read accessor is \"every requested
// key resolved\", not \"at least one key resolved\"."
// ADVERSARIAL
func TestReq20_OneUnreadableKeyRefusesTheWholeRead(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
	s.unreadable[keyProfile] = true
	e, _ := readerExec(t, s, keyStatus, keyProfile)

	got := e.Read(ctxOf(t), readerName)

	mustReadRefuse(t, got, accessor.ClassIncompleteRead)
	if len(got.Values) != 0 {
		t.Errorf("the refusal carries values %+v; %q resolved but the executor "+
			"MUST NOT return the keys that did resolve — the success predicate is "+
			"\"every requested key resolved\", not \"at least one\"",
			got.Values, keyStatus)
	}
}

// REQ-21: "A missing key MUST be distinguishable from an unread key only
// by which branch is taken — absence is a value, unreadability is a
// refusal."
// DOMAIN EDGE
//
// The single most load-bearing discriminator in this RDR. The two
// artifacts differ ONLY in why `profile` did not come back: one is
// unreadable, one is legitimately not carried. A thinning implementation
// produces one shape for both.
func TestReq21_AbsenceIsAValueAndUnreadabilityIsARefusal(t *testing.T) {
	requested := []string{keyStatus, keyProfile}

	t.Run("unreadable_key_takes_the_refusal_branch", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
		s.unreadable[keyProfile] = true
		e, _ := readerExec(t, s, requested...)

		got := e.Read(ctxOf(t), readerName)

		mustReadRefuse(t, got, accessor.ClassIncompleteRead)
	})

	t.Run("genuinely_absent_key_takes_the_values_branch", func(t *testing.T) {
		// Same requested set; the artifact simply does not carry `profile`.
		s := newStore(map[string]string{keyStatus: "Draft"})
		e, _ := readerExec(t, s, requested...)

		got := e.Read(ctxOf(t), readerName)

		values := mustReadSucceed(t, got)
		want := slices.Clone(requested)
		slices.Sort(want)
		if k := keysOf(values); !slices.Equal(k, want) {
			t.Fatalf("returned key set = %v; want exactly %v — absence is a VALUE "+
				"on the success branch", k, want)
		}
		v, ok := valueOf(values, keyProfile)
		if !ok {
			t.Fatalf("the success branch dropped %q; absence is a value, not an "+
				"omission from the accessor's own result", keyProfile)
		}
		if !v.Absent {
			t.Errorf("value for %q = %+v; want Absent=true — the binding read the "+
				"artifact successfully and the key was not there", keyProfile, v)
		}
	})
}

// REQ-25: "A read that cannot resolve every requested key MUST be
// reported as `incomplete_read`, its own refusal class, distinct from
// execution failure and from timeout."
// BOUNDARY
//
// The three classes are asserted as three DIFFERENT names over three
// inputs. An implementation collapsing them into one vague error passes
// nothing here.
func TestReq25_IncompleteReadIsDistinctFromExecutionFailureAndTimeout(t *testing.T) {
	t.Run("incomplete_read", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
		s.unreadable[keyProfile] = true
		e, _ := readerExec(t, s, keyStatus, keyProfile)

		mustReadRefuse(t, e.Read(ctxOf(t), readerName), accessor.ClassIncompleteRead)
	})

	t.Run("execution_failure", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		b := &readBinding{store: s, failWith: errBindingFailed}
		e := accessor.NewExecutor(registryOf(readerDef(b, keyStatus)), artifactsOf())

		mustReadRefuse(t, e.Read(ctxOf(t), readerName), accessor.ClassExecutionFailure)
	})

	t.Run("timeout", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		b := &readBinding{store: s, delay: 10 * fixtureTimeout}
		e := accessor.NewExecutor(registryOf(readerDef(b, keyStatus)), artifactsOf())

		mustReadRefuse(t, e.Read(ctxOf(t), readerName), accessor.ClassTimeout)
	})

	t.Run("the_three_classes_are_distinct_names", func(t *testing.T) {
		classes := []accessor.RefusalClass{
			accessor.ClassIncompleteRead,
			accessor.ClassExecutionFailure,
			accessor.ClassTimeout,
		}
		seen := map[accessor.RefusalClass]bool{}
		for _, c := range classes {
			if seen[c] {
				t.Errorf("refusal class %q is spelled the same as another; "+
					"incomplete_read, execution_failure, and timeout are three "+
					"distinct classes", c)
			}
			seen[c] = true
		}
	})
}

// REQ-26: "The refusal MUST name the requested keys it could not read." —
// TS 2 adds "the refusal carries no values and names the keys it could
// not read".
// BOUNDARY
//
// The payload assertion is what makes the refusal actionable. It also
// discriminates the "names EVERY requested key" mutant: only `profile`
// was unreadable, so only `profile` belongs in the payload.
func TestReq26_IncompleteReadNamesTheKeysItCouldNotRead(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
	s.unreadable[keyProfile] = true
	e, _ := readerExec(t, s, keyStatus, keyProfile)

	r := mustReadRefuse(t, e.Read(ctxOf(t), readerName), accessor.ClassIncompleteRead)

	got := slices.Clone(r.Keys)
	slices.Sort(got)
	if want := []string{keyProfile}; !slices.Equal(got, want) {
		t.Errorf("refusal names keys %v; want exactly %v — the keys it could not "+
			"read, not every requested key (%q resolved fine)",
			got, want, keyStatus)
	}
}

// REQ-27: "When a read exceeds its timeout before resolving every
// requested key, `timeout` takes precedence over `incomplete_read`."
// BOUNDARY
//
// ORA 7: the fixture "makes both classes true at once and asserts the
// name `timeout`". The spike could not witness this — it sleeps BEFORE
// its key loop so the two never overlap (REQ-124). Here the binding
// resolves one key, cannot read another, AND runs past its deadline.
//
// The negative control is the scenario-2 truncation fixture, which must
// still return `incomplete_read`: an implementation that always reports
// `timeout` fails that leg.
func TestReq27_TimeoutOutranksIncompleteReadWhenBothAreTrue(t *testing.T) {
	t.Run("both_true_at_once_reports_timeout", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
		s.unreadable[keyProfile] = true
		b := &overlapReadBinding{store: s, delay: 10 * fixtureTimeout}
		e := accessor.NewExecutor(
			registryOf(readerDef(b, keyStatus, keyProfile)), artifactsOf())

		got := e.Read(ctxOf(t), readerName)

		if !b.sawUnreadable {
			t.Fatal("the fixture did not make both classes true at once; the " +
				"unreadable key must be met BEFORE the deadline expires")
		}
		mustReadRefuse(t, got, accessor.ClassTimeout)
	})

	t.Run("negative_control_truncation_alone_reports_incomplete_read", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
		s.unreadable[keyProfile] = true
		e, _ := readerExec(t, s, keyStatus, keyProfile)

		mustReadRefuse(t, e.Read(ctxOf(t), readerName), accessor.ClassIncompleteRead)
	})
}

// REQ-28: "a key is *absent* only when the binding read the artifact
// successfully and the key was not there. Every other outcome — the
// artifact did not parse, the transport truncated, the key's value failed
// type coercion, the value read is the reserved `<clear>` literal,
// permission was denied — is *unreadable*, because in each the binding
// did not establish that the key is missing."
// REQ-48: "A read that yields `<clear>` as a value MUST treat that key as
// unreadable and refuse `incomplete_read`."
// DOMAIN EDGE
//
// The `<clear>` arm is A11's read half and is unwitnessed by the spike:
// no spike fixture carries the sentinel. A binding that hands back
// `<clear>` as a VALUE lets a downstream comparison match the literal.
func TestReq28_ReservedClearLiteralReadAsAValueIsUnreadable(t *testing.T) {
	s := newStore(map[string]string{keyStatus: table.ClearSentinel, keyProfile: "large"})
	e, _ := readerExec(t, s, keyStatus, keyProfile)

	got := e.Read(ctxOf(t), readerName)

	r := mustReadRefuse(t, got, accessor.ClassIncompleteRead)
	if !slices.Contains(r.Keys, keyStatus) {
		t.Errorf("refusal names keys %v; want %q named — a read yielding the "+
			"reserved literal %q must treat that key as UNREADABLE",
			r.Keys, keyStatus, table.ClearSentinel)
	}
	if v, ok := valueOf(got.Values, keyStatus); ok {
		t.Errorf("the result carries %q = %+v; the reserved literal must never "+
			"cross as a value", keyStatus, v)
	}
}

// REQ-29: "A binding that cannot tell the two apart at its own boundary
// MUST report `incomplete_read`; guessing absence is the failure this
// contract exists to prevent." — unreadable is the DEFAULT classification.
// ADVERSARIAL
//
// The binding returns a key in NEITHER the values nor the unreadable
// list — it cannot tell. The executor must default to unreadable rather
// than silently treating the gap as absence.
func TestReq29_UnclassifiedKeyDefaultsToUnreadableNotAbsent(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	b := &silentDropReadBinding{store: s, drop: keyProfile}
	e := accessor.NewExecutor(
		registryOf(readerDef(b, keyStatus, keyProfile)), artifactsOf())

	got := e.Read(ctxOf(t), readerName)

	r := mustReadRefuse(t, got, accessor.ClassIncompleteRead)
	if !slices.Contains(r.Keys, keyProfile) {
		t.Errorf("refusal names keys %v; want %q — the binding classified it as "+
			"neither read nor unreadable, and guessing ABSENCE is exactly the "+
			"failure this contract exists to prevent", r.Keys, keyProfile)
	}
}

// REQ-103: "No CONTRADICTION row, on two pairs." The C4 disjunction and
// the C6 completeness guarantee are jointly satisfiable "because
// completeness constrains *which* branch the disjunction takes rather
// than adding a third branch".
// BOUNDARY
//
// Totality plus exclusivity over the three read input classes: every read
// lands on exactly one of the two branches, never both and never neither.
func TestReq103_ReadHasExactlyTwoBranchesNeverAThird(t *testing.T) {
	cases := map[string]func() *store{
		"complete": func() *store {
			return newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
		},
		"genuine_absence": func() *store {
			return newStore(map[string]string{keyStatus: "Draft"})
		},
		"unreadable": func() *store {
			s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
			s.unreadable[keyProfile] = true
			return s
		},
	}

	for name, mk := range cases {
		t.Run(name, func(t *testing.T) {
			e, _ := readerExec(t, mk(), keyStatus, keyProfile)

			got := e.Read(ctxOf(t), readerName)

			refused := got.Refusal != nil
			valued := len(got.Values) != 0
			if refused && valued {
				t.Errorf("result carries BOTH a refusal (%q) and values %+v; the "+
					"disjunction is exclusive", got.Refusal.Class, got.Values)
			}
			if !refused && !valued {
				t.Error("result carries neither a refusal nor values; the " +
					"disjunction is total — completeness constrains WHICH branch " +
					"is taken, it does not add a third")
			}
			if got.Refused() != refused {
				t.Errorf("Refused() = %v; want %v", got.Refused(), refused)
			}
		})
	}
}

// --- bindings for the overlap and silent-drop oracles --------------------

// overlapReadBinding meets an unreadable key AND runs past its deadline
// in the same invocation, so `timeout` and `incomplete_read` are both
// true at once. This is the shape CA A9 says the spike cannot witness.
type overlapReadBinding struct {
	store *store
	delay time.Duration

	sawUnreadable bool
}

func (b *overlapReadBinding) Capability() accessor.Capability { return accessor.CapRead }

func (b *overlapReadBinding) Read(
	ctx context.Context, art accessor.Artifact, requested []string,
) ([]accessor.KeyValue, []string, error) {
	var values []accessor.KeyValue
	var unreadable []string
	for _, key := range requested {
		v, absent, readable := b.store.get(key)
		if !readable {
			b.sawUnreadable = true
			unreadable = append(unreadable, key)
			continue
		}
		values = append(values, accessor.KeyValue{Key: key, Value: v, Absent: absent})
	}
	// The truncation is already decided; NOW the deadline expires.
	select {
	case <-time.After(b.delay):
	case <-ctx.Done():
		return values, unreadable, ctx.Err()
	}
	return values, unreadable, nil
}

// silentDropReadBinding returns a key in neither list: it could not tell
// absence from unreadability at its own boundary.
type silentDropReadBinding struct {
	store *store
	drop  string
}

func (b *silentDropReadBinding) Capability() accessor.Capability { return accessor.CapRead }

func (b *silentDropReadBinding) Read(
	ctx context.Context, art accessor.Artifact, requested []string,
) ([]accessor.KeyValue, []string, error) {
	var values []accessor.KeyValue
	for _, key := range requested {
		if key == b.drop {
			continue
		}
		v, absent, readable := b.store.get(key)
		if !readable {
			continue
		}
		values = append(values, accessor.KeyValue{Key: key, Value: v, Absent: absent})
	}
	return values, nil, nil
}
