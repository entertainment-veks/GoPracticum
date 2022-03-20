package exit_analyzer

import (
	"errors"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var errUnsafePackage = errors.New("main contains os.Exit(), but must not")

var Analyzer = &analysis.Analyzer{
	Name:     "exit_check",
	Doc:      "check exit func in main files",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      exitCheck,
}

func exitCheck(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		for _, d := range f.Decls {
			switch x := d.(type) {
			case *ast.FuncDecl:
				if x.Name.Name != "main" {
					return nil, nil
				}
				for _, s := range x.Body.List {
					switch y := s.(type) {
					case *ast.ExprStmt:
						switch z := y.X.(type) {
						case *ast.CallExpr:
							switch c := z.Fun.(type) {
							case *ast.SelectorExpr:
								isFind := false
								switch b := c.X.(type) {
								case *ast.Ident:
									isFind = b.Name == "os"
								}
								if isFind {
									isFind = c.Sel.Name == "Exit"
								}

								if isFind {
									return c, errUnsafePackage
								}
							}
						}
					}
				}
			}
		}
	}

	return nil, nil
}
