Model: claude-opus-5

# In-repo prior art — how closure/growth is declared today (RDR 0029, Stage 2)

Gathered in the main authoring context (source reads, not a corpus search).
Every claim below is quoted from source at HEAD, not paraphrased.

## Finding 1 — the project already decides this per-vocabulary, in code comments

Four machine-readable vocabularies each declare their own growth discipline,
and the declarations are NOT consistent with each other.

`internal/graphlint/taxonomy.go:37` — the advisory tier, closed against growth:

> // The four advisory finding codes. The tier is CLOSED at these four
> // (`0006:C17`).

`internal/graphlint/taxonomy.go:56,68` — the `reason` set, closed but growing:

> // The closed, append-only `reason` set `graph-unprovable-coverage`
> // carries (Technical Design, reason table).
> ...
> // atoms. The set is closed and APPEND-ONLY, and this is the append.

`internal/table/category.go:89` and `:52` — the load categories, enumerable
and append-only, with the size explicitly disclaimed:

> // Categories returns the closed load-category set in declaration order.
> // A constant is not in the closed set until it is appended to
> // Categories() below

`internal/table/category.go` (0028 registration block) — the consumer rule:

> // The list's size is not a contract, so a consumer asserts RELATIVE
> // order and never a tail position or a count (`0028:C1.4` registration:).

`internal/cli/clierr/clierr.go:49` — the envelope's own field discipline:

> // Extend with new optional fields as needed — keep them
> // `omitempty` so the envelope stays append-only and stable for tools.

⇒ Forces here: the RDR is NOT inventing a policy on a blank slate. It is
**hoisting and reconciling four existing, divergent, code-comment-local
policies** into one stated promise. Two senses of the word "closed" are live
on adjacent surfaces — "will never grow" (graphlint advisory) vs "fully
enumerable but append-only" (table categories) — which is a live ambiguity a
consumer reading either comment would get wrong.

## Finding 2 — the load-category set has already grown four times

`internal/table/category.go` records each growth event with its RDR:

- RDR 0002's original 25 (`0002:C24`, "the load category floor")
- RDR 0024's 3 emit-vocabulary categories (`0024:D-naming`)
- RDR 0025's 6 command-carrier categories (`0025:C5`)
- RDR 0028's 6 `edit`-carrier categories (`0028:C1.4`)

⇒ Forces here: "a new code appears" is not hypothetical. It is this project's
**normal release event**, having happened four times pre-1.0. A policy that
classifies it as breaking would make nearly every release a major bump.

## Finding 3 — the severity ladder that makes non-breaking introduction possible ALREADY EXISTS

`internal/graphlint/taxonomy.go:48-53`:

> // The two severities. There is no third tier (`0006:C8`, `0006:C17`).
> const (
> 	SeverityBlocking = "blocking"
> 	SeverityInfo     = "info"
> )

`internal/graphlint/engine.go:32,35` partitions on it, and only the blocking
half reaches the refusal path — `internal/cli/lint.go:213-228`:

> 	blocking := report.Blocking()
> 	if len(blocking) > 0 {
> 		... Findings: blocking,   // refusal, non-zero exit
> 	}
> 	advisory := report.Advisory()
> 	... return respond.OK(cmd, respond.Success{Data: lintPayload{Findings: advisory}})

⇒ Forces here: a new finding code introduced at `info` severity rides the
SUCCESS envelope and cannot change any consumer's exit code or clean/dirty
verdict. The mechanism for "introduce a new check without breaking a
pipeline" is already built and shipping. The decision this RDR owes is the
*policy* governing promotion `info` -> `blocking`, not new machinery.
This is structurally the clippy allow-by-default -> warn-by-default ladder.

CAVEAT (demote to Resolve): the graphlint advisory tier is declared CLOSED at
four, so a new graphlint code cannot today be introduced at `info` without
amending `0006:C17`. Whether that closure is load-bearing or incidental is an
order/reachability-adjacent claim about an existing contract — Method: Peer
RDR, against `0006:C17`.

## Finding 4 — the authoritative wire-contract doc says nothing about versions

`grep -rn -i "release|upgrade|breaking" docs/cli-output-contract.md` returns
ZERO matches. The document `.rdr/resources.md` names as "authoritative for
verb I/O, error envelope, exit-code mapping" carries no statement about what
survives a version bump.

Meanwhile `llms.txt` actively directs agents to depend on the binary:

> The binary is authoritative and self-describing. Prefer asking it over
> reading any file here — a file is a snapshot, the binary is the build
> you actually have

⇒ Forces here: the promise has no home today, and the project is
simultaneously telling agents to bind tightly to the binary's output. The
RDR must name WHERE the promise lives, and `llms.txt`/`cli-output-contract.md`
are the two candidate homes already in the agent's read path.

## Finding 5 — `.goreleaser.yaml` already treats a surface as a compat contract

`.goreleaser.yaml:55`:

> # it is a compatibility surface, not cosmetics. Changing it breaks any

⇒ Forces here: precedent exists in-repo for naming a surface a compatibility
contract in the place that produces it. Consistent with putting the machine-
output promise beside the machine output rather than in a release doc.

## Finding 6 — exit-code classes are a coarse, stable surface

`internal/cli/clierr/clierr.go::ExitCodeFor` maps six groups onto four codes
(0, 2, 3, 130); `GroupUserEnv` and `GroupInternal` both map to 2.

⇒ Forces here: the exit surface is deliberately coarser than the group
vocabulary, so groups can be added without moving an exit code. That is a
third existing example of the same append-without-breaking shape, and it
suggests the promise should be stated per-surface with different strengths,
not as one blanket "JSON is stable" claim.
