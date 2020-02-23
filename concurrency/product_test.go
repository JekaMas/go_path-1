package concurrency

import (
	shop "go_path/struct"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestShop_AddProductRace(t *testing.T) {

	pn := shop.ProductNormal
	pp := shop.ProductPremium
	ps := shop.ProductSampled

	products := []shop.Product{
		NewProduct("P1", 10, pn),
		NewProduct("P2", 10, pp),
		NewProduct("P3", 10, ps),
		NewProduct("P4", 10, pn),
		NewProduct("P5", 10, pp),
		NewProduct("P6", 10, ps),
	}

	m := NewMarket()

	wg := sync.WaitGroup{}
	wg.Add(len(products))

	// test
	for _, product := range products {
		go func(product shop.Product) {
			err := m.AddProduct(product)
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}(product)
	}

	wg.Wait()
}

func TestShop_ModifyProductRace(t *testing.T) {

	m := NewMarket()
	_ = m.AddProduct(shop.Product{Name: "P1", Price: 100, Type: shop.ProductNormal})
	_ = m.AddProduct(shop.Product{Name: "P2", Price: 100, Type: shop.ProductPremium})

	products := []shop.Product{
		NewProduct("P1", 10, shop.ProductNormal),
		NewProduct("P2", 10, shop.ProductPremium),
		NewProduct("P1", 20, shop.ProductNormal),
		NewProduct("P1", 30, shop.ProductPremium),
		NewProduct("P1", 20, shop.ProductNormal),
		NewProduct("P2", 40, shop.ProductPremium),
		NewProduct("P1", 10, shop.ProductSampled),
		NewProduct("P2", 100, shop.ProductPremium),
	}

	wg := sync.WaitGroup{}
	wg.Add(len(products))

	// test
	for _, product := range products {
		go func(product shop.Product) {
			err := m.ModifyProduct(product)
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}(product)
	}

	wg.Wait()
}

func TestShop_RemoveProductRace(t *testing.T) {

	m := NewMarket()
	_ = m.AddProduct(shop.Product{Name: "P1", Price: 100, Type: shop.ProductNormal})
	_ = m.AddProduct(shop.Product{Name: "P2", Price: 100, Type: shop.ProductPremium})
	_ = m.AddProduct(shop.Product{Name: "P3", Price: 100, Type: shop.ProductNormal})
	_ = m.AddProduct(shop.Product{Name: "P4", Price: 100, Type: shop.ProductPremium})
	_ = m.AddProduct(shop.Product{Name: "P5", Price: 100, Type: shop.ProductNormal})
	_ = m.AddProduct(shop.Product{Name: "P6", Price: 100, Type: shop.ProductPremium})
	_ = m.AddProduct(shop.Product{Name: "P7", Price: 100, Type: shop.ProductNormal})
	_ = m.AddProduct(shop.Product{Name: "P8", Price: 100, Type: shop.ProductPremium})

	names := []string{
		"P1", "P2", "P3",
		"P4", "P5", "P6",
		"P7", "P8",
	}

	wg := sync.WaitGroup{}
	wg.Add(len(names))

	// test
	for _, name := range names {
		go func(name string) { // rand delay?
			delay := time.Millisecond * 2 * time.Duration(rand.Int31n(300))
			// t.Log("RUN ", name, delay)
			time.Sleep(delay)
			err := m.RemoveProduct(name)
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}(name)
	}

	wg.Wait()
}
