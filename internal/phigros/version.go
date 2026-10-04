package phigros

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/shogo82148/androidbinary"
)

// ParseVersion extracts android:versionName from a binary AndroidManifest.xml.
func ParseVersion(manifest []byte) (string, error) {
	f, err := androidbinary.NewXMLFile(bytes.NewReader(manifest))
	if err != nil {
		return "", fmt.Errorf("phigros: parse AndroidManifest: %w", err)
	}
	var m struct {
		XMLName     xml.Name `xml:"manifest"`
		VersionName string   `xml:"versionName,attr"`
	}
	if err := f.Decode(&m, nil, nil); err != nil {
		return "", fmt.Errorf("phigros: decode AndroidManifest: %w", err)
	}
	if m.VersionName == "" {
		return "", fmt.Errorf("phigros: versionName not found")
	}
	return m.VersionName, nil
}
