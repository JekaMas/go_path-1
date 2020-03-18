package time

import (
	"github.com/Kmortyk/go_path/shop"
	"time"
)

func NewOrder(products []shop.Product, bundles []shop.Bundle) shop.Order {
	return shop.NewOrder(products, bundles)
}

func (td *TimeoutDecorator) CalculateOrder(userName string, order shop.Order) (float32, error) {
	return td.timeoutFuncAmount(func(ch chan amountResult) {
		sum, err := td.Shop.CalculateOrder(userName, order)
		ch <- amountResult{sum, err}
	}, time.Millisecond*10)
}

func (td *TimeoutDecorator) PlaceOrder(userName string, order shop.Order) error {
	return td.timeoutFunc(func(ch chan error) {
		ch <- td.Shop.PlaceOrder(userName, order)
	}, time.Millisecond*10)
}
