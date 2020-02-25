package concurrency

import (
	shop "go_path/struct"
	"sync"
	"testing"
)

func TestShop_CalculateOrderRace(t *testing.T) {

	type order struct {
		accountName string
		order       shop.Order
	}

	m := NewMarket()

	accounts := []shop.Account{
		{"A", 0, shop.AccountNormal},
		{"B", 0, shop.AccountPremium},
	}

	// reg
	for _, a := range accounts {
		_ = m.Register(a.Name)
		_ = m.ModifyAccountType(a.Name, a.Type)
		_ = m.AddBalance(a.Name, a.Balance)
	}

	// tests
	orders := []order{
		{
			"A",
			NewOrder([]shop.Product{NewProduct("P", 90, shop.ProductNormal)}, nil),
		},
		{
			"A",
			NewOrder([]shop.Product{
				NewProduct("P", 1000, shop.ProductNormal),
				NewProduct("P", 100, shop.ProductNormal),
				NewProduct("P", 10, shop.ProductNormal),
				NewProduct("P", 9, shop.ProductNormal),
			}, nil),
		},
		{
			"B",
			NewOrder([]shop.Product{
				NewProduct("P", 9246, shop.ProductPremium),
				NewProduct("P", 932, shop.ProductPremium),
				NewProduct("P", 921, shop.ProductPremium),
				NewProduct("P", 9357, shop.ProductPremium),
			}, nil),
		},
		{
			"A",
			NewOrder([]shop.Product{
				NewProduct("P", 111, shop.ProductNormal),
				NewProduct("P", 11, shop.ProductPremium),
				NewProduct("P", 111, shop.ProductNormal),
				NewProduct("P", 111, shop.ProductNormal),
			}, nil),
		},
	}

	orders = append(orders, orders...)

	wg := sync.WaitGroup{}
	wg.Add(len(orders))

	for _, o := range orders { // fixme fall with timeout
		go func(o order) {
			defer wg.Done()
			_, err := m.CalculateOrder(o.accountName, o.order)
			if err != nil {
				t.Error(err)
			}
		}(o)
	}

	wg.Wait()
}

func TestShop_PlaceOrder(t *testing.T) {

	type order struct {
		accountName string
		order       shop.Order
	}

	m := NewMarket()

	accounts := []shop.Account{
		{"A", 100_000_000, shop.AccountNormal},
		{"B", 100_000, shop.AccountPremium},
	}

	// reg
	for _, a := range accounts {
		_ = m.Register(a.Name)
		_ = m.ModifyAccountType(a.Name, a.Type)
		_ = m.AddBalance(a.Name, a.Balance)
	}

	// tests
	orders := []order{
		{
			"A",
			NewOrder([]shop.Product{NewProduct("P", 90, shop.ProductNormal)}, nil),
		},
		{
			"A",
			NewOrder([]shop.Product{
				NewProduct("P", 1000, shop.ProductNormal),
				NewProduct("P", 100, shop.ProductNormal),
				NewProduct("P", 10, shop.ProductNormal),
				NewProduct("P", 9, shop.ProductNormal),
			}, nil),
		},
		{
			"B",
			NewOrder([]shop.Product{
				NewProduct("P", 9246, shop.ProductPremium),
				NewProduct("P", 932, shop.ProductPremium),
				NewProduct("P", 921, shop.ProductPremium),
				NewProduct("P", 9357, shop.ProductPremium),
			}, nil),
		},
		{
			"A",
			NewOrder([]shop.Product{
				NewProduct("P", 111, shop.ProductNormal),
				NewProduct("P", 11, shop.ProductPremium),
				NewProduct("P", 111, shop.ProductNormal),
				NewProduct("P", 111, shop.ProductNormal),
			}, nil),
		},
	}

	orders = append(orders, orders...)

	wg := sync.WaitGroup{}
	wg.Add(len(orders))

	for _, o := range orders { // fixme fall with timeout
		go func(o order) {
			defer wg.Done()
			err := m.PlaceOrder(o.accountName, o.order)
			if err != nil {
				t.Error(err)
			}
		}(o)
	}

	wg.Wait()
}
