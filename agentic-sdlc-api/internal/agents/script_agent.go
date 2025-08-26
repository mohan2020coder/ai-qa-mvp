package agents

import (
	"context"
	"fmt"
	"os/exec"
)

type ScriptAgent struct {
	name   string
	config map[string]any
	maxTok int
}

func NewScriptAgent(name string, config map[string]any) *ScriptAgent {
	return &ScriptAgent{name: name, config: config}
}

func (s *ScriptAgent) Name() string {
	return s.name
}

func (s *ScriptAgent) Run(ctx context.Context, input string) (string, error) {
	cmdStr, ok := s.config["command"].(string)
	if !ok || cmdStr == "" {
		return "", fmt.Errorf("missing command in config")
	}

	cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("script error: %v: %s", err, string(output))
	}
	return string(output), nil
}

func (s *ScriptAgent) MaxTokens() int {
	return s.maxTok
}

func (s *ScriptAgent) SetMaxTokens(n int) {
	s.maxTok = n
}
