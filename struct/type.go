package shop

const (
	ProductNormal ProductType = iota
	ProductPremium
	ProductSampled

	BundleNormal BundleType = iota
	BundleSample

	AccountNormal AccountType = iota
	AccountPremium

	SortByName AccountSortType = iota
	SortByNameReverse
	SortByBalance
)

var DiscountMap = map[ProductType]map[AccountType]float32{
	ProductPremium: {AccountPremium: -20, AccountNormal: -5},
	ProductNormal:  {AccountPremium: +50, AccountNormal: 0},
}

type ProductType uint8     // ProductNormal, ProductPremium, ProductSample
type BundleType uint8      // BundleNormal, BundleSample
type AccountType uint8     // AccountNormal, AccountSample
type AccountSortType uint8 // SortByName, SortByNameReverse, SortByBalance

type Product struct {
	Name  string
	Price float32
	Type  ProductType
}

type Order struct {
	Products []Product
	Bundles  []Bundle
	Account
}

type Bundle struct {
	Main       Product
	Additional []Product
	Type       BundleType
	Discount   float32
}

type Account struct {
	Name    string
	Balance float32
	Type    AccountType
}

type Market struct {
	Accounts map[string]*Account // map[username]Account
	Products map[string]*Product // map[productName]Product
	Bundles  map[string]*Bundle  // map[bundleName]Bundle
}
