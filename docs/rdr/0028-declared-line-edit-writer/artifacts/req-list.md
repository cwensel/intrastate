# REQ list — RDR 0028 declared-line-edit-writer

Phase 0 spec audit. Source: `docs/rdr/0028-declared-line-edit-writer.md`
(1632 lines; one durable contract C1 with six clauses C1.1–C1.6, MVV=1 (8 steps),
A=12 (A1–A12, all `Verified`), S=32 scenarios, F=8 failure modes, D=3 pre-seeded
deviations).

Element ids carried where a REQ derives from a labelled contract clause. Quotes
are exact bytes from the record's `normative` fence (lines 481–538) and from
testable prose outside it, reflowed only where a clause spans lines in the
source; no wording is changed.

Projection note: the record's `contracts_prose` flag is set and no addressable
element projection is available for this repo — `intrastate` is the *subject*
binary (a workflow-state CLI: `lint` / `flow`), not an RDR projector, and it has
no `inspect` command. The record was therefore read whole and the fence bytes
copied directly from `docs/rdr/0028-declared-line-edit-writer.md`. A low or zero
element count does NOT mean the record has no contracts: C1's six clauses are
fully normative and are the whole implementation surface below.

---

## C1.1 — carrier and grammar (`0028:C1.1`)

- [REQ-1] "[write.<id>.edit.<key>]           # exactly one table per member of `keys`; a table for a key not in `keys` is a defect (C1.4)" — (0028:C1.1, §normative-contracts)
- [REQ-2] "anchor  = \"<RE2 pattern>\"         # a line is SELECTED when the pattern matches anywhere in it (terminator excluded); authors pin `^…$`; must select exactly one line (C1.3)" — (0028:C1.1, §normative-contracts)
- [REQ-3] "replace = \"<template>\"            # the WHOLE replacement line, terminator excluded (C1.2) — never the matched span; a span-style anchor without `^…$` drops the unmatched prefix/suffix by design" — (0028:C1.1, §normative-contracts)
- [REQ-4] "clear   = \"line\"                  # optional; the disposition of a planned `<clear>` (C1.5); absent ⇒ `<clear>` refuses before mutation" — (0028:C1.1, §normative-contracts)
- [REQ-5] "role, keys, timeout               # unchanged (0002, 0004:C15)" — (0028:C1.1, §normative-contracts) — NEGATIVE REQ: these three fields are not redefined by this record; `timeout` stays unconditionally required for every entry regardless of carrier (A11).
- [REQ-6] "carrier: exactly one of `path` / `command` / `edit` per entry — 0025:C1's exactly-one rule with one member appended, not replaced" — (0028:C1.1, §normative-contracts)
- [REQ-7] "`edit` is admissible on WRITE entries only (a read or gate entry carrying `edit` is `edit_carrier_conflict`, C1.4)" — (0028:C1.1, §normative-contracts)
- [REQ-8] "The \"neither\" arm keeps its 0025:C5 wire string `command_and_path_conflict`, but its PREDICATE necessarily widens to \"none of the three\" — `internal/table/load.go::carrierDefect`'s `case !hasPath && !hasCommand:` gains `&& !hasEdit`, or an `edit`-only entry would be refused as carrier-less (S1 asserts it LOADS)." — (0028:C1.1, §normative-contracts)
- [REQ-9] "Its hardcoded message text, which names two carriers, is updated to name three; the wire string is the contract and the message is not (C1.4)." — (0028:C1.1, §normative-contracts)
- [REQ-10] "`edit_carrier_conflict` follows the \"both\" arm's discipline — keyed on the `edit` table being PRESENT, not on it being non-empty, matching `case a.Path != nil && a.Command != nil:` so an empty `edit` beside another carrier is a conflict rather than a silently ignored second carrier" — (0028:C1.1, §normative-contracts)
- [REQ-11] "in-memory: the registry's residue rule (0025:C1 runtime arm) is unchanged — an entry with no carrier builds a refusing binding, never a file binding" — (0028:C1.1, §normative-contracts)
- [REQ-12] "`registry.go::commandBacked` becomes a three-way carrier discriminator: `command` → command binding; `edit` → edit binding; `path` → file binding; residue → refusing binding, exactly as 0025:C1's runtime arm." — (0028 Load-Bearing Decisions §Selection/predicate, §proposed-solution)
- [REQ-13] "the accessor identity is 0004's `(flow, name, capability)`; the carrier does not enter it (0025 Identity LBD). A line rule's identity is `(entry, key)`." — (0028 Load-Bearing Decisions §Identity, §proposed-solution)
- [REQ-14] "the carrier is `edit`; the rule tables are `edit.<key>`; the fields are `anchor` / `replace` / `clear`." — (0028 Load-Bearing Decisions §Naming, §proposed-solution)

## C1.2 — substitution admission (`0028:C1.2`)

