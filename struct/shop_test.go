package shop

import (
	"reflect"
	"testing"
)

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

/* -- Additional ---------------------------------------------------------------------------------------------------- */

func TestMarketImplementsShop(t *testing.T) {

	m := NewMarket()
	var i interface{} = &m

	_, ok := i.(Shop)

	if !ok {
		t.Fatal("market doesn't implements shop interface")
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
		_ = m.AddBundle(b.Products[0].Name, b.Products[0], b.Discount, b.Products[1:]...)
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
			Products: []Product{
				{Name: "Bowling Ball", Price: 100, Type: ProductNormal},
				{Name: "Skittle", Price: 50, Type: ProductNormal},
				{Name: "Skittle", Price: 50, Type: ProductNormal},
				{Name: "Skittle", Price: 50, Type: ProductPremium},
			},
		},
		{
			Type:     BundleSample,
			Discount: 1,
			Products: []Product{
				Product{Name: "Toy Gun", Price: 0, Type: ProductPremium},
				{Name: "Battery", Price: 30, Type: ProductNormal},
				{Name: "Battery", Price: 30, Type: ProductNormal},
			},
		},
	}
}
