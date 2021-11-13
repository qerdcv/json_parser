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

	l, err = New("a")

	if err != nil {
		t.Fatalf("error %e", err)
	}

	if l.current != 'a' {
		t.Fatalf("expected %c, got %c", 'a', l.current)
	}
}
