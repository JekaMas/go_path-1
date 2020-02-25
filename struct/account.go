package shop

import (
	"github.com/pkg/errors"
	"sort"
	"strings"
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

func (m *Market) AddBalance(userName string, sum float32) error {
	if sum < 0 {
		return ErrorAccountAddNegativeBalance
	}

	return m.changeAccount(userName, addBalance(sum))
}

func (m *Market) ModifyAccountType(userName string, accountType AccountType) error {
	if _, ok := AccountTypeMap[accountType]; !ok {
		// check type itself
		return ErrorAccountInvalidType
	}

	return m.changeAccount(userName, changeType(accountType))
}

func (m *Market) Balance(userName string) (float32, error) {
	acc, err := m.GetAccount(userName)
	if err != nil {
		return 0, errors.Wrap(err, "can't get balance of the nil account")
	}

	return acc.Balance, nil
}

func (m *Market) GetAccounts(sortType AccountSortType) []Account {
	accs := m.clone()

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
