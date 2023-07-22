package interpreter

import (
	"fmt"
	"glox/lexer"
	"glox/parser"
	"math/rand"
	"strings"
	"testing"
)

func TestInterpret(t *testing.T) {

	x := rand.Float64()
	y := rand.Float64()
	fixtures := map[string]string{
		"1+1":                             "2",
		"!false":                          "true",
		"!true":                           "false",
		fmt.Sprintf("%v + %v", x, y):      fmt.Sprintf("%v", x+y),
		`"the number is"+1`:               "the number is1",
		`12+" is the number"`:             "12 is the number",
		`"hello world"`:                   "hello world",
		`(12+5*76/2)`:                     "202",
		"1==1.0000001":                    "false",
		`"yes" != "Yes"`:                  "true",
		`1>2?"1 is bigger":"2 is bigger"`: "2 is bigger",
	}

	for code, expected := range fixtures {
		lxr := lexer.New(code)
		prsr := parser.New(lxr.Tokenize())
		intrprtr := New()

		if expr, err := prsr.Parse(); err != nil {
			t.Fatalf("failed to parse code %q", code)
		} else {
			out := intrprtr.Interpret(expr)
			actual := fmt.Sprintf("%v", out)
			if actual != expected {
				t.Fatalf("failed to interpret %q. expected=%v got=%v", code, expected, actual)
			}
		}
	}

	failures := map[string][]string{
		"1+false": {
			"unsupported operands. This operation can only be performed with numbers and strings.",
			"+",
		},
		"1/0":        {"division by zero", "/"},
		"1-false":    {"Operator \"-\" only accepts number operands", "-"},
		"-\"yes\"":   {"Operator \"-\" only accepts number operands", "-"},
		"true*false": {"Operator \"*\" only accepts number operands", "*"},
		"false>true": {"Operator \">\" only accepts number operands", ">"},
		"false>=12":  {"Operator \">=\" only accepts number operands", ">="},
		"true<=12":   {"Operator \"<=\" only accepts number operands", "<="},
		"true<true":  {"Operator \"<\" only accepts number operands", "<"},
	}

	for code, chunks := range failures {
		lxr := lexer.New(code)
		prsr := parser.New(lxr.Tokenize())

		if expr, err := prsr.Parse(); err != nil {
			t.Fatalf("failed to parse code %q", code)
		} else {
			intrprtr := New()
			out := intrprtr.Interpret(expr)

			if err, isErr := out.(error); !isErr {
				t.Fatalf("failed to return error for=%q. expected=%q, got=%v", code, err, out)
			} else {
				actual := err.Error()
				errorMsgIsOk := true
				var notFound string
				for _, msg := range chunks {
					if !strings.Contains(actual, msg) {
						errorMsgIsOk = false
						notFound = msg
					}
				}
				if !errorMsgIsOk {
					t.Fatalf("%q => wrong error message. \nexpected contains=%q \ngot=%q", code, notFound, actual)
				}

			}
		}
	}

}
