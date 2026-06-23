#!/bin/sh
set -eu

fixture=${1:-guard-fixture.toml}

unknown=$(
  awk '
    /^[[]rule[.](match|guard[.]all|guard[.]unless)[.][^]]+[]]/ {
      section=$0
      next
    }
    /^[[]/ {
      section=""
      next
    }
    section != "" && /^[[:space:]]*[A-Za-z_][A-Za-z0-9_]*[[:space:]]*=/ {
      op=$1
      sub(/[[:space:]]*=.*/, "", op)
      if (op !~ /^(eq|in|lt|lte|gt|gte|exists|contains)$/) {
        print section " " op
      }
    }
  ' "$fixture"
)

if [ -n "$unknown" ]; then
  printf 'UNKNOWN OPERATORS\n%s\n' "$unknown"
  exit 1
fi

printf 'fixture=%s\n' "$fixture"
printf 'operators='
awk '
  /^[[]rule[.](match|guard[.]all|guard[.]unless)[.][^]]+[]]/ {
    section=$0
    next
  }
  /^[[]/ {
    section=""
    next
  }
  section != "" && /^[[:space:]]*[A-Za-z_][A-Za-z0-9_]*[[:space:]]*=/ {
    op=$1
    sub(/[[:space:]]*=.*/, "", op)
    seen[op]=1
  }
  END {
    first=1
    split("eq in lt gte exists", order, " ")
    for (i=1; i<=length(order); i++) {
      op=order[i]
      if (seen[op]) {
        if (!first) printf ","
        printf "%s", op
        first=0
      }
    }
    printf "\n"
  }
' "$fixture"
printf 'rules='
awk '/^[[][[]rule[]][]]$/ { count++ } END { print count + 0 }' "$fixture"
printf 'coverage=status/profile routing, cap-3 handling, prelock lens sets, cluster eligibility, rewind legality\n'
