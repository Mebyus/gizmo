package amd64

import "fmt"

// Command is something that can operate Encoder.
type Command interface {
	do(*Encoder)
}

// InsKind unique machine instruction identifier.
// Corresponds to concrete opcode (for CPU).
type InsKind uint32

const (
	InsNil InsKind = iota

	// Invoke system call.
	InsSystemCall

	InsDecFamily0Reg64

	// Perform bitwise not on Family0 64-bit register.
	// Result is placed into the same register.
	InsNotFamily0Reg64

	// Perform bitwise not on Family1 64-bit register.
	// Result is placed into the same register.
	InsNotFamily1Reg64

	// Perform relative jump by 32-bit immediate value offset.
	// Offset is relative to RIP (its would-be-value after jump instruction).
	// Offset is signed and encoded as 2's complement 32-bit integer.
	InsRelJumpImm32

	InsRelJumpNotZeroImm32

	// Place immediate 32-bit value into Family0 64-bit wide register.
	InsMoveFamily0Reg64Imm32

	// Place immediate 32-bit value into Family1 64-bit wide register.
	InsMoveFamily1Reg64Imm32

	// Add immediate 32-bit value to Family0 64-bit wide register.
	// Result is placed into the same register.
	InsAddFamily0Reg64Imm32

	// Add immediate 32-bit value to Family1 64-bit wide register.
	// Result is placed into the same register.
	InsAddFamily1Reg64Imm32

	InsXorFamily0Reg64Family0Reg64

	InsTestFamily0Reg64Family0Reg64

	// Copy Family0 64-bit register into Family0 64-bit.
	InsCopyFamily0Reg64Family0Reg64

	// Copy Family1 64-bit register into Family0 64-bit.
	InsCopyFamily0Reg64Family1Reg64

	// Copy Family0 64-bit register into Family1 64-bit.
	InsCopyFamily1Reg64Family0Reg64

	// Copy Family1 64-bit register into Family1 64-bit.
	InsCopyFamily1Reg64Family1Reg64
)

// Ins0 instruction with no operands, only opcode.
type Ins0 struct {
	Kind InsKind
}

// Ins1Reg instruction with one register operand.
type Ins1Reg struct {
	Kind InsKind
	Reg  Register
}

// Ins1Reg instruction with one immediate value operand.
type Ins1Imm struct {
	Val  uint64
	Kind InsKind
}

// Ins2RegImm instruction with two operands:
//
// Destination: register
// Source:      immediate value
type Ins2RegImm struct {
	Val  uint64
	Kind InsKind
	Reg  Register
}

// Ins2RegReg instruction with two operands:
//
// Destination: register 0
// Source:      register 1
type Ins2RegReg struct {
	Kind InsKind

	Reg0 Register
	Reg1 Register
}

func invalid(kind InsKind) string {
	return fmt.Sprintf("invalid instruction (%d)", kind)
}

func (g Ins0) do(e *Encoder) {
	switch g.Kind {
	case InsSystemCall:
		e.u8s(0x0F, 0x05)
	default:
		panic(invalid(g.Kind))
	}
}

func (g Ins1Reg) do(e *Encoder) {
	switch g.Kind {
	case InsNotFamily0Reg64:
		e.u8s(0x48, 0xF7, 0xD0|g.Reg.Num())
	case InsNotFamily1Reg64:
		e.u8s(0x49, 0xF7, 0xD0|g.Reg.Num())
	case InsDecFamily0Reg64:
		e.u8s(0x48, 0xFF, 0xC8|g.Reg.Num())
	default:
		panic(invalid(g.Kind))
	}
}

func (g Ins1Imm) do(e *Encoder) {
	switch g.Kind {
	case InsRelJumpImm32:
		e.u8(0xE9)

		// jump offset is late value: it may be unknown during
		// encoder first pass
		label := g.Val // at this stage g.Val holds label index for this instruction
		e.putBackRel32OffsetOrBackpatch(label)
	case InsRelJumpNotZeroImm32:
		e.u8s(0x0F, 0x85)

		// jump offset is late value: it may be unknown during
		// encoder first pass
		label := g.Val // at this stage g.Val holds label index for this instruction
		e.putBackRel32OffsetOrBackpatch(label)
	default:
		panic(invalid(g.Kind))
	}
}

func (g Ins2RegImm) do(e *Encoder) {
	switch g.Kind {
	case InsMoveFamily0Reg64Imm32:
		e.u8s(0x48, 0xC7, 0xC0|g.Reg.Num())
		e.u32(uint32(g.Val))
	case InsMoveFamily1Reg64Imm32:
		e.u8s(0x49, 0xC7, 0xC0|g.Reg.Num())
		e.u32(uint32(g.Val))
	case InsAddFamily0Reg64Imm32:
		e.u8s(0x48, 0x81, 0xC0|g.Reg.Num())
		e.u32(uint32(g.Val))
	case InsAddFamily1Reg64Imm32:
		e.u8s(0x49, 0x81, 0xC0|g.Reg.Num())
		e.u32(uint32(g.Val))
	default:
		panic(invalid(g.Kind))
	}
}

func (g Ins2RegReg) do(e *Encoder) {
	switch g.Kind {
	case InsXorFamily0Reg64Family0Reg64:
		e.u8s(0x48, 0x31, 0xC0|(g.Reg1.Num()<<3)|g.Reg0.Num())
	case InsTestFamily0Reg64Family0Reg64:
		e.u8s(0x48, 0x85, 0xC0|(g.Reg1.Num()<<3)|g.Reg0.Num())
	default:
		panic(invalid(g.Kind))
	}
}

// ComLabel command which places label at current Encoder position.
type ComLabel struct {
	Label uint64 // label index
}

func (g ComLabel) do(e *Encoder) {
	e.placeLabel(g.Label)
}
