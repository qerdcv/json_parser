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
	l, _ := New("ABCD")

	l.Cut(2)

	if got, exp := l.s, "CD"; got != exp {
		t.Fatalf("expected %s, got %s", exp, got)
	}
}

func TestLexer_LexBool(t *testing.T) {
	l, _ := New("truefalse")
	got := l.lexBool()
	exp := &Token{
		TokenType: JSONBool,
		Val: "true",
	}
	if got.Val != exp.Val && got.TokenType != exp.TokenType {
		t.Fatalf("expected %v, got %v", exp, got)
	}

	got = l.lexBool()
	exp = &Token{
		TokenType: JSONBool,
		Val: "false",
	}

	if got.Val != exp.Val && got.TokenType != exp.TokenType {
		t.Fatalf("expected %v, got %v", exp, got)
	}
}


func TestLexer_LexNull(t *testing.T) {
	l, _ := New("null")
	got := l.lexNull()
	exp := &Token{
		TokenType: JSONNull,
		Val: "null",
	}
	if got.Val != exp.Val && got.TokenType != exp.TokenType {
		t.Fatalf("expected %v, got %v", exp, got)
	}
}

func TestLexer_LexString(t *testing.T) {
	l, _ := New("\"Valera\"")
	got, _ := l.lexString()
	exp := &Token{
		TokenType: JSONString,
		Val: "Valera",
	}
	if got.Val != exp.Val && got.TokenType != exp.TokenType {
		t.Fatalf("expected %v, got %v", exp, got)
	}
}

func TestLexer_LexNumberInt(t *testing.T) {
	l, _ := New("12")
	got, _ := l.lexNumber()
	exp := &Token{
		TokenType: JSONInt,
		Val: "12",
	}
	if got.Val != exp.Val && got.TokenType != exp.TokenType {
		t.Fatalf("expected %v, got %v", exp, got)
	}
}


func TestLexer_LexNumberFloat(t *testing.T) {
	l, _ := New("12.12")
	got, _ := l.lexNumber()
	exp := &Token{
		TokenType: JSONInt,
		Val: "12.12",
	}
	if got.Val != exp.Val && got.TokenType != exp.TokenType {
		t.Fatalf("expected %v, got %v", exp, got)
	}
}
