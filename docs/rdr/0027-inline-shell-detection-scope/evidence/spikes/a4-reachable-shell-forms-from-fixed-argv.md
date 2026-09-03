Model: claude-sonnet-5

# A4 verification spike: reachable shell forms from a fixed argv (no shell)

Host: darwin/zsh. Go toolchain used to build a scratch harness that execs a
fixed argv vector with `os/exec` directly (no `/bin/sh -c` wrapper), wiring
stdout/stderr through and optionally feeding stdin.

## Harness source (scratch, deleted after run)

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	args := os.Args[1:]
	var stdinData string
	if len(args) > 0 && args[0] == "--stdin" {
		stdinData = args[1]
		args = args[2:]
	}
	if len(args) == 0 {
		fmt.Println("usage: harness [--stdin DATA] argv...")
		os.Exit(2)
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if stdinData != "" {
		cmd.Stdin = strings.NewReader(stdinData)
	}
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "HARNESS_ERROR: %v\n", err)
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		}
		os.Exit(127)
	}
}
```

Built with `go build -o harness main.go`. Executed via `./harness [--stdin DATA] argv...`
— every probe argv element below is a genuinely separate Go-level argv word
(no shell re-splitting on the harness's own command line, verified by
inspecting harness output where ambiguous).

## Probe results

### 1. `["sh","-c","echo HELLO_C"]` — baseline inline code, separate words

```
$ ./harness sh -c "echo HELLO_C"
HELLO_C
exit=0
```
Shell ran, executed caller code. This is what the existing lint already refuses (argv[i]=sh, argv[i+1]=-c as separate words — trivially caught).

### 2. `["env","-S","sh -c echo HELLO_S"]` — one-word shell string

Darwin's native `/usr/bin/env` (BSD env) does **not** support `-S`:
```
$ env --version
env: illegal option -- e
```
`./harness env -S "sh -c echo HELLO_S"` → BSD env rejects `-S` outright (illegal option), no shell runs, no caller code executes. So on stock darwin with `/usr/bin/env`, this exact spelling is NOT reachable.

However, GNU coreutils `env` (found at `/opt/local/bin/genv`, MacPorts naming; GNU env 9.11) **does** support `-S` and behaves as advertised:
```
$ /opt/local/bin/genv --debug -S "sh -c echo HELLO_S"
split -S:  ‘sh -c echo HELLO_S’
 into:    ‘sh’ & ‘-c’ & ‘echo’ & ‘HELLO_S’
