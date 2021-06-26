package simulator

import (
	"fmt"

	"example.com/gobot-simulator/src/common"
	"example.com/gobot-simulator/src/engine"
	"example.com/gobot-simulator/src/strategy"
	"example.com/gobot-simulator/src/worker"
	log "github.com/sirupsen/logrus"
)

type Simulator struct {
	symbolData      *common.SymbolData
	resultsFolder   string
	worker          worker.Worker
	exchange        engine.Exchange
	simulatorResult common.SimulatorResult
}

func NewSimulator(symbolData *common.SymbolData, resultsFolder string) *Simulator {
	exchange := engine.NewExchange()
	simulation := &Simulator{
		symbolData:      symbolData,
		resultsFolder:   resultsFolder,
		worker:          *worker.NewWorker(),
		exchange:        *exchange,
		simulatorResult: *common.NewSimulatorResult(),
	}

	// link worker, exchange and simulation through callbacks
	simulation.worker.SetExchangeAPI(simulation.exchange.GetAPI())
	simulation.exchange.NotifyPositionUpdateCallback = simulation.worker.HandlePositionUpdate
	simulation.exchange.UpdateSimulationStatusCallback = simulation.updateResult

	return simulation
}

// PUBLIC METHODS
func (s *Simulator) RunSingleSimulation(strategy strategy.StrategyWrapper) {
	info := s.start(strategy)
	fmt.Println(info)

	resultFile := s.resultsFolder + strategy.String() + ".csv"
	if err := s.writeResults(resultFile); err != nil {
		log.Errorf("Error writing results to file %s: %s", resultFile, err)
	} else {
		log.Infof("Simulation results saved to %s", resultFile)
	}
}

func (s *Simulator) RunMultipleSimulations(GOvec []uint, GSvec []float64, SFvec []float64, OFvec []float64, TSvec []float64) {
	symbol := ""

	N := len(GOvec) * len(GSvec) * len(SFvec) * len(OFvec) * len(TSvec)
	n := 0

	for _, GOi := range GOvec {
		for _, GSi := range GSvec {
			for _, SFi := range SFvec {
				for _, OFi := range OFvec {
					for _, TSi := range TSvec {
						n++
						fmt.Printf("Starting simulation %d/%d", n, N)
						pars := strategy.StrategyParameters{GO: GOi, GS: GSi, SF: SFi, OS: 1, OF: OFi, TS: TSi, SL: 0.3}
						strategy := strategy.NewStrategy(strategy.StrategyTypeAntiMartingala, symbol, engine.PositionSideLong, pars)
						info := s.start(*strategy)
						fmt.Println(info)
					}
				}
			}
		}
	}
}

func (s *Simulator) start(strategy strategy.StrategyWrapper) string {
	// Initialize exchange
	s.exchange.Init(1000, s.symbolData.Data[0])

	// Start strategy
	s.worker.SetStrategy(strategy)
	s.worker.StartStrategy()

	// Cycle over symbol data
	for _, tick := range s.symbolData.Data[1:] {
		s.exchange.Next(tick)

		// recreate grid if worker has still no position after some time
		// s.worker.HandleGridRecreation(tick.Time)
	}
	return strategy.String() + " -> " + fmt.Sprint(s.simulatorResult.Performance())
}

func (s *Simulator) writeResults(filepath string) error {
	return s.simulatorResult.WriteToFile(filepath)
}

// PRIVATE METHODS
func (s *Simulator) updateResult(status common.SimulatorStatus) {
	s.simulatorResult.Append(status)
}
