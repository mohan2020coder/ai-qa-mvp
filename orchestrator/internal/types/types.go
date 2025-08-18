package types

import (
	"time"

	"github.com/google/uuid"
)

type Run struct {
	ID          string     `json:"id" bson:"_id"`
	URL         string     `json:"url" bson:"url"`
	Status      string     `json:"status" bson:"status"`
	Error       string     `json:"error,omitempty" bson:"error,omitempty"`
	CreatedAt   time.Time  `json:"createdAt" bson:"createdAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty" bson:"completedAt,omitempty"`
	Analysis    *Analysis  `json:"analysis,omitempty" bson:"analysis,omitempty"`
	Report      string     `json:"report,omitempty" bson:"report,omitempty"`
}

type Analysis struct {
	RunID       string   `json:"run_id" bson:"run_id"`
	Steps       []Step   `json:"steps" bson:"steps"`
	Issues      []Issue  `json:"issues" bson:"issues"`
	Summary     string   `json:"summary" bson:"summary"`
	Screenshots []string `json:"screenshots" bson:"screenshots"`
}

type Step struct {
	Name, Description, Screenshot string
	Notes                         string `json:"notes,omitempty" bson:"notes,omitempty"`
}

type Issue struct {
	Severity, Title, Details string
}

func NewRun(url string) *Run {
	now := time.Now()
	return &Run{
		ID:        uuid.NewString(),
		URL:       url,
		Status:    "processing",
		CreatedAt: now,
	}
}
