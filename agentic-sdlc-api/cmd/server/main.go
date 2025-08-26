package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	// "agentic-sdlc-api/internal/agents"
	"agentic-sdlc-api/internal/api"
	// "agentic-sdlc-api/internal/orchestrator"
	"agentic-sdlc-api/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	var model string
	var addr string
	flag.StringVar(&model, "model", os.Getenv("OLLAMA_MODEL"), "Ollama model to use (e.g. llama3)")
	flag.StringVar(&addr, "addr", ":8080", "server bind address")
	flag.Parse()

	if model == "" {
		model = "llama3"
	}

	// Connect to Ollama
	llm, err := ollama.New(ollama.WithModel(model))
	if err != nil {
		log.Fatalf("failed to connect to ollama: %v", err)
	}

	// In-memory store
	st := store.NewInMemStore()

	// // Agents - product/design/code agents reuse same base for demo
	// prod := agents.NewProductAgent(llm)
	// design := agents.NewDesignAgent(llm)
	// code := agents.NewCodeAgent(llm)

	// agentsList := []agents.AgentInterface{prod, design, code}

	// Example: read from ENV for default workflow setup
	// agentNames := []string{"Product", "Design", "Code"}
	// if envAgents := os.Getenv("AGENTS"); envAgents != "" {
	// 	agentNames = strings.Split(envAgents, ",")
	// }

	// // Dynamically build
	// agentsList := agents.BuildAgents(llm, agentNames)

	// // Orchestrator
	// orc := orchestrator.New(agentsList, ".workspace")

	// Orchestrator
	// orc := orchestrator.New(agentsList, ".workspace")

	// API server
	r := gin.Default()
	api.RegisterRoutes(r, st, llm)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	fmt.Printf("server listening on %s (model=%s)", addr, model)

	// wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown: %v", err)
	}
	fmt.Println("Server exiting")
}
