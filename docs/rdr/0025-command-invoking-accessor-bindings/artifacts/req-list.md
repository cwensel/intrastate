# REQ list — RDR 0025 command-invoking-accessor-bindings

Phase 0 spec audit. Source: `docs/rdr/0025-command-invoking-accessor-bindings.md`
(1897 lines; C=7, MVV=1, A=15, F=10, G=1, D=4, S=7, BR=3, ALT=2).

Element ids carried where a REQ derives from a labelled contract. Quotes are
exact bytes from the record (via `rdr inspect --select`), reflowed only where a
clause spans lines in the source; no wording is changed.

---

## C1 — carrier and mutual exclusion (`0025:C1`)

- [REQ-1] "command       = [\"<argv0>\", \"<arg>\", ...]   # []string" — (0025:C1, §normative-contracts)
- [REQ-2] "output        = \"json\" | \"raw\"              # *string, read entries only; default \"json\" when omitted" — (0025:C1, §normative-contracts)
- [REQ-3] "exit_absent   = [<code>, ...]               # []int,            read entries only" — (0025:C1, §normative-contracts)
- [REQ-4] "exit_verdicts = { \"<code>\" = \"<verdict>\" }  # map[string]string, gate entries only" — (0025:C1, §normative-contracts)
- [REQ-5] "env           = { \"<KEY>\" = \"<value>\" }     # map[string]string" — (0025:C1, §normative-contracts)
- [REQ-6] "env_pass      = [\"<VAR>\", ...]              # []string" — (0025:C1, §normative-contracts)
- [REQ-7] "exactly one of `path` / `command` per entry. Load-time this is C5 `command_and_path_conflict`" — (0025:C1, §normative-contracts)
- [REQ-8] "at RUNTIME the constructor refuses a carrier-less entry with a refusing binding (`execution_failure`, Detail naming the malformed entry) and MUST NOT fall through to a `Path: \"\"` file binding, which would read every declared key as absent and confirm an unapplied write" — (0025:C1, §normative-contracts)
- [REQ-9] "`command` must be non-empty and contain no empty element." — (0025:C1 prose, §normative-contracts)
- [REQ-10] "`role`, `keys`, `timeout`, and `read_back` rules are unchanged from RDR 0002/0004." — (0025:C1 prose, §normative-contracts)
- [REQ-11] "The selection in `flowbind.Registry` must be **exhaustive, and fail closed on the residue**." ... "the constructor selects `command` first, `path` second, and **panics-free refuses** the third case" — (0025:C1 prose, §normative-contracts)
- [REQ-12] "The six fields above are the complete set this RDR adds to `internal/table/source.go::sourceAcc`, with the Go types named" — (0025:C1 prose, §normative-contracts)
- [REQ-13] "`output` is a pointer so that omitted and explicit-`\"json\"` are distinguishable at load, where the C5 `command_output_shape` arms are judged." — (0025:C1 prose, §normative-contracts)
- [REQ-14] "A command entry is now legal there, so the path-absent-and-command-absent case moves to `command_and_path_conflict`. Existing fixtures asserting the old category on a path-less entry change category, not verdict" — (0025:C1 prose, §normative-contracts)
- [REQ-15] "the Go field names and struct tags on `internal/table/source.go::sourceAcc` are unconstrained implementation choice, as is the decomposition of the binding into helpers and the package or file that holds it" — (0025:C1 prose, §normative-contracts) — NEGATIVE REQ: no test may pin Go field names, helper decomposition, or package layout.

## C2 — authority bound / placeholder vocabulary (`0025:C2`)

- [REQ-16] "{artifact}   # the closed placeholder vocabulary, v1 complete: replaced whole-element by the caller-bound artifact path for the entry's declared role" — (0025:C2, §normative-contracts)
- [REQ-17] "The executed argv is exactly the declared vector after whole-element placeholder substitution." — (0025:C2 prose, §normative-contracts)
- [REQ-18] "A placeholder is recognized only as a whole argv element. An element that contains a `{...}` token without being exactly a known placeholder is a load-time defect (C5) — never silently-literal text." — (0025:C2 prose, §normative-contracts)
- [REQ-19] "Execution invokes no shell and performs no other rewriting of the vector." — (0025:C2 prose, §normative-contracts)
- [REQ-20] "the **command binding** refuses the invocation (`execution_failure`, before spawn) when the caller-bound artifact path is relative or begins with `-`" — (0025:C2 prose, §normative-contracts)
- [REQ-21] "The check sites in the binding, not in `executor.go`" and it "**refuses rather than absolutizes**" — (0025:C2 prose, §normative-contracts)
- [REQ-22] "the refusal `Detail` says so (\"artifact path must be absolute for a command entry\")" — (0025:C2 prose, §normative-contracts)
- [REQ-23] "argv0 resolution is fixed at load, never cwd-relative: a bare name resolves through the parent's `PATH` at spawn (A2); a name containing a path separator resolves against the **model file's directory**" — (0025:C2 prose, §normative-contracts)
- [REQ-24] "an empty argv0 is a C5 defect, and the binding never restores implicit current-directory lookup" — (0025:C2 prose, §normative-contracts)
- [REQ-25] "the command-binding constructor receives `filepath.Dir` of the **absolutized** model path from the caller that opened the file, alongside the model." — (0025:C2 prose, §normative-contracts)
- [REQ-26] "This keeps `table` free of path resolution and leaves in-memory model construction untouched — no `table.Model` field, so no constructor or fixture changes." — (0025:C2 prose, §normative-contracts) — NEGATIVE REQ.
- [REQ-27] "`--model` is stored verbatim (it is not absolutized today), so the caller absolutizes before taking `Dir`; a relative `--model` must not make argv0 resolution cwd-dependent." — (0025:C2 prose, §normative-contracts)
- [REQ-28] "A separator-bearing argv0 in one [a model with no source file] is refused `execution_failure` at **binding construction**, before any spawn — not a load-time defect" ... "What it must never do is fall back to the process cwd" — (0025:C2 prose, §normative-contracts)

