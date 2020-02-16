package shop

import (
	"github.com/pkg/errors"
	"sort"
	"strings"
	"unicode"
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

	if _, err := m.GetAccount(userName); err == nil { // no error with get
		return ErrorAccountExists
	}

	acc := NewAccount(userName)
	return m.SetAccount(userName, acc)
}

func (m *Market) AddBalance(userName string, sum float32) error {

	if sum < 0 {
		return ErrorAccountAddNegativeBalance
	}

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

	acc, err := m.GetAccount(userName)
	if err != nil {
		return errors.Wrap(err, "can't modify nil account")
	}

	acc.Type = accountType
	return m.SetAccount(userName, acc)
}

func (m *Market) Balance(userName string) (float32, error) {

	acc, err := m.GetAccount(userName)
	if err != nil {
		return 0, errors.Wrap(err, "can't get balance of the nil account")
	}

	return acc.Balance, nil
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

func checkName(name string) error {
	if len(name) == 0 {
		return ErrorEmptyField
	}

	// TODO max chars count
	//if len(userName) > MAX_NAME_LENGTH {
	//	return ErrorAccountInvalidName
	//}

	for _, r := range name { // for each rune
		if !(unicode.IsLetter(r) || unicode.IsDigit(r)) {
			return ErrorAccountInvalidName
		}
	}

	return nil
}

func checkAccount(acc Account) error {
	if _, ok := AccountTypeMap[acc.Type]; !ok {
		return ErrorAccountInvalidType
	}

	if len(acc.Name) == 0 {
		return ErrorAccountInvalidName
	}

	return nil
}
