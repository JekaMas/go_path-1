package main

import (
	"fmt"
	"sort"
)

// 1. Добавить к каждому элементу единицу
func AddOne(slice []int) []int {
	for i := 0; i < len(slice); i++ {
		slice[i]++ // jekamas: можно использовать инкримент
	}
	return slice
}

// 2. Добавить в конец число 5
func AppendFive(slice []int) []int {
	return append(slice, 5)
}

// 3. Добавить в начало число 5
func PrependFive(slice []int) []int {
	return append([]int{5}, slice...)
}

// 4. Взять и удалить последний элемент
func PopLast(slice []int) (int, []int) {
	last := slice[len(slice)-1]
	return last, slice[:len(slice)-1]
}

// 5. Взять и удалить первый элемент */
func Pop(slice []int) (int, []int) {
	first := slice[0]
	return first, slice[1:]
}

// 6. Взять i-й элемент и удалить */
func PopIndex(slice []int, i int) (int, []int) {
	element := slice[i]
	copy(slice[i:], slice[i+1:len(slice)-1])
	return element, slice
}

// 7. Объединить 2 слайса
func Concat(slice1, slice2 []int) []int {
	return append(slice1, slice2...)
}

// 8. Удалить все элементы, которые есть во втором
// O((m)*n) реализация
func RemoveAll(slice, remove []int) (result []int) {
	for _, v := range slice {
		has := false // flag
		for _, r := range remove {
			if v == r {
				has = true
			}
		}
		if !has {
			result = append(result, v)
		}
	}
	return
}

// 8. in-place realization
func RemoveAllAlternative(slice, remove []int) []int {
	sort.Ints(remove) // because we're going to do N searches, let's presort the slice. It'll decrease the complexity from O(n*m) to O(m*log(m) + n*log(m))

	n := len(slice)
	var v int
	for i := 0; i < n; i++ {
		v = slice[i]

		if index := sort.SearchInts(remove, v); index < len(remove) && remove[index] == v {
			slice[i] = slice[len(slice)-1] // remove in-place
			slice = slice[:len(slice)-1]   // remove in-place

			i--
			n--
		}
	}

	return slice
}

// 9. Сдвинуть все элементы на 1 влево
func OffsetLeftOne(slice []int) []int {
	first := slice[0]
	copy(slice, slice[1:])
	slice[len(slice)-1] = first
	return slice
}

// 10. Сдвинуть все элементы на некоторое i влево
// jeksmas: а сможешь in-place реализацию? без создания нового массива
// с реализацией всё хорошо. я бы такой код принял и в прод. Что можно улучшить? offset = offset%len(slice) - вдруг дадут большое значение. и можно копировать не по одному элементу, а через copy
func OffsetLeft(slice []int, offset int) []int {
	var result = make([]int, len(slice))
	for i := 0; i < len(slice); i++ {
		j := i - offset
		if j < 0 {
			j += len(slice)
		}
		result[j] = slice[i]
	}
	return result
}

func OffsetLeft(slice []int, offset int) []int {
	var result = make([]int, len(slice))

	splitIndex := offset
	rightSide := slice[splitIndex:]
	leftSide := slice[:splitIndex]

	copy(result, rightSide)
	copy(result[len(rightSide):], leftSide)

	return result
}

// 11. Сдвинуть все элементы на 1 вправо
func OffsetRightOne(slice []int) []int {
	last := slice[len(slice)-1]
	for i := len(slice) - 1; i > 0; i-- {
		slice[i] = slice[i-1]
	}
	slice[0] = last
	return slice
}

// jekamas: OffsetRight, OffsetRightOne лучше реализовать через то, что уже сделано
func OffsetRight(slice []int, offset int) []int {
	return OffsetLeft(slice, len(slice)-offset)
}

func OffsetRightOne(slice []int) []int {
	return OffsetRight(slice, 1)
}

// 12. Сдвинуть все элементы на некоторое i вправо
func OffsetRight(slice []int, offset int) []int {
	var result = make([]int, len(slice))
	for i := 0; i < len(slice); i++ {
		result[(i+offset)%len(slice)] = slice[i]
	}
	return result
}

// 13. Копия слайса
// jekamas: все хорошо. больше решений https://github.com/golang/go/wiki/SliceTricks#copy
func Copy(slice []int) (sliceCopy []int) {
	sliceCopy = make([]int, len(slice))
	copy(sliceCopy, slice)
	return
}

/* 14. Поменять все чётные с ближайшими нечётными индексами */
func EvenOddSwap(slice []int) []int {
	for i := 0; i < len(slice)-1; i += 2 {
		slice[i], slice[i+1] = slice[i+1], slice[i]
	}
	return slice
}

// 15. Упорядочить слайс
// jekamas: всё хорошо. можно избавиться от else
func Sort(slice []int, reversed bool) []int {
	if reversed {
		sort.Sort(sort.Reverse(sort.IntSlice(slice)))
	} else {
		sort.Ints(slice)
	}
	return slice
}

func Sort(slice []int, reversed bool) []int {
	var toSort sort.Interface = sort.IntSlice(slice)
	if reversed {
		toSort = sort.Reverse(toSort)
	}
	sort.Sort(toSort)
	return slice
}

func SortLexical(slice []string) []string {
	sort.Strings(slice)
	return slice
}

func main() {
	slice := []int{0, 1, 2, 3}

	fmt.Println("Original: ", slice)

	fmt.Println(AddOne(slice))      // 1
	fmt.Println(AppendFive(slice))  // 2
	fmt.Println(PrependFive(slice)) // 3
	fmt.Println(PopLast(slice))     // 4
	fmt.Println(Pop(slice))         // 5
	fmt.Println(PopIndex(slice, 2)) // 6

	slice2 := []int{5, 6, 7}

	fmt.Println(Concat(slice, slice2))         // 7
	fmt.Println(RemoveAll(slice, []int{1, 3})) // 8
	fmt.Println(OffsetLeftOne(slice))          // 9
	fmt.Println(OffsetLeft(slice, 3))          // 10
	fmt.Println(OffsetRightOne(slice))         // 11
	fmt.Println(OffsetRight(slice, 3))         // 12
	fmt.Println(Copy(slice))                   // 13
	fmt.Println(EvenOddSwap(slice))            // 14

	fmt.Println(Sort(slice, false)) // 15
	fmt.Println(Sort(slice, true))  // 15

	sliceStr := []string{"hello", "beautiful", "world", "!"}

	fmt.Println(SortLexical(sliceStr)) // 15
}
