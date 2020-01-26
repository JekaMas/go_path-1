package shop

import (
	"reflect"
	"testing"
)

/**
1. Пробники, Бандлы доделать
2. Кеширование заказов, Бенчмарки
*/

/* -- ProductManager ------------------------------------------------------------------------------------------------ */

func TestShop_AddProduct(t *testing.T) {

	type test struct {
		name        string
		productName string
		price       float32
		ProductType
		wantErr bool
	}

	m := NewMarket()

	tests := []test{
		{"default", "A", 100, ProductNormal, false},
		{"default", "B", 100, ProductPremium, false},
		{"negative", "C", -10, ProductNormal, true},
		{"empty", "", 0, ProductNormal, true},
		{"empty", "", -10, ProductNormal, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := m.AddProduct(NewProduct(tt.productName, tt.price, tt.ProductType)); (err != nil) != tt.wantErr {
				t.Errorf("AddProduct() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				p := m.Products[tt.productName]
				if tt.price != p.Price || tt.ProductType != p.Type {
					t.Errorf("AddProduct() wrong product values want = %v get = %v",
						Product{tt.productName, tt.price, tt.ProductType}, p)
				}
			}
		})
	}
}

func TestShop_ModifyProduct(t *testing.T) {

	type test struct {
		name    string
		product Product
		wantErr bool
	}

	const (
		productName = "P1"
	)

	m := NewMarket()
	_ = m.AddProduct(Product{Name: "P1", Price: 100, Type: ProductNormal})

	tests := []test{
		{"default", Product{"P1", 120, ProductNormal}, false},
		{"default", Product{"P1", 24.435, ProductPremium}, false},

		{"negative", Product{"P1", -100, ProductNormal}, true},
		{"nil", Product{"P2", 100, ProductNormal}, true},
		{"nil2", Product{}, true},
		{"type", Product{"P1", 100, 42}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := m.ModifyProduct(tt.product); (err != nil) != tt.wantErr {
				t.Errorf("ModifyProduct() error = %v, wantErr %v", err, tt.wantErr)
			}

			if product := m.Products[productName]; !tt.wantErr && reflect.DeepEqual(product, tt.product) {
				t.Errorf("ModifyProduct() product = %v, want %v",
					product, tt.product)
			}
		})
	}
}

