// Package phigros extracts Phigros game data from Unity assets.
package phigros

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"math"

	"github.com/xeonds/phi-plug-go/internal/unity"
)

// typetree.json is taken from 7aGiven/Phigros_Resource (GPLv3) and describes
// the GameInformation MonoBehaviour layout.
//
//go:embed typetree.json
var typetreeJSON []byte

// infoMarker locates the GameInformation blob by the first song id.
var infoMarker = []byte("Glaciaxion.SunsetRay.0")

var gameInfoNode = func() *unity.TypeTreeNode {
	var specs struct {
		GameInformation []unity.NodeSpec `json:"GameInformation"`
	}
	if err := json.Unmarshal(typetreeJSON, &specs); err != nil {
		panic(fmt.Sprintf("phigros: parse typetree.json: %v", err))
	}
	return unity.BuildTypeTree(specs.GameInformation)
}()

type DifficultyRow struct {
	ID    string
	Ranks []float64
}

type InfoRow struct {
	ID          string
	Song        string
	Composer    string
	Illustrator string
	Charters    []string
}

// ExtractGameInfo parses a decompressed data.unity3d bundle and returns the
// difficulty and info tables.
func ExtractGameInfo(data []byte) ([]DifficultyRow, []InfoRow, error) {
	b, err := unity.ParseBundle(unity.NewReader(data))
	if err != nil {
		return nil, nil, err
	}

	objBytes := findGameInfoObject(b)
	if objBytes == nil {
		return nil, nil, fmt.Errorf("phigros: GameInformation object not found")
	}

	v := unity.ReadTypeTree(gameInfoNode, unity.NewReader(objBytes))
	m, ok := v.(map[string]any)
	if !ok {
		return nil, nil, fmt.Errorf("phigros: unexpected GameInformation type %T", v)
	}
	songBase, ok := m["song"].(map[string]any)
	if !ok {
		return nil, nil, fmt.Errorf("phigros: missing song field")
	}

	var diffs []DifficultyRow
	var infos []InfoRow
	for _, group := range []string{"mainSongs", "extraSongs", "sideStorySongs"} {
		list, _ := songBase[group].([]any)
		for _, it := range list {
			item, ok := it.(map[string]any)
			if !ok {
				continue
			}
			id := asString(item["songsId"])
			if len(id) < 2 {
				continue
			}
			id = id[:len(id)-2]

			ranks := asFloats(item["difficulty"])
			charters := asStrings(item["charter"])
			if len(ranks) == 5 {
				ranks = ranks[:4]
			}
			if len(ranks) > 0 && ranks[len(ranks)-1] == 0 {
				ranks = ranks[:len(ranks)-1]
				if len(charters) > 0 {
					charters = charters[:len(charters)-1]
				}
			}
			for i := range ranks {
				ranks[i] = math.Round(ranks[i]*10) / 10
			}

			diffs = append(diffs, DifficultyRow{ID: id, Ranks: ranks})
			infos = append(infos, InfoRow{
				ID:          id,
				Song:        asString(item["songsName"]),
				Composer:    asString(item["composer"]),
				Illustrator: asString(item["illustrator"]),
				Charters:    charters,
			})
		}
	}
	return diffs, infos, nil
}

// findGameInfoObject locates the raw MonoBehaviour bytes containing the
// GameInformation marker.
func findGameInfoObject(b *unity.Bundle) []byte {
	for _, n := range b.Nodes {
		nd := b.NodeData(n)
		if bytes.Index(nd, infoMarker) < 0 {
			continue
		}
		sf, err := unity.ParseSerializedFile(unity.NewReader(nd))
		if err != nil {
			continue
		}
		for _, o := range sf.Objects {
			if o.ClassID != unity.ClassMonoBehaviour {
				continue
			}
			ob := sf.ObjectBytes(nd, o)
			if bytes.Index(ob, infoMarker) >= 0 {
				return ob
			}
		}
	}
	return nil
}

func asString(v any) string {
	s, _ := v.(string)
	return s
}

func asInt(v any) int {
	switch x := v.(type) {
	case int64:
		return int(x)
	case uint64:
		return int(x)
	case float64:
		return int(x)
	}
	return 0
}

func asFloats(v any) []float64 {
	list, _ := v.([]any)
	out := make([]float64, 0, len(list))
	for _, x := range list {
		switch n := x.(type) {
		case float64:
			out = append(out, n)
		case int64:
			out = append(out, float64(n))
		}
	}
	return out
}

func asStrings(v any) []string {
	switch list := v.(type) {
	case []string:
		return list
	case []any:
		out := make([]string, 0, len(list))
		for _, x := range list {
			if s, ok := x.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}
