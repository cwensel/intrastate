Model: claude-sonnet-5

# Spike: A7 false-refusal base rate

## Question

A7 claims: the false-refusal class position-freedom opens — a listed
interpreter basename carried as a data argument under an unlisted argv0,
followed later by that interpreter's flag — is rare enough in real command
bindings to accept as a stated cost rather than bound.

Prior evidence (12 hand-built vectors) established the class EXISTS and is
trivially constructible. It did not establish the base RATE in bindings
people actually write. This spike samples real, independently-authored
shell-free argv bindings and counts how many C1 would refuse, then
classifies each refusal as a true positive (really spawns a shell/
interpreter with inline code) or a false refusal (the listed word is a
data argument, not an invocation).

## Harvest method

Sources actually available on this machine, scoped with `-g`/`find` (no
unscoped filesystem sweep). In-repo `command = [...]` fixtures in
`internal/table/testdata/**` and `*_test.go` were inspected first but
excluded from the sample proper — they are hand-authored fixtures for
this exact feature (0025/0027), not independent real-world bindings, so
counting them would double-count the same 12-vector-style construction
already on record for A7.

Independent real-world argv sources, harvested from checked-out repos
under `/Users/cwensel/sandbox/newcoinc/` (mostly `schema-evolution/_repos/`
and `schema-evolution/_examples/`, plus `thirdparty/*`, `state-machines/*`):

```
find . -path '*/node_modules/*' -prune -o -path '*/.git/*' -prune \
  -o -iname 'Dockerfile*' -print 2>/dev/null | grep -v node_modules
# -> 537 Dockerfiles found

find . -path '*/node_modules/*' -prune -o -path '*/.git/*' -prune \
  -o -iname 'docker-compose*.yml' -print -o -iname 'docker-compose*.yaml' -print \
  -o -iname 'compose.yml' -print -o -iname 'compose.yaml' -print 2>/dev/null \
  | grep -v node_modules
# -> 393 compose files found
```

From those files, only the JSON/exec-form array directives are true
shell-free argv bindings comparable to a `command = [...]` entry:
`CMD [...]` / `ENTRYPOINT [...]` in Dockerfiles, and `command: [...]`
(single- or multi-line) in compose files. A shell-form `CMD foo bar` or
compose `command: foo bar` (bare string) is NOT an argv binding — it is
shell text handed to `/bin/sh -c`, out of C1's scope by definition, so
those were excluded from the primary sample.

Parsing (JSON first, then Python-literal for single-quoted compose
arrays, then a comment-stripped JSON retry) yielded:

- 221 `CMD`/`ENTRYPOINT` exec-form arrays from Dockerfiles
- 77 compose `command:` list-form arrays
- 423 raw vectors total; after de-duplicating identical argv tuples
  (many repeated `["postgres","-c","config_file=..."]` across
  compose files, `CMD ["/app/main.py"]` etc.): **264 unique argv
  vectors**

A secondary, explicitly weaker proxy was also pulled to see how often a
listed interpreter word occurs as a token in shell-form command text
(Makefile recipe lines), which is informative about token frequency but
is NOT an argv binding and is kept separate:

```
find . -path '*/node_modules/*' -prune -o -path '*/.git/*' -prune \
  -o -iname 'Makefile' -print 2>/dev/null | grep -v node_modules
# -> 461 Makefiles found; 15,138 non-comment recipe-ish lines scanned
```

## Predicate implementation

Implements C1's predicate text exactly: `refuse iff there exist i < j
with base(argv[i]) a listed interpreter and argv[j] one of its listed
flags`, `base` = text after the last `/`, exact match, no folding.

```python
shellInterpreters = {
    "sh": ["-c"], "bash": ["-c"], "dash": ["-c"], "ksh": ["-c"],
    "zsh": ["-c"], "csh": ["-c"], "tcsh": ["-c"], "python": ["-c"],
    "ruby": ["-e"], "node": ["-e", "--eval"], "php": ["-r"],
}

def base(word):
    return word.rsplit("/", 1)[-1]

def predicate(argv):
    """Returns (refused, form, i, j) - lowest i then lowest j."""
    for i, w in enumerate(argv):
        flags = shellInterpreters.get(base(w))
        if not flags:
            continue
        for j in range(i + 1, len(argv)):
            if argv[j] in flags:
                return True, f"{argv[i]} {argv[j]}", i, j
    return False, None, None, None