- [REQ-15] "replace admits: literal text | ${N} and ${name} — the anchor's capture groups | {<key>} — the planned value of THE key this table is named for, and no other key" — (0028:C1.2, §normative-contracts)
- [REQ-16] "anchor  admits: literal RE2 | {tag.<key>} — a tag key the model declares with `observed` provenance, bound on the invocation's context (A1, A7); substituted regexp-quoted, so a bound value is a literal-match fragment and can never alter the pattern's structure." — (0028:C1.2, §normative-contracts)
- [REQ-17] "a declared key of either kind is structurally unbindable at invocation, and `{tag.<key>}` naming one is `edit_anchor_invalid` at LINT (C1.4) rather than a guaranteed runtime `execution_failure`" — (0028:C1.2, §normative-contracts) — the owned-key (`flow-tag-owned`) and recognized-key (`flow-tag-reserved`) refusals in `internal/cli/flow_input.go::parseTags` are the ground.
- [REQ-18] "The entry's own planned keys are NOT admissible in `anchor`: the value being written never decides where it is written" — (0028:C1.2, §normative-contracts) — NEGATIVE REQ.
- [REQ-19] "parse once: both templates are parsed at LOAD into segments (literal | group | placeholder). At apply each segment emits bytes; captured text and substituted values are never re-scanned for `${…}` or `{…}`." — (0028:C1.2, §normative-contracts)
- [REQ-20] "value shape: a planned value, or a bound tag value THIS ENTRY's rules reference, containing \"\\n\" or \"\\r\" refuses BEFORE mutation (Detail `edit_value_multiline`) — the single structural hazard of interpolation into line data." — (0028:C1.2, §normative-contracts)
- [REQ-21] "No word-splitting, option parsing, PATH resolution or shell exists on this path" — (0028:C1.2, §normative-contracts) — NEGATIVE REQ.
- [REQ-22] "escapes: `$$` emits a literal `$`; `{{` and `}}` emit literal braces in `replace`; a group that did not participate in the match expands to the empty string" — (0028:C1.2, §normative-contracts)
- [REQ-23] "the empty-group rule MATCHES Go `regexp.Expand`'s, but the parse is this clause's own and never a call to `Expand`, which re-scans its template at expansion time and would break `parse once:`" — (0028:C1.2, §normative-contracts) — NEGATIVE REQ: `regexp.Expand` must not be called on this path.
- [REQ-24] "In `replace` the vocabulary is closed — `{<key>}`, `${N}`/`${name}`, `$$`, `{{`/`}}` — and any other `{…}` or `$…` form is `edit_template_invalid`, including a bare `$1`/`$name` without braces (`Expand` accepts those; this clause does not)." — (0028:C1.2, §normative-contracts)
- [REQ-25] "In `anchor` ONLY the fixed prefix `{tag.` opens a placeholder, scanned to its closing `}`; every other `{`, `}` and `$` is passed to RE2 untouched, so `\\d{4}`, `a{2,3}` and `^…$` are ordinary pattern text and `{{`/`}}` is NOT an escape there." — (0028:C1.2, §normative-contracts)
- [REQ-26] "`edit_anchor_invalid` fires when a `{tag.…}` placeholder is malformed or names a key not declared `observed`, or when the pattern fails to compile — never for a brace RE2 itself accepts." — (0028:C1.2, §normative-contracts)
- [REQ-27] "A line whose CONTENT carries the literal text `{tag.` is anchored with RE2's own escape (`\\{tag\\.`): the placeholder scan is for the unescaped literal prefix, so an escaped brace is ordinary pattern text and reaches RE2 untouched." — (0028:C1.2, §normative-contracts)
- [REQ-28] "the escaping dialect in `anchor` is RE2's, in `replace` it is this clause's" — (0028:C1.2, §normative-contracts)
- [REQ-29] "out of scope for admission: a value that is well-formed line data but wrong for the DOCUMENT (a `|` inside a markdown table cell) is not a hazard this clause can name without document knowledge — the reader's read-back is the check (0004:C12), post-mutation by contract" — (0028:C1.2, §normative-contracts) — NEGATIVE REQ: no per-rule value guard in v1.
- [REQ-30] "a line is the bytes up to and excluding `\\n`, with a preceding `\\r` treated as part of the terminator; the file is bytes, not a decoded string, and RE2 matches over UTF-8 bytes. Rejected: normalising line endings on write (would touch bytes no rule selected)." — (0028 Load-Bearing Decisions §Wire/byte format, §proposed-solution)

## C1.3 — execution semantics (`0028:C1.3`)

