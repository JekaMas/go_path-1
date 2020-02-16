package shop

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
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

	m.accountsMutex.RLock()
	account, err := m.GetAccount(userName)
	m.accountsMutex.RUnlock()

	if err != nil {
		return 0, errors.Wrap(err, "cannot calculate on a nil account")
	}

	// cache key
	key := orderKey(account, order)
	// if exists, get from cache
	m.orderCacheMutex.RLock()
	if amount, ok := m.OrdersCache[key]; ok {
		return amount, nil
	}
	m.orderCacheMutex.RUnlock()

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

	amount := productsPrice + bundlesPrice

	m.orderCacheMutex.Lock()
	m.OrdersCache[key] = amount
	m.orderCacheMutex.Unlock()

	return amount, nil
}

func (m *Market) PlaceOrder(userName string, order Order) error {

	amount, err := m.CalculateOrder(userName, order)
	if err != nil {
		return errors.Wrap(err, "error during order calculation")
	}

	m.accountsMutex.Lock()
	defer m.accountsMutex.Unlock()

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

/* -- Util ---------------------------------------------------------------------------------------------------------- */

func orderKey(account Account, order Order) string {
	b := new(bytes.Buffer)
	for _, value := range order.Products { // FIXME handle errors
		_, _ = fmt.Fprintf(b, "%v", value)
	}
	for _, value := range order.Bundles {
		_, _ = fmt.Fprintf(b, "%v", value)
	}

	_, _ = fmt.Fprintf(b, "%v", account.Type)
	return b.String()
}
