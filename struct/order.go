package shop

import (
	"github.com/pkg/errors"
)

var (
	ErrorPriceNotNegative     = errors.New("total price not be negative/zero")
	ErrorIncorrectAccountType = errors.New("incorrect account type")
)

/* -- OrderManager -------------------------------------------------------------------------------------------------- */

func NewOrder(products []Product, bundles []Bundle) Order {
	return Order{
		Products: products,
		Bundles:  bundles,
	}
}

func (m *Market) CalculateOrder(userName string, order Order) (float32, error) {

	if order.Products == nil && order.Bundles == nil {
		return 0, errors.New("products in the order is nil")
	}

	account, err := m.GetAccount(userName)

	if err != nil {
		return 0, errors.Wrap(err, "cannot calculate on a nil account")
	}

	amount, ok := m.GetCached(account, order)
	if ok {
		return amount, nil
	}

	// products
	productsPrice := float32(0)
	for _, product := range order.Products {

		if product.Type == ProductSampled {
			return 0, errors.New("sampled product cannot be bought without bundle")
		}
		if product.Price < 0 {
			return 0, ErrorProductNegativePrice
		}

		discount := DiscountMap[product.Type][account.Type]
		productsPrice += product.Price * (1 - discount*0.01)
	}

	// bundles
	bundlesPrice := float32(0)
	for _, bundle := range order.Bundles {

		if bundle.Discount < 1 || bundle.Discount > 99 {
			return 0, ErrorInvalidDiscount
		}

		price := float32(0)
		for _, product := range bundle.Products {
			if product.Price < 0 {
				return 0, ErrorProductNegativePrice
			}

			price += product.Price
		}
		bundlesPrice += price * (1 - bundle.Discount*0.01)
	}

	amount = productsPrice + bundlesPrice
	m.PutCached(account, order, amount)

	return amount, nil
}

func (m *Market) PlaceOrder(userName string, order Order) error {

	amount, err := m.CalculateOrder(userName, order)
	if err != nil {
		return errors.Wrap(err, "error during order calculation")
	}

	acc, err := m.GetAccount(userName)
	if err != nil {
		return errors.Wrap(err, "can't place order to the nil account")
	}

	if acc.Balance < amount {
		return errors.New("insufficient funds")
	}

	acc.Balance -= amount
	return m.SetAccount(userName, acc)
}
