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
	LastExpectations []Expectation
	Text             string
	Position         int
	Rules            []Rule
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

func GetLastError(state *State) (string, bool) {
	if len(state.LastExpectations) < 1 {
		return "", true
	}

	var last_exp_position int = state.Position
	for _, v := range state.LastExpectations {
		last_exp_position = max(last_exp_position, v.position)
	}

	var dedupedExpectations []Expectation
	var lastExps []Expectation

	// filter last exps
	for i := 0; i < len(state.LastExpectations); i++ {
		if state.LastExpectations[i].position == last_exp_position {
			lastExps = append(lastExps, state.LastExpectations[i])
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

	lines := strings.Split(state.Text, "\n")

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

	if last_exp_position < len(state.Text) {
		unexpected_char = string(state.Text[last_exp_position])
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
	return unexpected + expected + "\n" + str_error_ln + ": " + extra, false

}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func String(rule string) ParserFunc {
	return func(state *State) (Ast, bool) {

		if state.Text[state.Position:state.Position+len(rule)] == rule {
			start := state.Position
			state.Position += len(rule)
			end := state.Position

			return Ast{
				typeData:       "string",
				match:          rule,
				start_position: start,
				end_position:   end,
			}, false
		}

		state.LastExpectations = []Expectation{
			Expectation{
				typeData: "string",
				rule:     rule,
				position: state.Position,
			}}
		return Ast{}, true // return err
	}
}

func RegexChar(rule string) ParserFunc {
	return func(state *State) (Ast, bool) {
		text := state.Text[state.Position:]
		isMatch, _ := regexp.MatchString(rule, text)

		if isMatch {
			r, _ := regexp.Compile(rule)
			match := r.FindString(text)

			start := state.Position
			state.Position += len(match)
			end := state.Position

			return Ast{
				typeData:       "regex_char",
				match:          match,
				start_position: start,
				end_position:   end,
			}, false
		}

		state.LastExpectations = []Expectation{
			Expectation{
				typeData: "regex_char",
				rule:     rule,
				position: state.Position,
			}}

		return Ast{}, true // return err
	}
}

func Sequence(parsers []ParserFunc) ParserFunc {
	return func(state *State) (Ast, bool) {
		var asts []Ast // Ast
		var expectations []Expectation
		var startPosition = state.Position

		for i := 0; i < len(parsers); i++ {
			var ast, err = parsers[i](state)
			expectations = append(expectations, state.LastExpectations...)

			if !err {
				asts = append(asts, ast)
			} else {
				state.LastExpectations = expectations
				return Ast{}, true
			}
		}
		state.LastExpectations = expectations

		var match = ""

		for i := 0; i < len(asts); i++ {
			match += asts[i].match
		}

		return Ast{
			typeData:       "sequence",
			match:          match,
			children:       asts,
			start_position: startPosition,
			end_position:   state.Position,
		}, false

	}
}

func OrderedChoice(parsers []ParserFunc) ParserFunc {
	return func(state *State) (Ast, bool) {
		var expectations []Expectation
		var initialState = State{
			Text:     state.Text,
			Position: state.Position,
		}

		for i := 0; i < len(parsers); i++ {
			var ast, err = parsers[i](state)

			if !err {
				return Ast{
					typeData:       "ordered_choice",
					match:          ast.match,
					children:       []Ast{ast},
					start_position: initialState.Position,
					end_position:   state.Position,
				}, false
			}
			state.Text = initialState.Text
			state.Position = initialState.Position
			expectations = append(expectations, state.LastExpectations...)
		}

		state.LastExpectations = expectations
		return Ast{}, true
	}
}

func ZeroOrMore(parser ParserFunc) ParserFunc {
	return func(state *State) (Ast, bool) {
		var hasAst = true
		var asts = []Ast{}
		var start_position = state.Position

		for hasAst {
			var state_position = state.Position
			ast, err := parser(state)
			hasAst = !err

			if !err {
				asts = append(asts, ast)
			} else {
				state.Position = state_position
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
			end_position:   state.Position,
		}, false
	}
}

func OneOrMore(parser ParserFunc) ParserFunc {
	return func(state *State) (Ast, bool) {
		var hasAst = true
		var asts = []Ast{}
		var start_position = state.Position

		for hasAst {
			var state_position = state.Position
			ast, err := parser(state)
			hasAst = !err

			if !err {
				asts = append(asts, ast)
			} else {
				state.Position = state_position
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
				end_position:   state.Position,
			}, false
		}

		return Ast{}, false
	}
}

func Optional(parser ParserFunc) ParserFunc {
	return func(state *State) (Ast, bool) {
		var start_position = state.Position
		var ast, err = parser(state)
		var asts = []Ast{}

		if !err {
			asts = append(asts, ast)
		} else {
			state.Position = start_position
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
			end_position:   state.Position,
		}, false
	}
}

func AndPredicate(parser ParserFunc) ParserFunc {
	return func(state *State) (Ast, bool) {
		var currentState = State{
			Text:     state.Text,
			Position: state.Position,
		}

		ast, err := parser(state)

		if !err {
			state.Text = currentState.Text
			state.Position = currentState.Position

			return Ast{
				typeData:       "and_predicate",
				children:       []Ast{ast},
				start_position: state.Position,
				end_position:   state.Position,
			}, false
		}

		return Ast{}, true
	}
}

func NotPredicate(parser ParserFunc) ParserFunc {
	return func(state *State) (Ast, bool) {
		var currentState = State{
			Text:     state.Text,
			Position: state.Position,
		}

		ast, err := parser(state)

		if !err {
			state.Text = currentState.Text
			state.Position = currentState.Position
			state.LastExpectations = []Expectation{
				Expectation{
					typeData: "not_predicate",
					children: []Ast{ast},
					position: state.Position,
				}}

			return Ast{}, true
		}

		state.LastExpectations = []Expectation{}

		return Ast{
			typeData:       "not_predicate",
			children:       []Ast{},
			start_position: state.Position,
			end_position:   state.Position,
		}, true
	}
}

func EndOfFile() ParserFunc {
	return func(state *State) (Ast, bool) {
		if len(state.Text) == state.Position {
			state.LastExpectations = []Expectation{}
			return Ast{
				typeData:       "end_of_file",
				children:       []Ast{},
				start_position: state.Position,
				end_position:   state.Position,
			}, false
		}
		state.LastExpectations = []Expectation{
			Expectation{
				typeData: "end_of_file",
				rule:     "EOF",
				position: state.Position,
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
		for _, r := range state.Rules {
			if r.name == name {
				rule = r
				break
			}
		}
		return rule.parser(state)
	}
}
