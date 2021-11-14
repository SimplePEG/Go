package speg_actions

//
//func Noop(node rd.Ast) rd.Ast {
//	return node
//}
//
//func Peg(node rd.Ast) rd.Ast {
//	return node.Children[3]
//}
//
//func ParsingBody(node rd.Ast) rd.Ast {
//	var children []rd.Ast
//	for i := 0; i < len(node.Children); i++ {
//		children = append(children, node.Children[i].Children[0])
//	}
//	node.Children = children
//
//	return node
//}
//
//func ParsingRule(node rd.Ast) rd.Rule {
//	var rule = *node.Children[4].Parser
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
//func ParsingExpression(node rd.Ast) rd.Ast {
//	return node.Children[0]
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
