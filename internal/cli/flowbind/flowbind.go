// Package flowbind is the production `accessor.Binding` implementation the
// `flow` command group invokes (RDR 0005, JDR 0001 §D8/§D11).
//
// RDR 0004 fixed the typed invocation seam and its execution safety rules
// but shipped no binding: until now the only implementations were that
// RDR's test fixtures. This package supplies one, and nothing more. It is
// deliberately the thinnest thing that satisfies the seam:
//
//   - a READ binding resolves a declared key set from one caller-supplied
//     artifact, reporting each key as present-with-value or
//     established-absent;
//   - a WRITE binding applies planned owned-tag values to that artifact,
//     treating the reserved `<clear>` value as a REMOVAL (`0004:C11`);
//   - a GATE binding answers allow, deny, or indeterminate.
//
// # The artifact format
//
// The on-disk artifact is a flat JSON object of tag key to tag value, both
// strings, written with the same non-HTML-escaping encoder every other CLI
// wire site uses. RDR 0005 fixes the CLI contract, not a file schema, and
// no REQ constrains this — the format is an IMPL-DECISION recorded in the
// deviations artifact.
//
// Two properties of the format ARE load-bearing, because contract clauses
// depend on them:
//
//  1. Key PRESENCE is what distinguishes an empty set from a cleared key.
//     `labels` holding `[]` is a present key whose value is the canonical
//     empty array; a cleared `labels` is absent from the object entirely.
//     REQ-66 and REQ-107 turn on exactly that distinction, and a format
//     that encoded absence as an empty string would collapse them.
//  2. Values are stored VERBATIM as the canonical strings that crossed the
//     CLI. A set value is already JDR 0001 §D13's canonical JSON array by
//     the time it reaches here, so storing it as a string and returning it
//     unchanged is what makes read-back equality BYTE equality (REQ-71,
//     `0004:C12`). Re-encoding the members here would be a second encoder,
//     which REQ-71 forbids.
//
// # Verdicts and failures are declared, not ambient
//
// A gate's verdict and an artifact's reachability come from the accessor's
// DECLARED `path` in the model, never from the process environment: RDR
// 0004's `0004:C3` forbids discovering authoritative artifacts ambiently,
// and REQ-31 restates it for this CLI. `path` is RDR 0002's carried
// per-accessor locator, so a model author names the behaviour the same way
// they name the role — which is what lets a fixture model declare a gate
// that denies, or a read-back that cannot complete, without the CLI
// consulting anything outside the model and the caller's `--artifact`
// bindings.
package flowbind

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/resolve"
)

// --- declared-path behaviour vocabulary ----------------------------------

const (
	// gatePrefix marks an accessor path as a gate locator. The segment
	// after it names the verdict the gate answers.
	gatePrefix = "flow.gate"
	// unreachableSuffix marks a locator the binding cannot reach: a gate
	// that cannot be consulted, or a write whose read-back cannot
	// complete. Both are exit-3 environment classes at the CLI (REQ-47,
	// REQ-55, REQ-104), never a deny and never a mismatch.
	unreachableSuffix = "unreachable"
)

// verdictFor reads the verdict a gate accessor's declared path names. A
// path with no verdict segment answers allow: an author who declares a
// gate without naming a verdict has declared a gate that permits.
func verdictFor(path string) (accessor.Verdict, bool) {
	suffix := strings.TrimPrefix(path, gatePrefix)
	switch strings.TrimPrefix(suffix, ".") {
	case string(accessor.VerdictDeny):
		return accessor.VerdictDeny, true
	case string(accessor.VerdictIndeterminate):
		return accessor.VerdictIndeterminate, true
	case unreachableSuffix:
		// Not a verdict at all: the gate could not be consulted.
		return "", false
	default:
		return accessor.VerdictAllow, true
	}
}

// unreachable reports whether a declared path names a locator this binding
// cannot reach.
func unreachable(path string) bool {
	return strings.HasSuffix(path, "."+unreachableSuffix) ||
		strings.HasSuffix(path, "-"+unreachableSuffix)
}

// --- the artifact store --------------------------------------------------

// store is the flat tag map one artifact carries. A key's PRESENCE is the
// artifact's own answer to "does this tag exist" — see the package comment.
type store map[string]string

// load reads the artifact at path. A file that does not exist yet is an
// EMPTY artifact, not an error: `flow set-state` establishes state on a
// path the caller names, and requiring the caller to pre-create it would
// make the first write of any flow impossible.
func load(path string) (store, error) {
	buf, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store{}, nil
	}
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(buf))) == 0 {
		return store{}, nil
	}
	var s store
	if err := json.Unmarshal(buf, &s); err != nil {
		return nil, err
	}
	if s == nil {
		s = store{}
	}
	return s, nil
}

