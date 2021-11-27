package speg

import (
	"github.com/SimplePEG/Go/rd"
)

type NodeVisit struct {
	Parsers  []rd.ParserFunc
	Node     *rd.Ast
	Children []NodeVisit
}

func ActionVisit(node *rd.Ast) (rd.ParserFunc, NodeVisit) {
	var visitNode = NodeVisit{
		Node: node,
	}

	if len(node.Children) > 0 {
		for i := 0; i < len(node.Children); i++ {
			ast := node.Children[i]
			p, n := ActionVisit(&ast)
			visitNode.Children = append(visitNode.Children, n)
			visitNode.Parsers = append(visitNode.Parsers, p)
			return p, n
		}
	}

	if node.Action != "" {
		switch node.Action {
		case "noop":
			{
				return visitNoop(node), visitNode
			}
		case "peg":
			{
				return visitPeg(node), visitNode
			}
		case "parsing_body":
			{
				return visitParsingBody(node), visitNode
			}
		case "parsing_expression":
			{
				return visitParsingExpression(node), visitNode
			}
		case "parsing_rule":
			{
				return visitParsingRule(node), visitNode
			}
		case "parsing_sequence":
			{
				return visitParsingSequence(node), visitNode
			}
		case "parsing_ordered_choice":
			{
				return visitParsingOrderedChoice(node), visitNode
			}
		case "parsing_sub_expression":
			{
				return visitParsingSubExpression(node), visitNode
			}
		case "parsing_group":
			{
				return visitParsingGroup(node), visitNode
			}
		case "parsing_atomic_expression":
			{
				return visitParsingAtomicExpression(node), visitNode
			}
		case "parsing_not_predicate":
			{
				return visitParsingNotPredicate(node), visitNode
			}
		case "parsing_and_predicate":
			{
				return visitParsingAndPredicate(node), visitNode
			}
		case "parsing_zero_or_more":
			{
				return visitParsingZeroOrMore(node), visitNode
			}
		case "parsing_optional":
			{
				return visitParsingOptional(node), visitNode
			}
		case "parsing_string":
			{
				return visitParsingString(node), visitNode
			}
		case "parsing_regex_char":
			{
				return visitParsingRegexChar(node), visitNode
			}
		case "parsing_rule_call":
			{
				return visitParsingRuleCall(node), visitNode
			}
		case "parsing_end_of_file":
			{
				return visitParsingEndOfFile(node), visitNode
			}
		}
	}

	return visitNoop(node), visitNode
}

func visitNoop(node *rd.Ast) rd.ParserFunc {
	return func(state *rd.State) (rd.Ast, bool) {
		return *node, false
	}
}

func visitTryParser(node *rd.Ast) rd.ParserFunc {
	return func(state *rd.State) (rd.Ast, bool) {
		return *node, false
	}
}

func visitPeg(node *rd.Ast) rd.ParserFunc {
	return func(state *rd.State) (rd.Ast, bool) {
		ast := node.Children[3]
		return ast, false
	}
}

func visitParsingBody(node *rd.Ast) rd.ParserFunc {
	var children []rd.Ast
	for i := 0; i < len(node.Children); i++ {
		children = append(children, node.Children[i].Children[0])
	}
	node.Children = children

	return visitNoop(node)
}

func visitParsingRule(node *rd.Ast) rd.ParserFunc {
	var rule = visitTryParser(&node.Children[4])
	//var ruleName = node.Children[0].Match

	var parser = func(state *rd.State) (rd.Ast, bool) {
		//var start = state.Position
		var ast, err = rule(state)

		//if(!err) {
		//	ast.Rule = ruleName
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

	return parser

	//return rd.Rule{
	//	Name: ruleName,
	//	Parser: parser,
	//}
}

//
//
func visitParsingExpression(node *rd.Ast) rd.ParserFunc {
	return visitTryParser(&node.Children[0])
}

func visitParsingSequence(node *rd.Ast) rd.ParserFunc {
	var head = []rd.ParserFunc{visitTryParser(&node.Children[0].Children[0])}

	for i := 0; i < len(node.Children[1].Children); i++ {
		head = append(head, visitTryParser(&node.Children[1].Children[i].Children[1].Children[0]))
	}

	return rd.Sequence(head)
}

//
func visitParsingOrderedChoice(node *rd.Ast) rd.ParserFunc {
	var head = []rd.ParserFunc{visitTryParser(&node.Children[0].Children[0])}

	for i := 0; i < len(node.Children[1].Children); i++ {
		head = append(head, visitTryParser(&node.Children[1].Children[i].Children[3]))
	}

	return rd.OrderedChoice(head)
}

func visitParsingSubExpression(node *rd.Ast) rd.ParserFunc {
	return func(state *rd.State) (rd.Ast, bool) {
		ast, err := visitTryParser(&node.Children[1].Children[0])(state)

		// TODO tags
		return ast, err
	}
}

//parsing_group
func visitParsingGroup(node *rd.Ast) rd.ParserFunc {
	return visitTryParser(&node.Children[2])
}

//parsing_atomic_expression
func visitParsingAtomicExpression(node *rd.Ast) rd.ParserFunc {
	return visitTryParser(&node.Children[0])
}

//parsing_not_predicate
func visitParsingNotPredicate(node *rd.Ast) rd.ParserFunc {
	return rd.NotPredicate(visitTryParser(&node.Children[1].Children[0]))
}

//parsing_and_predicate
func visitParsingAndPredicate(node *rd.Ast) rd.ParserFunc {
	return rd.AndPredicate(visitTryParser(&node.Children[1].Children[0]))
}

//parsing_zero_or_more
func visitParsingZeroOrMore(node *rd.Ast) rd.ParserFunc {
	return rd.ZeroOrMore(visitTryParser(&node.Children[0].Children[0]))
}

//parsing_one_or_more
func visitParsingOneOrMore(node *rd.Ast) rd.ParserFunc {
	return rd.OneOrMore(visitTryParser(&node.Children[0].Children[0]))
}

//parsing_optional
func visitParsingOptional(node *rd.Ast) rd.ParserFunc {
	return rd.Optional(visitTryParser(&node.Children[0].Children[0]))
}

//parsing_string
func visitParsingString(node *rd.Ast) rd.ParserFunc {
	// TODO match
	//node.children[1].match
	//.replace(/\\\\/g, '\\')
	//.replace(/\\"/g, '"')
	return rd.String(node.Children[1].Match)
}

//parsing_regex_char
func visitParsingRegexChar(node *rd.Ast) rd.ParserFunc {
	return rd.RegexChar(node.Children[0].Match)
}

//parsing_rule_call
func visitParsingRuleCall(node *rd.Ast) rd.ParserFunc {
	return rd.CallRuleByName(node.Match)
}

//parsing_end_of_file
func visitParsingEndOfFile(node *rd.Ast) rd.ParserFunc {
	return rd.EndOfFile()
}
