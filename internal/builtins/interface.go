package builtins

import (
	"bash/internal/session"
	"io"
)

type CommandFunc func(args []string, input io.Reader, output io.Writer, sess *session.Session) error
