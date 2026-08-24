#!/usr/bin/env python3
"""Regenerate every derived fixture from rdr-fixture.toml.

Each negative case is the RDR fixture with ONE mutation, so the category it
trips is the category it names (Testing Strategy scenario 3's
one-defect-per-fixture rule). Permutation cases are semantics-preserving.

Run: python3 gen-cases.py
"""
import os
import re

SRC = open('rdr-fixture.toml').read()

os.makedirs('neg', exist_ok=True)
os.makedirs('perm', exist_ok=True)


def neg(name, text):
    assert text != SRC, f'{name}: mutation was a no-op'
    open(os.path.join('neg', name + '.toml'), 'w').write(text)


def sub(old, new):
    assert old in SRC, f'anchor not found: {old[:60]!r}'
    return SRC.replace(old, new, 1)


# --- schema / version ------------------------------------------------------
neg('neg-unknown-field', sub('[model]\nid = "rdr"', '[model]\nid = "rdr"\nflavor = "spicy"'))
# A14's precedence oracle: v2-SHAPED -- version 2 PLUS a v2-only key the strict
# decoder would otherwise reject. Must refuse `unsupported version`, not
# `unknown schema field`.
neg('neg-v2-shaped', sub('version = 1', 'version = 2').replace(
    '[initial]\nstatus = "Draft"', '[schema_v2]\nlattice = "on"\n\n[initial]\nstatus = "Draft"', 1))
neg('neg-version', sub('version = 1', 'version = 3'))

# --- match-block operator restriction (§D6) --------------------------------
neg('neg-match-lt', sub(
    '[context.prelock.match.stage]\neq = "prelock"',
    '[context.prelock.match.stage]\neq = "prelock"\n[context.prelock.match.iter]\nlt = 3'))
neg('neg-match-exists', sub(
    '[rule.match.recognized]\neq = "reconcile-block"',
    '[rule.match.recognized]\neq = "reconcile-block"\n[rule.match.finalized_at]\nexists = true'))

# --- outcome binding -------------------------------------------------------
neg('neg-recognized-in-guard', sub(
    '[rule.guard.unless.profile]\neq = "small"',
    '[rule.guard.unless.profile]\neq = "small"\n[rule.guard.unless.recognized]\neq = "finalized"'))
neg('neg-outcome-outside', sub('eq = "reconcile-block"', 'eq = "not-in-alphabet"'))

# --- reserved <clear> (§D5) ------------------------------------------------
neg('neg-clear-as-write', sub('rewind_scope = "assumptions"', 'rewind_scope = "<clear>"'))
neg('neg-clear-in-initial', sub('stage = "propose"', 'stage = "<clear>"'))

# --- escape shape ----------------------------------------------------------
ESC = 'source = "rdr:escape"\nescape = ["no_match"]'
neg('neg-escape-with-gate', sub(ESC, 'source = "rdr:escape"\ngate = ["rdr-lock"]\nescape = ["no_match"]'))
neg('neg-escape-with-clear', sub(ESC, 'source = "rdr:escape"\nclear = ["prelock_lens"]\nescape = ["no_match"]'))
neg('neg-escape-with-write', sub(
    'escape = ["no_match"]\n[rule.match.recognized]\nin = ["round-clean", "reconcile-block"]',
    'escape = ["no_match"]\n[rule.match.recognized]\nin = ["round-clean", "reconcile-block"]\n[rule.write]\nstage = "archive"'))

# --- accessor bindings (provenance-scoped) ---------------------------------
READ_KEYS = 'keys = ["status", "stage", "profile", "iter", "rewind_scope", "prelock_lens", "cluster_ready"]\ntimeout = "2s"\n\n[read.rdr-finalized-at]'
neg('neg-owned-no-reader', sub(READ_KEYS,
    'keys = ["status", "stage", "profile", "iter", "rewind_scope", "prelock_lens"]\ntimeout = "2s"\n\n[read.rdr-finalized-at]'))
neg('neg-owned-two-readers', sub(
    '[read.rdr-finalized-at]\nrole = "rdr"\npath = "rdr.finalized_at"\nkeys = ["finalized_at"]',
    '[read.rdr-finalized-at]\nrole = "rdr"\npath = "rdr.finalized_at"\nkeys = ["finalized_at", "status"]'))
