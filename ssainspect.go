package analysisutil

import "golang.org/x/tools/go/ssa"

func InspectInstruction(b *ssa.BasicBlock, i int, f func(i int, instr ssa.Instruction) bool) {
	new(instrInspector).block(b, i, f)
}

type instrInspector struct {
	done map[*ssa.BasicBlock]bool
}

func (ins *instrInspector) block(b *ssa.BasicBlock, i int, f func(i int, instr ssa.Instruction) bool) {
	if ins.done == nil {
		ins.done = map[*ssa.BasicBlock]bool{}
	}

	if b == nil || ins.done[b] || len(b.Instrs) <= i {
		return
	}

	ins.done[b] = true
	ins.instrs(i+1, b.Instrs[i+1:], f)
	for _, s := range b.Succs {
		ins.block(s, 0, f)
	}

}

func (ins *instrInspector) instrs(offset int, instrs []ssa.Instruction, f func(i int, instr ssa.Instruction) bool) {
	for i, instr := range instrs {
		if !f(offset+i, instr) {
			break
		}
	}
}
