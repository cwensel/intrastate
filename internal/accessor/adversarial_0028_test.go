package accessor_test

// RDR 0028 Phase 3b — ADVERSARIAL failure-mode probes for the declared
// line-edit writer.
//
// Written independently of the Phase 1 suite, anchored in the record's
// `Trade-offs / Failure Modes` section. That section splits its arms into
// "Visible" and "Silent risk", and every probe below attacks one of the
// two seams where a Visible arm quietly becomes something else:
//
//   - a refusal the record files under "Visible … the artifact is
//     untouched" that reaches the caller in the WRONG EXIT GROUP, so an
//     automated caller reads "repair the environment and re-run
//     unchanged" over an invocation that can never succeed unchanged
//     (ADV-1, ADV-2);
//   - a mutation the record's `re-anchor:` invariant exists to prevent —
//     "This is what stops a replacement from de-anchoring itself or
//     POISONING A SIBLING RULE'S ANCHOR ON THE NEXT RUN" — that lands
//     green because the sibling was not in THIS invocation's plan
//     (ADV-3).
//
// All three are shapes a green Phase 1 suite does not catch: each test
// fixture LOADS, LINTS and produces a success- or refusal-shaped result of
// the right CLASS. What is wrong is the discriminator riding beside it, or
// the buffer the refusal never inspected.
//
// Helper and fixture identifiers here all carry an `adv0028` prefix so
// this file shares no package-level name with the Phase 1 suite.

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/flowbind"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// --- ADV-1 ---------------------------------------------------------------

// ADV-1 — the entry-level UNBOUND-`{tag.<key>}` precondition, which
// `0028:C1.3` `order:` names in the same breath as every other pre-write
// refusal, does NOT carry the typed `Err` that puts it in the exit-2
// group. It reaches the caller as an ordinary environment failure.
//
// FAILURE MODE (`0028:FM`, third Visible arm): "the entry-level
// preconditions refuse the same class before mutation, with no rule to
// name: an unbound `{tag.<key>}` and a `-`-prefixed bound value (Detail
// names the placeholder, C1.6) … the artifact is untouched."
//
// `0028:C1.3` `order:` enumerates exactly which refusals it is speaking
// about and then fixes their exit group in one sentence:
//
//	"…and the entry-level preconditions — the read-back gate below,
//	C1.6's unbound tag and C1.6's `-`-prefixed value — name the gate or
//	the placeholder instead, having no rule to name; … EXIT GROUP: THESE
//	REFUSALS are about the REQUEST, not the environment, so they take the
//	exit-2 group, not `execution_failure`'s default exit 3 — the
//	executor-facing typed `Err` discriminates at
//	`flow_exec.go::accessorFailureOf`…"
//
// "these refusals" is bound by the sentence immediately preceding it, and
// that sentence names the entry-level preconditions FIRST. The unbound tag
// is the very first item of `precedence:` step (1).
//
// `flowbind/edit.go::prepare` mints this one as a bare
// `&accessor.ExecError{Detail: …}` with no `Err` field, so
// `refusalOf`'s `errors.Is(err, ErrDeclaredRequest)` is false and
// `accessorFailureOf` routes it to `codeAccessorFailed` — exit 3.
//
// Why a green suite misses it: the refusal's CLASS is
// `execution_failure`, its Detail names the placeholder, and the artifact
// IS untouched. Every assertion the Visible arm invites passes. The only
// observable that differs is the discriminator, and a suite that checks
// class and Detail never reads it.
//
// Why it matters in production: exit 3 is the CLI's promise that the
// REQUEST was fine and the environment was not — "repair the environment
// and re-run the same request unchanged". A pipeline that retries on 3
// spins forever on an invocation that is missing a `--tag`, which is the
// one thing re-running unchanged can never supply. Exit 2 tells the caller
// to change the request, which is the truth here.
//
// Defends: `0028:C1.3` order:/EXIT GROUP (REQ-41, REQ-43), `0028:C1.6`
// binding: (REQ-88), `0028:FM` third Visible arm.
func TestAdv0028_1_UnboundAnchorTagTakesTheRequestExitGroup(t *testing.T) {
	path := adv0028Fixture(t, "- **Status**: Draft\n")

	// The anchor pins the record identity through `{tag.nnnn}`, which the
	// invocation does not bind. This is the consumer's own MVV shape: the
	// README/record anchors carry the record number as a bound tag.
	acc := adv0028Accessor(map[string]table.EditRule{
		"status": {
			Anchor:  `^- \*\*Status\*\*: \[{tag.nnnn}\].*$`,
			Replace: "- **Status**: {status}",
		},
	}, "status")

	writer := flowbind.NewEditWriter(acc, adv0028Entry)
	err := writer.Apply(
		t.Context(),
		// Context is EMPTY: `nnnn` is unbound on this invocation.
		accessor.Artifact{Role: adv0028Role, Path: path},
		[]resolve.Tag{{Key: "status", Value: "Final"}},
	)

	if err == nil {
		t.Fatalf("Apply accepted an anchor naming the unbound placeholder " +
			"`{tag.nnnn}`; `0028:C1.6` binding: requires a refusal before " +
			"mutation, since a placeholder is never passed through literally")
	}
	// The Visible arm's own assertions — these are what a green suite
	// checks, and they must keep passing: the finding below is that they
	// are not SUFFICIENT.
	if got := string(adv0028Read(t, path)); got != "- **Status**: Draft\n" {
		t.Fatalf("the artifact was mutated by a refused edit: %q", got)
	}

	// The finding. The typed `Err` is the ONLY carrier `0028:C1.3` gives
	// `accessorFailureOf` for the exit group, and this refusal omits it.
	if !errors.Is(err, accessor.ErrDeclaredRequest) {
		t.Fatalf("the unbound-`{tag.nnnn}` precondition does not wrap "+
			"`accessor.ErrDeclaredRequest`, so `refusalOf` records "+
			"DeclaredRequest() = false and `accessorFailureOf` routes it to "+
			"the environment exit group (exit 3) instead of the exit-2 group "+
			"`0028:C1.3` EXIT GROUP: fixes for it.\n"+
			"An unbound `--tag` is the definitional REQUEST defect: exit 3 "+
			"tells a caller to repair the environment and re-run the same "+
			"request unchanged, which can never bind the tag.\n"+
			"got err = %v", err)
	}
}

