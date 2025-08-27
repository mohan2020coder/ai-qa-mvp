package api

import (
	"agentic-sdlc-api/internal/workflow"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms/ollama"
)

// func RegisterRoutes(r *gin.Engine, s *store.InMemStore, llm *ollama.LLM) {
// 	deps := &ServerDeps{
// 		Store:         s,
// 		Orchestrators: make(map[string]*orchestrator.SimpleOrchestrator),
// 		LLM:           llm,
// 		DB:            &sql.DB{},
// 	}

// 	v1 := r.Group("/api/v1")
// 	v1.POST("/workflows", deps.CreateWorkflowv1)
// 	v1.GET("/workflows", deps.ListWorkflowsv1)
// 	v1.POST("/workflows/:id/executions", deps.StartExecutionv1)
// }

func RegisterWorkflowRoutes(r *gin.Engine, llm *ollama.LLM, db *sql.DB) {
	deps := &ServerDeps{
		// Store: s,
		// Orchestrators: make(map[string]*orchestrator.SimpleOrchestrator),
		LLM: llm,
		DB:  db,
	}

	v1 := r.Group("/api/v2")
	v1.POST("/workflows", deps.CreateWorkflowv2)
	v1.GET("/workflows", deps.ListWorkflowsv2)
	v1.GET("/workflows/:id", deps.GetWorkflowv2) // optional
	v1.POST("/workflows/:id/executions", func(c *gin.Context) {
		id := c.Param("id")
		wf, err := workflow.LoadWorkflowFromDB(db, llm, id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		var body struct {
			Input string `json:"input"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "invalid JSON"})
			return
		}

		results, err := wf.Run(c.Request.Context(), body.Input)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"results": results})
	})

}
