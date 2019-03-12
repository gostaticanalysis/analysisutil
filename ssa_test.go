package analysisutil_test

import (
	"reflect"
	"testing"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/ssa"
)

func Test_IfInstr(t *testing.T) {
	type I = ssa.Instruction
	var ifinst ssa.If
	tests := []struct {
		name   string
		instrs []ssa.Instruction
		hasIf  bool
	}{
		{"empty", []I{}, false},
		{"hasIf", []I{&ifinst}, true},
		{"hasIf and others", []I{nil, nil, &ifinst}, true},
		{"not tail", []I{&ifinst, nil}, false}, // actually not occur
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var b ssa.BasicBlock
			b.Instrs = tt.instrs
			got := analysisutil.IfInstr(&b)
			switch {
			case got == nil && tt.hasIf:
				t.Errorf("want *ssa.If but got nil")
			case got != nil && !tt.hasIf:
				t.Errorf("want nil but got *ssa.If")
			}
		})
	}
}

func Test_Phi(t *testing.T) {
	type I = ssa.Instruction
	var phi1, phi2 ssa.Phi
	tests := []struct {
		name   string
		instrs []ssa.Instruction
		want   []*ssa.Phi
	}{
		{"empty", []I{}, nil},
		{"P", []I{&phi1}, []*ssa.Phi{&phi1}},
		{"P_", []I{&phi1, nil}, []*ssa.Phi{&phi1}},
		{"_P", []I{nil, &phi1}, nil}, // actually not occur
		{"PP", []I{&phi1, &phi2}, []*ssa.Phi{&phi1, &phi2}},
		{"P_P", []I{&phi1, nil, &phi2}, []*ssa.Phi{&phi1}}, // actually not occur
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var b ssa.BasicBlock
			b.Instrs = tt.instrs
			got := analysisutil.Phi(&b)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want %#v but got %#v", tt.want, got)
			}
		})
	}
}
