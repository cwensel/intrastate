package cli

// RDR 0005 — the input families refused at the CLI BEFORE any accessor
// runs: model selection (REQ-22..REQ-25), state input (REQ-26..REQ-34),
// the write grammar (REQ-60..REQ-63, REQ-66, REQ-69), and the stable-code
// table rows those families own (REQ-80..REQ-90).
//
// The discriminating property in this file is ORDER as much as identity:
// "refused at the CLI before any accessor runs" is only proved by an
// invocation whose accessor would ALSO fail — if the input refusal wins,
// the accessor never ran.

import (
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
)

// REQ-22: "Model selection MUST accept exactly one of --flow <id> or
// --model <path>."
// REQ-90: `flow-model-not-found` — "neither / both of `--flow`, `--model`,
// or the selection does not resolve" — `GroupUserEnv` / 2 — `param`.
// INPUT EDGE — all three arms of the one code.
func TestReq22And90_ModelSelectionArityAndNonResolutionShareOneCode(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)

	for _, verb := range flowVerbs {
		t.Run(verb+"/neither", func(t *testing.T) {
			ce := requireRefusal(t, "flow-model-not-found", 2,
				"flow", verb, "--as=json")
			if ce.Param == "" {
				t.Error("`flow-model-not-found` carries no `param`; the code " +
					"table fixes `param` as its carrier")
			}
		})
		t.Run(verb+"/both", func(t *testing.T) {
			requireRefusal(t, "flow-model-not-found", 2,
				"flow", verb, "--model", model, "--flow", "mvvflow", "--as=json")
		})
	}

	// A `--flow <id>` config discovery cannot resolve is the SAME code —
	// a refusal, not a missing feature (ASSUMPTION A-1).
	t.Run("flow-id-does-not-resolve", func(t *testing.T) {
		requireRefusal(t, "flow-model-not-found", 2,
			"flow", "next", "--flow", "no-such-flow-id-exists", "--as=json")
	})
}

// REQ-23: "The CLI MUST perform the model file I/O and hand RDR 0002's
// loader bytes plus a source id." — the loader performs no file I/O.
// INPUT EDGE — an unreadable path is the CLI's I/O failure, and it is NOT
// a load-category refusal: nothing reached the loader.
func TestReq23_ModelFileIOBelongsToTheCLINotTheLoader(t *testing.T) {
	missing := newFlowArtifact(t, "definitely-absent.toml")

	stdout, _, err := runCmd(t, "flow", "next", "--model", missing, "--as=json")
	if err == nil {
		t.Fatalf("a nonexistent --model path succeeded:\n%s", stdout)
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("refusal is not a structured CLIError: %v", err)
	}
	// The CLI owns the read, so the failure is a selection failure — not
	// `flow-model-invalid`, which is reserved for LOAD categories the
	// loader produced from bytes it actually saw (REQ-24).
	if ce.Code == "flow-model-invalid" {
		t.Errorf("an unreadable path reported %q; that code is reserved for "+
			"load categories the loader produced from bytes, and the loader "+
			"never saw these", ce.Code)
	}
	if !strings.HasPrefix(ce.Code, "flow-") {
		t.Errorf("code = %q; a model-selection failure carries a `flow-*` "+
			"code", ce.Code)
	}
}

// REQ-24: "Every load-time category MUST map to flow-model-invalid with one
// findings[] entry per category hit carrying its locator."
// REQ-91: `flow-model-invalid` — "any RDR 0002 load category, RDR 0003
// predicate category, or RDR 0008 `reserved_tag_key`" — `GroupUserEnv` / 2
// — `findings[]` (category, locator).
// DOMAIN EDGE
func TestReq24And91_LoadCategoriesRideOneCodeWithPerCategoryFindings(t *testing.T) {
	invalid := writeFlowModel(t, flowInvalidModel)

	stdout, _, err := runCmd(t, "flow", "next", "--model", invalid, "--as=json")
	if err == nil {
		t.Fatalf("an unloadable model succeeded:\n%s", stdout)
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("refusal is not a structured CLIError: %v", err)
	}
	if ce.Code != "flow-model-invalid" {
		t.Fatalf("code = %q; every load category maps to the ONE code %q",
			ce.Code, "flow-model-invalid")
	}
	if clierr.ExitCodeFor(err) != 2 {
		t.Errorf("exit = %d; want 2", clierr.ExitCodeFor(err))
	}

	findings := flowFailureFindings(t, stdout)
	if len(findings) == 0 {
		t.Fatal("`flow-model-invalid` carried NO findings; the code table " +
			"fixes `findings[]` (category, locator) as its carrier — one " +
			"entry per category hit")
	}
	for i, f := range findings {
		// `code` is the category slug: the discriminator the caller
		// branches on, since the ONE CLI code cannot carry it.
		if f.Code == "" {
			t.Errorf("finding[%d] carries no `code`; it is the category slug", i)
		}
		if f.Code == "flow-model-invalid" {
			t.Errorf("finding[%d].code = %q; the finding's code is the LOAD "+
				"CATEGORY slug, not the CLI code it rode in on", i, f.Code)
		}
		// `locator` is file:line — the field REQ-14/A-7 adds to Finding.
		if f.Locator == "" {
			t.Errorf("finding[%d] carries no `locator`; REQ-24 fixes its "+
				"content as file:line", i)
		}
	}
}

