package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/drewhayward/glox/lox"
)

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
	rs := lox.NewRuntimeState()
	rs.Run(string(data))
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	rs := lox.NewRuntimeState()
	print("> ")
	for scanner.Scan() {
		line := scanner.Text()
		rs.Run(line)

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
