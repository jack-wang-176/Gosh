package builtins

var CodFunc = map[string]CommandFunc{
	"exit":    handleExit,
	"echo":    handleEcho,
	"type":    handleType,
	"pwd":     handlePwd,
	"cd":      handleCd,
	"history": handleHistory,
}
