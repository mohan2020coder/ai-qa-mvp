# AI QA Agent MVP (Golang Orchestrator + Python Analyzer + LLM Reports)

This is a minimal **working** MVP that coordinates tests in **Golang**, runs **screenshot + issue detection** in **Python** (Playwright + simple heuristics), and produces **LLM-generated reports**. A small web UI is included for triggering tests and viewing results.

## Quick start
```bash
unzip ai-qa-mvp.zip
cd ai-qa-mvp
export OPENAI_API_KEY=YOUR_KEY
# optional
export OPENAI_MODEL=gpt-4o-mini
export OPENAI_BASE_URL=https://api.openai.com/v1

docker compose up --build
```

Open **http://localhost:8080** and run a test against any public URL.

## Services
- **orchestrator** (Go): API, coordination, LLM report generation, static UI server (port 8080).
- **analyzer** (Python/FastAPI): Playwright-based screenshot + issue detection (port 8000 internal).
- **mongo** (MongoDB): Optional persistence (orchestrator falls back to JSON files if Mongo isn’t available).

## API
- `POST /api/runs` body: `{ "url": "https://example.com" }`
- `GET /api/runs`
- `GET /api/runs/{id}`
- Artifacts (screenshots) are served from: `/artifacts/runs/<id>/`

## Notes
- Screenshots and run artifacts are stored under `./data/runs/<id>/`.
- No LLM key? You’ll still get a **fallback report** so the demo works end-to-end.
- Extend flows by editing `analyzer/app.py`.
- Tweak the report prompt in `orchestrator/internal/report/llm.go`.



$env:OLLAMA_HOST="0.0.0.0"
 ollama serve



 # Create a virtual environment (Python 3.9+ recommended)
python3 -m venv venv

# Activate it
source venv/bin/activate   # on Linux / macOS
venv\Scripts\activate      # on Windows PowerShell


pip install fastapi uvicorn playwright pillow


playwright install chromium


uvicorn app:app --reload --port 8000


~ai-qa-mvp>docker-compose up -d --build orchestrator 