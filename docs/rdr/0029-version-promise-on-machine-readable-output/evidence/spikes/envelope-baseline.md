Model: claude-sonnet-5

## Build: go version

```
go version
```

```
go version go1.27.1 darwin/arm64
```

## Build: go build ./... (whole-module build, exit code)

```
cd /Users/cwensel/sandbox/newcoinc/intrastate && go build ./...
echo "exit=$?"
```

```
exit=0
```

## Build: binary to temp path

The repo's `Makefile` has a `build` target (`build: @mkdir -p $(BIN_DIR); $(GO) build -trimpath -ldflags '$(LDFLAGS)' -o $(BIN) $(PKG)` with `BIN_DIR=bin`, `PKG=./cmd/intrastate`), which stamps version/commit/date via `-ldflags`. For this spike a plain `go build` was used instead (no ldflags stamping) so the `version --as=json` capture below shows the fallback identity a stampless build produces, not a release identity.

```
cd /Users/cwensel/sandbox/newcoinc/intrastate && go build -o /tmp/intrastate-spike ./cmd/intrastate
echo "exit=$?"
```

```
exit=0
```

## intrastate --help (root usage, to find the json-emitting verb)

```
/tmp/intrastate-spike --help
```

```
intrastate makes workflow state transitions explicit, reviewable, and
deterministic. A flow is authored once as a transition model — a TOML
document of tags, rules, guards, and writes — and every caller navigates
it by asking this CLI, rather than reimplementing the flow in a skill,
script, or agent.

Two surfaces:

  lint    check a model against the graph invariants, before runtime
  flow    drive a model at runtime: what can happen, what does happen,
          and what state was written

Author a model, then check it:

  intrastate lint --model flow.toml

Ask what the current state can do, then commit one outcome:

  intrastate flow next     --model flow.toml --artifact state=state.json
  intrastate flow resolve  --model flow.toml --artifact state=state.json \
      --outcome approved
  intrastate flow set-state --model flow.toml --artifact state=state.json \
      --write status=approved

Given the same model, state, tags, and recognized outcome, the answer is
the same every time: one legal plan, or one typed refusal. intrastate
never guesses which of two matching rules you meant.

Global flags:

  --as text|json    output mode (default text)
  --help-all        the extended reference for any command

Run any subcommand with --help for its flags, or --help-all for its
extended reference. `intrastate help --all` prints the full reference for
every command at once.

Usage:
  intrastate [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  flow        Drive a transition model from a skill
  help        Help about any command
  lint        Check a transition model against the graph invariants
  version     Print build version, commit, and date

Flags:
      --as string   output mode: text | json (default "text")
  -h, --help        help for intrastate
      --help-all    show extended help (vocabulary, wire shapes, exit codes)
  -v, --version     version for intrastate

Use "intrastate [command] --help" for more information about a command.
```

## intrastate lint --help (flags, exit-code documentation)

```
/tmp/intrastate-spike lint --help
```

```
Check a transition model against the mandatory graph invariants.

A model with any blocking finding is refused: the command returns the
aggregate error graph-lint-failed at exit 2 and carries every
blocking finding in the machine-readable findings list.

Bounds enforced by this build (model-independent implementation
constants, not per-model inputs):

  product bound   2048
  node ceiling    4096

Usage:
  intrastate lint [flags]

Flags:
      --flow string    flow id to lint (reserved; this build resolves none — use --model)
  -h, --help           help for lint
      --help-all       show extended help (vocabulary, wire shapes, exit codes)
      --model string   path to the transition model to lint

Global Flags:
      --as string   output mode: text | json (default "text")
```

## intrastate version --help (states the version/json relationship in prose)

```
/tmp/intrastate-spike version --help
```

