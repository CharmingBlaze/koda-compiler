# Koda Reorganization Plan

## Current Structure
```
koda/
├── api/                    # Old API layer
├── cmd/
│   ├── build-release/      # Release scripts
│   ├── dist/               # Distribution artifacts
│   ├── examples/          # Example programs
│   ├── koda/               # Main compiler CLI
│   └── wrapgen/           # Wrapper generator (minimal)
├── internal/
│   ├── ast/               # AST definitions
│   ├── codegen/           # Minimal LLVM codegen
│   ├── koda/              # (empty)
│   ├── lexer/             # Lexer implementation
│   ├── parser/            # Parser implementation
│   ├── runtime/           # (empty - needs C runtime)
│   ├── sema/              # Semantic analysis
│   └── vm/                # Quarantined old VM (.wip files)
├── koda-ide/              # IDE (React frontend)
├── kodawrap/              # (empty)
├── runtime/               # (empty)
├── stdlib/                # Standard library
├── tests/                 # Tests
└── wrappers/              # Generated wrappers
```

## Proposed Production Structure
```
koda/
├── README.md
├── LICENSE
├── CHANGELOG.md
├── Makefile
├── go.mod
├── go.sum
│
├── cmd/
│   └── koda/
│       └── main.go        # CLI: koda build/run/version
│
├── internal/
│   ├── lexer/
│   │   ├── lexer.go
│   │   ├── token.go
│   │   └── lexer_test.go
│   │
│   ├── parser/
│   │   ├── parser.go
│   │   ├── ast.go
│   │   ├── precedence.go
│   │   └── parser_test.go
│   │
│   ├── sema/
│   │   ├── sema.go
│   │   ├── symbols.go
│   │   ├── scopes.go
│   │   ├── types.go
│   │   └── sema_test.go
│   │
│   ├── codegen/
│   │   ├── codegen.go
│   │   ├── types.go
│   │   ├── values.go
│   │   ├── functions.go
│   │   ├── closures.go
│   │   ├── objects.go
│   │   ├── optimize.go
│   │   ├── runtime.go
│   │   └── codegen_test.go
│   │
│   ├── platform/
│   │   ├── platform.go
│   │   ├── linux.go
│   │   ├── windows.go
│   │   └── darwin.go
│   │
│   └── version/
│       └── version.go
│
├── runtime/               # NEW: C runtime
│   ├── koda_runtime.c
│   ├── koda_runtime.h
│   ├── gc.c
│   ├── gc.h
│   ├── value.c
│   ├── value.h
│   ├── object.c
│   ├── object.h
│   ├── array.c
│   ├── string.c
│   ├── table.c
│   ├── natives.c
│   ├── natives.h
│   └── Makefile
│
├── kodawrap/              # NEW: Wrapper generator
│   ├── cmd/
│   │   └── kodawrap/
│   │       └── main.go
│   ├── internal/
│   │   ├── parser/
│   │   │   └── clang.go
│   │   ├── analyzer/
│   │   ├── codegen/
│   │   └── ffi/
│   ├── templates/
│   └── examples/
│
├── stdlib/
│   ├── prelude.koda
│   ├── array.koda
│   ├── string.koda
│   ├── math.koda
│   ├── time.koda
│   └── random.koda
│
├── examples/
│   ├── hello.koda
│   ├── fibonacci.koda
│   └── breakout/
│
├── tests/
│   ├── compiler/
│   ├── language/
│   ├── stdlib/
│   ├── gc/
│   └── benchmarks/
│
├── docs/
│   ├── language/
│   ├── compiler/
│   └── guides/
│
├── scripts/
│   ├── build.sh
│   ├── build_release.sh
│   ├── test.sh
│   └── cross_compile.sh
│
└── dist/                 # Release artifacts (gitignored)
```

## Migration Actions

### Phase 1: Internal Reorganization (Day 1-2)
1. **Move AST from internal/ast to internal/parser**
   - `internal/ast/ast.go` → `internal/parser/ast.go`
   - Update all imports

2. **Split codegen into modules**
   - Current `codegen.go` (273 lines) → split into:
     - `codegen.go` (main entry)
     - `types.go` (LLVM type definitions)
     - `values.go` (value compilation)
     - `functions.go` (function compilation)
     - `runtime.go` (runtime function declarations)

3. **Add platform abstraction**
   - Create `internal/platform/` package
   - Add `linux.go`, `windows.go`, `darwin.go`
   - Abstract linking commands

4. **Remove quarantined code**
   - Delete `internal/vm/*.go.wip` files (already quarantined)
   - Keep old codegen.go.wip for reference if needed

### Phase 2: Runtime Setup (Day 3-5)
1. **Create runtime/ directory structure**
2. **Port existing runtime concepts**
   - From old VM: value system, object model
   - Implement NaN-boxing in C
   - Implement basic object allocators

3. **Implement minimal GC**
   - Start with simple mark-sweep
   - Later upgrade to generational

### Phase 3: Codegen Completion (Week 2)
1. **Fix function call emission**
   - Implement function lookup in module
   - Handle function pointers correctly

2. **Implement variable lookup**
   - Track locals in stack slots
   - Handle upvalues for closures

3. **Implement object operations**
   - Object allocation
   - Property access
   - Method calls

### Phase 4: Linking & Executables (Week 2-3)
1. **Implement linker**
   - LLVM IR → object file
   - Object file + runtime → executable
   - Static linking for Linux (musl)
   - Static linking for Windows (MinGW)

2. **Add CLI commands**
   - `koda build file.koda -o output`
   - `koda run file.koda`
   - `koda version`

### Phase 5: Kodawrap (Week 4-5)
1. **Move cmd/wrapgen to kodawrap/**
2. **Implement libclang integration**
3. **Generate professional wrappers**
4. **Add documentation generation**

### Phase 6: Distribution (Week 6)
1. **Implement cross-platform builds**
2. **Create release scripts**
3. **Set up CI/CD**

## Breaking Changes
- AST import path changes (`internal/ast` → `internal/parser`)
- Codegen API changes (split into modules)
- CLI interface changes (unified `koda` command)
