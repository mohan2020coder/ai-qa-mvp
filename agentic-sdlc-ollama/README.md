# Agentic SDLC (Ollama + LangChainGo)

A **working agentic SDLC orchestrator** where each phase agent (Product, Design, Code, Test, Deploy) is powered by a **local LLM via Ollama** using `langchaingo`.

## Requirements
- Go 1.22+
- [Ollama](https://ollama.com/download)
- A local model (default: `llama3`). Pull with:
  ```bash
  ollama pull llama3
  ```
- Run the Ollama server:
  ```bash
  ollama serve
  ```

## Run
```bash
# from repo root
go run ./cmd/agentd   -spec examples/sample-spec.md   -model llama3   -out .workspace
```

Artifacts (agent outputs) will be written into `.workspace/outputs/` as markdown.



## Notes
- No import cycles: `main` wires agents into the orchestrator via `internal/contracts`.
- Swap models via `-model` flag or `OLLAMA_MODEL` env var.


-------------------------------

# 1) Make sure Ollama is running and a model is pulled
ollama pull llama3
ollama serve    # keep this running in a terminal

# 2) Run the agentic SDLC pipeline
unzip agentic-sdlc-ollama.zip
cd agentic-sdlc-ollama
go run ./cmd/agentd -spec examples/sample-spec.md -model llama3 -out .workspace

 ollama pull deepseek-r1
go run ./cmd/agentd -spec examples/sample-spec.md -model deepseek-r1 -out .workspace

 ollama pull codellama
go run ./cmd/agentd -spec examples/sample-spec.md -model codellama -out .workspace