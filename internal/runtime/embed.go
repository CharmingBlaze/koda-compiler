package kodaruntime

import _ "embed"

// This package embeds a large alternate C runtime (data/koda.c, data/koda.h).
//
// The default `koda build` / [internal/nativebuild.Build] pipeline does **not** link these
// files today: it links the C runtime built from [runtime/src] into
// runtime/libkoda_runtime.a. Keep LLVM `declare` names in [internal/codegen/runtime.go]
// aligned with that tree. Treat this embed as experimental / reserved unless a future
// backend switches to shipping monolithic koda.c again.

//go:embed data/koda.c
var KodaC []byte

//go:embed data/koda.h
var KodaH []byte
