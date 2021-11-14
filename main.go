package main

import (
	"fmt"
	"log"

	"github.com/qerdcv/json_parser/lexer"
	"github.com/qerdcv/json_parser/parser"
)

func main() {
	r := NewReader()
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
