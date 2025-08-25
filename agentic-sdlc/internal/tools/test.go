package tools

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
)

type TestTool struct {
	root string
}

func NewTestTool(root string) *TestTool { return &TestTool{root: root} }

func (t *TestTool) Name() string        { return "test" }
func (t *TestTool) Description() string { return "Run go tests and capture a simple report" }

func (t *TestTool) Call(ctx context.Context, input string, args map[string]any) (CallResult, error) {
	switch input {
	case "go_test":
		cmd := exec.CommandContext(ctx, "go", "test", "./...", "-count=1")
		cmd.Dir = t.root
		out, err := cmd.CombinedOutput()
		report := map[string]any{
			"cmd":    cmd.String(),
			"output": string(out),
			"passed": err == nil,
		}
		b, _ := json.MarshalIndent(report, "", "  ")
		rep := filepath.ToSlash(filepath.Join("reports", "test.json"))
		full := filepath.Join(t.root, rep)
		_ = os.MkdirAll(filepath.Dir(full), 0o755)
		_ = os.WriteFile(full, b, 0o644)
		return ok("tests run", map[string]any{"passed": err == nil}, []string{rep})
	default:
		return fail("unsupported op: " + input)
	}
}
