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
