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
		logrus.Info("sending...\n")
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
					{
						OfFunction: &openai.ChatCompletionFunctionToolParam{
							Type: "function",
							Function: shared.FunctionDefinitionParam{
								Name:        "Write",
								Description: param.Opt[string]{Value: "Write content to a file"},
								Parameters: shared.FunctionParameters{
									"type": "object",
									"properties": map[string]any{
										"file_path": map[string]any{
											"type":        "string",
											"description": "The path of the file to write to",
										},
										"content": map[string]any{
											"type":        "string",
											"description": "The content to write to the file",
										},
									},
									"required": []string{"file_path", "content"},
								},
							},
						},
					},
					{
						OfFunction: &openai.ChatCompletionFunctionToolParam{
							Type: "function",
							Function: shared.FunctionDefinitionParam{
								Name:        "Bash",
								Description: param.Opt[string]{Value: "Execute a shell command"},
								Parameters: shared.FunctionParameters{
									"type": "object",
									"properties": map[string]any{
										"file_path": map[string]any{
											"type":        "string",
											"description": "The command to execute",
										},
									},
									"required": []string{"command"},
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
		if len(resp.Choices[0].Message.ToolCalls) == 0 {
			fmt.Println(resp.Choices[0].Message.Content)
			break
		}

		readTool := tool.NewReadTool()
		writeTool := tool.NewWriteTool()
		bashTool := tool.NewBashTool()

		assistantMessageParam := resp.Choices[0].Message.ToAssistantMessageParam()
		messages = append(messages, openai.ChatCompletionMessageParamUnion{OfAssistant: &assistantMessageParam})
		for _, toolCall := range resp.Choices[0].Message.ToolCalls {
			var toolReturn = ""
			var err error
			switch toolCall.Function.Name {
			case "Read":
				toolReturn, err = readTool.Run(toolCall.Function.Arguments)
			case "Write":
				toolReturn, err = writeTool.Run(toolCall.Function.Arguments)
			case "Bash":
				toolReturn, err = bashTool.Run(toolCall.Function.Arguments)
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
