package codegen

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"koda/internal/lexer"
	"koda/internal/parser"
)

// emitStmt emits LLVM IR for statements.
func (g *Generator) emitStmt(stmt parser.Stmt) error {
	switch s := stmt.(type) {
	case *parser.BlockStmt:
		return g.emitBlockStmt(s)
	case *parser.ExpressionStmt:
		_, err := g.emitExpr(s.Expr)
		return err
	case *parser.ReturnStmt:
		return g.emitReturnStmt(s)
	case *parser.IfStmt:
		return g.emitIfStmt(s)
	case *parser.WhileStmt:
		return g.emitWhileStmt(s)
	case *parser.DoWhileStmt:
		return g.emitDoWhileStmt(s)
	case *parser.ForOfStmt:
		return g.emitForOfStmt(s)
	case *parser.ForStmt:
		return g.emitForStmt(s)
	case *parser.ForInStmt:
		return g.emitForInStmt(s)
	case *parser.BreakStmt:
		return g.emitBreakStmt(s)
	case *parser.ContinueStmt:
		return g.emitContinueStmt(s)
	case *parser.SwitchStmt:
		return g.emitSwitchStmt(s)
	case *parser.DeleteStmt:
		return g.emitDeleteStmt(s)
	case *parser.DeferStmt:
		return g.emitDeferStmt(s)
	default:
		return fmt.Errorf("unsupported statement type: %T", stmt)
	}
}

// emitBlockStmt emits LLVM IR for block statements.
func (g *Generator) emitBlockStmt(s *parser.BlockStmt) error {
	for _, decl := range s.Declarations {
		if err := g.emitDecl(decl); err != nil {
			return err
		}
	}
	return nil
}

// emitReturnStmt emits LLVM IR for return statements.
// Return value is evaluated before defers run (same order as Go), so defers observe the outgoing result slot only after their side effects.
func (g *Generator) emitReturnStmt(s *parser.ReturnStmt) error {
	if s.Value != nil {
		val, err := g.emitExpr(s.Value)
		if err != nil {
			return err
		}
		retSlot := g.entryAlloca(types.I64)
		g.block.NewStore(g.emitAsKodaI64(val), retSlot)
		if err := g.emitDefersForCurrentLayer(); err != nil {
			return err
		}
		g.emitCallTracePop()
		g.emitShadowPop()
		g.block.NewRet(g.block.NewLoad(types.I64, retSlot))
		return nil
	}
	if err := g.emitDefersForCurrentLayer(); err != nil {
		return err
	}
	g.emitCallTracePop()
	g.emitShadowPop()
	g.block.NewRet(constant.NewInt(types.I64, 0))
	return nil
}

func (g *Generator) pushDeferLayer() {
	g.deferLayers = append(g.deferLayers, nil)
}

func (g *Generator) popDeferLayer() {
	if len(g.deferLayers) == 0 {
		return
	}
	g.deferLayers = g.deferLayers[:len(g.deferLayers)-1]
}

func (g *Generator) emitDefersForCurrentLayer() error {
	if len(g.deferLayers) == 0 {
		return nil
	}
	i := len(g.deferLayers) - 1
	layer := g.deferLayers[i]
	for j := len(layer) - 1; j >= 0; j-- {
		if _, err := g.emitExpr(layer[j]); err != nil {
			return err
		}
	}
	g.deferLayers[i] = nil
	return nil
}

func (g *Generator) emitDeferStmt(s *parser.DeferStmt) error {
	if len(g.deferLayers) == 0 {
		return fmt.Errorf("defer outside function body")
	}
	idx := len(g.deferLayers) - 1
	g.deferLayers[idx] = append(g.deferLayers[idx], s.Expr)
	return nil
}

// emitIfStmt emits LLVM IR for if statements.
func (g *Generator) emitIfStmt(s *parser.IfStmt) error {
	cond, err := g.emitExpr(s.Condition)
	if err != nil {
		return err
	}
	g.shadowRewindTemps()

	g.tempN++
	thenBlock := g.block.Parent.NewBlock(fmt.Sprintf("then.%d", g.tempN))
	elseBlock := g.block.Parent.NewBlock(fmt.Sprintf("else.%d", g.tempN))
	mergeBlock := g.block.Parent.NewBlock(fmt.Sprintf("merge.%d", g.tempN))

	g.block.NewCondBr(g.emitTruthy(cond), thenBlock, elseBlock)

	g.block = thenBlock
	if err := g.emitStmt(s.Then); err != nil {
		return err
	}
	if g.block.Term == nil {
		g.block.NewBr(mergeBlock)
	}

	g.block = elseBlock
	g.shadowRewindTemps()
	if s.Else != nil {
		if err := g.emitStmt(s.Else); err != nil {
			return err
		}
	}
	if g.block.Term == nil {
		g.block.NewBr(mergeBlock)
	}

	g.block = mergeBlock
	g.shadowRewindTemps()
	return nil
}

