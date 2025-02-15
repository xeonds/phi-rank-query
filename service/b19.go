package service

import (
	"archive/zip"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
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
	"github.com/xeonds/phi-plug-go/model"
)

// api addresses
const BaseURL = "https://rak3ffdi.cloud.tds1.tapapis.cn/1.1"
const FileTokens = BaseURL + "/fileTokens"
const FileCallback = BaseURL + "/fileCallback"
const Save = BaseURL + "/classes/_GameSave"
const UserInfo = BaseURL + "/users/me"
const Files = BaseURL + "/files/"

// 获取用户信息
func GetuserInfo(config *config.Config, session string) ([]byte, error) {
	req, err := http.NewRequest("GET", UserInfo, nil)
	if err != nil {
		return nil, err
	}
	setHeader(req)
	req.Header.Set("X-LC-Session", session)

	// 修复因为不信任证书导致的无法访问
	resp, err := (&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.Server.InsecureSkipVerify,
			},
		},
	}).Do(req)
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
func GetB19Info(config *config.Config, session string) ([]byte, error) {
	req, err := http.NewRequest("GET", Save, nil)
	if err != nil {
		return nil, err
	}
	setHeader(req)
	req.Header.Set("X-LC-Session", session)

	// 修复因为不信任证书导致的无法访问
	resp, err := (&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.Server.InsecureSkipVerify,
			},
		},
	}).Do(req)
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
func GetSaveZip(config *config.Config, session, saveURL string) (*zip.Reader, error) {
	req, err := http.NewRequest("GET", saveURL, nil)
	if err != nil {
		return nil, err
	}

	// 修复因为不信任证书导致的无法访问
	resp, err := (&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.Server.InsecureSkipVerify,
			},
		},
	}).Do(req)
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

// 解析存档文件
func DecryptSaveZip(savezip *zip.Reader) *model.Game {
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

	return &model.Game{
		GameProgress: NewGameProcess(gameProgress),
		GameUser:     NewGameUser(gameuser),
		GameSettings: NewGameSettings(gamesettings),
		GameRecord:   NewGameRecord(record),
	}
}

// 计算BN信息
func CalcBNInfo(data *model.Game, config *config.Config) ([]model.Record, float64, []model.Record) {
	difficulty, err := lib.LoadCSV(config.Data.Difficulty)
	if err != nil {
		log.Fatal("reading difficulty: ", err)
	}
	songInfo, err := lib.LoadCSV(config.Data.Info)
	if err != nil {
		log.Fatal("reading songInfo: ", err)
	}
	comRks := 0.0
	var rksList []model.Record
	for title, song := range data.GameRecord.Record {
		// fix: difficulty数组访问下标错误导致定数数据获取失败的问题
		// 原因：疑似difficulty.csv下标变换导致map访问异常
		titleTrim := title[:len(title)-2]
		for level, tem := range song {
			difficulty_map := []string{"EZ", "HD", "IN", "AT", "Legacy"}
			if level == 4 {
				break
			}
			if tem == nil {
				continue
			}
			diff, err := strconv.ParseFloat(difficulty[titleTrim][difficulty_map[level]], 64)
			if err != nil {
				log.Println("parsing difficulty: ", err)
			}
			songRank := model.Record{
				Id:           titleTrim,
				Rks:          CalcSongRank(tem.Acc, diff),
				Score:        tem.Score,
				Difficulty:   difficulty[titleTrim][difficulty_map[level]],
				Level:        difficulty_map[level],
				Acc:          float64(tem.Acc),
				FullCombo:    tem.Fc,
				Song:         songInfo[titleTrim]["song"],
				Illustration: getIllustration(titleTrim),
			}
			rksList = append(rksList, songRank)
		}
	}

	sort.Slice(rksList, func(i, j int) bool {
		return rksList[i].Rks > rksList[j].Rks
	})

	var p3 []model.Record
	for _, record := range rksList {
		if record.Acc == 100 {
			p3 = append(p3, record)
			comRks += record.Rks
			if len(p3) == 3 {
				break
			}
		}
	}

	var userRks float64

	minUpRks := math.Floor(userRks*100)/100 + 0.005 - userRks
	if minUpRks < 0 {
		minUpRks += 0.01
	}

	for i := 0; i < len(rksList); i++ {
		if i < 27 {
			comRks += rksList[i].Rks
		}
		rksList[i].Rks = math.Floor(rksList[i].Rks*100) / 100
		rksList[i].Acc = math.Floor(rksList[i].Acc*100) / 100
	}

	return rksList, comRks / float64(30), p3
}

// 计算歌曲Rks
func CalcSongRank(acc float32, rank float64) float64 {
	if acc == 100 {
		return float64(rank)
	} else if acc < 70 {
		return 0
	} else {
		return rank * (((float64(acc) - 55) / 45) * ((float64(acc) - 55) / 45))
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
func NewGameProcess(data []byte) *model.GameProcess {
	reader := lib.NewByteReader(data)
	gameProcess := &model.GameProcess{}
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
func NewGameUser(data []byte) *model.GameUser {
	reader := lib.NewByteReader(data)
	gameUser := &model.GameUser{}
	gameUser.Name = "user"
	gameUser.Version = 1
	gameUser.ShowPlayerId = getBit(reader.GetByte(), 0)
	gameUser.SelfIntro = reader.GetString()
	gameUser.Avatar = reader.GetString()
	gameUser.Background = reader.GetString()
	return gameUser
}
func NewGameSettings(data []byte) *model.GameSettings {
	reader := lib.NewByteReader(data)
	gameSettings := &model.GameSettings{}
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
func NewGameRecord(data []byte) *model.GameRecord {
	gameRecord := &model.GameRecord{
		Name:    "gameRecord",
		Version: 1,
		Data:    lib.NewByteReader(data),
		Record:  make(map[string][]*model.LevelRecord),
	}
	gameRecord.Songsnum = int(gameRecord.Data.GetVarInt())
	for gameRecord.Data.Remaining() > 32 {
		key := gameRecord.Data.GetString()
		gameRecord.Data.SkipVarInt(0)
		length := gameRecord.Data.GetByte()
		fc := gameRecord.Data.GetByte()
		song := make([]*model.LevelRecord, 5)

		for level := 0; level < 5; level++ {
			if getBit(length, uint(level)) {
				song[level] = &model.LevelRecord{}
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
