package cli

// RDR 0023 — `0023:C2`: the declared ECHO/PLAN partition, its reflective
// completeness oracle, the always-keep core, the pointer-conversion scope,
// and C1's observed-reader-execution differential.
//
// The one thing C2 is most insistent about shapes every oracle here:
//
//	"an oracle that derives the ECHO set by observing what the projection
//	 drops is tautological — it restates the implementation and cannot
//	 fail."
//
// So the assignment lists this file asserts against
// (`echoGroup0023` / `planGroup0023`, in flow_fixtures_0023_test.go) are
// AUTHORED from the census, and the DECLARED table the implementation owes
// is asserted to AGREE with them — never consulted to learn what they are.
//
// Likewise, C1 rejects by name the vacuous form of the reader assertion:
// "Comparing invokedReaders(model, outcome) across the two runs is NOT an
// admissible form of this assertion … green by construction and cannot
// fail, which makes it worse than absent." The reader oracle below is
// written on OBSERVED EXECUTION — the run's own side effects — and never
// re-calls the planning function.

import (
	"encoding/json"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// --- reflective partition completeness (`0023:C2`, `0023:S3`) -----------

// REQ-60: "Per JDR 0002 §D1's enforcement rule, this verb's reflective
// oracle MUST assert every field of the resolve success payload is assigned
// to exactly one group, so an unassigned new field is a test failure, not a
// silent default — the projection MUST NOT be implemented as a bare
// omit-list whose complement is \"whatever else exists\"." —
// [JDR-0002:D1]
// REQ-61: "the ECHO group is the model reference, the observed tags, the
// assembled owned view and the invoked reader identities, and the requested
// outcome."
// REQ-62: "The PLAN group is rule identity, gate results, authored answers
// and their interpretations, planned next/writes/clear, the escape
// disposition, and revision."
// REQ-63: the per-field assignment fixed by the census.
// REQ-68: the oracle asserts three things: "every struct field has exactly
// one entry; the entry set and the field set are equal (neither a field
// without an entry nor an entry without a field); and the keys a projected
// run actually emits equal the declaration's `plan` side."
// REQ-69: "A new field with no entry then fails at the first assertion
// rather than defaulting into either width."
// REQ-79/REQ-80 (S2/S3 controls): "add a field to the payload without a C2
// side; S2 and S3 must both go red" / "add a field, declare it `plan`, but
// omit it from the projection — S3 must go red … while S2 stays green"
// ADVERSARIAL — the completeness oracle. Its whole force is BIDIRECTIONAL
// equality: a one-way containment check passes on a struct that grew a
// field, which is the exact drift `0023:RM` names.
func TestReq60And61And62And63And68And69_EveryPayloadFieldIsAssignedToExactlyOneC2Group(t *testing.T) {
	// The AUTHORED assignment, from the census (`0023:CEN`). Authored is
	// the point: a table derived from the projection's behaviour would
	// agree with any implementation.
	assignment := map[string]string{
		"model": "echo", "observed": "echo", "owned": "echo",
		"readers": "echo", "outcome": "echo",
		"revision": "plan", "rule": "plan", "gates": "plan", "emit": "plan",
		// RDR 0024 `0024:C4` appends `dispositions`; JDR 0002 §D1 assigns
		// it PLAN — a disposition token is the INTERPRETATION of an
		// authored answer, which this group names, never an echo of the
		// request. Registering it here is `0024:PH3`'s "whichever record
		// lands second" obligation on the 0023-first leg.
		"dispositions": "plan",
		"next":         "plan", "writes": "plan", "clear": "plan",
		"escaped": "plan", "escape_class": "plan",
	}

	// (1) every struct field has exactly one entry, and (2) the entry set
	// and the field set are EQUAL in both directions.
	fields := payloadWireKeys()
	for _, key := range fields {
		side, held := assignment[key]
		if !held {
			t.Errorf("`resolvePayload` declares the wire key %q, which the "+
				"C2 assignment does not carry; an UNASSIGNED new field is a "+
				"test failure, not a silent default — the projection is not "+
				"a bare omit-list whose complement is \"whatever else "+
				"exists\". Assign it a side in the RDR, then here", key)
			continue
		}
		if side != "echo" && side != "plan" {
			t.Errorf("wire key %q is assigned %q; the doctrine admits "+
				"exactly two groups", key, side)
		}
	}
	for key := range assignment {
		if !slices.Contains(fields, key) {
			t.Errorf("the C2 assignment carries an entry for %q, which "+
				"`resolvePayload` does not declare; the entry set and the "+
				"field set are equal in BOTH directions — neither a field "+
				"without an entry nor an entry without a field.\n"+
				"struct wire keys = %v", key, fields)
		}
	}

	// The two authored group lists must partition the assignment, so the
	// rest of the suite's use of `echoGroup0023`/`planGroup0023` is pinned
	// to the census rather than drifting independently.
	for _, key := range echoGroup0023 {
		if assignment[key] != "echo" {
			t.Errorf("%q is not assigned ECHO; `0023:C2` names the model "+
				"reference, the observed tags, the assembled owned view, "+
				"the invoked reader identities and the requested outcome",
				key)
		}
	}
	for _, key := range planGroup0023 {
		if assignment[key] != "plan" {
			t.Errorf("%q is not assigned PLAN; `0023:C2` names rule "+
				"identity, gate results, authored answers and their "+
				"interpretations, planned next/writes/clear, the escape "+
				"disposition, and revision", key)
		}
	}

	// (3) the keys a PROJECTED RUN actually emits equal the declaration's
	// `plan` side — modulo the producer presence rules, which REQ-25
	// requires reading off the same run's default output.
	def := rawFieldsOf(t, requireSuccess(t, append(pricingCall(t), "--as=json")...))
	proj := rawFieldsOf(t, requireSuccess(t,
		append(pricingCall(t, "--"+planOnlyFlag), "--as=json")...))

	var wantPlan []string
	for key, side := range assignment {
		if side != "plan" {
			continue
		}
		if _, carriedByDefault := def[key]; !carriedByDefault {
			// A presence-rule field the default mode omits is omitted
			// under the flag for the same reason (REQ-25).
			continue
		}
		wantPlan = append(wantPlan, key)
	}
	slices.Sort(wantPlan)

	gotPlan := sortedKeys(proj)
	if !slices.Equal(gotPlan, wantPlan) {
		t.Errorf("the keys a projected run emits = %v;\nthe declaration's "+
			"`plan` side (carried by the same request's default) = %v\n"+
			"A field declared `plan` but omitted from the projection fails "+
			"HERE while the key-set oracle stays green — that is the "+
			"discriminating control S3 owns and S2 does not",
			gotPlan, wantPlan)
	}
}

// REQ-65: "The oracle reflects the payload struct against an assignment
// source that is INDEPENDENT of the projection code … an oracle that
// derives the ECHO set by observing what the projection drops is
// tautological — it restates the implementation and cannot fail."
// REQ-66: "the assignment is DECLARED (a standalone table keyed by field
// name, one entry per field, carrying `echo` or `plan` — the carrier
// constrained below)"
// REQ-70: "The carrier is CONSTRAINED, not free: it MUST be a standalone
// table keyed by field name, in its own declaration site, NOT a per-field
// marker on `resolvePayload`."
// REQ-71: "Independence has to be structural — a separate site a projection
// edit does not open — or S3 asserts only that the implementation agrees
// with itself."
// A-4: "The declared assignment table lives in non-test production code,
// keyed by the Go FIELD name … one exported-or-package-level
// `map[string]group` (or equivalent) in `internal/cli`, in a file the
// projection edit does not open."
// ADVERSARIAL — the carrier constraint is what the clause exists for, and
// it is assertable two ways at once: the declaration must EXIST as a
// standalone production site, and `resolvePayload` must carry NO per-field
// group marker (which is the carrier the clause forbids, reachable in the
// very commit that changes the projection).
func TestReq65And66And70And71_TheAssignmentIsDeclaredInAStandaloneSiteNotAsAStructMarker(t *testing.T) {
	// --- the FORBIDDEN carrier: a per-field marker on the struct ------
	//
	// A2 rewrites the echo fields' types and tags on `resolvePayload`, so a
	// marker carried THERE is edited in the same commit that changes the
	// projection — the tautological oracle by the carrier the clause would
	// otherwise permit.
	rt := resolvePayloadType()
	for i := range rt.NumField() {
		f := rt.Field(i)
		for _, forbidden := range []string{
			"group", "partition", "echo", "plan", "project", "projection",
		} {
			if _, marked := f.Tag.Lookup(forbidden); marked {
				t.Errorf("`resolvePayload.%s` carries a `%s` struct tag; the "+
					"carrier MUST be a standalone table in its own "+
					"declaration site, NOT a per-field marker on the "+
					"payload struct. A2 retypes these very fields, so a "+
					"marker here moves the assertion and its subject "+
					"together — independence has to be STRUCTURAL",
					f.Name, forbidden)
			}
		}
	}

	// --- the REQUIRED carrier: a standalone production declaration -----
	//
	// Reached by reading `internal/cli`'s non-test sources rather than by
	// referencing a symbol, so this file compiles against a tree that does
	// not yet carry it and fails on its ABSENCE rather than taking the
	// whole package down with a build error.
	root := repoRootFor(t)
	dir := root + "/internal/cli"
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}

	var carrierFile string
	var projectionFile string
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") ||
			strings.HasSuffix(name, "_test.go") {
			continue
		}
		src, rerr := os.ReadFile(dir + "/" + name)
		if rerr != nil {
			t.Fatalf("read %s: %v", name, rerr)
		}
		text := string(src)

		// The declaration site: a standalone table keyed by field name
		// carrying the two group names, one entry per field. Recognised by
		// carrying an assignment for every echo AND every plan field.
		declares := true
		for _, key := range append(slices.Clone(echoGroup0023),
			planGroup0023...) {
			// Either the wire key or the Go field name may key the table
			// (A-4 reads "field name" as the Go identifier, with the JSON
			// tag reachable from it).
			if !strings.Contains(text, `"`+key+`"`) &&
				!strings.Contains(text, `"`+goFieldNameFor(key)+`"`) {
				declares = false
				break
			}
		}
		if declares && strings.Contains(text, "echo") && strings.Contains(text, "plan") {
			carrierFile = name
		}
		// The projection site: whatever file reads the flag.
		if strings.Contains(text, planOnlyFlag) {
			projectionFile = name
		}
	}

	if carrierFile == "" {
		t.Fatalf("no non-test file in internal/cli declares a standalone "+
			"ECHO/PLAN assignment table covering every payload field; C2 "+
			"requires the assignment DECLARED — one entry per field, "+
			"carrying `echo` or `plan` — and the projection implemented FROM "+
			"that declaration. Without it S3 asserts only that the "+
			"implementation agrees with itself.\nfields needing an entry: %v",
			payloadWireKeys())
	}
	if projectionFile != "" && carrierFile == projectionFile {
		t.Errorf("the assignment table and the projection both live in %s; "+
			"independence has to be STRUCTURAL — a separate site a "+
			"projection edit does not open. Sharing the file is the "+
			"tautological oracle reached by the carrier the clause forbids",
			carrierFile)
	}
}

