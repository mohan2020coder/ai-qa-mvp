package tools

import (
	"context"
	"os"
	"path/filepath"
)

type RepoTool struct {
	root string
}

func NewRepoTool(root string) *RepoTool {
	_ = os.MkdirAll(root, 0o755)
	return &RepoTool{root: root}
}

func (t *RepoTool) Name() string        { return "repo" }
func (t *RepoTool) Description() string { return "Read/write files in the workspace repo" }

func (t *RepoTool) Call(ctx context.Context, input string, args map[string]any) (CallResult, error) {
	_ = ctx
	switch input {
	case "write_file":
		path := args["path"].(string)
		content := args["content"].(string)
		full := filepath.Join(t.root, filepath.FromSlash(path))
		_ = os.MkdirAll(filepath.Dir(full), 0o755)
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			return fail(err.Error())
		}
		return ok("written", map[string]any{"path": path}, []string{path})
	case "read_file":
		path := args["path"].(string)
		full := filepath.Join(t.root, filepath.FromSlash(path))
		b, err := os.ReadFile(full)
		if err != nil {
			return fail(err.Error())
		}
		return ok("read", map[string]any{"path": path, "content": string(b)}, nil)
	default:
		return fail("unsupported op: " + input)
	}
}
