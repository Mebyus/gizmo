package lbl

// Kind indicates label kind
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly
	//
	// Mostly a trick to detect places where Kind is left unspecified
	empty Kind = iota

	Named

	Next
	Out
)

var text = [...]string{
	empty: "<nil>",

	Named: "named",

	Next: "next",
	Out:  "out",
}

func (k Kind) String() string {
	return text[k]
}
