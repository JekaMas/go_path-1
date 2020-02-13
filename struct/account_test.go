package shop

import (
	"errors"
	"reflect"
	"testing"
)

// AddBalanceMantiseError
func TestAddBalanceMantis(t *testing.T) {
	testShop := NewMarket()
	testShop.Accounts["Bred"] = Account{
		Name:    "Bred",
		Balance: 0,
		Type:    AccountNormal,
	}
	for i := 0; i < 1000; i++ {
		err := testShop.AddBalance("Bred", 0.1)
		if err != nil {
			t.Errorf("Test it should be success, error: %v", err)
		}
	}

	if testShop.Accounts["Bred"].Balance != 100 {
		t.Errorf("Error, balance not correct: %v != %v", testShop.Accounts["Bred"].Balance, 100)
	}
}

// ModifyAccountType
func TestModifyAccountTypeFailed(t *testing.T) {
	type accountTest struct {
		testName string
		Account  Account
		err      error
	}

	tests := []accountTest{
		{testName: "Error Type", Account: Account{Name: "Bred", Balance: 100, Type: 100}, err: ErrorAccountInvalidType},
	}
	testShop := NewMarket()
	testShop.Accounts["Bred"] = Account{
		Name:    "Bred",
		Balance: 0,
		Type:    AccountNormal,
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			err := testShop.ModifyAccountType(test.Account.Name, test.Account.Type)
			if err == nil {
				t.Errorf("Test it should be failed with error: %v, but get value type = %v", test.err, testShop.Accounts[test.Account.Name].Type)
			} else if err.Error() != test.err.Error() && !errors.Is(err, test.err) {
				t.Errorf("Values not equal: %v != %v", err, test.err)
			}
		})
	}
}

/* -- AccountManager ------------------------------------------------------------------------------------------------ */

func TestShop_Register(t *testing.T) {

	m := NewMarket()

	err := m.Register("Spike")
	if _, ok := m.Accounts["Spike"]; !ok {
		t.Fatalf("Register() error = %v", err)
	}

	if err := m.Register("Spike"); err == nil {
		t.Fatal("Register() registered twice")
	}

	_ = m.Register("")
	if _, ok := m.Accounts[""]; ok {
		t.Fatal("Register() registered with empty name")
	}
}

func TestShop_Balance(t *testing.T) {

	type test struct {
		name     string
		userName string
		want     float32
		wantErr  bool
	}

	m := NewMarket()
	_ = m.Register("John")
	_ = m.Register("Stan")
	_ = m.AddBalance("Stan", 100)

	tests := []test{
		{"default", "John", 0, false},
		{"default", "Stan", 100, false},
		{"not exists", "AAA", 0, true},
		{"empty", "", 0, true},
	}

	for _, tt := range tests {

		balance, err := m.Balance(tt.userName)
		if (err != nil) != tt.wantErr {
			t.Errorf("Balance() error = %v, wantErr %v", err, tt.wantErr)
		}

		if !tt.wantErr && balance != tt.want {
			t.Errorf("Balance() balance = %v, want %v",
				balance, tt.want)
		}
	}
}

func TestShop_AddBalance(t *testing.T) {

	type test struct {
		name     string
		userName string
		sum      float32
		want     float32
		wantErr  bool
	}

	m := NewMarket()
	_ = m.Register("John")
	_ = m.Register("Stan")
	_ = m.AddBalance("Stan", 100)

	tests := []test{
		{"default", "John", 100, 100, false},
		{"default", "Stan", 100, 200, false},
		{"negative", "Stan", -10, 200, true},
		{"not exists", "AAA", 0, 0, true},
		{"empty", "", 0, 0, true},
	}

	for _, tt := range tests {

		if err := m.AddBalance(tt.userName, tt.sum); (err != nil) != tt.wantErr {
			t.Errorf("AddBalance() error = %v, wantErr %v", err, tt.wantErr)
		}

		if balance, _ := m.Balance(tt.userName); !tt.wantErr && balance != tt.want {
			t.Errorf("AddBalance() balance = %v, want %v",
				balance, tt.want)
		}
	}
}

func TestShop_GetAccounts(t *testing.T) {

	type test struct {
		name     string
		sortType AccountSortType
		result   []Account
	}

	m := NewMarket()

	// init
	_ = m.Register("John")
	_ = m.Register("Tom")
	_ = m.Register("Stan")

	_ = m.AddBalance("John", 3)
	_ = m.AddBalance("Tom", 1)
	_ = m.AddBalance("Stan", 2)

	var nt = AccountNormal // normal type
	names := []Account{{"John", 3, nt}, {"Stan", 2, nt}, {"Tom", 1, nt}}
	rever := []Account{{"Tom", 1, nt}, {"Stan", 2, nt}, {"John", 3, nt}}
	balan := []Account{{"Tom", 1, nt}, {"Stan", 2, nt}, {"John", 3, nt}}

	tests := []test{
		{"userName", SortByName, names},
		{"userReverse", SortByNameReverse, rever},
		{"balance", SortByBalance, balan},
		{"err", 100, names},
	}

	// test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := m.GetAccounts(tt.sortType)
			// spew.Dump(a)
			if !reflect.DeepEqual(a, tt.result) {
				t.Errorf("GetAccounts() sort not working properly, want = %v get = %v",
					tt.result, a)
			}
		})
	}

	// empty test
	m = NewMarket()
	var empty []Account

	for _, tt := range tests {
		t.Run("empty", func(t *testing.T) {
			a := m.GetAccounts(tt.sortType)
			if !reflect.DeepEqual(a, empty) {
				t.Errorf("GetAccounts() sort not working properly, want = %v get = %v",
					[]Account{}, a)
			}
		})
	}
}
