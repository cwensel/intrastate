package accessor_test

// RDR 0004 — the boundary properties: bounded timeout, no package
// prints, the closed refusal-class set, and the disposition table's
// totality (REQ-68..REQ-82, REQ-96..REQ-100).

import (
	"bytes"
	"io"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/newcoinc/intrastate/internal/accessor"
	"github.com/newcoinc/intrastate/internal/resolve"
)

// REQ-68: "Every accessor invocation MUST have a bounded timeout." —
// every read, gate, and write invocation, INCLUDING the read-back
// re-read.
// REQ-71: the bounded-execution property is the non-functional check; no
// throughput target is set.
// BOUNDARY
//
// The oracle is that a binding which blocks forever returns within a
// bound derived from the DEFINITION'S DECLARED TIMEOUT, not from the
// test's patience and not merely from the blocking delay. Comparing
// `elapsed` against `blockFor` would pass a deadline wrong by five
// orders of magnitude — the bound must be a margin on `fixtureTimeout`,
// which is what the definitions declare.
func TestReq68_EveryInvocationIsBounded(t *testing.T) {
	// The binding blocks effectively forever, so nothing but the
	// declared timeout can end the invocation.
	const blockFor = time.Hour
	// The bound the assertions judge against: generous enough for a
	// loaded machine to schedule the timer and unwind, and still four
	// orders of magnitude below `blockFor`. Not tight enough to pin the
	// implementation's exact deadline arithmetic — REQ-68 requires a
	// bound derived from the declared timeout, not a specific one.
	const boundedBy = 20 * fixtureTimeout

	t.Run("read", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		b := &readBinding{store: s, delay: blockFor}
		e := accessor.NewExecutor(registryOf(readerDef(b, keyStatus)), artifactsOf())

		start := time.Now()
		got := e.Read(ctxOf(t), readerName)
		elapsed := time.Since(start)

		mustReadRefuse(t, got, accessor.ClassTimeout)
		if elapsed >= boundedBy {
			t.Errorf("the read took %v; the invocation is bounded by the "+
				"definition's declared timeout %v, not by the binding's %v delay",
				elapsed, fixtureTimeout, blockFor)
		}
	})

	t.Run("gate", func(t *testing.T) {
		e := gateExec(t, &gateBinding{verdict: accessor.VerdictAllow, delay: blockFor})

		start := time.Now()
		got := e.Gate(ctxOf(t), gateName)
		elapsed := time.Since(start)

		if !got.Refused() || got.Refusal.Class != accessor.ClassTimeout {
			t.Fatalf("gate result = %+v; want a %q refusal", got, accessor.ClassTimeout)
		}
		if elapsed >= boundedBy {
			t.Errorf("the gate took %v; the invocation is bounded by the "+
				"definition's declared timeout %v, not by the binding's %v delay",
				elapsed, fixtureTimeout, blockFor)
		}
	})

	t.Run("write", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		w := &writeBinding{store: s, delay: blockFor}
		e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

		start := time.Now()
		got := e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))
		elapsed := time.Since(start)

		mustWriteRefuse(t, got, accessor.ClassTimeout)
		if elapsed >= boundedBy {
			t.Errorf("the write took %v; the invocation is bounded by the "+
				"definition's declared timeout %v, not by the binding's %v delay",
				elapsed, fixtureTimeout, blockFor)
		}
	})

	t.Run("read_back_re_read_is_its_own_bounded_invocation", func(t *testing.T) {
		// The write lands immediately; the RE-READ blocks. If the re-read
		// borrowed the write's already-spent budget rather than running as
		// its own bounded invocation, this could not be classified.
		s := newStore(map[string]string{keyStatus: "Draft"})
		w := &writeBinding{store: s}
		r := &readBinding{store: s, delay: blockFor}
		e := writeExec(t, s, w, r, keyStatus)

		start := time.Now()
		got := e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))
		elapsed := time.Since(start)

		// The CLASS, not merely a refusal: the re-read expired against
		// its own deadline, so `timeout` is what the boundary reports.
		// Accepting any refusal would pass an executor that classified
		// the blocked re-read as an execution failure or as read-back
		// incomplete, which says the verification could not run rather
		// than that its invocation was bounded (REQ-68, REQ-69).
		mustWriteRefuse(t, got, accessor.ClassTimeout)
		if elapsed >= boundedBy {
			t.Errorf("the read-back took %v; every accessor invocation is bounded "+
				"by its declared timeout %v — the re-read included — not by the "+
				"binding's %v delay", elapsed, fixtureTimeout, blockFor)
		}
	})
}

