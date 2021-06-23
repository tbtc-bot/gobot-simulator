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
	symbolData      common.SymbolData
	resultsFolder   string
	worker          worker.Worker
	exchange        engine.Exchange
	simulatorResult common.SimulatorResult
}

func NewSimulator(dataFolder string, resultsFolder string) *Simulator {
	exchange := engine.NewExchange()
	simulation := &Simulator{
		symbolData:      common.NewSymbolData(dataFolder),
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
}

func (s *Simulator) RunMultipleSimulations() {
	symbol := ""

	for _, GOi := range []uint{4, 6} {
		for _, GSi := range []float64{0.3, 0.5} {
			for _, SFi := range []float64{1.3, 1.7} {
				for _, OFi := range []float64{1.5, 2.2} {
					for _, TSi := range []float64{0.2, 0.6} {
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
	log.Info("Simulation start")

	// Initialize exchange
	s.exchange.Init(1000, s.symbolData[0])

	// Start strategy
	s.worker.SetStrategy(strategy)
	s.worker.StartStrategy()

	// Cycle over symbol data
	for _, tick := range s.symbolData[1:] {
		s.exchange.Next(tick)

		// recreate grid if worker has still no position after some time
		// s.worker.HandleGridRecreation(tick.Time)
	}

	log.Info("Simulation done")

	resultFile := s.resultsFolder + strategy.String() + ".csv"
	if err := s.writeResults(resultFile); err != nil {
		log.Errorf("Error writing results to file %s: %s", resultFile, err)
	} else {
		log.Infof("Simulation results saved to %s", resultFile)
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
