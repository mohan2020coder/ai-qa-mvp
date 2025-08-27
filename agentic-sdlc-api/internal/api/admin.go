// in api/admin.go
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"

	"agentic-sdlc-api/internal/agents"
	"agentic-sdlc-api/internal/llmmanager"
)

func RegisterAdminRoutes(r *gin.Engine, mgr *llmmanager.Manager) {
	r.POST("/admin/switch-llm", func(c *gin.Context) {
		var req struct {
			Backend string `json:"backend"`
			Model   string `json:"model"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var newLLM agents.LLMInterface
		var err error

		switch req.Backend {
		case "ollama":
			newLLM, err = ollama.New(ollama.WithModel(req.Model))
		case "llama":
			newLLM, err = openai.New(
				openai.WithBaseURL("http://localhost:9090/v1"), // llama.cpp server
				openai.WithModel("phi-2.Q4_K_M"),
				openai.WithToken("dummy"), // model name
			)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown backend"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		mgr.Switch(newLLM)
		c.JSON(http.StatusOK, gin.H{"message": "LLM switched", "backend": req.Backend, "model": req.Model})
	})
}
