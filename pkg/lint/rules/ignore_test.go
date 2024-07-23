package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMissingIgnoreFile(t *testing.T) {
	t.Run("missing file should return an error", func(t *testing.T) {
		filePath := "made-up-file-path-here"
		_, err := ParseIgnoreFile(filePath)
		assert.Error(t, err)
	})
}
