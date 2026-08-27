package cli

// RDR 0005 — `flow read-state` and `flow set-state`: the two state verbs.
//
// They sit either side of the same artifact and are deliberately
// asymmetric about what they narrow by:
//
//   - `read-state` is DIAGNOSTIC. It has no candidate set to narrow by, so
//     it runs EVERY declared reader (REQ-37) — the single exception to the
//     narrowing rule the other two verbs follow. That is also why an
//     unbound role refuses here while the same argv succeeds under `next`.
//   - `set-state` MUTATES. Every input refusal precedes the accessor, and
//     success is reported only after read-back verification (REQ-65) —
//     "the write returned no error" is not the same claim as "the value is
//     there", and only the second is one a caller can build on.
//
// Neither verb runs a gate. A gate answers allow, deny, or indeterminate
// rather than a tag value, so coercing one into state would invent a fact
// the model never produced (REQ-59, REQ-64, REQ-115).

import (
	"slices"
	"strconv"
	"strings"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/cli/respond"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
	"github.com/spf13/cobra"
)

// --- flow read-state -----------------------------------------------------

// readStatePayload is `flow read-state`'s verb-specific `data` (REQ-77).
type readStatePayload struct {
	Model     string            `json:"model"`
	Revision  string            `json:"revision"`
	Artifacts map[string]string `json:"artifacts"`
	// Readers carries each reader's DECLARED key set beside the tags it
	// returned, so "absent from the artifact" and "not requested" are
	// distinguishable from the payload alone: a key in `keys` with no entry
	// in `tags` is absent, and a key in neither was never requested of this
	// reader.
	Readers []readerOutput `json:"readers"`
}

func newFlowReadStateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read-state",
		Short: "Report what the declared read accessors see",
		Long: `Report what the declared read accessors see.

read-state is diagnostic, so it runs EVERY declared reader rather than the
narrowed set flow next and flow resolve invoke. Every declared reader's
artifact role must therefore be bound.

Each reader reports its declared key set beside the tags it returned, so a
key that is absent from the artifact is distinguishable from one this
reader was never asked for.

It invokes no gate accessor: a gate answers allow, deny, or indeterminate
rather than a tag value.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE:          runFlowReadState,
	}
	registerSelectionFlags(cmd)
	withExtendedHelp(cmd, flowReadStateExtendedDesc)
	return cmd
}

const flowReadStateExtendedDesc = `read-state answers "what does the environment currently look like to
this model", and nothing else. It selects no rule, evaluates no guard,
and writes nothing.

Every declared reader must be bound

  next and resolve narrow to the readers their candidate rows demand, so
  an artifact no candidate needs may go unbound there. read-state is
  diagnostic: it runs EVERY declared reader, so every declared artifact
  role must be bound or the call refuses with ` + codeArtifactMissing + `.
  That is the intended asymmetry — a diagnostic that quietly skipped a
  reader would hide the very binding you ran it to check.

Reading the output

  Each reader reports its DECLARED key set beside the tags it actually
  returned. That is what makes the three states distinguishable:

    present    the key is in the artifact, with its value.
    absent     the reader was asked for the key and established that the
               artifact does not carry it.
    not asked  the key is outside this reader's declared set, so its
               absence here says nothing about the artifact.

  A reader that returns fewer keys than it declared without establishing
  their absence is ` + codeReadIncomplete + ` at exit 3 — an environment
  failure, not an empty result.

No gates

  read-state invokes no gate accessor. A gate answers allow, deny, or
  indeterminate; it does not return a tag value, so there is nothing for
  a state report to carry.

--tag is not accepted here: there is no verdict for observed context to
inform.

Shared refusals

  The codes above are the ones specific to this verb. Model selection,
  tag and artifact validation, and the accessor and environment failures
  are common to every flow verb and are listed once under
  "intrastate flow --help-all" rather than repeated here.

Exits

  0  every declared reader reported.
  2  the request or the model is wrong, or a role is unbound.
  3  a reader timed out, failed, or returned an incomplete key set.

