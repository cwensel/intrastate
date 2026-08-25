package accessor

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/newcoinc/intrastate/internal/resolve"
)

// Executor invokes validated accessor definitions against
// caller-supplied artifacts. It holds no ambient state and performs no
// scheduling: the caller decides WHEN a gate runs (JDR 0001 §D9 lands the
// gate SITE in RDR 0005; req-list Q2).
type Executor struct {
	Registry  Registry
	Artifacts Artifacts
}

// NewExecutor is THE executor constructor.
func NewExecutor(reg Registry, arts Artifacts) *Executor {
	return &Executor{Registry: reg, Artifacts: arts}
}

// --- selection -----------------------------------------------------------

// selects resolves one invocation to its definition and caller-supplied
// artifact, or to the refusal the failed selection mints. An unbound name
// is `unknown_accessor`; a name bound under another capability table is
// `capability_mismatch`; an unsupplied role is `execution_failure`, which
// is NOT a separate artifact-unavailability class (`0004:FM`).
func (e *Executor) selects(name string, capability Capability) (
	Definition, Artifact, time.Duration, *Refusal,
) {
	def, ok := e.Registry.Lookup(name, capability)
	if !ok {
		class := ClassUnknownAccessor
		if e.Registry.bound(name) {
			class = ClassCapabilityMismatch
		}
		return Definition{}, Artifact{}, 0, &Refusal{
			Class:      class,
			Accessor:   name,
			Capability: capability,
		}
	}

	timeout, ok := def.timeout()
	if !ok {
		// Validation rejects this before execution; at runtime a
		// definition that slipped through cannot be bounded, and an
		// unbounded invocation is not something this boundary performs.
		return def, Artifact{}, 0, refusalOf(def, 0, ClassExecutionFailure)
	}

	art, ok := e.Artifacts[def.Accessor.Role]
	if !ok {
		return def, Artifact{}, timeout, refusalOf(def, timeout, ClassExecutionFailure)
	}
	return def, art, timeout, nil
}

// refusalOf builds the diagnosis tuple every refusal carries: accessor
// identity, capability, artifact role, and declared timeout (`0004:FM`).
func refusalOf(def Definition, timeout time.Duration, class RefusalClass) *Refusal {
	return &Refusal{
		Class:      class,
		Accessor:   def.Identity.Name,
		Capability: def.Identity.Capability,
		Role:       def.Accessor.Role,
		Timeout:    timeout,
	}
}

// --- read ----------------------------------------------------------------

// Read invokes the named read accessor. The requested key set comes from
// the definition's validated metadata, never from the keys the read
// happened to resolve (`0004:C5`).
func (e *Executor) Read(ctx context.Context, name string) ReadResult {
	def, art, timeout, refusal := e.selects(name, CapRead)
	if refusal != nil {
		return ReadResult{Refusal: refusal}
	}

	// The requested set is the DEFINITION's declared metadata. Deriving
	// it from what came back makes completeness self-fulfilling.
	requested := def.RequestedKeys()

	raw := e.invokeRead(ctx, def, art, timeout, requested)
	if raw.class != "" {
		return ReadResult{Refusal: refusalWithKeys(def, timeout, raw.class, raw.unreadable)}
	}

	// A value that reads back as the reserved `<clear>` literal did not
	// establish the key's content: the binding cannot tell a removal from
	// a stored literal, so the key is UNREADABLE (LBD, Absent vs
	// unreadable; `0004:C11`). The read-back comparison deliberately does
	// not apply this rule — there the literal's presence is the defect.
	values, unread := raw.classify(requested, true)
	if len(unread) != 0 {
		return ReadResult{Refusal: refusalWithKeys(def, timeout, ClassIncompleteRead, unread)}
	}
	return ReadResult{Values: values}
}

func refusalWithKeys(def Definition, timeout time.Duration, class RefusalClass, keys []string) *Refusal {
	r := refusalOf(def, timeout, class)
	r.Keys = slices.Clone(keys)
	return r
}

// readOutcome is one raw binding invocation: what it resolved, what it
// could not read, and the refusal class the invocation itself minted.
type readOutcome struct {
	values     []KeyValue
	unreadable []string
	class      RefusalClass
}

