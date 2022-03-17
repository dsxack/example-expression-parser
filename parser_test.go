package parser

import (
	"fmt"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *Statement
	}{
		{
			name:  "literal",
			input: "some value 5", want: &Statement{
				Fragments: []Fragment{
					{Literal: stringP("some value 5")},
				},
			}},
		{
			name:  "variable",
			input: "${some.variable}", want: &Statement{
				Fragments: []Fragment{
					{Variable: stringP("some.variable")},
				},
			},
		},
		{
			name:  "if statement",
			input: `#{ true ? "yes" : "no" }`,
			want: &Statement{
				Fragments: []Fragment{
					{IfStatement: &IfStatement{
						Cond: Expr{Left: &Cond{Left: Term{Value: &Value{Boolean: boolP(true)}}}},
						Then: Term{Value: &Value{String: stringP("yes")}},
						Else: Term{Value: &Value{String: stringP("no")}},
					}},
				},
			},
		},
		{
			name:  "if statement variables",
			input: `#{ ${insecure} ? ${var1} : ${var2} }`,
			want: &Statement{
				Fragments: []Fragment{
					{IfStatement: &IfStatement{
						Cond: Expr{Left: &Cond{Left: Term{Variable: stringP("insecure")}}},
						Then: Term{Variable: stringP("var1")},
						Else: Term{Variable: stringP("var2")},
					}},
				},
			},
		},
		{
			name:  "if statement booleans",
			input: `#{ true ? true : false }`,
			want: &Statement{
				Fragments: []Fragment{
					{IfStatement: &IfStatement{
						Cond: Expr{Left: &Cond{Left: Term{Value: &Value{Boolean: boolP(true)}}}},
						Then: Term{Value: &Value{Boolean: boolP(true)}},
						Else: Term{Value: &Value{Boolean: boolP(false)}},
					}},
				},
			},
		},
		{
			name:  "if statement numbers",
			input: `#{ 5 ? 10 : 0 }`,
			want: &Statement{
				Fragments: []Fragment{
					{IfStatement: &IfStatement{
						Cond: Expr{Left: &Cond{Left: Term{Value: &Value{Number: float64P(5)}}}},
						Then: Term{Value: &Value{Number: float64P(10)}},
						Else: Term{Value: &Value{Number: float64P(0)}},
					}},
				},
			},
		},
		{
			name:  "if statement OR operator",
			input: `#{ ${cond1} or ${cond2} ? ${var1} : ${var2} }`,
			want: &Statement{
				Fragments: []Fragment{
					{IfStatement: &IfStatement{
						Cond: Expr{
							Left: &Cond{Left: Term{Variable: stringP("cond1")}},
							Right: []LogicExpr{
								{
									Operator: OrLogicOperator,
									Cond:     &Cond{Left: Term{Variable: stringP("cond2")}},
								},
							},
						},
						Then: Term{Variable: stringP("var1")},
						Else: Term{Variable: stringP("var2")},
					}},
				},
			},
		},
		{
			name:  "if statement OR operator ||",
			input: `#{ ${cond1} || ${cond2} ? ${var1} : ${var2} }`,
			want: &Statement{
				Fragments: []Fragment{
					{IfStatement: &IfStatement{
						Cond: Expr{
							Left: &Cond{Left: Term{Variable: stringP("cond1")}},
							Right: []LogicExpr{
								{
									Operator: OrLogicOperator,
									Cond:     &Cond{Left: Term{Variable: stringP("cond2")}},
								},
							},
						},
						Then: Term{Variable: stringP("var1")},
						Else: Term{Variable: stringP("var2")},
					}},
				},
			},
		},
		{
			name:  "if statement AND operator",
			input: `#{ ${cond1} and ${cond2} ? ${var1} : ${var2} }`,
			want: &Statement{
				Fragments: []Fragment{
					{IfStatement: &IfStatement{
						Cond: Expr{
							Left: &Cond{Left: Term{Variable: stringP("cond1")}},
							Right: []LogicExpr{
								{
									Operator: AndLogicOperator,
									Cond:     &Cond{Left: Term{Variable: stringP("cond2")}},
								},
							},
						},
						Then: Term{Variable: stringP("var1")},
						Else: Term{Variable: stringP("var2")},
					}},
				},
			},
		},
		{
			name:  "if statement AND operator &&",
			input: `#{ ${cond1} && ${cond2} ? ${var1} : ${var2} }`,
			want: &Statement{
				Fragments: []Fragment{
					{IfStatement: &IfStatement{
						Cond: Expr{
							Left: &Cond{Left: Term{Variable: stringP("cond1")}},
							Right: []LogicExpr{
								{
									Operator: AndLogicOperator,
									Cond:     &Cond{Left: Term{Variable: stringP("cond2")}},
								},
							},
						},
						Then: Term{Variable: stringP("var1")},
						Else: Term{Variable: stringP("var2")},
					}},
				},
			},
		},
		{
			name:  "if statement comparison operator",
			input: `#{ ${foo} == "bar" ? ${var1} : ${var2} }`,
			want: &Statement{
				Fragments: []Fragment{
					{IfStatement: &IfStatement{
						Cond: Expr{
							Left: &Cond{
								Left:     Term{Variable: stringP("foo")},
								Operator: EqualComparisonOperator,
								Right:    &Term{Value: &Value{String: stringP("bar")}},
							},
						},
						Then: Term{Variable: stringP("var1")},
						Else: Term{Variable: stringP("var2")},
					}},
				},
			},
		},
		{
			name:  "if statement left subexpression",
			input: `#{ (${cond1} or ${cond2}) and ${cond3} ? ${var1} : ${var2} }`,
			want: &Statement{
				Fragments: []Fragment{
					{IfStatement: &IfStatement{
						Cond: Expr{
							Sub: &Expr{
								Left: &Cond{Left: Term{Variable: stringP("cond1")}},
								Right: []LogicExpr{{
									Operator: OrLogicOperator,
									Cond:     &Cond{Left: Term{Variable: stringP("cond2")}},
								}},
							},
							Right: []LogicExpr{{
								Operator: AndLogicOperator,
								Cond:     &Cond{Left: Term{Variable: stringP("cond3")}},
							}},
						},
						Then: Term{Variable: stringP("var1")},
						Else: Term{Variable: stringP("var2")},
					}},
				},
			},
		},
		{
			name:  "if statement right subexpression",
			input: `#{ ${cond1} and (${cond2} or ${cond3}) ? ${var1} : ${var2} }`,
			want: &Statement{
				Fragments: []Fragment{
					{IfStatement: &IfStatement{
						Cond: Expr{
							Left: &Cond{Left: Term{Variable: stringP("cond1")}},
							Right: []LogicExpr{{
								Operator: AndLogicOperator,
								Sub: &Expr{
									Left: &Cond{Left: Term{Variable: stringP("cond2")}},
									Right: []LogicExpr{{
										Operator: OrLogicOperator,
										Cond:     &Cond{Left: Term{Variable: stringP("cond3")}},
									}},
								},
							}},
						},
						Then: Term{Variable: stringP("var1")},
						Else: Term{Variable: stringP("var2")},
					}},
				},
			},
		},
		{
			name:  "if statement comparison operator in subexpression",
			input: `#{ ${cond1} and (${cond2} > 5 or ${cond3}) ? ${var1} : ${var2} }`,
			want: &Statement{
				Fragments: []Fragment{
					{IfStatement: &IfStatement{
						Cond: Expr{
							Left: &Cond{Left: Term{Variable: stringP("cond1")}},
							Right: []LogicExpr{{
								Operator: AndLogicOperator,
								Sub: &Expr{
									Left: &Cond{
										Left:     Term{Variable: stringP("cond2")},
										Operator: GreaterThanComparisonOperator,
										Right:    &Term{Value: &Value{Number: float64P(5)}},
									},
									Right: []LogicExpr{{
										Operator: OrLogicOperator,
										Cond:     &Cond{Left: Term{Variable: stringP("cond3")}},
									}},
								},
							}},
						},
						Then: Term{Variable: stringP("var1")},
						Else: Term{Variable: stringP("var2")},
					}},
				},
			},
		},
		{
			name:  "combine different type of fragments",
			input: `#{ ${insecure} ? "http" : "https" }://${domain}/${basepath}`,
			want: &Statement{
				Fragments: []Fragment{
					{IfStatement: &IfStatement{
						Cond: Expr{Left: &Cond{Left: Term{Variable: stringP("insecure")}}},
						Then: Term{Value: &Value{String: stringP("http")}},
						Else: Term{Value: &Value{String: stringP("https")}},
					}},
					{Literal: stringP("://")},
					{Variable: stringP("domain")},
					{Literal: stringP("/")},
					{Variable: stringP("basepath")},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lex, _ := statementLexer.LexString("", test.input)
			tokens, _ := lexer.ConsumeAll(lex)
			fmt.Println(tokens)

			stmt := &Statement{}
			err := Parse(test.input, stmt)
			require.NoError(t, err)
			if test.want != nil {
				require.Equal(t, test.want, stmt)
			}
		})
	}
}

func boolP(value bool) *Boolean       { return (*Boolean)(&value) }
func float64P(value float64) *float64 { return &value }
func stringP(value string) *string    { return &value }