// REQ-69: "Timeout MUST be reported as its own refusal class, distinct
// from execution failure and read-back mismatch."
// REQ-109: "Make them separate refusal classes."
// BOUNDARY
func TestReq69_TimeoutIsDistinctFromExecutionFailureAndReadBackMismatch(t *testing.T) {
	classes := map[string]accessor.RefusalClass{
		"timeout":            accessor.ClassTimeout,
		"execution_failure":  accessor.ClassExecutionFailure,
		"read_back_mismatch": accessor.ClassReadBackMismatch,
	}
	seen := map[accessor.RefusalClass]string{}
	for name, c := range classes {
		if string(c) != name {
			t.Errorf("class %s spells %q; want %q", name, c, name)
		}
		if prior, dup := seen[c]; dup {
			t.Errorf("classes %s and %s share the spelling %q", name, prior, c)
		}
		seen[c] = name
		// All three must be MEMBERS of the enumerated set, or a mapper
		// cannot tell them apart at the boundary.
		if !slices.Contains(accessor.RefusalClasses(), c) {
			t.Errorf("class %q is not in RefusalClasses() = %v; three distinct "+
				"spellings are only distinguishable if all three are enumerated",
				c, accessor.RefusalClasses())
		}
	}
}

// REQ-70: "Missing or non-positive timeout metadata MUST fail validation
// before execution." [0002-delivered] —
// `internal/table/load.go::accessorTable` refuses an absent timeout, a
// non-duration timeout, and a non-positive timeout.
// INPUT EDGE
func TestReq70_MissingOrNonPositiveTimeoutFailsValidation(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})

	for name, timeout := range map[string]string{
		"absent":      "",
		"zero":        "0s",
		"negative":    "-1s",
		"non_numeric": "soon",
	} {
		t.Run(name, func(t *testing.T) {
			def := readerDef(&readBinding{store: s}, keyStatus)
			def.Accessor.Timeout = timeout

			fs := accessor.Validate(registryOf(def), []accessor.Identity{def.Identity})

			if !hasCode(fs, accessor.CodeMissingOrNonPositiveTimeout) {
				t.Errorf("validation codes = %v for timeout %q; want %q — the "+
					"defect must be caught BEFORE execution", findingCodes(fs),
					timeout, accessor.CodeMissingOrNonPositiveTimeout)
			}
		})
	}

	t.Run("control_a_positive_timeout_passes", func(t *testing.T) {
		def := readerDef(&readBinding{store: s}, keyStatus)

		fs := accessor.Validate(registryOf(def), []accessor.Identity{def.Identity})

		if hasCode(fs, accessor.CodeMissingOrNonPositiveTimeout) {
			t.Errorf("validation codes = %v for the positive timeout %q; want no "+
				"timeout finding", findingCodes(fs), def.Accessor.Timeout)
		}
	})
}