// REQ-25: "`revision` on every payload is the loaded model's own revision
// identity, carried through verbatim … The CLI never derives, hashes, or
// synthesizes it, so a model that declares none renders `revision` empty
// rather than a CLI-invented value."
// ADVERSARIAL — the discriminating oracle is the NEGATIVE: whatever the
// value is, it is never derived from the file path or its bytes.
func TestReq25_RevisionIsNeverCLIDerivedFromPathOrContent(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")

	data := flowData(t, requireSuccess(t, "flow", "next", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art), "--as=json"))

	raw, ok := data["revision"]
	if !ok {
		t.Fatalf("payload carries no `revision`; every payload carries it "+
			"(REQ-74..REQ-78). keys = %v", keysOf(data))
	}
	rev, ok := raw.(string)
	if !ok {
		t.Fatalf("`revision` = %#v; it is carried through verbatim as the "+
			"model's own identity", raw)
	}

	// Not the path, and not a path fragment: those would be CLI-derived.
	if rev != "" {
		if rev == model || strings.Contains(model, rev) && strings.Contains(rev, "/") {
			t.Errorf("`revision` = %q is derived from the --model PATH; the "+
				"CLI never derives it", rev)
		}
		// Not a hex digest of the bytes: a synthesized value.
		if isHexDigest(rev) {
			t.Errorf("`revision` = %q looks like a synthesized digest; the "+
				"CLI never hashes it", rev)
		}
	}
}

// isHexDigest reports whether s is a bare hex string of digest length — the
// shape a synthesized revision would take.
func isHexDigest(s string) bool {
	switch len(s) {
	case 32, 40, 64:
	default:
		return false
	}
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9', r >= 'a' && r <= 'f', r >= 'A' && r <= 'F':
		default:
			return false
		}
	}
	return true
}

// REQ-26: "--tag values MUST enter as observed context."
// REQ-114: "Caller-supplied `--tag` values are observed context and never
// satisfy an owned dependency."
// HAPPY PATH + DOMAIN EDGE
func TestReq26And114_TagEntersAsObservedAndNeverSatisfiesAnOwnedDependency(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	data := flowData(t, requireSuccess(t, "flow", "next", "--model", model,
		"--artifact", bind, "--tag", "profile=mid", "--as=json"))

	observed, ok := data["observed"].(map[string]any)
	if !ok {
		t.Fatalf("payload carries no `observed` object; --tag values enter "+
			"as OBSERVED context. keys = %v", keysOf(data))
	}
	if got := observed["profile"]; got != "mid" {
		t.Errorf("observed[profile] = %#v; want %q", got, "mid")
	}
	// And it did NOT land in `owned`: owned is assembled from readers only.
	if owned, ok := data["owned"].(map[string]any); ok {
		if _, leaked := owned["profile"]; leaked {
			t.Error("a --tag value appears in `owned`; caller-supplied tags " +
				"are observed context and NEVER satisfy an owned dependency")
		}
	}
}

