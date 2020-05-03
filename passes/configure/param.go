package configure

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/dogmatiq/dogmavet/internal/report"
	"golang.org/x/tools/go/analysis"
)

const expectedParamName = "c"

func checkParam(
	pass *analysis.Pass,
	decl *ast.FuncDecl,
	param *ast.Field,
) bool {
	if len(param.Names) != 1 {
		report.AtWithFix(
			pass,
			param,
			analysis.SuggestedFix{
				Message: "add a parameter name",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     param.Type.Pos(),
						End:     param.Type.Pos(),
						NewText: []byte(expectedParamName + " "),
					},
				},
			},
			`configurer parameter should be named '%s'`,
			expectedParamName,
		)

		return false
	}

	ident := param.Names[0]
	if ident.Name == expectedParamName {
		return true
	}

	fix := analysis.SuggestedFix{
		Message: fmt.Sprintf("rename parameter to '%s'", expectedParamName),
		TextEdits: []analysis.TextEdit{
			{
				Pos:     ident.NamePos,
				End:     ident.NamePos + token.Pos(len(ident.Name)),
				NewText: []byte(expectedParamName),
			},
		},
	}

	ast.Inspect(decl.Body, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Ident:
			if n.Obj == ident.Obj {
				fix.TextEdits = append(
					fix.TextEdits,
					analysis.TextEdit{
						Pos:     n.NamePos,
						End:     n.NamePos + token.Pos(len(n.Name)),
						NewText: []byte(expectedParamName),
					},
				)
			}
		}

		return true
	})

	report.AtWithFix(
		pass,
		param,
		fix,
		`configurer parameter should be named '%s'`,
		expectedParamName,
	)

	return true
}