```
Note: `-S` does dumb whitespace splitting, not shell-quote-aware splitting — `sh -c echo HELLO_S` splits into 4 words, so `-c`'s operand becomes just `echo` and `HELLO_S` lands as `$0`, producing a blank echo. Properly quoting the embedded code as its own shell-quoted token works as intended:
```
$ ./harness /opt/local/bin/genv -S "sh -c 'echo HELLO_S'"
HELLO_S
exit=0
```
Confirmed: with GNU env present under some name/PATH entry, `env -S "sh -c '<code>'"` is a single argv word that causes a shell to execute caller-supplied code. This is exactly C1's named form (a) — a one-word shell string, here sourced via env's own re-splitting rather than a tool doing the splitting itself. Classified as (a), not a new form. Portability note: not reachable via darwin's stock `/usr/bin/env`; reachable if a GNU-coreutils `env` is anywhere on PATH (common on machines with Homebrew/MacPorts coreutils or on Linux, where GNU env is the default).

### 3. `["sh","-s"]` with stdin `"echo HELLO_STDIN"` — stdin-fed shell

```
$ ./harness --stdin "echo HELLO_STDIN" sh -s
HELLO_STDIN
exit=0
```
Confirmed form (b).

### 4. `["sh"]` bare, with stdin `"echo HELLO_BARE"`

```
$ ./harness --stdin "echo HELLO_BARE" sh
HELLO_BARE
exit=0
```
Confirmed form (b), bare-sh variant.

### 5. `["sh","<","script.sh"]` — redirection as literal argv words

```
$ echo 'echo SHOULD_NOT_RUN' > script.sh
$ ./harness sh "<" script.sh
sh: <: No such file or directory
HARNESS_ERROR: exit status 127
exit=127
```
Confirmed UNREACHABLE: with no shell doing the exec, `<` is passed to `sh` as a literal filename argument, not interpreted as redirection. `sh` treats `<` as a script-file operand and fails immediately (`No such file or directory`), never touching `script.sh` or running any inline code. Matches the claim exactly.

### 6. `["sh","./gate.sh"]` — sanctioned wrapper-file form

```
$ printf '#!/bin/sh\necho HELLO_FILE\n' > gate.sh && chmod +x gate.sh
$ ./harness sh ./gate.sh
HELLO_FILE
exit=0
```
Confirmed: a shell runs, but it runs the declared file `./gate.sh`, not inline code supplied via argv. Sanctioned form, as claimed.

### 7. `["sh","-c"]` — degenerate, no code operand

```
$ ./harness sh -c
sh: -c: option requires an argument
exit=2
```
Fails cleanly; no code to run since there is none. No adversarial value.

## Adversarial round — hunting for a reachable third form

### ADV-1: bundled short flag glued to operand — `sh -cecho HELLO`

```go
cmd := exec.Command("sh", "-cecho HELLO_BUNDLED")
```
```
sh: - : invalid option
[bash getopts dump]
ERR: exit status 1
```
Darwin's `/bin/sh` (bash in posix mode) does NOT accept `-c<code>` glued as one token — it errors as an invalid combined option before reaching any code path that would run `echo HELLO_BUNDLED`. No shell code executes. Also, per the rule, the predicate matches argv[j] EXACTLY against the flag string `-c`; `-cecho HELLO_BUNDLED` as a single argv word is not string-equal to `-c`, so the predicate would NOT admit it even if it did work. Doubly non-viable: rejected by the predicate AND rejected by sh itself.

### ADV-2: `sh -es` — combined short flags including `-s`

```
$ ./harness --stdin 'echo HELLO_ES' sh -es
HELLO_ES
exit=0
```
Runs, and does execute stdin-fed code. But this is still stdin-fed — category (b), just spelled as a bundled `-es` instead of bare `-s`. Not a new *form* of code delivery (code still arrives via stdin, not via argv). Also note: whether the predicate's "argv[j] exactly one of that key's flags" matches `-es` depends on whether `-es` is in the enumerated flag set for `sh`; if the flag list only enumerates `-c`/`-s` exactly, `-es` would NOT be matched (predicate misses it) — but this doesn't matter for A4 since the delivery channel (stdin) is unchanged from listed form (b).

### ADV-3: `python -` / `node -` (stdin-fed script interpreters)

```
$ ./harness --stdin 'print("HELLO_PY_STDIN")' python3 -
HELLO_PY_STDIN
exit=0
```
Confirmed runnable, but this is stdin-fed code for a non-shell interpreter — same *category* as (b) (stdin-fed), just a different interpreter than `sh`. Not a new form under the argv-word-spelling axis A4 is about (no inline code appears as a separate argv word; code arrives over stdin exactly as in probe 3/4).

### ADV-4: `/bin/sh -c "..."` — absolute path

```
$ ./harness /bin/sh -c "echo HELLO_ABS"
HELLO_ABS
exit=0
```
Runs and executes caller code — but this is form (a)'s sibling, plain multi-word `-c` (same as probe 1), just with an absolute path instead of a bare name. `filepath.Base("/bin/sh") == "sh"`, so the base() rule in the predicate still catches it — confirmed not evasive, no new form.

### ADV-5: renamed/symlinked interpreter (unlisted spelling)

```
$ cp /bin/sh /tmp/notashell
$ ./harness /tmp/notashell -c "echo HELLO_RENAMED"
HARNESS_ERROR: signal: killed
exit=255
```
On this darwin host, copying `/bin/sh` produces a binary that is killed immediately on exec (SIP/code-signing enforcement blocks unsigned copies of protected system binaries). Could not directly execute this variant here. This is immaterial to the classification: per the task framing, an unlisted spelling of a real shell binary is the *known-open* case that C1's deny-list already admits by design (the deny-list is a fixed enumeration of names; anything not on it is out of scope for the check, not a counterexample to A4). Confirmed by inspection of the predicate structure (name-based matching against a fixed list) rather than by a runnable darwin repro.

## Classification table

| # | Candidate | Predicate admits? | Shell/interpreter executes? | Executes caller-supplied code? | Category |
|---|---|---|---|---|---|
| 1 | `sh -c "code"` (separate words) | No — this is exactly what the lint refuses | yes | yes | Already blocked (not "admitted") |
| 2 | `env -S "sh -c 'code'"` (GNU env only; not stock darwin) | admits (one word, no listed interpreter+flag pair as separate argv words) | yes | yes | (a) one-word shell string — named in C1 |
| 3 | `sh -s` + stdin | admits (no inline code word) | yes | yes | (b) stdin-fed — named in C1 |
| 4 | bare `sh` + stdin | admits | yes | yes | (b) stdin-fed — named in C1 |
| 5 | `sh "<" script.sh` (argv words) | n/a | no (errors, no such file `<`) | no | Unreachable, confirmed |
| 6 | `sh ./gate.sh` | admits | yes, runs a *file* | no (runs declared file, not inline code) | Sanctioned wrapper-file form |
| 7 | `sh -c` (no operand) | n/a | errors | no | Degenerate, non-viable |
| ADV-1 | `sh -cecho HELLO` (bundled) | NOT admitted (argv word != "-c" exactly) | no — sh itself rejects it | no | Non-viable (double failure) |
| ADV-2 | `sh -es` + stdin | admits (possibly, depending on flag enumeration) | yes | yes, but via stdin | Still (b), same channel |
| ADV-3 | `python -` / `node -` + stdin | admits (non-sh interpreter) | yes | yes, but via stdin | Still (b)-shaped (stdin channel), different interpreter |
| ADV-4 | `/bin/sh -c "code"` (absolute path) | admits — wait, no: base()="sh" so predicate treats it same as #1, i.e. it IS the pattern the lint targets | yes | yes | Same as (1)/(a)-adjacent — base() rule confirmed to still catch absolute paths |
| ADV-5 | renamed/symlinked `sh` binary | NOT admitted (unlisted name) | untestable here (SIP kill) on darwin; would run on a permissive host | would execute code if it ran | Known-open-set case by design, not a new form |

## Verdict

No reachable third form was found that is (i) admitted by the argv-word predicate, (ii) actually executes a shell running caller-supplied code, and (iii) is not already one of C1's named forms (a) one-word shell string or (b) stdin-fed shell, and is not the by-design open deny-list case (unlisted spelling). Every adversarial candidate collapsed into: already-blocked (1, ADV-4), unreachable (5), non-viable/rejected-by-the-tool (7, ADV-1), a same-channel variant of (b) stdin-fed (ADV-2, ADV-3), or the known-open unlisted-spelling case (ADV-5). The one portability wrinkle found: `env -S` as form (a) requires a GNU-coreutils `env`; darwin's stock BSD `env` does not support `-S` at all, so on a clean darwin host this specific sub-spelling of (a) is not reachable unless GNU env is installed and earlier on PATH — a fact worth noting for RDR scope/testing but it does not change the admitted-and-executes set relative to C1, since C1 already lists the *form*, not a per-platform reachability guarantee.

The admitted-and-actually-executes set equals C1's named list, modulo the by-design open deny-list. A4 is verified.
