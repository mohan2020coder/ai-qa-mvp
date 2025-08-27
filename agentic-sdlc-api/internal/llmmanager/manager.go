package llmmanager

import (
	"sync/atomic"

	"agentic-sdlc-api/internal/agents"
)

type Manager struct {
	current atomic.Value // holds agents.LLMInterface
}

func NewManager(llm agents.LLMInterface) *Manager {
	m := &Manager{}
	m.current.Store(llm)
	return m
}

func (m *Manager) Get() agents.LLMInterface {
	return m.current.Load().(agents.LLMInterface)
}

func (m *Manager) Switch(llm agents.LLMInterface) {
	m.current.Store(llm)
}
