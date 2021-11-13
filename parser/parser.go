package parser

import (
	"errors"
	"strconv"

	"gitlab.com/json_parser/lexer"
)

type Object struct {
	Keys   []string
	Values []interface{}
}

type Parser struct {
	tokens  []lexer.Token
	current lexer.Token
}

func (p *Parser) next() {
	if len(p.tokens) > 1 {
		p.tokens = p.tokens[1:]
		p.current = p.tokens[0]
	}
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: tokens[0],
	}
}



func parseObject(parser *Parser) map[string]interface{} {
	object := make(map[string]interface{}, 0)

	t := parser.current
	if t.Val == string(lexer.JSONBraceRight) {
		return object
	}

	for {
		key := parser.current

		if key.TokenType == lexer.JSONString {
			parser.next()
		} else {
			return object // TODO: add err
		}

		if parser.current.Val != string(lexer.JSONColon) {
			return object
		}
		parser.next()

		value, _ := Parse(parser, false)

		object[key.Val] = value

		parser.next()

		if parser.current.Val == string(lexer.JSONBraceRight) {
			return object
		}
		//else if parser.current.Val != string(lexer.JSONComma) {
		//	return object // TODO: Add error
		//}
	}
}

func parseList(parser *Parser) []interface{} {
	var list []interface{}
	token := parser.current
	if token.Val == string(lexer.JSONBracketRight) {
		return list
	}

	for {
		val, _ := Parse(parser, false)
		list = append(list, val)

		token = parser.current
		if token.Val == string(lexer.JSONBracketRight) {
			return list
		} else if token.Val != string(lexer.JSONComma) {
			return list // TODO: add error
		} else {
			parser.next()
		}
	}
}

func Parse(parser *Parser, isRoot bool) (interface{}, error) {
	token := parser.current
	parser.next()

	if isRoot && token.Val != string(lexer.JSONBraceLeft) {
		return nil, errors.New("JSON must starts with \"")
	}

	if token.Val == string(lexer.JSONBraceLeft) {
		return parseObject(parser), nil
	}

	if token.Val == string(lexer.JSONBracketLeft) {
		return parseList(parser), nil
	}

	switch token.TokenType {
	case lexer.JSONString:
		return token.Val, nil
	case lexer.JSONInt:
		return strconv.Atoi(token.Val)
	case lexer.JSONFloat:
		return strconv.ParseFloat(token.Val, 64)
	case lexer.JSONBool:
		if token.Val == "false" {
			return false, nil
		}
		return true, nil
	case lexer.JSONNull:
		return nil, nil

	}

	if token.TokenType == lexer.JSONInt {
		return strconv.Atoi(token.Val)
	}

	return parser.current, nil
}
