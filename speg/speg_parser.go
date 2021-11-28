package speg

import (
	"github.com/SimplePEG/Go/rd"
)

func peg() rd.ParserFunc {
	return rd.Action("peg", rd.Sequence([]rd.ParserFunc{
		rd.ZeroOrMore(noop()),
		parsingHeader(),
		rd.OneOrMore(noop()),
		parsingBody(),
		rd.EndOfFile(),
	}))
}

func parsingHeader() rd.ParserFunc {
	return rd.Action("noop", rd.Sequence([]rd.ParserFunc{
		rd.String("GRAMMAR"),
		rd.OneOrMore(noop()),
		rd.OneOrMore(parsingRuleName()),
	}))
}

func parsingBody() rd.ParserFunc {
	return rd.Action("parsing_body", rd.OneOrMore(rd.OrderedChoice([]rd.ParserFunc{
		rd.String("GRAMMAR"),
		parsingRule(),
		rd.OneOrMore(noop()),
	})))
}

func parsingRule() rd.ParserFunc {
	return rd.Action("parsing_rule", rd.Sequence([]rd.ParserFunc{
		parsingRuleName(),
		rd.ZeroOrMore(noop()),
		rd.String("->"),
		rd.ZeroOrMore(noop()),
		parsingExpression(),
		rd.ZeroOrMore(noop()),
		rd.String(";"),
		rd.ZeroOrMore(noop()),
	}))
}

func parsingRuleName() rd.ParserFunc {
	return rd.Action("noop", rd.Sequence([]rd.ParserFunc{
		rd.RegexChar("[a-zA-Z_]"),
		rd.ZeroOrMore(rd.RegexChar("[a-zA-Z0-9_]")),
	}))
}

func parsingExpression() rd.ParserFunc {
	return rd.Action("parsing_expression", rd.OrderedChoice([]rd.ParserFunc{
		parsingSequence(),
		parsingOrderedChoice(),
		parsingSubExpression(),
	}))
}

func parsingSequence() rd.ParserFunc {
	return rd.Action("parsing_sequence", rd.Sequence([]rd.ParserFunc{
		rd.OrderedChoice([]rd.ParserFunc{
			parsingOrderedChoice(),
			parsingSubExpression(),
		}),
		rd.OneOrMore(rd.Sequence([]rd.ParserFunc{
			rd.OneOrMore(noop()),
			rd.OrderedChoice([]rd.ParserFunc{
				parsingOrderedChoice(),
				parsingSubExpression(),
			}),
		})),
	}))
}

func parsingOrderedChoice() rd.ParserFunc {
	return rd.Action("parsing_ordered_choice", rd.Sequence([]rd.ParserFunc{
		parsingSubExpression(),
		rd.OneOrMore(rd.Sequence([]rd.ParserFunc{
			rd.ZeroOrMore(noop()),
			rd.String("/"), // rd.string('\/'),
			rd.ZeroOrMore(noop()),
			parsingSubExpression(),
		})),
	}))
}

func parsingSubExpression() rd.ParserFunc {
	return rd.Action("parsing_sub_expression", rd.Sequence([]rd.ParserFunc{
		rd.ZeroOrMore(rd.Sequence([]rd.ParserFunc{
			tag(),
			rd.String(":"),
		})),
		rd.OrderedChoice([]rd.ParserFunc{
			parsing_not_predicate(),
			parsingAndPredicate(),
			parsingOptional(),
			parsingOneOrMore(),
			parsingZeroOrMore(),
			parsingGroup(),
			parsingAtomicExpression(),
		}),
	}))
}

func tag() rd.ParserFunc {
	return rd.Action("noop", rd.Sequence([]rd.ParserFunc{
		rd.RegexChar("[a-zA-Z_]"),
		rd.ZeroOrMore(rd.RegexChar("[a-zA-Z0-9_]")),
	}))
}

