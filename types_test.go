package analysisutil_test

import (
	"go/token"
	"log"
	"testing"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/buildssa"
)

const pkg = "objectof"

var analyzer = &analysis.Analyzer{
	Name: "test_types",
	Run:  run_test_types,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

func Test_Types(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer, pkg)
}

func run_test_types(pass *analysis.Pass) (interface{}, error) {
	tests := []struct {
		path, name string
		found      bool
	}{
		{"fmt", "Println", true},
		{pkg, "A", true},
		{pkg, "EOF", true},
		{"io", "EOF", true},
		{"reflect", "Kind", false},
		{"a", "ok", false},
		{"vendored", "EOF", true},
		{"c", "EOF", false},
	}

	log.Println(analysisutil.LookupFromImports(pass.Pkg.Imports(), "vendored", "EOF"))

	for _, tt := range tests {
		tt := tt
		obj := analysisutil.ObjectOf(pass, tt.path, tt.name)

		if obj == nil && tt.found {
			pass.Reportf(token.NoPos, "objectof could not find %s.%s", tt.path, tt.name)
		}
		if obj != nil && !tt.found {
			pass.Reportf(token.NoPos, "objectof found %s.%s, which does not exist", tt.path, tt.name)
		}
	}

	return nil, nil
}
