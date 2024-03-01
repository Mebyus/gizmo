package source

import (
	"fmt"
	"sync"
)

// Loader caches loaded source files in memory.
// Safe for usage from multiple goroutines
type Loader struct {
	// maps path to file entry
	entries map[string]*Entry

	mu sync.Mutex
}

type State uint8

const (
	empty State = iota

	// failed to load file
	failed

	// file info loaded
	basic

	// file info + content loaded
	loaded
)

type Entry struct {
	// not nil only for failed state
	err error

	// always nil for failed state
	file *File

	state State
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

func newFailed(err error) *Entry {
	return &Entry{
		state: failed,
		err:   err,
	}
}

func newBasic(file *File) *Entry {
	return &Entry{
		state: basic,
		file:  file,
	}
}

func newLoaded(file *File) *Entry {
	return &Entry{
		state: loaded,
		file:  file,
	}
}

func NewLoader() *Loader {
	return &Loader{entries: make(map[string]*Entry)}
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
		entry = newFailed(err)
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
		entry = newFailed(err)
	} else {
		entry = newLoaded(file)
	}

	l.save(entry)
	return entry
}

func (l *Loader) save(entry *Entry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	path := entry.file.Path
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
