package report

import (
	"fmt"
	"os"
)

type LLMConfig struct {
	Provider   string
	APIKey     string
	BaseURL    string
	Model      string
	AuthHeader string
}

func LoadLLMConfig() (*LLMConfig, error) {
	provider := os.Getenv("LLM_PROVIDER")
	if provider == "" {
		provider = "openai"
	}

	cfg := &LLMConfig{Provider: provider}

	switch provider {
	case "openai":
		cfg.APIKey = os.Getenv("OPENAI_API_KEY")
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY not set")
		}
		cfg.BaseURL = os.Getenv("OPENAI_BASE_URL")
		if cfg.BaseURL == "" {
			cfg.BaseURL = "https://api.openai.com/v1"
		}
		cfg.Model = os.Getenv("OPENAI_MODEL")
		if cfg.Model == "" {
			cfg.Model = "gpt-4o-mini"
		}
		cfg.AuthHeader = "Bearer " + cfg.APIKey

	case "openrouter":
		cfg.APIKey = os.Getenv("OPENROUTER_API_KEY")
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("OPENROUTER_API_KEY not set")
		}
		cfg.BaseURL = os.Getenv("OPENROUTER_BASE_URL")
		if cfg.BaseURL == "" {
			cfg.BaseURL = "https://openrouter.ai/api/v1"
		}
		cfg.Model = os.Getenv("OPENROUTER_MODEL")
		if cfg.Model == "" {
			cfg.Model = "anthropic/claude-3.5-sonnet"
		}
		cfg.AuthHeader = "Bearer " + cfg.APIKey

	case "ollama":
		cfg.BaseURL = os.Getenv("OLLAMA_BASE_URL")
		if cfg.BaseURL == "" {
			cfg.BaseURL = "http://localhost:11434/api"
		}
		cfg.Model = os.Getenv("OLLAMA_MODEL")
		if cfg.Model == "" {
			cfg.Model = "llama3.1"
		}
		cfg.AuthHeader = "" // no key needed by default

	default:
		return nil, fmt.Errorf("unknown LLM_PROVIDER: %s", provider)
	}

	return cfg, nil
}
