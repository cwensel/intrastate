package cli

// RDR 0005 — model selection and the input grammar every `flow` verb
// shares: `--tag`, `--artifact`, `--write`, and `--clear`.
//
// Everything in this file runs BEFORE any accessor does. That is a
// contract obligation, not an optimisation: REQ-27 and REQ-61 both say
// "before any accessor runs", and the Phase 1 suite proves the ordering by
// leaving an artifact role unbound on invocations whose input is also
// wrong — if an accessor ran first, the refusal would name the wrong
// subject.

import (
	"encoding/json"
	"os"
	"slices"
	"strings"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
	"github.com/newcoinc/intrastate/internal/cli/flowbind"
	"github.com/newcoinc/intrastate/internal/table"
	"github.com/spf13/cobra"
)

// --- the stable code table (`0005:FM`) -----------------------------------
//
// The spellings are normative and the group fixes the exit. They are
// declared as constants in one place so a producer cannot spell one
// slightly differently at a second site.

const (
	codeTagInvalid      = "flow-tag-invalid"
	codeTagDuplicate    = "flow-tag-duplicate"
	codeTagReserved     = "flow-tag-reserved"
	codeTagOwned        = "flow-tag-owned"
	codeWriteInvalid    = "flow-write-invalid"
	codeWriteDuplicate  = "flow-write-duplicate"
	codeWriteUnbound    = "flow-write-unbound"
	codeClearUnbound    = "flow-clear-unbound"
	codeArtifactInvalid = "flow-artifact-invalid"
	codeArtifactMissing = "flow-artifact-missing"
	codeModelNotFound   = "flow-model-not-found"
	codeModelInvalid    = "flow-model-invalid"

	codeGateDenied        = "flow-gate-denied"
	codeGateIndeterminate = "flow-gate-indeterminate"

	codeAccessorTimeout  = "flow-accessor-timeout"
	codeAccessorFailed   = "flow-accessor-failed"
	codeReadIncomplete   = "flow-read-incomplete"
	codeAccessorUnknown  = "flow-accessor-unknown"
	codeCapabilityMismat = "flow-accessor-capability-mismatch"
	codeWriteNonOwned    = "flow-write-non-owned"

	codeReadBackMismatch   = "flow-write-readback-mismatch"
	codeReadBackIncomplete = "flow-write-readback-incomplete"
	codeReadBackTimeout    = "flow-write-readback-timeout"
)

// userErr builds a `GroupUserEnv` refusal: the request or the model is
// wrong, or the model said no. Exit 2 (REQ-19).
func userErr(code, param, message string) *clierr.CLIError {
	return &clierr.CLIError{
		Code: code, Param: param, Message: message,
		Group: clierr.GroupUserEnv,
	}
}

// envErr builds a `GroupEnvUnavailable` refusal: the environment could not
// be consulted and the same request may be re-run unchanged. Exit 3
// (REQ-19, REQ-20).
func envErr(code, param, message string) *clierr.CLIError {
	return &clierr.CLIError{
		Code: code, Param: param, Message: message,
		Group: clierr.GroupEnvUnavailable,
	}
}

// --- model selection -----------------------------------------------------

// selectModel resolves `--model <path>` or `--flow <id>` to a loaded model.
//
// The three arms of `flow-model-not-found` (REQ-22, REQ-90): neither flag,
// both flags, and a selection that does not resolve. Carrying the selection
// ARITY error under the same code as non-resolution is RDR-owned (CA A6) —
// from the caller's side both mean "this invocation names no one model".
//
// The CLI performs the file I/O and hands the loader BYTES plus a source id
// (REQ-23): `table.Load` does no file I/O, so an unreadable path is this
// command's failure and never a load category.
func selectModel(cmd *cobra.Command) (*table.Model, string, *clierr.CLIError) {
	path, _ := cmd.Flags().GetString("model")
	flowID, _ := cmd.Flags().GetString("flow")

	switch {
	case path != "" && flowID != "":
		return nil, "", userErr(codeModelNotFound, "model",
			"--model and --flow are mutually exclusive; give exactly one")
	case path == "" && flowID == "":
		return nil, "", userErr(codeModelNotFound, "model",
			"model selection requires exactly one of --model <path> or --flow <id>")
	case path == "":
		// `--flow <id>` config discovery. `--model <path>` ships first
		// (`0005:EIA`); an id this build cannot resolve is a REFUSAL under
		// the same code, not a missing feature (A-1).
		return nil, "", userErr(codeModelNotFound, "flow",
			"no model is registered for the flow id "+flowID+
				"; give --model <path> instead")
	}

	src, err := os.ReadFile(path)
	if err != nil {
		// The CLI owns the read. This is a selection failure — the loader
		// never saw these bytes, so it is emphatically not a load category
		// (REQ-23).
		return nil, "", userErr(codeModelNotFound, "model",
			"the model at "+path+" could not be read: "+err.Error())
	}

	model, err := table.Load(src, path)
	if err != nil {
		return nil, "", loadFailure(path, err)
	}
	return model, path, nil
}

// loadFailure maps a loader refusal onto `flow-model-invalid` with one
// findings entry per category hit (REQ-24, REQ-91).
//
// Every load-time category rides ONE CLI code, so the category slug is what
// a caller branches on and it travels as the finding's `code`. `locator` is
// the finding's file:line — the position the loader attributes the defect
// to — which is why REQ-14 needed a field the shipped record did not carry.
func loadFailure(path string, err error) *clierr.CLIError {
	category, ok := table.CategoryOf(err)
	if !ok {
		// A loader error carrying no category still refuses under this
		// code: it is a defect in the selected model, and inventing a
		// second code for it would break the one-code rule.
		category = "unknown"
	}
	return &clierr.CLIError{
		Code:    codeModelInvalid,
		Message: "the selected model could not be loaded",
		Group:   clierr.GroupUserEnv,
		Findings: []clierr.Finding{{
			Code:    string(category),
			Message: err.Error(),
			Locator: path + ":1",
		}},
	}
}

