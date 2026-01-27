package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// handleExit 退出 Shell
// args: 参数列表
// input: 标准输入 (虽然 exit 不用，但为了接口一致性)
// output: 标准输出
func handleExit(_ []string, _ io.Reader, _ io.Writer) error {
	os.Exit(0)
	return nil
}

// handleEcho 回显参数
func handleEcho(args []string, _ io.Reader, output io.Writer) error {
	// 使用 strings.Join 将参数拼接，并输出到 output (可能是屏幕，也可能是管道)
	// 注意：这里返回 err，如果管道破裂（如 | head），shell 需要知道写入失败
	_, err := fmt.Fprintln(output, strings.Join(args, " "))
	return err
}

// handleType 显示命令类型
func handleType(args []string, _ io.Reader, output io.Writer) error {
	// 注意：这个列表最好应该引用 map_and_interface.go 里的 CodFunc 键
	// 但为了解耦，这里暂时保持局部定义，或者你可以把它移到包级别的变量
	var builtins = map[string]bool{
		"type": true, "exit": true, "echo": true, "pwd": true, "cd": true, "history": true,
	}

	for _, cmd := range args {
		if builtins[cmd] {
			if _, err := fmt.Fprintln(output, cmd+" is a shell builtin"); err != nil {
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

// handlePwd 打印当前工作目录
func handlePwd(_ []string, _ io.Reader, output io.Writer) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(output, dir)
	return err
}

// handleCd 切换目录
func handleCd(args []string, _ io.Reader, _ io.Writer) error {
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

func handleHistory(args []string, _ io.Reader, output io.Writer) error {
	startIndex := 0
	historyLen := len(History)
	if len(args) == 1 {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("history: invalid argument")
		}
		if n < historyLen {
			startIndex = historyLen - n
		}
	} else if len(args) > 1 {
		cmd := args[0]
		file := args[1]
		if cmd == "-r" {
			openFile, err := os.OpenFile(file, os.O_RDONLY, 0)
			if err != nil {
				return fmt.Errorf("history: invalid argument")
			}
			defer func() {
				_ = openFile.Close()
			}()
			scanner := bufio.NewScanner(openFile)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" {
					History = append(History, line)
				}
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("history: error reading file: %v", err)
			}
			return nil
		}
		if cmd == "-w" {
			openFile, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				return fmt.Errorf("history: invalid argument")
			}
			defer func() {
				_ = openFile.Close()
			}()
			for _, h := range History {
				if _, err := fmt.Fprintln(openFile, h); err != nil {
					return err
				}
			}
			return nil
		}
		if cmd == "-a" {
			openFile, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				return fmt.Errorf("history: invalid argument")
			}
			defer func() {
				_ = openFile.Close()
			}()
			if HistoryIndex < len(History) {
				history := History[HistoryIndex:]
				for _, h := range history {
					if _, err := fmt.Fprintln(openFile, h); err != nil {
						return err
					}
				}
				HistoryIndex = len(History)
			}

			return nil
		}
	}
	for i := startIndex; i < historyLen; i++ {
		_, err := fmt.Fprintf(output, "%5d  %s\n", i+1, History[i])
		if err != nil {
			return err
		}
	}
	return nil
}
