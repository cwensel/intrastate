package table_test

// RDR 0002 sections D–E: accessor tables and their provenance-scoped
// bindings, and the root, stop set, and terminal contexts.

import (
	"reflect"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// REQ-19: "Every entry carries `role`, `path`, `keys` (a non-empty list of
// declared tag keys), and `timeout` (a Go duration string). **Each of the
// four is required, and any of them absent, empty, or ill-formed is a
// `malformed accessor declaration`** — an absent or empty `role` or
// `path`, an absent or empty `keys` list, a `keys` member naming an
// undeclared tag (`unknown tag`), an absent `timeout`, a `timeout` that
// does not parse as a Go duration, or one that parses non-positive."
// ADVERSARIAL
//
// The clause is stated per field because "carries" alone left three of
// the four with no refusal, so every field arm gets its own case. The
// `keys`-member arm carries a DIFFERENT category (unknown tag) and is
// asserted as such.
func TestReq19_AccessorDeclarationRequiresAllFourFields(t *testing.T) {
	base := string(readFixture(t, rdrFixture))
	const gateBlock = `[gate.rdr-lock]
role = "rdr"
path = "rdr.lock"
keys = ["status"]
timeout = "2s"`

	cases := []struct {
		name    string
		replace string
		want    table.Category
	}{
		{"absent role", `[gate.rdr-lock]
path = "rdr.lock"
keys = ["status"]
timeout = "2s"`, table.CatMalformedAccessorDeclaration},
		{"empty role", `[gate.rdr-lock]
role = ""
path = "rdr.lock"
keys = ["status"]
timeout = "2s"`, table.CatMalformedAccessorDeclaration},
		{"absent path", `[gate.rdr-lock]
role = "rdr"
keys = ["status"]
timeout = "2s"`, table.CatMalformedAccessorDeclaration},
		{"empty path", `[gate.rdr-lock]
role = "rdr"
path = ""
keys = ["status"]
timeout = "2s"`, table.CatMalformedAccessorDeclaration},
		{"absent keys", `[gate.rdr-lock]
role = "rdr"
path = "rdr.lock"
timeout = "2s"`, table.CatMalformedAccessorDeclaration},
		{"empty keys list", `[gate.rdr-lock]
role = "rdr"
path = "rdr.lock"
keys = []
timeout = "2s"`, table.CatMalformedAccessorDeclaration},
		{"keys member naming an undeclared tag", `[gate.rdr-lock]
role = "rdr"
path = "rdr.lock"
keys = ["nosuchtag"]
timeout = "2s"`, table.CatUnknownTag},
		{"absent timeout", `[gate.rdr-lock]
role = "rdr"
path = "rdr.lock"
keys = ["status"]`, table.CatMalformedAccessorDeclaration},
		{"timeout that does not parse as a Go duration", `[gate.rdr-lock]
role = "rdr"
path = "rdr.lock"
keys = ["status"]
timeout = "soon"`, table.CatMalformedAccessorDeclaration},
		{"timeout that parses non-positive", `[gate.rdr-lock]
role = "rdr"
path = "rdr.lock"
keys = ["status"]
timeout = "0s"`, table.CatMalformedAccessorDeclaration},
		{"negative timeout", `[gate.rdr-lock]
role = "rdr"
path = "rdr.lock"
keys = ["status"]
timeout = "-2s"`, table.CatMalformedAccessorDeclaration},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := strings.Replace(base, gateBlock, tc.replace, 1)
			if src == base {
				t.Fatal("gate block substitution did not apply")
			}
			_, err := table.Load([]byte(src), "accessor-decl.toml")
			if err == nil {
				t.Fatalf("an accessor entry with %s loaded clean", tc.name)
			}
			if cat, _ := table.CategoryOf(err); cat != tc.want {
				t.Errorf("category = %q; want %q", cat, tc.want)
			}
		})
	}

	t.Run("promoted neg-accessor-no-timeout", func(t *testing.T) {
		got := loadCategory(t, "neg/neg-accessor-no-timeout.toml")
		if got != table.CatMalformedAccessorDeclaration {
			t.Errorf("category = %q; want %q", got, table.CatMalformedAccessorDeclaration)
		}
	})
}

