package common

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type SymbolDataItem struct {
	Time  time.Time
	Price float64
}

type SymbolData struct {
	Symbol    string
	StartDate string
	EndDate   string
	Data      []SymbolDataItem
}

func (d *SymbolData) readFromFile(path string) {
	file, err := os.Open(path)
	failOnError(err, fmt.Sprintf("Could not open file %s", path))
	defer file.Close()

	d.Data = make([]SymbolDataItem, 0)
	scanner := bufio.NewScanner(file)
	scanner.Scan() // skip header
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Split(line, ",")

		timestamp, err := strconv.ParseFloat(values[1], 64)
		failOnError(err, "Error parsing timestamp")
		price, err := strconv.ParseFloat(values[2], 64)
		failOnError(err, "Error parsing price")

		d.Data = append(d.Data, SymbolDataItem{Time: time.Unix(int64(timestamp), 0), Price: price})
	}
	err = scanner.Err()
	failOnError(err, "Scanner error")

	// TODO set symbol
	d.StartDate = d.Data[0].Time.UTC().String()
	d.EndDate = d.Data[len(d.Data)-1].Time.UTC().String()
}

func (d *SymbolData) append(symbolData *SymbolData) {
	// TODO check that symbolData.startDate is after d.endDate
	d.Data = append(d.Data, symbolData.Data...)
	d.EndDate = symbolData.EndDate
}

func (d *SymbolData) writeToFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	datawriter := bufio.NewWriter(file)
	_, err = datawriter.WriteString("Timestamp,Price\n")
	if err != nil {
		log.Error("Error writing header of file")
	}

	for _, di := range d.Data {
		_, err = datawriter.WriteString(fmt.Sprintf("%d,%f\n", di.Time.Unix(), di.Price))
		if err != nil {
			log.Error("Error writing status to result file")
		}
	}
	datawriter.Flush()
	return nil
}

func NewSymbolDataFromProcessedFile(filePath string) *SymbolData {
	symbolData := &SymbolData{}
	symbolData.readFromFile(filePath)
	return symbolData
}

func NewSymbolDataFromTickDataFolder(folderPath string) *SymbolData {
	var files []string
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".csv" {
			files = append(files, path)
		}
		return nil
	})
	failOnError(err, fmt.Sprintf("Could not get files in folder %s", folderPath))

	symbolData := &SymbolData{}
	for _, file := range files {
		symbolData.append(processTickData(file))
	}

	// TODO set symbol
	symbolData.StartDate = symbolData.Data[0].Time.UTC().String()
	symbolData.EndDate = symbolData.Data[len(symbolData.Data)-1].Time.UTC().String()

	symbolData.writeToFile("../datasets/DOGE_1s.csv")
	return symbolData
}

func processTickData(filePath string) *SymbolData {
	log.Infof("Processing tick data file: %s", filePath)

	file, err := os.Open(filePath)
	failOnError(err, fmt.Sprintf("Could not open file %s", filePath))
	defer file.Close()

	symbolData := &SymbolData{}
	symbolData.Data = make([]SymbolDataItem, 0)
	scanner := bufio.NewScanner(file)
	scanner.Scan() // skip header

	lastTimestamp := 0
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Split(line, ",")

		timestamp, err := strconv.Atoi(values[5])
		timestamp /= 1000 // convert to seconds
		failOnError(err, "Error parsing timestamp")
		price, err := strconv.ParseFloat(values[2], 64)
		failOnError(err, "Error parsing price")

		if timestamp >= (lastTimestamp + 1) { // append new value only if 1 second has passed
			lastTimestamp = timestamp
			symbolData.Data = append(symbolData.Data, SymbolDataItem{Time: time.Unix(int64(lastTimestamp), 0), Price: price})
		}
	}
	err = scanner.Err()
	failOnError(err, "Scanner error")
	return symbolData
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panic(msg)
	}
}
