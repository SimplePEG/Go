package simplepeg

import (
	"strings"
)

type Expectation struct {
	position int
	rule     string
}

type State struct {
	lastExpectations []Expectation
	text             string
	position         int
}

func GetLastError(state State) interface{} {
	if len(state.lastExpectations) < 1 {
		return false
	}

	var last_exp_position int = state.position
	for _, v := range state.lastExpectations {
		last_exp_position = max(last_exp_position, v.position)
	}

	var dedupedExpectations []Expectation
	var lastExps []Expectation

	// filter last exps
	for i := 0; i < len(state.lastExpectations); i++ {
		if state.lastExpectations[i].position == last_exp_position {
			lastExps = append(lastExps, state.lastExpectations[i])
		}
	}
	// get dedupedExpectations
	for i := 0; i < len(lastExps); i++ {
		result := true

		for j := 0; j < len(lastExps); j++ {
			if lastExps[j].rule == lastExps[i].rule && j != i {
				result = false
				break
			}
		}

		if result == true {
			dedupedExpectations = append(dedupedExpectations, lastExps[i])
		}
	}

	var last_position = 0
	var line_of_error = ""
	var error_line_number int
	var position_of_error = 0

	lines := strings.Split(state.text, "\n")

	for i := 0; i < len(lines); i++ {
		line_lenght := len(lines[i]) + 1

		if last_exp_position >= last_position &&
			last_exp_position < (last_position+line_lenght) {
			line_of_error = lines[i]
			position_of_error = last_exp_position - last_position
			error_line_number = i + 1
			break
		}

		last_position += line_lenght
	}

	var str_error_ln = string(rune(error_line_number))
	var error_ln_length = len(str_error_ln)
	var unexpected_char = "EOF"

	if last_exp_position < len(state.text) {
		unexpected_char = string(state.text[last_exp_position])
	}

	var unexpected = "Unexpected '" + unexpected_char + "'"
	var rules []string

	for _, v := range dedupedExpectations {
		rules = append(rules, v.rule)
	}

	var expected = " expected (" + strings.Join(rules[:], ",") + ")"
	var pointer = ""

	pointer_count := position_of_error + 3 + error_ln_length

	for i := 0; i < pointer_count; i++ {
		pointer += "-"
	}
	pointer += "^"

	var extra = line_of_error + "\n" + pointer
	return unexpected + expected + "\n" + str_error_ln + ": " + extra

}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
