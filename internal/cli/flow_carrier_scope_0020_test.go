package cli

// RDR 0020 — the seam's SCOPE, and the two prose obligations that carry it.
//
// C1 is a one-seam change with a widening direction, and most of what the
// record obligates outside the fence is that the widening stops where it
// says it stops: `--write` and `--clear` do not move, `canonicalValue`
// keeps its copy of the arm for the `--write` carrier, and the two prose
// surfaces that currently promise the RETIRED behaviour are corrected.
//
// The scope tests are behavioural wherever the boundary is observable
// through the CLI. Two obligations are not — a comment's content and a
// doc's content — and those are asserted by reading the shipped text, the
// same construction RDR 0024's suite uses for its own doc-amendment
// clauses.

import (
	"strings"
	"testing"
)

// REQ-27 / TD (out of scope): "`--write`/`--clear` are out of scope:
// `internal/cli/flow_state.go::parseWrites` proves a writer binding (and
// therefore a declaration) before its `canonicalValue` call, so the
// zero-decl arm is reachable from `parseTags` alone."
// REQ-28 / IP Step 1: "`canonicalValue` keeps its copy for the `--write`
// carrier, which enters by its own path."
// ADVERSARIAL — the widening must NOT leak onto the write path. The same
// undeclared key that `--tag` now carries must still refuse on `--write`
// and on `--clear`, for both a scalar and an array value: an implementation
// that made `canonicalValue` carrier-aware unconditionally would let one of
// these through.
func TestReq27And28_0020_TheWidenIngDoesNotReachWriteOrClear(t *testing.T) {
	// RDR 0005's shipped MVV fixture: `status` is a declared owned enum
	// with a bound writer, and nothing is named `extra`.
	model := writeFlowModel(t, flowMVVModel)
	bind := artifactBinding(flowStateRole,
		seedArtifact(t, model, "status=draft"))

	t.Run("write/undeclared-key-still-unbound", func(t *testing.T) {
		for _, value := range []string{"plain", `["a","b"]`} {
			requireRefusal(t, "flow-write-unbound", 2,
				"flow", "set-state", "--model", model,
				"--artifact", bind,
				"--write", carrierUndeclaredScalarKey+"="+value, "--as=json")
		}
	})

	t.Run("clear/undeclared-key-still-unbound", func(t *testing.T) {
		requireRefusal(t, "flow-clear-unbound", 2,
			"flow", "set-state", "--model", model,
			"--artifact", bind,
			"--clear", carrierUndeclaredScalarKey, "--as=json")
	})

	// `--write` on a DECLARED key keeps every conformance arm, including
	// the empty-value one whose copy `canonicalValue` retains (REQ-28).
	// The empty case is the load-bearing one: it is the arm that MOVES on
	// the `--tag` path, so its `--write` twin proves the copy stayed.
	t.Run("write/declared-key-keeps-its-conformance-arms", func(t *testing.T) {
		for _, value := range []string{"", `["draft"]`, "not-in-domain"} {
			ce := requireRefusal(t, "flow-write-invalid", 2,
				"flow", "set-state", "--model", model,
				"--artifact", bind,
				"--write", carrierOwnedKey+"="+value, "--as=json")
			if ce.Param != carrierOwnedKey {
				t.Errorf("param = %q; want %q", ce.Param, carrierOwnedKey)
			}
		}
	})
}

