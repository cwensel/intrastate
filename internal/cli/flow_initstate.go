package cli

// RDR 0019 — `flow init-state`: the runtime bootstrap for owned state.
//
// `[initial]` is BOTH the lint-time reachability root (0006:C18 unchanged)
// AND the runtime bootstrap source. Materializing it is an EXPLICIT,
// PERSISTING act: no read verb, no artifact load, and no accessor read path
// synthesizes an `[initial]` value for an absent key — this verb is the one
// route that establishes them, and it does so through the declared write
// accessors with commit-time read-back, exactly as `set-state` does.
//
// Four gates run in a FIXED order before anything is planned (`0019:C1`
// ORDERING), and the order is observable rather than incidental:
//
//  1. CLASS. A `decision-table` model has no `[initial]` to materialize, so
//     initialization refuses rather than succeeding vacuously. It keys on
//     `table.IsDecisionTable` ALONE — never on `len(owned)`, which a
//     state-machine model declaring no owned tags shares.
//  2. CARRIER. Everything the predicate says about STORE keys, the seal, and
//     artifact bytes is scoped to the file-backed carrier, so a needed role
//     whose write OR read binding is not the file-backed type refuses:
//     "the store carries no key" is UNDEFINED for an edit-carried or
//     command-backed accessor. It is decided AT registry construction, which
//     is what makes it preempt the `--allow-commands` refusal — both of that
//     gate's raise sites are PAST construction.
//  3. ROLE BINDING. Every role a needed writer names must be bound, or the
//     verb refuses without reading any store. "Bound" is not a filter that
//     can shrink the emptiness quantifier's domain.
//  4. EMPTINESS. Seeding is ALL-OR-NOTHING over an EMPTY store: it seeds if
//     and only if EVERY bound artifact carries NO key, and then seeds every
//     `[initial]` key. The quantifier is ALL, not ANY and not per-artifact —
//     one surviving key anywhere blocks the whole seed, which is what makes
//     the cleared-key guarantee hold per model rather than per artifact.

import (
	"slices"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/cli/cmdbind"
	"github.com/cwensel/intrastate/internal/cli/flowbind"
	"github.com/cwensel/intrastate/internal/cli/respond"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
	"github.com/spf13/cobra"
)

// The two refusal classes NEW to this verb.
//
// Their SPELLINGS are non-normative, on the explicit `0028:C1.3`
// `codeWriteEditRefused` carve-out `0019:C1` cites: what the contract fixes
// is the exit GROUP (2 on both) and that the two are DISTINCT codes,
// distinct from each other and from the shared classes. No test may pin
// either string. The spellings below are C1's own sharpened ones, so the
// record and the wire agree without either binding the other.
const (
	// codeInitClassUnsupported: the model's class carries no `[initial]`
	// to materialize (`0019:C1` SEMANTICS).
	codeInitClassUnsupported = "flow-init-class-unsupported"
	// codeInitCarrierUnsupported: a needed role's write or read binding is
	// not the file-backed type, so store emptiness is undefined for it
	// (`0019:C1` CARRIER SCOPE).
	codeInitCarrierUnsupported = "flow-init-carrier-unsupported"
)

// initStatePayload is `flow init-state`'s verb-specific `data`.
//
// `0019:D-wire-byte-format` defers the field NAMES to implementation and
// fixes the CONTENT: the payload distinguishes the seeded-all case from the
// no-op case, its scope is exactly the `[initial]` key set, it reports keys
// by NAME and never echoes their values — `read-state` is the surface that
// reports values (RT1), and echoing them here would create a second place a
// seeded value can be read from with no invariant tying the two.
type initStatePayload struct {
	Model     string            `json:"model"`
	Revision  string            `json:"revision"`
	Artifacts map[string]string `json:"artifacts"`
	// Seeded names the `[initial]` keys this run wrote, and is non-empty on
	// exactly the seeding arm.
	Seeded []string `json:"seeded"`
	// Absent names the `[initial]` keys the store does not carry. It belongs
	// to the NO-OP arm, the only arm where the set can be non-empty: on the
	// seeded arm every `[initial]` key was just written, so it is
	// necessarily empty there and carries no information.
	Absent []string `json:"absent"`
}

func newFlowInitStateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init-state",
		Short: "Seed owned state from the model's [initial] root",
		Long: `Seed owned state from the model's [initial] root.

init-state writes the model's [initial] assignments through the declared
write accessors, verified by read-back — the same route set-state takes. It
carries no write grammar of its own: the model's own root is the plan.

Seeding is ALL-OR-NOTHING over an EMPTY store. It seeds if and only if every
bound artifact carries no key, and then it seeds every [initial] key. A store
that carries any key — torn, partially seeded, post-clear, or fully seeded
alike — is a no-op SUCCESS at exit 0: nothing is written, and the payload
reports which [initial] keys the store does not carry so the state is visible
without being repaired or resurrected.

The scope is the FILE-BACKED carrier. A model whose needed write or read
accessor is edit-carried or command-backed refuses: store emptiness is
undefined for a carrier with no JSON store.

A decision-table model declares no [initial] and refuses rather than
succeeding vacuously.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE:          runFlowInitState,
	}
	registerSelectionFlags(cmd)
	withExtendedHelp(cmd, flowInitStateExtendedDesc)
	return cmd
}

const flowInitStateExtendedDesc = `init-state answers "make this model's declared root real, once", and
nothing else. It selects no rule, evaluates no guard, and takes no
caller-authored value.

Where the plan comes from

  The plan is the model's [initial] block, taken from the LOADER-NORMALIZED
  assignments rather than transcribed to a --write spelling. Each value is
  rendered in the canonical form for its DECLARED KIND: a set is sorted,
  deduplicated and compact; every scalar kind is its single member verbatim.
  That is what makes read-back equality byte equality and leaves no third
  encoding.

  The loader admits [initial] spellings the argv route does not — a bare
  scalar for a set-valued tag, a single-member array for a scalar tag, an
  empty scalar. All three seed here, because conformance is already held by
  the loader and is not re-established at seed time. Transcribing through
  argv would refuse at seed time models that load clean.

The empty-store predicate

  Seeding is ALL-OR-NOTHING over an EMPTY store, never a per-key merge. The
  quantifier is ALL over every bound artifact, not ANY and not per-artifact:
  one surviving key ANYWHERE blocks the whole seed, so a cleared key is never
  re-established while the store retains at least one other.

  The count is over STORE keys, not owned keys. A read-back-sealed artifact
  carries the seal and no owned key: that is a one-key, NON-EMPTY store, and
  init-state declines to seed it — at exit 0, as a no-op, since no write and
  therefore no read-back happens on that arm.

  Clearing the LAST remaining key empties the store, and an emptied store is
  indistinguishable at the content level from a never-written one: a clear is
  a key REMOVAL, not a tombstone. A subsequent init-state therefore RESEEDS
  such a store. That is the accepted residual of rejecting tombstones, and it
  is the one qualification on the cleared-key guarantee above.

  A key added to [initial] after seeding is re-established by explicit
  set-state, not by init — the payload's absent-key report names it, and
  ` + "`--write`" + ` of the listed keys is the repair route.

Reading the output

  seeded[]  the [initial] keys this run wrote. Non-empty on exactly the
            seeding arm. Keys only: read-state is the surface that reports
            values.
  absent[]  the [initial] keys the store does not carry. It belongs to the
            no-op arm; on the seeded arm every key was just written, so it is
            necessarily empty.

The carrier scope

  Everything above about store keys and artifact bytes is scoped to the
  FILE-BACKED write carrier. An edit-carried or command-backed accessor has
  no JSON store, so "the store carries no key" is undefined for it, and a
  model whose needed write OR read binding is not file-backed refuses at
  exit 2 naming the capability and the accessor. That refusal is decided at
  registry construction, before any accessor is invoked, so it preempts the
  --allow-commands refusal: a command-backed accessor yields the carrier code
  whether or not the opt-in was passed. Extending the predicate to those
  carriers is out of scope here.

