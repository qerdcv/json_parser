package lexer

import (
	"errors"
	"fmt"
	"log"
	"os"
)

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
	Val string
}


type JsonToken int

const (
	JSONFloat JsonToken = iota
	JSONString
	JSONInt   
	JSONBool  
	JSONNull
	JSONSyntax
)

func contains(ch rune, slice []rune) bool {
	for _, el := range slice {
		if el == ch {
			return true
		}
	}
	return false
}

func lexString(str *JString) (*Token, error) {
	jsonString := ""

	if str.current == JSONQuote {
		str.next()
	} else {
		return nil, nil
	}

	for {
		ch := str.current
		if ch == JSONQuote {
			str.next()
			return &Token{TokenType: JSONString, Val: jsonString}, nil
		}
		jsonString += string(ch)

		_, err := str.next()
		if errors.Is(err, EndOfFileErrors) {
			break
		}
	}
	// log.Fatal("String must be ended with quote")
	return nil, fmt.Errorf("string must be ended with \";\n%w", EndOfFileErrors) // TODO: add error
}

func lexNumber(str *JString) (*Token, error) {
	jsonNumber := ""
	isFloat := false

	for {
		ch := str.current;
		if contains(ch, numberChr) {
			if ch == '.' {
				isFloat = true
			}
			jsonNumber += string(ch)
		} else {
			break
		}
		_, err := str.next()
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

func lexBool(str *JString) *Token {
	strLen := len(str.s)

	if strLen >= trueLen && str.s[:trueLen] == "true" {
		str.Cut(trueLen)
		return &Token{TokenType: JSONBool, Val: "true"}
	} else if strLen >= falseLen && str.s[:falseLen] == "false" {
		str.Cut(falseLen)
		return &Token{TokenType: JSONBool, Val: "false"}
	}
	return nil
}

func lexNull(str *JString) *Token {
	strLen := len(str.s)

	if strLen >= nullLen && str.s[:nullLen] == "null" {
		str.Cut(nullLen)
		return &Token{TokenType: JSONNull, Val: "null"}
	}

	return nil
}

func Lexer(str *JString) ([]Token, error) {
	tokens := make([]Token, 0)
	for len(str.s) != 0 {
		//log.Println(len(str.s), str.s)
		// ---- string ----
		jsonString, err := lexString(str)
		if err != nil {
			return nil, err
		}

		if jsonString != nil {
			tokens = append(tokens, *jsonString)
		}

		// ---- number ----
		jsonNumber, _ := lexNumber(str)
		if jsonNumber != nil {
			tokens = append(tokens, *jsonNumber)
		}

		// ---- bool -----
		jsonBool := lexBool(str)
		if jsonBool != nil {
			tokens = append(tokens, *jsonBool)
		}

		// ---- null ----
		jsonNull := lexNull(str)

		if jsonNull != nil {
			tokens = append(tokens, *jsonNull)
		}
	
		// ---- json syntax ----
		ch := str.current
		if contains(ch, jsonWhitespace) {
			str.next()
		} else if contains(ch, jsonSyntax) {
			tokens = append(tokens, Token{TokenType: JSONSyntax, Val: string(ch)})
			str.next()
		} else {
			log.Printf("Invalid char %c", ch)
			os.Exit(1)
		}
	}
	return tokens, nil
}
