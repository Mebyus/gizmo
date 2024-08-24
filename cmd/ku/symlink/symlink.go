package symlink

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/mebyus/gizmo/butler"
)

var Symlink = &butler.Lackey{
	Name: "symlink",

	Short: "create a symbolic link for ku",
	Usage: "ku symlink",

	Exec: execute,
}

func execute(r *butler.Lackey, _ []string) error {
	return symlink()
}

const symlinkPath = "/usr/local/bin/ku"

var ErrPermissionDeniedHint = fmt.Errorf("permission to \"%s\" denied, consider running command as root or with sudo", symlinkPath)

func symlink() error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	err = os.Remove(symlinkPath)
	if err != nil {
		switch {
		case errors.Is(err, fs.ErrPermission):
			return ErrPermissionDeniedHint
		case errors.Is(err, fs.ErrNotExist):
			// continue execution
		default:
			return err
		}
	}
	err = os.Symlink(execPath, symlinkPath)
	switch {
	case errors.Is(err, fs.ErrPermission):
		return ErrPermissionDeniedHint
	default:
		return err
	}
}
