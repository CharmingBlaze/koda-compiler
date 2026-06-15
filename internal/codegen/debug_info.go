package codegen

import (
	"path/filepath"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/metadata"
)

// debugInfo attaches DWARF-style metadata when EmitDebug is enabled.
type debugInfo struct {
	mod      *ir.Module
	files    map[string]*metadata.DIFile
	cu       *metadata.DICompileUnit
	subprogs map[string]*metadata.DISubprogram
}

func newDebugInfo(mod *ir.Module, sourcePath string) *debugInfo {
	d := &debugInfo{
		mod:      mod,
		files:    make(map[string]*metadata.DIFile),
		subprogs: make(map[string]*metadata.DISubprogram),
	}
	file := d.fileFor(sourcePath)
	d.cu = &metadata.DICompileUnit{
		Language: enum.DwarfLangC89,
		File:     file,
		Producer: "koda",
	}
	mod.MetadataDefs = append(mod.MetadataDefs, d.cu)
	return d
}

func (d *debugInfo) fileFor(path string) *metadata.DIFile {
	if path == "" {
		path = "main.koda"
	}
	key := strings.ToLower(path)
	if f, ok := d.files[key]; ok {
		return f
	}
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	f := &metadata.DIFile{
		Filename:  base,
		Directory: dir,
	}
	d.mod.MetadataDefs = append(d.mod.MetadataDefs, f)
	d.files[key] = f
	return f
}

func (d *debugInfo) subprogram(fn *ir.Func, name, file string, line int) {
	if fn == nil {
		return
	}
	sp := &metadata.DISubprogram{
		Name:        name,
		LinkageName: fn.Name(),
		File:        d.fileFor(file),
		Line:        int64(line),
		Unit:        d.cu,
	}
	d.mod.MetadataDefs = append(d.mod.MetadataDefs, sp)
	d.subprogs[strings.ToLower(name)] = sp
	fn.Metadata = ir.Metadata{&metadata.Attachment{Name: "dbg", Node: sp}}
}

func (d *debugInfo) loc(file string, line int) *metadata.DILocation {
	if line <= 0 {
		line = 1
	}
	sp := d.subprogs["user_main"]
	if sp == nil && len(d.subprogs) > 0 {
		for _, v := range d.subprogs {
			sp = v
			break
		}
	}
	loc := &metadata.DILocation{
		Line:  int64(line),
		Scope: sp,
	}
	d.mod.MetadataDefs = append(d.mod.MetadataDefs, loc)
	return loc
}

func attachDbg(inst interface{}, loc *metadata.DILocation) {
	if loc == nil {
		return
	}
	md := ir.Metadata{&metadata.Attachment{Name: "dbg", Node: loc}}
	switch i := inst.(type) {
	case *ir.InstCall:
		i.Metadata = md
	case *ir.InstStore:
		i.Metadata = md
	case *ir.TermRet:
		i.Metadata = md
	}
}
