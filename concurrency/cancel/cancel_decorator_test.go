package cancel

import (
	shop "go_path/struct"
	"sync"
	"testing"
	"time"
)

func TestCancelDecorator_Import(t *testing.T) {

	stub := CancelDecoratorStub{
		exportFunc: func() (bytes []byte, err error) {
			time.Sleep(time.Second)
			return []byte{}, nil
		},
	}

	hd := NewCancelDecorator(&stub)

	var err error
	wg := sync.WaitGroup{}
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			_, err = hd.Export()
		}()
	}

	_ = hd.CancelImportExport()
	time.Sleep(time.Second)

	if err == nil {
		t.Error("Cancel pass")
	}
	t.Log(err)

}

/* --- Stub --------------------------------------------------------------------------------------------------------- */

type CancelDecoratorStub struct {
	exportFunc func() ([]byte, error)
}

func (s CancelDecoratorStub) AddProduct(product shop.Product) error {
	panic("implement me")
}

func (CancelDecoratorStub) ModifyProduct(shop.Product) error {
	panic("implement me")
}

func (CancelDecoratorStub) RemoveProduct(name string) error {
	panic("implement me")
}

func (CancelDecoratorStub) GetProduct(name string) (shop.Product, error) {
	panic("implement me")
}

func (CancelDecoratorStub) SetProduct(name string, product shop.Product) error {
	panic("implement me")
}

func (CancelDecoratorStub) Register(userName string) error {
	panic("implement me")
}

func (CancelDecoratorStub) AddBalance(userName string, sum float32) error {
	panic("implement me")
}

func (CancelDecoratorStub) Balance(userName string) (float32, error) {
	panic("implement me")
}

func (CancelDecoratorStub) SetAccount(userName string, account shop.Account) error {
	panic("implement me")
}

func (CancelDecoratorStub) GetAccount(userName string) (shop.Account, error) {
	panic("implement me")
}

func (CancelDecoratorStub) GetAccounts(sortType shop.AccountSortType) []shop.Account {
	panic("implement me")
}

func (CancelDecoratorStub) ModifyAccountType(userName string, accountType shop.AccountType) error {
	panic("implement me")
}

func (CancelDecoratorStub) CalculateOrder(userName string, order shop.Order) (float32, error) {
	panic("implement me")
}

func (CancelDecoratorStub) PlaceOrder(userName string, order shop.Order) error {
	panic("implement me")
}

func (CancelDecoratorStub) AddBundle(name string, main shop.Product, discount float32, additional ...shop.Product) error {
	panic("implement me")
}

func (CancelDecoratorStub) ChangeDiscount(name string, discount float32) error {
	panic("implement me")
}

func (CancelDecoratorStub) RemoveBundle(name string) error {
	panic("implement me")
}

func (CancelDecoratorStub) SetBundle(name string, bundle shop.Bundle) error {
	panic("implement me")
}

func (CancelDecoratorStub) GetBundle(name string) (shop.Bundle, error) {
	panic("implement me")
}

func (CancelDecoratorStub) Import(data []byte) error {
	panic("implement me")
}

func (CancelDecoratorStub) ImportProductsCSV(data []byte) (errs []shop.ImportProductsError) {
	panic("implement me")
}

func (CancelDecoratorStub) ImportAccountsCSV(data []byte) (errs []shop.ImportAccountsError) {
	panic("implement me")
}

func (s CancelDecoratorStub) Export() ([]byte, error) {
	if s.exportFunc != nil {
		return s.exportFunc()
	}
	return []byte{}, nil
}

func (CancelDecoratorStub) ExportProductsCSV() ([]byte, error) {
	panic("implement me")
}

func (CancelDecoratorStub) ExportAccountsCSV() ([]byte, error) {
	panic("implement me")
}
