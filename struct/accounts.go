package shop

import (
	"sync"

	"github.com/pkg/errors"
)

type Accounts struct {
	Accounts map[string]Account
	mu       sync.RWMutex
}

func (m *Accounts) Register(userName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, err := m.getAccount(userName); err == nil { // no error with get
		return ErrorAccountExists
	}

	return m.setAccount(userName, NewAccount(userName))
}

func (m *Accounts) GetAccount(name string) (Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.getAccount(name)
}

func (m *Accounts) getAccount(name string) (Account, error) {
	acc, ok := m.Accounts[name]
	if !ok {
		return Account{}, ErrorAccountNotRegistered
	}

	return acc, nil
}

func (m *Accounts) SetAccount(userName string, account Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.setAccount(userName, account)
}

func (m *Accounts) setAccount(userName string, account Account) error {
	if err := checkName(userName); err != nil {
		return errors.Wrap(err, "can't set invalid account")
	}
	if err := checkAccount(account); err != nil {
		return errors.Wrap(err, "can't set invalid account")
	}

	m.Accounts[userName] = account
	return nil
}

func (m *Accounts) changeAccount(userName string, fns ...changeFunc) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	acc, err := m.getAccount(userName)
	if err != nil {
		return errors.Wrap(err, "can't add balance to the nil account")
	}

	for _, fn := range fns {
		fn(&acc)
	}

	m.Accounts[userName] = acc
	return nil
}

type changeFunc func(account *Account)

func changeType(t AccountType) changeFunc {
	return func(account *Account) {
		account.Type = t
	}
}

func addBalance(b float32) changeFunc {
	return func(account *Account) {
		account.Balance += b
	}
}

func (m *Accounts) clone() []Account {
	m.mu.RLock()
	defer m.mu.RUnlock()

	n := len(m.Accounts)
	if n == 0 {
		return nil
	}

	accs := make([]Account, n)
	i := 0
	for _, acc := range m.Accounts {
		accs[i] = acc
		i++
	}

	return accs
}
