package shop

import (
	"github.com/pkg/errors"
	"sort"
	"strings"
)

var (
	ErrorAccountNotRegistered = errors.New("account is not registered")
	ErrorAccountExists        = errors.New("account already exists")
	ErrorAccountInvalidType   = errors.New("account type is invalid")
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

	if len(userName) == 0 {
		return ErrorEmptyField
	}

	if _, ok := m.Accounts[userName]; ok {
		return ErrorAccountExists
	}

	acc := NewAccount(userName)
	m.Accounts[userName] = acc
	return nil
}

func (m *Market) AddBalance(userName string, sum float32) error {

	if sum < 0 {
		return errors.New("cannot add negative sum")
	}

	if _, ok := m.Accounts[userName]; !ok {
		return ErrorAccountNotRegistered
	}

	acc := m.Accounts[userName]
	acc.Balance += sum

	m.Accounts[userName] = acc
	return nil
}

func (m *Market) ModifyAccountType(userName string, accountType AccountType) error {

	if _, ok := m.Accounts[userName]; !ok {
		return ErrorAccountNotRegistered
	}
	if _, ok := AccountTypeMap[accountType]; !ok {
		return ErrorAccountInvalidType
	}

	acc := m.Accounts[userName]
	acc.Type = accountType

	m.Accounts[userName] = acc
	return nil
}

func (m *Market) Balance(userName string) (float32, error) {

	if _, ok := m.Accounts[userName]; !ok {
		return 0, ErrorAccountNotRegistered
	}

	return m.Accounts[userName].Balance, nil
}

func (m *Market) GetAccount(name string) (Account, error) {
	acc, ok := m.Accounts[name]

	if !ok {
		return Account{}, ErrorAccountNotRegistered
	}

	return acc, nil
}

func (m *Market) GetAccounts(sortType AccountSortType) []Account {
	var accs []Account
	for _, acc := range m.Accounts {
		accs = append(accs, acc)
	}
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
