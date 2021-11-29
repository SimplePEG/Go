package rd

import (
	"regexp"
	"strings"
)

type Expectation struct {
	TypeData string `type`
	Position int
	Children []Ast
	Rule     string
}

type Rule struct {
	Name   string
	Parser ParserFunc
}

type State struct {
	LastExpectations []Expectation
	Text             string
	Position         int
	Rules            []Rule
	SuccesfullRules  []Ast
	FailedRules      []Ast
}

type Ast struct {
	TypeData      string `type`
	Match         string
	Children      []Ast
	StartPosition int
	EndPosition   int
	Action        string
	Rule          string
}

//type Visitor interface {
//	visitAction(node *Ast) *Ast
//}
//
//func (node *Ast) Visit(v Visitor) *Ast {
//	for i := 0; i < len(node.Children); i++ {
//		ast := node.Children[i].Visit(v)
//		node.Children[i] = *ast
//	}
//
//	if node.Action != "" {
//		return v.visitAction(node)
//	}
//
//	return node
//}

type ParserFunc = func(state *State) (Ast, bool)

func GetLastError(state *State) (string, bool) {
	if len(state.LastExpectations) < 1 {
		return "", true
	}

	var last_exp_position int = state.Position
	for _, v := range state.LastExpectations {
		last_exp_position = max(last_exp_position, v.Position)
	}

	var dedupedExpectations []Expectation
	var lastExps []Expectation

	// filter last exps
	for i := 0; i < len(state.LastExpectations); i++ {
		if state.LastExpectations[i].Position == last_exp_position {
			lastExps = append(lastExps, state.LastExpectations[i])
		}
	}
	// get dedupedExpectations
	for i := 0; i < len(lastExps); i++ {
		result := true

		for j := 0; j < len(lastExps); j++ {
			if lastExps[j].Rule == lastExps[i].Rule && j != i {
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
		rules = append(rules, v.Rule)
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
		var textRule string
		if len(state.Text) >= state.Position+len(rule) {
			textRule = state.Text[state.Position : state.Position+len(rule)]
		}

		if textRule == rule {
			start := state.Position
			state.Position += len(rule)
			end := state.Position

			return Ast{
				TypeData:      "string",
				Match:         rule,
				StartPosition: start,
				EndPosition:   end,
			}, false
		}

		state.LastExpectations = []Expectation{
			Expectation{
				TypeData: "string",
				Rule:     rule,
				Position: state.Position,
			}}
		return Ast{}, true // return err
	}
}

func RegexChar(rule string) ParserFunc {
	return func(state *State) (Ast, bool) {
		text := state.Text[state.Position:]
		r, _ := regexp.Compile(rule)
		loc := r.FindStringIndex(text)

		if len(loc) > 0 && loc[0] == 0 {
			match := r.FindString(text)

			start := state.Position
			state.Position += len(match)
			end := state.Position

			return Ast{
				TypeData:      "regex_char",
				Match:         match,
				StartPosition: start,
				EndPosition:   end,
			}, false
		}

		state.LastExpectations = []Expectation{
			Expectation{
				TypeData: "regex_char",
				Rule:     rule,
				Position: state.Position,
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
			match += asts[i].Match
		}

		return Ast{
			TypeData:      "sequence",
			Match:         match,
			Children:      asts,
			StartPosition: startPosition,
			EndPosition:   state.Position,
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
					TypeData:      "ordered_choice",
					Match:         ast.Match,
					Children:      []Ast{ast},
					StartPosition: initialState.Position,
					EndPosition:   state.Position,
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
		//var hasAst = true
		var asts = []Ast{}
		var start_position = state.Position

		for {
			var state_position = state.Position

			ast, err := parser(state)

			if !err {
				asts = append(asts, ast)
			} else {
				state.Position = state_position
				break
			}
		}

		var match string

		for i := 0; i < len(asts); i++ {
			match += asts[i].Match
		}

		return Ast{
			TypeData:      "zero_or_more",
			Match:         match,
			Children:      asts,
			StartPosition: start_position,
			EndPosition:   state.Position,
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
				match += asts[i].Match
			}

			return Ast{
				TypeData:      "one_or_more",
				Match:         match,
				Children:      asts,
				StartPosition: start_position,
				EndPosition:   state.Position,
			}, false
		}

		return Ast{}, true
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
			match += asts[i].Match
		}

		return Ast{
			TypeData:      "optional",
			Match:         match,
			Children:      asts,
			StartPosition: start_position,
			EndPosition:   state.Position,
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
				TypeData:      "and_predicate",
				Children:      []Ast{ast},
				StartPosition: state.Position,
				EndPosition:   state.Position,
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
					TypeData: "not_predicate",
					Children: []Ast{ast},
					Position: state.Position,
				}}

			return Ast{}, true
		}

		state.LastExpectations = []Expectation{}

		return Ast{
			TypeData:      "not_predicate",
			StartPosition: state.Position,
			EndPosition:   state.Position,
		}, true
	}
}

func EndOfFile() ParserFunc {
	return func(state *State) (Ast, bool) {
		if len(state.Text) == state.Position {
			state.LastExpectations = []Expectation{}
			return Ast{
				TypeData:      "end_of_file",
				StartPosition: state.Position,
				EndPosition:   state.Position,
			}, false
		}
		state.LastExpectations = []Expectation{
			Expectation{
				TypeData: "end_of_file",
				Rule:     "EOF",
				Position: state.Position,
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
			ast.Action = name
		}

		return ast, err
	}
}

func CallRuleByName(name string) ParserFunc {
	return func(state *State) (Ast, bool) {
		var rule Rule
		for _, r := range state.Rules {
			if r.Name == name {
				rule = r
				break
			}
		}

		return rule.Parser(state)
	}
}
