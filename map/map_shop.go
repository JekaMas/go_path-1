package main

import (
	"errors"
	"fmt"
	"sort"
)

var prices = make(map[string]int) // prices[productName]price

var savedOrders = make(map[string]int) // savedOrders[orderKey]price

var users = make(map[string]int) // users[userName]amount

/* --- Product ------------------------------------------------------------------------------------------------------ */

// AddProduct add new product if not exists
func AddProduct(productName string, price int) error {
	_, ok := prices[productName]
	if ok {
		return errors.New("product already exists")
	}

	prices[productName] = price
	return nil
}

// UpdatePrice update the price of an existing product
func UpdatePrice(productName string, price int) error {
	_, ok := prices[productName]
	if !ok {
		return errors.New("can't update nil product")
	}

	prices[productName] = price
	return nil
}

/* --- Orders ------------------------------------------------------------------------------------------------------- */

//
func OrderPrice(order map[string]int) int { // map[name]count
	price := 0

	for name, count := range order {
		price += prices[name] * count
	}

	return price
}

func OrderPriceCached(order map[string]int) int { // map[name]count
	key := fmt.Sprint(order)

	cached, ok := savedOrders[key]
	if ok {
		return cached
	}

	price := OrderPrice(order)
	savedOrders[key] = price

	return price
}

/* --- User --------------------------------------------------------------------------------------------------------- */

func AddUser(userName string, amount int) error {
	_, ok := users[userName]
	if ok {
		return errors.New("user already exists")
	}

	users[userName] = amount
	return nil
}

// MakeOrder
func MakeOrder(userName string, order map[string]int) error {
	price := OrderPriceCached(order)
	amount := users[userName]

	if amount >= price {
		users[userName] -= price
		return nil
	}

	return errors.New("insufficient funds")
}

// enum of userNames possible sort orders
const (
	SORT_USERS_BY_NAME = iota
	SORT_USERS_BY_NAME_REVERSED
	SORT_USERS_BY_AMOUNT
)

// userNames return slice of userNames in some order
func userNames(sortType int) []string {

	var keys []string

	switch sortType {
	default:
		fallthrough
	case SORT_USERS_BY_NAME:
		keys = keysSortedByName(users, false)
	case SORT_USERS_BY_NAME_REVERSED:
		keys = keysSortedByName(users, true)
	case SORT_USERS_BY_AMOUNT:
		keys = keysSortedByValue(users)
	}

	return keys
}

// keysSortedByName return keys of the map sorted in lexicographical order
func keysSortedByName(myMap map[string]int, reversed bool) []string {
	var keys []string
	for k := range myMap {
		keys = append(keys, k)
	}

	var toSort sort.Interface = sort.StringSlice(keys)
	if reversed {
		toSort = sort.Reverse(toSort)
	}

	sort.Sort(toSort)

	return keys
}

// keysSortedByValue returns slice of keys sorted by the map's value
func keysSortedByValue(myMap map[string]int) []string {
	type KeyValue struct {
		key   string
		value int
	}

	var kvs []KeyValue
	for k, v := range myMap {
		kvs = append(kvs, KeyValue{k, v})
	}

	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].value > kvs[j].value
	})

	var keys []string
	for _, kv := range kvs {
		keys = append(keys, kv.key)
	}

	return keys
}
