package service_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/xeonds/phi-plug-go/config"
	"github.com/xeonds/phi-plug-go/lib"
	"github.com/xeonds/phi-plug-go/service"
)

var session = os.Getenv("SESSION")
var conf = lib.LoadConfig[config.Config]()

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
	data := new(service.GameSave)
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
	data := new(service.GameSave)
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
	data := new(service.Game)
	if err = json.Unmarshal(f, &data); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}
	config := lib.LoadConfig[config.Config]()
	b19, rks, phi := service.CalcBNInfo(data, config, 19)
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
