package guard_test

// Adversarial tests for RDR 0012's Failure Modes. Each plants the failure
// the record names and asserts the detection the record relies on fires.

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// zeroValueCheck is the S4 typed construction check the record's
// "Zero-value evaluator survives migration" failure mode names as the ONLY
// detection for a missed site whose guards use lt/lte/gt/gte or contains.
const zeroValueCheck = "^TestReq48_0012_NoZeroValueEvaluatorSurvivesOutsideTheOneAllowListedFunction$"

// ADV-1 (Failure Modes, "Zero-value evaluator survives migration"; S4):
// every zero-value form of guard.Evaluator Go admits in non-test code must
// fail the typed check. The record's S4 names the composite literal, the
// var, new, and the struct field; a composite literal spelled through a
// type alias or with its type elided is still a composite literal of type
// guard.Evaluator, and a named result is still a var of that type. A check
// that matches the selector `guard.Evaluator` syntactically rather than by
// type misses them, and a missed site then migrates with no signal at all.
//
// Each form is planted, alone, into a fresh copy of the module and the
// check is run there; a pristine copy must pass first, so a failure is the
// plant's and not the venue's.
func TestAdv1_0012_TheZeroValueCheckCatchesEveryZeroValueForm(t *testing.T) {
	if testing.Short() {
		t.Skip("copies the module and runs the S4 check in a subprocess")
	}
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("no go toolchain on PATH")
	}

	const cliHeader = "package cli\n\nimport (\n\t\"github.com/cwensel/intrastate/internal/guard\"\n\t\"github.com/cwensel/intrastate/internal/resolve\"\n)\n\nvar _ resolve.GuardEvaluator = guard.NewEvaluator(nil)\n\n"
	plants := []struct {
		name string
		path string
		src  string
	}{
		{"composite literal (control)", "internal/cli/adv0012_plant_composite.go",
			cliHeader + "func advPlantComposite() resolve.GuardEvaluator { return guard.Evaluator{} }\n"},
		{"var (control)", "internal/cli/adv0012_plant_var.go",
			cliHeader + "func advPlantVar() resolve.GuardEvaluator {\n\tvar ev guard.Evaluator\n\treturn ev\n}\n"},
		{"new (control)", "internal/cli/adv0012_plant_new.go",
			cliHeader + "func advPlantNew() resolve.GuardEvaluator { return *new(guard.Evaluator) }\n"},
		{"embedded struct field (control)", "internal/cli/adv0012_plant_field.go",
			cliHeader + "type advPlantField struct{ guard.Evaluator }\n\nvar _ = advPlantField{}\n"},
		{"in-package composite literal (control)", "internal/guard/adv0012_plant_inpkg.go",
			"package guard\n\nfunc advPlantInPackage() Evaluator { return Evaluator{} }\n\nvar _ = advPlantInPackage\n"},
		{"composite literal through a type alias", "internal/cli/adv0012_plant_alias.go",
			cliHeader + "type advPlantAlias = guard.Evaluator\n\nfunc advPlantAliased() resolve.GuardEvaluator { return advPlantAlias{} }\n"},
		{"elided composite literal", "internal/cli/adv0012_plant_elided.go",
			cliHeader + "func advPlantElided() resolve.GuardEvaluator { return []guard.Evaluator{{}}[0] }\n"},
		{"named result", "internal/cli/adv0012_plant_named.go",
			cliHeader + "func advPlantNamed() (ev guard.Evaluator) { return }\n"},
		{"make-allocated element", "internal/cli/adv0012_plant_make.go",
			cliHeader + "func advPlantMake() resolve.GuardEvaluator { return make([]guard.Evaluator, 1)[0] }\n"},
		{"composite literal through a dot import", "internal/cli/adv0012_plant_dot.go",
			"package cli\n\nimport (\n\t. \"github.com/cwensel/intrastate/internal/guard\"\n\t\"github.com/cwensel/intrastate/internal/resolve\"\n)\n\nfunc advPlantDot() resolve.GuardEvaluator { return Evaluator{} }\n"},
		{"in-package composite literal through a type alias", "internal/guard/adv0012_plant_inpkg_alias.go",
			"package guard\n\ntype advPlantAlias = Evaluator\n\nfunc advPlantInPackageAliased() Evaluator { return advPlantAlias{} }\n\nvar _ = advPlantInPackageAliased\n"},
		{"zero-value array var", "internal/cli/adv0012_plant_array.go",
			cliHeader + "func advPlantArray() resolve.GuardEvaluator {\n\tvar evs [1]guard.Evaluator\n\treturn evs[0]\n}\n"},
		{"zero-filled array literal element", "internal/cli/adv0012_plant_arraylit.go",
			cliHeader + "func advPlantArrayLit() resolve.GuardEvaluator { return [1]guard.Evaluator{}[0] }\n"},
		{"partially filled array literal element", "internal/cli/adv0012_plant_arraylit_partial.go",
			cliHeader + "func advPlantArrayLitPartial() resolve.GuardEvaluator {\n\treturn [2]guard.Evaluator{guard.NewEvaluator(nil)}[1]\n}\n"},
	}

	root := moduleRoot(t)
	venue := t.TempDir()
	copyModule(t, root, venue)

	run := func() (string, error) {
		cmd := exec.Command(goBin, "test", "-count=1", "-run", zeroValueCheck, "./internal/guard/")
		cmd.Dir = venue
		cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=")
		out, err := cmd.CombinedOutput()
		return string(out), err
	}

	if out, err := run(); err != nil {
		t.Fatalf("the S4 check fails on a pristine copy of the module, so the venue proves nothing:\n%s", out)
	}

	for _, p := range plants {
		t.Run(p.name, func(t *testing.T) {
			path := filepath.Join(venue, filepath.FromSlash(p.path))
			if err := os.WriteFile(path, []byte(p.src), 0o644); err != nil {
				t.Fatal(err)
			}
			defer os.Remove(path)

			out, err := run()
			if strings.Contains(out, "[build failed]") || strings.Contains(out, "[setup failed]") {
				t.Fatalf("the plant %s does not compile, so it proves nothing:\n%s", p.path, out)
			}
			if err == nil {
				t.Errorf("the S4 check PASSED with a zero-value guard.Evaluator planted (%s) in %s; "+
					"a missed construction site of this form migrates silently", p.name, p.path)
				return
			}
			if !strings.Contains(out, filepath.Base(p.path)) {
				t.Errorf("the S4 check failed but did not name the planted site %s:\n%s", p.path, out)
			}
		})
	}
}

// moduleRoot walks up from the package directory to the go.mod.
func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("no go.mod above the package directory")
		}
		dir = parent
	}
}

// copyModule copies the module's Go sources (go.mod, go.sum, cmd/,
// internal/) into dst.
func copyModule(t *testing.T, src, dst string) {
	t.Helper()
	for _, f := range []string{"go.mod", "go.sum"} {
		b, err := os.ReadFile(filepath.Join(src, f))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dst, f), b, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, top := range []string{"cmd", "internal"} {
		err := filepath.WalkDir(filepath.Join(src, top), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(src, path)
			if err != nil {
				return err
			}
			target := filepath.Join(dst, rel)
			if d.IsDir() {
				return os.MkdirAll(target, 0o755)
			}
			if !d.Type().IsRegular() {
				return nil
			}
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			return os.WriteFile(target, b, 0o644)
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}
