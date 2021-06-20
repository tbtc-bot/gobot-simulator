package engine

type ExchangeAPI struct {
	PlaceOrder  func(Order)
	CancelOrder func(Order) bool
	OpenOrders  func(PositionSideType) []Order
	MarkPrice   func() float64
	Balance     func() float64
	GridReached func(PositionSideType) int64
}
