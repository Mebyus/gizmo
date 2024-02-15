package token

import "fmt"

var words = map[string]Kind{
	"import":   Import,
	"fn":       Fn,
	"continue": Continue,
	"return":   Return,
	"break":    Break,
	"case":     Case,
	"for":      For,
	"else":     Else,
	"if":       If,
	"defer":    Defer,
	"bag":      Bag,
	"bind":     Bind,
	"unit":     Unit,
	"in":       In,
	"var":      Var,
	"type":     Type,
	"switch":   Switch,
	"enum":     Enum,
	"struct":   Struct,
	"union":    Union,
	"pub":      Pub,
	"never":    Never,
	"dirty":    Dirty,
	"nil":      Nil,
	"true":     True,
	"false":    False,
	"atr":      Atr,
	"declare":  Declare,
	"const":    Const,
	"method":   Method,
}

const (
	minKeywordLen = 2
	maxKeywordLen = 8
)

// Lookup finds token Kind (if any) by its literal
func Lookup(lit string) (Kind, bool) {
	if len(lit) < minKeywordLen || len(lit) > maxKeywordLen {
		return empty, false
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
