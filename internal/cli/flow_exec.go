package cli

// RDR 0005 — the shared execution path the four verbs sit on: which read
// accessors run, how their refusals become CLI codes, and how a kernel
// refusal is rendered.

import (
	"context"
	"path/filepath"
	"slices"
	"strings"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/cli/flowbind"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
	"github.com/spf13/cobra"
)

// flowRequest is one parsed, validated invocation, assembled before any
// accessor runs.
type flowRequest struct {
	model     *table.Model
	modelRef  string
	artifacts map[string]string
	observed  []resolveTag
	registry  accessor.Registry
}

// revision is the loaded model's own revision identity, carried through
// VERBATIM (REQ-25). The CLI never derives, hashes, or synthesizes it: a
// model that declares none renders `revision` empty rather than a
// CLI-invented value, because a synthesized revision would make two
// distinct models compare equal — or one model compare unequal to itself
// across a whitespace edit — and REQ-110 makes revision part of request
// identity.
//
// RDR 0002's `[model]` block admits `id`, `version`, `description`, and
// `metadata` and no revision field, so NO model can declare one and
// REQ-25's consequent governs every payload this CLI emits: empty (DEV-7).
//
// `KernelTable().Revision` is deliberately NOT forwarded here. RDR 0002
// sets it to the model id for the KERNEL's own use, where it is the
// table-identity the resolver compares — a different consumer with a
// different meaning. Standing it in as the model revision would be exactly
// the CLI-invented value REQ-25 forbids, and it would leave REQ-110's
// request identity unable to tell two revisions of one model apart, which
// is the discrimination the field is named for.
func (flowRequest) revision() string { return "" }

// --- the narrowed read-accessor set (REQ-35, REQ-36) ---------------------

// invokedReaders returns the read accessors that must run for `next` and
// `resolve`: exactly those serving an owned key some candidate row of the
// requested model requires — not every declared reader.
//
// `outcome` narrows the candidate rows for `resolve`; an empty outcome
// means `next`, which takes the union over ALL rows (DEV-1, reading (a)).
// Reading (a) is what keeps `flow-artifact-missing` a total, caller-
// independent refusal identity: under reading (b) the same argv would
// refuse or succeed depending on how much observed context the caller
// happened to pass, and REQ-110 does not list caller context as an input to
// reader selection.
//
// The demand set is `Row.RequiresOwned` UNIONED with the owned keys the
// row's guard atoms reference (DEV-8).
//
// `RequiresOwned` alone is not the demand set: RDR 0002 derives it at
// normalization from the write block and clear list only, and
// `resolve.Row`'s own doc adds that "the guard's key set is not required to
// be a subset of this one". A row that GUARDS on an owned key it never
// writes cannot be decided without that key, and REQ-35 fixes the invoked
// set as those readers "serving an owned key some candidate row REQUIRES" —
// a key a guard must evaluate is required by that row. Reading the demand
// as `RequiresOwned` alone left the serving reader uninvoked and turned a
// decidable transition into `flow-guard-unevaluable`, a refusal naming
// nothing the caller could repair. ASSUMPTION A-4 under-read the clause
// that way and is superseded by DEV-8.
//
// Both halves are dumpable NORMALIZED model data available before any
// evaluation (CA A7): `Row.Atoms` carries each atom's key and its authored
// block verbatim. Narrowing here reads KEYS and never operators or
// literals, so it evaluates nothing, invents no fact, and does not
// reimplement RDR 0003's typed operator semantics, which stay the kernel's
// seam (REQ-113).
func invokedReaders(m *table.Model, outcome string) []string {
	demanded := map[string]bool{}
	for _, row := range m.Rows {
		if outcome != "" && row.Outcome != outcome {
			// A row binding a different outcome cannot serve this request,
			// escape row included: the kernel skips any escape row whose
			// `Outcome != in.Recognized` before it can rescue anything
			// (`internal/resolve.escapeOrRefuse`), so a differing-outcome
			// escape row is unrescuable by construction. Retaining one used
			// to be harmless because escape rows carry no `RequiresOwned`;
			// once guard-owned keys joined the demand set (DEV-8) it stopped
			// being harmless, because an unbound or failing reader serving
			// only that unrescuable row now refuses the whole request at
			// exit 3 — masking a valid plan, or the `flow-unmodeled-outcome`
			// the caller should have seen. An escape row MATCHING the
			// requested outcome keeps its demands.
			continue
		}
		for _, key := range row.RequiresOwned {
			demanded[key] = true
		}
		for _, key := range guardOwnedKeys(m, row) {
			demanded[key] = true
		}
		// RDR 0011 C1: each row's MATCH-block owned keys join the demand
		// set, so the assembled view actually carries the keys `flow next`'s
		// candidate predicate reads. The term binds BOTH callers and is NOT
		// scoped to `next` (`0011:BR6`): a `next`-only term would have the
		// two verbs assemble different views over one model.
		//
		// It IS scoped to the rows the CALLER can consult, which for the
		// match term excludes an ESCAPE row under `next` (DEV-9). C1 fixes
		// `next`'s predicate over "each NON-ESCAPE row", and justifies the
		// term by "the assembled view MUST actually carry the keys THE
		// PREDICATE READS" — "a row cannot be match-decided without the
		// key". `runFlowNext`'s row loop skips every escape row
		// (`if len(row.Escape) != 0 { continue }`) before `probeRow`, so an
		// escape row's match atoms never reach a probe and `next` never
		// match-decides one; the RDR's own joint-check states it flatly:
		// "`runFlowNext` never lists escape rows". Demanding a reader for
		// such a key made `next` refuse `flow-artifact-missing` (exit 2)
		// above the row loop for a fact no reported row could consume — a
		// refusal no Failure Mode contemplates and one C1 leaves the caller
		// no remedy for, since the key is owned (`--tag` is refused
		// `flow-tag-owned`) and `--all` inherits the same pre-loop refusal,
		// denying F1 its two-run diagnostic.
		//
		// `resolve` keeps escape rows: `internal/resolve.escapeOrRefuse`
		// evaluates `view.matches(row.Match)` on every escape row binding
		// the requested outcome, and that phase's reachability is not
		// knowable before the kernel runs. Dropping them would let an
		// absent key silently fail a rescue through the kernel's two-valued
		// `matches` — reinstating the hidden-fact defect this term exists
		// to remove. The resulting exit 2/3 on a run that never reaches the
		// rescue phase is the class C1 NAMES AND ACCEPTS verbatim: "such a
		// run turns from a plan into exit 2/3, before the kernel and above
		// the escape phase".
		//
		// Mode-independence (C1) is untouched: `--all` and the default
		// consult the same row set — only the predicate over it differs.
		if outcome != "" || len(row.Escape) == 0 {
			for _, key := range matchOwnedKeys(m, row) {
				demanded[key] = true
			}
		}
	}

	var out []string
	for name, acc := range m.Readers {
		if slices.ContainsFunc(acc.Keys, func(k string) bool { return demanded[k] }) {
			out = append(out, name)
		}
	}
	slices.Sort(out)
	return out
}

