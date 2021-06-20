package engine

import (
	"fmt"
	"math"

	"example.com/gobot-simulator/src/common"
	log "github.com/sirupsen/logrus"
)

type PositionSideType string

const (
	PositionSideLong  PositionSideType = "LONG"
	PositionSideShort PositionSideType = "SHORT"
)

type Position struct {
	Symbol       string           `json:"symbol"`
	PositionSide PositionSideType `json:"positionSide"`
	EntryPrice   float64          `json:"entryPrice"`
	Size         float64          `json:"size"`
	MarkPrice    float64          `json:"markPrice"`
}

func (p *Position) Update(order Order) float64 {
	if order.PositionSide != p.PositionSide {
		log.Panic("Order position side is different from position side")
	}

	// Compute entry price
	var orderPrice float64
	if order.Type == OrderTypeMarket {
		orderPrice = p.MarkPrice
	} else {
		orderPrice = order.Price
	}
	p.EntryPrice = (p.EntryPrice*p.Size + orderPrice*order.Amount) / (p.Size + order.Amount)
	if math.IsNaN(p.EntryPrice) || math.IsInf(p.EntryPrice, 0) {
		log.Panic("Invalid entry price")
	}
	p.EntryPrice = common.RoundFloatWithPrecision(p.EntryPrice, 6)

	// Update position size
	if (p.PositionSide == PositionSideLong && order.Side == SideBuy) || (p.PositionSide == PositionSideShort && order.Side == SideSell) {
		p.Size += order.Amount
		p.Size = common.RoundFloatWithPrecision(p.Size, 6)
		return 0 // increase position: no realized profit
	} else if (p.PositionSide == PositionSideLong && order.Side == SideSell) || (p.PositionSide == PositionSideShort && order.Side == SideBuy) {
		pnl := p.PNL(order.Price)
		p.Size -= order.Amount
		p.Size = common.RoundFloatWithPrecision(p.Size, 6)
		if p.Size < 0 {
			log.Panic("Position size is negative")
		}
		return pnl // reduce position: return realized profit
	} else {
		log.Panic("Unexpected combination of position and order side")
		return -1
	}
}

func (p *Position) PNL(markPrice float64) float64 {
	if p.PositionSide == PositionSideLong {
		return (markPrice - p.EntryPrice) * p.Size
	} else {
		return (-markPrice + p.EntryPrice) * p.Size
	}
}

func (p *Position) String() string {
	return fmt.Sprintf("symbol %s, position side %s, entry price %.4f, size %.4f",
		p.Symbol, string(p.PositionSide), p.EntryPrice, p.Size)
}
