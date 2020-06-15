package evaluator

import (
	"testing"
	"yac/lexer"
	"yac/parser"
)

func TestEval(t *testing.T) {
	var l = lexer.New("16 + 4")
	var p = parser.New(l)
	var tree = p.Parse()
	if len(tree) != 1 {
		t.Errorf("Multiple nodes returned: %d", len(tree))
	}
	var eval = Eval(tree[0])
	if eval.Inspect() != "20" {
		t.Errorf("Expected Val to be 20, got: %s", eval.Inspect())
	}
}
