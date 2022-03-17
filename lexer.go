package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
	"strings"
)

var statementLexer = lexer.MustStateful(lexer.Rules{
	"Root": []lexer.Rule{
		{Name: "Ternary", Pattern: `\#\{`, Action: lexer.Push("Ternary")},
		{Name: "Variable", Pattern: `\$\{`, Action: lexer.Push("Variable")},
		{Name: "Literal", Pattern: `[^\$\#]+|[\$\#][^\{][^\$\#]+|[\$\#]`},
		lexer.Include("Common"),
	},
	"Ternary": []lexer.Rule{
		{Name: "TernaryThen", Pattern: `\?`},
		{Name: "TernaryElse", Pattern: `\:`},
		{Name: "TernaryEnd", Pattern: `}`, Action: lexer.Pop()},
		lexer.Include("Expr"),
		lexer.Include("Common"),
	},
	"Condition": []lexer.Rule{
		{Name: "Variable", Pattern: `\$\{`, Action: lexer.Push("Variable")},
		{Name: "ComparisonOperator", Pattern: strings.Join(comparisonOperatorKeywords(), "|")},
		{Name: "String", Pattern: `'[^']*'|"[^"]*"`},
		lexer.Include("Common"),
	},
	"Expr": []lexer.Rule{
		{Name: "LogicOperator", Pattern: `and|or|AND|OR|&&|\|\|`},
		{Name: "SubExpr", Pattern: `\(`, Action: lexer.Push("SubExpr")},
		lexer.Include("Condition"),
	},
	"SubExpr": []lexer.Rule{
		lexer.Include("Expr"),
		{Name: "SubExprEnd", Pattern: `\)`, Action: lexer.Pop()},
	},
	"Variable": []lexer.Rule{
		{Name: "VariableEnd", Pattern: `}`, Action: lexer.Pop()},
		lexer.Include("Common"),
	},
	"Common": []lexer.Rule{
		{Name: "whitespace", Pattern: `\s+`},
		{Name: "Number", Pattern: `(?:\d*\.)?\d+`},
		{Name: "Boolean", Pattern: `false|true|TRUE|FALSE"`},
		{Name: "Ident", Pattern: `[a-zA-Z][a-zA-Z0-9_\.]*`},
	},
})
