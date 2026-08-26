Model: claude-opus-5[1m]

# 3amigo iter-3 Persona 1 — Product Manager (RDR cli/0011)

Read: `0011:§problem-statement`, `0011:§approach`, `0011:§decision-rationale`, `0011:MVV` (owned set), then WIDENED to `0011:C1`, `0011:C2`, `0011:C3`, `0011:§consequences`, `0011:§failure-modes`, `0011:§briefly-rejected`, `0011:§load-bearing-decisions`, `0011:§mini-checks`, `0011:§existing-infrastructure-audit`, `0011:§phase-3-fixtures-and-docs`, `0011:§testing-strategy`, `0011:A5`, `0011:A6`, `0011:A11`. Grounded against `internal/cli/flow_next.go`, `internal/cli/flow.go`, `internal/cli/flow_state.go`, `internal/cli/flow_resolve.go`, `internal/cli/lint.go`, `internal/cli/root.go`, `README.md`, `docs/cli-output-contract.md`.

**What sent me outside the owned set.** The problem statement defines success as an outcome property ("no row is listed that the supplied state excludes, and every row listed that the supplied state could not decide names the undecided key"). Two halves of that promise are *silences*: what the caller sees when a row IS excluded, and where the user-facing documentation of the new default lives. Neither silence has a line range inside my owned set, so I widened to the contracts, Failure Modes, Briefly Rejected, and the phase plan, and then to the shipped payload struct and the verb list to check whether the recovery paths the RDR names actually exist.

---

## HIGH

### H1 — `0011:F1` names a diagnostic surface (`dump`) that does not exist, and the mode it recommends is contractually stripped of the fact the caller needs

`0011:F1` is the primary user-visible failure mode of the whole change: "a candidate the caller expected is absent." Its stated remedy is: *"Diagnose with `--all`: if it appears there, a match atom over a supplied key failed (compare the row's `match` block in `dump` against `owned`/`observed` in the payload)."* Both halves fail.

1. **There is no `dump` verb.** The shipped verb set is `version`, `lint`, `flow next`, `flow resolve`, `flow read-state`, `flow set-state` (`internal/cli/*.go` `Use:` fields). `0011:§background` and `0011:§existing-infrastructure-audit` name `internal/cli/flow.go`'s four-verb closure, and `0011:BR2` explicitly reaffirms "0005's four-verb closure stands." So the recovery instruction routes the user to a command that has never existed and that this RDR does not add.
2. **`--all` does not carry the missing fact.** The shipped `candidate` struct (`internal/cli/flow_next.go`) exposes `rule`, `outcome`, `required`, `unresolved`, `gates`, `next`, `writes`, `clear` — no `match` block, no expected value. `0011:C1` and `0011:C2` add `unknown` and `--all`; neither adds a match block. Worse, `0011:C2` states outright that under `--all` "a match atom is then neither an exclusion nor an entry in `unknown`, whatever its key's presence" and mandates the `--all` branch "MUST filter `BlockMatch` atoms out of `unknown`". So the RDR's own contract deliberately empties the diagnostic mode of match information.

Net user outcome: when the new default silently drops a row (present-and-unequal), the caller can observe the row's *existence* under `--all` and nothing about *which key failed or what value it wanted*. The only way to close the loop is to open the model TOML by hand — precisely the "re-deriving the transition table by hand" that `0011:§problem-statement` opens by declaring the defect.

**Blocks**: the Failure Modes / Recovery decision, and the disposition of `0011:BR3`. It also blocks locking `0011:C2`, whose scope is currently justified by a recovery path that does not work.

### H2 — `0011:BR3`'s rejection rests on a false equivalence, so the "why did my row vanish?" outcome was never actually decided

`0011:BR3` rejects reporting match-excluded rows with an `excluded_by` field on the grounds that it "grows the payload for the same information `--all` already yields." Per H1, `--all` does **not** yield that information: `0011:C2` strips match atoms from `unknown` under `--all`, and the payload has no match block in either mode. The rejected alternative is the only one on the list that would have delivered the excluded-row diagnosis; it was dismissed against a claim the record's own `0011:C2` contradicts.

This matters at the product level, not the design level: RDR 0005's own premortem quote carried in `0011:§key-discoveries` — "`next` omits enough condition detail to constrain a skill" — is the failure this RDR exists to close. This RDR narrows the *set* (a real improvement, A6's 21→3) but re-creates the same "not enough condition detail" gap one step further in: for the 18 rows that vanished, the caller is told nothing.

**Blocks**: the Alternatives disposition for `0011:BR3`, and the decision on whether `0011:C1` should require anything at all on an excluded row.

---

## MEDIUM

### M1 — `0011:§consequences` and `0011:§phase-3-fixtures-and-docs` contradict each other on whether the user-facing contract doc is edited

`0011:§consequences` states: "`docs/cli-output-contract.md` pins the envelope, not per-command payload fields, so this is a documented-surface change **without a contract-doc edit**; C3's help text and the test-file header are where a reader meets it."

`0011:§phase-3-fixtures-and-docs` states the opposite intent: "…**document `flow next`'s payload — including `unknown` and the `--all` invocation — in `docs/cli-output-contract.md`**, which currently shows only the invocation grammar."

