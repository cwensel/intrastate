package accessor

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
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
		return def, Artifact{}, 0, refusalOf(def, 0, ClassExecutionFailure, nil)
	}

	art, ok := e.Artifacts[def.Accessor.Role]
	if !ok {
		return def, Artifact{}, timeout, refusalOf(def, timeout, ClassExecutionFailure, nil)
	}
	return def, art, timeout, nil
}

// refusalOf builds the diagnosis tuple every refusal carries: accessor
// identity, capability, artifact role, and declared timeout (`0004:FM`),
// plus the bounded stderr tail an invocation error carries (`0025:C4`).
//
// `err` is the offending binding error or nil. A command binding returns
// a typed `*ExecError`, which this `errors.As`-es for its tail; a
// path-backed binding returns an untyped error and carries no tail, and a
// refusal raised without an invocation at all — timeout from `ctx.Err()`,
// read-back, the gate — gets an empty `Detail` by construction rather
// than by each call site remembering to leave it blank.
func refusalOf(def Definition, timeout time.Duration, class RefusalClass, err error) *Refusal {
	r := &Refusal{
		Class:      class,
		Accessor:   def.Identity.Name,
		Capability: def.Identity.Capability,
		Role:       def.Accessor.Role,
		Timeout:    timeout,
	}
	var ee *ExecError
	if errors.As(err, &ee) {
		r.Detail = ee.Detail
	}
	// The one bit RDR 0028's exit-group clause needs, carried on the
	// existing `Err` slot rather than through a wider seam. `declaredEdit`
	// narrows it to the refusals a line edit minted: `ErrDeclaredEdit`
	// WRAPS `ErrDeclaredRequest`, so the group answer is unchanged and
	// only the envelope's choice of subject turns on the narrower bit.
	r.declaredRequest = errors.Is(err, ErrDeclaredRequest)
	r.declaredEdit = errors.Is(err, ErrDeclaredEdit)
	return r
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
		return ReadResult{
			Refusal: refusalWithKeys(def, timeout, raw.class, raw.unreadable, raw.err),
		}
	}

	// A value that reads back as the reserved `<clear>` literal did not
	// establish the key's content: the binding cannot tell a removal from
	// a stored literal, so the key is UNREADABLE (LBD, Absent vs
	// unreadable; `0004:C11`). The read-back comparison deliberately does
	// not apply this rule — there the literal's presence is the defect.
	values, unread := raw.classify(requested, true)
	if len(unread) != 0 {
		return ReadResult{
			Refusal: refusalWithKeys(def, timeout, ClassIncompleteRead, unread, nil),
		}
	}
	return ReadResult{Values: values}
}

// refusalWithKeys wraps refusalOf for the read path, which is the ONE
// capability established tools bind directly and so the likeliest source
// of a stderr tail. The error parameter threads through BOTH constructors
// or reads silently drop their tail (`0025:C4`, A11).
func refusalWithKeys(
	def Definition, timeout time.Duration, class RefusalClass,
	keys []string, err error,
) *Refusal {
	r := refusalOf(def, timeout, class, err)
	r.Keys = slices.Clone(keys)
	return r
}

