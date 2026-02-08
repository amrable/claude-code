package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/shared"
	"github.com/sirupsen/logrus"
)

func main() {
	var prompt string
	flag.StringVar(&prompt, "p", "", "Prompt to send to LLM")
	flag.Parse()

	if prompt == "" {
		panic("Prompt must not be empty")
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	baseUrl := os.Getenv("OPENROUTER_BASE_URL")
	if baseUrl == "" {
		baseUrl = "https://openrouter.ai/api/v1"
	}

	if apiKey == "" {
		panic("Env variable OPENROUTER_API_KEY not found")
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
					{
						OfFunction: &openai.ChatCompletionFunctionToolParam{
							Type: "function",
							Function: shared.FunctionDefinitionParam{
								Name:        "Read",
								Description: param.Opt[string]{Value: "Read and return the contents of a file"},
								Parameters: shared.FunctionParameters{
									"type": "object",
									"properties": map[string]any{
										"file_path": map[string]any{
											"type":        "string",
											"description": "The path to the file to read",
										},
									},
									"required": []string{"file_path"},
								},
							},
						},
					},
				},
			},
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if len(resp.Choices) == 0 {
			panic("No choices in response")
		}

		// You can use print statements as follows for debugging, they'll be visible when running tests.
		// fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

		if resp.Choices[0].Message.Content != "" {
			fmt.Println(resp.Choices[0].Message.Content)
		}
		if len(resp.Choices[0].Message.ToolCalls) == 0 {
			break
		}

		assistantMessageParam := resp.Choices[0].Message.ToAssistantMessageParam()
		messages = append(messages, openai.ChatCompletionMessageParamUnion{OfAssistant: &assistantMessageParam})
		for _, toolCall := range resp.Choices[0].Message.ToolCalls {
			switch toolCall.Function.Name {
			case "Read":
				toolReturn := read(unmarhsalAndGet(toolCall.Function.Arguments))
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

func unmarhsalAndGet(payload string) string {
	var t struct {
		FilePath string `json:"file_path"`
	}
	json.Unmarshal([]byte(payload), &t)
	return t.FilePath
}

func read(filePath string) string {
	logrus.Infof("reading file path:%s\n", filePath)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err)
	}
	return string(content)
}
