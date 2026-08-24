Model: claude-fable-5

# Answer-vs-fences check — sibling 0004 (accessor-execution-safety-model), iteration 3

Date: 2026-08-24. Read-only. Inputs: home `docs/jdr/0001-resolve-kernel-seam.md`
§D5 (222-255), §D7 (300-380), Interface record JD-15 (534-539) and JD-17 (545-551);
sibling `docs/rdr/0004-accessor-execution-safety-model.md` at 526c481 (940 lines), read in full.
Line anchors are `0004:<n>` / `JDR:<n>`.

## Status line and Prerequisites

- `0004:9` — `Status: Final [joint decision → JDR 0001 §JD-15, §JD-17: \`<clear>\` write/read-back
  semantics; accessor metadata keys in 0002's layout]`. The qualifier names both entries as
  pending joint decisions; it does not say they are answered (both decided 2026-08-24).
- Prerequisites checklist `0004:771-780` — all four boxes still `- [ ]` (unchecked): A1-A8
  Verified / A9, A10 Pending-DOWNGRADED; RDR 0001 stateless; RDR 0002 carries accessor
  references / provenance / roles; RDR 0003 consumes without executing.
- The word "clear" occurs exactly once in 0004 — in the Status qualifier at `0004:9`. The home's
  premise at `JDR:227` ("0004 contains the word "clear" zero times") is true of the body; no
  clause in the body mentions `<clear>`.

## JD-15 (§D5) — `<clear>` reserved; write removes key; read-back asserts absence