// guardOwnedKeys returns the OWNED tag keys one row's guard atoms
// reference.
//
// Only `all` and `unless` are guard blocks: RDR 0007's row-verdict formula
// (`0007:C5`/`0007:C6`) names those two and nothing else, and a `match`
// atom is the kernel's selection pattern rather than a guard. Provenance
// comes from the model's own `[tags.<key>]` declaration, so an observed or
// recognized key a guard reads demands no reader — the caller supplies it
// through `--tag`, and no reader serves it.
func guardOwnedKeys(m *table.Model, row table.Row) []string {
	return ownedAtomKeys(m, row, table.BlockAll, table.BlockUnless)
}

// matchOwnedKeys returns the OWNED tag keys one row's MATCH atoms
// reference — RDR 0011 C1's demand-set term.
//
// `flow next` decides a match atom's PRESENCE against the assembled view,
// so a key no invoked reader establishes has its atom omitted from every
// probe and the row comes back a candidate carrying `{key, absent}`. Over a
// model that matches on an owned key it never writes, clears, or guards
// that degrades the default to `--all`. Adding the term is what makes the
// predicate real, and it invokes no reader the model did not already
// declare.
func matchOwnedKeys(m *table.Model, row table.Row) []string {
	return ownedAtomKeys(m, row, table.BlockMatch)
}

// ownedAtomKeys returns the OWNED tag keys the row's atoms reference from
// any of the named blocks. Provenance comes from the model's own
// `[tags.<key>]` declaration, so an observed or recognized key demands no
// reader — the caller supplies it through `--tag`, and no reader serves it.
func ownedAtomKeys(m *table.Model, row table.Row, blocks ...table.Block) []string {
	var out []string
	for _, atom := range row.Atoms {
		if !slices.Contains(blocks, atom.Block) {
			continue
		}
		if m.Tags[atom.Key].Provenance != table.ProvenanceOwned {
			continue
		}
		if !slices.Contains(out, atom.Key) {
			out = append(out, atom.Key)
		}
	}
	return out
}

// declaredReaders returns every declared reader, sorted. `flow read-state`
// is the single exception to narrowing: a diagnostic read has no candidate
// set to narrow by (REQ-37).
func declaredReaders(m *table.Model) []string {
	var out []string
	for name := range m.Readers {
		out = append(out, name)
	}
	slices.Sort(out)
	return out
}

// --- reading owned state -------------------------------------------------

// readerOutput is one reader's contribution: the keys it declares and the
// tags it returned.
type readerOutput struct {
	ID   string            `json:"id"`
	Keys []string          `json:"keys"`
	Tags map[string]string `json:"tags"`
}