```
Print this build's version, commit, and date.

Under --as=json the same identity is emitted as the structured Info
value under the terminal "ok" envelope, so a caller can pin the build
that produced any other output.

Usage:
  intrastate version [flags]

Flags:
  -h, --help       help for version
      --help-all   show extended help (vocabulary, wire shapes, exit codes)

Global Flags:
      --as string   output mode: text | json (default "text")
```

## SUCCESS envelope 1: intrastate lint --model models/rdr.toml --as=json (raw stdout)

The model fixture used is the repo's own `models/rdr.toml` — the same model `make graph-lint` lints (`MODEL ?= models/rdr.toml` in the Makefile), so this is not a spike-authored fixture, it is the model the project already lints in CI.

```
cd /Users/cwensel/sandbox/newcoinc/intrastate && /tmp/intrastate-spike lint --model models/rdr.toml --as=json
echo "exit=$?"
```

```
{"type":"ok","data":{"findings":[]}}
exit=0
```

## SUCCESS envelope 1: pretty-printed

```
/tmp/intrastate-spike lint --model models/rdr.toml --as=json | jq .
```

```
{
  "type": "ok",
  "data": {
    "findings": []
  }
}
```

## SUCCESS envelope 1: top-level key set

```
/tmp/intrastate-spike lint --model models/rdr.toml --as=json | jq -r 'keys|join(",")'
```

```
data,type
```

## SUCCESS envelope 2: intrastate version --as=json (raw stdout)

```
/tmp/intrastate-spike version --as=json
echo "exit=$?"
```

```
{"type":"ok","data":{"version":"v0.0.0-20260912002329-ba40b6174344","commit":"ba40b6174344","date":"2026-09-12T00:23:29Z"}}
exit=0
```

## SUCCESS envelope 2: pretty-printed

```
/tmp/intrastate-spike version --as=json | jq .
```

```
{
  "type": "ok",
  "data": {
    "version": "v0.0.0-20260912002329-ba40b6174344",
    "commit": "ba40b6174344",
    "date": "2026-09-12T00:23:29Z"
  }
}
```

Note: this build was `go build -o /tmp/intrastate-spike ./cmd/intrastate` without the Makefile's `-ldflags`, so `version`/`commit`/`date` are Go's own module-fallback pseudo-version, not a real release stamp. That fallback behavior — and the fact the build silently succeeds without one — is what `make release-check` exists to catch (see Makefile comment above `release-check`); it does not change the envelope shape being captured here.

## SUCCESS envelope 2: top-level key set

```
/tmp/intrastate-spike version --as=json | jq -r 'keys|join(",")'
```

```
data,type
```

## FAILURE envelope 1: intrastate lint --model <nonexistent path> --as=json (raw stdout)

```
/tmp/intrastate-spike lint --model /tmp/does-not-exist-spike.toml --as=json
echo "exit=$?"
```

```
{"code":"model-unreadable","message":"cannot read /tmp/does-not-exist-spike.toml","param":"model","detail":"open /tmp/does-not-exist-spike.toml: no such file or directory"}
exit=2
```

## FAILURE envelope 1: pretty-printed

```
/tmp/intrastate-spike lint --model /tmp/does-not-exist-spike.toml --as=json | jq .
```

```
{
  "code": "model-unreadable",
  "message": "cannot read /tmp/does-not-exist-spike.toml",
  "param": "model",
  "detail": "open /tmp/does-not-exist-spike.toml: no such file or directory"
}
```

## FAILURE envelope 1: top-level key set

```
/tmp/intrastate-spike lint --model /tmp/does-not-exist-spike.toml --as=json | jq -r 'keys|join(",")'
```

```
code,detail,message,param
```

## FAILURE envelope 2: intrastate lint --model <malformed TOML> --as=json (raw stdout)

Fixture written for this spike only (`/tmp/spike-malformed.toml`, contents `this is not valid toml [[[`); not committed to the repo.

```
echo "this is not valid toml [[[" > /tmp/spike-malformed.toml
/tmp/intrastate-spike lint --model /tmp/spike-malformed.toml --as=json
echo "exit=$?"
```

