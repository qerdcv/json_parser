package main

import (
	"fmt"
	"log"

	"gitlab.com/json_parser/lexer"
	"gitlab.com/json_parser/parser"
)

func main() {
	r := New()
	l, err := lexer.New(r.Read())
	if err != nil {
		log.Fatal(err)
	}
	tokens, _ := l.Lex()
	p, err := parser.New(tokens)
	if err != nil {
		log.Fatal(err)
	}

	obj, _ := p.Parse(true)
	fmt.Println(obj)
}
