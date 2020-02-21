package shop

import (
	"errors"
	"reflect"
	"testing"
)

func TestCalculateOrderSuccess(t *testing.T) {
	type calcTest struct {
		testName string
		order    Order
		total    float32
	}

	tests := []calcTest{
		{
			testName: "ProductsShouldBeInitialize",
			order: Order{
				Products: nil,
				Bundles: []Bundle{
					{
						Products: []Product{
							{
								Name:  "Main",
								Price: 100,
								Type:  ProductNormal,
							},
							{
								Name:  "Additional",
								Price: 100,
								Type:  ProductPremium,
							},
						},
						Type:     BundleNormal,
						Discount: 50,
					},
				},
			},
			total: 100,
		},
		{
			testName: "ProductsShouldBeInitialize",
			order: Order{
				Products: []Product{},
				Bundles: []Bundle{
					{
						Products: []Product{
							{
								Name:  "Main",
								Price: 100,
								Type:  ProductNormal,
							},
							{
								Name:  "Additional",
								Price: 100,
								Type:  ProductPremium,
							},
						},
						Type:     BundleNormal,
						Discount: 50,
					},
				},
			},
			total: 100,
		},
	}
	testShop := NewMarket()

	testShop.Accounts["Bred"] = Account{
		Name:    "Bred",
		Balance: 100_000,
		Type:    AccountNormal,
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			total, err := testShop.CalculateOrder("Bred", test.order)
			if err != nil {
				t.Fatalf("Test it should be success, error: %v", err)
			}
			if total != test.total {
				t.Errorf("Error, wrong sum: %v != %v", total, test.total)
			}
		})
	}
}

func TestPlaceOrderFailed(t *testing.T) {
	type productTest struct {
		testName string
		username string
		order    Order
		err      error
	}

	tests := []productTest{
		{testName: "PriceNotBeZero",
			username: "Bred",
			order: Order{
				Products: []Product{
					{Name: "Pineapple", Price: 0, Type: ProductNormal},
					{Name: "Pineapple", Price: 0, Type: ProductNormal},
					{Name: "Pineapple", Price: 0, Type: ProductNormal},
				},
				Bundles: nil,
			},
			err: ErrorProductNegativePrice,
		},
		{testName: "PriceNotBeNegative",
			username: "Bred",
			order: Order{
				Products: []Product{
					{Name: "Pineapple", Price: -100, Type: ProductNormal},
					{Name: "Pineapple", Price: -100, Type: ProductNormal},
					{Name: "Pineapple", Price: -100, Type: ProductNormal},
				},
				Bundles: nil,
			},
			err: ErrorProductNegativePrice,
		},
		{testName: "BuyWithoutAddForPremium",
			username: "Alfred",
			order: Order{
				Products: []Product{
					{Name: "Pineapple", Price: 100, Type: ProductPremium},
				},
				Bundles: nil,
			},
			err: ErrorAccountInvalidType,
		},
	}

	testShop := NewMarket()

	testShop.Accounts["Bred"] = Account{
		Name:    "Bred",
		Balance: 100_000,
		Type:    AccountNormal,
	}
	testShop.Accounts["Alfred"] = Account{
		Name:    "Alfred",
		Balance: 1000,
		Type:    100,
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			err := testShop.PlaceOrder(test.username, test.order)
			if err == nil {
				t.Errorf("Test it should be failed, error: %v, balance = %v", test.err, testShop.Accounts[test.username].Balance)
			} else if err.Error() != test.err.Error() && !errors.Is(err, test.err) {
				t.Errorf("Values not equal: %v != %v", err, test.err)
			}
		})
	}
}

// mantis error
// PlaceOrder
func TestCalculateOrderMantis(t *testing.T) {
	type calcTest struct {
		testName string
		order    Order
		total    float32
	}

	test := calcTest{
		testName: "MantisError",
		order: Order{
			Products: []Product{
				{Name: "Pineapple", Price: 0.1, Type: ProductNormal},
				{Name: "Pineapple", Price: 0.1, Type: ProductNormal},
				{Name: "Pineapple", Price: 0.1, Type: ProductNormal},
				{Name: "Pineapple", Price: 0.1, Type: ProductNormal},
				{Name: "Pineapple", Price: 0.1, Type: ProductNormal},

				{Name: "Pineapple", Price: 0.1, Type: ProductNormal},
				{Name: "Pineapple", Price: 0.1, Type: ProductNormal},
				{Name: "Pineapple", Price: 0.1, Type: ProductNormal},
				{Name: "Pineapple", Price: 0.1, Type: ProductNormal},
				{Name: "Pineapple", Price: 0.1, Type: ProductNormal},
			},
			Bundles: nil,
		},
		total: 1,
	}
	testShop := NewMarket()

	testShop.Accounts["Bred"] = Account{
		Name:    "Bred",
		Balance: 100_000,
		Type:    AccountNormal,
	}

	total, err := testShop.CalculateOrder("Bred", test.order)
	if err != nil {
		t.Errorf("Test it should be success, error: %v", err)
	}

	if total != test.total {
		t.Errorf("Error, price not correct: %v != %v", total, test.total)
	}
}

/* -- OrderManager -------------------------------------------------------------------------------------------------- */

