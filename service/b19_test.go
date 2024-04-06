package service_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/xeonds/phi-plug-go/service"
)

const session = "9r5shbrhcz9omano8qx5ahn3f"

func TestB19(t *testing.T) {
	b19Data, err := service.GetB19Info(session)
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
	savezip, err := service.GetSaveZip(latestModified.Gamefile.URL)
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
