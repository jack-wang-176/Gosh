package builtins

import (
	"bash/internal/session"
	"fmt"
	"io"
	"strings"
)

func handleEcho(args []string, _ io.Reader, output io.Writer, sess *session.Session) error {
	// 使用 strings.Join 将参数拼接，并输出到 output (可能是屏幕，也可能是管道)
	// 注意：这里返回 err，如果管道破裂（如 | head），session 需要知道写入失败
	_, err := fmt.Fprintln(output, strings.Join(args, " "))
	return err
}
