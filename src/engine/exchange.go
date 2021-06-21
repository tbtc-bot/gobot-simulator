package engine

import (
	"fmt"
	"time"

	"example.com/gobot-simulator/src/common"

	log "github.com/sirupsen/logrus"
)

type Exchange struct {
	time      time.Time
	markPrice float64
	balance   float64

	sessionLong  Session
	sessionShort Session

	NotifyPositionUpdateCallback   func(Position)
	UpdateSimulationStatusCallback func(common.SimulatorStatus)

	orderCounter int64 // used for order ID
}

func NewExchange() *Exchange {
	return &Exchange{
		sessionLong:  *NewSession(PositionSideLong),
		sessionShort: *NewSession(PositionSideShort),
		orderCounter: 0,
	}
}

// PUBLIC METHODS
func (e *Exchange) Init(balance float64, tick common.SymbolDataItem) {
	e.time = tick.Time
	e.markPrice = tick.Price
	e.balance = balance
	status := common.SimulatorStatus{
		Date:      e.time.String(),
		Timestamp: e.time.Unix(),
		MarkPrice: e.markPrice,
		Equity:    e.balance,
	}
	e.UpdateSimulationStatusCallback(status)
}

func (e *Exchange) Next(tick common.SymbolDataItem) {
	log.Debugf("### Next tick: date %s, mark price %.2f", tick.Time.String(), tick.Price)
	e.time = tick.Time
	markPrice := common.RoundFloatWithPrecision(tick.Price, 6)
	e.markPrice = markPrice
	e.sessionLong.position.MarkPrice = markPrice
	e.sessionShort.position.MarkPrice = markPrice

	// handle long order
	var longOrderAmount, longOrderPrice float64
	if orderLong := e.getOrderToExecuteLong(); orderLong != nil {
		e.executeOrder(*orderLong)
		longOrderAmount = orderLong.Amount
		longOrderPrice = orderLong.Price
	} else {
		longOrderAmount = 0
		longOrderPrice = 0
	}

	// handle short order
	var shortOrderAmount, shortOrderPrice float64
	if orderShort := e.getOrderToExecuteShort(); orderShort != nil {
		e.executeOrder(*orderShort)
		shortOrderAmount = orderShort.Amount
		shortOrderPrice = orderShort.Price
	} else {
		shortOrderAmount = 0
		shortOrderPrice = 0
	}

	// update status
	status := common.SimulatorStatus{
		Date:      e.time.String(),
		Timestamp: e.time.Unix(),
		MarkPrice: e.markPrice,
		Equity:    e.balance,

		LongPositionSize: e.sessionLong.position.Size,
		LongEntryPrice:   e.sessionLong.position.EntryPrice,
		LongGridReached:  e.sessionLong.gridReached,
		LongOrderAmount:  longOrderAmount,
		LongOrderPrice:   longOrderPrice,
		// LongFee           float64
		LongRealizedProfit: e.sessionLong.realizedProfit,
		// LongNetProfit     float64
		// LongROEPerc       float64
		LongUnrealizedPNL: e.sessionLong.position.PNL(e.markPrice),
		// LongDrawdownPerc  float64

		ShortPositionSize: e.sessionShort.position.Size,
		ShortEntryPrice:   e.sessionShort.position.EntryPrice,
		ShortGridReached:  e.sessionShort.gridReached,
		ShortOrderAmount:  shortOrderAmount,
		ShortOrderPrice:   shortOrderPrice,
		// ShortFee           float64
		ShortRealizedProfit: e.sessionShort.realizedProfit,
		// ShortNetProfit     float64
		// ShortROEPerc       float64
		ShortUnrealizedPNL: e.sessionShort.position.PNL(e.markPrice),
		// ShortDrawdownPerc  float64
	}
	e.UpdateSimulationStatusCallback(status)
}

func (e *Exchange) GetAPI() *ExchangeAPI {
	return &ExchangeAPI{
		PlaceOrder:  e.placeOrder,
		CancelOrder: e.cancelOrder,
		OpenOrders:  e.getOpenOrders,
		MarkPrice:   e.getMarkPrice,
		Balance:     e.getBalance,
		GridReached: e.getGridReached,
	}
}

