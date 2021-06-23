package strategy

import (
	"fmt"

	"example.com/gobot-simulator/src/engine"
)

type StrategyType string

const (
	StrategyTypeMartingala     StrategyType = "Martingala"
	StrategyTypeLogMartingala  StrategyType = "LogMartingala"
	StrategyTypeAntiMartingala StrategyType = "AntiMartingala"
)

type StrategyParameters struct {
	GO uint    `json:"GO"`
	GS float64 `json:"GS"`
	SF float64 `json:"SF"`
	OS float64 `json:"OS"`
	OF float64 `json:"OF"`
	TS float64 `json:"TS"`
	SL float64 `json:"SL"`
}

func (sp StrategyParameters) String() string {
	return fmt.Sprintf("GO %d, GS %.2f, SF %.2f, OS %.2f, OF %.2f, TS %.2f, SL %.2f", sp.GO, sp.GS, sp.SF, sp.OS, sp.OF, sp.TS, sp.SL)
}

type StrategyWrapper interface {
	GetType() StrategyType
	GetSymbol() string
	GetPositionSide() engine.PositionSideType
	GetStatus() string
	GetParameters() StrategyParameters
	SetSymbol(symbol string)
	SetPositionSide(positionSide engine.PositionSideType)
	SetStatus(status string)
	SetParameters(pars StrategyParameters)
	BuyGridOrders(balance float64, startPrice float64) []*engine.Order
	SellGridOrders(balance float64, startPrice float64) []*engine.Order
	TakeProfitOrder(position engine.Position, currentGrid int64) *engine.Order
	String() string
}

type Strategy struct {
	Type StrategyType `json:"type"`
}

func NewStrategy(strategyType StrategyType, symbol string, positionSide engine.PositionSideType, pars StrategyParameters) *StrategyWrapper {
	var strategy StrategyWrapper
	switch strategyType {
	case StrategyTypeMartingala:
		strategy = NewStrategyMartingala(symbol, positionSide, pars)
	case StrategyTypeLogMartingala:
		strategy = NewStrategyLogMartingala(symbol, positionSide, pars)
	case StrategyTypeAntiMartingala:
		strategy = NewStrategyAntiMartingala(symbol, positionSide, pars)
	}
	return &strategy
}
