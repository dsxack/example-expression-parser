# example-statement-parser

An statement parser with [participle](https://github.com/alecthomas/participle).

### Install

```sh
go get github.com/dsxack/example-statement-parser
```

### Usage example

```golang
package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/dsxack/example-statement-parser"
)

func Example() {
	stmt := parser.Statement{}
	err := parser.Parse(`#{ ${insecure} ? "http" : "https" }://${domain}/${basepath}`, &stmt)
	if err != nil {
		log.Fatal(err)
	}

	stringP := func(value string) *string { return &value }

	result := reflect.DeepEqual(stmt, parser.Statement{
		Fragments: []parser.Fragment{
			{IfStatement: &parser.IfStatement{
				Cond: parser.Expr{Left: &parser.Cond{Left: parser.Term{Variable: stringP("insecure")}}},
				Then: parser.Term{Value: &parser.Value{String: stringP("http")}},
				Else: parser.Term{Value: &parser.Value{String: stringP("https")}},
			}},
			{Literal: stringP("://")},
			{Variable: stringP("domain")},
			{Literal: stringP("/")},
			{Variable: stringP("basepath")},
		},
	})

	fmt.Println(result)

	// output: true
}
```

### Example of statements

* `some value 5"`
* `${some.variable}"`
* `#{ true ? "yes" : "no" }`
* `#{ ${insecure} ? ${var1} : ${var2} }`
* `#{ true ? true : false }`
* `#{ 5 ? 10 : 0 }`
* `#{ ${cond1} or ${cond2} ? ${var1} : ${var2} }`
* `#{ ${cond1} || ${cond2} ? ${var1} : ${var2} }`
* `#{ ${cond1} and ${cond2} ? ${var1} : ${var2} }`
* `#{ ${cond1} && ${cond2} ? ${var1} : ${var2} }`
* `#{ ${foo} == "bar" ? ${var1} : ${var2} }`
* `#{ (${cond1} or ${cond2}) and ${cond3} ? ${var1} : ${var2} }`
* `#{ ${cond1} and (${cond2} or ${cond3}) ? ${var1} : ${var2} }`
* `#{ ${cond1} and (${cond2} > 5 or ${cond3}) ? ${var1} : ${var2} }`
* `#{ ${insecure} ? "http" : "https" }://${domain}/${basepath}`

### Parser railway diagram

![](.github/images/railroad.png)
