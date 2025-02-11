package y

import "io"

type PageBuffer struct {
	pages        []*page
	length       int
	nextPageSize int
}

func NewPageBuffer(pageSize int) *PageBuffer {
	return &PageBuffer{}
}

func (b *PageBuffer) Write(data []byte) (int, error) {
	return 0, nil
}

func (b *PageBuffer) WriteByte(data byte) error {
	return nil
}

func (b *PageBuffer) Len() int {
	return 0
}

func (b *PageBuffer) pageForOffset(offset int) (int, int) {
	return 0, 0
}

func (b *PageBuffer) Truncate(n int) {

}

func (b *PageBuffer) Bytes() []byte {
	return nil
}

func (b *PageBuffer) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}

func (b *PageBuffer) NewReaderAt(offset int) *PageBufferReader {
	return &PageBufferReader{}
}
