package speg

import "github.com/SimplePEG/Go/rd"

func Parse(grammar string, text string) (rd.Ast, bool) {
	var spegParser = NewSPEGParser()

	var gAst, gErr = spegParser.ParseGrammar(grammar)

	if !gErr {
		return ParseText(gAst, text)
	}

	return rd.Ast{}, gErr
}

// GRules as global variable to
var GRules GrammarRules

func ParseText(ast rd.Ast, text string) (rd.Ast, bool) {
	parser, grule := GetParser(ast)

	result, err := parser(&rd.State{
		Text:     text,
		Position: 0,
		Rules:    grule.Rules,
	})

	return result, err
}

func GetParser(ast rd.Ast) (rd.ParserFunc, GrammarRules) {
	GRules = GrammarRules{}
	_, visitNode := actionVisit(&NodeVisit{Node: &ast})

	parser := visitNode.Parsers[3]

	return parser, GRules
}
