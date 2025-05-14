package main

import (
	"os"
	"sync"
	"wiscdb/options"
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

func (mf *manifestFile) AddChanges(changeParam []*pb.ManifestChange) error {
	return nil
}

func NewCreateChange(id uint64, level int, keyID uint64, c options.CompressionType) *pb.ManifestChange {
	return &pb.ManifestChange{
		Id:             id,
		Op:             pb.ManifestChange_CREATE,
		Level:          uint32(level),
		KeyId:          keyID,
		EncryptionAlgo: pb.EncryptionAlgo_aes,
		Compression:    uint32(c),
	}
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
