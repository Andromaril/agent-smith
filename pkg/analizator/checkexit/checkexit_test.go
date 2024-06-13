// Package checkexit содержит анализатор, запрещающий использовать прямой вызов os.Exit в функции main пакета main

package checkexit

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOsExitAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), OsExitAnalyzer, "./pkg1")
}
