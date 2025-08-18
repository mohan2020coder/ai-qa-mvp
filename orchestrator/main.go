package main

import (
	"ai-qa-mvp/orchestrator/internal/report"
	"ai-qa-mvp/orchestrator/internal/store"
	"ai-qa-mvp/orchestrator/internal/types"

	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

var (
	dataDir     = envOr("DATA_DIR", "/data")
	analyzerURL = envOr("ANALYZER_URL", "http://analyzer:8000")
)

type createRunReq struct {
	URL   string `json:"url"`
	RunID string `json:"run_id"`
}

func main() {
	// Load .env file if present
	_ = godotenv.Load()

	cfg, err := report.LoadLLMConfig()
	if err != nil {
		log.Fatalf("failed to load LLM config: %v", err)
	}

	log.Printf("Using LLM provider=%s model=%s base=%s", cfg.Provider, cfg.Model, cfg.BaseURL)

	s, err := store.NewStore(context.Background())
	if err != nil {
		log.Printf("store init warning: %v", err)
		s = store.NewFileStore(dataDir)
	}

	mux := http.NewServeMux()

	// Serve UI
	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fs)

	// API endpoints
	mux.HandleFunc("/api/runs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreateRun(s, w, r)
		case http.MethodGet:
			handleListRuns(s, w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/runs/", func(w http.ResponseWriter, r *http.Request) {
		id := filepath.Base(r.URL.Path)
		if id == "runs" || id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		run, err := s.GetRun(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		writeJSON(w, run)
	})

	// Serve artifacts (screenshots, reports, etc.)
	mux.Handle("/artifacts/", http.StripPrefix("/artifacts/", http.FileServer(http.Dir(dataDir))))

	log.Println("Orchestrator listening on :8080")
	// Wrap mux with CORS middleware
	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))
}

func handleCreateRun(s store.Store, w http.ResponseWriter, r *http.Request) {
	var req createRunReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.URL == "" {
		http.Error(w, "url required", http.StatusBadRequest)
		return
	}

	run := types.NewRun(req.URL)
	req.RunID = run.ID
	if err := s.SaveRun(r.Context(), run); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	runDir := filepath.Join(dataDir, "runs", run.ID)
	_ = os.MkdirAll(runDir, 0o755)

	analysis, err := callAnalyzer(req)

	if err != nil {
		log.Printf("analyzer error: %v", err)
		run.Status = "failed"
		run.Error = err.Error()
		_ = s.SaveRun(r.Context(), run)
		writeJSON(w, run)
		return
	}

	run.Analysis = analysis
	reportText, err := report.GenerateReport(run)
	if err != nil {
		log.Printf("llm error: %v (fallback)", err)
		reportText = report.FallbackReport(run)
	}
	run.Report = reportText
	run.Status = "completed"
	t := time.Now()
	run.CompletedAt = &t
	_ = s.SaveRun(r.Context(), run)
	writeJSON(w, run)
}

func handleListRuns(s store.Store, w http.ResponseWriter, r *http.Request) {
	runs, err := s.ListRuns(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, runs)
}

func callAnalyzer(req createRunReq) (*types.Analysis, error) {
	b, _ := json.Marshal(req)

	httpClient := &http.Client{Timeout: 15 * time.Minute}
	endpoint := fmt.Sprintf("%s/run", analyzerURL)
	resp, err := httpClient.Post(endpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("analyzer status %d", resp.StatusCode)
	}

	var analysis types.Analysis
	if err := json.NewDecoder(resp.Body).Decode(&analysis); err != nil {
		return nil, err
	}
	return &analysis, nil
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// ---- CORS middleware ----
func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow any origin (good for dev; restrict for production)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}
