package lib

import (
	"encoding/base64"
	"encoding/binary"
	"math"
)

type ByteReader struct {
	data     []byte
	position int
}

func NewByteReader(data []byte) *ByteReader {
	return &ByteReader{
		data:     data,
		position: 0,
	}
}

func (br *ByteReader) Remaining() int {
	return len(br.data) - br.position - 1
}

func (br *ByteReader) GetByte() byte {
	b := br.data[br.position]
	br.position++
	return b
}

func (br *ByteReader) GetAllByte() string {
	return base64.StdEncoding.EncodeToString(br.data[br.position:])
}

func (br *ByteReader) GetShort() uint16 {
	br.position += 2
	return binary.LittleEndian.Uint16(br.data[br.position-2 : br.position])
}

func (br *ByteReader) GetInt() uint32 {
	br.position += 4
	return binary.LittleEndian.Uint32(br.data[br.position-4 : br.position])
}

func (br *ByteReader) GetFloat() float32 {
	br.position += 4
	return math.Float32frombits(binary.LittleEndian.Uint32(br.data[br.position-4 : br.position]))
}

func (br *ByteReader) GetVarInt() uint64 {
	if br.data[br.position] > 127 {
		br.position += 2
		return uint64(0b01111111&br.data[br.position-2]) ^ (uint64(br.data[br.position-1]) << 7)
	}
	result := br.data[br.position]
	br.position++
	return uint64(result)
}

func (br *ByteReader) SkipVarInt(num int) {
	if num > 0 {
		for ; num > 0; num-- {
			br.SkipVarInt(0)
		}
	} else {
		br.position++
	}
}

func (br *ByteReader) GetBytes() []byte {
	length := br.GetByte()
	br.position += int(length)
	return br.data[br.position-int(length) : br.position]
}

func (br *ByteReader) GetString() string {
	length := br.GetVarInt()
	end := br.position
	br.position += int(length)
	return string(br.data[end:br.position])
}

func (br *ByteReader) SkipString() {
	br.position += int(br.GetByte()) + 1
}

func (br *ByteReader) InsertBytes(bytes []byte) {
	result := make([]byte, len(br.data)+len(bytes))
	copy(result[:br.position], br.data[:br.position])
	copy(result[br.position:], bytes)
	copy(result[br.position+len(bytes):], br.data[br.position:])
	br.data = result
}

func (br *ByteReader) ReplaceBytes(length int, bytes []byte) {
	if len(bytes) == length {
		copy(br.data[br.position:], bytes)
		return
	}
	result := make([]byte, len(br.data)+len(bytes)-length)
	copy(result[:br.position], br.data[:br.position])
	copy(result[br.position:], bytes)
	copy(result[br.position+len(bytes):], br.data[br.position+length:])
	br.data = result
}
