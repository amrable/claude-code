package tool

import "os"

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
