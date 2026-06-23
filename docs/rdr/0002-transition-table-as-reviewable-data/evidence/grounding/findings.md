# Grounding Findings

- **G-001 REFUTED** — A5 says `internal/cli/config/config.go::Load` "already uses stable parse/config error codes as the local pattern." Source confirms stable config load/read codes (`config-not-found`, `config-read-error`), but parsing is intentionally not wired and `config-invalid` is a TODO, not an emitted parse code. Evidence: `internal/cli/config/config.go::Load`.
