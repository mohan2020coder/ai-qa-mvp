package agents

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

type BaseLLMAgent struct {
	name   string
	llm    llms.Model
	system string
	maxTok int
}

func NewLLMAgent(name string, llm llms.Model, system string, maxTok int) *BaseLLMAgent {
	return &BaseLLMAgent{name: name, llm: llm, system: system, maxTok: maxTok}
}

func (b *BaseLLMAgent) Name() string { return b.name }

func (b *BaseLLMAgent) Run(ctx context.Context, input string) (string, error) {
	msgs := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, b.system),
		llms.TextParts(llms.ChatMessageTypeHuman, input),
	}

	resp, err := b.llm.GenerateContent(
		ctx,
		msgs,
		llms.WithMaxTokens(b.maxTok),
	)
	if err != nil {
		return "", err
	}
	if resp == nil || len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from model")
	}
	return resp.Choices[0].Content, nil
}

func (b *BaseLLMAgent) MaxTokens() int { return b.maxTok }

func (b *BaseLLMAgent) SetMaxTokens(n int) { b.maxTok = n }
