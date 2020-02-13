package shop

import (
	"github.com/pkg/errors"
)

var (
	ErrorBundleSampleAdditionalProducts = errors.New("bundle sample cannot have additional products")
)

/* -- BundleManager ------------------------------------------------------------------------------------------------- */

func NewBundle(main Product, discount float32, bundleType BundleType, additional ...Product) Bundle {
	return Bundle{
		Products: append(additional, main),
		Type:     bundleType,
		Discount: discount,
	}
}

func (m *Market) AddBundle(name string, main Product, discount float32, additional ...Product) error {

	if discount < 1 || discount > 99 {
		return ErrorInvalidDiscount
	}

	if main.Type == ProductSampled {
		return errors.New("main product cannot be sampled")
	}

	sampled := getProductsWithType(additional, ProductSampled)
	if len(sampled) > 1 {
		return errors.New("too many sampled products")
	}
	if len(sampled) == 1 && len(additional) > 1 {
		return ErrorBundleSampleAdditionalProducts
	}

	if _, ok := m.Bundles[name]; ok {
		return ErrorBundleExists
	}

	bundleType := BundleNormal
	if len(sampled) == 1 {
		bundleType = BundleSample
	}

	b := NewBundle(main, discount, bundleType, additional...)

	m.Bundles[name] = b
	return nil
}

func (m *Market) ChangeDiscount(name string, discount float32) error {

	if discount < 1 || discount > 99 {
		return ErrorInvalidDiscount
	}

	if _, ok := m.Bundles[name]; !ok {
		return ErrorBundleNotExists
	}

	acc := m.Bundles[name]
	acc.Discount = discount

	m.Bundles[name] = acc
	return nil
}

func (m *Market) RemoveBundle(name string) error {

	if _, ok := m.Bundles[name]; !ok {
		return ErrorBundleNotExists
	}

	delete(m.Bundles, name)
	return nil
}

func getProductsWithType(products []Product, productType ProductType) (result []Product) {
	for _, p := range products {
		if p.Type == productType {
			result = append(result, p)
		}
	}
	return
}
