package internal

import (
	"os"
	"sync"
	"wiscdb/pb"
)

type manifestFile struct {
	fp                        *os.File
	directory                 string
	externalMagic             uint16
	deletionsRewriteThreshold int
	appendLock                sync.Mutex
	manifest                  Manifest
	inMemory                  bool
}

const (
	ManifestFilename                 = "MANIFEST"
	manifestRewriteFilename          = "MANIFEST-REWRITE"
	manifestDeletionRewriteThreshold = 10000
	manifestDeletionsRatio           = 10
	wiscMagicVersion                 = 8
)

func (mf *manifestFile) close() error {
	return nil
}

func (mf *manifestFile) addChanges(changeParam []*pb.ManifestChange) error {
	return nil
}

func (mf *manifestFile) rewrite() error {
	return nil
}

func ReplayManifestFile(fp *os.File, extMagic uint16) (Manifest, int64, error) {
	return Manifest{}, 0, nil
}

func helpRewrite(dir string, m *Manifest, extMagic uint16) (*os.File, int, error) {
	return &os.File{}, 0, nil
}

func applyChangeSet(build *Manifest, changeSet *pb.ManifestChange) error {
	return nil
}

func applyManifestChange(build *Manifest, tc *pb.ManifestChange) error {
	return nil
}
