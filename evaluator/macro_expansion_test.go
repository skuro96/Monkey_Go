package evaluator

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func testParseProgram(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func TestDefineMacros(t *testing.T) {
	input := `
	let number = 1;
	let function = fn(x, y) { x + y; };
	let mymacro = macro(x, y) { x + y; };
	`

	env := object.NewEnvironment()
	program := testParseProgram(input)

	DefineMacros(program, env)

	if len(program.Statements) != 2 {
		t.Fatalf("Wrong number of statements. got=%d", len(program.Statements))
	}

	if _, ok := env.Get("number"); ok {
		t.Fatalf("number should not be defined")
	}
	if _, ok := env.Get("function"); ok {
		t.Fatalf("function should not be defined")
	}

	obj, ok := env.Get("mymacro")
	if !ok {
		t.Fatalf("macro not in environment")
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		t.Fatalf("object is not Macro. got=%T (%+v)", obj, obj)
	}

	if len(macro.Parameters) != 2 {
		t.Fatalf("Wrong number of macro parameters. got=%d", len(macro.Parameters))
	}

	if macro.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", macro.Parameters[0])
	}
	if macro.Parameters[1].String() != "y" {
		t.Fatalf("parameter is not 'y'. got=%q", macro.Parameters[1])
	}

	expectBody := "(x + y)"

	if macro.Body.String() != expectBody {
		t.Fatalf("body is not %q. got=%q", expectBody, macro.Body.String())
	}
}

func TestExpandMacros(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`
			let infixExpression = macro() { quote(1 + 2); };
			infixExpression();
			`,
			`(1 + 2)`,
		},
		{
			`
			let reverse = macro(a, b) { quote(unquote(b) - unquote(a)); };
			reverse(2 + 2, 10 - 5);
			`,
			`(10 - 5) - (2 + 2)`,
		},
	}

	for _, tt := range tests {
		expected := testParseProgram(tt.expected)
		program := testParseProgram(tt.input)

		env := object.NewEnvironment()
		DefineMacros(program, env)
		expanded := ExpandMacros(program, env)

		if expanded.String() != expected.String() {
			t.Errorf("not equal. want=%q, got=%q", expected.String(), expanded.String())
		}
	}
}
