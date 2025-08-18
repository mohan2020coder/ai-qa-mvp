import os
import sys
import asyncio

# --- Windows-specific asyncio fix (must run before Playwright is used) ---
if sys.platform.startswith("win"):
    # Use SelectorEventLoopPolicy on Windows so subprocess support works for Playwright.
    asyncio.set_event_loop_policy(asyncio.WindowsSelectorEventLoopPolicy())

import statistics
from fastapi import FastAPI
from pydantic import BaseModel
from typing import List, Optional
from playwright.sync_api import sync_playwright
from PIL import Image, ImageStat, UnidentifiedImageError

app = FastAPI()

DATA_DIR = os.environ.get("DATA_DIR", "/data")


# --- Models ---
class Issue(BaseModel):
    severity: str
    title: str
    details: str


class Step(BaseModel):
    name: str
    description: str
    screenshot: str
    notes: Optional[str] = None


class Analysis(BaseModel):
    run_id: str
    steps: List[Step]
    issues: List[Issue]
    summary: str
    screenshots: List[str]


class RunRequest(BaseModel):
    url: str
    run_id: str


# --- Helpers ---
def artifact_url(run_id: str, fn: str) -> str:
    return f"/artifacts/runs/{run_id}/{fn}"


def artifact_path(run_id: str, fn: str) -> str:
    d = os.path.join(DATA_DIR, "runs", run_id)
    os.makedirs(d, exist_ok=True)
    return os.path.join(d, fn)


def is_blank(path: str) -> bool:
    try:
        img = Image.open(path).convert("L")
        s = ImageStat.Stat(img)
        return s.stddev[0] < 2.0
    except (FileNotFoundError, UnidentifiedImageError) as e:
        # Treat unreadable/missing image as non-blank but record issue upstream.
        return False


# --- Main route ---
@app.post("/run", response_model=Analysis)
def run(req: RunRequest):
    steps: List[Step] = []
    issues: List[Issue] = []
    screenshots: List[str] = []

    # Defensive defaults so we can always return a valid Analysis object on error
    run_id = req.run_id

    # Use sync_playwright in a try/except and ensure cleanup
    try:
        with sync_playwright() as p:
            browser = None
            try:
                browser = p.chromium.launch(headless=True)
                page = browser.new_page(viewport={"width": 1366, "height": 900})

                # Navigate
                try:
                    resp = page.goto(req.url, wait_until="domcontentloaded", timeout=60000)
                except Exception as e:
                    issues.append(
                        Issue(
                            severity="critical",
                            title="Navigation failure",
                            details=str(e),
                        )
                    )
                    return Analysis(
                        run_id=run_id,
                        steps=steps,
                        issues=issues,
                        summary="Failed to load page",
                        screenshots=screenshots,
                    )

                # --- Screenshot 1 ---
                s1 = "01_initial.png"
                p1 = artifact_path(run_id, s1)
                try:
                    page.screenshot(path=p1, full_page=True)
                except Exception as e:
                    issues.append(
                        Issue(
                            severity="high",
                            title="Screenshot failed",
                            details=f"Initial screenshot error: {e}",
                        )
                    )
                else:
                    u1 = artifact_url(run_id, s1)
                    screenshots.append(u1)
                    steps.append(
                        Step(
                            name="Initial Load",
                            description=f"Navigated to {req.url}",
                            screenshot=u1,
                        )
                    )

                # --- Heuristics ---
                status = resp.status if resp else None
                try:
                    content = page.content()
                except Exception:
                    content = ""

                markers = [
                    "404",
                    "Not Found",
                    "500",
                    "Error",
                    "Exception",
                    "Bad Gateway",
                    "Gateway Timeout",
                ]

                if status and status >= 400:
                    issues.append(
                        Issue(
                            severity="critical",
                            title=f"HTTP {status}",
                            details="Page responded with HTTP error",
                        )
                    )
                elif any(m in content for m in markers):
                    issues.append(
                        Issue(
                            severity="high",
                            title="Error marker detected",
                            details="Page contains common error keywords",
                        )
                    )

                # blank-page heuristic (only if screenshot exists)
                if os.path.exists(p1) and is_blank(p1):
                    issues.append(
                        Issue(
                            severity="high",
                            title="Suspected blank page",
                            details="Initial screenshot variance is extremely low",
                        )
                    )

                # --- Scroll + Screenshot 2 ---
                try:
                    page.mouse.wheel(0, 800)
                    page.wait_for_timeout(500)
                    s2 = "02_scrolled.png"
                    p2 = artifact_path(run_id, s2)
                    page.screenshot(path=p2, full_page=True)
                    u2 = artifact_url(run_id, s2)
                    screenshots.append(u2)
                    steps.append(
                        Step(
                            name="Scroll",
                            description="Scrolled page",
                            screenshot=u2,
                        )
                    )
                except Exception as e:
                    issues.append(
                        Issue(
                            severity="medium",
                            title="Scroll/screenshot failed",
                            details=str(e),
                        )
                    )
            finally:
                # Ensure browser is closed even if something above fails
                try:
                    if browser:
                        browser.close()
                except Exception:
                    pass

    except NotImplementedError as e:
        # This is the specific Windows/asyncio-subprocess issue many see:
        # convert to a graceful Analysis response with an actionable note.
        issues.append(
            Issue(
                severity="critical",
                title="Playwright subprocess not supported by event loop",
                details=(
                    "Playwright attempted to spawn a subprocess but the current asyncio "
                    "event loop policy does not support subprocesses on Windows. "
                    "Try adding `asyncio.set_event_loop_policy(asyncio.WindowsSelectorEventLoopPolicy())` "
                    "before Playwright is used, run under WSL/Docker, or use Python 3.9+."
                ),
            )
        )
        return Analysis(
            run_id=run_id,
            steps=steps,
            issues=issues,
            summary="Playwright startup failed on this platform",
            screenshots=screenshots,
        )
    except Exception as e:
        # Catch any other unexpected exceptions and return a 500-like Analysis result
        issues.append(
            Issue(
                severity="critical",
                title="Analyzer exception",
                details=str(e),
            )
        )
        return Analysis(
            run_id=run_id,
            steps=steps,
            issues=issues,
            summary="Analyzer failed with an exception",
            screenshots=screenshots,
        )

    summary = (
        "Basic load + scroll completed. "
        "Heuristics: HTTP status, error keywords, blank-page variance."
    )

    return Analysis(
        run_id=run_id,
        steps=steps,
        issues=issues,
        summary=summary,
        screenshots=screenshots,
    )
