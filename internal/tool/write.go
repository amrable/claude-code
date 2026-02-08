package tool

import (
	"fmt"
	"os"
)

type WriteArgs struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

func NewWriteTool() Tool[WriteArgs] {
	return Tool[WriteArgs]{
		Fn: func(args WriteArgs) string {
			err := os.WriteFile(args.FilePath, []byte(args.Content), os.ModePerm)
			if err != nil {
				return err.Error()
			}
			return fmt.Sprintf("write operation is successfully done to %s", args.FilePath)
		},
	}
}
