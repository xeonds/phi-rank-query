package lib

import (
	"bufio"
	"os"
	"strings"
)

// LoadTSV reads a tab-separated file into a map keyed by the first column.
func LoadTSV(path string) (map[string]map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var rows [][]string
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		rows = append(rows, strings.Split(scanner.Text(), "\t"))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return map[string]map[string]string{}, nil
	}

	headers := rows[0]
	data := make(map[string]map[string]string, len(rows)-1)
	for _, row := range rows[1:] {
		values := make(map[string]string, len(row))
		for j := 1; j < len(row); j++ {
			values[headers[j]] = row[j]
		}
		data[row[0]] = values
	}
	return data, nil
}
