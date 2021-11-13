package parser

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"gitlab.com/json_parser/lexer"
)

var (
	EmptyTokensError = errors.New("expected non-empty list of tokens")
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

func New(tokens []lexer.Token) (*Parser, error) {
	if len(tokens) == 0 {
		return nil, EmptyTokensError
	}
	return &Parser{
		tokens:  tokens,
		current: tokens[0],
	}, nil
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

		current := p.current.Val
		p.next()

		if current == string(lexer.JSONBraceRight) {
			return object
		} else if current != string(lexer.JSONComma) {
			log.Fatal(fmt.Sprintf("expected comma have %s", current))
		}
	}
}

func (p *Parser) parseList() []interface{} {
	var list []interface{}
	token := p.current
	if token.Val == string(lexer.JSONBracketRight) {
		p.next()
		return list
	}

	for {
		val, _ := p.Parse(false)
		list = append(list, val)

		token = p.current

		p.next()
		if token.Val == string(lexer.JSONBracketRight) {
			return list
		} else if token.Val != string(lexer.JSONComma) {
			log.Fatal(fmt.Sprintf("expected comma have %s", token.Val))
		}
	}
}

func (p *Parser) Parse(isRoot bool) (interface{}, error) {
	token := p.current
	p.next()

	if isRoot && (token.Val != string(lexer.JSONBraceLeft) && token.Val != string(lexer.JSONBracketLeft)) {
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