Ordering

  The CLASS refusal precedes the CARRIER refusal, which precedes the
  ROLE-BINDING refusal, which precedes the emptiness read. An invocation that
  leaves a needed role unbound refuses at exit 2 regardless of what the bound
  stores contain: it does not reach the predicate and cannot take the no-op
  arm.

  The writer-arity half is not ordered here at all — a model routing one
  [initial] key to two writers refuses at LOAD as ` + codeModelInvalid + `,
  never in the writer-routing family.

Read-back is the commit-time check

  Every seed routes through the declared write accessor with commit-time
  read-back; nothing writes an artifact directly. A read-back that completed
  and DISAGREED is ` + codeReadBackMismatch + ` at exit 2, naming the key
  present-and-unverified. One that could not complete is
  ` + codeReadBackIncomplete + ` at exit 3. Both leave a NON-EMPTY store, so
  a re-run is a no-op that does NOT repair either — recovery is an explicit
  set-state, or discarding the artifact and re-running init.

Shared refusals

  The codes above are the ones specific to this verb. Model selection,
  artifact validation, and the accessor and environment failures are common
  to every flow verb and are listed once under "intrastate flow --help-all".

Exits

  0  every [initial] key seeded, or the store was non-empty and nothing was
     written.
  2  the request or the model is wrong, a role is unbound, the model's class
     or a needed accessor's carrier is unsupported, or the read-back
     disagreed.
  3  a writer or the read-back could not complete.

