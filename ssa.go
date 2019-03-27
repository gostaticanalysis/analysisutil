package analysisutil

import (
	"go/types"

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

// Called returns true when f is called in the call.
// If recv is not nil, Called also checks the receiver.
func Called(call ssa.CallInstruction, recv ssa.Value, f *types.Func) bool {
	common := call.Common()
	if common == nil {
		return false
	}

	callee := common.StaticCallee()
	if callee == nil {
		return false
	}

	fn, ok := callee.Object().(*types.Func)
	if !ok {
		return false
	}

	if recv != nil &&
		(len(common.Args) == 0 || common.Args[0] != recv) {
		return false
	}

	return fn == f
}
