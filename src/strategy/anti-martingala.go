package strategy

import (
	"log"
	"math"
	"time"

	"example.com/gobot-simulator/src/engine"
)

type StrategyAntiMartingala struct {
	Type         StrategyType            `json:"type"`
	Symbol       string                  `json:"symbol"`
	PositionSide engine.PositionSideType `json:"position_side"`
	Status       string                  `json:"status"`
	Parameters   StrategyParameters      `json:"parameters"`
}

const (
	// retry parameters
	ATTEMPTS = 3
	SLEEP    = 200 * time.Millisecond
)

func (s *StrategyAntiMartingala) BuyGridOrders(balance float64, startPrice float64) []*engine.Order {
	orders := []*engine.Order{}

	// start price and size
	s0 := (balance / startPrice) * (s.Parameters.OS / 100)

	// grids
	for grid := 1; grid < int(s.Parameters.GO)+1; grid++ {
		mf := 1 + s.Parameters.GS/100
		for ii := 1; ii < grid; ii++ {
			mf += (s.Parameters.GS / 100) * math.Pow(s.Parameters.SF, float64(ii))
		}
		p_i := startPrice * mf
		s_i := s0 * math.Pow(s.Parameters.OF, float64(grid-1))
		order := engine.NewOrderStop(s.Symbol, engine.SideBuy, engine.PositionSideLong, s_i, p_i, p_i)
		order.GridNumber = int64(grid)
		orders = append(orders, order)
	}

	return orders
}

func (s *StrategyAntiMartingala) SellGridOrders(balance float64, startPrice float64) []*engine.Order {
	orders := []*engine.Order{}

	// start price and size
	s0 := (balance / startPrice) * (s.Parameters.OS / 100)

	// grids
	for grid := 1; grid < int(s.Parameters.GO)+1; grid++ {
		ds := 1 + s.Parameters.GS/100
		for ii := 1; ii < grid; ii++ {
			ds += (s.Parameters.GS / 100) * math.Pow(s.Parameters.SF, float64(ii))
		}
		p_i := startPrice / ds
		s_i := s0 * math.Pow(s.Parameters.OF, float64(grid-1))
		order := engine.NewOrderStop(s.Symbol, engine.SideSell, engine.PositionSideShort, s_i, p_i, p_i)
		order.GridNumber = int64(grid)
		orders = append(orders, order)
	}

	return orders
}

func (s *StrategyAntiMartingala) TakeProfitOrder(position engine.Position, currentGrid int64) *engine.Order {
	markPrice := position.MarkPrice
	entryPrice := position.EntryPrice

	// if last grid has been reached
	// if currentGrid == int64(s.Parameters.GO) {
	// 	if s.PositionSide == engine.PositionSideLong {
	// 		takeProfitPrice := markPrice * (1 + 5*s.Parameters.SL/100)
	// 		order := engine.NewOrderLimit(s.Symbol, engine.SideSell, engine.PositionSideLong, position.Size, takeProfitPrice)
	// 		order.IsTP = true
	// 		return order
	// 	} else {
	// 		takeProfitPrice := markPrice * (1 - 5*s.Parameters.SL/100)
	// 		order := engine.NewOrderLimit(s.Symbol, engine.SideBuy, engine.PositionSideShort, math.Abs(position.Size), takeProfitPrice)
	// 		order.IsTP = true
	// 		return order
	// 	}

	// set stop loss
	if currentGrid < 3 {
		if s.PositionSide == engine.PositionSideLong {
			takeProfitPrice := entryPrice * (1 - s.Parameters.SL/100)
			order := engine.NewOrderStop(s.Symbol, engine.SideSell, engine.PositionSideLong, position.Size, takeProfitPrice, takeProfitPrice)
			order.IsTP = true
			return order
		} else {
			takeProfitPrice := entryPrice * (1 + s.Parameters.SL/100)
			order := engine.NewOrderStop(s.Symbol, engine.SideBuy, engine.PositionSideShort, math.Abs(position.Size), takeProfitPrice, takeProfitPrice)
			order.IsTP = true
			return order
		}

		// set take profit
	} else if currentGrid >= 3 {
		if s.PositionSide == engine.PositionSideLong {
			// TODO
			// takeProfitPrice := math.Pow(position.EntryPrice, s.Parameters.TS) * math.Pow(markPrice, 1-s.Parameters.TS)
			// takeProfitPrice := entryPrice * (1 + float64(currentGrid)/2*s.Parameters.SL/100)
			takeProfitPrice := entryPrice + (markPrice-entryPrice)*(float64(currentGrid)/float64(s.Parameters.GO))*0.9
			order := engine.NewOrderStop(s.Symbol, engine.SideSell, engine.PositionSideLong, position.Size, takeProfitPrice, takeProfitPrice)
			order.IsTP = true
			return order
		} else {
			// takeProfitPrice := math.Pow(position.EntryPrice, s.Parameters.TS) * math.Pow(markPrice, 1-s.Parameters.TS)
			// takeProfitPrice := entryPrice * (1 - float64(currentGrid)/2*s.Parameters.SL/100)
			takeProfitPrice := entryPrice - (entryPrice-markPrice)/2
			order := engine.NewOrderStop(s.Symbol, engine.SideBuy, engine.PositionSideShort, math.Abs(position.Size), takeProfitPrice, takeProfitPrice)
			order.IsTP = true
			return order
		}

	} else {
		log.Panicf("Current grid not valid %d", currentGrid)
		return nil
	}
}

func (s *StrategyAntiMartingala) String() string {
	return string(s.GetType()) + " " + string(s.GetPositionSide()) + " " + s.Parameters.String()
}

// GETTERS
func (s *StrategyAntiMartingala) GetType() StrategyType { return s.Type }

func (s *StrategyAntiMartingala) GetSymbol() string { return s.Symbol }

func (s *StrategyAntiMartingala) GetPositionSide() engine.PositionSideType { return s.PositionSide }

func (s *StrategyAntiMartingala) GetStatus() string { return s.Status }

func (s *StrategyAntiMartingala) GetParameters() StrategyParameters { return s.Parameters }

// SETTERS
func (s *StrategyAntiMartingala) SetSymbol(value string) { s.Symbol = value }

func (s *StrategyAntiMartingala) SetPositionSide(value engine.PositionSideType) {
	s.PositionSide = value
}

func (s *StrategyAntiMartingala) SetStatus(value string) { s.Status = value }

func (s *StrategyAntiMartingala) SetParameters(value StrategyParameters) { s.Parameters = value }

func NewStrategyAntiMartingala(symbol string, positionSide engine.PositionSideType, pars StrategyParameters) *StrategyAntiMartingala {
	return &StrategyAntiMartingala{
		Type:         StrategyTypeAntiMartingala,
		Symbol:       symbol,
		PositionSide: positionSide,
		Status:       "",
		Parameters:   pars,
	}
}
