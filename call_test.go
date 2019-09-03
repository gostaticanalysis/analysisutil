package analysisutil_test

import (
	"go/types"
	"log"
	"testing"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/buildssa"
)

var (
	st                     types.Type
	close                  *types.Func
	doSomethingAndReturnSt *types.Func
)

var Analyzer = &analysis.Analyzer{
	Name: "test_call",
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "b")
}

func run(pass *analysis.Pass) (interface{}, error) {
	st = analysisutil.LookupFromImports([]*types.Package{
		pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).Pkg.Pkg,
	}, "b", "st").Type().(*types.Named)
	close = analysisutil.MethodOf(st, "b.close")
	doSomethingAndReturnSt = analysisutil.MethodOf(st, "b.doSomethingAndReturnSt")

	funcs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs
	for _, f := range funcs {
		for _, b := range f.Blocks {
			for i, instr := range b.Instrs {
				if !analysisutil.Called(instr, nil, doSomethingAndReturnSt) {
					continue
				}
				called, ok := analysisutil.CalledFrom(b, i, st, close)
				log.Println(called, ok)
				if !(called && ok) {
					pass.Reportf(instr.Pos(), "close should be called after calling doSomething")
				}
			}
		}
	}

	return nil, nil
}
