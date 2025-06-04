package main

import (
	"os"
	"wiscdb/y"
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
	return os.Open(path)
}

func SyncDir(dir string) error {
	f, err := openDir(dir)
	if err != nil {
		return y.Wrapf(err, "While opening directory: %s.", dir)
	}
	err = f.Sync()
	closeErr := f.Close()
	if err != nil {
		return y.Wrapf(err, "While syncing directory: %s.", dir)
	}

	return y.Wrapf(closeErr, "While closing directory: %s.", dir)
}
