#!/bin/sh
# Build and run the A3/A7/A8 spike harness; tee to output.txt.
# Host paths are normalised ($SPIKE, $TMPDIR) so the transcript is portable.
set -u
cd "$(dirname "$0")"
SPIKE="$(pwd)"
T="${TMPDIR:-/tmp}"; T="${T%/}"
{
  echo "spike: rdr-0025 a3-a7-a8-real-tools"
  echo "date: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "go: $(go version)"
  echo "git: $(git --version)"
  echo "os: $(uname -sr)"
  echo "---"
  go build -o ./a3spike . && ./a3spike
  echo "harness_exit=$?"
} 2>&1 | sed -e "s#$SPIKE#\$SPIKE#g" -e "s#$T#\$TMPDIR#g" | tee output.txt
rm -f ./a3spike