// save writes the artifact atomically: the replacement is staged beside the
// target and renamed over it, so a reader never observes a half-written
// artifact. RDR 0004 makes read-back the commit-time check, and a torn
// write would make that check judge bytes no writer intended.
func save(path string, s store) error {
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return err
		}
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".flow-artifact-*")
	if err != nil {
		return err
	}
	staged := tmp.Name()
	defer func() { _ = os.Remove(staged) }()

	// The same non-HTML-escaping encoder every CLI wire site uses, so a set
	// member carrying `<` or `&` is stored as itself and read back
	// byte-identically (REQ-71).
	if err := clierr.WriteJSONLine(tmp, s); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(staged, 0o600); err != nil {
		return err
	}
	return os.Rename(staged, path)
}

// --- read ----------------------------------------------------------------

// Reader is the read binding over one caller-supplied artifact.
type Reader struct {
	// Path is the accessor's DECLARED locator, which may mark the artifact
	// unreachable. It is not the caller's artifact path — that arrives per
	// invocation on accessor.Artifact.
	Path string
}

// Capability reports read.
func (r Reader) Capability() accessor.Capability { return accessor.CapRead }

// Read resolves each requested key against the caller-supplied artifact.
//
// Every requested key is answered: one that the artifact carries comes back
// with its value, and one it does not comes back ESTABLISHED-ABSENT. That
// distinction is the whole point of `0004:C6` — the binding must not leave
// a key unclassified, because the executor's default for an unclassified
// key is `unreadable`, and guessing absence is the failure RDR 0004 exists
// to prevent.
func (r Reader) Read(_ context.Context, art accessor.Artifact, requested []string) (
	[]accessor.KeyValue, []string, error,
) {
	if unreachable(r.Path) {
		return nil, nil, errors.New("read accessor locator is unreachable: " + r.Path)
	}
	s, err := load(art.Path)
	if err != nil {
		return nil, nil, err
	}

	if _, sealed := s[sealedKey]; sealed {
		// The last write to this artifact declared an unreachable read-back
		// locator. Every requested key is reported UNREADABLE rather than
		// absent: the verification did not run, which is `incomplete_read`
		// (`0004:C7`), and reporting absence instead would be the guess RDR
		// 0004 exists to prevent — it would read as "the artifact does not
		// carry this", a positive claim nothing established.
		return nil, slices.Clone(requested), nil
	}

	values := make([]accessor.KeyValue, 0, len(requested))
	for _, key := range requested {
		v, held := s[key]
		values = append(values, accessor.KeyValue{Key: key, Value: v, Absent: !held})
	}
	return values, nil, nil
}

// --- write ---------------------------------------------------------------

// Writer is the write binding over one caller-supplied artifact.
type Writer struct {
	// Path is the accessor's DECLARED locator. A `-unreachable` suffix
	// makes the post-mutation read-back unable to complete, which is the
	// exit-3 "may have been applied" class (REQ-104), NOT a mismatch.
	Path string

	invocations int
}

// Capability reports write.
func (w *Writer) Capability() accessor.Capability { return accessor.CapWrite }

// Invocations reports how many times Apply ran, so a test can assert the
// accessor layer performed no retry or re-derivation (`0004:C14`).
func (w *Writer) Invocations() int { return w.invocations }

// Apply performs the planned mutation. A `<clear>` planned value is a
// REMOVAL — the key leaves the artifact — never an assignment of the
// literal (`0004:C11`). That is what makes a cleared key read back absent
// while an empty set reads back `[]`.
//
// A writer whose declared locator is unreachable APPLIES the mutation and
// then leaves the artifact unreadable, so the post-mutation read-back
// cannot complete. That ordering is the point: the command already ran, so
// the disposition must be "may have been applied and was not verified"
// (`0004:C14`, REQ-104) — never "the write did not occur", which is what
// returning an error here would claim, and never a read-back MISMATCH,
// which would assert the artifact is positively wrong when in truth nothing
// managed to look at it.
func (w *Writer) Apply(_ context.Context, art accessor.Artifact, planned []resolve.Tag) error {
	w.invocations++

	s, err := load(art.Path)
	if err != nil {
		return err
	}
	for _, t := range planned {
		if accessor.IsClear(t.Value) {
			delete(s, t.Key)
			continue
		}
		s[t.Key] = t.Value
	}
	if unreachable(w.Path) {
		// The mutation lands — `s` already carries it — AND the artifact is
		// marked unverifiable, in ONE write. The two must be atomic: the
		// applied-but-unverified sense `0004:C14` requires is a claim about
		// an artifact that really was mutated, so a seal that discarded the
		// mutation would make the refusal's own detail false.
		s[sealedKey] = sealedMarker
	} else {
		// This write's read-back CAN complete, so the artifact is no longer
		// unverifiable and the seal must go — in the same atomic save, for
		// the same reason. The seal names the LAST write's locator, not any
		// past one: a monotonic seal would leave an artifact whose author
		// corrected a mistyped locator permanently unreadable, with no
		// repair path through the CLI. REQ-104's disposition is about the
		// invocation that refused, and REQ-110 fixes identity over CURRENT
		// artifact contents, so nothing makes that refusal sticky.
		delete(s, sealedKey)
	}
	return save(art.Path, s)
}

