Model: claude-fable-5

# Prior art: model-declared external command invocation in peer state-machine / workflow tools

Scope: the 32 peer checkouts under the state-machines corpus. Question: when a MODEL or DECLARATION (not host-language code) names an external command, how is it carried, templated, enveloped, executed, and later hardened. Triage query per repo: `rg -l -i --max-count 1 'shell=True|sh -c|subprocess|os/exec|exec\.Command|child_process|spawn\('` with tests/node_modules/dist excluded. Repos whose only hits were build scripts, notebooks, or vendored highlight.js were dropped after one query. The corpus analysis docs (contrast/, research/, study/) mention shells only as a driver surface for the resolver ("CLI/shell drivable", `research/resolver-vs-roundtrip.md:148`); none summarize command carriers, so the per-repo sweep below is original.

Legend for each block: (1) carrier, (2) templating, (3) env, (4) output protocol / re-entry, (5) exit codes and timeouts, (6) reversal or hardening signal.

## awf-cli (Go; YAML state-machine workflows)

The strongest analogue: a YAML step type `command` whose value is a shell string.

1. Carrier: a single shell string, always. `internal/domain/workflow/step.go:77` `Command string // for command type`, `:78` `ScriptFile string // F064: external script file (XOR with Command)`. Executed via `internal/infrastructure/executor/shell_executor.go::Execute` — `execCmd = exec.CommandContext(ctx, e.shellPath, "-c", cmd.Program) //nolint:gosec // G204: intentional dynamic shell`. No argv form exists.
2. Templating: Go `text/template` syntax over the whole string (`docs/reference/interpolation.md:5` "AWF uses Go template syntax (`{{.var}}`) for variable interpolation"). Substitution is textual, in-string, then the result is handed to `$SHELL -c`. Load-time `TemplateValidator.ValidateStep` (`internal/domain/workflow/template_validation.go:146`) validates references, not shell safety.
3. Env: child inherits the full parent environment; declared step env is appended on top (`shell_executor.go` `execCmd.Env = os.Environ(); for k, v := range cmd.Env { ... append }`). Step-level `Env []string` is a "required environment variables" precondition, not an allowlist (`workflow.go:27`). Secrets are masked post hoc in captured stdout/stderr (`e.masker.MaskText(...)`, ADR-011).
4. Output: raw stdout captured as a string and re-enters as `{{.states.<step>.Output}}` (`docs/user-guide/agent-steps.md:442` "`{{.states.step_name.Output}}` - Previous step raw output"); oversized output spills to a file addressable as `.OutputPath` (`docs/user-guide/configuration.md:295`). Typed JSON post-processing (`output_format: json`) exists only for agent steps (`agent_config.go:38`), not command steps.
5. Exit codes: 0 = success, any non-zero routes to `on_failure`/transitions (`internal/application/execution_service.go:760` `if result.ExitCode != 0 { return s.handleNonZeroExit(...) }`). Retry can be scoped to `retryable_exit_codes: [1, 22]` (`docs/user-guide/examples.md:287`; `step.go:40`). Context cancellation kills the whole process group and reports `ExitCode = -1` (`shell_executor.go`). The CLI's own exit taxonomy is 1..4 (ADR-002: "3 | Execution error | Command failed, timeout, agent error").
6. ★ Reversals / hardening after use:
   - `examples/plugins/awf-plugin-security-validator/main.go` — a post-hoc lint plugin. `checkForDangerousCommands` (`:116`) string-matches argv0 against `"rm", "dd", "mkfs", "format", "fdisk", "kill", "killall", "pkill"` and `checkForCommandInjection` (`:141`) warns "Command contains unquoted variable references; ensure values are properly escaped" when `$` outnumbers `{{`. Both emit `SeverityWarning`, not errors; the check is heuristic because the carrier is an opaque string.
   - `docs/reference/interpolation.md:605` "### Shell Injection — User-provided values in commands can be dangerous. AWF provides `ShellEscape()`". Escaping is offered to plugin authors, not applied by the engine.
   - `docs/ADR/012-runtime-shell-detection.md` (2026-03-02): "`ShellExecutor` hardcodes `/bin/sh -c`, which maps to `dash` on Debian/Ubuntu. Bash-dependent workflow commands ... fail silently or with cryptic errors." Fix: use `$SHELL`. This is the shell-string carrier producing portability defects, resolved by binding to a *different* shell rather than leaving the shell.
   - `docs/ADR/014-shebang-execution-for-script-files.md` (2026-03-04): "`ShellExecutor.Execute()` unconditionally wraps it in `exec.CommandContext(ctx, e.shellPath, "-c", cmd.Program)`, bypassing the kernel's shebang mechanism entirely" and "passing large script content as a `-c` argument risks hitting `ARG_MAX`". Fix: temp file + direct exec when a shebang is present. Second correction in two days, both caused by `-c` as the universal carrier.

