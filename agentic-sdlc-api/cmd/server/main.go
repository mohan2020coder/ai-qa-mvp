package main

import (
	"context"
	"database/sql"
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
	_ "github.com/lib/pq"
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

	// DB connection for v2
	dsn := "postgres://postgres:postgres@localhost:5432/workflowdb?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	// API server
	r := gin.Default()
	api.RegisterRoutes(r, st, llm)

	// v2 dynamic workflows (from DB)
	api.RegisterWorkflowRoutes(r, st, llm, db)

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
