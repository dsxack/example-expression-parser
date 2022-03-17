package parser

import (
	"github.com/alecthomas/participle/v2"
)

var statementParser = participle.MustBuild(&Statement{},
	participle.Lexer(statementLexer),
	participle.Unquote("String"),
)

func Parse(input string, stmt *Statement) error {
	return statementParser.ParseString("", input, stmt)
}
