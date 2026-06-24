# Persona 2 — Implementer

1. In the first hour I would ask what concrete request structs and response
   structs each verb must expose to `respond.OK`. The RDR says text and JSON are
   derived from the same verb-specific result, but does not name the minimum
   fields for `next`, `resolve`, `read-state`, or `set-state`.

2. I would ask how to parse tags on the CLI. The draft uses `--tag stage=resolve`
   examples and mentions a structured input file, but does not state whether
   repeated `--tag key=value` is the required MVP form, how duplicate tags are
   refused, or whether structured input is in the first implementation slice.

3. I would ask which exact `CLIError.Code` values to add. The RDR lists
   categories such as unknown outcome, zero matching row, multiple matching
   rows, missing facts, accessor unavailable, gate indeterminate, and read-back
   mismatch, but leaves the code spelling to implementation taste.

4. I would ask what `flow next` returns for conditional rows. "Conditional
   summaries" could mean raw guard predicates, human text, row ids, missing fact
   names, or next tag-set candidates. That affects both model interfaces and
   tests.

5. I would ask what `flow set-state` takes as input. The examples pass only
   `--tag status=Final`, but the prose says it persists a decided owned-tag
   mutation through write accessors. It is unclear whether `--tag` here means
   expected target owned tags, observed context, or the already resolved next
   tag-set.