The record therefore both asserts the doc needs no edit and plans the edit. Grounded: `docs/cli-output-contract.md:127` carries the `flow next` invocation under the comment "**Enumerate** the legal outcomes and their candidate rules" — a user-facing sentence describing exactly the 0005 semantics this RDR inverts. A caller reading the shipped doc after this lands is told the verb enumerates.

The consequence is not cosmetic: no contract (`0011:C1`/`0011:C2`/`0011:C3`) makes the doc update normative — `0011:C3` scopes documentation obligations to the cobra Short/Long help and the 0005 test-file header only — and no MVV step or `0011:S1`–`0011:S7` scenario checks it. Under the `0011:§consequences` reading, the doc edit is optional and will be skipped; under Phase 3 it is intended but unenforced. A doc that still says "enumerate" is a live misdirection for the user this RDR is for.

**Blocks**: the Consequences/documentation-surface decision, and it prevents any test asserting the documented surface matches the shipped default.

### M2 — `README.md:43` is a shipped user-facing description of the verb and no element in the record accounts for it

`README.md:43` reads `intrastate flow next --flow <name>       # list legal next outcomes`. The record's documentation census (`0011:C3`'s enumeration of shipped strings, `0011:§consequences`, `0011:§phase-3-fixtures-and-docs`) names `internal/cli/flow_next.go:70-71`, the 0005 test-file header, and `docs/cli-output-contract.md` — never `README.md`. The line is arguably still true (the alphabet is unchanged), but it is the first thing a new skill author reads about the verb and it says nothing about the state-conditioning that is now the verb's defining behaviour or about `--all`.

**Blocks**: closing the user-facing-surface census; it prevents an oracle over "every shipped description of `flow next` agrees with C1's default."

### M3 — `0011:§problem-statement`'s success criterion is stated but never turned into an assertion, so the MVV cannot fail on the user outcome

`0011:§problem-statement` defines success as: "no row is listed that the supplied state excludes, and every row listed that the supplied state could not decide names the undecided key" — explicitly "not a target cardinality." But `0011:MVV` step 3 asserts precisely a cardinality-and-identity oracle over one hand-picked model (`candidates[]` is exactly `prelock`, `resolve-route-back`, `resolve-abandon`), and `0011:§mini-checks`' `oracle` table books that as "21→3 narrowing: a stripped-match build reports 21, so the count discriminates."

The gap: the stated criterion is universally quantified over rows and states; every MVV/`0011:S1`–`0011:S7` oracle is existential over fixtures. Nothing in the record checks the *general* property. That is defensible as a testing choice, but the record should not state a criterion in a form its validation cannot reach — a reader of `0011:MVV`'s "End-state" line will believe the outcome was validated as stated.

**Blocks**: the Scope Verification gate (`0011:G-scope`, still template placeholder text), which must state the specific test or proof for the MVV. As written the MVV proves the fixtures, not the criterion.

---

## LOW

### L1 — `0011:§problem-statement` concedes the 0010 decision-table class gets no narrowing at all, but the record never states what the user does about it

`0011:§problem-statement` is admirably direct: "where it is observed and unsupplied (0010's decision-table class with no `--tag`, A5) every row remains a candidate and each names the undecided key under `unknown`. That second case is the old symptom's shape with a diagnosis attached, and it is the intended behaviour." `0011:A5`'s spike confirms it live — 4 of 4 rows reported, no narrowing.

The record calls the remedy "named in the payload (bind the reader, supply the tag)". But `0011:C1`'s `unknown` entry is `{key, reason}` — a key name and the token `absent`. It does not tell the caller which of the two remedies applies (is this key observable via `--tag`, or does it want a reader binding?), nor which values the key can take. For a decision-table model the caller is handed N key names and left to work out the rest from the model file.

This is a scoping judgement the RDR is entitled to make, and it makes it explicitly, so this is LOW rather than higher. But `0011:§consequences`' positive claim — "the payload names which [remedy], so the caller is told what to do rather than only that something is unknown" — overstates what `{key, absent}` conveys for the class where narrowing does not happen at all. The differentiation the bullet claims exists between `absent` and `uncomparable`, not between "bind a reader" and "supply a tag", which is the fork the 0010 caller actually faces.

**Blocks**: nothing hard; it slightly overstates the Consequences bullet and would mislead a reader sizing the follow-on work for the 0010 class.

### L2 — `0011:D-naming`'s `--all` precedent survey does not cover the one asymmetry that matters to this user

`0011:D-naming` grounds `--all` well (`docker ps -a`, `git branch -a`, `gh workflow list --all`, `ps -A`). Every cited precedent shares a property this one does not: in `docker ps` and `git branch`, the *narrow* default and the *wide* `--all` return the same record shape, and the widened rows are self-explanatory (a stopped container still shows its status). Here `--all` returns a strictly *less* informative record for the rows it adds back (`0011:C2` strips their match atoms from `unknown`), which is the inverse of the precedent's ergonomics and is what makes H1 bite. The naming is fine; the precedent's implied ergonomic contract is not honoured and the record does not note the divergence.

**Blocks**: nothing on its own; it is the naming-decision-level restatement of H1 and would be resolved by resolving H1.
