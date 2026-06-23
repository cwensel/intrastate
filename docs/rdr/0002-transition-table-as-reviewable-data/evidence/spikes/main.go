package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type Model struct {
	Model struct {
		ID          string
		Version     int
		Description string
	}
	Tags      map[string]Tag
	Outcomes  []string
	Accessors map[string]Accessor
	Context   map[string]Context
	Rule      []Rule
	Dump      Dump
}

type Tag struct {
	Provenance string
	Kind       string
	Accessor   string
}

type Accessor struct {
	Mode string
	Path string
}

type Context struct {
	Inherits string
	Match    map[string]map[string]any
}

type Rule struct {
	ID     string
	Use    []string
	Match  map[string]map[string]any
	Guard  Guard
	Write  map[string]any
	Clear  []string
	Self   bool
	Source string
}

type Guard struct {
	All    map[string]map[string]any
	Unless map[string]map[string]any
}

type Dump struct {
	Order []string
}

type Row struct {
	Model  string
	Rule   string
	Source string
	Match  []string
	Write  []string
}

func main() {
	for _, path := range os.Args[1:] {
		data, err := os.ReadFile(path)
		must(err)
		var model Model
		must(toml.Unmarshal(data, &model))
		rows, err := normalize(model)
		must(err)
		fmt.Printf("MODEL %s rows=%d outcomes=%s\n", model.Model.ID, len(rows), strings.Join(model.Outcomes, ","))
		for _, row := range rows {
			fmt.Printf("%s.%s source=%s match=[%s] write=[%s]\n",
				row.Model, row.Rule, row.Source, strings.Join(row.Match, "; "), strings.Join(row.Write, "; "))
		}
	}
}

func normalize(model Model) ([]Row, error) {
	var rows []Row
	for _, rule := range model.Rule {
		predicates := map[string]string{}
		for _, contextID := range rule.Use {
			if err := mergeContext(model.Context, contextID, predicates, map[string]bool{}); err != nil {
				return nil, err
			}
		}
		mergePredicates(predicates, "", rule.Match)
		mergePredicates(predicates, "all:", rule.Guard.All)
		mergePredicates(predicates, "unless:", rule.Guard.Unless)

		writes := make([]string, 0, len(rule.Write)+len(rule.Clear))
		for tag, value := range rule.Write {
			writes = append(writes, fmt.Sprintf("%s=%v", tag, value))
		}
		for _, tag := range rule.Clear {
			writes = append(writes, fmt.Sprintf("%s=<clear>", tag))
		}
		sort.Strings(writes)

		match := make([]string, 0, len(predicates))
		for key, value := range predicates {
			match = append(match, fmt.Sprintf("%s%s", key, value))
		}
		sort.Strings(match)
		rows = append(rows, Row{
			Model:  model.Model.ID,
			Rule:   rule.ID,
			Source: rule.Source,
			Match:  match,
			Write:  writes,
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Model != rows[j].Model {
			return rows[i].Model < rows[j].Model
		}
		return rows[i].Rule < rows[j].Rule
	})
	return rows, nil
}

func mergeContext(contexts map[string]Context, id string, out map[string]string, seen map[string]bool) error {
	ctx, ok := contexts[id]
	if !ok {
		return fmt.Errorf("unknown context %q", id)
	}
	if seen[id] {
		return fmt.Errorf("cyclic context inheritance at %q", id)
	}
	seen[id] = true
	if ctx.Inherits != "" {
		if err := mergeContext(contexts, ctx.Inherits, out, seen); err != nil {
			return err
		}
	}
	mergePredicates(out, "", ctx.Match)
	return nil
}

func mergePredicates(out map[string]string, prefix string, block map[string]map[string]any) {
	for tag, ops := range block {
		for op, value := range ops {
			out[prefix+tag+"."+op+"="] = fmt.Sprint(value)
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
