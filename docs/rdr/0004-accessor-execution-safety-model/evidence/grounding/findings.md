# Grounding Findings

Findings:

- A3 `Method: Source Search` evidence cited file/line ranges instead of the
  required greppable `path::Symbol` form. The underlying source claims are
  confirmed, but the evidence shape would fail the RDR method rule for source
  searches.

Confirmed codebase claims:

- `internal/cli/clierr::CLIError`
- `internal/cli/clierr::ExitCodeFor`
- `internal/cli/respond::Fail`
- `internal/cli/respond::ValidateMode`
- `internal/cli/config::Load`
- `internal/cli/config::ConfigFileName`

Confirmed peer-RDR claims:

- RDR 0001 keeps resolver execution stateless and places accessor execution
  outside the resolver boundary.
- RDR 0002 owns table-carried accessor references and does not execute
  accessors.
- RDR 0003 consumes accessor-produced tag values after binding and does not
  execute accessors during guard evaluation.

Confirmed spike claims:

- `docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/main.go`
  defines declared read/gate/write capabilities, timeout execution,
  capability-mismatch refusal, gate-indeterminate refusal, read-back mismatch,
  and replay fixture behavior.
- `docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/output.txt`
  records the read, gate, timeout, execution failure, capability mismatch,
  write success, read-back mismatch, and replay dispositions.

Inverse search:

- Searched `internal/` and `cmd/` for accessor, capability, executor, refusal,
  read-back, gate-indeterminate, and shell/exec surfaces. No existing accessor
  executor, capability discriminator, or accessor refusal taxonomy exists in
  production code.
