package service

import (
	"archive/zip"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strconv"

	"github.com/xeonds/phi-plug-go/config"
	"github.com/xeonds/phi-plug-go/lib"
	"github.com/xeonds/phi-plug-go/model"
)

const (
	baseURL = "https://rak3ffdi.cloud.tds1.tapapis.cn/1.1"
	urlUser = baseURL + "/users/me"
	urlSave = baseURL + "/classes/_GameSave"
)

var difficultyLevels = []string{"EZ", "HD", "IN", "AT", "Legacy"}

// PlayerResult is the response of a player query.
type PlayerResult struct {
	Player    string         `json:"player"`
	Rks       float64        `json:"rks"`
	Challenge int            `json:"challenge"`
	Date      string         `json:"date"`
	Records   []model.Record `json:"records"`
	Phi       []model.Record `json:"phi"`
}

// Client talks to LeanCloud and computes rks from a session token.
type Client struct {
	cfg        *config.Config
	http       *http.Client
	difficulty map[string]map[string]string
	info       map[string]map[string]string
}

func New(cfg *config.Config) (*Client, error) {
	difficulty, err := lib.LoadTSV(cfg.Data.Difficulty)
	if err != nil {
		return nil, fmt.Errorf("load difficulty: %w", err)
	}
	info, err := lib.LoadTSV(cfg.Data.Info)
	if err != nil {
		return nil, fmt.Errorf("load song info: %w", err)
	}
	return &Client{
		cfg: cfg,
		http: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.Server.InsecureSkipVerify},
		}},
		difficulty: difficulty,
		info:       info,
	}, nil
}

func (c *Client) get(url, session string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	setHeader(req)
	req.Header.Set("X-LC-Session", session)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("leancloud status %d: %s", resp.StatusCode, body)
	}
	return body, nil
}

// Query fetches the player's save and computes rks + records.
func (c *Client) Query(session string) (*PlayerResult, error) {
	accountDump, err := c.get(urlUser, session)
	if err != nil {
		return nil, err
	}
	account := new(model.GameAccount)
	if err := json.Unmarshal(accountDump, account); err != nil {
		return nil, err
	}

	saveDump, err := c.get(urlSave, session)
	if err != nil {
		return nil, err
	}
	save := new(model.GameSave)
	if err := json.Unmarshal(saveDump, save); err != nil {
		return nil, err
	}
	if len(save.Results) == 0 {
		return nil, fmt.Errorf("no save found")
	}
	// LeanCloud returns saves oldest-first; pick the most recently modified.
	sort.Slice(save.Results, func(i, j int) bool {
		return save.Results[i].Modifiedat.Iso > save.Results[j].Modifiedat.Iso
	})
	latest := save.Results[0]

	saveZip, err := c.downloadZip(latest.Gamefile.URL)
	if err != nil {
		return nil, err
	}
	game, err := decryptSaveZip(saveZip)
	if err != nil {
		return nil, err
	}

	records, rks, phi := c.calc(game)
	return &PlayerResult{
		Player:    account.Nickname,
		Rks:       rks,
		Challenge: int(game.GameProgress.ChallengeModeRank),
		Date:      latest.Gamefile.Updatedat,
		Records:   records,
		Phi:       phi,
	}, nil
}

func (c *Client) downloadZip(url string) (*zip.Reader, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download save: %s", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return zip.NewReader(bytes.NewReader(data), int64(len(data)))
}

// calc mirrors the original CalcBNInfo: records sorted desc by rks, phi is the
// top 3 charts with 100% accuracy, rks is the average over the top 27 plus phi.
func (c *Client) calc(game *model.Game) ([]model.Record, float64, []model.Record) {
	var records []model.Record
	for title, song := range game.GameRecord.Record {
		if len(title) < 2 {
			continue
		}
		titleTrim := title[:len(title)-2]
		for level, tem := range song {
			if level == 4 || tem == nil {
				continue
			}
			diff, err := strconv.ParseFloat(c.difficulty[titleTrim][difficultyLevels[level]], 64)
			if err != nil {
				continue
			}
			records = append(records, model.Record{
				ID:           titleTrim,
				Song:         c.info[titleTrim]["song"],
				Level:        difficultyLevels[level],
				Difficulty:   c.difficulty[titleTrim][difficultyLevels[level]],
				Rks:          CalcSongRank(tem.Acc, diff),
				Score:        tem.Score,
				Acc:          float64(tem.Acc),
				FullCombo:    tem.Fc,
				Illustration: "/assets/illustrations/" + titleTrim + ".png",
			})
		}
	}
	sort.Slice(records, func(i, j int) bool { return records[i].Rks > records[j].Rks })

	var phi []model.Record
	sum := 0.0
	for _, r := range records {
		if r.Acc == 100 {
			phi = append(phi, r)
			sum += r.Rks
			if len(phi) == 3 {
				break
			}
		}
	}
	for i := range records {
		if i < 27 {
			sum += records[i].Rks
		}
		records[i].Rks = math.Floor(records[i].Rks*100) / 100
		records[i].Acc = math.Floor(records[i].Acc*100) / 100
	}
	return records, sum / 30, phi
}

