package service

import "sort"

// SongItem is one row of the public rank table.
type SongItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	EZ     string `json:"EZ"`
	HD     string `json:"HD"`
	IN     string `json:"IN"`
	AT     string `json:"AT"`
	Legacy string `json:"Legacy"`
}

// RankTable returns the song/difficulty table, sorted by song id.
func (c *Client) RankTable() []SongItem {
	out := make([]SongItem, 0, len(c.difficulty))
	for id, d := range c.difficulty {
		out = append(out, SongItem{
			ID:     id,
			Title:  c.info[id]["song"],
			EZ:     d["EZ"],
			HD:     d["HD"],
			IN:     d["IN"],
			AT:     d["AT"],
			Legacy: d["Legacy"],
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
