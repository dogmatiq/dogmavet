package configure

import (
	"go/ast"
	"go/types"
	"strings"

	"github.com/dogmatiq/dogmavet/internal/dogmatypes"
	"github.com/dogmatiq/dogmavet/internal/testhelper"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ssa"
)

// Analyzer reports issues with configuration of Dogma applications and handlers.
var Analyzer = &analysis.Analyzer{
	Name: "dogma_configure",
	Doc:  "check for misconfigured applications and handlers",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	if testhelper.Ignore(pass) {
		return nil, nil
	}

	var pkg *ssa.Package

	for _, f := range pass.Files {
		for _, d := range f.Decls {
			decl, ok := d.(*ast.FuncDecl)
			if !ok {
				continue
			}

			if !isConfigureMethod(pass, decl) {
				continue
			}

			if !checkParam(pass, decl) {
				continue
			}

			if pkg == nil {
				pkg = buildSSA(pass)
			}

			checkIdentity(pass, decl, pkg)
		}
	}

	return nil, nil
}

// buildSSA builds the SSA representation of the package under anaylsis.
func buildSSA(pass *analysis.Pass) *ssa.Package {
	var (
		prog    = ssa.NewProgram(pass.Fset, 0)
		visited = make(map[*types.Package]struct{})
		imports func(pkgs *types.Package)
	)

	imports = func(pkg *types.Package) {
		for _, p := range pkg.Imports() {
			if _, ok := visited[p]; ok {
				return
			}
			visited[p] = struct{}{}
			prog.CreatePackage(p, nil, nil, true)
			imports(p)
		}
	}

	imports(pass.Pkg)

	pkg := prog.CreatePackage(pass.Pkg, pass.Files, pass.TypesInfo, false)
	pkg.SetDebugMode(true)
	pkg.Build()

	return pkg
}

// isConfigureMethod returns true if fn is a method with a signature like:
//
//	func (T) Configure(dogma.[XXX]Configurer)
func isConfigureMethod(
	pass *analysis.Pass,
	decl *ast.FuncDecl,
) bool {
	if decl.Recv == nil {
		// This is a function, not a method.
		return false
	}

	if decl.Name.Name != "Configure" {
		// This method is not named Configure().
		return false
	}

	if decl.Type.Params.NumFields() != 1 {
		// This function does not accept exactly one parameter (not
		// including the receiver).
		return false
	}

	param := decl.Type.Params.List[0]

	nt, ok := pass.TypesInfo.TypeOf(param.Type).(*types.Named)
	if !ok {
		// The parameter does not have a named type (can't be dogma.Something).
		return false
	}

	if nt.Obj().Pkg().Path() != dogmatypes.PkgPath {
		// The parameter type is not in the Dogma package.
		return false
	}

	if !strings.HasSuffix(nt.Obj().Name(), "Configurer") {
		// The parameter type is not one of the Dogma configurer types.
		return false
	}

	return true
}
