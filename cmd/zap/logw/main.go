package main

import (
	"go-custom-linter/pkg/passes/zap/logw"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(logw.Analyzer)
}
