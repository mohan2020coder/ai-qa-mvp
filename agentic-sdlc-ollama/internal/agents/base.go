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