Worked call

  intrastate flow read-state --model flow.toml \
      --artifact state=state.json --as json`

func runFlowReadState(cmd *cobra.Command, _ []string) error {
	if ce := respond.ValidateMode(cmd); ce != nil {
		return respond.Fail(cmd, ce)
	}

	req, ce := buildRequest(cmd, false)
	if ce != nil {
		return respond.Fail(cmd, ce)
	}

	readers, _, ce := req.runReaders(cmd.Context(), declaredReaders(req.model))
	if ce != nil {
		return respond.Fail(cmd, ce)
	}

	return respond.OK(cmd, respond.Success{Data: readStatePayload{
		Model:     req.modelRef,
		Revision:  req.revision(),
		Artifacts: req.artifacts,
		Readers:   readers,
	}})
}

// --- flow set-state ------------------------------------------------------

// setStatePayload is `flow set-state`'s verb-specific `data` (REQ-78).
type setStatePayload struct {
	Model     string            `json:"model"`
	Revision  string            `json:"revision"`
	Artifacts map[string]string `json:"artifacts"`
	Writers   []string          `json:"writers"`
	Writes    map[string]string `json:"writes"`
	Clear     []string          `json:"clear"`
	// Owned is the READ-BACK-CONFIRMED owned tag values, not the request
	// echoed. A cleared key is absent from it: read-back verified the
	// removal, and a verified removal is not a written value (`0004:C11`).
	Owned map[string]string `json:"owned"`
}

func newFlowSetStateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-state",
		Short: "Apply planned owned-tag writes, verified by read-back",
		Long: `Apply planned owned-tag writes and verify them by read-back.

Writes are given as --write name=value and --clear <key>. A set value is a
JSON array literal; the sentinel <clear> is unauthorable as a value, so use
--clear.

Success is reported ONLY after the accessor layer's read-back confirms the
planned values and that non-owned tags are unchanged. A read-back that could
not complete exits 3 and says the write may have been applied.

--tag on set-state is context only and is never written. Nothing links a
set-state request to a prior resolve: the request stands on its own grammar
and the read-back is the only commit-time check.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE:          runFlowSetState,
	}
	registerSelectionFlags(cmd)
	registerTagFlag(cmd)
	cmd.Flags().StringArray("write", nil,
		"planned owned-tag write, as name=value (repeatable)")
	cmd.Flags().StringArray("clear", nil,
		"owned tag key to remove (repeatable)")
	withExtendedHelp(cmd, flowSetStateExtendedDesc)
	return cmd
}

const flowSetStateExtendedDesc = `The write grammar

  --write name=value   set an owned tag. Repeatable. A set value is a
                       JSON array literal, canonical on the wire: members
                       sorted, deduplicated, compact.
  --clear <key>        remove an owned tag. Repeatable.

  The two are separate flags because the reserved value <clear> is
  unauthorable: it can never cross as a write VALUE. A plan's removals
  therefore arrive in the plan's clear[] list, and you transcribe them
  with --clear — never as --write key=<clear>, which refuses.

  Only OWNED tags are writable. Naming a non-owned key refuses with
  ` + codeWriteNonOwned + `, and a key no declared writer serves refuses
  with ` + codeWriteUnbound + ` (` + codeClearUnbound + ` for --clear).

  --tag on set-state is CONTEXT ONLY and is never written. Its refusals
  still bite — an owned key through --tag is refused here as anywhere —
  but a parsed --tag value cannot become a write.

Read-back is the commit-time check

  Success is reported ONLY after the accessor layer reads the artifact
  back and confirms both that the planned values landed and that
  non-owned tags are unchanged. owned{} in the output is the
  READ-BACK-CONFIRMED state, not your request echoed: a cleared key is
  absent from it, because a verified removal is not a written value.

  Writes are grouped by the writer that serves each key. Each writer
  applies its own keys and reads back. No cross-writer atomicity is
  promised, and this says so rather than implying one: a multi-writer
  request can leave one writer's keys applied and another's not.

  ` + codeReadBackMismatch + ` at exit 2 means the read-back completed and
  disagreed. ` + codeReadBackIncomplete + ` and ` + codeReadBackTimeout + `
  at exit 3 mean it could not complete: the write MAY have been applied
  and was not verified. That is never reported as a write that did not
  occur — inspect the artifact before retrying.

Nothing links this to a prior resolve

  set-state validates its own request against the model's own grammar.
  It does not know a resolve call happened, and there is no token or
  session to carry. Transcribing a plan is a caller convenience; the
  read-back is the only guarantee.

Shared refusals

  The codes above are the ones specific to this verb. Model selection,
  tag and artifact validation, and the accessor and environment failures
  are common to every flow verb and are listed once under
  "intrastate flow --help-all" rather than repeated here.

Exits

  0  every planned write applied and verified by read-back.
  2  the request or the model is wrong, or the read-back disagreed.
  3  a writer or the read-back could not complete.

Worked call

  intrastate flow set-state --model flow.toml \
      --artifact state=state.json \
      --write status=approved --clear draft_note --as json`