// REQ-27: "A --tag naming an owned key or the reserved recognized key MUST
// be refused at the CLI before any accessor runs."
// REQ-82 / REQ-83: the codes are `flow-tag-reserved` and `flow-tag-owned`,
// both `GroupUserEnv` / 2 with carrier `param`.
// ADVERSARIAL — proved "before any accessor runs" by leaving the artifact
// role UNBOUND: were the accessor to run first, `flow-artifact-missing`
// would win instead.
func TestReq27And82And83_OwnedAndReservedTagsAreRefusedBeforeAnyAccessorRuns(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)

	t.Run("owned", func(t *testing.T) {
		ce := requireRefusal(t, "flow-tag-owned", 2,
			"flow", "resolve", "--model", model,
			"--outcome", "hold", "--tag", "status=draft", "--as=json")
		if ce.Param != "status" {
			t.Errorf("param = %q; want %q — the code table fixes `param` as "+
				"the carrier and the offending key is the subject",
				ce.Param, "status")
		}
	})

	t.Run("reserved-recognized", func(t *testing.T) {
		ce := requireRefusal(t, "flow-tag-reserved", 2,
			"flow", "next", "--model", model,
			"--tag", "recognized=hold", "--as=json")
		if ce.Param != "recognized" {
			t.Errorf("param = %q; want %q", ce.Param, "recognized")
		}
	})
}

// REQ-28 / REQ-81: "A duplicate tag name is `flow-tag-duplicate`" —
// `GroupUserEnv` / 2 — `param`.
// INPUT EDGE
func TestReq28And81_DuplicateTagNameIsItsOwnCode(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)

	ce := requireRefusal(t, "flow-tag-duplicate", 2,
		"flow", "next", "--model", model,
		"--tag", "profile=mid", "--tag", "profile=foundational", "--as=json")
	if ce.Param != "profile" {
		t.Errorf("param = %q; want %q", ce.Param, "profile")
	}
}

// REQ-29: "A set-kind tag's value is a JSON array literal
// (`'labels=[\"a\",\"b\"]'`) … a bare scalar for a set key, or an array for
// a scalar key, is `flow-tag-invalid`."
// REQ-80: `flow-tag-invalid` — "malformed, empty, or wrong-kind `--tag` /
// `--outcome` value" — `GroupUserEnv` / 2 — `param`.
// INPUT EDGE — both wrong-kind directions plus the empty value.
func TestReq29And80_WrongKindAndEmptyTagValuesAreFlowTagInvalid(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	// `--outcome` with an empty value (REQ-54).
	t.Run("empty-outcome", func(t *testing.T) {
		requireRefusal(t, "flow-tag-invalid", 2,
			"flow", "resolve", "--model", model, "--artifact", bind,
			"--outcome", "", "--as=json")
	})

	// An array handed to a SCALAR-kind observed key.
	t.Run("array-for-scalar-key", func(t *testing.T) {
		ce := requireRefusal(t, "flow-tag-invalid", 2,
			"flow", "next", "--model", model, "--artifact", bind,
			"--tag", `profile=["mid"]`, "--as=json")
		if ce.Param != "profile" {
			t.Errorf("param = %q; want %q", ce.Param, "profile")
		}
	})

	// A malformed `--tag` with no `=` at all.
	t.Run("no-equals", func(t *testing.T) {
		requireRefusal(t, "flow-tag-invalid", 2,
			"flow", "next", "--model", model, "--artifact", bind,
			"--tag", "profile", "--as=json")
	})
}

// REQ-31: "Nothing discovers artifacts ambiently; location comes only from
// explicit bindings."
// REQ-32 / REQ-88: "malformed `--artifact` binding" is
// `flow-artifact-invalid`, `GroupUserEnv` / 2, carrier `param`.
// ADVERSARIAL
func TestReq31And32And88_ArtifactLocationComesOnlyFromExplicitBindings(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)

	for _, malformed := range []string{
		"no-equals-sign",
		"=/tmp/orphaned-path",
		"role-with-no-path=",
	} {
		t.Run(malformed, func(t *testing.T) {
			requireRefusal(t, "flow-artifact-invalid", 2,
				"flow", "next", "--model", model,
				"--artifact", malformed, "--as=json")
		})
	}
}

// REQ-33 / REQ-89: "A role an invoked accessor needs that no binding
// supplies is `flow-artifact-missing`" — `GroupUserEnv` / 2 — "`param` =
// role".
// DOMAIN EDGE — the carrier is the ROLE, not the accessor id.
func TestReq33And89_MissingBindingForAnInvokedRoleCarriesTheRoleAsParam(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)

	// `resolve --outcome hold` needs owned `status`, served by `read.state`
	// on role `state`. With no binding at all, that role is missing.
	ce := requireRefusal(t, "flow-artifact-missing", 2,
		"flow", "resolve", "--model", model, "--outcome", "hold", "--as=json")
	if ce.Param != flowStateRole {
		t.Errorf("param = %q; want the ROLE %q — the code table fixes "+
			"`param` = role, not the accessor id", ce.Param, flowStateRole)
	}
}

