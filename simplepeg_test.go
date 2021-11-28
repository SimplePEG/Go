package simplepeg

import (
	"github.com/SimplePEG/Go/speg"
	"testing"
)

func TestParse(t *testing.T) {
	var grammarText = `GRAMMAR url

url       ->  scheme "://" host pathname search hash?;
scheme    ->  "http" "s"?;
host      ->  hostname port?;
hostname  ->  segment ("." segment)*;
segment   ->  [a-z0-9-]+;
port      ->  ":" [0-9]+;
pathname  ->  "/" [^ ?]*;
search    ->  ("?" [^ #]*)?;
hash      ->  "#" [^ ]*;
  `
	var spegParser = speg.NewSPEGParser()

	var ast, err = spegParser.ParseGrammar(grammarText)

	if ast.Match != grammarText {
		t.Error("Grammar not matched to original")
	}

	if !err {

		textAst, textErr := speg.ParseText(ast, "https://simplepeg.github.io/")

		if textErr {
			t.Error("Text not parsed")
		}

		if textAst.Match != "https://simplepeg.github.io/" {
			t.Error("Text not matched")
		}

		if textAst.StartPosition != 0 {
			t.Error("Text in not position")
		}

		if textAst.EndPosition != 28 {
			t.Error("Text in not position")
		}

	} else {
		t.Error("Grammar not parsed")
	}
}
