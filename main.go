package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/drewhayward/glox/lox"
)

func run(source string) {
	// Tokenize the source string
	tokens, err := lox.ScanTokens(source)

	if err == nil {
		fmt.Printf("%v\n", tokens)
	}
}

func runFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
	print(string(data))

	run(string(data))
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	print("> ")
	for scanner.Scan() {
		line := scanner.Text()
		run(line)

		print("> ")
	}
}

func main() {
	if len(os.Args) == 1 {
		println("Usage: glox [script]")
		runPrompt()
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	}
}
