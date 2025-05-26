package main

import (
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"hash/crc32"
	"os"
	"sync"
	"wiscdb/options"
	"wiscdb/pb"
	"wiscdb/y"
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
	if mf.inMemory {
		return nil
	}
	changes := pb.ManifestChangeSet{
		Changes: changeParam,
	}
	buf, err := proto.Marshal(&changes)
	if err != nil {
		return err
	}
	mf.appendLock.Lock()

	defer mf.appendLock.Unlock()

	if err := applyChangeSet(&mf.manifest, &changes); err != nil {
		return err
	}
	if mf.manifest.Deletions > mf.deletionsRewriteThreshold && mf.manifest.Deletions > manifestDeletionsRatio*(mf.manifest.Creations-mf.manifest.Deletions) {
		if err := mf.rewrite(); err != nil {
			return err
		}
	} else {
		var lenCrcBuf [8]byte
		binary.BigEndian.PutUint32(lenCrcBuf[0:4], uint32(len(buf)))
		binary.BigEndian.PutUint32(lenCrcBuf[4:8], crc32.Checksum(buf, y.CastTagNoLiCrcTable))
		buf = append(lenCrcBuf[:], buf...)
		if _, err := mf.fp.Write(buf); err != nil {
			return err
		}
	}

	// 将文件落盘
	return syncFunc(mf.fp)
}

// 将文件落盘
func syncFunc(f *os.File) error {
	return f.Sync()
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

func applyChangeSet(build *Manifest, changeSet *pb.ManifestChangeSet) error {
	for _, change := range changeSet.Changes {
		if err := applyManifestChange(build, change); err != nil {
			return err
		}
	}
	return nil
}

func applyManifestChange(build *Manifest, tc *pb.ManifestChange) error {
	switch tc.Op {
	case pb.ManifestChange_DELETE:
		tm, ok := build.Tables[tc.Id]
		if !ok {
			return fmt.Errorf("MANIFEST removes non-existing table %d", tc.Id)
		}
		delete(build.Levels[tm.Level].Tables, tc.Id)
		delete(build.Tables, tc.Id)
		build.Deletions++
	case pb.ManifestChange_CREATE:
		if _, ok := build.Tables[tc.Id]; ok {
			return fmt.Errorf("MANIFEST invalid, table %d exists", tc.Id)
		}
		build.Tables[tc.Id] = TableManifest{
			Level:       uint8(tc.Level),
			KeyID:       tc.KeyId,
			Compression: options.CompressionType(tc.Compression),
		}
		for len(build.Levels) <= int(tc.Level) {
			build.Levels = append(build.Levels, levelManifest{make(map[uint64]struct{})})
		}

		build.Levels[tc.Level].Tables[tc.Id] = struct{}{}
		build.Creations++
	default:
		return fmt.Errorf("MANIFEST file has invalid manifestChange op")
	}
	return nil
}
