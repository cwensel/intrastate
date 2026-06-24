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

type accessorRef struct {
	name       string
	capability capability
}

type accessorDefinition struct {
	accessor
	readBackRequired bool
	ambientDiscovery bool
	writes           []string
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

func validateDefinitions(refs []accessorRef, defs []accessorDefinition, ownedTags map[string]bool) []string {
	type key struct {
		name       string
		capability capability
	}

	bindings := make(map[key][]accessorDefinition, len(defs))
	for _, def := range defs {
		if def.timeout <= 0 {
			return []string{"missing_or_non_positive_timeout"}
		}
		if def.ambientDiscovery {
			return []string{"ambient_artifact_discovery"}
		}
		if def.capability == writeCap && !def.readBackRequired {
			return []string{"missing_write_read_back"}
		}
		for _, tag := range def.writes {
			if !ownedTags[tag] {
				return []string{"write_non_owned_tag"}
			}
		}
		k := key{name: def.name, capability: def.capability}
		bindings[k] = append(bindings[k], def)
	}

	var failures []string
	for _, ref := range refs {
		matched := bindings[key(ref)]
		if len(matched) == 1 {
			continue
		}
		if len(matched) == 0 {
			for k := range bindings {
				if k.name == ref.name {
					failures = append(failures, "capability_mismatch")
					goto nextRef
				}
			}
			failures = append(failures, "missing_accessor")
			continue
		}
		failures = append(failures, "multiply_bound_accessor")
	nextRef:
	}
	return failures
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

type writeCorruption string

const (
	corruptNone     writeCorruption = ""
	corruptOwned    writeCorruption = "owned"
	corruptObserved writeCorruption = "observed"
)

func write(r registry, artifacts map[string]*artifact, planned map[string]string, corrupt writeCorruption) result {
	name := "state.persist"
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
		before := clone(art.tags)
		maps.Copy(art.tags, planned)
		switch corrupt {
		case corruptOwned:
			art.tags["status"] = "corrupt"
		case corruptObserved:
			art.tags["profile"] = "small"
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
		for k, beforeValue := range before {
			if _, owned := planned[k]; owned {
				continue
			}
			if observed[k] != beforeValue {
				return &result{
					name:     name,
					cap:      writeCap,
					refusal:  refusalReadBackMismatch,
					expected: before,
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
		write(reg, artifacts, map[string]string{"status": "Final"}, corruptNone),
		write(reg, newArtifacts(), map[string]string{"status": "Final"}, corruptOwned),
		write(reg, newArtifacts(), map[string]string{"status": "Final"}, corruptObserved),
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

	for _, c := range validationCases() {
		failures := validateDefinitions(c.refs, c.defs, map[string]bool{"status": true})
		fmt.Printf("validation-%s=%s\n", c.name, strings.Join(failures, ","))
		if !reflect.DeepEqual(failures, c.want) {
			panic(fmt.Errorf("%s: got %v want %v", c.name, failures, c.want))
		}
	}
}

type validationCase struct {
	name string
	refs []accessorRef
	defs []accessorDefinition
	want []string
}

func validationCases() []validationCase {
	baseRefs := []accessorRef{
		{name: "state.read", capability: readCap},
		{name: "state.gate", capability: gateCap},
		{name: "state.persist", capability: writeCap},
	}
	baseDefs := []accessorDefinition{
		{accessor: accessor{name: "state.read", capability: readCap, role: "state", timeout: 10 * time.Millisecond}},
		{accessor: accessor{name: "state.gate", capability: gateCap, role: "state", timeout: 10 * time.Millisecond}},
		{
			accessor:         accessor{name: "state.persist", capability: writeCap, role: "state", timeout: 10 * time.Millisecond},
			readBackRequired: true,
			writes:           []string{"status"},
		},
	}

	return []validationCase{
		{name: "ok", refs: baseRefs, defs: baseDefs},
		{
			name: "missing",
			refs: append(baseRefs, accessorRef{name: "missing.read", capability: readCap}),
			defs: baseDefs,
			want: []string{"missing_accessor"},
		},
		{
			name: "multiply-bound",
			refs: baseRefs,
			defs: append(baseDefs, baseDefs[0]),
			want: []string{"multiply_bound_accessor"},
		},
		{
			name: "capability-mismatch",
			refs: append(baseRefs, accessorRef{name: "state.persist", capability: readCap}),
			defs: baseDefs,
			want: []string{"capability_mismatch"},
		},
		{
			name: "missing-timeout",
			refs: baseRefs,
			defs: []accessorDefinition{
				{accessor: accessor{name: "state.read", capability: readCap, role: "state"}},
			},
			want: []string{"missing_or_non_positive_timeout"},
		},
		{
			name: "missing-readback",
			refs: baseRefs,
			defs: []accessorDefinition{
				{accessor: accessor{name: "state.persist", capability: writeCap, role: "state", timeout: 10 * time.Millisecond}},
			},
			want: []string{"missing_write_read_back"},
		},
		{
			name: "ambient-discovery",
			refs: baseRefs,
			defs: []accessorDefinition{
				{
					accessor:         accessor{name: "state.persist", capability: writeCap, timeout: 10 * time.Millisecond},
					readBackRequired: true,
					ambientDiscovery: true,
					writes:           []string{"status"},
				},
			},
			want: []string{"ambient_artifact_discovery"},
		},
		{
			name: "non-owned-write",
			refs: baseRefs,
			defs: []accessorDefinition{
				{
					accessor:         accessor{name: "state.persist", capability: writeCap, role: "state", timeout: 10 * time.Millisecond},
					readBackRequired: true,
					writes:           []string{"observed"},
				},
			},
			want: []string{"write_non_owned_tag"},
		},
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
	return write(reg, artifacts, map[string]string{"status": "Final"}, corruptNone)
}