func runFlowSetState(cmd *cobra.Command, _ []string) error {
	if ce := respond.ValidateMode(cmd); ce != nil {
		return respond.Fail(cmd, ce)
	}

	// `--tag` is parsed so its refusals still bite (an owned key through
	// --tag is refused here as anywhere), but the parsed values are context
	// only and never become writes (REQ-34).
	req, ce := buildRequest(cmd, true)
	if ce != nil {
		return respond.Fail(cmd, ce)
	}

	planned, clears, ce := parseWrites(cmd, req.model)
	if ce != nil {
		return respond.Fail(cmd, ce)
	}

	// Group the planned mutations by the writer that serves them. Each
	// writer applies its OWN keys and reads back; no cross-writer atomicity
	// is promised, and the contract says so rather than implying one
	// (REQ-67).
	byWriter, ce := groupByWriter(req.model, planned)
	if ce != nil {
		return respond.Fail(cmd, ce)
	}

	for _, name := range slices.Sorted(maps(byWriter)) {
		def, ok := req.registry.Lookup(name, accessor.CapWrite)
		if !ok {
			return respond.Fail(cmd, internalErr(codeAccessorUnknown,
				"the write accessor `"+name+"` is not bound"))
		}
		if _, bound := req.artifacts[def.Accessor.Role]; !bound {
			return respond.Fail(cmd, userErr(codeArtifactMissing, def.Accessor.Role,
				"the artifact role `"+def.Accessor.Role+"` that the write "+
					"accessor `"+name+"` needs has no --artifact binding"))
		}
	}

	exec := accessor.NewExecutor(req.registry, req.artifactMap())
	writers := make([]string, 0, len(byWriter))
	confirmed := map[string]string{}

	for _, name := range slices.Sorted(maps(byWriter)) {
		result := exec.Write(cmd.Context(), name,
			resolve.Plan{Writes: byWriter[name]})
		if result.Refused() {
			return respond.Fail(cmd, accessorFailure(*result.Refusal, phaseWrite))
		}
		writers = append(writers, name)
		for _, t := range result.Written {
			confirmed[t.Key] = t.Value
		}
	}

	payload := setStatePayload{
		Model:     req.modelRef,
		Revision:  req.revision(),
		Artifacts: req.artifacts,
		Writers:   writers,
		Writes:    map[string]string{},
		Clear:     clears,
		Owned:     confirmed,
	}
	for _, t := range planned {
		if t.Value == table.ClearSentinel {
			continue
		}
		payload.Writes[t.Key] = t.Value
	}
	if payload.Clear == nil {
		payload.Clear = []string{}
	}

	return respond.OK(cmd, respond.Success{Data: payload})
}

