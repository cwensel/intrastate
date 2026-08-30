Model: claude-fable-5

# Prior art (batch B): subprocess command carriers in Go/Python CLIs

Scope: how seven projects declare, substitute, spawn, feed, read, and time out user-declared subprocesses, and which of those choices were later reversed. Paths are relative to each repo's checkout root.

---

## consul

1. Carrier shape: BOTH. `script`/`handler` = shell string; `args` = argv list, run directly. `agent/exec/exec_unix.go:15` `Script()` -> `exec.Command(shell, "-c", script)` where `shell` is `/bin/sh` unless `$SHELL` is set. `agent/exec/exec.go:12` `Subprocess(args)` -> `exec.Command(args[0], args[1:]...)`; empty argv is a load error ("need an executable to run"). Dispatch: `agent/checks/check.go::CheckMonitor.check` — `if len(c.ScriptArgs) > 0 { exec.Subprocess } else { exec.Script }`. Watches: `agent/watch_handler.go:33` `makeWatchHandler` switches on `string` (shell) vs `[]string` (argv). CLI `consul watch` has `-shell` flag, default **true** (`command/watch/watch.go:68`), so the CLI still defaults to a shell; the agent config path is what was deprecated.
2. Substitution: none. No placeholders in `args`; all variable data travels via env + stdin. `agent/watch_handler.go:216` rejects `handler`+`args` both set (mutually exclusive), and either with `type=http`.
3. Env policy: INHERIT + one var. `agent/watch_handler.go:62` `cmd.Env = append(os.Environ(), "CONSUL_INDEX="+idx)`. Script checks set no extra vars. Process attrs: `agent/exec/exec_unix.go:23` `SetSysProcAttr` -> `Setpgid: true`; `KillCommandSubtree` -> `syscall.Kill(-pid, SIGKILL)` (whole process group).
4. Output protocol: RAW combined stdout+stderr into a ring buffer. Watch handler stdin = JSON encoding of the watch data (`cmd.Stdin = &inp` after `json.NewEncoder(&inp).Encode(data)`, `agent/watch_handler.go:74-83`). Caps: `WatchBufSize = 4 * 1024` (`agent/watch_handler.go:29`); checks `check_output_max_size` default 4096 (`agent/config/builder.go:1055`), truncated with a prefix `"Captured %d of %d bytes\n...\n"`. Timeout: default 30s, then kill subtree, status Critical with output appended; after the kill it still `<-waitCh` "so we never start another instance concurrently".
5. Exit-code semantics (script checks, `agent/checks/check.go::CheckMonitor.check` tail): `err == nil` -> passing; `*exec.ExitError` with `WaitStatus.ExitStatus() == 1` -> warning; anything else -> critical. Watches: non-zero is only logged ("Failed to run watch handler").
6. REVERSALS (two, both strong):
   - Shell default -> argv. CHANGELOG (1.0.0 section, line ~6522): "**Support for Running Subproccesses Directly Without a Shell:** Consul agent checks and watches now support an `args` configuration which is a list of arguments to run for the subprocess, which runs the subprocess directly without a shell. The old `script` and `handler` configurations are now deprecated (specify a shell explicitly if you require one). A `-shell=false` option is also available on `consul lock`, `consul watch`, and `consul exec`" [GH-3509]. Runtime warning: `agent/watch_handler.go:201` "The 'handler' field in watches has been deprecated and replaced with the 'args' field."
   - Script checks off by default. CHANGELOG (0.9.0, line ~6619): "Added a new `enable_script_checks` configuration option that defaults to `false`, meaning that in order to allow an agent to run health checks that execute scripts, this will need to be configured and set to `true`. This provides a safer out-of-the-box configuration for Consul where operators must opt-in to allow script-based health checks." [GH-3087]. Later split (1.3.0, line ~6036): "Operators can now enable script checks from local config files only." [GH-4711] -> `enable_local_script_checks` (`agent/config/builder.go:826-827`: local defaults to the remote value). Enforcement `agent/agent.go:2947-2953`: "Scripts are disabled on this agent; to enable, configure 'enable_script_checks' or 'enable_local_script_checks' to true" / "...from remote calls; to enable, configure 'enable_script_checks' to true". Further tightening: warning when script checks enabled without ACLs ("this is not recommended for security reasons", `agent/agent.go:2956`; CHANGELOG GH-7437, GH-22877).

---

## opentofu

