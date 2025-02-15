package lib_test

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"
)

func TestReadCsv(t *testing.T) {
	t.Log("ReadCsv")
	LoadCSV("../build/dist/info.tsv")
}

func LoadCSV(filePath string) (map[string]map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Failed to open file: ", err)
		return nil, err
	}
	defer file.Close()

	var rows [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		rows = append(rows, strings.Split(line, "\t"))
	}
	if err := scanner.Err(); err != nil {
		log.Println("Failed to read file: ", err)
		return nil, err
	}

	headers := rows[0]
	log.Println(len(headers))
	data := make(map[string]map[string]string, len(rows)-1)
	for i := 1; i < len(rows); i++ {
		values := make(map[string]string)
		log.Println(len(rows[i]), i)
		// here, each row has its own length
		// so we need to loop according to each row's length
		for j := 1; j < len(rows[i]); j++ {
			values[headers[j]] = rows[i][j]
		}
		// use first column as key
		// may be error-prone but ok for this use case
		// because in this case the first column is unique
		data[rows[i][0]] = values
	}
	return data, nil
}
