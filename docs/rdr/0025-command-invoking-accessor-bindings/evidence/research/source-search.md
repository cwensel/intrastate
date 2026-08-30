Model: claude-fable-5

## A2 — os/exec passes argv verbatim (no shell, no glob)
Queries: grep exec.go for "does not invoke the system shell|Args holds|CmdLine|LookPath"; sed exec.go 5-20, 149-160, 385-405.
Accepted:
- $GOROOT/src/os/exec/exec.go:9-12 (package doc): "Unlike the "system" library call from C and other languages, the os/exec package intentionally does not invoke the system shell and does not expand any glob patterns or handle other expansions, pipelines, or redirections typically done by shells. The package behaves more like C's "exec" family of functions."
- exec.go:156-157 (Cmd.Args): "Args holds command line arguments, including the command as Args[0]. If the Args field is empty or nil, Run uses {Path}."
- exec.go:149-153 (Cmd.Path): "Path is the path of the command to run. This is the only field that must be set to a non-zero value. If Path is relative, it is evaluated relative to Dir."
- exec.go:387-389 (Command): "If name contains no path separators, Command uses [LookPath] to resolve name to a complete path if possible. Otherwise it uses name directly as Path."
- exec.go:396-403 (Windows caveat): "On Windows, processes receive the whole command line as a single string and do their own parsing. Command combines and quotes Args into a command line string with an algorithm compatible with applications using CommandLineToArgvW ... Notable exceptions are msiexec.exe and cmd.exe ... you can do the quoting yourself and provide the full command line in SysProcAttr.CmdLine, leaving Args empty."
Verdict: VERIFIED (Unix: argv verbatim; Windows: re-quoted per CommandLineToArgvW, overridable via SysProcAttr.CmdLine).

## A5 — Kubernetes ExecAction.Command is argv, no shell
Queries: grep -A12 "type ExecAction struct" k8s.io/api@v0.36.2/core/v1/types.go.
Accepted: /Users/cwensel/go/pkg/mod/k8s.io/api@v0.36.2/core/v1/types.go:2703-2711 — `type ExecAction struct` at 2703; Command doc at 2704-2707: "Command is the command line to execute inside the container, the working directory for the command  is root ('/') in the container's filesystem. The command is simply exec'd, it is not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use a shell, you need to explicitly call out to that shell."; field at 2711: `Command []string`.
Verdict: VERIFIED.

## A6 — 0004's code-set closure scopes 0004's own defects, not the seam
Queries: rdr inspect 0004 (element list); grep projector output for closed|eight|exhaust; rdr inspect --select 0004:A6/C16/§testing-strategy; rg ValidationCode internal/accessor/model.go + *_test.go; rg CatMalformedAccessorDeclaration internal/table.
Peer read (a): no 0004 contract element states the validation-code set is closed. The only "eight" wording is in 0004:§testing-strategy scenario 1 (projected line 940): "... a missing or empty read requested-key set also fail before resolution — eight validation arms, each asserting its own named code." The "closed" hits in 0004 (lines 819/825) refer to the kernel refusal kinds / refusal-class list, not to validation codes. Closure is therefore an own-code artifact scoped to the eight shapes 0004 defines; 0004:C1-C16 name no seam-wide closed enum.
Own code (b): internal/accessor/model.go:91-129 — `type ValidationCode string`; eight `Code*` constants (CodeMissingAccessor ... CodeMissingRequestedKeySet); `validationCodes` slice; `// ValidationCodes returns the closed eight-member validation-code set.` Breakage on a ninth member: internal/accessor/validation_0004_test.go:170 `TestReq84_TheValidationCodeSetIsClosedAtEight` compares `ValidationCodes()` to an exact eight-element list (lines 171-189) — adding a code to `validationCodes` fails it (a REQ-84 test, so amending it is a 0004-scope change).
Loader defect family today: internal/table/category.go:11 `type Category string`; accessor-entry defects use `CatMalformedAccessorDeclaration Category = "malformed_accessor_declaration"` (category.go:26), raised via `fail(CatMalformedAccessorDeclaration, ...)` in internal/table/load.go::accessorTable (load.go:942-995, e.g. "path is absent or empty" at 951-952); `CatMalformedAccessorBinding` also exists (load.go:975).
Verdict: VERIFIED — adding codes to accessor.ValidationCode does not contradict any 0004 contract, but it does collide with own-code closure (comment + TestReq84). Recommended C5 home: the table loader's `Category` family (a new `Cat*` code beside CatMalformedAccessorDeclaration, or that category itself with a detail string), since command/path are loader-side entry-shape defects, not accessor-execution defects.

## A4 — read-back reader is selected by role via readerFor; 0016:C4 owns it
Queries: rg readerFor|readBack|ReadBack internal/accessor non-test; sed model.go 225-240; rdr inspect 0016; rdr inspect --select 0016:C4 --select 0016:A5; rdr status --tags 0016.
Own code: internal/accessor/model.go:231 `func (reg Registry) readerFor(role string) (Definition, bool)` — predicate at 233: `if d.Identity.Capability == CapRead && d.Accessor.Role == role` (by role, first match); called from internal/accessor/executor.go:275 `reader, hasReader := e.Registry.readerFor(def.Accessor.Role)`. It returns a `Definition` and never inspects the binding's implementation — agnostic to how the reader is bound.
Peer: 0016:C4 (lines 322-335): "`0004:C12` and `0004:C13` stand as written; this RDR overrides nothing. Under C1, `internal/accessor/model.go::readerFor` resolves "the role's read definition" uniquely — and the resolution site MUST fail closed ..." 0016:A5: "`readerFor` is the sole role→reader resolution point" (Pending). rdr status --tags 0016: status=Draft, gate_stale=false, gate_written=false, ca=all-pending.
Verdict: VERIFIED for code + normative home; NOTE 0016 is still Draft (not Final), so the landing-order precondition on 0016 is not yet satisfied. Also note current readerFor is first-match, the fail-closed form 0016:C4 requires is not yet implemented.

## A9 — consumers of accessor `path`
Queries: rg '\.Path\b|"path"|\bpath\b' over internal/table/{load,normalize,dump,model,source}.go, internal/cli/*.go, internal/cli/flowbind/*.go; rg '\.Readers|\.Writers|\.Gates' internal non-test; rg '\.Path\b' internal excluding load.go/registry.go.
Sites reading/validating table.Accessor.Path (non-test):
1. internal/table/source.go:86 `Path *string toml:"path"` (source decl).
2. internal/table/load.go::accessorTable 951-952 (presence/empty validation) and 988 (copy into Accessor).
3. internal/cli/flowbind/registry.go:33,43,53 (constructors copy acc.Path into Reader/Writer/Gate bindings).
4. internal/cli/flowbind/flowbind.go:192-193, 264, 319-327 — the flowbind bindings' OWN Path field (copied at 3) is consulted via unreachable()/verdictFor() at run time.
Other accessor-map consumers (load.go:1012,1036; normalize.go:449; cli/flow_exec.go:154,215; cli/flow_state.go:448,472-473) read Role/Keys only, never Path. dump.go does not emit accessor path.
Verdict: PASS with one extra to enumerate — flowbind.go's runtime use of the binding's Path (downstream of the constructors, so a command-backed binding simply would not carry it, but the RDR should name it).

## Reuse audit
Query: rg -n 'os/exec|exec\.Command' --type go -g '!*_test.go' . — no hits (exit 1). No existing command-execution capability in own non-test code.