1. Carrier shape: SHELL STRING with an argv-list interpreter prefix. `local-exec` `command` is a single string; `interpreter` is a list. `internal/builtin/provisioners/local-exec/resource_provisioner.go:133-150`: if `interpreter` unset, `cmdargs = []string{"cmd", "/C"}` on windows else `{"/bin/sh", "-c"}`; then `cmdargs = append(cmdargs, command)`; `exec.CommandContext(p.ctx, cmdargs[0], cmdargs[1:]...)`. Docs (`website/docs/language/resources/provisioners/local-exec.mdx:41-44`): "It is evaluated in a shell, and can use environment variables or OpenTofu variables." Interpreter examples: `["perl","-e"]`, `["PowerShell","-Command"]`.
2. Substitution: HCL interpolation at plan time (`${self.private_ip}` style), i.e. the whole string is user-templated before the shell sees it — the classic injection surface; the docs offer `interpreter` but no quoting helper. No placeholder scheme inside the provisioner itself.
3. Env policy: INHERIT + `environment` map + auto `TRACEPARENT`. `resource_provisioner.go:172-186` `cmdEnv = os.Environ()`, append user map entries `k=v`, then inject `TRACEPARENT` when tracing active unless user set it (CHANGELOG line 22, #4014). Docs: "inherits the current process environment." `cmd.Dir = working_dir`.
4. Output protocol: RAW; stdout and stderr both go to one pipe streamed to UI (`cmd.Stderr = pw; cmd.Stdout = pw`); `quiet = true` suppresses only the "Executing:" echo of the argv. No cap; no stdin.
5. Exit-code semantics: any error from `cmd.Wait()` fails the provisioner (`on_failure = continue` is the escape hatch at the HCL layer). No exit-code map.
   Plugin protocol (go-plugin v1.7.0, `go.mod:58`): `internal/plugin/serve.go:28-38` `HandshakeConfig{ MagicCookieKey: "TF_PLUGIN_MAGIC_COOKIE", MagicCookieValue: "d602bf8f..." }` with the comment "The magic cookie values should NEVER be changed." The cookie is an env var set by the host so a provider binary refuses to run when executed by hand; the child's stdout carries the one-line handshake, and stderr is captured via `SyncStderr: logging.PluginOutputMonitor(...)` (`internal/command/meta_providers.go:373,444`). The go-plugin client source is not in this checkout, so the exact env-var list beyond the cookie is not quoted here.
6. REVERSALS: none found for `local-exec` shell semantics (CHANGELOG mentions only the TRACEPARENT addition). The design instead layers a second, explicit mechanism (the go-plugin provider protocol with magic cookie, argv exec, no shell) for anything trusted; `local-exec` docs steer users toward `interpreter` overrides rather than fixing the shell default.

---

## prometheus

No non-test `exec.Command` use in the repo and no textfile collector (that lives in node_exporter); promtool does not spawn subprocesses. Nothing to report for items 1-6.

---

## etcd

1. Carrier shape: ARGV LIST only. `etcdctl lock <name> [cmd args...]`: `etcdctl/ctlv3/command/lock_command.go:97` `exec.Command(cmdArgs[0], cmdArgs[1:]...)`. `etcdctl watch ... -- cmd args`: `watch_command.go:178` `exec.CommandContext(c.Ctx(), execArgs[0], execArgs[1:]...)`. No shell path exists.
2. Substitution: none in argv; event data goes to env only. `watch_command.go:255-297` errors `"bad number of arguments (found conflicting environment key)"` when interactive-mode args collide with env-supplied keys.
3. Env policy: INHERIT + explicit vars. `watch_command.go:179-183` `cmd.Env = os.Environ()` then `ETCD_WATCH_REVISION`, `ETCD_WATCH_EVENT_TYPE`, `ETCD_WATCH_KEY`, `ETCD_WATCH_VALUE` (%q-quoted). `lock_command.go:98` `cmd.Env = append(environLockResponse(m), os.Environ()...)` (lock key/header vars).
4. Output protocol: passthrough (`cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr`); no stdin feed; no cap.
5. Exit-code semantics: `lock_command.go:58-66` `getExitCodeFromError` unwraps `*exec.ExitError` -> `WaitStatus.ExitStatus()` and etcdctl exits with the child's code. `--command-timeout` (`etcdctl/ctlv3/ctl.go:60`, "timeout for short running command (excluding dial timeout)", default 5s) bounds the etcd RPC, not the child process.
6. REVERSALS: none found.

---

## roborev

1. Carrier shape: BOTH, by trust tier. AI agents: binary resolved by `PATH` lookup for a fixed agent name, args built in code (`internal/agent/cli_runner.go:36` `exec.CommandContext(ctx, spec.Command, spec.Args...)`; `claude.go:318,690`, `codex.go:233+`). User hooks: SHELL STRING — `internal/config/config.go:54` `Command string \`toml:"command"\` // shell command with {var} templates`; executed at `internal/daemon/hooks.go:330-333` as `exec.Command("powershell","-NoProfile","-Command", command)` on windows else `exec.Command("sh", "-c", command)`.
2. Substitution: hook `{repo}`, `{repo_name}`, `{sha}` etc. replaced via `strings.NewReplacer` with each value passed through `shellEscape` (`hooks.go:293-318`: single-quote wrap, `'` -> `'"'"'` on POSIX, `''` on PowerShell). The placeholder can appear inside a larger word — it is a text splice into shell source, mitigated by quoting rather than by argv boundaries. Agents: no placeholders; prompt goes on stdin (`claude.go:698` `cmd.Stdin = strings.NewReader(prompt)`).
3. Env policy: agents INHERIT + explicit. `internal/agent/process.go:35-38` `if cmd.Env == nil { cmd.Env = cmd.Environ() }; cmd.Env = append(cmd.Env, "GIT_OPTIONAL_LOCKS=0")` ("Prevent agents from taking .git/index.lock in the user's repo"); `claude.go:693` `buildClaudeEnv(cmd.Environ(), model, baseURL)`. Hooks: inherit (no `cmd.Env`), `cmd.Dir = workDir`. Kill mechanics `process.go:19-52`: `cmd.WaitDelay = 5s`, explicit `cmd.Cancel` because "Go's exec.CommandContext only provides a default Kill cancel when both Cancel==nil and WaitDelay==0"; a tracker records `canceledByContext` to distinguish timeout from real failure (`process.go:107` `errors.Is(runErr, exec.ErrWaitDelay)`). Process groups appear only in the test harness (`internal/testutil/integration/harness.go:160` `Setpgid: true`).
4. Output protocol: agents = streamed stdout, JSON stream for claude (`stream_json.go`) with trailing-text accumulation (`review_output.go`); hooks = `cmd.CombinedOutput()` logged raw, no cap.
5. Exit-code semantics: agents `*exec.ExitError` classified (`process.go:130`) plus quota/limit parsing (`limit_parse.go`); hooks: any error is logged ("Hook error (cmd=%q dir=%q)") and swallowed.
6. REVERSALS: none found in checkout; the shell-string hook remains. Note the split: the thing the daemon trusts (agent binary) is argv-only with a hardened kill path; the thing the user declares (hook) is shell with escaping.

---

## beads

1. Carrier shape: BOTH, by trust tier. Git/hooks: argv (`internal/beads/beads.go:287` `exec.CommandContext(ctx, "git", "-C", dir, args...)`; `internal/hooks/hooks_unix.go:54` `exec.CommandContext(ctx, hookPath, issue.ID, event)`). Credential helper: SHELL STRING — `internal/creds/command.go:86-90` `exec.CommandContext(ctx, "cmd.exe", "/C", command)` on windows else `exec.CommandContext(ctx, "sh", "-c", command)`, comment: "POSIX shells parse the helper command".
2. Substitution: none. Hook receives `(issueID, event)` as argv and the issue JSON on stdin (`hooks_unix.go:55` `cmd.Stdin = bytes.NewReader(issueJSON)`).
3. Env policy: INHERIT (no `cmd.Env` set on hooks or creds). Dolt credentials are set ONLY via `cmd.Env` on CLI subprocesses ("Credentials are set on the subprocess environment only via cmd.Env", `internal/storage/dolt/federation.go:549`; `credentials.go:429-477` copies `os.Environ()` filtering the credential prefix first). Kill: `hooks_unix.go:66-84` "Creating a process group (Setpgid) and sending a negative PID to syscall.Kill ensures the entire group" is killed on timeout (`syscall.Kill(-pid, SIGKILL)`, ESRCH tolerated).
4. Output protocol: creds = stdout parsed as JSON envelope if first byte is `{` ("A JSON object is read as the ExecCredential/getToken envelope; otherwise the trimmed output is taken as a bare token. A bare value containing whitespace is rejected", `creds/command.go:140-160`); empty output is an error. Hooks = stdout/stderr buffered for tracing only. Timeouts: creds 30s ("a helper that hangs must not wedge an open", `command.go:54`), git 5s, hooks configurable.
5. Exit-code semantics: creds non-zero -> error with trimmed stderr appended; hooks non-zero -> returned error; no code maps.
6. REVERSAL / trust gate: `internal/configfile/configfile.go:517-526` — the credential command "is deliberately read from the environment only, NOT metadata.json: a metadata-sourced command is arbitrary code run on open, so persisting it waits on a workspace-trust gate." I.e. a repo-persisted shell command was considered and refused pending a trust mechanism.

---

## semgrep

1-5 (semgrep's own subprocess use, `cli/src/semgrep/core_runner.py`): argv only — `asyncio.create_subprocess_exec(*self._cmd, stdout=PIPE, stderr=PIPE|None, limit=INPUT_BUFFER_LIMIT, preexec_fn=setrlimits_preexec_fn)` (`core_runner.py:519-528`); env inherited (no `env=`), telemetry injected before fork; timeouts passed to the core as flags (`-timeout`, `-timeout_threshold`, `-timeout_for_interfile_analysis`, `-secrets_timeout`, `core_runner.py:1107-1183`); on failure it waits 1s for the process to exit and otherwise logs "semgrep timed out waiting for the semgrep-core process to exit after an exception was raised" (`core_runner.py:549-555`); stdout is JSON parsed by pysemgrep; exit code surfaced as `CompletedProcess`.

Codified "don't do that" rules (registry snapshot under `perf/r2c-rules/` and `tests/precommit_dogfooding/`; the Go `exec.Command("sh","-c",…)` rule is not present in this checkout's snapshot):
- `python.lang.security.audit.subprocess-shell-true.subprocess-shell-true` (`tests/precommit_dogfooding/bandit.yml:232`): "Found 'subprocess' function '$FUNC' with 'shell=True'. This is dangerous because this call will spawn the command using a shell process. Doing so propagates current shell settings and variables, which makes it much easier for a malicious actor to execute commands. Use 'shell=False' instead." Note `pattern-not: subprocess.$FUNC("...", shell=True, ...)` — a fully literal string is exempted; only non-literal command text fires. Has an autofix `shell=True` -> `shell=False`.
- `python.lang.security.audit.dangerous-system-call.dangerous-system-call` (`perf/r2c-rules/r2c-security-audit.yml:6304`): "Found dynamic content used in a system call. This is dangerous if external data can reach this function call because it allows a malicious actor to execute commands. Use the 'subprocess' module instead, which is easier to use without accidentally exposing a command injection vulnerability." (CWE-78) — same literal-string exemption (`pattern-not: os.$W("...", ...)`).
- `python.lang.security.audit.dangerous-subprocess-use-tainted-env-args...` and `python.lang.security.dangerous-subprocess-use.dangerous-subprocess-use` (`tests/precommit_dogfooding/python.yml:2938, 6048`): taint variants that flag argv built from env/args even without a shell.
- `java.lang.security.audit.command-injection-formatted-runtime-call` (`perf/r2c-rules/java.yml:114`): formatted string into `Runtime.exec`.
6. REVERSALS: not applicable (rules, not a config surface). The shared shape of every rule: static literal command = OK; dynamic text reaching a shell = finding; dynamic argv element = taint-level warning only.

---

## Reversal signals

Ranked by strength (how directly a shipped design was later withdrawn, and how explicit the rationale is):

1. consul: shell-string `script`/`handler` deprecated in favour of argv `args` (1.0.0, GH-3509) — "runs the subprocess directly without a shell... The old `script` and `handler` configurations are now deprecated (specify a shell explicitly if you require one)". Runtime warning still emitted. STRONG: direct carrier reversal matching the C5 `command_shell_interpreter` defect.
2. consul: script checks flipped to off-by-default (0.9.0, GH-3087) — "defaults to `false`... This provides a safer out-of-the-box configuration for Consul where operators must opt-in"; later split into local vs remote opt-in (GH-4711) and a warning when enabled without ACLs (GH-7437, GH-22877). STRONG: a declared-command feature was re-gated after shipping, twice.
3. beads: repo-persisted credential command refused — "a metadata-sourced command is arbitrary code run on open, so persisting it waits on a workspace-trust gate"; only the env var is honoured. MEDIUM-STRONG: a pre-emptive refusal, not a post-hoc reversal, but the rationale is identical to intrastate's load-time lint concern (a TOML-declared command is code run on open).
4. roborev: agent subprocesses gained a hardened cancel path (`WaitDelay` + explicit `Cancel` + canceledByContext tracker) and `GIT_OPTIONAL_LOCKS=0` after observed contention — MEDIUM: behavioural fixes to timeout/kill semantics, no carrier change; user hooks remain `sh -c` with quoting.
5. opentofu/terraform: no reversal of `local-exec`'s shell default; mitigation is the `interpreter` argv-list escape hatch and the separate magic-cookie plugin protocol for trusted binaries. WEAK: absence of reversal despite known hazard.
6. etcd, prometheus, semgrep-core: argv-only from the start; no reversal needed. semgrep rules encode the community norm (shell + dynamic text = CWE-78 finding; literal-only exempt).
