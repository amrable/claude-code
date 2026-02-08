package tool

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/shared"
)

type WriteArgs struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

type Write struct{}

func NewWriteTool() Write {
	return Write{}
}

func (r *Write) Run(payload string) (string, error) {
	var args WriteArgs
	err := json.Unmarshal([]byte(payload), &args)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(args.FilePath, []byte(args.Content), os.ModePerm)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("write operation is successfully done to %s", args.FilePath), nil
}

func (t *Write) AsChatCompletionToolUnionParam() openai.ChatCompletionToolUnionParam {
	return openai.ChatCompletionToolUnionParam{
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
	}
}
