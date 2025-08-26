package api

import (
	"agentic-sdlc-api/internal/agents"
	"agentic-sdlc-api/internal/contracts"
	"agentic-sdlc-api/internal/orchestrator"
	"agentic-sdlc-api/internal/store"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tmc/langchaingo/llms/ollama"
)

type ServerDeps struct {
	Store         *store.InMemStore
	Orchestrators map[string]*orchestrator.SimpleOrchestrator
	LLM           *ollama.LLM
}

func RegisterRoutes(r *gin.Engine, s *store.InMemStore, llm *ollama.LLM) {
	deps := &ServerDeps{
		Store:         s,
		Orchestrators: make(map[string]*orchestrator.SimpleOrchestrator),
		LLM:           llm,
	}

	v1 := r.Group("/api/v1")
	v1.POST("/workflows", deps.createWorkflow)
	v1.GET("/workflows", deps.listWorkflows)
	v1.POST("/workflows/:id/executions", deps.startExecution)
}

func (d *ServerDeps) createWorkflow(c *gin.Context) {
	var w contracts.Workflow
	if err := c.BindJSON(&w); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate an ID for workflow
	if w.ID == "" {
		w.ID = uuid.NewString()
	}

	// Build agents dynamically from request
	agentsList := agents.BuildAgents(d.LLM, w.Agents)

	// Attach orchestrator to workflow
	orc := orchestrator.New(agentsList, ".workspace")

	// Save workflow
	created, _ := d.Store.CreateWorkflow(w)

	// Store orchestrator for later executions
	d.Orchestrators[w.ID] = orc

	c.JSON(http.StatusCreated, created)
}

func (d *ServerDeps) listWorkflows(c *gin.Context) {
	c.JSON(http.StatusOK, d.Store.ListWorkflows())
}

func (d *ServerDeps) startExecution(c *gin.Context) {
	id := c.Param("id")
	// fetch workflow
	var wf contracts.Workflow
	found := false
	for _, w := range d.Store.ListWorkflows() {
		if w.ID == id {
			wf = w
			found = true
			break
		}
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "workflow not found"})
		return
	}

	// get orchestrator for this workflow
	orc, ok := d.Orchestrators[wf.ID]
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "orchestrator not found for workflow"})
		return
	}

	// run orchestrator (sync for demo)
	report, err := orc.Execute(c.Request.Context(), wf.Spec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	exec := contracts.Execution{WorkflowID: wf.ID, Status: "succeeded"}
	created, _ := d.Store.CreateExecution(exec)
	c.JSON(http.StatusCreated, gin.H{"execution": created, "report": report})
}