// REQ-72: "The accessor package MUST return structured success/refusal
// values and MUST NOT write to stdout or stderr directly."
// REQ-73: "the no-print property is asserted by capturing stdout/stderr
// around the call and requiring both empty" and "the accessor package
// test asserts the returned structured value carries the refusal class".
// "Scenario 5 is the one absence-of-error oracle; the named capture
// control is what makes it discriminating and is normative for the
// implementation test."
// ADVERSARIAL
//
// Both halves are mandatory. The POSITIVE half (the returned value
// carries the refusal class) is what stops this passing vacuously, and
// the capture control is what makes the negative half discriminating —
// `TestReq73_TheCaptureControlItselfDiscriminates` proves the capture can
// actually fail.
func TestReq72_TheAccessorPackageNeverPrints(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
	s.unreadable[keyProfile] = true

	var got accessor.ReadResult
	stdout, stderr := captureOutput(t, func() {
		e, _ := readerExec(t, s, keyStatus, keyProfile)
		got = e.Read(ctxOf(t), readerName)
	})

	// Positive half: the structured value carries the refusal class.
	if got.Refusal == nil {
		t.Fatal("the read returned no refusal; the no-print property must be " +
			"asserted alongside a POSITIVE assertion on the returned value, or " +
			"it passes by absence of output alone")
	}
	if got.Refusal.Class != accessor.ClassIncompleteRead {
		t.Errorf("refusal class = %q; want %q carried on the returned value",
			got.Refusal.Class, accessor.ClassIncompleteRead)
	}

	// Negative half: the capture control.
	if stdout != "" {
		t.Errorf("the accessor package wrote %q to stdout; it returns "+
			"structured values and never prints", stdout)
	}
	if stderr != "" {
		t.Errorf("the accessor package wrote %q to stderr; `respond.Fail` is "+
			"CLI-only and is not the executor's API", stderr)
	}
}

// REQ-73, control half: "A test that deliberately prints must fail the
// capture assertion; without that control this scenario passes
// vacuously."
// ADVERSARIAL
func TestReq73_TheCaptureControlItselfDiscriminates(t *testing.T) {
	stdout, stderr := captureOutput(t, func() {
		os.Stdout.WriteString("deliberate stdout")
		os.Stderr.WriteString("deliberate stderr")
	})

	if stdout == "" {
		t.Error("the capture control saw no stdout for a deliberate write; the " +
			"no-print assertion would then pass vacuously for any implementation")
	}
	if stderr == "" {
		t.Error("the capture control saw no stderr for a deliberate write; the " +
			"no-print assertion would then pass vacuously")
	}
}

// REQ-74: "Accessor package must not print." — `respond.Fail` is CLI-only
// and is not the executor's API.
// REQ-75: "Expose accessor refusal classes to the CLI layer so RDR 0005
// can map them through `respond.Fail` and `clierr.ExitCodeFor` without
// package-level prints." — the mapping itself is RDR 0005's; this RDR
// ships the EXPOSURE.
// REQ-117: "The accessor package returns structured values only; the CLI
// layer can map them through `respond.Fail` and `clierr.ExitCodeFor`."
// BOUNDARY
//
// The exposure is asserted as a property of the returned value: the class
// is a stable, enumerable discriminator a mapper can switch on without
// inspecting error strings.
func TestReq75_RefusalClassesAreExposedForTheCLIToMap(t *testing.T) {
	classes := accessor.RefusalClasses()
	if len(classes) == 0 {
		t.Fatal("RefusalClasses() is empty; RDR 0005 cannot build a mapping " +
			"over a set it cannot enumerate")
	}

	s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
	s.unreadable[keyProfile] = true
	e, _ := readerExec(t, s, keyStatus, keyProfile)

	got := e.Read(ctxOf(t), readerName)

	if got.Refusal == nil {
		t.Fatal("no refusal to expose")
	}
	if !slices.Contains(classes, got.Refusal.Class) {
		t.Errorf("the returned class %q is outside the enumerated set %v; a "+
			"mapper switching on the set would fall through", got.Refusal.Class, classes)
	}
}

