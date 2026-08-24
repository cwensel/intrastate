Model: claude-opus-5[1m]

# Charted to successors — critique iter-3

- **M-5 — RDR 0009 quotes RDR 0002 text that no longer exists, and a `Verified`
  assumption rests on it.** RDR 0009 (`Final`) quotes, at `0009:244` and
  `0009:492`, that RDR 0002's normalization renders "a candidate row **with row
  kind `escape`**, its normal predicate set, source rule id, source locator, and
  modeled failure class list." RDR 0002 has since deleted and inverted that
  sentence: row kind is now "a **derived view property, never a row field**".
  The second site is inside RDR 0009 A4's **Consistency note**, which is marked
  `Verified` and reasons from the quoted text ("names **neither** writes **nor**
  next-tags … Stricter upstream, consistent — not a conflict"). The conclusion
  survives on the current text, but the evidence for it does not.

  Out of scope here: this RDR's own clause is correct and self-consistent — it
  forbids a kind field in the normalized value and kernel row while explicitly
  permitting a render-time view column, which is what satisfies RDR 0006's
  minimum-input contract (0006 consumes a graph view, not a kernel row). Nothing
  in RDR 0002 needs to change, and RDR 0002 cannot edit a `Final` peer.

  Suggested successor: **`/rdr-cluster-reconcile`** over 0002-0009. Note that
  both RDR 0006 and RDR 0009 carry "checked consistent at the 0002-0009
  iteration-3 gate" in their Status lines, so the existing gate did not catch
  this — the reconcile should re-run the citation check, not just re-affirm the
  stamp. This is the sharpest item this pass produced and the only one requiring
  a document other than RDR 0002 to change.

- **A citation-conformance check** (tooling-pass). Three of this pass's
  highest-value findings (M-2, M-3, M-4) were fabricated or inverted citations of
  `Final` peers — "RDR 0009 A4 fixes…", "RDR 0004 removes the key…", "RDR 0007
  spells…" — each stating as settled fact something the cited peer says the
  opposite of, or does not say at all. All three survived iter-1, iter-2, cove,
  and 3amigo, and all three were found in minutes by opening the peer file. The
  charted item is mechanical: for each `RDR NNNN <assumption-id>` or `RDR NNNN
  <verb>` citation in a draft, assert the cited document contains a
  correspondingly-worded claim. Cheap, scriptable, and currently absent from the
  Stage-7 tooling sweep. M-5 above is the same defect class pointing the other
  way (a peer citing *this* RDR staleley), so the check wants to run in both
  directions across a cluster.

- **A fixture-conformance check** (tooling-pass, partially discharged here).
  M-1 — the dump column vocabulary fixed to spellings the normative fixtures do
  not use — was created by the immediately preceding lens and would have been
  caught by loading the fixtures against the finished loader. A scenario-1
  control was added to this RDR for its own fixtures, which discharges the local
  case. The general item stays charted: any RDR that declares a fixture
  normative and then writes clauses about its contents owes an executed
  conformance assertion, not a review step. Related to the resolve prompt's
  "COMPUTE, DON'T ARGUE" rule, which this defect violated exactly.
