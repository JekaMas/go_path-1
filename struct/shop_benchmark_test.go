package shop

import (
	"testing"
)

func BenchmarkMarket_CalculateOrder(b *testing.B) {

	nt := ProductNormal  // normal type
	pt := ProductPremium // premium type

	acc := Account{"A", 0, AccountNormal}

	orderA := NewOrder([]Product{NewProduct("P1", 90, pt), NewProduct("P2", 10, pt)}, nil, acc)
	orderB := NewOrder([]Product{}, []Bundle{NewBundle(NewProduct("P1", 10, nt), -1, NewProduct("P2", 90, nt))}, acc)

	var orders []Order

	for i := 0; i < 100; i++ {
		orders = append(orders, orderA)
		orders = append(orders, orderB)
	}

	m := testMarket()

	b.ResetTimer()
	for _, order := range orders {
		_, _ = m.CalculateOrder(order)
	}
}