// REQ-76: "Visible failures are typed refusals: unknown accessor,
// capability mismatch, timeout, execution failure, incomplete read, gate
// denied, gate indeterminate, write attempted for a non-owned tag,
// read-back mismatch, and read-back incomplete." — the wire spellings
// appear in DISP.
// BOUNDARY
//
// The eight WIRE classes are asserted by exact spelling. "gate denied" is
// a typed result at this boundary (REQ-37) and "write attempted for a
// non-owned tag" is a VALIDATION code (REQ-41), not runtime classes —
// which is why the runtime set is eight, not ten.
func TestReq76_TheRefusalClassSetIsExactlyTheEightWireSpellings(t *testing.T) {
	want := []accessor.RefusalClass{
		"unknown_accessor",
		"capability_mismatch",
		"timeout",
		"execution_failure",
		"incomplete_read",
		"gate_indeterminate",
		"read_back_mismatch",
		"read_back_incomplete",
	}

	got := slices.Clone(accessor.RefusalClasses())
	slices.Sort(got)
	sortedWant := slices.Clone(want)
	slices.Sort(sortedWant)

	if !slices.Equal(got, sortedWant) {
		t.Errorf("RefusalClasses() = %v; want exactly %v", got, sortedWant)
	}
}

// REQ-77: "These are accessor refusal classes and are disjoint from the
// kernel's closed five-kind `resolve.RefusalKinds` set, which this RDR
// does not extend" — the implementation MUST NOT add a member to
// `RefusalKinds()`.
// ADVERSARIAL
func TestReq77_AccessorClassesAreDisjointFromTheKernelRefusalKinds(t *testing.T) {
	kernel := resolve.RefusalKinds()
	if len(kernel) != 5 {
		t.Errorf("resolve.RefusalKinds() has %d members: %v; the kernel set is "+
			"CLOSED at five and this RDR does not extend it", len(kernel), kernel)
	}

	for _, k := range kernel {
		for _, c := range accessor.RefusalClasses() {
			if string(k) == string(c) {
				t.Errorf("accessor class %q collides with kernel refusal kind %q; "+
					"the two sets are disjoint", c, k)
			}
		}
	}
}

// REQ-78: "`exit code` is out of scope: RDR 0005 owns the CLI mapping;
// this table stops at the structured value the accessor package returns."
// BOUNDARY
//
// The refusal carries no exit code and no CLI code. An implementation
// that pre-decides the exit here has taken RDR 0005's decision.
func TestReq78_TheRefusalCarriesNoExitCodeOrCLICode(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
	s.unreadable[keyProfile] = true
	e, _ := readerExec(t, s, keyStatus, keyProfile)

	r := mustReadRefuse(t, e.Read(ctxOf(t), readerName), accessor.ClassIncompleteRead)

	// The refusal's own vocabulary is the class plus the diagnosis tuple.
	// Nothing here names an exit status.
	if r.Class == "" {
		t.Fatal("the refusal carries no class")
	}
	for _, notAClass := range []string{"1", "2", "3", "exit"} {
		if string(r.Class) == notAClass {
			t.Errorf("refusal class = %q; the boundary stops at the structured "+
				"value — RDR 0005 owns the exit mapping", r.Class)
		}
	}
}

// REQ-82: "Diagnosis starts with the accessor identity, capability,
// artifact role, timeout, and expected versus observed tag values." — a
// refusal payload carries enough to name those.
// REQ-110: "Validation records capability, artifact role, and tag keys."
// BOUNDARY
func TestReq82_RefusalPayloadCarriesTheDiagnosisTuple(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
	s.unreadable[keyProfile] = true
	e, _ := readerExec(t, s, keyStatus, keyProfile)

	r := mustReadRefuse(t, e.Read(ctxOf(t), readerName), accessor.ClassIncompleteRead)

	if r.Accessor != readerName {
		t.Errorf("refusal accessor = %q; want %q", r.Accessor, readerName)
	}
	if r.Capability != accessor.CapRead {
		t.Errorf("refusal capability = %q; want %q", r.Capability, accessor.CapRead)
	}
	if r.Role != stateRole {
		t.Errorf("refusal role = %q; want %q", r.Role, stateRole)
	}
	if r.Timeout != fixtureTimeout {
		t.Errorf("refusal timeout = %v; want the declared %v", r.Timeout, fixtureTimeout)
	}
}