// goFieldNameFor maps a wire key to the Go field name `resolvePayload`
// declares it under, so the carrier probe accepts either keying (A-4).
func goFieldNameFor(wireKey string) string {
	rt := resolvePayloadType()
	for i := range rt.NumField() {
		f := rt.Field(i)
		if wireKeyForField(f) == wireKey {
			return f.Name
		}
	}
	return wireKey
}

// REQ-72: "Group membership is the declaration table's (C2), never a Go
// type's." — an implementation or oracle that treats "pointer + omitempty"
// as the mark of a projectable field is forbidden, because `escape_class`
// is a plain `string` with `,omitempty` and is always-keep core.
// REQ-88: "the mechanism applies to the five ECHO fields only, each `T` →
// `*T` with `T` unchanged from the census's Go-type column (`string`,
// `map[string]string`, `[]string`) — no other field on `resolvePayload`
// becomes pointer-valued."
// REQ-89: "`escape_class` stays a plain `string` whose `,omitempty` drops
// the key on `\"\"`. Its presence rule (`0005:A-3`) is therefore realised by
// a mechanism DISTINCT from the projection mechanism, and the two must not
// be conflated" — [0005-carried]
// A-10: "`*T` + `,omitempty` on exactly `model`, `observed`, `owned`,
// `readers`, `outcome`; the JSON key strings themselves are unchanged."
// BOUNDARY — the pointer conversion's exact scope. Its boundary is
// `escape_class`: a mechanism that swept "pointer + omitempty" would take
// an always-keep core field with it.
func TestReq72And88And89_ThePointerConversionCoversTheFiveEchoFieldsAndNothingElse(t *testing.T) {
	rt := resolvePayloadType()

	// The census's Go-type column: `T` UNCHANGED under the conversion.
	wantElem := map[string]reflect.Kind{
		"model":    reflect.String,
		"observed": reflect.Map,
		"owned":    reflect.Map,
		"readers":  reflect.Slice,
		"outcome":  reflect.String,
	}

	for i := range rt.NumField() {
		f := rt.Field(i)
		key := wireKeyForField(f)
		isEcho := slices.Contains(echoGroup0023, key)
		isPtr := f.Type.Kind() == reflect.Pointer

		if isEcho {
			if !isPtr {
				t.Errorf("echo field %q (`%s`) is %s, not a pointer; the "+
					"mechanism is pointer-valued echo fields with "+
					"`omitempty` — a bare non-pointer `omitempty` DROPS the "+
					"empty `{}`/`[]` and would change DEFAULT-mode bytes, "+
					"the hazard `0023:A2`'s spike reproduced",
					key, f.Name, f.Type)
				continue
			}
			if got, want := f.Type.Elem().Kind(), wantElem[key]; got != want {
				t.Errorf("echo field %q is *%s; the census pins `T` "+
					"UNCHANGED at %s — the conversion is `T` → `*T`, never a "+
					"retyping", key, f.Type.Elem(), want)
			}
			tag := f.Tag.Get("json")
			if !strings.Contains(tag, ",omitempty") {
				t.Errorf("echo field %q carries the json tag %q; the "+
					"mechanism is `*T` + `,omitempty` — without it a nilled "+
					"pointer renders `null` rather than dropping the key, "+
					"and `absent means absent` fails", key, tag)
			}
			if got := wireKeyForField(f); got != key {
				t.Errorf("echo field's wire key changed to %q; the JSON key "+
					"STRINGS are unchanged by the conversion, which is what "+
					"keeps the three shipped wire oracles green", got)
			}
			continue
		}

		// Every other field: NOT pointer-valued. `escape_class` is the
		// boundary case the clause names — a plain `string` with
		// `,omitempty`, always-keep core, whose presence rule is a
		// DIFFERENT mechanism.
		if isPtr {
			t.Errorf("field %q (`%s`) became pointer-valued; the conversion "+
				"applies to the FIVE ECHO FIELDS ONLY. Group membership is "+
				"the declaration table's, NEVER a Go type's — treating "+
				"\"pointer + omitempty\" as the mark of a projectable field "+
				"would sweep away `escape_class`, which the partition "+
				"requires CARRIED", key, f.Name)
		}
	}

	// `escape_class` explicitly: plain string, `,omitempty` kept, so its
	// `0005:A-3` presence rule is untouched and unconflated.
	ec, ok := rt.FieldByName("EscapeClass")
	if !ok {
		t.Fatal("`resolvePayload` declares no `EscapeClass` field")
	}
	if ec.Type.Kind() != reflect.String {
		t.Errorf("`EscapeClass` is %s; it stays a plain `string` whose "+
			"`,omitempty` drops the key on \"\" — a mechanism DISTINCT from "+
			"the projection mechanism, and the two must not be conflated",
			ec.Type)
	}
	if !strings.Contains(ec.Tag.Get("json"), ",omitempty") {
		t.Errorf("`EscapeClass`'s json tag is %q; `0005:A-3`'s presence rule "+
			"is realised by that `,omitempty` and this RDR does not move it",
			ec.Tag.Get("json"))
	}
}

