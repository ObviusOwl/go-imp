package lexer

import (
	"fmt"
	"strings"
)

type Token struct {
	Name  string
	Value string
	Line  int
}

type Position struct {
	File   string
	Index  int
	Line   int
	Column int
}

func (tok Token) String() string {
	return fmt.Sprintf("<%s '%s'>", tok.Name, tok.Value)
}

type TokenMatcher interface {
	Match(s string) (int, *Token)
}

type Lexer struct {
	Rules []TokenMatcher
	Pos   Position
	Data  string
}

func (lex *Lexer) Next() Token {
	for continueScanning := true; continueScanning; {
		// default to single pass, will be set to true if whitespace is found
		continueScanning = false

		if len(lex.Data) == 0 {
			// scanned until end of file
			return Token{Name: "eof", Value: "", Line: lex.Pos.Line}
		}

		for _, rule := range lex.Rules {
			consumed, tok := rule.Match(lex.Data)
			if consumed > 0 {
				// register data was consumed
				lex.Data = lex.Data[consumed:len(lex.Data)]
				lex.Pos.Index += consumed
				if consumed >= len(lex.Data) {
					lex.Pos.Line += strings.Count(lex.Data, "\n")
				} else {
					lex.Pos.Line += strings.Count(lex.Data[0:consumed], "\n")
				}
			}

			if tok != nil {
				// first matching rule wins
				tok.Line = lex.Pos.Line
				return *tok
			}
			if consumed > 0 && tok == nil {
				// whitespace: consumed data but got no token
				continueScanning = true
				break
			}
		}
	}
	// no rule matched, token not recognized
	return Token{Name: "error", Line: lex.Pos.Line}
}

func NewWithDefaultRules(code string) *Lexer {
	var lex = new(Lexer)
	lex.Data = code
	lex.Rules = []TokenMatcher{
		NewWhitespaceRule(`\s+|;`),
		NewRegexRule("keyword", `while|if|else|print`, ExtractFullMatch),
		NewRegexRule("int", `[0-9]+`, ExtractFullMatch),
		NewRegexRule("bool", `true|false`, ExtractFullMatch),
		NewRegexRule("op", `[=><+*!]|:=|\|\||&&|==`, ExtractFullMatch),
		NewRegexRule("par_open", `\(`, ExtractFullMatch),
		NewRegexRule("par_close", `\)`, ExtractFullMatch),
		NewRegexRule("brace_open", `\{`, ExtractFullMatch),
		NewRegexRule("brace_close", `\}`, ExtractFullMatch),
		NewRegexRule("identifier", `[a-z][a-zA-Z]*`, ExtractFullMatch),
	}
	return lex
}

func IsEof(t Token) bool {
	return t.Name == "eof" || t.Name == "error"
}
