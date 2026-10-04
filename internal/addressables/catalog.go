// Package addressables parses the legacy Unity Addressables binary catalog
// (catalog.json with m_KeyDataString/m_BucketDataString/m_EntryDataString).
package addressables

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type Entry struct {
	// Key is the addressable key, with any "Assets/Tracks/" prefix stripped.
	Key string
	// Bundle is the internal bundle file name ("<hash>_<name>.bundle" stripped
	// to "<name>.bundle").
	Bundle string
}

type bucketReader struct {
	b []byte
	p int
}

func (r *bucketReader) readInt() int {
	r.p += 4
	return int(r.b[r.p-4]) ^ int(r.b[r.p-3])<<8 ^ int(r.b[r.p-2])<<16
}

// Parse decodes the three base64 blobs into a key -> bundle table.
// Mirrors Phigros_Resource/resource.py's catalog handling.
func Parse(keyData, bucketData, entryData string) ([]Entry, error) {
	key, err := base64.StdEncoding.DecodeString(keyData)
	if err != nil {
		return nil, fmt.Errorf("addressables: key data: %w", err)
	}
	bucket, err := base64.StdEncoding.DecodeString(bucketData)
	if err != nil {
		return nil, fmt.Errorf("addressables: bucket data: %w", err)
	}
	entry, err := base64.StdEncoding.DecodeString(entryData)
	if err != nil {
		return nil, fmt.Errorf("addressables: entry data: %w", err)
	}

	br := &bucketReader{b: bucket}
	count := br.readInt()
	keys := make([]string, 0, count)
	idx := make([]int, 0, count)
	for x := 0; x < count; x++ {
		keyPos := br.readInt()
		keyType := key[keyPos]
		keyPos++
		var keyValue string
		switch keyType {
		case 0:
			length := int(key[keyPos])
			keyPos += 4
			keyValue = string(key[keyPos : keyPos+length])
		case 1:
			length := int(key[keyPos])
			keyPos += 4
			keyValue = utf16le(key[keyPos : keyPos+length])
		case 4:
			keyValue = fmt.Sprintf("%d", key[keyPos])
		default:
			return nil, fmt.Errorf("addressables: unknown key type %d", keyType)
		}
		entryValue := 0
		inner := br.readInt()
		for i := 0; i < inner; i++ {
			entryPos := br.readInt()
			off := 4 + 28*entryPos
			e := entry[off : off+28]
			entryValue = int(e[8]) ^ int(e[9])<<8
		}
		keys = append(keys, keyValue)
		idx = append(idx, entryValue)
	}

	// Resolve each entry's bundle to the key of the indexed row.
	bundles := make([]string, len(keys))
	for i := range keys {
		if idx[i] != 65535 && idx[i] >= 0 && idx[i] < len(keys) {
			bundles[i] = keys[idx[i]]
		}
	}

	out := make([]Entry, 0, len(keys))
	for i, k := range keys {
		b := bundles[i]
		if b == "" {
			continue
		}
		if strings.HasPrefix(k, "Assets/Tracks/") {
			k = k[len("Assets/Tracks/"):]
		} else if !strings.HasPrefix(k, "avatar.") {
			continue
		}
		if p := strings.IndexByte(b, '_'); p >= 0 {
			b = b[p+1:]
		}
		out = append(out, Entry{Key: k, Bundle: b})
	}
	return out, nil
}

func utf16le(b []byte) string {
	out := make([]rune, 0, len(b)/2)
	for i := 0; i+1 < len(b); i += 2 {
		out = append(out, rune(uint16(b[i])|uint16(b[i+1])<<8))
	}
	return string(out)
}
