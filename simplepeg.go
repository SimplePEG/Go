package simplepeg

import (
	"github.com/SimplePEG/Go/rd"
	"github.com/SimplePEG/Go/speg"
)

func Parse(grammar string, text string) (rd.Ast, bool) {
	return speg.Parse(grammar, text)
}
