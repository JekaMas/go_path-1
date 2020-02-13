package shop

import (
	"encoding/json"
	"github.com/pkg/errors"
)

var ErrorBundleExists = errors.New("bundle already exists")
var ErrorBundleNotExists = errors.New("bundle not exists")
var ErrorInvalidDiscount = errors.New("invalid discount")
var ErrorEmptyField = errors.New("empty field")
var ErrorNegativeProductPrice = errors.New("product price is negative")

func NewMarket() Market {
	return Market{
		Accounts:    make(map[string]Account),
		Products:    make(map[string]Product),
		Bundles:     make(map[string]Bundle),
		OrdersCache: make(map[string]float32),
	}
}

/* -- Importer, Exporter -------------------------------------------------------------------------------------------- */

func (m *Market) Import(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Market) Export() ([]byte, error) {
	return json.MarshalIndent(m, "", "")
}
