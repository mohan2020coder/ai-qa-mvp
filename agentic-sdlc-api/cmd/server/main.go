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

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	var backend string
	var model string
	var addr string
	flag.StringVar(&backend, "backend", "ollama", "LLM backend to use: ollama | llama")
	flag.StringVar(&model, "model", os.Getenv("LLM_MODEL"), "LLM model name")
	flag.StringVar(&addr, "addr", ":8080", "server bind address")

	flag.Parse()

	// if model == "" {
	// 	if backend == "ollama" {
	// 		model = "llama3"
	// 	} else {
	// 		model = "phi-2.Q4_K_M"
	// 	}
	// }

	if model == "" {
		if backend == "ollama" {
			model = "phi"
		} else {
			model = "phi-2.Q4_K_M"
		}
	}

	// var llm agents.LLMInterface
	var err error
	ctx := context.Background()

	// switch backend {
	// case "ollama":
	llm, err := ollama.New(
		ollama.WithModel(model),
		ollama.WithServerURL("http://localhost:9090"), // optional if default
	)
	if err != nil {
		log.Fatalf("failed to initialize Ollama LLM: %v", err)
	}

	// case "llama":
	// 	llm, err = openai.New(
	// 		openai.WithBaseURL("http://localhost:9090/v1"), // llama.cpp server
	// 		openai.WithModel(model),
	// 		openai.WithToken("dummy"), // model name
	// 	)
	// 	if err != nil {
	// 		log.Fatalf("failed to initialize llama.cpp LLM: %v", err)
	// 	}
	// default:
	// 	log.Fatalf("unknown backend: %s", backend)
	// }

	// In-memory store
	// st := store.NewInMemStore()

	// DB connection for v2
	dsn := "postgres://postgres:postgres@localhost:5432/workflowdb?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	// API server
	r := gin.Default()
	// api.RegisterRoutes(r, st, llm)

	// v2 dynamic workflows (from DB)
	api.RegisterWorkflowRoutes(r, llm, db)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	fmt.Printf("Server listening on %s (backend=%s, model=%s)\n", addr, backend, model)

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
