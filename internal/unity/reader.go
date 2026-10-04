package unity

import (
	"encoding/binary"
	"fmt"
	"math"
)

// Reader is a little/big-endian cursor over a byte slice, modelled after
// UnityPy's EndianBinaryReader. Unity data is little-endian by default.
type Reader struct {
	b  []byte
	p  int
	be bool
}

func NewReader(b []byte) *Reader { return &Reader{b: b} }

func (r *Reader) Len() int        { return len(r.b) }
func (r *Reader) Pos() int        { return r.p }
func (r *Reader) Remaining() int  { return len(r.b) - r.p }
func (r *Reader) Seek(p int)      { r.p = p }
func (r *Reader) SetBE(be bool)   { r.be = be }
func (r *Reader) BigEndian() bool { return r.be }

func (r *Reader) order() binary.ByteOrder {
	if r.be {
		return binary.BigEndian
	}
	return binary.LittleEndian
}

func (r *Reader) need(n int) {
	if n < 0 || r.p+n > len(r.b) {
		panic(fmt.Sprintf("unity: read out of range pos=%d n=%d len=%d", r.p, n, len(r.b)))
	}
}

func (r *Reader) Read(n int) []byte {
	r.need(n)
	s := r.b[r.p : r.p+n]
	r.p += n
	return s
}

func (r *Reader) U8() uint8 {
	r.need(1)
	v := r.b[r.p]
	r.p++
	return v
}

func (r *Reader) Bool() bool { return r.U8() != 0 }

func (r *Reader) U16() uint16 {
	r.need(2)
	v := r.order().Uint16(r.b[r.p:])
	r.p += 2
	return v
}

func (r *Reader) U32() uint32 {
	r.need(4)
	v := r.order().Uint32(r.b[r.p:])
	r.p += 4
	return v
}

func (r *Reader) U64() uint64 {
	r.need(8)
	v := r.order().Uint64(r.b[r.p:])
	r.p += 8
	return v
}

func (r *Reader) I8() int8   { return int8(r.U8()) }
func (r *Reader) I16() int16 { return int16(r.U16()) }
func (r *Reader) I32() int32 { return int32(r.U32()) }
func (r *Reader) I64() int64 { return int64(r.U64()) }

func (r *Reader) F32() float32 { return math.Float32frombits(r.U32()) }
func (r *Reader) F64() float64 { return math.Float64frombits(r.U64()) }

// CString reads a NUL-terminated UTF-8 string.
func (r *Reader) CString() string {
	start := r.p
	for r.p < len(r.b) && r.b[r.p] != 0 {
		r.p++
	}
	s := string(r.b[start:r.p])
	if r.p < len(r.b) {
		r.p++
	}
	return s
}

// Align advances to the next multiple of n.
func (r *Reader) Align(n int) {
	for r.p%n != 0 {
		r.p++
		r.need(0)
	}
}
