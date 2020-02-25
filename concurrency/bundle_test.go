package concurrency

import (
	shop "go_path/struct"
	"sync"
	"testing"
	"time"
)

func TestShop_AddBundleRace(t *testing.T) {

	type test struct {
		main       shop.Product
		discount   float32
		additional []shop.Product
	}

	// alias
	pn := shop.ProductNormal
	pp := shop.ProductPremium
	ps := shop.ProductSampled

	// bundles array
	bundles := []test{
		{NewProduct("P1", 10, pn), 1, []shop.Product{NewProduct("P2", 90, pn)}},
		{NewProduct("P2", 10, pn), 1, []shop.Product{NewProduct("P2", 90, pn)}},
		{NewProduct("P3", 10, pn), 1, []shop.Product{NewProduct("P2", 90, pn)}},
		{NewProduct("P4", 10, pn), 1, []shop.Product{NewProduct("P2", 90, pn)}},
		{NewProduct("P5", 10, pn), 1, []shop.Product{NewProduct("P2", 90, pp)}},
		{NewProduct("P6", 10, pn), 1, []shop.Product{NewProduct("P2", 90, ps)}},
	}

	m := NewMarket()

	wg := sync.WaitGroup{}
	wg.Add(len(bundles))

	// test
	for _, b := range bundles {
		go func(main shop.Product, discount float32, additional []shop.Product) {
			defer wg.Done()
			err := m.AddBundle(main.Name, main, discount, additional...)
			if err != nil {
				t.Error(err)
			}
		}(b.main, b.discount, b.additional)
	}

	wg.Wait()
}

func TestShop_ChangeDiscountRace(t *testing.T) {
	type test struct {
		bundleName string
		discount   float32
	}

	m := NewMarket()
	_ = m.AddBundle("B1", NewProduct("P1", 10, shop.ProductNormal), 1, NewProduct("P2", 90, shop.ProductPremium))
	_ = m.AddBundle("B2", NewProduct("P1", 10, shop.ProductNormal), 1, NewProduct("P2", 90, shop.ProductSampled))

	tests := []test{
		{"B1", 1},
		{"B1", 2},
		{"B2", 3},
		{"B1", 4},
		{"B1", 5},
		{"B2", 6},
		{"B2", 7},
		{"B2", 8},
	}

	wg := sync.WaitGroup{}
	wg.Add(len(tests))

	// test
	for _, tt := range tests {
		go func(name string, discount float32) {
			defer wg.Done()
			err := m.ChangeDiscount(name, discount)
			if err != nil {
				t.Error(err)
			}
		}(tt.bundleName, tt.discount)
	}

	wg.Wait()
}

func TestShop_RemoveBundleRace(t *testing.T) {

	m := NewMarket()
	_ = m.AddBundle("B1", NewProduct("P1", 10, shop.ProductNormal), 1, NewProduct("P2", 90, shop.ProductPremium))
	_ = m.AddBundle("B2", NewProduct("P1", 10, shop.ProductNormal), 1, NewProduct("P2", 90, shop.ProductSampled))
	_ = m.AddBundle("B3", NewProduct("P1", 10, shop.ProductNormal), 1, NewProduct("P2", 90, shop.ProductPremium))
	_ = m.AddBundle("B4", NewProduct("P1", 10, shop.ProductNormal), 1, NewProduct("P2", 90, shop.ProductSampled))
	_ = m.AddBundle("B5", NewProduct("P1", 10, shop.ProductNormal), 1, NewProduct("P2", 90, shop.ProductPremium))
	_ = m.AddBundle("B6", NewProduct("P1", 10, shop.ProductNormal), 1, NewProduct("P2", 90, shop.ProductSampled))

	names := []string{
		"B1", "B2",
		"B3", "B4",
		"B5", "B6",
	}

	wg := sync.WaitGroup{}
	wg.Add(len(names))

	// test
	for _, name := range names {
		go func(name string) {
			defer wg.Done()
			time.Sleep(time.Millisecond)
			err := m.RemoveBundle(name)
			if err != nil {
				t.Error(err)
			}
		}(name)
	}

	wg.Wait()
}
