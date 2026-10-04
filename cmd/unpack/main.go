// Command unpack fetches the Phigros APK and extracts difficulty/info tables
// and illustration images, replacing the old Python scripts.
package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xeonds/phi-plug-go/internal/addressables"
	"github.com/xeonds/phi-plug-go/internal/phigros"
	"github.com/xeonds/phi-plug-go/internal/taptap"
)

func httpGet(url string) (*http.Response, error) { return http.Get(url) }

func writePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "fetch":
		out := "build/base.apk"
		if len(os.Args) > 2 {
			out = os.Args[2]
		}
		fetch(out)
	case "run":
		if len(os.Args) < 4 {
			usage()
		}
		run(os.Args[2], os.Args[3])
	case "version":
		if len(os.Args) < 3 {
			usage()
		}
		printVersion(os.Args[2])
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage:
  unpack fetch [out.apk]        fetch the latest Phigros APK via TapTap
  unpack run <apk> <outdir>     extract difficulty.tsv/info.tsv + illustrations
  unpack version <apk>          print the game version from the APK`)
	os.Exit(2)
}

func printVersion(apkPath string) {
	zr, err := zip.OpenReader(apkPath)
	if err != nil {
		fatal(err)
	}
	defer zr.Close()
	manifest := readZip(&zr.Reader, "AndroidManifest.xml")
	if manifest == nil {
		fatal(fmt.Errorf("AndroidManifest.xml not found"))
	}
	v, err := phigros.ParseVersion(manifest)
	if err != nil {
		fatal(err)
	}
	fmt.Println(v)
}

func fetch(out string) {
	dl, err := taptap.Latest()
	if err != nil {
		fatal(err)
	}
	fmt.Printf("phigros %s (%d bytes)\n", dl.Name, dl.Size)
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		fatal(err)
	}
	f, err := os.Create(out)
	if err != nil {
		fatal(err)
	}
	defer f.Close()
	resp, err := httpGet(dl.URL)
	if err != nil {
		fatal(err)
	}
	defer resp.Body.Close()
	n, err := io.Copy(f, resp.Body)
	if err != nil {
		fatal(err)
	}
	fmt.Printf("wrote %s (%d bytes)\n", out, n)
}

func run(apkPath, outDir string) {
	zr, err := zip.OpenReader(apkPath)
	if err != nil {
		fatal(err)
	}
	defer zr.Close()

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fatal(err)
	}

	if manifest := readZip(&zr.Reader, "AndroidManifest.xml"); manifest != nil {
		if v, err := phigros.ParseVersion(manifest); err == nil {
			writeConfig(outDir, v)
		} else {
			fmt.Fprintln(os.Stderr, "warning: version:", err)
		}
	}

	if data := readZip(&zr.Reader, "assets/bin/Data/data.unity3d"); data != nil {
		writeInfo(data, outDir)
	} else {
		fmt.Fprintln(os.Stderr, "warning: data.unity3d not found; skipping info")
	}

	if cat := readZip(&zr.Reader, "assets/aa/catalog.json"); cat != nil {
		writeIllustrations(&zr.Reader, cat, outDir)
	} else {
		fmt.Fprintln(os.Stderr, "warning: catalog.json not found; skipping illustrations")
	}
}

func writeInfo(dataUnity3D []byte, outDir string) {
	diffs, infos, err := phigros.ExtractGameInfo(dataUnity3D)
	if err != nil {
		fatal(err)
	}
	fmt.Printf("extracted %d songs\n", len(diffs))

	var b strings.Builder
	b.WriteString("id\tEZ\tHD\tIN\tAT\tLegacy\n")
	for _, d := range diffs {
		b.WriteString(d.ID)
		for _, r := range d.Ranks {
			b.WriteByte('\t')
			b.WriteString(strconv.FormatFloat(r, 'f', 1, 64))
		}
		b.WriteByte('\n')
	}
	writeFile(filepath.Join(outDir, "difficulty.tsv"), b.String())

	b.Reset()
	b.WriteString("id\tsong\tcomposer\tillustrator\tEZ\tHD\tIN\tAT\tLegacy\n")
	for _, in := range infos {
		b.WriteString(in.ID + "\t" + in.Song + "\t" + in.Composer + "\t" + in.Illustrator)
		for _, c := range in.Charters {
			b.WriteByte('\t')
			b.WriteString(c)
		}
		b.WriteByte('\n')
	}
	writeFile(filepath.Join(outDir, "info.tsv"), b.String())
}

func writeIllustrations(zr *zip.Reader, catalogJSON []byte, outDir string) {
	var catalog struct {
		Key    string `json:"m_KeyDataString"`
		Bucket string `json:"m_BucketDataString"`
		Entry  string `json:"m_EntryDataString"`
	}
	if err := json.Unmarshal(catalogJSON, &catalog); err != nil {
		fatal(fmt.Errorf("catalog: %w", err))
	}
	entries, err := addressables.Parse(catalog.Key, catalog.Bucket, catalog.Entry)
	if err != nil {
		fatal(err)
	}

	dir := filepath.Join(outDir, "assets", "illustrations")
	if err := os.RemoveAll(dir); err != nil {
		fatal(err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fatal(err)
	}

	count := 0
	for _, e := range entries {
		i := strings.Index(e.Key, ".0/IllustrationLowRes")
		if i < 0 {
			continue
		}
		name := e.Key[:i]
		bundle := readZip(zr, "assets/aa/Android/"+e.Bundle)
		if bundle == nil {
			continue
		}
		img, err := phigros.DecodeTexture2D(bundle)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", name, err)
			continue
		}
		if err := writePNG(filepath.Join(dir, name+".png"), img); err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", name, err)
			continue
		}
		count++
	}
	fmt.Printf("extracted %d illustrations\n", count)
}

func readZip(zr *zip.Reader, name string) []byte {
	f, err := zr.Open(name)
	if err != nil {
		return nil
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		fatal(err)
	}
	return data
}

func writeConfig(outDir, version string) {
	cfg := map[string]string{
		"version":    version,
		"updateTime": time.Now().Format("2006-01-02"),
	}
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		fatal(err)
	}
	path := filepath.Join(outDir, "config.json")
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		fatal(err)
	}
	fmt.Printf("wrote %s\n", path)
}

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		fatal(err)
	}
	fmt.Printf("wrote %s\n", path)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
