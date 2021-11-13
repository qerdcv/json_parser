package main

import (
	"fmt"
	"gitlab.com/json_parser/lexer"
	"gitlab.com/json_parser/parser"
)

func main() {
	strToParse := "{\"name\": \"Vadym\", \"age\": 20, \"gender\": \"male\", \"married\": false, \"kids\": null, \"languages\": [\"Python\", \"JavaScript\", \"Golang\"]}"
	l := lexer.New(strToParse)
	tokens, _ := l.Lex()
	p := parser.New(tokens)
	obj, _ := p.Parse(true)
	fmt.Println(obj)
}
