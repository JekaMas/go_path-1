package concurrency

import (
	shop "go_path/struct"
	"time"
)

func (m *Market) AddBundle(name string, main shop.Product, discount float32, additional ...shop.Product) error {

	resChan := make(chan error, 1)

	go func() {
		resChan <- m.shop.AddBundle(name, main, discount, additional...)
	}()

	select {
	case res := <-resChan:
		return res
	case <-time.After(time.Second):
		return ErrorTimeout
	}
}

func (m *Market) ChangeDiscount(name string, discount float32) error {

	resChan := make(chan error, 1)

	go func() {
		resChan <- m.shop.ChangeDiscount(name, discount)
	}()

	select {
	case res := <-resChan:
		return res
	case <-time.After(time.Second):
		return ErrorTimeout
	}
}

func (m *Market) RemoveBundle(name string) error {

	resChan := make(chan error, 1)

	go func() {
		resChan <- m.shop.RemoveBundle(name)
	}()

	select {
	case res := <-resChan:
		return res
	case <-time.After(time.Second):
		return ErrorTimeout
	}
}

/* --- Interface ---------------------------------------------------------------------------------------------------- */

func (m *Market) GetBundle(name string) (shop.Bundle, error) {

	type result struct {
		bundle shop.Bundle
		err    error
	}

	resChan := make(chan result, 1)

	go func() {
		bundle, err := m.shop.GetBundle(name)
		resChan <- result{bundle, err}
	}()

	select {
	case res := <-resChan:
		return res.bundle, res.err
	case <-time.After(time.Second):
		return shop.Bundle{}, ErrorTimeout
	}
}

func (m *Market) SetBundle(name string, bundle shop.Bundle) error {

	resChan := make(chan error, 1)

	go func() {
		resChan <- m.shop.SetBundle(name, bundle)
	}()

	select {
	case res := <-resChan:
		return res
	case <-time.After(time.Second):
		return ErrorTimeout
	}
}
