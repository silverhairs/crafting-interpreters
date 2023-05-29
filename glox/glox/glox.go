package glox

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const PROMPT = "> "

func RunFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	bytes := make([]byte, info.Size())
	if _, err = bufio.NewReader(file).Read(bytes); err != nil {
		return err
	}
	run(string(bytes))
	return nil
}

func StartREPL(r io.Reader) {
	scanner := bufio.NewScanner(r)

	for {
		fmt.Print(PROMPT)
		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		run(line)
	}
}

func run(src string) {
	lxr := lexer.New(src)

	tokens := lxr.Tokenize()
	for _, tok := range tokens {
		fmt.Println(tok)
	}
}
