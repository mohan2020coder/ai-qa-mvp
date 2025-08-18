package report

import (
	"ai-qa-mvp/orchestrator/internal/types"
	"fmt"
	"os"
)

type PromptBuilder interface {
	SystemPrompt() string
	UserPrompt(run *types.Run) string
}

func NewPromptBuilderFromEnv() PromptBuilder {
	switch os.Getenv("REPORT_TYPE") {
	case "security":
		return SecurityAuditPrompt{}
	case "regression":
		return RegressionSummaryPrompt{}
	default:
		return QAReportPrompt{} // fallback
	}
}

// shared helper for building user prompts
func formatRunPrompt(title string, run *types.Run, instructions string) string {
	return fmt.Sprintf(`Create a %s for run %s.
URL: %s
Issues: %v
Steps: %v
Analysis Summary: %s
Please produce:
%s
`, title, run.ID, run.URL, run.Analysis.Issues, run.Analysis.Steps, run.Analysis.Summary, instructions)
}
