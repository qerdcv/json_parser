package lexer

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	l, err := New("")

	if l != nil {
		t.Fatal("lexer must be nil on empty string")
	}

	if !errors.Is(err, EmptyStringError) {
		t.Fatalf("error must by %e, got %e", EmptyStringError, err)
	}


	l, err = New("ABCD")
	if got, exp := l.s, "ABCD"; got != exp {
		t.Fatalf("expected %s, got %s", exp, got)
	}

	if got, exp := l.current, 'A'; got != exp {
		t.Fatalf("expected %c, got %c", exp, got)
	}
}


func TestLexer_Cut(t *testing.T) {
	l, err := New("ABCD")
	if err != nil {
		t.Fatalf("unexpected error %e", err)
	}

	l.Cut(2)

	if got, exp := l.s, "CD"; got != exp {
		t.Fatalf("expected %s, got %s", exp, got)
	}
}

func TestLexer_lexBool(t *testing.T) {
	l, _ := New("truefalse")
	got := l.lexBool()
	exp := &Token{
		TokenType: JSONBool,
		Val: "true",
	}
	if got.Val != exp.Val || got.TokenType != exp.TokenType {
		t.Fatalf("expected %v, got %v", exp, got)
	}

	got = l.lexBool()
	exp = &Token{
		TokenType: JSONBool,
		Val: "false",
	}

	if got.Val != exp.Val || got.TokenType != exp.TokenType {
		t.Fatalf("expected %v, got %v", exp, got)
	}
}

func TestLexer_lexBoolNil(t *testing.T) {
	l, _ := New("null")
	got := l.lexBool()
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}


func TestLexer_lexNull(t *testing.T) {
	l, _ := New("null")
	got := l.lexNull()
	exp := &Token{
		TokenType: JSONNull,
		Val: "null",
	}
	if got.Val != exp.Val || got.TokenType != exp.TokenType {
		t.Fatalf("expected %v, got %v", exp, got)
	}
}

func TestLexer_lexString(t *testing.T) {
	l, _ := New("\"Valera\"")
	got, _ := l.lexString()
	exp := &Token{
		TokenType: JSONString,
		Val: "Valera",
	}
	if got.Val != exp.Val || got.TokenType != exp.TokenType {
		t.Fatalf("expected %v, got %v", exp, got)
	}
}

func TestLexer_lexStringErr(t *testing.T) {
	l, _ := New("\"Valera")
	got, err := l.lexString()

	if got != nil || err == nil {
		t.Fatalf("expected error, got %v, nil", got)
	}
}

func TestLexer_lexStringNil(t *testing.T) {
	l, _ := New("true")
	got, _ := l.lexString()

	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestLexer_lexNumberInt(t *testing.T) {
	l, _ := New("12")
	got, _ := l.lexNumber()
	exp := &Token{
		TokenType: JSONInt,
		Val: "12",
	}
	if got.Val != exp.Val || got.TokenType != exp.TokenType {
		t.Fatalf("expected %v, got %v", exp, got)
	}
}

func TestLexer_lexNumberNil(t *testing.T) {
	l, _ := New("null")
	got, _ := l.lexNumber()
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}


func TestLexer_lexNumberFloat(t *testing.T) {
	l, _ := New("12.12")
	got, _ := l.lexNumber()
	exp := &Token{
		TokenType: JSONFloat,
		Val: "12.12",
	}
	if got.Val != exp.Val || got.TokenType != exp.TokenType {
		t.Fatalf("expected %v, got %v", exp, got)
	}
}


func TestLexer_next(t *testing.T) {
	l, _ := New("12")

	if l.current != '1' {
		t.Fatalf("expected '1', got %c", l.current)
	}

	c, _ := l.next()

	if c != l.current {
		t.Fatalf("expected %c, got %c", c, l.current)
	}

	newC, err := l.next()

	if newC != c {
		t.Fatalf("expected %c, got %c", newC, c)
	}

	if err == nil || !errors.Is(err, EndOfStringErrors) {
		t.Fatalf("expected %e, got %e", EndOfStringErrors, err)
	}
}

func TestLexer_Lex(t *testing.T) {
	expected := []Token{
		{TokenType: JSONBool, Val: "false"},
		{TokenType: JSONSyntax, Val: ","},
		{TokenType: JSONBool, Val: "true"},
		{TokenType: JSONSyntax, Val: ","},
		{TokenType: JSONNull, Val: "null"},
		{TokenType: JSONSyntax, Val: ","},
		{TokenType: JSONInt, Val: "12"},
		{TokenType: JSONSyntax, Val: ","},
		{TokenType: JSONFloat, Val: "12.12"},
		{TokenType: JSONSyntax, Val: ","},
		{TokenType: JSONString, Val: "String"},
		{TokenType: JSONSyntax, Val: ","},
	}
	l, _ := New("false,true,null,12,12.12,\"String\",")

	got, _ := l.Lex()

	for i := 0; i < len(expected); i++ {
		if expected[i].Val != got[i].Val || expected[i].TokenType != got[i].TokenType {
			t.Fatalf("expected %v, got %v", expected[i], got[i])
		}
	}

}