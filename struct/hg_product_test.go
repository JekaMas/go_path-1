package shop

import (
	"errors"
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
