package configure

import (
	"go/ast"
	"go/types"
	"strings"

	"github.com/dogmatiq/dogmavet/internal/dogmatypes"
	"github.com/dogmatiq/dogmavet/internal/report"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer reports issues with configuration of Dogma applications and handlers.
var Analyzer = &analysis.Analyzer{
	Name: "dogma_configure",
	Doc:  "check for misconfigured applications and handlers",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		ctrlflow.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	in := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	if report.Ignore(pass) {
		return nil, nil
	}

	// Look for methods that have a signature like:
	//
	// `func (T) Configure(dogma.[Something]Configurer)`
	in.Preorder(
		[]ast.Node{
			(*ast.FuncDecl)(nil),
		},
		func(x ast.Node) {
			decl := x.(*ast.FuncDecl)

			if decl.Recv == nil {
				// This is a function, not a method.
				return
			}

			if decl.Name.Name != "Configure" {
				// This method is not named Configure()
				return
			}

			if decl.Type.Params.NumFields() != 1 {
				// This does not accept exactly one parameter.
				return
			}

			param := decl.Type.Params.List[0]

			nt, ok := pass.TypesInfo.TypeOf(param.Type).(*types.Named)
			if !ok {
				// The parameter does not have a named type (can't be dogma.Something).
				return
			}

			if nt.Obj().Pkg().Path() != dogmatypes.PkgPath {
				// The parameter type is not in the Dogma package.
				return
			}

			if !strings.HasSuffix(nt.Obj().Name(), "Configurer") {
				// The parameter type is not one of the Dogma configurer types.
				return
			}

			if !checkParam(pass, decl, param) {
				return
			}

			checkIdentity(pass, decl, param.Names[0])
		},
	)

	return nil, nil
}