Worked call

  intrastate flow init-state --model flow.toml \
      --artifact state=state.json --as json`

func runFlowInitState(cmd *cobra.Command, _ []string) error {
	if ce := respond.ValidateMode(cmd); ce != nil {
		return respond.Fail(cmd, ce)
	}

	// `--tag` is not registered on this verb: there is no verdict for
	// observed context to inform, and the plan is the model's own root.
	req, ce := buildRequest(cmd, false)
	if ce != nil {
		return respond.Fail(cmd, ce)
	}

	// GATE 1 — CLASS. Answerable from the loaded model alone, which is what
	// makes class-first the cheaper order as well as the stated one: a
	// decision-table model refuses without a registry being consulted.
	if ce := initClassGate(req.model); ce != nil {
		return respond.Fail(cmd, ce)
	}

	// The NEEDED writers: those serving at least one `[initial]` key. A
	// writer serving no `[initial]` key is not needed, so neither its binding
	// nor its role takes part in the carrier gate, the role-binding check, or
	// the ALL quantifier.
	needed, ce := neededWriters(req.model)
	if ce != nil {
		return respond.Fail(cmd, ce)
	}

	// GATE 2 — CARRIER, over BOTH capabilities. Decided AFTER registry
	// construction (it reads the constructed binding) but BEFORE any accessor
	// is invoked: the gate reads `Identity`, `Accessor`, and `Binding` as
	// FIELDS and invokes no `Binding` method, so "no accessor ran" is true by
	// construction rather than by timing.
	if ce := initCarrierGate(req.registry, needed); ce != nil {
		return respond.Fail(cmd, ce)
	}

	// GATE 3 — ROLE BINDING, before any store is read.
	for _, name := range needed {
		def, ok := req.registry.Lookup(name, accessor.CapWrite)
		if !ok {
			return respond.Fail(cmd, internalErr(codeAccessorUnknown,
				"the write accessor `"+name+"` is not bound"))
		}
		if _, bound := req.artifacts[def.Accessor.Role]; !bound {
			return respond.Fail(cmd, userErr(codeArtifactMissing,
				def.Accessor.Role,
				"the artifact role `"+def.Accessor.Role+"` that the write "+
					"accessor `"+name+"` needs has no --artifact binding"))
		}
	}

	// GATE 4 — EMPTINESS, over the whole bound set. Every needed role's store
	// must carry NO key, or the invocation takes the no-op arm.
	stores, ce := initStoreKeys(req, needed)
	if ce != nil {
		return respond.Fail(cmd, ce)
	}

	payload := initStatePayload{
		Model:     req.modelRef,
		Revision:  req.revision(),
		Artifacts: req.artifacts,
		Seeded:    []string{},
		Absent:    []string{},
	}

	if !initStoresEmpty(stores) {
		// The NO-OP arm. Zero writes; the payload reports which `[initial]`
		// keys the store does not carry, so a torn or post-clear state is
		// VISIBLE without being repaired, resurrected, or failed on.
		absent, ce := initAbsentKeys(req, stores)
		if ce != nil {
			return respond.Fail(cmd, ce)
		}
		payload.Absent = absent
		return respond.OK(cmd, respond.Success{Data: payload})
	}

	// The SEEDING arm. Encode and plan EVERY `[initial]` assignment before
	// any accessor runs, then commit through the declared writers.
	byWriter := map[string][]resolve.Tag{}
	for _, assignment := range req.model.Initial {
		writer, ce := initWriterFor(req.model, assignment.Key)
		if ce != nil {
			return respond.Fail(cmd, ce)
		}
		value, ce := initSeedValue(req.model, assignment)
		if ce != nil {
			return respond.Fail(cmd, ce)
		}
		byWriter[writer] = append(byWriter[writer],
			resolve.Tag{Key: assignment.Key, Value: value})
		payload.Seeded = append(payload.Seeded, assignment.Key)
	}
	slices.Sort(payload.Seeded)

	exec := accessor.NewExecutor(req.registry, req.artifactMap())
	for _, name := range slices.Sorted(maps(byWriter)) {
		result := exec.Write(cmd.Context(), name,
			resolve.Plan{Writes: byWriter[name]})
		if result.Refused() {
			return respond.Fail(cmd, accessorFailure(*result.Refusal, phaseWrite))
		}
	}

	return respond.OK(cmd, respond.Success{Data: payload})
}

// initClassGate refuses a model whose CLASS carries no `[initial]` to
// materialize.
//
// It keys on `table.IsDecisionTable` ALONE. Never on "declares `[initial]`
// but is decision-table" — that state is unconstructible, since such a model
// refuses at LOAD — and never on a re-derivation from `len(owned)` or
// `len(Model.Initial)`, which a state-machine model declaring no owned tags
// shares with every admitted decision table. The refusal is a CLASS refusal,
// not an emptiness one: an ordinary state-machine model whose `[initial]` is
// empty is a different case, already blocked at lint by 0006:C18.
func initClassGate(m *table.Model) *clierr.CLIError {
	if !table.IsDecisionTable(m) {
		return nil
	}
	return userErr(codeInitClassUnsupported, "model",
		"the model declares class `"+table.ClassDecisionTable+"`, which has "+
			"no `[initial]` root to materialize; initialization is for "+
			"state-machine models")
}

// neededWriters returns the write accessors serving at least one `[initial]`
// key, sorted.
//
// "Needed" is what scopes the carrier gate, the role-binding check, and the
// emptiness quantifier alike: a write accessor serving no `[initial]` key
// contributes nothing to the plan, so its binding, its role, and its store
// take no part in any of the three.
//
// Arity is not re-checked here beyond selecting the single serving writer:
// `::checkAccessorBindings` folds every `[initial]` key into its `written`
// set and refuses `writerCount[key] != 1` at LOAD, so no model reaching this
// verb can route one key to two writers. `initWriterFor` still reports the
// count faithfully rather than picking, because a plan built on an arbitrary
// destination would leave the read-back confirming through the wrong writer.
func neededWriters(m *table.Model) ([]string, *clierr.CLIError) {
	seen := map[string]bool{}
	for _, assignment := range m.Initial {
		name, ce := initWriterFor(m, assignment.Key)
		if ce != nil {
			return nil, ce
		}
		seen[name] = true
	}
	return slices.Sorted(maps(seen)), nil
}

// initWriterFor names the single declared write accessor serving key.
//
// It reuses `::writerFor`'s COUNT semantics under the shipped
// `flow-write-unbound` code — a shared refusal class this verb reuses
// unchanged, at the group that code already carries. Zero writers and two
// writers are both refusals; neither is reachable through the loader for an
// `[initial]` key, and neither is silently resolved here.
func initWriterFor(m *table.Model, key string) (string, *clierr.CLIError) {
	if ce := writerFor(m, key, codeWriteUnbound); ce != nil {
		return "", ce
	}
	for _, name := range slices.Sorted(maps(m.Writers)) {
		if slices.Contains(m.Writers[name].Keys, key) {
			return name, nil
		}
	}
	// Unreachable: `writerFor` just established exactly one match.
	return "", internalErr(codeWriteNonOwned,
		"no write accessor serves the `[initial]` key `"+key+"`")
}

// initSeedValue renders one `[initial]` assignment in JDR 0001 §D13's
// canonical form, DISPATCHED ON THE DECLARED KIND.
//
// It is exactly the SET arm of `set-state`'s encoder and nothing more: for a
// `set` kind it is `::canonicalSet` — sort, deduplicate, compact, HTML
// escaping off — and for every scalar kind it is the single member VERBATIM,
// `members[0]`. `::canonicalSet` alone is NOT the whole encoder: it takes a
// `[]string` and renders a JSON array, so applying it to a scalar seed would
// persist `["draft"]` where `set-state --write status=draft` persists
// `draft`, breaking the very byte identity RT3 asserts.
//
// `::canonicalValue` is the correct REFERENCE for both arms' output but is
// deliberately NOT the call: it re-runs the argv admission checks, and the
// loader's `[initial]` admission set is a proper SUPERSET of the argv
// route's. Conformance is ALREADY HELD — `::loadInitial` conforms every
// assignment — so a re-conform pass cannot fail and buys nothing, while
// routing through argv would refuse at seed time models that load clean.
func initSeedValue(m *table.Model, assignment table.TagValue) (string, *clierr.CLIError) {
	decl := m.Tags[assignment.Key]
	if decl.Kind == "set" {
		return canonicalSet(assignment.Value), nil
	}
	if len(assignment.Value) != 1 {
		// Unreachable: `::loadInitial` refuses `decl.Kind != "set" &&
		// len(members) != 1` as a malformed `[initial]` declaration, so a
		// scalar assignment reaching here holds exactly one member.
		return "", internalErr(codeWriteInvalid,
			"the `[initial]` assignment for `"+assignment.Key+
				"` holds a member sequence for a non-set kind")
	}
	return assignment.Value[0], nil
}

