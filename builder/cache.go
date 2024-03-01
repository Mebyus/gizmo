package builder

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func (cfg *Config) Hash() uint64 {
	return 53582900
}

type Cache struct {
	// Base directory for current cache instance. Includes build
	// config seed
	dir string

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

func formatCacheSeed(seed uint64) string {
	// formats integer in hex with exactly displayed 8 bytes (16 characters)
	return fmt.Sprintf("%016x", seed)
}
