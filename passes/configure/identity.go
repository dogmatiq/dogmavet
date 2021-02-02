package configure

import (
	"go/ast"
	"go/constant"
	"go/types"
	"strconv"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/dogmavet/internal/report"
	"github.com/google/uuid"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ssa"
)

// checkIdentity checks that a Configure() method uses the configurer's
// Identity() method correctly.
func checkIdentity(
	pass *analysis.Pass,
	decl *ast.FuncDecl,
	pkg *ssa.Package,
) {
	obj := pass.TypesInfo.Defs[decl.Name].(*types.Func)
	fn := pkg.Prog.FuncValue(obj)
	param := fn.Params[1]

	visited := map[*ssa.BasicBlock]bool{}
	called := checkIdentityInBlock(
		pass,
		decl,
		param,
		fn.Blocks[0],
		visited,
		0,
	)

	if !called {
		report.At(
			pass,
			decl,
			"Configure() must call %s.Identity() exactly once",
			param.Name(),
		)
	}
}

// checkIdentityInBlock looks for a call to the Identity() function in a
// specific block, and its successors.
//
// It returns true if it is called at all, even if it is is not called on all
// execution paths.
func checkIdentityInBlock(
	pass *analysis.Pass,
	decl *ast.FuncDecl,
	param *ssa.Parameter,
	block *ssa.BasicBlock,
	visited map[*ssa.BasicBlock]bool,
	priorCalls int,
) (called bool) {
	if called, ok := visited[block]; ok {
		return called
	}
	defer func() {
		visited[block] = called
	}()

	var calls []*ast.CallExpr

	for _, i := range block.Instrs {
		if call, ok := i.(*ssa.Call); ok {
			if isConfigurerCall(param, call, "Identity") {
				called = true
				refs := *call.Referrers()
				expr := refs[0].(*ssa.DebugRef).Expr.(*ast.CallExpr)
				calls = append(calls, expr)
				checkIdentityCall(pass, expr)
			}
		}
	}

	if len(calls) > 0 {
		// If there is more than one call in this specific block, all but the
		// first is a duplicate.
		for _, call := range calls[1:] {
			report.At(
				pass,
				call,
				"%s.Identity() must be called exactly once",
				param.Name(),
			)
		}

		if priorCalls > 0 {
			// If there are any prior calls at all (from parent blocks), then
			// all of the calls in this block are duplicates.
			report.At(
				pass,
				calls[0],
				"%s.Identity() has already been called on at least one execution path",
				param.Name(),
			)
		}
	}

	if len(block.Succs) == 0 {
		return called
	}

	if len(block.Succs) == 1 {
		nextCalled := checkIdentityInBlock(
			pass,
			decl,
			param,
			block.Succs[0],
			visited,
			priorCalls+len(calls),
		)

		return called || nextCalled
	}

	thenBlock := block.Succs[0]
	thenCalled := checkIdentityInBlock(
		pass,
		decl,
		param,
		thenBlock,
		visited,
		priorCalls+len(calls),
	)

	elseBlock := block.Succs[1]
	elseCalled := checkIdentityInBlock(
		pass,
		decl,
		param,
		elseBlock,
		visited,
		priorCalls+len(calls),
	)

	if thenCalled != elseCalled {
		cond := block.Instrs[len(block.Instrs)-1].(*ssa.If).Cond
		refs := *cond.Referrers()
		expr := refs[0].(*ssa.DebugRef).Expr

		report.At(
			pass,
			expr,
			"this control-flow statement causes %s.Identity() to remain uncalled on some execution paths",
			param.Name(),
		)
	}

	return called || thenCalled || elseCalled
}

// isConfigureCall returns true if the given call expression is a call to a
// specific method on a Dogma configurer.
func isConfigurerCall(
	param *ssa.Parameter,
	call *ssa.Call,
	method string,
) bool {
	com := call.Common()

	if com.Value != param {
		return false
	}

	return com.Method.Name() == method
}

func checkIdentityCall(
	pass *analysis.Pass,
	call *ast.CallExpr,
) {
	checkIdentityName(pass, call.Args[0])
	checkIdentityKey(pass, call.Args[1])
}

func checkIdentityName(
	pass *analysis.Pass,
	expr ast.Expr,
) {
	v := pass.TypesInfo.Types[expr].Value
	if v == nil {
		return
	}

	name := constant.StringVal(v)

	if _, err := configkit.NewIdentity(name, "<placeholder>"); err != nil {
		report.At(
			pass,
			expr,
			"%s",
			err.Error(),
		)
	}
}

func checkIdentityKey(
	pass *analysis.Pass,
	expr ast.Expr,
) {
	v := pass.TypesInfo.Types[expr].Value
	if v == nil {
		return
	}

	key := constant.StringVal(v)

	if _, err := configkit.NewIdentity("<placeholder>", key); err != nil {
		report.At(
			pass,
			expr,
			"%s",
			err.Error(),
		)

		return
	}

	id, err := uuid.Parse(key)
	if err != nil {
		report.AtWithFix(
			pass,
			expr,
			analysis.SuggestedFix{
				Message: "generate a new UUID to use as the identity key",
				TextEdits: []analysis.TextEdit{
					{
						Pos: expr.Pos(),
						End: expr.End(),
						NewText: []byte(strconv.Quote(
							uuid.NewString(),
						)),
					},
				},
			},
			"identity keys should be UUIDs (%s)",
			err.Error(),
		)

		return
	}

	if len(key) != 36 {
		report.AtWithFix(
			pass,
			expr,
			analysis.SuggestedFix{
				Message: "reformat the UUID to use the standard format",
				TextEdits: []analysis.TextEdit{
					{
						Pos: expr.Pos(),
						End: expr.End(),
						NewText: []byte(strconv.Quote(
							id.String(),
						)),
					},
				},
			},
			"identity key UUIDs should use RFC-4122 hex notation (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)",
		)
	}
}