// initCarrierGate refuses when EITHER the bound write binding or the bound
// read binding for a needed role is not the FILE-BACKED type.
//
// Everything the predicate says about STORE keys, the seal, and artifact
// bytes is scoped to the file-backed carrier. An edit-carried or
// command-backed accessor has no JSON store, no seal, and no artifact bytes,
// so "the store carries no key" is UNDEFINED for it — and the emptiness
// answer and the commit-time read-back both travel the READ binding, so a
// non-file-backed READER leaves the predicate exactly as undefined as a
// non-file-backed writer does.
//
// The carrier comes from the CONSTRUCTED BINDING'S TYPE, never from a
// re-derivation over `table.Accessor`'s `Edit`/`Command` fields — that would
// duplicate a discriminator whose author is the registry, and would disagree
// with it on the residue the registry routes to a REFUSING binding.
//
// The POINTER/VALUE spelling is normative: the registry constructs writers as
// POINTERS and readers as VALUES, so a read-side case written
// `*flowbind.Reader` would match NOTHING and refuse every model, file-backed
// ones included.
//
// The gate reads `Identity`, `Accessor`, and `Binding` as FIELDS and invokes
// no `Binding` method — a type switch calls nothing — which is what makes the
// no-accessor-invoked property true by construction.
func initCarrierGate(reg accessor.Registry, needed []string) *clierr.CLIError {
	for _, name := range needed {
		def, ok := reg.Lookup(name, accessor.CapWrite)
		if !ok {
			return internalErr(codeAccessorUnknown,
				"the write accessor `"+name+"` is not bound")
		}
		if carrier := initWriteCarrier(def.Binding); carrier != "" {
			return initCarrierRefusal("write", name, carrier)
		}

		// The role's READER, selected as `::readerFor` selects it: the FIRST
		// definition in REGISTRY ORDER whose capability is read and whose
		// role matches. First match is normative — a gate that refused on
		// ambiguity would be WIDER than the executor's selector, and one
		// selecting the last would disagree with it — so `Definitions` is
		// scanned in the registry's own slice order, unsorted and unfiltered.
		reader, found := initReaderFor(reg, def.Accessor.Role)
		if !found {
			// A no-match is the flat `false` the executor's selector returns
			// and reaches the unbound-needed-role arm; it is NOT split into a
			// missing-reader code of its own.
			return userErr(codeArtifactMissing, def.Accessor.Role,
				"the artifact role `"+def.Accessor.Role+"` that the write "+
					"accessor `"+name+"` needs is served by no declared read "+
					"accessor, so its commit-time read-back has no reader")
		}
		if carrier := initReadCarrier(reader.Binding); carrier != "" {
			return initCarrierRefusal("read", reader.Identity.Name, carrier)
		}
	}
	return nil
}

