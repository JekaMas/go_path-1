package shop

import (
	"sync"

	"github.com/pkg/errors"
)

type Accounts struct {
	Accounts map[string]Account
	mu       sync.RWMutex
}

func NewAccounts() Accounts {
	return Accounts{
		Accounts: make(map[string]Account),
		mu:       sync.RWMutex{},
	}
}

func (a *Accounts) Register(userName string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, err := a.getAccount(userName); err == nil { // no error with get
		return ErrorAccountExists
	}

	return a.setAccount(userName, NewAccount(userName))
}

func (a *Accounts) GetAccount(name string) (Account, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.getAccount(name)
}

func (a *Accounts) getAccount(name string) (Account, error) {
	acc, ok := a.Accounts[name]
	if !ok {
		return Account{}, ErrorAccountNotRegistered
	}

	return acc, nil
}

func (a *Accounts) SetAccount(userName string, account Account) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.setAccount(userName, account)
}

func (a *Accounts) setAccount(userName string, account Account) error {
	if err := checkName(userName); err != nil {
		return errors.Wrap(err, "can't set invalid account")
	}
	if err := checkAccount(account); err != nil {
		return errors.Wrap(err, "can't set invalid account")
	}

	a.Accounts[userName] = account
	return nil
}

func (a *Accounts) changeAccount(userName string, fns ...changeAccountFunc) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	acc, err := a.getAccount(userName)
	if err != nil {
		return errors.Wrap(err, "can't add balance to the nil account")
	}

	for _, fn := range fns {
		fn(&acc)
	}

	a.Accounts[userName] = acc
	return nil
}

type changeAccountFunc func(account *Account)

func changeType(t AccountType) changeAccountFunc {
	return func(account *Account) {
		account.Type = t
	}
}

func addBalance(b float32) changeAccountFunc {
	return func(account *Account) {
		account.Balance += b
	}
}

func (a *Accounts) clone() []Account {
	a.mu.RLock()
	defer a.mu.RUnlock()

	n := len(a.Accounts)
	if n == 0 {
		return nil
	}

	accs := make([]Account, n)
	i := 0
	for _, acc := range a.Accounts {
		accs[i] = acc
		i++
	}

	return accs
}