// REQ-20: "A `[write.<id>]` entry additionally carries `read_back = true`;
// `read_back = false` or its absence on a write entry is the same
// category"
// ADVERSARIAL
func TestReq20_WriteEntryRequiresReadBackTrue(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	for name, replace := range map[string]string{
		"read_back = false": "read_back = false",
		"read_back absent":  "",
	} {
		t.Run(name, func(t *testing.T) {
			src := strings.Replace(base, "read_back = true", replace, 1)
			if src == base {
				t.Fatal("read_back substitution did not apply")
			}
			_, err := table.Load([]byte(src), "write-readback.toml")
			if err == nil {
				t.Fatalf("a [write.<id>] entry with %s loaded clean", name)
			}
			if cat, _ := table.CategoryOf(err); cat != table.CatMalformedAccessorDeclaration {
				t.Errorf("category = %q; want %q", cat, table.CatMalformedAccessorDeclaration)
			}
		})
	}
}

// REQ-21: "The same id MAY appear under two capability tables — RDR 0004's
// identity is `(flow, name, capability)`."
// HAPPY PATH
//
// The RDR fixture already authors `rdr-status` under both [read.*] and
// [write.*], so the positive control is the fixture loading clean and both
// entries surviving as distinct accessors.
func TestReq21_SameIDMayAppearUnderTwoCapabilityTables(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	if _, ok := m.Readers["rdr-status"]; !ok {
		t.Error("no reader named rdr-status")
	}
	if _, ok := m.Writers["rdr-status"]; !ok {
		t.Error("no writer named rdr-status")
	}
	// The two entries are distinct values, not one shared entry: only the
	// writer carries read_back.
	if !m.Writers["rdr-status"].ReadBack {
		t.Error("the writer entry did not carry read_back = true")
	}
}

