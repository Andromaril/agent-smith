// Package checkexit содержит анализатор, запрещающий использовать прямой вызов os.Exit в функции main пакета main
package checkexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OsExitAnalyzer анализатор
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for usage os.Exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, v := range pass.Files {
		if v.Name.Name != "main" {
			return nil, nil
		}
		ast.Inspect(v, func(n ast.Node) bool {
			// проверяем, какой конкретный тип лежит в узле
			if f, ok := n.(*ast.FuncDecl); ok && f.Name.Name != "main" {
				return false
			}
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			f, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if i, ok := f.X.(*ast.Ident); ok && i.Name == "os" && f.Sel.Name == "Exit" {
				pass.Reportf(n.(*ast.CallExpr).Pos(), "calling os.Exit")
			}
			return true

		})
	}
	return nil, nil
}
