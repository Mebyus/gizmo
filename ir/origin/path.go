package origin

import (
	"fmt"
	"hash/fnv"
	"sort"
)

// Path a combination of import string and origin is called "unit path".
// It uniquely identifies a unit (local or foreign) within a project
type Path struct {
	// Always not empty if origin is not empty
	ImpStr string

	Origin Origin
}

var Empty = Path{}

func (p Path) IsEmpty() bool {
	return p.Origin.IsEmpty()
}

func (p Path) String() string {
	return fmt.Sprintf("%s: %s", p.Origin.String(), p.ImpStr)
}

func Sort(p []Path) {
	sort.Slice(p, func(i, j int) bool {
		a := p[i]
		b := p[j]
		return Less(a, b)
	})
}

// Less returns true if a is less than b
func Less(a, b Path) bool {
	return a.Origin < b.Origin || (a.Origin == b.Origin && a.ImpStr < b.ImpStr)
}

func (p Path) Hash() uint64 {
	var buf [2]byte

	h := fnv.New64a()
	// leave second byte zeroed as a separator between origin and import string
	buf[0] = byte(p.Origin)
	h.Write(buf[:])
	h.Write([]byte(p.ImpStr))
	return h.Sum64()
}
