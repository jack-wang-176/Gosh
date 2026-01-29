package session

import (
	"io"
	"os"
)

type Session struct {
	History      []string
	HistoryIndex int
	input        io.Reader
	output       io.Writer
	stderr       io.Writer
}

func NewSession() *Session {
	return &Session{
		input:        os.Stdin,
		output:       os.Stdout,
		stderr:       os.Stderr,
		History:      []string{},
		HistoryIndex: 0,
	}
}
