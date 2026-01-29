package builtins

import (
	"bash/internal/session"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// handleCd 切换目录
func handleCd(args []string, _ io.Reader, _ io.Writer, sess *session.Session) error {
	// 1. 如果没有参数，默认为跳回用户主目录
	path := "~"
	if len(args) > 0 {
		path = args[0]
	}

	// 2. 处理波浪号 ~ 展开
	if path == "~" || strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		if path == "~" {
			path = homeDir
		} else {
			// 将 "~/Documents" 替换为 "/User/home/Documents"
			path = filepath.Join(homeDir, path[2:])

		}
	}

	// 3. 执行切换
	err := os.Chdir(path)
	if err != nil {
		return fmt.Errorf("cd: %s: No such file or directory", path)
	}
	return nil
}
