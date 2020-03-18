package time

import (
	"github.com/Kmortyk/go_path/shop"
	"time"
)

func NewProduct(productName string, price float32, productType shop.ProductType) shop.Product {
	return shop.NewProduct(productName, price, productType)
}

func (td *TimeoutDecorator) AddProduct(p shop.Product) error {
	return td.timeoutFunc(func(ch chan error) {
		ch <- td.Shop.AddProduct(p)
	}, time.Second)
}

func (td *TimeoutDecorator) ModifyProduct(p shop.Product) error {
	return td.timeoutFunc(func(ch chan error) {
		ch <- td.Shop.ModifyProduct(p)
	}, time.Second)
}

func (td *TimeoutDecorator) RemoveProduct(name string) error {
	return td.timeoutFunc(func(ch chan error) {
		ch <- td.Shop.RemoveProduct(name)
	}, time.Second)
}

/* --- Interface ---------------------------------------------------------------------------------------------------- */

func (td *TimeoutDecorator) GetProduct(name string) (shop.Product, error) {
	return td.timeoutFuncProduct(func(ch chan productResult) {
		prod, err := td.Shop.GetProduct(name)
		ch <- productResult{prod, err}
	}, time.Second)
}

func (td *TimeoutDecorator) SetProduct(name string, product shop.Product) error {
	return td.timeoutFunc(func(ch chan error) {
		ch <- td.Shop.SetProduct(name, product)
	}, time.Millisecond)
}
