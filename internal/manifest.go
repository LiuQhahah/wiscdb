package internal

import "wiscdb/pb"

type Manifest struct {
	Levels    []levelManifest
	Tables    map[uint64]TableManifest
	Creations int
	Deletions int
}

func createManifest() Manifest {
	return Manifest{}
}

func (m *Manifest) asChanges() []*pb.ManifestChange {
	return []*pb.ManifestChange{}
}

func (m *Manifest) clone() Manifest {
	return Manifest{}
}

func openOrCreateManifestFile(opt Options) (ret *manifestFile, result Manifest, err error) {
	return &manifestFile{}, Manifest{}, err
}

func helpOpenOrCreateManifestFile(dir string, readOnly bool, extMagic uint16, deletionsThreshold int) (*manifestFile, Manifest, error) {
	return &manifestFile{}, Manifest{}, nil
}
