package service_test

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"testing"

	"github.com/xeonds/phi-plug-go/config"
	"github.com/xeonds/phi-plug-go/lib"
	"github.com/xeonds/phi-plug-go/model"
	"github.com/xeonds/phi-plug-go/service"
)

var session = os.Getenv("SESSION")
var conf = lib.LoadConfig[config.Config]()

func TestUserInfo(t *testing.T) {
	userData, err := service.GetuserInfo(conf, session)
	if err != nil {
		t.Fatalf("Failed to get user data: %v", err)
	}
	err = os.WriteFile("user.json", userData, 0644)
	if err != nil {
		t.Fatalf("Failed to write JSON data to file: %v", err)
	}
}

func TestB19(t *testing.T) {
	b19Data, err := service.GetB19Info(conf, session)
	if err != nil {
		t.Fatalf("Failed to get B19 data: %v", err)
	}
	err = os.WriteFile("file.json", b19Data, 0644)
	if err != nil {
		t.Fatalf("Failed to write JSON data to file: %v", err)
	}
}

func TestGetRankUrl(t *testing.T) {
	f, err := os.ReadFile("file.json")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	data := new(model.GameSave)
	if err = json.Unmarshal(f, &data); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}
	latestModified := data.Results[0]
	fmt.Println(latestModified.Gamefile.URL)
}

func TestGetZipAndDecrypt(t *testing.T) {
	f, err := os.ReadFile("file.json")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	data := new(model.GameSave)
	if err = json.Unmarshal(f, &data); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}
	latestModified := data.Results[0]
	savezip, err := service.GetSaveZip(conf, session, latestModified.Gamefile.URL)
	if err != nil {
		t.Fatalf("Failed to get save zip: %v", err)
	}

	decrypted := service.DecryptSaveZip(savezip)
	tmp, err := json.Marshal(decrypted)
	if err != nil {
		t.Fatalf("Failed to marshal decrypted data: %v", err)
	}
	err = os.WriteFile("decrypted.json", tmp, 0644)
	if err != nil {
		t.Fatalf("Failed to write decrypted data to file: %v", err)
	}
}

func TestReadDifficulty(t *testing.T) {
	difficulty, err := lib.LoadCSV("../resources/info/difficulty.csv")
	if err != nil {
		log.Fatal(err)
	}
	// dump difficulty to json
	tmp, err := json.Marshal(difficulty)
	if err != nil {
		t.Fatalf("Failed to marshal difficulty data: %v", err)
	}
	err = os.WriteFile("difficulty.json", tmp, 0644)
	if err != nil {
		t.Fatalf("Failed to write difficulty data to file: %v", err)
	}
}

func TestCalcRks(t *testing.T) {
	f, err := os.ReadFile("decrypted.json")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	data := new(model.Game)
	if err = json.Unmarshal(f, &data); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}
	config := lib.LoadConfig[config.Config]()
	b19, rks, phi := service.CalcBNInfo(data, config)
	tmp, err := json.Marshal(map[string]interface{}{
		"b19": b19,
		"rks": rks,
		"phi": phi,
	})
	if err != nil {
		t.Fatalf("Failed to marshal rks data: %v", err)
	}
	err = os.WriteFile("rks.json", tmp, 0644)
	if err != nil {
		t.Fatalf("Failed to write rks data to file: %v", err)
	}
}

func TestCalcBN(t *testing.T) {
	config := new(config.Config)
	config.Server.InsecureSkipVerify = true
	accountInfoDump, err := service.GetuserInfo(config, "9r5shbrhcz9omano8qx5ahn3f")
	if err != nil {
		t.Fatalf("Failed to get user data: %v", err)
	}
	accountInfo := new(model.GameAccount)
	_ = json.Unmarshal(accountInfoDump, accountInfo)
	userInfoDump, err := service.GetB19Info(config, "9r5shbrhcz9omano8qx5ahn3f")
	if err != nil {
		t.Fatalf("Failed to get B19 data: %v", err)
	}
	userInfo := new(model.GameSave)
	_ = json.Unmarshal(userInfoDump, userInfo)
	saveZip, err := service.GetSaveZip(config, "9r5shbrhcz9omano8qx5ahn3f", userInfo.Results[0].Gamefile.URL)
	if err != nil {
		t.Fatalf("Failed to get save zip: %v", err)
		return
	}
	game := service.DecryptSaveZip(saveZip)
	CalcBNInfo(game)
}

// 计算BN信息
func CalcBNInfo(data *model.Game) ([]model.Record, float64, model.Record) {
	phi := model.Record{}
	difficulty, err := lib.LoadCSV("../build/dist/difficulty.tsv")
	if err != nil {
		log.Fatal("reading difficulty: ", err)
	}
	songInfo, err := lib.LoadCSV("../build/dist/info.tsv")
	if err != nil {
		log.Fatal("reading songInfo: ", err)
	}
	comRks := 0.0
	phi.Rks = 0.0
	var rksList []model.Record
	for title, song := range data.GameRecord.Record {
		titleTrim := title[:len(title)-2]
		for level, tem := range song {
			difficulty_map := []string{"EZ", "HD", "IN", "AT", "Legacy"}
			// if level == 4 {
			// 	break
			// }
			if tem == nil {
				continue
			}
			// fix: difficulty数组访问下标错误导致定数数据获取失败的问题
			// 原因：疑似difficulty.csv下标变换导致map访问异常
			diff, err := strconv.ParseFloat(difficulty[title][difficulty_map[level]], 64)
			if err != nil {
				log.Println("parsing difficulty: ", err)
			}
			songRank := model.Record{
				Id:           title,
				Rks:          service.CalcSongRank(tem.Acc, diff),
				Score:        tem.Score,
				Difficulty:   difficulty[title][difficulty_map[level]],
				Level:        difficulty_map[level],
				Acc:          float64(tem.Acc),
				FullCombo:    tem.Fc,
				Song:         songInfo[title]["song"],
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

	for i := 0; i < len(rksList); i++ {
		if i < 19 {
			comRks += rksList[i].Rks
		}
		rksList[i].Rks = math.Floor(rksList[i].Rks*100) / 100
		rksList[i].Acc = math.Floor(rksList[i].Acc*100) / 100
	}

	return rksList, comRks / float64(20), phi
}

func getIllustration(song string) string {
	return "/assets/illustrations/" + song + ".png"
}