Home text checked: `JDR:236-240` (option (a): "a `<clear>` write removes the key, that read-back
asserts the key is **absent** (already a value on its success branch), that clearing an absent key
succeeds, and that a read yielding `<clear>` as a value is unreadable") and `JDR:251-255`
(*Lands in 0004* list); JD-15 entry `JDR:534-539`.

| # | 0004 clause (anchor, fenced?) | Quote | Verdict | Note |
| --- | --- | --- | --- | --- |
| 15.1 | Write-scope fence `0004:323-326`, fenced | "A write accessor MUST apply only planned owned-tag writes produced by a successful transition. It MUST NOT write observed or recognized tags." | CONSISTENT | A rule-level `clear` is a planned owned-tag write in 0002's rendering (`Writes` carries `<clear>`); the fence is silent on removal vs. assignment and the answer fills it. |
| 15.2 | Read-back fence `0004:328-335`, fenced | "re-read the same caller-supplied artifact role ... and verify the expected owned-tag values and that observed and recognized tag values present before the write are unchanged ... A read-back mismatch MUST be reported as a write failure." | CONSISTENT | "expected owned-tag values" does not say *present*; an expectation of absence is representable. The home relies on exactly this ("already a value on its success branch"). Meaning is extended, not changed. |
| 15.3 | Completeness fence `0004:292-300`, fenced | "A missing key MUST be distinguishable from an unread key only by which branch is taken — absence is a value, unreadability is a refusal." | CONSISTENT | Supplies the "absence is a value" the read-back-absent clause needs. |
| 15.4 | Seam-omission fence `0004:310-316`, fenced | "what crosses the seam MUST leave the key absent from the owned snapshot" | CONSISTENT | A cleared key re-read after the write crosses as omission — same encoding the answer requires. |
| 15.5 | Read fence `0004:281-284`, fenced | "A read accessor MUST return typed tag values or a typed refusal." | CONSISTENT (silent) | Says nothing about a `<clear>` literal on the value side; "unreadable-on-read" would add a refusal cause, not a third branch. |
| 15.6 | Approach prose `0004:213-215`, **unfenced** | "then re-reads the same role and verifies that the expected owned-tag values are **present**." | CONTRADICTS (unfenced, wording) | For a clearing write the expected outcome is that the key is *absent*; "are present" is falsified for that case. Meaning changes only for the clear case (extension); non-clear writes unaffected. |
| 15.7 | Fidelity table `0004:511`, **unfenced** (Round-Trip/Inverse invariant restated) | "Every planned owned-tag key/value present in the re-read equals the transition plan's expected value \| value equality over the planned key set \| none" | CONTRADICTS (unfenced, wording) | "present in the re-read" cannot hold for a cleared key; the "none" lossy-exemption column would need an absence arm. Round-Trip prose at `0004:425-428` ("must return the expected owned-tag values") is silent and CONSISTENT. |
| 15.8 | Load-bearing "Absent vs unreadable" `0004:412-421`, unfenced | "the key's value failed type coercion ... is *unreadable*" | CONSISTENT (silent) | A reserved-sentinel read could be classed here by extension; nothing states it. |
| 15.9 | Disposition table `0004:440-457`, unfenced | rows "Write succeeds, owned tag matches / differs" | CONSISTENT (silent) | No row for a clearing write; nothing contradicts. |

Fenced text: all CONSISTENT. Unfenced: two wording sites (`0004:215`, `0004:511`) say "present"
where a clear requires "absent".

*Lands in 0004* items (`JDR:253-255`):

| Item | Status in 0004 | Evidence |
| --- | --- | --- |
| `<clear>` write removes the key | **absent** | no clause; `<clear>` appears only at `0004:9` |
| read-back asserts the key is absent | **absent** (and contradicted in unfenced wording at `0004:215`, `0004:511`) | fence `0004:329-334` silent |
| clearing an absent key succeeds (idempotent clear) | **absent** | no clause; Disposition table `0004:440-457` has no such row |
| a read yielding `<clear>` as a value is unreadable | **absent** | `0004:412-421` lists unreadable causes without it |
| MVV scenario for a clearing rule | **absent** | MVV `0004:784-797` and Scenarios 1-8 `0004:849-908` contain none |

## JD-17 (§D7) — capability tables, `keys` binding, read-back equality, "cannot be rebound"

Home text checked: `JDR:333-346` (ii Accessors), `JDR:356-358` (iv write-replaces), `JDR:361-370`
(*Lands in* list: "**0004** — its definition is the capability-table entry; read-back is equality;
writers carry `keys`; the "cannot be rebound" prose repaired"); JD-17 entry `JDR:545-551`.

| # | 0004 clause (anchor, fenced?) | Quote | Verdict | Note |
| --- | --- | --- | --- | --- |
| 17.1 | Capability fence `0004:264-268`, fenced | "Every accessor definition MUST declare exactly one capability: read, gate, or write." | CONSISTENT | Home: "0004's 'exactly one capability' holds structurally" (`JDR:336-338`). One table per capability satisfies the fence. |
| 17.2 | Identity fence `0004:270-274`, fenced | "each `(flow id, accessor name, capability)` identity MUST resolve to exactly one accessor binding." | CONSISTENT | Home cites it verbatim (`JDR:344-345`); same id in `[read.x]` and `[write.x]` is two identities. |
| 17.3 | Load-bearing Identity `0004:370-372`, **unfenced** | "The same name cannot be rebound to a different capability within one flow." | **CONTRADICTS** (unfenced; meaning changes) | Home: "The same id in `[read.x]` and `[write.x]` is legal ... its unfenced 'cannot be rebound' sentence is a citation repair" (`JDR:344-346`). The sentence also conflicts with 0004's own fence 17.2, which keys identity on capability. Repair **absent** — sentence still stands. |
| 17.4 | Requested-key fence `0004:286-290`, fenced | "Every read accessor definition MUST declare the requested key set as validated metadata." | CONSISTENT (silent on writers) | Home: "Readers and writers both declare `keys`" (`JDR:340`). The fence binds readers only; it does not forbid writer `keys`. |
| 17.5 | Technical Design `0004:226-229`, unfenced | "Each definition declares a stable name, capability, artifact role, expected tag keys, timeout policy, and whether read-back verification is required." | CONSISTENT | "Each definition ... expected tag keys" already covers writers in prose; the home's per-capability table is a shape refinement. |
| 17.6 | Timeout fences `0004:353-361`, fenced | "Missing or non-positive timeout metadata MUST fail validation before execution." | CONSISTENT | Home: "`timeout` (Go duration string; missing or non-positive refused)" (`JDR:342-343`). |
| 17.7 | Read-back fence `0004:328-335`, fenced | "verify the expected owned-tag values" | CONSISTENT | Silent on equality vs. containment; home's "read-back compares the held value for **equality**" (`JDR:358`) fills it. |
| 17.8 | Fidelity `0004:511`, MVV `0004:793-794`, Testing `0004:845-846`, unfenced | "equals the transition plan's expected value" | CONSISTENT | Equality already stated. |
| 17.9 | Write-scope fence `0004:323-326`, fenced | "MUST NOT write observed or recognized tags." | CONSISTENT | Home loader refuses "a writer key that is not owned" (`JDR:341`); 0004 already has `write_non_owned_tag` (`0004:492`). Two sites for one check (loader in 0002 vs. validator in 0004) — overlap, not contradiction. |
| 17.10 | Wire-format decision `0004:373-374`, unfenced | "RDR 0002 owns the TOML carrier. This RDR owns the accessor execution semantics embedded behind accessor references." | CONSISTENT | Home removes the tag-side `accessor` reference (`JDR:338-340`); 0004 never names `[accessors.<id>]`, `mode`, or the tag-side key, so nothing to falsify. |
| 17.11 | Selection decision `0004:378-380`, unfenced | "when an accessor reference names a capability, the executor selects only a binding with the same accessor identity and capability." | CONSISTENT (silent) | Under capability tables the table name is the capability; direction of reference (accessor→keys) is unstated here. |
| 17.12 | Scenario 1 arms `0004:849-855`, `0004:492`, unfenced | eight validation arms incl. `missing_accessor`, `multiply_bound_accessor` | CONSISTENT (silent) | The home's binding validations ("key served by zero or two readers, a written or cleared key not in exactly one writer") are 0002-loader rules; 0004 has no arm for them and none contradicting. |
| 17.13 | Read-back-required `0004:228-229` "whether read-back verification is required" vs fence `0004:329` "MUST re-read" | — | CONSISTENT | Home fixes `read_back = true` on writers; matches the fence (mandatory), the prose "whether ... required" is looser but not contradicting. |

Fenced text: all CONSISTENT. Unfenced: one CONTRADICTS (`0004:371`), pre-identified by the home as a
citation repair.

*Lands in 0004* items (`JDR:366-368`):

| Item | Status in 0004 | Evidence |
| --- | --- | --- |
| its definition is the capability-table entry (`[read.<id>]`/`[write.<id>]`/`[gate.<id>]`) | **absent** (not contradicted) | `0004:226-229`, `0004:373-374` say definitions are loader data owned by 0002; no layout named |
| read-back is equality | **present** (unfenced) | `0004:511`, `0004:793-794`, `0004:845-846` "equals"; fence `0004:331` silent |
| writers carry `keys` | **absent** (not contradicted) | fence `0004:287` reader-only; `0004:228` "each definition ... expected tag keys" is generic prose |
| the "cannot be rebound" prose repaired | **absent — contradicted** | `0004:371` still reads "The same name cannot be rebound to a different capability within one flow." |

## Observations (evidence only, not disposition)

1. Home §D7(ii) carries an internal tension worth the gate's eye: `JDR:336-337` "`keys` exists only
   on readers and `read_back` only on writers" vs `JDR:340` "Readers and writers both declare
   `keys`". 0004's fence `0004:287` agrees with the first sentence; the *Lands in* item ("writers
   carry `keys`") follows the second.
2. The home's D5 premise "read-back ... passes green on a tag that was never removed" (`JDR:228-230`)
   still describes 0004 as written: nothing at `0004:329-335` or `0004:511` would fail a literal
   `<clear>` string stored and read back.
3. Nothing in 0004's fenced text needs to move for either answer; every landing item is additive
   except the `0004:371` sentence, and the two "present" wordings at `0004:215`/`0004:511`.
