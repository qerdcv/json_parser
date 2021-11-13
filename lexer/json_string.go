package lexer

import "errors"


var (
	EndOfFileErrors = errors.New("End of file")
)


type JString struct {
	current rune
	s string
}

func New(str string) *JString {
	return &JString {
		current: rune(str[0]),
		s: str,
	}
}


func (jstring *JString) next() (rune, error) {
	if len(jstring.s) > 1 {
		jstring.s = jstring.s[1:]
		jstring.current = rune(jstring.s[0])
		return jstring.current, nil
	}
	jstring.current = rune(jstring.s[0])
	jstring.s = ""
	return jstring.current, EndOfFileErrors
}

func (jstring *JString) Cut(l int) {
	jstring.s = jstring.s[l:]
	if len(jstring.s) != 0 {
		jstring.current = rune(jstring.s[0])
	}
	// TODO: think about error
}
