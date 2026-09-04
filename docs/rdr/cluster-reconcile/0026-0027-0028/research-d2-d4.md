# Research grounding for JDR 0003 §D2–§D4 (2026-09-03)

Model: claude-sonnet (three delegated corpus passes); judged in the session model.
Corpora searched: DevRef, PapersFast, SpecDrivenDev, RetrofitRDR (arc text + semantic); Go stdlib docs; local man pages and module cache. No corpus named "langref" exists; the language-reference angle is the Go stdlib documentation.

## §D2 — invocation-bound argv words → (a) confirmed
- `internal/cli/cmdbind/cmdbind.go:415` — existing refusal `!filepath.IsAbs(path) || strings.HasPrefix(path, "-")` for `{artifact}`; `resolveArgv0` (434-453) never lets an unresolved bound path stand at argv0. (a) mirrors both.
- `go doc os/exec`, `os/exec.Command` — no shell, args passed verbatim: the class is CWE-88 argument injection.
- `go doc flag` — `--` terminator; bare `-` is a non-flag argument. The `--` alternative admits `-1`/`-` but only where every child honours it; the repo already chose refusal for `{artifact}`.
- `git rm --help` SYNOPSIS `[--] <pathspec>` — the canonical insert-`--` precedent; not adopted here for the same reason.
- DevRef/PapersFast/SpecDrivenDev: no on-point CWE-88 text; nothing contradicting.
- Known cost: legitimate `-1` / `-` values refuse; accepted already for `{artifact}`.

## §D3 — exit group of stale-model refusals → (b) confirmed
- `/usr/include/sysexits.h` — `EX_TEMPFAIL` (75, retry later) vs `EX_DATAERR` (65): a stale anchor is data-shaped.
- DevRef: RESTful Web APIs p.281, REST API Design Rulebook p.55 — HTTP 412 Precondition Failed: client re-reads and resubmits, never blind retry.
- DevRef: service_design_patterns p.245 — classify retryable before the retry loop (typed discriminator).
- RetrofitRDR `process/rdr/cli/0028-migrate-subcommand.md:258-266` + `coverage.md:163-169` — `translateDDLError` maps invocation/state refusals to the exit-2 group.
- RetrofitRDR `process/rdr/cli/0120-…:329-380` — stale anchor is a loud refusal, never silent-skip.
- JDR 0001 §D10 rule 7 and `docs/cli-output-contract.md:80-100` — the in-repo precedent and the "never exit 3 for a request refusal" rule.

## §D4 — applied sense at the envelope → (b) confirmed
- gRPC `codes/codes.go:41-67, 105-121, 170-180` — "may have completed" is a property of `Unknown`/`DeadlineExceeded`; `Aborted`/`Unavailable`/`FailedPrecondition` chosen by the client's remedy: one code per remedy.
- DevRef: Designing Data-Intensive Applications p.380 — 2PC in-doubt state resolves by querying status before retry (JD-1's stake).
- DevRef: build-apis-you-wont-hate p.51 — string-matching messages for control flow is the anti-pattern; errors tie to codes.
- JDR 0001 §D10 rule 7 — read-back-incomplete/timeout already distinct codes carrying the applied text; (b) extends the pattern.
- No source favours a boolean/enum field for this discriminator specifically.
