package analysisutil_test

import (
	"testing"

	"github.com/gostaticanalysis/analysisutil"
)

func TestRemoveVendor(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"vendor/path/to/pkg", "path/to/pkg"},
		{"path/to/pkg", "path/to/pkg"},
		{"", ""},
		{"a/vendor/path/to/pkg", "path/to/pkg"},
		{"pkg", "pkg"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()
			got := analysisutil.RemoveVendor(tt.path)
			if got != tt.want {
				t.Errorf("want %s but got %s", tt.want, got)
			}
		})
	}
}
