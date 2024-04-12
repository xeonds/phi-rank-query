package model

import (
	"time"

	"github.com/xeonds/phi-plug-go/lib"
)

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
	ACL struct {
		NAMING_FAILED struct {
			Write bool `json:"write"`
			Read  bool `json:"read"`
		} `json:"*"`
	} `json:"ACL"`
	AuthData struct {
		Taptap struct {
			AccessToken  string `json:"access_token"`
			Avatar       string `json:"avatar"`
			Kid          string `json:"kid"`
			MacAlgorithm string `json:"mac_algorithm"`
			MacKey       string `json:"mac_key"`
			Name         string `json:"name"`
			Openid       string `json:"openid"`
			TokenType    string `json:"token_type"`
			Unionid      string `json:"unionid"`
		} `json:"taptap"`
	} `json:"authData"`
	Avatar              string    `json:"avatar"`
	CreatedAt           time.Time `json:"createdAt"`
	EmailVerified       bool      `json:"emailVerified"`
	MobilePhoneVerified bool      `json:"mobilePhoneVerified"`
	Nickname            string    `json:"nickname"`
	ObjectID            string    `json:"objectId"`
	SessionToken        string    `json:"sessionToken"`
	ShortID             string    `json:"shortId"`
	UpdatedAt           time.Time `json:"updatedAt"`
	Username            string    `json:"username"`
}
type GameProcess struct {
	IsFirstRun                 bool
	LegacyChapterFinished      bool
	AlreadyShowCollectionTip   bool
	AlreadyShowAutoUnlockINTip bool
	Completed                  string
	SongUpdateInfo             int
	ChallengeModeRank          int16
	Money                      [5]int
	UnlockFlagOfSpasmodic      byte
	UnlockFlagOfIgallta        byte
	UnlockFlagOfRrharil        byte
	FlagOfSongRecordKey        byte
	RandomVersionUnlocked      byte
	Chapter8UnlockBegin        bool
	Chapter8UnlockSecondPhase  bool
	Chapter8Passed             bool
	Chapter8SongUnlocked       byte
}
type LevelRecord struct {
	Score uint32
	Acc   float32
	Fc    bool
}
type GameUser struct {
	Name         string
	Version      int
	ShowPlayerId bool
	SelfIntro    string
	Avatar       string
	Background   string
}
type GameSettings struct {
	ChordSupport      bool
	FcAPIndicator     bool
	EnableHitSound    bool
	LowResolutionMode bool
	DeviceName        string
	Bright            float32
	MusicVolume       float32
	EffectVolume      float32
	HitSoundVolume    float32
	SoundOffset       float32
	NoteScale         float32
}
type GameRecord struct {
	Name     string
	Version  int
	Data     *lib.ByteReader
	Record   map[string][]*LevelRecord
	Songsnum int
}

// 存档信息
type Game struct {
	GameProgress *GameProcess
	GameUser     *GameUser
	GameSettings *GameSettings
	GameRecord   *GameRecord
}

// 单曲Rks信息
type Record struct {
	UserID uint32

	Id           string
	Rks          float64
	Score        uint32
	Difficulty   string
	Level        string
	Acc          float64
	FullCombo    bool
	Song         string
	Illustration string
}
