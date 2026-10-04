// Package taptap fetches the current Phigros APK download URL from TapTap.
package taptap

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	appID     = 165287 // Phigros
	secret    = "PeCkE6Fu0B10Vm9BKfPfANwCUAn5POcs"
	apiHost   = "https://api.taptapdada.com"
	userAgent = "okhttp/3.12.1"
)

type Download struct {
	URL  string
	Name string
	Size int64
}

func xUA() string {
	return "V=1&PN=TapTap&VN=2.40.1-rel.100000&VN_CODE=240011000&LOC=CN&LANG=zh_CN&CH=default&UID=" +
		uuidString() + "&NT=1&SR=1080x2030&DEB=Xiaomi&DEM=Redmi+Note+5&OSV=9"
}

func uuidString() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func getJSON(rawurl string) (map[string]any, error) {
	req, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("taptap: %w", err)
	}
	return out, nil
}

// Latest returns the current APK download info for Phigros.
func Latest() (*Download, error) {
	ua := xUA()
	detail, err := getJSON(apiHost + "/app/v2/detail-by-id/" + strconv.Itoa(appID) + "?X-UA=" + url.QueryEscape(ua))
	if err != nil {
		return nil, err
	}
	data, _ := detail["data"].(map[string]any)
	dl, _ := data["download"].(map[string]any)
	apkID := int64(dl["apk_id"].(float64))

	nonce := randomNonce(5)
	t := time.Now().Unix()
	param := fmt.Sprintf("abi=arm64-v8a,armeabi-v7a,armeabi&id=%d&node=%s&nonce=%s&sandbox=1&screen_densities=xhdpi&time=%d",
		apkID, uuidString(), nonce, t)
	sum := md5.Sum([]byte("X-UA=" + ua + "&" + param + secret))
	body := param + "&sign=" + hex.EncodeToString(sum[:])

	req, err := http.NewRequest("POST", apiHost+"/apk/v1/detail?X-UA="+url.QueryEscape(ua), strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("taptap: %w", err)
	}
	d, _ := out["data"].(map[string]any)
	apk, _ := d["apk"].(map[string]any)
	return &Download{
		URL:  apk["download"].(string),
		Name: apk["name"].(string),
		Size: int64(apk["size"].(float64)),
	}, nil
}

func randomNonce(n int) string {
	const sample = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	seed := time.Now().UnixNano()
	for i := range b {
		seed = seed*6364136223846793005 + 1442695040888963407
		b[i] = sample[uint64(seed)%uint64(len(sample))]
	}
	return string(b)
}