```
{"code":"model-invalid","message":"model does not conform to the transition-model schema","param":"model","detail":"malformed_toml: toml: expected character =","findings":[{"code":"malformed_toml","message":"malformed_toml: toml: expected character =","locator":"/tmp/spike-malformed.toml:1"}]}
exit=2
```

## FAILURE envelope 2: pretty-printed

```
/tmp/intrastate-spike lint --model /tmp/spike-malformed.toml --as=json | jq .
```

```
{
  "code": "model-invalid",
  "message": "model does not conform to the transition-model schema",
  "param": "model",
  "detail": "malformed_toml: toml: expected character =",
  "findings": [
    {
      "code": "malformed_toml",
      "message": "malformed_toml: toml: expected character =",
      "locator": "/tmp/spike-malformed.toml:1"
    }
  ]
}
```

## FAILURE envelope 2: top-level key set

```
/tmp/intrastate-spike lint --model /tmp/spike-malformed.toml --as=json | jq -r 'keys|join(",")'
```

```
code,detail,findings,message,param
```

## FAILURE envelope 3 (additional data point, not requested): valid-TOML-but-schema-invalid model

Same shape as failure envelope 2 (schema-level `model-invalid`), captured only to confirm the "no `[model]` table" case does not surface a different envelope shape. Fixture written for this spike only (`/tmp/spike-blocking.toml`), not committed.

```
printf '[flow]\nid = "spike-blocking"\n' > /tmp/spike-blocking.toml
/tmp/intrastate-spike lint --model /tmp/spike-blocking.toml --as=json
echo "exit=$?"
```

```
{"code":"model-invalid","message":"model does not conform to the transition-model schema","param":"model","detail":"malformed_model_declaration: no [model] table","findings":[{"code":"malformed_model_declaration","message":"malformed_model_declaration: no [model] table","locator":"/tmp/spike-blocking.toml:1"}]}
exit=2
```

Key set is identical to failure envelope 2: `code,detail,findings,message,param`.

## Discriminator check: is a `type` key present on both success and failure records?

Quoting directly from the captured JSON above:

- SUCCESS envelope 1 (`lint`): `"type":"ok"` — present, value `ok`.
- SUCCESS envelope 2 (`version`): `"type":"ok"` — present, value `ok`.
- FAILURE envelope 1 (`lint`, unreadable path): top-level key set is `code,detail,message,param` — **no `type` key at all.** Not `"type":"failed"`, not any other value. The discriminator is simply absent from the failure shape.
- FAILURE envelope 2 (`lint`, malformed TOML): top-level key set is `code,detail,findings,message,param` — **no `type` key at all**, same absence.
- FAILURE envelope 3 (`lint`, no `[model]` table): identical key set to failure 2, same absence.

So the record is NOT `type: ok` / `type: failed` as a symmetric pair. It is: success envelopes wrap the payload as `{"type":"ok","data":{...}}`, and failure envelopes are a flat, undiscriminated error object `{"code":...,"message":...,"param":...,"detail":...[,"findings":[...]]}` with no `type` field of any kind. A caller distinguishing success from failure under `--as=json` cannot switch on a shared `type` key across both outcomes — it must either inspect exit code, or check for the presence/absence of the `type` key (or of `code`), since failure never sets `type` and success always does.

## grep-check: schema_version / format_version key present anywhere?

Quoting the captured top-level key sets:

- Success envelope 1 (`lint`): `data,type`
- Success envelope 2 (`version`): `data,type`
- Failure envelope 1: `code,detail,message,param`
- Failure envelope 2: `code,detail,findings,message,param`

None of the four contains `schema_version` or `format_version`. Answer: **no** — no such key exists in any captured envelope, at the top level or (for the `version` success envelope, the one place a version-shaped payload appears) inside `data`, where the keys are `version,commit,date` — i.e. intrastate's own build/release version, not an envelope/wire schema version.