// classify reduces a raw outcome to exactly the requested keys plus the
// keys that did not resolve. A key the binding classified as neither read
// nor unreadable defaults to UNREADABLE: guessing absence is the failure
// this contract exists to prevent (`0004:C6`, LBD Absent vs unreadable).
func (o readOutcome) classify(requested []string, clearIsUnreadable bool) (
	values []KeyValue, unread []string,
) {
	held := make(map[string]KeyValue, len(o.values))
	for _, v := range o.values {
		held[v.Key] = v
	}
	for _, key := range requested {
		if slices.Contains(o.unreadable, key) {
			unread = append(unread, key)
			continue
		}
		v, ok := held[key]
		if !ok {
			unread = append(unread, key)
			continue
		}
		if clearIsUnreadable && !v.Absent && IsClear(v.Value) {
			unread = append(unread, key)
			continue
		}
		values = append(values, KeyValue{Key: key, Value: v.Value, Absent: v.Absent})
	}
	return values, unread
}

// invokeRead performs one bounded read invocation. Every accessor
// invocation has a bounded timeout, the read-back re-read included
// (`0004:C15`); `timeout` outranks `incomplete_read` when both are true
// at once (`0004:C7`).
func (e *Executor) invokeRead(
	ctx context.Context, def Definition, art Artifact,
	timeout time.Duration, requested []string,
) readOutcome {
	binding, ok := def.Binding.(ReadBinding)
	if !ok {
		return readOutcome{class: ClassExecutionFailure}
	}

	bounded, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	values, unreadable, err := binding.Read(bounded, art, slices.Clone(requested))

	// Deadline first: a read that met an unreadable key AND ran past its
	// deadline makes both classes true, and `timeout` is the one reported.
	if errors.Is(bounded.Err(), context.DeadlineExceeded) {
		return readOutcome{class: ClassTimeout}
	}
	if err != nil {
		return readOutcome{class: ClassExecutionFailure}
	}
	return readOutcome{values: values, unreadable: unreadable}
}

// --- gate ----------------------------------------------------------------

// Gate invokes the named gate accessor. It REPORTS the verdict and never
// APPLIES a deny (JDR 0001 §D9; REQ-38). The gate SITE — when gates run,
// which gates run, and how several aggregate — is RDR 0005's.
func (e *Executor) Gate(ctx context.Context, name string) GateResult {
	def, art, timeout, refusal := e.selects(name, CapGate)
	if refusal != nil {
		return GateResult{Refusal: refusal}
	}

	binding, ok := def.Binding.(GateBinding)
	if !ok {
		return GateResult{Refusal: refusalOf(def, timeout, ClassExecutionFailure)}
	}

	bounded, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	verdict, reason, err := binding.Gate(bounded, art)

	// A gate's timeout or execution failure is an ACCESSOR refusal, not a
	// gate result — never folded into `gate_indeterminate` (JDR 0001 §D9).
	if errors.Is(bounded.Err(), context.DeadlineExceeded) {
		return GateResult{Refusal: refusalOf(def, timeout, ClassTimeout)}
	}
	if err != nil {
		return GateResult{Refusal: refusalOf(def, timeout, ClassExecutionFailure)}
	}

	switch verdict {
	case VerdictAllow, VerdictDeny:
		// Deny is a TYPED result carrying its reason at this boundary; it
		// is a refusal at the CLI (JDR 0001 §D9).
		return GateResult{Verdict: verdict, Reason: reason}
	case VerdictIndeterminate:
		// Refusal-class, never a false allow and never a false deny.
		r := refusalOf(def, timeout, ClassGateIndeterminate)
		r.Reason = reason
		return GateResult{Verdict: VerdictIndeterminate, Reason: reason, Refusal: r}
	default:
		return GateResult{Refusal: refusalOf(def, timeout, ClassExecutionFailure)}
	}
}

// --- write ---------------------------------------------------------------

