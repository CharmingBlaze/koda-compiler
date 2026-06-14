package kodahome

import "testing"

func TestInstallDirWritable(t *testing.T) {
	dir := t.TempDir()
	ok, detail := InstallDirWritable(dir)
	if !ok {
		t.Fatalf("expected temp dir writable, got %q", detail)
	}
	ok, detail = InstallDirWritable("")
	if ok || detail == "" {
		t.Fatalf("expected empty dir to fail, ok=%v detail=%q", ok, detail)
	}
}
