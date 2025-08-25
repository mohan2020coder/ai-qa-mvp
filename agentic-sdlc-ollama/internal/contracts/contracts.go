package contracts

import "context"

// Agent is a single role (Product, Design, Code, Test, Deploy) that
// consumes an input string (the working doc so far) and returns an output.
type Agent interface {
	Name() string
	Run(ctx context.Context, input string) (string, error)
}

// Orchestrator executes a set of agents in sequence and persists artifacts.
type Orchestrator interface {
	Execute(ctx context.Context, spec string) (Report, error)
}

type PhaseResult struct {
	Agent  string
	Output string
	Path   string // where it was written
}

type Report struct {
	SpecPath string
	Results  []PhaseResult
}