// runReaders invokes the named readers and assembles owned state from them.
//
// `flow-artifact-missing` is raised HERE and scoped to the roles the
// INVOKED set needs (REQ-33, REQ-36) — checking every declared reader's
// role instead would make an unneeded reader's unbound role refuse an
// invocation that never touches it.
func (r flowRequest) runReaders(ctx context.Context, names []string) (
	[]readerOutput, []resolve.Tag, *clierr.CLIError,
) {
	// Every needed role must be bound before ANY accessor runs, so the
	// refusal names the role rather than surfacing as an execution failure
	// deep inside a binding.
	for _, name := range names {
		def, ok := r.registry.Lookup(name, accessor.CapRead)
		if !ok {
			continue
		}
		if _, bound := r.artifacts[def.Accessor.Role]; !bound {
			return nil, nil, userErr(codeArtifactMissing, def.Accessor.Role,
				"the artifact role `"+def.Accessor.Role+"` that the read "+
					"accessor `"+name+"` needs has no --artifact binding")
		}
	}

	exec := accessor.NewExecutor(r.registry, r.artifactMap())

	outputs := make([]readerOutput, 0, len(names))
	var owned []resolve.Tag
	for _, name := range names {
		def, ok := r.registry.Lookup(name, accessor.CapRead)
		if !ok {
			continue
		}
		// RDR 0023 `0023:C1` — the reader-execution recording seam, wrapped
		// around the per-reader EXECUTION rather than around the name
		// computation.
		//
		// The position is normative even though the form is not: C1 fixes
		// the measurand as "OBSERVED READER EXECUTION, not a recomputed
		// derivation" and rejects `invokedReaders(model, outcome)` by name,
		// because neither of that function's arguments depends on anything
		// a report-width flag could touch — comparing it across two runs
		// returns equal sets on every implementation, including one that
		// never runs a reader at all. Recording HERE records what ran.
		//
		// The seam is FLAG-BLIND, necessarily: the projection flag is
		// readable at exactly one lexical site and this is not it. Both
		// report widths drive this hook identically, which is what makes
		// the two recordings comparable in the first place.
		recordReaderExecution(name)
		result := exec.Read(ctx, name)
		if result.Refused() {
			return nil, nil, accessorFailure(*result.Refusal, phaseRead)
		}

		tags := map[string]string{}
		for _, v := range result.Values {
			// An established-absent key is OMITTED from the reported tags
			// while STAYING in `keys`. That is what makes "absent from the
			// artifact" and "not requested" distinguishable from the
			// payload alone (REQ-77) — and it is the same omission
			// `0004:C8` requires at the resolver seam, where a placeholder
			// would silently retire an `owned_state_unavailable` refusal.
			if v.Absent {
				continue
			}
			tags[v.Key] = v.Value
		}
		outputs = append(outputs, readerOutput{
			ID: name, Keys: def.RequestedKeys(), Tags: tags,
		})
		owned = append(owned, result.OwnedSnapshot()...)
	}

	slices.SortFunc(owned, func(a, b resolve.Tag) int {
		return strings.Compare(a.Key, b.Key)
	})
	return outputs, owned, nil
}

// readerExecutionHook, when non-nil, is called with each reader's id at the
// moment that reader ACTUALLY RUNS — immediately before its `exec.Read`
// invocation inside `runReaders`, the one pass that yields `readers` and
// `owned` together.
//
// It is the seam RDR 0023 `0023:C1` requires the reader differential to be
// written against: "the oracle MUST instead observe which readers ACTUALLY
// RAN in each run — recording execution at the reader invocation site (a
// counting or recording seam around the reader pass, the run's own side
// effects, never a re-call of the planning function)".
//
// Nil in production, so the shipped path pays one nil check and records
// nothing; an observer installs it to collect a run's executed set. It
// takes the reader id and nothing else — in particular it takes no flag and
// no request, because a seam that could observe the projection flag would
// make the very execution it records flag-dependent, and that flag is
// readable at exactly one site, which is not this one.
var readerExecutionHook func(readerID string)

// recordReaderExecution reports that the named reader is about to run.
func recordReaderExecution(readerID string) {
	if readerExecutionHook != nil {
		readerExecutionHook(readerID)
	}
}

func (r flowRequest) artifactMap() accessor.Artifacts {
	// RDR 0028 A1 — the invocation's bound tags cross the seam on the
	// artifact, reaching both `Apply` (an `edit` writer's `{tag.<key>}`
	// anchors) and `Read` (a C1.6 command reader's `{tag.<key>}` argv).
	// ONE map serves both, so a shared artifact's reader and the writer
	// that addressed a row inside it agree on which row that was.
	//
	// Every role gets the SAME context, because the bindings are the
	// invocation's, not a role's: narrowing by role would make a tag
	// visible to the writer and invisible to the reader that verifies it.
	var context map[string]string
	if len(r.observed) != 0 {
		context = make(map[string]string, len(r.observed))
		for _, t := range r.observed {
			context[t.Key] = t.Value
		}
	}

	out := make(accessor.Artifacts, len(r.artifacts))
	for role, path := range r.artifacts {
		out[role] = accessor.Artifact{Role: role, Path: path, Context: context}
	}
	return out
}

// --- accessor refusal mapping (A-9) --------------------------------------

// phase distinguishes the invoking site, because two accessor refusal
// classes take different CLI codes depending on it: an incomplete READ is
// `flow-read-incomplete`, while an incomplete post-mutation READ-BACK is
// `flow-write-readback-incomplete`. Only the invoking phase separates them,
// and the code table names both at exit 3 with distinct descriptions.
type phase int

const (
	phaseRead phase = iota
	phaseGate
	phaseWrite
)

