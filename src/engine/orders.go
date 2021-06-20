package engine

import "fmt"

type OrderType string
type SideType string

const (
	OrderTypeLimit    OrderType = "LIMIT"
	OrderTypeMarket   OrderType = "MARKET"
	OrderTypeStop     OrderType = "STOP"
	OrderTypeTrailing OrderType = "TRAILING"

	SideBuy  SideType = "BUY"
	SideSell SideType = "SELL"
)

type Order struct {
	Type         OrderType        `json:"type"`
	Symbol       string           `json:"symbol"`
	Side         SideType         `json:"side"`
	PositionSide PositionSideType `json:"positionSide"`
	Amount       float64          `json:"amount"`
	ID           string           `json:"id"`
	Price        float64          `json:"price"`
	TriggerPrice float64          `json:"stopPrice"`
	CallbackRate float64          `json:"callbackRate"`
	GridNumber   int64            `json:"gridNumber"`
	IsTP         bool             `json:"isTP"`
}

func (order Order) String() string {
	switch order.Type {
	case OrderTypeLimit:
		return fmt.Sprintf("symbol %s, side %s, position side %s, price %.2f, amount %.2f", order.Symbol, order.Side, order.PositionSide, order.Price, order.Amount)
	case OrderTypeMarket:
		return fmt.Sprintf("symbol %s, side %s, position side %s, amount %.2f", order.Symbol, order.Side, order.PositionSide, order.Amount)
	case OrderTypeStop:
		return fmt.Sprintf("symbol %s, side %s, position side %s, trigger price %.2f, price %.2f, amount %.2f",
			order.Symbol, order.Side, order.PositionSide, order.TriggerPrice, order.Price, order.Amount)
	case OrderTypeTrailing:
		return fmt.Sprintf("symbol %s, side %s, position side %s, trigger price %.2f, amount %.2f, callback rate %.2f",
			order.Symbol, order.Side, order.PositionSide, order.TriggerPrice, order.Amount, order.CallbackRate)
	default:
		return "order type not recognized"
	}
}

func NewOrderLimit(symbol string, side SideType, positionSide PositionSideType, amount float64, price float64) *Order {
	return &Order{
		Type:         OrderTypeLimit,
		Symbol:       symbol,
		Side:         side,
		PositionSide: positionSide,
		Amount:       amount,
		Price:        price,
		ID:           "",
		IsTP:         false,
	}
}

func NewOrderMarket(symbol string, side SideType, positionSide PositionSideType, amount float64) *Order {
	return &Order{
		Type:         OrderTypeMarket,
		Symbol:       symbol,
		Side:         side,
		PositionSide: positionSide,
		Amount:       amount,
		IsTP:         false,
	}
}

func NewOrderStop(symbol string, side SideType, positionSide PositionSideType, amount float64, price float64, stopPrice float64) *Order {
	return &Order{
		Type:         OrderTypeStop,
		Symbol:       symbol,
		Side:         side,
		PositionSide: positionSide,
		Amount:       amount,
		Price:        price,
		TriggerPrice: stopPrice,
		ID:           "",
		IsTP:         false,
	}
}

func NewOrderTrailing(symbol string, side SideType, positionSide PositionSideType, amount float64, activationPrice float64, callbackRate float64) *Order {
	return &Order{
		Type:         OrderTypeStop,
		Symbol:       symbol,
		Side:         side,
		PositionSide: positionSide,
		Amount:       amount,
		TriggerPrice: activationPrice,
		CallbackRate: callbackRate,
		ID:           "",
		IsTP:         true,
	}
}