- [REQ-31] "input:    the whole file is read as bytes; lines are split on \"\\n\" and a preceding \"\\r\" stays with the terminator (CRLF preserved per line); a missing final terminator is preserved; the file's mode is preserved on write (A5)." — (0028:C1.3, §normative-contracts)
- [REQ-32] "A target that cannot be read refuses BEFORE mutation with the OS error in the Detail — it does NOT reuse `edit_anchor_unmatched`, and it does NOT take `flowbind.go::load`'s absent-is-empty precedent" — (0028:C1.3, §normative-contracts)
- [REQ-33] "target:   the caller-bound artifact path for the entry's role (0004:C3), symlinks resolved; the model names no path — an `edit` entry has exactly the authority a `path` entry has over the file the CALLER binds." — (0028:C1.3, §normative-contracts)
- [REQ-34] "The reuse of `flowbind.go::save`'s discipline stops at stage-and-rename and does NOT drag along `Writer`'s unreachable-locator seal (`unreachable(path)`/`sealedKey`)" ... "`edit` therefore has no seal affordance, by construction rather than by omission" — (0028:C1.3, §normative-contracts) — NEGATIVE REQ.
- [REQ-35] "select:   every rule's anchor is resolved against the PRE-EDIT content and selections are held as pre-edit line INDICES (a deletion never shifts a sibling rule's target)." — (0028:C1.3, §normative-contracts)
- [REQ-36] "Each must select exactly one line: 0 ⇒ `edit_anchor_unmatched`; ≥2 ⇒ `edit_anchor_ambiguous`; two rules selecting one line ⇒ `edit_anchor_collision`." — (0028:C1.3, §normative-contracts)
- [REQ-37] "Never last-match (Ansible `lineinfile`), never first-match, never insert or append (Puppet `append_on_no_match`, Ansible `insertafter`) — creation is fenced out, and an unmatched anchor is a stale model, not a missing line" — (0028:C1.3, §normative-contracts) — NEGATIVE REQ.
- [REQ-38] "re-anchor: after the buffer is rewritten in memory, every rule's anchor is run again over the POST-EDIT buffer and must select exactly its own rewritten line (or, for a deleted line, zero lines); otherwise refuse `edit_anchor_unstable` before any write." — (0028:C1.3, §normative-contracts)
- [REQ-39] "IDENTITY, not cardinality: the selected line must BE the rule's own — its held pre-edit index, shifted by the deletions of preceding sibling rules — and a rule that selects exactly one line which is a DIFFERENT line refuses." — (0028:C1.3, §normative-contracts)
- [REQ-40] "order:    every refusal in this clause and C1.2's `edit_value_multiline` is decided BEFORE any byte is written. A refused edit is NOT APPLIED and surfaces as 0004's `execution_failure`" — (0028:C1.3, §normative-contracts)
- [REQ-41] "a rule-scoped refusal carries a Detail naming the rule (`<id>.edit.<key>`) and the reason token (A8), and the entry-level preconditions — the read-back gate below, C1.6's unbound tag and C1.6's `-`-prefixed value — name the gate or the placeholder instead, having no rule to name" — (0028:C1.3, §normative-contracts)
- [REQ-42] "no new refusal class is introduced, and no pre-write refusal THIS BINDING MINTS ever carries 0004:C14's applied-but-unverified sense (applied sense per JDR 0003 §D1)" — (0028:C1.3, §normative-contracts) — NEGATIVE REQ.
- [REQ-43] "EXIT GROUP: these refusals are about the REQUEST, not the environment, so they take the exit-2 group, not `execution_failure`'s default exit 3 — the executor-facing typed `Err` discriminates at `flow_exec.go::accessorFailureOf` and a distinct CLI code carries the rule id and reason token in `findings[]` (JDR 0003 §D3 (b)" — (0028:C1.3, §normative-contracts)
- [REQ-44] "the class set and the Detail above are unchanged, and the code's spelling is Stage 8's, non-normative here" — (0028:C1.3, §normative-contracts) — NEGATIVE REQ: no test may pin the CLI code string (restated at S27).
- [REQ-45] "That arm stays 0004's, unchanged and out of this RDR's scope (A11, `Verified`): this clause's guarantee is scoped to refusals the `edit` binding itself mints, and does NOT assert the executor-minted timeout carries a false applied sense." — (0028:C1.3, §normative-contracts) — NEGATIVE REQ: `internal/accessor/executor.go`'s deadline arm is not amended.
- [REQ-46] "A rename that fails after a good staged write is also NOT APPLIED — the target is untouched by construction" — (0028:C1.3, §normative-contracts)
- [REQ-47] "write:    all rules of one entry rewrite ONE buffer and land in ONE write: staged beside the resolved target and renamed over it (`internal/cli/flowbind/flowbind.go::save`'s discipline, mode preserved, A5)." — (0028:C1.3, §normative-contracts)
- [REQ-48] "The atomicity boundary is ONE ENTRY over ONE file, and it is not transactional across entries" ... "a cross-entry transaction is not introduced" — (0028:C1.3, §normative-contracts) — NEGATIVE REQ.
- [REQ-49] "A post-edit buffer equal to the input is not written at all (no staging, no rename; S20 asserts it, witnessed by an unchanged inode" — (0028:C1.3, §normative-contracts)
- [REQ-50] "No lock and no compare-before-rename: a concurrent writer is out of scope, as it is for `path`" — (0028:C1.3, §normative-contracts) — NEGATIVE REQ.
- [REQ-51] "precedence: apply-time refusals fail-fast within one entry in this order, the mirror of C1.4's for load time: (1) the ENTRY-level preconditions — C1.6's unbound `{tag.<key>}`, then C1.6's `-`-prefixed bound value (both are properties of the binding, decided together before anything is spawned or read), then the gate-off command read-back above" — (0028:C1.3, §normative-contracts)
- [REQ-52] "(2) `edit_value_multiline` over the entry's planned values and the tag values its rules actually reference (C1.2 `value shape:`), the scope being per-USE-SITE, not per-invocation: a tag bound on the context but referenced by no `anchor` of this entry is never scanned" — (0028:C1.3, §normative-contracts)
- [REQ-53] "(3) `edit_clear_undeclared` (C1.5) — a `<clear>` plan on a rule that did not declare `clear` is decided on the RAW planned value, BEFORE selection and before any `replace` expansion, so the refusal never depends on whether that rule's anchor matched" — (0028:C1.3, §normative-contracts)
- [REQ-54] "(4) per-rule cardinality, `edit_anchor_unmatched` then `edit_anchor_ambiguous`, resolved for EVERY rule of the entry before (5) the cross-rule `edit_anchor_collision` sweep, which is only decidable once every rule holds a selection; (6) `edit_anchor_unstable`, necessarily last, being post-rewrite." — (0028:C1.3, §normative-contracts)
- [REQ-55] "Within one step, siblings are map-ranged and inherit C1.4's rule unchanged — which of two equally-defective rules is named is unspecified and no test may assert it." — (0028:C1.3, §normative-contracts) — NEGATIVE REQ.
- [REQ-56] "terminators: only \"\\n\" — optionally preceded by \"\\r\" — terminates a line for selection and rewriting; a bare \"\\r\", NEL or U+2028 is line content (C1.2 still refuses \"\\r\" in a VALUE)." — (0028:C1.3, §normative-contracts)
- [REQ-57] "Deleting the final line of a file that had no final terminator also removes the preceding line's terminator, so the file's final-terminator state is preserved either way" — (0028:C1.3, §normative-contracts)
- [REQ-58] "no subprocess: `edit` spawns nothing; `--allow-commands` (0025:C6) is not consulted by the write itself. `Invocations()` counts `Apply` calls exactly as `flowbind.go::Writer` does (0004:C14)" — (0028:C1.3, §normative-contracts)
- [REQ-59] "read-back: unchanged — 0004:C12/0004:C13 through the role's declared reader. When that reader is command-backed and the gate is off, the write MUST refuse BEFORE mutation (Detail naming the gate), because the reader is resolved before `Apply` (A2) — a forgotten flag must not produce `read_back_incomplete` for a write that ran no command." — (0028:C1.3, §normative-contracts)
- [REQ-60] "AUTHORITATIVE DETECTOR: this pre-check, not `Executor.Write`'s existing pre-`Apply` baseline read." ... "This clause does NOT amend that arm — it refuses earlier and unconditionally, so the swallow becomes unreachable for the gate-off case whether or not `protected` is empty. The baseline arm needs no change (A2, `Verified`)" — (0028:C1.3, §normative-contracts) — NEGATIVE REQ: `protectedKeys`/`baselineUnread` are untouched.
- [REQ-61] "SITE (A10, `Verified`): both halves of that predicate are answered at the executor without a new import or a `Binding` method — command-backed from `def.Accessor.Command` (already on `Definition`), and the gate from a field on `accessor.Registry`, set at `flowbind.go::Registry`, the single production construction site" — (0028:C1.3, §normative-contracts)
- [REQ-62] "`internal/accessor` does NOT import `cmdbind` (`cmdbind` imports `accessor`; the reverse is a cycle) and `AllowCommands` on `cmdbind.Config` stays where it is — the gate is carried to the accessor as state, never read across the seam (A10)" — (0028:C1.3, §normative-contracts) — NEGATIVE REQ.

