package builder

import "sort"

type CacheMap struct {
	Unit map[string]*CacheUnit `json:"unit"`
}

type UnitPath struct {
	Origin string `json:"origin"`
	Import string `json:"import"`
}

type EntryPoint struct {
	Name   string `json:"name"`
	Origin string `json:"origin"`
}

type BuildPart struct {
	Hash uint64 `json:"hash"`
	Size uint64 `json:"size"`

	// unix microseconds
	Timestamp int64 `json:"timestamp"`
}

type CachePartFile struct {
	Name string `json:"name"`
	Size uint64 `json:"size"`
	Hash uint64 `json:"hash"`

	// unix microseconds
	Timestamp int64 `json:"timestamp"`
}

type GizmoParts struct {
	Files []CachePartFile `json:"files,omitempty"`
	Gen   string          `json:"gen,omitempty"`

	// unix microseconds
	Timestamp int64 `json:"timestamp"`

	// file map, maps file name to index inside Files slice
	fm map[string]int `json:"-"`
}

type AsmParts struct {
	Files []CachePartFile `json:"files"`
	Obj   string          `json:"obj"`

	// unix microseconds
	Timestamp int64 `json:"timestamp"`
}

type Parts struct {
	Build BuildPart  `json:"build"`
	Gizmo GizmoParts `json:"gizmo"`

	Asm *AsmParts `json:"asm,omitempty"`
}

type CacheUnit struct {
	Path  UnitPath `json:"path"`
	Parts Parts    `json:"parts"`

	EntryPoint *EntryPoint `json:"entrypoint,omitempty"`
}

func (p *GizmoParts) SortFiles() {
	if len(p.fm) == len(p.Files) {
		return
	}
	sortCacheParts(p.Files)
}

func sortCacheParts(parts []CachePartFile) {
	sort.Slice(parts, func(i, j int) bool {
		return parts[i].Name < parts[j].Name
	})
}