func TestShop_RemoveProduct(t *testing.T) {

	type test struct {
		name        string
		productName string
		wantErr     bool
	}

	m := NewMarket()
	_ = m.AddProduct(Product{Name: "P1", Price: 100, Type: ProductNormal})

	tests := []test{
		{"default", "P1", false},
		{"netExist", "P2", true},
		{"emptyName", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := m.RemoveProduct(tt.productName); (err != nil) != tt.wantErr {
				t.Errorf("RemoveProduct() error = %v, wantErr %v, product %v", err, tt.wantErr, tt.name)
			}

			if product, ok := m.Products[tt.productName]; !tt.wantErr && ok {
				t.Errorf("RemoveProduct() product %v has not been removed", product)
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
		username string
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

		balance, err := m.Balance(tt.username)
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
		username string
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

		if err := m.AddBalance(tt.username, tt.sum); (err != nil) != tt.wantErr {
			t.Errorf("AddBalance() error = %v, wantErr %v", err, tt.wantErr)
		}

		if balance, _ := m.Balance(tt.username); !tt.wantErr && balance != tt.want {
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

/* -- OrderManager -------------------------------------------------------------------------------------------------- */

func TestShop_CalculateOrder(t *testing.T) {

	type test struct {
		name    string
		order   Order
		wantSum float32
		wantErr bool
	}

	nt := ProductNormal  // normal type
	pt := ProductPremium // premium type

	acc := Account{"A", 0, AccountNormal}
	premAcc := Account{"B", 0, AccountPremium}

	tests := []test{
		{"one",
			NewOrder([]Product{NewProduct("P", 90, nt)}, nil, acc),
			90, false,
		},
		{"two",
			NewOrder([]Product{NewProduct("P1", 90, nt), NewProduct("P2", 10, nt)}, nil, acc),
			100, false,
		},
		{"premProduct",
			NewOrder([]Product{NewProduct("P1", 90, pt), NewProduct("P2", 10, nt)}, nil, acc),
			95.5, false,
		},
		{"premProduct2",
			NewOrder([]Product{NewProduct("P1", 90, pt), NewProduct("P2", 10, pt)}, nil, acc),
			95, false,
		},
		{"premUserProduct",
			NewOrder([]Product{NewProduct("P1", 90, pt), NewProduct("P2", 10, nt)}, nil, premAcc),
			87, false,
		},
		{"premUserProduct",
			NewOrder([]Product{NewProduct("P1", 90, pt), NewProduct("P2", 10, pt)}, nil, premAcc),
			80, false,
		},
		{"bundle",
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 10, nt), 0, NewProduct("P2", 90, nt))}, acc),
			100, false,
		},
		{"bundle2",
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 10, nt), -10, NewProduct("P2", 90, nt))}, acc),
			90, false,
		},
		{"premBundle",
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 100, pt), -10, NewProduct("P2", 91, nt))}, acc),
			171.9, false,
		},
		{"premBundle2",
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 100, pt), -10, NewProduct("P2", 91, pt))}, acc),
			171.9, false,
		},
		{"premBundle3",
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 100, pt), -10, NewProduct("P2", 91, pt))}, premAcc),
			171.9, false,
		},
		{"zeroBundle",
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 10, nt), -100, NewProduct("P2", 90, nt))}, acc),
			0, false,
		},
		{"errBundle",
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 10, nt), -120, NewProduct("P2", 90, nt))}, acc),
			0, true,
		},
		{"zero",
			NewOrder([]Product{}, []Bundle{}, acc),
			0, false,
		},
		{"zero2",
			NewOrder([]Product{}, nil, acc),
			0, false,
		},
		{"err",
			NewOrder(nil, nil, acc),
			0, true,
		},
		{"err",
			NewOrder(nil, []Bundle{}, acc),
			0, true,
		},
	}

	m := testMarket()

	// test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sum, err := m.CalculateOrder(tt.order)

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
			NewOrder([]Product{NewProduct("P", 90, ProductNormal)}, nil, acc),
			10_000 - 90, false,
		},
		{"default2", acc,
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 10, ProductNormal), -10)}, acc),
			10_000 - 9, false,
		},
		{"default3", acc,
			NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 10, ProductNormal), -10, NewProduct("P2", 90, ProductNormal))}, acc),
			10_000 - 90, false,
		},
		{"zero", acc,
			NewOrder([]Product{NewProduct("P", 0, ProductNormal)}, nil, acc),
			10_000, false,
		},
		{"zero2", acc,
			NewOrder([]Product{}, nil, acc),
			10_000, false,
		},
		{"zero3", acc,
			NewOrder([]Product{}, nil, acc),
			10_000, false,
		},
		{"much", acc,
			NewOrder([]Product{NewProduct("P", 10_001, ProductNormal)}, nil, acc),
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

			if !tt.wantErr && !reflect.DeepEqual(m.GetAccount(tt.acc.Name).Balance, tt.wantBalance) {
				t.Errorf("PlaceOrder() wrong balance, want = %v get = %v",
					tt.wantBalance, m.GetAccount(tt.acc.Name).Balance)
			}
		})
	}

}

/* -- BundleManager ------------------------------------------------------------------------------------------------- */

