package logw_test

import (
	"go-custom-linter/pkg/passes/zap/logw"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()

	tests := []string{"a"}
	analysistest.Run(t, testdata, logw.Analyzer, tests...)
}
