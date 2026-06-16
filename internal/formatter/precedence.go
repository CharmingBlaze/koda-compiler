package formatter

import (
	"koda/internal/lexer"
)

const (
	_ int = iota
	precLowest
	precAssign
	precOr
	precAnd
	precBitOr
	precBitXor
	precBitAnd
	precEquals
	precLessGreater
	precSum
	precShift
	precProduct
	precPrefix
	precCall
	precIndex
)

var tokenPrec = map[lexer.TokenType]int{
	lexer.TokenEqual:                 precAssign,
	lexer.TokenQuestionQuestionEqual: precAssign,
	lexer.TokenPlusEqual:             precAssign,
	lexer.TokenMinusEqual:            precAssign,
	lexer.TokenStarEqual:             precAssign,
	lexer.TokenSlashEqual:            precAssign,
	lexer.TokenPercentEqual:          precAssign,
	lexer.TokenAndEqual:              precAssign,
	lexer.TokenOrEqual:               precAssign,
	lexer.TokenCaretEqual:            precAssign,
	lexer.TokenLessLessEqual:         precAssign,
	lexer.TokenGreaterGreaterEqual:   precAssign,
	lexer.TokenEqualEqual:  precEquals,
	lexer.TokenBangEqual:   precEquals,
	lexer.TokenLess:                  precLessGreater,
	lexer.TokenLessEqual:             precLessGreater,
	lexer.TokenGreater:               precLessGreater,
	lexer.TokenGreaterEqual:          precLessGreater,
	lexer.TokenPlus:                  precSum,
	lexer.TokenMinus:                 precSum,
	lexer.TokenLessLess:              precShift,
	lexer.TokenGreaterGreater:        precShift,
	lexer.TokenUnsignedShift:         precShift,
	lexer.TokenSlash:                 precProduct,
	lexer.TokenStar:                  precProduct,
	lexer.TokenPercent:               precProduct,
	lexer.TokenAndAnd:                precAnd,
	lexer.TokenOrOr:                  precOr,
	lexer.TokenAnd:                   precBitAnd,
	lexer.TokenOr:                    precBitOr,
	lexer.TokenCaret:                 precBitXor,
	lexer.TokenLParen:                precCall,
	lexer.TokenLBracket:              precIndex,
	lexer.TokenDot:                   precIndex,
}

func precOf(typ lexer.TokenType) int {
	if p, ok := tokenPrec[typ]; ok {
		return p
	}
	return precLowest
}
