package unity

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/pierrec/lz4/v4"
	"github.com/ulikunitz/xz/lzma"
)

// Decompress decompresses UnityFS block data given a compression flag
// (block flags or archive data flags; only the low 6 bits matter).
func Decompress(flag uint32, data []byte, uncompressedSize int) ([]byte, error) {
	switch flag & 0x3F {
	case 0:
		return data, nil
	case 2, 3: // LZ4 / LZ4HC share the block format
		dst := make([]byte, uncompressedSize)
		n, err := lz4.UncompressBlock(data, dst)
		if err != nil {
			return nil, fmt.Errorf("unity: lz4: %w", err)
		}
		return dst[:n], nil
	case 1: // LZMA (Unity stores raw LZMA1 with a 5-byte props header)
		return decompressLZMA(data)
	default:
		return nil, fmt.Errorf("unity: unsupported compression flag %d", flag&0x3F)
	}
}

func decompressLZMA(data []byte) ([]byte, error) {
	if len(data) < 5 {
		return nil, fmt.Errorf("unity: lzma data too short")
	}
	dictSize := binary.LittleEndian.Uint32(data[1:5])
	// Rebuild a classic .lzma header: props + dict size + unknown uncompressed size.
	hdr := make([]byte, 13)
	hdr[0] = data[0]
	copy(hdr[1:5], data[1:5])
	for i := 5; i < 13; i++ {
		hdr[i] = 0xFF
	}
	rc := lzma.ReaderConfig{DictCap: int(dictSize)}
	r, err := rc.NewReader(bytes.NewReader(append(hdr, data[5:]...)))
	if err != nil {
		return nil, fmt.Errorf("unity: lzma: %w", err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("unity: lzma: %w", err)
	}
	return out, nil
}
