package lexer

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type JsonToken int

const (
	JSONColon        = ":"
	JSONBracketRight = "]"
	JSONBracketLeft  = "["
	JSONBraceLeft    = "{"
	JSONBraceRight   = "}"
	JSONComma        = ","
	JSONQuote        = "\""
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
	jsonWhitespace = []string{
		" ",
		"\t",
		"\b",
		"\n",
		"\r",
	}
	jsonSyntax     = []string{
		JSONComma,
		JSONColon,
		JSONBracketLeft,
		JSONBracketRight,
		JSONBraceLeft,
		JSONBraceRight,
	}
	numberChr = []string{
		"0",
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"-",
		"e",
		".",
	}
)

type Token struct {
	TokenType JsonToken
	Val       string
}

var (
	ErrEndOfString = errors.New("end of string")

	ErrEmptyString = errors.New("empty string to parse")
)

type Lexer struct {
	current string
	s       string
}

func New(str string) (*Lexer, error) {
	if str == "" {
		return nil, ErrEmptyString
	}
	return &Lexer{
		current: string(str[0]),
		s:       str,
	}, nil
}

func (l *Lexer) next() error {
	if len(l.s) > 1 {
		l.s = l.s[1:]
		l.current = string(l.s[0])
		return nil
	}
	l.current = string(l.s[0])
	l.s = ""
	return ErrEndOfString
}

func (l *Lexer) Cut(lenToCut int) {
	l.s = l.s[lenToCut:]
	if len(l.s) != 0 {
		l.current = string(l.s[0])
	}
	// TODO: think about error
}

func contains(ch string, slice []string) bool {
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

		if err := l.next(); err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	for {
		ch := l.current
		if string(ch) == JSONQuote {
			if err := l.next(); err != nil {
				log.Println(err)
			}

			return &Token{TokenType: JSONString, Val: jsonString}, nil
		}
		jsonString += string(ch)

		err := l.next()
		if errors.Is(err, ErrEndOfString) {
			break
		}
	}
	return nil, fmt.Errorf("string must be ended with \";\n%w", ErrEndOfString)
}

func (l *Lexer) lexNumber() (*Token, error) {
	jsonNumber := ""
	isFloat := false

	for {
		ch := l.current
		if !contains(ch, numberChr) {
			break
		}
		if ch == "." {
			isFloat = true
		}
		jsonNumber += string(ch)
		if err := l.next(); err != nil {
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

func (l *Lexer) Lex() ([]Token, error) {
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
			if err := l.next(); err != nil {
				log.Println(err)
			}
		} else if contains(ch, jsonSyntax) {
			tokens = append(tokens, Token{TokenType: JSONSyntax, Val: string(ch)})
			if err := l.next(); err != nil {
				log.Println(err)
			}

		} else {
			log.Printf("Invalid char %s", ch)
			os.Exit(1)
		}
	}
	return tokens, nil
}
