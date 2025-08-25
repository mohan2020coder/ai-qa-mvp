package orchestrator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"agentic-sdlc-ollama/internal/contracts"
)

type simpleOrchestrator struct {
	agents  []contracts.Agent
	outRoot string
}

func New(agents []contracts.Agent, outRoot string) contracts.Orchestrator {
	return &simpleOrchestrator{agents: agents, outRoot: outRoot}
}

func (o *simpleOrchestrator) Execute(ctx context.Context, spec string) (contracts.Report, error) {
	if o.outRoot == "" {
		o.outRoot = ".workspace"
	}
	outDir := filepath.Join(o.outRoot, "outputs")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return contracts.Report{}, err
	}

	combined := "# Spec\n\n" + spec + "\n\n"
	report := contracts.Report{Results: []contracts.PhaseResult{}}
	var lastOut string

	for i, a := range o.agents {
		select {
		case <-ctx.Done():
			return report, ctx.Err()
		default:
		}

		// 👇 Instead of the full "combined", give spec + last output only
		input := "# Spec\n\n" + spec
		if lastOut != "" {
			input += "\n\n# Previous Output\n\n" + lastOut
		}

		start := time.Now()
		out, err := a.Run(ctx, input)
		if err != nil {
			out = fmt.Sprintf("_Agent %s failed: %v_", a.Name(), err)
		}

		filename := fmt.Sprintf("%02d-%s.md", i+1, a.Name())
		path := filepath.Join(outDir, filename)
		if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
			return report, err
		}

		report.Results = append(report.Results, contracts.PhaseResult{
			Agent:  a.Name(),
			Output: out,
			Path:   filepath.ToSlash(path),
		})

		// build final combined doc for user visibility
		combined += fmt.Sprintf(
			"\n\n# %s Output (generated in %s)\n\n%s\n",
			a.Name(),
			time.Since(start).Round(time.Millisecond),
			out,
		)

		lastOut = out
	}

	combinedPath := filepath.Join(outDir, "combined.md")
	if err := os.WriteFile(combinedPath, []byte(combined), 0o644); err != nil {
		return report, err
	}
	return report, nil
}
