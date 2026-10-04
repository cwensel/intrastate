package table

// Kata wchf — per-key `select` tables on a `command` read entry, so a
// tool whose `--json` output is nested is read through the MODEL rather
// than through a normalizing wrapper script lint cannot see.
//
// The tables are OPT-IN and narrow `0025:C3` rather than lifting it: an
// entry declaring no `select` reads exactly the flat JSON object of
// strings C3 fixes, and an entry declaring any must declare one per key.
//
// The locator is RFC 6901 JSON Pointer and nothing else. A jq subset has
// no standard subset and would be a second query dialect — the reason
// 0028 gave for one template parser — and a dotted path is ambiguous for
// a key that itself contains `.`. The parser lives in this package rather
// than beside the read binding for the reason `edit.go`'s does: the loader
// validates every pointer at load and the binding resolves them at read,
// and two parsers would be two dialects.

import (
	"errors"
	"maps"
	"slices"
	"strconv"
	"strings"
)

// The closed `absent` vocabulary of a `select.<key>` table. Each member
// turns one shape of "the tool reported no value" into ESTABLISHED
// absence, which is otherwise only `exit_absent`'s to declare (`0025:C3`).
const (
	// SelectAbsentNull: a JSON null at the pointer is established-absent.
	SelectAbsentNull = "null"
	// SelectAbsentMissing: a missing FINAL reference token under a parent
	// that exists is established-absent. A missing INTERMEDIATE node is
	// shape drift and stays unreadable whatever is declared.
	SelectAbsentMissing = "missing"
)

// selectAbsentMembers is the closed set, in the order a refusal names it.
var selectAbsentMembers = []string{SelectAbsentNull, SelectAbsentMissing}

// ParseJSONPointer splits an RFC 6901 JSON Pointer into its unescaped
// reference tokens. The empty string is the whole document and yields no
// tokens; any other pointer must begin with `/`. `~1` decodes to `/` and
// `~0` to `~`, in that order, and a `~` followed by anything else is a
// defect rather than literal text (RFC 6901 §3, §4).
func ParseJSONPointer(p string) ([]string, error) {
	if p == "" {
		return nil, nil
	}
	if p[0] != '/' {
		return nil, errors.New("pointer " + strconv.Quote(p) +
			" does not begin with `/`; an RFC 6901 pointer is empty or `/`-led")
	}
	raw := strings.Split(p[1:], "/")
	tokens := make([]string, 0, len(raw))
	for _, tok := range raw {
		for i := 0; i < len(tok); i++ {
			if tok[i] == '~' && (i+1 == len(tok) || (tok[i+1] != '0' && tok[i+1] != '1')) {
				return nil, errors.New("pointer " + strconv.Quote(p) +
					" carries a `~` that is not `~0` or `~1`")
			}
		}
		tok = strings.ReplaceAll(tok, "~1", "/")
		tok = strings.ReplaceAll(tok, "~0", "~")
		tokens = append(tokens, tok)
	}
	return tokens, nil
}

// selectDefect reports the FIRST `select` defect an entry carries, in
// registration order: placement, key mismatch, invalid table. Sibling
// tables are swept in SORTED key order so the reported key is
// deterministic, as `editTableDefect`'s are.
func selectDefect(
	a sourceAcc, capability, where string, keys []string,
	tags map[string]TagDecl,
) error {
	if a.Select == nil {
		return nil
	}
	rules := *a.Select

	// command_select_placement — a pointer selects out of a COMMAND'S json
	// stdout. A path entry has no stdout, a raw entry has no document, and
	// only a read entry's stdout is parsed into keys (`0025:C3`).
	switch {
	case capability != "read":
		return fail(CatCommandSelectPlacement,
			where+" declares `select`, which is a read-entry key")
	case a.Command == nil:
		return fail(CatCommandSelectPlacement,
			where+" declares `select` without a `command`; a selector reads "+
				"a command's json stdout")
	case a.Output != nil && *a.Output != "json":
		return fail(CatCommandSelectPlacement,
			where+" declares `select` with output = "+strconv.Quote(*a.Output)+
				`; a selector reads a json document, never raw stdout`)
	}

	// command_select_key_mismatch — `keys` and the tables in bijection,
	// mirroring `edit.<key>`: every requested key has a selector, so no
	// key falls back to the flat-object reading on one entry.
	for _, key := range keys {
		if _, ok := rules[key]; !ok {
			return fail(CatCommandSelectKeyMismatch,
				where+" declares the key "+key+" with no `select."+key+"` table")
		}
	}
	for _, key := range slices.Sorted(maps.Keys(rules)) {
		if !slices.Contains(keys, key) {
			return fail(CatCommandSelectKeyMismatch,
				where+" declares a `select."+key+"` table for a key not in `keys`")
		}
	}

	// command_select_invalid — each table's well-formedness. Lint cannot
	// see the tool's output, so what it proves is the declaration: the
	// pointers parse, the projection is shaped for the tag it fills, and
	// `absent` stays inside its closed set.
	for _, key := range slices.Sorted(maps.Keys(rules)) {
		r := rules[key]
		at := where + " `select." + key + "`"
		bad := func(detail string) error {
			return fail(CatCommandSelectInvalid, at+" "+detail)
		}
		if r.Pointer == nil {
			return bad("declares no `pointer`")
		}
		if _, err := ParseJSONPointer(*r.Pointer); err != nil {
			return bad(err.Error())
		}
		if r.Element != nil {
			if r.Prefix == nil {
				return bad("declares `element` without `prefix`; `element` " +
					"reaches into each member of the array a prefix projects")
			}
			if _, err := ParseJSONPointer(*r.Element); err != nil {
				return bad("element " + err.Error())
			}
		}
		if r.Prefix != nil {
			if *r.Prefix == "" {
				return bad("declares an empty `prefix`, which every member " +
					"matches")
			}
			// A prefix projects AT MOST ONE member, so the tag it fills
			// must hold one value. A set holds any subset, and reading
			// one member into it would silently drop the rest.
			if tags[key].Kind == "set" {
				return bad("declares `prefix` for the set-kind tag " + key +
					"; a prefix projects one member and a set holds many")
			}
		}
		if r.Absent != nil {
			seen := map[string]bool{}
			for _, m := range *r.Absent {
				if !slices.Contains(selectAbsentMembers, m) {
					return bad("declares the `absent` member " + strconv.Quote(m) +
						", outside the closed set {\"" +
						strings.Join(selectAbsentMembers, "\", \"") + "\"}")
				}
				if seen[m] {
					return bad("declares the `absent` member " + strconv.Quote(m) +
						" twice")
				}
				seen[m] = true
			}
		}
	}
	return nil
}

// selectRules converts the decoded tables into the model's typed form.
// PRESENCE survives: a bare `[read.x.select]` decodes to a non-nil empty
// map, as `edit` does.
func selectRules(src *map[string]sourceSelectRule) map[string]SelectRule {
	if src == nil {
		return nil
	}
	out := make(map[string]SelectRule, len(*src))
	for key, r := range *src {
		rule := SelectRule{}
		if r.Pointer != nil {
			rule.Pointer = *r.Pointer
		}
		if r.Element != nil {
			rule.Element = *r.Element
		}
		if r.Prefix != nil {
			rule.Prefix = *r.Prefix
		}
		if r.Absent != nil {
			rule.AbsentNull = slices.Contains(*r.Absent, SelectAbsentNull)
			rule.AbsentMissing = slices.Contains(*r.Absent, SelectAbsentMissing)
		}
		out[key] = rule
	}
	return out
}
