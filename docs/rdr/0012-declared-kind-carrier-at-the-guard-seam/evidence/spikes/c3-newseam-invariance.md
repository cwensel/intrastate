# Spike — C3's `newSeam` parameter and Go func return-invariance

Raised by: repeatability lens, iter-2 delta re-run (run-B, claude-fable-5-1).

## Question

C3 says `guard.Evaluator`'s suite call migrates to
`TestGuardEvaluatorContract(t, guard.NewEvaluator)`. C3's own parameter is
`newSeam func(kinds map[string]string) GuardEvaluator`, while C1 fixes
`func NewEvaluator(kinds map[string]string) Evaluator`. Does that assign?

## Method

Minimal reproduction of the shape (interface parameter, concrete-returning
constructor), compiled with the toolchain on `main`.

```go
type Iface interface{ M() int }
type Concrete struct{}
func (c Concrete) M() int { return 1 }
func NewConcrete(m map[string]string) Concrete { return Concrete{} }
func Suite(newSeam func(kinds map[string]string) Iface) { _ = newSeam }
func main() { Suite(NewConcrete) }
```

## Result

`go build` FAILS:

```
cannot use NewConcrete (value of type func(m map[string]string) Concrete)
  as func(kinds map[string]string) Iface value in argument to Suite
```

Go function types are invariant in the return position: a constructor
returning the concrete type does not satisfy a parameter typed to return the
interface, even though the concrete type implements it.

## Disposition

C3 amended: the migrating call sites wrap the constructor in a closure
(`func(k map[string]string) resolve.GuardEvaluator { return guard.NewEvaluator(k) }`)
rather than passing it bare. `NewEvaluator`'s return type is NOT widened to the
interface — C1 fixes the value shape deliberately, and widening it would make
every non-test construction site hold an interface where it holds a struct.