## C3 — invocation envelope (`0025:C3`)

- [REQ-29] "stdin (write): {\"<key>\": \"<value>\", ...}   # the planned tags; a planned `<clear>` crosses as the literal reserved value" — (0025:C3, §normative-contracts)
- [REQ-30] "stdin (read/gate): {}                       # nothing per-invocation beyond {artifact}; the object is always sent" — (0025:C3, §normative-contracts)
- [REQ-31] "stdout (read, output = \"json\", default): flat JSON object of strings; a declared key the object omits is UNREADABLE, never established-absent" — (0025:C3, §normative-contracts)
- [REQ-32] "stdout (read, output = \"raw\"):            the single declared key's value = stdout minus one trailing \"\\n\"; valid only when `keys` has exactly one entry" — (0025:C3, §normative-contracts)
- [REQ-33] "stdout (gate): {\"verdict\": \"allow\" | \"deny\" | \"indeterminate\", \"reason\": \"<text>\"}" — (0025:C3, §normative-contracts)
- [REQ-34] "exit_absent   = [<code>, ...]                # read entries: a listed exit with empty stdout establishes every declared key absent" — (0025:C3, §normative-contracts)
- [REQ-35] "empty stdout, no exit-map match: execution_failure, on BOTH read modes and on gate — never UNREADABLE and never established-absent." ... "a read whose exit is unlisted (exit 0 included) refuses" — (0025:C3, §normative-contracts)
- [REQ-36] "exit_verdicts = { \"<code>\" = \"allow\" | \"deny\" | \"indeterminate\", ... }   # gate entries: a listed exit with empty stdout is that verdict" — (0025:C3, §normative-contracts)
- [REQ-37] "Per-invocation data beyond the artifact path crosses only on stdin; read and gate results return only on stdout, in the shapes above." — (0025:C3 prose, §normative-contracts)
- [REQ-38] "A write command's success is never taken from its exit status alone — verification is read-back (`0004:C12`)." — (0025:C3 prose, §normative-contracts)
- [REQ-39] "the gate `verdict` strings are `internal/accessor/model.go::Verdict`'s (`0004:C9`)" — (0025:C3 prose, §normative-contracts)
- [REQ-40] "`<clear>` is carried unchanged as the planned value — the tool or wrapper performs the removal, and read-back verifies absence (`0004:C11`)" — (0025:C3 prose, §normative-contracts)
- [REQ-41] "the maps list **verdict** codes only — any unlisted exit, spawn failure, or malformed stdout is `execution_failure` (C4), so execution failure is never laundered into a verdict." — (0025:C3 prose, §normative-contracts)
- [REQ-42] "A spawn failure has no exit code at all (`exec.ErrNotFound`, fixture E1h), so no map entry can match it however the map is written — the maps are consulted only for a process that ran and exited." — (0025:C3 prose, §normative-contracts)
- [REQ-43] "A non-empty stdout is always parsed first (C4 ordering); the exit maps apply only to an empty stdout — including a stdout truncated at the 1 MiB cap, which is *non-empty* and therefore an `execution_failure` from the bound (C4), never an exit-map consultation." — (0025:C3 prose, §normative-contracts)
- [REQ-44] "absence reaches the executor as `KeyValue{Absent: true}` in `values`, never as omission from both slices, which `classify` reads as unreadable (C7)." — (0025:C3 prose, §normative-contracts)
- [REQ-45] "Under `output = \"raw\"` the single declared key is established from stdout, and an empty stdout with an unlisted exit is `execution_failure`, not an empty string." — (0025:C3 prose, §normative-contracts)
- [REQ-46] "The inherited reserved-literal rule composes unchanged **on the read path only**: a value that reads back as the literal `<clear>` is UNREADABLE" ... "`Executor.Read` passes `clearIsUnreadable = true`, while `Executor.Write`'s read-back passes `false`" ... "this RDR adds no arm and changes no site." — (0025:C3 prose, §normative-contracts) — NEGATIVE REQ on the executor side.
- [REQ-47] "The protocol's version rides out-of-band as `INTRASTATE_PROTOCOL` in the child env (C4 overlay), keeping the stdout map flat" — (0025:C3 prose, §normative-contracts)

## C4 — execution safety inheritance (`0025:C4`)