// emitWhileStmt emits LLVM IR for while loops.
func (g *Generator) emitWhileStmt(s *parser.WhileStmt) error {
	g.tempN++
	suf := fmt.Sprintf(".%d", g.tempN)
	condBlock := g.block.Parent.NewBlock("while.cond" + suf)
	bodyBlock := g.block.Parent.NewBlock("while.body" + suf)
	afterBlock := g.block.Parent.NewBlock("while.after" + suf)

	// Push loop context onto stack
	ctx := loopContext{condBlock: condBlock, incBlock: condBlock, afterBlock: afterBlock}
	g.loopStack = append(g.loopStack, ctx)

	g.block.NewBr(condBlock)

	g.block = condBlock
	g.shadowRewindTemps()
	cond, err := g.emitExpr(s.Condition)
	if err != nil {
		return err
	}
	g.block.NewCondBr(g.emitTruthy(cond), bodyBlock, afterBlock)

	g.block = bodyBlock
	g.shadowRewindTemps()
	if err := g.emitStmt(s.Body); err != nil {
		return err
	}
	if g.block.Term == nil {
		g.block.NewBr(condBlock)
	}

	g.block = afterBlock
	g.shadowRewindTemps()

	// Pop loop context from stack
	g.loopStack = g.loopStack[:len(g.loopStack)-1]

	return nil
}

// emitDoWhileStmt emits LLVM IR for do-while loops.
func (g *Generator) emitDoWhileStmt(s *parser.DoWhileStmt) error {
	g.tempN++
	suf := fmt.Sprintf(".%d", g.tempN)
	bodyBlock := g.block.Parent.NewBlock("dowhile.body" + suf)
	condBlock := g.block.Parent.NewBlock("dowhile.cond" + suf)
	afterBlock := g.block.Parent.NewBlock("dowhile.after" + suf)

	// Push loop context onto stack
	ctx := loopContext{condBlock: condBlock, incBlock: condBlock, afterBlock: afterBlock}
	g.loopStack = append(g.loopStack, ctx)

	g.block.NewBr(bodyBlock)

	g.block = bodyBlock
	g.shadowRewindTemps()
	if err := g.emitStmt(s.Body); err != nil {
		return err
	}
	if g.block.Term == nil {
		g.block.NewBr(condBlock)
	}

	g.block = condBlock
	g.shadowRewindTemps()
	cond, err := g.emitExpr(s.Condition)
	if err != nil {
		return err
	}
	g.block.NewCondBr(g.emitTruthy(cond), bodyBlock, afterBlock)

	g.block = afterBlock
	g.shadowRewindTemps()

	// Pop loop context from stack
	g.loopStack = g.loopStack[:len(g.loopStack)-1]

	return nil
}