// --- ADV-2 ---------------------------------------------------------------

// ADV-2 — the gate-off read-back pre-check, the third entry-level
// precondition and the one `0028:C1.3` calls the AUTHORITATIVE DETECTOR,
// lands in the environment exit group for the same reason.
//
// FAILURE MODE (`0028:FM`, third Visible arm): "…and a gate-off command
// read-back (Detail names the gate, C1.3) — the artifact is untouched.
// They are decided FIRST, ahead of every rule-scoped refusal above, since
// each condemns the whole entry."
//
// `0028:C1.3` `order:` lists "the read-back gate below" as the FIRST of
// the three entry-level preconditions its EXIT GROUP: sentence governs.
// `executor.go`'s pre-check builds its refusal with
// `refusalOf(def, timeout, ClassExecutionFailure, nil)` — a literal `nil`
// error — so `declaredRequest` is false by construction and no `Err` can
// ever reach it.
//
// This is a DISTINCT site from ADV-1: ADV-1's refusal is minted in the
// binding and could carry `Err`; this one is minted in the executor, which
// has no binding error to wrap, so the fix is a different edit. A suite
// that fixed one would leave the other.
//
// Why a green suite misses it: identical to ADV-1. The class is
// `execution_failure`, the Detail names the gate and the
// `--allow-commands` opt-in, and the artifact is byte-identical. MVV step
// 5 ("the same pipeline without `--allow-commands` … refuses before
// mutation and both files are byte-identical to before") passes.
//
// Why it matters in production: a forgotten `--allow-commands` is the
// archetypal request defect — the record says so itself ("a forgotten flag
// must not produce `read_back_incomplete`"). Having fenced out the
// misleading APPLIED sense, the refusal then hands the caller the
// misleading RETRY sense instead.
//
// Defends: `0028:C1.3` read-back:/order:/EXIT GROUP (REQ-41, REQ-43,
// REQ-59), `0028:FM` third Visible arm.
func TestAdv0028_2_GateOffReadBackPreCheckTakesTheRequestExitGroup(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "record.md")
	adv0028Write(t, path, "- **Status**: Draft\n")

	model := &table.Model{
		ID:   "adv0028",
		Tags: map[string]table.TagDecl{"status": {Provenance: table.ProvenanceOwned}},
		Readers: map[string]table.Accessor{
			// A COMMAND-backed reader for the same role: the pre-check's
			// left-hand predicate.
			"adv0028-reader": {
				Role:     adv0028Role,
				Keys:     []string{"status"},
				Command:  []string{"/bin/echo", "{artifact}"},
				Timeout:  "1s",
				ReadBack: true,
			},
		},
		Writers: map[string]table.Accessor{
			adv0028Entry: {
				Role:     adv0028Role,
				Keys:     []string{"status"},
				Timeout:  "1s",
				ReadBack: true,
				Edit: map[string]table.EditRule{
					"status": {
						Anchor:  `^- \*\*Status\*\*: .+$`,
						Replace: "- **Status**: {status}",
					},
				},
			},
		},
	}

	// The gate is OFF — `allowCommands` false at the single production
	// construction site.
	reg := flowbind.Registry(model, dir, false)
	exec := accessor.NewExecutor(reg, accessor.Artifacts{
		adv0028Role: {Role: adv0028Role, Path: path},
	})

	got := exec.Write(t.Context(), adv0028Entry, resolve.Plan{
		Writes: []resolve.Tag{{Key: "status", Value: "Final"}},
	})

	if got.Refusal == nil {
		t.Fatalf("Write ran with the gate OFF and a command-backed read-back "+
			"reader; `0028:C1.3` read-back: requires a refusal BEFORE "+
			"mutation. Written = %+v", got.Written)
	}
	// The Visible arm's assertions, which must keep passing.
	if got.Refusal.Class != accessor.ClassExecutionFailure {
		t.Fatalf("Refusal.Class = %q; want execution_failure — `0028:C1.3` "+
			"order: forbids a new refusal class", got.Refusal.Class)
	}
	if got.Refusal.Applied() {
		t.Fatalf("the gate pre-check carries the applied-but-unverified " +
			"sense; `0028:C1.3` order: forbids it for a pre-write refusal")
	}
	if body := string(adv0028Read(t, path)); body != "- **Status**: Draft\n" {
		t.Fatalf("the artifact was mutated by the gate-off pre-check: %q", body)
	}

	// The finding.
	if !got.Refusal.DeclaredRequest() {
		t.Fatalf("the gate-off read-back pre-check reports "+
			"DeclaredRequest() = false, so `accessorFailureOf` routes it to "+
			"the environment exit group (exit 3) rather than the exit-2 "+
			"group `0028:C1.3` EXIT GROUP: fixes for the entry-level "+
			"preconditions — of which `order:` names the read-back gate "+
			"FIRST.\n"+
			"`executor.go` mints this one as "+
			"`refusalOf(def, timeout, ClassExecutionFailure, nil)`: the "+
			"literal nil means no `Err` can ever set the discriminator, so "+
			"this is a second site from ADV-1's and needs its own fix.\n"+
			"A forgotten `--allow-commands` is a REQUEST defect; exit 3 "+
			"invites the caller to re-run it unchanged forever.\n"+
			"detail = %q", got.Refusal.Detail)
	}
}