// sealedKey marks an artifact whose last write declared an unreachable
// read-back locator. The verifying re-read refuses on it, so the write is
// "applied but not verified" (`0004:C13`, `0004:C14`, REQ-104) rather than
// a mismatch asserting the artifact is wrong.
//
// It lives IN the artifact rather than in binding state because RDR 0004
// selects the read-back reader by artifact ROLE, not by the writer's
// locator: the artifact is the only channel the writer and its verifying
// reader share. Keeping it in memory would also make the refusal depend on
// process history rather than on artifact contents, and REQ-110 fixes
// request identity over "the same model revision and the same artifact
// contents" — contents that a caller can inspect, and this key is part of
// them.
const sealedKey = "\x00flow.readback-unreachable"

const sealedMarker = "1"

// --- gate ----------------------------------------------------------------

// Gate is the gate binding. Its verdict comes from the accessor's declared
// path; see the package comment on why it is declared rather than ambient.
type Gate struct {
	Path string
}

// Capability reports gate.
func (g Gate) Capability() accessor.Capability { return accessor.CapGate }

// Gate answers allow, deny, or indeterminate. An unreachable locator is an
// EXECUTION FAILURE, not a verdict: a gate that could not be consulted must
// reach the CLI as exit 3 rather than as a deny (REQ-47, REQ-55), and
// returning a verdict here would launder it into a decision the model never
// made.
func (g Gate) Gate(_ context.Context, _ accessor.Artifact) (accessor.Verdict, string, error) {
	verdict, ok := verdictFor(g.Path)
	if !ok {
		return "", "", errors.New("gate accessor could not be consulted: " + g.Path)
	}
	switch verdict {
	case accessor.VerdictDeny:
		return verdict, "the gate " + g.Path + " denied the transition", nil
	case accessor.VerdictIndeterminate:
		return verdict, "the gate " + g.Path + " could not decide", nil
	default:
		return accessor.VerdictAllow, "", nil
	}
}

// --- cardinality (RDR 0019 A6) -------------------------------------------

// StoreKeys reports which keys the artifact at path carries, sorted.
//
// It is RDR 0019 A6's surviving carrier for the emptiness read: `init-state`
// seeds if and only if EVERY bound artifact carries NO key, and the shipped
// accessor seam is strictly KEY-SCOPED — `ReadBinding.Read` answers only the
// keys it is handed and reports no store cardinality — so no verb sitting on
// that seam as it ships can evaluate the predicate.
//
// It reads the store directly and DELIBERATELY does not enter `Reader.Read`'s
// sealed short-circuit. A read-back-SEALED artifact carries `::sealedKey` and
// nothing else once its owned keys are cleared: that is a ONE-key, NON-EMPTY
// store, and `0019:C1` fixes it as a no-op SUCCESS at exit 0 rather than the
// exit-3 refusal a key-scoped read returns. A carrier that inherited the
// unreadable short-circuit would degrade that arm from no-op success to a
// refusal. What the predicate needs is a COUNT, obtained without reading key
// VALUES, and what the no-op arm's report needs is which `[initial]` keys the
// store does not carry — both are answered by the key set alone, which is why
// one probe serves both and no key value crosses this seam.
//
// The seal is REPORTED, not filtered: `0019:C1` fixes the count as over STORE
// keys and not owned keys, and names the sealed store as the arm where the two
// differ. A probe that hid the seal would make that store read empty and seed
// beneath an unverifiable state.
//
// An absent file is an EMPTY artifact and carries no key, matching `::load` —
// a never-written store and one emptied by clears are the same store, since a
// clear is a key REMOVAL and not a tombstone (`0004:C11`).
func StoreKeys(path string) ([]string, error) {
	s, err := load(path)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(s))
	for key := range s {
		out = append(out, key)
	}
	slices.Sort(out)
	return out, nil
}