// accessorSubject names the accessor the way the invoking phase ran it,
// matching the spelling each phase's other refusals already use: the read
// arm says "the read accessor `x`" as `codeReadIncomplete` does, the gate
// arm says "the gate `x`" as `codeGateIndeterminate` does, and the write
// arm says "the write accessor `x`" as the read-back refusals do. Naming
// the role the accessor actually ran under is the whole point — the
// alternative is misattributing a failure to a role that never ran.
//
// The phase is the right axis for the SUBJECT, and only for the subject.
// Which ENVELOPE gets built is decided by the refusal's origin, because a
// write-phase refusal is not necessarily a line edit's.
func accessorSubject(at phase, id string) string {
	switch at {
	case phaseGate:
		return "the gate `" + id + "`"
	case phaseWrite:
		return "the write accessor `" + id + "`"
	default:
		return "the read accessor `" + id + "`"
	}
}

// accessorFailure maps one accessor refusal class onto its CLI code (A-9).
//
// The exit-3 population is exactly five classes (REQ-19, REQ-20): accessor
// timeout, execution failure, incomplete read, read-back incomplete, and
// post-mutation timeout. Everything else exits 2. The split matters to a
// caller's retry loop: exit 3 promises "repair the environment and re-run
// the same request unchanged", so a refusal about the REQUEST must never
// take it, or the caller spins forever on an input they must instead fix.
func accessorFailure(refusal accessor.Refusal, at phase) *clierr.CLIError {
	return withStderrTail(accessorFailureOf(refusal, at), refusal.Detail)
}

// withStderrTail lands a command-backed invocation's bounded stderr tail
// in the CLI's ONE `Detail` slot (`0025:C4`, REQ-65).
//
// There is exactly one slot — `clierr.CLIError` has a single `Detail` and
// `EmitText` renders a single `detail:` line — so "renders beneath" is not
// available and no second carrier is added. Where both senses are present
// they occupy that slot IN ORDER: the applied-sense text FIRST, because
// losing "the write may have applied" to a diagnostic tail would drop the
// more consequential fact, and the tail appended after it.
func withStderrTail(ce *clierr.CLIError, tail string) *clierr.CLIError {
	if ce == nil || tail == "" {
		return ce
	}
	if ce.Detail == "" {
		ce.Detail = tail
		return ce
	}
	ce.Detail = ce.Detail + "\n" + tail
	return ce
}

