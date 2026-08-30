Model: claude-fable-5

# RDR 0025 live spike — A3 (direct-binding class), A7 (env policy), A8 (exit-code mapping)

Host: macOS (Darwin 25.6.0), Go 1.26.6, git 2.55.0. Harness: `main.go` (module `a3spike`), invokes with
`exec.CommandContext`, no shell, `{artifact}` passed as one argv element, stdin JSON, stdout/stderr/exit captured,
C3 read envelope parsed BEFORE exit classification. `run.sh` builds + runs and tees to `output.txt`.
Artifact: git-config file `[state]\n\tphase = draft\n` in a temp dir.

## Per-scenario results (one line each)

- R1  `git config --file {a} --get state.phase` -> exit 0, stdout bytes exactly `"draft\n"`; C3 object parse FAILS; a RAW single-value mode (TrimRight `\n`, bind to the entry's one declared key) yields `state.phase=draft`.
- R2a `git config --file {a} --list` -> `"state.phase=draft\n"` (k=v lines); not a C3 object.
- R2b `wrapper-list.sh {a}` (1-line sh+awk over `--list`) -> `{"state.phase":"draft"}` ; C3 parse OK. Wrapper class: k=v -> JSON, trivial.
- R3a `git config --file {a} state.phase review` (value in argv) -> tool rewrites artifact to `phase = review`. But this argv is NOT expressible under C3: argv is fixed, only `{artifact}` is a placeholder, so the value cannot be in argv.
- R3b `git config --file {a} state.phase` with `{"state.phase":"review"}` on stdin -> git IGNORES stdin, treats the call as a get, prints `draft`, artifact unchanged. git config does not read stdin: a C3 WRITE to git-config is impossible without a wrapper.
- R3c `wrapper-write.sh {a}` reading stdin JSON, calling `git config --file $1 k v` per key -> artifact now `phase = review`; exit 0, empty stdout (WRITE success is not taken from exit; see R3d).
- R3d read-back via R1 -> `"review\n"`. Write verified through the role's reader, per C3.
- R4a P-1: `cat {a}` as WRITE with JSON stdin -> artifact intact ONLY because cat's stdout is the harness pipe, not the artifact (stdin is silently discarded / not even read).
- R4b P-1: `tee {a}` as WRITE with JSON stdin -> artifact body REPLACED by `{"state.phase":"review"}` (the envelope). Read-back R4c: exit 128 `fatal: bad config line 1`. This is the stdin-as-content corruption surface A3/F5 must name: any argv0 that sinks stdin into `{artifact}` will write the envelope into the artifact.
- R5 P-2: wrapper emitting a 2 MiB single-key JSON object -> capture works (2,097,164 bytes), parse OK, ~90 ms wall (88 ms this run; 386 ms on a cold first run). No bound exists today; propose a stdout cap (1 MiB) -> `execution_failure` when exceeded, applied at capture time (bytes.Buffer via io.LimitReader/ counting writer), not after parse.
- R6 read-absent: `git config --file {a} --get missing.key` -> exit 1, stdout `""`, stderr empty. Established-absent has a distinct (code, empty-stdout) signature.
- E1a `test -f {a}` exists -> 0.   E1b missing -> 1.   E1c `test --bogus` -> 2 (`unexpected operator`).
- E1d `git diff --quiet` clean -> 0.   E1e dirty -> 1.   E1f bad `-C` path -> 128 (`fatal: cannot change to`).
- E1g `git config --file <missing file> --get k` -> exit 1, empty stdout. NB: same signature as "key absent" (R6): git does not distinguish missing file from missing key at the exit-code level.
- E1h missing executable -> Go spawn error `exec: "...": executable file not found in $PATH`, `errors.Is(err, exec.ErrNotFound)==true`, NO exit code (ExitError not produced).
- V1  Env=[] (fully scrubbed), argv0 bare `git` -> works, `draft`. Go's exec.Command resolves argv0 with the PARENT's PATH at construction (LookPath), so an empty child env does not break spawn; the child itself has no PATH/HOME. `--list --show-origin` shows only `file:{a}` opened.
- V2  Env=[PATH] -> works, `draft`.
- V3a Env=[PATH,HOME=tmp] with a conflicting `~/.gitconfig` (`phase = FROM_HOME_GITCONFIG`) -> `--file` read returns `draft`: `--file` isolates from the global file. V3c control without `--file` returns `FROM_HOME_GITCONFIG`. V3b GIT_TRACE=1 shows only the built-in invocation (no hook/helper spawn). V3d `--show-origin` lists only the artifact.
- V4a GIT_CONFIG_GLOBAL=<tmpfile> -> `--file` read unaffected (`draft`); control without `--file` returns `FROM_GIT_CONFIG_GLOBAL`.
- V4b GIT_DIR=<repo>/.git -> `--file` read unaffected.
- V4c GIT_CONFIG_COUNT/KEY_0/VALUE_0 injection -> `--file --get` and `--get-all` unaffected (`draft`); control without `--file` returns `FROM_ENV_INJECTION`. `--file` short-circuits the config stack, including env-injected config.
- V4d GIT_CONFIG_PARAMETERS -> `--file` unaffected.  V4e GIT_CONFIG_SYSTEM -> unaffected.  V4f GIT_EXEC_PATH=<tmp> -> unaffected for a built-in (would matter for dashed externals / helpers).
- Net for A7: for THIS tool (`git config --file`) the environment does not change the result at all — but only because `--file` is an explicit isolation flag. Every control run (no `--file`) shows the same binary is fully steerable by inherited env (HOME, GIT_CONFIG_GLOBAL, GIT_CONFIG_COUNT...). The ambient-authority risk is real and is a property of the argv, not of the binary.

## Per tool class

| class (example)                          | direct (C3 object) | raw (single value) | exit-mapped | wrapped | what the wrapper did |
|------------------------------------------|--------------------|--------------------|-------------|---------|----------------------|
| VCS query, one key (`git config --get`)  | no (R1)            | YES: trim `\n`, bind to declared key; absent = exit 1 + empty stdout (R6) | absent only | not needed | — |
| VCS query, many keys (`git config --list`) | no (R2a)         | no (multi-line)    | —           | YES (R2b) | awk: `k=v` lines -> JSON object (1 line) |
| file probe (`test -f`)                   | no (empty stdout)  | no                 | YES: 0/1 verdict, 2 = execution failure (E1a-c) | not needed | — |
| status verb (`git diff --quiet`)         | no (empty stdout)  | no                 | YES: 0/1 verdict, 128 = execution failure (E1d-f) | not needed | — |
| VCS write (`git config k v`)             | no (R3b: ignores stdin) | no            | no          | REQUIRED (R3c) | parse stdin JSON, one `git config --file $1 k v` per key; 3 lines of sh (jq-free) |
| stdin-sinking write (`tee {artifact}`)   | CORRUPTS (R4b)     | —                  | —           | —       | must be refused/named as a hazard, not admitted |

A3 "wrapper class is small and mechanical": HOLDS for reads (R2b: one awk line) and for the write side the wrapper is
short (R3c) but it is NOT optional — every real-tool WRITE whose data is not in argv needs one, because C3 pins the
data to stdin and the surveyed tool does not read stdin. So the claim narrows to: "direct binding covers only
exit-mapped probes/status verbs and raw single-value reads; every write to a real tool is wrapped."

A8: evidence table of codes

| tool / condition                | code        | meaning         |
|---------------------------------|-------------|-----------------|
| `test -f` exists / missing      | 0 / 1       | verdict         |
| `test --bogus`                  | 2           | execution failure |
| `git diff --quiet` clean/dirty  | 0 / 1       | verdict         |
| `git -C <bad>`                  | 128         | execution failure |
| `git config --get` present/absent | 0 / 1    | value / established-absent |
| `git config --file <missing>`   | 1           | (ambiguous with absent) |
| `git config` bad config file    | 128         | execution failure |
| missing executable              | none (exec.ErrNotFound) | spawn failure |

A declared per-entry map `exit_verdicts = { "0" = "allow", "1" = "deny" }` with unlisted -> `execution_failure`
is SUFFICIENT for both probes: their error codes (2, 128) are outside the declared set and the spawn error has no
code at all, so nothing launders into a verdict. For reads, `exit_absent = [1]` cleanly expresses established-absent
for `git config --get` (exit 1 + empty stdout). Caveat (E1g): git returns 1 for a MISSING FILE too; if "artifact
missing" must differ from "key absent", the runtime must stat `{artifact}` before spawning (cheap, and it already
has the path) rather than expect the tool to tell them apart. Require empty stdout alongside the absent code;
non-empty stdout with an absent code is a contract violation -> `execution_failure`.

