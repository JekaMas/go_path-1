package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAddProduct(t *testing.T) {
	// ignore adding errors handling
	_ = AddProduct("Apple", 10)
	_ = AddProduct("Milk", 10)
	_ = AddProduct("Butter", 20)
	_ = AddProduct("Bread", 30)

	result := prices
	test := map[string]int{"Apple": 10, "Bread": 30, "Butter": 20, "Milk": 10}

	if !reflect.DeepEqual(result, test) {
		t.Fatalf("Products weren't added to the map\n "+
			"(%v vs %v).", result, test)
	}
}

func TestUpdatePrice(t *testing.T) {
	oldPrice := 10
	newPrice := 20

	_ = AddProduct("Apple", oldPrice)
	_ = UpdatePrice("Apple", newPrice)

	result := prices["Apple"]

	if result != newPrice {
		t.Fatalf("Can't updatem price of the product\n "+
			"(%v vs %v).", result, newPrice)
	}
}

func TestOrderPrice(t *testing.T) {
	_ = AddProduct("Apple", 10)
	_ = AddProduct("Milk", 10)
	_ = AddProduct("Butter", 20)
	_ = AddProduct("Bread", 30)

	test := 400
	result := OrderPrice(map[string]int{
		"Apple":  3,
		"Milk":   5,
		"Butter": 4,
		"Bread":  8})

	if result != test {
		if result != test {
			t.Fatalf("Incorrect order price\n "+
				"(%v vs %v).", result, test)
		}
	}
}

func BenchmarkOrderPrice(b *testing.B) {
	_ = AddProduct("Apple", 10)
	_ = AddProduct("Milk", 10)

	order := map[string]int{"Apple": 10, "Milk": 20}

	for i := 0; i < 1_000_000; i++ {
		_ = OrderPrice(order)
	}
}

func BenchmarkOrderPriceCached(b *testing.B) {
	_ = AddProduct("Apple", 10)
	_ = AddProduct("Milk", 10)

	order := map[string]int{"Apple": 10, "Milk": 20}

	for i := 0; i < 1_000_000; i++ {
		_ = OrderPriceCached(order)
	}
}

func TestAddUser(t *testing.T) {
	_ = AddUser("Vasya", 100)
	_ = AddUser("Vova", 1_000_000)

	result := users
	test := map[string]int{"Vasya": 100, "Vova": 1_000_000}

	if !reflect.DeepEqual(result, test) {
		t.Fatalf("Users weren't added to the map\n "+
			"(%v vs %v).", result, test)
	}
}

func TestMakeOrder(t *testing.T) {
	_ = AddProduct("Apple", 10)
	_ = AddProduct("Milk", 10)

	_ = AddUser("Vasya", 100)
	_ = AddUser("Vova", 1_000_000)

	users := users
	order := map[string]int{"Apple": 5, "Milk": 4}
	_ = MakeOrder("Vasya", order)

	if users["Vasya"] != 10 {
		t.Fatalf("Wrong user amount value\n "+
			"(%v vs %v).", users["Vasya"], 10)
	}

	ok := MakeOrder("Vasya", order)

	if ok == nil {
		t.Fatalf("The amount of money in the account "+
			"went into a negative value.\n "+
			"(%v).", users["Vasya"])
	}

	// t.Logf("Total amount '%v': %v", "Vasya", users["Vasya"])
}

func TestuserNames(t *testing.T) {
	_ = AddUser("Vova", 1_000_000)
	_ = AddUser("Katya", 1_999)
	_ = AddUser("Petya", 250)
	_ = AddUser("Vasya", 100)

	printUsers := func(keys []string) {
		for _, k := range keys {
			fmt.Println(k, " ", users[k])
		}
		fmt.Println()
	}

	names := userNames(SORT_USERS_BY_NAME)
	printUsers(names)
	names = userNames(SORT_USERS_BY_NAME_REVERSED)
	printUsers(names)
	names = userNames(SORT_USERS_BY_AMOUNT)
	printUsers(names)
}
