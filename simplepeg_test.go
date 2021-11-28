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

url       ->  scheme host;
scheme    ->  "http";
host      ->  "s";
  `
	var spegParser = speg.NewSPEGParser()

	var ast, err = spegParser.Parse(grammarText)

	if !err {

		textAst, textErr := speg.ParseText(ast, "https")

		println(textAst.Match, textErr)

	} else {
		println("Err")
	}
}
