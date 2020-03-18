package rps

import (
	"errors"
	"sync"

	"github.com/Kmortyk/go_path/shop"
)

type HardLimitDecorator struct {
	shop.Shop
	hardLock sync.Map
}

var (
	ErrorHardLimit = errors.New("hard limit error")
)

func NewHardLimitDecorator(s shop.Shop) *HardLimitDecorator {
	return &HardLimitDecorator{
		Shop: s,
	}
}

func (d *HardLimitDecorator) hardLimitFunc(f func() error, name string, max int) error {
	v, _ := d.hardLock.LoadOrStore(name, make(chan struct{}, max))
	ch, _ := v.(chan struct{})

	select {
	case ch <- struct{}{}:
		defer func() { <-ch }()
	default:
		// fix: you should not read from channel if you haven't sent to it.
		return ErrorHardLimit
	}

	return f()
}

func (d *HardLimitDecorator) hardLimitFuncProduct(f func() productResult, name string, max int) (shop.Product, error) {
	v, _ := d.hardLock.LoadOrStore(name, make(chan struct{}, max))
	ch := v.(chan struct{})

	select {
	case ch <- struct{}{}:
	default:
		defer func() { <-ch }()

		return shop.Product{}, ErrorHardLimit
	}

	res := f()
	return res.prod, res.err
}

func (d *HardLimitDecorator) hardLimitFuncBytes(f func() bytesResult, name string, max int) ([]byte, error) {
	v, _ := d.hardLock.LoadOrStore(name, make(chan struct{}, max))
	ch := v.(chan struct{})

	select {
	case ch <- struct{}{}:
	default:
		defer func() { <-ch }()
		return []byte{}, ErrorHardLimit
	}

	res := f()
	return res.bts, res.err
}
