package speg

import (
	"github.com/SimplePEG/Go/rd"
	"strings"
)

type GrammarRules struct {
	rules []rd.Rule
}

type NodeVisit struct {
	Parsers  []rd.ParserFunc
	Node     *rd.Ast
	Children []NodeVisit
}

func actionVisit(visitNode *NodeVisit) (rd.ParserFunc, *NodeVisit) {
	if len(visitNode.Node.Children) > 0 {
		for i := 0; i < len(visitNode.Node.Children); i++ {
			child := visitNode.Node.Children[i]
			p, n := actionVisit(&NodeVisit{Node: &child})

			visitNode.Children = append(visitNode.Children, *n)
			visitNode.Parsers = append(visitNode.Parsers, p)
		}
	}

	if visitNode.Node.Action != "" {
		switch visitNode.Node.Action {
		case "noop":
			{
				return visitNoop(visitNode), visitNode
			}
		case "peg":
			{
				return visitPeg(visitNode), visitNode
			}
		case "parsing_body":
			{
				return visitParsingBody(visitNode), visitNode
			}
		case "parsing_expression":
			{
				return visitParsingExpression(visitNode), visitNode
			}
		case "parsing_rule":
			{
				return visitParsingRule(visitNode), visitNode
			}
		case "parsing_sequence":
			{
				return visitParsingSequence(visitNode), visitNode
			}
		case "parsing_ordered_choice":
			{
				return visitParsingOrderedChoice(visitNode), visitNode
			}
		case "parsing_sub_expression":
			{
				return visitParsingSubExpression(visitNode), visitNode
			}
		case "parsing_group":
			{
				return visitParsingGroup(visitNode), visitNode
			}
		case "parsing_atomic_expression":
			{
				return visitParsingAtomicExpression(visitNode), visitNode
			}
		case "parsing_not_predicate":
			{
				return visitParsingNotPredicate(visitNode), visitNode
			}
		case "parsing_and_predicate":
			{
				return visitParsingAndPredicate(visitNode), visitNode
			}
		case "parsing_zero_or_more":
			{
				return visitParsingZeroOrMore(visitNode), visitNode
			}
		case "parsing_one_or_more":
			{
				return visitParsingOneOrMore(visitNode), visitNode
			}
		case "parsing_optional":
			{
				return visitParsingOptional(visitNode), visitNode
			}
		case "parsing_string":
			{
				return visitParsingString(visitNode), visitNode
			}
		case "parsing_regex_char":
			{
				return visitParsingRegexChar(visitNode), visitNode
			}
		case "parsing_rule_call":
			{
				return visitParsingRuleCall(visitNode), visitNode
			}
		case "parsing_end_of_file":
			{
				return visitParsingEndOfFile(visitNode), visitNode
			}
		}
	}

	return visitNoop(visitNode), visitNode
}

func visitNoop(node *NodeVisit) rd.ParserFunc {
	return func(state *rd.State) (rd.Ast, bool) {
		return *node.Node, false
	}
}

func visitPeg(node *NodeVisit) rd.ParserFunc {
	return node.Parsers[3]
}

func visitParsingBody(node *NodeVisit) rd.ParserFunc {
	var children []rd.ParserFunc
	for i := 0; i < len(node.Children); i++ {
		children = append(children, node.Children[i].Parsers[0])
	}

	return children[0]
}

func visitParsingRule(node *NodeVisit) rd.ParserFunc {
	var rule = node.Parsers[4]
	var ruleName = node.Node.Children[0].Match

	var parser = func(state *rd.State) (rd.Ast, bool) {
		var start = state.Position
		var ast, err = rule(state)

		if !err {
			ast.Rule = ruleName
			state.SuccesfullRules = append(state.SuccesfullRules, rd.Ast{
				Rule:          ast.Rule,
				Match:         ast.Match,
				StartPosition: ast.StartPosition,
				EndPosition:   ast.EndPosition,
			})
		} else {
			state.FailedRules = append(state.FailedRules, rd.Ast{
				Rule:          ruleName,
				StartPosition: start,
			})
		}

		return ast, err
	}

	GRules.rules = append(GRules.rules, rd.Rule{Name: ruleName, Parser: parser})

	return parser
}

func visitParsingExpression(node *NodeVisit) rd.ParserFunc {
	return node.Parsers[0]
}

func visitParsingSequence(node *NodeVisit) rd.ParserFunc {
	var head = []rd.ParserFunc{node.Children[0].Parsers[0]}

	for i := 0; i < len(node.Children[1].Children); i++ {
		head = append(head, node.Children[1].Children[i].Children[1].Parsers[0])
	}

	return rd.Sequence(head)
}

//
func visitParsingOrderedChoice(node *NodeVisit) rd.ParserFunc {
	var head = []rd.ParserFunc{node.Parsers[0]}

	for i := 0; i < len(node.Children[1].Children); i++ {
		head = append(head, node.Children[1].Children[i].Parsers[3])
	}

	return rd.OrderedChoice(head)
}

func visitParsingSubExpression(node *NodeVisit) rd.ParserFunc {
	return func(state *rd.State) (rd.Ast, bool) {
		ast, err := node.Children[1].Parsers[0](state)

		// TODO tags
		return ast, err
	}
}

//parsing_group
func visitParsingGroup(node *NodeVisit) rd.ParserFunc {
	return node.Parsers[2]
}

//parsing_atomic_expression
func visitParsingAtomicExpression(node *NodeVisit) rd.ParserFunc {
	return node.Parsers[0]
}

//parsing_not_predicate
func visitParsingNotPredicate(node *NodeVisit) rd.ParserFunc {
	return rd.NotPredicate(node.Children[1].Parsers[0])
}

//parsing_and_predicate
func visitParsingAndPredicate(node *NodeVisit) rd.ParserFunc {
	return rd.AndPredicate(node.Children[1].Parsers[0])
}

//parsing_zero_or_more
func visitParsingZeroOrMore(node *NodeVisit) rd.ParserFunc {
	return rd.ZeroOrMore(node.Children[0].Parsers[0])
}

//parsing_one_or_more
func visitParsingOneOrMore(node *NodeVisit) rd.ParserFunc {
	return rd.OneOrMore(node.Children[0].Parsers[0])
}

//parsing_optional
func visitParsingOptional(node *NodeVisit) rd.ParserFunc {
	return rd.Optional(node.Children[0].Parsers[0])
}

//parsing_string
func visitParsingString(node *NodeVisit) rd.ParserFunc {
	text := node.Node.Children[1].Match
	text = strings.ReplaceAll(text, "\\\\\\\\", "\\\\")
	text = strings.ReplaceAll(text, "\\\"", "\"")

	return rd.String(text)
}

//parsing_regex_char
func visitParsingRegexChar(node *NodeVisit) rd.ParserFunc {
	return rd.RegexChar(node.Node.Children[0].Match)
}

//parsing_rule_call
func visitParsingRuleCall(node *NodeVisit) rd.ParserFunc {
	return rd.CallRuleByName(node.Node.Match)
}

//parsing_end_of_file
func visitParsingEndOfFile(node *NodeVisit) rd.ParserFunc {
	return rd.EndOfFile()
}
