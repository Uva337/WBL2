package main

import "testing"

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true},
		{"", "", false},
		{`qwe\4\5`, "qwe45", false},
		{`qwe\45`, "qwe44444", false},
		{`qwe\\5`, `qwe\\\\\`, false}, // Экранированный бэкслеш повторяется 5 раз
		{`abc\`, "", true},            // Ошибка: висящий бэкслеш
	}

	for _, tc := range tests {
		res, err := Unpack(tc.input)
		if (err != nil) != tc.wantErr {
			t.Errorf("Unpack(%q) error = %v, wantErr %v", tc.input, err, tc.wantErr)
			continue
		}
		if res != tc.expected {
			t.Errorf("Unpack(%q) = %q, want %q", tc.input, res, tc.expected)
		}
	}
}
