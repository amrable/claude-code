package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/codecrafters-io/claude-code-starter-go/internal/tool"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/sirupsen/logrus"
)

var (
	ReadTool  = tool.NewReadTool()
	WriteTool = tool.NewWriteTool()
	BashTool  = tool.NewBashTool()
)

func main() {
	var prompt string
	flag.StringVar(&prompt, "p", "", "Prompt to send to LLM")
	flag.Parse()

	if prompt == "" {
		logrus.Fatal("Prompt must not be empty")
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	baseUrl := os.Getenv("OPENROUTER_BASE_URL")
	if baseUrl == "" {
		baseUrl = "https://openrouter.ai/api/v1"
	}

	if apiKey == "" {
		logrus.Fatal("Env variable OPENROUTER_API_KEY not found")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseUrl))
	var messages = []openai.ChatCompletionMessageParamUnion{
		{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String(prompt),
				},
			},
		},
	}

	for {
		resp, err := client.Chat.Completions.New(context.Background(),
			openai.ChatCompletionNewParams{
				Model:    "anthropic/claude-haiku-4.5",
				Messages: messages,
				Tools: []openai.ChatCompletionToolUnionParam{
					ReadTool.AsChatCompletionToolUnionParam(),
					WriteTool.AsChatCompletionToolUnionParam(),
					BashTool.AsChatCompletionToolUnionParam(),
				},
			},
		)
		if err != nil {
			logrus.Fatalf("error: %v\n", err)
		}
		if len(resp.Choices) == 0 {
			logrus.Fatal("No choices in response")
		}

		if len(resp.Choices[0].Message.ToolCalls) == 0 {
			fmt.Println(resp.Choices[0].Message.Content)
			break
		}

		assistantMessageParam := resp.Choices[0].Message.ToAssistantMessageParam()
		messages = append(messages, openai.ChatCompletionMessageParamUnion{OfAssistant: &assistantMessageParam})
		for _, toolCall := range resp.Choices[0].Message.ToolCalls {
			var toolReturn = ""
			var err error
			switch toolCall.Function.Name {
			case "Read":
				toolReturn, err = ReadTool.Run(toolCall.Function.Arguments)
			case "Write":
				toolReturn, err = WriteTool.Run(toolCall.Function.Arguments)
			case "Bash":
				toolReturn, err = BashTool.Run(toolCall.Function.Arguments)
			}
			if err != nil {
				logrus.Fatal(err.Error())
			}
			logrus.Infoln("toolReturn", toolReturn)
			messages = append(messages, openai.ChatCompletionMessageParamUnion{
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
