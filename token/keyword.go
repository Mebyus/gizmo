package token

import "fmt"

var words = map[string]Kind{
	"import": Import,
	"fun":    Fun,
	"test":   Test,
	"jump":   Jump,
	"ret":    Return,
	"for":    For,
	"else":   Else,
	"if":     If,
	"defer":  Defer,
	"bag":    Bag,
	"unit":   Unit,
	"in":     In,
	"var":    Var,
	"let":    Let,
	"type":   Type,
	"enum":   Enum,
	"struct": Struct,
	"union":  Union,
	"pub":    Pub,

	"never": Never,
	"stub":  Stub,
	"dirty": Dirty,
	"nil":   Nil,
	"true":  True,
	"false": False,

	"cast":  Cast,
	"tint":  Tint,
	"mcast": MemCast,
	"msize": MemSize,

	"any": Any,
}

const (
	minKeywordLen = 2
	maxKeywordLen = 6
)

// Lookup finds token Kind (if any) by its literal
func Lookup(lit string) (Kind, bool) {
	if len(lit) < minKeywordLen || len(lit) > maxKeywordLen {
		return 0, false
	}
	k, ok := words[lit]
	return k, ok
}

func init() {
	minLen := 1 << 10 // arbitrary large number
	maxLen := 0

	for word, kind := range words {
		if len(word) > maxLen {
			maxLen = len(word)
		}
		if len(word) < minLen {
			minLen = len(word)
		}

		lit := Literal[kind]
		if lit != word {
			panic(fmt.Sprintf("keyword \"%s\" has inconsistent literal", word))
		}
	}

	if minLen != minKeywordLen {
		panic(fmt.Sprintf("min keyword length should be %d, not %d", minLen, minKeywordLen))
	}
	if maxLen != maxKeywordLen {
		panic(fmt.Sprintf("max keyword length should be %d, not %d", maxLen, maxKeywordLen))
	}
}
