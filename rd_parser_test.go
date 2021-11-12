package simplepeg

import (
	"reflect"
	"testing"
)

// should successfully combine "string" methods with "sequence"'
func TestSequence(t *testing.T) {
	// arrange
	var text = "testsomething"
	var stringToMatch = "testsomething"
	var startPosition = 0

	// act
	var parser = Sequence([]ParserFunc{
		String("test"),
		String("something"),
	})

	var ast, err = parser(&State{
		text:     text,
		position: startPosition,
	})

	if err {
		t.Error("Sequence not parsed")
	}

	if ast.match != stringToMatch {
		t.Error("Expected ", stringToMatch, " but val is ", ast.match)
	}

	if ast.start_position != startPosition {
		t.Error("Expected ", startPosition, " but val is ", ast.start_position)
	}

	if ast.end_position != len(stringToMatch) {
		t.Error("Expected ", len(stringToMatch), " but val is ", ast.end_position)
	}

	var childrenToMath = []Ast{
		{
			typeData:       "string",
			match:          "test",
			start_position: 0,
			end_position:   4,
		},
		{
			typeData:       "string",
			match:          "something",
			start_position: 4,
			end_position:   13,
		},
	}

	if reflect.DeepEqual(ast.children, childrenToMath) == false {
		t.Error("Children asts is not equals")
	}
}

// should successfully parse simple math expression
func TestSimpleMath(t *testing.T) {
	// arrange
	var text = "1 + 2"
	var stringToMatch = "1 + 2"
	var startPosition = 0
	// var children_to_match []Ast

	var space = func() ParserFunc {
		return RegexChar("[\\s]")
	}
	var multiplicative = func() ParserFunc {
		return String("*")
	}
	var additive = func() ParserFunc {
		return String("+")
	}

	var exp func() ParserFunc

	var factor = func() ParserFunc {
		return OrderedChoice([]ParserFunc{
			Sequence([]ParserFunc{
				String("("),
				Rec(exp),
				String(")"),
			}),
			RegexChar("[0-9]"),
		})
	}

	var term = func() ParserFunc {
		return Sequence([]ParserFunc{
			factor(),
			ZeroOrMore(
				Sequence([]ParserFunc{
					space(),
					multiplicative(),
					space(),
					factor(),
				}),
			),
		})
	}

	exp = func() ParserFunc {
		return Sequence([]ParserFunc{
			term(),
			ZeroOrMore(
				Sequence([]ParserFunc{
					space(),
					additive(),
					space(),
					term(),
				}),
			),
		})
	}

	var math = func() ParserFunc {
		return Sequence([]ParserFunc{
			exp(),
			EndOfFile(),
		})

	}
	var parser = math()
	var ast, err = parser(&State{
		text:     text,
		position: startPosition,
	})

	if err == true {
		t.Error("should successfully parse simple math expression")
	}

	if ast.match != stringToMatch {
		t.Error("should successfully parse simple math expression")
	}
}

// should fail and return to position if ordered_choice failed
func TestOrderedChoiseFailed(t *testing.T) {
	// arrange
	var text = "AB"
	var startPosition = 0
	var state = &State{
		text:     text,
		position: startPosition,
	}
	// act
	var parser = OrderedChoice([]ParserFunc{
		Sequence([]ParserFunc{
			String("A"),
			String("A"),
		}),
		String("B"),
	})

	var _, err = parser(state)

	if err == false || state.position != 0 {
		t.Error("should fail and return to position if ordered_choice failed")
	}
}

// should successfully combine "regex_char" methods with "sequence"
func TestComboRegexCharSequence(t *testing.T) {
	// arrange
	var text = "1D"
	var stringToMatch = "1D"
	var startPosition = 0
	var childrenToMatch = []Ast{
		{
			end_position:   1,
			match:          "1",
			start_position: 0,
			typeData:       "regex_char",
		},
		{
			end_position:   2,
			match:          "D",
			start_position: 1,
			typeData:       "regex_char",
		},
	}

	// act
	var parser = Sequence([]ParserFunc{
		RegexChar("[0-9]"),
		RegexChar("[A-Z]"),
	})
	var ast, _ = parser(&State{
		text:     text,
		position: startPosition,
	})

	mockAst := Ast{
		match:          stringToMatch,
		children:       childrenToMatch,
		start_position: startPosition,
		end_position:   len(stringToMatch),
		typeData:       "sequence",
	}

	if reflect.DeepEqual(ast, mockAst) == false {
		t.Error("should successfully combine 'regex_char' methods with 'sequence'")
	}
}
