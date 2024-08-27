package lexer

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type TokenType int64

const (
	// Literal Type
	Number TokenType = iota
	String
	Identifier
	// Keywords
	Let
	Const
	Fn
	If
	Else
	// Grouping * Operators
	BinaryOperator      // + - * / %
	Equals              // =
	EqualsEquals        // ==
	NotEquals           // !=
	LessThanOrEquals    // <=
	GreaterThanOrEquals // >=
	LessThan            // <
	GreaterThan         // >
	Dot                 // .
	Coma                // ,
	Colon               // :
	Semicolon           // ;
	DoubleQuote         // "
	Not                 // !
	OpenParen           // (
	CloseParen          // )
	OpenBrace           // {
	CloseBrace          // }
	OpenBracket         // [
	CloseBracket        // ]
	EOF                 // Signified the end of file
)

type Token struct {
	Value     string
	TokenType TokenType
	Line      uint64
}

var KEYWORDS = map[string]TokenType{"let": Let, "const": Const, "fn": Fn, "if": If, "else": Else}
var currentLine uint64 = 1

func (tokenType TokenType) String() string {

	return []string{"Number", "String", "Identifier", "Let", "Const", "Fn", "If", "Else", "BinaryOperator", "Equals", "EqualsEquals", "NotEquals", "LessThanOrEquals", "GreaterThanOrEquals", "LessThan", "GreaterThan", "Dot", "Coma", "Colon", "Semicolon", "DoubleQuote", "Not", "OpenParen", "CloseParen", "OpenBrace", "CloseBrace", "OpenBracket", "CloseBracket", "EOF"}[tokenType]
}
func newToken(value string, tType TokenType) Token {
	return Token{Value: value, TokenType: tType, Line: currentLine}
}

func isInt(char string) bool {
	if _, err := strconv.ParseInt(char, 10, 64); err == nil {
		return true
	}

	return false
}
func isAlpha(str string) bool {
	var alpha = regexp.MustCompile("^[a-zA-Z_]+$")

	return alpha.MatchString(str)
}
func isAlphaNumeric(str string) bool {
	var alphanumeric = regexp.MustCompile("^[a-zA-Z_0-9]+$")

	return alphanumeric.MatchString(str)
}
func isSkippable(str string) bool {
	if str == "\n" {
		currentLine += 1
		return true
	}

	return str == " " || str == "\t" || str == "\r"
}

func Tokenize(sourceCode string) []Token {
	var tokens []Token

	src := strings.Split(sourceCode, "")
	for i := 0; i < len(src); i++ {

		if src[i] == "(" {
			tokens = append(tokens, newToken(src[i], OpenParen))
		} else if src[i] == ")" {
			tokens = append(tokens, newToken(src[i], CloseParen))
		} else if src[i] == "{" {
			tokens = append(tokens, newToken(src[i], OpenBrace))
		} else if src[i] == "}" {
			tokens = append(tokens, newToken(src[i], CloseBrace))
		} else if src[i] == "[" {
			tokens = append(tokens, newToken(src[i], OpenBracket))
		} else if src[i] == "]" {
			tokens = append(tokens, newToken(src[i], CloseBracket))
		} else if src[i] == "+" || src[i] == "-" || src[i] == "*" || src[i] == "/" || src[i] == "%" {
			tokens = append(tokens, newToken(src[i], BinaryOperator))
		} else if src[i] == "=" {
			if i+1 < len(src) && src[i+1] == "=" {
				tokens = append(tokens, newToken("==", EqualsEquals))
				i++
			} else {
				tokens = append(tokens, newToken(src[i], Equals))
			}
		} else if src[i] == "!" {
			if i+1 < len(src) && src[i+1] == "=" {
				tokens = append(tokens, newToken("!=", NotEquals))
				i++
			} else {
				tokens = append(tokens, newToken(src[i], Not))
			}

		} else if src[i] == "<" {
			if i+1 < len(src) && src[i+1] == "=" {
				tokens = append(tokens, newToken("<=", LessThanOrEquals))
				i++
			} else {
				tokens = append(tokens, newToken(src[i], LessThan))
			}
		} else if src[i] == ">" {
			if i+1 < len(src) && src[i+1] == "=" {
				tokens = append(tokens, newToken(">=", GreaterThanOrEquals))
				i++
			} else {
				tokens = append(tokens, newToken(src[i], GreaterThan))
			}
		} else if src[i] == ":" {
			tokens = append(tokens, newToken(src[i], Colon))
		} else if src[i] == ";" {
			tokens = append(tokens, newToken(src[i], Semicolon))
		} else if src[i] == "," {
			tokens = append(tokens, newToken(src[i], Coma))
		} else if src[i] == "." {
			tokens = append(tokens, newToken(src[i], Dot))
		} else if src[i] == `"` {
			i++
			str := ""
			for i < len(src) && src[i] != `"` {
				str += src[i]
				i++
			}

			tokens = append(tokens, newToken(str, String))
		} else {
			if isInt(src[i]) {
				var num string
				for i < len(src) && isInt(src[i]) {
					num += src[i]
					i++
				}
				tokens = append(tokens, newToken(num, Number))
				i--
				continue
			} else if isAlpha(src[i]) {
				indet := src[i]
				i++
				for i < len(src) && isAlphaNumeric(src[i]) {
					indet += src[i]
					i++
				}
				i--
				elem, ok := KEYWORDS[indet]
				if ok {
					tokens = append(tokens, newToken(indet, elem))
				} else {
					tokens = append(tokens, newToken(indet, Identifier))
				}
				continue
			} else if isSkippable(src[i]) {
				continue
			} else {
				fmt.Printf("syntaxError: unreconized character found in source: %v\n", src[i])
				os.Exit(1)
			}
		}
	}
	tokens = append(tokens, Token{Value: "EOF", TokenType: EOF})
	return tokens
}