// Write invokes the named write accessor for a SUCCESSFUL transition
// plan, then re-reads the same caller-supplied artifact role named by the
// write binding and verifies read-back (`0004:C12`).
//
// The plan's Writes are the ONLY tags applied; an escaped plan carries an
// empty write set and no write runs (deviations D1, RDR 0009
// `0009:1515-1527`).
func (e *Executor) Write(ctx context.Context, name string, plan resolve.Plan) WriteResult {
	def, art, timeout, refusal := e.selects(name, CapWrite)
	if refusal != nil {
		return WriteResult{Refusal: refusal}
	}
	if len(plan.Writes) == 0 {
		// An escaped plan reaches the write accessor with an empty write
		// set. Performing no write is not a failure.
		return WriteResult{}
	}

	binding, ok := def.Binding.(WriteBinding)
	if !ok {
		return WriteResult{Refusal: refusalOf(def, timeout, ClassExecutionFailure)}
	}

	planned := slices.Clone(plan.Writes)
	plannedKeys := make([]string, 0, len(planned))
	for _, t := range planned {
		plannedKeys = append(plannedKeys, t.Key)
	}

	// `0004:C10` binds what the accessor APPLIES, not only what a
	// definition declares. Validation's `write_non_owned_tag` arm inspects
	// declared metadata; the PLAN is caller-supplied at this boundary, so
	// a plan naming an observed or recognized tag would otherwise reach
	// `Apply` unfiltered and mutate the artifact under a green read-back
	// (REQ-40, REQ-41, REQ-80; FM lists "write attempted for a non-owned
	// tag" among the typed refusals). The check runs BEFORE the command:
	// the write must not reach the artifact at all (deviations D16).
	if nonOwned := nonOwnedPlanKeys(def, e.Registry, plannedKeys); len(nonOwned) != 0 {
		r := refusalOf(def, timeout, ClassExecutionFailure)
		r.Keys = nonOwned
		r.Expected = planned
		return WriteResult{Refusal: r}
	}

	// The read-back re-read goes through the READ path over the SAME
	// caller-supplied role the write binding names — never an ambient
	// artifact and never an unrelated role (`0004:C12`, `0004:C13`).
	reader, hasReader := e.Registry.readerFor(def.Accessor.Role)

	// Before the write, record the protected non-owned values this
	// boundary can observe: the reader's declared keys, minus the plan's
	// own owned keys, which are excluded from the comparison by
	// construction (FID, REQ-60). An empty protected set means there is
	// nothing to snapshot and no pre-write read runs.
	protected := protectedKeys(reader, hasReader, plannedKeys)
	before := map[string]string{}
	if len(protected) != 0 {
		rt, _ := reader.timeout()
		raw := e.invokeRead(ctx, reader, art, rt, protected)
		if raw.class == "" {
			vals, _ := raw.classify(protected, false)
			for _, v := range vals {
				if !v.Absent {
					// Tags absent before the write are unconstrained.
					before[v.Key] = v.Value
				}
			}
		}
	}

	// --- the write command ---------------------------------------------
	applyCtx, cancelApply := context.WithTimeout(ctx, timeout)
	err := binding.Apply(applyCtx, art, planned)
	appliedDeadline := errors.Is(applyCtx.Err(), context.DeadlineExceeded)
	cancelApply()

	if appliedDeadline {
		// The command already ran: "may have been applied and was not
		// verified", never "the write did not occur" (`0004:C14`).
		r := refusalOf(def, timeout, ClassTimeout)
		r.applied = true
		r.Expected = planned
		return WriteResult{Refusal: r}
	}
	if err != nil {
		// The command failed before mutating: this one did NOT occur.
		return WriteResult{Refusal: refusalOf(def, timeout, ClassExecutionFailure)}
	}

	// --- read-back verification ----------------------------------------
	if !hasReader {
		r := refusalOf(def, timeout, ClassReadBackIncomplete)
		r.applied = true
		r.Keys = plannedKeys
		r.Expected = planned
		return WriteResult{Refusal: r}
	}

	compared := slices.Clone(plannedKeys)
	for _, k := range protected {
		if !slices.Contains(compared, k) {
			compared = append(compared, k)
		}
	}
	readTimeout, _ := reader.timeout()
	raw := e.invokeRead(ctx, reader, art, readTimeout, compared)
	if raw.class == ClassTimeout {
		r := refusalOf(def, readTimeout, ClassTimeout)
		r.applied = true
		r.Expected = planned
		return WriteResult{Refusal: r}
	}
	if raw.class != "" {
		r := refusalOf(def, readTimeout, ClassReadBackIncomplete)
		r.applied = true
		r.Keys = compared
		r.Expected = planned
		return WriteResult{Refusal: r}
	}

	// The re-read is itself subject to read completeness. A key it cannot
	// read is `read_back_incomplete` — the verification did not run —
	// never `read_back_mismatch`, which asserts the artifact is wrong
	// (`0004:C13`). The `<clear>`-literal rule is NOT applied here: a
	// stored literal is exactly the presence a clear's read-back catches.
	observedValues, unread := raw.classify(compared, false)
	if len(unread) != 0 {
		r := refusalOf(def, readTimeout, ClassReadBackIncomplete)
		r.applied = true
		r.Keys = unread
		r.Expected = planned
		return WriteResult{Refusal: r}
	}

	observed := map[string]KeyValue{}
	for _, v := range observedValues {
		observed[v.Key] = v
	}

	if mismatch := verifyReadBack(planned, before, observed); mismatch {
		r := refusalOf(def, readTimeout, ClassReadBackMismatch)
		r.Expected = planned
		r.Observed = observedTags(observedValues)
		return WriteResult{Refusal: r}
	}

	return WriteResult{Written: planned}
}

