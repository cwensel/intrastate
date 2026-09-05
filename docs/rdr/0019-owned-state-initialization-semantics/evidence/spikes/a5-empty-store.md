Model: claude-opus-5

# A5 - the empty-store predicate and the emptied-vs-absent boundary

## Assumption under test

The load-bearing sub-claim: an artifact store emptied by clearing its last
key is INDISTINGUISHABLE at the content level from a never-written / absent
one, so an empty-store predicate can be stated purely as "the store carries
zero keys". Correspondingly, a store retaining at least one key must be
observably non-empty, so the predicate declines to write.

`flow init-state` does not exist yet. The spike therefore drove the SHIPPED
`flow set-state` / `flow read-state` verbs (the real CLI, built binary) AND
additionally exercised `flowbind.load` and `Writer.Apply` directly in a Go
test, to observe the store values the CLI only shows indirectly. Both were
run; they agree.

## What the source says

`internal/cli/flowbind/flowbind.go`

- `load` returns `store{}` (a non-nil, zero-key map) for `os.ErrNotExist`,
  and also for a whitespace-only file. No error, no sentinel.
- `Writer.Apply` treats a `<clear>` planned value as `delete(s, t.Key)` -
  key REMOVAL, not a tombstone assignment.
- `save` always writes the file, including for a zero-key store; the
  encoder is `clierr.WriteJSONLine`, so an emptied store persists as the
  three bytes `{}` plus a newline. There is no header, no version marker,
  no schema field, no tombstone list.
- One key CAN survive a full clear and is the only residue the format has:
  `sealedKey`, the NUL-prefixed `flow.readback-unreachable`, written by
  `Apply` when the writer's DECLARED locator is unreachable, and deleted by
  the next reachable-locator write.

## Commands and per-step results

Binary: `go build -o /tmp/a5-intrastate ./cmd/intrastate` (BUILD_OK,
v0.0.0-...-5f52ba5303f7+dirty).

Fixture model: `a5-fixture-model.toml` - the shipped
`flow_fixtures_0005_test.go` MVV model, trimmed to the state role, with
`write.state` serving three owned keys (`status`, `labels`, `stale`). Two
of them are written here.

### (i) baseline - ABSENT artifact  -> `a5-step-i-absent.txt`

    $ ls -la /tmp/a5run/absent.json
    ls: /tmp/a5run/absent.json: No such file or directory
    $ intrastate flow read-state --model model.toml \
        --artifact state=/tmp/a5run/absent.json --as json
    {"type":"ok","data":{...,"readers":[{"id":"state",
      "keys":["status","labels","stale"],"tags":{}}]}}
    exit=0

`read-state` succeeds over an absent artifact and reports `tags:{}` - zero
keys, every requested key established-absent. The file is NOT created by
the read.

Go-level: `load(absent) = flowbind.store{} len=0`.

