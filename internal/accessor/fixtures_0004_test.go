package accessor_test

// RDR 0004 — shared fixtures for the accessor execution safety suite.
//
// The fixture bindings here are deliberately NOT the Resolve spike's
// shapes. CA A9 names four spike shapes the implementation MUST NOT
// reproduce (REQ-124), and CA A10 names a fifth (REQ-125):
//
//  1. the spike derives its default key set from the artifact —
//     `expectedTagKeys`. Here the requested key set is pinned in the
//     definition's `keys` and the artifact contents never feed it.
//  2. the spike holds absence as an in-map sentinel `<absent>`. Here
//     absence is a typed KeyValue.Absent flag and the seam OMITS the key.
//  3. the spike sleeps before its key loop so timeout and truncation
//     never overlap. `overlapReadBinding` meets the unreadable key AND
//     runs past the deadline in one invocation, so both classes are true
//     at once (REQ-27).
//  4. the spike re-reads by cloning the tag map without going through
//     `read`. Here the write binding and the read-back READER are
//     separate objects, so the re-read can fail independently of the
//     write (REQ-61, REQ-125).

import (
	"context"
	"errors"
	"maps"
	"slices"
	"testing"
	"time"

	"github.com/newcoinc/intrastate/internal/accessor"
	"github.com/newcoinc/intrastate/internal/resolve"
	"github.com/newcoinc/intrastate/internal/table"
)

// --- fixture identities --------------------------------------------------

const (
	flowID = "rdr"

	readerName = "state.read"
	gateName   = "state.gate"
	writerName = "state.persist"

	stateRole = "state"
	statePath = "flows/state.toml"

	keyStatus  = "status"
	keyProfile = "profile"
	keyLabels  = "labels"

	fixtureTimeout = 50 * time.Millisecond
)

// --- the backing store ---------------------------------------------------

// store is one caller-supplied artifact role's tag content plus the
// per-key readability the fixture controls. It is the BINDING's private
// state, never the executor's: the executor reaches it only through the
// typed Binding methods.
type store struct {
	tags map[string]string
	// unreadable names keys whose backing read errors. This is DISTINCT
	// from a key the store legitimately lacks, which reads fine and is
	// established-absent (LBD, Absent vs unreadable).
	unreadable map[string]bool
}

func newStore(tags map[string]string) *store {
	return &store{tags: maps.Clone(tags), unreadable: map[string]bool{}}
}

func (s *store) get(key string) (value string, absent bool, readable bool) {
	if s.unreadable[key] {
		return "", false, false
	}
	v, held := s.tags[key]
	if !held {
		return "", true, true
	}
	return v, false, true
}

// --- read binding --------------------------------------------------------

type readBinding struct {
	store *store
	// delay makes the binding block, so a bounded invocation can expire.
	delay time.Duration
	// failWith makes the whole read an execution failure.
	failWith error
	// reads counts invocations, so a test can assert the re-read is a
	// SECOND, independently-failing invocation.
	reads int
	// lastRequested records the key set the executor asked for, so a
	// test can assert it came from the DEFINITION rather than from the
	// artifact's contents.
	lastRequested []string
}

func (b *readBinding) Capability() accessor.Capability { return accessor.CapRead }

