package time

import (
	"github.com/Kmortyk/go_path/shop"
	"time"
)

func (td *TimeoutDecorator) AddBundle(name string, main shop.Product, discount float32, additional ...shop.Product) error {
	return td.timeoutFunc(func(ch chan error) {
		ch <- td.Shop.AddBundle(name, main, discount, additional...)
	}, time.Second)
}

func (td *TimeoutDecorator) ChangeDiscount(name string, discount float32) error {
	return td.timeoutFunc(func(ch chan error) {
		ch <- td.Shop.ChangeDiscount(name, discount)
	}, time.Second)
}

func (td *TimeoutDecorator) RemoveBundle(name string) error {
	return td.timeoutFunc(func(ch chan error) {
		ch <- td.Shop.RemoveBundle(name)
	}, time.Second)
}

/* --- Interface ---------------------------------------------------------------------------------------------------- */

func (td *TimeoutDecorator) GetBundle(name string) (shop.Bundle, error) {
	return td.timeoutFuncBundle(func(ch chan bundleResult) {
		bun, err := td.Shop.GetBundle(name)
		ch <- bundleResult{bun, err}
	}, time.Second)
}

func (td *TimeoutDecorator) SetBundle(name string, bundle shop.Bundle) error {
	return td.timeoutFunc(func(ch chan error) {
		ch <- td.Shop.SetBundle(name, bundle)
	}, time.Second)
}
