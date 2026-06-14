package kodahome

import "os"

// InstallDirWritable reports whether the install directory allows creating a
// temporary file (needed for embedded toolchain extraction and writable stdlib trees).
func InstallDirWritable(dir string) (ok bool, detail string) {
	if dir == "" {
		return false, "install directory is empty"
	}
	f, err := os.CreateTemp(dir, ".koda_doctor_write_*")
	if err != nil {
		return false, err.Error()
	}
	_ = f.Close()
	_ = os.Remove(f.Name())
	return true, ""
}
