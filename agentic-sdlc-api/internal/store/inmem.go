package store

import (
    "sync"
    "agentic-sdlc-api/internal/contracts"
    "github.com/google/uuid"
)

type InMemStore struct {
    mu        sync.Mutex
    workflows map[string]contracts.Workflow
    execs     map[string]contracts.Execution
}

func NewInMemStore() *InMemStore {
    return &InMemStore{
        workflows: map[string]contracts.Workflow{},
        execs:     map[string]contracts.Execution{},
    }
}

func (s *InMemStore) CreateWorkflow(w contracts.Workflow) (contracts.Workflow, error) {
    s.mu.Lock(); defer s.mu.Unlock()
    if w.ID == "" { w.ID = uuid.New().String() }
    s.workflows[w.ID] = w
    return w, nil
}

func (s *InMemStore) ListWorkflows() []contracts.Workflow {
    s.mu.Lock(); defer s.mu.Unlock()
    out := make([]contracts.Workflow, 0, len(s.workflows))
    for _, v := range s.workflows { out = append(out, v) }
    return out
}

func (s *InMemStore) CreateExecution(e contracts.Execution) (contracts.Execution, error) {
    s.mu.Lock(); defer s.mu.Unlock()
    if e.ID == "" { e.ID = uuid.New().String() }
    s.execs[e.ID] = e
    return e, nil
}
