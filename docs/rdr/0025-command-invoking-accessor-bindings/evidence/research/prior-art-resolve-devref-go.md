Model: claude-fable-5

# RDR 0025 prior-art resolve: DevRef corpus + Go toolchain

Scope: `command = ["argv0", ..., "{artifact}"]` executed without a shell; closed whole-element
placeholder; stdin JSON string map; stdout JSON map or raw value; gate verdict as JSON or
declared exit-code map; child env policy (inherit vs allowlist); load-time lint (ban `sh -c`).

Corpus call shape used: `arc search semantic --corpus <C> --limit 8 --json "<query>"`.
Go source: Go 1.26.6 at `$(go env GOROOT)` = `/opt/local/lib/go-1.26`. All `file:line` below
are under `/opt/local/lib/go-1.26/src/`.

---

## T1 shell vs argv execution and injection

Queries (DevRef): (1) "command injection shell metacharacters execve argument vector";
(2) "system() versus execve safety untrusted input"; (3) "shell command injection untrusted
input passed to system() popen"; (4) "Set-user-ID programs should never use system() while
operating under the privileged identifier ..." (targeted re-open).

Accepted:
- DevRef: The Linux Programming Interface (Kerrisk), Ch. 27 "Program Execution",
  §27.6 "Avoid using system() in set-user-ID and set-group-ID programs", page 625:
  "Set-user-ID and set-group-ID programs should never use system() while operating under
  the program's privileged identifier. Even when such programs don't allow the user to
  specify the text of [the command...]" -- the hazard is the shell layer itself, not just
  user-controlled text; it settles "no shell" as the default for a privileged/trusted invoker.
- DevRef: TLPI Ch. 38 "Writing Secure Privileged Programs", page 835: "IFS specifies the
  delimiting characters that the shell interprets as separating the words of a command line.
  ... (Section 27.6 describes one vulnerability relating to IFS that appeared in older Bourne
  shells.)" -- the word-splitting rules of a shell are themselves an attack surface; an argv
  array has no IFS.
- DevRef: TLPI page 628 (system() caveat): "While the command is being executed, typing
  Control-C or Control-\ will kill only the child of system(), while the application
  (unexpectedly, to the user) continues to run." -- signal/exit-status semantics are muddied
  by the interposed shell.

Negatives: query (1) landed on interpreter-script (`#!`) mechanics and psql meta-commands,
not injection; no corpus passage discusses `execve` argv arrays as an injection *defense* in
those words. The argument is made only via the system()/setuid warning above.

## T2 environment as ambient authority / least privilege / hermetic

Queries (DevRef): (1) "environment variable inheritance child process security least
privilege"; (2) "sudo env_reset environment sanitized before running command";
(3) "hermetic build reproducible environment scrubbing"; (4) "erase the entire environment
and then rebuild it with selected values clearenv set-user-ID secure" (targeted re-open).

Accepted:
- DevRef: TLPI Ch. 6 §6.7 "Environment List", page 78: "[a child created] via fork()
  ... inherits a copy of its parent's environment. Thus, the environment provides a
  mechanism for a parent process to communicate information to a child process." --
  inheritance is the default and is a *communication channel*, i.e. ambient authority.
- DevRef: TLPI §6.7 (clearenv), page 173: "On occasion, it is useful to erase the entire
  environment, and then rebuild it with selected values. For example, we might do this in
  order to execute set-user-ID programs in a secure manner (Section 38.8)." -- this is the
  allowlist policy stated as a library idiom (`clearenv()` then `putenv()` selected keys).