func accessorFailureOf(refusal accessor.Refusal, at phase) *clierr.CLIError {
	id := refusal.Accessor

	switch refusal.Class {
	case accessor.ClassTimeout:
		if at == phaseWrite {
			ce := envErr(codeReadBackTimeout, id,
				"the post-mutation read-back for the write accessor `"+id+
					"` timed out")
			ce.Detail = detailMayHaveApplied
			return ce
		}
		return envErr(codeAccessorTimeout, id,
			"the accessor `"+id+"` timed out")

	case accessor.ClassExecutionFailure:
		// RDR 0028 `0028:C1.3` EXIT GROUP: — a refusal the binding minted
		// about the REQUEST takes exit 2, not `execution_failure`'s
		// default exit 3. Exit 3 promises "repair the environment and
		// re-run the same request unchanged", and a stale anchor, an
		// ambiguous one or a value carrying a newline is none of those:
		// re-running unchanged spins the caller forever on an input they
		// must instead fix.
		//
		// No refusal class is added. The class is still
		// `execution_failure` and the accessor-level Detail is unchanged;
		// what moves is the exit group and the carriage — the rule id and
		// the reason token reach the caller in `findings[]` rather than
		// only in `detail`, because a per-subject record is what
		// `docs/cli-output-contract.md` requires of a failure that names
		// a subject other than a flag.
		//
		// The marker is minted at three different KINDS of site and only
		// ONE of them is a line edit, so the envelope splits on the
		// refusal's ORIGIN — `DeclaredEdit()` — and not on the invoking
		// phase or the writer's carrier, neither of which can stand in for
		// it:
		//
		//   - `flowbind.EditWriter` mints the rule-scoped and entry-level
		//     refusals C1.3 was written about. These have a rule
		//     `<id>.edit.<key>` to name and are the only ones that may be
		//     reported as a line edit.
		//   - `cmdbind.substitute` mints C1.6's argv preconditions, which
		//     run for EVERY command-backed accessor. A reader's unbound
		//     `{tag.<key>}` arrives at `phaseRead`, a gate's at
		//     `phaseGate` — and a command-backed WRITER's at `phaseWrite`,
		//     which is why the phase is not the axis. That writer may
		//     carry no `edit` table at all.
		//   - the executor mints the read-back preconditions. Those are
		//     gated on `Edit != nil`, so the CARRIER is not the axis
		//     either: they arrive from an edit writer and are still about
		//     the READER's argv, with no rule to name.
		//
		// Reporting any of the latter two as a line edit invented a rule
		// subject that does not exist and put an argv complaint on a
		// caller's line-edit count. The exit GROUP is 2 on every arm, so
		// C1.3's exit-group sentence is satisfied throughout and no
		// contract moves.
		if refusal.DeclaredRequest() {
			if !refusal.DeclaredEdit() {
				// No `findings[]`: `editRefusalFindings` codes its entry as
				// a line-edit refusal and carries a Detail the caller reads
				// as a rule id `<id>.edit.<key>`. There is no such rule
				// here. The offending placeholder already reaches the
				// caller in `refusal.Detail`, which `withStderrTail`
				// appends, so the scalar subject is the accessor and rides
				// `param`.
				return userErr(codeRequestRefused, id,
					accessorSubject(at, id)+
						" refused this request before running")
			}
			return &clierr.CLIError{
				Code:  codeWriteEditRefused,
				Group: clierr.GroupUserEnv,
				Message: "the write accessor `" + id +
					"` refused the declared line edit before mutating the artifact",
				Findings: editRefusalFindings(refusal),
			}
		}

		// RDR 0026 (`0026:C1` `refusal:`). Checked after the request split
		// above. The two conditions are DISJOINT, not ordered: every
		// refusal 0028's marker names is decided BEFORE any byte is
		// written (C1.3 `order:`), so `Applied()` is false on all of them,
		// and swapping the two arms changes no outcome (verified by
		// mutation at the 0026/0028 rebase). The sequence is written this
		// way to keep that disjointness legible — a reader meets the
		// pre-mutation case first — not because either arm shadows the
		// other. A future change that lets a request refusal report
		// applied would make the order matter; it does not today.
		//
		// The key here is `Applied()`, never the phase. The neighbouring
		// `ClassTimeout` arm above keys on phase because a WRITE-phase
		// timeout is post-command by construction; an execution failure is
		// not — this arm is reached by every pre-mutation failure shape as
		// well, and labelling one of those "may have been applied" makes an
		// agent skip a retry that was safe (`0026:C1` `refusal:`, JDR 0003
		// §D4 (b)). Only the applied refusal takes the distinct code; the
		// rest keep `flow-accessor-failed`.
		if refusal.Applied() {
			ce := envErr(codeWriteFailedApplied, id,
				"the write accessor `"+id+"` failed after its command already ran")
			ce.Detail = detailMayHaveApplied
			return ce
		}
		return envErr(codeAccessorFailed, id,
			"the accessor `"+id+"` could not be executed")

	case accessor.ClassIncompleteRead:
		return envErr(codeReadIncomplete, id,
			"the read accessor `"+id+"` returned an incomplete key set: "+
				strings.Join(refusal.Keys, ", "))

	case accessor.ClassReadBackIncomplete:
		ce := envErr(codeReadBackIncomplete, id,
			"the post-mutation read-back for the write accessor `"+id+
				"` did not complete")
		ce.Detail = detailMayHaveApplied
		return ce

	case accessor.ClassReadBackMismatch:
		// Exit 2, not 3: the read-back RAN and found the artifact wrong.
		// Re-running the same request unchanged cannot help, so this is
		// emphatically not an environment class (`0004:C13`).
		return &clierr.CLIError{
			Code:  codeReadBackMismatch,
			Group: clierr.GroupUserEnv,
			Message: "the post-mutation read-back disagrees with the plan " +
				"for the write accessor `" + id + "`",
			Findings: readBackFindings(refusal),
		}

	case accessor.ClassGateIndeterminate:
		return &clierr.CLIError{
			Code:    codeGateIndeterminate,
			Group:   clierr.GroupUserEnv,
			Message: "the gate `" + id + "` could not decide",
			Findings: []clierr.Finding{{
				Code:    codeGateIndeterminate,
				Message: gateMessage(id, accessor.VerdictIndeterminate, refusal.Reason),
				Param:   id,
			}},
		}

	case accessor.ClassUnknownAccessor:
		return internalErr(codeAccessorUnknown,
			"the accessor `"+id+"` is not bound")

	case accessor.ClassCapabilityMismatch:
		return internalErr(codeCapabilityMismat,
			"the accessor `"+id+"` is bound under another capability")

	default:
		return internalErr(codeAccessorFailed,
			"the accessor `"+id+"` refused with an unmapped class "+
				string(refusal.Class))
	}
}

// detailMayHaveApplied is the sense `0004:C14` requires of a read-back that
// did not run: the mutation MAY have been applied and was NOT verified —
// never "the write did not occur". Without it a caller cannot tell whether
// to retry or to inspect.
const detailMayHaveApplied = "the write command already ran, so the mutation " +
	"may have been applied and was not verified; inspect the artifact before retrying"

// internalErr builds a `GroupInternal` refusal: exit 2, no carrier. These
// three codes are unreachable past load and the CLI checks (REQ-102), so
// reaching one is an invariant violation rather than a caller mistake —
// which is why it must not exit 3 and invite a retry.
func internalErr(code, message string) *clierr.CLIError {
	return &clierr.CLIError{Code: code, Message: message, Group: clierr.GroupInternal}
}

// editRefusalFindings carries a declared-line-edit refusal's rule id and
// reason token into `findings[]` (RDR 0028 `0028:C1.3` EXIT GROUP:).
//
// The output contract splits the two carriers by FAMILY, not by runtime
// cardinality: a scalar family names its one offending subject in
// `param`, and an aggregate family reports in `findings[]` even when it
// has one entry. This family's subject is a RULE — `<id>.edit.<key>` —
// which is not a flag or an argument, and an entry with several rules can
// name any of them, so it is an aggregate family and takes `findings[]`
// with no top-level `param`.
//
// The finding's Message is the accessor-level Detail VERBATIM, because
// that string is the contract's own: the reason token and the rule id, in
// the order and spelling C1.3 fixes. Re-rendering it here would fork the
// wording a caller matches on.
func editRefusalFindings(refusal accessor.Refusal) []clierr.Finding {
	return []clierr.Finding{{
		Code:    codeWriteEditRefused,
		Message: refusal.Detail,
	}}
}

