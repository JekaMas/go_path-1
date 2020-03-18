package rps

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Kmortyk/go_path/shop"
)

func TestAddProductSoftLimit(t *testing.T) {

	stub := LimitDecoratorStub{
		addProductFunc: func(product shop.Product) error {
			time.Sleep(time.Millisecond)
			return nil
		},
	}

	sd := NewSoftLimitDecorator(stub)
	wg := sync.WaitGroup{}
	const num = 2_000_000
	wg.Add(num)

	for i := 0; i < num/2; i++ {
		go func() {
			defer wg.Done()
			err := sd.AddProduct(shop.Product{})
			t.Logf("Passed %v", i)
			if err != nil {
				t.Errorf("err = %v", err)
			}
		}()
	}

	stub.addProductFunc = func(product shop.Product) error {
		return nil
	}

	for i := 0; i < num/2; i++ {
		go func() {
			defer wg.Done()
			err := sd.AddProduct(shop.Product{})
			t.Logf("Passed %v", i)
			if err != nil {
				t.Errorf("err = %v", err)
			}
		}()
	}

	wg.Wait()
}

func TestAddProductHardLimit(t *testing.T) {

	stub := LimitDecoratorStub{
		addProductFunc: func(product shop.Product) error {
			time.Sleep(time.Millisecond)
			return nil
		},
	}

	hd := NewHardLimitDecorator(&stub)

	var err error
	wg := sync.WaitGroup{}
	const num = 2_000_000
	wg.Add(num)

	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()
			err = hd.AddProduct(shop.NewProduct("a"+strconv.Itoa(i), 1, shop.ProductNormal))
		}()
	}
	wg.Wait()

	if err == nil {
		t.Error("Hard limit pass")
	}
	t.Log(err)
}

/* --- Stub --------------------------------------------------------------------------------------------------------- */

type LimitDecoratorStub struct {
	addProductFunc func(shop.Product) error
}

func (s LimitDecoratorStub) AddProduct(product shop.Product) error {
	if s.addProductFunc != nil {
		return s.addProductFunc(product)
	}
	return nil
}

func (LimitDecoratorStub) ModifyProduct(shop.Product) error {
	panic("implement me")
}

func (LimitDecoratorStub) RemoveProduct(name string) error {
	panic("implement me")
}

func (LimitDecoratorStub) GetProduct(name string) (shop.Product, error) {
	panic("implement me")
}

func (LimitDecoratorStub) SetProduct(name string, product shop.Product) error {
	panic("implement me")
}

func (LimitDecoratorStub) Register(userName string) error {
	panic("implement me")
}

func (LimitDecoratorStub) AddBalance(userName string, sum float32) error {
	panic("implement me")
}

func (LimitDecoratorStub) Balance(userName string) (float32, error) {
	panic("implement me")
}

func (LimitDecoratorStub) SetAccount(userName string, account shop.Account) error {
	panic("implement me")
}

func (LimitDecoratorStub) GetAccount(userName string) (shop.Account, error) {
	panic("implement me")
}

func (LimitDecoratorStub) GetAccounts(sortType shop.AccountSortType) []shop.Account {
	panic("implement me")
}

func (LimitDecoratorStub) ModifyAccountType(userName string, accountType shop.AccountType) error {
	panic("implement me")
}

func (LimitDecoratorStub) CalculateOrder(userName string, order shop.Order) (float32, error) {
	panic("implement me")
}

func (LimitDecoratorStub) PlaceOrder(userName string, order shop.Order) error {
	panic("implement me")
}

func (LimitDecoratorStub) AddBundle(name string, main shop.Product, discount float32, additional ...shop.Product) error {
	panic("implement me")
}

func (LimitDecoratorStub) ChangeDiscount(name string, discount float32) error {
	panic("implement me")
}

func (LimitDecoratorStub) RemoveBundle(name string) error {
	panic("implement me")
}

func (LimitDecoratorStub) SetBundle(name string, bundle shop.Bundle) error {
	panic("implement me")
}

func (LimitDecoratorStub) GetBundle(name string) (shop.Bundle, error) {
	panic("implement me")
}

func (LimitDecoratorStub) Import(data []byte) error {
	panic("implement me")
}

func (LimitDecoratorStub) ImportProductsCSV(data []byte) (errs []shop.ImportProductsError) {
	panic("implement me")
}

func (LimitDecoratorStub) ImportAccountsCSV(data []byte) (errs []shop.ImportAccountsError) {
	panic("implement me")
}

func (LimitDecoratorStub) Export() ([]byte, error) {
	panic("implement me")
}

func (LimitDecoratorStub) ExportProductsCSV() ([]byte, error) {
	panic("implement me")
}

func (LimitDecoratorStub) ExportAccountsCSV() ([]byte, error) {
	panic("implement me")
}
