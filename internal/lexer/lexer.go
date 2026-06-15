package lexer

import (
	"fmt"
	"strings"
)

type Lexer struct {
	source    []byte
	tokens    []Token
	start     int
	current   int
	line      int
	lineStart int
	file      string // absolute path; attached to every produced token
}

// SetFile sets the source path recorded on each token (typically absolute).
func (l *Lexer) SetFile(path string) {
	l.file = path
}

// NewLexer constructs a lexer for source. file is stored on every token for diagnostics (use "" if unknown).
func NewLexer(source, file string) *Lexer {
	src := []byte(source)
	// Handle BOM
	if len(src) >= 3 && src[0] == 0xEF && src[1] == 0xBB && src[2] == 0xBF {
		src = src[3:]
	}
	return &Lexer{
		source: src,
		line:   1,
		file:   file,
	}
}

func (l *Lexer) Tokenize() ([]Token, error) {
	for !l.isAtEnd() {
		l.start = l.current
		if err := l.scanToken(); err != nil {
			return nil, err
		}
	}

	l.addToken(TokenEOF)
	return l.tokens, nil
}

func (l *Lexer) scanToken() error {
	c := l.advance()
	switch c {
	case '(':
		l.addToken(TokenLParen)
	case ')':
		l.addToken(TokenRParen)
	case '{':
		l.addToken(TokenLBrace)
	case '}':
		l.addToken(TokenRBrace)
	case '[':
		l.addToken(TokenLBracket)
	case ']':
		l.addToken(TokenRBracket)
	case ',':
		l.addToken(TokenComma)
	case '.':
		if l.match('.') {
			if l.match('.') {
				l.addToken(TokenTripleDot)
			} else {
				l.addToken(TokenDotDot)
			}
		} else {
			l.addToken(TokenDot)
		}
	case ';':
		l.addToken(TokenSemicolon)
	case ':':
		l.addToken(TokenColon)
	case '?':
		if l.match('?') {
			if l.match('=') {
				l.addToken(TokenQuestionQuestionEqual)
			} else {
				l.addToken(TokenQuestionQuestion)
			}
		} else if l.match('.') {
			l.addToken(TokenOptionalDot)
		} else {
			l.addToken(TokenQuestion)
		}
	case '`':
		// Opening backtick already consumed above; avoid a second advance (EOF panic on "`").
		line := l.line
		col := l.current - l.lineStart
		l.tokens = append(l.tokens, Token{
			Type:   TokenTemplateStart,
			Lexeme: "`",
			Line:   line,
			Col:    col,
			File:   l.file,
		})
		return l.scanTemplateTail()
	case '-':
		if l.match('-') {
			l.addToken(TokenMinusMinus)
		} else if l.match('=') {
			l.addToken(TokenMinusEqual)
		} else {
			l.addToken(TokenMinus)
		}
	case '+':
		if l.match('+') {
			l.addToken(TokenPlusPlus)
		} else if l.match('=') {
			l.addToken(TokenPlusEqual)
		} else {
			l.addToken(TokenPlus)
		}
	case '*':
		if l.match('=') {
			l.addToken(TokenStarEqual)
		} else {
			l.addToken(TokenStar)
		}
	case '/':
		if l.match('/') {
			start := l.current
			for l.peek() != '\n' && !l.isAtEnd() {
				l.advance()
			}
			comment := string(l.source[start:l.current])
			trim := strings.TrimLeft(comment, " \t")
			if len(trim) >= 5 && strings.EqualFold(trim[:5], "koda:") {
				payload := strings.TrimSpace(trim[5:])
				l.addTokenWithLexeme(TokenComment, "koda:"+payload)
			}
		} else if l.match('*') {
			for !l.isAtEnd() {
				if l.peek() == '*' && l.peekNext() == '/' {
					l.advance()
					l.advance()
					break
				}
				if l.peek() == '\n' {
					l.line++
					l.lineStart = l.current + 1
				}
				l.advance()
			}
		} else if l.match('=') {
			l.addToken(TokenSlashEqual)
		} else {
			l.addToken(TokenSlash)
		}
	case '%':
		if l.match('=') {
			l.addToken(TokenPercentEqual)
		} else {
			l.addToken(TokenPercent)
		}
	case '&':
		if l.match('&') {
			l.addToken(TokenAndAnd)
		} else if l.match('=') {
			l.addToken(TokenAndEqual)
		} else {
			l.addToken(TokenAnd)
		}
	case '|':
		if l.match('|') {
			l.addToken(TokenOrOr)
		} else if l.match('=') {
			l.addToken(TokenOrEqual)
		} else {
			l.addToken(TokenOr)
		}
	case '^':
		if l.match('=') {
			l.addToken(TokenCaretEqual)
		} else {
			l.addToken(TokenCaret)
		}
	case '~':
		l.addToken(TokenTilde)
	case '!':
		if l.match('=') {
			if l.match('=') {
				l.addToken(TokenStrictNotEqual)
			} else {
				l.addToken(TokenBangEqual)
			}
		} else {
			l.addToken(TokenBang)
		}
	case '=':
		if l.match('=') {
			if l.match('=') {
				l.addToken(TokenStrictEqual)
			} else {
				l.addToken(TokenEqualEqual)
			}
		} else if l.match('>') {
			l.addToken(TokenArrow)
		} else {
			l.addToken(TokenEqual)
		}
	case '<':
		if l.match('=') {
			l.addToken(TokenLessEqual)
		} else if l.match('<') {
			if l.match('=') {
				l.addToken(TokenLessLessEqual)
			} else {
				l.addToken(TokenLessLess)
			}
		} else {
			l.addToken(TokenLess)
		}
	case '>':
		if l.match('=') {
			l.addToken(TokenGreaterEqual)
		} else if l.match('>') {
			if l.match('=') {
				l.addToken(TokenGreaterGreaterEqual)
			} else if l.match('>') {
				l.addToken(TokenUnsignedShift)
			} else {
				l.addToken(TokenGreaterGreater)
			}
		} else {
			l.addToken(TokenGreater)
		}
	case '"':
		return l.string()
	case ' ', '\r', '\t':
		// Ignore whitespace
	case '\n':
		l.line++
		l.lineStart = l.current
	case '#':
		return l.directive()
	default:
		if isDigit(c) {
			return l.number()
		} else if isAlpha(c) {
			return l.identifier()
		} else {
			return fmt.Errorf("unexpected character at %d:%d: %c", l.line, l.current-l.lineStart, c)
		}
	}
	return nil
}

