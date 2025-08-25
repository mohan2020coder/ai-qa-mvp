package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// Agent interface defines an SDLC agent
type Agent interface {
	Run(ctx context.Context, input string) (string, error)
	Name() string
}

// BaseAgent implements the Agent interface
type BaseAgent struct {
	llm            llms.Model
	name           string
	promptTemplate string
}

// Run sends input formatted into promptTemplate to the LLM and returns the output
func (a *BaseAgent) Run(ctx context.Context, input string) (string, error) {
	prompt := fmt.Sprintf(a.promptTemplate, input)
	fmt.Printf(">> [%s] Prompt length: %d\n", a.name, len(prompt))
	fmt.Printf(">> [%s] Starting LLM call...\n", a.name)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Minute) // increased timeout
	defer cancel()

	response, err := llms.GenerateFromSinglePrompt(ctxWithTimeout, a.llm, prompt, llms.WithTemperature(0.5))
	if err != nil {
		fmt.Printf(">> [%s] LLM error: %v\n", a.name, err)
		return "", err
	}
	fmt.Printf(">> [%s] LLM call completed. Response length: %d\n", a.name, len(response))
	return response, nil
}

func (a *BaseAgent) Name() string {
	return a.name
}

// loadPrompt loads the prompt template text from a file
func loadPrompt(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file %s: %w", filename, err)
	}
	return string(bytes), nil
}

// cleanLLM removes any <think> blocks or trailing reasoning text
func cleanLLM(output string) string {
	// Use regex to remove all <think>...</think> blocks (non-greedy)
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	cleaned := re.ReplaceAllString(output, "")
	return strings.TrimSpace(cleaned)
}

// writeToFile writes cleaned content to a file under baseDir
func writeToFile(baseDir, filename, content string) error {
	cleaned := cleanLLM(content)
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return err
	}
	fullPath := filepath.Join(baseDir, filename)
	fmt.Printf(">> Writing file: %s\n", fullPath)
	return os.WriteFile(fullPath, []byte(cleaned), 0644)
}

func main() {
	ctx := context.Background()

	llm, err := ollama.New(ollama.WithModel("deepseek-r1"))
	if err != nil {
		log.Fatal(err)
	}

	promptDir := "./prompts"
	agentNames := []string{
		"RequirementsAgent",
		"DesignAgent",
		"DevelopmentAgent",
		"TestingAgent",
		"DeploymentAgent",
	}

	var agents []Agent
	for _, name := range agentNames {
		promptFile := filepath.Join(promptDir, name+".txt")
		promptTemplate, err := loadPrompt(promptFile)
		if err != nil {
			log.Fatalf("Failed loading prompt for %s: %v", name, err)
		}
		agents = append(agents, &BaseAgent{
			llm:            llm,
			name:           name,
			promptTemplate: promptTemplate,
		})
	}

	projectDescription := "A CLI tool written in Go that manages tasks: create, list, delete. Store tasks in a local JSON file."
	baseDir := "./gen-project"

	input := projectDescription
	for _, agent := range agents {
		fmt.Printf("\n--- Running %s ---\n", agent.Name())
		result, err := agent.Run(ctx, input)
		if err != nil {
			log.Fatalf("%s failed: %v", agent.Name(), err)
		}

		// Print full output for debugging
		fmt.Printf("\n-- %s Output --\n%s\n\n", agent.Name(), result)

		cleanedResult := cleanLLM(result)

		switch agent.Name() {
		case "RequirementsAgent":
			if err := writeToFile(baseDir, "REQUIREMENTS.txt", cleanedResult); err != nil {
				log.Fatalf("Failed writing REQUIREMENTS.txt: %v", err)
			}
		case "DesignAgent":
			if err := writeToFile(baseDir, "DESIGN.md", cleanedResult); err != nil {
				log.Fatalf("Failed writing DESIGN.md: %v", err)
			}
		case "DevelopmentAgent":
			if err := writeToFile(baseDir, "main.go", cleanedResult); err != nil {
				log.Fatalf("Failed writing main.go: %v", err)
			}
		case "TestingAgent":
			if err := writeToFile(baseDir, "main_test.go", cleanedResult); err != nil {
				log.Fatalf("Failed writing main_test.go: %v", err)
			}
		case "DeploymentAgent":
			parts := strings.SplitN(cleanedResult, "\n---\n", 2)
			if len(parts) == 2 {
				if err := writeToFile(baseDir, "Dockerfile", parts[0]); err != nil {
					log.Fatalf("Failed writing Dockerfile: %v", err)
				}
				if err := writeToFile(baseDir, "README.md", parts[1]); err != nil {
					log.Fatalf("Failed writing README.md: %v", err)
				}
			} else {
				if err := writeToFile(baseDir, "README.md", cleanedResult); err != nil {
					log.Fatalf("Failed writing README.md: %v", err)
				}
			}
			mod := "module genproject\n\ngo 1.20"
			if err := writeToFile(baseDir, "go.mod", mod); err != nil {
				log.Fatalf("Failed writing go.mod: %v", err)
			}
		}

		input = cleanedResult
	}
}
