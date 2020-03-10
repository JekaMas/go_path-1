package time

import (
	shop "go_path/struct"
	"time"
)

func (td *TimeoutDecorator) Register(userName string) error {
	return td.timeoutFunc(func(ch chan error) {
		ch <- td.shop.Register(userName)
	}, time.Second)
}

func (td *TimeoutDecorator) AddBalance(userName string, sum float32) error {
	return td.timeoutFunc(func(ch chan error) {
		ch <- td.shop.AddBalance(userName, sum)
	}, time.Second)
}

func (td *TimeoutDecorator) ModifyAccountType(userName string, accountType shop.AccountType) error {
	return td.timeoutFunc(func(ch chan error) {
		ch <- td.shop.ModifyAccountType(userName, accountType)
	}, time.Second)
}

func (td *TimeoutDecorator) Balance(userName string) (float32, error) {
	return td.timeoutFuncAmount(func(ch chan amountResult) {
		sum, err := td.shop.Balance(userName)
		ch <- amountResult{sum, err}
	}, time.Second)
}

func (td *TimeoutDecorator) GetAccounts(sortType shop.AccountSortType) []shop.Account {
	return td.timeoutFuncAccounts(func(ch chan []shop.Account) {
		ch <- td.shop.GetAccounts(sortType)
	}, time.Second)
}

/* --- Interface ---------------------------------------------------------------------------------------------------- */

func (td *TimeoutDecorator) GetAccount(name string) (shop.Account, error) {
	return td.timeoutFuncAccount(func(ch chan accountResult) {
		acc, err := td.shop.GetAccount(name)
		ch <- accountResult{acc, err}
	}, time.Second)
}

func (td *TimeoutDecorator) SetAccount(userName string, account shop.Account) error {
	return td.timeoutFunc(func(ch chan error) {
		ch <- td.shop.SetAccount(userName, account)
	}, time.Second)
}