// PRIVATE METHODS
func (e *Exchange) getOrderToExecuteLong() *Order {
	// TODO check max 1 order per side is executed
	for _, o := range e.sessionLong.openOrders {
		switch o.Type {
		case OrderTypeMarket:
			return &o
		case OrderTypeLimit:
			if (o.Side == SideBuy && e.markPrice <= o.Price) || (o.Side == SideSell && e.markPrice >= o.Price) {
				return &o
			}
		case OrderTypeStop:
			if (o.Side == SideBuy && e.markPrice >= o.TriggerPrice) || (o.Side == SideSell && e.markPrice <= o.TriggerPrice) {
				return &o
			}
		case OrderTypeTrailing:
			// TODO
		}
	}
	return nil
}

func (e *Exchange) getOrderToExecuteShort() *Order {
	// TODO check max 1 order per side is executed
	for _, o := range e.sessionShort.openOrders {
		switch o.Type {
		case OrderTypeMarket:
			return &o
		case OrderTypeLimit:
			if (o.Side == SideSell && e.markPrice >= o.Price) || (o.Side == SideBuy && e.markPrice <= o.Price) {
				return &o
			}
		case OrderTypeStop:
			if (o.Side == SideSell && e.markPrice <= o.TriggerPrice) || (o.Side == SideBuy && e.markPrice >= o.TriggerPrice) {
				return &o
			}
		case OrderTypeTrailing:
			// TODO
		}
	}
	return nil
}

func (e *Exchange) executeOrder(order Order) {
	log.Debugf("Exchange: execute order %s", order.String())
	if order.PositionSide == PositionSideLong {
		if _, ok := e.sessionLong.openOrders[order.ID]; ok {
			delete(e.sessionLong.openOrders, order.ID)
		} else {
			log.Panic("Order id not found in open orders")
		}
		e.sessionLong.orderAmount = order.Amount
		e.sessionLong.orderPrice = order.Price
		realizedProfit := e.sessionLong.position.Update(order)
		e.sessionLong.realizedProfit = realizedProfit
		e.balance += realizedProfit
		if !order.IsTP { // don't update grid reached to 0 if is TP order, this is for the statistics
			e.sessionLong.gridReached = order.GridNumber
		}
		log.Debugf("Exchange: updated position %s", e.sessionLong.position.String())
		e.NotifyPositionUpdateCallback(e.sessionLong.position)
	} else {
		if _, ok := e.sessionShort.openOrders[order.ID]; ok {
			delete(e.sessionShort.openOrders, order.ID)
		} else {
			log.Panic("Order id not found in open orders")
		}
		e.sessionShort.orderAmount = order.Amount
		e.sessionShort.orderPrice = order.Price
		realizedProfit := e.sessionShort.position.Update(order)
		e.sessionShort.realizedProfit = realizedProfit
		e.balance += realizedProfit
		if !order.IsTP { // don't update grid reached to 0 if is TP order, this is for the statistics
			e.sessionShort.gridReached = order.GridNumber
		}
		log.Debugf("Exchange: updated position %s", e.sessionShort.position.String())
		e.NotifyPositionUpdateCallback(e.sessionShort.position)
	}
}

func (e *Exchange) placeOrder(order Order) {
	if order.Amount == 0 {
		log.Panic("Order amount is 0")
	}

	order.ID = fmt.Sprint(e.orderCounter)
	e.orderCounter++
	if order.PositionSide == PositionSideLong {
		e.sessionLong.openOrders[order.ID] = order
	} else {
		e.sessionShort.openOrders[order.ID] = order
	}
}

func (e *Exchange) cancelOrder(order Order) bool {
	if order.PositionSide == PositionSideLong {
		if _, ok := e.sessionLong.openOrders[order.ID]; ok {
			delete(e.sessionLong.openOrders, order.ID)
			return true
		} else {
			return false
		}
	} else {
		if _, ok := e.sessionShort.openOrders[order.ID]; ok {
			delete(e.sessionShort.openOrders, order.ID)
			return true
		} else {
			return false
		}
	}
}

func (e *Exchange) getOpenOrders(positionSide PositionSideType) []Order {
	var orders = make([]Order, 0)

	if positionSide == PositionSideLong {
		for _, order := range e.sessionLong.openOrders {
			orders = append(orders, order)
		}
	} else {
		for _, order := range e.sessionShort.openOrders {
			orders = append(orders, order)
		}
	}
	return orders
}

func (e *Exchange) getMarkPrice() float64 {
	return e.markPrice
}

func (e *Exchange) getBalance() float64 {
	return e.balance
}

func (e *Exchange) getGridReached(positionSide PositionSideType) int64 {
	if positionSide == PositionSideLong {
		return e.sessionLong.gridReached
	} else {
		return e.sessionShort.gridReached
	}
}
