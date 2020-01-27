package shop

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
)

var ErrorAccountNotRegistered = errors.New("account is not registered")
var ErrorAccountExists = errors.New("account already exists")
var ErrorBundleExists = errors.New("bundle already exists")
var ErrorBundleNotExists = errors.New("bundle not exists")
var ErrorInvalidDiscount = errors.New("invalid discount")
var ErrorEmptyField = errors.New("empty field")
var ErrorNegativeProductPrice = errors.New("product price is negative")

func NewMarket() Market {
	return Market{
		Accounts:    make(map[string]*Account),
		Products:    make(map[string]*Product),
		Bundles:     make(map[string]*Bundle),
		OrdersCache: make(map[string]float32),
	}
}

/* -- ProductManager ------------------------------------------------------------------------------------------------ */

func NewProduct(productName string, price float32, productType ProductType) Product {
	return Product{
		Name:  productName,
		Price: price,
		Type:  productType,
	}
}

func (m *Market) AddProduct(p Product) error {

	if len(p.Name) == 0 {
		return ErrorEmptyField
	}

	if p.Price < 0 {
		return ErrorNegativeProductPrice
	}

	if _, ok := m.Products[p.Name]; ok {
		return errors.New("product already exists")
	}

	m.Products[p.Name] = &p
	return nil
}

func (m *Market) ModifyProduct(p Product) error {

	if _, ok := m.Products[p.Name]; !ok {
		return errors.New("cannot modify nil product")
	}
	if !m.checkProduct(p) {
		return errors.New("invalid product check")
	}

	m.Products[p.Name] = &p
	return nil
}

func (m *Market) RemoveProduct(name string) error {

	if _, ok := m.Products[name]; !ok {
		return errors.New("cannot delete nil product")
	}

	delete(m.Products, name)
	return nil
}

func (m *Market) checkProduct(p Product) bool {
	return p.Price >= 0 && len(p.Name) > 0 && // FIXME check product type
		(p.Type == ProductSampled || p.Type == ProductPremium || p.Type == ProductNormal)
}

/* -- AccountManager ------------------------------------------------------------------------------------------------ */

func NewAccount(username string) Account {
	return Account{
		Name:    username,
		Balance: 0,
		Type:    AccountNormal,
	}
}

func (m *Market) Register(username string) error {

	if len(username) == 0 {
		return ErrorEmptyField
	}

	if _, ok := m.Accounts[username]; ok {
		return ErrorAccountExists
	}

	acc := NewAccount(username)
	m.Accounts[username] = &acc
	return nil
}

func (m *Market) AddBalance(username string, sum float32) error {

	if sum < 0 {
		return errors.New("cannot add negative sum")
	}

	if _, ok := m.Accounts[username]; !ok {
		return ErrorAccountNotRegistered
	}

	acc := m.Accounts[username]
	acc.Balance += sum
	return nil
}

func (m *Market) ModifyAccountType(username string, accountType AccountType) error {

	if _, ok := m.Accounts[username]; !ok {
		return ErrorAccountNotRegistered
	}

	acc := m.Accounts[username]
	acc.Type = accountType
	return nil
}

func (m *Market) Balance(username string) (float32, error) {

	if _, ok := m.Accounts[username]; !ok {
		return 0, ErrorAccountNotRegistered
	}

	return m.Accounts[username].Balance, nil
}

func (m *Market) GetAccount(name string) Account {
	return *m.Accounts[name]
}

func (m *Market) GetAccounts(sortType AccountSortType) []Account {
	var accs []Account
	for _, acc := range m.Accounts {
		accs = append(accs, *acc)
	}
	// compare function
	var less func(i, j int) bool

	switch sortType {
	default:
		fallthrough
	case SortByName:
		less = func(i, j int) bool {
			return strings.Compare(accs[i].Name, accs[j].Name) < 0
		}
	case SortByNameReverse:
		less = func(i, j int) bool {
			return strings.Compare(accs[i].Name, accs[j].Name) > 0
		}
	case SortByBalance:
		less = func(i, j int) bool {
			return accs[i].Balance < accs[j].Balance
		}
	}

	sort.Slice(accs, less)
	return accs
}

/* -- OrderManager -------------------------------------------------------------------------------------------------- */

func NewOrder(products []Product, bundles []Bundle, account Account) Order {
	return Order{
		Products: products,
		Bundles:  bundles,
		Account:  account,
	}
}

func (m *Market) CalculateOrder(order Order) (float32, error) {

	if order.Products == nil {
		return 0, errors.New("products in the order is nil")
	}

	account := order.Account
	key := orderKey(order)

	// if exists, get from cache
	if amount, ok := m.OrdersCache[key]; ok {
		return amount, nil
	}

	// products
	productsPrice := float32(0)
	for _, product := range order.Products {

		if product.Type == ProductSampled {
			return 0, errors.New("sampled product cannot be bought without bundle")
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
			price += product.Price
		}
		bundlesPrice += price * (1 - bundle.Discount*0.01)
	}

	amount := productsPrice + bundlesPrice
	m.OrdersCache[key] = amount

	return amount, nil
}

func (m *Market) PlaceOrder(username string, order Order) error {

	amount, err := m.CalculateOrder(order)
	if err != nil {
		return err
	}

	acc, ok := m.Accounts[username]
	if !ok {
		return ErrorAccountNotRegistered
	}

	if acc.Balance < amount {
		return errors.New("insufficient funds")
	}

	acc.Balance -= amount
	return nil
}

func orderKey(order Order) string {
	b := new(bytes.Buffer)
	for _, value := range order.Products { // FIXME handle errors
		_, _ = fmt.Fprintf(b, "%v", value)
	}
	for _, value := range order.Bundles {
		_, _ = fmt.Fprintf(b, "%v", value)
	}

	_, _ = fmt.Fprintf(b, "%v", order.Account.Type)
	return b.String()
}

/* -- BundleManager ------------------------------------------------------------------------------------------------- */

func NewBundle(main Product, discount float32, additional ...Product) Bundle {
	return Bundle{
		Products: append(additional, main),
		Type:     BundleNormal,
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

	if _, ok := m.Bundles[name]; ok {
		return ErrorBundleExists
	}

	b := NewBundle(main, discount, additional...)
	m.Bundles[name] = &b
	return nil
}

func (m *Market) ChangeDiscount(name string, discount float32) error {

	if discount < 1 || discount > 99 {
		return ErrorInvalidDiscount
	}

	if _, ok := m.Bundles[name]; !ok {
		return ErrorBundleNotExists
	}

	m.Bundles[name].Discount = discount
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

/* -- Importer, Exporter -------------------------------------------------------------------------------------------- */

func (m *Market) Import(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Market) Export() ([]byte, error) {
	return json.MarshalIndent(m, "", "")
}