// nonOwnedPlanKeys names the planned keys this write accessor has no
// authority to apply: a key outside the writer definition's own `keys`,
// or one the model does not carry as an owned tag. Both are required —
// the definition bounds what THIS accessor writes, and `Registry.OwnedTags`
// bounds what is an owned tag at all, which is what separates an owned tag
// from an observed or recognized one (`0004:C10`, REQ-40, REQ-41).
func nonOwnedPlanKeys(def Definition, reg Registry, plannedKeys []string) []string {
	owned := def.OwnedKeys()
	var out []string
	for _, k := range plannedKeys {
		if slices.Contains(out, k) {
			continue
		}
		if !slices.Contains(owned, k) || !slices.Contains(reg.OwnedTags, k) {
			out = append(out, k)
		}
	}
	return out
}

// protectedKeys is the pre-write snapshot's key set: the reader's
// declared keys minus the plan's own owned keys. RDR 0004 requires the
// executor to record "the same caller-supplied artifact role's observed
// and recognized tag values" before the write; the reader's declared
// `keys` is the widest set this boundary can observe, and the planned
// owned keys are excluded from the non-owned comparison by construction
// (TD; FID `write -> read`, non-owned tags; REQ-53, REQ-60).
func protectedKeys(reader Definition, hasReader bool, plannedKeys []string) []string {
	if !hasReader {
		return nil
	}
	var out []string
	for _, k := range reader.RequestedKeys() {
		if slices.Contains(plannedKeys, k) || slices.Contains(out, k) {
			continue
		}
		out = append(out, k)
	}
	return out
}

// verifyReadBack is the round-trip invariant: every planned owned tag
// held EXACTLY (absence for a `<clear>`), and every protected non-owned
// value present before the write unchanged (`0004:C12`, RT, FID).
func verifyReadBack(planned []resolve.Tag, before map[string]string, observed map[string]KeyValue) bool {
	for _, t := range planned {
		got, ok := observed[t.Key]
		if IsClear(t.Value) {
			// A clear is a REMOVAL: its read-back asserts ABSENCE. A
			// re-read still holding the key — including as the literal
			// `<clear>` — is a mismatch.
			if ok && !got.Absent {
				return true
			}
			continue
		}
		// A write replaces the whole value, so containment is not
		// equality; a set crosses as RDR 0002's canonical JSON array, so
		// this equality is byte equality (JDR 0001 §D7(iv), §D13).
		if !ok || got.Absent || got.Value != t.Value {
			return true
		}
	}
	for key, want := range before {
		got, ok := observed[key]
		if !ok || got.Absent || got.Value != want {
			return true
		}
	}
	return false
}

func observedTags(values []KeyValue) []resolve.Tag {
	out := make([]resolve.Tag, 0, len(values))
	for _, v := range values {
		if v.Absent {
			continue
		}
		out = append(out, resolve.Tag{Key: v.Key, Value: v.Value})
	}
	slices.SortFunc(out, func(a, b resolve.Tag) int {
		switch {
		case a.Key < b.Key:
			return -1
		case a.Key > b.Key:
			return 1
		default:
			return 0
		}
	})
	return out
}