func (b *readBinding) Read(ctx context.Context, art accessor.Artifact, requested []string) (
	[]accessor.KeyValue, []string, error,
) {
	b.reads++
	b.lastRequested = slices.Clone(requested)
	if b.failWith != nil {
		return nil, nil, b.failWith
	}
	if b.delay > 0 {
		select {
		case <-time.After(b.delay):
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
	}
	var values []accessor.KeyValue
	var unreadable []string
	for _, key := range requested {
		v, absent, readable := b.store.get(key)
		if !readable {
			unreadable = append(unreadable, key)
			continue
		}
		values = append(values, accessor.KeyValue{Key: key, Value: v, Absent: absent})
	}
	return values, unreadable, nil
}

// --- gate binding --------------------------------------------------------

type gateBinding struct {
	verdict accessor.Verdict
	reason  string
	delay   time.Duration
	failErr error
}

func (b *gateBinding) Capability() accessor.Capability { return accessor.CapGate }

func (b *gateBinding) Gate(ctx context.Context, art accessor.Artifact) (
	accessor.Verdict, string, error,
) {
	if b.failErr != nil {
		return "", "", b.failErr
	}
	if b.delay > 0 {
		select {
		case <-time.After(b.delay):
		case <-ctx.Done():
			return "", "", ctx.Err()
		}
	}
	return b.verdict, b.reason, nil
}

// --- write binding -------------------------------------------------------

// writeBinding applies planned tags to a store. It is a SEPARATE object
// from the read binding that performs the read-back, which is what lets a
// fixture make the re-read fail independently of the write (REQ-125).
type writeBinding struct {
	store *store
	// corrupt runs after the planned mutation, so the fixture can move an
	// owned or a non-owned tag out from under read-back.
	corrupt func(*store)
	// storeLiteral makes a `<clear>` write ASSIGN the literal instead of
	// removing the key — the defect `0004:C11` exists to catch.
	storeLiteral bool
	// rewritePlan makes the binding MUTATE the planned slice it was handed,
	// in place, before applying it — the shape that lets a binding rewrite
	// the expectation the read-back oracle judges it against, if the
	// executor hands over its own backing array.
	rewritePlan func([]resolve.Tag)
	// delay makes the write's own invocation expire AFTER the mutation.
	delay time.Duration
	// failErr makes the write command itself fail, before any mutation.
	failErr error

	applies int
}

func (b *writeBinding) Capability() accessor.Capability { return accessor.CapWrite }

func (b *writeBinding) Apply(ctx context.Context, art accessor.Artifact, planned []resolve.Tag) error {
	b.applies++
	if b.failErr != nil {
		return b.failErr
	}
	if b.rewritePlan != nil {
		b.rewritePlan(planned)
	}
	for _, t := range planned {
		switch {
		case t.Value == table.ClearSentinel && b.storeLiteral:
			b.store.tags[t.Key] = table.ClearSentinel
		case t.Value == table.ClearSentinel:
			delete(b.store.tags, t.Key)
		default:
			b.store.tags[t.Key] = t.Value
		}
	}
	if b.corrupt != nil {
		b.corrupt(b.store)
	}
	if b.delay > 0 {
		select {
		case <-time.After(b.delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func (b *writeBinding) Invocations() int { return b.applies }

// --- definition builders -------------------------------------------------

func readerDef(b accessor.Binding, keys ...string) accessor.Definition {
	return accessor.Definition{
		Identity: accessor.Identity{Flow: flowID, Name: readerName, Capability: accessor.CapRead},
		Accessor: table.Accessor{
			Role:    stateRole,
			Path:    statePath,
			Keys:    slices.Clone(keys),
			Timeout: fixtureTimeout.String(),
		},
		Binding: b,
	}
}

func gateDef(b accessor.Binding) accessor.Definition {
	return accessor.Definition{
		Identity: accessor.Identity{Flow: flowID, Name: gateName, Capability: accessor.CapGate},
		Accessor: table.Accessor{
			Role:    stateRole,
			Path:    statePath,
			Keys:    []string{keyStatus},
			Timeout: fixtureTimeout.String(),
		},
		Binding: b,
	}
}

func writerDef(b accessor.Binding, keys ...string) accessor.Definition {
	return accessor.Definition{
		Identity: accessor.Identity{Flow: flowID, Name: writerName, Capability: accessor.CapWrite},
		Accessor: table.Accessor{
			Role:     stateRole,
			Path:     statePath,
			Keys:     slices.Clone(keys),
			Timeout:  fixtureTimeout.String(),
			ReadBack: true,
		},
		Binding: b,
	}
}

// registryOf builds a validated-shape registry over the given definitions.
func registryOf(defs ...accessor.Definition) accessor.Registry {
	return accessor.Registry{
		Flow:        flowID,
		Definitions: defs,
		OwnedTags:   []string{keyStatus, keyLabels},
	}
}

func artifactsOf() accessor.Artifacts {
	return accessor.Artifacts{stateRole: {Role: stateRole, Path: statePath}}
}

// --- executor builders ---------------------------------------------------

// readerExec wires one read accessor over a store carrying tags, with the
// requested key set PINNED in the definition independently of the store's
// contents (REQ-16, REQ-18).
func readerExec(t *testing.T, s *store, requested ...string) (*accessor.Executor, *readBinding) {
	t.Helper()
	b := &readBinding{store: s}
	e := accessor.NewExecutor(registryOf(readerDef(b, requested...)), artifactsOf())
	return e, b
}

func gateExec(t *testing.T, b *gateBinding) *accessor.Executor {
	t.Helper()
	return accessor.NewExecutor(registryOf(gateDef(b)), artifactsOf())
}

// writeExec wires a write accessor whose read-back re-read goes through a
// SEPARATE read binding over the same store. The two are distinct objects
// on purpose: a re-read that cannot fail independently of the write has
// nothing to report (CA A10, REQ-125).
//
// The read-back reader declares the writer's owned keys PLUS the
// protected non-owned `profile`. `ReadBinding.Read` resolves exactly the
// keys it is asked for, so a non-owned tag the reader never declares is
// invisible at this boundary and `0004:C12`'s protected non-owned
// identity clause could not be witnessed at all (deviations D14).
func writeExec(t *testing.T, s *store, w *writeBinding, r *readBinding, owned ...string) *accessor.Executor {
	t.Helper()
	wd := writerDef(w, owned...)
	readKeys := slices.Clone(owned)
	if !slices.Contains(readKeys, keyProfile) {
		readKeys = append(readKeys, keyProfile)
	}
	rd := readerDef(r, readKeys...)
	return accessor.NewExecutor(registryOf(rd, wd), artifactsOf())
}

// --- plans ---------------------------------------------------------------

func planWriting(writes ...resolve.Tag) resolve.Plan {
	return resolve.Plan{
		RuleID:        "rdr.advance",
		SourceLocator: statePath + ":10",
		Revision:      "rev-0004",
		Writes:        slices.Clone(writes),
		NextTags:      slices.Clone(writes),
	}
}

// planAdvancing builds a plan whose NEXT-state tags and WRITE set differ.
//
// `planWriting` mirrors the two, which is the ordinary case the kernel
// produces, and no oracle built on it can tell the two fields apart: an
// executor applying `NextTags` passes every such test. RDR 0009's
// obligation (`0009:1515-1527`, carried here as deviations D1, Type
// TEST-FIXTURE) is that only `Writes` reaches the write binding —
// `NextTags` names the state the row transitions TO, which is the
// resolver's record, not an instruction to mutate the artifact
// (REQ-40, REQ-43).
func planAdvancing(writes, nextTags []resolve.Tag) resolve.Plan {
	return resolve.Plan{
		RuleID:        "rdr.advance",
		SourceLocator: statePath + ":10",
		Revision:      "rev-0004",
		Writes:        slices.Clone(writes),
		NextTags:      slices.Clone(nextTags),
	}
}

func escapedPlan() resolve.Plan {
	return resolve.Plan{
		RuleID:        "rdr.escape",
		SourceLocator: statePath + ":99",
		Revision:      "rev-0004",
		Escaped:       true,
	}
}

// --- assertion helpers ---------------------------------------------------

func mustReadRefuse(t *testing.T, r accessor.ReadResult, want accessor.RefusalClass) accessor.Refusal {
	t.Helper()
	if r.Refusal == nil {
		t.Fatalf("read succeeded with values %+v; want the refusal branch with class %q",
			r.Values, want)
	}
	if r.Refusal.Class != want {
		t.Fatalf("refusal class = %q; want %q", r.Refusal.Class, want)
	}
	return *r.Refusal
}

func mustReadSucceed(t *testing.T, r accessor.ReadResult) []accessor.KeyValue {
	t.Helper()
	if r.Refusal != nil {
		t.Fatalf("read refused with class %q; want the typed-values branch",
			r.Refusal.Class)
	}
	return r.Values
}

func mustWriteRefuse(t *testing.T, r accessor.WriteResult, want accessor.RefusalClass) accessor.Refusal {
	t.Helper()
	if r.Refusal == nil {
		t.Fatalf("write succeeded with %+v; want the refusal branch with class %q",
			r.Written, want)
	}
	if r.Refusal.Class != want {
		t.Fatalf("refusal class = %q; want %q", r.Refusal.Class, want)
	}
	return *r.Refusal
}

func mustWriteSucceed(t *testing.T, r accessor.WriteResult) {
	t.Helper()
	if r.Refusal != nil {
		t.Fatalf("write refused with class %q (applied=%v); want success",
			r.Refusal.Class, r.Refusal.Applied())
	}
}

// valueOf returns the read result's value for key and whether the key is
// carried at all.
func valueOf(values []accessor.KeyValue, key string) (accessor.KeyValue, bool) {
	for _, v := range values {
		if v.Key == key {
			return v, true
		}
	}
	return accessor.KeyValue{}, false
}

// keysOf renders the key set of a read result's values, sorted.
func keysOf(values []accessor.KeyValue) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		out = append(out, v.Key)
	}
	slices.Sort(out)
	return out
}

// tagKeysOf renders the key set of a seam snapshot, sorted.
func tagKeysOf(tags []resolve.Tag) []string {
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		out = append(out, t.Key)
	}
	slices.Sort(out)
	return out
}

func seamValueOf(tags []resolve.Tag, key string) (string, bool) {
	for _, t := range tags {
		if t.Key == key {
			return t.Value, true
		}
	}
	return "", false
}

func findingCodes(fs []accessor.Finding) []string {
	out := make([]string, 0, len(fs))
	for _, f := range fs {
		out = append(out, string(f.Code))
	}
	slices.Sort(out)
	return out
}

func hasCode(fs []accessor.Finding, code accessor.ValidationCode) bool {
	for _, f := range fs {
		if f.Code == code {
			return true
		}
	}
	return false
}

var errBindingFailed = errors.New("fixture binding failed")

// callerCtxKey types the sentinel `ctxOf` plants in every fixture
// context. Its own type makes the key unforgeable: nothing but `ctxOf`
// can put this value in a context, so a binding that reads it back has
// necessarily been handed a context DESCENDED from the caller's.
type callerCtxKey struct{}

// callerCtxValue is what `ctxOf` plants and `ctxWitnessReadBinding`
// looks for.
const callerCtxValue = "rdr-0004-caller-context"

// ctxOf is the caller context every fixture invocation passes in.
//
// It carries a sentinel value rather than being a bare `Background()`:
// "the invocation is context-bound" means the CALLER'S context reaches
// the binding, and an executor that discarded it and built a fresh
// `context.WithTimeout(context.Background(), timeout)` would still hand
// the binding a non-nil context carrying a deadline. Only a value that
// could have come from nowhere else distinguishes the two, and it is
// what makes a caller-side cancellation observable at the seam
// (REQ-92 / IP Phase 2).
func ctxOf(t *testing.T) context.Context {
	t.Helper()
	return context.WithValue(context.Background(), callerCtxKey{}, callerCtxValue)
}
