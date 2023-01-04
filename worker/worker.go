package worker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Result struct {
	Line    string
	LineNum int
	Path    string
}

type Results struct {
	Inner []Result
}

func NewResult(line string, num int, path string) Result {
	return Result{Line: line, LineNum: num, Path: path}
}

func FindInFile(filePath string, search string) *Results {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	results := Results{make([]Result, 0)}
	scanner := bufio.NewScanner(file)

	lineNum := 1
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, search) {
			results.Inner = append(results.Inner, NewResult(text, lineNum, filePath))
		}
		lineNum += 1
	}

	if len(results.Inner) == 0 {
		return nil
	}
	return &results
}
