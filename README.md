# Go

-----

Go version of SimplePEG. A very simple implementation of PEG parser generator.

```Go
import "github.com/SimplePEG/Go/SimplePeg"

func main() {
    var ast = simplepeg.Parse(`GRAMMAR test a->"A";`, 'A')
	
	println(ast.Match)
}

```

Or create parser function

```Go
import "github.com/SimplePEG/Go/SimplePeg/speg"

func main() {
    var spegParser = speg.NewSPEGParser()
    
    var gAst, gErr = spegParser.ParseGrammar(grammar)
    
    if !gErr {
            parser, grule := speg.GetParser(ast)
            
            result, err := parser(&rd.State{
                Text:     text,
                Position: 0,
                Rules:    grule.rules,
            })

			//...
	}
}

```