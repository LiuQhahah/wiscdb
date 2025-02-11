package internal

type header struct {
	kLen      uint32
	vLen      uint32
	expiresAt uint64
	meta      byte
	userMeta  byte
}

const (
	maxHeaderSize = 22
)

func (h header) Encode(out []byte) int {
	return 0
}

func (h *header) Decode(buf []byte) int {
	return 0
}

func (h *header) DecodeFrom(reader *hashReader) (int, error) {
	return 0, nil
}
