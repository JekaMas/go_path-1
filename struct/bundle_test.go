package shop

import (
	"errors"
	"reflect"
	"testing"
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
			err: ErrorBundleSampleAdditionalProducts},
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

/* -- BundleManager ------------------------------------------------------------------------------------------------- */

func TestShop_AddBundle(t *testing.T) {

	type test struct {
		name       string
		main       Product
		discount   float32
		bundleType BundleType
		additional []Product
		wantErr    bool
	}

	nt := ProductNormal  // normal type
	pt := ProductPremium // premium type
	st := ProductSampled

	tests := []test{
		{"default", NewProduct("P1", 10, nt), 1, BundleNormal, []Product{NewProduct("P2", 90, nt)}, false},
		{"default2", NewProduct("P1", 10, nt), 10, BundleNormal, []Product{NewProduct("P2", 90, nt)}, false},
		{"default3", NewProduct("P1", 10, nt), 55.2345, BundleNormal, []Product{NewProduct("P2", 90, nt)}, false},
		{"default4", NewProduct("P1", 10, nt), 76.546767, BundleNormal, []Product{NewProduct("P2", 90, nt)}, false},

		{"errDisc", NewProduct("P1", 10, nt), 0, BundleNormal, []Product{NewProduct("P2", 90, nt)}, true},
		{"errDisc2", NewProduct("P1", 10, nt), 100, BundleNormal, []Product{NewProduct("P2", 90, nt)}, true},
		{"errDisc3", NewProduct("P1", 10, nt), 101, BundleNormal, []Product{NewProduct("P2", 90, nt)}, true},
		{"errDisc4", NewProduct("P1", 10, nt), 99.12456, BundleNormal, []Product{NewProduct("P2", 90, nt)}, true},
		{"errDisc5", NewProduct("P1", 10, nt), 99.00001, BundleNormal, []Product{NewProduct("P2", 90, nt)}, true},
		{"errDisc6", NewProduct("P1", 10, nt), 200, BundleNormal, []Product{NewProduct("P2", 90, nt)}, true},
		{"negDiscount", NewProduct("P1", 10, nt), -1, BundleNormal, []Product{NewProduct("P2", 90, nt)}, true},

		{"nilProd", Product{}, 0, BundleNormal, []Product{NewProduct("P2", 90, nt)}, true},
		{"nilProd2", NewProduct("P1", 10, nt), 0, BundleNormal, []Product{{}}, true},
		{"nilProd3", NewProduct("P1", 10, nt), 0, BundleNormal, []Product{NewProduct("P2", 90, nt), {}}, true},
		{"nilProd4", NewProduct("P1", 10, nt), 0, BundleNormal, []Product{{}, NewProduct("P2", 90, nt)}, true},
		{"nilProd5", Product{}, 0, BundleNormal, []Product{NewProduct("P2", 90, nt), {}}, true},
		// prem
		{"defaultPrem", NewProduct("P1", 10, nt), 1, BundleNormal, []Product{NewProduct("P2", 90, pt)}, false},
		{"errDiscPrem2", NewProduct("P1", 10, pt), 200, BundleNormal, []Product{NewProduct("P2", 90, nt)}, true},
		{"nilProdPrem3", Product{}, 0, BundleNormal, []Product{NewProduct("P2", 90, pt), {}}, true},

		// sampled
		{"defaultSampled", NewProduct("P1", 10, nt), 1, BundleSample, []Product{NewProduct("P2", 90, st)}, false},
		{"errSampled", NewProduct("P1", 10, st), 1, BundleSample, []Product{NewProduct("P2", 90, nt)}, true},
		{"errSampled2", NewProduct("P1", 10, nt), 1, BundleSample, []Product{NewProduct("P2", 90, st), NewProduct("P2", 90, st)}, true},
	}

	// test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMarket()
			err := m.AddBundle(tt.name, tt.main, tt.discount, tt.additional...)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddBundle() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(NewBundle(tt.main, tt.discount, tt.bundleType, tt.additional...), m.Bundles[tt.name]) {
				t.Errorf("AddBundle() wrong bundle added = %v get = %v",
					NewBundle(tt.main, tt.discount, tt.bundleType, tt.additional...), m.Bundles[tt.name])
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
		{"default1", 2, false},
		{"default2", 50, false},
		{"err", 150, true},
		{"err2", 0, true},
		{"err2", 100, true},
		{"err2", 100.0001, true},
		{"err3", -10, true},
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
