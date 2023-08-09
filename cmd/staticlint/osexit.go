package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OsExitAnalyzer checks for the usage of os.Exit() in the main function of the main package
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osExitAnalyzer",
	Doc:  "Checks for the usage of os.Exit() in the main function of the main package",
	Run:  checkOsExit,
}

func checkOsExit(pass *analysis.Pass) (interface{}, error) {
	// Iterate through all the files in the package
	for _, file := range pass.Files {
		// Check if the package name is "main"
		if file.Name.Name != "main" {
			continue
		}

		// Walk through the AST and find os.Exit() call within the main function
		ast.Inspect(file, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if ok && fn.Name.Name == "main" {
				ast.Inspect(fn.Body, func(n ast.Node) bool {
					callExpr, ok := n.(*ast.CallExpr)
					if ok {
						selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
						if ok {
							ident, ok := selExpr.X.(*ast.Ident)
							if ok && ident.Name == "os" && selExpr.Sel.Name == "Exit" {
								pass.Reportf(callExpr.Pos(), "os.Exit() found in main function")
							}
						}
					}
					return true
				})
			}
			return true
		})
	}
	return nil, nil
}
