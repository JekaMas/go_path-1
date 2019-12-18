package main

import (
	"strings"
)

// CountWords counting words in text.
func CountWords(text string) map[string]int {
	fields := strings.Fields(text) // split by regex?
	counts := make(map[string]int)

	for _, field := range fields {
		counts[field]++
	}

	return counts
}

// AtLeastOne returns the map of numbers that found
// in the slice at least once.
func AtLeastOne(bigSlice []int) map[int]struct{} {
	nums := make(map[int]struct{})

	for _, i := range bigSlice {
		nums[i] = struct{}{}
	}

	return nums
}

// Intersection returns the map with the numbers that
// are in each of the arrays.
func Intersection(slice1 []int, slice2 []int) (result []int) {
	nums1 := AtLeastOne(slice1)
	nums2 := AtLeastOne(slice2)

	for k := range nums1 {
		if _, ok := nums2[k]; ok {
			result = append(result, k)
		}
	}

	return
}

// Fibonacci returns a function that calculates
// the nth fibonacci number with preserving
// the old values.
func Fibonacci() func(n int) int {
	saved := make(map[int]int)
	saved[0], saved[1] = 0, 1
	savedN := 1

	// замыкание функции
	return func(n int) int {
		if n < 0 {
			return -1
		}

		if n > savedN {
			for i := savedN + 1; i <= n; i++ {
				saved[i] = saved[i-1] + saved[i-2]
			}
			savedN = n
		}
		return saved[n]
	}
}
