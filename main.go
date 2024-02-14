package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/drewhayward/glox/lox"
)

func run(source string) {
	// Tokenize the source string
	tokens, err := lox.ScanTokens(source)
	if err != nil {
		fmt.Printf("Lexing Error %v\n", tokens)
		return
	}

    root, err := lox.Parse(tokens)
	if err != nil {
		fmt.Println(err)
		return
	}

    // Todo, init some kind of runtime environment
    // state to track variables and function calls
    pStmts := root.(lox.ProgramNode)
    for _, stmt := range pStmts.Statements {
        lox.Interpret(stmt)
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
    flag.Bool("v", false, "Verbose parsing and lexing")
    flag.Parse()

	if len(os.Args) == 1 {
		fmt.Println("Usage: glox [script]")
		runPrompt()
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	}
}
