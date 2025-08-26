package contracts

type Workflow struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Spec   string   `json:"spec"`
	Agents []string `json:"agents"`
}

type PhaseResult struct {
	Agent  string `json:"agent"`
	Output string `json:"output"`
	Path   string `json:"path"`
}

type Report struct {
	SpecPath string        `json:"spec_path"`
	Results  []PhaseResult `json:"results"`
}

type Execution struct {
	ID         string `json:"id"`
	WorkflowID string `json:"workflow_id"`
	Status     string `json:"status"`
}