// REQ-25 / TD: "One seam moves: the admission path in
// `internal/cli/flow_input.go`. `parseTags` switches its guarded lookup to
// the two-value form and routes an undeclared key around `canonicalValue`
// (or passes the declaredness bit into it — implementation latitude,
// bounded by C1's refusal list: whichever shape is taken, the empty-value
// arm must still fire for a carrier and must still leave a declared set on
// its conformance message), so the kind/shape/domain arms run only under a
// real declaration."
// BOUNDARY — the clause names implementation latitude explicitly, so what
// is asserted is the BOUND, not the shape: the two invariants that hold
// under either implementation, checked as a pair in one test so neither can
// be satisfied at the other's expense.
func TestReq25_0020_EitherImplementationShapeMustHoldBothBounds(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	// Bound 1 — the empty-value arm still fires for a carrier.
	emptyArgs := carrierResolve(model, withCarriers(
		carrierUndeclaredScalarKey+"=")...)
	emptyStdout, _ := runCarrier(t, emptyArgs...)
	requireRefusal(t, "flow-tag-invalid", 2, emptyArgs...)
	if got := refusalMessage(t, emptyStdout); got != carrierFixtureF4Message {
		t.Errorf("bound 1: message = %q; want %q — whichever shape is "+
			"taken, the empty-value arm must still fire for a carrier",
			got, carrierFixtureF4Message)
	}

	// Bound 2 — a declared set is still left on its conformance message.
	setArgs := carrierResolve(model,
		carrierDeclaredScalar+"=free", carrierDeclaredSet+"=")
	setStdout, _ := runCarrier(t, setArgs...)
	requireRefusal(t, "flow-tag-invalid", 2, setArgs...)
	if got := refusalMessage(t, setStdout); got != carrierFixtureGMessage {
		t.Errorf("bound 2: message = %q; want %q — whichever shape is "+
			"taken, a declared set must still be left on its conformance "+
			"message", got, carrierFixtureGMessage)
	}

	// And the kind/shape/domain arms run ONLY under a real declaration:
	// the same array literal passes undeclared and refuses declared.
	if _, err := runCarrier(t, carrierResolve(model, withCarriers(
		carrierUndeclaredArrayKey+"="+carrierArrayValue)...)...); err != nil {
		t.Errorf("bound 3: an undeclared array literal was refused; the "+
			"kind/shape/domain arms run only under a real declaration: %v",
			err)
	}
}

// REQ-39 / IP Step 2: "Rewrite the guarded-lookup comment to cite this
// RDR's carrier contract instead of promising a different decision
// elsewhere; `ConformValue`'s zero-decl doc stays true as written."
// ADVERSARIAL — the comment at the seam is one of the three in-code
// authorities the record adjudicates between, and it currently promises the
// REJECTED arm ("Refusing an undeclared key is a different decision …
// this is not the place that makes it"). Leaving it in place would ship the
// same comment/code contradiction this RDR exists to close, in the opposite
// direction. `ConformValue`'s doc is NOT edited, which is asserted too so
// the correction cannot overshoot.
func TestReq39_0020_TheGuardedLookupCommentCitesTheCarrierContract(t *testing.T) {
	src := readRepoFile(t, "internal/cli/flow_input.go")

	// The retired promise: this RDR IS the place that makes the decision.
	retired := []string{
		"Refusing an undeclared key is a different decision",
		"this is not the place that makes it",
	}
	for _, phrase := range retired {
		if strings.Contains(src, phrase) {
			t.Errorf("`internal/cli/flow_input.go` still says %q; the "+
				"guarded-lookup comment must cite THIS RDR's carrier "+
				"contract instead of promising a different decision "+
				"elsewhere", phrase)
		}
	}

	// The comment must actually cite the contract, not merely drop the
	// stale sentence: an uncited seam is how the contradiction re-accretes.
	if !strings.Contains(src, "0020:C1") {
		t.Error("`internal/cli/flow_input.go` never cites `0020:C1`; the " +
			"rewritten comment cites this RDR's carrier contract by id")
	}

	// `ConformValue`'s zero-decl doc stays true AS WRITTEN — the
	// correction is scoped to the admission comment and must not overshoot
	// into the table package.
	conform := readRepoFile(t, "internal/table/load.go")
	if !strings.Contains(conform, "A zero TagDecl conforms everything") {
		t.Error("`internal/table/load.go`'s `ConformValue` doc no longer " +
			"says \"A zero TagDecl conforms everything\"; C1 says that doc " +
			"stays true as written and is NOT edited")
	}
}

