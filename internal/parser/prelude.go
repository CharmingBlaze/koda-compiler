package parser

import (
	"strings"

	"koda/internal/lexer"
)

// InjectNativeMathPrelude prepends `let math = { … };` so members like `math.floor` resolve.
// Koda folds identifiers to lowercase (see normalizeIdentLexeme), so the binding is `math`, not `Math`.
func InjectNativeMathPrelude(bundle *ProgramBundle) {
	if bundle == nil || bundle.Entry == nil {
		return
	}
	for _, d := range bundle.Entry.Declarations {
		switch x := d.(type) {
		case *LetDecl:
			if strings.EqualFold(x.Name.Lexeme, "math") {
				return
			}
		case *FuncDecl:
			if strings.EqualFold(x.Name.Lexeme, "math") {
				return
			}
		}
	}

	src := `let math = {
floor: floor, ceil: ceil, round: round, trunc: trunc,
sin: sin, cos: cos, tan: tan, asin: asin, acos: acos, atan: atan, atan2: atan2,
pow: pow, exp: exp, log: log, log10: log10,
sqrt: sqrt, cbrt: cbrt, abs: abs, min: min, max: max, sign: sign,
random: random, randomint: randomint, randomchoice: randomchoice, randomseed: randomseed,
pi: pi, e: e, lerp: lerp, clamp: clamp,
hypot: hypot, fmod: fmod, degrees: degrees, radians: radians, wrap: wrap, approach: approach, smoothdamp: smoothdamp
};`
	l := lexer.NewLexer(src, "<builtin:math-prelude>")
	toks, err := l.Tokenize()
	if err != nil || len(toks) == 0 {
		return
	}
	p := NewParser(toks)
	prog, err := p.Parse()
	if err != nil || prog == nil || len(prog.Declarations) == 0 {
		return
	}
	bundle.Entry.Declarations = append(prog.Declarations, bundle.Entry.Declarations...)
}

// InjectColorPrelude prepends `let colors = { … };` with Raylib-ready named colors.
func InjectColorPrelude(bundle *ProgramBundle) {
	if bundle == nil || bundle.Entry == nil {
		return
	}
	for _, d := range bundle.Entry.Declarations {
		if let, ok := d.(*LetDecl); ok && strings.EqualFold(let.Name.Lexeme, "colors") {
			return
		}
	}

	src := `let colors = {
white: rgb(255, 255, 255),
black: rgb(0, 0, 0),
red: rgb(255, 0, 0),
green: rgb(0, 255, 0),
blue: rgb(0, 0, 255),
yellow: rgb(255, 255, 0),
orange: rgb(255, 165, 0),
purple: rgb(128, 0, 128),
pink: rgb(255, 192, 203),
cyan: rgb(0, 255, 255),
gray: rgb(128, 128, 128),
grey: rgb(128, 128, 128),
dark: rgb(25, 25, 25),
sky: rgba(255, 216, 168, 255),
grass: rgb(34, 139, 34),
forest: rgb(0, 100, 0),
dirt: rgb(139, 69, 19),
gold: rgb(255, 215, 0),
brown: rgb(85, 68, 34)
};`
	l := lexer.NewLexer(src, "<builtin:color-prelude>")
	toks, err := l.Tokenize()
	if err != nil || len(toks) == 0 {
		return
	}
	p := NewParser(toks)
	prog, err := p.Parse()
	if err != nil || prog == nil || len(prog.Declarations) == 0 {
		return
	}
	bundle.Entry.Declarations = append(prog.Declarations, bundle.Entry.Declarations...)
}
