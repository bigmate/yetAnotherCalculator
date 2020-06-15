package repl

import (
	"bufio"
	"fmt"
	"io"
	"yac/evaluator"
	"yac/lexer"
	"yac/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	var scanner = bufio.NewScanner(in)
	for {
		fmt.Print(PROMPT)
		if !scanner.Scan() {
			return
		}
		var line = scanner.Text()
		var l = lexer.New(line)
		var p = parser.New(l)
		var tree = p.Parse()
		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		if len(tree) != 1 {
			io.WriteString(out, "Multiple nodes returned")
		}
		var evaluated = evaluator.Eval(tree[0])
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Syntax Error:\n")
	for i := range errors {
		io.WriteString(out, "\t-" + errors[i] + "\n")
	}
}