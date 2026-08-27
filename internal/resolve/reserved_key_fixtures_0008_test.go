package resolve_test

// Fixture builders and structural probes for RDR 0008's kernel half.
//
// Nothing here mocks the unit under test. The builders produce plain Input
// values; the probes read the kernel package's own source and exported
// surface, which is the only honest way to assert a scope ceiling ("exactly
// one new exported predicate") and a doc-pointer obligation.

import (
	"go/ast"
	"go/token"
	"reflect"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// tableRecognizedTagKey is the table package's exported spelling of the
// reserved key. REQ-7's latitude permits two constants; REQ-7's test is that
// they agree, so this is compared against kernel BEHAVIOR, never against the
// kernel's own constant.
const tableRecognizedTagKey = table.RecognizedTagKey

// recognizedRequiresOwnedInput is scenario 9's fixture: a row naming the
// reserved key in RequiresOwned, over a view whose `recognized` key is
// present under ProvenanceRecognized.
func recognizedRequiresOwnedInput() resolve.Input {
	in := recognizedMatchInput()
	in.Table.Rows[0].RequiresOwned = []string{reservedKey}
	return in
}

// collidingReservedKeyInput is scenario 8's fixture: a hand-constructed
// non-conforming Input whose owned AND observed tags both carry the reserved
// key, alongside a genuine recognized outcome. All three provenances collide
// on one key, which is what D3's precedence must resolve.
func collidingReservedKeyInput() resolve.Input {
	in := recognizedMatchInput()
	in.Owned = append(in.Owned, resolve.Tag{Key: reservedKey, Value: "owned-wins"})
	in.Observed = append(in.Observed, resolve.Tag{Key: reservedKey, Value: "observed-loses"})
	return in
}

// reservedBinding is what the assembled view holds at the reserved key.
type reservedBinding struct {
	value      string
	provenance resolve.Provenance
	present    bool
}

// assembledReservedBinding reads the reserved key out of the assembled view
// on the package-internal path, bypassing the entry precondition.
func assembledReservedBinding(t *testing.T, in resolve.Input) reservedBinding {
	t.Helper()
	v, p, ok := resolve.AssembledBindingForTest(in, reservedKey)
	if !ok {
		t.Fatalf("the assembled view carries no %q key", reservedKey)
	}
	return reservedBinding{value: v, provenance: p, present: ok}
}

// resolveBypassingPrecondition runs the kernel pipeline from `assemble`,
// skipping the entry precondition, and reports the refusal it reached.
func resolveBypassingPrecondition(t *testing.T, in resolve.Input) (resolve.RefusalKind, []string) {
	t.Helper()
	got := resolve.ResolveBypassingPreconditionForTest(in)
	if got.Refusal == nil {
		t.Fatalf("the internal path yielded no refusal: %+v", got)
	}
	return got.Refusal.Kind, got.Refusal.MissingOwned
}

// --- structural probes ---------------------------------------------------

// constDoc returns the doc comment on a package-level constant of the kernel
// package, by name.
func constDoc(t *testing.T, name string) string {
	t.Helper()
	for _, f := range parseKernelPackage(t) {
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.CONST {
				continue
			}
			for _, spec := range gd.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, id := range vs.Names {
					if id.Name != name {
						continue
					}
					if vs.Doc != nil {
						return vs.Doc.Text()
					}
					if gd.Doc != nil {
						return gd.Doc.Text()
					}
				}
			}
		}
	}
	return ""
}

// mentionsRDR0008 reports whether a doc comment cites this RDR. The citation
// spelling is latitude; that a reader of the constant is pointed at 0008 is
// not.
func mentionsRDR0008(doc string) bool {
	flat := strings.Join(strings.Fields(doc), " ")
	for _, form := range []string{"RDR 0008", "rdr-0008", "RDR-0008", "0008:C1"} {
		if strings.Contains(flat, form) {
			return true
		}
	}
	return false
}

// exportedInputPredicates names every exported package-level function of the
// kernel whose sole parameter is an Input and whose sole result is an error.
// That is the shape `0008:C4` concedes, and REQ-69 caps the count at one.
func exportedInputPredicates(t *testing.T) []string {
	t.Helper()
	var out []string
	for _, f := range parseKernelPackage(t) {
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Recv != nil || !fd.Name.IsExported() {
				continue
			}
			if !singleParamNamed(fd.Type.Params, "Input") {
				continue
			}
			if !singleResultNamed(fd.Type.Results, "error") {
				continue
			}
			out = append(out, fd.Name.Name)
		}
	}
	return out
}

func singleParamNamed(fl *ast.FieldList, typeName string) bool {
	if fl == nil || len(fl.List) != 1 {
		return false
	}
	if n := len(fl.List[0].Names); n > 1 {
		return false
	}
	id, ok := fl.List[0].Type.(*ast.Ident)
	return ok && id.Name == typeName
}

func singleResultNamed(fl *ast.FieldList, typeName string) bool {
	if fl == nil || len(fl.List) != 1 {
		return false
	}
	if n := len(fl.List[0].Names); n > 1 {
		return false
	}
	id, ok := fl.List[0].Type.(*ast.Ident)
	return ok && id.Name == typeName
}

// resolveSignatureIsResultError reports whether Resolve still takes one Input
// and returns (Result, error). REQ-67 requires the check be added WITHOUT
// widening the signature.
func resolveSignatureIsResultError(t *testing.T) bool {
	t.Helper()
	ft := reflect.TypeOf(resolve.Resolve)
	if ft.Kind() != reflect.Func {
		return false
	}
	if ft.NumIn() != 1 || ft.NumOut() != 2 || ft.IsVariadic() {
		return false
	}
	return ft.In(0) == reflect.TypeOf(resolve.Input{}) &&
		ft.Out(0) == reflect.TypeOf(resolve.Result{}) &&
		ft.Out(1) == reflect.TypeOf((*error)(nil)).Elem()
}