// CalcSongRank is the single-chart rks formula.
func CalcSongRank(acc float32, rank float64) float64 {
	switch {
	case acc == 100:
		return rank
	case acc < 70:
		return 0
	default:
		x := (float64(acc) - 55) / 45
		return rank * x * x
	}
}

func setHeader(req *http.Request) {
	req.Header.Set("X-LC-Id", "rAK3FfdieFob2Nn8Am")
	req.Header.Set("X-LC-Key", "Qr9AEqtuoSVS3zeD6iVbM4ZC0AtkJcQ89tywVyi0")
	req.Header.Set("User-Agent", "LeanCloud-CSharp-SDK/1.0.3")
	req.Header.Set("Accept", "application/json")
}

// --- save decryption ---

func decryptSaveZip(savezip *zip.Reader) (*model.Game, error) {
	read := func(name string) ([]byte, error) {
		f, err := savezip.Open(name)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return io.ReadAll(f)
	}
	progress, err := read("gameProgress")
	if err != nil {
		return nil, err
	}
	record, err := read("gameRecord")
	if err != nil {
		return nil, err
	}
	return &model.Game{
		GameProgress: parseGameProgress(decrypt(progress[1:])),
		GameRecord:   parseGameRecord(decrypt(record[1:])),
	}, nil
}

func decrypt(ciphertext []byte) []byte {
	key, _ := base64.StdEncoding.DecodeString("6Jaa0qVAJZuXkZCLiOa/Ax5tIZVu+taKUN1V1nqwkks=")
	iv, _ := base64.StdEncoding.DecodeString("Kk/wisgNYwcAV8WVGMgyUw==")
	block, _ := aes.NewCipher(key)
	plain := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plain, ciphertext)
	// Strip PKCS#7 padding so parsers can rely on an exact byte length.
	if n := int(plain[len(plain)-1]); n > 0 && n <= aes.BlockSize && n <= len(plain) {
		plain = plain[:len(plain)-n]
	}
	return plain
}

func parseGameProgress(data []byte) *model.GameProcess {
	r := lib.NewByteReader(data)
	p := &model.GameProcess{}
	r.GetByte() // flags
	p.Completed = r.GetString()
	r.GetVarInt() // songUpdateInfo
	p.ChallengeModeRank = int16(r.GetShort())
	for i := 0; i < 5; i++ {
		r.GetVarInt() // money
	}
	r.GetByte() // unlockFlagOfSpasmodic
	r.GetByte() // unlockFlagOfIgallta
	r.GetByte() // unlockFlagOfRrharil
	r.GetByte() // flagOfSongRecordKey
	r.GetByte() // randomVersionUnlocked
	r.GetByte() // chapter8 unlock bits
	r.GetByte() // chapter8SongUnlocked
	return p
}

func parseGameRecord(data []byte) *model.GameRecord {
	gr := &model.GameRecord{Data: lib.NewByteReader(data), Record: map[string][]*model.LevelRecord{}}
	gr.Songsnum = int(gr.Data.GetVarInt())
	for gr.Data.Remaining() > 0 {
		key := gr.Data.GetString()
		gr.Data.SkipVarInt(0)
		length := gr.Data.GetByte()
		fc := gr.Data.GetByte()
		song := make([]*model.LevelRecord, 5)
		for level := 0; level < 5; level++ {
			if length&(1<<uint(level)) != 0 {
				song[level] = &model.LevelRecord{
					Score: gr.Data.GetInt(),
					Acc:   gr.Data.GetFloat(),
					Fc:    fc&(1<<uint(level)) != 0,
				}
			}
		}
		gr.Record[key] = song
	}
	return gr
}