- [REQ-48] "order:    parse the stdout envelope, THEN classify the exit code" — (0025:C4, §normative-contracts)
- [REQ-49] "deadline: Setpgid; at ctx deadline Cancel = SIGKILL to -pgid; WaitDelay = 500ms bounds the stdin write and the pipe drain; timeout is classified from ctx.Err()" — (0025:C4, §normative-contracts)
- [REQ-50] "bounds:   stdout capped at 1 MiB (overflow = execution_failure); a NEW `Detail string` field on `accessor.Refusal` carries the last 4 KiB of stderr" — (0025:C4, §normative-contracts)
- [REQ-51] "env:      child env = allowlisted parent vars {PATH, HOME, TMPDIR, LANG, LC_*} + named `env_pass = [\"VAR\", ...]` vars + the entry's literal `env = { KEY = \"value\" }` + the overlay {INTRASTATE_ROLE, INTRASTATE_CAPABILITY, INTRASTATE_ACCESSOR, INTRASTATE_PROTOCOL=1}; nothing else is inherited." — (0025:C4, §normative-contracts)
- [REQ-52] "On a key collision the LATER layer wins, in exactly that order: overlay > entry `env` > `env_pass` > parent allowlist." — (0025:C4, §normative-contracts)
- [REQ-53] "`LC_*` is a literal prefix match on `LC_` (the one prefix rule; `env_pass` names whole variables and admits no pattern)." — (0025:C4, §normative-contracts)
- [REQ-54] "An `env` key or `env_pass` name matching the `INTRASTATE_` prefix is a C5 `command_env_conflict` defect at load, so the overlay is never shadowed silently" — (0025:C4, §normative-contracts)
- [REQ-55] "detail:   the stderr tail reaches the refusal as a typed `*accessor.ExecError` (Detail string) returned by the command binding through the existing `error` return; `executor.go::refusalOf` gains an error parameter and `errors.As`-es it." — (0025:C4, §normative-contracts)
- [REQ-56] "`executor.go::refusalWithKeys` wraps it and is what `Executor.Read` routes every read refusal through, so the error parameter threads through BOTH or the read path — the capability that binds directly — silently drops its tail." — (0025:C4, §normative-contracts)
- [REQ-57] "`Detail` is empty on every refusal carrying no error (timeout, read-back, gate-off) and set only where a binding returned one." — (0025:C4, §normative-contracts)
- [REQ-58] "`ExecError` is `struct { Detail string; Err error }` with `Error() string` and `Unwrap() error`: it WRAPS the offending `os/exec` error rather than flattening it, so `errors.Is(err, exec.ErrNotFound)` survives the trip to the refusal site and the argv0-not-found case stays distinguishable from a non-zero exit. `Detail` is the 4 KiB stderr tail, not `Err.Error()`" — (0025:C4, §normative-contracts)
- [REQ-59] "write:    read_back required; verified only through the role's reader, never by exit status" — (0025:C4, §normative-contracts)
- [REQ-60] "platform: a runtime `runtime.GOOS` check in the command binding refuses `execution_failure` before spawn on non-Unix — the predicate is `goos == \"windows\" || goos == \"js\" || goos == \"plan9\"` (refuse-listed, not allow-listed)" — (0025:C4, §normative-contracts)
- [REQ-61] "The check reads an injectable package-level `goos` var (defaulting to `runtime.GOOS`) so the refusal is testable on Unix CI" — (0025:C4, §normative-contracts)
- [REQ-62] "Command bindings run under RDR 0004's executor with its refusal *classes* unchanged and one additive type change — a `Detail string` field on `accessor.Refusal`" — (0025:C4 prose, §normative-contracts)
- [REQ-63] "`Reason` is not reused because `0004:C7` reserves it for a gate deny." — (0025:C4 prose, §normative-contracts) — NEGATIVE REQ.
- [REQ-64] "`readOutcome` drops the error today ... so it gains an `err` field for the read path to reach `refusalWithKeys`" — (0025:A11 evidence, cited by C4/Capability Dependencies) — additive `err` field on `readOutcome`.
- [REQ-65] "There is exactly **one** slot to land in — `internal/cli/clierr/clierr.go::CLIError` has a single `Detail string` ... they occupy that one slot in order: the applied-sense text **first** ..., the stderr tail appended after it." — (0025:C4 prose, §normative-contracts)
- [REQ-66] "The deadline mints `ClassTimeout`; spawn failure, a malformed or oversized stdout envelope, and — by default — a non-zero exit are `ClassExecutionFailure`, whose `Detail` carries the bounded stderr tail." — (0025:C4 prose, §normative-contracts)
- [REQ-67] "the stdout envelope is parsed **before** exit-code classification, so a well-formed deny envelope with a non-zero exit is a deny, not a failure. A non-zero exit is a gate verdict or an established absence **only** when the entry's C3 exit map lists that code and stdout is empty." — (0025:C4 prose, §normative-contracts)
- [REQ-68] "A write entry carries `read_back = true` (mandatory for every write since RDR 0004), its success is never taken from exit status, and it is verified through the role's declared reader — the re-read runs under its own bounded timeout (`0004:C15`), never the write's residue — with the applied-but-unverified sense (`0004:C14`) preserved when read-back cannot complete." — (0025:C4 prose, §normative-contracts)
- [REQ-69] "Path-backed bindings return untyped errors and are unaffected, so every production `NewExecutor` caller and every existing binding behave unchanged" — (0025:C4 prose, §normative-contracts) — NEGATIVE REQ.
- [REQ-70] "the binding builds on `exec.CommandContext` and carries no timer of its own (A1)." — (§technical-design, Binding family)

## C5 — static validation / load-time defect classes (`0025:C5`)

