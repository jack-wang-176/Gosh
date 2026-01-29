package parser

import (
	"fmt"
	"strings"
)

type Command struct {
	Args       []string
	InputFile  string
	OutputFile string
	ErrFile    string
	AppendMode bool
}
type Pipeline struct {
	Commands []*Command
}

func Parse(line string) (*Pipeline, error) {
	split := strings.Split(line, "|")
	pipeline := &Pipeline{}
	for _, part := range split {
		cmd, err := parseSingleCommand(part)
		if err != nil {
			return nil, err
		}
		if len(cmd.Args) > 0 {
			pipeline.Commands = append(pipeline.Commands, cmd)
		}
	}
	if len(pipeline.Commands) == 0 {
		return nil, nil
	}
	return pipeline, nil
}
func parseSingleCommand(cmd string) (*Command, error) {
	tokens, err := parseToken(cmd)
	if err != nil {
		return nil, err
	}
	tokenCmd := &Command{
		Args:       make([]string, 0),
		AppendMode: false,
	}
	for n := 0; n < len(tokens); n++ {
		switch tokens[n] {
		case "<":
			// 输入重定向
			if n+1 >= len(tokens) {
				return nil, fmt.Errorf("syntax error: missing filename after '<'")
			}
			n++ // 跳过下一个 token（文件名）
			tokenCmd.InputFile = tokens[n]
		case ">", "1>":
			if n+1 >= len(tokens[n]) {
				return tokenCmd, fmt.Errorf(`invalid character "%s" in input`, tokens[n])
			}
			n++
			tokenCmd.OutputFile = tokens[n+1]
		case "2>":
			if n+1 >= len(tokens[n]) {
				return tokenCmd, fmt.Errorf(`invalid character "%s" in input`, tokens[n])
			}
			n++
			tokenCmd.ErrFile = tokens[n+1]
		case ">>", "1>>":
			if n+1 >= len(tokens[n]) {
				return tokenCmd, fmt.Errorf(`invalid character "%s" in input`, tokens[n])
			}
			n++
			tokenCmd.OutputFile = tokens[n+1]
			tokenCmd.AppendMode = true
		case "2>>":
			if n+1 >= len(tokens[n]) {
				return tokenCmd, fmt.Errorf(`invalid character "%s" in input`, tokens[n])
			}
			n++
			tokenCmd.ErrFile = tokens[n+1]
			tokenCmd.AppendMode = true
		default:
			tokenCmd.Args = append(tokenCmd.Args, tokens[n])

		}
	}
	return tokenCmd, nil
}
func parseToken(input string) ([]string, error) {
	var token strings.Builder
	var args []string
	isSingleQuote := false
	isDoubleQuote := false
	isBlackSlash := false
	for _, char := range input {
		//必须首先判断是不是前面有有效的/，必须将标志位即时逆转
		if isBlackSlash {
			if char != '\\' && char != '$' && char != '"' && char != '\n' && char != '`' && isDoubleQuote {
				token.WriteRune('\\')
			}
			token.WriteRune(char)
			isBlackSlash = false
			continue
		}
		if char == '\\' {
			if isSingleQuote {
				token.WriteRune(char)
			} else if isDoubleQuote {
				isBlackSlash = true
			} else {
				isBlackSlash = true
			}
			continue
		}
		switch char {
		case '"':
			if isSingleQuote {
				token.WriteRune(char)
			} else {
				isDoubleQuote = !isDoubleQuote
			}
		case '\'':
			if isDoubleQuote {
				token.WriteRune(char)
			} else {
				isSingleQuote = !isSingleQuote
			}

		case ' ', '\t', '\n':
			if isSingleQuote || isDoubleQuote {
				token.WriteRune(char)
			} else {
				if token.Len() > 0 {
					args = append(args, token.String())
					token.Reset()
				}
			}
		default:
			token.WriteRune(char)
		}
	}
	if token.Len() > 0 {
		args = append(args, token.String())
		token.Reset()
	}

	if isSingleQuote || isDoubleQuote {
		return args, fmt.Errorf("singal Qoute Error")
	}
	return args, nil
}
