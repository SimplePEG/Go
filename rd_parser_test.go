package simplepeg

import (
	"reflect"
	"testing"
)

func TestSequence(t *testing.T) {
	// arrange
	var text = "testsomething"
	var string_to_match = "testsomething"
	var start_position = 0

	// act
	var parser = Sequence([]ParserFunc{
		String("test"),
		String("something"),
	})

	var ast, err = parser(&State{
		text:     text,
		position: start_position,
	})

	if err {
		t.Error("Sequence not parsed")
	}

	if ast.match != string_to_match {
		t.Error("Expected ", string_to_match, " but val is ", ast.match)
	}

	if ast.start_position != start_position {
		t.Error("Expected ", start_position, " but val is ", ast.start_position)
	}

	if ast.end_position != len(string_to_match) {
		t.Error("Expected ", len(string_to_match), " but val is ", ast.end_position)
	}

	var children_to_match = []Ast{
		Ast{
			typeData:       "string",
			match:          "test",
			start_position: 0,
			end_position:   4,
		},
		Ast{
			typeData:       "string",
			match:          "something",
			start_position: 4,
			end_position:   13,
		},
	}

	if reflect.DeepEqual(ast.children, children_to_match) == false {
		t.Error("Children asts is not equals")
	}
}
