package builtins

import (
	"bash/internal/session"
	"io"
	"os"
)

func handleExit(_ []string, _ io.Reader, _ io.Writer, sess *session.Session) error {
	os.Exit(0)
	return nil
}
