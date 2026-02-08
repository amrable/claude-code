package tool

import (
	"fmt"
	"os/exec"
)

type BashArgs struct {
	Command string `json:"command"`
}

func NewBashTool() Tool[BashArgs] {
	return Tool[BashArgs]{
		Fn: func(args BashArgs) string {
			cmd := exec.Command("/bin/sh", "-c", args.Command)
			output, err := cmd.CombinedOutput()
			res := string(output)
			if err != nil {
				return fmt.Sprintf("error: %v\noutput: %s", err, res)
			}
			if res != "" {
				return res
			}
			return fmt.Sprintf("executed %s successfully", args.Command)
		},
	}
}