- [REQ-71] "command_and_path_conflict     # both or neither of path/command declared" — (0025:C5, §normative-contracts)
- [REQ-72] "command_empty                 # empty vector or empty argv element" — (0025:C5, §normative-contracts)
- [REQ-73] "command_unknown_placeholder   # unknown or non-whole-element {…} token" — (0025:C5, §normative-contracts)
- [REQ-74] "command_shell_interpreter     # argv0 + inline-code flag (sh -c, bash -c, python -c, env chains); no opt-in in v1" — (0025:C5, §normative-contracts)
- [REQ-75] "command_output_shape          # output = \"raw\" with keys ≠ 1; exit_absent on a non-read; exit_verdicts on a non-gate or naming a non-verdict" — (0025:C5, §normative-contracts)
- [REQ-76] "command_env_conflict          # an `env` key or `env_pass` name matching the reserved `INTRASTATE_` prefix (C4 overlay)" — (0025:C5, §normative-contracts)
- [REQ-77] "registration: all six are appended to `table.Categories()`, whose hand-maintained list IS the closed set; appended in the order declared above, after the existing members, so a consumer enumerating the list sees additions only at the tail." — (0025:C5, §normative-contracts)
- [REQ-78] "Each gets a typed `Category` constant of the existing `Cat…` form beside the others; the wire STRINGS above are the contract and the constant identifiers are not." — (0025:C5, §normative-contracts)
- [REQ-79] "The list's total size is not a contract at any point — it is append-only, and no clause or test may assert a count" — (0025:C5, §normative-contracts) — NEGATIVE REQ.
- [REQ-80] "precedence: WITHIN one entry — load is fail-fast (`0002:C24` — one categorized error for the whole document, never a list), so an entry carrying several of these defects reports the first in the order declared above." — (0025:C5, §normative-contracts)
- [REQ-81] "ACROSS entries in the SAME table there is no order: `accessorTable` ranges a Go map, so which of two defective entries is reported is unspecified, and no test may assert it." — (0025:C5, §normative-contracts) — NEGATIVE REQ.
- [REQ-82] "ACROSS tables the order IS fixed — `load.go::loadAccessors` runs read, then write, then gate, returning on the first error ... a test may rely on it only as 0002's contract" — (0025:C5, §normative-contracts)
- [REQ-83] "interpreter set: OPEN (deny-listed, not closed) — an unlisted interpreter is admitted, so the list grows by amendment" — (0025:C5, §normative-contracts)
- [REQ-84] "these are `internal/table/category.go::Category` constants with the spellings above ... not `accessor.ValidationCode`, whose eight-member closure is RDR 0004's own test contract." — (0025:C5 prose, §normative-contracts)
- [REQ-85] "A `Category` constant is not in the closed set until it is appended to `table.Categories()` ... so S1 must assert membership, not merely the refusal." — (0025:C5 prose, §normative-contracts)
- [REQ-86] "the defect's report says so (\"inline shell is not a declared command; put it in a script and declare the script as argv0\")" — (0025:C5 prose, §normative-contracts)
- [REQ-87] "the loader decodes strictly (`source.go::decodeStrict` sets `DisallowUnknownFields`), so a field absent from the struct is a **document-level** `CatUnknownSchemaField`, not one of C5's per-entry categories." — (0025:C1 prose, bears on C5, §normative-contracts)

## C6 — execution gate (`0025:C6`)

- [REQ-88] "--allow-commands        # v1's ONLY opt-in: a per-invocation flag; never the model file; `lint` does not carry it" — (0025:C6, §normative-contracts)
- [REQ-89] "registration: ONE registration, on the `flow` GROUP — `internal/cli/flow.go::newFlowCmd`'s `PersistentFlags()`, whose subtree is exactly the four `buildRequest` call sites across three verbs (`flow_next.go`, `flow_resolve.go`, `flow_state.go` ×2) and nothing else." — (0025:C6, §normative-contracts)
- [REQ-90] "`lint` sits at ROOT, outside the group, so it does not carry the flag" — (0025:C6, §normative-contracts)
- [REQ-91] "lookup: `buildRequest` reads the flag with the error CHECKED, not discarded. Because persistent registration guarantees presence, a lookup miss can only be a wiring bug and MUST surface as one" — (0025:C6, §normative-contracts)
- [REQ-92] "absent ⇒ every command invocation refuses execution_failure before spawn, Detail naming the gate; lint (C1/C5) validates regardless" — (0025:C6, §normative-contracts)
- [REQ-93] "gate site: ONE — `buildRequest` reads the flag and passes it to `flowbind.Registry`, the single production construction site (`flow_exec.go:830`), which builds refusing command bindings when it is unset." — (0025:C6, §normative-contracts)
- [REQ-94] "signature: ... It gains BOTH new inputs and no others, in this order: `func Registry(m *table.Model, baseDir string, allowCommands bool) accessor.Registry`." — (0025:C6, §normative-contracts)
- [REQ-95] "the return type is unchanged and `Registry` does not gain an error return — a carrier-less or gate-refused entry yields a refusing binding (C1), not a construction failure" — (0025:C6, §normative-contracts)
- [REQ-96] "The gate sites **in the command binding's constructor**, not in the executor" — (0025:C6 prose, §normative-contracts)
- [REQ-97] "With the gate off, a command invocation refuses before spawn (`execution_failure`, `Detail` naming `allow_commands`), while `intrastate lint` still validates the entries" — (0025:C6 prose, §normative-contracts)

## C7 — the implemented seam (`0025:C7`)