### (ii) write two owned keys into a fresh artifact -> `a5-step-ii-write-two.txt`

    $ intrastate flow set-state --model model.toml --artifact state=state.json \
        --write status=final --write stale=keepme --as json
    {"type":"ok","data":{...,"writes":{"stale":"keepme","status":"final"},
      "clear":[],"owned":{"stale":"keepme","status":"final"}}}
    exit=0
    $ xxd state.json
    00000000: 7b22 7374 616c 6522 3a22 6b65 6570 6d65  {"stale":"keepme
    00000010: 222c 2273 7461 7475 7322 3a22 6669 6e61  ","status":"fina
    00000020: 6c22 7d0a                                l"}.

36 bytes, mode 0600, flat JSON object. read-state reports
`tags:{"stale":"keepme","status":"final"}`.

### (iii) clear ONE key; one key remains -> `a5-step-iii-clear-one.txt`

    $ intrastate flow set-state ... --clear stale --as json
    {"type":"ok","data":{...,"writes":{},"clear":["stale"],"owned":{}}}
    exit=0
    $ xxd state.json
    00000000: 7b22 7374 6174 7573 223a 2266 696e 616c  {"status":"final
    00000010: 227d 0a                                  "}.
    $ read-state -> "tags":{"status":"final"}

The cleared key is GONE from the bytes - no tombstone, no null, no empty
string. It reads absent. `status` remains, so the store is observably
non-empty: `len(load(art)) == 1`. An empty-store predicate correctly
DECLINES here. This is the "preserves cleared keys" case, and it holds.

### (iv) clear the REMAINING key - the critical measurement -> `a5-step-iv-clear-last.txt`

    $ intrastate flow set-state ... --clear status --as json
    {"type":"ok","data":{...,"writes":{},"clear":["status"],"owned":{}}}
    exit=0
    $ ls -la state.json
    -rw-------  1 cwensel  wheel  3 ... /tmp/a5run/state.json
    $ cat state.json
    {}
    $ xxd state.json
    00000000: 7b7d 0a                                  {}.
    $ wc -c state.json
           3
    $ read-state -> "tags":{}

The emptied artifact is EXACTLY `{}` plus a newline, three bytes. No
residue: no header, no version marker, no empty table, no tombstone, no
retained key names. `load` on it yields `flowbind.store{}` - the same
zero-key store `load` yields for an absent file. The file still EXISTS
(clear does not unlink), but its CONTENT carries nothing that distinguishes
it.

### (v) diff emptied-store vs absent-file -> `a5-step-v-diff.txt`

    $ diff <(sed 's@.../state.json@ARTIFACT@'  read-emptied.json) \
           <(sed 's@.../absent.json@ARTIFACT@' read-absent.json)
    diff-exit=0

Modulo the artifact path echoed in the envelope, the two `read-state`
results are BYTE-IDENTICAL. Same reader id, same requested key set, same
`"tags":{}`.

Writer path, same question:

    $ set-state --write status=draft  over the EMPTIED artifact  -> {"status":"draft"}
    $ set-state --write status=draft  over a FRESH  artifact     -> {"status":"draft"}
    $ cmp emptied-copy.json fresh2.json
    BYTE-IDENTICAL

So the equivalence holds for every reader AND every writer path over this
store: writing into an emptied store produces the same artifact as writing
into a never-written one.

### direct load / Apply confirmation -> `a5-step-go-load-apply.txt`

Source: `a5-load-spike_test.go.txt`, run as an in-package test in
`internal/cli/flowbind/` and then removed - no repo file was left changed.

    (i)  load(absent)  = flowbind.store{}  len=0
    (ii)   bytes={"stale":"keepme","status":"final"}  len=36 keys=2
    (iii)  bytes={"status":"final"}                   len=19 keys=1
    (iv)   bytes={}                                   len=3  keys=0
    (v)  absent==emptied: flowbind.store{} == flowbind.store{} -> SAME ZERO-KEY STORE
    (iii-check) one-key store len=1 -> empty-store predicate DECLINES
    --- PASS

`reflect.DeepEqual(load(absent), load(emptied))` is TRUE.

## Verdict

VERIFIES, in the affirming direction. An emptied store and an absent store
are the same zero-key `store` to `load`, and a store retaining at least one
key is observably non-empty. The RDR's stated exception is exact as
written: clearing the LAST key empties the store, and a later `init-state`
guarded by an empty-store predicate would reseed it. Unconditional
automated re-invocation is safe for as long as at least one key remains.

## One boundary caveat the RDR should state -> `a5-step-boundary-sealed.txt`

The predicate's boundary sits at "the store carries zero KEYS", NOT at "the
author's owned keys are all cleared", and those differ in exactly one
shipped case: the read-back seal.

When a `write` accessor declares an unreachable read-back locator, `Apply`
writes `sealedKey` (the NUL-prefixed `flow.readback-unreachable`) into the
same store. Clearing every owned key then leaves a ONE-KEY artifact
carrying only that marker - 39 bytes, the NUL JSON-escaped; see the raw
capture for the exact bytes. `load` reports `keys=1`: a NON-empty store
with zero author-visible keys. A verb whose predicate is `len(store) == 0`
reads this as "already initialized" and declines to seed, even though no
owned key is present.

This does not falsify A5. Such an artifact is already exit-3 unreadable to
every reader - `flow read-state` over it refuses `flow-read-incomplete`,
exit 3 - so it is not a state any correct flow rests in, and the seal is
cleared by the next write with a reachable locator. But the RDR should say
the predicate is over the STORE's key count, and that a sealed artifact is
not an empty store, rather than implying the boundary is "no owned keys
remain".
