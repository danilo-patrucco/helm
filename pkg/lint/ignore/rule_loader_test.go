package ignore

//
//import (
//	"fmt"
//	"github.com/stretchr/testify/assert"
//	"path/filepath"
//	"testing"
//)
//
//func TestNewIgnorer(t *testing.T) {
//	chartPath := "../rules/testdata/withsubchartlintignore"
//	ignoreFilePath := filepath.Join(chartPath, ".helmlintignore")
//	ignorer, err := NewRuleLoader(chartPath, ignoreFilePath, func(format string, args ...interface{}) {
//		t.Logf(format, args...)
//	})
//	assert.NoError(t, err)
//	assert.NotNil(t, ignorer, "RuleLoader should not be nil")
//}
