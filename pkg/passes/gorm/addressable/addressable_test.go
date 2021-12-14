package addressable_test

import (
	"go-custom-linter/pkg/passes/gorm/addressable"
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()

	tests := []string{"a"}
	analysistest.Run(t, testdata, addressable.Analyzer, tests...)
}
