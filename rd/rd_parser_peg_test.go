package rd

import (
	"reflect"
	"testing"
)

// should implement "string" method'
func TestString(t *testing.T) {
	// arrange
	var text = "test"
	var string_to_match = "test"
	var start_position = 0
	var rule = "test"
	// act
	var parser = String(rule)
	var state = &State{
		Text:     text,
		Position: start_position,
	}
	var ast, err = parser(state)

	if err == true {
		t.Error("Not parsed")
	}

	// we need to assert without children
	if len(state.LastExpectations) != 0 {
		t.Error("Not parsed")
	}

	var mock = Ast{
		match:          string_to_match,
		start_position: start_position,
		end_position:   len(string_to_match),
		typeData:       "string",
	}

	if reflect.DeepEqual(ast, mock) == false {
		t.Error("asts is not equals")
	}
}

// should correctly handle wrong Text for "string" method
func TestStringCorrect(t *testing.T) {
	// arrange
	var text = "asda"
	//var string_to_match = "test"
	var start_position = 0
	var rule = "test"
	// act
	var parser = String(rule)
	var state = &State{
		Text:     text,
		Position: start_position,
	}
	var _, err = parser(state)

	if err == false {
		t.Error("Parsed, by should be not parse")
	}

	var mock = []Expectation{
		{
			rule:     rule,
			position: start_position,
			typeData: "string",
		},
	}

	if reflect.DeepEqual(state.LastExpectations, mock) == false {
		t.Error("asts is not equals")
	}
}

// should implement "regex_char" method
func TestRegexChar(t *testing.T) {
	// arrange
	var text = "8"
	var stringToMatch = "8"
	var rule = "[0-9]"
	var startPosition = 0
	// act
	var parser = RegexChar(rule)
	var state = &State{
		Text:     text,
		Position: startPosition,
	}
	var ast, err = parser(state)

	if err == true {
		t.Error("Not parsed")
	}

	if ast.match != stringToMatch || ast.typeData != "regex_char" || ast.start_position != startPosition || ast.end_position != len(stringToMatch) {
		t.Error("Not parsed implement")
	}
}
