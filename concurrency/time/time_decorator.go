package time

import (
	"errors"
	shop "go_path/struct"
	"time"
)

var (
	ErrorTimeout = errors.New("timeout")
)

type amountResult struct {
	sum float32
	err error
}

type accountResult struct {
	acc shop.Account
	err error
}

type productResult struct {
	prod shop.Product
	err  error
}

type bundleResult struct {
	bun shop.Bundle
	err error
}

type TimeoutDecorator struct {
	shop shop.Shop
}

func NewTimeoutDecorator() TimeoutDecorator {
	m := shop.NewMarket()
	return TimeoutDecorator{&m}
}

func (td *TimeoutDecorator) timeoutFunc(f func(ch chan error), timeout time.Duration) error {
	ch := make(chan error, 1)
	go f(ch)

	select {
	case err := <-ch:
		return err
	case <-time.After(timeout):
		return ErrorTimeout
	}
}

func (td *TimeoutDecorator) timeoutFuncAmount(f func(ch chan amountResult), timeout time.Duration) (float32, error) {
	ch := make(chan amountResult, 1)
	go f(ch)

	select {
	case res := <-ch:
		return res.sum, res.err
	case <-time.After(timeout):
		return 0, ErrorTimeout
	}
}

func (td *TimeoutDecorator) timeoutFuncAccounts(f func(ch chan []shop.Account), timeout time.Duration) []shop.Account {
	ch := make(chan []shop.Account, 1)
	go f(ch)

	select {
	case res := <-ch:
		return res
	case <-time.After(timeout):
		return []shop.Account{} // empty list
	}
}

func (td *TimeoutDecorator) timeoutFuncAccount(f func(ch chan accountResult), timeout time.Duration) (shop.Account, error) {
	ch := make(chan accountResult, 1)
	go f(ch)

	select {
	case res := <-ch:
		return res.acc, res.err
	case <-time.After(timeout):
		return shop.Account{}, ErrorTimeout // empty list
	}
}

func (td *TimeoutDecorator) timeoutFuncProduct(f func(ch chan productResult), timeout time.Duration) (shop.Product, error) {
	ch := make(chan productResult, 1)
	go f(ch)

	select {
	case res := <-ch:
		return res.prod, res.err
	case <-time.After(timeout):
		return shop.Product{}, ErrorTimeout // empty list
	}
}

func (td *TimeoutDecorator) timeoutFuncBundle(f func(ch chan bundleResult), timeout time.Duration) (shop.Bundle, error) {
	ch := make(chan bundleResult, 1)
	go f(ch)

	select {
	case res := <-ch:
		return res.bun, res.err
	case <-time.After(timeout):
		return shop.Bundle{}, ErrorTimeout // empty list
	}
}
