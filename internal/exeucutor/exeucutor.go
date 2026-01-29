package exeucutor

import (
	"bash/internal/builtins"
	"bash/internal/parser"
	"bash/internal/session"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type WaitFunc func() error

func Execute(pl *parser.Pipeline, session session.Session) error {
	var prePipeReader *os.File = nil
	var waitFunctions []WaitFunc
	var resourcesToClose []io.Closer
	defer func() {
		for _, closer := range resourcesToClose {
			_ = closer.Close()
		}
	}()

	for n, cmd := range pl.Commands {
		var input io.Reader = os.Stdin
		var output io.Writer = os.Stdout
		var errOutput io.Writer = os.Stderr
		var nextPipeReader *os.File
		var nextPipeWriter *os.File

		if prePipeReader != nil {
			input = prePipeReader
		}
		if cmd.InputFile != "" {
			inputFile, err := os.Open(cmd.InputFile)
			if err != nil {
				return fmt.Errorf("Error opening input file: %s ", err)
			}
			resourcesToClose = append(resourcesToClose, inputFile)
			input = inputFile
		}
		if cmd.ErrFile != "" {
			flag := os.O_TRUNC | os.O_CREATE | os.O_WRONLY
			if cmd.AppendMode { // 注意：parser里没区分 stdout append 和 stderr append，这里假设共用标记
				flag |= os.O_APPEND
			}
			// 打开文件
			file, err := os.OpenFile(cmd.ErrFile, flag, 0644)
			if err != nil {
				return fmt.Errorf("Error opening error file: %s ", err)
			}
			// 加入自动关闭列表
			resourcesToClose = append(resourcesToClose, file)
			errOutput = file
		}

		isLastCmd := n == len(pl.Commands)-1
		if !isLastCmd {
			var err error
			nextPipeReader, nextPipeWriter, err = os.Pipe()
			if err != nil {
				return fmt.Errorf("Error creating pipe: %s ", err)
			}
			output = nextPipeWriter
		} else {
			if cmd.OutputFile != "" {
				flag := os.O_TRUNC | os.O_CREATE | os.O_WRONLY
				if cmd.AppendMode {
					flag |= os.O_APPEND
				}
				file, err := os.OpenFile(cmd.OutputFile, flag, 0644)
				if err != nil {
					return fmt.Errorf("Error opening output file: %s ", err)
				}
				resourcesToClose = append(resourcesToClose, file)
				output = file
			}
		}
		wait, err := startSingleCommand(cmd.Args, input, output, errOutput, &session)
		if err != nil {
			// 如果启动失败，别忘了关掉刚才创建的管道，防止泄漏
			if nextPipeWriter != nil {
				_ = nextPipeWriter.Close()
			}
			if nextPipeReader != nil {
				_ = nextPipeReader.Close()
			}
			return err
		}
		waitFunctions = append(waitFunctions, wait)
		if nextPipeReader != nil {
			_ = nextPipeReader.Close()
		}
		if prePipeReader != nil {
			_ = prePipeReader.Close()
		}
		prePipeReader = nextPipeReader
	}
	var finalErr error
	for _, wait := range waitFunctions {
		if err := wait(); err != nil {
			finalErr = err
		}
	}

	return finalErr
}
func startSingleCommand(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer, sess *session.Session) (WaitFunc, error) {
	name := args[0]
	codFunc, flag := builtins.CodFunc[name]
	if flag {
		done := make(chan error, 1)
		go func() {
			err := codFunc(args[1:], stdin, stdout, sess)
			done <- err
			if w, ok := stdout.(io.WriteCloser); ok && stdout != os.Stdout && stdout != os.Stderr {
				_ = w.Close()
			}
			if r, ok := stdin.(io.Closer); ok && stdin != os.Stdin {
				_ = r.Close()
			}
		}()

		return func() error {
			return <-done
		}, nil
	}

	cmd := exec.Command(name, args[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("%s: command not found", name)
	}

	return cmd.Wait, nil
}
