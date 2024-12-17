package amd64

import (
	"fmt"
	"strconv"
)

// Operation acts as a prototype for concrete machine instruction.
// It is an intermediate representation between assembly "instruction"
// (written as text line in file) and machine instruction with binary
// encoding.
//
// This abstraction is born from the fact that amd64 assembly has wildly
// different encodings for "instructions" with the same mnemonic.
type Operation struct {
	// Operands. Unused Operands are set to nil.
	Ops [3]bagOp

	Me Me
}

// bagOp bag for operation operand (argument).
type bagOp interface {
	bagOp()
}

// Me assembly operation mnemonic acts as Operation discriminator.
type Me uint32

const (
	// Zero value of Mo. Should not be used explicitly.
	//
	// Mostly a trick to detect places where OpKind is left unspecified.
	MeNil Me = iota

	MeMov
	MeSysCall
	MeAdd
	MeJump
	MeTest
	MeInc
	MeXor
	MeAnd
	MeNot
)

var meText = [...]string{
	MeMov:     "mov",
	MeSysCall: "syscall",
	MeAdd:     "add",
	MeJump:    "jump",
	MeTest:    "test",
	MeInc:     "inc",
	MeXor:     "xor",
	MeAnd:     "and",
	MeNot:     "not",
}

var meLookup map[string]Me

func init() {
	meLookup = make(map[string]Me, len(meText))
	for m, s := range meText {
		if s == "" {
			panic(fmt.Errorf("empty mnemonic (%d) text", m))
		}
		meLookup[s] = Me(m)
	}
}

const (
	PropPtr = 0
	PropLen = 1
)

var propText = [...]string{
	PropPtr: "ptr",
	PropLen: "len",
}

const (
	FlagZero    = 0
	FlagNotZero = 1
)

var flagText = [...]string{
	FlagZero:    "z",
	FlagNotZero: "nz",
}

type Register uint16

const (
	RAX Register = 0x0
	RCX Register = 0x1
	RDX Register = 0x2
	RBX Register = 0x3
	RSP Register = 0x4
	RBP Register = 0x5
	RSI Register = 0x6
	RDI Register = 0x7

	R8  Register = 0x10
	R9  Register = 0x11
	R10 Register = 0x12
	R11 Register = 0x13
	R12 Register = 0x14
	R13 Register = 0x15
	R14 Register = 0x16
	R15 Register = 0x17
)

// Num returns register number inside its family.
func (r Register) Num() uint8 {
	return uint8(r & 0xF)
}

func (r Register) Family() uint8 {
	return uint8((r >> 4) & 0xF)
}

var family0RegText = [...]string{
	RAX: "rax",
	RCX: "rcx",
	RDX: "rdx",
	RBX: "rbx",
	RSP: "rsp",
	RBP: "rbp",
	RSI: "rsi",
	RDI: "rdi",
}

func (r Register) String() string {
	switch f := r.Family(); f {
	case 0:
		return family0RegText[r.Num()]
	case 1:
		return "r" + strconv.FormatUint(8+uint64(r.Num()), 10)
	default:
		panic(fmt.Errorf("unknown family (%d)", f))
	}
}
