package cloner

import (
	"testing"
)

func Test_normaliseURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "bare path gets https prefix",
			input: "github.com/nxlak/test",
			want:  "https://github.com/nxlak/test",
		},
		{
			name:  "https URL is unchanged",
			input: "https://github.com/stretchr/testify",
			want:  "https://github.com/stretchr/testify",
		},
		{
			name:  "http URL is unchanged",
			input: "http://github.com/sirupsen/logrus",
			want:  "http://github.com/sirupsen/logrus",
		},
		{
			name:  "ssh git@ URL is unchanged",
			input: "git@github.com:golang-jwt/jwt.git",
			want:  "git@github.com:golang-jwt/jwt.git",
		},
		{
			name:  "ssh:// URL is unchanged",
			input: "ssh://git@github.com/gorilla/mux.git",
			want:  "ssh://git@github.com/gorilla/mux.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeAndValidateURL(tt.input)
			if err != nil {
				t.Fatalf("normalizeAndValidateURL(%q): unexpected error: %v", tt.input, err)
			}

			if got != tt.want {
				t.Errorf("normaliseURL(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}
