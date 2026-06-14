# Build script for Koda using the new codegen and runtime

param(
    [Parameter(Mandatory=$true)]
    [string]$KodaFile,
    [string]$Output = "output.exe"
)

Write-Host "Building $KodaFile..."

# Parse and generate LLVM IR
go run cmd/koda/main.go check $KodaFile
if ($LASTEXITCODE -ne 0) {
    exit 1
}

# Generate LLVM IR
go run cmd/koda/main.go disasm $KodaFile | Out-File -Encoding ASCII output.ll
if ($LASTEXITCODE -ne 0) {
    exit 1
}

# Compile LLVM IR to object
llc -filetype=obj output.ll -o output.o
if ($LASTEXITCODE -ne 0) {
    exit 1
}

# Link with new runtime library
gcc -static -O3 -s output.o runtime/libkoda_runtime.a -lm -o $Output
if ($LASTEXITCODE -ne 0) {
    exit 1
}

# Clean up
Remove-Item -Force output.ll, output.o -ErrorAction SilentlyContinue

Write-Host "Built: $Output"
