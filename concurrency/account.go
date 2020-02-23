package concurrency

import (
	"errors"
	shop "go_path/struct"
	"time"
)

var (
	ErrorTimeout = errors.New("timeout")
)

func (m *Market) Register(userName string) error {

	resChan := make(chan error, 1)

	go func() {
		resChan <- m.shop.Register(userName)
	}()

	select {
	case res := <-resChan:
		return res
	case <-time.After(time.Second):
		return ErrorTimeout
	}
}

func (m *Market) AddBalance(userName string, sum float32) error {

	resChan := make(chan error, 1)

	go func() {
		resChan <- m.shop.AddBalance(userName, sum)
	}()

	select {
	case res := <-resChan:
		return res
	case <-time.After(time.Second):
		return ErrorTimeout
	}
}

func (m *Market) ModifyAccountType(userName string, accountType shop.AccountType) error {

	resChan := make(chan error, 1)

	go func() {
		resChan <- m.shop.ModifyAccountType(userName, accountType)
	}()

	select {
	case res := <-resChan:
		return res
	case <-time.After(time.Second):
		return ErrorTimeout
	}
}

func (m *Market) Balance(userName string) (float32, error) {

	type result struct {
		sum float32
		err error
	}

	resChan := make(chan result, 1)

	go func() {
		sum, err := m.shop.Balance(userName)
		resChan <- result{sum, err}
	}()

	select {
	case res := <-resChan:
		return res.sum, res.err
	case <-time.After(time.Second):
		return 0, ErrorTimeout
	}
}

func (m *Market) GetAccounts(sortType shop.AccountSortType) []shop.Account {

	resChan := make(chan []shop.Account, 1)

	go func() {
		resChan <- m.shop.GetAccounts(sortType)
	}()

	select {
	case res := <-resChan:
		return res
	case <-time.After(time.Second):
		return nil
	}
}

/* --- Interface ---------------------------------------------------------------------------------------------------- */

func (m *Market) GetAccount(name string) (shop.Account, error) {

	type result struct {
		acc shop.Account
		err error
	}

	resChan := make(chan result, 1)

	go func() {
		sum, err := m.shop.GetAccount(name)
		resChan <- result{sum, err}
	}()

	select {
	case res := <-resChan:
		return res.acc, res.err
	case <-time.After(time.Second):
		return shop.Account{}, ErrorTimeout
	}
}

func (m *Market) SetAccount(userName string, account shop.Account) error {

	resChan := make(chan error, 1)

	go func() {
		resChan <- m.shop.SetAccount(userName, account)
	}()

	select {
	case res := <-resChan:
		return res
	case <-time.After(time.Second):
		return ErrorTimeout
	}
}
