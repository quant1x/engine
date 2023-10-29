package runtime

import (
	"os"
	"path/filepath"
)

func ApplicationName() string {
	path, _ := os.Executable()
	_, exec := filepath.Split(path)
	return exec
}
