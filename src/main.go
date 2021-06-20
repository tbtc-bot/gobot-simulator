package main

import (
	"fmt"
	"io"
	"os"

	"example.com/gobot-simulator/src/engine"
	"example.com/gobot-simulator/src/simulator"
	"example.com/gobot-simulator/src/strategy"
	"example.com/gobot-simulator/src/worker"

	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	file, _ := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&easy.Formatter{
		// TimestampFormat: "2006-01-02 15:04:05",
		LogFormat: "[%lvl%] %msg%\n",
	})

	symbol := ""
	pars := strategy.StrategyParameters{
		GO: 7,
		GS: 0.3,
		SF: 1.5,
		OS: 1,
		OF: 2,
		TS: 0.3,
		SL: 0.3,
	}
	strategy := strategy.NewStrategy(strategy.StrategyTypeAntiMartingala, symbol, engine.PositionSideLong, pars)
	worker := worker.NewWorker(*strategy)

	dataFolder := "../datasets/LTCUSDT_1s.csv"
	simulator := simulator.NewSimulator(dataFolder, *worker)
	simulator.Start()

	resultFile := "../results/" + (*strategy).String() + ".csv"
	if err := simulator.WriteResults(resultFile); err != nil {
		log.Errorf("Error writing results to file %s: %s", resultFile, err)
	} else {
		log.Infof("Simulation results saved to %s", resultFile)
	}
}
