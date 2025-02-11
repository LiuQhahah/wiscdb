package y

type PageBufferReader struct {
	buf      *PageBuffer
	pageIdx  int
	startIdx int
}

func (r *PageBufferReader) Read(p []byte) (int, error) {

	return 0, nil
}