## C1.4 — load-time categories (`intrastate lint`) (`0028:C1.4`)

- [REQ-63] "edit_carrier_conflict    # `edit` beside `path` or `command`; or `edit` on a read/gate entry" — (0028:C1.4, §normative-contracts)
- [REQ-64] "edit_key_mismatch        # a `keys` member with no `edit.<key>` table, or an `edit.<key>` table for a key not in `keys`" — (0028:C1.4, §normative-contracts)
- [REQ-65] "edit_anchor_invalid      # anchor fails to compile as RE2 (checked with every `{tag.<key>}` replaced by a quoted probe), names an undeclared tag key, or carries any other `{…}` form" — (0028:C1.4, §normative-contracts)
- [REQ-66] "edit_template_invalid    # replace carries an unknown placeholder, another key's placeholder, a group reference the anchor does not define, or a malformed `${…}`/`{…}` token" — (0028:C1.4, §normative-contracts)
- [REQ-67] "edit_clear_invalid       # `clear` outside the closed set {\"line\"}" — (0028:C1.4, §normative-contracts)
- [REQ-68] "edit_tag_argv0           # a `{tag.<key>}` element at argv0 of a read, gate or write entry's `command` (C1.6)" — (0028:C1.4, §normative-contracts)
- [REQ-69] "registration: appended to `table.Categories()` after 0025:C5's six, in the order above; typed `Cat…` constants beside the others; the wire strings are the contract, the identifiers are not; the list's size is not a contract (0025:C5)." — (0028:C1.4, §normative-contracts)
- [REQ-70] "`edit_tag_argv0` carries the `edit_` prefix because it is this record's category, minted with the `{tag.<key>}` family C1.6 introduces — it fires on `command` entries of every kind, including entries carrying no `edit` table" — (0028:C1.4, §normative-contracts)
- [REQ-71] "precedence: within one entry, fail-fast in the order above, evaluated after 0025:C5's clauses 1–6 (an `edit` entry reaches only the clauses that are not `command`-argv-specific — 0025:C5 owns which those are, and this clause neither restates nor narrows its map)" — (0028:C1.4, §normative-contracts)
- [REQ-72] "across entries and tables, 0025:C5's rules apply unchanged — including across the sibling `edit.<key>` tables of ONE entry, which are map-ranged like `accessorTable`'s entries: which of two equally-defective tables is reported is unspecified and no test may assert it." — (0028:C1.4, §normative-contracts) — NEGATIVE REQ.
- [REQ-73] "The first five categories fire on the `edit` table; `edit_tag_argv0` fires on an entry's `command` argv, so an entry carrying `command` reaches it while an `edit` entry never does — the two sets are disjoint by carrier and never race" — (0028:C1.4, §normative-contracts)
- [REQ-74] "what lint proves: the carrier is one, every anchor compiles, every template parses, every placeholder is closed over the entry's keys and the model's declared tags, every `keys` member has one rule, and no `{tag.<key>}` sits at argv0." — (0028:C1.4, §normative-contracts)
- [REQ-75] "What lint does NOT prove: that an anchor matches exactly one line of a particular file, or that a bound tag value is not flag-shaped — those are C1.3's and C1.6's apply-time refusals, by design (a model is linted without its artifacts and without the invocation's bindings)" — (0028:C1.4, §normative-contracts) — NEGATIVE REQ.
- [REQ-76] "The assertion is RELATIVE order, never a tail position or a count — the shipped `TestReq77` pins 0025:C5's six as `Categories()`'s tail and this append necessarily rewrites it (deviations D1)." — (0028 Testing Strategy S6, §validation) — see DEVIATION note below.

## C1.5 — `<clear>` on a line rule (`0028:C1.5`)

- [REQ-77] "clear = \"line\" ⇒ a planned `<clear>` deletes the anchored line, terminator included; read-back asserts the key ABSENT through the role's reader (0004:C11 unchanged)." — (0028:C1.5, §normative-contracts)
- [REQ-78] "An anchor matching ZERO lines on a `<clear>` plan is SUCCESS with no write — \"clearing a key the artifact does not hold MUST succeed\" (0004:C11); ≥2 matches still refuses `edit_anchor_ambiguous`" — (0028:C1.5, §normative-contracts)
- [REQ-79] "clear absent ⇒ a planned `<clear>` refuses BEFORE mutation, Detail `edit_clear_undeclared`, decided on the RAW planned value ahead of selection (C1.3 `precedence:`) — so an unmatched anchor on such a rule reports the undeclared `<clear>`, never `edit_anchor_unmatched`." — (0028:C1.5, §normative-contracts)
- [REQ-80] "one-way: a deleted line cannot be re-established by `edit` (no append, C1.3), so after a `clear` the next non-clear write refuses `edit_anchor_unmatched`" — (0028:C1.5, §normative-contracts)
- [REQ-81] "A blank-the-cell disposition (`clear = { replace = … }`) is deferred: it is a document-shape decision, and v1's closed set is {\"line\"}" — (0028:C1.5, §normative-contracts) — NEGATIVE REQ.
- [REQ-82] "the literal string `<clear>` is never substituted into `replace` (0004:C11: a re-read holding the literal is a mismatch). The sentinel is a property of the WHOLE planned value, tested before expansion: a `<clear>` authored into a `replace` template's LITERAL segment is ordinary bytes and emits as such, since only the planned value — never a template segment — is read as the sentinel" — (0028:C1.5, §normative-contracts)

## C1.6 — `{tag.<key>}` as a command placeholder (0025:C2 extended) (`0028:C1.6`)

