package colorhex

import "testing"

func TestParseHexColor(t *testing.T) {
	tests := []struct {
		in   string
		want uint32
	}{
		{"F00", 0xFF0000FF},
		{"FF0000", 0xFF0000FF},
		{"101018", 0x101018FF},
		{"FF000080", 0xFF000080},
		{"F008", 0xFF000088},
	}
	for _, tc := range tests {
		got, err := Parse(tc.in)
		if err != nil {
			t.Fatalf("Parse(%q): %v", tc.in, err)
		}
		if uint32(got) != tc.want {
			t.Fatalf("Parse(%q) = 0x%X, want 0x%X", tc.in, uint32(got), tc.want)
		}
	}
}

func TestParseHexColorInvalid(t *testing.T) {
	if _, err := Parse("GG"); err == nil {
		t.Fatal("expected error for invalid chars")
	}
	if _, err := Parse("12"); err == nil {
		t.Fatal("expected error for wrong length")
	}
}
