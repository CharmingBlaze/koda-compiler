package parser

import (
	"fmt"
	"strings"

	"koda/internal/lexer"
)

func (p *Parser) parseDeclaration() ([]Decl, error) {
	if p.match(lexer.TokenComment) {
		comment := p.previous().Lexeme
		if strings.HasPrefix(strings.ToLower(comment), "koda:") {
			body := strings.TrimSpace(comment[len("koda:"):])
			parts := strings.Fields(body)
			if len(parts) >= 3 && strings.EqualFold(parts[0], "extern") {
				arity := 0
				if len(parts) >= 4 {
					_, _ = fmt.Sscanf(parts[3], "%d", &arity)
				}
				p.lastDirective = &NativeDirective{
					BindingName: strings.ToLower(parts[1]),
					Symbol:      parts[2],
					Arity:       arity,
				}
			}
		}
		return p.parseDeclaration()
	}
	if p.match(lexer.TokenInclude) {
		d, err := p.parseIncludeDeclaration()
		if err != nil {
			return nil, err
		}
		return []Decl{d}, nil
	}
	if p.match(lexer.TokenVar) {
		tok := p.previous()
		return nil, fmt.Errorf("%d:%d: 'var' is reserved; use 'let' to declare a variable", tok.Line, tok.Col)
	}
	if p.match(lexer.TokenLet) {
		decls, err := p.parseLetDeclarations(false)
		if err != nil {
			return nil, err
		}
		if len(decls) > 0 {
			if let, ok := decls[0].(*LetDecl); ok {
				let.Native = p.lastDirective
			}
		}
		p.lastDirective = nil
		return decls, nil
	}
	if p.match(lexer.TokenConst) {
		decls, err := p.parseLetDeclarations(true)
		if err != nil {
			return nil, err
		}
		p.lastDirective = nil
		return decls, nil
	}
	if p.match(lexer.TokenFunc) {
		decl, err := p.parseFuncDeclaration()
		if err != nil {
			return nil, err
		}
		if f, ok := decl.(*FuncDecl); ok {
			f.Native = p.lastDirective
		}
		p.lastDirective = nil
		return []Decl{decl}, nil
	}
	if p.match(lexer.TokenStruct) {
		decl, err := p.parseStructDeclaration()
		if err != nil {
			return nil, err
		}
		p.lastDirective = nil
		return []Decl{decl}, nil
	}
	if p.match(lexer.TokenEnum) {
		decl, err := p.parseEnumDeclaration()
		if err != nil {
			return nil, err
		}
		p.lastDirective = nil
		return []Decl{decl}, nil
	}
	p.lastDirective = nil
	s, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	return []Decl{s}, nil
}

func (p *Parser) parseIncludeDeclaration() (Decl, error) {
	token := p.previous()
	path, err := p.consume(lexer.TokenString, "expected include path")
	if err != nil {
		return nil, err
	}
	return &IncludeDecl{Token: token, Path: path}, nil
}

