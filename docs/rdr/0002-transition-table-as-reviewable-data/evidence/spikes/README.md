# RDR 0002 normalizer spike

`iter-2/` is the **current** and normative iteration: the RDR and kata fixtures
authored to the JDR 0001 §D7 closed layout, with §D6 block-keyed routing and
general match-block expansion. Testing Strategy scenarios 1/2/3/5 cite it, and
`iter-2/output.txt` is the approved normative fixture.

The loose files beside this README are **iteration 1, superseded**. They author
`[accessors.<id>]` with `mode` and a tag-side `accessor` reference — two shapes
RDR 0002's Load-Bearing Decisions now list as *Rejected* — so the current
contracts refuse them outright:

```
$ (cd iter-2 && go run . ../rdr-fixture.toml)
refused: unknown schema field: accessors.rdr-finalized-at, accessors.rdr-status,
         tags.finalized_at.accessor, tags.status.accessor
```

They are retained as design history only. Do not promote them, and do not cite
them as evidence.
