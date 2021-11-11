package simplepeg

import (
	"regexp"
	"strings"
)

type Expectation struct {
	typeData string `type`
	position int
	rule     string
}

type State struct {
	lastExpectations []Expectation
	text             string
	position         int
}

type Ast struct {
	typeData       string `type`
	match          string
	children       []Ast
	start_position int
	end_position   int
}

type ParserFunc = func(state *State) (Ast, bool)

func GetLastError(state *State) interface{} {
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

func String(rule string) ParserFunc {
	return func(state *State) (Ast, bool) {

		if state.text[state.position:state.position+len(rule)] == rule {
			start := state.position
			state.position += len(rule)
			end := state.position

			return Ast{
				typeData:       "string",
				match:          rule,
				start_position: start,
				end_position:   end,
			}, false
		}

		state.lastExpectations = []Expectation{
			Expectation{
				typeData: "string",
				rule:     rule,
				position: state.position,
			}}
		return Ast{}, true // return err
	}
}

func RegexChar(rule string) ParserFunc {
	return func(state *State) (Ast, bool) {
		text := state.text[state.position:]
		isMatch, _ := regexp.MatchString(rule, text)

		if isMatch {
			r, _ := regexp.Compile(rule)
			match := r.FindString(text)

			start := state.position
			state.position += len(rule)
			end := state.position

			return Ast{
				typeData:       "regex_char",
				match:          match,
				start_position: start,
				end_position:   end,
			}, false
		}

		state.lastExpectations = []Expectation{
			Expectation{
				typeData: "regex_char",
				rule:     rule,
				position: state.position,
			}}

		return Ast{}, true // return err
	}
}

func Sequence(parsers []ParserFunc) ParserFunc {
	return func(state *State) (Ast, bool) {
		var asts []Ast // Ast
		var expectations []Expectation
		var startPosition = state.position

		for i := 0; i < len(parsers); i++ {
			var ast, err = parsers[i](state)
			expectations = append(expectations, state.lastExpectations...)

			if !err {
				asts = append(asts, ast)
			} else {
				state.lastExpectations = expectations
				return Ast{}, true
			}
		}
		state.lastExpectations = expectations

		var match = ""

		for i := 0; i < len(asts); i++ {
			match += asts[i].match
		}

		return Ast{
			typeData:       "sequence",
			match:          match,
			children:       asts,
			start_position: startPosition,
			end_position:   state.position,
		}, false

	}
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
