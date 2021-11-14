package rd

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
		Text:     text,
		Position: startPosition,
	})

	if err {
		t.Error("Sequence not parsed")
	}

	if ast.Match != stringToMatch {
		t.Error("Expected ", stringToMatch, " but val is ", ast.Match)
	}

	if ast.StartPosition != startPosition {
		t.Error("Expected ", startPosition, " but val is ", ast.StartPosition)
	}

	if ast.EndPosition != len(stringToMatch) {
		t.Error("Expected ", len(stringToMatch), " but val is ", ast.EndPosition)
	}

	var childrenToMath = []Ast{
		{
			TypeData:      "string",
			Match:         "test",
			StartPosition: 0,
			EndPosition:   4,
		},
		{
			TypeData:      "string",
			Match:         "something",
			StartPosition: 4,
			EndPosition:   13,
		},
	}

	if reflect.DeepEqual(ast.Children, childrenToMath) == false {
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
		Text:     text,
		Position: startPosition,
	})

	if err == true {
		t.Error("should successfully parse simple math expression")
	}

	if ast.Match != stringToMatch {
		t.Error("should successfully parse simple math expression")
	}
}

// should fail and return to Position if ordered_choice failed
func TestOrderedChoiseFailed(t *testing.T) {
	// arrange
	var text = "AB"
	var startPosition = 0
	var state = &State{
		Text:     text,
		Position: startPosition,
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

	if err == false || state.Position != 0 {
		t.Error("should fail and return to Position if ordered_choice failed")
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
			EndPosition:   1,
			Match:         "1",
			StartPosition: 0,
			TypeData:      "regex_char",
		},
		{
			EndPosition:   2,
			Match:         "D",
			StartPosition: 1,
			TypeData:      "regex_char",
		},
	}

	// act
	var parser = Sequence([]ParserFunc{
		RegexChar("[0-9]"),
		RegexChar("[A-Z]"),
	})
	var ast, _ = parser(&State{
		Text:     text,
		Position: startPosition,
	})

	mockAst := Ast{
		Match:         stringToMatch,
		Children:      childrenToMatch,
		StartPosition: startPosition,
		EndPosition:   len(stringToMatch),
		TypeData:      "sequence",
	}

	if reflect.DeepEqual(ast, mockAst) == false {
		t.Error("should successfully combine 'regex_char' methods with 'sequence'")
	}
}
