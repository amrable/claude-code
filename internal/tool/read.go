package tool

import (
	"encoding/json"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/shared"
)

type ReadArgs struct {
	FilePath string `json:"file_path"`
}

type Read struct{}

func NewReadTool() Read {
	return Read{}
}

func (r *Read) Run(payload string) (string, error) {
	var args ReadArgs
	err := json.Unmarshal([]byte(payload), &args)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(args.FilePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (r *Read) AsChatCompletionToolUnionParam() openai.ChatCompletionToolUnionParam {
	return openai.ChatCompletionToolUnionParam{
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
	}
}
