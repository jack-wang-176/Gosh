package main

import (
	"bash/commands"
	"fmt"
	"io"
	"os"

	"os/exec"
)

type WaitFunc func() error

func startCommands(args []string, stdin io.Reader, stdout io.Writer) (WaitFunc, bool, error) {
	name := args[0]
	codFunc, flag := commands.CodFunc[name]
	if flag {
		done := make(chan error, 1)
		go func() {
			err := codFunc(args[1:], stdin, stdout)
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
		}, true, nil
	}

	cmd := exec.Command(name, args[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, false, fmt.Errorf("%s: command not found", name)
	}

	return cmd.Wait, false, nil
}