// REQ-64: "The set is FOURTEEN fields — `reflect.TypeOf(resolvePayload{})
// .NumField()` is 14, the count `internal/cli/decision_table_0010_test.go`
// already pins" — [0010-carried]
// REQ-90: "The `resolvePayload` struct keeps exactly its current field
// count … so A2's pointer conversion changes field TYPES only." —
// [0010-carried]
// REQ-91: "Adding a field here would falsify A3's \"no predecessor oracle
// moves\" claim, so a field addition is out of scope for this RDR by
// construction, not by preference." — negative REQ
// REQ-135: "no 0005/0010/0011 test moves (A3)" — [0005/0010/0011-carried]
// DOMAIN EDGE — the negative that keeps A3 true: the conversion changes
// TYPES, never the count, so the shipped 0010 oracle passes verbatim.
func TestReq64And90And91And135_TheConversionChangesFieldTypesNeverTheFieldCount(t *testing.T) {
	rt := resolvePayloadType()

	// 15 on THIS tree. `0023:JC1`/REQ-87 record the ordering tolerance with
	// cli/0024 and fix whose cost the move is: "if that RDR lands first its
	// own additive append moves this count, and that cost is 0024's, not
	// this projection's". 0023 landed FIRST here, so it is `0024:C4`'s
	// append of `dispositions` — one field, immediately after `Emit` — that
	// took the count from 14 to 15.
	//
	// The subject of THIS assertion is unchanged and still fails on its own
	// terms: the pointer conversion changes field TYPES only, so the count
	// moves exactly once per additive RDR and never as a side effect of the
	// conversion.
	if rt.NumField() != 15 {
		t.Errorf("`resolvePayload` declares %d fields; the pointer "+
			"conversion changes field TYPES only, and the one field "+
			"addition on this tree is `0024:C4`'s `dispositions` append "+
			"(14 -> 15). A further addition is out of scope for BOTH "+
			"records",
			rt.NumField())
	}

	// The shipped 0010 oracle's other two probes — `FieldByName` existence
	// and `Emit`/`Gates` adjacency — must also pass verbatim, since
	// `0023:A3`'s type-inspection sweep rests on the conversion preserving
	// count, name and position.
	gates, gok := rt.FieldByName("Gates")
	emit, eok := rt.FieldByName("Emit")
	if !gok || !eok {
		t.Fatalf("`resolvePayload` no longer declares Gates (%v) and Emit "+
			"(%v); the conversion preserves NAMES", gok, eok)
	}
	if emit.Index[0] != gates.Index[0]+1 {
		t.Errorf("`Emit` is at index %d and `Gates` at %d; `0010:C4` fixes "+
			"`emit` IMMEDIATELY after `Gates` and the conversion preserves "+
			"POSITION", emit.Index[0], gates.Index[0])
	}
}

