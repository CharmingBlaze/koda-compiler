//go:build !release

package kodahome

func embeddedToolchain() (*Toolchain, error) {
	return nil, nil
}
