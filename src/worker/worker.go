package worker

import (
	"math"
	"time"

	"example.com/gobot-simulator/src/common"
	"example.com/gobot-simulator/src/engine"
	"example.com/gobot-simulator/src/strategy"

	log "github.com/sirupsen/logrus"
)

type Worker struct {
	strategy            strategy.StrategyWrapper
	exchangeAPI         *engine.ExchangeAPI
	lastTimeCreatedGrid time.Time // this is used in AntiMartingala to recreate the grid if after some time has still no position
}

func NewWorker() *Worker {
	return &Worker{}
}

// PUBLIC METHODS
func (w *Worker) SetStrategy(strategy strategy.StrategyWrapper) {
	w.strategy = strategy
}

func (w *Worker) SetExchangeAPI(api *engine.ExchangeAPI) {
	w.exchangeAPI = api
}

func (w *Worker) HandlePositionUpdate(position engine.Position) {
	if position.Size == 0 {
		w.StartStrategy()
	} else {
		w.setTakeProfit(position)
	}
}

func (w *Worker) StartStrategy() {
	log.Debug("Worker: start strategy")
	// w.lastTimeCreatedGrid = w.exchangeAPI.CurrentTime()
	symbol := w.strategy.GetSymbol()
	balance := w.exchangeAPI.Balance()
	markPrice := w.exchangeAPI.MarkPrice()
	s0 := (balance / markPrice) * (w.strategy.GetParameters().OS / 100)
	if math.IsNaN(s0) {
		log.Panic("Order size is NaN")
	}

	// TODO change this: worker should not know the strategy
	positionSide := w.strategy.GetPositionSide()
	w.createGrid(positionSide, balance, markPrice)

	// TODO this is valid only for Martingala
	var order engine.Order
	if positionSide == engine.PositionSideLong {
		order = *engine.NewOrderMarket(symbol, engine.SideBuy, engine.PositionSideLong, s0)
	} else {
		order = *engine.NewOrderMarket(symbol, engine.SideSell, engine.PositionSideShort, s0)
	}
	w.placeOrder(order)
}

func (w *Worker) HandleGridRecreation(time time.Time) {
	position := w.exchangeAPI.Position(w.strategy.GetPositionSide())
	if position.Size > 0 {
		w.lastTimeCreatedGrid = time
		return
	}

	diffMinutes := time.Sub(w.lastTimeCreatedGrid).Minutes()
	if diffMinutes > 3 {
		log.Debug("Worker recreates grid because more than 3 mins have passed without having a position")
		w.StartStrategy()
	}
}

// PRIVATE METHODS
func (w *Worker) createGrid(positionSide engine.PositionSideType, balance float64, startPrice float64) {
	w.cancelGrid(positionSide)

	// get orders and execute
	var orders []*engine.Order
	if positionSide == engine.PositionSideLong {
		orders = w.strategy.BuyGridOrders(balance, startPrice)
	} else {
		orders = w.strategy.SellGridOrders(balance, startPrice)
	}

	for _, order := range orders {
		w.placeOrder(*order)
	}
}

func (w *Worker) cancelGrid(positionSide engine.PositionSideType) {
	openOrders := w.exchangeAPI.OpenOrders(positionSide)
	for _, o := range openOrders {
		w.cancelOrder(o)
	}
}

func (w *Worker) setTakeProfit(position engine.Position) {
	w.cancelLastTakeProfit(position.PositionSide)
	gridReached := w.exchangeAPI.GridReached(position.PositionSide)
	order := w.strategy.TakeProfitOrder(position, gridReached)
	log.Debugf("Worker: set take profit order %s", order.String())
	w.placeOrder(*order)
}

func (w *Worker) cancelLastTakeProfit(positionSide engine.PositionSideType) {
	openOrders := w.exchangeAPI.OpenOrders(positionSide)
	for _, order := range openOrders {
		if order.IsTP {
			w.cancelOrder(order)
		}
	}
}

func (w *Worker) cancelOrder(order engine.Order) bool {
	return w.exchangeAPI.CancelOrder(order)
}

func (w *Worker) placeOrder(order engine.Order) {
	// round order price and amount to 6 digits
	order.Amount = common.RoundFloatWithPrecision(order.Amount, 6)
	order.Price = common.RoundFloatWithPrecision(order.Price, 6)

	w.exchangeAPI.PlaceOrder(order)
}
