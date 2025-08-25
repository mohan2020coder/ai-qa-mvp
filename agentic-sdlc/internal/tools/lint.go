package tools

import (
	"context"
	"os"
	"path/filepath"
)

type LintTool struct {
	root string
}

func NewLintTool(root string) *LintTool { return &LintTool{root: root} }

func (t *LintTool) Name() string        { return "lint" }
func (t *LintTool) Description() string { return "Trivial linter placeholder that always succeeds" }

func (t *LintTool) Call(ctx context.Context, input string, args map[string]any) (CallResult, error) {
	_ = ctx; _ = args
	switch input {
	case "lint":
		rep := filepath.ToSlash(filepath.Join("reports", "lint.txt"))
		full := filepath.Join(t.root, rep)
		_ = os.MkdirAll(filepath.Dir(full), 0o755)
		_ = os.WriteFile(full, []byte("lint: ok\n"), 0o644)
		return ok("lint ok", map[string]any{}, []string{rep})
	default:
		return fail("unsupported op: " + input)
	}
}