func parsingGroup() rd.ParserFunc {
	return rd.Action("parsing_group", rd.Sequence([]rd.ParserFunc{
		rd.String("("),
		rd.ZeroOrMore(noop()),
		rd.Rec(parsingExpression),
		rd.ZeroOrMore(noop()),
		rd.String(")"),
	}))
}

func parsingAtomicExpression() rd.ParserFunc {
	return rd.Action("parsing_atomic_expression", rd.OrderedChoice([]rd.ParserFunc{
		parsingString(),
		parsingRegexChar(),
		parsingEof(),
		parsingRuleCall(),
	}))
}

func parsing_not_predicate() rd.ParserFunc {
	return rd.Action("parsing_not_predicate", rd.Sequence([]rd.ParserFunc{
		rd.String("!"),
		rd.OrderedChoice([]rd.ParserFunc{
			parsingGroup(),
			parsingAtomicExpression(),
		}),
	}))
}

func parsingAndPredicate() rd.ParserFunc {
	return rd.Action("parsing_and_predicate", rd.Sequence([]rd.ParserFunc{
		rd.String("&"),
		rd.OrderedChoice([]rd.ParserFunc{
			parsingGroup(),
			parsingAtomicExpression(),
		}),
	}))
}

func parsingZeroOrMore() rd.ParserFunc {
	return rd.Action("parsing_zero_or_more", rd.Sequence([]rd.ParserFunc{
		rd.OrderedChoice([]rd.ParserFunc{
			parsingGroup(),
			parsingAtomicExpression(),
		}),
		rd.String("*"),
	}))
}

func parsingOneOrMore() rd.ParserFunc {
	return rd.Action("parsing_one_or_more", rd.Sequence([]rd.ParserFunc{
		rd.OrderedChoice([]rd.ParserFunc{
			parsingGroup(),
			parsingAtomicExpression(),
		}),
		rd.String("+"),
	}))
}

func parsingOptional() rd.ParserFunc {
	return rd.Action("parsing_optional", rd.Sequence([]rd.ParserFunc{
		rd.OrderedChoice([]rd.ParserFunc{
			parsingGroup(),
			parsingAtomicExpression(),
		}),
		rd.String("?"),
	}))
}

func parsingRuleCall() rd.ParserFunc {
	return rd.Action("parsing_rule_call", parsingRuleName())
}

func parsingString() rd.ParserFunc {
	return rd.Action("parsing_string", rd.Sequence([]rd.ParserFunc{
		rd.String("\""),
		rd.OneOrMore(rd.OrderedChoice([]rd.ParserFunc{
			rd.String("\\\\"),
			rd.String("\\\""),
			rd.RegexChar("[^\"]"),
		})),
		rd.String("\""),
	}))
}

func parsingRegexChar() rd.ParserFunc {
	return rd.Action("parsing_regex_char", rd.OrderedChoice([]rd.ParserFunc{
		rd.Sequence([]rd.ParserFunc{
			rd.String("["),
			rd.Optional(rd.String("^")),
			rd.OneOrMore(rd.OrderedChoice([]rd.ParserFunc{
				rd.String("\\]"),
				rd.String("\\["),
				rd.RegexChar("[^\\]]"),
			})),
			rd.String("]"),
		}),
		rd.String("."),
	}))
}

func parsingEof() rd.ParserFunc {
	return rd.Action("parsing_end_of_file", rd.String("EOF"))
}

func noop() rd.ParserFunc {
	return rd.Action("noop", rd.RegexChar("[\\s]"))
}

type SPEGParser struct {
	state  rd.State
	parser rd.ParserFunc
}

func NewSPEGParser() SPEGParser {
	speg := SPEGParser{parser: peg()}
	return speg
}

func (sp SPEGParser) ParseGrammar(text string) (rd.Ast, bool) {
	state := &rd.State{
		Text:     text,
		Position: 0,
	}

	//sp.state = state

	return sp.parser(state)
}

func (sp SPEGParser) GetLastExpectations() []rd.Expectation {
	return sp.state.LastExpectations
}

func (sp SPEGParser) GetLastError() (string, bool) {
	return rd.GetLastError(&sp.state)
}
