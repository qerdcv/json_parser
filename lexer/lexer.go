package lexer

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type JsonToken int

const (
	JSONColon        = ':'
	JSONBracketRight = ']'
	JSONBracketLeft  = '['
	JSONBraceLeft    = '{'
	JSONBraceRight   = '}'
	JSONComma        = ','
	JSONQuote        = '"'
	falseLen         = 5
	trueLen          = 4
	nullLen          = 4
)

const (
	JSONFloat JsonToken = iota
	JSONString
	JSONInt
	JSONBool
	JSONNull
	JSONSyntax
)

var (
	jsonWhitespace = []rune{' ', '\t', '\b', '\n', '\r'}
	jsonSyntax     = []rune{
		JSONComma,
		JSONColon,
		JSONBracketLeft,
		JSONBracketRight,
		JSONBraceLeft,
		JSONBraceRight,
	}
	numberChr = []rune{
		'0',
		'1',
		'2',
		'3',
		'4',
		'5',
		'6',
		'7',
		'8',
		'9',
		'-',
		'e',
		'.',
	}
)

type Token struct {
	TokenType JsonToken
	Val       string
}

var (
	EndOfFileErrors = errors.New("end of file")

	EmptyStringError = errors.New("empty string to parse")
)

type Lexer struct {
	current rune
	s       string
}

func New(str string) (*Lexer, error) {
	if len(str) == 0 {
		return nil, EmptyStringError
	}
	return &Lexer{
		current: rune(str[0]),
		s:       str,
	}, nil
}

func (l *Lexer) next() (rune, error) {
	if len(l.s) > 1 {
		l.s = l.s[1:]
		l.current = rune(l.s[0])
		return l.current, nil
	}
	l.current = rune(l.s[0])
	l.s = ""
	return l.current, EndOfFileErrors
}

func (l *Lexer) Cut(lenToCut int) {
	l.s = l.s[lenToCut:]
	if len(l.s) != 0 {
		l.current = rune(l.s[0])
	}
	// TODO: think about error
}

func contains(ch rune, slice []rune) bool {
	for _, el := range slice {
		if el == ch {
			return true
		}
	}
	return false
}

func (l *Lexer) lexString() (*Token, error) {
	jsonString := ""

	if l.current == JSONQuote {
		l.next()
	} else {
		return nil, nil
	}

	for {
		ch := l.current
		if ch == JSONQuote {
			l.next()
			return &Token{TokenType: JSONString, Val: jsonString}, nil
		}
		jsonString += string(ch)

		_, err := l.next()
		if errors.Is(err, EndOfFileErrors) {
			break
		}
	}
	// log.Fatal("String must be ended with quote")
	return nil, fmt.Errorf("string must be ended with \";\n%w", EndOfFileErrors) // TODO: add error
}

func (l *Lexer) lexNumber() (*Token, error) {
	jsonNumber := ""
	isFloat := false

	for {
		ch := l.current
		if contains(ch, numberChr) {
			if ch == '.' {
				isFloat = true
			}
			jsonNumber += string(ch)
		} else {
			break
		}
		_, err := l.next()
		if err != nil {
			break
		}
	}

	if jsonNumber == "" {
		return nil, nil
	}

	if isFloat {
		return &Token{TokenType: JSONFloat, Val: jsonNumber}, nil
	}

	return &Token{TokenType: JSONInt, Val: jsonNumber}, nil
}

func (l *Lexer) lexBool() *Token {
	strLen := len(l.s)

	if strLen >= trueLen && l.s[:trueLen] == "true" {
		l.Cut(trueLen)
		return &Token{TokenType: JSONBool, Val: "true"}
	} else if strLen >= falseLen && l.s[:falseLen] == "false" {
		l.Cut(falseLen)
		return &Token{TokenType: JSONBool, Val: "false"}
	}
	return nil
}

func (l *Lexer) lexNull() *Token {
	strLen := len(l.s)

	if strLen >= nullLen && l.s[:nullLen] == "null" {
		l.Cut(nullLen)
		return &Token{TokenType: JSONNull, Val: "null"}
	}

	return nil
}

func (l *Lexer) lex() ([]Token, error) {
	tokens := make([]Token, 0)
	for len(l.s) != 0 {
		//log.Println(len(str.s), str.s)
		// ---- string ----
		jsonString, err := l.lexString()
		if err != nil {
			return nil, err
		}

		if jsonString != nil {
			tokens = append(tokens, *jsonString)
		}

		// ---- number ----
		jsonNumber, _ := l.lexNumber()
		if jsonNumber != nil {
			tokens = append(tokens, *jsonNumber)
		}

		// ---- bool -----
		jsonBool := l.lexBool()
		if jsonBool != nil {
			tokens = append(tokens, *jsonBool)
		}

		// ---- null ----
		jsonNull := l.lexNull()

		if jsonNull != nil {
			tokens = append(tokens, *jsonNull)
		}

		// ---- json syntax ----
		ch := l.current
		if contains(ch, jsonWhitespace) {
			l.next()
		} else if contains(ch, jsonSyntax) {
			tokens = append(tokens, Token{TokenType: JSONSyntax, Val: string(ch)})
			l.next()
		} else {
			log.Printf("Invalid char %c", ch)
			os.Exit(1)
		}
	}
	return tokens, nil
}
