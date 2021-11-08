package lexer

import (
	"log"
	"os"
	"strconv"
	"time"
)

const (
	jsonColon        = ':'
	jsonBracketRight = ']'
	jsonBracketLeft  = '['
	jsonBraceLeft    = '{'
	jsonBraceRight   = '}'
	jsonComma        = ','
	jsonQuote        = '"'
	falseLen         = 5
	trueLen          = 4
	nullLen          = 4
)

var (
	jsonWhitespace = []rune{' ', '\t', '\b', '\n', '\r'}
	jsonSyntax     = []rune{
		jsonComma,
		jsonColon,
		jsonBracketLeft,
		jsonBracketRight,
		jsonBraceLeft,
		jsonBraceRight,
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
	JSONFloat  float64
	JSONString string
	JSONInt    int
	JSONBool   bool
	JSONNull   interface{}
	JSONArray  []Token
	JSONSyntax rune
}

func contains(ch rune, slice []rune) bool {
	for _, el := range slice {
		if el == ch {
			return true
		}
	}
	return false
}

func lexString(str string) (*Token, string) {
	jsonString := ""

	if str[0] == jsonQuote {
		str = str[1:]
	} else {
		return nil, str
	}

	for _, ch := range str {
		if ch == jsonQuote {
			return &Token{JSONString: jsonString}, str[len(jsonString)+1:]
		}
		jsonString += string(ch)
	}
	log.Fatal("String must be ended with quote")
	return nil, str
}

func lexNumber(str string) (*Token, string) {
	jsonNumber := ""
	isFloat := false

	for _, ch := range str {
		if contains(ch, numberChr) {
			if ch == '.' {
				isFloat = true
			}
			jsonNumber += string(ch)
		} else {
			break
		}
	}

	rest := str[len(jsonNumber):]

	if jsonNumber == "" {
		return nil, str
	}

	if isFloat {
		number, _ := strconv.ParseFloat(jsonNumber, 64)
		return &Token{JSONFloat: number}, rest
	}
	number, _ := strconv.Atoi(jsonNumber)
	return &Token{JSONInt: number}, rest
}

func lexBool(str string) (*Token, string) {
	strLen := len(str)

	if strLen >= trueLen && str[:trueLen] == "true" {
		return &Token{JSONBool: true}, str[trueLen:]
	} else if strLen >= falseLen && str[:falseLen] == "false" {
		return &Token{JSONBool: false}, str[falseLen:]
	}

	return nil, str
}

func lexNull(str string) (*Token, string) {
	strLen := len(str)

	if strLen >= nullLen && str[:nullLen] == "null" {
		return &Token{JSONNull: nil}, str[trueLen:]
	}

	return nil, str
}

func Lexer(str string) []Token {
	tokens := make([]Token, 0)
	for {
		log.Println(str)
		// ---- string ----
		jsonString, str := lexString(str)
		if jsonString != nil {
			tokens = append(tokens, *jsonString)
		}

		// ---- number ----
		//jsonNumber, str := lexNumber(str)
		//if jsonNumber != nil {
		//	tokens = append(tokens, *jsonNumber)
		//}

		// ---- bool -----
		jsonBool, str := lexBool(str)
		if jsonBool != nil {
			tokens = append(tokens, *jsonBool)
		}

		// ---- null ----
		jsonNull, str := lexNull(str)

		if jsonNull != nil {
			tokens = append(tokens, *jsonNull)
		}

		if str == "" {
			break
		}

		ch := rune(str[0])
		log.Println(contains(ch, jsonSyntax))
		if contains(ch, jsonWhitespace) {
			str = str[1:]
		} else if contains(ch, jsonSyntax) {
			tokens = append(tokens, Token{JSONSyntax: ch})
			str = str[1:]
			log.Println(str, str[1:])
		} else {
			log.Printf("Invalid char %c", ch)
			os.Exit(1)
		}
		log.Println(tokens, " - tokens")
		time.Sleep(1 * time.Second)
	}
	return tokens
}
