package report

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// --- Chat Request/Response Structs ---

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatReq struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream,omitempty"`
}

type chatResp struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type ollamaResp struct {
	Results []struct {
		Output string `json:"output"`
	} `json:"results"`
}

// --- Call LLM Function ---

func callLLM(cfg *LLMConfig, req chatReq) (string, error) {
	b, _ := json.Marshal(req)

	var url string
	if cfg.Provider == "ollama" {
		url = fmt.Sprintf("%s/chat", cfg.BaseURL)
	} else {
		url = fmt.Sprintf("%s/chat/completions", cfg.BaseURL)
	}

	fmt.Printf("Calling LLM at %s with model %s\n", url, req.Model)

	httpClient := &http.Client{Timeout: 60 * time.Second}
	reqHTTP, _ := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewReader(b))
	if cfg.AuthHeader != "" {
		reqHTTP.Header.Set("Authorization", cfg.AuthHeader)
	}
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(reqHTTP)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("LLM raw response:", string(body)) // debug raw JSON

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("llm status %d", resp.StatusCode)
	}

	switch cfg.Provider {
	case "ollama":
		var out ollamaResp
		if err := json.Unmarshal(body, &out); err != nil {
			return "", err
		}
		if len(out.Results) == 0 {
			return "", fmt.Errorf("no ollama results")
		}
		return out.Results[0].Output, nil

	default: // openai/openrouter
		var out chatResp
		if err := json.Unmarshal(body, &out); err != nil {
			return "", err
		}
		if len(out.Choices) == 0 {
			return "", fmt.Errorf("no llm choices")
		}
		return out.Choices[0].Message.Content, nil
	}
}

// Call Ollama LLM with streaming
// Call Ollama LLM with streaming and debug logs
func callOllamaLLM(cfg *LLMConfig, req chatReq) (string, error) {
	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshaling request:", err)
		return "", err
	}

	url := fmt.Sprintf("%s/chat", cfg.BaseURL)
	fmt.Println("Ollama URL:", url)
	fmt.Println("Request body:", string(b))

	httpClient := &http.Client{
		Timeout: 0, // No timeout for streaming
	}

	reqHTTP, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewReader(b))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return "", err
	}
	reqHTTP.Header.Set("Content-Type", "application/json")
	if cfg.AuthHeader != "" {
		reqHTTP.Header.Set("Authorization", cfg.AuthHeader)
	}

	resp, err := httpClient.Do(reqHTTP)
	if err != nil {
		fmt.Println("HTTP request error:", err)
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("HTTP Status Code:", resp.StatusCode)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("ollama status %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	var output string
	lineCount := 0

	for scanner.Scan() {
		lineCount++
		line := scanner.Text()
		fmt.Printf("Line %d: %s\n", lineCount, line)

		if line == "" {
			continue
		}

		var msg struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			Done bool `json:"done"`
		}

		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			fmt.Println("Error unmarshaling line:", err)
			continue
		}

		fmt.Printf("Appending content: %q\n", msg.Message.Content)
		output += msg.Message.Content

		if msg.Done {
			fmt.Println("Done flag received. Stopping scan.")
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Scanner error:", err)
		return "", err
	}

	fmt.Println("Final output:", output)
	return output, nil
}
