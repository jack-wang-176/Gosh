package main

import (
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
	// add_ -> add_dog
	//offset = 4
	newLine, offset := b.bellLine.Do(line, pos)
	currentString := string(line)
	//one match command
	if len(newLine) <= 1 {
		b.Reset()
		if len(newLine) == 0 {
			fmt.Printf("\x07")
		}
		return newLine, offset
	}
	//above one match command
	lcp := lcpMatch(newLine)
	// lcp = dog
	if len(lcp) > 0 {
		// there exist the shared part
		b.Reset()
		return [][]rune{lcp}, offset
	}
	if currentString == b.currentString {
		//two or more press tab button
		b.tabNum++
	} else {
		//initialize first tab
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
	var exec []string
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
			if info.Mode()&0111 != 0 {
				name := entry.Name()
				if !flag[name] {
					exec = append(exec, name)
					flag[name] = true
				}
			}
		}
	}
	return exec
}

func initExecList() *readline.PrefixCompleter {
	var items []readline.PrefixCompleterInterface
	command := make(map[string]bool)
	var addCod = func(name string) {
		if !command[name] {
			command[name] = true
			items = append(items, readline.PcItem(name))
		}
	}
	addCod("cd")
	addCod("exit")
	addCod("echo")
	addCod("type")
	addCod("pwd")

	execs := getAllPathExec()
	for _, exec := range execs {
		if !command[exec] {
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
