package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type DeployTool struct {
	artifacts string
}

func NewDeployTool(artifacts string) *DeployTool {
	_ = os.MkdirAll(artifacts, 0o755)
	return &DeployTool{artifacts: artifacts}
}

func (t *DeployTool) Name() string        { return "deploy" }
func (t *DeployTool) Description() string { return "Simulated deploy to staging environment" }

func (t *DeployTool) Call(ctx context.Context, input string, args map[string]any) (CallResult, error) {
	_ = ctx
	switch input {
	case "deploy_staging":
		image := fmt.Sprintf("%v", args["image"])
		rel := filepath.ToSlash(filepath.Join("deploy", "staging.txt"))
		full := filepath.Join(t.artifacts, "staging.txt")
		msg := "staging deploy ok with image: " + image + "\n"
		_ = os.WriteFile(full, []byte(msg), 0o644)
		return ok("deployed", map[string]any{"env": "staging", "image": image}, []string{rel})
	default:
		return fail("unsupported op: " + input)
	}
}