func TestShop_AddBundle(t *testing.T) {

	type test struct {
		name    string
		bundle  Bundle
		wantErr bool
	}

	nt := ProductNormal  // normal type
	pt := ProductPremium // premium type

	tests := []test{
		{"default", NewBundle(NewProduct("P1", 10, nt), 1, NewProduct("P2", 90, nt)), false},
		{"default2", NewBundle(NewProduct("P1", 10, nt), -10, NewProduct("P2", 90, nt)), false},
		{"default3", NewBundle(NewProduct("P1", 10, nt), -55.2345, NewProduct("P2", 90, nt)), false},
		{"default4", NewBundle(NewProduct("P1", 10, nt), 76.546767, NewProduct("P2", 90, nt)), false},

		{"errDisc", NewBundle(NewProduct("P1", 10, nt), 0, NewProduct("P2", 90, nt)), true},
		{"errDisc2", NewBundle(NewProduct("P1", 10, nt), 100, NewProduct("P2", 90, nt)), true},
		{"errDisc3", NewBundle(NewProduct("P1", 10, nt), 101, NewProduct("P2", 90, nt)), true},
		{"errDisc4", NewBundle(NewProduct("P1", 10, nt), 99.12456, NewProduct("P2", 90, nt)), true},
		{"errDisc5", NewBundle(NewProduct("P1", 10, nt), 99.00001, NewProduct("P2", 90, nt)), true},
		{"errDisc6", NewBundle(NewProduct("P1", 10, nt), 200, NewProduct("P2", 90, nt)), true},

		{"nilProd", NewBundle(Product{}, 0, NewProduct("P2", 90, nt)), true},
		{"nilProd2", NewBundle(NewProduct("P1", 10, nt), 0, Product{}), true},
		{"nilProd3", NewBundle(NewProduct("P1", 10, nt), 0, NewProduct("P2", 90, nt), Product{}), true},
		{"nilProd4", NewBundle(NewProduct("P1", 10, nt), 0, Product{}, NewProduct("P2", 90, nt)), true},
		{"nilProd5", NewBundle(Product{}, 0, NewProduct("P2", 90, nt), Product{}), true},
		// prem
		{"defaultPrem", NewBundle(NewProduct("P1", 10, nt), 1, NewProduct("P2", 90, pt)), false},
		{"errDiscPrem2", NewBundle(NewProduct("P1", 10, pt), 200, NewProduct("P2", 90, nt)), true},
		{"nilProdPrem3", NewBundle(Product{}, 0, NewProduct("P2", 90, pt), Product{}), true},
	}

	// test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMarket()
			err := m.AddBundle(tt.name, tt.bundle.Main, tt.bundle.Discount, tt.bundle.Additional...)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddBundle() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.bundle, *m.Bundles[tt.name]) {
				t.Errorf("AddBundle() wrong bundle added = %v get = %v",
					tt.bundle, m.Bundles[tt.name])
			}
		})
	}

}

func TestShop_ChangeDiscount(t *testing.T) {
	type test struct {
		name     string
		discount float32
		wantErr  bool
	}

	m := NewMarket()
	_ = m.AddBundle("default", NewProduct("P1", 10, ProductNormal), 1, NewProduct("P2", 90, ProductPremium))

	tests := []test{
		{"default", 2, false},
		{"default1", -2, false},
		{"default2", 50, false},
		{"err", 150, true},
		{"err2", 0, true},
		{"err2", 100, true},
		{"err2", 100.0001, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := m.ChangeDiscount("default", tt.discount); (err != nil) != tt.wantErr {
				t.Errorf("ChangeDiscount() error = %v, wantErr %v", err, tt.wantErr)
			}

			if discount := m.Bundles["default"].Discount; !tt.wantErr && discount != tt.discount {
				t.Errorf("ChangeDiscount() discount = %v, want %v",
					discount, tt.discount)
			}
		})
	}

	// nil
	if err := m.ChangeDiscount("aaa", 10); err == nil {
		t.Errorf("ChangeDiscount() modified nil bundle")
	}
}

func TestShop_RemoveBundle(t *testing.T) {

	type test struct {
		name       string
		bundleName string
		wantErr    bool
	}

	m := NewMarket()
	_ = m.AddBundle("B1", Product{"P1", 100, ProductPremium}, 10, Product{"P2", 100, ProductNormal})

	tests := []test{
		{"default", "B1", false},
		{"notExist", "P1", true},
		{"emptyName", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := m.RemoveBundle(tt.bundleName); (err != nil) != tt.wantErr {
				t.Errorf("RemoveBundle() error = %v, wantErr %v, product %v", err, tt.wantErr, tt.name)
			}

			if product, ok := m.Bundles[tt.bundleName]; !tt.wantErr && ok {
				t.Errorf("RemoveBundle() product %v has not been removed", product)
			}
		})
	}

}

