package unity

import (
	"fmt"
	"strconv"
	"strings"
)

type BlockInfo struct {
	UncompressedSize uint32
	CompressedSize   uint32
	Flags            uint16
}

type Node struct {
	Offset int64
	Size   int64
	Flags  uint32
	Path   string
}

type Bundle struct {
	Signature     string
	Version       uint32
	VersionPlayer string
	VersionEngine string
	DataFlags     uint32
	Blocks        []BlockInfo
	Nodes         []Node
	data          []byte
}

// Data returns the decompressed concatenation of all blocks.
func (b *Bundle) Data() []byte { return b.data }

// NodeData returns the bytes of a single internal file.
func (b *Bundle) NodeData(n Node) []byte {
	return b.data[n.Offset : n.Offset+n.Size]
}

func parseUnityVersion(s string) (major, minor, patch int) {
	// e.g. "2022.3.62f2" or "5.x.x"
	parts := strings.SplitN(s, ".", 3)
	n := func(x string) int {
		// strip trailing letters/build info
		for i, c := range x {
			if c < '0' || c > '9' {
				x = x[:i]
				break
			}
		}
		v, _ := strconv.Atoi(x)
		return v
	}
	if len(parts) > 0 {
		major = n(parts[0])
	}
	if len(parts) > 1 {
		minor = n(parts[1])
	}
	if len(parts) > 2 {
		patch = n(parts[2])
	}
	return
}

// ParseBundle parses a UnityFS AssetBundle header, block info and file nodes.
func ParseBundle(r *Reader) (*Bundle, error) {
	r.SetBE(true)
	b := &Bundle{}
	b.Signature = r.CString()
	b.Version = r.U32()
	b.VersionPlayer = r.CString()
	b.VersionEngine = r.CString()

	switch b.Signature {
	case "UnityFS":
		// supported
	case "UnityWeb", "UnityRaw":
		return nil, fmt.Errorf("unity: bundle signature %q not supported", b.Signature)
	default:
		return nil, fmt.Errorf("unity: unknown bundle signature %q", b.Signature)
	}

	r.U64() // total size (unused)
	compressedSize := r.U32()
	uncompressedSize := r.U32()
	dataflags := r.U32()
	b.DataFlags = dataflags

	major, minor, patch := parseUnityVersion(b.VersionEngine)
	oldFlags := major < 2020 ||
		(major == 2020 && (minor < 3 || (minor == 3 && patch < 34))) ||
		(major == 2021 && (minor < 3 || (minor == 3 && patch < 2))) ||
		(major == 2022 && minor < 1)

	encFlag := uint32(0x1400)
	if oldFlags {
		encFlag = 0x200
	}
	if dataflags&encFlag != 0 {
		return nil, fmt.Errorf("unity: encrypted bundles (Unity CN) not supported")
	}
	if b.Version >= 7 || (major == 2019 && (minor > 4 || (minor == 4 && patch >= 15))) {
		r.Align(16)
	}

	start := r.Pos()
	var blocksInfo []byte
	if dataflags&0x80 != 0 { // BlocksInfoAtTheEnd
		r.Seek(r.Len() - int(compressedSize))
		blocksInfo = r.Read(int(compressedSize))
		r.Seek(start)
	} else {
		blocksInfo = r.Read(int(compressedSize))
	}

	dec, err := Decompress(dataflags, blocksInfo, int(uncompressedSize))
	if err != nil {
		return nil, fmt.Errorf("unity: blocks info: %w", err)
	}
	br := NewReader(dec)
	br.SetBE(true)
	br.Read(16) // uncompressed data hash
	blockCount := int(br.I32())
	for i := 0; i < blockCount; i++ {
		b.Blocks = append(b.Blocks, BlockInfo{br.U32(), br.U32(), br.U16()})
	}
	nodeCount := int(br.I32())
	for i := 0; i < nodeCount; i++ {
		b.Nodes = append(b.Nodes, Node{br.I64(), br.I64(), br.U32(), br.CString()})
	}

	if !oldFlags && dataflags&0x200 != 0 { // BlockInfoNeedPaddingAtStart
		r.Align(16)
	}

	var out []byte
	for i, blk := range b.Blocks {
		raw := r.Read(int(blk.CompressedSize))
		d, err := Decompress(uint32(blk.Flags), raw, int(blk.UncompressedSize))
		if err != nil {
			return nil, fmt.Errorf("unity: block %d: %w", i, err)
		}
		out = append(out, d...)
	}
	b.data = out
	return b, nil
}