- [REQ-98] "the three command bindings implement `internal/accessor/binding.go`'s interfaces UNCHANGED — this RDR adds no method, changes no signature" — (0025:C7, §normative-contracts) — NEGATIVE REQ.
- [REQ-99] "Read(ctx context.Context, art Artifact, requested []string) (values []KeyValue, unreadable []string, err error)" — (0025:C7, §normative-contracts)
- [REQ-100] "Gate(ctx context.Context, art Artifact) (Verdict, string, error)" — (0025:C7, §normative-contracts)
- [REQ-101] "Apply(ctx context.Context, art Artifact, planned []resolve.Tag) error" — (0025:C7, §normative-contracts)
- [REQ-102] "Invocations() int          # on WriteBinding, beside Apply" — (0025:C7, §normative-contracts)
- [REQ-103] "the command write binding counts `Apply` ENTRIES — one per call — and performs NO internal respawn: one `Apply`, one spawn, one increment." — (0025:C7, §normative-contracts)
- [REQ-104] "the write method is `Apply`, NOT `Write`. `Gate` returns verdict + reason + error, so C3's `reason` field crosses back through the Go return and is not dropped." — (0025:C7, §normative-contracts)
- [REQ-105] "C3's read envelope maps onto `Read`'s SPLIT return: a parsed key becomes a `KeyValue` in `values`, an omitted declared key a name in `unreadable` — UNREADABLE is that second slice, never an error and never an absent value." — (0025:C7, §normative-contracts)
- [REQ-106] "`exit_absent` absence crosses as `KeyValue{Key: k, Absent: true}` IN `values` — a present record carrying the flag, the same channel `internal/cli/flowbind/flowbind.go::Reader.Read` already uses for the path-backed binding." — (0025:C7, §normative-contracts)
- [REQ-107] "Returning the key in NEITHER slice does NOT establish absence ... so an omitted key refuses the whole read `incomplete_read`. ... A command binding that signalled absence by omission would be refused, not believed" — (0025:C7, §normative-contracts)
- [REQ-108] "The two boundaries are distinct and both hold: at the BINDING return absence is the flag; at the RESOLVER seam it becomes omission, converted by `internal/accessor/model.go::ReadResult.OwnedSnapshot`, which is what `0004:C8` governs." — (0025:C7, §normative-contracts)

## MVV (`0025:MVV`)

- [REQ-MVV] `0025:MVV` — §minimum-viable-validation, lines 1472–1501. Five numbered steps plus the end-state paragraph, quoted in full:

  1. "Author a model declaring one command-backed read and one command-backed write (`read_back = true`) over a single artifact role, delegating to an established tool present in CI: the read binds `git config --file {artifact} --get <key>` directly in raw mode; the write is a declared wrapper over `git config` (the A3 spike's `wrapper-write.sh` shape)."
  2. "`intrastate lint` accepts it; six mutated copies (path+command conflict, empty element, unknown placeholder, `[\"sh\", \"-c\", …]` interpreter form, `output = \"raw\"` with two keys, an `env` key shadowing `INTRASTATE_ROLE`) are each rejected with their C5 defect."
  3. "With `--allow-commands` on the invocation (C6; without it, the same invocation refuses `execution_failure` naming the gate, before spawn), drive a state change end to end through the existing write path (`internal/cli/flow_state.go`'s set-state verb, whose planned tags come from its `--write` flags): the declared write **command** — not intrastate — applies the edit to the artifact. Automatic resolve→apply wiring is not this RDR's scope and no Phase builds it (see Problem Statement)"
  4. "Read-back runs through the declared command reader and verifies the written value; the result reports the written tags."
  5. "A command that sleeps past its declared timeout refuses `timeout` and leaves no orphan process."

  End state: "a linted model applies and verifies a real state change through a declared command, with the executed argv readable in the model."

## Load-Bearing Decisions (`0025:D-*`)

- [REQ-109] "the accessor identity stays RDR 0004's `(flow, name, capability)` triple; whether an entry is path- or command-backed does not enter identity." — (0025:D-identity, §load-bearing-decisions)
- [REQ-110] "The read envelope is deliberately the same flat map-of-strings, with the same presence-is-the-answer rule, that `flowbind.go::store` already uses ... They stay separate decoders because they decode different things" — (0025:D-wire-byte-format, §load-bearing-decisions)
- [REQ-111] "Rejected: carrying the requested key set on stdin (a list is not a string value, and the keys are already declared in the model)." — (0025:D-wire-byte-format, §load-bearing-decisions) — NEGATIVE REQ.
- [REQ-112] "the field is `command`; the family is \"command-backed accessors\". Rejected: `exec` ..., `run` ..., `argv` ..." — (0025:D-naming, §load-bearing-decisions)
- [REQ-113] "the registry selects the binding constructor by carrier field: `path` → file binding (today's `flowbind`), `command` → command binding. C1's exactly-one rule makes the selection total; no precedence order exists to get wrong." — (0025:D-selection-predicate, §load-bearing-decisions)
- [REQ-114] "The file binding keeps its magic-suffix vocabulary (`verdictFor`, `unreachable`) **unchanged**" — (0025:D-selection-predicate, §load-bearing-decisions) — NEGATIVE REQ.

## Implementation phases (§implementation-plan)

