package main

import (
	"bash/commands"
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chzyer/readline"
)

func main() {
	environment := os.Getenv("HISTFILE")
	if environment != "" {
		file, err := os.Open(environment)
		if err == nil {
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" {
					commands.History = append(commands.History, line)
				}
			}
			commands.HistoryIndex = len(commands.History)
		}
		err = file.Close()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
	}
	ex, err := readline.NewEx(&readline.Config{
		Prompt:          "$ ",
		AutoComplete:    myCompleter,
		EOFPrompt:       "exit",
		InterruptPrompt: "^C",
	})
	if err != nil {
		panic(err)
	}
	for {
		command, err := ex.Readline()
		if err != nil {
			os.Exit(0)
		}
		commands.History = append(commands.History, command)
		command = strings.TrimSpace(command)
		if command == "exit" {
			theEnvironment, err := os.OpenFile(environment, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
			if err != nil {
				os.Exit(0)
			}
			for _, line := range commands.History {
				_, err = fmt.Fprintln(theEnvironment, line)
				if err != nil {
					os.Exit(0)
				}
			}
			err = theEnvironment.Close()
			if err != nil {
				os.Exit(0)
			}
			os.Exit(0)
		}
		split, err := parseInput(command)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}

		if len(split) == 0 {
			continue
		}

		currentOut := os.Stdout
		currentErr := os.Stderr
		pipe, err := detectPipe(split)
		if err != nil {
			// 这里必须过滤掉 ExitError (比如 SIGPIPE 或 exit code 1)
			// 否则像 ls | type 这种导致 ls 收到 SIGPIPE 的情况会打印错误
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				// 这是一个运行时的退出状态错误 (如 broken pipe)，静默忽略
			} else {
				// 其他严重错误 (如 pipe 创建失败)，打印到 Stderr
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
		}
		if pipe {
			continue
		}
		err, split, file := detectChara(split)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			continue
		}

		waitFunc, _, err := startCommands(split, os.Stdin, os.Stdout)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		if waitFunc != nil {
			if err := waitFunc(); err != nil {
				// 准备一个空指针，用来接收解包后的 ExitError
				var exitErr *exec.ExitError

				// errors.As 会递归地解包 err，看能否找到 *exec.ExitError 类型的错误
				// 如果找到了，返回 true，并将找到的错误赋值给 exitErr
				if errors.As(err, &exitErr) {
					// 这是一个退出状态码错误 (比如 exit status 1)，静默忽略
				} else {
					// 这是一个真正的错误 (比如 I/O 错误，管道断裂等)，打印出来
					_, _ = fmt.Fprintln(os.Stderr, err)
				}
			}
		}

		if file != nil {
			_ = file.Close()
			os.Stdout = currentOut
			os.Stderr = currentErr
		}
	}
}
