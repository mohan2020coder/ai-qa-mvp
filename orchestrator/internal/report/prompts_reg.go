package report

import "ai-qa-mvp/orchestrator/internal/types"

type RegressionSummaryPrompt struct{}

func (RegressionSummaryPrompt) SystemPrompt() string {
	return "You are a QA lead. Summarize regression test results clearly for product managers."
}

func (RegressionSummaryPrompt) UserPrompt(run *types.Run) string {
	return formatRunPrompt("Regression Summary", run, `
- Overview of passed/failed steps
- Key regressions detected
- Possible root causes
- Recommended next steps
`)
}
