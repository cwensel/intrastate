# Flight drain report — batch:rdr-0030

Invocation: `/kata-flight --label batch:rdr-0030 --drain` (review gate on)

```
flight: 0 shipped, 0 stopped, 0 skipped over 1 waves
  shipped: (none)
  stopped: (none)
  skipped: (none)
```

Wave 1 resolved `[c8ge]` (1 kata, 0 dropped by eligibility).

```
scope-review: c8ge  reviewed: 1
  in-scope:   (none)
  rdr-shaped: c8ge (-> kind:rdr-seed; batch:rdr-0030 stripped)
```

- c8ge: RDR-SHAPED. A match `eq` pinning a key that the same group's guard also
  gates yields coverage gaps on claimed cells. The code conforms to 0003:C13 /
  0006:C7 (a guard dimension spans its full declared domain; match keys contribute
  no assignment). Narrowing that dimension to the pinned member is the general
  form of option (b), which deviation D13 rejected because it re-decides
  C13/C7. It is a design fork for the record owner, not a flight-shippable fix.

After review, 0 eligible survivors remained. No spin-offs were minted, so no
re-sweep was needed and the drain ended after wave 1.
