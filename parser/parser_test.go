package parser

import (
	"errors"
	"testing"

	"gitlab.com/json_parser/lexer"
)

var (
	objFixture = []lexer.Token{
		{TokenType: lexer.JSONSyntax, Val: "{"},
		{TokenType: lexer.JSONString, Val: "name"},
		{TokenType: lexer.JSONSyntax, Val: ":"},
		{TokenType: lexer.JSONString, Val: "Test"},
		{TokenType: lexer.JSONSyntax, Val: "}"},
	}
	listFixture = []lexer.Token{
		{TokenType: lexer.JSONSyntax, Val: "["},
		{TokenType: lexer.JSONInt, Val: "123"},
		{TokenType: lexer.JSONSyntax, Val: ","},
		{TokenType: lexer.JSONFloat, Val: "12312.3"},
		{TokenType: lexer.JSONSyntax, Val: ","},
		{TokenType: lexer.JSONBool, Val: "true"},
		{TokenType: lexer.JSONSyntax, Val: "]"},
	}
)

func TestNew(t *testing.T) {
	p, err := New([]lexer.Token{})
	if p != nil {
		t.Fatal("expected nil, got ", p)
	}

	if !errors.Is(err, EmptyTokensError) {
		t.Fatalf("expected %e, got %e", EmptyTokensError, err)
	}


	p, err = New([]lexer.Token{
		{},
	})
	if err != nil {
		t.Fatalf("expected nil error, got %e", err)
	}

	if p == nil {
		t.Fatalf("expected parser, got nil")
	}
}

func TestParser_next(t *testing.T) {
	p, _ := New([]lexer.Token{
		{TokenType: lexer.JSONSyntax, Val: "{"},
		{TokenType: lexer.JSONString, Val: "name"},
	})
	exp := lexer.Token{
		TokenType: lexer.JSONSyntax,
		Val: "{",
	}
	got := p.current;

	if exp.TokenType != got.TokenType || exp.Val != got.Val {
		t.Fatalf("expected %v, got %v", exp, got)
	}
	p.next()

	got = p.current
	exp = lexer.Token{
		TokenType: lexer.JSONString,
		Val: "name",
	}

	if exp.TokenType != got.TokenType || exp.Val != got.Val {
		t.Fatalf("expected %v, got %v", exp, got)
	}
	p.next()

	exp = got
	got = p.current

	if exp.TokenType != got.TokenType || exp.Val != got.Val {
		t.Fatalf("expected %v, got %v", exp, got)
	}
}

func	TestParser_parseObject(t *testing.T) {
	p, _ := New(objFixture)

	exp := map[string]interface{}{
		"name": "Test",
	}
	p.next()
	got := p.parseObject()

	for k, v := range exp {
		if val, ok := got[k]; !ok || val != v {
			t.Fatalf("expected %v, got %v", exp, got)
		}
	}
}

func TestParse_parseList(t *testing.T) {
	p, _ := New(listFixture)
	exp := []interface{}{
		123,
		12312.3,
		true,
	}

	p.next()

	got := p.parseList()

	for i := range exp {
		if exp[i] != got[i] {
			t.Fatalf("expected %v, got %v", exp, got)
		}
	}
}

func TestParse_ParseStrign(t *testing.T) {

	p, _ := New([]lexer.Token{
		{TokenType: lexer.JSONString, Val: "Test"},
	})

	exp := "Test"
	got, _ := p.Parse(false)

	if exp != got {
		t.Fatalf("expected %s, got %s", exp, got)
	}
}

func TestParse_ParseInt(t *testing.T) {

	p, _ := New([]lexer.Token{
		{TokenType: lexer.JSONInt, Val: "12312"},
	})

	exp := 12312
	got, _ := p.Parse(false)

	if exp != got {
		t.Fatalf("expected %d, got %d", exp, got)
	}
}

func TestParse_ParseFloat(t *testing.T) {

	p, _ := New([]lexer.Token{
		{TokenType: lexer.JSONFloat, Val: "123.12"},
	})

	exp := 123.12
	got, _ := p.Parse(false)

	if exp != got {
		t.Fatalf("expected %f, got %f", exp, got)
	}
}

func TestParse_ParseNull(t *testing.T) {

	p, _ := New([]lexer.Token{
		{TokenType: lexer.JSONNull, Val: "null"},
	})

	got, _ := p.Parse(false)

	if got != nil {
		t.Fatalf("expected nil, got %f", got)
	}
}

func TestParse_ParseBool(t *testing.T) {

	p, _ := New([]lexer.Token{
		{TokenType: lexer.JSONBool, Val: "true"},
		{TokenType: lexer.JSONBool, Val: "false"},
	})

	exp := true
	got, _ := p.Parse(false)

	if got != exp {
		t.Fatalf("expected %t, got %t", exp, got)
	}

	exp = false
	got, _ = p.Parse(false)


	if got != exp {
		t.Fatalf("expected %t, got %t", exp, got)
	}
}
