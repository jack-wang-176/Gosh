package builtins

import (
	"bash/internal/session"
	"fmt"
	"io"
	"os/exec"
)

// handleType 显示命令类型
func handleType(args []string, _ io.Reader, output io.Writer, _ *session.Session) error {
	// 但为了解耦，这里暂时保持局部定义，或者你可以把它移到包级别的变量
	var builtins = map[string]bool{
		"type": true, "exit": true, "echo": true, "pwd": true, "cd": true, "history": true,
	}

	for _, cmd := range args {
		if builtins[cmd] {
			if _, err := fmt.Fprintln(output, cmd+" is a session builtin"); err != nil {
				return err
			}
			continue
		}

		path, err := exec.LookPath(cmd)
		if err == nil {
			if _, err := fmt.Fprintln(output, cmd+" is "+path); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintln(output, cmd+": not found"); err != nil {
				return err
			}
		}
	}
	return nil
}
