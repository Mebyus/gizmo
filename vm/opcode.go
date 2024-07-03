package vm

type Opcode uint8

const (
	// Reserved empty instruction. Does nop (no operation).
	Nop Opcode = 0x0004

	// Halts program execution.
	Halt Opcode = 0x0001

	// Call builtin procedure.
	SysCall Opcode = 0x0002

	// Trap causes immediate abnormal program exit.
	// This instruction is much like halt, but intented
	// to be placed as a trap for code which must not be executed
	// thus causing runtime error.
	//
	// Most notable usecase for this is to detect illegal
	// control flow execution of unreachable code.
	Trap Opcode = 0x0000

	// Set register value to zero.
	ClearReg Opcode = 0x0008

	// Set 8 consecutive bytes in memory to zero.
	// Address specified by constant pointer.
	ClearAddr Opcode = 0x0009

	// Set 8 consecutive bytes in memory to zero.
	// Address specified by pointer stored in register.
	ClearPtr Opcode = 0x000A

	// Load constant value to register.
	LoadValReg Opcode = 0x0010

	// Load constant value to system register.
	LoadValSysReg Opcode = 0x001A

	// Load register value to another register.
	LoadRegReg Opcode = 0x0011

	// Load value stored in memory to register.
	// Address specified by constant pointer.
	LoadAddrReg Opcode = 0x0012

	// Load value stored in register to memory.
	// Address specified by constant pointer.
	LoadRegAddr Opcode = 0x0013

	// Load constant value to memory.
	// Address specified by constant pointer.
	LoadValAddr Opcode = 0x0014

	// Load value stored in memory to register.
	// Address specified by pointer stored in register.
	LoadPtrReg Opcode = 0x0015

	// Load value stored in memory to register.
	// Address specified by constant offset added with pointer stored in register.
	LoadPtrOffReg Opcode = 0x0016

	// Load value stored in register to memory.
	// Address specified by pointer stored in register.
	LoadRegPtr Opcode = 0x0017

	// Load value stored in register to memory.
	// Address specified by constant offset added with pointer stored in register.
	LoadRegPtrOff Opcode = 0x0018

	// Add values stored in two registers.
	// Result is stored in destination register.
	AddRegReg Opcode = 0x0020

	// Add constant to value stored in register.
	// Result is stored in destination register.
	AddValReg Opcode = 0x0021

	// Add one to value stored in register.
	IncReg Opcode = 0x0022

	// Subtract constant value from value stored in register.
	// Result is stored in destination register.
	SubRegVal Opcode = 0x0028

	// Subtract value stored in register from constant value.
	// Result is stored in destination register.
	SubValReg Opcode = 0x0029

	// Subtract value stored in register from value stored in register.
	// Result is stored in destination register.
	SubRegReg Opcode = 0x002A

	// Compare constant value with value stored in register.
	TestValReg Opcode = 0x0030

	// Compare value stored in register with constant value.
	TestRegVal Opcode = 0x0031

	// Compare value stored in register with value stored in register.
	CompRegReg Opcode = 0x0032

	// Jump to constant address.
	JumpAddr Opcode = 0x0040

	// Jump to constant address if comparison flags register indicates not zero.
	JumpAddrNotZero = 0x0041

	// Push constant value to stack.
	PushVal Opcode = 0x0070

	// Push value stored in register to stack.
	PushReg Opcode = 0x0071

	// Pop value stored on top of the stack to register.
	PopReg Opcode = 0x0078

	// Call procedure.
	// Address specified by constant value.
	CallVal Opcode = 0x0080

	// Call procedure.
	// Address specified by pointer stored in register.
	CallPtr Opcode = 0x0081

	// Return from procedure.
	Ret Opcode = 0x0088
)

// Size allows to obtain instruction size denoted
// by specific opcode. Result gives full instruction
// size i.e. total length in bytes of opcode + data.
// If stored value for opcode is 0 then it means that
// vm does not recognize the given opcode, and in that
// case a runtime error must be generated.
//
// In addition layout of each instruction is documented here.
// Layout includes byte sizes of each instruction segment
// and their order.
//
//	OP - 1 byte  - opcode
//	SR - 1 byte  - source register index
//	DR - 1 byte  - destination register index
//	MR - 1 byte  - system register index
//	TA - 4 bytes - text address (for jumps)
//	CV - 8 bytes - constant value
//	CA - 8 bytes - constant address in memory
var Size = [...]uint8{
	// OP
	Nop: 1,

	// OP
	Halt: 1,

	// OP
	SysCall: 1,

	// OP
	Trap: 1,

	// OP + DR
	ClearReg: 1 + 1,

	// OP + DR + CV
	LoadValReg: 1 + 1 + 8,

	// OP + CV
	LoadValSysReg: 1 + 4,

	// OP + DR + SR,
	LoadRegReg: 1 + 1 + 1,

	// OP + DR + CA,
	LoadAddrReg: 1 + 1 + 8,

	// OP + DR + SR,
	AddRegReg: 1 + 1 + 1,

	// OP + DR + CV
	AddValReg: 1 + 1 + 8,

	// OP + DR
	IncReg: 1 + 1,

	// OP + SR + CV
	TestValReg: 1 + 1 + 8,

	// OP + SR + CV
	TestRegVal: 1 + 1 + 8,

	// OP + TA
	JumpAddr: 1 + 4,

	// OP + TA
	JumpAddrNotZero: 1 + 4,
}