// readBackFindings names one finding per key the read-back disagreed on
// (REQ-103), so a caller sees every divergence rather than the first.
func readBackFindings(refusal accessor.Refusal) []clierr.Finding {
	observed := map[string]string{}
	for _, t := range refusal.Observed {
		observed[t.Key] = t.Value
	}
	out := make([]clierr.Finding, 0, len(refusal.Expected))
	for _, want := range refusal.Expected {
		got, held := observed[want.Key]
		if held && got == want.Value {
			continue
		}
		out = append(out, clierr.Finding{
			Code:    codeReadBackMismatch,
			Key:     want.Key,
			Message: readBackMismatchMessage(want, got, held),
		})
	}
	if len(out) == 0 {
		// A protected non-owned tag changed: the plan's own keys all hold,
		// so no per-key entry above fired, but the artifact is still wrong.
		out = append(out, clierr.Finding{
			Code: codeReadBackMismatch,
			Message: "a tag the write did not plan to mutate changed during " +
				"the write",
		})
	}
	return out
}

func readBackMismatchMessage(want resolve.Tag, got string, held bool) string {
	if accessor.IsClear(want.Value) {
		return "the key `" + want.Key + "` was planned for removal but reads " +
			"back holding " + got
	}
	if !held {
		return "the key `" + want.Key + "` was planned as " + want.Value +
			" but reads back absent"
	}
	return "the key `" + want.Key + "` was planned as " + want.Value +
		" but reads back as " + got
}

// --- kernel refusal mapping (REQ-4, REQ-49) ------------------------------

// kernelFailure maps one kernel refusal onto its mirrored CLI code:
// `flow-<kind>` with underscores hyphenated, one-to-one over the closed
// five-kind set (REQ-4, REQ-49). The CLI mints no kind of its own.
// kernelCode is the one place a kernel refusal kind becomes a CLI code.
// kernelFailure emits through it and the help text spells its codes
// through it, so a kind renamed in the kernel renames itself on both the
// wire and in `--help-all` — the published vocabulary cannot drift from
// the one a caller actually receives.
func kernelCode(kind resolve.RefusalKind) string {
	return "flow-" + strings.ReplaceAll(string(kind), "_", "-")
}

// The kernel-mirrored codes `flow resolve` can refuse under, named so the
// help text and the emitter share one derivation.
var (
	codeUnmodeledOutcome      = kernelCode(resolve.KindUnmodeledOutcome)
	codeNoMatch               = kernelCode(resolve.KindNoMatch)
	codeAmbiguousMatch        = kernelCode(resolve.KindAmbiguousMatch)
	codeOwnedStateUnavailable = kernelCode(resolve.KindOwnedStateUnavailable)
	codeGuardUnevaluable      = kernelCode(resolve.KindGuardUnevaluable)
)

func kernelFailure(refusal resolve.Refusal) *clierr.CLIError {
	code := kernelCode(refusal.Kind)

	ce := &clierr.CLIError{Code: code, Group: clierr.GroupUserEnv}

	switch refusal.Kind {
	case resolve.KindUnmodeledOutcome:
		// Single subject: the offending outcome (REQ-79, REQ-92).
		ce.Param = refusal.Recognized
		ce.Message = "the recognized outcome `" + refusal.Recognized +
			"` is outside the model's declared alphabet"

	case resolve.KindNoMatch:
		ce.Message = "no rule matches the recognized outcome `" +
			refusal.Recognized + "` over the assembled state"
		ce.Findings = rowFindings(refusal.Rows, code,
			"this rule was considered and did not match")
		if len(ce.Findings) == 0 {
			// The kernel named no rows, but `findings[]` is this code's
			// carrier and a caller needs something to act on.
			ce.Findings = []clierr.Finding{{
				Code: code,
				Message: "no rule in the model responds to the recognized " +
					"outcome `" + refusal.Recognized + "`",
			}}
		}

	case resolve.KindAmbiguousMatch:
		ce.Message = "more than one rule matches the recognized outcome `" +
			refusal.Recognized + "`; resolve does not choose among them"
		// EVERY conflicting row is named. Reporting one would be exactly
		// the tie-break REQ-53 forbids.
		ce.Findings = rowFindings(refusal.Rows, code,
			"this rule matches and conflicts with the others named here")

	case resolve.KindOwnedStateUnavailable:
		ce.Message = "the transition requires owned state the declared " +
			"readers did not establish"
		for _, key := range refusal.MissingOwned {
			ce.Findings = append(ce.Findings, clierr.Finding{
				Code: code,
				Key:  key,
				Message: "the owned tag `" + key + "` is required by a " +
					"candidate rule and no reader established it",
			})
		}

	case resolve.KindGuardUnevaluable:
		ce.Message = "a guard predicate could not be decided over the " +
			"assembled state"
		ce.Findings = undecidedFindings(refusal.Undecided, code)
	}

	return ce
}

