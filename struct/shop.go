package shop

import (
	"encoding/json"
	"github.com/pkg/errors"
)

var (
	ErrorInvalidDiscount = errors.New("invalid discount")
	ErrorEmptyField      = errors.New("empty field")
)

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