func TestShop_CalculateOrder(t *testing.T) {

	type test struct {
		name     string
		userName string
		order    Order
		wantSum  float32
		wantErr  bool
	}

	nt := ProductNormal  // normal type
	pt := ProductPremium // premium type
	st := ProductSampled // sampled type

	acc := Account{"A", 0, AccountNormal}
	premAcc := Account{"B", 0, AccountPremium}

	m := testMarket()
	// reg
	for _, a := range []Account{acc, premAcc} {
		_ = m.Register(a.Name)
		_ = m.ModifyAccountType(a.Name, a.Type)
		_ = m.AddBalance(a.Name, a.Balance)
	}

	// tests
	tests := []test{
		{"one", acc.Name,
			NewOrder([]Product{NewProduct("P", 90, nt)}, nil),
			90, false,
		},
		{"two", acc.Name,
			NewOrder([]Product{NewProduct("P1", 90, nt), NewProduct("P2", 10, nt)}, nil),
			100, false,
		},
		{"premProduct", acc.Name,
			NewOrder([]Product{NewProduct("P1", 90, pt), NewProduct("P2", 10, nt)}, nil),
			95.5, false,
		},
		{"premProduct2", acc.Name,
			NewOrder([]Product{NewProduct("P1", 90, pt), NewProduct("P2", 10, pt)}, nil),
			95, false,
		},
		{"premUserProduct", premAcc.Name,
			NewOrder([]Product{NewProduct("P1", 90, pt), NewProduct("P2", 10, nt)}, nil),
			87, false,
		},
		{"premUserProduct1", premAcc.Name,
			NewOrder([]Product{NewProduct("P1", 90, pt), NewProduct("P2", 10, pt)}, nil),
			80, false,
		},
		{"bundle", acc.Name,
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 10, nt), 1, BundleNormal, NewProduct("P2", 90, nt))}),
			99, false,
		},
		{"bundle2", acc.Name,
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 10, nt), 10, BundleNormal, NewProduct("P2", 90, nt))}),
			90, false,
		},
		{"premBundle", acc.Name,
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 100, pt), 10, BundleNormal, NewProduct("P2", 91, nt))}),
			171.9, false,
		},
		{"premBundle2", acc.Name,
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 100, pt), 10, BundleNormal, NewProduct("P2", 91, pt))}),
			171.9, false,
		},
		{"premBundle3", premAcc.Name,
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 100, pt), 10, BundleNormal, NewProduct("P2", 91, pt))}),
			171.9, false,
		},
		//{"nineBundle", // FIXME precision
		//	NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 10, nt), -99, NewProduct("P2", 90, nt))}, acc),
		//	1, false,
		//},
		{"errBundle", acc.Name,
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 10, nt), 100, BundleNormal, NewProduct("P2", 90, nt))}),
			0, true,
		},
		{"errBundle", acc.Name,
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 10, nt), 120, BundleNormal, NewProduct("P2", 90, nt))}),
			0, true,
		},
		{"zero", acc.Name,
			NewOrder([]Product{}, []Bundle{}),
			0, false,
		},
		{"zero2", acc.Name,
			NewOrder([]Product{}, nil),
			0, false,
		},
		{"err", acc.Name,
			NewOrder(nil, nil),
			0, true,
		},
		{"nilProducts", acc.Name,
			NewOrder(nil, []Bundle{}),
			0, false,
		},
		{"errSampled", acc.Name,
			NewOrder([]Product{NewProduct("P1", 90, st), NewProduct("P2", 10, nt)}, nil),
			100, true,
		},
		{"sampled", acc.Name,
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 10, nt), 10, BundleSample, NewProduct("P2", 90, st))}),
			90, false,
		},
		//{"sampledErr", // FIXME case
		//	NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 10, nt), 10, NewProduct("P2", 90, st), NewProduct("P2", 90, st))}, acc),
		//	0, true,
		//},
	}

	// test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sum, err := m.CalculateOrder(tt.userName, tt.order)

			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.wantSum, sum) {
				t.Errorf("CalculateOrder() wrong sum, want = %v get = %v",
					tt.wantSum, sum)
			}
		})
	}
}

func TestShop_PlaceOrder(t *testing.T) {

	type test struct {
		name        string
		acc         Account
		order       Order
		wantBalance float32
		wantErr     bool
	}

	acc := Account{"A", 10_000, AccountNormal}

	tests := []test{
		{"default", acc,
			NewOrder([]Product{NewProduct("P", 90, ProductNormal)}, nil),
			10_000 - 90, false,
		},
		{"default2", acc,
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 10, ProductNormal), 10, BundleNormal)}),
			10_000 - 9, false,
		},
		{"default3", acc,
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 10, ProductNormal), 10, BundleNormal, NewProduct("P2", 90, ProductNormal))}),
			10_000 - 90, false,
		},
		{"zero", acc,
			NewOrder([]Product{NewProduct("P", 0, ProductNormal)}, nil),
			10_000, false,
		},
		{"zero2", acc,
			NewOrder([]Product{}, nil),
			10_000, false,
		},
		{"zero3", acc,
			NewOrder([]Product{}, nil),
			10_000, false,
		},
		{"much", acc,
			NewOrder([]Product{NewProduct("P", 10_001, ProductNormal)}, nil),
			0, true,
		},
	}

	// test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := testMarket()
			_ = m.Register(tt.acc.Name)
			_ = m.AddBalance(tt.acc.Name, tt.acc.Balance)

			err := m.PlaceOrder(tt.acc.Name, tt.order)

			if (err != nil) != tt.wantErr {
				t.Errorf("PlaceOrder() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				acc, _ := m.GetAccount(tt.acc.Name)
				if !tt.wantErr && !reflect.DeepEqual(acc.Balance, tt.wantBalance) {
					t.Errorf("PlaceOrder() wrong balance, want = %v get = %v",
						tt.wantBalance, acc.Balance)
				}
			}
		})
	}

}
