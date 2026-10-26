package model

import "testing"

func TestNormalizePhone(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"bare digits", "380671234567", "+380671234567"},
		{"already e164", "+380671234567", "+380671234567"},
		{"formatted", "+38 (067) 123-45-67", "+380671234567"},
		{"dotted", "380.67.123.45.67", "+380671234567"},
		{"international 00 prefix", "00380671234567", "+380671234567"},
		{"surrounding spaces", "  380671234567 ", "+380671234567"},
		{"letters", "abc", ""},
		{"too short", "123", ""},
		{"too long", "3806712345678901", ""},
		{"leading zero", "0671234567", ""},
		{"plus in the middle", "380+671234567", ""},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizePhone(tt.in); got != tt.want {
				t.Errorf("NormalizePhone(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
