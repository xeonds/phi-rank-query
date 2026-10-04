package model

import "github.com/xeonds/phi-plug-go/lib"

// Record is one chart result.
type Record struct {
	ID           string  `json:"id"`
	Song         string  `json:"song"`
	Level        string  `json:"level"`
	Difficulty   string  `json:"difficulty"`
	Rks          float64 `json:"rks"`
	Score        uint32  `json:"score"`
	Acc          float64 `json:"acc"`
	FullCombo    bool    `json:"fullCombo"`
	Illustration string  `json:"illustration"`
}

// --- LeanCloud save payloads ---

type GameSave struct {
	Results []struct {
		Createdat string `json:"createdat"`
		Gamefile  struct {
			Type      string `json:"__type"`
			Bucket    string `json:"bucket"`
			Createdat string `json:"createdat"`
			Key       string `json:"key"`
			Metadata  struct {
				Checksum string `json:"_checksum"`
				Prefix   string `json:"prefix"`
				Size     int    `json:"size"`
			} `json:"metadata"`
			MimeType  string `json:"mime_type"`
			Name      string `json:"name"`
			Objectid  string `json:"objectid"`
			Provider  string `json:"provider"`
			Updatedat string `json:"updatedat"`
			URL       string `json:"url"`
		} `json:"gamefile"`
		Modifiedat struct {
			Type string `json:"__type"`
			Iso  string `json:"iso"`
		} `json:"modifiedat"`
		Name      string `json:"name"`
		Objectid  string `json:"objectid"`
		Summary   string `json:"summary"`
		Updatedat string `json:"updatedat"`
		User      struct {
			Type      string `json:"__type"`
			Classname string `json:"classname"`
			Objectid  string `json:"objectid"`
		} `json:"user"`
	} `json:"results"`
}

type GameAccount struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type Game struct {
	GameProgress *GameProcess
	GameUser     *GameUser
	GameSettings *GameSettings
	GameRecord   *GameRecord
}

type GameProcess struct {
	Completed         string
	ChallengeModeRank int16
}

type GameUser struct {
	Name         string
	SelfIntro    string
	Avatar       string
	Background   string
	ShowPlayerID bool
}

type GameSettings struct {
	DeviceName string
}

type GameRecord struct {
	Songsnum int
	Data     *lib.ByteReader
	Record   map[string][]*LevelRecord
}

type LevelRecord struct {
	Score uint32
	Acc   float32
	Fc    bool
}