// --- always-keep core (`0023:C2`) ---------------------------------------

// REQ-73: "This verb's ALWAYS-KEEP core (JDR 0002 §D1) is rule, escaped,
// escape_class, revision: if the boolean axis ever generalizes to an enum
// or a field list, no mode may PROJECT them away" — [JDR-0002:D1]
// REQ-74: "Always-keep is projection-invariance, not unconditional
// presence: a core field whose producer already has a presence rule
// (escape_class) appears under the flag exactly when it appears by default"
// REQ-76: "revision rides the PLAN side deliberately … Projecting it costs
// 14 bytes and keeps the wire shape stable for the day 0002 admits the
// key." — carried under the flag even though constant-empty today.
// REQ-75: "This RDR therefore does NOT mint a separate always-keep oracle
// (it would today assert exactly what S2 asserts)" — negative REQ: the
// enforcement is S2's literal, and this test asserts the core's membership
// in that literal rather than minting a second oracle over modes.
// DOMAIN EDGE — the safety invariant: a projected plan can never launder a
// rescued plan into an ordinary one, and can never detach a plan from the
// revision that produced it.
func TestReq73And74And76_TheAlwaysKeepCoreIsProjectionInvariant(t *testing.T) {
	escapeModel := writeFlowModel(t, dtEscapeModel0010)

	for _, arm := range []struct {
		name string
		args []string
	}{
		{"ordinary-plan", pricingCall(t)},
		{"rescued-plan", resolveArgs(escapeModel, "", "decide",
			"--tag", "a=y", "--tag", "b=q")},
	} {
		t.Run(arm.name, func(t *testing.T) {
			def := rawFieldsOf(t,
				requireSuccess(t, append(slices.Clone(arm.args), "--as=json")...))
			proj := rawFieldsOf(t, requireSuccess(t,
				append(slices.Clone(arm.args), "--"+planOnlyFlag, "--as=json")...))

			for _, key := range alwaysKeepCore0023 {
				defRaw, defHeld := def[key]
				projRaw, projHeld := proj[key]

				// PROJECTION-INVARIANCE, not unconditional presence: the
				// core field appears under the flag EXACTLY when it appears
				// by default.
				if defHeld != projHeld {
					t.Errorf("always-keep core field %q is %s by default and "+
						"%s under the flag; always-keep is "+
						"PROJECTION-INVARIANCE — a core field whose producer "+
						"has a presence rule appears under the flag exactly "+
						"when it appears by default",
						key, presence(defHeld), presence(projHeld))
					continue
				}
				if defHeld && string(defRaw) != string(projRaw) {
					t.Errorf("always-keep core field %q = %s by default and "+
						"%s under the flag", key, defRaw, projRaw)
				}
			}

			// The laundering guard actually rests on `escaped`, which is
			// UNCONDITIONALLY present: a projected payload can never read as
			// an ordinary plan when it was a rescue.
			raw, held := proj["escaped"]
			if !held {
				t.Fatalf("the projected payload carries no `escaped`; that " +
					"is the field the laundering guard rests on, and it is " +
					"unconditionally present in both widths")
			}
			if arm.name == "rescued-plan" && strings.TrimSpace(string(raw)) != "true" {
				t.Errorf("`escaped` = %s on the rescued arm; a projected "+
					"payload must never launder a rescued plan into an "+
					"ordinary one", raw)
			}

			// `revision` rides the PLAN side and is carried under the flag
			// even though it is constant-empty today — the forward-binding
			// half of the core.
			rev, held := proj["revision"]
			if !held {
				t.Errorf("the projected payload carries no `revision`; it " +
					"rides the PLAN side deliberately (loader-produced, " +
					"never restated from the request) and its always-keep " +
					"clause is FORWARD-BINDING — the model-identity slot " +
					"cannot be dropped at the moment 0002 gives it a value")
			} else if string(rev) != `""` {
				// Not a defect if 0002 later admits the key; stated so a
				// reader knows what today's witness is.
				t.Logf("`revision` = %s (constant-empty today per `0023:C2`)", rev)
			}
		})
	}
}

