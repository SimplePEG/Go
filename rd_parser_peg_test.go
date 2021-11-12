package simplepeg

import "testing"

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
		text:     text,
		position: startPosition,
	}
	var ast, err = parser(state)

	if err == true {
		t.Error("Not parsed")
	}

	if ast.match != stringToMatch || ast.typeData != "regex_char" || ast.start_position != startPosition || ast.end_position != len(stringToMatch) {
		t.Error("Not parsed implement")
	}
}
