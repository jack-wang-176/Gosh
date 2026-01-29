package builtins

import (
	"bash/internal/session"
	"fmt"
	"io"
	"os"
)

// handlePwd 打印当前工作目录
func handlePwd(_ []string, _ io.Reader, output io.Writer, sess *session.Session) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(output, dir)
	return err
}