- [REQ-83] "[read.<id> | gate.<id> | write.<id>]  command = [..., \"{tag.<key>}\", ...]" — (0028:C1.6, §normative-contracts)
- [REQ-84] "admission: `{tag.<key>}` joins 0025:C2's vocabulary as a family, WHOLE-ELEMENT only under 0025:C2's substitution rule, replaced by the bound value of a tag key the model DECLARES" — (0028:C1.6, §normative-contracts)
- [REQ-85] "`<key>` undeclared is `command_unknown_placeholder` (0025:C5, unchanged wire string); a `{…}` element that is neither `{artifact}` nor a declared `{tag.<key>}` stays `command_unknown_placeholder`." — (0028:C1.6, §normative-contracts)
- [REQ-86] "POSITION: a `{tag.<key>}` element at argv0 is `edit_tag_argv0` at LINT (C1.4) — the executable is the one word a reviewer must be able to read off the model" — (0028:C1.6, §normative-contracts)
- [REQ-87] "Statically decidable, so it is refused where it is visible rather than at spawn; `{artifact}` is unaffected, having no such rule and no such reviewer promise to break" — (0028:C1.6, §normative-contracts) — NEGATIVE REQ: `{artifact}` at argv0 stays admitted (S7).
- [REQ-88] "binding: the value comes from the invocation's context (the same channel as C1.2's anchor tags, A1); an unbound tag at invocation refuses `execution_failure` BEFORE spawn, Detail naming the placeholder — a placeholder is never passed through literally (0025:C2)." — (0028:C1.6, §normative-contracts)
- [REQ-89] "VALUE: a bound value beginning with `-` refuses `execution_failure` BEFORE spawn, Detail naming the placeholder, mirroring `cmdbind.go::substitute`'s existing `{artifact}` rule verbatim rather than inventing a second policy" — (0028:C1.6, §normative-contracts)
- [REQ-90] "No `--` is inserted: only a child that honours the separator would be helped, and the model cannot know which do." — (0028:C1.6, §normative-contracts) — NEGATIVE REQ.
- [REQ-91] "The refusal is per-USE-SITE like C1.2's `value shape:` — a tag bound on the context but named by no argv element of this entry is never scanned." — (0028:C1.6, §normative-contracts)
- [REQ-92] "why here: 0004:C12 read-back re-reads the SAME role; a shared artifact (an index) has one reader for many records, and without an identity in its argv that reader cannot say which row it read. The consumer's own projector stays the reader (no line-oriented READ carrier is introduced — Briefly Rejected)" — (0028:C1.6, §normative-contracts) — NEGATIVE REQ: no `edit` read carrier.
- [REQ-93] "no other change to 0025:C1–C1.6: stdin envelopes, exit maps, env overlay and the gate are untouched." — (0028:C1.6, §normative-contracts) — NEGATIVE REQ.
- [REQ-94] "The shell-interpreter deny-list stays STATIC — the two rules above are rules on this placeholder family, not entries on that list and not a widening of it, which is what keeps 0027:C1's promise (\"argv WORDS only\", two named admitted forms) true with no third form to disclose (JDR 0003 §D2 (a); 0027 unchanged)" — (0028:C1.6, §normative-contracts) — NEGATIVE REQ.

## Seam extension (A1) — the one interface change this RDR owns

- [REQ-95] "Context tags crossing the write seam | This RDR | Introduced | A1 — the one seam extension" — (0028 Capability Dependencies, §proposed-solution)
- [REQ-96] "context tags must reach `Apply` for `{tag.<key>}` anchors, which today's `Artifact{Role, Path}` cannot carry ⇒ this RDR owns that field, and its form is sharpened at Resolve." — (0028 Technical Design, §proposed-solution)
- [REQ-97] "the SAME `art` value reaches both carriers, so one context map on `Artifact` serves the writer's anchor tags and C1.6's reader" — (0028:A1 Evidence, §critical-assumptions)
- [REQ-98] "Context tags cross to `Apply` and `Read` (A1)" — (0028 Phase 3, §implementation-plan)

## Data flow / phase ordering (non-normative but load-bearing on structure)

- [REQ-99] "Data flow at apply: `Executor.Write` → non-owned check (unchanged) → resolve the role's reader (unchanged; C1.3's gate pre-check hangs here, A2) → `Apply`: read file → split lines (terminators kept) → for each rule: substitute `{tag.*}` into the anchor (quoted), compile, select exactly one line → for each rule: expand the parsed template segments (group text, planned value) → refuse or rewrite → stage + rename → executor read-back through the role's reader." — (0028 Technical Design, §proposed-solution)
- [REQ-100] "Extract `stageAndRename(dir, write, mode)` for both callers, or copy the discipline; do not extend `save`" — (0028 Existing Infrastructure Audit, §proposed-solution) — NEGATIVE REQ: `flowbind.go::save` is not extended (its `MkdirAll(0o700)` creation, `WriteJSONLine` store coupling and fixed `Chmod(0o600)` are each incompatible).
- [REQ-101] "Admit `edit` on write entries — `table.Accessor.Edit`, the `edit.<key>` tables, C1.4's categories in `carrierDefect`'s order, template and anchor parsed at load (C1, C1.2, C1.4); dump/normalize round-trip carries it (A6)." — (0028 Phase 1, §implementation-plan)
- [REQ-102] "`docs/cli-output-contract.md` names the new categories and Detail tokens" — (0028 Phase 4, §implementation-plan)
- [REQ-103] "the consumer's `rdr-write.toml` migration is the consumer's kata (rdr#yjye), not this plan's" — (0028 Phase 4, §implementation-plan) — NEGATIVE REQ: out of scope.

## REQ-MVV — the acceptance scenario (`0028:MVV`)

The MVV is ONE element with eight steps; each is quoted below as a sub-step of
REQ-MVV so a later stage can trace an end-to-end failure to its step.

