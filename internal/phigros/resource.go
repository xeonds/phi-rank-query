package phigros

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/xeonds/phi-plug-go/internal/unity"
)

// DecodeTexture2D finds the Texture2D in a decompressed UnityFS bundle and
// decodes its pixel data into an image. Only uncompressed formats used by
// Phigros illustration assets are supported.
func DecodeTexture2D(bundleData []byte) (image.Image, error) {
	b, err := unity.ParseBundle(unity.NewReader(bundleData))
	if err != nil {
		return nil, err
	}

	var tex map[string]any
	var resS []byte
	for _, n := range b.Nodes {
		nd := b.NodeData(n)
		if strings.HasSuffix(n.Path, ".resS") || strings.HasSuffix(n.Path, ".resource") {
			resS = nd
			continue
		}
		sf, err := unity.ParseSerializedFile(unity.NewReader(nd))
		if err != nil {
			continue
		}
		for _, o := range sf.Objects {
			if o.ClassID != unity.ClassTexture2D {
				continue
			}
			st := sf.Types[o.TypeID]
			if st.Node == nil {
				continue
			}
			v := unity.ReadTypeTree(st.Node, unity.NewReader(sf.ObjectBytes(nd, o)))
			if m, ok := v.(map[string]any); ok {
				tex = m
			}
		}
	}
	if tex == nil {
		return nil, fmt.Errorf("phigros: Texture2D not found")
	}

	w := asInt(tex["m_Width"])
	h := asInt(tex["m_Height"])
	format := asInt(tex["m_TextureFormat"])

	var pixels []byte
	if stream, ok := tex["m_StreamData"].(map[string]any); ok {
		size := asInt(stream["size"])
		offset := asInt(stream["offset"])
		if size > 0 {
			if resS == nil || offset+size > len(resS) {
				return nil, fmt.Errorf("phigros: streamed texture data unavailable")
			}
			pixels = resS[offset : offset+size]
		}
	}
	if pixels == nil {
		if raw, ok := tex["image data"].([]byte); ok {
			pixels = raw
		}
	}
	if pixels == nil {
		return nil, fmt.Errorf("phigros: empty texture data")
	}

	return decodePixels(pixels, w, h, format)
}

func decodePixels(data []byte, w, h, format int) (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	n := w * h
	switch format {
	case 3: // RGB24
		if len(data) < n*3 {
			return nil, fmt.Errorf("phigros: short RGB24 data")
		}
		for i := 0; i < n; i++ {
			setRow(img, i, w, h, data[i*3], data[i*3+1], data[i*3+2], 255)
		}
	case 4: // RGBA32
		if len(data) < n*4 {
			return nil, fmt.Errorf("phigros: short RGBA32 data")
		}
		for i := 0; i < n; i++ {
			setRow(img, i, w, h, data[i*4], data[i*4+1], data[i*4+2], data[i*4+3])
		}
	case 5: // ARGB32
		if len(data) < n*4 {
			return nil, fmt.Errorf("phigros: short ARGB32 data")
		}
		for i := 0; i < n; i++ {
			setRow(img, i, w, h, data[i*4+1], data[i*4+2], data[i*4+3], data[i*4])
		}
	case 1: // Alpha8
		if len(data) < n {
			return nil, fmt.Errorf("phigros: short Alpha8 data")
		}
		for i := 0; i < n; i++ {
			setRow(img, i, w, h, 255, 255, 255, data[i])
		}
	case 7: // RGB565
		if len(data) < n*2 {
			return nil, fmt.Errorf("phigros: short RGB565 data")
		}
		for i := 0; i < n; i++ {
			v := uint16(data[i*2]) | uint16(data[i*2+1])<<8
			r := uint8((v>>11)&0x1F) << 3
			g := uint8((v>>5)&0x3F) << 2
			b := uint8(v&0x1F) << 3
			setRow(img, i, w, h, r, g, b, 255)
		}
	default:
		return nil, fmt.Errorf("phigros: unsupported texture format %d", format)
	}
	return img, nil
}

// setRow writes pixel i, flipping vertically so the image is top-down.
func setRow(img *image.RGBA, i, w, h int, r, g, b, a uint8) {
	x := i % w
	y := h - 1 - i/w
	img.SetRGBA(x, y, color.RGBA{r, g, b, a})
}
