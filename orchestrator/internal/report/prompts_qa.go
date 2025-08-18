package report

import "ai-qa-mvp/orchestrator/internal/types"

type QAReportPrompt struct{}

func (QAReportPrompt) SystemPrompt() string {
	return "You are a senior QA engineer. Produce a concise, structured bug report from the provided test run."
}

func (QAReportPrompt) UserPrompt(run *types.Run) string {
	return formatRunPrompt("QA Report", run, `
- Executive summary (2–3 sentences)
- Detected issues (with severity and clear reproduction steps if apparent)
- Impact and suggestions
- Attachments (just link paths provided; do not embed images)
`)
}
