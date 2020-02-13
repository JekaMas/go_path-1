package shop

// Shop - сборный интерфейс магазина.
type Shop interface {
	ProductModifier
	AccountManager
	OrderManager
	BundleManager
	Importer
	Exporter
}

// ProductModifier - интерфейс дя работы со списком продуктов магазина.
type ProductModifier interface {
	AddProduct(Product) error
	ModifyProduct(Product) error
	RemoveProduct(name string) error
}

// AccountManager - интерфейс для работы с пользователями.
type AccountManager interface {
	Register(userName string) error
	AddBalance(userName string, sum float32) error
	Balance(userName string) (float32, error)
	GetAccounts(sortType AccountSortType) []Account
}

// OrderManager - интерфейс для работы с заказами.
type OrderManager interface {
	CalculateOrder(userName string, order Order) (float32, error)
	PlaceOrder(userName string, order Order) error
}

// BundleManager - интерфейс для работы с наборами.
type BundleManager interface {
	AddBundle(name string, main Product, discount float32, additional ...Product) error
	ChangeDiscount(name string, discount float32) error
	RemoveBundle(name string) error
}

// Exporter - интерфейс для получения полного состояния магазина.
type Exporter interface {
	Export() ([]byte, error)
}

// Importer - интерфейс для загрузки состояния магазина
// принимает формат, который возвращает Exporter.
type Importer interface {
	Import(data []byte) error
}