- DevRef: Software Engineering at Google (O'Reilly 2020), Ch. 14 "Larger Testing",
  "Hermetic servers", page 519: "what goes into the test doesn't change based on outside
  dependencies, so when you run a test twice with the same application and test code, you
  should get the same results." -- names the determinism/isolation payoff of a scrubbed env.

Negatives: "sudo env_reset" -- no corpus evidence (hits were `su` subshell env leakage in
Unix Power Tools p602, which corroborates the hazard but not the sudo mechanism).

## T3 exit-code conventions and output protocols

Queries (DevRef): (1) "exit status conventions silence is golden Unix philosophy";
(2) "Rule of Silence Rule of Repair Unix"; (3) "structured output JSON versus text streams
pipes composability"; (4) "text streams are a universal interface textuality protocols".

Accepted:
- DevRef: The Art of Unix Programming (Raymond), §1.6.11 "Rule of Silence", page 53:
  "When a program has nothing surprising to say, it should say nothing. ... Well-behaved
  Unix programs do their jobs unobtrusively, with a minimum of fuss and bother." --
  supports: a read result on stdout is the *only* thing on stdout; diagnostics go to stderr.
- DevRef: Unix Power Tools (O'Reilly), §44.7/"Exit Status", page 1199: "Many (but not all)
  UNIX commands return a status of zero if everything was okay or non-zero (1, 2, etc.) if
  something went wrong. A few commands, like grep and diff, return a different non-zero
  status for different kinds of problems; see your online manual pages to find out." --
  exit codes are per-program conventions, not a universal vocabulary; this is the argument
  for a *declared* exit-code map rather than an implicit one.
- DevRef: TAOUP Ch. 5 "Textuality", §5.3 "Application Protocol Design", page 156: "When your
  application protocol is textual and easily parsed by eyeball, many good things become
  easier. Transaction dumps become much easier to interpret. Test loads become easier to
  write." -- supports JSON-on-stdio over a binary/opaque protocol.
- DevRef: TAOUP index page 554 lists "Rule of Repair, 13, 21-22, 147, 154, 252"; page 187
  case study: "SMTP failures are noisy, usefully so." -- fail loudly on protocol violation.

Negatives: "structured output JSON vs text pipes" landed on database/JSON-type chapters
(Silberschatz p397, PostgreSQL docs) and jq usage (advanced-microservices p171-172) -- no
corpus passage argues JSON-vs-plain-text for CLI composability specifically.

## T4 configuration-as-data vs configuration-as-code

Queries (DevRef): (1) "configuration complexity clock hardcoded values config file rules
engine DSL"; (2) "declarative configuration templating pitfalls config file becomes
programming language"; (3) "separate policy from mechanism interfaces from engines";
(4) "minilanguage accretion Turing-complete temptation add features control structures
configuration language".

Accepted:
- DevRef: TAOUP Ch. 8 "Minilanguages", page 217: "Sadly, most people do their first
  minilanguage the wrong way, and only realize later what a mess it is. ... Another
  notorious example of language-by-feature creep was the editor TECO, which grew first
  macros and then loops ..." (sendmail.cf is the example "every Unix guru will ... shudder
  over"). -- the "config grows into a language" hazard, by name.
- DevRef: TAOUP Ch. 8, page 240: "a lot of bad designs have been botched by designers who
  failed to face up to the fact that they really needed a minilanguage rather than a
  data-file format. Too often, language-like features get pasted on as an afterthought."
  -- the converse: decide up front that `command` is *data* (an argv array with one closed
  placeholder), not a template language.
- DevRef: TAOUP §1.6.8/Ch.1, page 49: "hardwiring policy and mechanism together has two
  bad effects: It makes policy rigid and harder to change in response to user requirements,
  and it means that trying to change policy has a strong tendency to destabilize the
  mechanisms." -- the exec mechanism (argv, stdio protocol) stays fixed; what to run is policy.
- DevRef: Continuous Delivery (Humble & Farley), Ch. 2 "Managing Software Configuration",
  page 74: "it is an enduring myth that configuration information is somehow less risky to
  change than source code. Our bet is that, given access to both, we can stop your system at
  least as easily by changing the configuration as by changing the source code." --
  justifies load-time lint of `command` entries as seriously as code review.
- DevRef: Domain-Driven Design (Evans), Ch. 10 "Declarative Design", page 166: "Generating
  a running program from a declaration of model properties is a kind of Holy Grail ... but it
  does have its pitfalls in practice. ... A declaration language not expressive enough to do
  everything needed, but a framework that makes it very difficult to [extend]" -- the
  escape-hatch requirement: a closed vocabulary must pair with an obvious exit (write a
  wrapper script and point argv0 at it).

Negatives: "configuration complexity clock" (Mike Hadlow's essay) -- no corpus evidence
under that name.

## T5 DX: defaults, progressive disclosure, escape hatches, least astonishment

Queries (DevRef): (1) "progressive disclosure escape hatch tool design defaults";
(2) "principle of least astonishment command line interface design"; (3) "closed
vocabulary placeholder substitution versus general template language"; (4) "sensible
defaults most users should never need to touch configuration".

Accepted:
- DevRef: TAOUP §11.1 "Applying the Rule of Least Surprise", page 287 (printed 255):
  "Try to find functional similarities between your program and programs they are likely to
  already know about. Then mimic the relevant parts of the existing interfaces." -- argues
  for argv-array + `{placeholder}` shapes users already know (systemd, pre-commit, k8s) and
  for `0 = pass / non-zero = fail` as the default gate map.
- DevRef: TAOUP Ch. 11, page 297: "The flip side of configurability is an urgent need for
  good defaults and an easy way to set everything to the default. The flip side of
  expressivity is a need for guidance ... on where to get started." -- progressive
  disclosure: `command` alone should work; env policy / exit maps are opt-in layers.
- DevRef: TAOUP §10.1 "What Should Be Configurable?", page 266: "Proliferating unnecessary
  options has many bad effects. One of the subtlest but most serious is what it will do to
  your test coverage." -- caution against adding a general template engine "just in case".
- DevRef: TAOUP Ch. 6 "Transparency", page 185: "programs that are opaque about what they
  are doing tend to have a lot of assumptions baked into them, and to be frustrating or
  brittle or both in any use case not anticipated by the designer." -- supports an `-x`-style
  echo of the resolved argv (cf. `go generate -x`, G1 below).

Negatives: "progressive disclosure" and "escape hatch" as terms -- no corpus evidence
(nearest is TAOUP discoverability vs transparency p167). "closed vocabulary placeholders vs
general templates" -- no corpus evidence (hits were PostgreSQL text-search templates and
a type-theory evaluation-context chapter). StateMachineRes returned only headings
(yaml-schema.md "Template Syntax", authoring-commands.md "Command Structure"); not opened.

---

## G1 `go generate` directive execution

File: `cmd/go/internal/generate/generate.go` (doc mirrored in `cmd/go/alldocs.go`).

No shell / argument splitting:
- generate.go:60-66 (alldocs.go ~577-583): "The arguments to the directive are
  space-separated tokens or double-quoted strings passed to the generator as individual
  arguments when it is run. Quoted strings use Go syntax and are evaluated before
  execution; a quoted string appears as a single argument to the generator."
- generate.go:100-102 (alldocs.go:618-620): "Other than variable substitution and
  quoted-string evaluation, no special processing such as "globbing" is performed on the
  command line."
- generate.go:382-385: `// split breaks the line into words, evaluating quoted strings and
  evaluating environment variables.` ... `func (g *Generator) split(line string) []string`
- generate.go:506: `cmd := exec.Command(path, words[1:]...)` -- argv array, no `sh -c`.

Variable vocabulary (declared, but NOT closed -- see caveat):
- generate.go:77-98: "Go generate sets several variables when it runs the generator:
  $GOARCH $GOOS $GOFILE $GOLINE $GOPACKAGE $GOROOT $DOLLAR $PATH". `$DOLLAR` -- "A dollar
  sign." `$PATH` -- "The $PATH of the parent process, with $GOROOT/bin placed at the
  beginning."
- generate.go:104-110 (alldocs.go:622-628): "As a last step before running the command,
  any invocations of any environment variables with alphanumeric names, such as $GOFILE or
  $HOME, are expanded throughout the command line. ... Due to the order of evaluation,
  variables are expanded even inside quoted strings. If the variable NAME is not set,
  $NAME expands to the empty string."
- generate.go:444: `words[i] = os.Expand(word, g.expandVar)`; generate.go:462-470
  `expandVar` checks `g.env` first, then falls back to `return os.Getenv(word)`.
  CAVEAT for RDR 0025: the vocabulary is *open* (any parent env var expands, unset -> ""),
  expansion is in-string, and `$DOLLAR` exists precisely because there is no other escape.
  This is the design RDR 0025's closed whole-element `{artifact}` deliberately avoids.

Child environment (sets specific vars AND inherits the rest):
- generate.go:365-380 `setEnv`: builds `GOROOT=`, `GOARCH=`, `GOOS=`, `GOFILE=`, `GOLINE=`,
  `GOPACKAGE=`, `DOLLAR=$`, then `base.AppendPATH(env)` and `base.AppendPWD(env, g.dir)`.
- generate.go:513: `cmd.Env = str.StringList(cfg.OrigEnv, g.env)` -- parent env
  (`cfg.OrigEnv`) first, overrides appended; i.e. inherit + overlay, not allowlist.
- generate.go:498-505: resolves bare command names against `$GOROOT/bin` first
  ("If a generator says '//go:generate go run <blah>' it almost certainly intends to use
  the same 'go' as 'go generate' itself"), then `cmd.Args[0] = words[0]`.
- generate.go:509-511: `cmd.Stdout = os.Stdout; cmd.Stderr = os.Stderr` and
  `cmd.Dir = g.dir` ("Run the command in the package directory").

`-command` alias:
- generate.go:114-123: "//go:generate -command xxx args... specifies, for the remainder of
  this source file only, that the string xxx represents the command identified by the
  arguments. This can be used to create aliases or to handle multiword generators."
- generate.go:342-345 dispatch; generate.go:472-482 `setShorthand`: errors on
  "no command specified for -command" and "command %q multiply defined"; stores
  `g.commands[command] = slices.Clip(words[2:])`.
- generate.go:347-352: with `-n`/`-x` the resolved words are echoed to stderr
  (`strings.Join(words, " ")`) before/instead of exec -- the transparency affordance.

## G2 `os/exec` hardening history

File: `os/exec/exec.go`, `os/exec/lp_unix.go`.

ErrDot (Go 1.19 change, with a GODEBUG reversal knob):
- exec.go:30-38: "Modern practice is that including the current directory is usually
  unexpected and often leads to security problems. To avoid those security problems, as of
  Go 1.19, this package will not resolve a program using an implicit or explicit path entry
  relative to the current directory. ... these functions return an error err satisfying
  errors.Is(err, ErrDot)."
- exec.go:83-86: "Setting the environment variable GODEBUG=execerrdot=0 disables generation
  of ErrDot entirely, temporarily restoring the pre-Go 1.19 behavior for programs that are
  unable to apply more targeted fixes. A future version of Go may remove support for this
  variable."
- exec.go:1331-1338: `// ErrDot indicates that a path lookup resolved to an executable in
  the current directory ...` `var ErrDot = errors.New("cannot run executable found relative
  to current directory")`.
- exec.go:1340-1348: `validateLookPath` -- "excludes paths that can't be valid executable
  names. See issue #74466 and CVE-2025-47906." rejects `""`, `"."`, `".."` (a second,
  later hardening in the same spot).
- lp_unix.go:44-47: `// NOTE(rsc): I wish we could use the Plan 9 behavior here (only bypass
  the path if file begins with / or ./ or ../) but that would not match all the Unix
  shells.`; lp_unix.go:60-63: `if dir == "" { // Unix shell semantics: path element "" means
  "." dir = "." }`; lp_unix.go:67-72: `if !filepath.IsAbs(path) { if execerrdot.Value() !=
  "0" { return path, &Error{file, ErrDot} } execerrdot.IncNonDefault() }`.
  Plain statement: Go shipped PATH-relative-to-cwd resolution for ~10 years (Go 1.0-1.18),
  then in Go 1.19 decided it was a security defect and reversed it, keeping an explicit,
  deprecated opt-out. This is a "decided X, reverted after use" signal in favour of
  resolving argv0 strictly (absolute path, or PATH without `.`), and of refusing `""`.

`Cmd.Env`, `Environ()`, `Dir`:
- exec.go:162-172: "Env specifies the environment of the process. Each entry is of the form
  "key=value". If Env is nil, the new process uses the current process's environment. If Env
  contains duplicate environment keys, only the last value in the slice for each duplicate
  key is used. ... See also the Dir field, which may set PWD in the environment."
  -> inherit is the zero-value default; allowlist must be an explicit non-nil slice
  (and `[]string{}` is *empty*, not inherit -- a footgun for RDR 0025's policy encoding).
- exec.go:1238-1244: "Environ returns a copy of the environment in which the command would
  be run as it is currently configured." (added Go 1.19 so callers can do
  `append(cmd.Environ(), "K=V")` and keep the PWD/Dir coupling).
- exec.go:174-191: "Dir specifies the working directory of the command. If Dir is the empty
  string, Run runs the command in the calling process's current directory. On Unix systems,
  the value of Dir also determines the child process's PWD environment variable if not
  otherwise specified."
- exec.go:385-394 (`Command` doc): "If name contains no path separators, Command uses
  LookPath to resolve name to a complete path if possible. Otherwise it uses name directly
  as Path. ... Args[0] is always name, not the possibly resolved Path." Note: a *relative*
  `Path` containing a separator (e.g. `./tool`) is resolved by the OS against the parent's
  cwd at exec time, not against `Cmd.Dir` -- a long-standing documented gotcha; RDR 0025
  should resolve relative argv0 explicitly against a declared base before exec.

## G3 cgo flag allowlist (`cmd/go/internal/work/security.go`)

- security.go:5-10: `// Checking of compiler and linker flags. // We must avoid flags like
  -fplugin=, which can allow // arbitrary code execution during the build. // Do not make
  changes here without carefully // considering the implications. // (That's why the code is
  isolated in a file named security.go.)`
- security.go:18-23: `// Note also that GNU binutils accept any argument @foo // as meaning
  "read more flags from the file foo", so we must // guard against any command-line argument
  beginning with @, // even things like "-I @foo". // We use load.SafeArg (which is even more
  conservative) // to reject these.`
- security.go:44-59: `var validCompilerFlags = []*lazyregexp.Regexp{ re(`-D([A-Za-z_]...`),
  ... re(`-I([^@\-].*)`), ...`; security.go:92: `re(`-fsanitize=(.+)`)`; security.go:170
  (linker): `re(`-fsanitize=([^@\-].*)`)`.
- security.go:342-362 `checkFlags`: `// Let users override rules with $CGO_CFLAGS_ALLOW,
  $CGO_CFLAGS_DISALLOW, etc.` -- `cfg.Getenv("CGO_" + name + "_ALLOW")` compiled as a
  regexp; parse failure is a hard error `"parsing $CGO_%s_ALLOW: %v"`.
  Precedent: a positive regexp allowlist over *individual argv elements*, a documented
  escape hatch (env-provided regexp), and a note that argument-level rejection must consider
  how the callee re-splits (`-Wl,a,b`, `@file`). Directly analogous to RDR 0025's lint that
  inspects argv0 and each element, and to the "closed placeholder is a whole element" rule.

## G4 how `go get`/modfetch run `git`/`hg`

- vcs/vcs.go:521-525: `cmd := exec.Command(v.Cmd, args...)`; `cmd.Dir = dir`;
  `if v.Env != nil { cmd.Env = append(cmd.Environ(), v.Env...) }` -- argv array; inherit +
  overlay.
- vcs/vcs.go:145-151 (hg): `// HGPLAIN=+strictflags turns off additional output that a user
  may have enabled via config options or certain extensions.` `Env:
  []string{"HGPLAIN=+strictflags"}` -- an env var set *to neutralise the user's own
  config* so the child's output protocol is parseable.
- vcs/vcs.go:98-104: `// GIT_ALLOW_PROTOCOL is an environment variable defined by Git. It is
  a colon-separated list of schemes that are allowed to be used with git fetch/clone. Any
  scheme not mentioned will be considered insecure.` -- env read by the *parent* as policy.
- vcs/vcs.go:752-765: `// ... Bazaar, Fossil, and Subversion have primarily been used in
  trusted, authenticated environments and are not as well scrutinized as attack surfaces.
  // See golang.org/issue/41730 for details.` `var defaultGOVCS = govcsConfig{{"private",
  []string{"all"}}, {"public", []string{"git", "hg"}}}` -- an allowlist of *which
  executables may be invoked*, added after the fact (Go 1.16, issue 41730) once the
  attack surface was recognised: a "shipped open, then narrowed" signal.
- modload/init.go:505-510: `// If user has explicitly set GIT_TERMINAL_PROMPT=1, keep` ...
  `if os.Getenv("GIT_TERMINAL_PROMPT") == "" { os.Setenv("GIT_TERMINAL_PROMPT", "0") }` --
  sets a child-affecting env var in the parent, only if unset (defaults, but user wins).
- work/buildid.go:224: `cmd.Env = append(os.Environ(), "LC_ALL=C")` -- pin locale so
  compiler output is parseable (output-protocol stabilisation via env).
- modfetch/codehost/codehost.go:375-385: `// TODO: Impose limits on command output size.
  // TODO: Set environment to get English error messages.` ... `c :=
  exec.CommandContext(ctx, cmd[0], cmd[1:]...)`; `c.Cancel = func() error { return
  c.Process.Signal(os.Interrupt) }`; `c.Env = append(c.Environ(), args.env...)`.
- modfetch/codehost/git.go:981-984: `// Manually supply GIT_DIR so Git works with
  safe.bareRepository=explicit set.` `args.env = []string{"GIT_DIR=" + r.dir}`.
- vcweb/script.go:136 (test harness only): `"GIT_CONFIG_NOSYSTEM=1"`.
  Negative: no `GIT_SSH_COMMAND`, `GIT_CONFIG_*` unsetting, or general env *scrubbing* in
  the production `go` command; the toolchain's consistent posture is inherit + overlay with
  a handful of protocol-stabilising keys (HGPLAIN, LC_ALL, GIT_TERMINAL_PROMPT, GIT_DIR)
  and allowlisting applied to *what runs* (GOVCS) and *which flags pass* (security.go),
  not to the env.

---

## Precedent classes (unverified)

All items below are stated from memory to name the class only; none were verified against a
local checkout.

- systemd `ExecStart=` (unverified -- from memory): argv split by systemd, no shell unless
  you spell `/bin/sh -c`; `%i`/`%n` specifiers are a closed vocabulary with `%%` escape;
  `Environment=`/`PassEnvironment=` opt-in over a near-empty default env; `SuccessExitStatus=`
  is a declared exit-code map; `ExecStart=-cmd` prefix ignores failure.
- sudo `env_reset` + `env_keep` (unverified -- from memory): env is reset to a minimal set by
  default with an explicit keep-list, made the default after a series of env-driven CVEs
  (LD_PRELOAD, PS4/bash, PERLLIB-class issues).
- Bazel `--incompatible_strict_action_env` (unverified -- from memory): actions once
  inherited the client env; the flag flipped the default to a fixed PATH and empty env, with
  `--action_env=KEY` as the allowlist escape hatch.
- git `difftool.trustExitCode` (unverified -- from memory): defaults to false; git ignores
  the tool's exit code unless the user opts in, because tools' exit conventions are unknown.
- Nagios plugin exit codes (unverified -- from memory): 0 OK / 1 WARNING / 2 CRITICAL /
  3 UNKNOWN, a fixed 4-value verdict map plus first-line stdout as the human reason.
- Kubernetes `$(VAR)` in `command`/`args` (unverified -- from memory): in-string expansion
  limited to the container's declared `env` names, `$$` escape, unresolvable references left
  literal; a declared-vocabulary in-string template.
- pre-commit `entry` (unverified -- from memory): shlex-split, executed without a shell,
  filenames appended as trailing argv elements; `language: script/system`.
- Ansible `command` vs `shell` (unverified -- from memory): `command` is argv/no-shell and
  the documented default; `shell` is the explicit escape hatch; `failed_when`/`changed_when`
  let the play declare its own rc map instead of trusting rc != 0.
- Docker credential helpers / client-go exec credential plugins (unverified -- from memory):
  action on argv, payload on stdin, JSON on stdout; client-go adds `apiVersion` in the
  `ExecCredential` envelope and passes `KUBERNETES_EXEC_INFO` via env; the `v1alpha1`
  version was later removed in favour of `v1beta1`/`v1`.

## Reversal signals

Ranked strongest to weakest as "decided X, then reverted/narrowed after use":

1. Go `os/exec` ErrDot (VERIFIED, exec.go:30-38, 83-86; lp_unix.go:67-72): cwd-relative
   PATH resolution shipped for ~10 years, reversed in Go 1.19 with a deprecated
   `GODEBUG=execerrdot=0` opt-out; then further narrowed (`validateLookPath`,
   CVE-2025-47906). Lesson: resolve argv0 strictly at load time; never let `""`/`.` through.
2. sudo `env_reset` default (unverified -- from memory): inherit-everything was the shipped
   default; scrubbed-with-keep-list became the default after CVEs. Lesson: allowlist is the
   safe default; inherit is the opt-in.
3. Bazel `--incompatible_strict_action_env` (unverified -- from memory): same shape as #2
   in a build tool; inherit-client-env default flipped to scrubbed.
4. Go `GOVCS` allowlist (VERIFIED, vcs.go:752-765, golang.org/issue/41730): `go get` ran any
   supported VCS binary for years; narrowed to `git|hg` for public modules in Go 1.16 once the
   "client of an untrusted server" surface was recognised. Lesson: allowlist *which
   executables* may be named, or at least lint argv0.
5. git `difftool.trustExitCode` (unverified -- from memory): exit code explicitly *not*
   trusted by default because tool conventions vary. Lesson: gate exit-code maps must be
   declared, not inferred (corroborated by Unix Power Tools p1199).
6. client-go exec credential `v1alpha1` -> `v1beta1`/`v1` (unverified -- from memory):
   protocol envelope needed a version field and a deprecation path. Lesson: put a version
   or kind discriminator in the stdin/stdout JSON from day one.
7. `go generate` `$DOLLAR` (VERIFIED, generate.go:92-93, 104-110): an escape variable had
   to be *added* because in-string `$NAME` expansion over an open vocabulary made a literal
   `$` unwritable. Not a reversal, but a scar showing why RDR 0025's whole-element closed
   placeholder needs no escape mechanism.
