package build

import (
	"fmt"
	"strings"
)

type Kind uint8

const (
	// debug-friendly optimizations + debug information in binaries + safety checks
	Debug Kind = iota + 1

	// moderate-level optimizations + debug information in binaries + safety checks
	Test

	// most optimizations enabled + safety checks
	Safe

	// all optimizations enabled + disabled safety checks
	Fast
)

var buildKindText = [...]string{
	0: "<nil>",

	Debug: "debug",
	Test:  "test",
	Safe:  "safe",
	Fast:  "fast",
}

func (k Kind) String() string {
	return buildKindText[k]
}

func Parse(s string) (Kind, error) {
	s = strings.TrimSpace(s)

	var k Kind
	switch s {
	case "":
		return 0, fmt.Errorf("empty build kind")
	case "debug":
		k = Debug
	case "test":
		k = Test
	case "safe":
		k = Safe
	case "fast":
		k = Fast
	default:
		return 0, fmt.Errorf("unknown \"%s\" build kind", s)
	}

	if k == 0 {
		panic("empty build kind")
	}

	return k, nil
}
