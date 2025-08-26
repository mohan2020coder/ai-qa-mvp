Agentic SDLC API - Starter (Complete)

This starter provides a minimal API server wired to Ollama via langchaingo (pinned to v0.1.13).
It is a scaffold to extend into the configuration-driven workflow system described earlier.

Requirements:
- Go 1.22+
- Ollama running locally (ollama serve) and a model pulled, e.g. `ollama pull llama3`
- Set OLLAMA_MODEL env var or pass -model flag

Run:
    go run ./cmd/server -model llama3 -addr :8080

Endpoints (minimal):
- POST /api/v1/workflows          -> create workflow (JSON body)
- GET  /api/v1/workflows         -> list workflows
- POST /api/v1/workflows/{id}/executions -> start execution

This is an in-memory prototype (no DB). Extend to persist storage and secure auth for production.

----------------------------------------------------

v1 - this will not work
curl -X POST http://localhost:8080/api/v1/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Shopping App",
    "agents": ["Product", "Design", "Code"]
  }'

{"id":"f4499b93-4966-45ed-83c2-c267a72023ad","name":"New Shopping App","spec":""}

curl http://localhost:8080/api/v1/workflows

[{"id":"f4499b93-4966-45ed-83c2-c267a72023ad","name":"New Shopping App","spec":""}]


curl -X POST http://localhost:8080/api/v1/workflows/f4499b93-4966-45ed-83c2-c267a72023ad/executions

---------------------------
v2  -


curl -X POST http://localhost:8080/api/v1/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Shopping App",
    "agents": ["Code"],
    "spec": "Build a simple task management web app with backend in Go and frontend in React."
  }'

{
    "id":"89856e9c-5855-44b1-bf82-295476957d90",
    "name":"New Shopping App","spec":"Build a simple task management web app with backend in Go and frontend in React.",
    "agents":["Product","Design","Code"]
}

 curl -X POST http://localhost:8080/api/v1/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Shopping App",
    "agents": ["Code"],
    "spec": "Build a simple task management web app with backend in Go and frontend in React."
  }'

{"id":"93345d2f-e2fa-4b48-9a68-9fa71e0f5a6c","name":"New Shopping App","spec":"Build a simple task management web app with backend in Go and frontend in React.","agents":["Code"]}


curl -X POST http://localhost:8080/api/v1/workflows/1211e0e7-8f4a-41d9-af6f-3191dcfa1e82/executions


[{"id":"9c1ec936-661e-4439-9bd8-69bfd5a1a3c9","name":"New Shopping App","spec":"Build a simple task management web app with backend in Go and frontend in React.","agents":["Product","Design","Code","Test","Deploy"]}]


curl -X POST http://localhost:8080/api/v1/workflows/9c1ec936-661e-4439-9bd8-69bfd5a1a3c9/executions