package builder

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/source/origin"
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

	cfg *Config

	src *source.Loader

	cmap *CacheMap

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
		cfg: cfg,

		srcdir: cfg.BaseSourceDir,
	}

	err := c.initBaseDir()
	if err != nil {
		return nil, fmt.Errorf("init build cache dir: %w", err)
	}
	c.loadMap()

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

func (c *Cache) loadMap() {
	c.loadMapFromFile()
	if c.cmap != nil && len(c.cmap.Unit) != 0 {
		if debug {
			c.debug("loaded map with %d stored units", len(c.cmap.Unit))
		}
		return
	}
	c.cmap = &CacheMap{Unit: make(map[string]*CacheUnit)}
	if debug {
		c.debug("using new blank map")
	}
}

func (c *Cache) loadMapFromFile() {
	if c.init {
		return
	}

	file, err := os.Open(filepath.Join(c.dir, "map.json"))
	if err != nil {
		if debug {
			c.debug("warn: failed to load map file: %s", err)
		}
		return
	}
	defer file.Close()

	var m CacheMap
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&m)
	if err != nil {
		if debug {
			c.debug("warn: failed to decode map file: %s", err)
		}
		return
	}
	c.cmap = &m
}

func (c *Cache) SaveMap() {
	path := filepath.Join(c.dir, "map.json")
	file, err := os.Create(path)
	if err != nil {
		if debug {
			c.debug("warn: failed to create map file: %s", err)
		}
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(c.cmap)
	if err != nil {
		if debug {
			c.debug("warn: failed to encode map file: %s", err)
		}
		return
	}

	if debug {
		c.debug("saved map with %d stored units as \"%s\"", len(c.cmap.Unit), path)
	}
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

func (c *Cache) InitUnitDir(p origin.Path) (string, error) {
	hash := formatHash(p.Hash())
	dir := filepath.Join(c.dir, "unit", hash)
	err := os.MkdirAll(dir, 0o775)
	if err != nil {
		return "", err
	}

	cachedUnit := c.cmap.Unit[hash]
	if cachedUnit == nil {
		c.cmap.Unit[hash] = &CacheUnit{
			Path: UnitPath{
				Origin: p.Origin.String(),
				Import: p.ImpStr,
			},
		}
	}

	return dir, nil
}

func (c *Cache) SaveFileGenout(p origin.Path, file *source.File, data []byte) {
	hash := formatHash(p.Hash())
	dir := filepath.Join(c.dir, "unit", hash)
	path := filepath.Join(dir, genPartName(file))
	err := os.WriteFile(path, data, 0o664)
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

func (c *Cache) SaveModGenout(mod string, code *PartsBuffer) (string, error) {
	name := filepath.Base(mod) + ".cpp"
	dir := filepath.Join(c.dir, "mod", "gen", filepath.Dir(mod))
	err := os.MkdirAll(dir, 0o775)
	if err != nil {
		return "", err
	}

	path := filepath.Join(dir, name)
	err = SavePartsBuffer(path, code)
	if err != nil {
		return "", err
	}
	return path, nil
}

// LoadUnitGenout load bytes stored in cache for combined generated unit output.
// If it is saved in cache and no older than supplied mod time than it will be
// returned with (<data>, true) otherwise (<nil>, false)
func (c *Cache) LoadUnitGenout(p origin.Path, mod time.Time) ([]byte, bool) {
	if c.init {
		return nil, false
	}
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
	if c.init {
		return nil, false
	}
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
	return "part_" + formatHash(file.Hash) + "_" + strconv.FormatUint(file.Size, 10) + ".gen.cpp"
}

// Generates hashed file name of the form:
//
//	"part_obj_<hash1>_<hash2>.o"
//	<hash1> = hash produced from file paths
//	<hash2> = hash produced from file content hashes
//
// slice argument must always contain at least one element
func genPartObjName(files []*source.File) string {
	hash1 := formatHash(filePathsHash(files))
	hash2 := formatHash(fileHashesHash(files))
	return "part_obj_" + hash1 + "_" + hash2 + ".o"
}

func filePathsHash(files []*source.File) uint64 {
	// calculate buffer size for this set of files
	var size int
	size += len(files[0].Path)
	for _, file := range files[1:] {
		size += 1 // for separator between paths
		size += len(file.Path)
	}

	buf := make([]byte, size)
	var i int // buf write index

	i += copy(buf[i:], files[0].Path)
	for _, file := range files[1:] {
		// write zero separator between paths
		buf[i] = 0
		i += 1

		i += copy(buf[i:], file.Path)
	}

	h := fnv.New64a()
	h.Write(buf)
	return h.Sum64()
}

func fileHashesHash(files []*source.File) uint64 {
	size := 8 * len(files) // 8 bytes for each File.Hash uint64 number

	buf := make([]byte, size)
	var i int // buf write index

	for _, file := range files {
		binary.LittleEndian.PutUint64(buf[i:], file.Hash)
		i += 8
	}

	h := fnv.New64a()
	h.Write(buf)
	return h.Sum64()
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

func formatHash(h uint64) string {
	// formats integer in hex with exactly displayed 8 bytes (16 characters)
	return fmt.Sprintf("%016x", h)
}

func formatCacheSeed(seed uint64) string {
	return formatHash(seed)
}
