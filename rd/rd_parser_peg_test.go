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

	// we need to assert without Children
	if len(state.LastExpectations) != 0 {
		t.Error("Not parsed")
	}

	var mock = Ast{
		Match:         string_to_match,
		StartPosition: start_position,
		EndPosition:   len(string_to_match),
		TypeData:      "string",
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
			Rule:     rule,
			Position: start_position,
			TypeData: "string",
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

	if ast.Match != stringToMatch || ast.TypeData != "regex_char" || ast.StartPosition != startPosition || ast.EndPosition != len(stringToMatch) {
		t.Error("Not parsed implement")
	}
}

// should implement "regex_char" method by string
func TestRegexCharString(t *testing.T) {
	// arrange
	var text = " "
	var stringToMatch = " "
	var rule = "[\\s]"
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

	println(state.Position)

	if ast.Match != stringToMatch || ast.TypeData != "regex_char" || ast.StartPosition != startPosition || ast.EndPosition != len(stringToMatch) {
		t.Error("Not parsed implement")
	}
}

// should implement "regex_char" method by string failed
func TestRegexCharStringWrong(t *testing.T) {
	// arrange
	var text = "1  1"
	//var stringToMatch = " "
	var rule = "[\\s]"
	var startPosition = 0
	// act
	var parser = RegexChar(rule)
	var state = &State{
		Text:     text,
		Position: startPosition,
	}
	var _, err = parser(state)

	if err == false {
		t.Error("Should be not parsed")
	}
}

// should correctly handle wrong text for "regex_char" method
func TestRegexCharWrong(t *testing.T) {
	// arrange
	var text = "asda"
	//var stringToMatch = "8"
	var startPosition = 0
	var rule = "[0-9]"

	// act
	var parser = RegexChar(rule)
	var state = &State{
		Text:     text,
		Position: startPosition,
	}
	var _, err = parser(state)

	if err == false {
		t.Error("Parsed, by should be not parse")
	}

	var exp = state.LastExpectations[0]

	if exp.Rule == rule && exp.Position == startPosition && exp.TypeData == "regex_char" {
	} else {
		t.Error("Not parsed implement")
	}
}

// should implement "sequence" method for string/regex
func TestSequenceRegexString(t *testing.T) {
	// arrange
	var text = "test2"
	var stringToMatch = "test2"
	var children_to_match = []Ast{
		{
			Match:         "test",
			StartPosition: 0,
			EndPosition:   4,
			TypeData:      "string",
		},
		{
			Match:         "2",
			StartPosition: 4,
			EndPosition:   5,
			TypeData:      "regex_char",
		},
	}

	var parsing_expressions = []ParserFunc{
		String("test"),
		RegexChar("[0-9]"),
	}

	var start_position = 0
	var end_position = len(text)

	var parser = Sequence(parsing_expressions)

	var state = State{
		Text:     text,
		Position: start_position,
	}

	var ast, err = parser(&state)

	if err == true {
		t.Error("Should be parse")
	}

	if len(state.LastExpectations) != 0 {
		t.Error("LastExpectations, should bu empty")
	}

	var mock = Ast{
		Match:         stringToMatch,
		Children:      children_to_match,
		StartPosition: start_position,
		EndPosition:   end_position,
		TypeData:      "sequence",
	}

	if reflect.DeepEqual(ast, mock) == false {
		t.Error("asts is not equals")
	}
}

// should implement "zero_or_more" method for string
func TestZeroOrMore(t *testing.T) {
	// arrange
	var text = "22"
	var stringToMatch = "22"
	var children_to_match = []Ast{
		{
			Match:         "2",
			StartPosition: 0,
			EndPosition:   1,
			TypeData:      "regex_char",
		},
		{
			Match:         "2",
			StartPosition: 1,
			EndPosition:   2,
			TypeData:      "regex_char",
		},
	}

	var parsing_expression = RegexChar("[0-9]")

	var start_position = 0
	var end_position = len(text)

	var parser = ZeroOrMore(parsing_expression)

	var state = State{
		Text:     text,
		Position: start_position,
	}

	var ast, err = parser(&state)

	if err == true {
		t.Error("Should be parse")
	}

	var mockLastE = []Expectation{
		{
			Position: 2,
			Rule:     "[0-9]",
			TypeData: "regex_char",
		},
	}

	if reflect.DeepEqual(state.LastExpectations, mockLastE) == false {
		t.Error("LastExpectations, should bu empty")
	}

	var mock = Ast{
		Match:         stringToMatch,
		Children:      children_to_match,
		StartPosition: start_position,
		EndPosition:   end_position,
		TypeData:      "zero_or_more",
	}

	if reflect.DeepEqual(ast, mock) == false {
		t.Error("asts is not equals")
	}
}