// REQ-77: "Under JDR 0002 §D1, 0024's dispositions, if that RDR lands, is a
// PLAN-group field … and its never-omitted clause is untouched by this
// projection." — [JDR-0002:D1]
// REQ-87: "Ordering tolerance with cli/0024 confirmed (A4): either RDR may
// land first; the second lands with `dispositions` already/newly in the
// plan group and no contract in either moves."
// DOMAIN EDGE — the joint-decision seam, asserted conditionally so it is
// live whichever RDR lands first and vacuous-by-design where the field does
// not exist yet.
func TestReq77And87_IfDispositionsExistsItIsCarriedUnderTheFlag(t *testing.T) {
	def := rawFieldsOf(t, requireSuccess(t, append(pricingCall(t), "--as=json")...))
	if _, held := def["dispositions"]; !held {
		t.Skip("`dispositions` is not on this tree; cli/0024 has not landed. " +
			"The ordering tolerance is confirmed either way (`0023:A4`) — " +
			"this assertion becomes live when that RDR lands, in either order")
	}

	proj := rawFieldsOf(t, requireSuccess(t,
		append(pricingCall(t, "--"+planOnlyFlag), "--as=json")...))
	projRaw, held := proj["dispositions"]
	if !held {
		t.Fatalf("`dispositions` is absent under --%s; it is a PLAN-group "+
			"field — each token is assigned by the [emit] declaration to the "+
			"selected row's authored emit value, derived from plan-group "+
			"inputs ONLY, never from `observed` or `owned` — so this "+
			"projection strips no derivation input and its never-omitted "+
			"clause is untouched", planOnlyFlag)
	}
	if string(projRaw) != string(def["dispositions"]) {
		t.Errorf("`dispositions` = %s by default and %s under the flag; "+
			"neither contract moves in either landing order",
			def["dispositions"], projRaw)
	}
}

// --- observed reader execution (`0023:C1`, `0023:S1`) -------------------

