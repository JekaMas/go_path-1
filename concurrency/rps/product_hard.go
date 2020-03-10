package rps

import (
	shop "go_path/struct"
)

func (d *HardLimitDecorator) GetProduct(name string) (shop.Product, error) {
	return d.hardLimitFuncProduct(func(ch chan productResult) {
		prod, err := d.Shop.GetProduct(name)
		ch <- productResult{prod, err}
	}, "get", 10000)
}

func (d *HardLimitDecorator) AddProduct(p shop.Product) error {
	return d.hardLimitFunc(func(ch chan error) {
		ch <- d.Shop.AddProduct(p)
	}, "add", 5000)
}

func (d *HardLimitDecorator) ModifyProduct(p shop.Product) error {
	return d.hardLimitFunc(func(ch chan error) {
		ch <- d.Shop.ModifyProduct(p)
	}, "add", 5000)
}

func (d *HardLimitDecorator) PlaceOrder(userName string, o shop.Order) error {
	return d.hardLimitFunc(func(ch chan error) {
		ch <- d.Shop.PlaceOrder(userName, o)
	}, "order", 10000)
}

func (d *HardLimitDecorator) Import(data []byte) error {
	return d.hardLimitFunc(func(ch chan error) {
		ch <- d.Shop.Import(data)
	}, "import", 100)
}

func (d *HardLimitDecorator) Export() ([]byte, error) {
	return d.hardLimitFuncBytes(func(ch chan bytesResult) {
		bts, err := d.Shop.Export()
		ch <- bytesResult{bts, err}
	}, "export", 1000)
}
