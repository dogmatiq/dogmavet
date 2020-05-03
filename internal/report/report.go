package report

import (
	"fmt"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// At reports a diagnostic at the specified range.
func At(
	pass *analysis.Pass,
	r analysis.Range,
	f string,
	v ...interface{},
) {
	report(
		pass,
		r.Pos(), r.End(),
		nil,
		f, v,
	)
}

// AtWithFix reports a diagnostic at the specified range and includes a
// suggested fix.
func AtWithFix(
	pass *analysis.Pass,
	r analysis.Range,
	fix analysis.SuggestedFix,
	f string,
	v ...interface{},
) {
	report(
		pass,
		r.Pos(), r.End(),
		&fix,
		f, v,
	)
}

// AtLineOf reports a diagnostic at the start of the line of the given range.
func AtLineOf(
	pass *analysis.Pass,
	r analysis.Range,
	f string,
	v ...interface{},
) {
	column := pass.Fset.Position(r.Pos()).Column
	pos := r.Pos() - token.Pos(column) + 1

	report(
		pass,
		pos,
		r.End(),
		nil,
		f, v,
	)
}

func report(
	pass *analysis.Pass,
	pos, end token.Pos,
	fix *analysis.SuggestedFix,
	f string,
	v []interface{},
) {
	d := analysis.Diagnostic{
		Pos:      pos,
		End:      end,
		Category: "dogma",
		Message:  "dogma: " + fmt.Sprintf(f, v...),
	}

	if fix != nil {
		d.SuggestedFixes = []analysis.SuggestedFix{
			*fix,
		}
	}

	pass.Report(d)
}
