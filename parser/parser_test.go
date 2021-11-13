package parser

import (
	"errors"
	"gitlab.com/json_parser/lexer"
	"testing"
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