## ms-conductor (Python; YAML multi-agent workflow)

The closest match to the RDR's shape: argv list, per-element templating, exec without shell, JSON-on-stdout re-entry.

1. Carrier: `command: str` plus `args: list[str]` (`src/conductor/config/schema.py:579-583`), executed with `asyncio.create_subprocess_exec(rendered_command, *rendered_args, ...)` (`src/conductor/executor/script.py`). No shell; a `FileNotFoundError` becomes `ExecutionError("command not found ...")`. Schema forbids `prompt`, `provider`, `model`, `tools` etc. on `type: script` (`schema.py:928-946`, `extra="forbid"`).
2. Templating: Jinja2 rendered per element — `rendered_args = [self.renderer.render(arg, context) for arg in agent.args]`. Substitution is intra-element (a template may expand inside an arg), not whole-element; `working_dir` is templated too. `${VAR:-default}` in `env:` is resolved by the config loader at parse time (`script.py` comment).
3. Env: `base_env = {**os.environ, "PYTHONUTF8": "1"}; env = {**base_env, **agent.env}` — inherit-all plus declared overrides.
4. Output: always structured — `schema.py:475` "Output is always `{stdout, stderr, exit_code}` with parsed-JSON keys merged on top when `stdout` is valid JSON." Engine (`src/conductor/engine/workflow.py:2679-2703`) builds `output_content = {"stdout", "stderr", "exit_code"}`, `json.loads(stdout)`, and `output_content.update(parsed_json)`; keys that collide with the built-ins log "Script '%s' JSON output shadows built-in fields". If the step declares `output:` schema, `_validate_script_output_schema` (`:951`) hard-fails on non-JSON stdout ("When 'output:' is declared, the script must write a JSON object to stdout. Write logs to stderr.") and on JSON arrays/scalars ("Arrays and scalars are not accepted."). Traced to "issue #118" in the docstring, i.e. added after use.
5. Exit codes: non-zero emits a `script_failed` event with `exit_code` (`workflow.py:2719-2730`); no per-code map. Per-script `timeout` seconds -> `process.kill()` + `ExecutionError("timed out after Ns")` (`script.py:128-137`), plus a workflow-level timeout wrapper (`workflow.py:821`).
6. ★ Hardening: the typed-stdout contract and the built-in-field shadowing warning are both later additions (issue #118 reference; `assert process.returncode is not None` with the comment "Do NOT use `process.returncode or 0` — 0 is falsy" records a real bug). No shell was ever offered, so no shell-related reversal exists.

## ralph-tui (TypeScript/Bun; agent orchestration loop with TOML config)

1. Carrier: `command: string` in `[[agents]]` and at top level (`src/config/schema.ts:112,166,185`), documented as an executable path or wrapper (`@example "ccr code"`). Spawned with `spawn(command, args, ...)` argv style (`src/utils/process.ts:49`); `shell: true` only on Windows for `.cmd`/`.bat` wrappers, with `quoteForWindowsShell` and `shouldUseWindowsShell` in `src/plugins/agents/base.ts:93-125`.
2. Templating: none in the command; flags come from `defaultFlags: string[]`.
3. Env: ★ default deny-list — `src/plugins/agents/base.ts:180-190` "Default environment variable exclusion patterns. These patterns are always excluded from agent subprocesses to prevent accidental API key leakage (e.g., from .env files auto-loaded by Bun) which can cause unexpected billing." `DEFAULT_ENV_EXCLUDE_PATTERNS = ['*_API_KEY', '*_SECRET_KEY', '*_SECRET']`, with `envExclude` / `envPassthrough` config knobs (`schema.ts:120-121`).
4. Output: agent stream parsed per plugin; not a model-level result protocol.
5. Exit codes: plugin-specific; `timeout` per agent.
6. ★ Reversal: `src/config/schema.ts:183-186` — `command: z.string().refine((cmd) => !/[;&|`$()]/.test(cmd), 'Command cannot contain shell metacharacters (;&|`$()). Use a wrapper script instead.')`. A load-time schema lint that bans shell syntax in the command carrier and tells the author to move logic into a script file. This is the same shape as the RDR's proposed `sh -c` argv0 ban. (Repo is a single-commit snapshot; the introducing commit cannot be dated locally.)

## statewright (Rust + plugins; JSON state machine constraining coding agents)

Does not shell out from the model. The model instead *constrains* the agent's own shell tool per state.

1. Carrier: `allowed_commands: ["pytest", "cargo test", "npm test"]` on a state (`README.md:169`), a prefix/glob allow-list, not an invocation.
2-5. Enforcement is a PreToolUse hook: `plugins/claude-code/hook.sh:250-275` — "if it's Bash, classify the command to prevent bypass of Write/Edit/Destructive restrictions via shell redirects"; destructive verbs (`rm|rmdir|shred|truncate|unlink`) are denied "always ... regardless of allowed_commands", then `case "$COMMAND" in $pattern)` glob matching against the state's list; the gateway exposes the list via `crates/mcp-gateway/src/custom_tools.rs:120,170`. README `:129` "Bash discernment | Blocks `echo > file`, `rm -rf`, `sed -i`, and scripting interpreters (`python`, `node`) when Write/Edit aren't allowed. Even if Bash itself is permitted."
6. ★ Signal: a whole guardrail layer (`hook.sh`, "Write-via-redirect and destructive ops are still blocked even when Bash is allowed", `README.md:17`) exists because a free-form shell string cannot be reasoned about statically; the `&&|;` chained-command check is the tell that string-level heuristics were patched after bypasses.

## temporal-cli (Go)

Model does not declare commands; the CLI dispatches unknown subcommands to `temporal-<name>` extensions. `internal/temporalcli/commands.extension.go:74` `cmd := exec.CommandContext(ctx, extPath, append(cliPassArgs, extArgs...)...)` — argv only, stdin/stdout/stderr passed through, `--command-timeout` via context; a child `*exec.ExitError` is swallowed (`return nil, true`) so the extension's exit code propagates unchanged. Prior art for argv-only + pass-through I/O, not for a workflow-declared command.

## Archon (TypeScript; workflow engine over Claude Code SDK)

Model does not declare shell commands; git operations use `execFile` argv (`packages/git/src/exec.ts:8-13` `execFileAsync(cmd, args, ...)`). The security note is about the agent's shell, not the model's: `docs/reference/security.md:14` "Archon runs the Claude Code SDK in `bypassPermissions` mode ... can read, write, and execute files without interactive confirmation". Env hardening: `:144` passes `--no-env-file` so Bun does not autoload `.env` into the spawned agent; per-codebase env is merged explicitly (`:132`). Negative for carrier prior art; positive for env-leak prevention.

## dbos-transact-py (Python; durable workflows)

Model (`dbos-config.yaml`) declares `database.migrate: [<shell strings>]`; `dbos/cli/cli.py:322` `subprocess.run(command, shell=True, text=True)`, non-zero prints "Migration command failed". A shell-string carrier with no templating, no env policy, and no output protocol — the migration list is operator config, not a workflow step. Workflow steps themselves are Python callables.

## AgentSpec / Pro2Guard (Python; rule-based agent guardrails)

No model-declared commands. Rules pattern-match the *agent's* proposed code for shell reach: `AgentSpec/src/rules/manual/pythonrepl.py:159` `re.compile(r'(?:os\.system|subprocess\.run)\s*\(\s*[\'"](.+?)[\'"]\s*\)')` and `:176-179` a deny list `"os.system", "subprocess.run", "subprocess.call", "subprocess.Popen"`. Pro2Guard's `GitCommand.py` shells out itself but only with fixed argv (`subprocess.check_output(["git", "branch", "--show-current"])`). Prior art for "shell reach is what the guardrail looks for", i.e. the RDR's lint target.

## BehaviorTree.CPP

Refuses the shell by design: `<Script code=" A:=42; B:=3.14 " />` runs an internal expression language over the blackboard (`include/behaviortree_cpp/actions/script_node.h:22-29`, "can use the BT++ scripting"). No process spawn anywhere in `src/` or `include/`. External work is a host-language `ActionNode`.

## state-machine-cat (TypeScript)

Rendering only: `src/render/vector/dot-to-vector-native.mts:33` `spawnSync(lOptions.exec, [...])` argv to `dot`; not part of the model.

## langgraph

`libs/cli/langgraph_cli/exec.py` wraps `docker` argv for the dev CLI; graph nodes and tools are Python callables. Negative.

## py_trees

`demos/action.py` / `demos/dot_graphs.py` use `subprocess` for demo rendering only; behaviours are Python classes. Negative.

## restate

Rust; `Command::new` only for `$EDITOR` (`cli/src/cli_env.rs:235`), lambda test harness, and the local cluster runner. Service handlers are SDK code. Negative.

## inngest

Only `tools/` build helpers and a redis test; steps are SDK closures. Negative.

## One-line negatives

- async-fsm, fsm, javascript-state-machine, jfsm, qmuntal-stateless, qpc, scxmlcc, xgrammar, outlines: zero triage hits; transitions bind host callbacks only.
- fizz, re2c, ragel-colm-suite, sml, StateSmith, lmql, guidance: hits are build scripts, benchmarks, docs tooling, or notebooks; the model/DSL has no command carrier.
- xstate: hits are `spawn(` for actors (not processes) and example/CHANGELOG text; no process exec.

## Reversal signals

Ranked by strength (how directly a project changed course on a command carrier after shipping it):

1. **ralph-tui `command` metacharacter ban** — `src/config/schema.ts:183-186`: `refine((cmd) => !/[;&|`$()]/.test(cmd), 'Command cannot contain shell metacharacters (;&|`$()). Use a wrapper script instead.')`. A load-time schema lint that forbids shell syntax in a string carrier and redirects authors to a script file. Direct analogue of the RDR's `sh -c` argv0 defect.
2. **awf-cli ADR-012 and ADR-014** — two accepted ADRs (2026-03-02, 2026-03-04) correcting the `$SHELL -c <string>` executor: dash-vs-bash silent failures, ignored shebangs, `ARG_MAX` truncation. Both keep the shell string but carve exceptions around it; the carrier itself was never replaced.
3. **awf-cli security-validator plugin** — `examples/plugins/awf-plugin-security-validator/main.go:116,141`: post-hoc heuristic lint (argv0 deny words, `$`-vs-`{{` counting) emitting warnings only, because an opaque shell string cannot be validated structurally.
4. **ms-conductor stdout contract (issue #118)** — `src/conductor/engine/workflow.py:951-1003`: after use, stdout must be a JSON *object* when `output:` is declared; arrays/scalars rejected; built-in key shadowing warned. The re-entry protocol was tightened, not the carrier (argv was correct from the start).
5. **statewright Bash discernment** — `plugins/claude-code/hook.sh:250-275` plus `README.md:129`: incremental string heuristics (`rm` at start, then after `&&|;`, redirects, interpreters) layered on top of `allowed_commands` — the signature of chasing bypasses in a free-form shell string.
6. **ralph-tui / Archon env deny-lists** — `base.ts:180-190` `*_API_KEY, *_SECRET*` excluded by default "to prevent accidental API key leakage ... unexpected billing"; Archon `--no-env-file`. Both are reactions to inherit-everything env policy.
7. **awf-cli `ShellEscape()` doc** — `docs/reference/interpolation.md:605`: acknowledges "Shell Injection" but leaves escaping to the caller; weakest signal (advice, not enforcement).

## Family disposition

Of 32 peers, only 4 let a model or declaration name an external command at all: awf-cli (shell string, always `$SHELL -c`), dbos-transact-py (shell string, `shell=True`, config-level migrations only), ms-conductor (argv `command` + `args`, `create_subprocess_exec`, no shell), and ralph-tui (executable string spawned argv-style with a metacharacter ban; `shell: true` only for Windows wrappers). So 2 carry shell strings, 2 carry argv, and both shell-string carriers accumulated corrective ADRs or lint afterwards while neither argv carrier did. Templating splits the same way: awf substitutes Go templates into the string before `-c`; ms-conductor renders Jinja per argv element (intra-element, so still not whole-element); ralph-tui and dbos do not template. Only ms-conductor has a typed output protocol (stdout JSON object merged over `{stdout, stderr, exit_code}`, schema-validated when declared); awf re-enters raw stdout as a string; none has a declared exit-code map beyond zero/non-zero plus awf's `retryable_exit_codes`. Env authority is bound in 1 of 4 (ralph-tui default deny patterns); awf and ms-conductor inherit the whole parent env plus overrides. A further 4 (statewright, AgentSpec, Pro2Guard, Archon) bind authority only over the *agent's* shell via allow-lists or deny-regexes, and BehaviorTree.CPP deliberately substitutes an internal expression language. The remaining 23 never leave the host language. Net: the family has no example of a whole-element placeholder, a closed argv0 lint, or a declared exit-code-to-verdict map; the RDR would be first in this set on all three, and the shell-string members are the ones that later paid for their carrier.
