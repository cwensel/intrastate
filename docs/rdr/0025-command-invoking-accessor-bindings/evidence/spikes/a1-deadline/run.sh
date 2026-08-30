#!/bin/sh
# Build and run the A1 deadline spike; tee to output.txt.
set -eu
cd "$(dirname "$0")"
{
  echo "--- go vet"
  go vet ./... && echo "vet: ok"
  echo "--- build"
  go build -o a1spike . && echo "build: ok"
  echo "--- run"
  ./a1spike
  echo "--- host check (pkill -f A1SPIKE_ as belt-and-braces)"
  pkill -f 'A1SPIKE_' && echo "pkill: killed leftovers" || echo "pkill: nothing to kill"
} 2>&1 | tee output.txt
rm -f a1spike
