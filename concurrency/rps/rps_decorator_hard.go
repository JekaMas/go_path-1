package rps

import (
	"errors"
	shop "go_path/struct"
	"sync"
)

type HardLimitDecorator struct {
	shop.Shop
	hardLock sync.Map
}

var (
	ErrorHardLimit = errors.New("hard limit error")
)

func NewHardLimitDecorator(s shop.Shop) HardLimitDecorator {
	return HardLimitDecorator{
		Shop:     s,
		hardLock: sync.Map{},
	}
}

func (d *HardLimitDecorator) hardLimitFunc(f func(chan error), name string, max int) error {
	v, _ := d.hardLock.LoadOrStore(name, make(chan struct{}, max))
	ch := v.(chan struct{})
	resCh := make(chan error, 1)

	defer func() { <-ch }()

	select {
	case ch <- struct{}{}:
	default:
		return ErrorHardLimit
	}

	f(resCh)
	return <-resCh
}

func (d *HardLimitDecorator) hardLimitFuncProduct(f func(chan productResult), name string, max int) (shop.Product, error) {
	v, _ := d.hardLock.LoadOrStore(name, make(chan struct{}, max))
	ch := v.(chan struct{})
	resCh := make(chan productResult, 1)

	defer func() { <-ch }()

	select {
	case ch <- struct{}{}:
	default:
		return shop.Product{}, ErrorHardLimit
	}

	f(resCh)

	res := <-resCh
	return res.prod, res.err
}

func (d *HardLimitDecorator) hardLimitFuncBytes(f func(chan bytesResult), name string, max int) ([]byte, error) {
	v, _ := d.hardLock.LoadOrStore(name, make(chan struct{}, max))
	ch := v.(chan struct{})
	resCh := make(chan bytesResult, 1)

	defer func() { <-ch }()

	select {
	case ch <- struct{}{}:
	default:
		return []byte{}, ErrorHardLimit
	}

	f(resCh)
	res := <-resCh
	return res.bts, res.err
}