// REQ-79: "No input class exits silently **at this boundary**: every row
// either returns a value, mints a named refusal, or — for a
// genuinely-absent key — omits that key from the owned snapshot."
// REQ-80: "A silent failure is any outcome where the accessor returns a
// **success-shaped** result over state it did not establish" — four named
// shapes, each with a mandatory guard.
// BOUNDARY
//
// The Disposition Table's totality claim, walked as a table: every input
// class the boundary can meet lands on a named outcome, and no
// success-shaped result covers state the accessor did not establish.
func TestReq79_EveryInputClassLandsOnANamedOutcome(t *testing.T) {
	cases := map[string]struct {
		run  func(t *testing.T) disposition
		want disposition
	}{
		"all_requested_keys_read": {
			run: func(t *testing.T) disposition {
				s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
				e, _ := readerExec(t, s, keyStatus, keyProfile)
				return dispositionOfRead(e.Read(ctxOf(t), readerName))
			},
			want: disposition{refused: false},
		},
		"requested_key_unreadable": {
			run: func(t *testing.T) disposition {
				s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
				s.unreadable[keyProfile] = true
				e, _ := readerExec(t, s, keyStatus, keyProfile)
				return dispositionOfRead(e.Read(ctxOf(t), readerName))
			},
			want: disposition{refused: true, class: accessor.ClassIncompleteRead},
		},
		"requested_key_absent_from_artifact": {
			run: func(t *testing.T) disposition {
				s := newStore(map[string]string{keyStatus: "Draft"})
				e, _ := readerExec(t, s, keyStatus, keyProfile)
				return dispositionOfRead(e.Read(ctxOf(t), readerName))
			},
			want: disposition{refused: false},
		},
		"accessor_exceeds_its_timeout": {
			run: func(t *testing.T) disposition {
				s := newStore(map[string]string{keyStatus: "Draft"})
				b := &readBinding{store: s, delay: time.Hour}
				e := accessor.NewExecutor(registryOf(readerDef(b, keyStatus)), artifactsOf())
				return dispositionOfRead(e.Read(ctxOf(t), readerName))
			},
			want: disposition{refused: true, class: accessor.ClassTimeout},
		},
		"accessor_name_not_bound": {
			run: func(t *testing.T) disposition {
				s := newStore(map[string]string{keyStatus: "Draft"})
				e, _ := readerExec(t, s, keyStatus)
				return dispositionOfRead(e.Read(ctxOf(t), "nope.read"))
			},
			want: disposition{refused: true, class: accessor.ClassUnknownAccessor},
		},
		"accessor_used_off_capability": {
			run: func(t *testing.T) disposition {
				s := newStore(map[string]string{keyStatus: "Draft"})
				e := accessor.NewExecutor(
					registryOf(writerDef(&writeBinding{store: s}, keyStatus)), artifactsOf())
				return dispositionOfRead(e.Read(ctxOf(t), writerName))
			},
			want: disposition{refused: true, class: accessor.ClassCapabilityMismatch},
		},
		"artifact_role_not_supplied": {
			run: func(t *testing.T) disposition {
				s := newStore(map[string]string{keyStatus: "Draft"})
				e := accessor.NewExecutor(
					registryOf(readerDef(&readBinding{store: s}, keyStatus)),
					accessor.Artifacts{})
				return dispositionOfRead(e.Read(ctxOf(t), readerName))
			},
			want: disposition{refused: true, class: accessor.ClassExecutionFailure},
		},
		"write_succeeds_owned_tag_matches": {
			run: func(t *testing.T) disposition {
				s := newStore(map[string]string{keyStatus: "Draft"})
				e := writeExec(t, s, &writeBinding{store: s}, &readBinding{store: s}, keyStatus)
				return dispositionOfWrite(e.Write(ctxOf(t), writerName,
					planWriting(resolve.Tag{Key: keyStatus, Value: "Final"})))
			},
			want: disposition{refused: false},
		},
		"write_succeeds_owned_tag_differs": {
			run: func(t *testing.T) disposition {
				s := newStore(map[string]string{keyStatus: "Draft"})
				w := &writeBinding{store: s, corrupt: func(st *store) {
					st.tags[keyStatus] = "corrupt"
				}}
				e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)
				return dispositionOfWrite(e.Write(ctxOf(t), writerName,
					planWriting(resolve.Tag{Key: keyStatus, Value: "Final"})))
			},
			want: disposition{refused: true, class: accessor.ClassReadBackMismatch},
		},
		"write_clears_a_held_key": {
			run: func(t *testing.T) disposition {
				s := newStore(map[string]string{keyStatus: "Draft"})
				e := writeExec(t, s, &writeBinding{store: s}, &readBinding{store: s}, keyStatus)
				return dispositionOfWrite(e.Write(ctxOf(t), writerName,
					planWriting(resolve.Tag{Key: keyStatus, Value: accessor.ClearSentinel})))
			},
			want: disposition{refused: false},
		},
		"write_clears_a_key_not_held": {
			run: func(t *testing.T) disposition {
				s := newStore(map[string]string{keyProfile: "large"})
				e := writeExec(t, s, &writeBinding{store: s}, &readBinding{store: s}, keyStatus)
				return dispositionOfWrite(e.Write(ctxOf(t), writerName,
					planWriting(resolve.Tag{Key: keyStatus, Value: accessor.ClearSentinel})))
			},
			want: disposition{refused: false},
		},
		"read_yields_the_literal_clear": {
			run: func(t *testing.T) disposition {
				s := newStore(map[string]string{keyStatus: accessor.ClearSentinel})
				e, _ := readerExec(t, s, keyStatus)
				return dispositionOfRead(e.Read(ctxOf(t), readerName))
			},
			want: disposition{refused: true, class: accessor.ClassIncompleteRead},
		},
		"write_read_back_cannot_read_a_compared_key": {
			run: func(t *testing.T) disposition {
				s := newStore(map[string]string{keyStatus: "Draft"})
				w := &writeBinding{store: s}
				e := accessor.NewExecutor(
					registryOf(readerDef(&readBackFailure{store: s, failOn: keyStatus}, keyStatus),
						writerDef(w, keyStatus)), artifactsOf())
				return dispositionOfWrite(e.Write(ctxOf(t), writerName,
					planWriting(resolve.Tag{Key: keyStatus, Value: "Final"})))
			},
			want: disposition{refused: true, class: accessor.ClassReadBackIncomplete},
		},
		"write_command_runs_then_the_invocation_times_out": {
			run: func(t *testing.T) disposition {
				s := newStore(map[string]string{keyStatus: "Draft"})
				w := &writeBinding{store: s, delay: time.Hour}
				e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)
				return dispositionOfWrite(e.Write(ctxOf(t), writerName,
					planWriting(resolve.Tag{Key: keyStatus, Value: "Final"})))
			},
			want: disposition{refused: true, class: accessor.ClassTimeout},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := c.run(t)
			if got.refused != c.want.refused {
				t.Fatalf("refused = %v (class %q); want refused = %v (class %q) — "+
					"no input class exits silently at this boundary",
					got.refused, got.class, c.want.refused, c.want.class)
			}
			if got.class != c.want.class {
				t.Errorf("refusal class = %q; want %q", got.class, c.want.class)
			}
		})
	}
}