- [REQ-115] "Phase 1: Carrier — Extend `internal/table` (`sourceAcc`, `Accessor`, `accessorTable`) with `command` and the C1/C5 load rules." — (§implementation-plan)
- [REQ-116] "Phase 2: Binding family — Implement command-backed `ReadBinding`/`GateBinding`/`WriteBinding` over `exec.CommandContext` honoring C2–C4 (a sibling package to `flowbind`)." — (§implementation-plan)
- [REQ-117] "Phase 3: Selection — Discriminate the constructor by carrier field at `internal/cli/flowbind/registry.go::Registry` ... Register `--allow-commands` on the `flow` group's `PersistentFlags()` and thread it (C6) from `buildRequest` — reading it with the lookup error checked — in the SAME phase" — (§implementation-plan)
- [REQ-118] "Phase 4: Surface and proof — Wire lint reporting for the C5 defects and land the MVV scenario plus the timeout and read-back failure scenarios as tests." — (§implementation-plan)
- [REQ-119] "Illustrative — intent only; tests must not assert it literally." — (§illustrative-code) — NEGATIVE REQ.
- [REQ-120] "0025 lands after 0016 or inherits first-match until it does (Prerequisites)" — (0025:A4 / §prerequisites) — 0016 is Draft; this build inherits first-match `readerFor`.

## Testing Strategy scenarios (`0025:S1`–`S7`, §testing-strategy)

- [REQ-121] "Coverage goal: every C1–C6 clause has a test that fails when its rule is dropped; the MVV scenario is the integration proof." — (§testing-strategy)
- [REQ-122] S1 — "load-time validation (C1/C5) — a valid command entry plus the six MVV mutants ... Expected: the valid entry loads; each mutant is rejected with its own named C5 `table.Category` — asserted by category, never by \"validation returned non-empty\"; all six categories are members of `table.Categories()`, at the tail and in clause order ...; `accessor.ValidationCodes` stays at eight (REQ-84). Plus the **open-deny-list negative**: an entry whose argv0 is an unlisted interpreter (`perl -e`) **loads**" — (0025:S1, §testing-strategy)
- [REQ-123] S2 — "substitution guard (C2) — `{artifact}` bound to an absolute path, a relative path, and a `-`-prefixed path; an element embedding the token mid-string. Expected: the absolute path is substituted whole-element and the argv the child observes equals the declared vector otherwise byte-for-byte; the relative and `-` cases refuse `execution_failure` before spawn; the mid-string case is a C5 defect at load." — (0025:S2, §testing-strategy)
- [REQ-124] S3 — "envelope-before-exit ordering (C3/C4)" — nine cases, "Expected, in the order listed: deny; `execution_failure`; `execution_failure`; deny; `execution_failure`; established-absent (asserted as `KeyValue{Absent: true}` in `values` for every declared key, NOT as omission ...); `execution_failure`; `execution_failure`; `execution_failure`." Empty-stdout arms follow normative fixture **FX-exit-codes**; non-empty-stdout arms are contract-derived and "have **no fixture oracle and need none**". "Each refusal carries the 0004 diagnosis tuple and, per C4, the 4 KiB stderr tail in `Refusal.Detail` — asserted non-empty on case 3" — (0025:S3, §testing-strategy)
- [REQ-125] S4 — "deadline (C4, A1) — a child sleeping past `timeout`, a wrapper whose grandchild holds the stdout pipe, and a child that never reads stdin (1 MiB stdin) ... Expected: `timeout` refusal within `timeout + WaitDelay` in all three, no hang, and no process from the child's group surviving the refusal — normative fixture **FX-deadline**. The two ablations are **not test arms**" — (0025:S4, §testing-strategy) — NEGATIVE REQ: no injection seam that can disable the deadline triple.
- [REQ-126] S4b — "raw read and env policy (C3/C4, A3/A7) ... Expected: the value is stdout minus one trailing newline — normative fixture **FX-raw-read** ...; the injection never reaches the child ... Plus the composition arms C4 now orders: a variable named in `env_pass` reaches the child; an unnamed one does not; a key set in both the parent allowlist and the entry's `env` arrives with the entry's value; an entry `env` key shadowing `INTRASTATE_ROLE` is a C5 defect at load; `LC_ALL` passes and `LCFOO` does not ... Asserted on the child's observed environment, not on the read result" — (0025:S4b, §testing-strategy)
- [REQ-127] S5 — "write read-back (C3/C4) ... Expected: success reporting the written tags; `read_back_mismatch`; `read_back_incomplete` (the corrupted artifact no longer parses for the reader); the key reads back absent — the write's exit status decides none of them. Plus the inherited reserved-literal arm (C3) ... Asserted on the **read** path, where `clearIsUnreadable` is true; the write's read-back passes false ... the two senses are asserted separately and neither test stands in for the other." — (0025:S5, §testing-strategy)
- [REQ-128] S5b — "read-back reader selection under two readers on one role (A4, Prerequisites) ... Expected: two arms. (a) The owned-key arm: a second reader serving a command write's owned key is REFUSED at load (`malformed_accessor_binding`) ... (b) The disjoint-keys arm: pending 0016's fail-closed uniqueness, the selection is `registry.go`'s name-sorted first match — asserted **by the selected reader's identity**, never by \"read-back succeeded\" ... arm (b) asserts **both halves**: that the registry emits readers name-sorted, and that `readerFor` returns the first." — (0025:S5b, §testing-strategy)
- [REQ-129] S6 — "the MVV end to end against an established tool present in CI. Expected: the declared write command, not intrastate, edits the artifact, and the declared command reader verifies it." — (0025:S6, §testing-strategy)
- [REQ-130] S6b — "platform refusal (C4 `platform`, F7) — a valid command entry invoked with the binding's `goos` set to `windows`. Expected: `execution_failure` naming the unsupported platform, with no child spawned (asserted by absence of a spawn ...); the same model passes `intrastate lint` under that setting" — (0025:S6b, §testing-strategy)
- [REQ-131] S7 — "execution gate (C6) — a valid command model invoked without `--allow-commands`; with it; `intrastate lint` under both. ... the scenario derives the set of verbs reaching `buildRequest` and asserts each RESOLVES the flag through the group's persistent set — a structural `Flags().Lookup(\"allow-commands\") != nil` walk, never a text match on pflag's error. ... The complementary arm asserts `lint`, at root, does NOT resolve the flag ... Plus the C1 residue arm: an entry with neither carrier builds a **refusing** binding, never a `Path: \"\"` file binding. Expected: refusal `execution_failure` naming `allow_commands` with no child process spawned — asserted by **absence of a spawn** (a sentinel argv0 that would leave an observable trace if executed), not merely a non-zero exit; normal execution under the flag; lint passes in both ... The oracle derives both sets rather than restating a literal four." — (0025:S7, §testing-strategy)

