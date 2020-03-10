package time

import (
	shop "go_path/struct"
	"time"
)

func NewOrder(products []shop.Product, bundles []shop.Bundle) shop.Order {
	return shop.NewOrder(products, bundles)
}

func (td *TimeoutDecorator) CalculateOrder(userName string, order shop.Order) (float32, error) {
	return td.timeoutFuncAmount(func(ch chan amountResult) {
		sum, err := td.shop.CalculateOrder(userName, order)
		ch <- amountResult{sum, err}
	}, time.Millisecond*10)
}

func (td *TimeoutDecorator) PlaceOrder(userName string, order shop.Order) error {
	return td.timeoutFunc(func(ch chan error) {
		ch <- td.shop.PlaceOrder(userName, order)
	}, time.Millisecond*10)
}