// REQ-96: "The executor is not a shell runner and not a state-machine
// action callback surface. It invokes typed bindings selected by accessor
// name and capability, enforces a per-accessor timeout, converts
// execution errors into stable refusal classes, and never prints
// directly."
// REQ-97: "The resolver stays stateless: it receives the accessor-read
// values and write plan disposition" — the executor MUST NOT make the
// resolver stateful or orchestrating.
// BOUNDARY
//
// The typed-binding property is asserted structurally: a definition's
// reach into the world is a `Binding` value, and there is no command
// string, script, or host-callback field for authority to hide in.
func TestReq96_AuthorityTravelsAsATypedBindingNotACommandString(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	def := readerDef(&readBinding{store: s}, keyStatus)

	if def.Binding == nil {
		t.Fatal("the definition carries no typed binding")
	}
	if def.Binding.Capability() != accessor.CapRead {
		t.Errorf("binding capability = %q; want %q — the executor selects by "+
			"name and capability, never by an executable name",
			def.Binding.Capability(), accessor.CapRead)
	}
	// RDR 0002 carries `path` as artifact location, not as a command to
	// run: the binding, not the path, is what executes.
	if def.Accessor.Path != statePath {
		t.Errorf("definition path = %q; want the artifact path %q",
			def.Accessor.Path, statePath)
	}
}

