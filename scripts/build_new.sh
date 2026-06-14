#!/bin/bash
# Build script for Koda using the new codegen and runtime

KODA_FILE=$1
OUTPUT=${2:-"output"}

if [ -z "$KODA_FILE" ]; then
    echo "Usage: $0 <file.koda> [output]"
    exit 1
fi

echo "Building $KODA_FILE..."

# Parse and generate LLVM IR
go run cmd/koda/main.go check "$KODA_FILE" || exit 1

# Generate LLVM IR
go run cmd/koda/main.go disasm "$KODA_FILE" > output.ll || exit 1

# Compile LLVM IR to object
llc -filetype=obj output.ll -o output.o || exit 1

# Link with new runtime library
gcc -static -O3 -s output.o runtime/libkoda_runtime.a -lm -o "$OUTPUT" || exit 1

# Clean up
rm -f output.ll output.o

echo "Built: $OUTPUT"
