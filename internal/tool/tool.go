package tool

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/shared"
)

type Tool interface {
	Run(payload string) (string, error)
	AsChatCompletionToolUnionParam() openai.ChatCompletionToolUnionParam
}

type RunFn[Args any] func(Args) (string, error)

type BaseTool[Args any] struct {
	Name        string
	Description string
	RunFn       RunFn[Args]
	Parameters  map[string]any
	Required    []string
}

func NewBaseTool[Args any](name string, description string, runFn RunFn[Args], parameters map[string]any, required []string) BaseTool[Args] {
	return BaseTool[Args]{
		Name:        name,
		Description: description,
		RunFn:       runFn,
		Parameters:  parameters,
		Required:    required,
	}
}

func (t *BaseTool[Args]) Run(payload string) (string, error) {
	var args Args
	err := json.Unmarshal([]byte(payload), &args)
	if err != nil {
		return "", err
	}
	return t.RunFn(args)
}

func (t *BaseTool[Args]) AsChatCompletionToolUnionParam() openai.ChatCompletionToolUnionParam {
	return openai.ChatCompletionToolUnionParam{
		OfFunction: &openai.ChatCompletionFunctionToolParam{
			Type: "function",
			Function: shared.FunctionDefinitionParam{
				Name:        t.Name,
				Description: param.Opt[string]{Value: t.Description},
				Parameters: shared.FunctionParameters{
					"type":       "object",
					"properties": t.Parameters,
					"required":   t.Required,
				},
			},
		},
	}
}

type ReadArgs struct {
	FilePath string `json:"file_path"`
}

func NewReadTool() Tool {
	tool := NewBaseTool(
		"Read",
		"Read and return the contents of a file",
		func(args ReadArgs) (string, error) {
			content, err := os.ReadFile(args.FilePath)
			if err != nil {
				return "", err
			}
			return string(content), nil
		},
		map[string]any{
			"file_path": map[string]any{
				"type":        "string",
				"description": "The path to the file to read",
			},
		},
		[]string{"file_path"},
	)
	return &tool
}

type WriteArgs struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

func NewWriteTool() Tool {
	tool := NewBaseTool(
		"Write",
		"Write content to a file",
		func(args WriteArgs) (string, error) {
			err := os.WriteFile(args.FilePath, []byte(args.Content), 0644)
			if err != nil {
				return "", err
			}
			return "write operation is successfully done to " + args.FilePath, nil
		},
		map[string]any{
			"file_path": map[string]any{
				"type":        "string",
				"description": "The path of the file to write to",
			},
			"content": map[string]any{
				"type":        "string",
				"description": "The content to write to the file",
			},
		},
		[]string{"file_path", "content"},
	)
	return &tool
}

type BashArgs struct {
	Command string `json:"command"`
}

func NewBashTool() Tool {
	tool := NewBaseTool(
		"Bash",
		"Execute a shell command",
		func(args BashArgs) (string, error) {
			cmd := exec.Command("/bin/sh", "-c", args.Command)
			output, err := cmd.CombinedOutput()
			res := string(output)
			if err != nil {
				return "", err
			}
			if res != "" {
				return res, nil
			}
			return "executed " + args.Command + " successfully", nil
		},
		map[string]any{
			"command": map[string]any{
				"type":        "string",
				"description": "The command to execute",
			},
		},
		[]string{"command"},
	)
	return &tool
}
