package report

import (
	"ai-qa-mvp/orchestrator/internal/types"
	"fmt"
)

func GenerateReport(run *types.Run) (string, error) {
	cfg, err := LoadLLMConfig()
	if err != nil {
		return "", err
	}

	builder := NewPromptBuilderFromEnv()

	req := chatReq{
		Model: cfg.Model,
		Messages: []chatMessage{
			{"system", builder.SystemPrompt()},
			{"user", builder.UserPrompt(run)},
		},
		Stream: false,
	}

	if cfg.Provider == "ollama" {
		return callOllamaLLM(cfg, req)
	} else {
		return callLLM(cfg, req)
	}
}

func FallbackReport(run *types.Run) string {
	return fmt.Sprintf(`QA Report (Fallback)
Run: %s
URL: %s
Summary:
- Auto-generated fallback report (no LLM key present).
Detected Issues:
%v
Steps:
%v
Attachments:
%v
`, run.ID, run.URL, run.Analysis.Issues, run.Analysis.Steps, run.Analysis.Screenshots)
}

/*
example usage:
report, err := report.GenerateReport(run)
// REPORT_TYPE=qa|security|regression decides which prompt is used
*/
