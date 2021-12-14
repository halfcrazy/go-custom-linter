package main

import (
	"go-custom-linter/pkg/passes/gorm/addressable"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(addressable.Analyzer)
}
