# Koda Implementation Roadmap

**Execution model:** this repository targets **one pipeline** — LLVM IR emission (`internal/codegen`) plus the C runtime (`runtime/src`). There is no bytecode interpreter in-tree.

## Critical Path Analysis

### Dependencies
```
Runtime (C) → Codegen (LLVM) → Linking → Executables
     ↓
   Kodawrap (depends on working compiler)
     ↓
Distribution
```

### Risk Areas
- **High Risk**: C runtime implementation (GC correctness, performance)
- **High Risk**: LLVM codegen (function calls, closures, object model)
- **Medium Risk**: Cross-platform linking (musl, MinGW, lipo)
- **Low Risk**: Kodawrap (can iterate independently)

## Week-by-Week Plan

### Week 1: Codegen Completion
**Goal**: Generate working LLVM IR that compiles to executables

**Day 1-2: Function Calls & Variable Lookup**
- Implement function lookup in module
- Fix identifier emission to return function pointers
- Implement local variable stack allocation
- Handle parameter passing
- Test: simple function calls work

**Day 3-4: Object Model**
- Implement object allocation IR
- Implement property access (hash table lookup)
- Implement object literal emission
- Test: object creation and property access

**Day 5: Closures & Upvalues**
- Implement closure allocation
- Implement upvalue capture
- Implement upvalue read/write
- Test: closures capture variables correctly

**Deliverable**: Codegen can emit complete LLVM IR for basic programs

### Week 2: C Runtime Implementation
**Goal**: Implement runtime that can execute LLVM IR

**Day 1-2: Value System**
- Implement NaN-boxing in C
- Implement value type predicates (IS_OBJ, IS_NUMBER, etc.)
- Implement value conversion functions
- Test: value operations work correctly

**Day 3-4: Object System**
- Implement object allocator
- Implement hash table for properties
- Implement string and array objects
- Test: object allocation and property access

**Day 5: Basic GC**
- Implement simple mark-sweep GC
- Implement root marking (stack, globals)
- Test: GC doesn't crash, frees memory

**Deliverable**: Runtime can execute basic LLVM IR

### Week 3: Linking & Executables
**Goal**: Generate standalone executables

**Day 1-2: LLVM IR → Object File**
- Integrate llc (LLVM compiler) to compile IR to object
- Handle platform-specific object formats
- Test: IR compiles to .o file

**Day 3-4: Linker Implementation**
- Implement linking with runtime library
- Handle static linking (musl on Linux, MinGW on Windows)
- Implement platform-specific linker commands
- Test: object file links to executable

**Day 5: CLI Integration**
- Implement `koda build` command
- Implement `koda run` command (native temp binary)
- Implement `koda version` command
- Test: `koda build hello.koda -o hello` produces working executable

**Deliverable**: Can compile and run Koda programs

### Week 4: GC Optimization
**Goal**: GC pauses <5ms for 99% of collections

**Day 1-2: Generational GC**
- Implement nursery (Gen 0) with bump allocator
- Implement young generation (Gen 1)
- Implement old generation (Gen 2)
- Test: minor collections work

**Day 3-4: Write Barrier**
- Implement card marking or remembered set
- Track old→young pointers
- Test: write barrier prevents leaks

**Day 5: Performance Tuning**
- Measure GC pause times
- Tune generation sizes
- Add GC statistics API
- Test: GC meets performance targets

**Deliverable**: GC meets performance requirements

### Week 5: Standard Library
**Goal**: Complete standard library

**Day 1-2: Core Functions**
- Implement print, type, typeof
- Implement assert, error
- Test: core functions work

**Day 3-4: Array & String Methods**
- Implement array methods (push, pop, slice, map, filter)
- Implement string methods (split, join, trim, substring)
- Test: stdlib methods work correctly

**Day 5: Math & Time**
- Implement math functions (sin, cos, tan, sqrt, abs, min, max)
- Implement time functions (time, sleep, deltaTime)
- Test: math and time functions work

**Deliverable**: Complete standard library

### Week 6: Kodawrap Implementation
**Goal**: Generate professional C/C++ bindings

**Day 1-2: libclang Integration**
- Implement header parsing with libclang
- Extract functions, structs, enums, macros
- Test: can parse C headers

**Day 3-4: Code Generation**
- Generate .koda wrapper files
- Generate .md documentation
- Implement type mapping (C → Koda)
- Test: generates working wrapper

**Day 5: Examples**
- Generate Raylib bindings
- Generate SDL2 bindings
- Test: bindings compile and work

**Deliverable**: Kodawrap generates working bindings

### Week 7: Distribution & Release
**Goal**: Ship cross-platform binaries

**Day 1-2: Cross-Platform Builds**
- Implement Linux build (musl static)
- Implement Windows build (MinGW static)
- Implement macOS universal binary (lipo)
- Test: binaries run on all platforms

**Day 3-4: Release Automation**
- Implement build scripts
- Set up GitHub Actions CI/CD
- Automate release process
- Test: release builds work

**Day 5: Documentation & Polish**
- Write getting started guide
- Add examples
- Update README
- Test: documentation is clear

**Deliverable**: v1.0.0 release

## Immediate Next Steps (Today)

1. **Complete codegen function calls** (Critical Path)
   - Implement function lookup in module
   - Fix identifier emission
   - Test function calls work

2. **Implement variable lookup** (Critical Path)
   - Track locals in stack slots
   - Handle parameter passing
   - Test variable access works

3. **Add runtime stub** (Critical Path)
   - Create runtime/ directory
   - Add basic C runtime files
   - Implement minimal value system

## Testing Strategy

### Unit Tests
- Each package has comprehensive unit tests
- Fuzz testing for lexer/parser
- Property-based testing for runtime

### Integration Tests
- End-to-end: source → executable → run
- Standard library tests
- GC stress tests

### Benchmarks
- Fibonacci (recursion)
- Binary trees (GC stress)
- N-body (performance)

### Conformance Tests
- Language feature tests
- Edge case tests
- Error handling tests

## Success Criteria

### v1.0.0 Requirements
- ✅ Lexer, Parser, AST, Sema working
- ⏳ Codegen generates working LLVM IR
- ⏳ Runtime executes IR correctly
- ⏳ Can build standalone executables
- ⏳ GC pauses <5ms (99th percentile)
- ⏳ Standard library complete
- ⏳ Kodawrap generates bindings
- ⏳ Cross-platform distribution
- ⏳ Comprehensive documentation

### Performance Targets
- Compile time: <1s for hello world
- Runtime overhead: <10% vs C
- GC pause: <5ms (99th percentile)
- Binary size: ~2-3MB for typical program
