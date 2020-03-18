package rps

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/Kmortyk/go_path/concurrency/cancel"
	"github.com/Kmortyk/go_path/concurrency/rps"
	"github.com/Kmortyk/go_path/concurrency/time"
	"github.com/Kmortyk/go_path/shop"
)

func TestDebugAddProductHardLimit(t *testing.T) {
	productTypes := [...]shop.ProductType{
		shop.ProductNormal,
		shop.ProductPremium,
		shop.ProductSampled,
	}

	const n = 8000 //5000
	products := make([]shop.Product, n)

	for i := range products {
		products[i].Type = productTypes[i%len(productTypes)]
		products[i].Name = "p" + strconv.Itoa(rand.Int())
		products[i].Price = float32(rand.Intn(100_000))
	}

	/*
		hd := rps.NewHardLimitDecorator(
			cancel.NewCancelDecorator(
				time.NewTimeoutDecorator()))
	*/

	hd := rps.NewHardLimitDecorator(shop.NewMarket())

	wg := sync.WaitGroup{}
	wg.Add(len(products))

	// test
	for _, product := range products {
		go func(product shop.Product) {
			defer func() {
				wg.Done()
			}()

			hd.AddProduct(product)
		}(product)
	}

	wg.Wait()
}

func TestAddProductAllDecorators(t *testing.T) {
	productTypes := [...]shop.ProductType{
		shop.ProductNormal,
		shop.ProductPremium,
		shop.ProductSampled,
	}

	const n = 8000 //5000
	products := make([]shop.Product, n)

	for i := range products {
		products[i].Type = productTypes[i%len(productTypes)]
		products[i].Name = "p" + strconv.Itoa(rand.Int())
		products[i].Price = float32(rand.Intn(100_000))
	}

	hd := rps.NewHardLimitDecorator(
		cancel.NewCancelDecorator(
			time.NewTimeoutDecorator()))

	wg := sync.WaitGroup{}
	wg.Add(len(products))

	// test
	for _, product := range products {
		go func(product shop.Product) {
			defer func() {
				wg.Done()
			}()

			hd.AddProduct(product)
		}(product)
	}

	wg.Wait()
}