neg('neg-recognized-in-keys', sub('keys = ["finalized_at"]', 'keys = ["finalized_at", "recognized"]'))
neg('neg-accessor-no-timeout', sub(
    '[gate.rdr-lock]\nrole = "rdr"\npath = "rdr.lock"\nkeys = ["status"]\ntimeout = "2s"',
    '[gate.rdr-lock]\nrole = "rdr"\npath = "rdr.lock"\nkeys = ["status"]'))
neg('neg-gate-is-reader', sub('gate = ["rdr-lock"]', 'gate = ["rdr-status"]'))
# writer-provenance arm: a writer's keys list names an observed tag.
neg('neg-write-observed', sub(
    '[write.rdr-status]\nrole = "rdr"\npath = "rdr.status"\nkeys = ["status", "stage", "profile", "iter", "rewind_scope", "prelock_lens", "cluster_ready"]',
    '[write.rdr-status]\nrole = "rdr"\npath = "rdr.status"\nkeys = ["status", "stage", "profile", "iter", "rewind_scope", "prelock_lens", "cluster_ready", "finalized_at"]'))

# --- [initial] -------------------------------------------------------------
neg('neg-initial-unknown', sub('[initial]\nstatus = "Draft"', '[initial]\nnosuchtag = "x"\nstatus = "Draft"'))
# the VALUE arm -- the only arm of `malformed initial declaration` that is
# reachable (the provenance arm has no authorable form; see the two probes).
neg('neg-initial-bad-value', sub('stage = "propose"', 'stage = "nosuchstage"'))
neg('neg-initial-int-out-of-range', sub('iter = 0', 'iter = 99'))
# Enumeration witnesses for the unauthorable provenance arm: an observed key in
# [initial] has exactly two authorings, and BOTH refuse under a writer-binding
# category before the owned-tag predicate is reached.
neg('probe-a-initial-observed-no-writer',
    sub('[initial]\nstatus = "Draft"', '[initial]\nfinalized_at = "x"\nstatus = "Draft"'))
neg('probe-b-initial-observed-with-writer',
    sub('[initial]\nstatus = "Draft"', '[initial]\nfinalized_at = "x"\nstatus = "Draft"').replace(
        'keys = ["status", "stage", "profile", "iter", "rewind_scope", "prelock_lens", "cluster_ready"]\ntimeout = "2s"\nread_back = true',
        'keys = ["status", "stage", "profile", "iter", "rewind_scope", "prelock_lens", "cluster_ready", "finalized_at"]\ntimeout = "2s"\nread_back = true', 1))

# --- domain / type model ---------------------------------------------------
neg('neg-literal-outside-domain', sub(
    '[context.draft.match.status]\neq = "Draft"', '[context.draft.match.status]\neq = "Archived"'))

# --- identity / separator --------------------------------------------------
neg('neg-ruleid-hash', sub('id = "reconcile-rewind"', 'id = "reconcile#rewind"'))
neg('neg-alphabet-hash', sub('"reconcile-block", "finalized"]', '"reconcile-block", "final#ized"]'))
neg('neg-in-member-hash', sub('in = ["large", "foundational"]', 'in = ["large", "found#ational"]'))
neg('neg-duplicate-rule-id', sub('id = "continue-prelock-cluster"', 'id = "continue-prelock"'))

# --- alphabet --------------------------------------------------------------
neg('neg-empty-alphabet-member', sub('outcomes = ["round-clean"', 'outcomes = ["", "round-clean"'))
neg('neg-dup-alphabet-member', sub('outcomes = ["round-clean"', 'outcomes = ["round-clean", "round-clean"'))

# --- rule shape / references ----------------------------------------------
neg('neg-no-write-block', sub(
    '[rule.match.recognized]\nin = ["finalized", "verdict-flapping"]\n[rule.write]\nstage = "archive"',
    '[rule.match.recognized]\nin = ["finalized", "verdict-flapping"]'))