## Failure Modes (`0025:F1`–`F10`, §failure-modes) — observable-behaviour REQs

- [REQ-132] "load-time C5 defects name the entry and the offending element; runtime refusals carry the 0004 diagnosis tuple (accessor, capability, role, timeout) with class `timeout` or `execution_failure`." — (§failure-modes)
- [REQ-133] "a gate command with a malformed envelope refuses `execution_failure` — never a deny the model didn't decide." — (§failure-modes)
- [REQ-134] "a command entry invoked on Windows refuses `execution_failure` naming the unsupported platform, before spawn (C4 `platform`) ... lint stays platform-neutral, so a model authored on Windows still validates." — (§failure-modes)
- [REQ-135] "a relative `--artifact` path that works for a path-backed entry refuses `execution_failure` for a command entry ... Diagnosis: the refusal `Detail` says the artifact path must be absolute for a command entry." — (§failure-modes)
- [REQ-136] "It is not statically detectable [the stdin-sinking write tool]: whether a tool reads stdin is not visible in its argv, so no C5 arm can catch it and none is claimed" — (§failure-modes) — NEGATIVE REQ: do not build a C5 arm for it; a declared `stdin = "none" | "envelope"` field is charted, not built.
- [REQ-137] "bindings hold no state between invocations (`0004:C14` — no retry, no undo)." — (§failure-modes)

## Cross-Cutting Concerns (`0025:G-cross-cutting`)

- [REQ-138] "credentials never appear in `command`, in `env` literals, or in stdin values ... A lint advisory on credential-shaped argv elements is **not** in v1" — (0025:G-cross-cutting) — NEGATIVE REQ (policy; no v1 lint arm).
- [REQ-139] "Every existing model keeps working untouched, and no path-backed entry changes meaning." — (0025:G-cross-cutting) — NEGATIVE REQ.
- [REQ-140] "The `raw` mode is defined bytewise: the value is the child's stdout minus **exactly one** trailing `\n` — no trimming, no case folding, no whitespace normalization, so a value with meaningful leading or internal whitespace survives." — (0025:G-cross-cutting)
- [REQ-141] "This RDR claims **no** byte-identical output, content-addressed identity, or replay-stable hash, so the hash/pre-image checklist does not apply." — (0025:G-cross-cutting) — NEGATIVE REQ.
- [REQ-142] "The placeholder vocabulary is versioned as a closed set: v1 is `{artifact}` and is complete (C2)." — (0025:G-cross-cutting)

---

## ASSUMPTIONS

Implicit choices made where wording was imprecise but a single reading is
defensible.

- ASSUMPTION: REQ-1 — `command` decodes as a TOML array of strings on the same
  entry table as `path`/`role`/`keys`/`timeout`; no nested table form is
  admitted. C1's grammar block shows only the flat spelling.
- ASSUMPTION: REQ-2/REQ-13 — `output` accepts exactly the two literals `"json"`
  and `"raw"`; any other string is `command_output_shape` (C5), by the same
  reading that makes `exit_verdicts` "naming a non-verdict" that category.
- ASSUMPTION: REQ-3/REQ-4 — `exit_absent`/`exit_verdicts`/`output` declared on a
  command-less (path-backed) entry is also `command_output_shape`; the clause
  says "read entries only"/"gate entries only", and C5's arm names capability
  mismatch, so the capability test is the discriminator, not carrier presence.
- ASSUMPTION: REQ-8/REQ-95 — the "refusing binding" is a per-capability type
  implementing the same three interfaces, returning `execution_failure`-shaped
  errors (`*accessor.ExecError` so the `Detail` reaches the refusal) on every
  method. C1 fixes the behaviour, not the type name.
- ASSUMPTION: REQ-20 — "begins with `-`" is tested on the substituted artifact
  path value only, not on other declared argv elements (a declared `--flag` is
  legitimate and unaffected).
- ASSUMPTION: REQ-32 — "minus one trailing `\n`" removes at most one `\n`, and
  only if present; a stdout with no trailing newline is used verbatim, and a
  stdout ending `\n\n` keeps one. G-cross-cutting's "exactly one" pins this.
- ASSUMPTION: REQ-34/REQ-106 — under `exit_absent` every key in the entry's
  declared `keys` is returned as `KeyValue{Key: k, Absent: true}`; the clause
  says "every declared key absent" and C7 fixes the channel.
