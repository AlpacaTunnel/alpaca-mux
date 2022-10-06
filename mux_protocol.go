package alpacamux

import (
	"encoding/binary"
)

const (
	HEADER_LEN = 16
	MAX_MTU    = 1450
	MAGIC      = 20202
)

type Header struct {
	Magic     uint32
	Length    uint16
	PathID    uint16
	Timestamp uint32
	Sequence  uint32
}

func (h *Header) FromNetwork(data []byte) {
	if len(data) < HEADER_LEN {
		h.Magic, h.Length, h.Timestamp = 0, 0, 0
		return
	}
	h.Magic = binary.BigEndian.Uint32(data[0:4])
	h.Length = binary.BigEndian.Uint16(data[4:6])
	h.PathID = binary.BigEndian.Uint16(data[6:8])
	h.Timestamp = binary.BigEndian.Uint32(data[8:12])
	h.Sequence = binary.BigEndian.Uint32(data[12:16])
}

func (h *Header) ToNetwork() []byte {
	data := make([]byte, HEADER_LEN)

	binary.BigEndian.PutUint32(data[0:4], h.Magic)
	binary.BigEndian.PutUint16(data[4:6], h.Length)
	binary.BigEndian.PutUint16(data[6:8], h.PathID)
	binary.BigEndian.PutUint32(data[8:12], h.Timestamp)
	binary.BigEndian.PutUint32(data[12:16], h.Sequence)

	return data
}
