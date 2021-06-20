package simulator

import (
	"example.com/gobot-simulator/src/common"
	"example.com/gobot-simulator/src/engine"
	"example.com/gobot-simulator/src/worker"
	log "github.com/sirupsen/logrus"
)

type Simulator struct {
	symbolData      common.SymbolData
	worker          worker.Worker
	exchange        engine.Exchange
	simulatorResult common.SimulatorResult
}

func NewSimulator(dataFolder string, worker worker.Worker) *Simulator {
	exchange := engine.NewExchange()
	simulation := &Simulator{
		symbolData:      common.NewSymbolData(dataFolder),
		worker:          worker,
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
func (s *Simulator) Start() {
	log.Info("Simulation start")

	// Initialize exchange
	s.exchange.Init(1000, s.symbolData[0])

	// Start strategy
	s.worker.StartStrategy()

	// Cycle over symbol data
	for _, tick := range s.symbolData[1:] {
		s.exchange.Next(tick)
	}

	log.Info("Simulation done")
}

func (s *Simulator) WriteResults(filepath string) error {
	return s.simulatorResult.WriteToFile(filepath)
}

// PRIVATE METHODS
func (s *Simulator) updateResult(status common.SimulatorStatus) {
	s.simulatorResult.Append(status)
}