// REQ-34: "`--tag` on `set-state` is context only and is never written." —
// restated fenced: "It MUST NOT treat --tag values as writes."
// ADVERSARIAL — the observable proof: a `--tag` naming an OBSERVED key on
// `set-state` must not appear in the read-back-confirmed owned values, and
// a `--tag` naming an owned key is refused outright (REQ-27).
func TestReq34_SetStateNeverTreatsTagValuesAsWrites(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	// An OWNED key through --tag is refused, so it cannot become a write.
	requireRefusal(t, "flow-tag-owned", 2,
		"flow", "set-state", "--model", model, "--artifact", bind,
		"--write", "status=final", "--tag", "labels=[\"plain\"]", "--as=json")

	// An OBSERVED key through --tag is accepted as context and does not
	// become a write: it is absent from the confirmed owned values.
	data := flowData(t, requireSuccess(t, "flow", "set-state",
		"--model", model, "--artifact", bind,
		"--write", "status=final", "--tag", "profile=mid", "--as=json"))
	if owned, ok := data["owned"].(map[string]any); ok {
		if _, leaked := owned["profile"]; leaked {
			t.Error("a `--tag` on set-state reached the owned tag values; " +
				"--tag is context only and is NEVER written")
		}
	}
	if writes, ok := data["writes"].(map[string]any); ok {
		if _, leaked := writes["profile"]; leaked {
			t.Error("a `--tag` on set-state was reported as a planned write")
		}
	}
}

// REQ-61: "Before any accessor runs it MUST refuse a --write or --clear key
// that is not a declared owned tag served by exactly one [write.<id>] whose
// keys list names it, a --write value malformed for its declared kind, the
// literal value <clear>, and a key given twice."
// REQ-86 / REQ-87: the unbound codes are `flow-write-unbound` and
// `flow-clear-unbound`, both `GroupUserEnv` / 2 with carrier `param`.
// ADVERSARIAL — again proved "before any accessor runs" by leaving the
// artifact role unbound, so an accessor-first order would refuse
// `flow-artifact-missing` instead.
func TestReq61And86And87_UnboundWriteAndClearKeysAreRefusedBeforeAccessors(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)

	t.Run("write-unbound", func(t *testing.T) {
		// `note` is a declared OWNED tag, but `[write.state].keys` does not
		// name it — no writer serves it.
		ce := requireRefusal(t, "flow-write-unbound", 2,
			"flow", "set-state", "--model", model,
			"--write", "note=x", "--as=json")
		if ce.Param != "note" {
			t.Errorf("param = %q; want %q", ce.Param, "note")
		}
	})

	t.Run("clear-unbound", func(t *testing.T) {
		ce := requireRefusal(t, "flow-clear-unbound", 2,
			"flow", "set-state", "--model", model,
			"--clear", "note", "--as=json")
		if ce.Param != "note" {
			t.Errorf("param = %q; want %q", ce.Param, "note")
		}
	})

	t.Run("undeclared-key", func(t *testing.T) {
		requireRefusal(t, "flow-write-unbound", 2,
			"flow", "set-state", "--model", model,
			"--write", "no-such-tag=x", "--as=json")
	})
}

// REQ-62: "The sentinel `<clear>` is unauthorable: `--write k=<clear>` is
// refused `flow-write-invalid` with the hint \"use `--clear`\"."
// REQ-84: `flow-write-invalid` — "malformed or wrong-kind `--write` value,
// or the literal `<clear>`" — `GroupUserEnv` / 2 — `param`.
// ADVERSARIAL
func TestReq62And84_ClearSentinelIsUnauthorableAndHintsAtTheFlag(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)

	ce := requireRefusal(t, "flow-write-invalid", 2,
		"flow", "set-state", "--model", model,
		"--write", "stale=<clear>", "--as=json")
	if ce.Param != "stale" {
		t.Errorf("param = %q; want %q", ce.Param, "stale")
	}
	if !strings.Contains(ce.Hint, "--clear") {
		t.Errorf("hint = %q; REQ-62 fixes the hint as \"use `--clear`\" — it "+
			"must name the flag that does authorize a clear", ce.Hint)
	}
}