## Recommendations

1. RAW read mode grammar (A3). Add a per-entry `read_mode = "raw"` (default `"json"`). Raw is legal only when the
   entry declares exactly one key; the runtime takes stdout, strips one trailing `\n` (and `\r\n`), binds it to
   that key; empty stdout with exit in `exit_absent` = not-carried; empty stdout with exit 0 = carried empty
   string. Multi-line stdout in raw mode = `execution_failure` (it is `--list`-shaped, use a wrapper). Evidence:
   R1 (`"draft\n"`), R6 (exit 1, `""`).
2. Exit-mapping shape (A8). `exit_verdicts = { "<code>" = "<verdict>" }` for GATE entries and
   `exit_absent = [<codes>]` for READ entries; any code not listed, any signal death, and any spawn error
   (`exec.ErrNotFound`, EACCES) -> `execution_failure`. Envelope on stdout is still parsed first; a GATE that emits
   a verdict envelope AND has `exit_verdicts` is a config error (pick one). Evidence: E1a-h (verdict codes 0/1;
   error codes 2/128/none are disjoint for both tools surveyed).
3. Env policy (A7). Allowlist, not inherit: child env = {PATH, HOME, TMPDIR, LANG, LC_*} copied from the parent,
   plus explicit per-entry `env = { K = "v" }` and `env_pass = ["GIT_CONFIG_GLOBAL", ...]` opt-ins. PATH is needed
   for wrappers that call tools (R2b/R3c call `git` by name) though not for spawn itself (V1: Go resolves argv0
   with the parent PATH). HOME is harmless for `--file` (V3a) and needed by tools that consult it legitimately.
   Ground: V4a/V4c controls show the same binary flips its answer on GIT_CONFIG_GLOBAL and GIT_CONFIG_COUNT when
   the argv lacks `--file`; an inherit-as-is policy makes every entry's determinism depend on the operator's
   shell, whereas an allowlist makes the ambient-authority surface explicit and reviewable in the TOML.
   Also name the P-1 hazard in F5: an argv0 that sinks stdin to `{artifact}` (tee, `cp /dev/stdin`, `dd of=`)
   overwrites the artifact with the envelope (R4b, then R4c fatal 128).

