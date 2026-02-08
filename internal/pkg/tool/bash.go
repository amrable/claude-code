package tool

import "os/exec"

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