/* -- Importer, Exporter -------------------------------------------------------------------------------------------- */

func TestShop_Export(t *testing.T) {

	type test struct {
		name    string
		m       Market
		wantErr bool
	}

	tests := []test{
		{"default", testMarket(), false},
		{"empty", NewMarket(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.m
			bytes, err := m.Export()

			if (err != nil) != tt.wantErr {
				t.Errorf("Export() error while export shop, err= %v", err)
			}

			if bytes == nil || len(bytes) == 0 {
				t.Errorf("Export() empty export")
			}
		})
	}

}

func TestShop_Import(t *testing.T) {

	m := testMarket()
	bytes, _ := m.Export()

	m2 := NewMarket() // empty
	err := m2.Import(bytes)

	if err != nil {
		t.Errorf("Import() error = %v", err)
	}

	if !reflect.DeepEqual(m, m2) {
		t.Errorf("Import() imported shop is invalid")
	}

	err = m2.Import(nil)

	if err == nil {
		t.Errorf("Import() imported nil bytes")
	}

	err = m2.Import([]byte{0})

	if err == nil {
		t.Errorf("Import() imported invalid characters")
	}

}

/* --- new ---------------------------------------------------------------------------------------------------------- */

func testMarket() Market {

	m := NewMarket()

	for _, acc := range testAccounts() {
		_ = m.Register(acc.Name)
		_ = m.AddBalance(acc.Name, acc.Balance)
		_ = m.ModifyAccountType(acc.Name, acc.Type)
	}

	for _, p := range testProducts() {
		_ = m.AddProduct(p)
	}

	for _, b := range testBundles() {
		_ = m.AddBundle(b.Main.Name, b.Main, b.Discount, b.Additional...)
	}

	return m
}

func testAccounts() []Account {
	return []Account{
		{
			Name:    "Sofia",
			Balance: 1_000_000_000,
			Type:    AccountNormal,
		},
		{
			Name:    "Larry",
			Balance: 1_000_000,
			Type:    AccountPremium,
		},
		{
			Name:    "Stan",
			Balance: 100,
			Type:    AccountNormal,
		},
		{
			Name:    "Mary",
			Balance: 1,
			Type:    AccountPremium,
		},
		{
			Name:    "Colin",
			Balance: 0,
			Type:    AccountPremium,
		},
		{
			Name:    "John",
			Balance: 0,
			Type:    AccountNormal,
		},
	}
}

func testProducts() []Product {
	return []Product{
		{
			Name:  "Mouse",
			Price: 599,
			Type:  ProductNormal,
		},
		{
			Name:  "Gaming Mouse",
			Price: 1999,
			Type:  ProductPremium,
		},
		{
			Name:  "Keyboard",
			Price: 399,
			Type:  ProductNormal,
		},
		{
			Name:  "Keyboard",
			Price: 399,
			Type:  ProductNormal,
		},
		{
			Name:  "Plastic Bag",
			Price: 0,
			Type:  ProductNormal,
		},
		{
			Name:  "Product A",
			Price: 1_999_999_999,
			Type:  ProductNormal,
		},
		{
			Name:  "Product B",
			Price: 1_999_999_999,
			Type:  ProductPremium,
		},
	}
}

func testBundles() []Bundle {
	return []Bundle{
		{
			Type:     BundleNormal,
			Discount: 20,
			Main:     Product{Name: "Bowling Ball", Price: 100, Type: ProductNormal},
			Additional: []Product{
				{Name: "Skittle", Price: 50, Type: ProductNormal},
				{Name: "Skittle", Price: 50, Type: ProductNormal},
				{Name: "Skittle", Price: 50, Type: ProductPremium},
			},
		},
		{
			Type:     BundleSample,
			Discount: 1,
			Main:     Product{Name: "Toy Gun", Price: 0, Type: ProductPremium},
			Additional: []Product{
				{Name: "Battery", Price: 30, Type: ProductNormal},
				{Name: "Battery", Price: 30, Type: ProductNormal},
			},
		},
	}
}
