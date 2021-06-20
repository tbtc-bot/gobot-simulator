package engine

type Session struct {
	position       Position
	openOrders     map[string]Order
	orderAmount    float64
	orderPrice     float64
	realizedProfit float64
	gridReached    int64
}

func NewSession(positionSide PositionSideType) *Session {
	return &Session{
		position:   Position{PositionSide: positionSide},
		openOrders: make(map[string]Order),
	}
}
