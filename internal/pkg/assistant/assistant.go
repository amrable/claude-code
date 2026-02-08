package assistant

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/codecrafters-io/claude-code-starter-go/internal/pkg/logger"
	"github.com/codecrafters-io/claude-code-starter-go/internal/pkg/tool"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
)

type Assistant struct {
	client   openai.Client
	model    string
	tools    map[string]tool.Tool
	messages []openai.ChatCompletionMessageParamUnion
	scanner  *bufio.Scanner
}

func New() (*Assistant, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	baseUrl := os.Getenv("BASE_URL")
	model := os.Getenv("MODEL")

	if apiKey == "" {
		return nil, fmt.Errorf("env variable OPENROUTER_API_KEY not found")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseUrl),
	)

	tools := map[string]tool.Tool{
		"Read":  tool.NewReadTool(),
		"Write": tool.NewWriteTool(),
		"Bash":  tool.NewBashTool(),
	}

	return &Assistant{
		client:   client,
		model:    model,
		tools:    tools,
		messages: []openai.ChatCompletionMessageParamUnion{},
		scanner:  bufio.NewScanner(os.Stdin),
	}, nil
}

func (a *Assistant) Prompt(ctx context.Context) (string, error) {
	type result struct {
		text string
		err  error
	}

	resultChan := make(chan result, 1)

	go func() {
		if a.scanner.Scan() {
			resultChan <- result{text: a.scanner.Text()}
		} else {
			if err := a.scanner.Err(); err != nil {
				resultChan <- result{err: err}
			} else {
				resultChan <- result{err: fmt.Errorf("EOF")}
			}
		}
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case r := <-resultChan:
		return r.text, r.err
	}
}

func (a *Assistant) Process(ctx context.Context, prompt string) error {
	a.messages = append(a.messages, openai.ChatCompletionMessageParamUnion{
		OfUser: &openai.ChatCompletionUserMessageParam{
			Content: openai.ChatCompletionUserMessageParamContentUnion{
				OfString: openai.String(prompt),
			},
		},
	})

	toolParams := []openai.ChatCompletionToolUnionParam{}
	for _, t := range a.tools {
		toolParams = append(toolParams, t.AsChatCompletionToolUnionParam())
	}

	for {
		resp, err := a.client.Chat.Completions.New(ctx,
			openai.ChatCompletionNewParams{
				Model:    a.model,
				Messages: a.messages,
				Tools:    toolParams,
			},
		)
		if err != nil {
			return fmt.Errorf("error: %v", err)
		}
		if len(resp.Choices) == 0 {
			return fmt.Errorf("no choices in response")
		}

		if len(resp.Choices[0].Message.ToolCalls) == 0 {
			result := resp.Choices[0].Message.Content
			logger.AssistantResponse(result)
			return nil
		}

		assistantMessageParam := resp.Choices[0].Message.ToAssistantMessageParam()
		a.messages = append(a.messages, openai.ChatCompletionMessageParamUnion{OfAssistant: &assistantMessageParam})

		for _, toolCall := range resp.Choices[0].Message.ToolCalls {
			logger.ToolCall(toolCall.Function.Name)

			t, ok := a.tools[toolCall.Function.Name]
			if !ok {
				return fmt.Errorf("unknown tool: %s", toolCall.Function.Name)
			}

			toolReturn, err := t.Run(toolCall.Function.Arguments)
			if err != nil {
				return err
			}

			_ = toolReturn // Debug logging disabled for cleaner output
			a.messages = append(a.messages, openai.ChatCompletionMessageParamUnion{
				OfTool: &openai.ChatCompletionToolMessageParam{
					Content: openai.ChatCompletionToolMessageParamContentUnion{
						OfString: param.Opt[string]{Value: toolReturn},
					},
					Role:       "tool",
					ToolCallID: toolCall.ID,
				},
			})
		}
	}
}
