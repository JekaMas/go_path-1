package concurrency

import (
	shop "go_path/struct"
	"time"
)

func NewProduct(productName string, price float32, productType shop.ProductType) shop.Product {
	return shop.NewProduct(productName, price, productType)
}

func (m *Market) AddProduct(p shop.Product) error {

	resChan := make(chan error, 1)

	go func() {
		resChan <- m.shop.AddProduct(p)
	}()

	select {
	case res := <-resChan:
		return res
	case <-time.After(time.Second):
		return ErrorTimeout
	}
}

func (m *Market) ModifyProduct(p shop.Product) error {

	resChan := make(chan error, 1)

	go func() {
		resChan <- m.shop.ModifyProduct(p)
	}()

	select {
	case res := <-resChan:
		return res
	case <-time.After(time.Second):
		return ErrorTimeout
	}
}

func (m *Market) RemoveProduct(name string) error {

	resChan := make(chan error, 1)

	go func() {
		resChan <- m.shop.RemoveProduct(name)
	}()

	select {
	case res := <-resChan:
		return res
	case <-time.After(time.Second):
		return ErrorTimeout
	}
}

/* --- Interface ---------------------------------------------------------------------------------------------------- */

func (m *Market) GetProduct(name string) (shop.Product, error) {

	type result struct {
		product shop.Product
		err     error
	}

	resChan := make(chan result, 1)

	go func() {
		product, err := m.shop.GetProduct(name)
		resChan <- result{product, err}
	}()

	select {
	case res := <-resChan:
		return res.product, res.err
	case <-time.After(time.Second):
		return shop.Product{}, ErrorTimeout
	}
}

func (m *Market) SetProduct(name string, product shop.Product) error {

	resChan := make(chan error, 1)

	go func() {
		resChan <- m.shop.SetProduct(name, product)
	}()

	select {
	case res := <-resChan:
		return res
	case <-time.After(time.Second):
		return ErrorTimeout
	}
}
