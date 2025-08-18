package report

import "ai-qa-mvp/orchestrator/internal/types"

type SecurityAuditPrompt struct{}

func (SecurityAuditPrompt) SystemPrompt() string {
	return "You are a security analyst. Produce a structured vulnerability report from the provided test run."
}

func (SecurityAuditPrompt) UserPrompt(run *types.Run) string {
	return formatRunPrompt("Security Audit", run, `
- Summary of potential vulnerabilities
- Risk assessment with severity
- Exploitation paths (if apparent)
- Recommendations for mitigation
`)
}