// REQ-63 / REQ-85: "The same key under both flags, or twice under either, is
// `flow-write-duplicate`" — `GroupUserEnv` / 2 — `param`.
// INPUT EDGE — all three collision shapes.
func TestReq63And85_DuplicateKeyAcrossEitherFlagIsOneCode(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)

	cases := []struct {
		name string
		args []string
	}{
		{"write-twice", []string{"--write", "status=draft", "--write", "status=final"}},
		{"clear-twice", []string{"--clear", "stale", "--clear", "stale"}},
		{"write-and-clear", []string{"--write", "stale=x", "--clear", "stale"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"flow", "set-state", "--model", model},
				tc.args...)
			requireRefusal(t, "flow-write-duplicate", 2,
				append(args, "--as=json")...)
		})
	}
}

// REQ-69: "A carried plan artifact (`--plan <file|->`) is deferred as a
// seed, not part of this contract."
//
// The DEFERRAL IS DISCHARGED (kata c3xz): `--plan` now ships, on `set-state`
// alone. REQ-69 was a scoping note about what 0005 itself decided, not a
// permanent prohibition — 0005 names the flag, its `<file|->` argument, and
// the byte-identical copy-through that motivates it, and defers only the
// deciding. So this oracle is RETARGETED rather than deleted, onto the fence
// that survives the discharge: the plan carrier rides the one verb that can
// APPLY a plan, and the three verbs that cannot must not acquire it.
//
// That is the same shape `0023:C1`'s `--plan-only` fence takes on `resolve`,
// and it is asserted the same way — three negative arms plus the POSITIVE
// registration, because "absent from `next`" is vacuously true of a tree
// where the flag exists nowhere and would keep passing if the feature were
// reverted. Only the contrast discriminates.
//
// BOUNDARY
func TestReq69_ThePlanCarrierRidesSetStateAlone(t *testing.T) {
	// Vacuously true of an unregistered group, so the verbs must exist
	// before their flag sets mean anything.
	if names := flowGroupNames(t); len(names) != len(flowVerbs) {
		t.Fatalf("the `flow` group registers %v; all four verbs must exist "+
			"before the placement of `--plan` across them is assertable", names)
	}

	var checked bool
	for _, c := range NewRootCmd().Commands() {
		if c.Name() != "flow" {
			continue
		}
		if c.Flags().Lookup(planFlagName) != nil ||
			c.PersistentFlags().Lookup(planFlagName) != nil {
			t.Error("the `flow` group registers --plan of its own; the plan " +
				"carrier is VERB-LOCAL to `set-state`, and a group-level " +
				"registration would hand it to all four verbs")
		}
		for _, sub := range c.Commands() {
			if sub.Name() == "set-state" {
				checked = true
				// The POSITIVE half. Without it the negatives below hold
				// against a tree that never registered the flag at all.
				if sub.Flags().Lookup(planFlagName) == nil {
					t.Error("`flow set-state` does not register --plan; it is " +
						"the one verb that APPLIES a plan, and until it " +
						"accepts the flag the absences below are cobra's " +
						"unknown-flag path rather than this placement")
				}
				continue
			}
			if sub.Flags().Lookup(planFlagName) != nil {
				t.Errorf("`flow %s` registers --plan; only `set-state` "+
					"applies a plan — `next` and `read-state` write nothing, "+
					"and `resolve` PRODUCES the plan rather than consuming "+
					"one", sub.Name())
			}
		}
	}
	if !checked {
		t.Fatal("the `flow` group registers no `set-state` verb")
	}
}

