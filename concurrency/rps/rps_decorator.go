package rps

import (
	"errors"
	"github.com/Kmortyk/go_path/shop"
	"sync"
)

type SoftLimitDecorator struct {
	shop.Shop
	softLock sync.Map
}

func NewSoftLimitDecorator(s shop.Shop) *SoftLimitDecorator {
	return &SoftLimitDecorator{
		Shop:     s,
		softLock: sync.Map{},
	}
}

var (
	ErrorInvalidQueueMap = errors.New("can't create channel in the sync map")
)

type productResult struct {
	prod shop.Product
	err  error
}

type bytesResult struct {
	bts []byte
	err error
}

func (d *SoftLimitDecorator) softLimitFunc(f func(chan error), name string, max int) error {
	v, _ := d.softLock.LoadOrStore(name, make(chan struct{}, max))
	ch := v.(chan struct{})
	resCh := make(chan error, 1)

	go func() { // not block the main thread
		ch <- struct{}{} // try add yourself to the queue
		defer func() { <-ch }()

		f(resCh)
	}()

	return <-resCh
}

func (d *SoftLimitDecorator) softLimitFuncProduct(f func(chan productResult), name string, max int) (shop.Product, error) {
	v, ok := d.softLock.LoadOrStore(name, make(chan struct{}, max))
	ch := v.(chan struct{})
	resCh := make(chan productResult, 1)
	if !ok {
		return shop.Product{}, ErrorInvalidQueueMap
	}

	go func() { // not block the main thread
		ch <- struct{}{} // try add yourself to the queue
		defer func() { <-ch }()

		f(resCh)
	}()

	res := <-resCh
	return res.prod, res.err
}

func (d *SoftLimitDecorator) softLimitFuncBytes(f func(chan bytesResult), name string, max int) ([]byte, error) {
	v, ok := d.softLock.LoadOrStore(name, make(chan struct{}, max))
	ch := v.(chan struct{})
	resCh := make(chan bytesResult, 1)
	if !ok {
		return []byte{}, ErrorInvalidQueueMap
	}

	go func() { // not block the main thread
		ch <- struct{}{} // try add yourself to the queue
		defer func() { <-ch }()

		f(resCh)
	}()

	res := <-resCh
	return res.bts, res.err
}
