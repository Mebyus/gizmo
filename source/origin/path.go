package origin

import (
	"fmt"
	"hash/fnv"
	"sort"
)

// Path a combination of import string and origin is called "unit path".
// It uniquely identifies a unit (local or foreign) within a project
type Path struct {
	// Always empty if origin is empty. Always not empty if origin is not empty
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

func Local(s string) Path {
	return Path{
		ImpStr: s,
		Origin: Loc,
	}
}

func Locals(ss []string) []Path {
	if len(ss) == 0 {
		return nil
	}

	paths := make([]Path, 0, len(ss))
	for _, s := range ss {
		paths = append(paths, Local(s))
	}
	return paths
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
