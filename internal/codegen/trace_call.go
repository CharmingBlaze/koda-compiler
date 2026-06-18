package codegen

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// cStringPtr returns an i8* to a private, immutable, null-terminated global for C interop.
func (g *Generator) cStringPtr(s string) value.Value {
	arr := constant.NewCharArrayFromString(s)
	gl := g.mod.NewGlobalDef("", arr)
	gl.Immutable = true
	gl.Linkage = enum.LinkagePrivate
	zero := constant.NewInt(types.I32, 0)
	return g.block.NewGetElementPtr(arr.Type(), gl, zero, zero)
}

func (g *Generator) emitCallTracePush(fnName, fileName string, line int) {
	if g.ctx == nil || !g.ctx.EmitDebug {
		return
	}
	if g.runtimePushCall == nil || g.block == nil {
		return
	}
	fnPtr := g.cStringPtr(fnName)
	filePtr := g.cStringPtr(fileName)
	g.block.NewCall(g.runtimePushCall, fnPtr, filePtr, constant.NewInt(types.I32, int64(line)))
}

func (g *Generator) emitCallTracePop() {
	if g.ctx == nil || !g.ctx.EmitDebug {
		return
	}
	if g.runtimePopCall == nil || g.block == nil {
		return
	}
	g.block.NewCall(g.runtimePopCall)
}
