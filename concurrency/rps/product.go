package rps

import (
	shop "go_path/struct"
)

func (d *SoftLimitDecorator) GetProduct(name string) (shop.Product, error) {
	return d.softLimitFuncProduct(func(ch chan productResult) {
		prod, err := d.Shop.GetProduct(name)
		ch <- productResult{prod, err}
	}, "get", 5000)
}

func (d *SoftLimitDecorator) AddProduct(p shop.Product) error {
	return d.softLimitFunc(func(ch chan error) {
		ch <- d.Shop.AddProduct(p)
	}, "add", 1000)
}

func (d *SoftLimitDecorator) ModifyProduct(p shop.Product) error {
	return d.softLimitFunc(func(ch chan error) {
		ch <- d.Shop.ModifyProduct(p)
	}, "add", 1000)
}

func (d *SoftLimitDecorator) PlaceOrder(userName string, o shop.Order) error {
	return d.softLimitFunc(func(ch chan error) {
		ch <- d.Shop.PlaceOrder(userName, o)
	}, "order", 4000)
}

func (d *SoftLimitDecorator) Import(data []byte) error {
	return d.softLimitFunc(func(ch chan error) {
		ch <- d.Shop.Import(data)
	}, "import", 1)
}

func (d *SoftLimitDecorator) Export() ([]byte, error) {
	return d.softLimitFuncBytes(func(ch chan bytesResult) {
		bts, err := d.Shop.Export()
		ch <- bytesResult{bts, err}
	}, "export", 10)
}