func (p *Parser) parseLetDeclarations(isConst bool) ([]Decl, error) {
	token := p.previous()

	if p.check(lexer.TokenLBrace) {
		p.advance()
		var keys []lexer.Token
		if !p.check(lexer.TokenRBrace) {
			for {
				keyTok, err := p.consume(lexer.TokenIdentifier, "expected property name in binding pattern")
				if err != nil {
					return nil, err
				}
				keyTok = normalizeIdentLexeme(keyTok)
				keys = append(keys, keyTok)
				if !p.match(lexer.TokenComma) {
					break
				}
			}
		}
		if _, err := p.consume(lexer.TokenRBrace, "expected '}' after binding pattern"); err != nil {
			return nil, err
		}
		if _, err := p.consume(lexer.TokenEqual, "expected '=' after destructuring pattern"); err != nil {
			return nil, err
		}
		init, err := p.parseExpression(PrecedenceLowest)
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(lexer.TokenSemicolon, "expected ';' after variable declaration"); err != nil {
			return nil, err
		}
		tmpLexeme := fmt.Sprintf("__koda_destruct_%d", p.destructTmp)
		p.destructTmp++
		tmpTok := lexer.Token{Type: lexer.TokenIdentifier, Lexeme: tmpLexeme, Line: token.Line, Col: token.Col, File: token.File}

		out := []Decl{
			&LetDecl{Token: token, Name: tmpTok, IsConst: isConst, Init: init},
		}
		for _, kt := range keys {
			idTok := kt
			objIdent := &IdentifierExpr{Token: tmpTok, Name: tmpTok}
			idxLit := &LiteralExpr{Token: kt, Value: kt.Lexeme}
			initIx := &IndexExpr{Token: kt, Object: objIdent, Index: idxLit}
			out = append(out, &LetDecl{Token: token, Name: idTok, IsConst: isConst, Init: initIx})
		}
		return out, nil
	}

	name, err := p.consume(lexer.TokenIdentifier, "expected variable name")
	if err != nil {
		if p.check(lexer.TokenVar) {
			tok := p.peek()
			p.advance()
			return nil, fmt.Errorf("%d:%d: 'var' is reserved; use 'let' to declare a variable", tok.Line, tok.Col)
		}
		return nil, err
	}
	name = normalizeIdentLexeme(name)

	var typeAnnot string
	if p.match(lexer.TokenColon) {
		typeTok, err := p.consume(lexer.TokenIdentifier, "expected type name after ':'")
		if err != nil {
			return nil, err
		}
		typeAnnot = typeTok.Lexeme
	}

	var init Expr
	if p.match(lexer.TokenEqual) {
		init, err = p.parseExpression(PrecedenceLowest)
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(lexer.TokenSemicolon, "expected ';' after variable declaration"); err != nil {
		return nil, err
	}

	return []Decl{&LetDecl{Token: token, Name: name, TypeAnnot: typeAnnot, IsConst: isConst, Init: init}}, nil
}

func (p *Parser) parseFuncDeclaration() (Decl, error) {
	token := p.previous()
	name, err := p.consume(lexer.TokenIdentifier, "expected function name")
	if err != nil {
		if p.check(lexer.TokenVar) {
			tok := p.peek()
			p.advance()
			return nil, fmt.Errorf("%d:%d: 'var' is reserved; use 'let' to declare a variable", tok.Line, tok.Col)
		}
		return nil, err
	}
	name = normalizeIdentLexeme(name)

	if _, err := p.consume(lexer.TokenLParen, "expected '(' after function name"); err != nil {
		return nil, err
	}

	params := []Param{}
	if !p.check(lexer.TokenRParen) {
		for {
			isRest := p.match(lexer.TokenTripleDot)
			paramName, err := p.consume(lexer.TokenIdentifier, "expected parameter name")
			if err != nil {
				if p.check(lexer.TokenVar) {
					tok := p.peek()
					p.advance()
					return nil, fmt.Errorf("%d:%d: 'var' is reserved; use 'let' to declare a variable", tok.Line, tok.Col)
				}
				return nil, err
			}
			paramName = normalizeIdentLexeme(paramName)
			for _, q := range params {
				if q.Name == paramName.Lexeme {
					return nil, p.error(paramName, "duplicate parameter name")
				}
			}
			param := Param{Name: paramName.Lexeme, IsRest: isRest}
			if p.match(lexer.TokenEqual) {
				if isRest {
					return nil, p.error(paramName, "rest parameter cannot have a default")
				}
				param.Default, err = p.parseExpression(PrecedenceLowest)
				if err != nil {
					return nil, err
				}
			}
			params = append(params, param)
			if isRest && !p.check(lexer.TokenRParen) {
				return nil, p.error(paramName, "rest parameter must be last")
			}

			if !p.match(lexer.TokenComma) {
				break
			}
		}
	}

	if _, err := p.consume(lexer.TokenRParen, "expected ')' after parameters"); err != nil {
		return nil, err
	}

	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return &FuncDecl{Token: token, Name: name, Params: params, Body: body}, nil
}

func (p *Parser) parseStatement() (Stmt, error) {
	if p.match(lexer.TokenDelete) {
		return p.parseDeleteStatement()
	}
	if p.match(lexer.TokenIf) {
		return p.parseIfStatement()
	}
	if p.match(lexer.TokenWhile) {
		return p.parseWhileStatement()
	}
	if p.match(lexer.TokenDo) {
		return p.parseDoWhileStatement()
	}
	if p.match(lexer.TokenFor) {
		return p.parseForStatement()
	}
	if p.match(lexer.TokenReturn) {
		return p.parseReturnStatement()
	}
	if p.match(lexer.TokenDefer) {
		return p.parseDeferStatement()
	}
	if p.match(lexer.TokenBreak) {
		token := p.previous()
		if _, err := p.consume(lexer.TokenSemicolon, "expected ';' after break"); err != nil {
			return nil, err
		}
		return &BreakStmt{Token: token}, nil
	}
	if p.match(lexer.TokenContinue) {
		token := p.previous()
		if _, err := p.consume(lexer.TokenSemicolon, "expected ';' after continue"); err != nil {
			return nil, err
		}
		return &ContinueStmt{Token: token}, nil
	}
	if p.match(lexer.TokenSwitch) {
		return p.parseSwitchStatement()
	}
	if p.check(lexer.TokenLBrace) {
		return p.parseBlockStatement()
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseBlockStatement() (*BlockStmt, error) {
	token, err := p.consume(lexer.TokenLBrace, "expected '{' at start of block")
	if err != nil {
		return nil, err
	}
	declarations := []Decl{}

	for !p.check(lexer.TokenRBrace) && !p.isAtEnd() {
		decls, err := p.parseDeclaration()
		if err != nil {
			return nil, err
		}
		for _, decl := range decls {
			if _, ok := decl.(Stmt); !ok {
				return nil, p.error(p.previous(), "expected statement in block")
			}
			declarations = append(declarations, decl)
		}
	}

	if _, err := p.consume(lexer.TokenRBrace, "expected '}' after block"); err != nil {
		return nil, err
	}

	return &BlockStmt{Token: token, Declarations: declarations}, nil
}

func (p *Parser) parseIfStatement() (Stmt, error) {
	token := p.previous()
	if _, err := p.consume(lexer.TokenLParen, "expected '(' after 'if'"); err != nil {
		return nil, err
	}

	condition, err := p.parseExpression(PrecedenceLowest)
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(lexer.TokenRParen, "expected ')' after condition"); err != nil {
		return nil, err
	}

	thenBranch, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	var elseBranch Stmt
	if p.match(lexer.TokenElse) {
		elseBranch, err = p.parseStatement()
		if err != nil {
			return nil, err
		}
	}

	return &IfStmt{Token: token, Condition: condition, Then: thenBranch, Else: elseBranch}, nil
}

func (p *Parser) parseWhileStatement() (Stmt, error) {
	token := p.previous()
	if _, err := p.consume(lexer.TokenLParen, "expected '(' after 'while'"); err != nil {
		return nil, err
	}

	condition, err := p.parseExpression(PrecedenceLowest)
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(lexer.TokenRParen, "expected ')' after condition"); err != nil {
		return nil, err
	}

	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	return &WhileStmt{Token: token, Condition: condition, Body: body}, nil
}

func (p *Parser) parseDoWhileStatement() (Stmt, error) {
	token := p.previous()
	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(lexer.TokenWhile, "expected 'while' after do body"); err != nil {
		return nil, err
	}
	if _, err := p.consume(lexer.TokenLParen, "expected '(' after 'while'"); err != nil {
		return nil, err
	}
	condition, err := p.parseExpression(PrecedenceLowest)
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(lexer.TokenRParen, "expected ')' after do-while condition"); err != nil {
		return nil, err
	}
	if _, err := p.consume(lexer.TokenSemicolon, "expected ';' after do-while"); err != nil {
		return nil, err
	}
	return &DoWhileStmt{Token: token, Body: body, Condition: condition}, nil
}

func (p *Parser) parseSwitchStatement() (Stmt, error) {
	token := p.previous()
	if _, err := p.consume(lexer.TokenLParen, "expected '(' after 'switch'"); err != nil {
		return nil, err
	}
	subject, err := p.parseExpression(PrecedenceLowest)
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(lexer.TokenRParen, "expected ')' after switch subject"); err != nil {
		return nil, err
	}
	if _, err := p.consume(lexer.TokenLBrace, "expected '{' before switch body"); err != nil {
		return nil, err
	}

	var cases []SwitchCase
	var def []Decl

	for !p.check(lexer.TokenRBrace) && !p.isAtEnd() {
		if p.match(lexer.TokenCase) {
			val, err := p.parseExpression(PrecedenceLowest)
			if err != nil {
				return nil, err
			}
			if _, err := p.consume(lexer.TokenColon, "expected ':' after case value"); err != nil {
				return nil, err
			}
			var body []Decl
			for !p.check(lexer.TokenCase) && !p.check(lexer.TokenDefault) && !p.check(lexer.TokenRBrace) && !p.isAtEnd() {
				ds, err := p.parseDeclaration()
				if err != nil {
					return nil, err
				}
				body = append(body, ds...)
			}
			cases = append(cases, SwitchCase{Value: val, Body: body})
			continue
		}
		if p.match(lexer.TokenDefault) {
			if _, err := p.consume(lexer.TokenColon, "expected ':' after default"); err != nil {
				return nil, err
			}
			for !p.check(lexer.TokenRBrace) && !p.isAtEnd() {
				ds, err := p.parseDeclaration()
				if err != nil {
					return nil, err
				}
				def = append(def, ds...)
			}
			continue
		}
		return nil, p.error(p.peek(), "expected case or default in switch")
	}

	if _, err := p.consume(lexer.TokenRBrace, "expected '}' after switch"); err != nil {
		return nil, err
	}
	return &SwitchStmt{Token: token, Subject: subject, Cases: cases, Default: def}, nil
}

func (p *Parser) parseForStatement() (Stmt, error) {
	token := p.previous()
	if _, err := p.consume(lexer.TokenLParen, "expected '(' after 'for'"); err != nil {
		return nil, err
	}
	// for-in / for-of: for (let name in/of …) or for (let [k, v] of …)
	if p.match(lexer.TokenLet) {
		if p.match(lexer.TokenLBracket) {
			keyTok, err := p.consume(lexer.TokenIdentifier, "expected key name in [k, v]")
			if err != nil {
				return nil, err
			}
			keyTok = normalizeIdentLexeme(keyTok)
			if _, err := p.consume(lexer.TokenComma, "expected ',' between key and value in [k, v]"); err != nil {
				return nil, err
			}
			valTok, err := p.consume(lexer.TokenIdentifier, "expected value name in [k, v]")
			if err != nil {
				return nil, err
			}
			valTok = normalizeIdentLexeme(valTok)
			if _, err := p.consume(lexer.TokenRBracket, "expected ']' after [k, v]"); err != nil {
				return nil, err
			}
			if !p.match(lexer.TokenOf) {
				return nil, p.error(p.peek(), "expected 'of' after [k, v]")
			}
			iter, err := p.parseExpression(PrecedenceLowest)
			if err != nil {
				return nil, err
			}
			if _, err := p.consume(lexer.TokenRParen, "expected ')' after for-of iterable"); err != nil {
				return nil, err
			}
			body, err := p.parseStatement()
			if err != nil {
				return nil, err
			}
			v := valTok
			return &ForOfStmt{Token: token, VarName: keyTok, ValueVar: &v, Iterable: iter, Body: body}, nil
		}
		name, err := p.consume(lexer.TokenIdentifier, "expected loop variable name")
		if err != nil {
			return nil, err
		}
		name = normalizeIdentLexeme(name)
		if p.match(lexer.TokenIn) {
			iter, err := p.parseExpression(PrecedenceLowest)
			if err != nil {
				return nil, err
			}
			if _, err := p.consume(lexer.TokenRParen, "expected ')' after for-in iterable"); err != nil {
				return nil, err
			}
			body, err := p.parseStatement()
			if err != nil {
				return nil, err
			}
			n := name
			return &ForInStmt{Token: token, KeyVar: &n, Iterable: iter, Body: body}, nil
		}
		if p.match(lexer.TokenOf) {
			iter, err := p.parseExpression(PrecedenceLowest)
			if err != nil {
				return nil, err
			}
			if _, err := p.consume(lexer.TokenRParen, "expected ')' after for-of iterable"); err != nil {
				return nil, err
			}
			body, err := p.parseStatement()
			if err != nil {
				return nil, err
			}
			return &ForOfStmt{Token: token, VarName: name, Iterable: iter, Body: body}, nil
		}
		// Classic for: for (let name [= expr]? [, let name [= expr]?]*; cond; incr)
		letTok := p.tokens[p.current-2]
		var init Expr
		if p.match(lexer.TokenEqual) {
			init, err = p.parseExpression(PrecedenceLowest)
			if err != nil {
				return nil, err
			}
		}
		inits := []Decl{&LetDecl{Token: letTok, Name: name, Init: init}}
		for p.match(lexer.TokenComma) {
			if !p.match(lexer.TokenLet) {
				return nil, p.error(p.peek(), "expected 'let' after ',' in for-loop initializer")
			}
			letTok2 := p.previous()
			name2, err := p.consume(lexer.TokenIdentifier, "expected variable name in for-loop initializer")
			if err != nil {
				return nil, err
			}
			name2 = normalizeIdentLexeme(name2)
			var init2 Expr
			if p.match(lexer.TokenEqual) {
				init2, err = p.parseExpression(PrecedenceLowest)
				if err != nil {
					return nil, err
				}
			}
			inits = append(inits, &LetDecl{Token: letTok2, Name: name2, Init: init2})
		}
		if _, err := p.consume(lexer.TokenSemicolon, "expected ';' after for-loop initializer"); err != nil {
			return nil, err
		}
		return p.finishClassicFor(token, inits)
	}
	// Classic for with empty initializer: for (; cond; incr)
	if p.match(lexer.TokenSemicolon) {
		return p.finishClassicFor(token, nil)
	}
	return nil, p.error(p.peek(), "expected 'let' or ';' after '(' in for-loop")
}

func (p *Parser) finishClassicFor(forTok lexer.Token, inits []Decl) (*ForStmt, error) {
	var cond Expr
	var err error
	if !p.check(lexer.TokenSemicolon) {
		cond, err = p.parseExpression(PrecedenceLowest)
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(lexer.TokenSemicolon, "expected ';' after for-loop condition"); err != nil {
		return nil, err
	}

	var increments []Expr
	if !p.check(lexer.TokenRParen) {
		for {
			inc, err := p.parseExpression(PrecedenceLowest)
			if err != nil {
				return nil, err
			}
			increments = append(increments, inc)
			if !p.match(lexer.TokenComma) {
				break
			}
		}
	}
	if _, err := p.consume(lexer.TokenRParen, "expected ')' after for-loop clauses"); err != nil {
		return nil, err
	}
	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	return &ForStmt{Token: forTok, Inits: inits, Condition: cond, Increments: increments, Body: body}, nil
}

func (p *Parser) parseReturnStatement() (Stmt, error) {
	token := p.previous()
	var value Expr
	var err error

	if !p.check(lexer.TokenSemicolon) {
		value, err = p.parseExpression(PrecedenceLowest)
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(lexer.TokenSemicolon, "expected ';' after return value"); err != nil {
		return nil, err
	}

	return &ReturnStmt{Token: token, Value: value}, nil
}

func (p *Parser) parseDeferStatement() (Stmt, error) {
	token := p.previous()
	expr, err := p.parseExpression(PrecedenceLowest)
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(lexer.TokenSemicolon, "expected ';' after defer expression"); err != nil {
		return nil, err
	}
	return &DeferStmt{Token: token, Expr: expr}, nil
}

func (p *Parser) parseDeleteStatement() (Stmt, error) {
	tok := p.previous()
	target, err := p.parseExpression(PrecedenceLowest)
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(lexer.TokenSemicolon, "expected ';' after delete"); err != nil {
		return nil, err
	}
	return &DeleteStmt{Token: tok, Target: target}, nil
}

func (p *Parser) parseExpressionStatement() (Stmt, error) {
	token := p.peek()
	expr, err := p.parseExpression(PrecedenceLowest)
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(lexer.TokenSemicolon, "expected ';' after expression"); err != nil {
		return nil, err
	}

	return &ExpressionStmt{Token: token, Expr: expr}, nil
}

func (p *Parser) parseStructDeclaration() (Decl, error) {
	kwTok := p.previous()
	name, err := p.consume(lexer.TokenIdentifier, "expected struct name")
	if err != nil {
		return nil, err
	}
	name = normalizeIdentLexeme(name)
	if _, err := p.consume(lexer.TokenLBrace, "expected '{' before struct fields"); err != nil {
		return nil, err
	}
	var fields []lexer.Token
	var methods []*FuncDecl
	for !p.check(lexer.TokenRBrace) {
		if p.match(lexer.TokenFunc) {
			fd, err := p.parseFuncDeclaration()
			if err != nil {
				return nil, err
			}
			methods = append(methods, fd.(*FuncDecl))
			_ = p.match(lexer.TokenSemicolon)
			continue
		}
		f, err := p.consume(lexer.TokenIdentifier, "expected field name or func declaration")
		if err != nil {
			return nil, err
		}
		f = normalizeIdentLexeme(f)
		fields = append(fields, f)
		if p.check(lexer.TokenRBrace) {
			break
		}
		if p.check(lexer.TokenFunc) {
			continue
		}
		_ = p.match(lexer.TokenSemicolon)
		if !p.match(lexer.TokenComma) {
			if p.check(lexer.TokenRBrace) {
				break
			}
			if p.check(lexer.TokenFunc) {
				continue
			}
		}
	}
	if _, err := p.consume(lexer.TokenRBrace, "expected '}' after struct fields"); err != nil {
		return nil, err
	}
	_ = p.match(lexer.TokenSemicolon)
	return &StructDecl{Token: kwTok, Name: name, Fields: fields, Methods: methods}, nil
}

func (p *Parser) parseEnumDeclaration() (Decl, error) {
	kwTok := p.previous()
	name, err := p.consume(lexer.TokenIdentifier, "expected enum name")
	if err != nil {
		return nil, err
	}
	name = normalizeIdentLexeme(name)
	if _, err := p.consume(lexer.TokenLBrace, "expected '{' before enum members"); err != nil {
		return nil, err
	}
	var members []lexer.Token
	if !p.check(lexer.TokenRBrace) {
		for {
			m, err := p.consume(lexer.TokenIdentifier, "expected enum member name")
			if err != nil {
				return nil, err
			}
			m = normalizeIdentLexeme(m)
			members = append(members, m)
			if !p.match(lexer.TokenComma) {
				break
			}
		}
	}
	if _, err := p.consume(lexer.TokenRBrace, "expected '}' after enum members"); err != nil {
		return nil, err
	}
	_ = p.match(lexer.TokenSemicolon)
	return &EnumDecl{Token: kwTok, Name: name, Members: members}, nil
}
