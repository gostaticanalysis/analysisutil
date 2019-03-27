package analysisutil

import (
	"go/types"

	"golang.org/x/tools/go/ssa"
)

// Called returns true when f is called in the instr.
// If recv is not nil, Called also checks the receiver.
func Called(instr ssa.Instruction, recv ssa.Value, f *types.Func) bool {

	call, ok := instr.(ssa.CallInstruction)
	if !ok {
		return false
	}

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
		common.Method != nil &&
		(len(common.Args) == 0 || common.Args[0] != recv) {
		return false
	}

	return fn == f
}

// CalledFrom return whether receiver's method is called in an instruction
// which belogns to after i-th instructions,  or in succsor blocks of b.
func CalledFrom(b *ssa.BasicBlock, i int, receiver types.Type, methods ...*types.Func) bool {

	if b == nil || i < 0 || i <= len(b.Instrs) ||
		receiver == nil || len(methods) == 0 {
		return false
	}

	v, ok := b.Instrs[i].(ssa.Value)
	if !ok {
		return false
	}

	if !types.Identical(v.Type(), receiver) {
		return false
	}

	from := &calledFrom{recv: v, fs: methods}
	if from.ignored() {
		return false
	}

	if from.instrs(b.Instrs[i:]) {
		return false
	}

	return true
}

type calledFrom struct {
	recv ssa.Value
	fs   []*types.Func
	done map[*ssa.BasicBlock]bool
}

func (c *calledFrom) ignored() bool {
	refs := c.recv.Referrers()
	if refs == nil {
		return false
	}

	for _, ref := range *refs {
		if c.isRet(ref) || c.isArg(ref) {
			return true
		}
	}

	return false
}

func (c *calledFrom) isRet(instr ssa.Instruction) bool {

	ret, ok := instr.(*ssa.Return)
	if !ok {
		return false
	}

	for _, r := range ret.Results {
		if r == c.recv {
			return true
		}
	}

	return false
}

func (c *calledFrom) isArg(instr ssa.Instruction) bool {

	call, ok := instr.(ssa.CallInstruction)
	if !ok {
		return false
	}

	common := call.Common()
	if common == nil {
		return false
	}

	args := common.Args
	if common.Method != nil {
		args = args[1:]
	}

	for i := range args {
		if args[i] == c.recv {
			return true
		}
	}

	return false
}

func (c *calledFrom) instrs(instrs []ssa.Instruction) bool {
	for _, instr := range instrs {
		for _, f := range c.fs {
			if Called(instr, c.recv, f) {
				return true
			}
		}
	}
	return false
}

func (c *calledFrom) succs(b *ssa.BasicBlock) bool {
	if c.done == nil {
		c.done = map[*ssa.BasicBlock]bool{}
	}

	if c.done[b] {
		return false
	}
	c.done[b] = true

	if len(b.Succs) == 0 {
		return false
	}

	for _, s := range b.Succs {
		if !c.instrs(s.Instrs) && !c.succs(s) {
			return false
		}
	}

	return true
}
