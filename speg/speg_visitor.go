package speg

import (
	"github.com/SimplePEG/Go/rd"
)

type GrammarRules struct {
	rules []rd.Rule
}

type NodeVisit struct {
	Parsers  []rd.ParserFunc
	Node     *rd.Ast
	Children []NodeVisit
}

// Global
var GRules GrammarRules

func ParseText(ast rd.Ast, text string) (rd.Ast, bool) {
	_, child := ActionVisit(&NodeVisit{Node: &ast})

	parser := child.Parsers[3]

	result, err := parser(&rd.State{
		Text:     text,
		Position: 0,
		Rules:    GRules.rules,
	})

	return result, err
}

func ActionVisit(visitNode *NodeVisit) (rd.ParserFunc, NodeVisit) {
	if len(visitNode.Node.Children) > 0 {
		for i := 0; i < len(visitNode.Node.Children); i++ {
			child := visitNode.Node.Children[i]
			p, n := ActionVisit(&NodeVisit{Node: &child})

			visitNode.Children = append(visitNode.Children, n)
			visitNode.Parsers = append(visitNode.Parsers, p)
		}
	}

	if visitNode.Node.Action != "" {
		switch visitNode.Node.Action {
		case "noop":
			{
				return visitNoop(visitNode)
			}
		case "peg":
			{
				return visitPeg(visitNode)
			}
		case "parsing_body":
			{
				return visitParsingBody(visitNode)
			}
		case "parsing_expression":
			{
				return visitParsingExpression(visitNode)
			}
		case "parsing_rule":
			{
				return visitParsingRule(visitNode)
			}
		case "parsing_sequence":
			{
				return visitParsingSequence(visitNode)
			}
		case "parsing_ordered_choice":
			{
				return visitParsingOrderedChoice(visitNode)
			}
		case "parsing_sub_expression":
			{
				return visitParsingSubExpression(visitNode)
			}
		case "parsing_group":
			{
				return visitParsingGroup(visitNode)
			}
		case "parsing_atomic_expression":
			{
				return visitParsingAtomicExpression(visitNode)
			}
		case "parsing_not_predicate":
			{
				return visitParsingNotPredicate(visitNode)
			}
		case "parsing_and_predicate":
			{
				return visitParsingAndPredicate(visitNode)
			}
		case "parsing_zero_or_more":
			{
				return visitParsingZeroOrMore(visitNode)
			}
		case "parsing_one_or_more":
			{
				return visitParsingOneOrMore(visitNode)
			}
		case "parsing_optional":
			{
				return visitParsingOptional(visitNode)
			}
		case "parsing_string":
			{
				return visitParsingString(visitNode)
			}
		case "parsing_regex_char":
			{
				return visitParsingRegexChar(visitNode)
			}
		case "parsing_rule_call":
			{
				return visitParsingRuleCall(visitNode)
			}
		case "parsing_end_of_file":
			{
				return visitParsingEndOfFile(visitNode)
			}
		}
	}

	return visitNoop(visitNode)
}

func visitNoop(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	return func(state *rd.State) (rd.Ast, bool) {
		return *node.Node, false
	}, *node
}

func visitPeg(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	//return func(state *rd.State) (rd.Ast, bool) {
	//	ast := node.Children[3]
	//	return ast, false
	//}
	return node.Parsers[3], *node
}

func visitParsingBody(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	var children []rd.ParserFunc
	for i := 0; i < len(node.Children); i++ {
		children = append(children, node.Children[i].Parsers[0])
	}

	//node.Node.Children = children

	//return visitTryParser(node.Node).(rd.ParserFunc), *node
	return children[0], *node
}

func visitParsingRule(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	var rule = node.Parsers[4]
	var ruleName = node.Node.Children[0].Match

	var parser = func(state *rd.State) (rd.Ast, bool) {
		//var start = state.Position
		var ast, err = rule(state)

		if !err {
			ast.Rule = ruleName
		}
		//	state.SuccesfullRules = append(state.SuccesfullRules, rd.Ast{
		//		Rule: ast.Rule,
		//		Match: ast.Match,
		//		StartPosition: ast.StartPosition,
		//		EndPosition: ast.EndPosition,
		//	})
		//} else {
		//	state.FailedRules = append(state.FailedRules, rd.Ast{
		//		Rule: ruleName,
		//		StartPosition: start,
		//	})
		//}

		return ast, err
	}

	GRules.rules = append(GRules.rules, rd.Rule{Name: ruleName, Parser: parser})

	return parser, *node

	//return rd.Rule{
	//	Name: ruleName,
	//	Parser: parser,
	//}
}

//
//
func visitParsingExpression(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	return node.Parsers[0], *node
}

func visitParsingSequence(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	var head = []rd.ParserFunc{node.Children[0].Parsers[0]}

	for i := 0; i < len(node.Children[1].Children); i++ {
		head = append(head, node.Children[1].Children[i].Children[1].Parsers[0])
	}

	return rd.Sequence(head), *node
}

//
func visitParsingOrderedChoice(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	var head = []rd.ParserFunc{node.Children[0].Parsers[0]}

	for i := 0; i < len(node.Children[1].Children); i++ {
		head = append(head, node.Children[1].Children[i].Parsers[3])
	}

	return rd.OrderedChoice(head), *node
}

func visitParsingSubExpression(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	return func(state *rd.State) (rd.Ast, bool) {
		ast, err := node.Children[1].Parsers[0](state)

		// TODO tags
		return ast, err
	}, *node
}

//parsing_group
func visitParsingGroup(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	return node.Parsers[2], *node
}

//parsing_atomic_expression
func visitParsingAtomicExpression(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	return node.Parsers[0], *node
}

//parsing_not_predicate
func visitParsingNotPredicate(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	return rd.NotPredicate(node.Children[1].Parsers[0]), *node
}

//parsing_and_predicate
func visitParsingAndPredicate(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	return rd.AndPredicate(node.Children[1].Parsers[0]), *node
}

//parsing_zero_or_more
func visitParsingZeroOrMore(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	return rd.ZeroOrMore(node.Children[0].Parsers[0]), *node
}

//parsing_one_or_more
func visitParsingOneOrMore(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	return rd.OneOrMore(node.Children[0].Parsers[0]), *node
}

//parsing_optional
func visitParsingOptional(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	return rd.Optional(node.Children[0].Parsers[0]), *node
}

//parsing_string
func visitParsingString(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	// TODO match
	//node.children[1].match
	//.replace(/\\\\/g, '\\')
	//.replace(/\\"/g, '"')
	return rd.String(node.Node.Children[1].Match), *node
}

//parsing_regex_char
func visitParsingRegexChar(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	return rd.RegexChar(node.Node.Children[0].Match), *node
}

//parsing_rule_call
func visitParsingRuleCall(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	return rd.CallRuleByName(node.Node.Match), *node
}

//parsing_end_of_file
func visitParsingEndOfFile(node *NodeVisit) (rd.ParserFunc, NodeVisit) {
	return rd.EndOfFile(), *node
}
