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

func (m *Market) ChangeDiscount(name string, discount float32) error {
	if discount < 1 || discount > 99 {
		return ErrorInvalidDiscount
	}

	return m.changeBundle(name, changeDiscountFunc(discount))
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