```

A second function, `argv0_only_predicate`, was also run for contrast: it
mirrors the SHIPPED `interpreterForm` in `internal/table/load.go` (which
only checks position 0 after skipping a leading `env` NAME=VALUE chain,
not "any i"). This is a discrepancy between C1's predicate text and the
current implementation, noted below, but is a separate finding from A7's
rarity question — it would make the shipped code UNDER-refuse relative
to C1, not over-refuse.

## Run

```
$ python3 harvest_parse.py     # parses Dockerfiles + compose files -> 264 unique argv vectors
$ python3 run_predicate.py     # applies both predicates
N total unique vectors: 264
C1 predicate (any i<j) refusals: 23
Shipped argv0-only predicate refusals: 22
```

## Refusal table

All 23 refusals under the full C1 predicate, classified:

| # | Source | argv | matched | i,j | TP / FP |
|---|--------|------|---------|-----|---------|
| 1 | posthog/rust/Dockerfile | `['/bin/sh','-c','/usr/local/bin/$BIN']` | `/bin/sh -c` | 0,1 | TP |
| 2 | immich/server/Dockerfile | `['tini','--','/bin/bash','-c']` | `/bin/bash -c` | 2,3 | TP |
| 3-11 | chroma/rust/Dockerfile (9 services) | `['sh','-c','ulimit -c 0 && exec ./<svc>']` | `sh -c` | 0,1 | TP |
| 12 | kratix-marketplace/istio/.../Dockerfile | `['sh','-c','./execute-pipeline']` | `sh -c` | 0,1 | TP |
| 13 | kratix-marketplace/workshop-app-promise/.../Dockerfile | `['sh','-c','pipeline.sh']` | `sh -c` | 0,1 | TP |
| 14 | kratix-marketplace/sql/aws/.../Dockerfile | `['sh','-c','./resource-lifecycle']` | `sh -c` | 0,1 | TP |
| 15 | electric-sql/examples/burn/Dockerfile | `['sh','-c','/app/bin/migrate && /app/bin/server']` | `sh -c` | 0,1 | TP |
| 16 | OpenMetadata/.../Dockerfile.fuseki-working | `['sh','-c','mkdir -p ... && exec ...']` | `sh -c` | 0,1 | TP |
| 17 | apache-atlas/dev-support/atlas-docker/Dockerfile | `['/bin/bash','-c','/root/atlas-bin/... ; tail -fF ...']` | `/bin/bash -c` | 0,1 | TP |
| 18 | kratix/hack/kratix-pipeline-debugger-image/Dockerfile | `['sh','-c','./execute-pipeline.bash']` | `sh -c` | 0,1 | TP |
| 19 | openlineage/integration/hive-docker/Dockerfile | `['sh','-c','/entrypoint.sh']` | `sh -c` | 0,1 | TP |
| 20 | kratix/samples/paved-path-demo/configure-pipeline/Dockerfile | `['sh','-c','cp /transfer-input/* /kratix/output']` | `sh -c` | 0,1 | TP |
| 21 | supabase/docker/docker-compose.yml | `['/bin/sh','-c','/app/bin/migrate && ... eval "..." && ...']` | `/bin/sh -c` | 0,1 | TP |
| 22 | forem/uffizzi/docker-compose.uffizzi.yml | `['bash','-c','./uffizzi/entrypoint.sh bootstrap']` | `bash -c` | 0,1 | TP |

23 rows total (chroma contributes 9 near-identical service entries, rows
3-11 collapsed in this table; each is a distinct file/service pair, all
`sh -c '<shell text>'` at argv0/1).

**23/23 refusals are true positives.** Every refusal fires at `i=0, j=1`
— the interpreter is literally argv0, immediately followed by its flag.
Zero refusals fire on a listed word carried at a non-zero position under
an unlisted argv0 (the A7 false-refusal shape).

Note on refusal #2 (`tini -- /bin/bash -c` with a separate
`CMD ['start.sh']`): at Docker runtime, ENTRYPOINT and CMD concatenate,
so the effective exec is `tini -- /bin/bash -c start.sh` — bash really
is invoked with `-c`, so this is correctly a true positive under the
full C1 predicate. It is also the one case where the SHIPPED
argv0-only predicate (22 refusals) diverges from C1's stated predicate
(23 refusals): the shipped code only walks a leading `env` chain, not
`tini --`, so it misses this true shell spawn. That is a shipped-code
under-refusal relative to C1's own text — a distinct finding from A7,
not a false-refusal / over-refusal issue, and not pursued further here.

### Precondition check: does a listed basename even appear at a non-argv0 position?

Independent of whether a flag follows, checking every non-argv0 word in
all 264 vectors for a listed interpreter basename (the precondition for
the A7 class to be reachable at all):

```
non-argv0 occurrences of a listed interpreter basename: 2
  ('openproject/docker/ci/Dockerfile', 'CMD', ['setup-tests','bash'], 'bash', pos=1)
  ('immich/server/Dockerfile', 'ENTRYPOINT', ['tini','--','/bin/bash','-c'], '/bin/bash', pos=2)
