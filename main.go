package main

import (
	"fmt"
	"gitlab.com/json_parser/lexer"
	"gitlab.com/json_parser/parser"
	"gitlab.com/json_parser/reader"
	"log"
)

func main() {
	r := reader.New()
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
