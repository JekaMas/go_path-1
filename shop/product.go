package shop

import "github.com/pkg/errors"

var (
	ErrorProductNegativePrice   = errors.New("product price is negative")
	ErrorProductInvalidType     = errors.New("no such product type")
	ErrorProductAlreadyExists   = errors.New("product already exists")
	ErrorProductNotExists       = errors.New("product not exists")
	ErrorProductExist           = errors.New("product already exists")
	ErrorProductSampleWithPrice = errors.New("sample not has price")
)

func NewProduct(productName string, price float32, productType ProductType) Product {
	return Product{
		Name:  productName,
		Price: price,
		Type:  productType,
	}
}

func (m *Market) ModifyProduct(p Product) error {

	if err := checkProduct(p); err != nil {
		return errors.Wrap(err, "invalid product check")
	}

	return m.changeProduct(p.Name, func(product *Product) {
		product.Name = p.Name
		product.Type = p.Type
		product.Price = p.Price
	})
}
