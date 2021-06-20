package common

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type SymbolData []SymbolDataItem

type SymbolDataItem struct {
	Time  time.Time
	Price float64
}

func NewSymbolData(filePath string) SymbolData {
	file, err := os.Open(filePath)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	log.Info("Reading file: ", filePath)

	symbolData := make(SymbolData, 0)
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Split(line, ",")

		timestamp, err := strconv.ParseFloat(values[1], 64)
		if err != nil {
			log.Panic("Error parsing timestamp")
		}
		price, err := strconv.ParseFloat(values[2], 64)
		if err != nil {
			log.Panic("Error parsing price")
		}
		symbolData = append(symbolData, SymbolDataItem{Time: time.Unix(int64(timestamp), 0), Price: price})
	}

	return symbolData
}

// func NewSymbolData(dataFolder string) SymbolData {
// 	symbolData := make(SymbolData, 0)

// 	var files []string
// 	err := filepath.Walk(dataFolder, func(path string, info os.FileInfo, err error) error {
// 		if filepath.Ext(path) == ".csv" {
// 			files = append(files, path)
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	for _, f := range files {
// 		file, err := os.Open(f)
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 		defer file.Close()

// 		scanner := bufio.NewScanner(file)

// 		scanner.Scan()
// 		log.Info("Reading file: ", f)

// 		rows := 0
// 		last_timestamp := 0
// 		for scanner.Scan() {
// 			line := scanner.Text()
// 			values := strings.Split(line, ",")

// 			price, _ := strconv.ParseFloat(values[2], 64)
// 			timestamp, _ := strconv.Atoi(values[5])
// 			if timestamp >= (last_timestamp + 1000) {
// 				rows++
// 				last_timestamp = timestamp
// 				symbolData = append(symbolData, SymbolDataItem{Time: time.Unix(int64(timestamp/1000), 0), Price: price})
// 			}
// 		}

// 		if err := scanner.Err(); err != nil {
// 			log.Panic(err)
// 		}
// 	}

// 	return symbolData
// }
