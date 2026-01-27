package commands

import "io"

var CodFunc = map[string]CommandFunc{
	"exit":    handleExit,
	"echo":    handleEcho,
	"type":    handleType,
	"pwd":     handlePwd,
	"cd":      handleCd,
	"history": handleHistory,
}
var History []string
var HistoryIndex = 0

type CommandFunc func(args []string, input io.Reader, output io.Writer) error
