package speg

import (
	"github.com/SimplePEG/Go/rd"
)

func SPEGActionVisit(node interface{}) interface{} {
	for i := 0; i < len(node.(rd.Ast).Children); i++ {
		ast := node.(rd.Ast).Children[i]
		SPEGActionVisit(ast)
	}

	if node.(rd.Ast).Action != "" {
		switch node.(type) {
		case rd.Ast:
			{
				switch node.(rd.Ast).Action {
				case "noop":
					{
						return visitNoop(&node)
					}
				case "peg":
					{
						return visitPeg(node)
					}
					//case "parsing_body":
					//	{
					//		return visitParsingBody(node)
					//	}
					//case "parsing_expression":
					//	{
					//		return visitParsingExpression(node)
					//	}
				}
			}
		}
	}

	return node

}

func visitNoop(node *interface{}) *interface{} {
	return node
}

func visitPeg(node interface{}) interface{} {
	return node.(rd.Ast).Children[3]
}

//func visitParsingBody(node interface{}) interface{} {
//	var children []rd.Ast
//	for i := 0; i < len(node.Children); i++ {
//		children = append(children, node.Children[i].Children[0])
//	}
//	node.Children = children
//
//	return node
//}
//
//
//func ParsingRule(v *SPEGActionVisitor, node *rd.Ast) *interface{} {
//	var rule = *node.Children[4]
//	var ruleName = node.Children[0].Match
//
//	var parser = func(state *rd.State) (rd.Ast, bool) {
//		var start = state.Position
//		var ast, err = rule(state)
//
//		if(!err) {
//			ast.Rule = ruleName
//			state.SuccesfullRules = append(state.SuccesfullRules, rd.Ast{
//				Rule: ast.Rule,
//				Match: ast.Match,
//				StartPosition: ast.StartPosition,
//				EndPosition: ast.EndPosition,
//			})
//		} else {
//			state.FailedRules = append(state.FailedRules, rd.Ast{
//				Rule: ruleName,
//				StartPosition: start,
//			})
//		}
//
//		return ast, false
//
//	}
//
//	return rd.Rule{
//		Name: ruleName,
//		Parser: parser,
//	}
//}
//
//
//func visitParsingExpression(v *SPEGActionVisitor, node *rd.Ast) *rd.Ast {
//	return &node.Children[0]
//}
//
//func ParsingSequence(node rd.Ast) rd.ParserFunc {
//	var head = []rd.ParserFunc{ *node.Children[0].Children[0].Parser}
//
//	for i := 0; i < len(node.Children[1].Children); i++ {
//		head = append(head, *node.Children[1].Children[i].Children[1].Children[0].Parser)
//	}
//
//	return rd.Sequence(head)
//}
//
//func ParsingOrderedChoice(node rd.Ast) rd.ParserFunc {
//	var head = []rd.ParserFunc{ *node.Children[0].Children[0].Parser}
//
//	for i := 0; i < len(node.Children[1].Children); i++ {
//		head = append(head, *node.Children[1].Children[i].Children[3].Parser)
//	}
//
//	return rd.OrderedChoice(head)
//}
