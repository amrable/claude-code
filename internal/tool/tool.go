package tool

import (
	"encoding/json"

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
