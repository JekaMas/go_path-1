package shop

import (
	"bytes"
	"fmt"
	"sync"
)

type OrdersCache struct {
	OrdersCache map[string]float32 // map[key]sum
	mu          sync.RWMutex
}

func NewOrdersCache() OrdersCache {
	return OrdersCache{
		OrdersCache: make(map[string]float32),
		mu:          sync.RWMutex{},
	}
}

func (c OrdersCache) GetCached(account Account, order Order) (float32, bool) {
	// cache key
	key := orderKey(account, order)
	// if exists, get from cache
	c.mu.RLock()
	defer c.mu.RUnlock()

	if amount, ok := c.OrdersCache[key]; ok {
		return amount, true
	}

	return 0, false
}

func (c OrdersCache) PutCached(account Account, order Order, amount float32) {
	key := orderKey(account, order) // fixme too long?

	c.mu.Lock()
	defer c.mu.Unlock()

	c.OrdersCache[key] = amount
}

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
