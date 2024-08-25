package env

import (
	"fmt"
	"os"
	"path/filepath"
)

// RootDir returns compiler root directory.
func RootDir() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	root := filepath.Join(path, "./../../..")
	if root == "." || root == "/" || root == "" {
		return "", fmt.Errorf("bad root dir \"%s\"", root)
	}
	return root, nil
}
