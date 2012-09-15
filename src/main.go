package main

import (
	"os"
	"bufio"
	"./calc"
)

func main() {
	bufin := bufio.NewReader(os.Stdin)
	repl := calc.NewRepl(bufin, os.Stdout, os.Stderr)
	repl.ReadDefault()
	return
}