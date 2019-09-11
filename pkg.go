package analysisutil

import (
	"go/types"
	"strings"
	"path/filepath"
)

// RemoVendor removes vendoring infomation from import path.
func RemoveVendor(path string) string {
	unixVendorPath := "vendor/"
	i := strings.Index(path, filepath.FromSlash(unixVendorPath))
	if i >= 0 {
		return path[i+len(unixVendorPath):]
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