// REQ-22: "every key in a `keys` list MUST be a declared tag (`unknown
// tag`), and MUST NOT be `recognized`"
// ADVERSARIAL
func TestReq22_KeysMustBeDeclaredAndNeverRecognized(t *testing.T) {
	t.Run("undeclared key in a reader's keys", func(t *testing.T) {
		base := string(readFixture(t, rdrFixture))
		src := strings.Replace(base, `keys = ["finalized_at"]`, `keys = ["nosuchtag"]`, 1)
		_, err := table.Load([]byte(src), "keys-undeclared.toml")
		if err == nil {
			t.Fatal("a keys list naming an undeclared tag loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatUnknownTag {
			t.Errorf("category = %q; want %q", cat, table.CatUnknownTag)
		}
	})

	t.Run("promoted neg-recognized-in-keys", func(t *testing.T) {
		got := loadCategory(t, "neg/neg-recognized-in-keys.toml")
		if got != table.CatMalformedAccessorBinding {
			t.Errorf("category = %q; want %q — the kernel supplies the "+
				"recognized key and no accessor reads or writes it",
				got, table.CatMalformedAccessorBinding)
		}
	})
}

// REQ-23: "every **owned** tag MUST be served by exactly one reader; an
// **observed** tag MAY be served by at most one reader, and zero is legal
// because an observed key may arrive from the caller (JDR 0001 §JD-9)"
// BOUNDARY
//
// The asymmetry is the whole clause, so both arms are asserted: the owned
// arm refuses at zero AND at two, while the observed arm accepts zero.
func TestReq23_ReaderArityIsProvenanceScoped(t *testing.T) {
	t.Run("owned tag served by zero readers refuses", func(t *testing.T) {
		got := loadCategory(t, "neg/neg-owned-no-reader.toml")
		if got != table.CatMalformedAccessorBinding {
			t.Errorf("category = %q; want %q", got, table.CatMalformedAccessorBinding)
		}
	})
	t.Run("owned tag served by two readers refuses", func(t *testing.T) {
		got := loadCategory(t, "neg/neg-owned-two-readers.toml")
		if got != table.CatMalformedAccessorBinding {
			t.Errorf("category = %q; want %q", got, table.CatMalformedAccessorBinding)
		}
	})
	t.Run("observed tag served by zero readers loads (JD-9)", func(t *testing.T) {
		// The kata fixture's `owner` tag is observed and served by no
		// reader. Refusing it would make every --tag-supplied key
		// unauthorable.
		m := mustLoad(t, kataFixture)
		decl, ok := m.Tags["owner"]
		if !ok {
			t.Fatal("kata fixture has no `owner` tag declaration")
		}
		if decl.Provenance != table.ProvenanceObserved {
			t.Fatalf("`owner` provenance = %q; want observed", decl.Provenance)
		}
		for id, r := range m.Readers {
			for _, k := range r.Keys {
				if k == "owner" {
					t.Fatalf("reader %q serves `owner`; the fixture's point is "+
						"that it is served by none", id)
				}
			}
		}
	})
}

// REQ-24: "every key named by any rule's write block or clear list, and
// every key assigned in `[initial]`, MUST be served by exactly one writer
// (`malformed accessor binding`)"
// ADVERSARIAL
func TestReq24_EveryWrittenKeyIsServedByExactlyOneWriter(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	t.Run("written key with no writer", func(t *testing.T) {
		// Drop `stage` from the writer's keys; three rules write it.
		src := strings.Replace(base,
			`[write.rdr-status]
role = "rdr"
path = "rdr.status"
keys = ["status", "stage", "profile", "iter", "rewind_scope", "prelock_lens", "cluster_ready"]`,
			`[write.rdr-status]
role = "rdr"
path = "rdr.status"
keys = ["status", "profile", "iter", "rewind_scope", "prelock_lens", "cluster_ready"]`, 1)
		if src == base {
			t.Fatal("writer keys substitution did not apply")
		}
		_, err := table.Load([]byte(src), "no-writer.toml")
		if err == nil {
			t.Fatal("a written key served by zero writers loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatMalformedAccessorBinding {
			t.Errorf("category = %q; want %q", cat, table.CatMalformedAccessorBinding)
		}
	})

	t.Run("[initial] key served by zero writers", func(t *testing.T) {
		got := loadCategory(t, "neg/probe-a-initial-observed-no-writer.toml")
		if got != table.CatMalformedAccessorBinding {
			t.Errorf("category = %q; want %q", got, table.CatMalformedAccessorBinding)
		}
	})

	t.Run("cleared key is a written key", func(t *testing.T) {
		// `prelock_lens` reaches the writer-binding rule only through
		// reconcile-rewind's clear list, so dropping it from the writer
		// keys must refuse.
		src := strings.Replace(base,
			`keys = ["status", "stage", "profile", "iter", "rewind_scope", "prelock_lens", "cluster_ready"]
timeout = "2s"
read_back = true`,
			`keys = ["status", "stage", "profile", "iter", "rewind_scope", "cluster_ready"]
timeout = "2s"
read_back = true`, 1)
		if src == base {
			t.Fatal("writer keys substitution did not apply")
		}
		_, err := table.Load([]byte(src), "clear-no-writer.toml")
		if err == nil {
			t.Fatal("a cleared key served by zero writers loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatMalformedAccessorBinding {
			t.Errorf("category = %q; want %q", cat, table.CatMalformedAccessorBinding)
		}
	})
}

// REQ-25: "every key in a `[write.<id>]` entry's `keys` list MUST be an
// **owned** tag (`write to non-owned tag`)"
// ADVERSARIAL
func TestReq25_WriterKeysMustBeOwned(t *testing.T) {
	t.Run("promoted neg-write-observed", func(t *testing.T) {
		got := loadCategory(t, "neg/neg-write-observed.toml")
		if got != table.CatWriteToNonOwnedTag {
			t.Errorf("category = %q; want %q", got, table.CatWriteToNonOwnedTag)
		}
	})
	t.Run("promoted probe-b: an observed [initial] key WITH a writer", func(t *testing.T) {
		got := loadCategory(t, "neg/probe-b-initial-observed-with-writer.toml")
		if got != table.CatWriteToNonOwnedTag {
			t.Errorf("category = %q; want %q", got, table.CatWriteToNonOwnedTag)
		}
	})
}

// REQ-26: "every id in a rule's `gate` list MUST resolve to a
// `[gate.<id>]` entry (`unknown accessor`)"
// ADVERSARIAL
//
// The discriminating case is an id that exists under a DIFFERENT
// capability table: a gate list naming a reader id must still refuse,
// because gate references resolve only against [gate.*].
func TestReq26_GateIDsResolveOnlyAgainstGateEntries(t *testing.T) {
	t.Run("gate names a reader id", func(t *testing.T) {
		got := loadCategory(t, "neg/neg-gate-is-reader.toml")
		if got != table.CatUnknownAccessor {
			t.Errorf("category = %q; want %q", got, table.CatUnknownAccessor)
		}
	})

	t.Run("gate names nothing at all", func(t *testing.T) {
		base := string(readFixture(t, rdrFixture))
		src := strings.Replace(base, `gate = ["rdr-lock"]`, `gate = ["nosuchgate"]`, 1)
		_, err := table.Load([]byte(src), "unknown-gate.toml")
		if err == nil {
			t.Fatal("a gate list naming no entry loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatUnknownAccessor {
			t.Errorf("category = %q; want %q", cat, table.CatUnknownAccessor)
		}
	})
}

// REQ-27: "A violation of the second or third bullet is a `malformed
// accessor binding`; the fourth carries `write to non-owned tag`."
// BOUNDARY
//
// The categories are stated separately on purpose — the Testing Strategy
// asserts on the category, not the message — so the discriminating
// assertion is that the two bullets do NOT collapse into one category.
func TestReq27_ReaderArityAndWriterProvenanceCarryDistinctCategories(t *testing.T) {
	binding := loadCategory(t, "neg/neg-owned-no-reader.toml")
	provenance := loadCategory(t, "neg/neg-write-observed.toml")

	if binding == provenance {
		t.Fatalf("both bullets refused with %q; the writer-arity and "+
			"writer-provenance obligations carry separate categories", binding)
	}
	if binding != table.CatMalformedAccessorBinding {
		t.Errorf("arity category = %q; want %q", binding, table.CatMalformedAccessorBinding)
	}
	if provenance != table.CatWriteToNonOwnedTag {
		t.Errorf("provenance category = %q; want %q", provenance, table.CatWriteToNonOwnedTag)
	}
}

// REQ-28: "A rule's `gate` list names the gate accessors whose `allow` the
// executor requires before applying that rule's plan; the list is carried
// on the normalized row and is part of its value."
// HAPPY PATH
func TestReq28_GateListIsCarriedOnTheNormalizedRow(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	for _, id := range []string{"rdr.continue-prelock#large", "rdr.continue-prelock#foundational"} {
		row := rowByID(t, m, id)
		if !reflect.DeepEqual(row.Gate, []string{"rdr-lock"}) {
			t.Errorf("%s gate list = %v; want [rdr-lock]", id, row.Gate)
		}
	}
	// A rule with no gate list carries an empty one, not a missing field.
	cluster := rowByID(t, m, "rdr.continue-prelock-cluster#large")
	if len(cluster.Gate) != 0 {
		t.Errorf("continue-prelock-cluster gate list = %v; want empty", cluster.Gate)
	}
}

// REQ-29: "An escape rule MUST NOT carry a `gate` list (it yields no plan
// to gate) — `malformed escape declaration`."
// ADVERSARIAL
func TestReq29_EscapeRuleMustNotCarryAGateList(t *testing.T) {
	got := loadCategory(t, "neg/neg-escape-with-gate.toml")
	if got != table.CatMalformedEscapeDeclaration {
		t.Errorf("category = %q; want %q", got, table.CatMalformedEscapeDeclaration)
	}
}

// REQ-30: "`[initial]` is a table of `tag = value` assignments; every key
// MUST be a declared owned tag and every value MUST be well-formed for
// that tag's declared kind and domain (`malformed initial declaration`
// otherwise; the `<clear>` sentinel is refused under the reserved-value
// rule)."
// ADVERSARIAL
//
// The category is scoped by SITE, not by defect: an [initial] value
// outside its declared domain is `malformed initial declaration`, while
// the same defect on a predicate literal is `malformed predicate atom`
// (req-list Q2, proceeding on the TS-3 fixture names). The `<clear>`
// sentinel is the one [initial] value that leaves this category.
func TestReq30_InitialAssignmentsAreValidatedBySite(t *testing.T) {
	cases := map[string]table.Category{
		"neg/neg-initial-bad-value.toml":        table.CatMalformedInitialDeclaration,
		"neg/neg-initial-int-out-of-range.toml": table.CatMalformedInitialDeclaration,
		// An [initial] key naming no declared tag is `unknown tag`, not
		// this category: the defect is the undeclared tag rather than the
		// declaration's shape.
		"neg/neg-initial-unknown.toml": table.CatUnknownTag,
		// The reserved-value rule takes precedence over the value arm.
		"neg/neg-clear-in-initial.toml": table.CatReservedTagValue,
	}
	for rel, want := range cases {
		t.Run(rel, func(t *testing.T) {
			if got := loadCategory(t, rel); got != want {
				t.Errorf("category = %q; want %q", got, want)
			}
		})
	}

	t.Run("well-formed [initial] loads and is carried", func(t *testing.T) {
		m := mustLoad(t, rdrFixture)
		want := map[string][]string{
			"status": {"Draft"},
			"stage":  {"propose"},
			"iter":   {"0"},
		}
		if len(m.Initial) != len(want) {
			t.Fatalf("[initial] carries %d assignments; want %d (%+v)",
				len(m.Initial), len(want), m.Initial)
		}
		for k, v := range want {
			got, ok := tagValue(m.Initial, k)
			if !ok {
				t.Errorf("[initial] lost the %q assignment", k)
				continue
			}
			if !reflect.DeepEqual(got, v) {
				t.Errorf("[initial] %q = %v; want %v", k, got, v)
			}
		}
	})
}

// REQ-31: "Root `terminal` is a list of context ids, each of which MUST
// resolve (`unknown context`)"
// ADVERSARIAL
func TestReq31_TerminalIDsMustResolve(t *testing.T) {
	got := loadCategory(t, "neg/neg-terminal-unknown.toml")
	if got != table.CatUnknownContext {
		t.Errorf("category = %q; want %q", got, table.CatUnknownContext)
	}
}

// REQ-32: "**A terminal context id is a reference, not a state name:
// normalization MUST dereference each to the context's explicit predicate
// set over owned tags**, in the same atom shape rules use, and the
// normalized value carries the predicate sets — never the bare ids."
// BOUNDARY
//
// The discriminating assertion is the shape: a normalizer carrying bare
// ids passes any test that only checks the terminal list is non-empty.
// The RDR fixture's `archived` context is `stage.eq=archive`.
func TestReq32_TerminalIsDereferencedToPredicateSets(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	if len(m.Terminal) != 1 {
		t.Fatalf("terminal carries %d predicate sets; want 1 (%+v)",
			len(m.Terminal), m.Terminal)
	}
	want := []table.Atom{{
		Key:      "stage",
		Block:    table.BlockMatch,
		Operator: "eq",
		Literal:  []string{"archive"},
	}}
	if !reflect.DeepEqual(m.Terminal[0], want) {
		t.Errorf("terminal predicate set = %+v; want %+v — never the bare id",
			m.Terminal[0], want)
	}
}

// REQ-33: "Whether a model *omits* `[initial]` or `terminal`, and whether
// a terminal context matches a non-owned tag, are RDR 0006's blocking
// findings, not load failures: the loader accepts the absence"
// INPUT EDGE
func TestReq33_AbsentInitialOrTerminalIsNotALoadFailure(t *testing.T) {
	// Each arm asserts the absence is ACCEPTED *and* that the rest of the
	// document still normalized: a loader returning an empty model with no
	// error would otherwise pass by doing nothing.
	t.Run("no [initial]", func(t *testing.T) {
		m := mustLoad(t, "pos-no-initial.toml")
		if len(m.Initial) != 0 {
			t.Errorf("[initial] = %+v; want empty", m.Initial)
		}
		if len(m.Rows) != 9 {
			t.Errorf("%d candidate rows; want 9 — the absence is accepted and "+
				"the rest of the document still normalizes", len(m.Rows))
		}
		if len(m.Terminal) == 0 {
			t.Error("the stop set was dropped along with [initial]")
		}
	})
	t.Run("no terminal", func(t *testing.T) {
		m := mustLoad(t, "pos-no-terminal.toml")
		if len(m.Terminal) != 0 {
			t.Errorf("terminal = %+v; want empty", m.Terminal)
		}
		if len(m.Rows) != 9 {
			t.Errorf("%d candidate rows; want 9", len(m.Rows))
		}
		if len(m.Initial) == 0 {
			t.Error("the root was dropped along with `terminal`")
		}
	})
	t.Run("terminal context matching a non-owned tag", func(t *testing.T) {
		// The kata fixture's `owner` is observed. A terminal context over
		// it is RDR 0006's finding, not a load failure.
		base := string(readFixture(t, kataFixture))
		src := strings.Replace(base,
			"[context.shipped.match.phase]\neq = \"ship\"",
			"[context.shipped.match.owner]\neq = \"current-session\"", 1)
		if src == base {
			t.Fatal("terminal context substitution did not apply")
		}
		m, err := table.Load([]byte(src), "terminal-observed.toml")
		if err != nil {
			t.Fatalf("a terminal context over a non-owned tag refused with %v; "+
				"it is RDR 0006's blocking finding, not a load failure", err)
		}
		// And the non-owned predicate is dereferenced and carried, not
		// silently dropped.
		// The whole atom is the oracle, not just its key: a dereference
		// producing the right key under the wrong block, operator, or
		// literal is still a defect.
		want := table.Atom{
			Key:      "owner",
			Block:    table.BlockMatch,
			Operator: "eq",
			Literal:  []string{"current-session"},
		}
		if len(m.Terminal) != 1 || len(m.Terminal[0]) != 1 ||
			!reflect.DeepEqual(m.Terminal[0][0], want) {
			t.Errorf("terminal predicate sets = %+v; want the dereferenced "+
				"`owner.eq=current-session` atom %+v", m.Terminal, want)
		}
	})
}