- [REQ-MVV] "The consumer's acceptance scenario, carried from intrastate#zyh0." Fixture: "a `rdr-write.toml` re-authored as a state-machine over (a) one 0025 command reader `[\"rdr\",\"status\",\"-json\",\"-filter\",\"status\",\"{artifact}\"]` on role `record`, (b) an `edit` writer for the record's Status line on role `record`, (c) a C1.6 command reader over the README row on role `readme` (`{tag.nnnn}` in its argv) and an `edit` writer for that row anchored by `{tag.nnnn}`" — (0028:MVV, §implementation-plan)
  - [REQ-MVV.1] "`intrastate lint --model rdr-write.toml` passes; no wrapper script exists anywhere in the fixture."
  - [REQ-MVV.2] "ONE INVOCATION PER RECORD — `{tag.nnnn}` binds one record identity, so it addresses one record file and one README row; the three fixture records are three sequential runs of the pipeline below, not one." ... "`--tag` is bound on BOTH halves and this is load-bearing, not symmetry: `resolve` runs the narrowed readers before its kernel call, so the C1.6 readme reader's `{tag.nnnn}` must already be bound there or the pipeline refuses at `resolve`, before `set-state` is ever reached."
  - [REQ-MVV.3] "End state: the record's Status line reads `Final` (the joint-decision record reads `Final [joint decision → …]`, qualifier kept by backreference; the wrapped record keeps its continuation line byte-identical and its reader reports the qualifier whole); the README row's status cell reads `Final`; both confirmed by read-back through the command reader; every other byte of both files is unchanged." ... "the invariant is one line per file PER INVOCATION, and pinning a whole-run total would assert the fixture's size instead."
  - [REQ-MVV.4] "A `none` or `stopped:*` resolve row applies nothing and exits 0 carrying `dispositions`."
  - [REQ-MVV.5] "Negative: the same pipeline without `--allow-commands` on `set-state` refuses before mutation and both files are byte-identical to before."
  - [REQ-MVV.6] "Negative: a README with the record's row duplicated refuses `edit_anchor_ambiguous`; files untouched."
  - [REQ-MVV.7] "Negative: a `replace` whose output no longer matches its own anchor refuses `edit_anchor_unstable`; files untouched."
  - [REQ-MVV.8] "Negative, and required rather than optional: a README carrying one row in the live corpus's unlinked `| NNNN |` form (F6's 0001/0006 shape) refuses `edit_anchor_unmatched` for that record, names the rule in the Detail, and leaves both files byte-identical. The other rows in the same README still flip."
- [REQ-MVV-BLOCK] "BLOCKING the MVV (not the unit work): the consumer's row-addressed projector verb the C1.6 reader invokes must exist. The Illustrative Code marks it \"verb illustrative; the consumer's to add\" and no such verb ships today, so Phase 4's end-to-end proof cannot run until it does. Phases 1–3 do not depend on it" — (0028 Prerequisites, §implementation-plan)

## Test-surface REQs that no contract clause states alone

These are testable obligations the Testing Strategy states over and above the
clause text; each is cited to its scenario.

- [REQ-104] "exactly one line selected per in-scope file; no decoy (template comment, quoted example) scores ≥2." ... "The invariant is PER FILE and the cardinalities are descriptive, never asserted — the corpus grows with every RDR seeded ... so a test pinning a total is a scheduled false positive" — (0028 S8, §validation)
- [REQ-105] "A record whose Status VALUE sits on the continuation line. **Expected**: zero matches → `edit_anchor_unmatched`, not a rewrite." — (0028 S13, §validation)
- [REQ-106] "A bound tag value containing RE2 metacharacters, used in an anchor. **Expected**: regexp-quoted — matches literally, cannot alter the pattern's structure (C1.2)." — (0028 S16, §validation)
- [REQ-107] "The witness is the target's INODE, unchanged across the call: file content cannot separate \"did not write\" from \"wrote identical bytes\", but stage-and-rename always replaces the inode (A5), so an unchanged inode is the observable that proves absence of a write." — (0028 S20, §validation)
- [REQ-108] "Mode preservation on a 0644 target and on a 0755 target. **Expected**: `git diff --summary` empty in both cases." — (0028 S21, §validation)
- [REQ-109] "A symlinked target. **Expected**: the symlink survives and points at the new content" ... the unresolved-symlink contrast "is stated here as the reason the clause resolves first, NOT as an assertion this scenario runs" — (0028 S22, §validation)
- [REQ-110] "`Invocations()` after an `edit` apply, on an entry whose role's reader is FILE-backed, run with the gate OFF. **Expected**: counts `Apply` calls exactly as `flowbind.go::Writer` does (0004:C14); the write succeeds." — (0028 S23, §validation)
- [REQ-111] "the same value bound but named by no argv element of the entry; the same value reaching an `anchor` only." ... "bound-but-unreferenced is never scanned (per-use-site); anchor-only is ADMITTED — `-` is ordinary regexp-quoted line data there" — (0028 S24, §validation)
- [REQ-112] "Run it BOTH ways over `protectedKeys`: once with the reader declaring only the planned key (`protected` empty ...) and once with the reader declaring a second key (`protected` non-empty ...). C1.3's pre-check must refuse identically in both" — (0028 S25, §validation)
- [REQ-113] "At the envelope the refusal exits **2**, not 3 ... The code's SPELLING is Stage 8's to choose and no test may pin the string; what this scenario pins is the exit group, the discriminator and the `findings[]` carriage — an exit-3 assertion here is the regression this scenario exists to catch" — (0028 S27, §validation)
- [REQ-114] "Refusal reporting follows `docs/cli-output-contract.md`: a scalar failure carries `param`, an aggregate carries `findings[]`, regardless of runtime cardinality." — (0028 S27, §validation)
- [REQ-115] "A bound target that cannot be read — absent (ENOENT), a directory (EISDIR), and unreadable (EACCES). **Expected**: all three refuse BEFORE mutation with `execution_failure` carrying the OS error in the Detail" — (0028 S27b, §validation)
- [REQ-116] "One entry carrying TWO simultaneous load-time defects ... **Expected**: the earlier category in C1.4's registration order is the one reported, and only it" plus "TWO `edit.<key>` tables in one entry each carrying a defect of the SAME category. The assertion is that exactly ONE category is reported and it is that category — NOT which table is named." — (0028 S28, §validation)
- [REQ-117] "Two invocations against one artifact — a `clear = \"line\"` rule applying a planned `<clear>`, then a non-clear write through the same rule. **Expected**: the first deletes the line and read-back reports the key ABSENT; the second refuses `edit_anchor_unmatched`" — (0028 S29, §validation)
- [REQ-118] "An entry with NO carrier at all, exercised at runtime rather than at load. **Expected**: the registry builds a REFUSING binding, never a file binding" ... "`registry.go::commandBacked` is `len(acc.Command) != 0 || acc.Path == \"\"` today, and its `acc.Path == \"\"` arm IS the residue path, so adding the `edit` arm could silently reroute residue to a file binding." — (0028 S30, §validation)
- [REQ-119] "A decoy fixture where the anchor's exactly-one match is the WRONG line ... **Expected**: the write applies, and the read-back through the role's reader then refuses `read_back_mismatch`. `Applied()` is FALSE on this arm" ... "a test asserting `Applied()` TRUE here would fail for the correct reason and must not be relaxed into deleting the check." — (0028 S31, §validation)
- [REQ-120] "One entry carrying TWO simultaneous APPLY-time defects ... **Expected**: the earlier category in C1.3 `precedence:` is the one reported, and only it" plus "an entry whose rule A matches zero lines while rules B and C select the same line refuses `edit_anchor_unmatched`, not `edit_anchor_collision`, because every rule's cardinality resolves before the cross-rule sweep." — (0028 S32, §validation)
- [REQ-121] "Escapes and group semantics, per template. In `replace`: `$$` → `$`; `{{`/`}}` → literal braces; a group that did not participate expands empty; a bare `$1`. In `anchor`: `\\d{4}`, `a{2,3}`, a literal `{{`." ... "The two templates are NOT symmetric, which is the point of the scenario." — (0028 S14, §validation)
- [REQ-122] "Read-back of a planned `Final` against a record carrying a bracketed qualifier. **Expected**: byte-equal — the reader reports `status=Final` with the qualifier on `status_form`/`status.qualifier` (A4's normative fixture)." — (0028 S26, §validation)
- [REQ-123] "Done = every C1–C1.6 clause has a test, every refusal category fires on a fixture that earns it, and the MVV's seven steps pass end to end." — (0028 Testing Strategy, §validation) — note the MVV as written carries EIGHT steps; see ASSUMPTION A-7 below.

