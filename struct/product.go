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

/* -- ProductManager ------------------------------------------------------------------------------------------------ */

func NewProduct(productName string, price float32, productType ProductType) Product {
	return Product{
		Name:  productName,
		Price: price,
		Type:  productType,
	}
}

func (m *Market) AddProduct(p Product) error {

	if err := checkProduct(p); err != nil {
		return errors.Wrap(err, "invalid check product")
	}

	if _, err := m.GetProduct(p.Name); err == nil {
		return ErrorProductAlreadyExists
	}
	return m.SetProduct(p.Name, p)
}

func (m *Market) ModifyProduct(p Product) error {

	if err := checkProduct(p); err != nil {
		return errors.Wrap(err, "invalid product check")
	}

	if _, err := m.GetProduct(p.Name); err != nil {
		return errors.Wrap(err, "cannot modify nil product")
	}
	return m.SetProduct(p.Name, p)
}

func (m *Market) RemoveProduct(name string) error {

	if _, err := m.GetProduct(name); err != nil {
		return errors.Wrap(err, "cannot delete nil product")
	}

	delete(m.Products, name)
	return nil
}

/* --- Interface ---------------------------------------------------------------------------------------------------- */

func (m *Market) GetProduct(name string) (Product, error) {
	product, ok := m.Products[name]

	if !ok {
		return Product{}, ErrorProductNotExists
	}

	return product, nil
}

func (m *Market) SetProduct(name string, product Product) error {

	if err := checkName(name); err != nil {
		return errors.Wrap(err, "can't set invalid product")
	}
	if err := checkProduct(product); err != nil {
		return errors.Wrap(err, "can't set invalid product")
	}

	m.Products[name] = product
	return nil
}

/* --- Checks ------------------------------------------------------------------------------------------------------  */

func checkProduct(product Product) error {

	if product.Price < 0 {
		return ErrorProductNegativePrice
	}

	if len(product.Name) == 0 {
		return ErrorEmptyField
	}

	if _, ok := ProductTypeMap[product.Type]; !ok {
		return ErrorProductInvalidType
	}

	return nil
}
