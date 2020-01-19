package shop

type Shop interface {
	ProductModifier
	AccountManager
	OrderManager
	BundleManager
	Importer
	Exporter
}

type ProductModifier interface {
	AddProduct(Product) error
	ModifyProduct(Product) error
	RemoveProduct(name string) error
}

type AccountManager interface {
	Register(username string) error
	AddBalance(username string, sum float32) error
	Balance(username string) (float32, error)
	GetAccounts(sortType AccountSortType) []Account
}

type OrderManager interface {
	CalculateOrder(order Order) (float32, error)
	PlaceOrder(username string, order Order) error
}

type BundleManager interface {
	AddBundle(name string, main Product, discount float32, additional ...Product) error
	ChangeDiscount(name string, discount float32) error
	RemoveBundle(name string) error
}

type Exporter interface {
	Export() ([]byte, error)
}

type Importer interface {
	Import(data []byte) error
}
