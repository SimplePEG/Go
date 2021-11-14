package simplepeg

import (
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

	//var ast, err = rd.RegexChar("[\\s]")(&rd.State{Text: "1  ", Position: 0})

	if !err {
		println(123)
		println(ast.EndPosition)
	} else {
		println("Err")
	}
}
