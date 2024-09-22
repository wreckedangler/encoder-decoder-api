package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveAllExtensions(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"document.pdf", "document"},
		{"archive.tar.gz", "archive"},
		{"my.file.name.txt", "my"},
		{"README", "README"},
		{"noextension", "noextension"},
		{"multi.part.name.docx", "multi"},
		{"complex.name.with.many.dots.tar.gz", "complex"},
		{"file.with.multiple.dots.", "file"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := removeAllExtensions(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