// emitForOfStmt emits LLVM IR for for-of loops over arrays and tables (slot order).
func (g *Generator) emitForOfStmt(s *parser.ForOfStmt) error {
	if s.ValueVar == nil {
		if r, ok := s.Iterable.(*parser.RangeExpr); ok {
			if from, to, ok := rangeConstBounds(r); ok {
				return g.emitForOfConstRange(s, from, to)
			}
			return g.emitForOfDynamicRange(s, r)
		}
	}

	iterable, err := g.emitExpr(s.Iterable)
	if err != nil {
		return err
	}
	iterI := g.emitAsKodaI64(iterable)
	g.shadowRewindTemps()

	lenVal := g.block.NewCall(g.runtimeForOfLength, iterI)

	idxSlot := g.entryAlloca(types.I64)
	zeroNum := g.emitAsKodaI64(constant.NewFloat(types.Double, 0))
	g.block.NewStore(zeroNum, idxSlot)

	valSlot := g.entryAlloca(types.I64)
	var keySlot value.Value
	if s.ValueVar != nil {
		keySlot = g.entryAlloca(types.I64)
		g.locals[s.VarName.Lexeme] = keySlot
		g.locals[s.ValueVar.Lexeme] = valSlot
	} else {
		g.locals[s.VarName.Lexeme] = valSlot
	}

	g.tempN++
	condBlock := g.block.Parent.NewBlock(fmt.Sprintf("forof.cond.%d", g.tempN))
	bodyBlock := g.block.Parent.NewBlock(fmt.Sprintf("forof.body.%d", g.tempN))
	incBlock := g.block.Parent.NewBlock(fmt.Sprintf("forof.inc.%d", g.tempN))
	afterBlock := g.block.Parent.NewBlock(fmt.Sprintf("forof.after.%d", g.tempN))

	ctx := loopContext{condBlock: condBlock, incBlock: incBlock, afterBlock: afterBlock}
	g.loopStack = append(g.loopStack, ctx)

	g.block.NewBr(condBlock)

	g.block = condBlock
	idxV := g.block.NewLoad(types.I64, idxSlot)
	idxD := g.block.NewBitCast(idxV, types.Double)
	lenD := g.block.NewBitCast(lenVal, types.Double)
	cmp := g.block.NewFCmp(enum.FPredOLT, idxD, lenD)
	g.block.NewCondBr(cmp, bodyBlock, afterBlock)

	g.block = bodyBlock
	if s.ValueVar != nil {
		key := g.block.NewCall(g.runtimeForOfKeyAt, iterI, idxV)
		val := g.block.NewCall(g.runtimeForOfValueAt, iterI, idxV)
		g.block.NewStore(key, keySlot)
		g.block.NewStore(val, valSlot)
	} else {
		val := g.block.NewCall(g.runtimeForOfValueAt, iterI, idxV)
		g.block.NewStore(val, valSlot)
	}
	if err := g.emitStmt(s.Body); err != nil {
		return err
	}
	g.block.NewBr(incBlock)

	g.block = incBlock
	idxV2 := g.block.NewLoad(types.I64, idxSlot)
	idxD2 := g.block.NewBitCast(idxV2, types.Double)
	nextD := g.block.NewFAdd(idxD2, constant.NewFloat(types.Double, 1))
	nextI := g.block.NewBitCast(nextD, types.I64)
	g.block.NewStore(nextI, idxSlot)
	g.block.NewBr(condBlock)

	g.block = afterBlock

	g.loopStack = g.loopStack[:len(g.loopStack)-1]
	return nil
}

func rangeConstBounds(r *parser.RangeExpr) (int64, int64, bool) {
	fromLit, ok := r.From.(*parser.LiteralExpr)
	if !ok {
		return 0, 0, false
	}
	toLit, ok := r.To.(*parser.LiteralExpr)
	if !ok {
		return 0, 0, false
	}
	var from int64
	switch v := fromLit.Value.(type) {
	case int:
		from = int64(v)
	case float64:
		from = int64(v)
	default:
		return 0, 0, false
	}
	var to int64
	switch v := toLit.Value.(type) {
	case int:
		to = int64(v)
	case float64:
		to = int64(v)
	default:
		return 0, 0, false
	}
	return from, to, true
}

func (g *Generator) emitForOfConstRange(s *parser.ForOfStmt, from int64, to int64) error {
	idxSlot := g.entryAlloca(types.I64)
	valSlot := g.entryAlloca(types.I64)
	g.locals[s.VarName.Lexeme] = valSlot
	g.block.NewStore(constant.NewInt(types.I64, from), idxSlot)

	g.tempN++
	condBlock := g.block.Parent.NewBlock(fmt.Sprintf("forof.range.cond.%d", g.tempN))
	bodyBlock := g.block.Parent.NewBlock(fmt.Sprintf("forof.range.body.%d", g.tempN))
	incBlock := g.block.Parent.NewBlock(fmt.Sprintf("forof.range.inc.%d", g.tempN))
	afterBlock := g.block.Parent.NewBlock(fmt.Sprintf("forof.range.after.%d", g.tempN))
	ctx := loopContext{condBlock: condBlock, incBlock: incBlock, afterBlock: afterBlock}
	g.loopStack = append(g.loopStack, ctx)

	g.block.NewBr(condBlock)

	g.block = condBlock
	cur := g.block.NewLoad(types.I64, idxSlot)
	cmp := g.block.NewICmp(enum.IPredSLT, cur, constant.NewInt(types.I64, to))
	g.block.NewCondBr(cmp, bodyBlock, afterBlock)

	g.block = bodyBlock
	curD := g.block.NewSIToFP(cur, types.Double)
	curBoxed := g.block.NewCall(g.runtimeBoxNumber, curD)
	g.block.NewStore(curBoxed, valSlot)
	if err := g.emitStmt(s.Body); err != nil {
		return err
	}
	g.block.NewBr(incBlock)

	g.block = incBlock
	next := g.block.NewAdd(g.block.NewLoad(types.I64, idxSlot), constant.NewInt(types.I64, 1))
	g.block.NewStore(next, idxSlot)
	g.block.NewBr(condBlock)

	g.block = afterBlock
	g.loopStack = g.loopStack[:len(g.loopStack)-1]
	return nil
}