// REQ-35: the differential asserts "an identical invoked-reader set"
// REQ-37: "A later change that skips work whose only consumer is a
// projected-away field breaches this clause." — the reader pass is the
// reachable surface.
// REQ-40: "The measurand is therefore OBSERVED READER EXECUTION, not a
// recomputed derivation. Comparing invokedReaders(model, outcome) across
// the two runs is NOT an admissible form of this assertion"
// REQ-41: "The oracle MUST instead observe which readers ACTUALLY RAN in
// each run — recording execution at the reader invocation site (a counting
// or recording seam around the reader pass, the run's own side effects,
// never a re-call of the planning function) — and assert the two observed
// sets are equal."
// REQ-42: "The invocation site the seam wraps is the per-reader `exec.Read`
// call inside `internal/cli/flow_exec.go::(flowRequest).runReaders`, the one
// pass that yields `readers` and `owned` together"
// REQ-43: "The DEFAULT-mode run of the same request supplies the expected
// set. The discriminating property is that a projected run which skips the
// reader pass MUST turn this assertion red"
// REQ-44: reading the set "from the projected payload" is forbidden and
// unwritable: "`readers` is an ECHO field this flag projects away, and the
// shipped helper reading it fails hard on absence".
// REQ-46: "because one call yields both fields, no implementation may
// satisfy the projection by suppressing `owned` and `readers`
// independently."
// REQ-78 (S1's discriminating control): "SKIP the reader pass under the
// flag and confirm the reader assertion goes red."
// ADVERSARIAL — the clause's whole point. The measurand is the RUN'S OWN
// SIDE EFFECTS: a reader that actually executes leaves observable traces
// (its established owned state reaches the selected plan, and its accessor
// touches the artifact), and a projected run that skipped the pass produces
// a DIFFERENT DECISION. Recomputing `invokedReaders` is rejected by name;
// reading `readers` off the projected payload is unwritable.
func TestReq35And37And40And41And42And43And44And46_TheInvokedReaderSetIsObservedExecutionNotARecomputation(t *testing.T) {
	// A reader-backed model whose owned state DECIDES the rule: the reader
	// pass is not decorative here, so skipping it is observable in the
	// plan itself and not merely in a field the flag projects away.
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft", "stale=x")
	bind := artifactBinding(flowStateRole, art)
	args := resolveArgs(model, bind, "advance")

	def := rawFieldsOf(t, requireSuccess(t, append(slices.Clone(args), "--as=json")...))
	proj := rawFieldsOf(t, requireSuccess(t,
		append(slices.Clone(args), "--"+planOnlyFlag, "--as=json")...))

	// Control 1: the DEFAULT run supplies the expected set, and it is
	// NON-EMPTY. Without this the equality below is satisfiable by two runs
	// that both ran nothing — the vacuous form the clause rejects.
	readers, held := def["readers"]
	if !held || strings.TrimSpace(string(readers)) == "[]" {
		t.Fatalf("the default run reports no invoked reader: %s\nThe "+
			"expected set comes from the DEFAULT-mode run of the same "+
			"request, and an empty expected set makes the differential "+
			"unfailable — two runs that both ran nothing compare equal",
			readers)
	}

	// Control 2: `readers` really is projected away, so the assertion
	// CANNOT be written from the projected payload. This is REQ-44's
	// "unwritable", asserted rather than asserted-about.
	if _, held := proj["readers"]; held {
		t.Fatalf("`readers` is PRESENT under the flag; it is an ECHO field " +
			"this flag projects away, and the fact that it is unreadable " +
			"there is why the measurand has to be observed execution")
	}

	// The assertion, on OBSERVED EXECUTION. A reader that ran established
	// owned state; that state is what the kernel selected over, so it
	// reaches the PLAN group — `rule`, `next`, `writes`, `clear` — which
	// the flag carries. A projected run that skipped the reader pass
	// cannot reproduce these, because it would have resolved over an empty
	// owned view and selected differently (or refused).
	//
	// This is the run's own side effects, never a re-call of
	// `invokedReaders`: that function takes `(*table.Model, outcome)`,
	// neither of which the flag touches, so comparing it across the two
	// runs returns equal sets on EVERY implementation — including one that
	// never runs a reader at all.
	for _, key := range []string{"rule", "next", "writes", "clear"} {
		defRaw, defHeld := def[key]
		projRaw, projHeld := proj[key]
		if !defHeld || !projHeld {
			t.Errorf("plan-group field %q is %s by default and %s under the "+
				"flag", key, presence(defHeld), presence(projHeld))
			continue
		}
		if string(defRaw) != string(projRaw) {
			t.Errorf("plan-group field %q = %s by default and %s under the "+
				"flag; this field is DECIDED OVER the owned state the "+
				"reader pass established, so a divergence here is the "+
				"observable signature of a projected run that skipped work "+
				"whose only consumer is a projected-away field",
				key, defRaw, projRaw)
		}
	}

	// The reader pass yields `readers` and `owned` TOGETHER from one call
	// (`0023:CEN`), so no implementation may satisfy the projection by
	// suppressing them independently: both are absent under the flag or
	// the projection is not the partition it claims to be.
	_, ownedHeld := proj["owned"]
	_, readersHeld := proj["readers"]
	if ownedHeld != readersHeld {
		t.Errorf("under the flag `owned` is %s and `readers` is %s; ONE "+
			"call yields both fields, so no implementation may satisfy the "+
			"projection by suppressing them independently",
			presence(ownedHeld), presence(readersHeld))
	}

	// And the artifact is untouched by either run: `resolve` performs reads
	// only, in both widths, so a projected run that "optimised away" the
	// read pass would be the only way these could diverge.
	//
	// The read-back binds the `orphan` role as well as `state`. That is not
	// a loosening: `flow read-state` is `0005:REQ-37`'s deliberate exception
	// and runs EVERY declared reader, because a diagnostic read has no
	// candidate set to narrow by — so the role `resolve` correctly leaves
	// unbound is one `read-state` requires. The 0005 MVV binds both roles to
	// the same artifact for exactly this reason.
	readBack := func() string {
		t.Helper()
		return requireSuccess(t, "flow", "read-state", "--model", model,
			"--artifact", bind,
			"--artifact", artifactBinding(flowOrphanRole, art), "--as=json")
	}
	before := readBack()
	requireSuccess(t, append(slices.Clone(args), "--"+planOnlyFlag, "--as=json")...)
	if after := readBack(); after != before {
		t.Errorf("the artifact's observed state changed across a projected "+
			"`flow resolve`; --%s is REPORT-ONLY\nbefore:\n%s\nafter:\n%s",
			planOnlyFlag, before, after)
	}
}

