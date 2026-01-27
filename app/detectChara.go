package main

import (
	"fmt"
	"io"
	"os"
)

// detectPipe deal command in this function
func detectPipe(split []string) (bool, error) {

	// 阶段一：解析 (Parsing Phase)

	// 1. 快速检查是否有管道符
	hasPipe := false
	for _, s := range split {
		if s == "|" {
			hasPipe = true
			break
		}
	}
	if !hasPipe {
		return false, nil
	}

	// 2. 拆分命令：将 ["ls", "|", "wc"] 转换为 [["ls"], ["wc"]]
	var commands [][]string
	var currentCommand []string
	for _, s := range split {
		if s == "|" {
			if len(currentCommand) > 0 {
				commands = append(commands, currentCommand)
				currentCommand = []string{} // 清空，准备下一个命令
			}
		} else {
			currentCommand = append(currentCommand, s)
		}
	}
	// 别忘了追加最后一个命令
	if len(currentCommand) > 0 {
		commands = append(commands, currentCommand)
	}

	// 阶段二：执行 (Execution Phase)

	var waitFunctions []WaitFunc
	var prevPipeReader *os.File = nil // 上一个管道的读取端（即当前命令的输入）

	for i, cmd := range commands {
		// 1. 确定输入源
		var input io.Reader = os.Stdin
		if prevPipeReader != nil {
			input = prevPipeReader
		}

		// 2. 确定输出目标
		var output io.Writer = os.Stdout
		var nextPipeReader *os.File = nil
		var nextPipeWriter *os.File = nil

		// 如果不是最后一个命令，需要创建管道
		if i < len(commands)-1 {
			var err error
			nextPipeReader, nextPipeWriter, err = os.Pipe()
			if err != nil {
				return true, fmt.Errorf("pipe create error: %v", err)
			}
			output = nextPipeWriter
		}

		// 3. 启动命令
		wait, isBuiltin, err := startCommands(cmd, input, output)
		if err != nil {
			// 启动失败，清理所有打开的句柄
			if prevPipeReader != nil {
				_ = prevPipeReader.Close()
			}
			if nextPipeReader != nil {
				_ = nextPipeReader.Close()
			}
			if nextPipeWriter != nil {
				_ = nextPipeWriter.Close()
			}
			return true, fmt.Errorf("start command error: %v", err)
		}
		waitFunctions = append(waitFunctions, wait)

		// 4. 资源清理 (关键步骤)

		// A. 处理【输出端】(nextPipeWriter)
		// 无论是外部命令还是内置命令，父进程都需要关闭它，
		// 因为父进程不需要写数据，只负责传递句柄。
		// 注意：如果内置命令是协程，它会持有副本或共享引用，父进程这里关了不影响协程写。
		//  为了安全起见，通常对于 Write 端，父进程必须 Close，否则 Reader 永远收不到 EOF)
		if nextPipeWriter != nil {
			if !isBuiltin {
				_ = nextPipeWriter.Close()
			}
		}

		// B. 处理【输入端】(prevPipeReader) - 也就是刚刚用完的那个 input
		// 外部命令：父进程必须关闭，因为子进程已经继承了。
		// 内置命令：协程正在使用，父进程不能关。
		if prevPipeReader != nil {
			if !isBuiltin {
				_ = prevPipeReader.Close()
			}
		}

		// C. 传递【下个读取端】
		// 不要关闭 nextPipeReader！它是下一个循环的 input。
		prevPipeReader = nextPipeReader
	}

	// 阶段三：等待 (Wait Phase)

	var finalErr error
	for _, wait := range waitFunctions {
		if err := wait(); err != nil {
			finalErr = err
		}
	}

	return true, finalErr
}
func detectChara(split []string) (error, []string, *os.File) {
	found := false
	isStdError := false
	isAppend := false
	for n := range split {
		switch split[n] {
		case ">", "1>":
			found = true
			isAppend = false
			isStdError = false
		case "2>":
			found = true
			isAppend = false
			isStdError = true
		case ">>", "1>>":
			found = true
			isAppend = true
			isStdError = false
		case "2>>":
			found = true
			isAppend = true
			isStdError = true
		case "|":
			found = true
			isAppend = false
			isStdError = false
		}
		if found {
			if n < 1 {
				return fmt.Errorf(`invalid character "%s" in input`, split[n]), split, nil
			}
			if n+1 >= len(split) {
				return fmt.Errorf(`invalid character "%s" in input`, split[n]), split, nil
			}
			filename := split[n+1]
			flag := os.O_CREATE | os.O_WRONLY
			if isAppend {
				flag |= os.O_APPEND
			} else {
				flag |= os.O_TRUNC
			}
			file, err := os.OpenFile(filename, flag, 0644)
			if err != nil {
				return err, split, nil
			}
			newSplit := split[:n]
			if isStdError {
				os.Stderr = file
			} else {
				os.Stdout = file
			}
			return nil, newSplit, file
		}
	}
	return nil, split, nil
}