// emitForOfDynamicRange lowers `for (let i of from..to)` when bounds are not both numeric literals,
// using a counted i64 loop and no KODA_range allocation (half-open [from, to) like emitRange).
func (g *Generator) emitForOfDynamicRange(s *parser.ForOfStmt, r *parser.RangeExpr) error {
	fromVal, err := g.emitExpr(r.From)
	if err != nil {
		return err
	}
	toVal, err := g.emitExpr(r.To)
	if err != nil {
		return err
	}
	// Bounds are Koda values (NaN-boxed); the loop counter uses plain i64 like emitForOfConstRange.
	fromBoxed := g.emitAsKodaI64(fromVal)
	toBoxed := g.emitAsKodaI64(toVal)
	fromInt := g.block.NewFPToSI(g.block.NewCall(g.runtimeUnboxNumber, fromBoxed), types.I64)
	toInt := g.block.NewFPToSI(g.block.NewCall(g.runtimeUnboxNumber, toBoxed), types.I64)

	fromSlot := g.entryAlloca(types.I64)
	toSlot := g.entryAlloca(types.I64)
	g.block.NewStore(fromInt, fromSlot)
	g.block.NewStore(toInt, toSlot)

	idxSlot := g.entryAlloca(types.I64)
	valSlot := g.entryAlloca(types.I64)
	g.locals[s.VarName.Lexeme] = valSlot
	g.block.NewStore(g.block.NewLoad(types.I64, fromSlot), idxSlot)
	g.shadowRewindTemps()

	g.tempN++
	condBlock := g.block.Parent.NewBlock(fmt.Sprintf("forof.dyn.cond.%d", g.tempN))
	bodyBlock := g.block.Parent.NewBlock(fmt.Sprintf("forof.dyn.body.%d", g.tempN))
	incBlock := g.block.Parent.NewBlock(fmt.Sprintf("forof.dyn.inc.%d", g.tempN))
	afterBlock := g.block.Parent.NewBlock(fmt.Sprintf("forof.dyn.after.%d", g.tempN))
	ctx := loopContext{condBlock: condBlock, incBlock: incBlock, afterBlock: afterBlock}
	g.loopStack = append(g.loopStack, ctx)

	g.block.NewBr(condBlock)

	g.block = condBlock
	cur := g.block.NewLoad(types.I64, idxSlot)
	limit := g.block.NewLoad(types.I64, toSlot)
	cmp := g.block.NewICmp(enum.IPredSLT, cur, limit)
	g.block.NewCondBr(cmp, bodyBlock, afterBlock)

	g.block = bodyBlock
	curD := g.block.NewSIToFP(cur, types.Double)
	curBoxed := g.block.NewCall(g.runtimeBoxNumber, curD)
	g.block.NewStore(curBoxed, valSlot)
	if err := g.emitStmt(s.Body); err != nil {
		return err
	}
	g.block.NewBr(incBlock)

	g.block = incBlock
	next := g.block.NewAdd(g.block.NewLoad(types.I64, idxSlot), constant.NewInt(types.I64, 1))
	g.block.NewStore(next, idxSlot)
	g.block.NewBr(condBlock)

	g.block = afterBlock
	g.loopStack = g.loopStack[:len(g.loopStack)-1]
	return nil
}