// The `--plan` carrier must not disturb RDR 0023's `--plan-only` fence, and
// the fence must not have been widened to cover it. The two flags share a
// prefix and ride different verbs, which is exactly the arrangement in which
// a parser that did prefix matching would silently conflate them.
//
// ADVERSARIAL — the collision this pair makes newly reachable.
func TestReq69_ThePlanCarrierDoesNotDisturbThePlanOnlyFence(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	// `--plan-only` on `set-state` still fails the PARSE, even though
	// `set-state` now registers a flag `--plan-only` is a prefix-extension
	// of. pflag does no long-flag prefix matching, so `0023:C1`'s
	// non-registration fence holds structurally.
	_, _, err := runCmd(t, "flow", "set-state", "--model", model,
		"--artifact", bind, "--write", "status=draft",
		"--"+planOnlyFlag, "--as=json")
	if err == nil {
		t.Fatalf("`flow set-state --%s` SUCCEEDED; registering --%s must not "+
			"make --%s reachable on this verb — `0023:C1` places the "+
			"projection axis on `resolve` ALONE", planOnlyFlag, planFlagName,
			planOnlyFlag)
	}
	if code := clierr.ErrorCode(err); code != "command-error" {
		t.Errorf("code = %q; want %q — `--%s` is unregistered here and pflag "+
			"fails the parse before any RunE; a `flow-*` code would mean the "+
			"flag was accepted and refused downstream", code, "command-error",
			planOnlyFlag)
	}

	// And the converse: `--plan` is not reachable on `resolve`, which is
	// where `--plan-only` lives. A caller who typed the wrong one of the
	// pair is told so by the parse rather than having it silently absorbed.
	if _, _, err := runCmd(t, "flow", "resolve", "--model", model,
		"--artifact", bind, "--outcome", "hold",
		"--"+planFlagName, "-", "--as=json"); err == nil {
		t.Fatalf("`flow resolve --%s` SUCCEEDED; `resolve` PRODUCES a plan "+
			"and consumes none", planFlagName)
	}
}

// REQ-112: "No new third-party dependency. Cobra, `respond`, and `clierr`
// already exist; JSON array literals parse with the standard library. Any
// TOML dependency is RDR 0002's."
// BOUNDARY
func TestReq112_NoNewThirdPartyDependencyIsIntroduced(t *testing.T) {
	// The allowed direct-require set as it stands before this RDR.
	allowed := map[string]bool{
		"github.com/pelletier/go-toml/v2":      true,
		"github.com/spf13/cobra":               true,
		"github.com/spf13/pflag":               true,
		"github.com/inconshreveable/mousetrap": true,
	}

	// The clause is about what THIS RDR adds, so it only bites once this
	// RDR's surface exists: an unimplemented `flow` group trivially adds
	// no dependency.
	if names := flowGroupNames(t); len(names) != len(flowVerbs) {
		t.Fatalf("the `flow` group registers %v; the no-new-dependency "+
			"clause is about the surface this RDR lands, so that surface "+
			"must exist before the check discriminates", names)
	}

	for _, line := range strings.Split(readRepoFile(t, "go.mod"), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "github.com/") &&
			!strings.HasPrefix(line, "golang.org/") &&
			!strings.HasPrefix(line, "gopkg.in/") {
			continue
		}
		mod := strings.Fields(line)[0]
		if !allowed[mod] {
			t.Errorf("go.mod requires %q; this RDR introduces NO new "+
				"third-party dependency", mod)
		}
	}
}

