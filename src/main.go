package main

import (
	"io"
	"os"

	"example.com/gobot-simulator/src/common"
	"example.com/gobot-simulator/src/engine"
	"example.com/gobot-simulator/src/simulator"
	"example.com/gobot-simulator/src/strategy"

	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func main() {
	file, _ := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&easy.Formatter{
		// TimestampFormat: "2006-01-02 15:04:05",
		LogFormat: "[%lvl%] %msg%\n",
	})

	// symbolData := common.NewSymbolDataFromProcessedFile("../datasets/LTCUSDT_1s.csv")
	symbolData := common.NewSymbolDataFromTickDataFolder("../datasets/test_doge")
	resultsFolder := "../results/"
	simulator := simulator.NewSimulator(symbolData, resultsFolder)

	pars := strategy.StrategyParameters{GO: 5, GS: 0.3, SF: 1.5, OS: 1, OF: 2, TS: 0.3, SL: 0.3}
	strategy := strategy.NewStrategy(strategy.StrategyTypeAntiMartingala, "", engine.PositionSideLong, pars)
	simulator.RunSingleSimulation(*strategy)

	// TODO parameters combination
	// simulator.RunMultipleSimulations()
}
