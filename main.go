package main

import (
	"contrib-linter/passes/zap/logw"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(logw.Analyzer)
}