func (l *Lexer) string() error {
	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\n' {
			l.line++
			l.lineStart = l.current + 1
		}
		if l.peek() == '\\' {
			l.advance()
			if l.isAtEnd() {
				return fmt.Errorf("unterminated string at %d", l.line)
			}
			l.advance()
			continue
		}
		l.advance()
	}

	if l.isAtEnd() {
		return fmt.Errorf("unterminated string at %d", l.line)
	}

	l.advance() // The closing "
	l.addToken(TokenString)
	return nil
}

func (l *Lexer) number() error {
	// Leading digit was consumed before this call; l.source[l.start] is the first digit.
	if l.source[l.start] == '0' && (l.peek() == 'x' || l.peek() == 'X') {
		l.advance() // x / X
		startHex := l.current
		for isHexDigit(l.peek()) {
			l.advance()
		}
		if l.current == startHex {
			return fmt.Errorf("invalid hex literal at %d:%d", l.line, l.start-l.lineStart+1)
		}
		l.addToken(TokenNumber)
		return nil
	}
	if l.source[l.start] == '0' && (l.peek() == 'b' || l.peek() == 'B') {
		l.advance()
		startBin := l.current
		for l.peek() == '0' || l.peek() == '1' {
			l.advance()
		}
		if l.current == startBin {
			return fmt.Errorf("invalid binary literal at %d:%d", l.line, l.start-l.lineStart+1)
		}
		l.addToken(TokenNumber)
		return nil
	}

	for isDigit(l.peek()) {
		l.advance()
	}

	if l.peek() == '.' && isDigit(l.peekNext()) {
		l.advance() // Consume the "."
		for isDigit(l.peek()) {
			l.advance()
		}
	}

	if l.peek() == 'e' || l.peek() == 'E' {
		l.advance()
		if l.peek() == '+' || l.peek() == '-' {
			l.advance()
		}
		if !isDigit(l.peek()) {
			return fmt.Errorf("invalid exponent in number at %d:%d", l.line, l.start-l.lineStart+1)
		}
		for isDigit(l.peek()) {
			l.advance()
		}
	}

	l.addToken(TokenNumber)
	return nil
}

func (l *Lexer) identifier() error {
	for isAlphaNumeric(l.peek()) {
		l.advance()
	}

	raw := string(l.source[l.start:l.current])
	typ := l.lookupKeyword(strings.ToLower(raw))
	if typ == TokenIdentifier {
		l.addToken(typ)
	} else {
		l.addTokenWithLexeme(typ, strings.ToLower(raw))
	}
	return nil
}

func (l *Lexer) directive() error {
	for isAlpha(l.peek()) {
		l.advance()
	}
	text := string(l.source[l.start:l.current])
	if strings.EqualFold(text, "#include") {
		l.addTokenWithLexeme(TokenInclude, strings.ToLower(text))
	} else {
		l.addToken(TokenError)
	}
	return nil
}