---

## ASSUMPTIONS

Implicit choices taken where the wording is imprecise but one reading is
defensible. Each is grounded against the record's own evidence base
(A1–A12, the fixtures they name, and the predecessor records 0025/0004)
before being taken.

- **A-1 — `edit` is admissible on write entries only, and `edit_tag_argv0` is
  the sole `edit_*` category reachable from a non-write entry.** C1.1 says
  `edit` on a read/gate entry is `edit_carrier_conflict`; C1.4 says
  `edit_tag_argv0` "fires on `command` entries of every kind". Taken: the
  loader evaluates `edit_tag_argv0` over every entry's `command` argv
  regardless of capability, and the other five `edit_*` categories only where
  an `edit` table is present. Grounded on C1.4 `precedence:`'s "the two sets
  are disjoint by carrier and never race" (REQ-73).

- **A-2 — "the six wire strings" of C1.4 append in the listed order; the
  deviations file's "five" is a stale count.** `deviations.md` D1/D2 each say
  "five" `edit_*` categories, while C1.4 lists six and its `registration:`
  clause and S6 both say six. Taken: SIX categories append, in C1.4's stated
  order. Grounded on the fence being the contract (C1.4) and on D1/D2 being
  pre-seeded 7.1 notes written before `edit_tag_argv0` was folded in at
  finalize (the Joint-check line records the 2026-09-03 extension that added
  it). This is a bookkeeping discrepancy in an artifact, not a contract
  ambiguity — recorded here rather than as a QUESTION.

- **A-3 — the `{tag.<key>}` anchor placeholder is scanned for the unescaped
  literal prefix `{tag.` only, and a malformed placeholder is a LINT defect,
  not a pass-through.** C1.2 states both halves in prose (REQ-25, REQ-27) but
  does not give the scanner's stopping rule for an unterminated `{tag.`.
  Taken: an unterminated `{tag.` (no closing `}` before end of pattern) is
  `edit_anchor_invalid` — "a `{tag.…}` placeholder is malformed" (REQ-26).
  Grounded on C1.2's "a placeholder is never passed through literally"
  posture, inherited from 0025:C2 (0025 REQ-18: an element carrying a `{...}`
  token without being exactly a known placeholder is a load-time defect,
  "never silently-literal text").

- **A-4 — the anchor's capture groups are the source of truth for
  `edit_template_invalid`'s "a group reference the anchor does not define".**
  C1.4 REQ-66 names the defect but not how group arity is known at load.
  Taken: the anchor is compiled at load (already required by REQ-65, with
  `{tag.<key>}` replaced by a quoted probe) and the compiled `*regexp.Regexp`'s
  `NumSubexp()`/`SubexpNames()` decide it. Grounded on C1.4's "checked with
  every `{tag.<key>}` replaced by a quoted probe" — the probe substitution
  cannot change group arity, since `regexp.QuoteMeta` emits no group syntax.

- **A-5 — "an entry's planned values" in C1.3 `precedence:` step (2) means the
  values planned for THIS entry's keys, scanned whether or not the rule that
  owns each key reached selection.** C1.3 words step (2) as "over the entry's
  planned values and the tag values its rules actually reference", giving the
  per-use-site qualifier only to the TAG half. Taken: planned values are
  scanned per-entry (every planned key of the entry), tag values per-use-site.
  Grounded on the asymmetry being explicit in the clause text and on C1.2
  `value shape:` naming "a planned value, or a bound tag value THIS ENTRY's
  rules reference" — the restriction attaches to the tag half only.

- **A-6 — the C1.6 argv0 rule is positional over the DECLARED argv, evaluated
  before any substitution.** C1.6 says the check is "statically decidable, so
  it is refused where it is visible rather than at spawn". Taken: `argv[0]`
  of the declared `command` vector is tested for being exactly a
  `{tag.<key>}` element at LOAD; no runtime argv0 check is added. Grounded on
  REQ-86/REQ-87 and on A9's finding that `cmdbind.go::substitute` compares
  whole elements at one site (so a runtime check would be a second policy the
  clause forbids).

