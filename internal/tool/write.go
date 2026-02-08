package tool

import "os"

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
