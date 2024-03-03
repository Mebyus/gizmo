package source

import (
	"fmt"
	"sync"
)

const debug = true

// Loader caches loaded source files in memory.
// Safe for usage from multiple goroutines
type Loader struct {
	// maps path to file entry
	entries map[string]*Entry

	mu sync.Mutex
}

type state uint8

const (
	empty state = iota

	// failed to load file
	failed

	// file info loaded
	basic

	// file info + content loaded
	loaded
)

var stateText = [...]string{
	empty: "<nil>",

	failed: "failed",
	basic:  "basic",
	loaded: "loaded",
}

func (s state) String() string {
	return stateText[s]
}

type Entry struct {
	path string

	// not nil only for failed state
	err error

	// always nil for failed state
	file *File

	state state
}

func (e *Entry) failed(err error) {
	e.err = err
	e.file = nil
	e.state = failed
}

func (e *Entry) loaded() {
	e.state = loaded
}

func (e *Entry) loadFromInfo() {
	err := e.file.Load()
	if err != nil {
		e.failed(err)
		return
	}
	e.loaded()
}

func newFailed(path string, err error) *Entry {
	return &Entry{
		path:  path,
		state: failed,
		err:   err,
	}
}

func newBasic(file *File) *Entry {
	return &Entry{
		path:  file.Path,
		state: basic,
		file:  file,
	}
}

func newLoaded(file *File) *Entry {
	return &Entry{
		path:  file.Path,
		state: loaded,
		file:  file,
	}
}

func NewLoader() *Loader {
	return &Loader{entries: make(map[string]*Entry)}
}

func (l *Loader) debug(format string, args ...any) {
	fmt.Print("[debug] loader  | ")
	fmt.Printf(format, args...)
	fmt.Println()
}

// Load source file specified by path. Argument should be
// normalized by filepath.Clean to avoid duplication due to
// mapping of different path strings to the same file
//
// File is loaded and then cached for subsequent Load and Info calls
func (l *Loader) Load(path string) (*File, error) {
	entry := l.lookup(path)
	if entry == nil {
		entry = l.loadAndSave(path)
		return entry.file, entry.err
	}

	switch entry.state {
	case empty:
		panic("empty state")
	case failed, loaded:
		return entry.file, entry.err
	case basic:
		entry.loadFromInfo()
		return entry.file, entry.err
	default:
		panic(fmt.Sprintf("unexpected state: %d", entry.state))
	}
}

// Info works by the same logic as Load but does not load file contents,
// only basic information about file
func (l *Loader) Info(path string) (*File, error) {
	entry := l.lookup(path)
	if entry == nil {
		entry = l.infoAndSave(path)
		return entry.file, entry.err
	}

	switch entry.state {
	case empty:
		panic("empty state")
	case failed, loaded, basic:
		return entry.file, entry.err
	default:
		panic(fmt.Sprintf("unexpected state: %d", entry.state))
	}
}

func (l *Loader) infoAndSave(path string) *Entry {
	var entry *Entry

	file, err := Info(path)
	if err != nil {
		entry = newFailed(path, err)
	} else {
		entry = newBasic(file)
	}

	l.save(entry)
	return entry
}

func (l *Loader) loadAndSave(path string) *Entry {
	var entry *Entry

	file, err := Load(path)
	if err != nil {
		entry = newFailed(path, err)
	} else {
		entry = newLoaded(file)
	}

	l.save(entry)
	return entry
}

func (l *Loader) save(entry *Entry) {
	if debug {
		if entry.state == failed {
			l.debug("save \"%s\" (state=%s) %s",
				entry.path, entry.state.String(), entry.err)
		} else if entry.state == basic {
			l.debug("save \"%s\" (state=%s size=%d)",
				entry.path, entry.state.String(), entry.file.Size)
		} else {
			l.debug("save \"%s\" (state=%s hash=%016x size=%d)",
				entry.path, entry.state.String(), entry.file.Hash, entry.file.Size)
		}
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	path := entry.path
	e := l.entries[path]
	if e == nil {
		l.entries[path] = entry
		return
	}
	if e.state == loaded {
		return
	}
	if e.state < entry.state {
		l.entries[path] = entry
	}
}

func (l *Loader) lookup(path string) *Entry {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.entries[path]
}
