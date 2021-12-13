package logw_test

import (
	"testing"

	"contrib-linter/passes/zap/logw"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()

	tests := []string{"a"}
	analysistest.Run(t, testdata, logw.Analyzer, tests...)
}
