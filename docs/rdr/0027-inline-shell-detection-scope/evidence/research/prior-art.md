Model: claude-fable-5

# Prior art — inline-shell detection scope on command bindings (Stage 2 read, not a spike)

Problem class: a declarative config carries a fixed argv; a lint tries to
tell, from argv alone, whether that argv hands inline code to a shell or
interpreter. Instance question: what does each named peer do for THIS
operator — does any peer DETECT shell-ness from argv, or is it declared /
always-on / never-on?

## Accepted citations (load-bearing; quoted from source)

### Go `os/exec` — the runtime the command carrier is built on

`$(go env GOROOT)/src/os/exec/exec.go` package doc:

> Unlike the "system" library call from C and other languages, the
> os/exec package intentionally does not invoke the system shell and
> does not expand any glob patterns or handle other expansions,
> pipelines, or redirections typically done by shells.

⇒ From a fixed argv, a shell runs ONLY if some argv word names one (directly
or through a wrapper that execs its trailing words). Redirection (`sh
<script`) is not a reachable form: without a shell `<script` is a literal
filename argument. So the forms an argv-level check can ever see are argv
WORDS; a shell string carried inside one word (`env -S "sh -c …"`) and a
stdin-fed shell (`sh -s`) are a different axis by construction.

### gh-cli — shell-ness is DECLARED, never detected

`gh-cli/pkg/cmd/alias/set/set.go::NewCmdSet`:

> If the expansion starts with `!` or if `--shell` was given, the expansion
> is a shell expression that will be evaluated through the `sh` interpreter
> when the alias is invoked.

    cmd.Flags().BoolVarP(&opts.IsShell, "shell", "s", false, "Declare an alias to be passed through a shell interpreter")

⇒ The only peer with BOTH modes marks the shell mode by author declaration
(`!` prefix / `--shell`) and never inspects the expansion for `sh -c`.
Detection-by-inspection is not the peer pattern; declaration is.

### roborev — hooks are ALWAYS a shell string

`roborev/internal/daemon/hooks.go::(*HookRunner).runHook`:

    cmd = exec.Command("sh", "-c", command)

### beads — credential-helper commands are ALWAYS a shell string

`beads/internal/creds/command.go::(CommandSource).Resolve`:

    cmd = exec.CommandContext(ctx, "sh", "-c", command)

⇒ Two mid-size Go CLIs that take a user-configured command choose
always-shell and say so in the config's shape (a string, not a vector).
Neither attempts to classify. Across the peer set there is NO instance of
"argv vector + a detector that refuses shell forms" — intrastate's
deny-list is its own construction, so its promise has to be bounded by its
own predicate, not borrowed from a peer.

### RDR 0004 — the two rejected authority shapes this check sits between

`0004:ALT2` (Raw Shell-Out Accessors), reason for rejection:

> makes static review nearly impossible — the shell script becomes the
> real contract.

`0004:ALT3` (Global Allowlisted Commands), reason for rejection:

> an executable allowlist does not prove a command has read, gate, or
> write authority over the intended artifact role … It constrains
> mechanism but not semantic power.

⇒ Name-based allow/deny lists are already on record as unable to be the
authority; 0025:C5 accepted that and scoped the deny-list as defense in
depth. This RDR must not re-promote it to a barrier (Alt 2's wrapper table
would), and must not delete the visibility it does give (Alt 3 would).

### RDR 0025 — what C5 promised and what its prose says the promise is

`0025:C5`:

> command_shell_interpreter     # argv0 + inline-code flag (sh -c, bash -c, python -c, env chains); no opt-in in v1
> interpreter set: OPEN (deny-listed, not closed) — an unlisted interpreter is admitted, so the list grows by amendment

`0025:§normative-contracts` prose beside C5:

> the check raises the cost of an inline-shell carrier without claiming to
> make one impossible. The guarantee C2/C5 actually enforce is that the
> argv is fixed and reviewable; the interpreter deny-list is defense in
> depth over that, not the barrier itself. Adding a form is an amendment to
> this clause, not a lint-rule tweak.

`0025` `artifacts/verification.md` §Phase 3a (C5 row):

> the `env`-chain interpreter form (`env FOO=1 /bin/bash -c`) is caught.

`0025` `evidence/critique/Charted.md` (stdin appetite):

> The spike-proven `tee {artifact}` hazard destroys the user's artifact and
> is not statically detectable from argv, so no C5 arm can catch it.

⇒ C5 names `env chains` in scope, so option-flag blindness is a
conformance gap, not an OPEN-set case; the verification PASS text asserts
only the assignment form. Charted.md already concedes argv is not the
authority for what a command CONSUMES (stdin) — which is exactly the `sh
-s` axis — so that axis has a named successor and this RDR routes it there
rather than duplicating it.

