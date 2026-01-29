package main

import (
	"bash/internal/builtins"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/chzyer/readline"
)

type bellCompleter struct {
	bellLine      *readline.PrefixCompleter
	currentString string
	tabNum        int
}

// 初始化补全器
var standCompleter = initExecList()
var myCompleter = &bellCompleter{
	standCompleter,
	"",
	0,
}

func (b *bellCompleter) Reset() {
	b.tabNum = 0
	b.currentString = ""
}

func (b *bellCompleter) Do(line []rune, pos int) (possible [][]rune, length int) {
	newLine, offset := b.bellLine.Do(line, pos)
	currentString := string(line)

	if len(newLine) <= 1 {
		b.Reset()
		if len(newLine) == 0 {
			fmt.Printf("\x07") // Beep
		}
		return newLine, offset
	}

	lcp := lcpMatch(newLine)
	if len(lcp) > 0 {
		b.Reset()
		return [][]rune{lcp}, offset
	}
	if currentString == b.currentString {
		b.tabNum++
	} else {
		b.tabNum = 1
		b.currentString = currentString
	}
	if b.tabNum == 1 {
		fmt.Printf("\x07")
		return nil, 0
	} else {
		prefix := currentString
		var candidates []string
		for _, runes := range newLine {
			fullLine := prefix + string(runes)
			candidates = append(candidates, fullLine)
		}
		sort.Strings(candidates)
		fmt.Printf("\n")
		allPossible := strings.Join(candidates, " ")
		fmt.Println(allPossible)
		fmt.Print("$ " + currentString)
		return nil, 0
	}
}

func getAllPathExec() []string {
	pathAll := os.Getenv("PATH")
	list := filepath.SplitList(pathAll)
	var execList []string
	flag := make(map[string]bool)
	for _, dir := range list {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			info, err := entry.Info()
			if err != nil {
				continue
			}
			// 检查是否可执行 (Linux/Mac)
			if info.Mode()&0111 != 0 {
				name := entry.Name()
				if !flag[name] {
					execList = append(execList, name)
					flag[name] = true
				}
			}
		}
	}
	return execList
}

func initExecList() *readline.PrefixCompleter {
	var items []readline.PrefixCompleterInterface
	commandMap := make(map[string]bool)

	var addCod = func(name string) {
		if !commandMap[name] {
			commandMap[name] = true
			items = append(items, readline.PcItem(name))
		}
	}

	// 1. 动态添加内建命令 (从 builtins 包读取)
	// 这样以后你加新命令，补全自动就有了！
	for name := range builtins.CodFunc {
		addCod(name)
	}

	// 2. 添加环境变量 PATH 里的外部命令
	// 注意：为了启动速度，这里可以考虑放到 goroutine 里或者懒加载
	// 但目前保持原样即可
	execs := getAllPathExec()
	for _, exec := range execs {
		if !commandMap[exec] {
			addCod(exec)
		}
	}

	var completer = readline.NewPrefixCompleter(items...)
	return completer
}

func lcpMatch(candidates [][]rune) []rune {
	if len(candidates) == 0 {
		return nil
	}
	prefix := candidates[0]
	for _, candidate := range candidates[1:] {
		match := 0
		length := len(prefix)
		if len(candidate) < length {
			length = len(candidate)
		}
		for i := 0; i < length; i++ {
			if prefix[i] != candidate[i] {
				break
			}
			match++
		}
		prefix = prefix[:match]
		if len(prefix) == 0 {
			return nil
		}
	}
	return prefix
}
