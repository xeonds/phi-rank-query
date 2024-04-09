package service

import (
	"archive/zip"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"

	"github.com/xeonds/phi-plug-go/config"
	"github.com/xeonds/phi-plug-go/lib"
)

// api addresses
const BaseURL = "https://rak3ffdi.cloud.tds1.tapapis.cn/1.1"
const FileTokens = BaseURL + "/fileTokens"
const FileCallback = BaseURL + "/fileCallback"
const Save = BaseURL + "/classes/_GameSave"
const UserInfo = BaseURL + "/users/me"
const Files = BaseURL + "/files/"

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

// 获取用户信息
func GetuserInfo(session string) ([]byte, error) {
	req, err := http.NewRequest("GET", UserInfo, nil)
	if err != nil {
		return nil, err
	}
	setHeader(req)
	req.Header.Set("X-LC-Session", session)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// 获取存档信息
func GetB19Info(session string) ([]byte, error) {
	req, err := http.NewRequest("GET", Save, nil)
	if err != nil {
		return nil, err
	}
	setHeader(req)
	req.Header.Set("X-LC-Session", session)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// 获取存档文件
func GetSaveZip(saveURL string) (*zip.Reader, error) {
	resp, err := http.Get(saveURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch save zip: %s", resp.Status)
	}

	saveBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	saveZip, err := zip.NewReader(bytes.NewReader(saveBytes), int64(len(saveBytes)))
	if err != nil {
		return nil, err
	}

	return saveZip, nil
}

// 存档信息
type Game struct {
	GameProgress *GameProcess
	GameUser     *GameUser
	GameSettings *GameSettings
	GameRecord   *GameRecord
}

// 解析存档文件
func DecryptSaveZip(savezip *zip.Reader) *Game {
	gameProgressFile, err := savezip.Open("gameProgress")
	if err != nil {
		log.Fatal(err)
	}
	defer gameProgressFile.Close()
	gameProgressData, err := io.ReadAll(gameProgressFile)
	if err != nil {
		log.Fatal(err)
	}
	gameProgress := Decrypt(gameProgressData[1:])

	userFile, err := savezip.Open("user")
	if err != nil {
		log.Fatal(err)
	}
	defer userFile.Close()
	userData, err := io.ReadAll(userFile)
	if err != nil {
		log.Fatal(err)
	}
	gameuser := Decrypt(userData[1:])

	settingsFile, err := savezip.Open("settings")
	if err != nil {
		log.Fatal(err)
	}
	defer settingsFile.Close()
	settingsData, err := io.ReadAll(settingsFile)
	if err != nil {
		log.Fatal(err)
	}
	gamesettings := Decrypt(settingsData[1:])

	gameRecordFile, err := savezip.Open("gameRecord")
	if err != nil {
		log.Fatal(err)
	}
	defer gameRecordFile.Close()
	gameRecordData, err := io.ReadAll(gameRecordFile)
	if err != nil {
		log.Fatal(err)
	}
	record := Decrypt(gameRecordData[1:])

	return &Game{
		GameProgress: NewGameProcess(gameProgress),
		GameUser:     NewGameUser(gameuser),
		GameSettings: NewGameSettings(gamesettings),
		GameRecord:   NewGameRecord(record),
	}
}

// 单曲Rks信息
type Record struct {
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

// 计算BN信息
func CalcBNInfo(data *Game, config *config.Config, nnum int) ([]Record, float64, Record) {
	phi := Record{}
	difficulty, err := lib.LoadCSV(config.Data.Difficulty)
	if err != nil {
		log.Fatal("reading difficulty: ", err)
	}
	songInfo, err := lib.LoadCSV(config.Data.Info)
	if err != nil {
		log.Fatal("reading songInfo: ", err)
	}
	comRks := 0.0
	phi.Rks = 0.0
	var rksList []Record
	for title, song := range data.GameRecord.Record {
		titleTrim := title[:len(title)-2]
		for level, tem := range song {
			difficulty_map := []string{"EZ", "HD", "IN", "AT", "LEGACY"}
			if level == 4 {
				break
			}
			if tem == nil {
				continue
			}
			songRank := Record{
				Id:           titleTrim,
				Rks:          CalcSongRank(tem.Acc, difficulty[titleTrim][difficulty_map[level]]),
				Score:        tem.Score,
				Difficulty:   difficulty[titleTrim][difficulty_map[level]],
				Level:        difficulty_map[level],
				Acc:          float64(tem.Acc),
				FullCombo:    tem.Fc,
				Song:         songInfo[titleTrim]["song"],
				Illustration: getIllustration(titleTrim),
			}
			if tem.Acc >= 100 {
				if songRank.Rks > phi.Rks {
					phi.Id = titleTrim
					phi.Rks = songRank.Rks
					phi.Acc = songRank.Acc
					phi.Score = songRank.Score
					phi.Song = songRank.Song
					phi.Illustration = songRank.Illustration
					phi.Difficulty = songRank.Difficulty
					phi.FullCombo = songRank.FullCombo
					phi.Level = songRank.Level
				}
			}
			rksList = append(rksList, songRank)
		}
	}

	if phi.Rks != 0 {
		comRks += phi.Rks
		phi.Rks = math.Floor(phi.Rks*100) / 100
		phi.Acc = math.Floor(phi.Acc*100) / 100
	}

	var userRks float64

	minUpRks := math.Floor(userRks*100)/100 + 0.005 - userRks
	if minUpRks < 0 {
		minUpRks += 0.01
	}

	sort.Slice(rksList, func(i, j int) bool {
		return rksList[i].Rks > rksList[j].Rks
	})

	var b19List []Record
	for i := 0; i < nnum && i < len(rksList); i++ {
		if i < 19 {
			comRks += rksList[i].Rks
		}
		rksList[i].Rks = math.Floor(rksList[i].Rks*100) / 100
		rksList[i].Acc = math.Floor(rksList[i].Acc*100) / 100
		b19List = append(b19List, rksList[i])
	}

	return b19List, comRks / float64(nnum+1), phi
}

// 计算歌曲Rks
func CalcSongRank(acc float32, rank string) float64 {
	rankValue, _ := strconv.ParseFloat(rank, 64)
	if acc == 100 {
		return float64(rankValue)
	} else if acc < 70 {
		return 0
	} else {
		return rankValue * (((float64(acc) - 55) / 45) * ((float64(acc) - 55) / 45))
	}
}

func getIllustration(song string) string {
	return "/assets/illustrations/" + song + ".png"
}

// utils
func Decrypt(ciphertext []byte) []byte {
	key, _ := base64.StdEncoding.DecodeString("6Jaa0qVAJZuXkZCLiOa/Ax5tIZVu+taKUN1V1nqwkks=")
	iv, _ := base64.StdEncoding.DecodeString("Kk/wisgNYwcAV8WVGMgyUw==")

	decrypted := func(ciphertext, key, iv []byte) []byte {
		block, _ := aes.NewCipher(key)
		mode := cipher.NewCBCDecrypter(block, iv)

		plaintext := make([]byte, len(ciphertext))
		mode.CryptBlocks(plaintext, ciphertext)

		return plaintext
	}(ciphertext, key, iv)
	return decrypted
}
func setHeader(req *http.Request) {
	req.Header.Set("X-LC-Id", "rAK3FfdieFob2Nn8Am")
	req.Header.Set("X-LC-Key", "Qr9AEqtuoSVS3zeD6iVbM4ZC0AtkJcQ89tywVyi0")
	req.Header.Set("User-Agent", "LeanCloud-CSharp-SDK/1.0.3")
	req.Header.Set("Accept", "application/json")
}
func NewGameProcess(data []byte) *GameProcess {
	reader := lib.NewByteReader(data)
	gameProcess := &GameProcess{}
	tem := reader.GetByte()
	gameProcess.IsFirstRun = getBit(tem, 0)
	gameProcess.LegacyChapterFinished = getBit(tem, 1)
	gameProcess.AlreadyShowCollectionTip = getBit(tem, 2)
	gameProcess.AlreadyShowAutoUnlockINTip = getBit(tem, 3)
	gameProcess.Completed = reader.GetString()
	gameProcess.SongUpdateInfo = int(reader.GetVarInt())
	gameProcess.ChallengeModeRank = int16(reader.GetShort())
	for i := 0; i < 5; i++ {
		gameProcess.Money[i] = int(reader.GetVarInt())
	}
	gameProcess.UnlockFlagOfSpasmodic = reader.GetByte()
	gameProcess.UnlockFlagOfIgallta = reader.GetByte()
	gameProcess.UnlockFlagOfRrharil = reader.GetByte()
	gameProcess.FlagOfSongRecordKey = reader.GetByte()
	gameProcess.RandomVersionUnlocked = reader.GetByte()
	tem = reader.GetByte()
	gameProcess.Chapter8UnlockBegin = getBit(tem, 0)
	gameProcess.Chapter8UnlockSecondPhase = getBit(tem, 1)
	gameProcess.Chapter8Passed = getBit(tem, 2)
	gameProcess.Chapter8SongUnlocked = reader.GetByte()
	return gameProcess
}
func NewGameUser(data []byte) *GameUser {
	reader := lib.NewByteReader(data)
	gameUser := &GameUser{}
	gameUser.Name = "user"
	gameUser.Version = 1
	gameUser.ShowPlayerId = getBit(reader.GetByte(), 0)
	gameUser.SelfIntro = reader.GetString()
	gameUser.Avatar = reader.GetString()
	gameUser.Background = reader.GetString()
	return gameUser
}
func NewGameSettings(data []byte) *GameSettings {
	reader := lib.NewByteReader(data)
	gameSettings := &GameSettings{}
	tem := reader.GetByte()
	gameSettings.ChordSupport = getBit(tem, 0)
	gameSettings.FcAPIndicator = getBit(tem, 1)
	gameSettings.EnableHitSound = getBit(tem, 2)
	gameSettings.LowResolutionMode = getBit(tem, 3)
	gameSettings.DeviceName = reader.GetString()
	gameSettings.Bright = reader.GetFloat()
	gameSettings.MusicVolume = reader.GetFloat()
	gameSettings.EffectVolume = reader.GetFloat()
	gameSettings.HitSoundVolume = reader.GetFloat()
	gameSettings.SoundOffset = reader.GetFloat()
	gameSettings.NoteScale = reader.GetFloat()
	return gameSettings
}
func NewGameRecord(data []byte) *GameRecord {
	gameRecord := &GameRecord{
		Name:    "gameRecord",
		Version: 1,
		Data:    lib.NewByteReader(data),
		Record:  make(map[string][]*LevelRecord),
	}
	gameRecord.Songsnum = int(gameRecord.Data.GetVarInt())
	for gameRecord.Data.Remaining() > 32 {
		key := gameRecord.Data.GetString()
		gameRecord.Data.SkipVarInt(0)
		length := gameRecord.Data.GetByte()
		fc := gameRecord.Data.GetByte()
		song := make([]*LevelRecord, 5)

		for level := 0; level < 5; level++ {
			if getBit(length, uint(level)) {
				song[level] = &LevelRecord{}
				song[level].Score = gameRecord.Data.GetInt()
				song[level].Acc = gameRecord.Data.GetFloat()
				song[level].Fc = getBit(fc, uint(level))
			}
		}

		gameRecord.Record[key] = song
	}
	return gameRecord
}
func getBit(data byte, index uint) bool {
	return (data & (1 << index)) != 0
}