func (l *Lexer) lookupKeyword(text string) TokenType {
	switch text {
	case "break":
		return TokenBreak
	case "case":
		return TokenCase
	case "continue":
		return TokenContinue
	case "default":
		return TokenDefault
	case "defer":
		return TokenDefer
	case "delete":
		return TokenDelete
	case "do":
		return TokenDo
	case "else":
		return TokenElse
	case "false":
		return TokenFalse
	case "for":
		return TokenFor
	case "func":
		return TokenFunc
	case "if":
		return TokenIf
	case "import":
		return TokenImport
	case "in":
		return TokenIn
	case "let":
		return TokenLet
	case "const":
		return TokenConst
	case "var":
		return TokenVar
	case "null":
		return TokenNull
	case "of":
		return TokenOf
	case "return":
		return TokenReturn
	case "switch":
		return TokenSwitch
	case "struct":
		return TokenStruct
	case "test":
		return TokenTest
	case "enum":
		return TokenEnum
	case "this":
		return TokenThis
	case "true":
		return TokenTrue
	case "while":
		return TokenWhile
	case "typeof":
		return TokenTypeof
	default:
		return TokenIdentifier
	}
}

func unescapeTemplateRun(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case '`', '$', '\\':
				b.WriteByte(s[i+1])
				i++
				continue
			}
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

// scanTemplateTail runs after the opening ` token has been consumed and TemplateStart appended.
func (l *Lexer) scanTemplateTail() error {
	if l.peek() == '`' {
		l.advance()
		l.tokens = append(l.tokens, Token{
			Type:   TokenTemplateClose,
			Lexeme: "`",
			Line:   l.line,
			Col:    l.current - l.lineStart,
			File:   l.file,
		})
		return nil
	}

	for {
		if l.isAtEnd() {
			return fmt.Errorf("unterminated template literal at %d:%d", l.line, l.current-l.lineStart+1)
		}
		if l.peek() == '`' {
			l.advance()
			l.tokens = append(l.tokens, Token{
				Type:   TokenTemplateClose,
				Lexeme: "`",
				Line:   l.line,
				Col:    l.current - l.lineStart,
				File:   l.file,
			})
			return nil
		}
		if l.peek() == '$' && l.peekNext() == '{' {
			l.advance()
			l.advance()
			exprStart := l.current
			depth := 1
			for depth > 0 && !l.isAtEnd() {
				c := l.peek()
				if c == '{' {
					depth++
					l.advance()
					continue
				}
				if c == '}' {
					depth--
					l.advance()
					continue
				}
				if c == '\n' {
					l.line++
					l.lineStart = l.current + 1
				}
				l.advance()
			}
			if depth != 0 {
				return fmt.Errorf("unterminated '${}' in template literal at %d:%d", l.line, l.current-l.lineStart+1)
			}
			body := string(l.source[exprStart : l.current-1])
			l.tokens = append(l.tokens, Token{
				Type:   TokenTemplateInterp,
				Lexeme: body,
				Line:   l.line,
				Col:    exprStart - l.lineStart + 1,
				File:   l.file,
			})
			continue
		}

		textStart := l.current
		for l.peek() != '`' && !(l.peek() == '$' && l.peekNext() == '{') && !l.isAtEnd() {
			if l.peek() == '\\' {
				l.advance()
				if !l.isAtEnd() {
					l.advance()
				}
				continue
			}
			if l.peek() == '\n' {
				l.line++
				l.lineStart = l.current + 1
			}
			l.advance()
		}
		raw := string(l.source[textStart:l.current])
		if raw != "" {
			l.tokens = append(l.tokens, Token{
				Type:   TokenTemplateString,
				Lexeme: unescapeTemplateRun(raw),
				Line:   l.line,
				Col:    textStart - l.lineStart + 1,
				File:   l.file,
			})
		}
	}
}

func (l *Lexer) advance() byte {
	c := l.source[l.current]
	l.current++
	return c
}

func (l *Lexer) match(expected byte) bool {
	if l.isAtEnd() {
		return false
	}
	if l.source[l.current] != expected {
		return false
	}
	l.current++
	return true
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}
	return l.source[l.current]
}

func (l *Lexer) peekNext() byte {
	if l.current+1 >= len(l.source) {
		return 0
	}
	return l.source[l.current+1]
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) addToken(typ TokenType) {
	text := string(l.source[l.start:l.current])
	l.addTokenWithLexeme(typ, text)
}

func (l *Lexer) addTokenWithLexeme(typ TokenType, lexeme string) {
	l.tokens = append(l.tokens, Token{
		Type:   typ,
		Lexeme: lexeme,
		Line:   l.line,
		Col:    l.start - l.lineStart + 1,
		File:   l.file,
	})
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isHexDigit(c byte) bool {
	return isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}
