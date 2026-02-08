package main

import (
	"bufio"
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

var prompt string
var messages = []openai.ChatCompletionMessageParamUnion{}

func main() {
	flag.StringVar(&prompt, "p", "", "Prompt to send to LLM")
	flag.Parse()

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	baseUrl := "https://api.moonshot.ai/v1"

	if apiKey == "" {
		logrus.Fatal("Env variable OPENROUTER_API_KEY not found")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseUrl))
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		scanner.Scan()
		prompt = scanner.Text()

		for {
			messages = append(messages, openai.ChatCompletionMessageParamUnion{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String(prompt),
					},
				},
			})

			resp, err := client.Chat.Completions.New(context.Background(),
				openai.ChatCompletionNewParams{
					Model:    "kimi-k2-0905-preview",
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
				logrus.Info("Assistant: ", resp.Choices[0].Message.Content)
				break
			}

			assistantMessageParam := resp.Choices[0].Message.ToAssistantMessageParam()
			messages = append(messages, openai.ChatCompletionMessageParamUnion{OfAssistant: &assistantMessageParam})
			for _, toolCall := range resp.Choices[0].Message.ToolCalls {
				logrus.Info("calling ", toolCall.Function.Name)
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
				logrus.Debug("toolReturn", toolReturn)
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
}