// --- ADV-3 ---------------------------------------------------------------

// ADV-3 — the re-anchor pass skips every rule of the entry the CURRENT
// invocation did not plan, so one rule's replacement can silently poison
// an unplanned sibling's anchor. The write lands green; the NEXT
// invocation refuses on an artifact this one corrupted.
//
// FAILURE MODE (`0028:FM`, "Silent risk" arms): the record files exactly
// two silent risks, both of the shape "the wrong line is rewritten and
// read-back is the only defence". This is a THIRD, and read-back cannot
// see it: the poisoned key is not in the plan, so the read-back oracle
// never compares it. The corruption is invisible at BOTH ends.
//
// `0028:C1.3` `re-anchor:` states the invariant with no plan qualifier and
// then names this exact consequence as the reason it exists:
//
//	"after the buffer is rewritten in memory, EVERY RULE'S anchor is run
//	again over the POST-EDIT buffer and must select exactly its own
//	rewritten line (or, for a deleted line, zero lines); otherwise refuse
//	`edit_anchor_unstable` before any write. … This is what stops a
//	replacement from de-anchoring itself or POISONING A SIBLING RULE'S
//	ANCHOR ON THE NEXT RUN (premortem P-4, P-13)."
//
// "on the next run" is decisive: a sibling whose anchor is poisoned only
// matters on a LATER invocation, which is precisely the invocation that
// did not plan it now. Restricting the pass to the planned rules removes
// the clause's stated purpose while keeping its letter for the trivial
// case.
//
// `flowbind/edit.go::prepare` `continue`s past any key with no planned
// value, so such a rule never enters `plans` and `reAnchor` never sees it.
//
// The fixture is the record's own `select:` grammar, not a contrivance:
// `0028:C1.1` states "a line is SELECTED when the pattern matches ANYWHERE
// in it (terminator excluded); authors pin `^…$`" — pinning is ADVICE, and
// an unpinned anchor is admitted by contract. The `owner` rule here is
// unpinned; the `status` rule's planned value happens to contain the text
// that anchor matches.
//
// Why a green suite misses it: the entry LINTS, the planned rule's own
// re-anchor passes, the write applies, read-back verifies `status` — the
// only key in the plan — and the result is SUCCESS. Nothing in the run
// touches `owner`.
//
// Why it matters in production: the artifact is now in a state no
// invocation can repair through this carrier. The next plan naming `owner`
// refuses `edit_anchor_ambiguous`, pointing the author at the `owner`
// rule, which is correct; the line that broke it was written by the
// `status` rule an unknown number of runs earlier, and `0028:FM`'s
// Recovery arm ("fix the anchor or the artifact") cannot be followed
// without that history.
//
// Defends: `0028:C1.3` re-anchor: (REQ-38, REQ-39), `0028:C1.3` write:
// ("ALL RULES of one entry rewrite ONE buffer"), `0028:FM` Silent risk.
func TestAdv0028_3_ReAnchorCoversRulesThisPlanDidNotName(t *testing.T) {
	const before = "- **Status**: Draft\n- **Owner**: alice\n"
	path := adv0028Fixture(t, before)

	acc := adv0028Accessor(map[string]table.EditRule{
		"status": {
			Anchor:  `^- \*\*Status\*\*: .+$`,
			Replace: "- **Status**: {status}",
		},
		// UNPINNED by design and admitted by `0028:C1.1`: it selects the
		// Owner line today because that is the only line carrying `alice`.
		"owner": {
			Anchor:  `alice`,
			Replace: "- **Owner**: {owner}",
		},
	}, "status", "owner")

	writer := flowbind.NewEditWriter(acc, adv0028Entry)

	// The plan names ONLY `status`. `owner` keeps whatever the artifact
	// holds — and its rule is still a rule of this entry.
	err := writer.Apply(
		t.Context(),
		accessor.Artifact{Role: adv0028Role, Path: path},
		[]resolve.Tag{{Key: "status", Value: "Final (handover from alice)"}},
	)

	after := string(adv0028Read(t, path))

	if err != nil {
		// The contract-honouring outcome: `edit_anchor_unstable` before
		// any write, because the `owner` rule's anchor no longer selects
		// exactly its own line on the post-edit buffer.
		if after != before {
			t.Fatalf("the edit refused but the artifact changed: %q", after)
		}
		return
	}

	// The finding. The write landed, and the `owner` rule's anchor now
	// selects TWO lines of the file this entry just produced.
	hits := 0
	for _, line := range adv0028Lines(after) {
		if adv0028Contains(line, "alice") {
			hits++
		}
	}
	if hits > 1 {
		t.Fatalf("the `status` rule's replacement poisoned the UNPLANNED "+
			"sibling rule `%s.edit.owner`: its anchor now selects %d lines "+
			"of the post-edit buffer, and the write applied anyway.\n"+
			"`0028:C1.3` re-anchor: runs EVERY RULE'S anchor over the "+
			"post-edit buffer and refuses `edit_anchor_unstable` before any "+
			"write; `edit.go::prepare` drops a rule whose key this plan did "+
			"not name, so `reAnchor` never sees it.\n"+
			"The clause's own stated purpose is what this defeats: it exists "+
			"to stop a replacement \"poisoning a sibling rule's anchor ON "+
			"THE NEXT RUN\" — and the next run is exactly the invocation "+
			"that plans the sibling.\n"+
			"Read-back cannot catch it either: `owner` is not in the plan, "+
			"so the oracle never compares it. Silent at both ends.\n"+
			"before = %q\nafter  = %q", adv0028Entry, hits, before, after)
	}
}

