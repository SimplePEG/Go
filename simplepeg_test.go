package simplepeg

import (
	"github.com/SimplePEG/Go/rd"
	"github.com/SimplePEG/Go/speg"
	"testing"
)

func TestHello(t *testing.T) {
	v := Hello()

	if v != "World" {
		t.Error("Expected World, got", v)
	}

	var grammarText = `GRAMMAR url
	
	url ->  "1";`
	var spegParser = speg.NewSPEGParser()

	var ast, err = spegParser.Parse(grammarText)

	if !err {
		println(len(ast.Children))
		//ast.Visit()
		_, child := speg.ActionVisit(&ast)
		parser := child.Parsers[0]

		textAst, textErr := parser(&rd.State{
			Text:     "1",
			Position: 0,
		})

		println(textAst.Match, textErr)
	} else {
		println("Err")
	}
}
