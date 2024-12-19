package amd64

import (
	"fmt"
	"strconv"
)

// MneKind assembly operation mnemonic acts as instruction prototype discriminator.
type MneKind uint32

const (
	// Zero value of MneKind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where MneKind is left unspecified.
	MneNil MneKind = iota

	MneMov
	MneSysCall
	MneAdd
	MneJump
	MneTest
	MneInc
	MneDec
	MneXor
	MneAnd
	MneNot
)

var meText = [...]string{
	MneNil: "<nil>",

	MneMov:     "mov",
	MneSysCall: "syscall",
	MneAdd:     "add",
	MneJump:    "jump",
	MneTest:    "test",
	MneInc:     "inc",
	MneDec:     "dec",
	MneXor:     "xor",
	MneAnd:     "and",
	MneNot:     "not",
}

var meLookup map[string]MneKind

func init() {
	meLookup = make(map[string]MneKind, len(meText))
	for m, s := range meText {
		if s == "" {
			panic(fmt.Errorf("empty mnemonic (%d) text", m))
		}
		meLookup[s] = MneKind(m)
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
