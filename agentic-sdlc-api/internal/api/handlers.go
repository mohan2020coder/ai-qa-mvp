package api

import (
	"agentic-sdlc-api/internal/workflow"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms/ollama"
)

type ServerDeps struct {
	// Store         *store.InMemStore
	// Orchestrators map[string]*orchestrator.SimpleOrchestrator
	LLM *ollama.LLM
	DB  *sql.DB
}

// Create a new workflow
func (s *ServerDeps) CreateWorkflowv2(c *gin.Context) {
	var req struct {
		Name  string                   `json:"name"`
		Nodes []map[string]interface{} `json:"nodes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	tx, err := s.DB.Begin()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	var workflowID string
	err = tx.QueryRow(`INSERT INTO workflows(name) VALUES($1) RETURNING id`, req.Name).Scan(&workflowID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Insert nodes into DB
	for _, n := range req.Nodes {
		id := n["id"].(string)
		nodeType := n["type"].(string)

		config, _ := json.Marshal(n["config"])
		nextIDs, _ := json.Marshal(n["next"])
		condition, _ := n["condition"].(string)

		_, err = tx.Exec(`
			INSERT INTO workflow_nodes (id, workflow_id, name, type, config, next_ids, condition)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
		`, id, workflowID, n["name"], nodeType, config, nextIDs, condition)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"id": workflowID})
}

// List workflows
func (s *ServerDeps) ListWorkflowsv2(c *gin.Context) {
	rows, err := s.DB.Query(`SELECT id, name FROM workflows`)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var workflows []map[string]any
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			continue
		}
		workflows = append(workflows, map[string]any{
			"id":   id,
			"name": name,
		})
	}

	c.JSON(200, workflows)
}

// (Optional) Get a single workflow
func (s *ServerDeps) GetWorkflowv2(c *gin.Context) {
	id := c.Param("id")
	wf, err := workflow.LoadWorkflowFromDB(s.DB, s.LLM, id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	// Build a JSON-safe response (exclude Agent runtime)
	nodes := []map[string]any{}
	for _, n := range wf.Nodes {
		nodes = append(nodes, map[string]any{
			"id":        n.ID,
			"config":    n.Config,
			"next":      n.Next,
			"condition": n.Condition,
		})
	}

	c.JSON(200, gin.H{
		"id":    wf.ID,
		"name":  wf.Name,
		"nodes": nodes,
	})
}

/*
version 1 - which is depreceted
*/

// func (d *ServerDeps) CreateWorkflowv1(c *gin.Context) {
// 	var w contracts.Workflow
// 	if err := c.BindJSON(&w); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Generate an ID for workflow
// 	if w.ID == "" {
// 		w.ID = uuid.NewString()
// 	}

// 	// Build agents dynamically from request
// 	agentsList := agents.BuildAgents(d.LLM, w.Agents)

// 	// Attach orchestrator to workflow
// 	orc := orchestrator.New(agentsList, ".workspace")

// 	// Save workflow
// 	created, _ := d.Store.CreateWorkflow(w)

// 	// Store orchestrator for later executions
// 	d.Orchestrators[w.ID] = orc

// 	c.JSON(http.StatusCreated, created)
// }

// func (d *ServerDeps) ListWorkflowsv1(c *gin.Context) {
// 	c.JSON(http.StatusOK, d.Store.ListWorkflows())
// }

// func (d *ServerDeps) StartExecutionv1(c *gin.Context) {
// 	id := c.Param("id")
// 	// fetch workflow
// 	var wf contracts.Workflow
// 	found := false
// 	for _, w := range d.Store.ListWorkflows() {
// 		if w.ID == id {
// 			wf = w
// 			found = true
// 			break
// 		}
// 	}
// 	if !found {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "workflow not found"})
// 		return
// 	}

// 	// get orchestrator for this workflow
// 	orc, ok := d.Orchestrators[wf.ID]
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "orchestrator not found for workflow"})
// 		return
// 	}

// 	// run orchestrator (sync for demo)
// 	report, err := orc.Execute(c.Request.Context(), wf.Spec)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	exec := contracts.Execution{WorkflowID: wf.ID, Status: "succeeded"}
// 	created, _ := d.Store.CreateExecution(exec)
// 	c.JSON(http.StatusCreated, gin.H{"execution": created, "report": report})
// }
