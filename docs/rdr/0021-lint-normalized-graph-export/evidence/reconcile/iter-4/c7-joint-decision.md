Model: claude-sonnet-5

# C-7 joint decision — re-check (iter-4)

## Q1 — 0029:C4 sequencing (verbatim, from the `cli/0021:C2` paragraph added to C4)

"A consumer reads the envelope's major to decide whether it can parse at
all, and the document's `schema` to decide what the payload means. That
ordering is normative, and it settles the two mixed cases: the envelope
is the FRAME, so an unsupported envelope major is rejected under C1
before `data` is read at all — the document's marker is not consulted
and cannot rescue it, because a consumer that cannot parse the frame has
no warranted reading of what it contains. In the inverse case — envelope
major supported, document `schema` marker unrecognized — the envelope
parses, its non-`data` members (including `schema_version`, `type` and
any advisory fields) are trustworthy, and `data` alone is opaque: the
consumer MUST NOT infer the document's shape from the envelope's version,
since C1's rule is that `schema_version` is never projected into `data`
and therefore says nothing about the payload's own evolution. A consumer
pinning the document marker refuses `data` and keeps the envelope; one
that does not pin it MAY pass `data` through untouched."

Verdict: **ANSWERED**. C4 now states the read order (envelope major
first, document `schema` second) and covers both mixed-support cases
(unsupported envelope major; supported envelope + unrecognized document
marker).

## Q2 — 0029 status / open-joint tracker

`rdr status --tags 0029`:
- status=Implemented
- status_form=none
- open_joint_decisions=[]
- joint_checks=1
- joint_check_home=homed

`rdr index --open-joint`:
```
total 0 open joint decisions over 9 records
```
This is a clean-tracker artifact, not evidence C-7 was unanswered at the
time of tracking — the tolerance lived in re-entry notes (deleted at
re-lock), never in a Status qualifier, exactly as the cluster-reconcile
report anticipated.

## Q3 — cluster-reconcile report, pair 0021-0029

Path: docs/rdr/cluster-reconcile/0021-0029/reconcile-report.md:64,115-138;
open.txt.

Finding C-7 recorded: "Neither record sequences the two version markers:
which does an agent read first, and what happens when one is supported
and the other is not" — classed `gap`/`risks-impl`/**JOINT-DECISION**,
Home `cli/0029 §Normative Contracts` C4.

Tolerance recorded at demotion: "both records are demoted by this report,
so neither Status line can carry the `Final [joint decision → …]`
qualifier — that slot holds the demotion qualifier. The tolerance is
recorded in both re-entry notes instead, and the qualifier applies at
re-lock if C4 has not answered by then. 0029 answers it in C4 as one
paragraph; 0021 cites the home rather than restating it."

`open.txt` carried `C-7` as the sole undispositioned row at that
iteration, pending exactly this C4 paragraph.

## Conclusion

C4 as it now reads discharges the tolerance the reconcile report set:
"0029 answers it in C4 as one paragraph" — confirmed, that paragraph
exists and sequences both reads plus both mixed cases. 0021's citation of
C4 as HOME resolves. This is PASS: the obligation was at the HOME (0029),
0021 correctly claims nothing about read order itself, and the home has
now answered it. Not a blocker for 0021.
