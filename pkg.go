package analysisutil

import (
	"go/types"
	"strconv"
	"strings"

	"github.com/Matts966/refsafe/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ssa"
)

// RemoVendor removes vendoring infomation from import path.
func RemoveVendor(path string) string {
	i := strings.Index(path, "vendor")
	if i >= 0 {
		return path[i+len("vendor")+1:]
	}
	return path
}

// LookupFromImports finds an object from import paths.
func LookupFromImports(imports []*types.Package, path, name string) types.Object {
	path = RemoveVendor(path)
	for i := range imports {
		if path == RemoveVendor(imports[i].Path()) {
			return imports[i].Scope().Lookup(name)
		}
	}
	return nil
}

// PkgUsedInFunc returns true when the given f imports the pkg.
func PkgUsedInFunc(pass *analysis.Pass, pkgPath string, f *ssa.Function) bool {
	if f == nil {
		return false
	}
	fo := f.Object()
	if fo == nil {
		return false
	}

	ff := analysisutil.File(pass, fo.Pos())
	if ff == nil {
		return false
	}
	for _, i := range ff.Imports {
		path, err := strconv.Unquote(i.Path.Value)
		if err != nil {
			continue
		}
		if analysisutil.RemoveVendor(path) == pkgPath {
			return true
		}
	}
	return false
}

// PkgUsedInFunc returns true when the given f imports the pkg.
func PkgUsedIn(pass *analysis.Pass, pkgPath string, f *ssa.Function) bool {
	if f == nil {
		return false
	}
	fo := f.Object()
	if fo == nil {
		return false
	}

	ff := analysisutil.File(pass, fo.Pos())
	if ff == nil {
		return false
	}
	for _, i := range ff.Imports {
		path, err := strconv.Unquote(i.Path.Value)
		if err != nil {
			continue
		}
		if analysisutil.RemoveVendor(path) == pkgPath {
			return true
		}
	}
	return false
}