// parseWrites parses `--write name=value` and `--clear <key>` into one
// planned mutation set.
//
// Every refusal here precedes any accessor (REQ-61). The four the contract
// names, in the order a caller meets them:
//
//   - the literal `<clear>` as a value is `flow-write-invalid` with a hint
//     naming `--clear`. The sentinel is unauthorable because a stored
//     literal and a removal are indistinguishable on read-back, so allowing
//     it would make the verification unable to tell them apart (REQ-62);
//   - a key given twice, under either flag or across both, is
//     `flow-write-duplicate` — two mutations for one key have no defined
//     order, so the request is ambiguous rather than merely redundant;
//   - a value malformed for its declared kind is `flow-write-invalid`;
//   - a key no single `[write.<id>].keys` serves is `flow-write-unbound`
//     or `flow-clear-unbound`.
func parseWrites(cmd *cobra.Command, m *table.Model) ([]resolve.Tag, []string, *clierr.CLIError) {
	rawWrites, _ := cmd.Flags().GetStringArray("write")
	rawClears, _ := cmd.Flags().GetStringArray("clear")

	seen := map[string]bool{}

	var planned []resolve.Tag
	for _, entry := range rawWrites {
		key, value, ok := strings.Cut(entry, "=")
		if !ok || key == "" {
			return nil, nil, userErr(codeWriteInvalid, "write",
				"--write takes name=value; got "+entry)
		}
		if value == table.ClearSentinel {
			ce := userErr(codeWriteInvalid, key,
				"`"+table.ClearSentinel+"` is a reserved sentinel and cannot "+
					"be written as a value")
			ce.Hint = "use `--clear " + key + "` to remove the key"
			return nil, nil, ce
		}
		if seen[key] {
			return nil, nil, userErr(codeWriteDuplicate, key,
				"the key `"+key+"` is given more than once across --write "+
					"and --clear")
		}
		seen[key] = true

		// The BINDING check precedes the value check: a key no writer serves
		// has no business being told its value is out of domain, and a
		// caller who mistyped the key must hear about the key.
		if ce := writerFor(m, key, codeWriteUnbound); ce != nil {
			return nil, nil, ce
		}
		// A `--write` key always has a declaration — `writerFor` just proved
		// a writer names it, and a writer's `keys` are declared tags — so
		// the lookup here is total, unlike `parseTags`'s.
		canonical, ce := canonicalValue(key, value, m.Tags[key], "write")
		if ce != nil {
			return nil, nil, ce
		}
		planned = append(planned, resolve.Tag{Key: key, Value: canonical})
	}

	clears := make([]string, 0, len(rawClears))
	for _, key := range rawClears {
		if key == "" {
			return nil, nil, userErr(codeWriteInvalid, "clear",
				"--clear takes a tag key; got an empty value")
		}
		if seen[key] {
			return nil, nil, userErr(codeWriteDuplicate, key,
				"the key `"+key+"` is given more than once across --write "+
					"and --clear")
		}
		seen[key] = true

		if ce := writerFor(m, key, codeClearUnbound); ce != nil {
			return nil, nil, ce
		}
		clears = append(clears, key)
		// A clear travels to the accessor layer as the reserved sentinel —
		// RDR 0002's normalized spelling for a removal — and the binding
		// turns it into an actual deletion (`0004:C11`).
		planned = append(planned, resolve.Tag{Key: key, Value: table.ClearSentinel})
	}

	return planned, clears, nil
}

// writerFor refuses a key no single `[write.<id>]` serves, under the code
// the flag it arrived on takes (REQ-61, REQ-86, REQ-87).
//
// The check is against the WRITER's declared keys rather than against the
// model's owned tag set, because a key can be owned and still served by no
// writer — `note` in the MVV fixture is exactly that. Refusing on ownership
// alone would report the wrong reason, and refusing later would let the
// request reach an accessor it has no authority to invoke.
//
// "No single" is a COUNT, not an existence check: two writers naming one key
// is as unbound as none. Returning on the first match would leave
// `groupByWriter` to pick the lexicographically first of the candidates, and
// the read-back would confirm the value through that same writer — so an
// arbitrary write destination would surface nowhere. Both spellings reuse
// the published code; RDR 0005's refusal table is closed, and REQ-86 /
// REQ-87's "served by no single `[write.<id>].keys`" already covers the
// two-writer case.
func writerFor(m *table.Model, key, code string) *clierr.CLIError {
	n := 0
	for _, acc := range m.Writers {
		if slices.Contains(acc.Keys, key) {
			n++
		}
	}
	switch {
	case n == 1:
		return nil
	case n == 0:
		return userErr(code, key,
			"no declared write accessor's `keys` list names the tag `"+key+"`")
	default:
		return userErr(code, key,
			"the tag `"+key+"` is named by "+strconv.Itoa(n)+
				" write accessors; want exactly one")
	}
}

// groupByWriter assigns each planned mutation to the writer whose declared
// keys name it.
func groupByWriter(m *table.Model, planned []resolve.Tag) (map[string][]resolve.Tag, *clierr.CLIError) {
	out := map[string][]resolve.Tag{}
	for _, t := range planned {
		var owner string
		for _, name := range slices.Sorted(maps(m.Writers)) {
			if slices.Contains(m.Writers[name].Keys, t.Key) {
				owner = name
				break
			}
		}
		if owner == "" {
			// Unreachable: parseWrites already refused an unserved key.
			return nil, internalErr(codeWriteNonOwned,
				"no write accessor serves the planned key `"+t.Key+"`")
		}
		out[owner] = append(out[owner], t)
	}
	return out, nil
}

func maps[V any](m map[string]V) func(func(string) bool) {
	return func(yield func(string) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}