## Appendix: output.txt

```
spike: rdr-0025 a3-a7-a8-real-tools
date: 2026-08-30T18:38:10Z
go: go version go1.26.6 darwin/arm64
git: git version 2.55.0
os: Darwin 25.6.0
---
artifact=$TMPDIR/a3spike1497112283/state.cfg
git=/opt/local/bin/git

################ A3 direct-binding class
---
R1 read: git config --get (single value)
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "--get" "state.phase"]
  exit=0
  stdout="draft\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'd' looking for beginning of value)
  dur=14.408ms
  RAW-mode bind: key=state.phase value="draft" (TrimRight \n)
---
R2a read: git config --list (raw k=v lines)
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "--list"]
  exit=0
  stdout="state.phase=draft\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 's' looking for beginning of value)
  dur=11.182ms
---
R2b read: wrapper-list.sh (k=v -> JSON)
  argv=["$SPIKE/wrapper-list.sh" "$TMPDIR/a3spike1497112283/state.cfg"]
  exit=0
  stdout="{\"state.phase\":\"draft\"}\n"
  c3_read_envelope: OK keys=1
  dur=15.354ms
---
R3a write: git config set (value in argv -- NOT expressible in fixed argv)
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "state.phase" "review"]
  exit=0
  stdout=""
  c3_read_envelope: NOT PARSEABLE (unexpected end of JSON input)
  dur=8.741ms
  artifact_after="[state]\n\tphase = review\n"
---
R3b write: git config with data on stdin only (does git read stdin?)
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "state.phase"]
  stdin="{\"state.phase\":\"review\"}"
  exit=0
  stdout="draft\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'd' looking for beginning of value)
  dur=7.307ms
  artifact_after="[state]\n\tphase = draft\n"
---
R3c write: wrapper-write.sh (stdin JSON -> git config per key)
  argv=["$SPIKE/wrapper-write.sh" "$TMPDIR/a3spike1497112283/state.cfg"]
  stdin="{\"state.phase\":\"review\"}"
  exit=0
  stdout=""
  c3_read_envelope: NOT PARSEABLE (unexpected end of JSON input)
  dur=16.778ms
  artifact_after="[state]\n\tphase = review\n"
---
R3d read-back via R1
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "--get" "state.phase"]
  exit=0
  stdout="review\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'r' looking for beginning of value)
  dur=6.56ms
---
R4a P-1: cat {artifact} as WRITE with JSON stdin
  argv=["cat" "$TMPDIR/a3spike1497112283/state.cfg"]
  stdin="{\"state.phase\":\"review\"}"
  exit=0
  stdout="[state]\n\tphase = draft\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 's' looking for beginning of value)
  dur=2.002ms
  artifact_after="[state]\n\tphase = draft\n"
---
R4b P-1: tee {artifact} as WRITE with JSON stdin
  argv=["tee" "$TMPDIR/a3spike1497112283/state.cfg"]
  stdin="{\"state.phase\":\"review\"}"
  exit=0
  stdout="{\"state.phase\":\"review\"}"
  c3_read_envelope: OK keys=1
  dur=1.964ms
  artifact_after="{\"state.phase\":\"review\"}"
---
R4c read-back after tee
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "--get" "state.phase"]
  exit=128
  stdout=""
  stderr="fatal: bad config line 1 in file $TMPDIR/a3spike1497112283/state.cfg\n"
  c3_read_envelope: NOT PARSEABLE (unexpected end of JSON input)
  dur=6.449ms
---
R5 P-2: 2 MiB JSON envelope
  argv=["$SPIKE/wrapper-big.sh" "$TMPDIR/a3spike1497112283/state.cfg"]
  exit=0
  stdout_len=2097164 stdout_head="{\"blob\":\"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  c3_read_envelope: OK keys=1
  dur=87.695ms
---
R6 read-absent: git config --get missing.key
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "--get" "missing.key"]
  exit=1
  stdout=""
  c3_read_envelope: NOT PARSEABLE (unexpected end of JSON input)
  dur=6.939ms

################ A8 exit-code mapping
---
E1a test -f (exists)
  argv=["test" "-f" "$TMPDIR/a3spike1497112283/state.cfg"]
  exit=0
  stdout=""
  c3_read_envelope: NOT PARSEABLE (unexpected end of JSON input)
  dur=1.985ms
---
E1b test -f (missing)
  argv=["test" "-f" "$TMPDIR/a3spike1497112283/nope"]
  exit=1
  stdout=""
  c3_read_envelope: NOT PARSEABLE (unexpected end of JSON input)
  dur=1.767ms
---
E1c test with bogus option (execution failure)
  argv=["test" "--bogus" "$TMPDIR/a3spike1497112283/state.cfg"]
  exit=2
  stdout=""
  stderr="test: --bogus: unexpected operator\n"
  c3_read_envelope: NOT PARSEABLE (unexpected end of JSON input)
  dur=1.787ms
---
E1d git diff --quiet (clean)
  argv=["git" "-C" "$TMPDIR/a3spike1497112283/repo" "diff" "--quiet"]
  exit=0
  stdout=""
  c3_read_envelope: NOT PARSEABLE (unexpected end of JSON input)
  dur=6.924ms
---
E1e git diff --quiet (dirty)
  argv=["git" "-C" "$TMPDIR/a3spike1497112283/repo" "diff" "--quiet"]
  exit=1
  stdout=""
  c3_read_envelope: NOT PARSEABLE (unexpected end of JSON input)
  dur=6.005ms
---
E1f git with bad path (execution failure)
  argv=["git" "-C" "$TMPDIR/a3spike1497112283/norepo" "diff" "--quiet"]
  exit=128
  stdout=""
  stderr="fatal: cannot change to '$TMPDIR/a3spike1497112283/norepo': No such file or directory\n"
  c3_read_envelope: NOT PARSEABLE (unexpected end of JSON input)
  dur=5.327ms
---
E1g git config --file <missing file> --get
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/nofile.cfg" "--get" "state.phase"]
  exit=1
  stdout=""
  c3_read_envelope: NOT PARSEABLE (unexpected end of JSON input)
  dur=6.34ms
---
E1h missing executable (spawn error)
  argv=["definitely-not-a-binary-xyz" "$TMPDIR/a3spike1497112283/state.cfg"]
  spawn_error=exec: "definitely-not-a-binary-xyz": executable file not found in $PATH  is_ErrNotFound=true  exit=(none)
  stdout=""
  c3_read_envelope: NOT PARSEABLE (unexpected end of JSON input)
  dur=0s

################ A7 env policy
---
V1 Env=[] (scrubbed; bare argv0 'git')
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "--get" "state.phase"]
  env=[]
  exit=0
  stdout="draft\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'd' looking for beginning of value)
  dur=8.342ms
---
V1b Env=[] --list --show-origin (which files does git open?)
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "--list" "--show-origin"]
  env=[]
  exit=0
  stdout="file:$TMPDIR/a3spike1497112283/state.cfg\tstate.phase=draft\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'i' in literal false (expecting 'a'))
  dur=7.955ms
---
V2 Env=[PATH]
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "--get" "state.phase"]
  env=["PATH=<parent PATH>"]
  exit=0
  stdout="draft\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'd' looking for beginning of value)
  dur=7.762ms
---
V3a Env=[PATH,HOME=tmp w/ conflicting ~/.gitconfig] --get
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "--get" "state.phase"]
  env=["PATH=<parent PATH>" "HOME=$TMPDIR/a3spike1497112283/home"]
  exit=0
  stdout="draft\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'd' looking for beginning of value)
  dur=8.524ms
---
V3b same + GIT_TRACE=1
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "--get" "state.phase"]
  env=["PATH=<parent PATH>" "HOME=$TMPDIR/a3spike1497112283/home" "GIT_TRACE=1"]
  exit=0
  stdout="draft\n"
  stderr="11:38:11.690948 git.c:502               trace: built-in: git config --file $TMPDIR/a3spike1497112283/state.cfg --get state.phase\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'd' looking for beginning of value)
  dur=7.769ms
---
V3c control: NO --file, HOME=tmp -> reads ~/.gitconfig
  argv=["git" "config" "--get" "state.phase"]
  env=["PATH=<parent PATH>" "HOME=$TMPDIR/a3spike1497112283/home"]
  exit=0
  stdout="FROM_HOME_GITCONFIG\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'F' looking for beginning of value)
  dur=7.682ms
---
V3d --file --list --show-origin under HOME=tmp
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "--list" "--show-origin"]
  env=["PATH=<parent PATH>" "HOME=$TMPDIR/a3spike1497112283/home"]
  exit=0
  stdout="file:$TMPDIR/a3spike1497112283/state.cfg\tstate.phase=draft\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'i' in literal false (expecting 'a'))
  dur=7.724ms
---
V4a GIT_CONFIG_GLOBAL set, --file --get
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "--get" "state.phase"]
  env=["PATH=<parent PATH>" "HOME=$TMPDIR/a3spike1497112283/home" "GIT_CONFIG_GLOBAL=$TMPDIR/a3spike1497112283/global.cfg"]
  exit=0
  stdout="draft\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'd' looking for beginning of value)
  dur=8.333ms
---
V4a' GIT_CONFIG_GLOBAL set, NO --file --get (control)
  argv=["git" "config" "--get" "state.phase"]
  env=["PATH=<parent PATH>" "HOME=$TMPDIR/a3spike1497112283/home" "GIT_CONFIG_GLOBAL=$TMPDIR/a3spike1497112283/global.cfg"]
  exit=0
  stdout="FROM_GIT_CONFIG_GLOBAL\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'F' looking for beginning of value)
  dur=8.504ms
---
V4b GIT_DIR set, --file --get
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "--get" "state.phase"]
  env=["PATH=<parent PATH>" "HOME=$TMPDIR/a3spike1497112283/home" "GIT_DIR=$TMPDIR/a3spike1497112283/repo/.git"]
  exit=0
  stdout="draft\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'd' looking for beginning of value)
  dur=7.942ms
---
V4c GIT_CONFIG_COUNT injection, --file --get
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "--get" "state.phase"]
  env=["PATH=<parent PATH>" "HOME=$TMPDIR/a3spike1497112283/home" "GIT_CONFIG_COUNT=1" "GIT_CONFIG_KEY_0=state.phase" "GIT_CONFIG_VALUE_0=FROM_ENV_INJECTION"]
  exit=0
  stdout="draft\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'd' looking for beginning of value)
  dur=8.092ms
---
V4c' GIT_CONFIG_COUNT injection, --file --get-all
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "--get-all" "state.phase"]
  env=["PATH=<parent PATH>" "HOME=$TMPDIR/a3spike1497112283/home" "GIT_CONFIG_COUNT=1" "GIT_CONFIG_KEY_0=state.phase" "GIT_CONFIG_VALUE_0=FROM_ENV_INJECTION"]
  exit=0
  stdout="draft\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'd' looking for beginning of value)
  dur=7.94ms
---
V4c'' GIT_CONFIG_COUNT injection, NO --file (control)
  argv=["git" "config" "--get" "state.phase"]
  env=["PATH=<parent PATH>" "HOME=$TMPDIR/a3spike1497112283/home" "GIT_CONFIG_COUNT=1" "GIT_CONFIG_KEY_0=state.phase" "GIT_CONFIG_VALUE_0=FROM_ENV_INJECTION"]
  exit=0
  stdout="FROM_ENV_INJECTION\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'F' looking for beginning of value)
  dur=7.983ms
---
V4d GIT_CONFIG_PARAMETERS, --file --get
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "--get" "state.phase"]
  env=["PATH=<parent PATH>" "HOME=$TMPDIR/a3spike1497112283/home" "GIT_CONFIG_PARAMETERS='state.phase=FROM_GIT_CONFIG_PARAMETERS'"]
  exit=0
  stdout="draft\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'd' looking for beginning of value)
  dur=7.614ms
---
V4e GIT_CONFIG_SYSTEM, --file --get
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "--get" "state.phase"]
  env=["PATH=<parent PATH>" "HOME=$TMPDIR/a3spike1497112283/home" "GIT_CONFIG_SYSTEM=$TMPDIR/a3spike1497112283/global.cfg" "GIT_CONFIG_NOSYSTEM=0"]
  exit=0
  stdout="draft\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'd' looking for beginning of value)
  dur=7.871ms
---
V4f GIT_EXEC_PATH=tmp (hijack helper dir), --file --get
  argv=["git" "config" "--file" "$TMPDIR/a3spike1497112283/state.cfg" "--get" "state.phase"]
  env=["PATH=<parent PATH>" "HOME=$TMPDIR/a3spike1497112283/home" "GIT_EXEC_PATH=$TMPDIR/a3spike1497112283"]
  exit=0
  stdout="draft\n"
  c3_read_envelope: NOT PARSEABLE (invalid character 'd' looking for beginning of value)
  dur=7.497ms
---
done
harness_exit=0
```
