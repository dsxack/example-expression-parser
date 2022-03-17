package parser

import (
	"fmt"
	"strings"
)

type Statement struct {
	Fragments []Fragment `parser:"@@+"`
}

type Fragment struct {
	IfStatement *IfStatement `parser:"'#{' @@ '}' "`
	Variable    *string      `parser:"| '${' @Ident '}'"`
	Literal     *string      `parser:"| @!(Ternary Variable)"`
}

type IfStatement struct {
	Cond Expr `parser:"@@"`
	Then Term `parser:"'?' @@"`
	Else Term `parser:"':' @@"`
}

type Expr struct {
	Left  *Cond       `parser:"( @@"`
	Sub   *Expr       `parser:"| '(' @@ ')' )"`
	Right []LogicExpr `parser:"@@*"`
}

type LogicExpr struct {
	Operator LogicOperator `parser:"@LogicOperator"`
	Cond     *Cond         `parser:"( @@"`
	Sub      *Expr         `parser:"| '(' @@ ')' )"`
}

type Cond struct {
	Left     Term               `parser:"@@"`
	Operator ComparisonOperator `parser:"( @ComparisonOperator"`
	Right    *Term              `parser:"@@ )?"`
}

type Term struct {
	Variable *string `parser:"'${' @Ident '}'"`
	Value    *Value  `parser:"| @@"`
}

type Value struct {
	Number  *float64 `parser:"@Number"`
	String  *string  `parser:"| @String"`
	Boolean *Boolean `parser:"| @Boolean"`
	Array   []Value  `parser:"| '[' @@ (',' @@)* ']'"`
}

type Boolean bool

func (b *Boolean) Capture(s []string) error {
	switch strings.ToUpper(s[0]) {
	case "TRUE":
		*b = true
	case "FALSE":
		*b = false
	default:
		return fmt.Errorf("unexpected string: %s", s[0])
	}
	return nil
}

type LogicOperator string

const (
	AndLogicOperator LogicOperator = "AND"
	OrLogicOperator  LogicOperator = "OR"
)

func (operator *LogicOperator) Capture(s []string) error {
	switch strings.ToUpper(s[0]) {
	case "AND", "&&":
		*operator = AndLogicOperator
	case "OR", "||":
		*operator = OrLogicOperator
	default:
		return fmt.Errorf("unexpected string: %s", s[0])
	}
	return nil
}

type ComparisonOperator string

const (
	EqualComparisonOperator            ComparisonOperator = "=="
	NotEqualComparisonOperator         ComparisonOperator = "!="
	GreaterThanComparisonOperator      ComparisonOperator = ">"
	GreaterThanEqualComparisonOperator ComparisonOperator = ">="
	LessThanComparisonOperator         ComparisonOperator = "<"
	LessThanEqualComparisonOperator    ComparisonOperator = "<="
	InComparisonOperator               ComparisonOperator = "IN"
)

func comparisonOperatorKeywords() []string {
	return []string{
		string(EqualComparisonOperator),
		string(NotEqualComparisonOperator),
		string(GreaterThanComparisonOperator),
		string(GreaterThanEqualComparisonOperator),
		string(LessThanComparisonOperator),
		string(LessThanEqualComparisonOperator),
		string(InComparisonOperator),
	}
}
