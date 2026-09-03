#!/bin/sh
# Build and run the same two spike programs inside the linux container.
set -e
cd /w
mkdir -p /tmp/b && cp fixture.go /tmp/b/fixture_src.go && cp spike.go /tmp/b/spike_src.go
cd /tmp/b
export GOFLAGS=-mod=mod GOCACHE=/tmp/gocache
go build -o /tmp/b/fixture fixture_src.go
go build -o /tmp/b/spike   spike_src.go
exec /tmp/b/spike /tmp/b/fixture