### Source under decision

`internal/table/load.go::interpreterForm` (read in full): walks `env` +
`NAME=VALUE` words, stops at the first other word, requires THAT word's
basename to be listed, then scans ALL later words for a listed flag.
⇒ The "flag anywhere later" laxity already exists (`ruby tool.rb -e prod`
is refused today); only the interpreter POSITION is pinned. Freeing the
position closes every wrapper at once with no wrapper table, at no new
false-positive class beyond the one already accepted.

`internal/table/load.go::carrierDefect` clause 3 comment: "That exemption
belongs to the INTERPRETER FORM, not to whitespace as such: clause 4 fires
only when argv0 (after the `env` walk) names a listed interpreter".
⇒ Today `["nice","sh","-c","cat {artifact}"]` reports
`command_unknown_placeholder` — the wrong defect, masking the interpreter
form (0025's D6 failure re-opened under a wrapper). A position-free
predicate reaches clause 3 through the same `isInterp` read, so the fix is
one predicate, not two.

`interpreterForm` consumers: `grep -rn 'interpreterForm' --include='*.go'
internal cmd` → only `internal/table/load.go` (definition + one call in
`carrierDefect`). Tests: `internal/table/command_carrier_0025_test.go::
TestReq74_KnownInterpreterWithInlineCodeIsALoadTimeDefect` (4 probes:
`sh -c`, `bash -c`, `python -c`, `env sh -c`), plus the negative fixture
`neg/neg-command-shell-interpreter.toml` (`["sh", "-c", "rdr-gate"]`).

Sibling-path check (step 5): `internal/cli/cmdbind/cmdbind.go::resolveArgv0`
reasons about argv0's LOCATION (abs / bare / relative), not its identity;
`internal/cli/flowbind/registry.go` selects by CARRIER (`command` vs
`path`). Searched for any wrapper / command-position classifier in the
repo: none exists.

## Queries (corpora + checkouts)

- DevRef semantic ×2: "allowlist of commands bypassed by wrapper programs
  such as env nice timeout xargs that re-exec a shell" → Unix Power Tools
  `nice` chapters (noise); "default deny: whitelist known-good commands
  rather than blacklist known-bad …" → TCP-wrappers / PAM "Implicit Deny"
  (class principle only, no wrapper-bypass content). REJECTED as
  load-bearing.
- PapersFast semantic ×1: "command allowlist bypass via wrapper programs
  env nice timeout xargs that exec a shell; argv-based policy is
  insufficient" → FreeBSD firewall rc.conf, file permissions (noise).
  REJECTED.
- ⚠ no prior-art coverage in the project's corpora for the
  wrapper-re-exec bypass class (the sudoers/GTFOBins-style problem); the
  class frame below rests on the peer-checkout instance reads and 0004/0025
  above, not on the corpora.
- Peer checkouts semantic ×1: "run a user-configured command string through
  the shell (sh -c) when a shell flag or prefix marks it, otherwise split
  into argv and exec directly" → roborev `hooks.go::runHook` (accepted),
  consul `agent/exec` (Windows/Unix split; not opened), semgrep Dockerfile
  AST `Argv | Sh_command` (exec-form vs shell-form is a DECLARED
  distinction in Dockerfile too; not opened).
- Peer checkouts literal sweep `"sh", "-c"` (non-test Go) → beads
  `creds/command.go` (accepted), beads `init_git_hooks.go` (a printed
  example), kubebuilder e2e helpers (test scaffolding; rejected).
- gh-cli: `rg -n shell pkg/cmd/alias/set/set.go` → `--shell` flag +
  `!` prefix (accepted).

## Rejected branches

- Sudoers-style per-wrapper option grammars (Alt 2): every wrapper is a
  mini-parser (`timeout -k 5 10 sh -c`, `xargs -I{} sh -c`, `env -S`), the
  wrapper set is as open as the interpreter set, and a mis-walk is a
  SILENT admission. Rejected once the position-free scan was seen to close
  the class without a table.
- Declared `shell` opt-in (gh's model): the peer pattern, but 0025:C5 says
  "no opt-in in v1" and 0004:ALT2 rejected inline script as the contract;
  visibility is served equally by the wrapper-file remediation C5 already
  names. Briefly rejected in the record.
- Exec-time binary identification (`#!` / ELF read of the resolved argv0):
  wrappers exec whatever follows, so the resolved binary is the wrapper;
  and it moves a LINT promise to run time, where the reviewer is absent.
