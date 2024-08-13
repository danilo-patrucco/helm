package ignore

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestNewIgnorer(t *testing.T) {
	chartPath := "../rules/testdata/withsubchartlintignore"
	ignoreFilePath := filepath.Join(chartPath, ".helmlintignore")
	ignorer, err := NewCmdIgnorer(chartPath, ignoreFilePath, func(format string, args ...interface{}) {
		t.Logf(format, args...)
	})
	assert.NoError(t, err)
	assert.NotNil(t, ignorer, "CmdIgnorer should not be nil")
	assert.NotEmpty(t, ignorer.Patterns, "Expected patterns to be loaded from the file, but none were found")
	if len(ignorer.Patterns) == 0 {
		t.Errorf("Expected patterns to be loaded from the file, but none were found")
	}
}

func TestDebug(t *testing.T) {
	var captured string
	debugFn := func(format string, args ...interface{}) {
		captured = fmt.Sprintf(format, args...)
	}
	ignorer := &CmdIgnorer{
		debugFnOverride: debugFn,
	}
	ignorer.Debug("test %s", "debug")
	assert.Equal(t, "test debug", captured)
}