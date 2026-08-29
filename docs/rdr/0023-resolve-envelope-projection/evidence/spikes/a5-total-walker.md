Model: claude-opus-5[1m]

# A5 spike: total cobra command-tree walker

Date: 2026-08-29

## Claim under test

A5's implementability half: a TOTAL command-tree walker — no name-based
skip, no `Hidden` gate — is writable against the real root, enumerates
the auto-generated `help` AND `completion` commands, and can assert that
the set of commands registering the flag `plan-only` is exactly
`{flow resolve}`. Today that set is empty, since the flag is
unimplemented; the spike's job is to show the walker reaches the scope
C1 needs, not to find the flag.

## Cobra facts

Cobra version: `github.com/spf13/cobra v1.10.2` (`go.mod`).

Both auto-generated subcommands have EXPORTED initializers, so a test can
materialize them without going through `Execute`:

- `github.com/spf13/cobra@v1.10.2::(*Command).InitDefaultCompletionCmd`
  (`completions.go`) — `ExecuteC` merely calls it; a test can call it
  directly. Its removal branch fires only when the receiver has no
  subcommands. `NewRootCmd` adds `version`, `lint`, `flow`, `docs`
  (`internal/cli/root.go`), so `completion` persists.
- `github.com/spf13/cobra@v1.10.2::(*Command).InitDefaultHelpCmd` —
  likewise exported, and already called by the repo at
  `internal/cli/help_all.go::registerHelpAllOnTree`, which force-inits
  `help` to wire `--all` onto it.

Consequence: a bare `NewRootCmd()` carries `help` (force-inited by
`help_all.go`) but has no `completion` until `InitDefaultCompletionCmd`
runs. The spike's bare-root probe confirms this: children are
`[docs flow help lint version]`.

## Spike program

Throwaway test in package `cli` (needed for `NewRootCmd` and for the
contrast call into unexported `walkCommandTree`); deleted after the run.

```go
// totalWalk recurses into every child from Commands() with NO name-based
// skip and NO Hidden gate.
func totalWalk(root *cobra.Command, fn func(*cobra.Command)) {
	fn(root)
	for _, sub := range root.Commands() {
		totalWalk(sub, fn)
	}
}

func TestA5TotalWalker(t *testing.T) {
	root := NewRootCmd()
	root.InitDefaultHelpCmd()
	root.InitDefaultCompletionCmd()

	var total []string
	totalWalk(root, func(c *cobra.Command) {
		total = append(total, c.CommandPath())
	})
	sort.Strings(total)

	fmt.Printf("total walked set (%d):\n", len(total))
	for _, p := range total {
		fmt.Printf("  %s\n", p)
	}

	in := func(set []string, want string) bool {
		for _, s := range set {
			if s == want {
				return true
			}
		}
		return false
	}

	fmt.Printf("help in set:       %v\n", in(total, "intrastate help"))
	fmt.Printf("completion in set: %v\n", in(total, "intrastate completion"))
	for _, shell := range []string{"bash", "zsh", "fish", "powershell"} {
		p := "intrastate completion " + shell
		fmt.Printf("completion %-10s reached: %v\n", shell, in(total, p))
	}

	var planOnly []string
	totalWalk(root, func(c *cobra.Command) {
		if c.Flags().Lookup("plan-only") != nil {
			planOnly = append(planOnly, c.CommandPath())
		}
	})
	sort.Strings(planOnly)
	fmt.Printf("commands registering --plan-only (%d): %v\n", len(planOnly), planOnly)

	// Contrast: the repo's existing partial walker on the SAME materialized root.
	var partial []string
	walkCommandTree(root, func(c *cobra.Command) {
		partial = append(partial, c.CommandPath())
	})
	sort.Strings(partial)
	fmt.Printf("help_all.go::walkCommandTree set (%d):\n", len(partial))
	for _, p := range partial {
		fmt.Printf("  %s\n", p)
	}
	fmt.Printf("partial reaches completion: %v\n", in(partial, "intrastate completion"))
	fmt.Printf("partial reaches help:       %v\n", in(partial, "intrastate help"))

	// Bare root, before initializers, for the "no completion until init" fact.
	bare := NewRootCmd()
	var bareSet []string
	for _, sub := range bare.Commands() {
		bareSet = append(bareSet, sub.Name())
	}
	sort.Strings(bareSet)
	fmt.Printf("bare NewRootCmd() children: %v\n", bareSet)
}
```

## Output

```
=== RUN   TestA5TotalWalker
total walked set (15):
  intrastate
  intrastate completion
  intrastate completion bash
  intrastate completion fish
  intrastate completion powershell
  intrastate completion zsh
  intrastate docs
  intrastate flow
  intrastate flow next
  intrastate flow read-state
  intrastate flow resolve
  intrastate flow set-state
  intrastate help
  intrastate lint
  intrastate version
help in set:       true
completion in set: true
completion bash       reached: true
completion zsh        reached: true
completion fish       reached: true
completion powershell reached: true
commands registering --plan-only (0): []
help_all.go::walkCommandTree set (9):
  intrastate
  intrastate docs
  intrastate flow
  intrastate flow next
  intrastate flow read-state
  intrastate flow resolve
  intrastate flow set-state
  intrastate lint
  intrastate version
partial reaches completion: false
partial reaches help:       false
bare NewRootCmd() children: [docs flow help lint version]
--- PASS: TestA5TotalWalker (0.00s)
PASS
```

## Total vs partial walker

Both walkers ran against the SAME materialized root in the same test.

| command path | total walker | `help_all.go::walkCommandTree` |
|---|:--:|:--:|
| `intrastate` | yes | yes |
| `intrastate completion` | yes | no |
| `intrastate completion bash` | yes | no |
| `intrastate completion fish` | yes | no |
| `intrastate completion powershell` | yes | no |
| `intrastate completion zsh` | yes | no |
| `intrastate docs` | yes | yes |
| `intrastate flow` | yes | yes |
| `intrastate flow next` | yes | yes |
| `intrastate flow read-state` | yes | yes |
| `intrastate flow resolve` | yes | yes |
| `intrastate flow set-state` | yes | yes |
| `intrastate help` | yes | no |
| `intrastate lint` | yes | yes |
| `intrastate version` | yes | yes |
| **count** | **15** | **9** |

The existing walker at `internal/cli/help_all.go::walkCommandTree` skips
any child named `help` or `completion` — deliberately, since those carry
no user-authored extended body. That skip drops six commands, including
the entire `completion` subtree. Reusing it for a flag-scope assertion
would silently exempt exactly the commands a total sweep exists to cover,
so C1 needs its own walker rather than the existing one.

## Reading

A5's implementability half is CONFIRMED. The total walker is writable in
a handful of lines against the real `NewRootCmd()`, needs no cobra fork
or unexported access — both initializers are exported — and enumerates 15
commands including `help`, `completion`, and all four `completion`
shell children, proving the recursion is total rather than one level
deep. The `plan-only` scope assertion is expressible and evaluates to the
empty set today, which is the correct pre-implementation reading: once
`flow resolve` registers the flag the same assertion pins the set to
`{intrastate flow resolve}` and fails if any other command picks it up.
The contrast run shows the assertion must NOT be built on
`help_all.go::walkCommandTree`, whose name-based skip would hide six of
the fifteen commands.
