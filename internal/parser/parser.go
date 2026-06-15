package parser

import (
	"fmt"
	"koda/internal/lexer"
)

type Parser struct {
	tokens        []lexer.Token
	current       int
	lastDirective *NativeDirective
	destructTmp   int
	testSeq       int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() (*Program, error) {
	program := &Program{Declarations: []Decl{}}
	for !p.isAtEnd() {
		decls, err := p.parseDeclaration()
		if err != nil {
			return nil, err
		}
		program.Declarations = append(program.Declarations, decls...)
	}
	return program, nil
}

// Utility methods

func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(typ lexer.TokenType, message string) (lexer.Token, error) {
	if p.check(typ) {
		return p.advance(), nil
	}
	return lexer.Token{}, p.error(p.peek(), message)
}

func (p *Parser) check(typ lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == typ
}

func (p *Parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == lexer.TokenEOF
}

func (p *Parser) peek() lexer.Token {
	return p.tokens[p.current]
}

func (p *Parser) peekNext() lexer.Token {
	if p.current+1 >= len(p.tokens) {
		return lexer.Token{Type: lexer.TokenEOF}
	}
	return p.tokens[p.current+1]
}

func (p *Parser) previous() lexer.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) error(token lexer.Token, message string) error {
	return fmt.Errorf("[line %d:%d] error at '%s': %s", token.Line, token.Col, token.Lexeme, message)
}
