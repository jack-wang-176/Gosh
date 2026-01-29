package builtins

import (
	"bash/internal/session"
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func handleHistory(args []string, _ io.Reader, output io.Writer, sess *session.Session) error {
	startIndex := 0
	historyLen := len(sess.History)
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
					sess.History = append(sess.History, line)
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
			for _, h := range sess.History {
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
			if sess.HistoryIndex < len(sess.History) {
				history := sess.History[sess.HistoryIndex:]
				for _, h := range history {
					if _, err := fmt.Fprintln(openFile, h); err != nil {
						return err
					}
				}
				sess.HistoryIndex = len(sess.History)
			}

			return nil
		}
	}
	for i := startIndex; i < historyLen; i++ {
		_, err := fmt.Fprintf(output, "%5d  %s\n", i+1, sess.History[i])
		if err != nil {
			return err
		}
	}
	return nil
}
