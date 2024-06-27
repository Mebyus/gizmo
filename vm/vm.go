package vm

import (
	"encoding/binary"
	"fmt"
)

type Machine struct {
	// Instruction pointer. Index in text memory.
	ip uint64

	// Stack pointer. Index in stack memory.
	sp uint64

	// Frame pointer. Index in stack memory.
	fp uint64

	// Comparison flags.
	cf uint64

	// General-purpose registers.
	r [64]uint64

	// Code of the program being executed, size cannot change during execution.
	text []byte

	// Static, read-only program data. Loaded at program start.
	data []byte

	// Memory for global variables. Loaded and initialized at program start.
	global []byte

	// Stack memory, size cannot change during execution.
	stack []byte

	// Heap memory, size can change during execution.
	heap []byte

	// Runtime error occured while executing current instruction.
	err error

	// Indicates if jump occured while executing current instruction.
	jump bool

	// Indicates if vm was halted by instruction or runtime error.
	halt bool
}

type Prog struct {
	text   []byte
	data   []byte
	global []byte
}

type Config struct {
	StackSize    int
	InitHeapSize int
}

func (m *Machine) Init(cfg *Config) {
	m.stack = make([]byte, cfg.StackSize)
	m.heap = make([]byte, cfg.InitHeapSize)
}

func (m *Machine) Exec(prog *Prog) *Exit {
	m.text = prog.text
	m.data = prog.data
	m.global = prog.global

	// reset vm state
	m.halt = false
	m.ip = 0
	m.sp = 0
	m.fp = 0

	for !m.halt {
		m.step()
	}

	return m.exit()
}

func (m *Machine) step() {
	m.jump = false

	ip := m.ip
	if ip >= uint64(len(m.text)) {
		m.stop(fmt.Errorf("end of program text reached"))
		return
	}

	op := m.text[ip]
	size := Size[op]

	if size == 0 {
		m.stop(fmt.Errorf("unknown opcode 0x%02X", op))
		return
	}
	if ip+uint64(size) > uint64(len(m.text)) {
		m.stop(fmt.Errorf("not enough code to read %d bytes for %s instruction (0x%02X)", size, "", op))
		return
	}

	var err error
	switch Opcode(op) {
	case Nop:
		// no operation
	case Halt:
		m.halt = true
		return
	case LoadValReg:
		err = m.loadValReg()
	case LoadRegReg:
		err = m.loadRegReg()
	default:
		// unknown opcodes should be handled via size check
		panic(fmt.Sprintf("unhandled instruction 0x%02X", op))
	}
	if err != nil {
		m.stop(err)
		return
	}

	if !m.jump {
		m.ip += uint64(size)
	}
}

// switch to halt state with runtime error
func (m *Machine) stop(err error) {
	m.err = err
	m.halt = true
}

// get n bytes of current instruction data (opcode not included)
func (m *Machine) instd(n uint64) []byte {
	ip := m.ip
	return m.text[ip+1 : ip+1+n]
}

// get register value
func (m *Machine) get(r uint8) (uint64, error) {
	if r >= 64 {
		return 0, fmt.Errorf("register index %d out of range", r)
	}
	v := m.r[r]
	return v, nil
}

// set register value
func (m *Machine) set(r uint8, v uint64) error {
	if r >= 64 {
		return fmt.Errorf("register index %d out of range", r)
	}
	m.r[r] = v
	return nil
}

// set register value
func (m *Machine) unsafeSet(r uint8, v uint64) {
	m.r[r] = v
}

func val64(buf []byte) uint64 {
	return binary.LittleEndian.Uint64(buf)
}

// Exit describes vm exit state after program execution.
// Includes both normal and abnormal
type Exit struct {
	// Runtime error for abnormal exit.
	Error error

	// Value of instruction pointer register.
	IP uint64

	// Exit status of the program. Obtained from first general-purpose register
	// upon program exit. Valid only for normal exit.
	Status uint64

	// True for normal exit. Occurs via explicit halt instruction.
	Normal bool
}

func (m *Machine) exit() *Exit {
	e := &Exit{IP: m.ip}

	if m.err != nil {
		e.Error = m.err
		return e
	}

	e.Normal = true
	e.Status = m.r[0]
	return e
}
