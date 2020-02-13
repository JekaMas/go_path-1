package shop

import "github.com/pkg/errors"

/* -- ProductManager ------------------------------------------------------------------------------------------------ */

func NewProduct(productName string, price float32, productType ProductType) Product {
	return Product{
		Name:  productName,
		Price: price,
		Type:  productType,
	}
}

func (m *Market) AddProduct(p Product) error {

	if err := m.checkProduct(p); err != nil {
		return errors.Wrap(err, "invalid check product")
	}

	if _, ok := m.Products[p.Name]; ok {
		return errors.New("product already exists")
	}

	m.Products[p.Name] = p
	return nil
}

func (m *Market) ModifyProduct(p Product) error {

	if err := m.checkProduct(p); err != nil {
		return errors.Wrap(err, "invalid product check")
	}

	if _, ok := m.Products[p.Name]; !ok {
		return errors.New("cannot modify nil product")
	}

	m.Products[p.Name] = p
	return nil
}

func (m *Market) RemoveProduct(name string) error {

	if _, ok := m.Products[name]; !ok {
		return errors.New("cannot delete nil product")
	}

	delete(m.Products, name)
	return nil
}

func (m *Market) checkProduct(p Product) error {

	if p.Price < 0 {
		return ErrorNegativeProductPrice
	}

	if len(p.Name) == 0 {
		return ErrorEmptyField
	}

	if !(p.Type == ProductSampled || p.Type == ProductPremium || p.Type == ProductNormal) {
		return errors.New("no such product type")
	}

	return nil
}
