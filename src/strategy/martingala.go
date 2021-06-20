package strategy

import (
	"math"

	"example.com/gobot-simulator/src/engine"
)

type StrategyMartingala struct {
	Type         StrategyType            `json:"type"`
	Symbol       string                  `json:"symbol"`
	PositionSide engine.PositionSideType `json:"positionSide"`
	Status       string                  `json:"status"`
	Parameters   StrategyParameters      `json:"parameters"`
}

func (s *StrategyMartingala) BuyGridOrders(balance float64, startPrice float64) []*engine.Order {
	orders := []*engine.Order{}

	// start price and size
	p0 := startPrice * (1 - s.Parameters.GS/100)
	s0 := (balance / startPrice) * (s.Parameters.OS / 100)

	// first grid
	p_1 := p0 * (1 - s.Parameters.GS/100)
	s_1 := s0
	p_2 := p0
	order := engine.NewOrderLimit(s.Symbol, engine.SideBuy, engine.PositionSideLong, s_1, p_1)
	order.GridNumber = 1
	orders = append(orders, order)

	// other grids
	for i := 2; i < int(s.Parameters.GO)+1; i++ {
		p_i := p_1 - (p_2-p_1)*s.Parameters.SF
		s_i := s_1 * s.Parameters.OF
		order := engine.NewOrderLimit(s.Symbol, engine.SideBuy, engine.PositionSideLong, s_i, p_i)
		order.GridNumber = int64(i)
		orders = append(orders, order)
		p_2 = p_1
		p_1 = p_i
		s_1 = s_i
	}

	return orders
}

func (s *StrategyMartingala) SellGridOrders(balance float64, startPrice float64) []*engine.Order {
	orders := []*engine.Order{}

	// start price and size
	p0 := startPrice * (1 + s.Parameters.GS/100)
	s0 := (balance / startPrice) * (s.Parameters.OS / 100)

	// first grid
	p_1 := p0 * (1 + s.Parameters.GS/100)
	s_1 := s0
	p_2 := p0
	order := engine.NewOrderLimit(s.Symbol, engine.SideSell, engine.PositionSideShort, s_1, p_1)
	order.GridNumber = 1
	orders = append(orders, order)

	// other grids
	for i := 2; i < int(s.Parameters.GO)+1; i++ {
		p_i := p_1 + (p_1-p_2)*s.Parameters.SF
		s_i := s_1 * s.Parameters.OF
		order := engine.NewOrderLimit(s.Symbol, engine.SideSell, engine.PositionSideShort, s_i, p_i)
		order.GridNumber = int64(i)
		orders = append(orders, order)
		p_2 = p_1
		p_1 = p_i
		s_1 = s_i
	}

	return orders
}

func (s *StrategyMartingala) TakeProfitOrder(position engine.Position, currentGrid int64) *engine.Order {
	switch s.PositionSide {
	case engine.PositionSideLong:
		takeProfitPrice := position.EntryPrice * (1 + s.Parameters.TS/100)
		order := engine.NewOrderLimit(s.Symbol, engine.SideSell, engine.PositionSideLong, position.Size, takeProfitPrice)
		order.IsTP = true
		return order
	case engine.PositionSideShort:
		takeProfitPrice := position.EntryPrice * (1 - s.Parameters.TS/100)
		order := engine.NewOrderLimit(s.Symbol, engine.SideBuy, engine.PositionSideShort, math.Abs(position.Size), takeProfitPrice)
		order.IsTP = true
		return order
	default:
		return nil
	}
}

func (s *StrategyMartingala) String() string {
	return string(s.GetType()) + " " + string(s.GetPositionSide()) + " " + s.Parameters.String()
}

// GETTERS
func (s *StrategyMartingala) GetType() StrategyType { return s.Type }

func (s *StrategyMartingala) GetSymbol() string { return s.Symbol }

func (s *StrategyMartingala) GetPositionSide() engine.PositionSideType { return s.PositionSide }

func (s *StrategyMartingala) GetStatus() string { return s.Status }

func (s *StrategyMartingala) GetParameters() StrategyParameters { return s.Parameters }

// SETTERS
func (s *StrategyMartingala) SetSymbol(symbol string) { s.Symbol = symbol }

func (s *StrategyMartingala) SetPositionSide(positionSide engine.PositionSideType) {
	s.PositionSide = positionSide
}

func (s *StrategyMartingala) SetStatus(status string) { s.Status = status }

func (s *StrategyMartingala) SetParameters(pars StrategyParameters) { s.Parameters = pars }

func NewStrategyMartingala(symbol string, positionSide engine.PositionSideType, pars StrategyParameters) *StrategyMartingala {
	strat := &StrategyMartingala{
		Type:         StrategyTypeMartingala,
		Symbol:       symbol,
		PositionSide: positionSide,
		Status:       "",
		Parameters:   pars,
	}
	return strat
}
