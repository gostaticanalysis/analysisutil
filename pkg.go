package analysisutil

import "strings"

// RemoVendor removes vendoring infomation from import path.
func RemoveVendor(path string) string {
	i := strings.Index(path, "vendor")
	if i >= 0 {
		return path[i+len("vendor")+1:]
	}
	return path
}