// --- artifact bindings ---------------------------------------------------

// parseArtifacts parses `--artifact role=path` bindings. Nothing discovers
// artifacts ambiently; location comes ONLY from these (REQ-31, `0004:C3`).
func parseArtifacts(cmd *cobra.Command) (map[string]string, *clierr.CLIError) {
	raw, _ := cmd.Flags().GetStringArray("artifact")
	out := make(map[string]string, len(raw))
	for _, binding := range raw {
		role, path, ok := strings.Cut(binding, "=")
		if !ok || role == "" || path == "" {
			return nil, userErr(codeArtifactInvalid, "artifact",
				"--artifact takes role=path; got "+binding)
		}
		out[role] = path
	}
	return out, nil
}

// --- observed tags -------------------------------------------------------

// parseTags parses `--tag name=value` into observed context.
//
// Three refusals precede any accessor (REQ-26, REQ-27, REQ-28):
// `flow-tag-reserved` for the reserved `recognized` key, `flow-tag-owned`
// for a key the model declares owned, and `flow-tag-duplicate` for a
// repeated name. The first two protect provenance: owned state is assembled
// from readers, never from the caller, so a `--tag` that could satisfy an
// owned dependency would let a caller assert state the artifact does not
// hold (REQ-114).
func parseTags(cmd *cobra.Command, m *table.Model) ([]resolveTag, *clierr.CLIError) {
	raw, _ := cmd.Flags().GetStringArray("tag")
	owned := flowbind.OwnedTags(m)
	sets := flowbind.SetKeys(m)

	var out []resolveTag
	seen := map[string]bool{}
	for _, entry := range raw {
		key, value, ok := strings.Cut(entry, "=")
		if !ok || key == "" {
			return nil, userErr(codeTagInvalid, "tag",
				"--tag takes name=value; got "+entry)
		}
		switch {
		case key == table.RecognizedTagKey:
			return nil, userErr(codeTagReserved, key,
				"the tag key `"+key+"` is reserved; the recognized outcome "+
					"enters through --outcome")
		case slices.Contains(owned, key):
			return nil, userErr(codeTagOwned, key,
				"the tag key `"+key+"` is owned; owned state is read from "+
					"the declared read accessors, never supplied by the caller")
		case seen[key]:
			return nil, userErr(codeTagDuplicate, key,
				"the tag `"+key+"` is given more than once")
		}
		seen[key] = true

		canonical, ce := canonicalValue(key, value, slices.Contains(sets, key), "tag")
		if ce != nil {
			return nil, ce
		}
		out = append(out, resolveTag{Key: key, Value: canonical})
	}
	return out, nil
}

// resolveTag is one parsed key/value pair. It mirrors `resolve.Tag` without
// importing the kernel into the parsing layer.
type resolveTag struct {
	Key   string
	Value string
}

// canonicalValue validates a value against its declared kind and renders it
// in the form the seam carries (REQ-29, REQ-70).
//
// A set value crosses as JDR 0001 §D13's canonical JSON array — members
// sorted, duplicate-free, compact — and is RE-CANONICALISED here rather
// than passed through: the caller may legitimately hand back a plan's own
// literal, but may equally hand an unsorted or duplicated one, and read-back
// equality is byte equality over the canonical form. A scalar handed an
// array, or a set handed a bare scalar, is the wrong-kind refusal.
func canonicalValue(key, value string, isSet bool, flag string) (string, *clierr.CLIError) {
	looksArray := strings.HasPrefix(strings.TrimSpace(value), "[")

	if !isSet {
		if looksArray {
			return "", userErr(codeInvalidFor(flag), key,
				"the tag `"+key+"` is not set-valued; got the array literal "+value)
		}
		if value == "" {
			return "", userErr(codeInvalidFor(flag), key,
				"the tag `"+key+"` was given an empty value")
		}
		return value, nil
	}

	if !looksArray {
		return "", userErr(codeInvalidFor(flag), key,
			"the tag `"+key+"` is set-valued and takes a JSON array literal; got "+value)
	}
	var members []string
	if err := json.Unmarshal([]byte(value), &members); err != nil {
		return "", userErr(codeInvalidFor(flag), key,
			"the set literal for `"+key+"` is not a JSON array of strings: "+err.Error())
	}
	return canonicalSet(members), nil
}

// codeInvalidFor selects the wrong-kind code for the flag the value arrived
// on. `--tag` and `--outcome` take `flow-tag-invalid`; `--write` takes
// `flow-write-invalid` (REQ-80, REQ-84).
func codeInvalidFor(flag string) string {
	if flag == "write" {
		return codeWriteInvalid
	}
	return codeTagInvalid
}

// canonicalSet renders members in JDR 0001 §D13's canonical form: sorted,
// duplicate-free, compact, HTML escaping disabled.
//
// This is THE encoder for a set literal at the CLI. Every site that emits
// or compares one calls it, which is what makes plan-to-request
// copy-through and read-back equality byte equality (REQ-71).
func canonicalSet(members []string) string {
	canonical := slices.Compact(slices.Sorted(slices.Values(members)))
	if canonical == nil {
		canonical = []string{}
	}
	var buf strings.Builder
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(canonical); err != nil {
		// Encoding a []string cannot fail.
		return "[]"
	}
	return strings.TrimRight(buf.String(), "\n")
}
