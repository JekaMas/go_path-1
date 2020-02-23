package concurrency

import (
	shop "go_path/struct"
	"time"
)

func NewOrder(products []shop.Product, bundles []shop.Bundle) shop.Order {
	return shop.NewOrder(products, bundles)
}

func (m *Market) CalculateOrder(userName string, order shop.Order) (float32, error) {

	type result struct {
		sum float32
		err error
	}

	resChan := make(chan result, 1)

	go func() {
		sum, err := m.shop.CalculateOrder(userName, order)
		resChan <- result{sum, err}
	}()

	select {
	case res := <-resChan:
		return res.sum, res.err
	case <-time.After(time.Second):
		return 0, ErrorTimeout
	}
}

func (m *Market) PlaceOrder(userName string, order shop.Order) error {

	resChan := make(chan error, 1)

	go func() {
		resChan <- m.shop.PlaceOrder(userName, order)
	}()

	select {
	case res := <-resChan:
		return res
	case <-time.After(time.Second):
		return ErrorTimeout
	}
}
