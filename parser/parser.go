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



func (p *Parser) parseObject() map[string]interface{} {
	object := make(map[string]interface{}, 0)

	t := p.current
	if t.Val == string(lexer.JSONBraceRight) {
		return object
	}

	for {
		key := p.current

		if key.TokenType == lexer.JSONString {
			p.next()
		} else {
			return object // TODO: add err
		}

		if p.current.Val != string(lexer.JSONColon) {
			return object
		}
		p.next()

		value, _ := p.Parse(false)

		object[key.Val] = value

		p.next()

		if p.current.Val == string(lexer.JSONBraceRight) {
			return object
		}
		//else if parser.current.Val != string(lexer.JSONComma) {
		//	return object // TODO: Add error
		//}
	}
}

func (p *Parser) parseList() []interface{} {
	var list []interface{}
	token := p.current
	if token.Val == string(lexer.JSONBracketRight) {
		return list
	}

	for {
		val, _ := p.Parse(false)
		list = append(list, val)

		token = p.current
		if token.Val == string(lexer.JSONBracketRight) {
			return list
		} else if token.Val != string(lexer.JSONComma) {
			return list // TODO: add error
		} else {
			p.next()
		}
	}
}

func (p *Parser) Parse(isRoot bool) (interface{}, error) {
	token := p.current
	p.next()

	if isRoot && token.Val != string(lexer.JSONBraceLeft) {
		return nil, errors.New("JSON must starts with \"")
	}

	if token.Val == string(lexer.JSONBraceLeft) {
		return p.parseObject(), nil
	}

	if token.Val == string(lexer.JSONBracketLeft) {
		return p.parseList(), nil
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

	return p.current, nil
}