// readOutcome is one raw binding invocation: what it resolved, what it
// could not read, and the refusal class the invocation itself minted.
type readOutcome struct {
	values     []KeyValue
	unreadable []string
	class      RefusalClass
	// err is the offending binding error, carried so the read path can
	// hand it to `refusalWithKeys` (`0025:C4`).
	err error
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
		return readOutcome{class: ClassExecutionFailure, err: err}
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
		return GateResult{Refusal: refusalOf(def, timeout, ClassExecutionFailure, nil)}
	}

	bounded, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	verdict, reason, err := binding.Gate(bounded, art)

	// A gate's timeout or execution failure is an ACCESSOR refusal, not a
	// gate result — never folded into `gate_indeterminate` (JDR 0001 §D9).
	if errors.Is(bounded.Err(), context.DeadlineExceeded) {
		return GateResult{Refusal: refusalOf(def, timeout, ClassTimeout, nil)}
	}
	if err != nil {
		return GateResult{Refusal: refusalOf(def, timeout, ClassExecutionFailure, err)}
	}

	switch verdict {
	case VerdictAllow, VerdictDeny:
		// Deny is a TYPED result carrying its reason at this boundary; it
		// is a refusal at the CLI (JDR 0001 §D9).
		return GateResult{Verdict: verdict, Reason: reason}
	case VerdictIndeterminate:
		// Refusal-class, never a false allow and never a false deny.
		r := refusalOf(def, timeout, ClassGateIndeterminate, nil)
		r.Reason = reason
		return GateResult{Verdict: VerdictIndeterminate, Reason: reason, Refusal: r}
	default:
		return GateResult{Refusal: refusalOf(def, timeout, ClassExecutionFailure, nil)}
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
		return WriteResult{Refusal: refusalOf(def, timeout, ClassExecutionFailure, nil)}
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
		r := refusalOf(def, timeout, ClassExecutionFailure, nil)
		r.Keys = nonOwned
		r.Expected = planned
		return WriteResult{Refusal: r}
	}

	// The read-back re-read goes through the READ path over the SAME
	// caller-supplied role the write binding names — never an ambient
	// artifact and never an unrelated role (`0004:C12`, `0004:C13`).
	reader, hasReader := e.Registry.readerFor(def.Accessor.Role)

	// RDR 0028 `0028:C1.3` read-back: — the gate pre-check, and the
	// AUTHORITATIVE detector for this case.
	//
	// An `edit` write spawns nothing, so the gate is not consulted by the
	// write itself; but its read-back goes through the role's declared
	// reader, and a command-backed reader with the gate OFF cannot run.
	// Refusing here — before `Apply`, unconditionally — is what stops a
	// forgotten `--allow-commands` from mutating the artifact and then
	// reporting `read_back_incomplete`, which carries the
	// applied-but-unverified sense for a write that ran no command.
	//
	// It does NOT amend `Executor.Write`'s existing pre-`Apply` baseline
	// read below. That baseline runs only when `protected` is non-empty
	// and SWALLOWS its failure into `baselineUnread`; this refuses
	// earlier and unconditionally, so the swallow becomes unreachable for
	// the gate-off case whether or not `protected` is empty (A2).
	//
	// The scope is this carrier's. A command-backed WRITER already
	// refuses at its own pre-spawn ladder (0025:C6) and is not this
	// clause's to re-decide.
	// C1.3 `precedence:` step (1) puts C1.6's two ARGV preconditions —
	// an unbound `{tag.<key>}`, then a `-`-prefixed bound value — ahead
	// of everything else, "decided together before anything is spawned
	// or read". For the read-back reader those refusals otherwise fire
	// inside `invokeRead`, which `Write` reaches only AFTER `Apply` has
	// already mutated the artifact: the caller is then told the write
	// "may have been applied and was not verified" for what is purely a
	// defect of the request. C1.3 `order:` forbids exactly that — "every
	// refusal in this clause and C1.2's `edit_value_multiline` is decided
	// BEFORE any byte is written".
	//
	// The gate arm below cannot stand in for this one: it is reachable
	// only when the gate is OFF, and these refusals require the gate ON
	// (a passed gate is what spawns the child that would raise them).
	// Scoped to this carrier, like the gate arm: a command-backed WRITER
	// keeps its own pre-spawn ladder (0025:C6).
	if len(def.Accessor.Edit) != 0 && hasReader &&
		len(reader.Accessor.Command) != 0 && e.Registry.AllowCommands {
		// TWO passes, not one. The clause orders ALL unbound tags ahead
		// of ANY `-`-prefixed value, so a single argv walk would report
		// whichever defect sits earlier in the vector — reversing the
		// precedence whenever a prefixed value precedes an unbound key.
		// Argv order is the author's; the reported category is the
		// contract's.
		readBackRefusal := func(key, why string) WriteResult {
			r := refusalOf(def, timeout, ClassExecutionFailure, ErrDeclaredRequest)
			r.Detail = "the read-back for the write accessor `" + name +
				"` goes through the command-backed reader `" +
				reader.Identity.Name + "`, whose argv carries `{tag." +
				key + "}` " + why + "; refusing before mutation rather " +
				"than writing and reporting an unverified read-back"
			return WriteResult{Refusal: r}
		}
		for _, el := range reader.Accessor.Command {
			key, ok := table.CommandTagKey(el)
			if !ok {
				continue
			}
			if _, bound := art.Context[key]; !bound {
				return readBackRefusal(key,
					"and this invocation binds no `"+key+"`")
			}
		}
		for _, el := range reader.Accessor.Command {
			key, ok := table.CommandTagKey(el)
			if !ok {
				continue
			}
			if strings.HasPrefix(art.Context[key], "-") {
				return readBackRefusal(key, "bound to a `-`-prefixed value")
			}
		}
	}

	if len(def.Accessor.Edit) != 0 && hasReader &&
		len(reader.Accessor.Command) != 0 && !e.Registry.AllowCommands {
		// `ErrDeclaredRequest`, not nil: this is an ENTRY-level
		// precondition and C1.3's EXIT GROUP: sentence puts all three of
		// them in the exit-2 group. A forgotten `--allow-commands` is a
		// property of the REQUEST — exit 3's "re-run the same request
		// unchanged" is advice that can never succeed here. There is no
		// binding error to wrap, so the sentinel is passed on its own;
		// it is not an `*ExecError`, so it contributes no Detail and the
		// gate text below remains the whole of it.
		r := refusalOf(def, timeout, ClassExecutionFailure, ErrDeclaredRequest)
		// The Detail names the GATE, having no rule to name: this is an
		// ENTRY-level precondition, not a rule-scoped refusal
		// (`0028:C1.3` order:).
		r.Detail = "the read-back for the write accessor `" + name +
			"` goes through the command-backed reader `" +
			reader.Identity.Name + "`, which requires the allow-commands " +
			"opt-in (--allow-commands); refusing before mutation rather " +
			"than writing and reporting an unverified read-back"
		r.Expected = planned
		return WriteResult{Refusal: r}
	}

	// Before the write, record the protected non-owned values this
	// boundary can observe: the reader's declared keys, minus the plan's
	// own owned keys, which are excluded from the comparison by
	// construction (FID, REQ-60). An empty protected set means there is
	// nothing to snapshot and no pre-write read runs.
	protected := protectedKeys(reader, hasReader, plannedKeys)
	before := map[string]string{}
	// baselineUnread names protected keys whose PRE-write value could not
	// be established. Their "unchanged" claim is unverifiable, so the
	// verification did not run for them (`0004:C13`). The write still
	// proceeds and the refusal is raised post-command: `read_back_incomplete`
	// is by contract reported after the write command already ran and
	// carries the applied-but-unverified sense (`0004:C14`, REQ-63).
	var baselineUnread []string
	if len(protected) != 0 {
		rt, _ := reader.timeout()
		raw := e.invokeRead(ctx, reader, art, rt, protected)
		if raw.class != "" {
			// A timed-out, errored, or incomplete pre-write read
			// establishes NO baseline. Discarding it leaves `before`
			// empty, which makes `0004:C12`'s protected-tag clause
			// vacuously true (REQ-53, REQ-80).
			baselineUnread = slices.Clone(protected)
		} else {
			vals, unread := raw.classify(protected, false)
			// A PARTIAL snapshot protects only the keys it read; the
			// unread remainder is unverifiable, not unconstrained.
			baselineUnread = unread
			for _, v := range vals {
				if !v.Absent {
					// Tags absent before the write are unconstrained.
					before[v.Key] = v.Value
				}
			}
		}
	}

	// --- the write command ---------------------------------------------
	// The binding gets its OWN copy. `planned` is the expectation the
	// read-back oracle judges the artifact against, and it is also every
	// refusal's `Expected` and the success record's `Written`. Handing the
	// binding the same backing array would let it rewrite the expectation
	// it is judged against: `verifyReadBack` would compare the artifact to
	// whatever the binding chose, and `Written` would report the mutated
	// value as verified — a self-referential oracle, defeatable by the very
	// component it verifies. Bindings are not trusted at this boundary
	// (`0004:465-474`, `0004:801`); D16 already enforces the caller-supplied
	// plan at runtime rather than trusting validation, and this is the same
	// failure shape one hop later. `resolve.Tag` is a value struct of
	// strings, so a shallow clone fully severs the aliasing.
	applyCtx, cancelApply := context.WithTimeout(ctx, timeout)
	err := binding.Apply(applyCtx, art, slices.Clone(planned))
	appliedDeadline := errors.Is(applyCtx.Err(), context.DeadlineExceeded)
	cancelApply()

	if appliedDeadline {
		// The command already ran: "may have been applied and was not
		// verified", never "the write did not occur" (`0004:C14`).
		r := refusalOf(def, timeout, ClassTimeout, nil)
		r.applied = true
		r.Expected = planned
		return WriteResult{Refusal: r}
	}
	if err != nil {
		// The command failed before mutating: this one did NOT occur.
		return WriteResult{Refusal: refusalOf(def, timeout, ClassExecutionFailure, err)}
	}

	// --- read-back verification ----------------------------------------
	if !hasReader {
		r := refusalOf(def, timeout, ClassReadBackIncomplete, nil)
		r.applied = true
		r.Keys = plannedKeys
		r.Expected = planned
		return WriteResult{Refusal: r}
	}

	// An unestablished pre-write baseline does NOT short-circuit the
	// re-read. `0004:C12` is one conjunctive obligation with two conjuncts
	// — each planned owned tag equals its held value, and the protected
	// non-owned values present before the write are unchanged — and only
	// the second depends on the baseline. `planned` alone supplies the
	// first conjunct's expectation, so it stays verifiable and REQ-56
	// makes the read-back mandatory, not optional. The baseline-incomplete
	// refusal is raised below, AFTER the re-read has had its chance to
	// evaluate the conjunct it can (`0004:C13`, REQ-50, REQ-56).
	compared := slices.Clone(plannedKeys)
	for _, k := range protected {
		if !slices.Contains(compared, k) {
			compared = append(compared, k)
		}
	}
	readTimeout, _ := reader.timeout()
	raw := e.invokeRead(ctx, reader, art, readTimeout, compared)
	if raw.class == ClassTimeout {
		r := refusalOf(def, readTimeout, ClassTimeout, nil)
		r.applied = true
		r.Expected = planned
		return WriteResult{Refusal: r}
	}
	if raw.class != "" {
		// The read-back's OWN invocation error carries the tail: the CLI
		// composes it after the applied-sense text in its one slot
		// (`0025:C4` detail, REQ-65).
		r := refusalOf(def, readTimeout, ClassReadBackIncomplete, raw.err)
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
		r := refusalOf(def, readTimeout, ClassReadBackIncomplete, nil)
		r.applied = true
		r.Keys = unread
		r.Expected = planned
		return WriteResult{Refusal: r}
	}

	observed := map[string]KeyValue{}
	for _, v := range observedValues {
		observed[v.Key] = v
	}

	// Every conjunct of `0004:C12` that IS evaluable, evaluated BEFORE any
	// incompleteness refusal: each planned owned tag, and every protected
	// value the pre-write baseline did establish. `before` holds exactly
	// the latter — a partial snapshot contributes the keys it read and
	// omits the rest — so this one call spans the whole verifiable surface
	// whether the baseline is complete, partial, or wholly unestablished.
	//
	// The re-read that judges these completed over every compared key, so
	// `0004:C13`'s "the verification did not run" predicate is false by
	// construction here and a demonstrated inequality is
	// `read_back_mismatch`. C13 scopes read completeness to the RE-READ
	// ("If IT cannot read a key it must compare"), which the three arms
	// above already enforce; the pre-write baseline is not the re-read, so
	// it may not swallow a conjunct the re-read did evaluate. Ordering
	// matters in both directions: a key in `baselineUnread` must not
	// suppress a mismatch demonstrated on a DIFFERENT key, and REQ-67
	// forbids collapsing "the artifact is wrong" into "unverified".
	//
	// `applied` stays UNSET: `0004:C14` scopes the applied-but-unverified
	// sense to `read_back_incomplete` and a post-mutation `timeout`, and a
	// mismatch is not unverified — its verification ran and found the
	// artifact wrong (REQ-67). No `Keys` either: `Keys` names what could
	// not be READ, and carrying `baselineUnread` here would conflate that
	// with "read and wrong".
	if mismatch := verifyReadBack(planned, before, observed); mismatch {
		r := refusalOf(def, readTimeout, ClassReadBackMismatch, nil)
		r.Expected = planned
		r.Observed = observedTags(observedValues)
		return WriteResult{Refusal: r}
	}

	// No evaluable conjunct failed. A protected key with no established
	// pre-write value still cannot be compared, so `0004:C12`'s "unchanged"
	// clause did not run for it: `read_back_incomplete` — never a mismatch,
	// which asserts the artifact is wrong, and never success (`0004:C13`,
	// REQ-62), carrying the applied-but-unverified sense (`0004:C14`).
	if len(baselineUnread) != 0 {
		r := refusalOf(def, timeout, ClassReadBackIncomplete, nil)
		r.applied = true
		r.Keys = baselineUnread
		r.Expected = planned
		return WriteResult{Refusal: r}
	}

	return WriteResult{Written: verifiedWritten(planned)}
}

// verifiedWritten renders what the read-back actually VERIFIED as held.
// A cleared key is OMITTED: read-back verified it ABSENT, so echoing the
// reserved `<clear>` literal as a written value would record the artifact
// as HOLDING a key it does not hold — the placeholder-for-absence
// encoding this RDR exists to forbid, here in the record rather than at
// the seam (`0004:C11`, REQ-32, REQ-45, REQ-46, REQ-80). It is the same
// rule `ReadResult.OwnedSnapshot` applies to an absent key.
func verifiedWritten(planned []resolve.Tag) []resolve.Tag {
	out := make([]resolve.Tag, 0, len(planned))
	for _, t := range planned {
		if IsClear(t.Value) {
			continue
		}
		out = append(out, t)
	}
	return out
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
