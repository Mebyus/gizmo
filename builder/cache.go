package builder

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mebyus/gizmo/ir/origin"
	"github.com/mebyus/gizmo/source"
)

func (cfg *Config) Hash() uint64 {
	return 53582900
}

type Cache struct {
	// Base directory for source files lookup
	srcdir string

	// Base directory for current cache instance. Includes build
	// config seed
	dir string

	src *source.Loader

	// If this flag is true than disk cache was empty before
	// build cache object was created. Thus by using the flag
	// we can quickly determine is it worth looking for anything
	// in cache
	//
	// In other words when cache is in init mode it cannot give
	// us anything (because it was not stored yet) and we only use
	// cache for storing for future use
	init bool
}

func NewCache(cfg *Config) (*Cache, error) {
	seed := cfg.Hash()
	dir := filepath.Join(cfg.BaseCacheDir, formatCacheSeed(seed))

	c := &Cache{
		dir: dir,
		src: source.NewLoader(),

		srcdir: cfg.BaseSourceDir,
	}

	err := c.initBaseDir()
	if err != nil {
		return nil, fmt.Errorf("init build cache dir: %w", err)
	}

	return c, nil
}

func (c *Cache) initBaseDir() error {
	info, err := os.Stat(c.dir)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		err = os.MkdirAll(c.dir, 0o775)
		if err != nil {
			return err
		}

		c.init = true
		return nil
	}

	if !info.IsDir() {
		return fmt.Errorf("file \"%s\" is not a directory", c.dir)
	}

	return nil
}

func (c *Cache) LookupBuild(unit string) {

}

func (c *Cache) LoadSourceFile(p origin.Path, name string) (*source.File, error) {
	switch p.Origin {
	case origin.Std:
		panic("not implemented for std")
	case origin.Pkg:
		panic("not implemented for pkg")
	case origin.Loc:
		return c.src.Load(filepath.Join(c.srcdir, p.ImpStr, name))
	default:
		panic("unexpected import origin: " + strconv.FormatInt(int64(p.Origin), 10))
	}
}

func (c *Cache) InfoSourceFile(p origin.Path, name string) (*source.File, error) {
	switch p.Origin {
	case origin.Std:
		panic("not implemented for std")
	case origin.Pkg:
		panic("not implemented for pkg")
	case origin.Loc:
		return c.src.Info(filepath.Join(c.srcdir, p.ImpStr, name))
	default:
		panic("unexpected import origin: " + strconv.FormatInt(int64(p.Origin), 10))
	}
}

func (c *Cache) SaveUnitGenout(p origin.Path, data []byte) {
	dir := filepath.Join(c.dir, "unit", formatHash(p.Hash()))
	err := os.MkdirAll(dir, 0o775)
	if err != nil {
		if debug {
			c.debug("warn: failed to create dir \"%s\": %s", dir, err)
		}
		return
	}

	path := filepath.Join(dir, "stapled_unit.cpp")
	err = os.WriteFile(path, data, 0o664)
	if err != nil {
		if debug {
			c.debug("warn: failed to save file \"%s\": %s", path, err)
		}
		return
	}

	if debug {
		c.debug("saved stapled unit \"%s\" (%d bytes)", p.String(), len(data))
	}
}

func (c *Cache) SaveFileGenout(p origin.Path, file *source.File, data []byte) {
	dir := filepath.Join(c.dir, "unit", formatHash(p.Hash()))
	err := os.MkdirAll(dir, 0o775)
	if err != nil {
		if debug {
			c.debug("warn: failed to create dir \"%s\": %s", dir, err)
		}
		return
	}

	path := filepath.Join(dir, genPartName(file))
	err = os.WriteFile(path, data, 0o664)
	if err != nil {
		if debug {
			c.debug("warn: failed to save file \"%s\": %s", path, err)
		}
		return
	}

	if debug {
		c.debug("saved file genout \"%s/%s\" (%d bytes)", p.String(), file.Name, len(data))
	}
}

// LoadUnitGenout load bytes stored in cache for combined generated unit output.
// If it is saved in cache and no older than supplied mod time than it will be
// returned with (<data>, true) otherwise (<nil>, false)
func (c *Cache) LoadUnitGenout(p origin.Path, mod time.Time) ([]byte, bool) {
	path := filepath.Join(c.dir, "unit", formatHash(p.Hash()), "stapled_unit.cpp")
	m, ok := getFileModTime(path)
	if !ok {
		return nil, false
	}
	if mod.After(m) {
		return nil, false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if debug {
			c.debug("warn: read existing file \"%s\": %s", path, err)
		}
		return nil, false
	}
	if debug {
		c.debug("stapled unit \"%s\" is up to date (%d bytes)", p.String(), len(data))
	}
	return data, true
}

// LoadFileGenout load bytes stored in chache for generated output of a file in particular
// unit denoted by origin path. If stored output is no older than supplied file than
// this method returns with (<data>, true) otherwise (<nil>, false)
func (c *Cache) LoadFileGenout(p origin.Path, file *source.File) ([]byte, bool) {
	path := filepath.Join(c.dir, "unit", formatHash(p.Hash()), genPartName(file))
	m, ok := getFileModTime(path)
	if !ok {
		return nil, false
	}
	if file.ModTime.After(m) {
		return nil, false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if debug {
			c.debug("warn: read existing file \"%s\": %s", path, err)
		}
		return nil, false
	}
	if debug {
		c.debug("file genout \"%s/%s\" is up to date (%d bytes)", p.String(), file.Name, len(data))
	}
	return data, true
}

func genPartName(file *source.File) string {
	return "part_" + formatHash(file.Hash) + "_" + strconv.FormatUint(file.Size, 10) + ".cpp"
}

func getFileModTime(path string) (time.Time, bool) {
	info, err := os.Lstat(path)
	if err != nil {
		return time.Time{}, false
	}
	if info.IsDir() {
		return time.Time{}, false
	}
	if !info.Mode().IsRegular() {
		return time.Time{}, false
	}
	return info.ModTime(), true
}

func (c *Cache) debug(format string, args ...any) {
	fmt.Print("[debug] cache   | ")
	fmt.Printf(format, args...)
	fmt.Println()
}

func formatHash(seed uint64) string {
	// formats integer in hex with exactly displayed 8 bytes (16 characters)
	return fmt.Sprintf("%016x", seed)
}

func formatCacheSeed(seed uint64) string {
	return formatHash(seed)
}
