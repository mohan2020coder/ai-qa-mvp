package agents

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

type baseLLMAgent struct {
	name   string
	llm    llms.Model
	system string
	maxTok int
}

func (b *baseLLMAgent) Name() string { return b.name }

func (b *baseLLMAgent) Run(ctx context.Context, input string) (string, error) {
	msgs := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, b.system),
		llms.TextParts(llms.ChatMessageTypeHuman, input),
	}

	resp, err := b.llm.GenerateContent(
		ctx,
		msgs,
		llms.WithMaxTokens(2048), // increase output budget
	)
	if err != nil {
		return "", err
	}
	if resp == nil || len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from model")
	}
	return resp.Choices[0].Content, nil
}

func (b *baseLLMAgent) MaxTokens() int {
	return b.maxTok
}

func CountTokens(s string) int {
	return len(s) / 4 // rough estimate: 1 token ~ 4 chars
}
func (b *baseLLMAgent) SetMaxTokens(n int) {
	b.maxTok = n
}
