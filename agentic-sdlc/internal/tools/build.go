package tools

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
)

type BuildTool struct {
	root      string
	artifacts string
}

func NewBuildTool(root, artifacts string) *BuildTool {
	_ = os.MkdirAll(artifacts, 0o755)
	return &BuildTool{root: root, artifacts: artifacts}
}

func (t *BuildTool) Name() string        { return "build" }
func (t *BuildTool) Description() string { return "Build a docker image (simulated if docker not available)" }

func (t *BuildTool) Call(ctx context.Context, input string, args map[string]any) (CallResult, error) {
	switch input {
	case "docker_build":
		name := args["name"].(string)
		tag := args["tag"].(string)
		image := name + ":" + tag

		// Try docker build (ignore errors in demo)
		cmd := exec.CommandContext(ctx, "docker", "build", "-t", image, ".")
		cmd.Dir = t.root
		_ = cmd.Run()

		rel := filepath.ToSlash(filepath.Join("artifacts", "image.txt"))
		full := filepath.Join(t.artifacts, "image.txt")
		_ = os.WriteFile(full, []byte(image+"\n"), 0o644)
		return ok("image built", map[string]any{"image": image}, []string{rel})
	default:
		return fail("unsupported op: " + input)
	}
}
