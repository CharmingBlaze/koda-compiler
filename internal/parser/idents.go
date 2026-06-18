package parser

import (
	"strings"

	"koda/internal/lexer"
)

// normalizeIdentLexeme folds an identifier to ASCII lowercase for names,
// property keys, and lookups (Koda is case-insensitive for identifiers).
func normalizeIdentLexeme(t lexer.Token) lexer.Token {
	t.Lexeme = strings.ToLower(t.Lexeme)
	return t
}

// consumeParamName reads a function parameter name. Receiver aliases self/this are allowed.
func (p *Parser) consumeParamName() (lexer.Token, error) {
	switch p.peek().Type {
	case lexer.TokenIdentifier, lexer.TokenSelf, lexer.TokenThis:
		return p.advance(), nil
	default:
		return p.consume(lexer.TokenIdentifier, "expected parameter name")
	}
}