// rowFindings names one finding per implicated row.
func rowFindings(rows []resolve.RowRef, code, message string) []clierr.Finding {
	out := make([]clierr.Finding, 0, len(rows))
	for _, row := range rows {
		out = append(out, clierr.Finding{
			Code:    code,
			Rule:    row.RuleID,
			Locator: row.SourceLocator,
			Message: "rule `" + row.RuleID + "`: " + message,
		})
	}
	return out
}

// undecidedFindings names every undecidable row and, for each, its
// unevaluable atoms. The atom's four fields ride FLAT on the record, never
// nested under an `atom` object (REQ-15).
func undecidedFindings(rows []resolve.UndecidedRow, code string) []clierr.Finding {
	var out []clierr.Finding
	for _, row := range rows {
		if len(row.Atoms) == 0 {
			out = append(out, clierr.Finding{
				Code: code, Rule: row.RuleID, Locator: row.SourceLocator,
				Message: "rule `" + row.RuleID + "` carries a guard that " +
					"could not be decided",
			})
			continue
		}
		for _, atom := range row.Atoms {
			out = append(out, clierr.Finding{
				Code:     code,
				Rule:     row.RuleID,
				Locator:  row.SourceLocator,
				Key:      atom.Key,
				Block:    string(atom.Block),
				Operator: atom.Operator,
				Literal:  atom.Literal,
				Message: "rule `" + row.RuleID + "`: the atom on `" +
					atom.Key + "` could not be decided (" +
					string(atom.Reason) + ")",
			})
		}
	}
	return out
}

// --- gates (JDR 0001 §D9) ------------------------------------------------

// gateResult is one gate's reported disposition on a candidate or a plan.
type gateResult struct {
	ID     string `json:"id"`
	Result string `json:"result"`
	Reason string `json:"reason,omitempty"`
}

// runGates invokes every gate on one row and reports each result.
//
// No gate result is ever dropped silently, and a gate that could not be
// CONSULTED is exit 3 rather than a deny (REQ-47, REQ-55). Conflating the
// two is the silent-failure mode the contract names by name: a denied
// transition and an unreachable gate call for opposite caller responses,
// and reporting the second as the first would tell a caller their model
// refused when in fact nothing decided.
func (r flowRequest) runGates(ctx context.Context, gates []string) (
	[]gateResult, *clierr.CLIError,
) {
	exec := accessor.NewExecutor(r.registry, r.artifactMap())

	out := make([]gateResult, 0, len(gates))
	for _, id := range gates {
		def, ok := r.registry.Lookup(id, accessor.CapGate)
		if !ok {
			return nil, internalErr(codeAccessorUnknown,
				"the gate accessor `"+id+"` is not bound")
		}
		if _, bound := r.artifacts[def.Accessor.Role]; !bound {
			return nil, userErr(codeArtifactMissing, def.Accessor.Role,
				"the artifact role `"+def.Accessor.Role+"` that the gate "+
					"accessor `"+id+"` needs has no --artifact binding")
		}

		result := exec.Gate(ctx, id)
		if result.Refused() &&
			result.Refusal.Class != accessor.ClassGateIndeterminate {
			// Timeout or execution failure: an ACCESSOR refusal, exit 3,
			// never folded into a verdict.
			return nil, accessorFailure(*result.Refusal, phaseGate)
		}
		out = append(out, gateResult{
			ID:     id,
			Result: string(result.Verdict),
			Reason: result.Reason,
		})
	}
	return out, nil
}

// gateVerdictFailure turns a set of reported gate results into the refusal
// they imply for `flow resolve`, or nil when every gate allowed.
//
// Deny OVERRIDES allow and indeterminate, and `flow-gate-indeterminate`
// fires only when no gate denied and at least one is indeterminate
// (REQ-50, REQ-51, `0005:FM`). Precedence runs that way because a deny is a
// decision the model MADE, while indeterminate is one it could not make —
// reporting the weaker disposition would understate what the model said.
//
// Precedence selects the CODE only. `findings[]` carries ONE ENTRY PER
// GATE ON THE ROW — the ALLOWING gates included — because REQ-50 requires
// that "every gate result MUST be reported" and REQ-47 that "no gate result
// is ever dropped silently". Partitioning the results and reporting only
// the matching partition dropped every allow, which is the same row
// `flow next --evaluate-gates` reports in full through its candidate
// `gates[]`; the two verbs must agree about what the gates said.
func gateVerdictFailure(results []gateResult) *clierr.CLIError {
	var denied, indeterminate int
	for _, g := range results {
		switch accessor.Verdict(g.Result) {
		case accessor.VerdictDeny:
			denied++
		case accessor.VerdictIndeterminate:
			indeterminate++
		}
	}

	switch {
	case denied > 0:
		return gateRefusal(codeGateDenied,
			"a gate on the selected rule denied the transition", results)
	case indeterminate > 0:
		return gateRefusal(codeGateIndeterminate,
			"a gate on the selected rule could not decide", results)
	default:
		return nil
	}
}

