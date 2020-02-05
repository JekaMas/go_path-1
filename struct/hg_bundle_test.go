package shop

import (
	"errors"
	"testing"
)

var (
	ErrorBundleSampleNotHasOtherProduct = errors.New("bundle sample not has other product")
)

func TestAddBundleSuccess(t *testing.T) {
	type bundleTest struct {
		testName    string
		mainProduct Product
		discount    float32
		addProduct  []Product
		bundleType  BundleType
	}

	tests := []bundleTest{
		{testName: "SampleType",
			mainProduct: Product{Name: "Main", Price: 1000, Type: ProductNormal},
			discount:    50,
			addProduct: []Product{{
				Name:  "Second",
				Price: 1000,
				Type:  ProductSampled,
			}},
			bundleType: BundleSample},
	}

	testShop := NewMarket()
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			err := testShop.AddBundle(test.testName, test.mainProduct, test.discount, test.addProduct...)
			if err != nil {
				t.Errorf("Test it should be success, error: %v", err)
			}
			if testShop.Bundles[test.testName].Type != test.bundleType {
				t.Errorf("Error, type not correct: %v != %v", testShop.Bundles[test.testName].Type, test.bundleType)
			}
		})
	}

}

//Boris, [28.01.20 00:25]
//Например если покупаешь шампунь, кондиционер и расческу - скидка 20%
//Boris, [28.01.20 00:25]
//Но все обычные товары
//Дмитрий Майоров, [28.01.20 00:25]
//а в случает типа sampleBundle в наборе хранится только шампунь и пробник?
//Boris, [28.01.20 00:26]
//Да
func TestAddBundleFailed(t *testing.T) {
	type bundleTest struct {
		testName    string
		mainProduct Product
		discount    float32
		addProduct  []Product
		err         error
	}

	tests := []bundleTest{
		{testName: "SampleWithNotSample",
			mainProduct: Product{Name: "Main", Price: 1000, Type: ProductNormal},
			discount:    50,
			addProduct: []Product{{
				Name:  "Second",
				Price: 1000,
				Type:  ProductSampled,
			},
				{
					Name:  "Third",
					Price: 1000,
					Type:  ProductPremium,
				}},
			err: ErrorBundleSampleNotHasOtherProduct},
	}

	testShop := NewMarket()
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			err := testShop.AddBundle(test.testName, test.mainProduct, test.discount, test.addProduct...)
			if err == nil {
				t.Errorf("Test it should be failed with error: %v", test.err)
			} else if err.Error() != test.err.Error() && !errors.Is(err, test.err) {
				t.Errorf("Values not equal: %v != %v", err, test.err)
			}
		})
	}

}
