package shop

import (
	"errors"
	"reflect"
	"testing"
)

var (
	ErrorProductExist           = errors.New("product already exists")
	ErrorProductSampleWithPrice = errors.New("sample not has price")
)

func TestAddProductFailed(t *testing.T) {
	type productTest struct {
		testName string
		Product  Product
		err      error
	}

	tests := []productTest{
		{testName: "FreeProduct", Product: Product{Name: "Free", Price: 0, Type: ProductNormal}, err: ErrorNegativeProductPrice},
		{testName: "SampleWithPrice", Product: Product{Name: "Sample", Price: 100, Type: ProductSampled}, err: ErrorProductSampleWithPrice},
	}

	testShop := NewMarket()

	testShop.Products["Pineapple"] = Product{Name: "Pineapple", Price: 100, Type: ProductNormal}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			err := testShop.AddProduct(test.Product)
			if err == nil {
				t.Errorf("Test it should be failed, error: %v", test.Product)
			} else if err.Error() != test.err.Error() && !errors.Is(err, test.err) {
				t.Errorf("Values not equal: %v != %v", err, test.err)
			}
		})
	}
}

func TestModifyProductFailed(t *testing.T) {
	type productTest struct {
		testName string
		Product  Product
		err      error
	}

	tests := []productTest{
		{testName: "FreeProduct", Product: Product{Name: "Pineapple", Price: 0, Type: ProductNormal}, err: ErrorNegativeProductPrice},
		{testName: "SampleWithPrice", Product: Product{Name: "Pineapple", Price: 100, Type: ProductSampled}, err: errors.New("sample not has price")},
	}

	testShop := NewMarket()

	testShop.Products["Pineapple"] = Product{Name: "Pineapple", Price: 100, Type: ProductNormal}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			err := testShop.ModifyProduct(test.Product)
			if err == nil {
				t.Errorf("Test it should be failed with error: %v, but get value %v", test.err, test.Product)
			} else if err.Error() != test.err.Error() && !errors.Is(err, test.err) {
				t.Errorf("Values not equal: %v != %v", err, test.err)
			}
		})
	}
}

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

			if product := m.Products[productName]; !tt.wantErr && !reflect.DeepEqual(product, tt.product) {
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