// emitForStmt emits LLVM IR for C-style for loops.
func (g *Generator) emitForStmt(s *parser.ForStmt) error {
	// Create blocks
	g.tempN++
	suf := fmt.Sprintf(".%d", g.tempN)
	condBlock := g.block.Parent.NewBlock("for.cond" + suf)
	bodyBlock := g.block.Parent.NewBlock("for.body" + suf)
	incBlock := g.block.Parent.NewBlock("for.inc" + suf)
	afterBlock := g.block.Parent.NewBlock("for.after" + suf)

	// Push loop context onto stack
	ctx := loopContext{condBlock: condBlock, incBlock: incBlock, afterBlock: afterBlock}
	g.loopStack = append(g.loopStack, ctx)

	// Emit initialization
	for _, init := range s.Inits {
		if err := g.emitDecl(init); err != nil {
			return err
		}
	}

	// Jump to condition block
	g.block.NewBr(condBlock)

	// Condition block
	g.block = condBlock
	g.shadowRewindTemps()
	if s.Condition != nil {
		cond, err := g.emitExpr(s.Condition)
		if err != nil {
			return err
		}
		g.block.NewCondBr(g.emitTruthy(cond), bodyBlock, afterBlock)
	} else {
		// No condition means always enter loop
		g.block.NewBr(bodyBlock)
	}

	// Body block
	g.block = bodyBlock
	g.shadowRewindTemps()
	if err := g.emitStmt(s.Body); err != nil {
		return err
	}
	g.block.NewBr(incBlock)

	// Increment block
	g.block = incBlock
	g.shadowRewindTemps()
	for _, inc := range s.Increments {
		_, err := g.emitExpr(inc)
		if err != nil {
			return err
		}
	}
	g.block.NewBr(condBlock)

	// After block
	g.block = afterBlock
	g.shadowRewindTemps()

	// Pop loop context from stack
	g.loopStack = g.loopStack[:len(g.loopStack)-1]

	return nil
}

// emitForInStmt emits LLVM IR for for-in loops (keys in insertion order for tables).
func (g *Generator) emitForInStmt(s *parser.ForInStmt) error {
	if s.ValueVar != nil {
		return fmt.Errorf("native codegen: for-in with two variables is not supported")
	}
	if s.KeyVar == nil {
		return fmt.Errorf("native codegen: for-in missing loop variable")
	}
	if r, ok := s.Iterable.(*parser.RangeExpr); ok {
		fake := &parser.ForOfStmt{
			Token:    s.Token,
			VarName:  *s.KeyVar,
			Iterable: r,
			Body:     s.Body,
		}
		return g.emitForOfStmt(fake)
	}
	iterable, err := g.emitExpr(s.Iterable)
	if err != nil {
		return err
	}
	iterI := g.emitAsKodaI64(iterable)
	g.shadowRewindTemps()

	lenVal := g.block.NewCall(g.runtimeForOfLength, iterI)

	idxSlot := g.entryAlloca(types.I64)
	zeroNum := g.emitAsKodaI64(constant.NewFloat(types.Double, 0))
	g.block.NewStore(zeroNum, idxSlot)

	keySlot := g.entryAlloca(types.I64)
	g.locals[s.KeyVar.Lexeme] = keySlot

	g.tempN++
	condBlock := g.block.Parent.NewBlock(fmt.Sprintf("forin.cond.%d", g.tempN))
	bodyBlock := g.block.Parent.NewBlock(fmt.Sprintf("forin.body.%d", g.tempN))
	incBlock := g.block.Parent.NewBlock(fmt.Sprintf("forin.inc.%d", g.tempN))
	afterBlock := g.block.Parent.NewBlock(fmt.Sprintf("forin.after.%d", g.tempN))

	ctx := loopContext{condBlock: condBlock, incBlock: incBlock, afterBlock: afterBlock}
	g.loopStack = append(g.loopStack, ctx)

	g.block.NewBr(condBlock)

	g.block = condBlock
	idxV := g.block.NewLoad(types.I64, idxSlot)
	idxD := g.block.NewBitCast(idxV, types.Double)
	lenD := g.block.NewBitCast(lenVal, types.Double)
	cmp := g.block.NewFCmp(enum.FPredOLT, idxD, lenD)
	g.block.NewCondBr(cmp, bodyBlock, afterBlock)

	g.block = bodyBlock
	key := g.block.NewCall(g.runtimeForOfKeyAt, iterI, idxV)
	g.block.NewStore(key, keySlot)
	if err := g.emitStmt(s.Body); err != nil {
		return err
	}
	g.block.NewBr(incBlock)

	g.block = incBlock
	idxV2 := g.block.NewLoad(types.I64, idxSlot)
	idxD2 := g.block.NewBitCast(idxV2, types.Double)
	nextD := g.block.NewFAdd(idxD2, constant.NewFloat(types.Double, 1))
	nextI := g.block.NewBitCast(nextD, types.I64)
	g.block.NewStore(nextI, idxSlot)
	g.block.NewBr(condBlock)

	g.block = afterBlock

	g.loopStack = g.loopStack[:len(g.loopStack)-1]
	return nil
}

