package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/larzconwell/micro/parser"
	"github.com/larzconwell/micro/scanner"
)

const (
	maxErrors = 5
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "Usage: micro <path|->")
		os.Exit(2)
	}

	err := run(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(path string) error {
	var err error
	var file *os.File

	if path == "-" {
		path = "stdin"
		file = os.Stdin
	} else {
		file, err = os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	reader := bufio.NewReader(file)

	scanner := scanner.New(reader, maxErrors)
	tokens, err := scanner.Scan()
	if err != nil {
		return err
	}

	parser := parser.New(tokens, maxErrors)
	err = parser.Parse()
	if err != nil {
		return err
	}

	return nil
}