// REQ-41/REQ-42, the seam's POSITION: "a counting or recording seam around
// the reader pass … The invocation site the seam wraps is the per-reader
// `exec.Read` call inside `internal/cli/flow_exec.go::(flowRequest)
// .runReaders` … the seam's FORM is left to the implementer, but its
// position is not — a seam placed at the name-computation
// (`invokedReaders`) rather than the execution records intent, not
// execution, and is the vacuous form this clause rejects."
// A-3: "a flag-blind recording hook (or an observable side effect of the
// real reader run) that both the default and the projected run drive
// identically, with the oracle comparing the two recordings."
// ADVERSARIAL — the seam must be at the EXECUTION, and it must be
// FLAG-BLIND (REQ-49 forbids the flag being readable there). Asserted
// structurally over the production source: a recording seam that sits in
// `invokedReaders` records intent, and one that reads the flag violates the
// single-lexical-site rule.
func TestReq41And42_TheReaderRecordingSeamSitsAtTheExecutionAndIsFlagBlind(t *testing.T) {
	src, err := os.ReadFile(repoRootFor(t) + "/internal/cli/flow_exec.go")
	if err != nil {
		t.Fatalf("read flow_exec.go: %v", err)
	}
	text := string(src)

	// The flag is NEVER readable in the reader pass: C1 fixes exactly ONE
	// lexical site for it, the projection branch, and A-3 reads the seam as
	// necessarily flag-blind.
	if strings.Contains(text, planOnlyFlag) {
		t.Errorf("`flow_exec.go` mentions %q; the flag is read at exactly "+
			"ONE lexical site — the projection branch — and is never "+
			"consulted upstream of payload assembly. A flag-conditional "+
			"reader seam would make the very execution it records "+
			"flag-dependent", planOnlyFlag)
	}

	// The seam is at the EXECUTION, not the name computation. Split the
	// source at the two function bodies and require the recording to sit
	// with `exec.Read`.
	execAt := strings.Index(text, "exec.Read(")
	if execAt < 0 {
		t.Fatal("`flow_exec.go` no longer calls `exec.Read`; that per-reader " +
			"call inside `runReaders` is the invocation site the seam wraps, " +
			"and C1 fixes its POSITION even while leaving the seam's form to " +
			"the implementer")
	}
	invokedAt := strings.Index(text, "func invokedReaders(")
	if invokedAt < 0 {
		t.Fatal("`flow_exec.go` no longer declares `invokedReaders`")
	}

	// A recording seam exists and is nearer the execution than the
	// name-computation. Recognised by a recorder hook the tests can drive;
	// its FORM is the implementer's, so this looks for any of the shapes
	// the clause admits ("counting or recording").
	var seamAt = -1
	for _, marker := range []string{
		"readerRecorder", "recordReaderExecution", "observeReaderExecution",
		"readerExecutionHook", "recordRead(",
	} {
		if i := strings.Index(text, marker); i >= 0 {
			seamAt = i
			break
		}
	}
	if seamAt < 0 {
		t.Fatalf("no reader-execution recording seam is present in " +
			"`flow_exec.go`; C1 requires the oracle to OBSERVE which readers " +
			"actually ran — \"a counting or recording seam around the reader " +
			"pass, the run's own side effects, never a re-call of the " +
			"planning function\". Without it the differential's reader " +
			"assertion is either unwritable (from the projected payload, " +
			"where `readers` is absent) or vacuous (by recomputing " +
			"`invokedReaders`, whose two arguments the flag does not touch)")
	}

	// Position: the seam belongs with the `exec.Read` pass, not with the
	// name computation. `invokedReaders` returns names; recording there
	// records INTENT, which is the vacuous form.
	invokedEnd := strings.Index(text[invokedAt:], "\n}\n")
	if invokedEnd > 0 && seamAt > invokedAt && seamAt < invokedAt+invokedEnd {
		t.Errorf("the recording seam sits inside `invokedReaders`; a seam " +
			"placed at the NAME-COMPUTATION rather than the execution " +
			"records intent, not execution, and is the vacuous form this " +
			"clause rejects. Its position is the per-reader `exec.Read` call " +
			"inside `runReaders`")
	}
}

// --- A9's Phase-1 empty-container oracle ---------------------------------