// --- fixtures ------------------------------------------------------------

const (
	adv0028Role  = "record"
	adv0028Entry = "adv0028-writer"
)

// adv0028Accessor builds a write accessor carrying `edit` and nothing
// else, so the registry's carrier discriminator takes the `edit` arm.
func adv0028Accessor(rules map[string]table.EditRule, keys ...string) table.Accessor {
	return table.Accessor{
		Role:     adv0028Role,
		Keys:     keys,
		Timeout:  "1s",
		ReadBack: true,
		Edit:     rules,
	}
}

func adv0028Fixture(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "record.md")
	adv0028Write(t, path, body)
	return path
}

func adv0028Write(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("writing the fixture %s: %v", path, err)
	}
}

func adv0028Read(t *testing.T, path string) []byte {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the fixture %s: %v", path, err)
	}
	return body
}

// adv0028Lines splits on "\n" for the ASSERTION only — the binding's own
// splitter is what the contract governs and this one deliberately does not
// stand in for it.
func adv0028Lines(body string) []string {
	var out []string
	start := 0
	for i := range len(body) {
		if body[i] == '\n' {
			out = append(out, body[start:i])
			start = i + 1
		}
	}
	if start < len(body) {
		out = append(out, body[start:])
	}
	return out
}

func adv0028Contains(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
