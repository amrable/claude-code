package tool

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/shared"
)

type BashArgs struct {
	Command string `json:"command"`
}

type Bash struct {
}

func NewBashTool() Bash {
	return Bash{}
}

func (b *Bash) Run(payload string) (string, error) {
	var args BashArgs
	err := json.Unmarshal([]byte(payload), &args)
	if err != nil {
		return "", err
	}
	cmd := exec.Command("/bin/sh", "-c", args.Command)
	output, err := cmd.CombinedOutput()
	res := string(output)
	if err != nil {
		return "", err
	}
	if res != "" {
		return res, nil
	}
	return fmt.Sprintf("executed %s successfully", args.Command), nil
}

func (b *Bash) AsChatCompletionToolUnionParam() openai.ChatCompletionToolUnionParam {
	return openai.ChatCompletionToolUnionParam{
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
	}
}
