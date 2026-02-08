package tool

import (
	"fmt"
	"os"
)

type ReadArgs struct {
	FilePath string `json:"file_path"`
}

func NewReadTool() Tool[ReadArgs] {
	return Tool[ReadArgs]{
		Fn: func(args ReadArgs) string {
			content, err := os.ReadFile(args.FilePath)
			if err != nil {
				return fmt.Sprintf("Error reading file: %v", err)
			}
			return string(content)
		},
	}
}
