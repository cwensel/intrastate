package main

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"sort"
	"strings"
	"time"
)

type capability string

const (
	readCap  capability = "read"
	gateCap  capability = "gate"
	writeCap capability = "write"
)

type refusal string

const (
	refusalCapabilityMismatch refusal = "capability_mismatch"
	refusalExecutionFailure   refusal = "execution_failure"
	refusalGateIndeterminate  refusal = "gate_indeterminate"
	refusalReadBackMismatch   refusal = "read_back_mismatch"
	refusalTimeout            refusal = "timeout"
	refusalUnknownAccessor    refusal = "unknown_accessor"
)

type artifact struct {
	role string
	tags map[string]string
}

type accessor struct {
	name       string
	capability capability
	role       string
	timeout    time.Duration
}

type result struct {
	name     string
	cap      capability
	refusal  refusal
	tags     map[string]string
	gate     string
	expected map[string]string
	observed map[string]string
}

type registry map[string]accessor

func (r registry) lookup(name string, cap capability) (accessor, *result) {
	a, ok := r[name]
	if !ok {
		return accessor{}, &result{name: name, cap: cap, refusal: refusalUnknownAccessor}
	}
	if a.capability != cap {
		return accessor{}, &result{name: name, cap: cap, refusal: refusalCapabilityMismatch}
	}
	if a.timeout <= 0 {
		panic("fixture accessors must declare a timeout")
	}
	return a, nil
}

func withTimeout(a accessor, run func(context.Context) *result) *result {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()
	return run(ctx)
}

func read(r registry, artifacts map[string]*artifact, name string) result {
	a, failure := r.lookup(name, readCap)
	if failure != nil {
		return *failure
	}
	return *withTimeout(a, func(ctx context.Context) *result {
		if name == "slow.read" {
			select {
			case <-time.After(25 * time.Millisecond):
			case <-ctx.Done():
				return &result{name: name, cap: readCap, refusal: refusalTimeout}
			}
		}
		art, ok := artifacts[a.role]
		if !ok {
			return &result{name: name, cap: readCap, refusal: refusalExecutionFailure}
		}
		return &result{name: name, cap: readCap, tags: clone(art.tags)}
	})
}

func gate(r registry, mode string) result {
	name := "state.gate"
	a, failure := r.lookup(name, gateCap)
	if failure != nil {
		return *failure
	}
	return *withTimeout(a, func(ctx context.Context) *result {
		select {
		case <-ctx.Done():
			return &result{name: name, cap: gateCap, refusal: refusalTimeout}
		default:
		}
		switch mode {
		case "allow", "deny":
			return &result{name: name, cap: gateCap, gate: mode}
		case "indeterminate":
			return &result{name: name, cap: gateCap, refusal: refusalGateIndeterminate}
		default:
			return &result{name: name, cap: gateCap, refusal: refusalExecutionFailure}
		}
	})
}

func write(r registry, artifacts map[string]*artifact, name string, planned map[string]string, corrupt bool) result {
	a, failure := r.lookup(name, writeCap)
	if failure != nil {
		return *failure
	}
	return *withTimeout(a, func(ctx context.Context) *result {
		select {
		case <-ctx.Done():
			return &result{name: name, cap: writeCap, refusal: refusalTimeout}
		default:
		}
		art, ok := artifacts[a.role]
		if !ok {
			return &result{name: name, cap: writeCap, refusal: refusalExecutionFailure}
		}
		maps.Copy(art.tags, planned)
		if corrupt {
			art.tags["status"] = "corrupt"
		}
		observed := clone(art.tags)
		for k, v := range planned {
			if observed[k] != v {
				return &result{
					name:     name,
					cap:      writeCap,
					refusal:  refusalReadBackMismatch,
					expected: clone(planned),
					observed: observed,
				}
			}
		}
		return &result{name: name, cap: writeCap, tags: clone(planned)}
	})
}

func clone(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	maps.Copy(out, in)
	return out
}

func newArtifacts() map[string]*artifact {
	return map[string]*artifact{
		"state": {role: "state", tags: map[string]string{"status": "Draft", "profile": "large"}},
	}
}

func stableLine(r result) string {
	parts := []string{fmt.Sprintf("%s/%s", r.name, r.cap)}
	if r.refusal != "" {
		parts = append(parts, "refusal="+string(r.refusal))
	}
	if r.gate != "" {
		parts = append(parts, "gate="+r.gate)
	}
	if len(r.tags) > 0 {
		parts = append(parts, "tags="+formatMap(r.tags))
	}
	if len(r.expected) > 0 {
		parts = append(parts, "expected="+formatMap(r.expected))
	}
	if len(r.observed) > 0 {
		parts = append(parts, "observed="+formatMap(r.observed))
	}
	return strings.Join(parts, " ")
}

func formatMap(in map[string]string) string {
	keys := make([]string, 0, len(in))
	for k := range in {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(in[k])
	}
	b.WriteByte('}')
	return b.String()
}

func main() {
	reg := registry{
		"state.read":    {name: "state.read", capability: readCap, role: "state", timeout: 10 * time.Millisecond},
		"state.gate":    {name: "state.gate", capability: gateCap, role: "state", timeout: 10 * time.Millisecond},
		"state.persist": {name: "state.persist", capability: writeCap, role: "state", timeout: 10 * time.Millisecond},
		"slow.read":     {name: "slow.read", capability: readCap, role: "state", timeout: 1 * time.Millisecond},
	}

	artifacts := newArtifacts()
	cases := []result{
		read(reg, artifacts, "state.read"),
		gate(reg, "allow"),
		gate(reg, "indeterminate"),
		gate(reg, "network-error"),
		read(reg, artifacts, "slow.read"),
		read(reg, artifacts, "missing.read"),
		read(reg, artifacts, "state.persist"),
		write(reg, artifacts, "state.persist", map[string]string{"status": "Final"}, false),
		write(reg, newArtifacts(), "state.persist", map[string]string{"status": "Final"}, true),
	}
	for _, c := range cases {
		fmt.Println(stableLine(c))
	}

	first := replay(reg, false)
	second := replay(reg, false)
	injected := replay(reg, true)
	fmt.Printf("replay-identical=%v first=%s second=%s\n", reflect.DeepEqual(first, second), stableLine(first), stableLine(second))
	fmt.Printf("replay-injected=%s\n", stableLine(injected))

	if !reflect.DeepEqual(first, second) {
		panic(errors.New("expected identical replay disposition"))
	}
	if injected.refusal != refusalGateIndeterminate {
		panic(errors.New("expected stable injected refusal"))
	}
}

func replay(reg registry, injectRefusal bool) result {
	artifacts := newArtifacts()
	readResult := read(reg, artifacts, "state.read")
	if readResult.refusal != "" {
		return readResult
	}
	mode := "allow"
	if injectRefusal {
		mode = "indeterminate"
	}
	gateResult := gate(reg, mode)
	if gateResult.refusal != "" {
		return gateResult
	}
	if gateResult.gate != "allow" {
		return gateResult
	}
	return write(reg, artifacts, "state.persist", map[string]string{"status": "Final"}, false)
}
