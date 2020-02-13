package shop

import (
	"errors"
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
		{testName: "Error Type", Account: Account{Name: "Bred", Balance: 100, Type: 5}, err: errors.New("incorrect type")},
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
