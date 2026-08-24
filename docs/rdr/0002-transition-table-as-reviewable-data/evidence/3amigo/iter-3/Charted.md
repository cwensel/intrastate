Model: claude-opus-5[1m]

# Charted to successor — 3amigo iteration 3, RDR 0002

- **A re-readable expanded-table dump grammar.** Origin: PM3-002. The rendered
  dump is lossy at one named site (set-valued atom literals and multi-entry
  fields render with unescaped separators, so `["a,b","c"]` and `["a","b,c"]`
  render identically). RDR 0002 answers the *blocked* question in place — Phase
  2's dump is deliberately not required to carry an escaping grammar, because
  defining one would make the dump a re-readable format whose inverse this RDR
  explicitly does not claim, and every consumer needing recoverability has the
  normalized value. **Why out of scope**: a re-readable dump is a new format
  with its own grammar, escaping rules, and round-trip invariant — a contract
  surface, not a clause. Absorbing it here would be the scope-expansion wormhole.
  **Suggested successor**: a dump-format RDR owning the rendered grammar and its
  read-back invariant. It did not exist while the `fidelity` table charted
  read-back "to a successor dump-format RDR"; that row now says so plainly
  instead of naming a phantom dependency. Substantial enough to warrant
  `/rdr-seed` — surfaced in the close packet.

- **Spike artifacts carry a stray `Model:` evidence stamp.** Origin: QA persona
  observation (recorded as a non-finding in `consolidation.md`, not a draft
  defect). `evidence/spikes/iter-2/output.txt` and `negative-cases.txt` both open
  with `Model: claude-opus-5[1m]` — the lens-evidence convention (rdr-common
  §model-stamp) leaking into spike artifacts. Harmless to the digest as computed
  (it is taken over the row block, not the whole file), but it travels into the
  production test tree when these fixtures are promoted. **Why out of scope**:
  evidence hygiene, touching no contract in this draft. **Suggested successor**:
  strip the stamp at promotion, or a kata against the spike-artifact convention.
