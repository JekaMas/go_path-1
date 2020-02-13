package shop

import (
	"errors"
	"testing"
)

var (
	ErrorPriceNotNegative     = errors.New("total price not be negative/zero")
	ErrorIncorrectAccountType = errors.New("incorrect account type")
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
			total: 125,
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
			total: 125,
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
				t.Errorf("Error, type not correct: %v != %v", total, test.total)
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
			err: ErrorNegativeProductPrice,
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
			err: ErrorNegativeProductPrice,
		},
		{testName: "BuyWithoutAddForPremium",
			username: "Alfred",
			order: Order{
				Products: []Product{
					{Name: "Pineapple", Price: 100, Type: ProductPremium},
				},
				Bundles: nil,
			},
			err: ErrorIncorrectAccountType,
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
