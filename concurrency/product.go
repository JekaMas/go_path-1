package concurrency

import (
	shop "go_path/struct"
	"time"
)

func NewProduct(productName string, price float32, productType shop.ProductType) shop.Product {
	return shop.NewProduct(productName, price, productType)
}

func (td *TimeoutDecorator) AddProduct(p shop.Product) error {
	return td.timeoutFunc(func(ch chan error) {
		ch <- td.shop.AddProduct(p)
	}, time.Second)
}

func (td *TimeoutDecorator) ModifyProduct(p shop.Product) error {
	return td.timeoutFunc(func(ch chan error) {
		ch <- td.shop.ModifyProduct(p)
	}, time.Second)
}

func (td *TimeoutDecorator) RemoveProduct(name string) error {
	return td.timeoutFunc(func(ch chan error) {
		ch <- td.shop.RemoveProduct(name)
	}, time.Second)
}

/* --- Interface ---------------------------------------------------------------------------------------------------- */

func (td *TimeoutDecorator) GetProduct(name string) (shop.Product, error) {
	return td.timeoutFuncProduct(func(ch chan productResult) {
		prod, err := td.shop.GetProduct(name)
		ch <- productResult{prod, err}
	}, time.Second)
}

func (td *TimeoutDecorator) SetProduct(name string, product shop.Product) error {
	return td.timeoutFunc(func(ch chan error) {
		ch <- td.shop.SetProduct(name, product)
	}, time.Second)
}
