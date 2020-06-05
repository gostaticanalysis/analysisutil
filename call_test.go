package analysisutil_test

import (
	"errors"
	"testing"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/buildssa"
)

var callAnalyzer = &analysis.Analyzer{
	Name: "test_call",
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, callAnalyzer, "call")
}

func run(pass *analysis.Pass) (interface{}, error) {
	resTyp := analysisutil.TypeOf(pass, "call", "*res")
	if resTyp == nil {
		return nil, errors.New("analyzer does not find *call.res type")
	}

	close := analysisutil.MethodOf(resTyp, "close")
	if close == nil {
		return nil, errors.New("analyzer does not find (call.res).close method")
	}

	funcs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs
	for _, f := range funcs {
		for _, b := range f.Blocks {
			for i, instr := range b.Instrs {
				called, ok := analysisutil.CalledFrom(b, i, resTyp, close)
				if ok && !called {
					pass.Reportf(instr.Pos(), "NG")
				}
			}
		}
	}

	return nil, nil
}
