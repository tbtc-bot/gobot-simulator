package strategy

import (
	"math"

	"example.com/gobot-simulator/src/engine"
)

type StrategyLogMartingala struct {
	Type         StrategyType            `json:"type"`
	Symbol       string                  `json:"symbol"`
	PositionSide engine.PositionSideType `json:"positionSide"`
	Status       string                  `json:"status"`
	Parameters   StrategyParameters      `json:"parameters"`
}

func (s *StrategyLogMartingala) BuyGridOrders(balance float64, startPrice float64) []*engine.Order {
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
		order := engine.NewOrderLimit(s.Symbol, engine.SideBuy, engine.PositionSideLong, s_i, p_i)
		order.GridNumber = int64(grid)
		orders = append(orders, order)
	}

	return orders
}

func (s *StrategyLogMartingala) SellGridOrders(balance float64, startPrice float64) []*engine.Order {
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
		order := engine.NewOrderLimit(s.Symbol, engine.SideSell, engine.PositionSideShort, s_i, p_i)
		order.GridNumber = int64(grid)
		orders = append(orders, order)
	}

	return orders
}

func (s *StrategyLogMartingala) TakeProfitOrder(position engine.Position, currentGrid int64) *engine.Order {
	switch s.PositionSide {
	case engine.PositionSideLong:
		takeProfitPrice := position.EntryPrice * (1 + s.Parameters.TS/100)
		order := engine.NewOrderLimit(s.Symbol, engine.SideSell, engine.PositionSideLong, position.Size, takeProfitPrice)
		order.IsTP = true
		return order
	case engine.PositionSideShort:
		takeProfitPrice := position.EntryPrice / (1 + s.Parameters.TS/100)
		order := engine.NewOrderLimit(s.Symbol, engine.SideBuy, engine.PositionSideShort, math.Abs(position.Size), takeProfitPrice)
		order.IsTP = true
		return order
	}
	return nil
}

func (s *StrategyLogMartingala) String() string {
	return string(s.GetType()) + " " + string(s.GetPositionSide()) + " " + s.Parameters.String()
}

// GETTERS
func (s *StrategyLogMartingala) GetType() StrategyType { return s.Type }

func (s *StrategyLogMartingala) GetSymbol() string { return s.Symbol }

func (s *StrategyLogMartingala) GetPositionSide() engine.PositionSideType { return s.PositionSide }

func (s *StrategyLogMartingala) GetStatus() string { return s.Status }

func (s *StrategyLogMartingala) GetParameters() StrategyParameters { return s.Parameters }

// SETTERS
func (s *StrategyLogMartingala) SetSymbol(value string) { s.Symbol = value }

func (s *StrategyLogMartingala) SetPositionSide(value engine.PositionSideType) {
	s.PositionSide = value
}

func (s *StrategyLogMartingala) SetStatus(value string) { s.Status = value }

func (s *StrategyLogMartingala) SetParameters(value StrategyParameters) { s.Parameters = value }

func NewStrategyLogMartingala(symbol string, positionSide engine.PositionSideType, pars StrategyParameters) *StrategyLogMartingala {
	return &StrategyLogMartingala{
		Type:         StrategyTypeLogMartingala,
		Symbol:       symbol,
		PositionSide: positionSide,
		Status:       "",
		Parameters:   pars,
	}
}