// REQ-61: "Before any accessor runs it MUST refuse … a --write value
// malformed for its declared kind".
// REQ-84: `flow-write-invalid` — "malformed or wrong-kind `--write` value" —
// `GroupUserEnv` / 2 — `param`.
//
// "Malformed for its DECLARED KIND" reaches the declaration's domain, not
// only the set-vs-scalar shape. The declaration carries an enum's `domain`,
// an int's `min`/`max`, a bool's two literals, and a set's `elements`, and
// RDR 0002's loader already hard-refuses each of those from a RULE's
// authored literal — so a `--write` that may persist what a rule may not
// author would make the declaration advisory on exactly the path that
// writes.
//
// ADVERSARIAL — each case pairs the refusal with an ACCEPTED sibling under
// the same key, so the test cannot pass by refusing that key wholesale.
func TestReq61And84_AWriteValueIsHeldToItsDeclaredDomainAndBounds(t *testing.T) {
	model := writeFlowModel(t, flowDomainModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	set := func(t *testing.T, write string) (string, string, error) {
		t.Helper()
		return runCmd(t, "flow", "set-state", "--model", model,
			"--artifact", bind, "--write", write, "--as=json")
	}

	for _, tc := range []struct {
		name string
		// bad is refused `flow-write-invalid`; good is ACCEPTED. The pair
		// is what makes the case discriminate: a fix that refused the key
		// unconditionally would fail on `good`.
		bad, good string
		// member, when set, must appear in the refusal's message — a set
		// refusal that named only the key would leave a caller hunting
		// through the literal for which member offended.
		member string
	}{
		{
			name: "enum outside its domain",
			bad:  "status=notARealState", good: "status=final",
		},
		{
			name: "int above max",
			bad:  "iter=10", good: "iter=9",
		},
		{
			name: "int below min",
			bad:  "iter=-1", good: "iter=0",
		},
		{
			name: "int that is not a number at all",
			bad:  "iter=seven", good: "iter=7",
		},
		{
			name: "bool spelled some other way",
			bad:  "ready=yes", good: "ready=true",
		},
		{
			name: "set member outside its elements",
			bad:  `labels=["alpha","gamma"]`, good: `labels=["alpha","beta"]`,
			member: "gamma",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ce := requireRefusal(t, "flow-write-invalid", 2,
				"flow", "set-state", "--model", model,
				"--artifact", bind, "--write", tc.bad, "--as=json")

			key, _, _ := strings.Cut(tc.bad, "=")
			if ce.Param != key {
				t.Errorf("refusal param = %q; want %q — the code table fixes "+
					"`param` as this code's carrier, and it must name the "+
					"OFFENDING key", ce.Param, key)
			}
			if tc.member != "" && !strings.Contains(ce.Message, tc.member) {
				t.Errorf("the refusal message does not name the offending "+
					"member %q:\n%s\nA set refusal that names only the key "+
					"leaves the caller to find which member offended",
					tc.member, ce.Message)
			}

			// The conforming sibling under the SAME key must still be
			// accepted, so the refusal above is the DECLARATION biting.
			if _, _, err := set(t, tc.good); err != nil {
				t.Errorf("the conforming value %q was refused: %v\nHolding a "+
					"value to its declaration must not reject what the "+
					"declaration admits", tc.good, err)
			}
		})
	}

	// The UNCONSTRAINED control: a scalar declares no domain, so it keeps
	// taking anything. If this refused, the change would be a blanket
	// tightening of `--write` rather than declaration conformance.
	t.Run("an unconstrained scalar still takes any value", func(t *testing.T) {
		if _, _, err := set(t, "free=anything at all $%^"); err != nil {
			t.Errorf("a `scalar` write was refused %v; a scalar declares no "+
				"domain and there is nothing for it to violate", err)
		}
	})
}

// REQ-61 / REQ-84 — ORDERING. An unbound key and an out-of-domain value in
// the SAME argument: the caller must hear that the key is bound to no
// writer, because that is the mistake they made. Reporting a domain
// complaint about a key no writer serves sends them to fix the wrong thing.
// ADVERSARIAL — the two refusals are both reachable, so only their ORDER
// distinguishes a correct implementation.
func TestReq61_AnUnboundKeyIsReportedBeforeItsValuesDomain(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")

	// `profile` is an OBSERVED enum with domain ["mid","foundational"], so
	// `nonsense` is out of domain AND no writer's keys name the key. The
	// binding refusal must win.
	requireRefusal(t, "flow-write-unbound", 2,
		"flow", "set-state", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--write", "profile=nonsense", "--as=json")
}

// REQ-80: `flow-tag-invalid` — "malformed `--tag`" — `GroupUserEnv` / 2.
// REQ-26: `--tag` carries OBSERVED context.
//
// An observed key the model DOES declare is held to its domain, on the same
// grounds as `--write`. An observed key the model does NOT declare is not:
// `--tag` legitimately carries context the model never modelled, and
// refusing it would be a different decision under a different code.
// DOMAIN EDGE — the two halves are the whole point, so both are asserted.
func TestReq80And26_ADeclaredTagIsDomainCheckedAndAnUndeclaredOneIsNot(t *testing.T) {
	model := writeFlowModel(t, flowDomainModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	t.Run("a declared observed enum is held to its domain", func(t *testing.T) {
		ce := requireRefusal(t, "flow-tag-invalid", 2,
			"flow", "next", "--model", model, "--artifact", bind,
			"--tag", "hint=sideways", "--as=json")
		if ce.Param != "hint" {
			t.Errorf("refusal param = %q; want %q", ce.Param, "hint")
		}

		// …and a value INSIDE the domain still passes.
		requireSuccess(t, "flow", "next", "--model", model,
			"--artifact", bind, "--tag", "hint=high", "--as=json")
	})

	t.Run("an undeclared observed key still passes", func(t *testing.T) {
		// The model declares no `sidechannel` tag at all. There is no
		// declaration to violate, so the shape-only behaviour a caller
		// already relies on must survive: this is NOT the kata that
		// decides whether an unmodeled `--tag` key should be refused.
		requireSuccess(t, "flow", "next", "--model", model,
			"--artifact", bind,
			"--tag", "sidechannel=whatever-the-caller-likes", "--as=json")
	})
}
