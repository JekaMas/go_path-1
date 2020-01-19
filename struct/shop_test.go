package shop

import (
	"reflect"
	"testing"
)

/**
1. Доделать тесты
2. Пробники
3. Сохранение/чтение из файла
4. Кеширование заказов
5. Бандлы доделать
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

// TODO
func TestShop_ModifyProduct(t *testing.T) {
	type fields struct {
		Accounts map[string]*Account
		Products map[string]*Product
		Bundles  map[string]*Bundle
	}
	type args struct {
		p Product
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Market{
				Accounts: tt.fields.Accounts,
				Products: tt.fields.Products,
				Bundles:  tt.fields.Bundles,
			}
			if err := m.ModifyProduct(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("ModifyProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TODO
func TestShop_RemoveProduct(t *testing.T) {
	type fields struct {
		Accounts map[string]*Account
		Products map[string]*Product
		Bundles  map[string]*Bundle
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Market{
				Accounts: tt.fields.Accounts,
				Products: tt.fields.Products,
				Bundles:  tt.fields.Bundles,
			}
			if err := m.RemoveProduct(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("RemoveProduct() error = %v, wantErr %v", err, tt.wantErr)
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

// TODO
func TestShop_GetAccounts(t *testing.T) {
	type fields struct {
		Accounts map[string]*Account
		Products map[string]*Product
		Bundles  map[string]*Bundle
	}
	type args struct {
		sortType AccountSortType
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Account
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Market{
				Accounts: tt.fields.Accounts,
				Products: tt.fields.Products,
				Bundles:  tt.fields.Bundles,
			}
			if got := m.GetAccounts(tt.args.sortType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAccounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

/* -- OrderManager -------------------------------------------------------------------------------------------------- */

// TODO
func TestShop_CalculateOrder(t *testing.T) {
	type fields struct {
		Accounts map[string]*Account
		Products map[string]*Product
		Bundles  map[string]*Bundle
	}
	type args struct {
		order Order
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float32
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Market{
				Accounts: tt.fields.Accounts,
				Products: tt.fields.Products,
				Bundles:  tt.fields.Bundles,
			}
			got, err := m.CalculateOrder(tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CalculateOrder() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO
func TestShop_PlaceOrder(t *testing.T) {
	type fields struct {
		Accounts map[string]*Account
		Products map[string]*Product
		Bundles  map[string]*Bundle
	}
	type args struct {
		username string
		order    Order
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Market{
				Accounts: tt.fields.Accounts,
				Products: tt.fields.Products,
				Bundles:  tt.fields.Bundles,
			}
			if err := m.PlaceOrder(tt.args.username, tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("PlaceOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

/* -- BundleManager ------------------------------------------------------------------------------------------------- */

// TODO
func TestShop_AddBundle(t *testing.T) {
	type fields struct {
		Accounts map[string]*Account
		Products map[string]*Product
		Bundles  map[string]*Bundle
	}
	type args struct {
		name       string
		main       Product
		discount   float32
		additional []Product
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Market{
				Accounts: tt.fields.Accounts,
				Products: tt.fields.Products,
				Bundles:  tt.fields.Bundles,
			}
			if err := m.AddBundle(tt.args.name, tt.args.main, tt.args.discount, tt.args.additional...); (err != nil) != tt.wantErr {
				t.Errorf("AddBundle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TODO
func TestShop_ChangeDiscount(t *testing.T) {
	type fields struct {
		Accounts map[string]*Account
		Products map[string]*Product
		Bundles  map[string]*Bundle
	}
	type args struct {
		name     string
		discount float32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Market{
				Accounts: tt.fields.Accounts,
				Products: tt.fields.Products,
				Bundles:  tt.fields.Bundles,
			}
			if err := m.ChangeDiscount(tt.args.name, tt.args.discount); (err != nil) != tt.wantErr {
				t.Errorf("ChangeDiscount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TODO
func TestShop_RemoveBundle(t *testing.T) {
	type fields struct {
		Accounts map[string]*Account
		Products map[string]*Product
		Bundles  map[string]*Bundle
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Market{
				Accounts: tt.fields.Accounts,
				Products: tt.fields.Products,
				Bundles:  tt.fields.Bundles,
			}
			if err := m.RemoveBundle(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("RemoveBundle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

/* -- Importer, Exporter -------------------------------------------------------------------------------------------- */

// TODO
func TestShop_Export(t *testing.T) {
	type fields struct {
		Accounts map[string]*Account
		Products map[string]*Product
		Bundles  map[string]*Bundle
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Market{
				Accounts: tt.fields.Accounts,
				Products: tt.fields.Products,
				Bundles:  tt.fields.Bundles,
			}
			got, err := m.Export()
			if (err != nil) != tt.wantErr {
				t.Errorf("Export() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Export() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO
func TestShop_Import(t *testing.T) {
	type fields struct {
		Accounts map[string]*Account
		Products map[string]*Product
		Bundles  map[string]*Bundle
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Market{
				Accounts: tt.fields.Accounts,
				Products: tt.fields.Products,
				Bundles:  tt.fields.Bundles,
			}
			if err := m.Import(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Import() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

/* --- new ---------------------------------------------------------------------------------------------------------- */

func newMarket() Market {

	m := NewMarket()

	for _, acc := range newAccounts() {
		_ = m.Register(acc.Name)
		_ = m.AddBalance(acc.Name, acc.Balance)
		_ = m.ModifyAccountType(acc.Name, acc.Type)
	}

	for _, p := range newProducts() {
		_ = m.AddProduct(p)
	}

	for _, b := range newBundles() {
		_ = m.AddBundle(b.Main.Name, b.Main, b.Discount, b.Additional...)
	}

	return m
}

func newAccounts() []Account {
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

func newProducts() []Product {
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

func newBundles() []Bundle {
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
