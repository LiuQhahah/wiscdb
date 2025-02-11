package internal

import (
	"os"
)

type directoryLockGuard struct {
	r        *os.File
	path     string
	readOnly bool
}

func acquireDirectoryLock(dirPath string, pidFileName string, readOnly bool) (*directoryLockGuard, error) {
	return nil, nil
}

func (guard *directoryLockGuard) release() error {
	return nil
}

func openDir(path string) (*os.File, error) {

	return nil, nil
}

func syncDir(dir string) error {
	return nil
}