// REQ-84: "Two CONSTRUCTION obligations carry into the build … and A9's
// empty-container default-mode assertion (Phase 1, BEFORE the pointer
// conversion lands, riding S2's absence control)."
// REQ-85: "a default-mode assertion over a request whose containers are
// empty (no `--tag`, no owned keys, no readers), asserting the rendered
// bytes carry `{}`/`[]` and no `null`"
// REQ-86: "the five pointers must be non-nil whenever a success payload is
// emitted, and that is now an implementation invariant stated here rather
// than an implicit one."
// REQ-123: "`observed`/`owned` render as `{}` when empty in default mode,
// which bare-map `omitempty` would wrongly drop (A2)."
// INPUT EDGE — the nil arm. This is the arm the three shipped wire oracles
// do NOT cover (they run populated or decision-table shapes), and it is the
// one where the pointer conversion can silently change DEFAULT-mode bytes
// for callers who never opted in. It rides S2's absence control and is
// deliberately written to be meaningful BEFORE the conversion lands.
func TestReq84And85And86And123_EmptyEchoContainersRenderAsBracesAndBracketsNeverNull(t *testing.T) {
	// The pricing decision table declares ZERO owned tags and ZERO
	// accessors (that is what makes it a decision table), so on any
	// successful call `owned` is `{}` and `readers` is `[]` — two of the
	// three echo containers empty, which is A9's arm. The observed tags are
	// supplied because the model's guards need them to select a row: a
	// refusal never reaches the success payload this test is about.
	//
	// The DEFAULT-mode half below runs and must PASS today: it pins
	// existing, already-correct behaviour that this RDR explicitly
	// preserves, and it exists to be assertable BEFORE the pointer
	// conversion lands. That ordering is the whole point — once the
	// conversion is in, a nil container renders `null` under `omitempty`
	// and silently changes default-mode bytes for callers who never opted
	// in, and S1-S5 cannot see it because they compare the new build
	// against itself.
	stdout := requireSuccess(t, append(pricingCall(t), "--as=json")...)

	line := emittedLine(t, stdout)
	fields := rawFieldsOf(t, stdout)

	// Control: the containers really ARE empty on this request. Without it
	// the assertion below passes on a populated shape and never reaches the
	// nil arm — which is precisely why the three shipped wire oracles
	// (`decision_table_0010_test.go::wireKeyOrder`, `TestReq41_...`,
	// `flow_resolve_0005_test.go::TestReq76_...`) do NOT discharge this
	// one: each runs a populated or decision-table shape and none is
	// written as an empty-container assertion.
	for _, empty := range []struct {
		key  string
		want string
	}{
		{"owned", "{}"},
		{"readers", "[]"},
	} {
		raw, held := fields[empty.key]
		if !held {
			t.Fatalf("DEFAULT mode carries no %q; the five echo keys ride "+
				"the default width unconditionally, and a missing one here "+
				"is exactly the silent default-mode byte change the pointer "+
				"conversion risks — a bare non-pointer `omitempty` DROPS an "+
				"empty `{}`/`[]`", empty.key)
		}
		if got := strings.TrimSpace(string(raw)); got != empty.want {
			t.Errorf("%q renders as %s in DEFAULT mode on an "+
				"empty-container request; it must render %s. A nil container "+
				"reached by any writer renders `null` under `omitempty`; a "+
				"bare non-pointer `omitempty` drops the key entirely. Both "+
				"silently change default-mode bytes for callers who never "+
				"opted in", empty.key, got, empty.want)
		}
	}

	// `observed` is POPULATED on this call, and that is deliberate: it
	// witnesses that the same conversion keeps a non-empty container
	// rendering unchanged, so the empty assertions above are not passing
	// because the whole group went missing.
	if raw, held := fields["observed"]; !held ||
		strings.TrimSpace(string(raw)) == "{}" {
		t.Errorf("`observed` = %s; this call supplies two --tag values, so "+
			"a populated container is the control that keeps the empty "+
			"assertions above honest", raw)
	}

	// No `null` anywhere on the default-mode line. This is the assertion
	// S2's absence control rides, whose negative control renders one echo
	// key as `null`.
	if strings.Contains(line, ":null") {
		t.Errorf("the DEFAULT-mode emitted line carries a `null` value:\n"+
			"%s\nEvery container writer feeding the echo group returns a "+
			"non-nil value, so a pointer to it renders `{}`/`[]` and never "+
			"`null`. Under A2's conversion that incidental property becomes "+
			"LOAD-BEARING: the five pointers must be non-nil whenever a "+
			"success payload is emitted", line)
	}

	// The same request under the flag: ABSENCE, not a `null` or `{}`
	// placeholder. Empty and absent are different claims, and this is where
	// the two meet — an empty container that becomes absent rather than
	// `null` is the mechanism working correctly.
	t.Run("projected", func(t *testing.T) {
		proj := requireSuccess(t,
			append(pricingCall(t, "--"+planOnlyFlag), "--as=json")...)
		projFields := rawFieldsOf(t, proj)
		for _, key := range echoGroup0023 {
			if raw, held := projFields[key]; held {
				t.Errorf("projected key %q is present as %s on the "+
					"empty-container request; empty and absent are DIFFERENT "+
					"claims and the projection makes the key ABSENT", key, raw)
			}
		}
	})
}

// --- negative REQs the build must not violate ---------------------------

// REQ-125: "*Version marker* — none is minted, deliberately. The wire shape
// is identified by the flag the caller passed, not by an in-payload
// marker" — negative REQ.
// REQ-114: "Every C1 MUST maps to at least one oracle below or an MVV step;
// A1's byte table is recorded evidence, not a test assertion." — negative
// REQ: no oracle asserts A1's byte percentages.
// REQ-128: "No new encoding surface. Tag literals are refused rather than
// coerced by the shipped parser … the sole rewrite is `canonicalSet`'s
// re-encode of an already-valid set literal, which is value-preserving." —
// [0005-carried], negative REQ.
// DOMAIN EDGE — the three negatives, asserted as absences on the wire and
// on the tag-parsing surface.
func TestReq114And125And128_NoVersionMarkerNoNewEncodingSurfaceNoByteThresholdOnTheWire(t *testing.T) {
	def := rawFieldsOf(t, requireSuccess(t, append(pricingCall(t), "--as=json")...))
	proj := rawFieldsOf(t, requireSuccess(t,
		append(pricingCall(t, "--"+planOnlyFlag), "--as=json")...))

	// No in-payload version/width marker is minted, in either width.
	for _, marker := range []string{
		"version", "width", "projected", "plan_only", "planOnly", "mode",
		"shape", "projection",
	} {
		for width, fields := range map[string]map[string]json.RawMessage{
			"default": def, "projected": proj,
		} {
			if _, held := fields[marker]; held {
				t.Errorf("the %s payload carries a %q key; NO version marker "+
					"is minted, deliberately — the wire shape is identified "+
					"by the FLAG THE CALLER PASSED, not by an in-payload "+
					"marker. `revision` is the plan-side identity slot and is "+
					"contracted-but-vacant", width, marker)
			}
		}
	}

	// No new encoding surface: a set literal refused-not-coerced still
	// refuses under the flag, with the same code and exit.
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	for _, name := range []string{"default", "projected"} {
		args := resolveArgs(model, bind, "hold", "--tag", "status=draft", "--as=json")
		if name == "projected" {
			args = resolveArgs(model, bind, "hold",
				"--tag", "status=draft", "--"+planOnlyFlag, "--as=json")
		}
		t.Run(name, func(t *testing.T) {
			// `status` is OWNED, so supplying it as an observed tag is
			// refused by the shipped parser — refused, never coerced.
			requireRefusal(t, "flow-tag-owned", 2, args...)
		})
	}
}
