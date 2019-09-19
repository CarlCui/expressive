package codegen

import (
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

type ResultType int

const (
	VOID    ResultType = iota // fragment will not produce a value
	VALUE                     // fragment will produce a value
	POINTER                   // fragment will produce a pointer value
)

type FragmentType int

const (
	MODULE FragmentType = iota
	FUNCTIONS
	BLOCKS
	INSTRUCTIONS
)

// Fragment represents a fragment of instructions in asm
type Fragment interface {
	Append(fragment Fragment)
	AppendWithResultPropogation(fragment Fragment)
	GetResultType() ResultType
	GetResult() value.Value
}

type ModuleFragment struct {
	Module *ir.Module
}

func NewModuleFragment(module *ir.Module) *ModuleFragment {
	if module == nil {
		module = ir.NewModule()
	}

	return &ModuleFragment{Module: module}
}

func (moduleFragment *ModuleFragment) Append(fragment Fragment) {
	switch f := fragment.(type) {
	case *ModuleFragment:
		panic("Cannot append module to module")
	case *FunctionsFragment:
		moduleFragment.Module.Funcs = append(moduleFragment.Module.Funcs, f.Functions...)

		for _, function := range f.Functions {
			function.Parent = moduleFragment.Module
		}
	case *BlocksFragment:
		panic("Cannot append blocks to module")
	case *InstructionsFragment:
		panic("Cannot append instructions to module")
	}
}

func (moduleFragment *ModuleFragment) AppendWithResultPropogation(fragment Fragment) {
	moduleFragment.Append(fragment)
}

func (moduleFragment *ModuleFragment) GetResultType() ResultType {
	return VOID
}

func (moduleFragment *ModuleFragment) GetResult() value.Value {
	return nil
}

type FunctionsFragment struct {
	Functions       []*ir.Func
	currentFunction *ir.Func
}

func NewFunctionsFragment() *FunctionsFragment {
	functions := make([]*ir.Func, 0)

	return &FunctionsFragment{Functions: functions, currentFunction: nil}
}

func (functionsFragment *FunctionsFragment) AddFunc(function *ir.Func) {
	if function == nil {
		panic("Cannot add nil function")
	}

	functionsFragment.Functions = append(functionsFragment.Functions, function)
	functionsFragment.currentFunction = function
}

func (functionsFragment *FunctionsFragment) Append(fragment Fragment) {
	switch f := fragment.(type) {
	case *ModuleFragment:
		panic("Cannot append module to function")
	case *FunctionsFragment:
		functionsFragment.Functions = append(functionsFragment.Functions, f.Functions...)

		functionsFragment.currentFunction = f.currentFunction
	case *BlocksFragment:
		blocks := functionsFragment.currentFunction.Blocks

		// Chain new blocks to existing blocks
		numberOfBlocks := len(blocks)
		if numberOfBlocks > 0 && blocks[numberOfBlocks-1].Term == nil && len(f.Blocks) > 0 {
			blocks[numberOfBlocks-1].NewBr(f.Blocks[0])
		}

		functionsFragment.currentFunction.Blocks = append(blocks, f.Blocks...)

		for _, block := range f.Blocks {
			block.Parent = functionsFragment.currentFunction
		}
	case *InstructionsFragment:
		blockLength := len(functionsFragment.currentFunction.Blocks)
		lastBlock := functionsFragment.currentFunction.Blocks[blockLength-1]
		instructions := lastBlock.Insts

		lastBlock.Insts = append(instructions, f.Instructions...)
	}
}

func (functionsFragment *FunctionsFragment) AppendWithResultPropogation(fragment Fragment) {
	functionsFragment.Append(fragment)
}

func (functionsFragment *FunctionsFragment) GetResultType() ResultType {
	return VOID
}

func (functionsFragment *FunctionsFragment) GetResult() value.Value {
	return nil
}

type BlocksFragment struct {
	Blocks       []*ir.Block
	CurrentBlock *ir.Block
	resultType   ResultType
	resultValue  value.Value
}

func NewBlocksFragment(resultType ResultType) *BlocksFragment {
	blocks := make([]*ir.Block, 0)

	return &BlocksFragment{Blocks: blocks, CurrentBlock: nil, resultType: resultType, resultValue: nil}
}

func (blocksFragment *BlocksFragment) Append(fragment Fragment) {
	switch f := fragment.(type) {
	case *ModuleFragment:
		panic("Cannot append blocks to module")
	case *FunctionsFragment:
		panic("Cannot append blocks to functions")
	case *BlocksFragment:
		if len(f.Blocks) > 0 {
			blocksFragment.ChainBlocks(f.Blocks...)
			blocksFragment.Blocks = append(blocksFragment.Blocks, f.Blocks...)
			blocksFragment.CurrentBlock = f.CurrentBlock
		}

	case *InstructionsFragment:
		blocksFragment.CurrentBlock.Insts = append(blocksFragment.CurrentBlock.Insts, f.Instructions...)
	}
}

func (blocksFragment *BlocksFragment) AppendWithResultPropogation(fragment Fragment) {
	switch f := fragment.(type) {
	case *ModuleFragment:
		panic("Cannot append blocks to module")
	case *FunctionsFragment:
		panic("Cannot append blocks to functions")
	case *BlocksFragment:
		if len(f.Blocks) > 0 {
			blocksFragment.ChainBlocks(f.Blocks...)
			blocksFragment.Blocks = append(blocksFragment.Blocks, f.Blocks...)
			blocksFragment.CurrentBlock = f.CurrentBlock

			blocksFragment.resultType = f.resultType
			blocksFragment.resultValue = f.resultValue
		}

	case *InstructionsFragment:
		blocksFragment.CurrentBlock.Insts = append(blocksFragment.CurrentBlock.Insts, f.Instructions...)
	}
}

func (blocksFragment *BlocksFragment) GetResultType() ResultType {
	return blocksFragment.resultType
}

func (blocksFragment *BlocksFragment) GetResult() value.Value {
	switch blocksFragment.resultType {
	case VOID:
		return nil
	case VALUE, POINTER:
		return blocksFragment.resultValue
	default:
		panic("does not support other result type" + strconv.Itoa((int)(blocksFragment.resultType)))
	}
}

func (blocksFragment *BlocksFragment) NewBlock(name string) {
	newBlock := ir.NewBlock(name)

	blocksFragment.ChainBlocks(newBlock)
	blocksFragment.Blocks = append(blocksFragment.Blocks, newBlock)
	blocksFragment.CurrentBlock = newBlock
}

func (blocksFragment *BlocksFragment) AddBlock(block *ir.Block) {
	if block != nil {
		blocksFragment.ChainBlocks(block)
		blocksFragment.Blocks = append(blocksFragment.Blocks, block)
		blocksFragment.CurrentBlock = block
	}
}

func (blocksFragment *BlocksFragment) ChainBlocks(blocks ...*ir.Block) {
	if blocksFragment.CurrentBlock != nil && blocksFragment.CurrentBlock.Term == nil &&
		blocks != nil && len(blocks) > 0 {
		blocksFragment.CurrentBlock.NewBr(blocks[0])
	}
}

type InstructionsFragment struct {
	Instructions       []ir.Instruction
	currentInstruction *ir.Instruction
	resultType         ResultType
	resultValue        value.Value
}

func NewInstructionsFragment(resultType ResultType) *InstructionsFragment {
	instructions := make([]ir.Instruction, 0)

	return &InstructionsFragment{Instructions: instructions, currentInstruction: nil, resultType: resultType, resultValue: nil}
}

func (instructionsFragment *InstructionsFragment) Append(fragment Fragment) {
	switch f := fragment.(type) {
	case *ModuleFragment:
		panic("Cannot append instructions to module")
	case *FunctionsFragment:
		panic("Cannot append instructions to functions")
	case *BlocksFragment:
		panic("Cannot append instructions to blocks")
	case *InstructionsFragment:
		instructionsFragment.Instructions = append(instructionsFragment.Instructions, f.Instructions...)
		instructionsFragment.currentInstruction = f.currentInstruction
	}
}

func (instructionsFragment *InstructionsFragment) AppendWithResultPropogation(fragment Fragment) {

}

func (instructionsFragment *InstructionsFragment) GetResultType() ResultType {
	return instructionsFragment.resultType
}

func (instructionsFragment *InstructionsFragment) GetResult() value.Value {
	switch instructionsFragment.resultType {
	case VOID:
		return nil
	case VALUE:
		return instructionsFragment.resultValue
	default:
		panic("does not support other result type" + strconv.Itoa((int)(instructionsFragment.resultType)))
	}
}

func (instructionsFragment *InstructionsFragment) AddInstruction(instr ir.Instruction) {
	instructionsFragment.Instructions = append(instructionsFragment.Instructions, instr)
}
