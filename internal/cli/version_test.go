package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
)

// runCmd builds a fresh root tree, captures stdout/stderr, and runs it
// against args. It is the harness future verb tests copy: drive
// ExecuteAndEmit (the production emission path), then assert on the
// captured streams and the returned error's structured code.
func runCmd(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	cmd := NewRootCmd()
	var out, errBuf bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errBuf)
	err = ExecuteAndEmit(cmd, args)
	return out.String(), errBuf.String(), err
}

func TestVersion_Text(t *testing.T) {
	stdout, _, err := runCmd(t, "version")
	if err != nil {
		t.Fatalf("version: unexpected error: %v", err)
	}
	if !strings.Contains(stdout, "dev") {
		t.Errorf("stdout = %q; want it to contain the dev version", stdout)
	}
}

func TestVersion_JSON(t *testing.T) {
	stdout, _, err := runCmd(t, "version", "--as=json")
	if err != nil {
		t.Fatalf("version --as=json: unexpected error: %v", err)
	}
	var env struct {
		Type string `json:"type"`
		Data struct {
			Version string `json:"version"`
		} `json:"data"`
	}
	if jerr := json.Unmarshal([]byte(stdout), &env); jerr != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", jerr, stdout)
	}
	if env.Type != "ok" {
		t.Errorf("type = %q; want %q", env.Type, "ok")
	}
	if env.Data.Version == "" {
		t.Errorf("data.version is empty")
	}
}

func TestUnknownCommand_IsStructured(t *testing.T) {
	_, _, err := runCmd(t, "bogus")
	if err == nil {
		t.Fatal("expected an error for an unknown command")
	}
	if code := clierr.ErrorCode(err); code != "command-error" {
		t.Errorf("code = %q; want %q", code, "command-error")
	}
	if got := clierr.ExitCodeFor(err); got != 2 {
		t.Errorf("exit code = %d; want 2", got)
	}
}

func TestInvalidMode_IsRefused(t *testing.T) {
	_, _, err := runCmd(t, "version", "--as=yaml")
	if err == nil {
		t.Fatal("expected an error for --as=yaml")
	}
	if code := clierr.ErrorCode(err); code != "flag-invalid-value" {
		t.Errorf("code = %q; want %q", code, "flag-invalid-value")
	}
}
