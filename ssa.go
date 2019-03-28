package analysisutil

import (
	"golang.org/x/tools/go/ssa"
)

// IfInstr returns *ssa.If which is contained in the block b.
// If the block b has not any if instruction, IfInstr returns nil.
func IfInstr(b *ssa.BasicBlock) *ssa.If {
	if len(b.Instrs) == 0 {
		return nil
	}

	ifinstr, ok := b.Instrs[len(b.Instrs)-1].(*ssa.If)
	if !ok {
		return nil
	}

	return ifinstr
}

// Phi returns phi values which are contained in the block b.
func Phi(b *ssa.BasicBlock) (phis []*ssa.Phi) {
	for _, instr := range b.Instrs {
		if phi, ok := instr.(*ssa.Phi); ok {
			phis = append(phis, phi)
		} else {
			// no more phi
			break
		}
	}
	return
}

// Returns returns a slice of *ssa.Return in the function.
func Returns(v ssa.Value) []*ssa.Return {
	fn, ok := v.(*ssa.Function)
	if !ok {
		return nil
	}

	var rets []*ssa.Return
	done := map[*ssa.BasicBlock]bool{}
	for _, b := range fn.Blocks {
		if done[b] && len(b.Instrs) != 1 {
			continue
		}
		switch instr := b.Instrs[0].(type) {
		case *ssa.Return:
			rets = append(rets, instr)
		}
	}

	return rets
}
