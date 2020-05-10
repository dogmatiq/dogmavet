package configure

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/dogmatiq/dogmavet/internal/report"
	"golang.org/x/tools/go/analysis"
)

// preferredParamName is the preferred name for the configurer parameter to a
// Configure() method.
const preferredParamName = "c"

// checkParam performs checks on the configurer parameter on a Configure()
// method.
//
// It returns true if the parameter has a name.
func checkParam(
	pass *analysis.Pass,
	decl *ast.FuncDecl,
) bool {
	param := decl.Type.Params.List[0]

	if len(param.Names) == 0 {
		report.AtWithFix(
			pass,
			param,
			analysis.SuggestedFix{
				Message: "add a parameter name",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     param.Type.Pos(),
						End:     param.Type.Pos(),
						NewText: []byte(preferredParamName + " "),
					},
				},
			},
			`configurer parameter should be named '%s'`,
			preferredParamName,
		)

		return false
	}

	ident := param.Names[0]
	if ident.Name == preferredParamName {
		return true
	}

	fix := analysis.SuggestedFix{
		Message: fmt.Sprintf("rename parameter to '%s'", preferredParamName),
		TextEdits: []analysis.TextEdit{
			{
				Pos:     ident.NamePos,
				End:     ident.NamePos + token.Pos(len(ident.Name)),
				NewText: []byte(preferredParamName),
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
						NewText: []byte(preferredParamName),
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
		preferredParamName,
	)

	return true
}