// gateRefusal builds a gate refusal carrying ONE findings entry per gate on
// the selected row (REQ-51, REQ-97, REQ-98), whatever that gate answered.
//
// Each finding's own `code` names THAT gate's disposition, so an allowing
// gate is not mislabelled as the refusal that the denying one caused: the
// envelope `code` is the row's verdict, and the per-finding codes are the
// per-gate results the caller needs in order to know which gate to fix.
func gateRefusal(code, message string, gates []gateResult) *clierr.CLIError {
	findings := make([]clierr.Finding, 0, len(gates))
	for _, g := range gates {
		verdict := accessor.Verdict(g.Result)
		findings = append(findings, clierr.Finding{
			Code:    gateFindingCode(verdict, code),
			Param:   g.ID,
			Message: gateMessage(g.ID, verdict, g.Reason),
		})
	}
	return &clierr.CLIError{
		Code: code, Message: message, Group: clierr.GroupUserEnv,
		Findings: findings,
	}
}

// gateFindingCode names one gate's own disposition within a gate refusal.
// A deny takes `flow-gate-denied` and an indeterminate
// `flow-gate-indeterminate` regardless of which one carried the envelope,
// so the finding set stays a faithful transcript of what every gate
// answered (REQ-47, REQ-50). A gate that ALLOWED reports under the
// envelope's code — it is a result on that refusal, and its `message`
// states the verdict in prose (REQ-17), so no invented code is needed.
func gateFindingCode(verdict accessor.Verdict, envelope string) string {
	switch verdict {
	case accessor.VerdictDeny:
		return codeGateDenied
	case accessor.VerdictIndeterminate:
		return codeGateIndeterminate
	default:
		return envelope
	}
}

// gateMessage renders a gate's disposition as prose that stands ALONE —
// REQ-17 requires the message render the failure readably with no
// structured field consulted, so the gate id and its verdict are in the
// sentence rather than only in `param`.
func gateMessage(id string, verdict accessor.Verdict, reason string) string {
	var what string
	switch verdict {
	case accessor.VerdictDeny:
		what = "the gate `" + id + "` denied the transition"
	case accessor.VerdictIndeterminate:
		what = "the gate `" + id + "` could not decide the transition"
	default:
		what = "the gate `" + id + "` answered " + string(verdict)
	}
	if reason != "" {
		return what + ": " + reason
	}
	return what
}

// --- tag rendering -------------------------------------------------------

// tagMap renders tags as a JSON object. A set value is already the
// canonical literal by the time it reaches here, so it is carried VERBATIM
// as a string rather than re-encoded — re-encoding would be a second
// encoder, which REQ-71 forbids, and would break byte equality with the
// request that produced it.
func tagMap(tags []resolve.Tag) map[string]string {
	out := make(map[string]string, len(tags))
	for _, t := range tags {
		out[t.Key] = t.Value
	}
	return out
}

func observedTagMap(tags []resolveTag) map[string]string {
	out := make(map[string]string, len(tags))
	for _, t := range tags {
		out[t.Key] = t.Value
	}
	return out
}

func kernelTags(tags []resolveTag) []resolve.Tag {
	out := make([]resolve.Tag, 0, len(tags))
	for _, t := range tags {
		out = append(out, resolve.Tag{Key: t.Key, Value: t.Value})
	}
	return out
}

// buildRequest performs the whole pre-accessor input phase for a verb.
func buildRequest(cmd *cobra.Command, withTags bool) (flowRequest, *clierr.CLIError) {
	model, ref, ce := selectModel(cmd)
	if ce != nil {
		return flowRequest{}, ce
	}
	artifacts, ce := parseArtifacts(cmd)
	if ce != nil {
		return flowRequest{}, ce
	}

	var observed []resolveTag
	if withTags {
		observed, ce = parseTags(cmd, model)
		if ce != nil {
			return flowRequest{}, ce
		}
	}

	// `0025:C6` — the flag is read with the error CHECKED, not discarded.
	// Persistent registration on the `flow` group guarantees presence, so
	// a lookup miss can only be a WIRING BUG and must surface as one:
	// pflag's own text is "flag accessed but not defined", a programmer
	// error, never a statement about user intent. A default-deny resting
	// on the zero value would make fail-closed a coincidence between two
	// independent decisions, and the dangerous edit — a parent registering
	// the flag, or a verb registering a different default — silently
	// converts the miss into a hit and fails OPEN.
	allowCommands, ferr := cmd.Flags().GetBool("allow-commands")
	if ferr != nil {
		return flowRequest{}, internalErr(codeAccessorFailed,
			"the --allow-commands gate is not registered on this command: "+
				ferr.Error())
	}

	// `0025:C2` — the root a separator-bearing argv0 resolves against is
	// the MODEL FILE's directory, absolutized here because `--model` is
	// stored verbatim and a relative one must not make argv0 resolution
	// cwd-dependent. The loader performs no path resolution (`0002:EIA`),
	// so the base dir is carried at this layer and never on `table.Model`.
	baseDir := ""
	if abs, aerr := filepath.Abs(ref); aerr == nil {
		baseDir = filepath.Dir(abs)
	}

	return flowRequest{
		model:     model,
		modelRef:  ref,
		artifacts: artifacts,
		observed:  observed,
		registry:  flowbind.Registry(model, baseDir, allowCommands),
	}, nil
}
