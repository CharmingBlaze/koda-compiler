package lexer

import "fmt"

type TokenType uint8

const (
	// Special
	TokenError TokenType = iota
	TokenEOF

	// Single-character tokens
	TokenLParen    // (
	TokenRParen    // )
	TokenLBrace    // {
	TokenRBrace    // }
	TokenLBracket  // [
	TokenRBracket  // ]
	TokenComma     // ,
	TokenDot       // .
	TokenMinus     // -
	TokenPlus      // +
	TokenSemicolon // ;
	TokenSlash     // /
	TokenStar      // *
	TokenPercent   // %
	TokenColon     // :
	TokenQuestion       // ?
	TokenQuestionQuestion // ??
	TokenOptionalDot    // ?.
	TokenCaret     // ^
	TokenTilde     // ~
	TokenAnd       // &
	TokenOr        // |
	TokenBang      // !

	// One or two character tokens
	TokenEqual           // =
	TokenEqualEqual      // ==
	TokenStrictEqual     // ===
	TokenBangEqual       // !=
	TokenStrictNotEqual  // !==
	TokenGreater        // >
	TokenGreaterEqual   // >=
	TokenLess           // <
	TokenLessEqual      // <=
	TokenPlusPlus       // ++
	TokenMinusMinus     // --
	TokenAndAnd         // &&
	TokenOrOr           // ||
	TokenArrow          // =>
	TokenDotDot         // ..
	TokenTripleDot      // ...
	TokenLessLess       // <<
	TokenGreaterGreater // >>
	TokenUnsignedShift  // >>>

	// Assignment operators
	TokenPlusEqual           // +=
	TokenMinusEqual          // -=
	TokenStarEqual           // *=
	TokenSlashEqual          // /=
	TokenPercentEqual        // %=
	TokenAndEqual            // &=
	TokenOrEqual             // |=
	TokenCaretEqual          // ^=
	TokenLessLessEqual       // <<=
	TokenGreaterGreaterEqual // >>=
	TokenQuestionQuestionEqual // ??=

	// Literals
	TokenIdentifier
	TokenString
	TokenNumber

	// Keywords
	TokenBreak
	TokenCase
	TokenContinue
	TokenDefault
	TokenDefer
	TokenDelete
	TokenDo
	TokenElse
	TokenFalse
	TokenFor
	TokenFunc
	TokenIf
	TokenImport
	TokenIn
	TokenLet
	TokenConst
	TokenNull
	TokenOf
	TokenReturn
	TokenSwitch
	TokenStruct
	TokenTest
	TokenEnum
	TokenThis
	TokenTrue
	TokenWhile
	TokenVar // reserved keyword; use 'let' for variable declarations
	TokenTypeof

	// Template literals (` ... ${ expr } ... `)
	TokenTemplateStart
	TokenTemplateString
	TokenTemplateInterp
	TokenTemplateClose

	// Directives
	TokenInclude
	TokenComment
)

func (t TokenType) String() string {
	switch t {
	case TokenComment:
		return "COMMENT"
	case TokenEOF:
		return "EOF"
	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	case TokenLBrace:
		return "{"
	case TokenRBrace:
		return "}"
	case TokenLBracket:
		return "["
	case TokenRBracket:
		return "]"
	case TokenComma:
		return ","
	case TokenDot:
		return "."
	case TokenMinus:
		return "-"
	case TokenPlus:
		return "+"
	case TokenSemicolon:
		return ";"
	case TokenSlash:
		return "/"
	case TokenStar:
		return "*"
	case TokenPercent:
		return "%"
	case TokenColon:
		return ":"
	case TokenQuestion:
		return "?"
	case TokenQuestionQuestion:
		return "??"
	case TokenOptionalDot:
		return "?."
	case TokenCaret:
		return "^"
	case TokenTilde:
		return "~"
	case TokenAnd:
		return "&"
	case TokenOr:
		return "|"
	case TokenBang:
		return "!"
	case TokenEqual:
		return "="
	case TokenEqualEqual:
		return "=="
	case TokenStrictEqual:
		return "==="
	case TokenBangEqual:
		return "!="
	case TokenStrictNotEqual:
		return "!=="
	case TokenGreater:
		return ">"
	case TokenGreaterEqual:
		return ">="
	case TokenLess:
		return "<"
	case TokenLessEqual:
		return "<="
	case TokenPlusPlus:
		return "++"
	case TokenMinusMinus:
		return "--"
	case TokenAndAnd:
		return "&&"
	case TokenOrOr:
		return "||"
	case TokenArrow:
		return "=>"
	case TokenDotDot:
		return ".."
	case TokenTripleDot:
		return "..."
	case TokenLessLess:
		return "<<"
	case TokenGreaterGreater:
		return ">>"
	case TokenUnsignedShift:
		return ">>>"
	case TokenPlusEqual:
		return "+="
	case TokenMinusEqual:
		return "-="
	case TokenStarEqual:
		return "*="
	case TokenSlashEqual:
		return "/="
	case TokenPercentEqual:
		return "%="
	case TokenAndEqual:
		return "&="
	case TokenOrEqual:
		return "|="
	case TokenCaretEqual:
		return "^="
	case TokenLessLessEqual:
		return "<<="
	case TokenGreaterGreaterEqual:
		return ">>="
	case TokenQuestionQuestionEqual:
		return "??="
	case TokenIdentifier:
		return "IDENTIFIER"
	case TokenString:
		return "STRING"
	case TokenNumber:
		return "NUMBER"
	case TokenBreak:
		return "break"
	case TokenCase:
		return "case"
	case TokenContinue:
		return "continue"
	case TokenDefault:
		return "default"
	case TokenDefer:
		return "defer"
	case TokenDelete:
		return "delete"
	case TokenDo:
		return "do"
	case TokenElse:
		return "else"
	case TokenFalse:
		return "false"
	case TokenFor:
		return "for"
	case TokenFunc:
		return "func"
	case TokenIf:
		return "if"
	case TokenImport:
		return "import"
	case TokenIn:
		return "in"
	case TokenLet:
		return "let"
	case TokenConst:
		return "const"
	case TokenNull:
		return "null"
	case TokenOf:
		return "of"
	case TokenReturn:
		return "return"
	case TokenSwitch:
		return "switch"
	case TokenStruct:
		return "struct"
	case TokenTest:
		return "test"
	case TokenEnum:
		return "enum"
	case TokenThis:
		return "this"
	case TokenTrue:
		return "true"
	case TokenWhile:
		return "while"
	case TokenTypeof:
		return "typeof"
	case TokenVar:
		return "var"
	case TokenTemplateStart:
		return "TEMPLATE_START"
	case TokenTemplateString:
		return "TEMPLATE_STRING"
	case TokenTemplateInterp:
		return "TEMPLATE_EXPR"
	case TokenTemplateClose:
		return "TEMPLATE_CLOSE"
	case TokenInclude:
		return "#include"
	default:
		return "UNKNOWN"
	}
}

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
	Col    int
	// File is the absolute path of the source file this token was lexed from, if known.
	File string
}

func (t Token) String() string {
	if t.File != "" {
		return fmt.Sprintf("%s:%d:%d %s '%s'", t.File, t.Line, t.Col, t.Type, t.Lexeme)
	}
	return fmt.Sprintf("%d:%d %s '%s'", t.Line, t.Col, t.Type, t.Lexeme)
}