// initReaderFor reproduces `internal/accessor::Registry.readerFor`'s
// semantics over the EXPORTED `Definitions` slice. That selector is
// unexported with its only non-test caller inside its own package, so the
// gate reproduces it rather than calling it; the equivalence rides on sharing
// the registry's own slice order, which is why nothing here sorts, filters,
// or re-orders before scanning.
func initReaderFor(reg accessor.Registry, role string) (accessor.Definition, bool) {
	for _, def := range reg.Definitions {
		if def.Identity.Capability == accessor.CapRead &&
			def.Accessor.Role == role {
			return def, true
		}
	}
	return accessor.Definition{}, false
}

// initWriteCarrier names the write binding's carrier when it is NOT the
// file-backed one, and "" when it is. Only `*flowbind.Writer` is admitted.
func initWriteCarrier(binding accessor.Binding) string {
	switch binding.(type) {
	case *flowbind.Writer:
		return ""
	case *flowbind.EditWriter:
		return "edit"
	case *cmdbind.Writer:
		return "command"
	default:
		return "unsupported"
	}
}

// initReadCarrier names the read binding's carrier when it is NOT the
// file-backed one, and "" when it is. Only `flowbind.Reader` is admitted —
// a VALUE, not a pointer: the registry constructs readers as values, so a
// case written `*flowbind.Reader` would match nothing and refuse every model.
func initReadCarrier(binding accessor.Binding) string {
	switch binding.(type) {
	case flowbind.Reader:
		return ""
	case cmdbind.Reader:
		return "command"
	default:
		return "unsupported"
	}
}

// initCarrierRefusal builds the carrier refusal, naming WHICH CAPABILITY and
// WHICH ACCESSOR failed — one distinct terminal refusal over both
// capabilities, not two gates minting two codes.
func initCarrierRefusal(capability, name, carrier string) *clierr.CLIError {
	ce := userErr(codeInitCarrierUnsupported, name,
		"the "+capability+" accessor `"+name+"` is "+carrier+"-carried, not "+
			"file-backed; init-state's empty-store predicate is defined over "+
			"the file-backed carrier's JSON store alone, so store emptiness "+
			"is undefined for this accessor")
	ce.Hint = "seed this model's owned state with `flow set-state --write`, " +
		"which carries no emptiness predicate"
	return ce
}

// initStoreKeys reads every needed role's store key set, keyed by role.
//
// It is ONE read of the bound artifacts, serving both the emptiness predicate
// and the no-op arm's absent-key report — the two questions C1 poses about a
// store, both answerable from the key set alone and neither requiring a key
// VALUE.
//
// The carrier is `flowbind.StoreKeys`, which answers the store's key set
// WITHOUT entering the key-scoped reader's unreadable short-circuit. A
// carrier that inherited that short-circuit would degrade the sealed-store arm
// from no-op success to an exit-3 refusal.
func initStoreKeys(req flowRequest, needed []string) (map[string][]string, *clierr.CLIError) {
	out := map[string][]string{}
	for _, role := range slices.Sorted(maps(initNeededRoles(req, needed))) {
		keys, err := flowbind.StoreKeys(req.artifacts[role])
		if err != nil {
			return nil, userErr(codeArtifactInvalid, role,
				"the artifact bound to role `"+role+"` could not be read: "+
					err.Error())
		}
		out[role] = keys
	}
	return out, nil
}

