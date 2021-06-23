package common

import (
	"bufio"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type SimulatorStatus struct {
	Date      string
	Timestamp int64
	MarkPrice float64
	Equity    float64

	LongPositionSize   float64
	LongEntryPrice     float64
	LongGridReached    int64
	LongOrderAmount    float64
	LongOrderPrice     float64
	LongFee            float64
	LongRealizedProfit float64
	LongUnrealizedPNL  float64

	ShortPositionSize   float64
	ShortEntryPrice     float64
	ShortGridReached    int64
	ShortOrderAmount    float64
	ShortOrderPrice     float64
	ShortFee            float64
	ShortRealizedProfit float64
	ShortUnrealizedPNL  float64
}

type SimulatorResult struct {
	statusHistory []SimulatorStatus
}

func NewSimulatorResult() *SimulatorResult {
	return &SimulatorResult{
		statusHistory: make([]SimulatorStatus, 0),
	}
}

func (s *SimulatorResult) Append(status SimulatorStatus) {
	s.statusHistory = append(s.statusHistory, status)
}

func (s *SimulatorResult) WriteToFile(filepath string) error {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	datawriter := bufio.NewWriter(file)
	_, err = datawriter.WriteString("Date,Timestamp,MarkPrice,Equity,PositionSize-L,EntryPrice-L,GridReached-L,OrderSize-L,OrderPrice-L,GrossProfit-L,PNL-L,PositionSize-S,EntryPrice-S,GridReached-S,OrderSize-S,OrderPrice-S,GrossProfit-S,PNL-S\n")
	if err != nil {
		log.Error("Error writing header of file")
	}

	for N, st := range s.statusHistory {
		// reduce output file size with 1 min discretization
		if N%60 != 0 && st.LongOrderAmount == 0 && st.ShortOrderAmount == 0 {
			continue
		}

		_, err = datawriter.WriteString(fmt.Sprintf("%s,%d,%f,%f,%f,%f,%d,%f,%f,%f,%f,%f,%f,%d,%f,%f,%f,%f\n",
			st.Date, st.Timestamp, st.MarkPrice, st.Equity,
			st.LongPositionSize, st.LongEntryPrice, st.LongGridReached, st.LongOrderAmount, st.LongOrderPrice, st.LongRealizedProfit, st.LongUnrealizedPNL,
			st.ShortPositionSize, st.ShortEntryPrice, st.ShortGridReached, st.ShortOrderAmount, st.ShortOrderPrice, st.ShortRealizedProfit, st.ShortUnrealizedPNL))
		if err != nil {
			log.Error("Error writing status to result file")
		}
	}
	datawriter.Flush()
	return nil
}

func (s *SimulatorResult) Performance() float64 {
	first := s.statusHistory[0]
	last := s.statusHistory[len(s.statusHistory)-1]
	return last.Equity - first.Equity
}
