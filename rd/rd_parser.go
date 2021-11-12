package rd

import (
	"regexp"
	"strings"
)

type Expectation struct {
	typeData string `type`
	position int
	children []Ast
	rule     string
}

type Rule struct {
	name   string
	parser ParserFunc
}

type State struct {
	lastExpectations []Expectation
	text             string
	position         int
	rules            []Rule
}

type Ast struct {
	typeData       string `type`
	match          string
	children       []Ast
	start_position int
	end_position   int
	action         string
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

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
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
			state.position += len(match)
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

func OrderedChoice(parsers []ParserFunc) ParserFunc {
	return func(state *State) (Ast, bool) {
		var expectations []Expectation
		var initialState = State{
			text:     state.text,
			position: state.position,
		}

		for i := 0; i < len(parsers); i++ {
			var ast, err = parsers[i](state)

			if !err {
				return Ast{
					typeData:       "ordered_choice",
					match:          ast.match,
					children:       []Ast{ast},
					start_position: initialState.position,
					end_position:   state.position,
				}, false
			}
			state.text = initialState.text
			state.position = initialState.position
			expectations = append(expectations, state.lastExpectations...)
		}

		state.lastExpectations = expectations
		return Ast{}, true
	}
}

func ZeroOrMore(parser ParserFunc) ParserFunc {
	return func(state *State) (Ast, bool) {
		var hasAst = true
		var asts = []Ast{}
		var start_position = state.position

		for hasAst {
			var state_position = state.position
			ast, err := parser(state)
			hasAst = !err

			if !err {
				asts = append(asts, ast)
			} else {
				state.position = state_position
			}
		}

		var match string

		for i := 0; i < len(asts); i++ {
			match += asts[i].match
		}

		return Ast{
			typeData:       "zero_or_more",
			match:          match,
			children:       asts,
			start_position: start_position,
			end_position:   state.position,
		}, false
	}
}

func OneOrMore(parser ParserFunc) ParserFunc {
	return func(state *State) (Ast, bool) {
		var hasAst = true
		var asts = []Ast{}
		var start_position = state.position

		for hasAst {
			var state_position = state.position
			ast, err := parser(state)
			hasAst = !err

			if !err {
				asts = append(asts, ast)
			} else {
				state.position = state_position
			}
		}

		if len(asts) > 0 {

			var match string

			for i := 0; i < len(asts); i++ {
				match += asts[i].match
			}

			return Ast{
				typeData:       "one_or_more",
				match:          match,
				children:       asts,
				start_position: start_position,
				end_position:   state.position,
			}, false
		}

		return Ast{}, false
	}
}

func Optional(parser ParserFunc) ParserFunc {
	return func(state *State) (Ast, bool) {
		var start_position = state.position
		var ast, err = parser(state)
		var asts = []Ast{}

		if !err {
			asts = append(asts, ast)
		} else {
			state.position = start_position
		}

		var match string

		for i := 0; i < len(asts); i++ {
			match += asts[i].match
		}

		return Ast{
			typeData:       "optional",
			match:          match,
			children:       asts,
			start_position: start_position,
			end_position:   state.position,
		}, false
	}
}

func AndPredicate(parser ParserFunc) ParserFunc {
	return func(state *State) (Ast, bool) {
		var currentState = State{
			text:     state.text,
			position: state.position,
		}

		ast, err := parser(state)

		if !err {
			state.text = currentState.text
			state.position = currentState.position

			return Ast{
				typeData:       "and_predicate",
				children:       []Ast{ast},
				start_position: state.position,
				end_position:   state.position,
			}, false
		}

		return Ast{}, true
	}
}

func NotPredicate(parser ParserFunc) ParserFunc {
	return func(state *State) (Ast, bool) {
		var currentState = State{
			text:     state.text,
			position: state.position,
		}

		ast, err := parser(state)

		if !err {
			state.text = currentState.text
			state.position = currentState.position
			state.lastExpectations = []Expectation{
				Expectation{
					typeData: "not_predicate",
					children: []Ast{ast},
					position: state.position,
				}}

			return Ast{}, true
		}

		state.lastExpectations = []Expectation{}

		return Ast{
			typeData:       "not_predicate",
			children:       []Ast{},
			start_position: state.position,
			end_position:   state.position,
		}, true
	}
}

func EndOfFile() ParserFunc {
	return func(state *State) (Ast, bool) {
		if len(state.text) == state.position {
			state.lastExpectations = []Expectation{}
			return Ast{
				typeData:       "end_of_file",
				children:       []Ast{},
				start_position: state.position,
				end_position:   state.position,
			}, false
		}
		state.lastExpectations = []Expectation{
			Expectation{
				typeData: "end_of_file",
				rule:     "EOF",
				position: state.position,
			},
		}

		return Ast{}, false
	}
}

func Rec(f func() ParserFunc) ParserFunc {
	return func(state *State) (Ast, bool) {
		var ast, err = f()(state)
		return ast, err
	}
}

func Action(name string, parser ParserFunc) ParserFunc {
	return func(state *State) (Ast, bool) {
		ast, err := parser(state)

		if !err {
			ast.action = name
		}

		return ast, false
	}
}

func CallRuleByName(name string) ParserFunc {
	return func(state *State) (Ast, bool) {
		var rule Rule
		for _, r := range state.rules {
			if r.name == name {
				rule = r
				break
			}
		}
		return rule.parser(state)
	}
}
