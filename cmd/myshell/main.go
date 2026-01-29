package main

import (
	"bash/internal/exeucutor"
	"bash/internal/parser"
	"bash/internal/session"
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

func main() {
	// 1. 初始化 Session (代替了原来的全局变量)
	sess := session.NewSession()
	var resourcesToClose []io.Closer
	defer func() {
		for _, closer := range resourcesToClose {
			_ = closer.Close()
		}
	}()

	// 2. 加载历史记录 (如果有)
	// 这里的逻辑和你原来的一样，只是把 commands.History 换成了 sess.History
	envFile := os.Getenv("HISTFILE")
	if envFile != "" {
		file, err := os.Open(envFile)
		if err == nil {
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" {
					sess.History = append(sess.History, line)
				}
			}
			sess.HistoryIndex = len(sess.History)
			resourcesToClose = append(resourcesToClose, file)
		}
	}

	// 3. 初始化 Readline (自动补全和输入UI)
	// 注意：myCompleter 是你在 readline.go 里定义的，
	// 只要 readline.go 也在 main 包下，这里就能直接用。
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "$ ",
		AutoComplete:    myCompleter,
		EOFPrompt:       "exit",
		InterruptPrompt: "^C",
	})
	if err != nil {
		panic(err)
	}
	resourcesToClose = append(resourcesToClose, rl)

	// 4. 主循环 (REPL: Read-Eval-Print Loop)
	for {
		// --- READ (读取) ---
		line, err := rl.Readline()
		if err != nil { // Ctrl+D 或 错误
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 记录历史
		sess.History = append(sess.History, line)
		sess.HistoryIndex = len(sess.History)

		// 特殊处理 exit，因为要在退出前保存历史
		// (也可以做一个内置命令 handleExit，但在这里处理能确保保存逻辑执行)
		if line == "exit" {
			saveHistory(sess, envFile)
			break
		}

		// --- PARSE (解析) ---
		// 调用我们新的 Parser，把字符串变成 Pipeline 结构体
		pipeline, err := parser.Parse(line)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Syntax error:", err)
			continue
		}
		if pipeline == nil {
			continue
		}

		// --- EXECUTE (执行) ---
		// 调用我们新的 Executor，把 Pipeline 跑起来
		// 注意：这里传入 *sess (指针)，因为命令可能会修改 Session (比如 cd)
		err = exeucutor.Execute(pipeline, *sess)
		if err != nil {
			// 如果是 ExitError (比如 grep 没找到内容返回 1)，通常不报错
			// 只有真正的错误才打印
			var exitErr *os.PathError // 举例，可根据需要调整
			if !errors.As(err, &exitErr) {
				// 简单的错误打印
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
		}
	}

	// 退出前保存历史
	saveHistory(sess, envFile)
}

// 辅助函数：保存历史记录
func saveHistory(sess *session.Session, envFile string) {

	if envFile == "" {
		return
	}
	f, err := os.OpenFile(envFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()
	for _, line := range sess.History {
		_, _ = fmt.Fprintln(f, line)
	}
}