// initStoresEmpty answers the ALL quantifier: whether EVERY bound artifact
// the needed writers name carries NO key.
//
// The count is over STORE keys, not owned keys, and the difference is
// reachable: a read-back-SEALED artifact carries the seal and nothing else
// once its owned keys are cleared, which is a ONE-key, NON-EMPTY store that
// init-state declines to seed — at exit 0, as a no-op, since no write and
// therefore no read-back happens on that arm.
//
// No implementation may substitute "every `[initial]` key reads absent" for
// this: that is the per-key variant the contract rejects, since it cannot see
// a non-`[initial]` key — the seal included — and so would seed into a
// non-empty store.
func initStoresEmpty(stores map[string][]string) bool {
	for _, keys := range stores {
		if len(keys) != 0 {
			return false
		}
	}
	return true
}

// initNeededRoles is the set of roles the needed writers name — the exact
// domain the emptiness quantifier ranges over. "Bound" is not a filter that
// can shrink it: gate 3 already established that every one of these roles IS
// bound, so an unbound needed role is a REFUSAL rather than an artifact
// treated as empty.
func initNeededRoles(req flowRequest, needed []string) map[string]bool {
	roles := map[string]bool{}
	for _, name := range needed {
		if def, ok := req.registry.Lookup(name, accessor.CapWrite); ok {
			roles[def.Accessor.Role] = true
		}
	}
	return roles
}

// initAbsentKeys reports which `[initial]` keys the bound stores do not
// carry — the no-op arm's informational report.
//
// Its scope is exactly the `[initial]` key set: the verb claims nothing about
// owned keys `[initial]` does not assign, so a store carrying such a key is
// not described here and a key outside `[initial]` never appears.
//
// It answers from the same store key sets the predicate read, which is what
// keeps the sealed store a no-op SUCCESS: routing this report through the
// declared readers would re-enter the unreadable short-circuit and turn the
// arm into an exit-3 refusal.
//
// The question is asked PER ROLE, like every other quantifier in this verb:
// a key is present only if the store bound to the role its OWN declared
// writer names carries it. Answering from the UNION of stores lets a key
// carried by an unrelated role suppress its own absent entry, which makes
// `0019:F2`'s "it names precisely the keys to pass" false — the repair route
// is `set-state` of the listed keys, and `set-state` routes each key to that
// same writer's artifact.
func initAbsentKeys(req flowRequest, stores map[string][]string) ([]string, *clierr.CLIError) {
	absent := []string{}
	for _, assignment := range req.model.Initial {
		role, ce := initRoleFor(req, assignment.Key)
		if ce != nil {
			return nil, ce
		}
		if !slices.Contains(stores[role], assignment.Key) {
			absent = append(absent, assignment.Key)
		}
	}
	slices.Sort(absent)
	return absent, nil
}

// initRoleFor names the artifact role serving key: the role its single
// declared write accessor names.
//
// It is `initNeededRoles`' resolution for ONE key — `initWriterFor` to the
// writer, then the registry to that writer's role — and it is deliberately
// the same two steps, so the domain the absent report ranges over cannot
// drift from the domain the emptiness predicate read. Both refusals are
// unreachable from the no-op arm: `neededWriters` already resolved every
// `[initial]` key's writer, and gate 3 already looked every needed writer up.
func initRoleFor(req flowRequest, key string) (string, *clierr.CLIError) {
	name, ce := initWriterFor(req.model, key)
	if ce != nil {
		return "", ce
	}
	def, ok := req.registry.Lookup(name, accessor.CapWrite)
	if !ok {
		return "", internalErr(codeAccessorUnknown,
			"the write accessor `"+name+"` is not bound")
	}
	return def.Accessor.Role, nil
}
