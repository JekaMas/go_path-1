package shop

import (
	"github.com/pkg/errors"
	"sync"
)

type Products struct {
	Products map[string]Product
	mu       sync.RWMutex
}

type changeProductFunc func(product *Product)

func NewProducts() Products {
	return Products{
		Products: make(map[string]Product),
		mu:       sync.RWMutex{},
	}
}

func (p *Products) AddProduct(product Product) error {

	if err := checkProduct(product); err != nil {
		return errors.Wrap(err, "invalid check product")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, err := p.getProduct(product.Name); err == nil {
		return ErrorProductAlreadyExists
	}
	return p.setProduct(product.Name, product)
}

func (p *Products) RemoveProduct(name string) error {

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, err := p.getProduct(name); err != nil {
		return errors.Wrap(err, "cannot delete nil product")
	}

	delete(p.Products, name)
	return nil
}

func (p *Products) changeProduct(name string, fns ...changeProductFunc) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	acc, err := p.getProduct(name)
	if err != nil {
		return errors.Wrap(err, "can't change the nil product")
	}

	for _, fn := range fns {
		fn(&acc)
	}

	p.Products[name] = acc
	return nil
}

/* --- Interface ---------------------------------------------------------------------------------------------------- */

func (p *Products) GetProduct(name string) (Product, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.getProduct(name)
}

func (p *Products) getProduct(name string) (Product, error) {
	product, ok := p.Products[name]

	if !ok {
		return Product{}, ErrorProductNotExists
	}

	return product, nil
}

func (p *Products) SetProduct(name string, product Product) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.setProduct(name, product)
}

func (p *Products) setProduct(name string, product Product) error {

	if err := checkName(name); err != nil {
		return errors.Wrap(err, "can't set invalid product")
	}
	if err := checkProduct(product); err != nil {
		return errors.Wrap(err, "can't set invalid product")
	}

	p.Products[name] = product
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
