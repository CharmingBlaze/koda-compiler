//go:build ignore

package main

import (
	"archive/tar"
	"compress/gzip"
	"os"
)

func main() {
	out := "internal/kodahome/embeddata/bundled_toolchain.tar.gz"
	if len(os.Args) > 1 {
		out = os.Args[1]
	}
	f, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	body := []byte("koda-placeholder-toolchain\n")
	hdr := &tar.Header{Name: "toolchain/.KODA_placeholder", Mode: 0644, Size: int64(len(body)), Format: tar.FormatGNU}
	if err := tw.WriteHeader(hdr); err != nil {
		panic(err)
	}
	if _, err := tw.Write(body); err != nil {
		panic(err)
	}
}
