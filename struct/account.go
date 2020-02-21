package shop

import (
	"github.com/pkg/errors"
	"sort"
	"strings"
	"time"
)

var (
	ErrorAccountAddNegativeBalance = errors.New("can't add negative balance")
	ErrorAccountNotRegistered      = errors.New("account is not registered")
	ErrorAccountInvalidType        = errors.New("account type is invalid")
	ErrorAccountInvalidName        = errors.New("account is invalid")
	ErrorAccountExists             = errors.New("account already exists")
)

/* -- AccountManager ------------------------------------------------------------------------------------------------ */

func NewAccount(userName string) Account {
	return Account{
		Name:    userName,
		Balance: 0,
		Type:    AccountNormal,
	}
}

func (m *Market) Register(userName string) error {

	resChan := make(chan error)

	go func() {
		if _, err := m.GetAccount(userName); err == nil { // no error with get
			resChan <- ErrorAccountExists
		}
		acc := NewAccount(userName)

		m.accountsMutex.Lock()
		defer m.accountsMutex.Unlock()

		resChan <- m.SetAccount(userName, acc)
	}()

	select {
	case res := <-resChan:
		return res
	case <-time.After(time.Second):
		return ErrorTimeout
	}
}

func (m *Market) AddBalance(userName string, sum float32) error {

	if sum < 0 {
		return ErrorAccountAddNegativeBalance
	}

	m.accountsMutex.Lock()
	defer m.accountsMutex.Unlock()

	acc, err := m.GetAccount(userName)
	if err != nil {
		return errors.Wrap(err, "can't add balance to the nil account")
	}

	acc.Balance += sum
	return m.SetAccount(userName, acc)
}

func (m *Market) ModifyAccountType(userName string, accountType AccountType) error {

	if _, ok := AccountTypeMap[accountType]; !ok { // check type itself
		return ErrorAccountInvalidType
	}

	m.accountsMutex.Lock()
	defer m.accountsMutex.Unlock()

	acc, err := m.GetAccount(userName)
	if err != nil {
		return errors.Wrap(err, "can't modify nil account")
	}

	acc.Type = accountType
	return m.SetAccount(userName, acc)
}

func (m *Market) Balance(userName string) (float32, error) {

	m.accountsMutex.RLock()
	acc, err := m.GetAccount(userName)
	m.accountsMutex.RUnlock()

	if err != nil {
		return 0, errors.Wrap(err, "can't get balance of the nil account")
	}

	return acc.Balance, nil
}

func (m *Market) GetAccounts(sortType AccountSortType) []Account {
	var accs []Account

	m.accountsMutex.RLock()
	for _, acc := range m.Accounts {
		accs = append(accs, acc)
	}
	m.accountsMutex.RUnlock()

	// compare function
	var less func(i, j int) bool

	switch sortType {
	default:
		fallthrough
	case SortByName:
		less = func(i, j int) bool {
			return strings.Compare(accs[i].Name, accs[j].Name) < 0
		}
	case SortByNameReverse:
		less = func(i, j int) bool {
			return strings.Compare(accs[i].Name, accs[j].Name) > 0
		}
	case SortByBalance:
		less = func(i, j int) bool {
			return accs[i].Balance < accs[j].Balance
		}
	}

	sort.Slice(accs, less)
	return accs
}

/* --- Interface ---------------------------------------------------------------------------------------------------- */

func (m *Market) GetAccount(name string) (Account, error) {
	acc, ok := m.Accounts[name]

	if !ok {
		return Account{}, ErrorAccountNotRegistered
	}

	return acc, nil
}

func (m *Market) SetAccount(userName string, account Account) error {

	if err := checkName(userName); err != nil {
		return errors.Wrap(err, "can't set invalid account")
	}
	if err := checkAccount(account); err != nil {
		return errors.Wrap(err, "can't set invalid account")
	}

	m.Accounts[userName] = account
	return nil
}

/* --- Checks ------------------------------------------------------------------------------------------------------- */

func checkAccount(acc Account) error {
	if _, ok := AccountTypeMap[acc.Type]; !ok {
		return ErrorAccountInvalidType
	}

	if len(acc.Name) == 0 {
		return ErrorAccountInvalidName
	}

	return nil
}
