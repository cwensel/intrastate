Model: claude-opus-5[1m]

# A3 verification spike: wrapper argv-word transparency + widened-predicate coverage

Run 2026-09-03 in the main Resolve context (no sub-agent).

## Claim under test (A3, first half)

Every wrapper form the Problem Statement names places the interpreter name and
its inline-code flag as SEPARATE argv words, in that order, so the position-free
predicate refuses each.

## Part 1 — argv-word transparency, probed against a real no-shell exec

Harness: `exec.Command(argv[0], argv[1:]...)`, no shell, printing the argv vector
as constructed plus the wrapper's outcome. The DECISIVE column is the argv vector:
the predicate reads words, so whether the wrapper binary exists on this darwin host
is immaterial to A3.

```
argv=["nice" "sh" "-c" "echo MARKER_RAN"]
  -> ok | out="MARKER_RAN"
argv=["nohup" "sh" "-c" "echo MARKER_RAN"]
  -> ok | out="MARKER_RAN"
argv=["setsid" "sh" "-c" "echo MARKER_RAN"]
  -> exec: "setsid": executable file not found in $PATH | out=""
argv=["stdbuf" "-o0" "sh" "-c" "echo MARKER_RAN"]
  -> ok | out="MARKER_RAN"
argv=["timeout" "5" "sh" "-c" "echo MARKER_RAN"]
  -> exec: "timeout": executable file not found in $PATH | out=""
argv=["xargs" "sh" "-c" "echo MARKER_RAN"]
  -> ok | out=""
argv=["env" "-i" "sh" "-c" "echo MARKER_RAN"]
  -> ok | out="MARKER_RAN"
argv=["env" "-u" "FOO" "sh" "-c" "echo MARKER_RAN"]
  -> ok | out="MARKER_RAN"
argv=["chpst" "sh" "-c" "echo MARKER_RAN"]
  -> exec: "chpst": executable file not found in $PATH | out=""
argv=["doas" "sh" "-c" "echo MARKER_RAN"]
  -> exec: "doas": executable file not found in $PATH | out=""
```

Every vector presents `sh` and `-c` as separate, ordered words. Four wrappers
(`setsid`, `timeout`, `chpst`, `doas`) are Linux/BSD tools absent from stock darwin
and report `executable file not found`; their argv layout is nonetheless correct and
is all the predicate reads. `xargs` ran with empty output because it consumed empty
stdin — word layout still correct.

## Part 2 — old vs widened predicate over those exact vectors

Both predicates transcribed verbatim from `internal/table/load.go` (old) and
RDR 0027 C1 (new); `shellInterpreters` and `filepathBase` copied as-is.

```
argv | old | new
---
["env"                                          "-i"                                           "sh"                                           "-c"                                           "echo hi"                                     ] old=(false,"") new=(true,"sh -c")   <== CHANGED
["env"                                          "-u"                                           "FOO"                                          "sh"                                           "-c"                                           "echo hi"                                     ] old=(false,"") new=(true,"sh -c")   <== CHANGED
["nice"                                         "sh"                                           "-c"                                           "echo hi"                                     ] old=(false,"") new=(true,"sh -c")   <== CHANGED
["timeout"                                      "5"                                            "sh"                                           "-c"                                           "echo hi"                                     ] old=(false,"") new=(true,"sh -c")   <== CHANGED
["xargs"                                        "sh"                                           "-c"                                           "echo hi"                                     ] old=(false,"") new=(true,"sh -c")   <== CHANGED
["nohup"                                        "sh"                                           "-c"                                           "echo hi"                                     ] old=(false,"") new=(true,"sh -c")   <== CHANGED
["setsid"                                       "sh"                                           "-c"                                           "echo hi"                                     ] old=(false,"") new=(true,"sh -c")   <== CHANGED
["stdbuf"                                       "-o0"                                          "sh"                                           "-c"                                           "echo hi"                                     ] old=(false,"") new=(true,"sh -c")   <== CHANGED
["chpst"                                        "sh"                                           "-c"                                           "echo hi"                                     ] old=(false,"") new=(true,"sh -c")   <== CHANGED
["doas"                                         "sh"                                           "-c"                                           "echo hi"                                     ] old=(false,"") new=(true,"sh -c")   <== CHANGED
["nice"                                         "sh"                                           "-c"                                           "cat {artifact}"                              ] old=(false,"") new=(true,"sh -c")   <== CHANGED
["sh"                                           "./gate.sh"                                   ] old=(false,"") new=(false,"")
["env"                                          "-S"                                           "sh -c echo"                                  ] old=(false,"") new=(false,"")
["sh"                                           "-s"                                          ] old=(false,"") new=(false,"")
```

## Finding

All ten named wrapper forms flip `false -> true` under the widened predicate, each
reporting the form `sh -c`. The `{artifact}`-bearing mutant `["nice","sh","-c","cat
{artifact}"]` also flips to `true`, which is what withdraws it from clause 3's
placeholder arm (`isInterp` becomes true, so the whitespace exemption applies and
clause 4 owns it) — the wrong-defect masking A3's second half names.

The three admitted forms — `["sh","./gate.sh"]`, `["env","-S","sh -c echo"]`,
`["sh","-s"]` — stay `false` under both, as C1 states.

## Verdict

A3 first half: VERIFIED by derivation over the argv vectors (not deferred to the MVV).
A3 second half (the reported CATEGORY under a wrapper, i.e. that clause 3 yields to
clause 4): verified structurally here — `isInterp` is the exemption's gate at
`internal/table/load.go:1114` and the refusal's trigger at `:1124`, both fed by the
single `interpreterForm` call at `:1095` — and confirmed end-to-end by the MVV's
step-2 `{artifact}` mutant at implementation.