neg('neg-unknown-tag-match', sub('[context.draft.match.status]', '[context.draft.match.nosuchtag]'))
neg('neg-unknown-context', sub('use = ["large-prelock"]\nsource = "rdr:prelock"', 'use = ["nope"]\nsource = "rdr:prelock"'))
neg('neg-terminal-unknown', sub('terminal = ["archived"]', 'terminal = ["nonexistent"]'))
neg('neg-recognized-misnamed', sub('[tags.recognized]\nprovenance = "recognized"', '[tags.outcome]\nprovenance = "recognized"'))
neg('neg-bad-exists', sub('exists = false', 'exists = "nope"'))
# byte-exact key identity negative control: a normalizer that case-folds passes
# every positive test, so the control must be a mis-cased REFERENCE.
neg('neg-case-folded-tag', sub('[context.draft.match.status]', '[context.draft.match.Status]'))

# --- positive controls (MUST load) ----------------------------------------
open('pos-no-initial.toml', 'w').write(re.sub(r'\[initial\]\n(?:[a-z_]+ = .*\n)+\n', '', SRC))
open('pos-no-terminal.toml', 'w').write(SRC.replace('terminal = ["archived"]\n', '', 1))

# --- permutations (semantics-preserving) ----------------------------------
key_perm = SRC.replace(
    '[tags.status]\nprovenance = "owned"\nkind = "enum"\ndomain = ["Draft", "Final"]\nsingle_valued = true\nrequired = true',
    '[tags.status]\nrequired = true\nsingle_valued = true\ndomain = ["Final", "Draft"]\nkind = "enum"\nprovenance = "owned"', 1).replace(
    '[read.rdr-status]\nrole = "rdr"\npath = "rdr.status"\nkeys = ["status", "stage", "profile", "iter", "rewind_scope", "prelock_lens", "cluster_ready"]\ntimeout = "2s"',
    '[read.rdr-status]\ntimeout = "2s"\nkeys = ["cluster_ready", "prelock_lens", "rewind_scope", "iter", "profile", "stage", "status"]\npath = "rdr.status"\nrole = "rdr"', 1).replace(
    '[rule.write]\nstage = "prelock"\niter = 2', '[rule.write]\niter = 2\nstage = "prelock"', 1)
assert key_perm != SRC
open('perm/rdr-keyorder.toml', 'w').write(key_perm)

dump_idx = SRC.index('[dump]')
body, tail = SRC[:dump_idx], SRC[dump_idx:]
first = body.index('[[rule]]')
blocks = ['[[rule]]' + p for p in body[first:].split('[[rule]]')[1:]]
assert len(blocks) == 5, len(blocks)
open('perm/rdr-ruleorder.toml', 'w').write(
    body[:first] + ''.join([blocks[4], blocks[2], blocks[0], blocks[3], blocks[1]]) + tail)

open('perm/rdr-eq-as-in.toml', 'w').write(
    SRC.replace('[rule.match.recognized]\neq = "reconcile-block"',
                '[rule.match.recognized]\nin = ["reconcile-block"]', 1))

# --- merge semantics (contexts clause; A13) -------------------------------
# idempotence half: rule + inherited context contribute a BYTE-IDENTICAL atom.
open('merge-idempotent.toml', 'w').write(SRC.replace(
    '[rule.match.recognized]\neq = "reconcile-block"',
    '[rule.match.recognized]\neq = "reconcile-block"\n[rule.match.status]\neq = "Draft"', 1))
# distinct-literal half: two contexts, same (key, block, operator), different
# literals -- BOTH must survive. Authored in both `use` orders.
distinct = SRC.replace(
    '[context.archived.match.stage]\neq = "archive"',
    '[context.archived.match.stage]\neq = "archive"\n\n[context.draft-final.match.status]\neq = "Final"', 1).replace(
    'use = ["draft"]\nsource = "rdr:reconcile"', 'use = ["draft", "draft-final"]\nsource = "rdr:reconcile"', 1)
open('merge-distinct.toml', 'w').write(distinct)
open('merge-distinct-rev.toml', 'w').write(
    distinct.replace('use = ["draft", "draft-final"]', 'use = ["draft-final", "draft"]', 1))

print('regenerated:', len(os.listdir('neg')), 'negative,',
      len(os.listdir('perm')), 'permutation, 2 positive, 3 merge')
