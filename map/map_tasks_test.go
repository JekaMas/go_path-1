package main

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestCountWordAndInText(t *testing.T) {
	text := "All the world is made of faith, and trust, and pixie dust."

	result := CountWords(text)

	if result["and"] != 2 {
		t.Fatalf("Counting the word 'and' is incorrect "+
			"(%v vs %v).", result["and"], 2)
	}
}

func TestAtLeastOne(t *testing.T) {
	slice := []int{1, 1, 1, 2, 2, 3}
	test := map[int]struct{}{1: {}, 2: {}, 3: {}}

	result := AtLeastOne(slice)

	if !reflect.DeepEqual(result, test) {
		t.Fatalf("The resulting map is not equal to the test one \n "+
			"(%v vs %v).", result, test)
	}
}

func BenchmarkAtLeastOne(b *testing.B) {
	slice := make([]int, 1_000_000)
	rand.Seed(time.Now().UnixNano())
	// fill slice with rand values
	for i := range slice {
		slice[i] = rand.Intn(100)
	}

	b.ResetTimer() // test only func
	AtLeastOne(slice)
}

func TestIntersection(t *testing.T) {
	slice1 := []int{1, 1, 1, 7, 2, 2, 3, 4, 4, 5}
	slice2 := []int{9, 2, 0, 3, 21, -1}
	test := []int{2, 3}

	result := Intersection(slice1, slice2)

	if !reflect.DeepEqual(result, test) {
		t.Fatalf("The resulting map is not equal to the test one \n "+
			"(%v vs %v).", result, test)
	}
}

func BenchmarkIntersection(b *testing.B) {
	slice1 := make([]int, 1_000_000)
	slice2 := make([]int, 1_000_00)
	// seed the random
	rand.Seed(time.Now().UnixNano())
	// generate first slice values
	for i := range slice1 {
		slice1[i] = rand.Intn(100)
	}
	// generate second slice values
	for i := range slice2 {
		slice2[i] = rand.Intn(100)
	}

	b.ResetTimer()
	Intersection(slice1, slice2)
}

func TestFibonacci(t *testing.T) {
	numbers := []int{3, 5, 11}
	values := []int{2, 5, 89}
	fibonacci := Fibonacci()

	for i := range numbers {
		res := fibonacci(numbers[i])
		fib := values[i]

		if res != fib {
			t.Fatalf("Wrong fibonacci value (%v vs %v).", res, fib)
		}
	}
}
