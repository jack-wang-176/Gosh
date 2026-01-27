package main

import (
	"fmt"
	"strings"
)

func parseInput(input string) ([]string, error) {
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
		case ' ':
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
