package rps

import (
	"github.com/Kmortyk/go_path/shop"
)

const (
	getProductLimit = 10000
	addProductLimit = 5000
	modifyProductLimit = 5000
	placeOrderLimit = 5000
	importLimit = 100
	exportLimit = 1000
)

// fixme it's hard to test due to the hard limit is hard coded
func (d *HardLimitDecorator) GetProduct(name string) (shop.Product, error) {
	return d.hardLimitFuncProduct(func() productResult {
		prod, err := d.Shop.GetProduct(name)
		return productResult{prod, err}
	}, "get", getProductLimit)
}

func (d *HardLimitDecorator) AddProduct(p shop.Product) error {
	return d.hardLimitFunc(func() error {
		return d.Shop.AddProduct(p)
	}, "add", addProductLimit)
}

func (d *HardLimitDecorator) ModifyProduct(p shop.Product) error {
	return d.hardLimitFunc(func() error {
		return d.Shop.ModifyProduct(p)
	}, "add", modifyProductLimit)
}

func (d *HardLimitDecorator) PlaceOrder(userName string, o shop.Order) error {
	return d.hardLimitFunc(func() error {
		return d.Shop.PlaceOrder(userName, o)
	}, "order", placeOrderLimit)
}

func (d *HardLimitDecorator) Import(data []byte) error {
	return d.hardLimitFunc(func() error {
		return d.Shop.Import(data)
	}, "import", importLimit)
}

func (d *HardLimitDecorator) Export() ([]byte, error) {
	return d.hardLimitFuncBytes(func() bytesResult {
		bts, err := d.Shop.Export()
		return bytesResult{bts, err}
	}, "export", exportLimit)
}