// REQ-40 / IP Step 3: "the `docs/model-schema.md` line for authors
// composing a producer's whole tag output (Background's obligation)"
// REQ-41 / Background: "under arm (1) or (3), a `docs/model-schema.md` line
// for authors composing a producer's whole tag output."
// ASSUMPTION-1 / Q1 (proceeding under reading (b)): the named file does not
// exist at HEAD; the obligation lands in `docs/model-authoring.md`, whose
// §"Passing a set-valued tag" states "**An undeclared tag whose value is an
// array is refused.**" and "**every set-valued key the caller passes must
// be declared, even one no row guards on.**" — both falsified by C1.
// ADVERSARIAL — the doc is the surface the record's own Problem Statement
// names as the failure mode ("undiscoverable from the error"). Shipping C1
// while leaving the doc asserting the retired behaviour would recreate the
// defect in prose. Asserted as the absence of the falsified claims PLUS the
// presence of the composition guidance the obligation owes.
func TestReq40And41_0020_TheAuthoringDocNoLongerRefusesUndeclaredArrays(t *testing.T) {
	doc := readRepoFile(t, "docs/model-authoring.md")

	falsified := []struct {
		phrase string
		why    string
	}{
		{
			phrase: "An undeclared tag whose value is an array is refused",
			why: "C1 makes an undeclared array a pure carrier admitted " +
				"verbatim (REQ-1, REQ-31)",
		},
		{
			phrase: "every set-valued key the caller passes must be declared",
			why: "C1 removes exactly that requirement — composing a " +
				"producer's whole tag output is the motivating case",
		},
		{
			phrase: "so an array literal there returns",
			why:    "the zero declaration no longer kind-checks anything",
		},
	}
	for _, f := range falsified {
		if strings.Contains(doc, f.phrase) {
			t.Errorf("`docs/model-authoring.md` still says %q; %s",
				f.phrase, f.why)
		}
	}

	// The obligation is a LINE FOR AUTHORS composing a producer's whole tag
	// output, so deleting the falsified paragraph does not discharge it:
	// the doc must say what happens now.
	if !strings.Contains(doc, "undeclared") {
		t.Error("`docs/model-authoring.md` never discusses an undeclared " +
			"tag key; the obligation is a LINE for authors composing a " +
			"producer's whole tag output, not a deletion")
	}
}

// REQ-59 / `0020:G-cross-cutting` (scoping): "This RDR makes no
// byte-identical output, content-addressed identity, or replay-stable hash
// claim, so the determinism checklist does not apply. The two
// \"byte-identical\" uses in C1 are narrower: they assert message-string
// equality between two refusal sites across the hoist"
// DOMAIN EDGE — a scoping clause, so what it obligates is that the SUITE
// does not overreach: "byte-identical" here means message-string equality,
// and a hash/golden-digest test built on it would be asserting a claim the
// record explicitly disclaims. Asserted by construction — the two refusal
// sites the hoist spans are compared as STRINGS, and their equality is what
// the clause actually says.
func TestReq59_0020_ByteIdenticalMeansMessageStringEqualityNotAHash(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	// The two refusal sites the hoist spans: an UNDECLARED key (fixture F4,
	// which after the hoist comes from the admission path) and a DECLARED
	// SCALAR (fixture E, which after the hoist may come from either site).
	// The clause's whole content is that these two strings agree modulo the
	// key name.
	undeclaredStdout, _ := runCarrier(t, carrierResolve(model, withCarriers(
		carrierUndeclaredScalarKey+"=")...)...)
	declaredStdout, _ := runCarrier(t, carrierResolve(model,
		carrierDeclaredScalar+"=",
		carrierDeclaredSet+"="+carrierDeclaredSetValue)...)

	const template = "the tag `%s` was given an empty value"
	undeclaredMsg := refusalMessage(t, undeclaredStdout)
	declaredMsg := refusalMessage(t, declaredStdout)

	if undeclaredMsg != carrierFixtureF4Message {
		t.Errorf("undeclared-site message = %q; want %q (fixture F4)",
			undeclaredMsg, carrierFixtureF4Message)
	}
	if declaredMsg != carrierFixtureEMessage {
		t.Errorf("declared-scalar-site message = %q; want %q (fixture E)",
			declaredMsg, carrierFixtureEMessage)
	}

	// The two differ ONLY in the key name — which is what "byte-identical
	// between two refusal sites across the hoist" asserts, and all it
	// asserts. No digest is computed anywhere in this suite.
	undeclaredShape := strings.Replace(undeclaredMsg,
		carrierUndeclaredScalarKey, "%s", 1)
	declaredShape := strings.Replace(declaredMsg,
		carrierDeclaredScalar, "%s", 1)
	if undeclaredShape != declaredShape || undeclaredShape != template {
		t.Errorf("the two refusal sites do not share one message template\n"+
			"  undeclared: %q\n  declared:   %q\n  want both:  %q",
			undeclaredShape, declaredShape, template)
	}
}
