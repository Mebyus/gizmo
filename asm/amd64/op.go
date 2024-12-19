package amd64

type ArgKind uint8

const (
	ArgNil ArgKind = iota

	// Mnemonic argument.
	ArgMne

	ArgInteger
	ArgRegister
	ArgLabel
	ArgFlag
)

type OpArg struct {
	Val  uint64
	Kind ArgKind
}

// Operation represents text line from assembly source text in compact way.
type Operation struct {
	Args [2]OpArg

	// Meaning depends on Op field.
	Val uint64

	Op OpKind

	// Number of arguments inside array.
	Num uint8
}

// OpKind specifies what assembly line does.
type OpKind uint8

const (
	OpNil OpKind = iota

	// Operation acts as a prototype for concrete machine instruction.
	// It is an intermediate representation between assembly "instruction"
	// (written as text line in file) and machine instruction with binary
	// encoding.
	//
	// This abstraction is born from the fact that amd64 assembly has wildly
	// different encodings for "instructions" with the same mnemonic.
	OpProto

	// Label placement.
	OpLabel
)

func translate(o *Operation) Command {
	switch o.Op {
	case OpProto:
		switch MneKind(o.Val) {
		case MneSysCall:
			return translateSysCall(o)
		case MneMov:
			return translateMov(o)
		case MneXor:
			return translateXor(o)
		case MneDec:
			return translateDec(o)
		case MneTest:
			return translateTest(o)
		case MneJump:
			return translateJump(o)
		default:
			panic("invalid mnemonic")
		}
	case OpLabel:
		return ComLabel{Label: o.Val}
	default:
		panic("invalid operation")
	}
}

func translateJump(o *Operation) Command {
	if o.Num == 1 {
		panic("not implemented")
	}

	if o.Num != 2 {
		panic("not enough args")
	}

	f := o.Args[0]
	d := o.Args[1]
	if f.Kind != ArgFlag {
		panic("invalid arg")
	}
	if d.Kind != ArgLabel {
		panic("invalid arg")
	}

	if f.Val != FlagNotZero {
		panic("not implemented")
	}

	return Ins1Imm{
		Kind: InsRelJumpNotZeroImm32,
		Val:  d.Val, // label index
	}
}

func translateTest(o *Operation) Command {
	if o.Num != 2 {
		panic("not enough args")
	}

	d := o.Args[0]
	s := o.Args[1]
	switch {
	case d.Kind == ArgRegister && s.Kind == ArgRegister:
		return Ins2RegReg{
			Kind: InsTestFamily0Reg64Family0Reg64,
			Reg0: Register(d.Val),
			Reg1: Register(s.Val),
		}
	default:
		panic("not implemented")
	}
}

func translateXor(o *Operation) Command {
	if o.Num != 2 {
		panic("not enough args")
	}

	d := o.Args[0]
	s := o.Args[1]
	switch {
	case d.Kind == ArgRegister && s.Kind == ArgRegister:
		return Ins2RegReg{
			Kind: InsXorFamily0Reg64Family0Reg64,
			Reg0: Register(d.Val),
			Reg1: Register(s.Val),
		}
	default:
		panic("not implemented")
	}
}

func translateMov(o *Operation) Command {
	if o.Num != 2 {
		panic("not enough args")
	}

	d := o.Args[0]
	s := o.Args[1]
	switch {
	case d.Kind == ArgRegister && s.Kind == ArgInteger:
		return Ins2RegImm{
			Kind: InsMoveFamily0Reg64Imm32,
			Reg:  Register(d.Val),
			Val:  s.Val,
		}
	default:
		panic("not implemented")
	}
}

func translateDec(o *Operation) Command {
	if o.Num != 1 {
		panic("not enough args")
	}

	d := o.Args[0]
	return Ins1Reg{
		Kind: InsDecFamily0Reg64,
		Reg:  Register(d.Val),
	}
}

func translateSysCall(o *Operation) Command {
	return Ins0{Kind: InsSystemCall}
}
