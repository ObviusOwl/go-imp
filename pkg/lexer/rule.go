package lexer

import "regexp"

func ExtractFullMatch(m []string) string { return m[0] }

type RegexRule struct {
	tokenName string
	matchReg  *regexp.Regexp
	generator func(match []string) string
}

func NewRegexRule(name string, reg string, generator func(match []string) string) RegexRule {
	return RegexRule{
		tokenName: name,
		matchReg:  regexp.MustCompile("^(?:" + reg + ")"),
		generator: generator,
	}
}

func (rule RegexRule) Match(s string) (int, *Token) {
	if m := rule.matchReg.FindStringSubmatch(s); m != nil {
		return len(m[0]), &Token{Name: rule.tokenName, Value: rule.generator(m)}
	}
	return 0, nil
}

type WhitespaceRule struct {
	matchReg *regexp.Regexp
}

func (rule WhitespaceRule) Match(s string) (int, *Token) {
	return len(rule.matchReg.FindString(s)), nil
}

func NewWhitespaceRule(reg string) WhitespaceRule {
	return WhitespaceRule{matchReg: regexp.MustCompile("^(?:" + reg + ")")}
}
