package concurrency

import shop "go_path/struct"

type Market struct {
	shop *shop.Market
}

func NewMarket() Market {
	m := shop.NewMarket()
	return Market{&m}
}
