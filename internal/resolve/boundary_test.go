package resolve_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"strconv"
	"strings"
	"testing"
)

// This file holds the package-boundary inspection helpers used by the
// negative-contract tests. RDR 0001's prohibitions ("MUST NOT print
// output, inspect CLI flags, … execute persistence side effects") are
// claims about the kernel package's dependency surface, so they are
// checked against the package's own source rather than against runtime
// behaviour that a stub could fake.

const kernelModulePath = "github.com/newcoinc/intrastate"

// kernelImportPath is the import path of the resolver kernel package
// under test.
func kernelImportPath(t *testing.T) string {
	t.Helper()
	return kernelModulePath + "/internal/resolve"
}

// isInternalPackage reports whether path lives under an internal/ tree.
func isInternalPackage(path string) bool {
	return strings.Contains(path, "/internal/") || strings.HasSuffix(path, "/internal")
}

// isCLIPackage reports whether path is part of the CLI tree.
func isCLIPackage(path string) bool {
	return path == kernelModulePath+"/internal/cli" ||
		strings.HasPrefix(path, kernelModulePath+"/internal/cli/")
}

// isThirdParty reports whether path is neither a standard-library package
// nor a package of this module.
func isThirdParty(path string) bool {
	if strings.HasPrefix(path, kernelModulePath) {
		return false
	}
	first := path
	if i := strings.Index(path, "/"); i >= 0 {
		first = path[:i]
	}
	// Standard-library first path elements never contain a dot.
	return strings.Contains(first, ".")
}

// parseKernelPackage parses the non-test sources of the kernel package.
func parseKernelPackage(t *testing.T) map[string]*ast.File {
	t.Helper()
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", func(fi fs.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, parser.ParseComments)
	if err != nil {
		t.Fatalf("parsing the resolver kernel package: %v", err)
	}
	for name, pkg := range pkgs {
		if strings.HasSuffix(name, "_test") {
			continue
		}
		return pkg.Files
	}
	t.Fatal("no non-test source found for the resolver kernel package")
	return nil
}

// kernelPackageImports returns the set of packages the kernel imports.
func kernelPackageImports(t *testing.T) map[string]bool {
	t.Helper()
	out := map[string]bool{}
	for _, f := range parseKernelPackage(t) {
		for _, imp := range f.Imports {
			p, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				t.Fatalf("unquoting import %s: %v", imp.Path.Value, err)
			}
			out[p] = true
		}
	}
	return out
}

// exportedKernelSymbols returns every exported top-level name the kernel
// package declares: types, funcs, methods, consts, and vars.
func exportedKernelSymbols(t *testing.T) []string {
	t.Helper()
	var names []string
	add := func(n string) {
		if n != "" && ast.IsExported(n) {
			names = append(names, n)
		}
	}
	for _, f := range parseKernelPackage(t) {
		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				add(d.Name.Name)
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						add(s.Name.Name)
					case *ast.ValueSpec:
						for _, n := range s.Names {
							add(n.Name)
						}
					}
				}
			}
		}
	}
	return names
}