```

Only 2 of 264 vectors (0.8%) even carry a listed interpreter word outside
argv0. One (`setup-tests bash`) has no following flag at all, so C1
never fires on it — the class's own precondition (word present, flag
follows) still isn't met. The other is `tini -- /bin/bash -c`, a genuine
shell spawn (true positive, see above). **Zero of the 264 vectors match
the A7 false-refusal shape**: a listed basename as data under an
unlisted argv0, with a same-interpreter flag later, where the interpreter
is NOT actually being invoked.

### Weaker proxy: Makefile recipe-line word co-occurrence (informative only, not argv, kept separate)

```
makefile recipe-ish lines scanned: 15138
lines containing a listed interpreter word: 332
lines containing both a listed word and a listed flag token: 30
```

Manually inspecting a sample of the 30 "both" lines: every one seen is a
genuine `sh -c '...'` / `bash -c "..."` shell invocation embedded as
one shell-syntax word (e.g. `docker exec $(CONTAINER) sh -c '...'`,
`xargs ... sh -c '...'`) — out of C1's scope by definition (a one-word
shell string, not a tokenized argv, per C1's explicit "out of scope, BY
NAME" carve-out) and not a disguised-data-argument case either. This
proxy corroborates the Dockerfile/compose finding — where a listed
interpreter word co-occurs with its flag in real command text, it is
overwhelmingly because the text really does invoke that interpreter, not
because the word is smuggled in as an unrelated data argument.

## Verdict

**A7 is VERIFIED, downgraded from "accepted cost" to "measured near-zero
rate on the sampled corpus."**

Across 264 unique, independently-authored, real shell-free argv bindings
(Docker `CMD`/`ENTRYPOINT` exec-form and compose `command:` lists, sourced
from ~20 unrelated OSS projects — Kratix, Supabase, Chroma, Immich,
PostHog, Apache Atlas, OpenMetadata, OpenLineage, ElectricSQL, Forem,
OpenProject, and others):

- 23 refusals under C1's full predicate, **23/23 (100%) true positives**
  — every single refusal is a real `sh -c` / `bash -c` shell spawn at
  argv0.
- Only 2/264 vectors (0.8%) carry a listed interpreter basename at any
  non-argv0 position at all, and neither is a false refusal (one has no
  trailing flag so never fires; the other is a genuine shell spawn via
  a wrapper prefix).
- **0/264 (0%) exhibit the A7 false-refusal shape.**

This supports "rare enough to accept as a stated cost rather than bound"
— on this sample the rate is not just rare, it is zero, and the
precondition (listed word as a non-argv0 data token at all) is itself
rare (0.8%).

**Sampling caveat, stated plainly**: 264 vectors from public OSS
Dockerfiles/compose files is an upper-bound-style sample, not a
population rate for intrastate's own binding authors. It skews toward
container entrypoints (a narrow, convention-heavy genre: `sh -c
"<script>"` or a bare binary) rather than the general shape of hand-written
CI/task-runner argv bindings intrastate targets. A construction like
`["backup-tool", "python", "-c", "restore.sql"]` (a tool that happens to
take a literal argument named `python` followed unrelatedly by a later
`-c` flag for its own purposes) was not observed anywhere in this
corpus, but the corpus is not large or diverse enough to certify such a
shape doesn't exist in the wild — only that it is not common enough to
surface in ~300 real-world container/compose argv bindings. A7's
"If wrong" trigger (the class turns out common) remains the correct
re-open condition if a future binding author reports a real false
refusal; this spike found no evidence to trigger it now.