- ASSUMPTION: REQ-50 — the 1 MiB stdout cap is enforced by bounding the read
  (reading up to cap+1 bytes and refusing on overflow), so an unbounded child
  cannot exhaust memory; "overflow = execution_failure" implies detection, not
  silent truncation.
- ASSUMPTION: REQ-51 — a parent var named in the allowlist but unset in the
  parent env is simply not passed (no empty-string synthesis).
- ASSUMPTION: REQ-51 — `env_pass` naming a variable unset in the parent env
  passes nothing; it is not a load defect.
- ASSUMPTION: REQ-54 — the `INTRASTATE_` prefix check is case-sensitive and
  literal, matching the `LC_` prefix rule's stated literalness.
- ASSUMPTION: REQ-58 — `ExecError` lives in package `accessor`
  (`*accessor.ExecError` is the spelling C4 uses), so the command binding
  imports `internal/accessor` to construct it — the package already being
  imported by `flowbind`.
- ASSUMPTION: REQ-60/REQ-61 — the `goos` package-level var lives in the command
  binding package and is unexported; tests in that package set it directly.
  Exporting it would put a runtime weakening seam in the public surface.
- ASSUMPTION: REQ-77 — "in clause order" is the C5 block's declared order:
  conflict, empty, unknown_placeholder, shell_interpreter, output_shape,
  env_conflict.
- ASSUMPTION: REQ-89/REQ-94 — `baseDir` is passed to `Registry` even when no
  command entry exists (the signature is unconditional); a path-backed-only
  model ignores it.
- ASSUMPTION: REQ-91 — "surface as one [a wiring bug]" means `buildRequest`
  returns an error (its existing error return), not a panic; the repo has no
  panic-on-wiring-bug idiom.
- ASSUMPTION: REQ-103 — `Invocations()` counts `Apply` entries even when the
  spawn fails, since the clause says "one `Apply`, one spawn, one increment"
  and the existing assertions are about layer re-entry, not process success.
- ASSUMPTION: REQ-116 — the sibling package name is implementation choice
  (C1 explicitly leaves "the package or file that holds it" unconstrained);
  the audit does not pin one.
- ASSUMPTION: REQ-30 — stdin is closed after the empty object is written, so a
  child that reads to EOF does not block.
- ASSUMPTION: REQ-124 — "case 3" in S3's Detail assertion is the third listed
  case (malformed envelope with exit zero → `execution_failure`), counting the
  Expected list in the stated order.

---

## QUESTIONS

Genuinely ambiguous clauses — two readings would produce materially different
behaviour and no in-record precedent settles them. Recorded, not blocking.

- Q1 — REQ-51 / REQ-52 vs REQ-42: `PATH` is in the C4 allowlist and A2 says
  `argv0` "resolves through the parent's `PATH` before spawn". If an entry's
  `env` or `env_pass` sets a different `PATH`, does `LookPath` use the parent's
  `PATH` (A2's stated mechanism, so the composed child `PATH` affects only the
  child's own subprocesses) or the composed child `PATH`? The record states
  the parent-`PATH` mechanism as a fact but never says whether a declared
  `PATH` override is honoured for resolution. Reading taken for the audit: the
  parent's `PATH` resolves argv0 (A2/C4 prose "argv0 still resolves through
  the parent's `PATH` before spawn (A2), which the allowlist passes on
  unchanged"), and the composed env is what the child receives. Flagged
  because the opposite reading changes which binary runs.

- Q2 — REQ-31 vs REQ-35: under `output = "json"`, a well-formed but EMPTY JSON
  object `{}` on stdout is non-empty stdout (so not exit-map territory, C3/REQ-43)
  and every declared key is omitted ⇒ all keys UNREADABLE ⇒ `incomplete_read`.
  A `{}` body is also two bytes, so "empty stdout" cannot mean it. Reading
  taken: `{}` parses, all declared keys are unreadable, the read refuses
  `incomplete_read`, not `execution_failure`. The alternative reading (a
  no-keys envelope is a broken tool ⇒ `execution_failure`) is defensible from
  "A silent tool is a broken tool" but contradicts the parse-first ordering.

- Q3 — REQ-59 / REQ-68: C4 `write:` says "read_back required". Is a write
  command entry declaring `read_back = false` (or omitting it) a LOAD-TIME
  defect, and if so under which category? C5 registers no arm for it, and the
  record says read_back is "mandatory for every write since RDR 0004" — i.e.
  inherited from 0004, not established here. Reading taken: the 0004 rule
  already governs and this RDR adds no new category; if 0004's loader does not
  enforce it for command entries the gap is a deviation, not a new C5 arm.

- Q4 — REQ-89 / REQ-131: C6 pins registration to `newFlowCmd`'s
  `PersistentFlags()` and S7 asserts containment structurally. Checked against
  `main`: `newFlowCmd` adds exactly four verbs (`next`, `resolve`,
  `read-state`, `set-state`) and `buildRequest` has four call sites, so today
  the subtree IS the caller set. The ambiguity is forward-looking: a future
  in-group verb that does not call `buildRequest` still inherits the flag and
  passes S7's "every in-group verb resolves it", while the flag is a false
  affordance on it — the exact objection C6 raises against `lint`. No clause
  requires the converse assertion. Reading taken: assert containment as S7
  words it; do not add a stricter "every in-group verb calls buildRequest"
  arm.