// emitBreakStmt emits LLVM IR for break statements.
func (g *Generator) emitBreakStmt(_ *parser.BreakStmt) error {
	// Break statement - jump to the loop's after block
	if len(g.loopStack) == 0 {
		return fmt.Errorf("break statement outside of loop")
	}
	ctx := g.loopStack[len(g.loopStack)-1]
	g.block.NewBr(ctx.afterBlock)
	return nil
}

// emitContinueStmt emits LLVM IR for continue statements.
func (g *Generator) emitContinueStmt(_ *parser.ContinueStmt) error {
	// Continue statement - jump to the loop's increment or condition block
	if len(g.loopStack) == 0 {
		return fmt.Errorf("continue statement outside of loop")
	}
	ctx := g.loopStack[len(g.loopStack)-1]
	g.block.NewBr(ctx.incBlock)
	return nil
}

// emitSwitchStmt emits LLVM IR for switch/match statements.
func (g *Generator) emitSwitchStmt(s *parser.SwitchStmt) error {
	subject, err := g.emitExpr(s.Subject)
	if err != nil {
		return err
	}
	subjI := g.emitAsKodaI64(subject)
	g.shadowRewindTemps()

	g.tempN++
	suf := fmt.Sprintf(".%d", g.tempN)
	mergeBlock := g.block.Parent.NewBlock("switch.merge" + suf)
	isMatch := s.Token.Type == lexer.TokenMatch

	swCtx := loopContext{condBlock: mergeBlock, incBlock: mergeBlock, afterBlock: mergeBlock}
	g.loopStack = append(g.loopStack, swCtx)
	defer func() { g.loopStack = g.loopStack[:len(g.loopStack)-1] }()

	caseBodyBlocks := make([]*ir.Block, len(s.Cases))
	for i := range s.Cases {
		caseBodyBlocks[i] = g.block.Parent.NewBlock(fmt.Sprintf("switch.case.%d.body%s", i, suf))
	}

	var defaultBodyBlock *ir.Block
	if len(s.Default) > 0 || len(s.Cases) == 0 {
		defaultBodyBlock = g.block.Parent.NewBlock("switch.default.body" + suf)
	} else {
		defaultBodyBlock = mergeBlock
	}

	if len(s.Cases) == 0 {
		g.block.NewBr(defaultBodyBlock)
	} else {
		cmpBlock := g.block
		for i, caseStmt := range s.Cases {
			caseVal, err := g.emitExpr(caseStmt.Value)
			if err != nil {
				return err
			}
			cmp := cmpBlock.NewICmp(enum.IPredEQ, subjI, g.emitAsKodaI64(caseVal))
			var failBlock *ir.Block
			if i < len(s.Cases)-1 {
				failBlock = cmpBlock.Parent.NewBlock(fmt.Sprintf("switch.cmp.%d%s", i+1, suf))
			} else {
				failBlock = defaultBodyBlock
			}
			cmpBlock.NewCondBr(cmp, caseBodyBlocks[i], failBlock)
			cmpBlock = failBlock
		}
	}

	if defaultBodyBlock != mergeBlock {
		g.block = defaultBodyBlock
		g.shadowRewindTemps()
		for _, decl := range s.Default {
			if err := g.emitDecl(decl); err != nil {
				return err
			}
		}
		g.block.NewBr(mergeBlock)
	}

	for i, caseStmt := range s.Cases {
		g.block = caseBodyBlocks[i]
		g.shadowRewindTemps()
		for _, decl := range caseStmt.Body {
			if err := g.emitDecl(decl); err != nil {
				return err
			}
		}
		if isMatch {
			g.block.NewBr(mergeBlock)
		} else if i < len(s.Cases)-1 {
			g.block.NewBr(caseBodyBlocks[i+1])
		} else {
			g.block.NewBr(mergeBlock)
		}
	}

	g.block = mergeBlock
	g.shadowRewindTemps()
	return nil
}

func (g *Generator) emitDeleteStmt(s *parser.DeleteStmt) error {
	ix, ok := s.Target.(*parser.IndexExpr)
	if !ok {
		return fmt.Errorf("delete expects a property access expression such as obj[\"key\"]")
	}
	obj, err := g.emitExpr(ix.Object)
	if err != nil {
		return err
	}
	key, err := g.emitExpr(ix.Index)
	if err != nil {
		return err
	}
	g.block.NewCall(g.runtimeObjRemove, g.emitAsKodaI64(obj), g.emitAsKodaI64(key))
	return nil
}