// REQ-98: "No new third-party dependency is selected at Propose. Resolve
// may choose a Go test helper or TOML library only if RDR 0002 has not
// already selected one." — RDR 0002 has already selected the TOML
// carrier, so no new dependency is expected.
// REQ-100: "Accessor/table path discovery belongs to CLI integration." —
// the executor does not load config.
// BOUNDARY
//
// The executor's inputs are a validated registry and caller-supplied
// artifacts. It receives no path to discover and no config to load.
func TestReq100_TheExecutorLoadsNoConfigAndDiscoversNoPath(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	// Constructed entirely from caller-supplied values: no path lookup,
	// no config load, no ambient discovery.
	b := &readBinding{store: s}
	e := accessor.NewExecutor(registryOf(readerDef(b, keyStatus)), artifactsOf())

	got := e.Read(ctxOf(t), readerName)

	values := mustReadSucceed(t, got)
	// Positive precondition: the read reached the caller-supplied
	// artifact and returned its value — no path was discovered, and none
	// needed to be.
	if b.reads != 1 {
		t.Fatalf("the read binding ran %d times; want exactly 1", b.reads)
	}
	if v, ok := valueOf(values, keyStatus); !ok || v.Value != "Draft" {
		t.Errorf("read returned %+v; want %q = %q from the caller-supplied "+
			"artifact, reached without loading config or discovering a path",
			values, keyStatus, "Draft")
	}
}

// --- helpers -------------------------------------------------------------

// disposition is one boundary outcome reduced to the Disposition
// Table's two columns: which branch was taken and, on a refusal, the
// class minted.
type disposition struct {
	refused bool
	class   accessor.RefusalClass
}

func dispositionOfRead(r accessor.ReadResult) disposition {
	if r.Refusal == nil {
		return disposition{}
	}
	return disposition{refused: true, class: r.Refusal.Class}
}

func dispositionOfWrite(r accessor.WriteResult) disposition {
	if r.Refusal == nil {
		return disposition{}
	}
	return disposition{refused: true, class: r.Refusal.Class}
}

// captureOutput redirects the process's stdout and stderr around fn and
// returns what was written to each. This is ORA 5's normative capture
// control.
func captureOutput(t *testing.T, fn func()) (stdout, stderr string) {
	t.Helper()

	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	errR, errW, err := os.Pipe()
	if err != nil {
		t.Fatalf("stderr pipe: %v", err)
	}

	origOut, origErr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = outW, errW

	outCh := make(chan string, 1)
	errCh := make(chan string, 1)
	go func() {
		var b bytes.Buffer
		_, _ = io.Copy(&b, outR)
		outCh <- b.String()
	}()
	go func() {
		var b bytes.Buffer
		_, _ = io.Copy(&b, errR)
		errCh <- b.String()
	}()

	func() {
		defer func() {
			os.Stdout, os.Stderr = origOut, origErr
			_ = outW.Close()
			_ = errW.Close()
		}()
		fn()
	}()

	return <-outCh, <-errCh
}
