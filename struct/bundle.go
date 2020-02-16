package shop

import (
	"github.com/pkg/errors"
)

var (
	ErrorBundleSampleAdditionalProducts = errors.New("bundle sample cannot have additional products")
	ErrorBundleTooManySampled           = errors.New("too many sampled products")
	ErrorBundleMainSample               = errors.New("main product cannot be sampled")
	ErrorBundleExists                   = errors.New("bundle already exists")
	ErrorBundleNotExists                = errors.New("bundle not exists")
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

	if main.Type == ProductSampled {
		return ErrorBundleMainSample
	}

	sampled := getProductsWithType(additional, ProductSampled)

	bundleType := BundleNormal
	if len(sampled) == 1 {
		bundleType = BundleSample
	}

	b := NewBundle(main, discount, bundleType, additional...)

	if _, err := m.GetBundle(name); err == nil {
		return ErrorBundleExists
	}
	return m.SetBundle(name, b)
}

func (m *Market) ChangeDiscount(name string, discount float32) error {

	if discount < 1 || discount > 99 {
		return ErrorInvalidDiscount
	}

	bundle, err := m.GetBundle(name)
	if err != nil {
		return errors.Wrap(err, "can't change discount of the nil bundle")
	}

	bundle.Discount = discount
	return m.SetBundle(name, bundle)
}

func (m *Market) RemoveBundle(name string) error {

	_, err := m.GetBundle(name)
	if err != nil {
		return errors.Wrap(err, "can't delete nil bundle")
	}

	delete(m.Bundles, name)
	return nil
}

/* --- Interface ---------------------------------------------------------------------------------------------------- */

func (m *Market) GetBundle(name string) (Bundle, error) {
	bundle, ok := m.Bundles[name]

	if !ok {
		return Bundle{}, ErrorBundleNotExists
	}

	return bundle, nil
}

func (m *Market) SetBundle(name string, bundle Bundle) error {

	if err := checkName(name); err != nil {
		return errors.Wrap(err, "can't set invalid bundle")
	}
	if err := checkBundle(bundle); err != nil {
		return errors.Wrap(err, "can't set invalid bundle")
	}

	m.Bundles[name] = bundle
	return nil
}

/* --- Check -------------------------------------------------------------------------------------------------------- */

func checkBundle(bundle Bundle) error {

	if bundle.Discount < 1 || bundle.Discount > 99 {
		return ErrorInvalidDiscount
	}

	sampled := getProductsWithType(bundle.Products, ProductSampled)
	if len(sampled) > 1 {
		return ErrorBundleTooManySampled
	}

	if len(sampled) == 1 &&
		len(bundle.Products)-1 > 1 { // except one sample
		return ErrorBundleSampleAdditionalProducts
	}

	return nil
}

/* --- Util --------------------------------------------------------------------------------------------------------- */

func getProductsWithType(products []Product, productType ProductType) (result []Product) {
	for _, p := range products {
		if p.Type == productType {
			result = append(result, p)
		}
	}
	return
}
