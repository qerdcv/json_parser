package main

import (
	"fmt"

	"gitlab.com/json_parser/lexer"
)

func fromString(str string) []lexer.Token {
	tokens := lexer.Lexer(str)
	return tokens
}

func main() {
	fmt.Println(fromString("{\"Jovani Jorjo\": 1231212, \"adult\": true, \"hasChild\": null}"))
}