- **A-7 — the MVV has EIGHT steps; Testing Strategy's "the MVV's seven steps"
  is a stale count.** The Implementation Plan's MVV enumerates 1–8, and step 8
  is marked "required rather than optional". Testing Strategy's Done line and
  the Acceptance paragraph both say "seven". Taken: all eight steps are the
  acceptance gate. Grounded on step 8's own text ("Without this step the MVV
  proves the composed scenario only against a fixture curated to exclude the
  one class of drift the real corpus has") — a step the record argues is
  load-bearing cannot be the one the count drops.

- **A-8 — `stageAndRename` staging happens in the resolved target's directory.**
  C1.3 `write:` says "staged beside the resolved target and renamed over it"
  and the audit names the helper `stageAndRename(dir, write, mode)`. Taken:
  the temp file is created in `filepath.Dir(resolvedTarget)` so the rename is
  same-filesystem and atomic. Grounded on A5's spike (`os.Rename` semantics)
  and on `flowbind.go::save`'s existing `CreateTemp`-beside-target discipline
  that C1.3 names as the reuse.

- **A-9 — a rule whose `<clear>` plan hits a zero-match anchor contributes no
  deletion to the re-anchor index arithmetic.** C1.5 REQ-78 makes it "SUCCESS
  with no write" and C1.3 `re-anchor:` shifts by "the deletions of preceding
  sibling rules". Taken: only rules that actually deleted a line count toward
  the shift. Grounded on A12's "the shift is a count, not a diff" and on
  `clear` being the only deletion.

- **A-10 — the no-op arm (REQ-49) is evaluated AFTER the re-anchor pass.**
  C1.3 orders `edit_anchor_unstable` "necessarily last, being post-rewrite"
  (REQ-54) and the no-op arm suppresses the write. Taken: re-anchor runs on
  the post-edit buffer first; only if it passes is the buffer compared to the
  input and the write skipped when equal. Grounded on `precedence:` step (6)
  being the last REFUSAL step, with the write (or its suppression) after it.

- **A-11 — "the exit-2 group" is the CLI-level exit code; the accessor-level
  class remains `execution_failure`.** C1.3 REQ-43 and S27 REQ-113 both state
  this, and A8/JDR 0003 §D1 forbid a new class. Taken: `Err` stays
  `execution_failure` with `Applied()` false at the accessor seam; the exit
  code is discriminated at `flow_exec.go::accessorFailureOf`. Grounded on A8's
  verified evidence and on REQ-42's "no new refusal class is introduced".

- **A-12 — the six `edit_*` wire strings are `edit_carrier_conflict`,
  `edit_key_mismatch`, `edit_anchor_invalid`, `edit_template_invalid`,
  `edit_clear_invalid`, `edit_tag_argv0`; the apply-time Detail reason tokens
  are `edit_anchor_unmatched`, `edit_anchor_ambiguous`, `edit_anchor_collision`,
  `edit_anchor_unstable`, `edit_value_multiline`, `edit_clear_undeclared` — two
  disjoint sets.** The record uses `edit_*` spellings on both sides of the
  load/apply line without labelling the split. Taken: only the six in C1.4
  register in `table.Categories()`; the six apply-time tokens ride the Detail
  and register nowhere. Grounded on C1.4 listing exactly six and on C1.3/C1.5
  naming their tokens as "Detail" content (REQ-41), plus F2's "a Detail naming
  `<id>.edit.<key>` and the reason token".

---

## QUESTIONS

Genuinely ambiguous clauses — two readings would produce materially different
behaviour and no predecessor precedent settles them. Recorded, not blocking.

- **Q1 — Does `replace`'s `{<key>}` placeholder have to appear at all?**
  C1.2 admits `{<key>}` as "the planned value of THE key this table is named
  for, and no other key", and C1.4's `edit_template_invalid` fires on "another
  key's placeholder", but no clause requires the OWNED placeholder to be
  PRESENT. Reading A: a `replace` with no `{<key>}` is legal (it writes a
  constant line), and read-back then decides whether the planned value landed
  — which for a non-matching constant refuses `read_back_mismatch` (S31's
  class). Reading B: a `replace` that never interpolates its own key is
  `edit_template_invalid` at lint, because the entry declares it writes that
  key. The two differ materially: B refuses at lint what A ships and lets fail
  at read-back. No predecessor resolves it — 0025:C2's vocabulary has no
  analogous "must use the placeholder" rule, and C1.4's `what lint proves:`
  list (REQ-74) says "every placeholder is closed over the entry's keys",
  which is a closure rule, not a presence rule. **Reading taken for the audit:
  A** (presence not required), because `what lint does NOT prove` (REQ-75)
  explicitly leaves value-correctness to apply-time read-back, and B would be
  a lint rule the closed list does not name. Flagged because B is the safer
  reading and a reviewer may expect it.

- **Q2 — On a `clear = "line"` rule whose planned value is NOT `<clear>`, does
  `clear` change anything?** C1.5 defines `clear` only as "the disposition of a
  planned `<clear>`". Reading A: `clear = "line"` is inert for a non-clear
  plan — the rule behaves exactly as if `clear` were absent, expanding
  `replace` normally. Reading B: declaring `clear` makes the rule
  deletion-capable and a non-clear write through it is still an ordinary
  rewrite, but C1.5's `one-way:` warning ("declare `clear` only where the LINE
  is the key") hints the two dispositions are not meant to coexist on one rule
  in practice. The behaviours differ only if B is read as forbidding the
  combination. **Reading taken: A** — `clear` is a per-rule disposition
  consulted only when the planned value IS the `<clear>` sentinel; S29
  (REQ-117) exercises exactly the sequence "clear, then non-clear write
  through the same rule" and expects `edit_anchor_unmatched`, not a load-time
  or apply-time rejection of the combination, which settles it in A's favour.
  Recorded because the `one-way:` prose reads as advice, not as a rule, and a
  reader could take it as normative.

---

## DEVIATION NOTES (pre-existing, see `deviations.md`)

- D1 (`TestReq77`'s tail-slice assertion) is REQ-76's direct consequence and is
  already OPEN in `deviations.md`. Phase 1 will trip it.
- D2 (per-category description text) depends on 0027 landing; carries the same
  stale "five" count as A-2 above.
- D3 (merged clause-3 behaviour of `carrierDefect`) is owed by whichever of
  0027/0028 lands second.
