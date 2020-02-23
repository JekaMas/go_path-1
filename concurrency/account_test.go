package concurrency

import (
	shop "go_path/struct"
	"sync"
	"testing"
)

/**
Note: run this tests only with

			--race

	  flag.
*/

/* -- AccountManager ------------------------------------------------------------------------------------------------ */

func TestShop_RegisterRace(t *testing.T) {

	names := []string{
		"Spike",
		"Walter",
		"Stefan",
		"Colin",
		"Mary",
		"Ann",
	}

	m := NewMarket()
	wg := sync.WaitGroup{}

	wg.Add(len(names))

	for _, name := range names {
		// add each name in goroutine
		go func(name string) {
			err := m.Register(name)
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}(name) // copy value
	}

	wg.Wait()
}

func TestShop_BalanceRace(t *testing.T) {

	m := NewMarket()
	_ = m.Register("John")
	_ = m.Register("Stan")
	_ = m.AddBalance("Stan", 100)

	names := []string{
		"John", "John",
		"Stan",
		"John", "John", "John",
		"Stan",
	}

	wg := sync.WaitGroup{}

	wg.Add(len(names))

	for _, name := range names {
		go func(name string) {
			_, err := m.Balance(name)
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}(name)
	}

	wg.Wait()
}

func TestShop_AddBalanceRace(t *testing.T) {

	type test struct {
		name string
		sum  float32
	}

	m := NewMarket()
	_ = m.Register("John")
	_ = m.Register("Stan")
	_ = m.AddBalance("Stan", 100)

	tests := []test{
		{"John", 1},
		{"John", 1},
		{"John", 1},
		{"Stan", 1},
		{"Stan", 1},
		{"John", 1},
		{"John", 1},
		{"John", 1},
		{"Stan", 1},
		{"John", 1},
	}

	wg := sync.WaitGroup{}

	wg.Add(len(tests))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func(name string, sum float32) {
				err := m.AddBalance(name, sum)
				if err != nil {
					t.Error(err) // fixme race in tests lol
				}
				wg.Done()
			}(tt.name, tt.sum)
		})
	}

	wg.Wait()
}

func TestShop_GetAccountsRace(t *testing.T) {

	types := []shop.AccountSortType{
		shop.SortByBalance,
		shop.SortByNameReverse,
		shop.SortByName,
		shop.SortByBalance,
		shop.SortByName,
		shop.SortByName,
	}

	m := NewMarket()

	// init
	_ = m.Register("John")
	_ = m.Register("Tom")
	_ = m.Register("Stan")

	_ = m.AddBalance("John", 3)
	_ = m.AddBalance("Tom", 1)
	_ = m.AddBalance("Stan", 2)

	wg := sync.WaitGroup{}
	wg.Add(len(types))

	// test
	for _, typ := range types {
		go func(typ shop.AccountSortType) {
			accs := m.GetAccounts(typ)
			if len(accs) == 0 {
				t.Error("empty accounts")
			}
			wg.Done()
		}(typ)
	}

	wg.Wait()
}
